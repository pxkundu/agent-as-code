#!/bin/bash
# Agent-as-Code Installation Script
# Downloads and installs the latest agent CLI binary

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REGISTRY_URL="${AGENT_REGISTRY_URL:-https://api.myagentregistry.com}"
INSTALL_DIR="${AGENT_INSTALL_DIR:-/usr/local/bin}"
VERSION="${AGENT_VERSION:-latest}"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect platform and architecture
detect_platform() {
    local platform=""
    local arch=""
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     platform="linux";;
        Darwin*)    platform="darwin";;
        CYGWIN*|MINGW*|MSYS*) platform="windows";;
        *)          
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64";;
        arm64|aarch64)  arch="arm64";;
        *)              
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    
    echo "${platform}_${arch}"
}

# Get latest version from API
get_latest_version() {
    log_info "Fetching latest version from registry..."
    
    local versions_url="${REGISTRY_URL}/binary/releases/agent-as-code/versions"
    local versions
    
    if command -v curl >/dev/null 2>&1; then
        versions=$(curl -s "$versions_url" | grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' | tail -1 | tr -d '"')
    elif command -v wget >/dev/null 2>&1; then
        versions=$(wget -qO- "$versions_url" | grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' | tail -1 | tr -d '"')
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    if [ -z "$versions" ]; then
        log_error "Could not determine latest version"
        exit 1
    fi
    
    echo "$versions"
}

# Parse version into major.minor for URL
parse_version() {
    local version="$1"
    local major=$(echo "$version" | cut -d. -f1)
    local minor=$(echo "$version" | cut -d. -f2)
    echo "${major}/${minor}"
}

# Download binary
download_binary() {
    local platform_arch="$1"
    local version="$2"
    local platform=$(echo "$platform_arch" | cut -d_ -f1)
    local arch=$(echo "$platform_arch" | cut -d_ -f2)
    
    # Determine binary name
    local binary_name="agent"
    if [ "$platform" = "windows" ]; then
        binary_name="agent.exe"
    fi
    
    # Create filename
    local filename="agent_as_code_${version}_${platform}_${arch}.zip"
    local version_path=$(parse_version "$version")
    local download_url="${REGISTRY_URL}/binary/releases/agent-as-code/${version_path}/${filename}"
    
    log_info "Downloading agent CLI v${version} for ${platform}/${arch}..."
    log_info "URL: ${download_url}"
    
    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    local zip_file="${temp_dir}/${filename}"
    
    # Download
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "$zip_file" "$download_url"; then
            log_error "Download failed"
            rm -rf "$temp_dir"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$zip_file" "$download_url"; then
            log_error "Download failed"
            rm -rf "$temp_dir"
            exit 1
        fi
    else
        log_error "Neither curl nor wget found"
        exit 1
    fi
    
    # Extract
    if ! command -v unzip >/dev/null 2>&1; then
        log_error "unzip not found. Please install unzip."
        rm -rf "$temp_dir"
        exit 1
    fi
    
    log_info "Extracting binary..."
    if ! unzip -q "$zip_file" -d "$temp_dir"; then
        log_error "Failed to extract binary"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Find extracted binary
    local extracted_binary
    extracted_binary=$(find "$temp_dir" -name "$binary_name" -type f | head -1)
    
    if [ -z "$extracted_binary" ]; then
        log_error "Binary not found in downloaded archive"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    echo "$extracted_binary"
}

# Install binary
install_binary() {
    local binary_path="$1"
    local install_path="${INSTALL_DIR}/agent"
    
    log_info "Installing agent CLI to ${install_path}..."
    
    # Check if install directory exists and is writable
    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Creating install directory: ${INSTALL_DIR}"
        if ! sudo mkdir -p "$INSTALL_DIR"; then
            log_error "Failed to create install directory"
            exit 1
        fi
    fi
    
    # Install binary
    if ! sudo cp "$binary_path" "$install_path"; then
        log_error "Failed to install binary"
        exit 1
    fi
    
    # Make executable
    if ! sudo chmod +x "$install_path"; then
        log_error "Failed to make binary executable"
        exit 1
    fi
    
    log_success "Agent CLI installed successfully!"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command -v agent >/dev/null 2>&1; then
        local version_output
        version_output=$(agent --version 2>&1)
        log_success "Installation verified: ${version_output}"
    else
        log_warning "agent command not found in PATH"
        log_info "You may need to add ${INSTALL_DIR} to your PATH"
        log_info "Add this to your shell profile:"
        log_info "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Show usage
usage() {
    cat << EOF
Agent-as-Code Installation Script

Usage: $0 [OPTIONS]

Options:
    --version VERSION    Install specific version (default: latest)
    --install-dir DIR    Installation directory (default: /usr/local/bin)
    --registry URL       Registry URL (default: https://api.myagentregistry.com)
    --help              Show this help

Environment Variables:
    AGENT_VERSION        Version to install
    AGENT_INSTALL_DIR    Installation directory
    AGENT_REGISTRY_URL   Registry URL

Examples:
    $0                          # Install latest version
    $0 --version 1.2.3          # Install specific version
    $0 --install-dir ~/.local/bin # Install to user directory

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --registry)
            REGISTRY_URL="$2"
            shift 2
            ;;
        --help)
            usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Main installation process
main() {
    echo
    log_info "🚀 Agent-as-Code Installation"
    log_info "Registry: ${REGISTRY_URL}"
    log_info "Install Directory: ${INSTALL_DIR}"
    echo
    
    # Detect platform
    local platform_arch
    platform_arch=$(detect_platform)
    log_info "Detected platform: ${platform_arch}"
    
    # Get version
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(get_latest_version)
        log_info "Latest version: ${VERSION}"
    else
        log_info "Installing version: ${VERSION}"
    fi
    
    # Download binary
    local binary_path
    binary_path=$(download_binary "$platform_arch" "$VERSION")
    
    # Install binary
    install_binary "$binary_path"
    
    # Cleanup
    rm -rf "$(dirname "$binary_path")"
    
    # Verify installation
    verify_installation
    
    echo
    log_success "🎉 Agent-as-Code installed successfully!"
    log_info "Get started with: agent --help"
    echo
}

# Run main function
main
