# Quick Start Guide

This guide will help you get started with Agent-as-Code in just a few minutes. You'll learn how to install the framework, create your first agent, and deploy it.

## Prerequisites

- Python 3.8 or later
- Docker (optional, for container-based deployment)
- An OpenAI API key or other LLM provider key

## Installation

Choose your preferred installation method:

### Option 1: Python Package (Recommended)
```bash
pip install agent-as-code
```

### Option 2: Direct Binary Download
```bash
# Linux/macOS
curl -L https://api.myagentregistry.com/install.sh | sh

# Windows (PowerShell)
iwr -useb https://api.myagentregistry.com/install.ps1 | iex
```

### Option 3: Build from Source
```bash
git clone https://github.com/pxkundu/agent-as-code
cd agent-as-code
make install
```

## Creating Your First Agent

### 1. Initialize a New Agent Project

```bash
# Create a new agent from the chatbot template
agent init my-chatbot --template chatbot
cd my-chatbot
```

### 2. Configure Your Agent

Edit `agent.yaml`:

```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: my-chatbot
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  environment:
    - name: OPENAI_API_KEY
      from: secret
```

### 3. Build and Run

```bash
# Build the agent
agent build -t my-chatbot:latest .

# Run locally
agent run my-chatbot:latest

# Test it
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how can you help me?"}'
```

### 4. Deploy to Production

```bash
# Push to registry
agent push my-chatbot:latest

# Deploy to cloud
agent deploy my-chatbot:latest --cloud aws --replicas 3
```

## Next Steps

- [Agent Configuration Guide](./agent-configuration.md) - Learn about all configuration options
- [Templates Guide](./templates.md) - Explore available templates and create your own
- [Deployment Guide](./deployment-strategies.md) - Learn about different deployment options
- [Examples](./example-agents.md) - See more example agents and use cases

## Common Commands

| Command         | Description        | Example                                  |
|-----------------|--------------------|------------------------------------------|
| `agent init`    | Create new agent   | `agent init my-bot --template chatbot`   |
| `agent build`   | Build agent        | `agent build -t my-bot:latest .`         |
| `agent run`     | Run locally        | `agent run my-bot:latest`                |
| `agent push`    | Push to registry   | `agent push my-bot:latest`               |
| `agent pull`    | Pull from registry | `agent pull my-bot:latest`               |
| `agent deploy`  | Deploy to cloud    | `agent deploy my-bot:latest --cloud aws` |

## Troubleshooting

### Common Issues

1. **Installation Problems**
   ```bash
   # Verify installation
   agent --version
   
   # Check Python version
   python --version
   ```

2. **Build Failures**
   ```bash
   # Check build logs
   agent build -t my-bot:latest . --verbose
   ```

3. **Runtime Errors**
   ```bash
   # Check agent logs
   agent logs my-bot:latest
   ```

### Getting Help

- Visit our [Community Forum](https://github.com/pxkundu/agent-as-code/discussions)
- Check the [FAQ](./faq.md)
- Report issues on [GitHub](https://github.com/pxkundu/agent-as-code/issues)

## Security Best Practices

1. Never commit API keys or secrets to version control
2. Use environment variables or secret management systems
3. Keep your Agent-as-Code installation updated
4. Follow the principle of least privilege when configuring agents

## What's Next?

After completing this quickstart guide, you can:

1. Explore more [advanced features](./advanced-patterns.md)
2. Learn about [local LLM integration](./local-llm.md)
3. Set up [monitoring and observability](./monitoring.md)
4. Join our [community](https://github.com/pxkundu/agent-as-code/discussions)

Remember to check our [documentation](./README.md) for detailed information about all features and capabilities.
