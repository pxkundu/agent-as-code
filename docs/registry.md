# Registry

The Agent Registry is a centralized storage and distribution system for AI agents. It enables sharing, versioning, and discovery of agents across teams and organizations.

## Overview

The registry provides:
- **Agent Storage**: Secure storage for agent packages and metadata
- **Version Management**: Semantic versioning and release tracking
- **Distribution**: Easy sharing and discovery of agents
- **Access Control**: Authentication and authorization
- **Search & Discovery**: Find agents by capabilities, tags, and metadata

## Architecture

The registry integrates with:
- **Docker Registry**: Standard Docker registry compatibility
- **Agent Registry**: Custom agent-specific registry
- **Authentication**: Token-based authentication
- **Storage**: Distributed storage backend

## Core Components

### Registry Struct
```go
type Registry struct {
    dockerClient *client.Client
    registryURL  string
    authToken    string
}
```

### PushOptions
```go
type PushOptions struct {
    Image    string // Image to push
    Registry string // Target registry
    AllTags  bool   // Push all tags
}
```

### PullOptions
```go
type PullOptions struct {
    Image    string // Image to pull
    Registry string // Source registry
    Quiet    bool   // Quiet mode
}
```

### ListOptions
```go
type ListOptions struct {
    Filter []string // Filter criteria
    All    bool     // Show all images
}
```

### PushResult
```go
type PushResult struct {
    Repository  string // Image repository
    Tag         string // Image tag
    Digest      string // Image digest
    Size        string // Image size
    RegistryURL string // Registry URL
}
```

### PullResult
```go
type PullResult struct {
    ImageID     string // Docker image ID
    Size        string // Image size
    Digest      string // Image digest
    RegistryURL string // Registry URL
}
```

### ImageInfo
```go
type ImageInfo struct {
    ID         string    // Image ID
    Repository string    // Repository name
    Tag        string    // Image tag
    Created    time.Time // Creation time
    Size       int64     // Image size in bytes
}
```

## Registry Operations

### Image Push

#### Basic Push
```bash
# Push to default registry
agent push my-agent:latest

# Push to specific registry
agent push my-agent:latest --registry myregistry.com

# Push all tags
agent push my-agent --all-tags
```

#### Push with Options
```go
package main

import (
    "fmt"
    "github.com/pxkundu/agent-as-code/internal/registry"
)

func main() {
    // Create registry instance
    r := registry.New()
    
    // Push options
    options := &registry.PushOptions{
        Image:    "my-agent:latest",
        Registry: "myregistry.com",
        AllTags:  false,
    }
    
    // Push image
    result, err := r.Push(options)
    if err != nil {
        fmt.Printf("Push failed: %v\n", err)
        return
    }
    
    fmt.Printf("Push successful: %s\n", result.Repository)
    fmt.Printf("Tag: %s\n", result.Tag)
    fmt.Printf("Digest: %s\n", result.Digest)
}
```

#### Push Process
1. **Validation**: Verify local image exists
2. **Authentication**: Authenticate with registry
3. **Upload**: Stream image layers to registry
4. **Verification**: Confirm successful upload
5. **Metadata**: Update registry metadata

### Image Pull

#### Basic Pull
```bash
# Pull from default registry
agent pull my-agent:latest

# Pull from specific registry
agent pull my-agent:latest --registry myregistry.com

# Quiet pull (minimal output)
agent pull my-agent:latest --quiet
```

#### Pull with Options
```go
// Pull options
options := &registry.PullOptions{
    Image:    "my-agent:latest",
    Registry: "myregistry.com",
    Quiet:    false,
}

// Pull image
result, err := r.Pull(options)
if err != nil {
    fmt.Printf("Pull failed: %v\n", err)
    return
}

fmt.Printf("Pull successful: %s\n", result.ImageID)
fmt.Printf("Size: %s\n", result.Size)
```

#### Pull Process
1. **Authentication**: Authenticate with registry
2. **Download**: Stream image layers from registry
3. **Verification**: Verify image integrity
4. **Storage**: Store locally with metadata
5. **Cleanup**: Remove temporary files

### Image Listing

#### List Local Images
```bash
# List all images
agent images

# List with filters
agent images --filter "my-agent"

# Show all images (including untagged)
agent images --all
```

