# Parser

The AaC Parser is the core component responsible for reading, validating, and processing agent.yaml configurations. It ensures that agent definitions are syntactically correct, semantically valid, and ready for building.

## Overview

The parser performs several critical functions:
- **Syntax Validation**: Ensures agent.yaml follows correct YAML format
- **Semantic Validation**: Validates configuration logic and relationships
- **Dependency Resolution**: Checks package compatibility and versions
- **Security Validation**: Ensures secure configurations
- **Schema Compliance**: Validates against the AaC specification

## Architecture

The parser is built using Go's `gopkg.in/yaml.v3` package and provides a clean, structured API for parsing agent configurations.

## Core Data Structures

### AgentSpec
The main structure representing a complete agent configuration:

```go
type AgentSpec struct {
    APIVersion string            `yaml:"apiVersion"`
    Kind       string            `yaml:"kind"`
    Metadata   AgentMetadata     `yaml:"metadata"`
    Spec       AgentSpecDetails  `yaml:"spec"`
}
```

### AgentMetadata
Contains agent identification and metadata:

```go
type AgentMetadata struct {
    Name        string            `yaml:"name"`
    Version     string            `yaml:"version,omitempty"`
    Description string            `yaml:"description,omitempty"`
    Author      string            `yaml:"author,omitempty"`
    Tags        []string          `yaml:"tags,omitempty"`
    Labels      map[string]string `yaml:"labels,omitempty"`
}
```

### AgentSpecDetails
Contains the core agent specification:

```go
type AgentSpecDetails struct {
    Runtime      string                 `yaml:"runtime"`
    Model        ModelConfig            `yaml:"model"`
    Capabilities []string               `yaml:"capabilities,omitempty"`
    Dependencies []string               `yaml:"dependencies,omitempty"`
    Environment  []EnvironmentVar       `yaml:"environment,omitempty"`
    Ports        []PortConfig           `yaml:"ports,omitempty"`
    Volumes      []VolumeConfig         `yaml:"volumes,omitempty"`
    HealthCheck  *HealthCheckConfig     `yaml:"healthCheck,omitempty"`
    Resources    *ResourceConfig        `yaml:"resources,omitempty"`
    Config       map[string]interface{} `yaml:"config,omitempty"`
}
```

### ModelConfig
Defines the AI model configuration:

```go
type ModelConfig struct {
    Provider string                 `yaml:"provider"`
    Name     string                 `yaml:"name"`
    Config   map[string]interface{} `yaml:"config,omitempty"`
}
```

### EnvironmentVar
Represents environment variables:

```go
type EnvironmentVar struct {
    Name  string `yaml:"name"`
    Value string `yaml:"value,omitempty"`
    From  string `yaml:"from,omitempty"` // For secrets/configmaps
}
```

### PortConfig
Defines port mappings:

```go
type PortConfig struct {
    Container int    `yaml:"container"`
    Host      int    `yaml:"host,omitempty"`
    Protocol  string `yaml:"protocol,omitempty"`
}
```

### VolumeConfig
Defines volume mounts:

```go
type VolumeConfig struct {
    Source string `yaml:"source"`
    Target string `yaml:"target"`
    Type   string `yaml:"type,omitempty"`
}
```

### HealthCheckConfig
Defines health check configuration:

```go
type HealthCheckConfig struct {
    Command     []string `yaml:"command"`
    Interval    string   `yaml:"interval,omitempty"`
    Timeout     string   `yaml:"timeout,omitempty"`
    Retries     int      `yaml:"retries,omitempty"`
    StartPeriod string   `yaml:"startPeriod,omitempty"`
}
```

### ResourceConfig
Defines resource constraints:

```go
type ResourceConfig struct {
    Limits   ResourceLimits `yaml:"limits,omitempty"`
    Requests ResourceLimits `yaml:"requests,omitempty"`
}

type ResourceLimits struct {
    CPU    string `yaml:"cpu,omitempty"`
    Memory string `yaml:"memory,omitempty"`
}
```

## Supported Runtimes

The parser validates against these supported runtime environments:

