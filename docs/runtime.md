# Runtime

The AaC Runtime manages the execution of agent containers. It provides a comprehensive container orchestration system with support for different execution strategies and resource management.

## Overview

The runtime system handles:
- **Container Orchestration**: Docker-based container management
- **Runtime Configuration**: Dynamic configuration and environment setup
- **Monitoring and Logging**: Real-time container monitoring and log streaming
- **Resource Management**: CPU, memory, and network resource allocation
- **Scaling Strategies**: Horizontal and vertical scaling capabilities

## Architecture

The runtime integrates with:
- **Docker Engine**: For container lifecycle management
- **Parser**: For configuration validation
- **Builder**: For image availability verification

## Core Components

### Runtime Struct
```go
type Runtime struct {
    dockerClient *client.Client
}
```

### RunOptions
```go
type RunOptions struct {
    Image       string   // Container image name
    Ports       []string // Port mappings (host:container)
    Environment []string // Environment variables
    Detach      bool     // Run in background
    Name        string   // Container name
    Volumes     []string // Volume mounts
    Interactive bool     // Interactive mode
}
```

### ContainerInfo
```go
type ContainerInfo struct {
    ID    string        // Container ID
    Name  string        // Container name
    Ports []PortMapping // Port mappings
}
```

### PortMapping
```go
type PortMapping struct {
    Host      string // Host port
    Container string // Container port
    Protocol  string // Protocol (tcp/udp)
}
```

## Runtime Strategies

### Docker Runtime (Default)
- **Container Engine**: Docker Engine
- **Orchestration**: Native Docker API
- **Networking**: Docker bridge network
- **Storage**: Docker volumes and bind mounts

### Kubernetes Runtime (Experimental)
- **Orchestration**: Kubernetes API
- **Scaling**: Horizontal Pod Autoscaler
- **Networking**: Service mesh integration
- **Storage**: Persistent volumes

### Local Process Runtime
- **Execution**: Direct process execution
- **Isolation**: Process-level isolation
- **Networking**: Local port binding
- **Storage**: File system access

## Container Management

### Starting Containers
```go
package main

import (
    "fmt"
    "github.com/pxkundu/agent-as-code/internal/runtime"
)

func main() {
    // Create runtime instance
    r := runtime.New()
    
    // Run options
    options := &runtime.RunOptions{
        Image: "my-agent:latest",
        Ports: []string{"8080:8080", "9090:9090"},
        Environment: []string{
            "OPENAI_API_KEY=your-key",
            "LOG_LEVEL=INFO",
        },
        Name: "my-agent-instance",
        Volumes: []string{"./data:/app/data"},
        Interactive: false,
    }
    
    // Start container
    container, err := r.Run(options)
    if err != nil {
        fmt.Printf("Failed to start container: %v\n", err)
        return
    }
    
    fmt.Printf("Container started: %s\n", container.ID)
    fmt.Printf("Container name: %s\n", container.Name)
}
```

### Container Lifecycle
1. **Validation**: Verify image exists locally
2. **Creation**: Create container with configuration
3. **Startup**: Start container and wait for readiness
4. **Monitoring**: Track container status and health
5. **Cleanup**: Handle container termination

## Port Management

### Port Mapping Syntax
```bash
# Basic port mapping
8080:8080

# Host port different from container
80:8080

# Protocol specification
8080:8080/tcp
8080:8080/udp

# Multiple ports
8080:8080,9090:9090
```

### Port Configuration Examples
```yaml
# agent.yaml
spec:
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
    - container: 9090
      host: 9090
      protocol: tcp
```

### Port Validation
- Port numbers: 1-65535
- Protocol support: TCP, UDP
- Host port availability check
- Container port exposure

## Environment Configuration

### Environment Variables
```yaml
# agent.yaml
spec:
  environment:
    - name: OPENAI_API_KEY
      value: "your-api-key"
    - name: LOG_LEVEL
      value: "INFO"
    - name: MODEL_NAME
      from: "config"
```

### Runtime Environment
```bash
# Set environment variables
export OPENAI_API_KEY="your-key"
export LOG_LEVEL="DEBUG"

# Run with environment
agent run -e OPENAI_API_KEY=${OPENAI_API_KEY} my-agent:latest
```

