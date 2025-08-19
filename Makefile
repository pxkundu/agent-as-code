# Agent as Code - Makefile
# Build automation for the revolutionized hybrid Go + Python approach

.PHONY: help build build-go build-python test clean install dev-setup release upload download

# Default target
help: ## Show this help message
	@echo "Agent as Code - Build System"
	@echo "============================="
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build configuration
VERSION ?= 1.0.0
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build configuration
GOOS_TARGETS = linux darwin windows
GOARCH_TARGETS = amd64 arm64
GO_LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build directories
BUILD_DIR = bin
PYTHON_DIR = python
PYTHON_BIN_DIR = $(PYTHON_DIR)/agent_as_code/bin

# Registry configuration
REGISTRY_URL ?= https://api.myagentregistry.com
AGENT_REGISTRY_TOKEN ?= $(shell echo $$AGENT_REGISTRY_TOKEN)

# Build all targets
build: build-go build-python ## Build both Go binaries and Python package

# Build Go binaries for all platforms
build-go: ## Build Go binaries for all platforms
	@echo "🔨 Building Go binaries..."
	@mkdir -p $(BUILD_DIR) $(PYTHON_BIN_DIR)
	@for os in $(GOOS_TARGETS); do \
		for arch in $(GOARCH_TARGETS); do \
			if [ "$$os" = "windows" ]; then \
				binary_name="agent-$$os-$$arch.exe"; \
			else \
				binary_name="agent-$$os-$$arch"; \
			fi; \
			echo "  Building $$binary_name..."; \
			GOOS=$$os GOARCH=$$arch go build $(GO_LDFLAGS) -o $(BUILD_DIR)/$$binary_name ./cmd/agent; \
			cp $(BUILD_DIR)/$$binary_name $(PYTHON_BIN_DIR)/; \
		done; \
	done
	@echo "✅ Go binaries built successfully!"

# Build Python package
build-python: build-go ## Build Python package (requires Go binaries)
	@echo "🐍 Building Python package..."
	        @cd $(PYTHON_DIR) && python3 -m build
	@echo "✅ Python package built successfully!"

# Test Go code
test-go: ## Run Go tests
	@echo "🧪 Running Go tests..."
	@go test -v ./...

# Test Python code
test-python: ## Run Python tests
	@echo "🧪 Running Python tests..."
	        @cd $(PYTHON_DIR) && python3 -m pytest tests/ -v

