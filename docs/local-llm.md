# Local LLM Configuration Guide

This guide covers how to set up and use local Large Language Models (LLMs) with Agent-as-Code, enabling AI agent development without API costs or internet dependencies.

## Overview

Agent-as-Code integrates with Ollama to provide local LLM capabilities. This allows you to:
- Run AI agents completely offline
- Avoid API costs and rate limits
- Maintain complete data privacy
- Develop and test agents without external dependencies

## Prerequisites

### 1. Install Ollama

**macOS/Linux:**
```bash
curl -fsSL https://ollama.ai/install.sh | sh
```

**Windows:**
Download from [https://ollama.ai](https://ollama.ai)

**Docker:**
```bash
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

### 2. Start Ollama Service

```bash
# Start the Ollama service
ollama serve

# Verify it's running
curl http://localhost:11434/api/tags
```

## Quick Setup

Use the built-in setup command to get started:

```bash
# Interactive setup guide
agent llm setup
```

This command provides step-by-step instructions for:
1. Installing Ollama
2. Starting the service
3. Pulling your first model
4. Testing the setup
5. Creating your first local AI agent

## Local LLM Commands

### List Available Models

```bash
# List all local models
agent llm list

# Example output:
# 🤖 Available Local Models
# =========================
# 
# llama2
#   Size:     3.8GB
#   Backend:  ollama
#   Status:   available
#   Modified: 2024-01-15T10:30:00Z
# 
# mistral:7b
#   Size:     4.1GB
#   Backend:  ollama
#   Status:   available
#   Modified: 2024-01-15T11:45:00Z
```

### Pull Models

```bash
# Pull a specific model
agent llm pull llama2

# Pull with specific tag
agent llm pull llama2:7b
agent llm pull mistral:7b
agent llm pull codellama:13b

# Pull code-specific models
agent llm pull codellama
agent llm pull wizardcoder
```

### Test Models

```bash
# Test if a model is working
agent llm test llama2

# Example output:
# 🧪 Testing model: llama2
# ✅ Model test successful. Response: Test successful
```

### Get Model Information

```bash
# Detailed model info
agent llm info llama2

# Example output:
# 📋 Model Information: llama2
# ========================
# Size:       3.8GB
# Backend:    ollama
# Status:     available
# Modified:   2024-01-15T10:30:00Z
# Digest:     sha256:abc123...
```

### Remove Models

```bash
# Remove a model to free disk space
agent llm remove llama2
```

## Model Recommendations

Get intelligent recommendations based on your use case:

```bash
# Get recommendations for different use cases
agent llm recommend chatbot
agent llm recommend code
agent llm recommend general
agent llm recommend fast
```

### Use Case Categories

| Use Case | Description | Recommended Models |
|----------|-------------|-------------------|
| `chatbot` | Conversational AI agents | llama2, llama2:7b, llama2:13b, mistral, mistral:7b |
| `code` | Code generation and analysis | codellama, codellama:7b, codellama:13b, wizardcoder |
| `general` | General-purpose tasks | llama2, mistral, neural-chat, orca-mini |
| `fast` | Quick responses, lower resource usage | llama2:7b, mistral:7b, orca-mini:3b, phi |

## Creating Agents with Local LLMs

### 1. Initialize Agent with Local Model

```bash
# Create a new agent using a local model
agent init my-chatbot --template chatbot --model local/llama2

# For code-focused agents
agent init my-code-assistant --template python-agent --model local/codellama
```

### 2. Configure agent.yaml for Local LLM

```yaml
apiVersion: agent.as.code/v1
kind: Agent
metadata:
  name: my-local-agent
  version: 1.0.0
  description: "Local AI agent using Ollama"

spec:
  runtime: python
  model:
    provider: local
    name: llama2
    config:
      temperature: 0.7
      max_tokens: 2048
      # Ollama-specific settings
      ollama_host: "http://localhost:11434"
      ollama_model: "llama2"
  
  capabilities:
    - chat
    - text-generation
  
  dependencies:
    - ollama-python
  
  environment:
    - name: OLLAMA_HOST
      value: "http://localhost:11434"
    - name: OLLAMA_MODEL
      value: "llama2"
```

### 3. Build and Run Local Agent

```bash
# Build the agent
agent build -t my-local-agent:latest .

# Run the agent
agent run my-local-agent:latest

# Test the agent
agent test my-local-agent:latest
```

## Local LLM Integration in Python Agents

### Basic Ollama Integration

```python
import requests
import json

class LocalLLMClient:
    def __init__(self, host="http://localhost:11434", model="llama2"):
        self.host = host
        self.model = model
    
    def generate(self, prompt, temperature=0.7, max_tokens=2048):
        """Generate text using local Ollama model"""
        url = f"{self.host}/api/generate"
        
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
            "options": {
                "temperature": temperature,
                "num_predict": max_tokens
            }
        }
        
        response = requests.post(url, json=payload)
        response.raise_for_status()
        
        return response.json()["response"]
    
    def chat(self, messages, temperature=0.7):
        """Chat completion using local model"""
        url = f"{self.host}/api/chat"
        
        payload = {
            "model": self.model,
            "messages": messages,
            "stream": False,
            "options": {
                "temperature": temperature
            }
        }
        
        response = requests.post(url, json=payload)
        response.raise_for_status()
        
        return response.json()["message"]["content"]

