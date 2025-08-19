// Package api provides binary upload functionality
package api

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Uploader handles binary uploads to the API
type Uploader struct {
	client  *Client
	version string
}

// NewUploader creates a new binary uploader
func NewUploader(baseURL, authToken, version string) *Uploader {
	client := NewClient(baseURL)
	client.SetAuthToken(authToken)

	return &Uploader{
		client:  client,
		version: version,
	}
}

// UploadOptions represents options for binary upload
type UploadOptions struct {
	Platform     string
	Architecture string
	FilePath     string
	Force        bool // Overwrite existing binary
}

// UploadResult represents the result of a binary upload
type UploadResult struct {
	Success      bool
	Platform     string
	Architecture string
	Version      string
	DownloadURL  string
	Error        error
}

// UploadBinary uploads a single binary
func (u *Uploader) UploadBinary(opts UploadOptions) *UploadResult {
	result := &UploadResult{
		Platform:     opts.Platform,
		Architecture: opts.Architecture,
		Version:      u.version,
	}

	// Validate file exists
	if _, err := os.Stat(opts.FilePath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("binary file not found: %s", opts.FilePath)
		return result
	}

	// Upload binary
	resp, err := u.client.UploadBinary(opts.FilePath, u.version, opts.Platform, opts.Architecture)
	if err != nil {
		result.Error = fmt.Errorf("upload failed: %w", err)
		return result
	}

	result.Success = resp.Success
	result.DownloadURL = resp.Release.DownloadURL

	return result
}

// UploadAllPlatforms uploads binaries for all supported platforms
func (u *Uploader) UploadAllPlatforms(binDir string) []*UploadResult {
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

	var results []*UploadResult

	for _, platform := range platforms {
		// Determine binary filename
		binaryName := "agent"
		if platform.OS == "windows" {
			binaryName += ".exe"
		}

		// Construct binary path
		binaryPath := filepath.Join(binDir, fmt.Sprintf("%s-%s-%s", binaryName, platform.OS, platform.Arch))
		if platform.OS == "windows" {
			binaryPath = strings.TrimSuffix(binaryPath, ".exe") + ".exe"
		}

		// Check if binary exists
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			results = append(results, &UploadResult{
				Platform:     platform.OS,
				Architecture: platform.Arch,
				Version:      u.version,
				Error:        fmt.Errorf("binary not found: %s", binaryPath),
			})
			continue
		}

		// Upload binary
		opts := UploadOptions{
			Platform:     platform.OS,
			Architecture: platform.Arch,
			FilePath:     binaryPath,
		}

		result := u.UploadBinary(opts)
		results = append(results, result)
	}

	return results
}

// UploadCurrentPlatform uploads binary for current platform only
func (u *Uploader) UploadCurrentPlatform(binaryPath string) *UploadResult {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	opts := UploadOptions{
		Platform:     platform,
		Architecture: arch,
		FilePath:     binaryPath,
	}

	return u.UploadBinary(opts)
}

// ValidateUpload validates a binary upload by downloading and comparing
func (u *Uploader) ValidateUpload(platform, arch string) error {
	// Download the binary we just uploaded
	data, err := u.client.DownloadBinary(u.version, platform, arch)
	if err != nil {
		return fmt.Errorf("failed to download binary for validation: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("downloaded binary is empty")
	}

	return nil
}

// GetUploadSummary returns a summary of upload results
func GetUploadSummary(results []*UploadResult) string {
	var summary strings.Builder

	successful := 0
	failed := 0

	summary.WriteString("📦 Binary Upload Summary:\n\n")

	for _, result := range results {
		if result.Success {
			successful++
			summary.WriteString(fmt.Sprintf("✅ %s/%s - %s\n",
				result.Platform, result.Architecture, result.DownloadURL))
		} else {
			failed++
			summary.WriteString(fmt.Sprintf("❌ %s/%s - %s\n",
				result.Platform, result.Architecture, result.Error.Error()))
		}
	}

	summary.WriteString(fmt.Sprintf("\n📊 Results: %d successful, %d failed\n", successful, failed))

	if successful > 0 {
		summary.WriteString("\n🎉 Binaries are now available for download!\n")
	}

	return summary.String()
}
