#!/bin/bash
# Agent-as-Code Build Script
# Builds binaries for all supported platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="agent-as-code"
BINARY_NAME="agent"
BUILD_DIR="bin"
PYTHON_DIR="python"

# Supported platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

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

# Get version from git or environment
get_version() {
    if [ -n "$VERSION" ]; then
        echo "$VERSION"
    elif [ -n "$GITHUB_REF_NAME" ]; then
        echo "$GITHUB_REF_NAME"
    elif git describe --tags --exact-match 2>/dev/null; then
        git describe --tags --exact-match
    else
        echo "dev-$(git rev-parse --short HEAD)"
    fi
}

# Get commit hash
get_commit() {
    if [ -n "$GITHUB_SHA" ]; then
        echo "$GITHUB_SHA"
    else
        git rev-parse HEAD
    fi
}

# Get build date
get_date() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
}

# Build binary for specific platform
build_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    
    local output_name="${BINARY_NAME}-${os}-${arch}"
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${output_name}"
    
    log_info "Building ${platform}..."
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0
    
    # Build with ldflags for version info
    go build \
        -ldflags="-X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.date=${BUILD_DATE}' -w -s" \
        -o "${output_path}" \
        ./cmd/agent
    
    if [ $? -eq 0 ]; then
        log_success "Built ${output_path}"
        
        # Get file size
        if command -v ls >/dev/null 2>&1; then
            size=$(ls -lh "${output_path}" | awk '{print $5}')
            log_info "Size: ${size}"
        fi
    else
        log_error "Failed to build ${platform}"
        return 1
    fi
}

# Create build directory
prepare_build() {
    log_info "Preparing build environment..."
    
    # Clean and create build directory
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    
    # Verify Go installation
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    log_info "Go version: $(go version)"
    
    # Get build information
    VERSION=$(get_version)
    COMMIT=$(get_commit)
    BUILD_DATE=$(get_date)
    
    log_info "Version: $VERSION"
    log_info "Commit: $COMMIT"
    log_info "Date: $BUILD_DATE"
    
    # Export for use in build functions
    export VERSION COMMIT BUILD_DATE
}

# Build all platforms
build_all() {
    log_info "Building binaries for all platforms..."
    
    local failed_builds=()
    
    for platform in "${PLATFORMS[@]}"; do
        if ! build_platform "$platform"; then
            failed_builds+=("$platform")
        fi
    done
    
    # Report results
    echo
    log_info "Build Summary:"
    
    total=${#PLATFORMS[@]}
    failed=${#failed_builds[@]}
    successful=$((total - failed))
    
    log_info "Total platforms: $total"
    log_success "Successful builds: $successful"
    
    if [ $failed -gt 0 ]; then
        log_error "Failed builds: $failed"
        for platform in "${failed_builds[@]}"; do
            log_error "  - $platform"
        done
        return 1
    fi
    
    log_success "All binaries built successfully!"
}

# Copy binaries to Python package
copy_to_python() {
    if [ ! -d "$PYTHON_DIR" ]; then
        log_warning "Python directory not found, skipping binary copy"
        return 0
    fi
    
    log_info "Copying binaries to Python package..."
    
    local python_bin_dir="${PYTHON_DIR}/agent_as_code/bin"
    mkdir -p "$python_bin_dir"
    
    # Copy all binaries
    cp -r "$BUILD_DIR"/* "$python_bin_dir/"
    
    log_success "Binaries copied to Python package"
}

# Create checksums
create_checksums() {
    log_info "Creating checksums..."
    
    cd "$BUILD_DIR"
    
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum * > checksums.txt
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 * > checksums.txt
    else
        log_warning "No checksum utility found, skipping checksums"
        cd ..
        return 0
    fi
    
    log_success "Checksums created"
    cd ..
}

# Main execution
main() {
    echo
    log_info "🚀 Starting Agent-as-Code build process..."
    echo
    
    prepare_build
    build_all
    copy_to_python
    create_checksums
    
    echo
    log_success "🎉 Build completed successfully!"
    log_info "Binaries available in: $BUILD_DIR"
    
    # List all built binaries
    echo
    log_info "Built binaries:"
    ls -la "$BUILD_DIR"
}

# Run main function
main "$@"
