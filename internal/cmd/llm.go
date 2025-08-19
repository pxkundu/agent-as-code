package cmd

import (
	"fmt"
	"strings"

	"github.com/pxkundu/agent-as-code/internal/llm"
	"github.com/spf13/cobra"
)

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Manage local LLM models",
	Long: `Manage local LLM models for AI agent development.

This command provides tools to work with local LLM models, including
Ollama integration for running models locally without API costs.

Examples:
  agent llm list                    # List available local models
  agent llm pull llama2             # Pull a model from Ollama
  agent llm test llama2             # Test a local model
  agent llm recommend chatbot       # Get recommended models for chatbots`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var llmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available local models",
	Long: `List all available local LLM models.

This command shows all models that are currently available on your
local system through Ollama or other local backends.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listLocalModels()
	},
}

var llmPullCmd = &cobra.Command{
	Use:   "pull [MODEL]",
	Short: "Pull a model from Ollama",
	Long: `Pull a model from Ollama to your local system.

This command downloads and installs a model locally, making it available
for AI agent development without API costs.

Examples:
  agent llm pull llama2
  agent llm pull llama2:7b
  agent llm pull mistral:7b`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return pullLocalModel(modelName)
	},
}

var llmTestCmd = &cobra.Command{
	Use:   "test [MODEL]",
	Short: "Test a local model",
	Long: `Test a local model to ensure it's working correctly.

This command runs a simple test prompt through the specified model
to verify it's functioning properly.

Examples:
  agent llm test llama2
  agent llm test mistral:7b`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return testLocalModel(modelName)
	},
}

var llmRemoveCmd = &cobra.Command{
	Use:   "remove [MODEL]",
	Short: "Remove a local model",
	Long: `Remove a local model to free up disk space.

This command removes the specified model from your local system.

Examples:
  agent llm remove llama2
  agent llm remove mistral:7b`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return removeLocalModel(modelName)
	},
}

var llmRecommendCmd = &cobra.Command{
	Use:   "recommend [USE_CASE]",
	Short: "Get recommended models for specific use cases",
	Long: `Get recommended models for specific use cases.

This command suggests appropriate models based on your intended use case,
helping you choose the right model for your AI agent.

Use cases: chatbot, code, general, fast

Examples:
  agent llm recommend chatbot
  agent llm recommend code
  agent llm recommend fast`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		useCase := args[0]
		return recommendModels(useCase)
	},
}

var llmInfoCmd = &cobra.Command{
	Use:   "info [MODEL]",
	Short: "Show detailed information about a local model",
	Long: `Show detailed information about a local model.

This command displays comprehensive information about the specified
model, including size, modification date, and other details.

Examples:
  agent llm info llama2
  agent llm info mistral:7b`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return showModelInfo(modelName)
	},
}

var llmSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup local LLM environment",
	Long: `Setup local LLM environment for AI agent development.

This command helps you set up Ollama and other local LLM backends
for running AI agents locally without API costs.

Examples:
  agent llm setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupLocalLLM()
	},
}

func init() {
	// LLM command
	rootCmd.AddCommand(llmCmd)

	// LLM subcommands
	llmCmd.AddCommand(llmListCmd)
	llmCmd.AddCommand(llmPullCmd)
	llmCmd.AddCommand(llmTestCmd)
	llmCmd.AddCommand(llmRemoveCmd)
	llmCmd.AddCommand(llmRecommendCmd)
	llmCmd.AddCommand(llmInfoCmd)
	llmCmd.AddCommand(llmSetupCmd)
}

func listLocalModels() error {
	manager := llm.NewLocalLLMManager()

	// Check if Ollama is available
	if err := manager.CheckOllamaAvailability(); err != nil {
		fmt.Printf("⚠️  %v\n", err)
		fmt.Println("\n💡 To get started with local LLMs:")
		fmt.Println("   1. Install Ollama: https://ollama.ai")
		fmt.Println("   2. Start Ollama: ollama serve")
		fmt.Println("   3. Pull a model: agent llm pull llama2")
		return nil
	}

	models, err := manager.ListLocalModels()
	if err != nil {
		return fmt.Errorf("failed to list models: %v", err)
	}

	if len(models) == 0 {
		fmt.Println("ℹ️  No local models found")
		fmt.Println("\n💡 To get started:")
		fmt.Println("   agent llm pull llama2")
		return nil
	}

	fmt.Println("🤖 Available Local Models")
	fmt.Println("=========================")

	for _, model := range models {
		fmt.Printf("\n%s\n", model.Name)
		fmt.Printf("  Size:     %s\n", model.Size)
		fmt.Printf("  Backend:  %s\n", model.Backend)
		fmt.Printf("  Status:   %s\n", model.Status)
		if model.ModifiedAt != "" {
			fmt.Printf("  Modified: %s\n", model.ModifiedAt)
		}
	}

	return nil
}

