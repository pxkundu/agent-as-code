package llm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// LocalLLMManager handles local LLM operations
type LocalLLMManager struct {
	ollamaURL string
	timeout   time.Duration
}

// LocalModel represents a local LLM model
type LocalModel struct {
	Name        string            `json:"name"`
	Size        string            `json:"size"`
	ModifiedAt  string            `json:"modified_at"`
	Digest      string            `json:"digest"`
	Details     map[string]string `json:"details,omitempty"`
	Backend     string            `json:"backend"`
	Status      string            `json:"status"`
}

// LocalModelResponse represents Ollama API response
type LocalModelResponse struct {
	Models []LocalModel `json:"models"`
}

// NewLocalLLMManager creates a new local LLM manager
func NewLocalLLMManager() *LocalLLMManager {
	return &LocalLLMManager{
		ollamaURL: "http://localhost:11434",
		timeout:   30 * time.Second,
	}
}

// CheckOllamaAvailability checks if Ollama is running
func (m *LocalLLMManager) CheckOllamaAvailability() error {
	client := &http.Client{Timeout: m.timeout}
	
	resp, err := client.Get(fmt.Sprintf("%s/api/tags", m.ollamaURL))
	if err != nil {
		return fmt.Errorf("Ollama is not running. Please start Ollama first: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama is running but not responding correctly (status: %d)", resp.StatusCode)
	}
	
	return nil
}

// ListLocalModels lists all available local models
func (m *LocalLLMManager) ListLocalModels() ([]LocalModel, error) {
	if err := m.CheckOllamaAvailability(); err != nil {
		return nil, err
	}
	
	client := &http.Client{Timeout: m.timeout}
	resp, err := client.Get(fmt.Sprintf("%s/api/tags", m.ollamaURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %v", err)
	}
	defer resp.Body.Close()
	
	var modelResp LocalModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	
	// Add backend information
	for i := range modelResp.Models {
		modelResp.Models[i].Backend = "ollama"
		modelResp.Models[i].Status = "available"
	}
	
	return modelResp.Models, nil
}

// PullModel pulls a model from Ollama
func (m *LocalLLMManager) PullModel(modelName string) error {
	if err := m.CheckOllamaAvailability(); err != nil {
		return err
	}
	
	fmt.Printf("📥 Pulling model: %s\n", modelName)
	
	// Use ollama CLI to pull the model
	cmd := exec.Command("ollama", "pull", modelName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull model '%s': %v", modelName, err)
	}
	
	fmt.Printf("✅ Model '%s' pulled successfully\n", modelName)
	return nil
}

// RemoveModel removes a local model
func (m *LocalLLMManager) RemoveModel(modelName string) error {
	if err := m.CheckOllamaAvailability(); err != nil {
		return err
	}
	
	fmt.Printf("🗑️  Removing model: %s\n", modelName)
	
	cmd := exec.Command("ollama", "rm", modelName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove model '%s': %v", modelName, err)
	}
	
	fmt.Printf("✅ Model '%s' removed successfully\n", modelName)
	return nil
}

// TestModel tests if a local model is working
func (m *LocalLLMManager) TestModel(modelName string) error {
	if err := m.CheckOllamaAvailability(); err != nil {
		return err
	}
	
	fmt.Printf("🧪 Testing model: %s\n", modelName)
	
	// Simple test prompt
	testPrompt := "Hello, this is a test. Please respond with 'Test successful' if you can see this message."
	
	// Use ollama CLI to test the model
	cmd := exec.Command("ollama", "run", modelName, testPrompt)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("model test failed: %v", err)
	}
	
	response := strings.TrimSpace(string(output))
	fmt.Printf("✅ Model test successful. Response: %s\n", response)
	
	return nil
}

// GetModelInfo gets detailed information about a local model
func (m *LocalLLMManager) GetModelInfo(modelName string) (*LocalModel, error) {
	models, err := m.ListLocalModels()
	if err != nil {
		return nil, err
	}
	
	for _, model := range models {
		if model.Name == modelName {
			return &model, nil
		}
	}
	
	return nil, fmt.Errorf("model '%s' not found", modelName)
}

// IsModelAvailable checks if a specific model is available
func (m *LocalLLMManager) IsModelAvailable(modelName string) bool {
	_, err := m.GetModelInfo(modelName)
	return err == nil
}

// GetRecommendedModels returns a list of recommended models for different use cases
func (m *LocalLLMManager) GetRecommendedModels() map[string][]string {
	return map[string][]string{
		"chatbot": {
			"llama2",
			"llama2:7b",
			"llama2:13b",
			"mistral",
			"mistral:7b",
		},
		"code": {
			"codellama",
			"codellama:7b",
			"codellama:13b",
			"wizardcoder",
		},
		"general": {
			"llama2",
			"mistral",
			"neural-chat",
			"orca-mini",
		},
		"fast": {
			"llama2:7b",
			"mistral:7b",
			"orca-mini:3b",
			"phi",
		},
	}
}

// GetModelSize gets the size of a model in human-readable format
func (m *LocalLLMManager) GetModelSize(modelName string) string {
	info, err := m.GetModelInfo(modelName)
	if err != nil {
		return "unknown"
	}
	return info.Size
}

// ValidateModelName validates if a model name is valid for Ollama
func (m *LocalLLMManager) ValidateModelName(modelName string) error {
	if modelName == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	
	// Check for basic format
	if strings.Contains(modelName, " ") {
		return fmt.Errorf("model name cannot contain spaces")
	}
	
	// Check if it's a valid Ollama model format
	parts := strings.Split(modelName, ":")
	if len(parts) > 2 {
		return fmt.Errorf("invalid model name format. Use 'model' or 'model:tag'")
	}
	
	return nil
}
