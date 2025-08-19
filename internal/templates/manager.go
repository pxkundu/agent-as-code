package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Template directory structure embedded in binary
//
//go:embed chatbot/* sentiment/*
var templateFS embed.FS

// AgentConfig represents the configuration for generating an agent
type AgentConfig struct {
	Name     string
	Template string
	Runtime  string
	Model    string
}

// Manager handles template operations
type Manager struct {
	templatesDir string
}

// New creates a new template manager
func New() *Manager {
	return &Manager{
		templatesDir: getTemplatesDir(),
	}
}

// NewWithDir creates a new template manager with custom templates directory
func NewWithDir(dir string) *Manager {
	return &Manager{
		templatesDir: dir,
	}
}

// Generate generates a new agent project from a template
func (m *Manager) Generate(projectDir string, config *AgentConfig) error {
	// Create agent.yaml
	if err := m.generateAgentYAML(projectDir, config); err != nil {
		return fmt.Errorf("failed to generate agent.yaml: %w", err)
	}

	// Copy template files (now handles embedded templates and fallback)
	if err := m.copyTemplateFiles("", projectDir, config); err != nil {
		return fmt.Errorf("failed to copy template files: %w", err)
	}

	return nil
}

// generateAgentYAML generates the agent.yaml file
func (m *Manager) generateAgentYAML(projectDir string, config *AgentConfig) error {
	// Parse model provider and name
	modelProvider, modelName := parseModel(config.Model)

	// Build template based on provider
	var agentYAMLTemplate string
	if modelProvider == "ollama" {
		agentYAMLTemplate = `apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: {{ .Name }}
  version: 1.0.0
  description: {{ .Name }} agent generated from {{ .Template }} template
spec:
  runtime: {{ .Runtime }}
  model:
    provider: {{ .ModelProvider }}
    name: {{ .ModelName }}
    config:
      temperature: 0.7
      max_tokens: 500
      base_url: "http://localhost:11434"
  capabilities:
    - {{ .Template }}
  dependencies:
    - requests==2.31.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
    - pydantic==2.5.0
  ports:
    - container: 8080
      host: 8080
  environment:
    - name: LOG_LEVEL
      value: INFO
    - name: OLLAMA_BASE_URL
      value: "http://localhost:11434"
    - name: MODEL_NAME
      value: "{{ .ModelName }}"
  healthCheck:
    command: ["curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
`
	} else {
		agentYAMLTemplate = `apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: {{ .Name }}
  version: 1.0.0
  description: {{ .Name }} agent generated from {{ .Template }} template
spec:
  runtime: {{ .Runtime }}
  model:
    provider: {{ .ModelProvider }}
    name: {{ .ModelName }}
    config:
      temperature: 0.7
      max_tokens: 500
  capabilities:
    - {{ .Template }}
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
    - pydantic==2.5.0
  ports:
    - container: 8080
      host: 8080
  environment:
    - name: LOG_LEVEL
      value: INFO
    - name: OPENAI_API_KEY
      from: secret
  healthCheck:
    command: ["curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
`
	}

	// Template data
	data := struct {
		Name          string
		Template      string
		Runtime       string
		ModelProvider string
		ModelName     string
	}{
		Name:          config.Name,
		Template:      config.Template,
		Runtime:       config.Runtime,
		ModelProvider: modelProvider,
		ModelName:     modelName,
	}

	// Parse template
	tmpl, err := template.New("agent.yaml").Parse(agentYAMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outputPath := filepath.Join(projectDir, "agent.yaml")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create agent.yaml: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// copyTemplateFiles copies template files to the project directory
func (m *Manager) copyTemplateFiles(templateDir, projectDir string, config *AgentConfig) error {
	// Use embedded templates
	templatePrefix := config.Template

	// Check if template directory exists in embedded FS
	entries, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read embedded templates: %w", err)
	}

	templateExists := false
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == templatePrefix {
			templateExists = true
			break
		}
	}

	if !templateExists {
		// For now, create basic template files
		return m.createBasicTemplate(projectDir, config)
	}

	// Walk through embedded template files
	return fs.WalkDir(templateFS, templatePrefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templatePrefix, path)
		if err != nil {
			return err
		}

		// Skip agent.yaml (we generate our own)
		if relPath == "agent.yaml" || relPath == "agent.yml" {
			return nil
		}

		// Create destination path
		destPath := filepath.Join(projectDir, relPath)

		// Create destination directory if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}

		// Read file from embedded FS
		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Write to destination
		return os.WriteFile(destPath, content, 0644)
	})
}