### Secret Management
```yaml
# agent.yaml
spec:
  environment:
    - name: DATABASE_PASSWORD
      from: "secret"
      secretName: "db-secret"
      secretKey: "password"
```

## Volume Management

### Volume Mounts
```yaml
# agent.yaml
spec:
  volumes:
    - source: ./data
      target: /app/data
      type: bind
    - source: model-cache
      target: /app/models
      type: volume
```

### Volume Types
- **Bind Mounts**: Host directory to container
- **Named Volumes**: Docker-managed volumes
- **Temporary Volumes**: Ephemeral storage
- **Config Maps**: Configuration data

### Volume Examples
```bash
# Bind mount
agent run -v $(pwd)/data:/app/data my-agent:latest

# Named volume
agent run -v model-cache:/app/models my-agent:latest

# Multiple volumes
agent run -v ./data:/app/data -v ./config:/app/config my-agent:latest
```

## Health Checks

### Health Check Configuration
```yaml
# agent.yaml
spec:
  healthCheck:
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
```

### Health Check Types
- **HTTP**: HTTP endpoint health checks
- **Command**: Custom command execution
- **TCP**: Port availability checks
- **File**: File existence checks

### Health Check Implementation
```python
# Python health check example
from flask import Flask

app = Flask(__name__)

@app.route('/health')
def health_check():
    return {"status": "healthy", "timestamp": time.time()}

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
```

## Resource Management

### Resource Limits
```yaml
# agent.yaml
spec:
  resources:
    limits:
      cpu: "2"
      memory: "2Gi"
    requests:
      cpu: "500m"
      memory: "512Mi"
```

### Resource Monitoring
```go
// Monitor container resources
func monitorResources(containerID string) {
    // Get container stats
    stats, err := dockerClient.ContainerStats(context.Background(), containerID, false)
    if err != nil {
        return
    }
    
    // Parse resource usage
    var resourceUsage struct {
        CPU    float64
        Memory int64
    }
    
    // Process stats...
}
```

## Container Operations

### Starting Containers
```bash
# Basic start
agent run my-agent:latest

# With port mapping
agent run -p 8080:8080 my-agent:latest

# With environment
agent run -e API_KEY=value my-agent:latest

# With volumes
agent run -v ./data:/app/data my-agent:latest

# Interactive mode
agent run -i -t my-agent:latest
```

### Stopping Containers
```go
// Stop container
err := r.Stop(containerID)
if err != nil {
    fmt.Printf("Failed to stop container: %v\n", err)
    return
}
```

### Container Logs
```go
// Stream container logs
err := r.StreamLogs(containerID)
if err != nil {
    fmt.Printf("Failed to stream logs: %v\n", err)
    return
}
```

### Container Listing
```go
// List running containers
containers, err := r.List()
if err != nil {
    fmt.Printf("Failed to list containers: %v\n", err)
    return
}

for _, container := range containers {
    fmt.Printf("ID: %s, Name: %s\n", container.ID, container.Name)
}
```

## Networking

### Network Configuration
```yaml
# agent.yaml
spec:
  network:
    mode: bridge
    ports:
      - 8080:8080
    dns:
      - 8.8.8.8
      - 8.8.4.4
```

### Network Modes
- **Bridge**: Default Docker networking
- **Host**: Host network namespace
- **None**: No networking
- **Custom**: User-defined networks

### Service Discovery
```yaml
# agent.yaml
spec:
  services:
    - name: database
      port: 5432
      protocol: tcp
    - name: cache
      port: 6379
      protocol: tcp
```

## Monitoring and Observability

### Container Metrics
- **CPU Usage**: Per-container CPU consumption
- **Memory Usage**: Memory allocation and usage
- **Network I/O**: Network traffic statistics
- **Disk I/O**: Storage access patterns

### Log Management
```go
// Configure logging
logOptions := types.ContainerLogsOptions{
    ShowStdout: true,
    ShowStderr: true,
    Follow:     true,
    Timestamps: true,
    Tail:       "100",
}

// Get logs
reader, err := dockerClient.ContainerLogs(ctx, containerID, logOptions)
```

