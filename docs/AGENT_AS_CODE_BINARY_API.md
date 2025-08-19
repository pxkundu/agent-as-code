# Agent as Code Binary API

This document describes the API for uploading and downloading Agent as Code tool binaries, following the Terraform release pattern.

## Overview

The Agent as Code Binary API provides a Terraform-like interface for distributing binary releases of the agent-as-code CLI tool. It supports versioned releases with platform and architecture specific binaries.

## URL Pattern

Following the Terraform release pattern:

**Terraform:**
```
https://releases.hashicorp.com/terraform/1.12.2/terraform_1.12.2_windows_386.zip
```

**Agent as Code:** We have our own binary creating process as a common build for all OS.
```
https://api.myagentregistry.com/binary/releases/agent-as-code/MAJOR/MINOR/agent_as_code_VERSION_PLATFORM_ARCH.zip
```

## Supported Platforms and Architectures

### Platforms
- `windows` - Microsoft Windows
- `linux` - Linux distributions
- `darwin` - macOS

### Architectures
- `amd64` - 64-bit x86 processors
- `arm64` - 64-bit ARM processors
- `386` - 32-bit x86 processors (legacy)

## API Endpoints

### 1. List All Versions

Get a list of all available versions.

```http
GET /binary/releases/agent-as-code/versions
```

**Response:**
```json
{
  "success": true,
  "versions": ["1.0.0", "1.1.0", "1.12.2"],
  "count": 3
}
```

### 2. List Version Files

Get all files available for a specific major.minor version.

```http
GET /binary/releases/agent-as-code/{major}/{minor}/
```

**Example:**
```http
GET /binary/releases/agent-as-code/1/12/
```

**Response:**
```json
{
  "success": true,
  "major": 1,
  "minor": 12,
  "files": [
    {
      "filename": "agent_as_code_1.12.2_windows_amd64.zip",
      "version": "1.12.2",
      "platform": "windows",
      "architecture": "amd64",
      "size": 15728640,
      "last_modified": "2024-01-15T10:30:00Z",
      "download_url": "/binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_windows_amd64.zip"
    },
    {
      "filename": "agent_as_code_1.12.2_linux_amd64.zip",
      "version": "1.12.2",
      "platform": "linux",
      "architecture": "amd64",
      "size": 14857216,
      "last_modified": "2024-01-15T10:30:00Z",
      "download_url": "/binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_linux_amd64.zip"
    },
    {
      "filename": "agent_as_code_1.12.2_darwin_arm64.zip",
      "version": "1.12.2",
      "platform": "darwin",
      "architecture": "arm64",
      "size": 14234624,
      "last_modified": "2024-01-15T10:30:00Z",
      "download_url": "/binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_darwin_arm64.zip"
    }
  ],
  "count": 3
}
```

### 3. Download Binary

Download a specific binary release.

```http
GET /binary/releases/agent-as-code/{major}/{minor}/agent_as_code_{version}_{platform}_{arch}.zip
```

**Examples:**
```http
GET /binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_windows_amd64.zip
GET /binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_linux_amd64.zip
GET /binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_darwin_arm64.zip
```

**Response:**
```json
{
  "success": true,
  "filename": "agent_as_code_1.12.2_windows_amd64.zip",
  "content_type": "application/zip",
  "content_length": 15728640,
  "metadata": {
    "version": "1.12.2",
    "platform": "windows",
    "architecture": "amd64",
    "checksum": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "uploaded_at": "2024-01-15T10:30:00Z"
  },
  "file_data": "UEsDBBQAAAAIAA...", // Base64 encoded zip file
  "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```

### 4. Upload Binary (Authenticated)

Upload a new binary release. Requires authentication via Cognito JWT Bearer token.

```http
POST /binary/releases/agent-as-code/{major}/{minor}/upload
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "version": "1.12.2",
  "platform": "windows",
  "architecture": "amd64",
  "file_data": "UEsDBBQAAAAIAA...", // Base64 encoded zip file
  "filename": "agent_as_code_1.12.2_windows_amd64.zip", // Optional, auto-generated if not provided
  "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // Optional, auto-calculated if not provided
}
```

**Response:**
```json
{
  "success": true,
  "message": "Binary uploaded successfully",
  "release": {
    "version": "1.12.2",
    "major": 1,
    "minor": 12,
    "patch": 2,
    "platform": "windows",
    "architecture": "amd64",
    "filename": "agent_as_code_1.12.2_windows_amd64.zip",
    "s3_key": "releases/1/12/2/agent_as_code_1.12.2_windows_amd64.zip",
    "file_size": 15728640,
    "content_type": "application/zip",
    "uploaded_at": "2024-01-15T10:30:00Z",
    "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "download_url": "/binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_windows_amd64.zip"
  }
}
```

## Storage Structure

### S3 Bucket
```
Bucket: agent-as-code-{account}-{region}
```

### S3 Key Pattern
```
releases/{major}/{minor}/{patch}/agent_as_code_{version}_{platform}_{arch}.zip
```

**Examples:**
```
releases/1/12/2/agent_as_code_1.12.2_windows_amd64.zip
releases/1/12/2/agent_as_code_1.12.2_linux_amd64.zip
releases/1/12/2/agent_as_code_1.12.2_darwin_arm64.zip
```

## Authentication

- **Downloads**: Public access (no authentication required)
- **Uploads**: Requires authentication via Cognito JWT token: `Bearer <cognito-jwt-token>`

## Error Responses

### 400 Bad Request
```json
{
  "error": "Bad Request",
  "message": "Missing required fields: version, platform, architecture"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Authentication token required for binary uploads"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Binary not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "Failed to upload binary: <details>"
}
```

## Usage Examples

### CLI Download Example
```bash
# Download latest version for current platform
curl -o agent_as_code.zip \
  "https://api.myagentregistry.com/binary/releases/agent-as-code/1/12/agent_as_code_1.12.2_linux_amd64.zip"

# Extract and install
unzip agent_as_code.zip
chmod +x agent_as_code
sudo mv agent_as_code /usr/local/bin/
```

### Upload Example (with authentication)
```bash
# Upload new binary release
curl -X POST \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.12.2",
    "platform": "linux",
    "architecture": "amd64",
    "file_data": "'$(base64 -w 0 agent_as_code_1.12.2_linux_amd64.zip)'"
  }' \
  "https://api.myagentregistry.com/binary/releases/agent-as-code/1/12/upload"
```

## Integration with CLI Tools

This API is designed to be consumed by:

1. **Agent as Code CLI**: For self-updating functionality
2. **Package managers**: Homebrew, apt, yum, etc.
3. **CI/CD pipelines**: For automated builds and releases
4. **Installation scripts**: For easy setup on various platforms

## Security Considerations

1. **Public Downloads**: All binary downloads are publicly accessible to support easy installation
2. **Authenticated Uploads**: Only authenticated users can upload new releases
3. **Checksum Verification**: All binaries include SHA256 checksums for integrity verification
4. **S3 Security**: Proper IAM policies and bucket policies ensure secure storage
5. **Rate Limiting**: API Gateway provides built-in rate limiting and throttling