// ListTemplates returns available templates
func (m *Manager) ListTemplates() ([]string, error) {
	var templates []string

	// Read from embedded FS
	entries, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}

	return templates, nil
}

// GetTemplateInfo returns information about a template
func (m *Manager) GetTemplateInfo(templateName string) (*TemplateInfo, error) {
	// Check if template exists in embedded FS
	entries, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		// If embedded FS fails, assume template is supported via fallback
		return &TemplateInfo{
			Name:        templateName,
			Description: fmt.Sprintf("%s agent template", templateName),
			Runtimes:    []string{"python"}, // Default
		}, nil
	}

	templateExists := false
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == templateName {
			templateExists = true
			break
		}
	}

	// If template not found in embedded FS, check if it's a supported template
	supportedTemplates := []string{"chatbot", "sentiment", "summarizer", "translator", "data-analyzer", "content-gen"}
	if !templateExists {
		for _, supported := range supportedTemplates {
			if templateName == supported {
				return &TemplateInfo{
					Name:        templateName,
					Description: fmt.Sprintf("%s agent template (fallback)", templateName),
					Runtimes:    []string{"python"}, // Default
				}, nil
			}
		}
		return nil, fmt.Errorf("template '%s' not found", templateName)
	}

	// Read template metadata (if exists)
	metadataPath := filepath.Join(templateName, "template.yaml")
	if _, err := fs.Stat(templateFS, metadataPath); err == nil {
		return m.parseTemplateMetadata(metadataPath)
	}

	// Return basic info
	return &TemplateInfo{
		Name:        templateName,
		Description: fmt.Sprintf("%s agent template", templateName),
		Runtimes:    []string{"python"}, // Default
	}, nil
}

// TemplateInfo represents template information
type TemplateInfo struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author,omitempty"`
	Version     string   `yaml:"version,omitempty"`
	Runtimes    []string `yaml:"runtimes"`
	Tags        []string `yaml:"tags,omitempty"`
}

// parseTemplateMetadata parses template metadata file
func (m *Manager) parseTemplateMetadata(path string) (*TemplateInfo, error) {
	// This would parse template.yaml metadata file
	// For now, return basic info
	return &TemplateInfo{
		Name:        "template",
		Description: "Agent template",
		Runtimes:    []string{"python"},
	}, nil
}

// Helper functions
func getTemplatesDir() string {
	// For embedded templates, we don't need a directory path
	return ""
}

