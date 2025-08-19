#!/bin/bash
# Agent-as-Code Release Script
# Complete release workflow: build, test, upload, and publish

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION=""
REGISTRY_URL="${REGISTRY_URL:-https://api.myagentregistry.com}"
SKIP_TESTS=""
SKIP_UPLOAD=""
DRY_RUN=""

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

log_step() {
    echo -e "\n${BLUE}==== $1 ====${NC}"
}

# Show usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] VERSION

Complete release workflow for Agent-as-Code.

Arguments:
    VERSION               Version to release (e.g., 1.2.3)

Options:
    --skip-tests         Skip running tests
    --skip-upload        Skip uploading binaries
    --dry-run           Show what would be done without executing
    -h, --help          Show this help message

Examples:
    $0 1.2.3
    $0 1.2.3 --skip-tests
    $0 1.2.3 --dry-run

Environment Variables:
    AGENT_REGISTRY_TOKEN     Authentication token for binary uploads
    REGISTRY_URL             Registry base URL (default: https://api.myagentregistry.com)
EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-tests)
                SKIP_TESTS=true
                shift
                ;;
            --skip-upload)
                SKIP_UPLOAD=true
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
                if [ -z "$VERSION" ]; then
                    VERSION="$1"
                else
                    log_error "Unknown argument: $1"
                    usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
}

# Validate version format
validate_version() {
    if [ -z "$VERSION" ]; then
        log_error "Version is required"
        usage
        exit 1
    fi
    
    if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "Invalid version format. Expected: MAJOR.MINOR.PATCH (e.g., 1.2.3)"
        exit 1
    fi
    
    log_info "Release version: $VERSION"
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking Prerequisites"
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository"
        exit 1
    fi
    
    # Check for uncommitted changes
    if [ -n "$(git status --porcelain)" ]; then
        log_warning "There are uncommitted changes in the repository"
        if [ "$DRY_RUN" != true ]; then
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_error "Aborting due to uncommitted changes"
                exit 1
            fi
        fi
    fi
    
    # Check if tag already exists
    if git tag | grep -q "^v$VERSION$"; then
        log_error "Tag v$VERSION already exists"
        exit 1
    fi
    
    # Check required tools
    local required_tools=("go" "python3" "make")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            log_error "Required tool not found: $tool"
            exit 1
        fi
    done
    
    log_success "Prerequisites check passed"
}

# Clean build environment
clean_build() {
    log_step "Cleaning Build Environment"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: make clean"
        return 0
    fi
    
    if ! make clean; then
        log_error "Failed to clean build environment"
        exit 1
    fi
    
    log_success "Build environment cleaned"
}

# Run tests
run_tests() {
    if [ "$SKIP_TESTS" = true ]; then
        log_warning "Skipping tests (--skip-tests specified)"
        return 0
    fi
    
    log_step "Running Tests"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: make test"
        return 0
    fi
    
    if ! make test; then
        log_error "Tests failed"
        exit 1
    fi
    
    log_success "All tests passed"
}

# Build binaries
build_binaries() {
    log_step "Building Binaries"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: make build VERSION=$VERSION"
        return 0
    fi
    
    if ! make build VERSION="$VERSION"; then
        log_error "Build failed"
        exit 1
    fi
    
    log_success "Binaries built successfully"
}

# Upload binaries
upload_binaries() {
    if [ "$SKIP_UPLOAD" = true ]; then
        log_warning "Skipping binary upload (--skip-upload specified)"
        return 0
    fi
    
    log_step "Uploading Binaries"
    
    if [ -z "$AGENT_REGISTRY_TOKEN" ]; then
        log_error "AGENT_REGISTRY_TOKEN environment variable required for uploading"
        exit 1
    fi
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: make upload VERSION=$VERSION"
        return 0
    fi
    
    if ! make upload VERSION="$VERSION"; then
        log_error "Binary upload failed"
        exit 1
    fi
    
    log_success "Binaries uploaded successfully"
}

