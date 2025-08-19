package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [TAG]",
	Short: "Show agent details",
	Long: `Show detailed information about an agent image.

This command displays comprehensive information about the specified agent,
including configuration details, runtime settings, capabilities, and metadata.

Examples:
  agent inspect my-agent:latest
  agent inspect my-agent:v1.0.0
  agent inspect --format json my-agent:latest`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]
		format, _ := cmd.Flags().GetString("format")
		
		fmt.Printf("🔍 Inspecting agent: %s\n", tag)
		
		// Check if the agent image exists
		if !imageExists(tag) {
			return fmt.Errorf("agent image '%s' not found", tag)
		}
		
		// Get agent information
		info, err := getAgentInfo(tag)
		if err != nil {
			return fmt.Errorf("failed to inspect agent: %v", err)
		}
		
		// Display the information
		return displayAgentInfo(info, format)
	},
}

func init() {
	inspectCmd.Flags().String("format", "table", "output format (table, json)")
	rootCmd.AddCommand(inspectCmd)
}

type AgentInfo struct {
	Tag         string            `json:"tag"`
	ImageID     string            `json:"image_id"`
	Created     string            `json:"created"`
	Size        string            `json:"size"`
	Config      AgentConfig       `json:"config"`
	Runtime     RuntimeInfo       `json:"runtime"`
	Health      HealthInfo        `json:"health"`
	Ports       []PortMapping     `json:"ports"`
	Environment []EnvVariable     `json:"environment"`
	Labels      map[string]string `json:"labels"`
}

type AgentConfig struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Capabilities []string `json:"capabilities"`
	Model       ModelInfo `json:"model"`
}

type ModelInfo struct {
	Provider string            `json:"provider"`
	Name     string            `json:"name"`
	Config   map[string]string `json:"config"`
}

type RuntimeInfo struct {
	Type      string `json:"type"`
	BaseImage string `json:"base_image"`
	WorkDir   string `json:"work_dir"`
}

type HealthInfo struct {
	Command     []string `json:"command"`
	Interval    string   `json:"interval"`
	Timeout     string   `json:"timeout"`
	Retries    int      `json:"retries"`
	StartPeriod string  `json:"start_period"`
}

type PortMapping struct {
	Host      string `json:"host"`
	Container string `json:"container"`
	Protocol  string `json:"protocol"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	From  string `json:"from,omitempty"`
}

func getAgentInfo(tag string) (*AgentInfo, error) {
	// In a real implementation, this would query Docker and parse the agent.yaml
	// For now, we'll return mock data based on the tag
	
	info := &AgentInfo{
		Tag:     tag,
		ImageID: "sha256:1234567890abcdef",
		Created: "2025-08-16T10:30:00Z",
		Size:    "45.2MB",
		Config: AgentConfig{
			Name:        strings.Split(tag, ":")[0],
			Version:     "1.0.0",
			Description: fmt.Sprintf("%s agent", strings.Split(tag, ":")[0]),
			Capabilities: []string{"conversation", "api"},
			Model: ModelInfo{
				Provider: "openai",
				Name:     "gpt-4",
				Config: map[string]string{
					"temperature": "0.7",
					"max_tokens":  "500",
				},
			},
		},
		Runtime: RuntimeInfo{
			Type:      "python",
			BaseImage: "python:3.11-slim",
			WorkDir:   "/app",
		},
		Health: HealthInfo{
			Command:     []string{"curl", "-f", "http://localhost:8080/health"},
			Interval:    "30s",
			Timeout:     "10s",
			Retries:    3,
			StartPeriod: "5s",
		},
		Ports: []PortMapping{
			{
				Host:      "8080",
				Container: "8080",
				Protocol:  "tcp",
			},
		},
		Environment: []EnvVariable{
			{
				Name:  "LOG_LEVEL",
				Value: "INFO",
			},
		},
		Labels: map[string]string{
			"maintainer": "Agent as Code Team",
			"version":    "1.0.0",
		},
	}
	
	return info, nil
}

func displayAgentInfo(info *AgentInfo, format string) error {
	switch format {
	case "json":
		return displayJSON(info)
	default:
		return displayTable(info)
	}
}

func displayJSON(info *AgentInfo) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	
	fmt.Println(string(data))
	return nil
}

func displayTable(info *AgentInfo) error {
	fmt.Printf("\n📋 Agent Information\n")
	fmt.Printf("===================\n")
	fmt.Printf("Tag:         %s\n", info.Tag)
	fmt.Printf("Image ID:    %s\n", info.ImageID)
	fmt.Printf("Created:     %s\n", info.Created)
	fmt.Printf("Size:        %s\n", info.Size)
	
	fmt.Printf("\n🔧 Configuration\n")
	fmt.Printf("================\n")
	fmt.Printf("Name:        %s\n", info.Config.Name)
	fmt.Printf("Version:     %s\n", info.Config.Version)
	fmt.Printf("Description: %s\n", info.Config.Description)
	fmt.Printf("Capabilities: %s\n", strings.Join(info.Config.Capabilities, ", "))
	
	fmt.Printf("\n🤖 Model\n")
	fmt.Printf("========\n")
	fmt.Printf("Provider:    %s\n", info.Config.Model.Provider)
	fmt.Printf("Name:        %s\n", info.Config.Model.Name)
	fmt.Printf("Config:      %v\n", info.Config.Model.Config)
	
	fmt.Printf("\n⚙️  Runtime\n")
	fmt.Printf("==========\n")
	fmt.Printf("Type:        %s\n", info.Runtime.Type)
	fmt.Printf("Base Image:  %s\n", info.Runtime.BaseImage)
	fmt.Printf("Work Dir:    %s\n", info.Runtime.WorkDir)
	
	fmt.Printf("\n🏥 Health Check\n")
	fmt.Printf("===============\n")
	fmt.Printf("Command:     %s\n", strings.Join(info.Health.Command, " "))
	fmt.Printf("Interval:    %s\n", info.Health.Interval)
	fmt.Printf("Timeout:     %s\n", info.Health.Timeout)
	fmt.Printf("Retries:     %d\n", info.Health.Retries)
	fmt.Printf("Start Period: %s\n", info.Health.StartPeriod)
	
	fmt.Printf("\n🌐 Ports\n")
	fmt.Printf("========\n")
	for _, port := range info.Ports {
		fmt.Printf("  %s:%s (%s)\n", port.Host, port.Container, port.Protocol)
	}
	
	fmt.Printf("\n🔑 Environment\n")
	fmt.Printf("==============\n")
	for _, env := range info.Environment {
		if env.From != "" {
			fmt.Printf("  %s (from: %s)\n", env.Name, env.From)
		} else {
			fmt.Printf("  %s=%s\n", env.Name, env.Value)
		}
	}
	
	fmt.Printf("\n🏷️  Labels\n")
	fmt.Printf("==========\n")
	for key, value := range info.Labels {
		fmt.Printf("  %s: %s\n", key, value)
	}
	
	return nil
}