# Usage in your agent
def main():
    llm = LocalLLMClient(model="llama2")
    
    # Simple text generation
    response = llm.generate("Write a Python function to calculate fibonacci numbers")
    print(response)
    
    # Chat completion
    messages = [
        {"role": "user", "content": "Hello, how are you?"}
    ]
    response = llm.chat(messages)
    print(response)

if __name__ == "__main__":
    main()
```

### Advanced Integration with Streaming

```python
import requests
import json
import sys

class StreamingLocalLLM:
    def __init__(self, host="http://localhost:11434", model="llama2"):
        self.host = host
        self.model = model
    
    def stream_generate(self, prompt, temperature=0.7, max_tokens=2048):
        """Stream text generation from local model"""
        url = f"{self.host}/api/generate"
        
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": True,
            "options": {
                "temperature": temperature,
                "num_predict": max_tokens
            }
        }
        
        response = requests.post(url, json=payload, stream=True)
        response.raise_for_status()
        
        for line in response.iter_lines():
            if line:
                data = json.loads(line.decode('utf-8'))
                if 'response' in data:
                    yield data['response']
                if data.get('done', False):
                    break

# Usage
def main():
    llm = StreamingLocalLLM(model="llama2")
    
    print("Generating response...")
    for chunk in llm.stream_generate("Explain quantum computing in simple terms"):
        print(chunk, end='', flush=True)
    print()

if __name__ == "__main__":
    main()
```

## Model Configuration Options

### Temperature and Sampling

```yaml
spec:
  model:
    provider: local
    name: llama2
    config:
      temperature: 0.1    # Low for consistent outputs
      temperature: 0.7    # Balanced creativity
      temperature: 1.0    # High for creative tasks
      top_p: 0.9         # Nucleus sampling
      top_k: 40          # Top-k sampling
      repeat_penalty: 1.1 # Prevent repetition
```

### Context and Memory

```yaml
spec:
  model:
    provider: local
    name: llama2
    config:
      context_length: 4096    # Model context window
      max_tokens: 2048        # Max generation length
      num_ctx: 4096          # Ollama context size
```

## Performance Optimization

### Model Selection by Hardware

| Hardware | Recommended Models | Memory Usage |
|----------|-------------------|--------------|
| **Low-end** (4GB RAM) | llama2:7b, mistral:7b, phi | 2-4GB |
| **Mid-range** (8GB RAM) | llama2:13b, codellama:13b | 4-8GB |
| **High-end** (16GB+ RAM) | llama2:70b, codellama:34b | 8-16GB |

### Memory Management

```bash
# Check available models and their sizes
agent llm list

