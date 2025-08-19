# Builder

The AaC Builder is responsible for transforming agent configurations into deployable container images. It handles the complete build process from parsing agent.yaml to generating Docker images.

## Overview

The builder performs these key functions:
- **Configuration Parsing**: Loads and validates agent.yaml
- **Dockerfile Generation**: Creates optimized Dockerfiles based on runtime
- **Image Building**: Uses Docker BuildKit for efficient builds
- **Context Management**: Creates optimized build contexts
- **Multi-Platform Support**: Handles different target platforms

## Architecture

The builder integrates with:
- **Parser**: For configuration validation
- **Docker Engine**: For image building
- **BuildKit**: For optimized builds

## Core Components

### Builder Struct
```go
type Builder struct {
    parser       *parser.Parser
    dockerClient *client.Client
}
```

### BuildOptions
```go
type BuildOptions struct {
    Path     string  // Build context path
    Tag      string  // Image tag
    NoCache  bool    // Disable build cache
    Push     bool    // Push after build
    Platform string  // Target platform
}
```

### BuildResult
```go
type BuildResult struct {
    ImageID string   // Docker image ID
    Size    string   // Image size (human readable)
    Tags    []string // Applied tags
}
```

## Build Process

### 1. Context Validation
The builder first validates the build context:
- Checks for agent.yaml existence
- Parses and validates configuration
- Ensures all required files are present

### 2. Dockerfile Generation
Generates runtime-specific Dockerfiles:

#### Python Runtime
```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Environment variables
ENV OPENAI_API_KEY=your_key

# Expose ports
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD ["curl", "http://localhost:8080/health"]

# Run the application
CMD ["python", "main.py"]
```

#### Node.js Runtime
```dockerfile
FROM node:18-slim

WORKDIR /app

# Install Node.js dependencies
COPY package*.json .
RUN npm ci --only=production

# Copy application code
COPY . .

# Environment variables
ENV NODE_ENV=production

# Expose ports
EXPOSE 8080

# Run the application
CMD ["node", "index.js"]
```

#### Go Runtime
```dockerfile
FROM golang:1.21-alpine AS builder
FROM alpine:latest

WORKDIR /app

# Copy application code
COPY . .

# Expose ports
EXPOSE 8080

# Run the application
CMD ["./app"]
```

### 3. Build Context Creation
Creates optimized tar archives for Docker builds:
- Excludes hidden files and directories
- Includes all source code and dependencies
- Optimizes for layer caching

### 4. Image Building
Uses Docker API for building:
- Streams build output in real-time
- Handles build errors gracefully
- Provides progress feedback

## Supported Runtimes

### Python
- **Base Image**: `python:3.11-slim`
- **Dependencies**: `requirements.txt`
- **Entry Point**: `python main.py`
- **Package Manager**: pip

### Node.js
- **Base Image**: `node:18-slim`
- **Dependencies**: `package.json`
- **Entry Point**: `node index.js`
- **Package Manager**: npm

### Go
- **Base Image**: `golang:1.21-alpine` (build), `alpine:latest` (runtime)
- **Dependencies**: Go modules
- **Entry Point**: `./app`
- **Build Process**: Multi-stage build

## Build Configuration

### Environment Variables
```yaml
spec:
  environment:
    - name: OPENAI_API_KEY
      value: "your-api-key"
    - name: LOG_LEVEL
      value: "INFO"
```

### Port Configuration
```yaml
spec:
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
```

### Health Checks
```yaml
spec:
  healthCheck:
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
```

### Dependencies
```yaml
spec:
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
```

## Usage Examples

### Basic Build
```bash
# Build from current directory
agent build -t my-agent:latest .

# Build with specific tag
agent build -t my-agent:v1.0.0 .

# Build without cache
agent build -t my-agent:latest . --no-cache
```

### Build with Options
```go
package main

import (
    "fmt"
    "github.com/pxkundu/agent-as-code/internal/builder"
)

func main() {
    // Create builder instance
    b := builder.New()
    
    // Build options
    options := &builder.BuildOptions{
        Path:     "./my-agent",
        Tag:      "my-agent:latest",
        NoCache:  false,
        Push:     false,
        Platform: "linux/amd64",
    }
    
    // Build the agent
    result, err := b.Build(options)
    if err != nil {
        fmt.Printf("Build failed: %v\n", err)
        return
    }
    
    fmt.Printf("Build successful: %s\n", result.ImageID)
    fmt.Printf("Image size: %s\n", result.Size)
}
```

