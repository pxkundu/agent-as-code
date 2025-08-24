# Python Examples

This directory contains Python examples demonstrating the enhanced LLM commands available in Agent as Code v1.1.0+.

## 🚀 **Enhanced LLM Commands Demo**

### `enhanced_llm_demo.py`

A comprehensive demonstration of all the new enhanced LLM features:

- **Intelligent Agent Creation**: Create fully functional agents with AI-powered generation
- **Model Optimization**: Optimize models for specific use cases
- **Benchmarking**: Comprehensive model performance testing
- **Agent Deployment**: Automated deployment and testing
- **Model Analysis**: Deep insights into model capabilities

### Running the Demo

```bash
# Install the package
pip install agent-as-code>=1.1.0

# Run the demo
python enhanced_llm_demo.py
```

### Prerequisites

- Python 3.8+
- Agent as Code v1.1.0+
- Ollama running locally (optional, for full functionality)
- Local LLM models (optional, for full functionality)

## 📋 **Example Use Cases**

### 1. **Chatbot Agent**
```python
from agent_as_code import AgentCLI

cli = AgentCLI()
cli.create_agent('chatbot')
cli.deploy_agent('chatbot-agent')
```

### 2. **Sentiment Analyzer**
```python
from agent_as_code import AgentCLI

cli = AgentCLI()
cli.create_agent('sentiment-analyzer')
cli.deploy_agent('sentiment-analyzer-agent')
```

### 3. **Workflow Automation**
```python
from agent_as_code import AgentCLI

cli = AgentCLI()
cli.create_agent('workflow-automation')
cli.deploy_agent('workflow-automation-agent')
```

### 4. **Model Optimization**
```python
from agent_as_code import AgentCLI

cli = AgentCLI()
cli.optimize_model('llama2', 'chatbot')
cli.optimize_model('mistral:7b', 'code-generation')
```

### 5. **Comprehensive Benchmarking**
```python
from agent_as_code import AgentCLI

cli = AgentCLI()
cli.benchmark_models(['chatbot', 'code-generation', 'analysis'])
```

## 🔧 **Configuration**

### Environment Variables

```bash
# Set log level
export LOG_LEVEL=DEBUG

# Set custom binary path (optional)
export AGENT_BINARY_PATH=/path/to/agent
```

### Custom Binary Path

```python
from agent_as_code import AgentCLI

# Use custom binary path
cli = AgentCLI(binary_path="/custom/path/to/agent")
```

## 🧪 **Testing**

### Run Tests

```bash
# Install test dependencies
pip install -r requirements.txt

# Run tests
pytest tests/
```

### Test Coverage

```bash
# Run with coverage
pytest --cov=agent_as_code tests/
```

## 📚 **Documentation**

For more information about the enhanced LLM commands:

- **Package Documentation**: `python -c "import agent_as_code; help(agent_as_code)"`
- **CLI Help**: `agent llm --help`
- **Command Help**: `agent llm create-agent --help`

## 🚨 **Troubleshooting**

### Common Issues

1. **Binary Not Found**
   ```bash
   # Ensure the Go binary is available
   which agent
   ```

2. **Ollama Not Running**
   ```bash
   # Start Ollama
   ollama serve
   ```

3. **Models Not Available**
   ```bash
   # Pull required models
   ollama pull llama2:7b
   ollama pull mistral:7b
   ```

4. **Permission Issues**
   ```bash
   # Make binary executable
   chmod +x /path/to/agent
   ```

### Getting Help

- Check the main package documentation
- Review the Go binary help: `agent --help`
- Check the Python package help: `python -c "import agent_as_code; help(agent_as_code.AgentCLI)"`

## 🔄 **Version Compatibility**

This example requires:
- **Agent as Code**: v1.1.0+
- **Python**: 3.8+
- **Go Binary**: Compatible version

For older versions, some features may not be available.