func pullLocalModel(modelName string) error {
	manager := llm.NewLocalLLMManager()

	// Validate model name
	if err := manager.ValidateModelName(modelName); err != nil {
		return err
	}

	// Check if model is already available
	if manager.IsModelAvailable(modelName) {
		fmt.Printf("ℹ️  Model '%s' is already available\n", modelName)
		return nil
	}

	// Pull the model
	return manager.PullModel(modelName)
}

func testLocalModel(modelName string) error {
	manager := llm.NewLocalLLMManager()

	// Check if model is available
	if !manager.IsModelAvailable(modelName) {
		return fmt.Errorf("model '%s' is not available. Pull it first with 'agent llm pull %s'", modelName, modelName)
	}

	// Test the model
	return manager.TestModel(modelName)
}

func removeLocalModel(modelName string) error {
	manager := llm.NewLocalLLMManager()

	// Check if model is available
	if !manager.IsModelAvailable(modelName) {
		return fmt.Errorf("model '%s' is not available", modelName)
	}

	// Remove the model
	return manager.RemoveModel(modelName)
}

func recommendModels(useCase string) error {
	manager := llm.NewLocalLLMManager()

	recommendations := manager.GetRecommendedModels()

	models, ok := recommendations[strings.ToLower(useCase)]
	if !ok {
		validUseCases := make([]string, 0, len(recommendations))
		for uc := range recommendations {
			validUseCases = append(validUseCases, uc)
		}
		return fmt.Errorf("invalid use case '%s'. Valid use cases: %s", useCase, strings.Join(validUseCases, ", "))
	}

	fmt.Printf("🎯 Recommended Models for: %s\n", useCase)
	fmt.Println("=================================")

	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}

	fmt.Printf("\n💡 To pull a model: agent llm pull <model_name>\n")
	fmt.Printf("   Example: agent llm pull %s\n", models[0])

	return nil
}

func showModelInfo(modelName string) error {
	manager := llm.NewLocalLLMManager()

	// Check if model is available
	if !manager.IsModelAvailable(modelName) {
		return fmt.Errorf("model '%s' is not available. Pull it first with 'agent llm pull %s'", modelName, modelName)
	}

	info, err := manager.GetModelInfo(modelName)
	if err != nil {
		return fmt.Errorf("failed to get model info: %v", err)
	}

	fmt.Printf("📋 Model Information: %s\n", info.Name)
	fmt.Println("========================")
	fmt.Printf("Size:       %s\n", info.Size)
	fmt.Printf("Backend:    %s\n", info.Backend)
	fmt.Printf("Status:     %s\n", info.Status)
	fmt.Printf("Modified:   %s\n", info.ModifiedAt)
	fmt.Printf("Digest:     %s\n", info.Digest)

	if len(info.Details) > 0 {
		fmt.Println("\nDetails:")
		for key, value := range info.Details {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}

func setupLocalLLM() error {
	fmt.Println("🚀 Setting up Local LLM Environment")
	fmt.Println("===================================")

	fmt.Println("\n1️⃣  Installing Ollama...")
	fmt.Println("   Visit: https://ollama.ai")
	fmt.Println("   Or run: curl -fsSL https://ollama.ai/install.sh | sh")

	fmt.Println("\n2️⃣  Starting Ollama...")
	fmt.Println("   Run: ollama serve")

	fmt.Println("\n3️⃣  Pulling your first model...")
	fmt.Println("   Run: agent llm pull llama2")

	fmt.Println("\n4️⃣  Testing the setup...")
	fmt.Println("   Run: agent llm test llama2")

	fmt.Println("\n5️⃣  Creating your first local AI agent...")
	fmt.Println("   Run: agent init my-chatbot --template chatbot --model local/llama2")

	fmt.Println("\n✅ You're all set for local AI development!")
	fmt.Println("\n💡 Benefits of local LLMs:")
	fmt.Println("   • No API costs")
	fmt.Println("   • Complete privacy")
	fmt.Println("   • No rate limits")
	fmt.Println("   • Works offline")

	return nil
}
