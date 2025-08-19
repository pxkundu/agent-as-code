# Templates Guide

Agent-as-Code provides a powerful template system for quickly creating common types of AI agents. This guide explains how to use and create templates.

## Available Templates

| Template | Description | Use Case |
|----------|-------------|----------|
| `chatbot` | Conversational agent | Customer service, chat interfaces |
| `sentiment` | Sentiment analysis | Text analysis, feedback processing |
| `summarizer` | Text summarization | Document processing |
| `translator` | Language translation | Multi-language support |
| `data-analyzer` | Data analysis | Business intelligence |
| `content-gen` | Content generation | Marketing, creative writing |

## Using Templates

### List Available Templates

```bash
agent init --list-templates
```

### Create Agent from Template

```bash
agent init my-agent --template chatbot
```

## Template Details

### Chatbot Template

A conversational agent template with session management and API endpoints.

**Features:**
- REST API with FastAPI
- Session management
- Health checks
- Conversation history

**Structure:**
```
chatbot/
├── main.py              # Main agent code
├── requirements.txt     # Dependencies
└── README.md           # Documentation
```

**API Endpoints:**
```python
# main.py
class ChatRequest(BaseModel):
    message: str
    session_id: Optional[str] = None

class ChatResponse(BaseModel):
    response: str
    timestamp: str

@app.post("/chat")
async def chat(request: ChatRequest):
    # Chat endpoint implementation
```

### Sentiment Analysis Template

A template for analyzing text sentiment.

**Features:**
- Sentiment analysis API
- Confidence scores
- Batch processing
- Health monitoring

**Structure:**
```
sentiment/
├── main.py              # Main agent code
├── requirements.txt     # Dependencies
└── README.md           # Documentation
```

**API Endpoints:**
```python
# main.py
class SentimentRequest(BaseModel):
    text: str
    include_confidence: Optional[bool] = True

class SentimentResponse(BaseModel):
    sentiment: str
    confidence: float
    timestamp: str

@app.post("/analyze")
async def analyze_sentiment(request: SentimentRequest):
    # Sentiment analysis implementation
```

## Template Configuration

Each template includes:

### 1. agent.yaml Configuration

```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: template-name
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  capabilities:
    - template-specific-capability
```

### 2. Dependencies

```text
# requirements.txt
fastapi==0.104.0
uvicorn==0.24.0
pydantic==2.5.0
```

### 3. Documentation

```markdown
# Template Name

Description of the template and its capabilities.

## Usage
1. Installation steps
2. Configuration guide
3. Example usage
```

## Creating Custom Templates

Templates are embedded in the binary and follow a specific structure:

```go
//go:embed chatbot/* sentiment/*
var templateFS embed.FS
```

### Template Structure

```
template-name/
├── main.py              # Main implementation
├── requirements.txt     # Dependencies
├── README.md           # Documentation
└── template.yaml       # Template metadata
```

### Template Metadata

```yaml
# template.yaml
name: template-name
description: Template description
author: Your Name
version: 1.0.0
runtimes:
  - python
tags:
  - tag1
  - tag2
```

## Template Best Practices

### 1. Code Organization
- Clear file structure
- Modular design
- Type hints
- Documentation

### 2. Error Handling
```python
@app.exception_handler(Exception)
async def global_exception_handler(request, exc):
    return JSONResponse(
        status_code=500,
        content={"error": str(exc)}
    )
```

### 3. Health Checks
```python
@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat()
    }
```

### 4. Configuration
```python
class Settings(BaseSettings):
    model_provider: str
    model_name: str
    api_key: str
    
    class Config:
        env_file = ".env"
```

### 5. Logging
```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
```

## See Also

- [Development Workflow](./development-workflow.md)
- [Configuration Guide](./agent-configuration.md)
- [API Reference](./api-reference.md)
- [Examples](./example-agents.md)
