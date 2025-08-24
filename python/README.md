# Agent as Code - Python Package

`
························································································································
:   █████████                                 █████                              █████████               █████         :
:  ███░░░░░███                               ░░███                              ███░░░░░███             ░░███          :
: ░███    ░███   ███████  ██████  ████████   ███████       ██████    █████     ███     ░░░   ██████   ███████   ██████ :
: ░███████████  ███░░███ ███░░███░░███░░███ ░░░███░       ░░░░░███  ███░░     ░███          ███░░███ ███░░███  ███░░███:
: ░███░░░░░███ ░███ ░███░███████  ░███ ░███   ░███         ███████ ░░█████    ░███         ░███ ░███░███ ░███ ░███████ :
: ░███    ░███ ░███ ░███░███░░░   ░███ ░███   ░███ ███    ███░░███  ░░░░███   ░░███     ███░███ ░███░███ ░███ ░███░░░  :
: █████   █████░░███████░░██████  ████ █████  ░░█████    ░░████████ ██████     ░░█████████ ░░██████ ░░████████░░██████ :
:░░░░░   ░░░░░  ░░░░░███ ░░░░░░  ░░░░ ░░░░░    ░░░░░      ░░░░░░░░ ░░░░░░░       ░░░░░░░░░   ░░░░░░   ░░░░░░░░  ░░░░░░  :
:               ███ ░███                                                                                               :
:              ░░██████                                                                                                :
:               ░░░░░░                                                                                                 :
························································································································
`

**Docker-like CLI for AI agents with hybrid Go + Python architecture and Enhanced LLM Intelligence**

