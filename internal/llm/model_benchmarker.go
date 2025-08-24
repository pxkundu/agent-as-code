package llm

import (
	"fmt"
	"time"
)

// ModelBenchmarker runs comprehensive benchmarks on models
type ModelBenchmarker struct {
	modelManager *LocalLLMManager
}

// BenchmarkResult represents the result of a model benchmark
type BenchmarkResult struct {
	ModelName           string
	AverageResponseTime string
	MemoryUsage         string
	Throughput          string
	QualityScore        string
	CostEfficiency      string
	Tasks               []TaskResult
}

// TaskResult represents the result of a specific benchmark task
type TaskResult struct {
	TaskName     string
	ResponseTime time.Duration
	Accuracy     float64
	MemoryUsed   string
	Success      bool
	Error        string
}

// BenchmarkTask represents a benchmark task
type BenchmarkTask struct {
	Name        string
	Prompt      string
	Expected    string
	MaxTokens   int
	Temperature float64
}

// NewModelBenchmarker creates a new model benchmarker
func NewModelBenchmarker() *ModelBenchmarker {
	return &ModelBenchmarker{
		modelManager: NewLocalLLMManager(),
	}
}

// GetAvailableModels gets all available models for benchmarking
func (b *ModelBenchmarker) GetAvailableModels() ([]string, error) {
	models, err := b.modelManager.ListLocalModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %v", err)
	}

	var modelNames []string
	for _, model := range models {
		modelNames = append(modelNames, model.Name)
	}

	return modelNames, nil
}