# Run all tests
test: test-go test-python ## Run all tests

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/agent-*
	@rm -rf $(PYTHON_BIN_DIR)/agent-*
	@rm -rf $(PYTHON_DIR)/build/
	@rm -rf $(PYTHON_DIR)/dist/
	@rm -rf $(PYTHON_DIR)/*.egg-info/
	@rm -rf downloads/
	@rm -f Dockerfile.agent
	@go clean
	@echo "✅ Clean completed!"

# Install Python package locally
install: build-python ## Install Python package locally
	@echo "📦 Installing Python package locally..."
	@cd $(PYTHON_DIR) && pip install -e .
	@echo "✅ Package installed! Try: agent --help"

# Development setup
dev-setup: ## Set up development environment
	@echo "🛠️  Setting up development environment..."
	@go mod download
	@cd $(PYTHON_DIR) && pip install -e ".[dev]"
	@echo "✅ Development environment ready!"

# Format code
fmt: ## Format Go and Python code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@cd $(PYTHON_DIR) && black agent_as_code/ tests/
	@echo "✅ Code formatted!"

# Lint code
lint: ## Lint Go and Python code
	@echo "🔍 Linting code..."
	@golangci-lint run ./...
	@cd $(PYTHON_DIR) && flake8 agent_as_code/ tests/
	@cd $(PYTHON_DIR) && mypy agent_as_code/
	@echo "✅ Linting completed!"

# Security scan
security: ## Run security scans
	@echo "🔒 Running security scans..."
	@gosec ./...
	@cd $(PYTHON_DIR) && safety check
	@echo "✅ Security scan completed!"

# Generate templates
generate-templates: ## Generate template files
	@echo "📝 Generating template files..."
	@mkdir -p templates/basic/python
	@echo "# Basic Python agent template" > templates/basic/python/main.py
	@echo "✅ Templates generated!"

# Release preparation
release-prep: clean build test lint security ## Prepare for release
	@echo "🚀 Release preparation completed!"
	@echo "   Version: $(VERSION)"
	@echo "   Commit: $(COMMIT)"
	@echo "   Date: $(DATE)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review build artifacts in $(PYTHON_DIR)/dist/"
	@echo "  2. Test installation: make install"
	@echo "  3. Create release: git tag v$(VERSION) && git push --tags"

# Local testing
test-install: build-python ## Test local installation
	@echo "🧪 Testing local installation..."
	@cd $(PYTHON_DIR) && pip install dist/*.whl --force-reinstall
	@agent --version
	@agent --help
	@echo "✅ Local installation test passed!"

# Docker build (for testing)
docker-build: ## Build Docker image for testing
	@echo "🐳 Building Docker test image..."
	@docker build -t agent-as-code:test .
	@echo "✅ Docker image built!"

# Benchmark
benchmark: ## Run performance benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks completed!"

# Documentation
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@go doc ./...
	@cd $(PYTHON_DIR) && python -m sphinx.cmd.build -b html docs/ docs/_build/
	@echo "✅ Documentation generated!"

# Check dependencies
deps-check: ## Check for dependency updates
	@echo "🔍 Checking dependencies..."
	@go list -u -m all
	@cd $(PYTHON_DIR) && pip list --outdated
	@echo "✅ Dependency check completed!"

# Update dependencies
deps-update: ## Update dependencies
	@echo "⬆️  Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@cd $(PYTHON_DIR) && pip install --upgrade pip setuptools wheel
	@echo "✅ Dependencies updated!"

# Cross-platform test
test-cross-platform: build-go ## Test binaries on current platform
	@echo "🌍 Testing cross-platform binaries..."
	@for binary in $(PYTHON_BIN_DIR)/agent-*; do \
		if [[ "$$binary" == *"$(shell go env GOOS)-$(shell go env GOARCH)"* ]] || \
		   [[ "$$binary" == *"$(shell go env GOOS)-$(shell go env GOARCH).exe"* ]]; then \
			echo "  Testing $$binary..."; \
			$$binary --version || echo "    ❌ Failed"; \
		fi; \
	done
	@echo "✅ Cross-platform test completed!"

# Show build info
info: ## Show build information
	@echo "Agent as Code - Build Information"
	@echo "================================="
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"
	@echo ""
	@echo "Go Information:"
	@echo "  Version: $(shell go version)"
	@echo "  GOOS:    $(shell go env GOOS)"
	@echo "  GOARCH:  $(shell go env GOARCH)"
	@echo ""
	@echo "Python Information:"
	@echo "  Version: $(shell python --version)"
	@echo "  Pip:     $(shell pip --version)"
	@echo ""
	@echo "Build Targets:"
	@echo "  Platforms: $(GOOS_TARGETS)"
	@echo "  Architectures: $(GOARCH_TARGETS)"

# Quick development cycle
dev: clean build-go install ## Quick development cycle (clean, build Go, install)
	@echo "🚀 Development cycle completed!"
	@echo "Try: agent --help"

# Binary API Operations
# ====================

# Upload agent CLI binaries for distribution
upload: build-go ## Upload agent CLI binaries to registry for user installation
	@echo "📤 Uploading agent CLI binaries to registry..."
	@if [ -z "$(AGENT_REGISTRY_TOKEN)" ]; then \
		echo "❌ Error: AGENT_REGISTRY_TOKEN environment variable required"; \
		echo "   Set your token: export AGENT_REGISTRY_TOKEN=your_token_here"; \
		exit 1; \
	fi
	@cd tools/binary-uploader && go run main.go --version $(VERSION) --all-platforms --bin-dir ../../$(BUILD_DIR)
	@echo "✅ Agent CLI binaries uploaded successfully!"

# Upload specific platform binary
upload-platform: build-go ## Upload agent CLI binary for specific platform (requires PLATFORM and ARCH)
	@echo "📤 Uploading agent CLI binary for $(PLATFORM)/$(ARCH)..."
	@if [ -z "$(PLATFORM)" ] || [ -z "$(ARCH)" ]; then \
		echo "❌ Error: PLATFORM and ARCH variables required"; \
		echo "   Usage: make upload-platform PLATFORM=linux ARCH=amd64"; \
		exit 1; \
	fi
	@if [ -z "$(AGENT_REGISTRY_TOKEN)" ]; then \
		echo "❌ Error: AGENT_REGISTRY_TOKEN environment variable required"; \
		exit 1; \
	fi
	@cd tools/binary-uploader && go run main.go --version $(VERSION) --platform $(PLATFORM) --arch $(ARCH) --bin-dir ../../$(BUILD_DIR)

# Test installation (for development)
test-install: ## Test agent CLI installation using install script
	@echo "🧪 Testing agent CLI installation..."
	@./scripts/install.sh --install-dir ./test-install --registry $(REGISTRY_URL)
	@echo "✅ Installation test completed!"

# Create install script for distribution
create-install-script: ## Create install.sh script for distribution
	@echo "📝 Creating installation script..."
	@cp scripts/install.sh install.sh
	@echo "✅ install.sh created for distribution"

# Release workflow - Upload agent CLI binaries for distribution
release: clean test-go build upload create-install-script ## Complete release workflow (clean, test, build, upload agent CLI)
	@echo "🚀 Agent CLI Release $(VERSION) completed successfully!"
	@echo ""
	@echo "Release Summary:"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Date:    $(DATE)"
	@echo ""
	@echo "Agent CLI is now available for installation:"
	@echo "  • Direct install: curl -L $(REGISTRY_URL)/install.sh | sh"
	@echo "  • Python package: pip install agent-as-code==$(VERSION)"
	@echo "  • Registry: $(REGISTRY_URL)/binary/releases/agent-as-code/"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Test installation: make test-install"
	@echo "  2. Create git tag: git tag v$(VERSION) && git push --tags"
	@echo "  3. Update PyPI package"
	@echo "  4. Announce release"

# Dry run upload (for testing)
upload-dry-run: build-go ## Dry run upload (shows what would be uploaded)
	@echo "🔍 Dry run upload (no actual upload)..."
	@cd tools/binary-uploader && go run main.go --version $(VERSION) --all-platforms --bin-dir ../../$(BUILD_DIR) --dry-run

# Build and test scripts
build-scripts: ## Make all scripts executable
	@echo "🔧 Making scripts executable..."
	@chmod +x scripts/*.sh
	@echo "✅ Scripts are now executable!"