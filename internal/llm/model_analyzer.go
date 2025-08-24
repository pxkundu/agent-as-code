package llm

import (
	"fmt"
	"strings"
)

// ModelAnalyzer analyzes model capabilities and limitations
type ModelAnalyzer struct {
	modelManager *LocalLLMManager
}

// ModelAnalysis represents a comprehensive model analysis
type ModelAnalysis struct {
	ModelName        string
	Architecture     ModelArchitecture
	Performance      ModelPerformance
	Capabilities     []string
	Limitations      []string
	BestUseCases     []string
	OptimizationTips []string
}

// ModelArchitecture represents model architecture information
type ModelArchitecture struct {
	ModelType     string
	Parameters    string
	ContextWindow string
	TrainingData  string
}

// ModelPerformance represents model performance characteristics
type ModelPerformance struct {
	ResponseTime string
	MemoryUsage  string
	Throughput   string
}

// NewModelAnalyzer creates a new model analyzer
func NewModelAnalyzer() *ModelAnalyzer {
	return &ModelAnalyzer{
		modelManager: NewLocalLLMManager(),
	}
}

// IsModelAvailable checks if a model is available
func (a *ModelAnalyzer) IsModelAvailable(modelName string) bool {
	return a.modelManager.IsModelAvailable(modelName)
}

// AnalyzeModel performs comprehensive analysis of a model
func (a *ModelAnalyzer) AnalyzeModel(modelName string) (*ModelAnalysis, error) {
	// Get model info
	modelInfo, err := a.modelManager.GetModelInfo(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get model info: %v", err)
	}

	// Analyze the model
	analysis := &ModelAnalysis{
		ModelName:        modelName,
		Architecture:     a.analyzeArchitecture(modelName, modelInfo),
		Performance:      a.analyzePerformance(modelName),
		Capabilities:     a.analyzeCapabilities(modelName),
		Limitations:      a.analyzeLimitations(modelName),
		BestUseCases:     a.analyzeBestUseCases(modelName),
		OptimizationTips: a.generateOptimizationTips(modelName),
	}

	return analysis, nil
}

// analyzeArchitecture analyzes the model architecture
func (a *ModelAnalyzer) analyzeArchitecture(modelName string, modelInfo *LocalModel) ModelArchitecture {
	arch := ModelArchitecture{
		ModelType:     "Transformer",
		Parameters:    "Unknown",
		ContextWindow: "Unknown",
		TrainingData:  "Unknown",
	}

	// Determine model size from name
	if strings.Contains(modelName, "7b") {
		arch.Parameters = "7B parameters"
		arch.ContextWindow = "4K tokens"
	} else if strings.Contains(modelName, "13b") {
		arch.Parameters = "13B parameters"
		arch.ContextWindow = "8K tokens"
	} else if strings.Contains(modelName, "30b") {
		arch.Parameters = "30B parameters"
		arch.ContextWindow = "16K tokens"
	} else if strings.Contains(modelName, "65b") {
		arch.Parameters = "65B parameters"
		arch.ContextWindow = "32K tokens"
	} else if strings.Contains(modelName, "70b") {
		arch.Parameters = "70B parameters"
		arch.ContextWindow = "32K tokens"
	} else {
		arch.Parameters = "Unknown size"
		arch.ContextWindow = "Unknown"
	}

	// Determine model type
	if strings.Contains(modelName, "llama") {
		arch.ModelType = "LLaMA"
		arch.TrainingData = "Public datasets, code, conversations"
	} else if strings.Contains(modelName, "mistral") {
		arch.ModelType = "Mistral"
		arch.TrainingData = "High-quality web data, code, conversations"
	} else if strings.Contains(modelName, "codellama") {
		arch.ModelType = "Code Llama"
		arch.TrainingData = "Code repositories, documentation, conversations"
	} else if strings.Contains(modelName, "neural-chat") {
		arch.ModelType = "Neural Chat"
		arch.TrainingData = "Conversations, web data, books"
	} else if strings.Contains(modelName, "orca") {
		arch.ModelType = "Orca"
		arch.TrainingData = "Instruction-following data, conversations"
	}

	return arch
}