#### List with Options
```go
// List options
options := &registry.ListOptions{
    Filter: []string{"my-agent"},
    All:    false,
}

// List images
images, err := r.ListLocal(options)
if err != nil {
    fmt.Printf("List failed: %v\n", err)
    return
}

for _, image := range images {
    fmt.Printf("ID: %s\n", image.ID)
    fmt.Printf("Repository: %s\n", image.Repository)
    fmt.Printf("Tag: %s\n", image.Tag)
    fmt.Printf("Size: %s\n", formatSize(image.Size))
    fmt.Printf("Created: %s\n", image.Created.Format("2006-01-02 15:04:05"))
    fmt.Println("---")
}
```

## Registry Types

### Docker Registry
Standard Docker registry compatibility:
- **Docker Hub**: Public Docker registry
- **ECR**: AWS Elastic Container Registry
- **ACR**: Azure Container Registry
- **GCR**: Google Container Registry

### Agent Registry
Custom agent-specific registry:
- **myagentregistry.com**: Official agent registry
- **Custom Registries**: Self-hosted registries
- **Enterprise Registries**: Corporate registries

## Authentication

### Environment Variables
```bash
# Registry URL
export AGENT_REGISTRY_URL="https://myregistry.com"

# Authentication token
export AGENT_REGISTRY_TOKEN="your-auth-token"

# Docker credentials (for Docker registries)
export DOCKER_USERNAME="your-username"
export DOCKER_PASSWORD="your-password"
```

### Authentication Methods
1. **Token Authentication**: Bearer token for API access
2. **Username/Password**: Basic authentication
3. **OAuth2**: OAuth2 token flow
4. **Service Account**: Kubernetes service accounts

### Security Best Practices
- Use environment variables for secrets
- Rotate tokens regularly
- Implement least privilege access
- Monitor authentication attempts

## Image Management

### Image Tagging
```bash
# Tag local image
docker tag my-agent:latest myregistry.com/my-agent:v1.0.0

# Tag with specific registry
docker tag my-agent:latest myregistry.com/team/my-agent:latest

# Tag multiple versions
docker tag my-agent:latest my-agent:v1.0.0
docker tag my-agent:latest my-agent:latest
```

### Image Versioning
```bash
# Semantic versioning
agent push my-agent:1.0.0
agent push my-agent:1.0.1
agent push my-agent:1.1.0
agent push my-agent:2.0.0

# Latest tag
agent push my-agent:latest

# Development tags
agent push my-agent:dev
agent push my-agent:staging
```

### Image Cleanup
```bash
# Remove local images
docker rmi my-agent:old-version

# Remove untagged images
docker image prune

# Remove all unused images
docker image prune -a

# Remove specific registry images
docker rmi myregistry.com/my-agent:old-version
```

## Registry Configuration

### Registry Profiles
```yaml
# ~/.aac/config.yaml
registries:
  default:
    url: "https://myregistry.com"
    auth:
      type: "token"
      token: "${AGENT_REGISTRY_TOKEN}"
  
  production:
    url: "https://prod.registry.com"
    auth:
      type: "oauth2"
      client_id: "${OAUTH_CLIENT_ID}"
      client_secret: "${OAUTH_CLIENT_SECRET}"
  
  development:
    url: "https://dev.registry.com"
    auth:
      type: "username"
      username: "${DEV_USERNAME}"
      password: "${DEV_PASSWORD}"
```

### Registry Selection
```bash
# Use specific registry profile
agent push my-agent:latest --profile production

# Override registry URL
agent push my-agent:latest --registry https://custom.registry.com

# Use default registry
agent push my-agent:latest
```

## Search and Discovery

### Image Search
```bash
# Search by name
agent search my-agent

# Search by tag
agent search my-agent:latest

# Search by capability
agent search --capability "text-generation"

# Search by author
agent search --author "ai-team"
```

### Search Filters
- **Name**: Exact or partial name matching
- **Tag**: Specific version or tag
- **Capability**: Agent capabilities
- **Author**: Image creator
- **Date**: Creation or update date
- **Size**: Image size range

### Search Results
```json
{
  "images": [
    {
      "name": "my-agent",
      "tag": "latest",
      "description": "Advanced text generation agent",
      "author": "ai-team",
      "capabilities": ["text-generation", "chat"],
      "size": "150MB",
      "created": "2024-01-01T00:00:00Z",
      "downloads": 150,
      "rating": 4.8
    }
  ],
  "total": 1,
  "page": 1,
  "per_page": 10
}
```

## Performance Optimization

### Caching Strategies
```yaml
# Registry caching configuration
cache:
  enabled: true
  max_size: "10GB"
  ttl: "24h"
  strategies:
    - layer_caching: true
    - metadata_caching: true
    - search_caching: true
```