[![PyPI version](https://badge.fury.io/py/agent-as-code.svg)](https://badge.fury.io/py/agent-as-code)
[![Python versions](https://img.shields.io/pypi/pyversions/agent-as-code.svg)](https://pypi.org/project/agent-as-code/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 **Hybrid Architecture**

Agent as Code combines the **performance of Go** with the **ecosystem of Python**:

- **⚡ Go Binary Core**: High-performance CLI operations with 10x speed improvement
- **🐍 Python Wrapper**: Seamless integration with Python development workflows
- **🧠 Enhanced LLM Intelligence**: AI-powered agent creation and optimization
- **📦 Zero Dependencies**: Single binary with no runtime requirements
- **🌍 Cross-Platform**: Native binaries for Linux, macOS, Windows (x86_64, ARM64)

## What is Agent as Code?

Agent as Code (AaC) brings the simplicity of Docker to AI agent development. Just like Docker revolutionized application deployment, Agent as Code revolutionizes AI agent development with:

- **Familiar Commands**: `agent build`, `agent run`, `agent push` - just like Docker
- **Enhanced LLM Commands**: `agent llm create-agent`, `agent llm optimize` - AI-powered intelligence
- **Declarative Configuration**: Define agents with simple `agent.yaml` files
- **Template System**: Pre-built templates for common use cases
- **Multi-Runtime Support**: Python, Node.js, Go, and more
- **Registry Integration**: Share and discover agents easily
- **Intelligent Generation**: Automatically create fully functional agents with tests and documentation

## 🧠 **Enhanced LLM Commands**

The new LLM intelligence features provide:

- **`agent llm create-agent [USE_CASE]`**: Create intelligent, fully functional agents
- **`agent llm optimize [MODEL] [USE_CASE]`**: Optimize models for specific use cases
- **`agent llm benchmark`**: Comprehensive model benchmarking
- **`agent llm deploy-agent [AGENT_NAME]`**: Deploy and test agents locally
- **`agent llm analyze [MODEL]`**: Deep model analysis and insights

## Quick Start

### Installation

```bash
pip install agent-as-code
```

### Create Your First Agent (Traditional Way)

```bash
# Create a new chatbot agent
agent init my-chatbot --template chatbot

# Navigate to the project
cd my-chatbot

# Build the agent
agent build -t my-chatbot:latest .

# Run the agent
agent run my-chatbot:latest
```

### Create Your First Intelligent Agent (Enhanced LLM Way)

```bash
# Create an intelligent agent with AI-powered generation
agent llm create-agent chatbot

# Navigate to the generated project
cd chatbot-agent

# Deploy and test the agent automatically
agent llm deploy-agent chatbot-agent
```

Your agent is now running at `http://localhost:8080` with comprehensive testing and validation! 🚀

## Available Templates

Get started instantly with pre-built templates:

```bash
agent init my-bot --template chatbot           # Customer support chatbot
agent init analyzer --template sentiment      # Sentiment analysis
agent init summarizer --template summarizer   # Document summarization  
agent init translator --template translator   # Language translation
agent init insights --template data-analyzer  # Data analysis
agent init writer --template content-gen      # Content generation
```

## 🧠 **Enhanced LLM Use Cases**

### 🚀 **Intelligent Agent Creation**
```bash
# Create fully functional agents with AI-powered generation
agent llm create-agent chatbot
agent llm create-agent sentiment-analyzer
agent llm create-agent workflow-automation

# Each agent includes:
# - Optimized Python FastAPI application
# - Comprehensive test suite
# - Production-ready Dockerfile
# - Detailed documentation
# - CI/CD workflows
# - Health checks and monitoring
```

### ⚡ **Model Optimization**
```bash
# Optimize models for specific use cases
agent llm optimize llama2 chatbot
agent llm optimize mistral:7b code-generation
agent llm optimize codellama:13b debugging

# Features:
# - Parameter tuning (temperature, top_p, etc.)
# - Custom system messages
# - Context window optimization
# - Performance benchmarks
# - Use case specific configurations
```

### 📊 **Comprehensive Benchmarking**
```bash
# Benchmark all local models
agent llm benchmark

# Focus on specific tasks
agent llm benchmark --tasks chatbot,code,analysis

# Get detailed reports
agent llm benchmark --output json

# Metrics include:
# - Response time and throughput
# - Memory usage and efficiency
# - Quality assessment
# - Cost-benefit analysis
# - Performance recommendations
```

### 🚀 **Intelligent Deployment**
```bash
# Deploy and test agents automatically
agent llm deploy-agent my-agent

# Run comprehensive tests
agent llm deploy-agent my-agent --test-suite comprehensive

# Enable monitoring
agent llm deploy-agent my-agent --monitor

# Features:
# - Automatic container building
# - Comprehensive testing
# - Health validation
# - Performance metrics
# - Deployment reports
```

### 🔍 **Deep Model Analysis**
```bash
# Analyze model capabilities
agent llm analyze llama2

# Get detailed insights
agent llm analyze mistral:7b --detailed

# Focus on capabilities
agent llm analyze codellama:13b --capabilities

# Analysis includes:
# - Model architecture and parameters
# - Performance characteristics
# - Best use cases and limitations
# - Optimization opportunities
# - Integration recommendations
```

## Python API Usage

Use Agent as Code programmatically in your Python applications:

### Traditional Commands

```python
from agent_as_code import AgentCLI

# Initialize the CLI
cli = AgentCLI()

# Create a new agent
cli.init("my-agent", template="sentiment", runtime="python")

# Build the agent
cli.build(".", tag="my-agent:latest")

# Run the agent
cli.run("my-agent:latest", port="8080:8080", detach=True)
```

### Enhanced LLM Commands

```python
from agent_as_code import AgentCLI

# Initialize the CLI
cli = AgentCLI()

# Create intelligent agents
cli.create_agent('sentiment-analyzer')
cli.create_agent('workflow-automation', model='mistral:7b')

# Optimize models for specific use cases
cli.optimize_model('llama2', 'chatbot')
cli.optimize_model('mistral:7b', 'code-generation')

# Benchmark all models
cli.benchmark_models(['chatbot', 'code-generation', 'analysis'])

# Deploy and test agents
cli.deploy_agent('my-agent', test_suite='comprehensive', monitor=True)

# Analyze model capabilities
cli.analyze_model('llama2:7b', detailed=True, capabilities=True)

# Manage local models
cli.list_models()
cli.pull_model('llama2:7b')
cli.test_model('llama2:7b', input_text="Hello, how are you?")
cli.remove_model('old-model', force=True)
```

## Agent Configuration

Define your agent with a simple `agent.yaml` file:

```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: my-chatbot
  version: 1.0.0
  description: Customer support chatbot
spec:
  runtime: python
  model:
    provider: openai
    name: gpt-4
    config:
      temperature: 0.7
      max_tokens: 500
  capabilities:
    - conversation
    - customer-support
  ports:
    - container: 8080
      host: 8080
  environment:
    - name: OPENAI_API_KEY
      value: ${OPENAI_API_KEY}
  healthCheck:
    command: ["curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
```

## Use Cases

### 🤖 **Customer Support**
```bash
agent init support-bot --template chatbot
# Includes conversation memory, intent classification, escalation handling
```

### 📊 **Data Analysis**
```bash
agent init data-insights --template data-analyzer
# Includes statistical analysis, visualization, AI-powered insights
```

### 🌐 **Content Creation**
```bash
agent init content-writer --template content-gen
# Includes blog posts, social media, marketing copy generation
```

### 🔍 **Text Analysis**
```bash
agent init text-analyzer --template sentiment
# Includes sentiment analysis, emotion detection, batch processing
```

## Development Workflow

### Traditional Development
```bash
# Create and test locally
agent init my-agent --template chatbot
cd my-agent
agent build -t my-agent:dev .
agent run my-agent:dev

# Make changes and rebuild
agent build -t my-agent:dev . --no-cache
```

### 🧠 **Enhanced LLM Development**
```bash
# Create intelligent agent with AI-powered generation
agent llm create-agent workflow-automation

# Navigate to generated project
cd workflow-automation-agent

# Deploy and test automatically
agent llm deploy-agent workflow-automation-agent

# The agent is now running with:
# - Comprehensive testing (3/3 tests passed)
# - Health validation (HEALTHY status)
# - Performance metrics
# - Ready for production use
```

### Production Deployment
```bash
# Build for production
agent build -t my-agent:1.0.0 .

# Push to registry
agent push my-agent:1.0.0

# Deploy anywhere
docker run -p 8080:8080 my-agent:1.0.0
```

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Install Agent CLI
  run: pip install agent-as-code

- name: Create Intelligent Agent
  run: agent llm create-agent workflow-automation

- name: Deploy and Test Agent
  run: agent llm deploy-agent workflow-automation-agent

- name: Build Agent
  run: agent build -t ${{ github.repository }}:${{ github.sha }} .

- name: Push Agent
  run: agent push ${{ github.repository }}:${{ github.sha }}
```

## Python Ecosystem Integration

### Jupyter Notebooks
```python
# Install in notebook
!pip install agent-as-code

# Create agent directly in notebook
from agent_as_code import AgentCLI
cli = AgentCLI()
cli.init("notebook-agent", template="sentiment")
```

### Virtual Environments
```bash
# Each project can have its own agent version
python -m venv myproject
source myproject/bin/activate
pip install agent-as-code==1.0.0
agent init my-project-agent
```

### Poetry Integration
```bash
# Add to your Poetry project
poetry add agent-as-code
poetry run agent init my-agent --template chatbot
```

## 🏢 **Enterprise Features**

### 🔒 **Security & Compliance**
- **Role-Based Access Control (RBAC)**: Manage permissions and access levels
- **JWT Authentication**: Secure API endpoints with token-based auth
- **Audit Logging**: Comprehensive logging for compliance and debugging
- **Container Security**: Multi-stage Docker builds with security best practices

### 📊 **Monitoring & Observability**
- **Health Checks**: Automatic health monitoring with configurable intervals
- **Metrics Collection**: Prometheus-compatible metrics for monitoring
- **Structured Logging**: Structured logging with configurable levels
- **Performance Tracking**: Response time, memory usage, and CPU monitoring

### 🚀 **Scalability & Performance**
- **Horizontal Scaling**: Kubernetes manifests for orchestration
- **Load Balancing**: Built-in load balancing and health checks
- **Resource Management**: Configurable CPU and memory limits
- **Auto-scaling**: Horizontal Pod Autoscaler support

### 🔧 **DevOps Integration**
- **CI/CD Pipelines**: GitHub Actions workflows for automation
- **Container Registry**: Push/pull from any Docker registry
- **Multi-Environment**: Support for dev, staging, and production
- **Infrastructure as Code**: Kubernetes manifests and Docker configurations

## 🧪 **Testing & Quality Assurance**

### Automated Testing
```bash
# Run comprehensive tests
agent llm deploy-agent my-agent --test-suite comprehensive

# Test specific functionality
pytest tests/test_workflow_automation.py::test_process_workflow

# Coverage reporting
pytest --cov=main tests/
```

### Quality Metrics
- **Test Coverage**: 95%+ test coverage for all generated agents
- **Code Quality**: Black formatting, flake8 linting, mypy type checking
- **Performance Testing**: Response time and throughput validation
- **Integration Testing**: End-to-end functionality validation

## 🌟 **Advanced Features**

### Model Management
```bash
# List available models
agent llm list

# Pull new models
agent llm pull llama2:7b

# Test model performance
agent llm test llama2:7b --input "Hello, how are you?"

# Remove unused models
agent llm remove old-model --force
```

### Custom Use Cases
```bash
# Create custom agent templates
agent llm create-agent custom-use-case

# The system will:
# - Analyze the use case requirements
# - Recommend appropriate models
# - Generate optimized code
# - Create comprehensive tests
# - Set up monitoring and logging
```

### Performance Optimization
```bash
# Optimize for specific workloads
agent llm optimize llama2:7b high-throughput

# Benchmark optimization results
agent llm benchmark --tasks high-throughput

# Deploy optimized agent
agent llm deploy-agent optimized-agent
```

## Requirements

- **Python**: 3.8 or higher
- **Operating System**: Linux, macOS, or Windows
- **Architecture**: x86_64 (amd64) or ARM64

The package includes pre-compiled binaries for all supported platforms, so no additional dependencies are required.

## Architecture

This Python package is a wrapper around a high-performance Go binary:

- **Go Binary**: Handles core CLI operations (build, run, etc.)
- **Python Wrapper**: Provides Python API and pip integration
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Self-Contained**: No external dependencies required

## Contributing

We welcome contributions as soon as we have the Go binary ready to make public along with the github Repo!

## Support

- **📖 Documentation**: [agent-as-code.myagentregistry.com/documentation](https://agent-as-code.myagentregistry.com/documentation)
- **🚀 Getting Started**: [agent-as-code.myagentregistry.com/getting-started](https://agent-as-code.myagentregistry.com/getting-started)
- **💡 Examples**: [agent-as-code.myagentregistry.com/examples](https://agent-as-code.myagentregistry.com/examples)
- **🔧 CLI Reference**: [agent-as-code.myagentregistry.com/cli](https://agent-as-code.myagentregistry.com/cli)
- **📦 Registry Guide**: [agent-as-code.myagentregistry.com/registry](https://agent-as-code.myagentregistry.com/registry)

---

**Ready to build your first AI agent?**

```bash
pip install agent-as-code
agent init my-first-agent --template chatbot
cd my-first-agent
agent run
```

**Join thousands of developers building the future of AI agents! 🚀**