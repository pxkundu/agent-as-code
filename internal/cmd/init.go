package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pxkundu/agent-as-code/internal/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [NAME]",
	Short: "Initialize a new agent project",
	Long: `Initialize a new agent project with the specified name and template.

This command creates a new directory with the agent name and sets up
the basic project structure including agent.yaml configuration file
and template-specific implementation files.

Examples:
  agent init my-chatbot --template chatbot
  agent init sentiment-analyzer --template sentiment
  agent init my-agent --runtime python`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

var (
	initTemplate string
	initRuntime  string
	initModel    string
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "template to use (chatbot, sentiment, summarizer, translator, data-analyzer, content-gen)")
	initCmd.Flags().StringVarP(&initRuntime, "runtime", "r", "python", "runtime environment (python, nodejs, go)")
	initCmd.Flags().StringVarP(&initModel, "model", "m", "openai/gpt-4", "default model to use (supports local models like 'local/llama2')")
}

func runInit(cmd *cobra.Command, args []string) error {
	agentName := args[0]

	// Validate agent name
	if agentName == "" {
		return fmt.Errorf("agent name cannot be empty")
	}

	// Check if directory already exists
	if _, err := os.Stat(agentName); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", agentName)
	}

	// Create agent directory
	if err := os.MkdirAll(agentName, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Initialize template manager
	templateManager := templates.New()

	// Template validation is now handled by the template manager with fallback logic

	// Validate local model if specified
	if strings.HasPrefix(initModel, "local/") {
		localModelName := strings.TrimPrefix(initModel, "local/")
		if err := validateLocalModel(localModelName); err != nil {
			return fmt.Errorf("local model validation failed: %v", err)
		}
	}

	// Determine template to use
	template := initTemplate
	if template == "" {
		template = "basic" // Default template
	}

	// Create agent configuration
	config := &templates.AgentConfig{
		Name:     agentName,
		Template: template,
		Runtime:  initRuntime,
		Model:    initModel,
	}

	// Generate project files
	if err := templateManager.Generate(agentName, config); err != nil {
		// Clean up on error
		os.RemoveAll(agentName)
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Success message
	fmt.Printf("✅ Agent project '%s' created successfully!\n\n", agentName)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", agentName)
	fmt.Printf("  agent build -t %s:latest .\n", agentName)
	fmt.Printf("  agent run %s:latest\n", agentName)

	if template != "basic" {
		fmt.Printf("\n📖 Check the README.md for template-specific instructions.\n")
	}

	return nil
}

func validateTemplate(template string) error {
	validTemplates := []string{"basic", "chatbot", "sentiment", "summarizer", "translator", "data-analyzer", "content-gen"}

	for _, valid := range validTemplates {
		if template == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid template '%s'. Valid templates: %v", template, validTemplates)
}

func validateLocalModel(modelName string) error {
	// Import the local LLM manager to validate local models
	// For now, we'll do basic validation
	if modelName == "" {
		return fmt.Errorf("local model name cannot be empty")
	}

	if strings.Contains(modelName, " ") {
		return fmt.Errorf("local model name cannot contain spaces")
	}

	// Check if it's a valid Ollama model format
	parts := strings.Split(modelName, ":")
	if len(parts) > 2 {
		return fmt.Errorf("invalid local model name format. Use 'model' or 'model:tag'")
	}

	return nil
}

func isValidTemplate(template string) bool {
	validTemplates := getValidTemplates()
	for _, valid := range validTemplates {
		if template == valid {
			return true
		}
	}
	return false
}

func getValidTemplates() []string {
	return []string{"basic", "chatbot", "sentiment", "summarizer", "translator", "data-analyzer", "content-gen"}
}
