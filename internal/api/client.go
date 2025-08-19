// Package api provides client functionality for the Agent-as-Code Binary API
package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Client represents the Binary API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewClient creates a new Binary API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAuthToken sets the authentication token for API requests
func (c *Client) SetAuthToken(token string) {
	c.AuthToken = token
}

// BinaryInfo represents metadata about a binary release
type BinaryInfo struct {
	Filename     string `json:"filename"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
	DownloadURL  string `json:"download_url"`
}

// VersionsResponse represents the response from the versions endpoint
type VersionsResponse struct {
	Success  bool     `json:"success"`
	Versions []string `json:"versions"`
	Count    int      `json:"count"`
}

// FilesResponse represents the response from the files endpoint
type FilesResponse struct {
	Success bool         `json:"success"`
	Major   int          `json:"major"`
	Minor   int          `json:"minor"`
	Files   []BinaryInfo `json:"files"`
	Count   int          `json:"count"`
}

// UploadRequest represents a binary upload request
type UploadRequest struct {
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	FileData     string `json:"file_data"` // Base64 encoded
	Filename     string `json:"filename"`  // Optional
	Checksum     string `json:"checksum"`  // Optional
}

// UploadResponse represents the response from binary upload
type UploadResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Release Release `json:"release"`
}

// Release represents a binary release
type Release struct {
	Version      string `json:"version"`
	Major        int    `json:"major"`
	Minor        int    `json:"minor"`
	Patch        int    `json:"patch"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	Filename     string `json:"filename"`
	S3Key        string `json:"s3_key"`
	FileSize     int64  `json:"file_size"`
	ContentType  string `json:"content_type"`
	UploadedAt   string `json:"uploaded_at"`
	Checksum     string `json:"checksum"`
	DownloadURL  string `json:"download_url"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ListVersions lists all available binary versions
func (c *Client) ListVersions() (*VersionsResponse, error) {
	url := fmt.Sprintf("%s/binary/releases/agent-as-code/versions", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var versionsResp VersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &versionsResp, nil
}

// ListFiles lists all files for a specific major.minor version
func (c *Client) ListFiles(major, minor int) (*FilesResponse, error) {
	url := fmt.Sprintf("%s/binary/releases/agent-as-code/%d/%d/", c.BaseURL, major, minor)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var filesResp FilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&filesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &filesResp, nil
}

// DownloadBinary downloads a specific binary release
func (c *Client) DownloadBinary(version, platform, arch string) ([]byte, error) {
	major, minor, err := parseVersion(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}

	filename := fmt.Sprintf("agent_as_code_%s_%s_%s.zip", version, platform, arch)
	url := fmt.Sprintf("%s/binary/releases/agent-as-code/%d/%d/%s", c.BaseURL, major, minor, filename)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// UploadBinary uploads a binary release
func (c *Client) UploadBinary(filePath, version, platform, arch string) (*UploadResponse, error) {
	if c.AuthToken == "" {
		return nil, fmt.Errorf("authentication token required for binary uploads")
	}

	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate checksum
	hash := sha256.Sum256(fileData)
	checksum := hex.EncodeToString(hash[:])

	// Encode file data to base64
	base64Data := base64.StdEncoding.EncodeToString(fileData)

	// Create filename
	filename := fmt.Sprintf("agent_as_code_%s_%s_%s.zip", version, platform, arch)

	// Create upload request
	uploadReq := UploadRequest{
		Version:      version,
		Platform:     platform,
		Architecture: arch,
		FileData:     base64Data,
		Filename:     filename,
		Checksum:     checksum,
	}

	// Parse version for URL
	major, minor, err := parseVersion(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}

	// Create request
	reqBody, err := json.Marshal(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/binary/releases/agent-as-code/%d/%d/upload", c.BaseURL, major, minor)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &uploadResp, nil
}

// GetLatestBinary gets the latest binary for the current platform
func (c *Client) GetLatestBinary() (*BinaryInfo, error) {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Get all versions
	versions, err := c.ListVersions()
	if err != nil {
		return nil, err
	}

	if len(versions.Versions) == 0 {
		return nil, fmt.Errorf("no versions available")
	}

	// Get the latest version (assuming they're sorted)
	latestVersion := versions.Versions[len(versions.Versions)-1]
	major, minor, err := parseVersion(latestVersion)
	if err != nil {
		return nil, err
	}

	// Get files for latest version
	files, err := c.ListFiles(major, minor)
	if err != nil {
		return nil, err
	}

	// Find binary for current platform
	for _, file := range files.Files {
		if file.Platform == platform && file.Architecture == arch {
			return &file, nil
		}
	}

	return nil, fmt.Errorf("no binary found for platform %s/%s", platform, arch)
}

// parseVersion parses a semantic version string and returns major, minor
func parseVersion(version string) (int, int, error) {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	return major, minor, nil
}

// handleErrorResponse handles API error responses
func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("API error: %s - %s", errorResp.Error, errorResp.Message)
}

// SaveBinaryToFile saves binary data to a file
func SaveBinaryToFile(data []byte, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