// analyzePerformance analyzes model performance characteristics
func (a *ModelAnalyzer) analyzePerformance(modelName string) ModelPerformance {
	perf := ModelPerformance{
		ResponseTime: "Unknown",
		MemoryUsage:  "Unknown",
		Throughput:   "Unknown",
	}

	// Estimate performance based on model size
	if strings.Contains(modelName, "7b") {
		perf.ResponseTime = "2-5 seconds"
		perf.MemoryUsage = "4-8 GB RAM"
		perf.Throughput = "10-20 requests/min"
	} else if strings.Contains(modelName, "13b") {
		perf.ResponseTime = "5-10 seconds"
		perf.MemoryUsage = "8-16 GB RAM"
		perf.Throughput = "5-10 requests/min"
	} else if strings.Contains(modelName, "30b") {
		perf.ResponseTime = "10-20 seconds"
		perf.MemoryUsage = "16-32 GB RAM"
		perf.Throughput = "2-5 requests/min"
	} else if strings.Contains(modelName, "65b") || strings.Contains(modelName, "70b") {
		perf.ResponseTime = "20-40 seconds"
		perf.MemoryUsage = "32-64 GB RAM"
		perf.Throughput = "1-3 requests/min"
	}

	return perf
}

// analyzeCapabilities analyzes model capabilities
func (a *ModelAnalyzer) analyzeCapabilities(modelName string) []string {
	var capabilities []string

	// Base capabilities for all models
	capabilities = append(capabilities, "Text generation")
	capabilities = append(capabilities, "Language understanding")
	capabilities = append(capabilities, "Context awareness")

	// Model-specific capabilities
	if strings.Contains(modelName, "llama") {
		capabilities = append(capabilities, "General conversation")
		capabilities = append(capabilities, "Creative writing")
		capabilities = append(capabilities, "Problem solving")
		capabilities = append(capabilities, "Multi-language support")
	}

	if strings.Contains(modelName, "mistral") {
		capabilities = append(capabilities, "High-quality text generation")
		capabilities = append(capabilities, "Efficient reasoning")
		capabilities = append(capabilities, "Code understanding")
		capabilities = append(capabilities, "Instruction following")
	}

	if strings.Contains(modelName, "codellama") {
		capabilities = append(capabilities, "Code generation")
		capabilities = append(capabilities, "Code debugging")
		capabilities = append(capabilities, "Code explanation")
		capabilities = append(capabilities, "Programming languages")
		capabilities = append(capabilities, "Software architecture")
	}

	if strings.Contains(modelName, "neural-chat") {
		capabilities = append(capabilities, "Conversational AI")
		capabilities = append(capabilities, "Emotional intelligence")
		capabilities = append(capabilities, "Personality consistency")
		capabilities = append(capabilities, "Multi-turn conversations")
	}

	if strings.Contains(modelName, "orca") {
		capabilities = append(capabilities, "Instruction following")
		capabilities = append(capabilities, "Task completion")
		capabilities = append(capabilities, "Reasoning")
		capabilities = append(capabilities, "Problem solving")
	}

	// Size-based capabilities
	if strings.Contains(modelName, "13b") || strings.Contains(modelName, "30b") ||
		strings.Contains(modelName, "65b") || strings.Contains(modelName, "70b") {
		capabilities = append(capabilities, "Complex reasoning")
		capabilities = append(capabilities, "Detailed analysis")
		capabilities = append(capabilities, "Long-form content")
		capabilities = append(capabilities, "Advanced problem solving")
	}

	return capabilities
}

// analyzeLimitations analyzes model limitations
func (a *ModelAnalyzer) analyzeLimitations(modelName string) []string {
	var limitations []string

	// Base limitations for all models
	limitations = append(limitations, "No real-time information")
	limitations = append(limitations, "Training data cutoff")
	limitations = append(limitations, "Potential hallucinations")
	limitations = append(limitations, "Context window limits")

	// Size-based limitations
	if strings.Contains(modelName, "7b") {
		limitations = append(limitations, "Limited reasoning complexity")
		limitations = append(limitations, "Shorter context retention")
		limitations = append(limitations, "Less nuanced understanding")
		limitations = append(limitations, "Faster but less accurate")
	}

	if strings.Contains(modelName, "13b") {
		limitations = append(limitations, "Moderate reasoning capability")
		limitations = append(limitations, "Balanced performance")
		limitations = append(limitations, "Memory constraints")
	}

	if strings.Contains(modelName, "30b") || strings.Contains(modelName, "65b") ||
		strings.Contains(modelName, "70b") {
		limitations = append(limitations, "High memory requirements")
		limitations = append(limitations, "Slower response times")
		limitations = append(limitations, "Resource intensive")
		limitations = append(limitations, "Higher computational cost")
	}

	// Model-specific limitations
	if strings.Contains(modelName, "codellama") {
		limitations = append(limitations, "Specialized for code")
		limitations = append(limitations, "Less general knowledge")
		limitations = append(limitations, "Code-specific context needed")
	}

	if strings.Contains(modelName, "neural-chat") {
		limitations = append(limitations, "Conversation-focused")
		limitations = append(limitations, "Less technical depth")
		limitations = append(limitations, "Personality consistency challenges")
	}

	return limitations
}