- `python` - Python runtime
- `nodejs` - Node.js runtime  
- `go` - Go runtime
- `rust` - Rust runtime
- `java` - Java runtime

## Validation Rules

### Required Fields
The parser enforces these required fields:

1. **apiVersion**: Must be present (e.g., "agent.dev/v1")
2. **kind**: Must be exactly "Agent"
3. **metadata.name**: Agent name is required
4. **spec.runtime**: Runtime environment is required
5. **spec.model.provider**: Model provider is required
6. **spec.model.name**: Model name is required

### Port Validation
- Container ports must be between 1-65535
- Host ports (if specified) must be between 1-65535
- Port numbers must be valid integers

### Runtime Validation
- Runtime must be one of the supported values
- Case-sensitive matching

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/pxkundu/agent-as-code/internal/parser"
)

func main() {
    // Create parser instance
    p := parser.New()
    
    // Parse from file
    spec, err := p.ParseFile("agent.yaml")
    if err != nil {
        fmt.Printf("Parse error: %v\n", err)
        return
    }
    
    // Use the parsed specification
    fmt.Printf("Agent: %s\n", spec.Metadata.Name)
    fmt.Printf("Runtime: %s\n", spec.Spec.Runtime)
}
```

### Parse from Bytes

```go
// Parse from byte data
yamlData := []byte(`
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: my-agent
spec:
  runtime: python
  model:
    provider: openai
    name: gpt-4
`)

spec, err := p.Parse(yamlData)
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
    return
}
```

### Find Agent File

```go
// Find agent.yaml in a directory
agentFile, err := p.FindAgentFile("/path/to/agent")
if err != nil {
    fmt.Printf("No agent.yaml found: %v\n", err)
    return
}

fmt.Printf("Found agent file: %s\n", agentFile)
```

## File Discovery

The parser automatically searches for agent configuration files in the following order:

1. `agent.yaml`
2. `agent.yml` 
3. `Agent.yaml`
4. `Agent.yml`

## Error Handling

The parser provides detailed error messages for common issues:

- **File not found**: "no agent.yaml file found in {directory}"
- **Invalid YAML**: "failed to parse YAML: {yaml error}"
- **Missing required field**: "{field} is required"
- **Invalid runtime**: "invalid runtime '{runtime}'. Valid runtimes: {list}"
- **Invalid port**: "invalid container port {port} at index {index}"

## Best Practices

### 1. Configuration Structure
- Use consistent indentation
- Group related configurations together
- Use meaningful names and descriptions

### 2. Validation
- Always validate configurations before building
- Test with different runtime environments
- Verify port configurations

### 3. Error Handling
- Check for parse errors in your code
- Provide user-friendly error messages
- Log validation failures for debugging

## Integration

The parser integrates with other AaC components:

- **Builder**: Provides validated configurations for building
- **Runtime**: Supplies runtime configuration for execution
- **Registry**: Validates metadata for publishing

## Performance Considerations

- **File I/O**: Uses efficient file reading with `ioutil.ReadFile`
- **YAML Parsing**: Leverages optimized YAML parsing library
- **Validation**: Performs validation in a single pass
- **Memory**: Minimal memory footprint for large configurations

## Troubleshooting

### Common Issues

1. **YAML Syntax Errors**
   - Check indentation (use spaces, not tabs)
   - Verify YAML syntax with online validators
   - Ensure proper quoting for special characters

2. **Missing Required Fields**
   - Verify all required fields are present
   - Check field names for typos
   - Ensure proper nesting structure

3. **Invalid Runtime**
   - Use exact runtime names (case-sensitive)
   - Check supported runtime list
   - Verify runtime version format

### Debug Mode

Enable verbose logging to see detailed parsing information:

```bash
export AAC_LOG_LEVEL=debug
agent build -t my-agent .
```

## See Also

- [Agent Configuration](./agent-configuration.md) - Complete configuration reference
- [Builder](./builder.md) - Building agents from configurations
- [Runtime](./runtime.md) - Executing parsed configurations
- [CLI Overview](./cli-overview.md) - Command-line interface usage 