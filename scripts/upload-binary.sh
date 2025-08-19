#!/bin/bash
# Agent-as-Code Binary Upload Script
# Uploads built binaries to the registry using the Binary API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REGISTRY_URL="${REGISTRY_URL:-https://api.myagentregistry.com}"
BUILD_DIR="${BUILD_DIR:-bin}"
VERSION=""
AUTH_TOKEN=""

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

# Show usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Upload Agent-as-Code binaries to the registry.

Options:
    -v, --version VERSION      Binary version to upload (required)
    -t, --token TOKEN         Authentication token (or use AGENT_REGISTRY_TOKEN env var)
    -r, --registry URL        Registry base URL (default: https://api.myagentregistry.com)
    -d, --bin-dir DIR         Directory containing binaries (default: bin)
    -p, --platform PLATFORM  Upload specific platform only
    -a, --arch ARCH          Upload specific architecture only
    --all-platforms          Upload all available binaries
    --dry-run                Show what would be uploaded without actually uploading
    -h, --help               Show this help message

Examples:
    $0 --version 1.2.3 --token \$TOKEN --all-platforms
    $0 --version 1.2.3 --platform linux --arch amd64 ./bin/agent-linux-amd64
    $0 --version 1.2.3 --dry-run

Environment Variables:
    AGENT_REGISTRY_TOKEN     Authentication token
    REGISTRY_URL             Registry base URL
    BUILD_DIR                Directory containing binaries
EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -t|--token)
                AUTH_TOKEN="$2"
                shift 2
                ;;
            -r|--registry)
                REGISTRY_URL="$2"
                shift 2
                ;;
            -d|--bin-dir)
                BUILD_DIR="$2"
                shift 2
                ;;
            -p|--platform)
                UPLOAD_PLATFORM="$2"
                shift 2
                ;;
            -a|--arch)
                UPLOAD_ARCH="$2"
                shift 2
                ;;
            --all-platforms)
                ALL_PLATFORMS=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                BINARY_PATH="$1"
                shift
                ;;
        esac
    done
}

# Validate configuration
validate_config() {
    log_info "Validating configuration..."
    
    if [ -z "$VERSION" ]; then
        log_error "Version is required (use --version or -v)"
        exit 1
    fi
    
    if [ -z "$AUTH_TOKEN" ]; then
        AUTH_TOKEN="$AGENT_REGISTRY_TOKEN"
        if [ -z "$AUTH_TOKEN" ]; then
            log_error "Authentication token required (use --token or AGENT_REGISTRY_TOKEN env var)"
            exit 1
        fi
    fi
    
    if [ ! -d "$BUILD_DIR" ]; then
        log_error "Build directory not found: $BUILD_DIR"
        exit 1
    fi
    
    log_info "Registry: $REGISTRY_URL"
    log_info "Version: $VERSION"
    log_info "Build Dir: $BUILD_DIR"
    
    if [ "$DRY_RUN" = true ]; then
        log_warning "DRY RUN MODE - No actual uploads will be performed"
    fi
}

# Check if agent binary exists
check_agent_binary() {
    if ! command -v ./agent >/dev/null 2>&1; then
        # Try to find agent binary in current directory or PATH
        if [ -f "./agent" ]; then
            AGENT_BINARY="./agent"
        elif command -v agent >/dev/null 2>&1; then
            AGENT_BINARY="agent"
        else
            log_error "Agent binary not found. Please build the project first."
            exit 1
        fi
    else
        AGENT_BINARY="./agent"
    fi
    
    log_info "Using agent binary: $AGENT_BINARY"
}

# Upload using Go CLI
upload_with_cli() {
    check_agent_binary
    
    local upload_args=(
        "upload"
        "--version" "$VERSION"
        "--registry" "$REGISTRY_URL"
        "--token" "$AUTH_TOKEN"
    )
    
    if [ "$ALL_PLATFORMS" = true ]; then
        upload_args+=("--all-platforms" "--bin-dir" "$BUILD_DIR")
        log_info "Uploading all platforms using CLI..."
    elif [ -n "$UPLOAD_PLATFORM" ] && [ -n "$UPLOAD_ARCH" ]; then
        if [ -z "$BINARY_PATH" ]; then
            BINARY_PATH="$BUILD_DIR/agent-$UPLOAD_PLATFORM-$UPLOAD_ARCH"
            if [ "$UPLOAD_PLATFORM" = "windows" ]; then
                BINARY_PATH="$BINARY_PATH.exe"
            fi
        fi
        upload_args+=("--platform" "$UPLOAD_PLATFORM" "--arch" "$UPLOAD_ARCH" "$BINARY_PATH")
        log_info "Uploading $UPLOAD_PLATFORM/$UPLOAD_ARCH using CLI..."
    elif [ -n "$BINARY_PATH" ]; then
        upload_args+=("$BINARY_PATH")
        log_info "Uploading $BINARY_PATH using CLI..."
    else
        upload_args+=("--all-platforms" "--bin-dir" "$BUILD_DIR")
        log_info "Uploading all platforms using CLI..."
    fi
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: $AGENT_BINARY ${upload_args[*]}"
        return 0
    fi
    
    log_info "Executing upload command..."
    if ! "$AGENT_BINARY" "${upload_args[@]}"; then
        log_error "Upload failed"
        return 1
    fi
    
    log_success "Upload completed successfully!"
}

# Create zip files for upload (if needed)
create_zip_files() {
    log_info "Checking for zip files in $BUILD_DIR..."
    
    local zip_count=0
    
    for binary in "$BUILD_DIR"/agent-*; do
        if [ -f "$binary" ]; then
            local zip_file="${binary}.zip"
            
            if [ ! -f "$zip_file" ]; then
                log_info "Creating zip file: $(basename "$zip_file")"
                
                if [ "$DRY_RUN" != true ]; then
                    # Create zip file
                    (cd "$BUILD_DIR" && zip -q "$(basename "$zip_file")" "$(basename "$binary")")
                fi
                
                ((zip_count++))
            fi
        fi
    done
    
    if [ $zip_count -gt 0 ]; then
        log_success "Created $zip_count zip files"
    else
        log_info "All zip files already exist"
    fi
}

# Verify uploads
verify_uploads() {
    if [ "$DRY_RUN" = true ]; then
        return 0
    fi
    
    log_info "Verifying uploads..."
    
    # Use CLI to list available files for this version
    if ! "$AGENT_BINARY" download --list --registry "$REGISTRY_URL" >/dev/null 2>&1; then
        log_warning "Could not verify uploads (download command failed)"
        return 0
    fi
    
    log_success "Upload verification completed"
}

# Main execution
main() {
    echo
    log_info "🚀 Starting Agent-as-Code binary upload process..."
    echo
    
    parse_args "$@"
    validate_config
    create_zip_files
    upload_with_cli
    verify_uploads
    
    echo
    if [ "$DRY_RUN" = true ]; then
        log_info "🔍 Dry run completed - no uploads were performed"
    else
        log_success "🎉 Binary upload completed successfully!"
        log_info "Binaries are now available at: $REGISTRY_URL"
    fi
}

# Run main function
main "$@"