### Network Optimization
- **CDN Integration**: Content delivery networks
- **Mirror Registries**: Geographic mirrors
- **Connection Pooling**: Reuse connections
- **Compression**: Layer compression

### Storage Optimization
- **Layer Deduplication**: Shared layers
- **Compression**: Image compression
- **Cleanup Policies**: Automatic cleanup
- **Storage Tiers**: Hot/cold storage

## Monitoring and Metrics

### Registry Metrics
```bash
# View registry metrics
agent registry metrics

# Monitor specific metrics
agent registry metrics --metric "push_count" --period "24h"

# Custom metric queries
agent registry metrics --query "rate(push_operations_total[5m])"
```

### Key Metrics
- **Push Operations**: Images pushed per time period
- **Pull Operations**: Images pulled per time period
- **Storage Usage**: Registry storage consumption
- **Authentication**: Login attempts and failures
- **Performance**: Response times and throughput

### Health Monitoring
```bash
# Check registry health
agent registry health

# Detailed health check
agent registry health --detailed

# Health check with timeout
agent registry health --timeout 30s
```

## Error Handling

### Common Errors
1. **Authentication Failed**: "authentication failed"
2. **Image Not Found**: "image not found in registry"
3. **Network Error**: "failed to connect to registry"
4. **Storage Full**: "insufficient storage space"
5. **Rate Limited**: "rate limit exceeded"

### Error Recovery
```go
// Retry logic for registry operations
func pushWithRetry(options *PushOptions, maxRetries int) (*PushResult, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := r.Push(options)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Exponential backoff
        backoff := time.Duration(i+1) * time.Second
        time.Sleep(backoff)
    }
    
    return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}
```

### Troubleshooting
```bash
# Check registry connectivity
curl -I https://myregistry.com/v2/

# Verify authentication
agent registry auth test

# Check image existence
agent registry image info my-agent:latest

# View registry logs
agent registry logs --level error --since 1h
```

## Security Features

### Image Signing
```bash
# Sign image before push
agent push my-agent:latest --sign

# Verify image signature on pull
agent pull my-agent:latest --verify-signature

# View image signatures
agent registry image signatures my-agent:latest
```

### Vulnerability Scanning
```bash
# Scan image for vulnerabilities
agent push my-agent:latest --security-scan

# View scan results
agent registry image scan-results my-agent:latest

# Block vulnerable images
agent registry policy set --block-vulnerabilities high
```

### Access Control
```bash
# Set image visibility
agent push my-agent:latest --visibility private

# Manage permissions
agent registry permission add my-agent:latest --user john --permission read

# View access logs
agent registry access-logs --user john --since 24h
```

## Integration

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Build and Push Agent
  run: |
    agent build -t my-agent:${{ github.sha }} .
    agent push my-agent:${{ github.sha }}
    agent push my-agent:latest
```

### Kubernetes Integration
```yaml
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-agent
spec:
  template:
    spec:
      containers:
      - name: agent
        image: myregistry.com/my-agent:latest
        imagePullPolicy: Always
```

### Docker Compose Integration
```yaml
# docker-compose.yml
version: '3.8'
services:
  agent:
    image: myregistry.com/my-agent:latest
    ports:
      - "8080:8080"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
```

## Best Practices

### 1. Image Management
- Use semantic versioning
- Tag images consistently
- Clean up old images
- Monitor image sizes

### 2. Security
- Sign all images
- Scan for vulnerabilities
- Use private registries
- Implement access controls

### 3. Performance
- Use caching effectively
- Optimize image sizes
- Monitor registry performance
- Implement cleanup policies

### 4. Operations
- Monitor registry health
- Set up alerting
- Document procedures
- Regular backups

## Troubleshooting

### Registry Issues
1. **Authentication Problems**
   - Verify token validity
   - Check environment variables
   - Test with curl commands

2. **Network Issues**
   - Check firewall rules
   - Verify DNS resolution
   - Test connectivity

3. **Storage Issues**
   - Check available space
   - Verify storage permissions
   - Review cleanup policies

### Debug Commands
```bash
# Enable debug mode
export AAC_LOG_LEVEL=debug

# Test registry connection
agent registry ping

# View registry configuration
agent registry config show

# Test authentication
agent registry auth test

# View detailed logs
agent registry logs --level debug
```

## See Also

- [Parser](./parser.md) - Configuration parsing
- [Builder](./builder.md) - Container building
- [Runtime](./runtime.md) - Agent execution
- [CLI Overview](./cli-overview.md) - Command-line usage
- [Agent Configuration](./agent-configuration.md) - Configuration reference 