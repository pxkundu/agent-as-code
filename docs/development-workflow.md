# Development Workflow

This guide explains the complete workflow for developing AI agents using Agent-as-Code.

## Prerequisites

- Go 1.21 or later
- Python 3.8 or later
- Docker (optional, for container-based deployment)
- Git

## Project Structure

A typical agent project has the following structure:

```
my-agent/
├── agent.yaml           # Agent configuration
├── src/
│   ├── main.py         # Main agent code
│   ├── capabilities/   # Agent capabilities
│   ├── models/        # Model integrations
│   └── utils/         # Utility functions
├── tests/             # Test files
├── requirements.txt   # Python dependencies
└── README.md         # Project documentation
```

## Development Steps

### 1. Initialize Project

```bash
# Create new agent from template
agent init my-agent --template chatbot

# Or create custom agent
agent init my-agent --runtime python
```

### 2. Configure Agent

Edit `agent.yaml`:

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
  capabilities:
    - text-generation
    - chat
```

### 3. Implement Agent Logic

Example `src/main.py`:

```python
from agent_as_code import Agent

class MyAgent(Agent):
    async def process(self, input: str) -> str:
        # Implement agent logic
        response = await self.model.generate(input)
        return response

if __name__ == "__main__":
    agent = MyAgent()
    agent.run()
```

### 4. Local Development

```bash
# Build agent
agent build -t my-agent:dev .

# Run locally
agent run my-agent:dev

# View logs
agent logs my-agent:dev

# Test agent
agent test my-agent:dev
```

### 5. Testing

Create tests in `tests/` directory:

```python
# tests/test_agent.py
from my_agent import MyAgent

async def test_process():
    agent = MyAgent()
    response = await agent.process("Hello")
    assert response is not None
```

Run tests:
```bash
agent test my-agent:dev
```

### 6. Debugging

```bash
# Run with debug logging
agent run my-agent:dev --log-level debug

# Inspect agent
agent inspect my-agent:dev

# Check health
agent health my-agent:dev
```

### 7. Version Control

```bash
# Initialize git
git init

# Add files
git add agent.yaml src/ tests/

# Commit
git commit -m "Initial agent implementation"
```

### 8. Deployment

```bash
# Build production version
agent build -t my-agent:1.0.0 .

# Push to registry
agent push my-agent:1.0.0

# Deploy
agent deploy my-agent:1.0.0
```

## Development Best Practices

### 1. Configuration Management

- Use version control for `agent.yaml`
- Keep secrets in environment variables
- Document configuration options

### 2. Code Organization

- Follow template structure
- Separate concerns
- Use type hints
- Add documentation

### 3. Testing

- Write unit tests
- Test capabilities separately
- Use mock models for testing
- Test error handling

### 4. Error Handling

```python
from agent_as_code import Agent, AgentError

class MyAgent(Agent):
    async def process(self, input: str) -> str:
        try:
            response = await self.model.generate(input)
            return response
        except Exception as e:
            raise AgentError(f"Processing failed: {e}")
```

### 5. Logging

```python
import logging
from agent_as_code import Agent

class MyAgent(Agent):
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        
    async def process(self, input: str) -> str:
        self.logger.info(f"Processing input: {input}")
        response = await self.model.generate(input)
        self.logger.info(f"Generated response: {response}")
        return response
```

### 6. Resource Management

- Set appropriate resource limits
- Monitor memory usage
- Handle cleanup properly

### 7. Security

- Validate inputs
- Sanitize outputs
- Use secure dependencies
- Keep dependencies updated

## Common Development Tasks

### Adding New Capabilities

1. Create capability module:
```python
# src/capabilities/summarization.py
class SummarizationCapability:
    async def summarize(self, text: str) -> str:
        # Implement summarization
        pass
```

2. Update agent.yaml:
```yaml
spec:
  capabilities:
    - summarization
```

### Integrating New Models

1. Add model configuration:
```yaml
spec:
  model:
    provider: anthropic
    name: claude-3
```

2. Implement model integration:
```python
from agent_as_code.models import ModelIntegration

class ClaudeIntegration(ModelIntegration):
    async def generate(self, prompt: str) -> str:
        # Implement Claude integration
        pass
```

### Adding Health Checks

```yaml
spec:
  healthCheck:
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
```

## Troubleshooting

### Common Issues

1. Build Failures
```bash
# Check build logs
agent build -t my-agent:dev . --verbose
```

2. Runtime Errors
```bash
# Check agent logs
agent logs my-agent:dev
```

3. Model Issues
```bash
# Test model connection
agent test my-agent:dev --test-model
```

## See Also

- [CLI Reference](./cli-overview.md)
- [Configuration Guide](./agent-configuration.md)
- [Templates Guide](./templates.md)
- [API Reference](./api-reference.md)
