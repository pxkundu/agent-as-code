// Package api provides binary download functionality
package api

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Downloader handles binary downloads from the API
type Downloader struct {
	client *Client
}

// NewDownloader creates a new binary downloader
func NewDownloader(baseURL string) *Downloader {
	return &Downloader{
		client: NewClient(baseURL),
	}
}

// DownloadOptions represents options for binary download
type DownloadOptions struct {
	Version      string
	Platform     string
	Architecture string
	OutputDir    string
	OutputFile   string
}

// DownloadResult represents the result of a binary download
type DownloadResult struct {
	Success      bool
	Platform     string
	Architecture string
	Version      string
	FilePath     string
	Size         int64
	Error        error
}

// DownloadBinary downloads a specific binary version
func (d *Downloader) DownloadBinary(opts DownloadOptions) *DownloadResult {
	result := &DownloadResult{
		Platform:     opts.Platform,
		Architecture: opts.Architecture,
		Version:      opts.Version,
	}

	// Download binary data
	data, err := d.client.DownloadBinary(opts.Version, opts.Platform, opts.Architecture)
	if err != nil {
		result.Error = fmt.Errorf("download failed: %w", err)
		return result
	}

	// Determine output file path
	outputFile := opts.OutputFile
	if outputFile == "" {
		filename := fmt.Sprintf("agent_as_code_%s_%s_%s.zip", opts.Version, opts.Platform, opts.Architecture)
		outputFile = filepath.Join(opts.OutputDir, filename)
	}

	// Save to file
	if err := SaveBinaryToFile(data, outputFile); err != nil {
		result.Error = fmt.Errorf("failed to save file: %w", err)
		return result
	}

	// Get file info
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		result.Error = fmt.Errorf("failed to get file info: %w", err)
		return result
	}

	result.Success = true
	result.FilePath = outputFile
	result.Size = fileInfo.Size()

	return result
}

// DownloadLatest downloads the latest binary for current platform
func (d *Downloader) DownloadLatest(outputDir string) *DownloadResult {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Get latest binary info
	latest, err := d.client.GetLatestBinary()
	if err != nil {
		return &DownloadResult{
			Platform:     platform,
			Architecture: arch,
			Error:        fmt.Errorf("failed to get latest binary info: %w", err),
		}
	}

	opts := DownloadOptions{
		Version:      latest.Version,
		Platform:     platform,
		Architecture: arch,
		OutputDir:    outputDir,
	}

	return d.DownloadBinary(opts)
}

// DownloadAllPlatforms downloads binaries for all supported platforms
func (d *Downloader) DownloadAllPlatforms(version, outputDir string) []*DownloadResult {
	platforms := []struct {
		OS   string
		Arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
		{"windows", "arm64"},
	}

	var results []*DownloadResult

	for _, platform := range platforms {
		opts := DownloadOptions{
			Version:      version,
			Platform:     platform.OS,
			Architecture: platform.Arch,
			OutputDir:    outputDir,
		}

		result := d.DownloadBinary(opts)
		results = append(results, result)
	}

	return results
}

// ListAvailableVersions lists all available versions
func (d *Downloader) ListAvailableVersions() ([]string, error) {
	resp, err := d.client.ListVersions()
	if err != nil {
		return nil, err
	}

	return resp.Versions, nil
}

// ListAvailableBinaries lists all available binaries for a version
func (d *Downloader) ListAvailableBinaries(version string) ([]BinaryInfo, error) {
	major, minor, err := parseVersion(version)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.ListFiles(major, minor)
	if err != nil {
		return nil, err
	}

	return resp.Files, nil
}

// GetBinaryInfo gets information about a specific binary
func (d *Downloader) GetBinaryInfo(version, platform, arch string) (*BinaryInfo, error) {
	binaries, err := d.ListAvailableBinaries(version)
	if err != nil {
		return nil, err
	}

	for _, binary := range binaries {
		if binary.Platform == platform && binary.Architecture == arch {
			return &binary, nil
		}
	}

	return nil, fmt.Errorf("binary not found for %s/%s version %s", platform, arch, version)
}

// InstallBinary downloads and installs a binary to the system
func (d *Downloader) InstallBinary(version, installDir string) *DownloadResult {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "agent-install-")
	if err != nil {
		return &DownloadResult{
			Platform:     platform,
			Architecture: arch,
			Version:      version,
			Error:        fmt.Errorf("failed to create temp directory: %w", err),
		}
	}
	defer os.RemoveAll(tempDir)

	// Download binary
	opts := DownloadOptions{
		Version:      version,
		Platform:     platform,
		Architecture: arch,
		OutputDir:    tempDir,
	}

	result := d.DownloadBinary(opts)
	if !result.Success {
		return result
	}

	// Extract zip and install binary
	if result.Success {
		// Line 217 in internal/api/downloader.go
		if err := d.extractAndInstallBinary(result.FilePath, installDir, version, platform, arch); err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to install binary: %w", err)
		}
	}

	return result
}

// extractAndInstallBinary extracts the downloaded zip file and installs the binary
func (d *Downloader) extractAndInstallBinary(zipPath, installDir, version, platform, arch string) error {
	// Import archive/zip at the top of the file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Determine binary name
	binaryName := "agent"
	if platform == "windows" {
		binaryName += ".exe"
	}

	// Extract binary from zip
	var binaryFound bool
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, binaryName) {
			// Extract this file
			rc, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			// Create installation directory
			if err := os.MkdirAll(installDir, 0755); err != nil {
				return fmt.Errorf("failed to create install directory: %w", err)
			}

			// Create destination file
			destPath := filepath.Join(installDir, binaryName)
			destFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %w", err)
			}
			defer destFile.Close()

			// Copy binary content
			if _, err := io.Copy(destFile, rc); err != nil {
				return fmt.Errorf("failed to copy binary: %w", err)
			}

			// Set executable permissions on Unix systems
			if platform != "windows" {
				if err := os.Chmod(destPath, 0755); err != nil {
					return fmt.Errorf("failed to set executable permissions: %w", err)
				}
			}

			binaryFound = true
			break
		}
	}

	if !binaryFound {
		return fmt.Errorf("binary not found in zip file")
	}

	// Clean up downloaded zip file
	os.Remove(zipPath)

	fmt.Printf("✅ Binary installed successfully to %s\n", filepath.Join(installDir, binaryName))
	return nil
}