# Create git tag
create_tag() {
    log_step "Creating Git Tag"
    
    local tag_message="Release v$VERSION"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would execute: git tag -a v$VERSION -m \"$tag_message\""
        log_info "Would execute: git push origin v$VERSION"
        return 0
    fi
    
    if ! git tag -a "v$VERSION" -m "$tag_message"; then
        log_error "Failed to create git tag"
        exit 1
    fi
    
    if ! git push origin "v$VERSION"; then
        log_error "Failed to push git tag"
        exit 1
    fi
    
    log_success "Git tag v$VERSION created and pushed"
}

# Generate changelog
generate_changelog() {
    log_step "Generating Changelog"
    
    local changelog_file="CHANGELOG-v$VERSION.md"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would generate changelog: $changelog_file"
        return 0
    fi
    
    # Get commits since last tag
    local last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    local commit_range=""
    
    if [ -n "$last_tag" ]; then
        commit_range="$last_tag..HEAD"
        log_info "Generating changelog from $last_tag to HEAD"
    else
        commit_range="HEAD"
        log_info "Generating changelog for all commits"
    fi
    
    # Generate changelog
    cat > "$changelog_file" << EOF
# Release v$VERSION

Release Date: $(date -u +"%Y-%m-%d")

## Changes

$(git log $commit_range --pretty=format:"- %s" --no-merges)

## Binary Downloads

All binaries are available at: $REGISTRY_URL/binary/releases/agent-as-code/

### Supported Platforms

- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

### Installation

\`\`\`bash
# Download latest release
curl -L $REGISTRY_URL/install.sh | sh

# Or install via pip
pip install agent-as-code==$VERSION
\`\`\`

EOF
    
    log_success "Changelog generated: $changelog_file"
}

# Verify release
verify_release() {
    log_step "Verifying Release"
    
    if [ "$DRY_RUN" = true ]; then
        log_info "Would verify binary downloads and functionality"
        return 0
    fi
    
    # Verify binaries are downloadable
    if [ "$SKIP_UPLOAD" != true ]; then
        log_info "Verifying binary downloads..."
        
        # Test download for current platform
        local platform=$(go env GOOS)
        local arch=$(go env GOARCH)
        
        if ! timeout 30 make download-version VERSION="$VERSION" >/dev/null 2>&1; then
            log_warning "Could not verify binary download (this may be expected)"
        else
            log_success "Binary download verified"
        fi
    fi
    
    log_success "Release verification completed"
}

# Print release summary
print_summary() {
    log_step "Release Summary"
    
    cat << EOF

🎉 Release v$VERSION completed successfully!

📦 Release Information:
   Version: $VERSION
   Tag: v$VERSION
   Registry: $REGISTRY_URL
   
📥 Installation:
   pip install agent-as-code==$VERSION
   curl -L $REGISTRY_URL/install.sh | sh

🔗 Links:
   • Binaries: $REGISTRY_URL/binary/releases/agent-as-code/
   • Repository: https://github.com/pxkundu/agent-as-code
   • Documentation: https://agent-as-code.myagentregistry.com/documentation

📋 Next Steps:
   1. Verify installation: agent --version
   2. Update documentation if needed
   3. Announce the release
   4. Monitor for any issues

EOF

    if [ "$DRY_RUN" = true ]; then
        echo -e "${YELLOW}NOTE: This was a dry run - no actual changes were made${NC}"
    fi
}

# Main execution
main() {
    echo
    log_info "🚀 Starting Agent-as-Code release process..."
    echo
    
    parse_args "$@"
    validate_version
    check_prerequisites
    clean_build
    run_tests
    build_binaries
    upload_binaries
    create_tag
    generate_changelog
    verify_release
    print_summary
    
    echo
    log_success "🎉 Release process completed successfully!"
}

# Run main function
main "$@"
