package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AgentSpec represents the agent.yaml specification
type AgentSpec struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   AgentMetadata     `yaml:"metadata"`
	Spec       AgentSpecDetails  `yaml:"spec"`
}

// AgentMetadata contains agent metadata
type AgentMetadata struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Author      string            `yaml:"author,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// AgentSpecDetails contains the agent specification
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

// ModelConfig represents model configuration
type ModelConfig struct {
	Provider string                 `yaml:"provider"`
	Name     string                 `yaml:"name"`
	Config   map[string]interface{} `yaml:"config,omitempty"`
}

// EnvironmentVar represents an environment variable
type EnvironmentVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
	From  string `yaml:"from,omitempty"` // For secrets/configmaps
}

// PortConfig represents port configuration
type PortConfig struct {
	Container int    `yaml:"container"`
	Host      int    `yaml:"host,omitempty"`
	Protocol  string `yaml:"protocol,omitempty"`
}

// VolumeConfig represents volume configuration
type VolumeConfig struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Type   string `yaml:"type,omitempty"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Command     []string `yaml:"command"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"startPeriod,omitempty"`
}

// ResourceConfig represents resource constraints
type ResourceConfig struct {
	Limits   ResourceLimits `yaml:"limits,omitempty"`
	Requests ResourceLimits `yaml:"requests,omitempty"`
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// Parser handles agent.yaml parsing
type Parser struct{}

// New creates a new parser instance
func New() *Parser {
	return &Parser{}
}

// ParseFile parses an agent.yaml file
func (p *Parser) ParseFile(path string) (*AgentSpec, error) {
	// Read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read agent.yaml: %w", err)
	}
	
	return p.Parse(data)
}

// Parse parses agent.yaml content
func (p *Parser) Parse(data []byte) (*AgentSpec, error) {
	var spec AgentSpec
	
	// Parse YAML
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Validate the spec
	if err := p.Validate(&spec); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	return &spec, nil
}

// Validate validates the agent specification
func (p *Parser) Validate(spec *AgentSpec) error {
	// Check required fields
	if spec.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}
	
	if spec.Kind == "" {
		return fmt.Errorf("kind is required")
	}
	
	if spec.Kind != "Agent" {
		return fmt.Errorf("kind must be 'Agent', got '%s'", spec.Kind)
	}
	
	if spec.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}
	
	if spec.Spec.Runtime == "" {
		return fmt.Errorf("spec.runtime is required")
	}
	
	// Validate runtime
	validRuntimes := []string{"python", "nodejs", "go", "rust", "java"}
	if !contains(validRuntimes, spec.Spec.Runtime) {
		return fmt.Errorf("invalid runtime '%s'. Valid runtimes: %v", spec.Spec.Runtime, validRuntimes)
	}
	
	// Validate model configuration
	if spec.Spec.Model.Provider == "" {
		return fmt.Errorf("spec.model.provider is required")
	}
	
	if spec.Spec.Model.Name == "" {
		return fmt.Errorf("spec.model.name is required")
	}
	
	// Validate ports
	for i, port := range spec.Spec.Ports {
		if port.Container <= 0 || port.Container > 65535 {
			return fmt.Errorf("invalid container port %d at index %d", port.Container, i)
		}
		
		if port.Host != 0 && (port.Host <= 0 || port.Host > 65535) {
			return fmt.Errorf("invalid host port %d at index %d", port.Host, i)
		}
	}
	
	return nil
}

// FindAgentFile finds agent.yaml in the given directory
func (p *Parser) FindAgentFile(dir string) (string, error) {
	candidates := []string{"agent.yaml", "agent.yml", "Agent.yaml", "Agent.yml"}
	
	for _, candidate := range candidates {
		path := filepath.Join(dir, candidate)
		if fileExists(path) {
			return path, nil
		}
	}
	
	return "", fmt.Errorf("no agent.yaml file found in %s", dir)
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := ioutil.ReadFile(path)
	return err == nil
}