# Remove unused models to free space
agent llm remove llama2:70b

# Pull optimized models for your hardware
agent llm pull llama2:7b  # Smaller, faster
```

### GPU Acceleration

Ollama automatically uses GPU acceleration when available:

```bash
# Check GPU support
ollama list

# Force CPU-only (if needed)
OLLAMA_HOST=0.0.0.0:11434 OLLAMA_ORIGINS=* ollama serve
```

## Troubleshooting

### Common Issues

**1. Ollama not running**
```bash
# Error: Ollama is not running. Please start Ollama first
# Solution:
ollama serve
```

**2. Model not found**
```bash
# Error: model 'llama2' not found
# Solution:
agent llm pull llama2
```

**3. Out of memory**
```bash
# Error: CUDA out of memory
# Solution: Use smaller model or increase swap space
agent llm pull llama2:7b  # Instead of llama2:70b
```

**4. Slow performance**
```bash
# Solutions:
# 1. Use smaller model
agent llm pull llama2:7b

# 2. Enable GPU acceleration
nvidia-smi  # Check GPU availability

# 3. Adjust model parameters
# Lower temperature, reduce max_tokens
```

### Health Checks

```bash
# Check Ollama status
curl http://localhost:11434/api/tags

# Test model availability
agent llm test llama2

# Verify agent can connect
agent llm list
```

## Best Practices

### 1. Model Selection
- Start with smaller models (7B) for development
- Use larger models (13B+) for production quality
- Match model size to your hardware capabilities

### 2. Resource Management
- Monitor memory usage with `agent llm list`
- Remove unused models regularly
- Use appropriate model sizes for your use case

### 3. Development Workflow
```bash
# Development setup
agent llm pull llama2:7b
agent init dev-agent --template python-agent --model local/llama2:7b

# Production setup
agent llm pull llama2:13b
agent init prod-agent --template python-agent --model local/llama2:13b
```

### 4. Testing Strategy
```bash
# Test model functionality
agent llm test llama2

# Test agent with local model
agent test my-local-agent:latest

# Validate performance
time agent llm test llama2
```

## Integration Examples

### Chatbot Agent with Local LLM

```yaml
# agent.yaml
apiVersion: agent.as.code/v1
kind: Agent
metadata:
  name: local-chatbot
  version: 1.0.0

spec:
  runtime: python
  model:
    provider: local
    name: llama2
    config:
      temperature: 0.7
      max_tokens: 1024
  
  capabilities:
    - chat
    - conversation
  
  dependencies:
    - requests
    - flask
  
  environment:
    - name: OLLAMA_HOST
      value: "http://localhost:11434"
    - name: OLLAMA_MODEL
      value: "llama2"
  
  ports:
    - host: 8080
      container: 8080
  
  healthCheck:
    path: /health
    port: 8080
```

### Code Assistant Agent

```yaml
# agent.yaml
apiVersion: agent.as.code/v1
kind: Agent
metadata:
  name: code-assistant
  version: 1.0.0

spec:
  runtime: python
  model:
    provider: local
    name: codellama
    config:
      temperature: 0.1
      max_tokens: 2048
  
  capabilities:
    - code-generation
    - code-review
    - debugging
  
  dependencies:
    - requests
    - pygments
  
  environment:
    - name: OLLAMA_HOST
      value: "http://localhost:11434"
    - name: OLLAMA_MODEL
      value: "codellama"
```

## See Also

- [LLM Integration Guide](./llm.md) - Cloud LLM providers
- [Agent Configuration](./agent-configuration.md) - Complete agent.yaml reference
- [CLI Overview](./cli-overview.md) - Command-line interface
- [Runtime Guide](./runtime.md) - Container execution