// analyzeBestUseCases analyzes best use cases for the model
func (a *ModelAnalyzer) analyzeBestUseCases(modelName string) []string {
	var useCases []string

	// Size-based use cases
	if strings.Contains(modelName, "7b") {
		useCases = append(useCases, "Fast prototyping")
		useCases = append(useCases, "Simple Q&A")
		useCases = append(useCases, "Basic text generation")
		useCases = append(useCases, "Resource-constrained environments")
		useCases = append(useCases, "Real-time applications")
	}

	if strings.Contains(modelName, "13b") {
		useCases = append(useCases, "Production applications")
		useCases = append(useCases, "Moderate complexity tasks")
		useCases = append(useCases, "Balanced performance needs")
		useCases = append(useCases, "General AI assistants")
		useCases = append(useCases, "Content creation")
	}

	if strings.Contains(modelName, "30b") || strings.Contains(modelName, "65b") ||
		strings.Contains(modelName, "70b") {
		useCases = append(useCases, "Complex reasoning tasks")
		useCases = append(useCases, "Research and analysis")
		useCases = append(useCases, "High-quality content generation")
		useCases = append(useCases, "Advanced AI applications")
		useCases = append(useCases, "Enterprise solutions")
	}

	// Model-specific use cases
	if strings.Contains(modelName, "llama") {
		useCases = append(useCases, "General AI applications")
		useCases = append(useCases, "Conversational AI")
		useCases = append(useCases, "Content generation")
		useCases = append(useCases, "Language tasks")
	}

	if strings.Contains(modelName, "mistral") {
		useCases = append(useCases, "High-quality text generation")
		useCases = append(useCases, "Reasoning tasks")
		useCases = append(useCases, "Instruction following")
		useCases = append(useCases, "Efficient AI applications")
	}

	if strings.Contains(modelName, "codellama") {
		useCases = append(useCases, "Code generation")
		useCases = append(useCases, "Software development")
		useCases = append(useCases, "Programming assistance")
		useCases = append(useCases, "Code review")
		useCases = append(useCases, "Technical documentation")
	}

	if strings.Contains(modelName, "neural-chat") {
		useCases = append(useCases, "Chatbots")
		useCases = append(useCases, "Customer service")
		useCases = append(useCases, "Conversational AI")
		useCases = append(useCases, "Social applications")
	}

	if strings.Contains(modelName, "orca") {
		useCases = append(useCases, "Task completion")
		useCases = append(useCases, "Instruction following")
		useCases = append(useCases, "Workflow automation")
		useCases = append(useCases, "Process optimization")
	}

	return useCases
}

// generateOptimizationTips generates optimization tips for the model
func (a *ModelAnalyzer) generateOptimizationTips(modelName string) []string {
	var tips []string

	// General optimization tips
	tips = append(tips, "Use appropriate temperature settings for your use case")
	tips = append(tips, "Optimize context window usage")
	tips = append(tips, "Implement proper error handling")
	tips = append(tips, "Monitor memory usage and performance")

	// Size-specific tips
	if strings.Contains(modelName, "7b") {
		tips = append(tips, "Keep prompts concise and focused")
		tips = append(tips, "Use streaming for real-time responses")
		tips = append(tips, "Implement caching for repeated queries")
		tips = append(tips, "Consider batch processing for efficiency")
	}

	if strings.Contains(modelName, "13b") {
		tips = append(tips, "Balance between speed and quality")
		tips = append(tips, "Use appropriate batch sizes")
		tips = append(tips, "Implement request queuing")
		tips = append(tips, "Monitor resource utilization")
	}

	if strings.Contains(modelName, "30b") || strings.Contains(modelName, "65b") ||
		strings.Contains(modelName, "70b") {
		tips = append(tips, "Implement proper resource management")
		tips = append(tips, "Use async processing for long operations")
		tips = append(tips, "Consider model sharding if possible")
		tips = append(tips, "Implement intelligent caching strategies")
		tips = append(tips, "Use load balancing for multiple instances")
	}

	// Model-specific tips
	if strings.Contains(modelName, "codellama") {
		tips = append(tips, "Provide clear code context")
		tips = append(tips, "Use code-specific prompts")
		tips = append(tips, "Implement code validation")
		tips = append(tips, "Consider security implications")
	}

	if strings.Contains(modelName, "neural-chat") {
		tips = append(tips, "Maintain conversation context")
		tips = append(tips, "Implement personality consistency")
		tips = append(tips, "Use conversation history effectively")
		tips = append(tips, "Handle emotional context appropriately")
	}

	return tips
}
