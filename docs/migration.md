# Migration Guide

This guide helps you migrate existing AI agents and projects to Agent-as-Code.

## Overview

Agent-as-Code provides a structured way to containerize and manage AI agents. This guide covers:
1. Converting existing agents to Agent-as-Code format
2. Migrating configuration
3. Updating dependencies
4. Testing and validation

## Converting Existing Agents

### 1. Project Structure

Convert your project to the Agent-as-Code structure:

Before:
```
my-agent/
├── app.py
├── requirements.txt
└── config.json
```

After:
```
my-agent/
├── agent.yaml           # Agent configuration
├── src/
│   └── main.py         # Main agent code
├── requirements.txt     # Dependencies
└── README.md           # Documentation
```

### 2. Configuration Migration

Convert configuration to `agent.yaml`:

Before (`config.json`):
```json
{
    "model": {
        "type": "gpt-4",
        "api_key": "sk-..."
    },
    "server": {
        "port": 8080
    }
}
```

After (`agent.yaml`):
```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: my-agent
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  environment:
    - name: OPENAI_API_KEY
      from: secret
  ports:
    - container: 8080
```

### 3. Code Migration

Update your code to use Agent-as-Code patterns:

Before (`app.py`):
```python
import openai
import json

with open("config.json") as f:
    config = json.load(f)

openai.api_key = config["model"]["api_key"]

def process_input(text):
    response = openai.ChatCompletion.create(
        model=config["model"]["type"],
        messages=[{"role": "user", "content": text}]
    )
    return response.choices[0].message.content
```

After (`src/main.py`):
```python
from agent_as_code import Agent

class MyAgent(Agent):
    async def process(self, input: str) -> str:
        response = await self.model.generate(input)
        return response

if __name__ == "__main__":
    agent = MyAgent()
    agent.run()
```

## Dependency Migration

### 1. Update Requirements

Update `requirements.txt`:

Before:
```text
openai==0.28.0
flask==2.0.1
```

After:
```text
agent-as-code==1.0.0
openai==1.0.0
fastapi==0.104.0
```

### 2. Environment Variables

Move from config files to environment variables:

Before:
```python
with open("config.json") as f:
    config = json.load(f)
```

After:
```yaml
# agent.yaml
spec:
  environment:
    - name: OPENAI_API_KEY
      from: secret
    - name: LOG_LEVEL
      value: INFO
```

## API Migration

### 1. HTTP Endpoints

Update API endpoints to use FastAPI:

Before:
```python
from flask import Flask, request

app = Flask(__name__)

@app.route("/process", methods=["POST"])
def process():
    text = request.json["text"]
    result = process_input(text)
    return {"result": result}
```

After:
```python
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class ProcessRequest(BaseModel):
    text: str

class ProcessResponse(BaseModel):
    result: str

@app.post("/process")
async def process(request: ProcessRequest) -> ProcessResponse:
    result = await agent.process(request.text)
    return ProcessResponse(result=result)
```

### 2. Health Checks

Add health check endpoint:

```python
@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "version": "1.0.0",
        "timestamp": datetime.now().isoformat()
    }
```

## Testing Migration

### 1. Update Test Structure

Before:
```
tests/
└── test_app.py
```

After:
```
tests/
├── test_agent.py
├── test_models.py
└── test_api.py
```

### 2. Update Test Code

Before:
```python
def test_process():
    result = process_input("Hello")
    assert result is not None
```

After:
```python
async def test_process():
    agent = MyAgent()
    result = await agent.process("Hello")
    assert result is not None
```

## Validation Steps

1. **Configuration Check**
```bash
# Validate agent.yaml
agent validate .
```

2. **Build Test**
```bash
# Build agent
agent build -t my-agent:test .
```

3. **Runtime Test**
```bash
# Run agent
agent run my-agent:test

# Test endpoints
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello"}'
```

## Common Issues

### 1. Configuration Issues

Problem:
```yaml
spec:
  model: gpt-4  # Incorrect format
```

Solution:
```yaml
spec:
  model:
    provider: openai
    name: gpt-4
```

### 2. Runtime Issues

Problem:
```python
# Direct model calls
response = openai.ChatCompletion.create(...)
```

Solution:
```python
# Use agent model interface
response = await self.model.generate(...)
```

### 3. Environment Issues

Problem:
```python
# Hardcoded configuration
api_key = "sk-..."
```

Solution:
```yaml
# agent.yaml
spec:
  environment:
    - name: OPENAI_API_KEY
      from: secret
```

## Best Practices

### 1. Configuration
- Use environment variables for secrets
- Keep configuration in agent.yaml
- Use proper versioning

### 2. Code Organization
- Follow template structure
- Use type hints
- Implement proper error handling

### 3. Testing
- Write unit tests
- Add integration tests
- Test with different configurations

### 4. Documentation
- Update README.md
- Document API endpoints
- Include configuration options

## See Also

- [Development Guide](./development-workflow.md)
- [Configuration Guide](./agent-configuration.md)
- [API Reference](./api-reference.md)
- [Examples](./example-agents.md)