### Health Monitoring
```go
// Health check loop
func healthMonitor(containerID string) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // Check container health
        if !isHealthy(containerID) {
            // Handle unhealthy container
            handleUnhealthyContainer(containerID)
        }
    }
}
```

## Scaling and Orchestration

### Horizontal Scaling
```yaml
# agent.yaml
spec:
  scaling:
    replicas: 3
    minReplicas: 1
    maxReplicas: 10
    metrics:
      - type: Resource
        resource:
          name: cpu
          target:
            type: Utilization
            averageUtilization: 70
```

### Vertical Scaling
```yaml
# agent.yaml
spec:
  resources:
    limits:
      cpu: "4"
      memory: "8Gi"
    requests:
      cpu: "1"
      memory: "2Gi"
```

### Auto-scaling
```yaml
# agent.yaml
spec:
  autoscaling:
    enabled: true
    targetCPUUtilizationPercentage: 70
    targetMemoryUtilizationPercentage: 80
    scaleDownDelay: 300s
    scaleUpDelay: 60s
```

## Security

### Security Context
```yaml
# agent.yaml
spec:
  security:
    runAsUser: 1000
    runAsGroup: 1000
    readOnlyRootFilesystem: true
    capabilities:
      drop:
        - ALL
      add:
        - NET_BIND_SERVICE
```

### Network Security
```yaml
# agent.yaml
spec:
  network:
    security:
      allowExternalConnections: false
      allowedHosts:
        - "api.openai.com"
        - "api.anthropic.com"
```

### Secret Management
```yaml
# agent.yaml
spec:
  secrets:
    - name: api-key
      from: "kubernetes-secret"
      secretName: "agent-secrets"
      secretKey: "openai-api-key"
```

## Error Handling

### Common Runtime Errors
1. **Image Not Found**: "image 'my-agent:latest' not found locally"
2. **Port Conflicts**: "port 8080 is already in use"
3. **Resource Limits**: "insufficient memory"
4. **Network Issues**: "failed to create network"

### Error Recovery
```go
// Retry logic for container operations
func runWithRetry(options *RunOptions, maxRetries int) (*ContainerInfo, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        container, err := r.Run(options)
        if err == nil {
            return container, nil
        }
        
        lastErr = err
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}
```

## Performance Optimization

### Container Optimization
- **Image Size**: Use minimal base images
- **Layer Caching**: Optimize Docker layers
- **Resource Allocation**: Right-size resource limits
- **Network Optimization**: Use appropriate network modes

### Monitoring Optimization
- **Metrics Collection**: Efficient metrics gathering
- **Log Rotation**: Implement log rotation
- **Resource Tracking**: Monitor resource usage
- **Performance Profiling**: Profile container performance

## Best Practices

### 1. Container Design
- Use minimal base images
- Implement proper health checks
- Handle graceful shutdowns
- Optimize layer caching

### 2. Resource Management
- Set appropriate resource limits
- Monitor resource usage
- Implement auto-scaling
- Use resource quotas

### 3. Security
- Run containers as non-root
- Implement network policies
- Use secrets for sensitive data
- Regular security updates

### 4. Monitoring
- Implement comprehensive logging
- Monitor container health
- Track performance metrics
- Set up alerting

## Troubleshooting

### Common Issues
1. **Container Won't Start**
   - Check image availability
   - Verify port availability
   - Check resource limits
   - Review container logs

2. **Performance Issues**
   - Monitor resource usage
   - Check network performance
   - Review storage I/O
   - Analyze container metrics

3. **Network Issues**
   - Verify port mappings
   - Check firewall rules
   - Test network connectivity
   - Review DNS configuration

### Debug Commands
```bash
# Check container status
docker ps -a

# View container logs
docker logs <container-id>

# Inspect container
docker inspect <container-id>

# Check resource usage
docker stats <container-id>

# Execute in container
docker exec -it <container-id> /bin/bash
```

## See Also

- [Parser](./parser.md) - Configuration parsing
- [Builder](./builder.md) - Container building
- [Registry](./registry.md) - Image distribution
- [CLI Overview](./cli-overview.md) - Command-line usage
- [Agent Configuration](./agent-configuration.md) - Configuration reference