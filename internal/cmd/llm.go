package cmd

import (
	"fmt"
	"strings"

	"github.com/pxkundu/agent-as-code/internal/llm"
	"github.com/spf13/cobra"
)

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Manage local LLM models and create intelligent agents",
	Long: `Manage local LLM models and create intelligent, fully functional AI agents.

This command provides advanced tools to work with local LLM models, including
Ollama integration, intelligent agent generation, automated testing, and
optimization for specific use cases.

Examples:
  agent llm list                    # List available local models
  agent llm pull llama2             # Pull a model from Ollama
  agent llm test llama2             # Test a local model
  agent llm recommend chatbot       # Get recommended models for chatbots
  agent llm create-agent chatbot    # Create intelligent chatbot agent
  agent llm optimize llama2         # Optimize model for specific use case
  agent llm benchmark               # Benchmark all local models
  agent llm deploy-agent my-agent   # Deploy and test agent locally`,
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

var llmCreateAgentCmd = &cobra.Command{
	Use:   "create-agent [USE_CASE]",
	Short: "Create an intelligent, fully functional agent",
	Long: `Create an intelligent, fully functional AI agent optimized for a specific use case.

This command uses LLM intelligence to:
- Generate optimized code based on the use case
- Create comprehensive test suites
- Set up proper error handling and logging
- Configure optimal model parameters
- Generate deployment configurations
- Create detailed documentation

Use cases: chatbot, sentiment-analyzer, code-assistant, data-analyzer, 
          content-generator, translator, qa-system, workflow-automation

Examples:
  agent llm create-agent chatbot
  agent llm create-agent sentiment-analyzer --model local/llama2
  agent llm create-agent code-assistant --optimize --test`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		useCase := args[0]
		return createIntelligentAgent(useCase)
	},
}

var llmOptimizeCmd = &cobra.Command{
	Use:   "optimize [MODEL] [USE_CASE]",
	Short: "Optimize a model for specific use case",
	Long: `Optimize a local LLM model for a specific use case.

This command analyzes the model and use case to:
- Adjust model parameters (temperature, top_p, etc.)
- Create custom prompts and system messages
- Optimize context window usage
- Generate performance benchmarks
- Create use case specific configurations

Examples:
  agent llm optimize llama2 chatbot
  agent llm optimize mistral:7b code-generation
  agent llm optimize codellama:13b debugging`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		useCase := args[1]
		return optimizeModelForUseCase(modelName, useCase)
	},
}

var llmBenchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark all local models",
	Long: `Run comprehensive benchmarks on all local LLM models.

This command tests models across multiple dimensions:
- Response time and throughput
- Memory usage and efficiency
- Quality assessment for different tasks
- Cost-benefit analysis
- Performance recommendations

Examples:
  agent llm benchmark
  agent llm benchmark --tasks chatbot,code,analysis
  agent llm benchmark --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return benchmarkAllModels()
	},
}

var llmDeployAgentCmd = &cobra.Command{
	Use:   "deploy-agent [AGENT_NAME]",
	Short: "Deploy and test an agent locally",
	Long: `Deploy and comprehensively test an agent on your local machine.

This command:
- Builds the agent container
- Deploys it locally
- Runs automated tests
- Validates functionality
- Provides performance metrics
- Generates deployment report

Examples:
  agent llm deploy-agent my-chatbot
  agent llm deploy-agent sentiment-analyzer --test-suite comprehensive
  agent llm deploy-agent code-assistant --monitor`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return deployAndTestAgent(agentName)
	},
}

var llmAnalyzeCmd = &cobra.Command{
	Use:   "analyze [MODEL]",
	Short: "Analyze model capabilities and limitations",
	Long: `Analyze a local LLM model's capabilities and limitations.

This command provides deep insights into:
- Model architecture and parameters
- Performance characteristics
- Best use cases and limitations
- Optimization opportunities
- Integration recommendations