### Build Validation
```go
// Validate build context before building
err := b.ValidateContext("./my-agent")
if err != nil {
    fmt.Printf("Invalid context: %v\n", err)
    return
}
```

## Build Optimization

### Layer Caching
- Dependencies installed first (cached layer)
- Source code copied last (changes frequently)
- Multi-stage builds for complex agents

### Context Optimization
- Excludes unnecessary files
- Minimizes build context size
- Optimizes for Docker layer caching

### Dependency Management
- Pins dependency versions
- Uses production-only dependencies
- Minimizes image size

## Error Handling

### Common Build Errors
1. **Missing agent.yaml**: "no agent.yaml found"
2. **Invalid configuration**: "invalid agent.yaml"
3. **Docker not available**: "Docker client not available"
4. **Build context errors**: "failed to create build context"

### Troubleshooting
```bash
# Check Docker status
docker info

# Verify agent.yaml
agent build --validate-only .

# Enable verbose logging
export AAC_LOG_LEVEL=debug
agent build -t my-agent .
```

## Build Artifacts

### Generated Files
- `Dockerfile.agent`: Generated Dockerfile
- Container image with specified tag
- Build logs and progress information

### Image Information
- Image ID and digest
- Size and layer information
- Applied tags and labels

## Advanced Features

### Multi-Platform Builds
```bash
# Build for specific platform
agent build -t my-agent:latest . --platform linux/arm64

# Build for multiple platforms
agent build -t my-agent:latest . --platform linux/amd64,linux/arm64
```

### Build Profiles
```yaml
# build-profiles.yaml
profiles:
  development:
    base_image: python:3.11-slim
    debug: true
    
  production:
    base_image: python:3.11-slim
    debug: false
    optimization: true
```

### Custom Base Images
```yaml
spec:
  runtime: python
  build:
    base_image: custom/python:3.11
    optimization: true
```

## Integration

### With Registry
```go
// Build and push
result, err := b.Build(options)
if err != nil {
    return err
}

// Push to registry
err = b.Push(options.Tag)
if err != nil {
    return err
}
```

### With Runtime
```go
// Build for runtime execution
result, err := b.Build(options)
if err != nil {
    return err
}

// Use with runtime
runtime := runtime.New()
runOptions := &runtime.RunOptions{
    Image: result.ImageID,
    Ports: []string{"8080:8080"},
}
```

## Performance Considerations

### Build Time Optimization
- Use layer caching effectively
- Minimize dependency changes
- Optimize build context

### Memory Usage
- Stream large files
- Use efficient tar creation
- Monitor Docker daemon resources

### Network Optimization
- Use local Docker registry for development
- Optimize dependency downloads
- Cache base images locally

## Best Practices

### 1. Build Optimization
- Use multi-stage builds for complex agents
- Minimize layer count
- Optimize dependency installation

### 2. Security
- Use minimal base images
- Scan for vulnerabilities
- Validate dependencies

### 3. Reproducibility
- Pin dependency versions
- Use deterministic builds
- Document build environment

### 4. Monitoring
- Track build times
- Monitor image sizes
- Log build failures

## Troubleshooting

### Build Failures
1. **Docker Issues**
   - Check Docker daemon status
   - Verify Docker version compatibility
   - Check available disk space

2. **Configuration Issues**
   - Validate agent.yaml syntax
   - Check required fields
   - Verify runtime support

3. **Dependency Issues**
   - Check dependency versions
   - Verify package compatibility
   - Test dependency installation

### Debug Information
```bash
# Enable debug mode
export AAC_LOG_LEVEL=debug

# Show Docker info
docker info

# Check build context
ls -la

# Validate configuration
agent build --validate-only .
```

## See Also

- [Parser](./parser.md) - Configuration parsing
- [Runtime](./runtime.md) - Agent execution
- [Registry](./registry.md) - Image distribution
- [CLI Overview](./cli-overview.md) - Command-line usage
- [Agent Configuration](./agent-configuration.md) - Configuration reference 