// RunBenchmarks runs comprehensive benchmarks on all models
func (b *ModelBenchmarker) RunBenchmarks(modelNames []string) ([]*BenchmarkResult, error) {
	var results []*BenchmarkResult

	for _, modelName := range modelNames {
		fmt.Printf("🏃 Benchmarking %s...\n", modelName)

		result, err := b.benchmarkModel(modelName)
		if err != nil {
			fmt.Printf("⚠️  Failed to benchmark %s: %v\n", modelName, err)
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// benchmarkModel benchmarks a single model
func (b *ModelBenchmarker) benchmarkModel(modelName string) (*BenchmarkResult, error) {
	// Define benchmark tasks
	tasks := b.getBenchmarkTasks()

	var taskResults []TaskResult
	var totalResponseTime time.Duration
	var totalMemory int64
	var successfulTasks int

	// Run each task
	for _, task := range tasks {
		result, err := b.runTask(modelName, task)
		if err != nil {
			result.Error = err.Error()
			result.Success = false
		} else {
			result.Success = true
			successfulTasks++
			totalResponseTime += result.ResponseTime
		}

		taskResults = append(taskResults, result)
	}

	// Calculate metrics
	avgResponseTime := "N/A"
	if successfulTasks > 0 {
		avgResponseTime = fmt.Sprintf("%.2fs", totalResponseTime.Seconds()/float64(successfulTasks))
	}

	memoryUsage := "N/A"
	if totalMemory > 0 {
		memoryUsage = b.formatBytes(totalMemory)
	}

	throughput := "N/A"
	if successfulTasks > 0 {
		throughput = fmt.Sprintf("%.1f tasks/min", float64(successfulTasks)/totalResponseTime.Minutes())
	}

	qualityScore := b.calculateQualityScore(taskResults)
	costEfficiency := b.calculateCostEfficiency(modelName, qualityScore, avgResponseTime)

	return &BenchmarkResult{
		ModelName:           modelName,
		AverageResponseTime: avgResponseTime,
		MemoryUsage:         memoryUsage,
		Throughput:          throughput,
		QualityScore:        qualityScore,
		CostEfficiency:      costEfficiency,
		Tasks:               taskResults,
	}, nil
}

// getBenchmarkTasks returns the benchmark tasks to run
func (b *ModelBenchmarker) getBenchmarkTasks() []BenchmarkTask {
	return []BenchmarkTask{
		{
			Name:        "Simple Question",
			Prompt:      "What is the capital of France?",
			Expected:    "Paris",
			MaxTokens:   50,
			Temperature: 0.7,
		},
		{
			Name:        "Code Generation",
			Prompt:      "Write a Python function to calculate fibonacci numbers",
			Expected:    "def fibonacci",
			MaxTokens:   200,
			Temperature: 0.3,
		},
		{
			Name:        "Sentiment Analysis",
			Prompt:      "Analyze the sentiment of: 'I love this product, it's amazing!'",
			Expected:    "positive",
			MaxTokens:   100,
			Temperature: 0.2,
		},
		{
			Name:        "Translation",
			Prompt:      "Translate 'Hello, how are you?' to Spanish",
			Expected:    "Hola",
			MaxTokens:   50,
			Temperature: 0.4,
		},
		{
			Name:        "Creative Writing",
			Prompt:      "Write a short story about a robot learning to paint",
			Expected:    "story",
			MaxTokens:   300,
			Temperature: 0.8,
		},
	}
}

// runTask runs a single benchmark task
func (b *ModelBenchmarker) runTask(modelName string, task BenchmarkTask) (TaskResult, error) {
	start := time.Now()

	// Simulate running the task (in a real implementation, this would call the actual model)
	time.Sleep(100 * time.Millisecond) // Simulate processing time

	responseTime := time.Since(start)

	// Simulate results (in a real implementation, this would be actual model output)
	accuracy := 0.85 + (0.1 * float64(time.Now().UnixNano()%100) / 100) // Random accuracy between 0.85-0.95
	memoryUsed := "128MB"                                               // Simulated memory usage

	return TaskResult{
		TaskName:     task.Name,
		ResponseTime: responseTime,
		Accuracy:     accuracy,
		MemoryUsed:   memoryUsed,
		Success:      true,
	}, nil
}

// calculateQualityScore calculates the overall quality score
func (b *ModelBenchmarker) calculateQualityScore(taskResults []TaskResult) string {
	if len(taskResults) == 0 {
		return "N/A"
	}

	var totalAccuracy float64
	var successfulTasks int

	for _, task := range taskResults {
		if task.Success {
			totalAccuracy += task.Accuracy
			successfulTasks++
		}
	}

	if successfulTasks == 0 {
		return "0%"
	}

	avgAccuracy := totalAccuracy / float64(successfulTasks)
	return fmt.Sprintf("%.1f%%", avgAccuracy*100)
}

// calculateCostEfficiency calculates the cost efficiency score
func (b *ModelBenchmarker) calculateCostEfficiency(modelName, qualityScore, responseTime string) string {
	// Simple heuristic based on model name and performance
	if containsSubstring(modelName, "7b") {
		return "High"
	} else if containsSubstring(modelName, "13b") {
		return "Medium"
	} else if containsSubstring(modelName, "30b") || containsSubstring(modelName, "65b") || containsSubstring(modelName, "70b") {
		return "Low"
	}
	return "Medium"
}

// contains checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr)))
}

// formatBytes formats bytes into human-readable format
func (b *ModelBenchmarker) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GenerateRecommendations generates recommendations based on benchmark results
func (b *ModelBenchmarker) GenerateRecommendations(results []*BenchmarkResult) []string {
	var recommendations []string

	if len(results) == 0 {
		return []string{"No models available for benchmarking"}
	}

	// Find best performing models
	var bestResponseTime, bestQuality, bestEfficiency *BenchmarkResult

	for _, result := range results {
		if bestResponseTime == nil || result.AverageResponseTime < bestResponseTime.AverageResponseTime {
			bestResponseTime = result
		}
		if bestQuality == nil || result.QualityScore > bestQuality.QualityScore {
			bestQuality = result
		}
		if bestEfficiency == nil || result.CostEfficiency == "High" {
			bestEfficiency = result
		}
	}

	// Generate recommendations
	if bestResponseTime != nil {
		recommendations = append(recommendations,
			fmt.Sprintf("Fastest model: %s (%s response time)",
				bestResponseTime.ModelName, bestResponseTime.AverageResponseTime))
	}

	if bestQuality != nil {
		recommendations = append(recommendations,
			fmt.Sprintf("Highest quality: %s (%s quality score)",
				bestQuality.ModelName, bestQuality.QualityScore))
	}

	if bestEfficiency != nil {
		recommendations = append(recommendations,
			fmt.Sprintf("Most cost-effective: %s (%s efficiency)",
				bestEfficiency.ModelName, bestEfficiency.CostEfficiency))
	}

	// General recommendations
	recommendations = append(recommendations,
		"Consider model size vs. performance trade-offs for your use case")
	recommendations = append(recommendations,
		"Test models with your specific prompts and requirements")
	recommendations = append(recommendations,
		"Monitor memory usage and response times in production")

	return recommendations
}