Examples:
  agent llm analyze llama2
  agent llm analyze mistral:7b --detailed
  agent llm analyze codellama:13b --capabilities`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return analyzeModelCapabilities(modelName)
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

	// New intelligent commands
	llmCmd.AddCommand(llmCreateAgentCmd)
	llmCmd.AddCommand(llmOptimizeCmd)
	llmCmd.AddCommand(llmBenchmarkCmd)
	llmCmd.AddCommand(llmDeployAgentCmd)
	llmCmd.AddCommand(llmAnalyzeCmd)
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

func createIntelligentAgent(useCase string) error {
	fmt.Printf("🧠 Creating intelligent agent for: %s\n", useCase)
	fmt.Println("=====================================")

	// Initialize intelligent agent creator
	creator := llm.NewIntelligentAgentCreator()

	// Validate use case
	if err := creator.ValidateUseCase(useCase); err != nil {
		return fmt.Errorf("invalid use case: %v", err)
	}

	// Get recommended model for the use case
	recommendedModel, err := creator.GetRecommendedModel(useCase)
	if err != nil {
		return fmt.Errorf("failed to get recommended model: %v", err)
	}

	fmt.Printf("📋 Use Case: %s\n", useCase)
	fmt.Printf("🤖 Recommended Model: %s\n", recommendedModel)
	fmt.Printf("🔧 Capabilities: %s\n", strings.Join(creator.GetCapabilities(useCase), ", "))

	// Create intelligent agent
	agentConfig, err := creator.CreateAgent(useCase, recommendedModel)
	if err != nil {
		return fmt.Errorf("failed to create agent: %v", err)
	}

	fmt.Printf("\n✅ Intelligent agent created successfully!\n")
	fmt.Printf("📁 Project Directory: %s\n", agentConfig.Name)
	fmt.Printf("🐍 Runtime: %s\n", agentConfig.Runtime)
	fmt.Printf("🧠 Model: %s\n", agentConfig.Model)
	fmt.Printf("📚 Dependencies: %d packages\n", len(agentConfig.Dependencies))
	fmt.Printf("🧪 Test Coverage: %s\n", agentConfig.TestCoverage)

	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", agentConfig.Name)
	fmt.Printf("   agent build -t %s:latest .\n", agentConfig.Name)
	fmt.Printf("   agent llm deploy-agent %s\n", agentConfig.Name)

	return nil
}

func optimizeModelForUseCase(modelName, useCase string) error {
	fmt.Printf("⚡ Optimizing %s for %s\n", modelName, useCase)
	fmt.Println("=================================")

	// Initialize model optimizer
	optimizer := llm.NewModelOptimizer()

	// Check if model is available
	if !optimizer.IsModelAvailable(modelName) {
		return fmt.Errorf("model '%s' is not available. Pull it first with 'agent llm pull %s'", modelName, modelName)
	}

	// Optimize model for use case
	optimization, err := optimizer.OptimizeForUseCase(modelName, useCase)
	if err != nil {
		return fmt.Errorf("optimization failed: %v", err)
	}

	fmt.Printf("✅ Model optimization completed!\n\n")
	fmt.Printf("📊 Performance Improvements:\n")
	fmt.Printf("  Response Time: %s\n", optimization.ResponseTimeImprovement)
	fmt.Printf("  Memory Usage: %s\n", optimization.MemoryOptimization)
	fmt.Printf("  Quality Score: %s\n", optimization.QualityImprovement)

	fmt.Printf("\n🔧 Optimized Parameters:\n")
	for param, value := range optimization.Parameters {
		fmt.Printf("  %s: %v\n", param, value)
	}

	fmt.Printf("\n📝 System Message:\n")
	fmt.Printf("  %s\n", optimization.SystemMessage)

	fmt.Printf("\n💾 Configuration saved to: %s\n", optimization.ConfigPath)

	return nil
}

func benchmarkAllModels() error {
	fmt.Println("🏁 Running comprehensive model benchmarks")
	fmt.Println("=======================================")

	// Initialize benchmark runner
	benchmarker := llm.NewModelBenchmarker()

	// Get all available models
	models, err := benchmarker.GetAvailableModels()
	if err != nil {
		return fmt.Errorf("failed to get models: %v", err)
	}

	if len(models) == 0 {
		fmt.Println("ℹ️  No models available for benchmarking")
		fmt.Println("💡 Pull some models first:")
		fmt.Println("   agent llm pull llama2")
		fmt.Println("   agent llm pull mistral:7b")
		return nil
	}

	// Run benchmarks
	results, err := benchmarker.RunBenchmarks(models)
	if err != nil {
		return fmt.Errorf("benchmarking failed: %v", err)
	}

	// Display results
	fmt.Printf("\n📊 Benchmark Results\n")
	fmt.Println("===================")

	for _, result := range results {
		fmt.Printf("\n🤖 %s\n", result.ModelName)
		fmt.Printf("  ⏱️  Response Time: %s\n", result.AverageResponseTime)
		fmt.Printf("  🧠 Memory Usage: %s\n", result.MemoryUsage)
		fmt.Printf("  📈 Throughput: %s\n", result.Throughput)
		fmt.Printf("  🎯 Quality Score: %s\n", result.QualityScore)
		fmt.Printf("  💰 Cost Efficiency: %s\n", result.CostEfficiency)
	}

	// Generate recommendations
	recommendations := benchmarker.GenerateRecommendations(results)
	fmt.Printf("\n💡 Recommendations:\n")
	for _, rec := range recommendations {
		fmt.Printf("  • %s\n", rec)
	}

	return nil
}

func deployAndTestAgent(agentName string) error {
	fmt.Printf("🚀 Deploying and testing agent: %s\n", agentName)
	fmt.Println("=====================================")

	// Initialize deployment manager
	deployer := llm.NewAgentDeployer()

	// Check if agent project exists
	if !deployer.AgentExists(agentName) {
		return fmt.Errorf("agent project '%s' not found. Create it first with 'agent init %s'", agentName, agentName)
	}

	// Build agent
	fmt.Printf("🔨 Building agent...\n")
	if err := deployer.BuildAgent(agentName); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	// Deploy agent
	fmt.Printf("📦 Deploying agent...\n")
	container, err := deployer.DeployAgent(agentName)
	if err != nil {
		return fmt.Errorf("deployment failed: %v", err)
	}

	// Run tests
	fmt.Printf("🧪 Running tests...\n")
	testResults, err := deployer.RunTests(agentName)
	if err != nil {
		return fmt.Errorf("testing failed: %v", err)
	}

	// Validate functionality
	fmt.Printf("✅ Validating functionality...\n")
	validation, err := deployer.ValidateAgent(agentName)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Display results
	fmt.Printf("\n🎉 Agent deployment successful!\n")
	fmt.Printf("🐳 Container: %s\n", container.Name)
	fmt.Printf("🔗 Access: http://localhost:%s\n", container.Port)
	fmt.Printf("🧪 Tests: %d/%d passed\n", testResults.Passed, testResults.Total)
	fmt.Printf("✅ Validation: %s\n", validation.Status)

	if validation.Issues > 0 {
		fmt.Printf("⚠️  Issues found: %d\n", validation.Issues)
		for _, issue := range validation.IssueDetails {
			fmt.Printf("   • %s\n", issue)
		}
	}

	fmt.Printf("\n📊 Performance Metrics:\n")
	fmt.Printf("  Response Time: %s\n", validation.ResponseTime)
	fmt.Printf("  Memory Usage: %s\n", validation.MemoryUsage)
	fmt.Printf("  CPU Usage: %s\n", validation.CPUUsage)

	fmt.Printf("\n💡 Management commands:\n")
	fmt.Printf("  View logs: agent logs %s\n", container.Name)
	fmt.Printf("  Stop agent: agent stop %s\n", container.Name)
	fmt.Printf("  Restart: agent restart %s\n", container.Name)

	return nil
}

func analyzeModelCapabilities(modelName string) error {
	fmt.Printf("🔍 Analyzing model: %s\n", modelName)
	fmt.Println("=========================")

	// Initialize model analyzer
	analyzer := llm.NewModelAnalyzer()

	// Check if model is available
	if !analyzer.IsModelAvailable(modelName) {
		return fmt.Errorf("model '%s' is not available. Pull it first with 'agent llm pull %s'", modelName, modelName)
	}

	// Analyze model
	analysis, err := analyzer.AnalyzeModel(modelName)
	if err != nil {
		return fmt.Errorf("analysis failed: %v", err)
	}

	// Display analysis results
	fmt.Printf("✅ Model analysis completed!\n\n")

	fmt.Printf("🏗️  Architecture:\n")
	fmt.Printf("  Model Type: %s\n", analysis.Architecture.ModelType)
	fmt.Printf("  Parameters: %s\n", analysis.Architecture.Parameters)
	fmt.Printf("  Context Window: %s\n", analysis.Architecture.ContextWindow)
	fmt.Printf("  Training Data: %s\n", analysis.Architecture.TrainingData)

	fmt.Printf("\n📊 Performance:\n")
	fmt.Printf("  Response Time: %s\n", analysis.Performance.ResponseTime)
	fmt.Printf("  Memory Usage: %s\n", analysis.Performance.MemoryUsage)
	fmt.Printf("  Throughput: %s\n", analysis.Performance.Throughput)

	fmt.Printf("\n🎯 Capabilities:\n")
	for _, capability := range analysis.Capabilities {
		fmt.Printf("  ✅ %s\n", capability)
	}

	fmt.Printf("\n⚠️  Limitations:\n")
	for _, limitation := range analysis.Limitations {
		fmt.Printf("  ❌ %s\n", limitation)
	}

	fmt.Printf("\n💡 Best Use Cases:\n")
	for _, useCase := range analysis.BestUseCases {
		fmt.Printf("  🎯 %s\n", useCase)
	}

	fmt.Printf("\n🔧 Optimization Tips:\n")
	for _, tip := range analysis.OptimizationTips {
		fmt.Printf("  💡 %s\n", tip)
	}

	return nil
}