func parseModel(model string) (provider, name string) {
	// Parse model string like "openai/gpt-4" or "local/llama2"
	parts := strings.Split(model, "/")
	if len(parts) >= 2 {
		provider = parts[0]
		name = strings.Join(parts[1:], "/") // Handle cases like "local/models/llama2"

		// Map local to ollama provider
		if provider == "local" {
			provider = "ollama"
		}

		return provider, name
	}

	// Default to OpenAI for backward compatibility
	return "openai", model
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// createBasicTemplate creates a basic template when embedded templates are not available
func (m *Manager) createBasicTemplate(projectDir string, config *AgentConfig) error {
	// Create basic files based on template type
	switch config.Template {
	case "chatbot":
		return m.createChatbotTemplate(projectDir, config)
	case "sentiment":
		return m.createSentimentTemplate(projectDir, config)
	default:
		return m.createGenericTemplate(projectDir, config)
	}
}

// createChatbotTemplate creates a basic chatbot template
func (m *Manager) createChatbotTemplate(projectDir string, config *AgentConfig) error {
	// Create main.py
	mainPy := `#!/usr/bin/env python3
"""
Chatbot Agent - Generated by Agent-as-Code
"""

import os
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="Chatbot Agent")

class ChatRequest(BaseModel):
    message: str

class ChatResponse(BaseModel):
    response: str

@app.post("/chat")
async def chat(request: ChatRequest):
    # Basic echo response - replace with your logic
    return ChatResponse(response=f"Echo: {request.message}")

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
`

	// Create requirements.txt
	requirements := `fastapi==0.104.0
uvicorn==0.24.0
pydantic==2.5.0
`

	// Create README.md
	readme := "# " + config.Name + "\n\nA chatbot agent generated by Agent-as-Code.\n\n## Usage\n\n1. Install dependencies: `pip install -r requirements.txt`\n2. Run the agent: `python main.py`\n3. Test: `curl -X POST http://localhost:8080/chat -H \"Content-Type: application/json\" -d '{\"message\": \"Hello\"}'`"

	// Write files
	files := map[string]string{
		"main.py":          mainPy,
		"requirements.txt": requirements,
		"README.md":        readme,
	}

	for filename, content := range files {
		path := filepath.Join(projectDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// createSentimentTemplate creates a basic sentiment analysis template
func (m *Manager) createSentimentTemplate(projectDir string, config *AgentConfig) error {
	// Create main.py
	mainPy := `#!/usr/bin/env python3
"""
Sentiment Analysis Agent - Generated by Agent-as-Code
"""

import os
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="Sentiment Analysis Agent")

class SentimentRequest(BaseModel):
    text: str

class SentimentResponse(BaseModel):
    sentiment: str
    confidence: float

@app.post("/analyze")
async def analyze_sentiment(request: SentimentRequest):
    # Basic sentiment analysis - replace with your logic
    sentiment = "positive" if "good" in request.text.lower() else "negative"
    return SentimentResponse(sentiment=sentiment, confidence=0.8)

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
`

	// Create requirements.txt
	requirements := `fastapi==0.104.0
uvicorn==0.24.0
pydantic==2.5.0
`

	// Create README.md
	readme := "# " + config.Name + "\n\nA sentiment analysis agent generated by Agent-as-Code.\n\n## Usage\n\n1. Install dependencies: `pip install -r requirements.txt`\n2. Run the agent: `python main.py`\n3. Test: `curl -X POST http://localhost:8080/analyze -H \"Content-Type: application/json\" -d '{\"text\": \"This is good\"}'`"

	// Write files
	files := map[string]string{
		"main.py":          mainPy,
		"requirements.txt": requirements,
		"README.md":        readme,
	}

	for filename, content := range files {
		path := filepath.Join(projectDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// createGenericTemplate creates a generic template
func (m *Manager) createGenericTemplate(projectDir string, config *AgentConfig) error {
	// Create main.py
	mainPy := `#!/usr/bin/env python3
"""
` + config.Name + ` Agent - Generated by Agent-as-Code
"""

import os
from fastapi import FastAPI

app = FastAPI(title="` + config.Name + ` Agent")

@app.get("/")
async def root():
    return {"message": "` + config.Name + ` Agent is running"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
`

	// Create requirements.txt
	requirements := `fastapi==0.104.0
uvicorn==0.24.0
`

	// Create README.md
	readme := "# " + config.Name + "\n\nAn AI agent generated by Agent-as-Code.\n\n## Usage\n\n1. Install dependencies: `pip install -r requirements.txt`\n2. Run the agent: `python main.py`\n3. Access: `http://localhost:8080`"

	// Write files
	files := map[string]string{
		"main.py":          mainPy,
		"requirements.txt": requirements,
		"README.md":        readme,
	}

	for filename, content := range files {
		path := filepath.Join(projectDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}
