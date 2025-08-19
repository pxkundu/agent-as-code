# Agent Configuration Guide

Agents are configured via `agent.yaml` files using a declarative YAML syntax. This file defines all aspects of an agent including its runtime environment, AI model configuration, capabilities, and deployment settings.

## Overview

The `agent.yaml` file serves as the single source of truth for agent configuration, providing:
- **Runtime Definition**: Execution environment and dependencies
- **Model Configuration**: AI model provider and parameters
- **Capability Specification**: Agent abilities and features
- **Deployment Settings**: Ports, volumes, and resource limits
- **Environment Configuration**: Variables and secrets

## File Structure

### Basic Structure
```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: my-agent
  version: 1.0.0
  description: "My AI Agent"
  author: "your-email@example.com"
  tags:
    - chatbot
    - ai
spec:
  runtime: python
  model:
    provider: openai
    name: gpt-4
  capabilities:
    - text-generation
    - chat
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
  environment:
    - name: OPENAI_API_KEY
      from: secret
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
```

## Configuration Sections

### 1. API Version and Kind

**Required fields** that identify the configuration type:

```yaml
apiVersion: agent.dev/v1  # API version (required)
kind: Agent               # Must be "Agent" (required)
```

### 2. Metadata

**Agent identification and metadata**:

```yaml
metadata:
  name: my-agent                # Agent name (required)
  version: 1.0.0               # Semantic version (optional)
  description: "My AI Agent"   # Description (optional)
  author: "name@email.com"     # Author contact (optional)
  tags:                        # Search tags (optional)
    - chatbot
    - ai
    - text-generation
  labels:                      # Key-value labels (optional)
    environment: production
    team: ai-team
    project: chatbot-v2
```

### 3. Specification (spec)

**Core agent configuration**:

```yaml
spec:
  runtime: python              # Runtime environment (required)
  model:                       # AI model configuration (required)
    provider: openai
    name: gpt-4
    config:
      temperature: 0.7
      max_tokens: 200
  capabilities:                # Agent capabilities (optional)
    - text-generation
    - chat
    - sentiment-analysis
  dependencies:                # Package dependencies (optional)
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
  environment:                 # Environment variables (optional)
    - name: OPENAI_API_KEY
      from: secret
    - name: LOG_LEVEL
      value: INFO
  ports:                       # Port configuration (optional)
    - container: 8080
      host: 8080
      protocol: tcp
  volumes:                     # Volume mounts (optional)
    - source: ./data
      target: /app/data
      type: bind
  healthCheck:                 # Health check (optional)
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
  resources:                   # Resource limits (optional)
    limits:
      cpu: "1"
      memory: "1Gi"
    requests:
      cpu: "500m"
      memory: "512Mi"
  config:                      # Custom configuration (optional)
    custom_setting: "value"
    feature_flags:
      beta_features: true
```

## Runtime Environments

### Supported Runtimes

The parser validates against these supported runtime environments:

| Runtime | Base Image | Package Manager | Entry Point | Status |
|---------|------------|-----------------|-------------|---------|
| `python` | `python:3.11-slim` | pip | `python main.py` | ✅ Stable |
| `nodejs` | `node:18-slim` | npm | `node index.js` | ✅ Stable |
| `go` | `golang:1.21-alpine` | Go modules | `./app` | ✅ Stable |
| `rust` | `rust:1.70-alpine` | Cargo | `./target/release/app` | 🔄 Beta |
| `java` | `openjdk:17-slim` | Maven/Gradle | `java -jar app.jar` | 🔄 Beta |

### Runtime-Specific Configuration

#### Python Runtime
```yaml
spec:
  runtime: python
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
    - pydantic==2.5.0
  environment:
    - name: PYTHONPATH
      value: "/app"
    - name: PYTHONUNBUFFERED
      value: "1"
```

#### Node.js Runtime
```yaml
spec:
  runtime: nodejs
  dependencies:
    - express@4.18.2
    - openai@4.20.1
    - dotenv@16.3.1
  environment:
    - name: NODE_ENV
      value: "production"
    - name: NODE_OPTIONS
      value: "--max-old-space-size=2048"
```

#### Go Runtime
```yaml
spec:
  runtime: go
  dependencies:
    - github.com/gin-gonic/gin@v1.9.1
    - github.com/sashabaranov/go-openai@v1.17.9
  environment:
    - name: GIN_MODE
      value: "release"
    - name: GO_ENV
      value: "production"
```

## Model Configuration

### Model Providers

#### OpenAI
```yaml
spec:
  model:
    provider: openai
    name: gpt-4
    config:
      temperature: 0.7
      max_tokens: 200
      top_p: 0.9
      frequency_penalty: 0.0
      presence_penalty: 0.0
```

#### Anthropic
```yaml
spec:
  model:
    provider: anthropic
    name: claude-3-opus
    config:
      temperature: 0.3
      max_tokens: 1000
      top_p: 0.9
```

#### Local Models (Ollama)
```yaml
spec:
  model:
    provider: local
    name: llama2
    config:
      base_url: "http://localhost:11434"
      temperature: 0.7
      top_k: 40
      top_p: 0.9
```

#### Custom Models
```yaml
spec:
  model:
    provider: custom
    name: my-custom-model
    config:
      endpoint: "https://api.example.com/v1"
      api_key: "${CUSTOM_API_KEY}"
      model_params:
        temperature: 0.5
        max_tokens: 500
```

### Model Configuration Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `temperature` | float | Randomness (0.0-2.0) | 1.0 |
| `max_tokens` | integer | Maximum response length | Model default |
| `top_p` | float | Nucleus sampling (0.0-1.0) | 1.0 |
| `top_k` | integer | Top-k sampling | Model default |
| `frequency_penalty` | float | Frequency penalty (-2.0-2.0) | 0.0 |
| `presence_penalty` | float | Presence penalty (-2.0-2.0) | 0.0 |

## Capabilities

### Supported Capabilities

| Capability | Description | Requirements |
|------------|-------------|--------------|
| `text-generation` | Generate text content | Text model |
| `chat` | Conversational AI | Chat model |
| `sentiment-analysis` | Analyze text sentiment | NLP model |
| `summarization` | Summarize text content | Text model |
| `translation` | Translate between languages | Multilingual model |
| `code-generation` | Generate code | Code model |
| `data-analysis` | Analyze structured data | Data model |
| `image-generation` | Generate images | Image model |
| `speech-to-text` | Convert speech to text | Audio model |
| `text-to-speech` | Convert text to speech | Audio model |

### Capability Configuration
```yaml
spec:
  capabilities:
    - text-generation
    - chat
    - sentiment-analysis
  
  # Capability-specific configuration
  capability_config:
    text-generation:
      max_length: 1000
      style: "creative"
    chat:
      context_window: 4096
      memory_enabled: true
    sentiment-analysis:
      languages: ["en", "es", "fr"]
      confidence_threshold: 0.8
```

## Dependencies

### Package Dependencies

#### Python Dependencies
```yaml
spec:
  runtime: python
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
    - pydantic==2.5.0
    - requests==2.31.0
    - python-dotenv==1.0.0
```

#### Node.js Dependencies
```yaml
spec:
  runtime: nodejs
  dependencies:
    - express@4.18.2
    - openai@4.20.1
    - dotenv@16.3.1
    - cors@2.8.5
    - helmet@7.1.0
```

#### Go Dependencies
```yaml
spec:
  runtime: go
  dependencies:
    - github.com/gin-gonic/gin@v1.9.1
    - github.com/sashabaranov/go-openai@v1.17.9
    - github.com/joho/godotenv@v1.5.1
    - github.com/gin-contrib/cors@v1.4.0
```

### Dependency Management

#### Version Pinning
```yaml
spec:
  dependencies:
    # Exact versions (recommended for production)
    - openai==1.0.0
    - fastapi==0.104.0
    
    # Version ranges
    - requests>=2.28.0,<3.0.0
    
    # Latest compatible
    - python-dotenv
```

#### Development Dependencies
```yaml
spec:
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
  
  dev_dependencies:
    - pytest==7.4.0
    - black==23.7.0
    - flake8==6.0.0
    - mypy==1.5.0
```

## Environment Configuration

### Environment Variables

#### Direct Values
```yaml
spec:
  environment:
    - name: LOG_LEVEL
      value: "INFO"
    - name: DEBUG
      value: "false"
    - name: PORT
      value: "8080"
```

#### From Secrets
```yaml
spec:
  environment:
    - name: OPENAI_API_KEY
      from: "secret"
      secretName: "openai-secret"
      secretKey: "api-key"
    - name: DATABASE_PASSWORD
      from: "secret"
      secretName: "db-secret"
      secretKey: "password"
```

#### From ConfigMaps
```yaml
spec:
  environment:
    - name: APP_CONFIG
      from: "configmap"
      configMapName: "app-config"
      configMapKey: "config.yaml"
```

#### From Environment
```yaml
spec:
  environment:
    - name: NODE_ENV
      from: "env"
      envKey: "NODE_ENV"
    - name: DEBUG
      from: "env"
      envKey: "DEBUG"
      defaultValue: "false"
```

### Environment Variable Types

| Type | Description | Example |
|------|-------------|---------|
| `value` | Direct value assignment | `value: "production"` |
| `from: "secret"` | Kubernetes secret | `secretName: "api-key"` |
| `from: "configmap"` | Kubernetes ConfigMap | `configMapName: "config"` |
| `from: "env"` | Host environment | `envKey: "DEBUG"` |

## Port Configuration

### Port Mappings

#### Basic Port Mapping
```yaml
spec:
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
```

#### Multiple Ports
```yaml
spec:
  ports:
    - container: 8080
      host: 80
      protocol: tcp
    - container: 9090
      host: 9090
      protocol: tcp
    - container: 5432
      host: 5432
      protocol: tcp
```

#### Port Configuration Options

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `container` | integer | ✅ | Container port (1-65535) | `8080` |
| `host` | integer | ❌ | Host port (1-65535) | `80` |
| `protocol` | string | ❌ | Protocol (tcp/udp) | `tcp` |

### Port Validation Rules

- **Container ports**: Must be between 1-65535
- **Host ports**: Optional, defaults to container port
- **Protocol**: Optional, defaults to "tcp"
- **Port conflicts**: Host ports must be unique

## Volume Configuration

### Volume Mounts

#### Bind Mounts
```yaml
spec:
  volumes:
    - source: ./data
      target: /app/data
      type: bind
      readOnly: false
```

#### Named Volumes
```yaml
spec:
  volumes:
    - source: model-cache
      target: /app/models
      type: volume
      readOnly: false
```

#### Temporary Volumes
```yaml
spec:
  volumes:
    - source: temp-data
      target: /app/temp
      type: tmpfs
      size: "100Mi"
```

### Volume Configuration Options

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `source` | string | ✅ | Source path or volume name | `./data` |
| `target` | string | ✅ | Container mount path | `/app/data` |
| `type` | string | ❌ | Volume type | `bind` |
| `readOnly` | boolean | ❌ | Read-only mount | `false` |

### Volume Types

| Type | Description | Use Case |
|------|-------------|----------|
| `bind` | Host directory mount | Development, data access |
| `volume` | Docker-managed volume | Persistent data |
| `tmpfs` | Temporary filesystem | Temporary data |
| `configmap` | Kubernetes ConfigMap | Configuration data |

## Health Check Configuration

### Health Check Types

#### HTTP Health Check
```yaml
spec:
  healthCheck:
    type: http
    path: "/health"
    port: 8080
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
```

#### Command Health Check
```yaml
spec:
  healthCheck:
    type: command
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
```

#### TCP Health Check
```yaml
spec:
  healthCheck:
    type: tcp
    port: 8080
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
```

### Health Check Options

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `type` | string | ❌ | Check type (http/command/tcp) | `command` |
| `command` | array | ❌ | Command to execute | `["curl", "/health"]` |
| `path` | string | ❌ | HTTP path for health check | `/health` |
| `port` | integer | ❌ | Port for health check | Container port |
| `interval` | string | ❌ | Check interval | `30s` |
| `timeout` | string | ❌ | Check timeout | `10s` |
| `retries` | integer | ❌ | Retry count | `3` |
| `startPeriod` | string | ❌ | Start delay | `5s` |

## Resource Configuration

### Resource Limits

#### CPU and Memory Limits
```yaml
spec:
  resources:
    limits:
      cpu: "2"
      memory: "2Gi"
      ephemeral-storage: "1Gi"
    requests:
      cpu: "500m"
      memory: "512Mi"
      ephemeral-storage: "100Mi"
```

#### Resource Units

| Resource | Unit | Examples |
|----------|------|----------|
| `cpu` | Cores | `"1"`, `"500m"`, `"0.5"` |
| `memory` | Bytes | `"512Mi"`, `"1Gi"`, `"2GB"` |
| `ephemeral-storage` | Bytes | `"100Mi"`, `"1Gi"` |

#### Resource Validation

- **CPU**: Must be positive number or millicores
- **Memory**: Must be positive number with unit
- **Storage**: Must be positive number with unit
- **Requests**: Cannot exceed limits

### Resource Calculation Examples

| CPU Request | Memory Request | Use Case |
|-------------|----------------|----------|
| `"100m"` | `"128Mi"` | Lightweight agent |
| `"500m"` | `"512Mi"` | Standard agent |
| `"1"` | `"1Gi"` | Heavy agent |
| `"2"` | `"4Gi"` | Resource-intensive agent |

## Advanced Configuration

### Custom Configuration

#### Arbitrary Configuration
```yaml
spec:
  config:
    # Custom settings
    custom_setting: "value"
    feature_flags:
      beta_features: true
      experimental_mode: false
    
    # Model-specific config
    model_config:
      temperature: 0.7
      max_tokens: 1000
    
    # Agent behavior
    behavior:
      max_conversations: 100
      timeout: 30s
      retry_attempts: 3
```

#### Configuration Inheritance
```yaml
spec:
  config:
    base_config: "production"
    overrides:
      development:
        debug: true
        log_level: "DEBUG"
      production:
        debug: false
        log_level: "INFO"
```

### Scaling Configuration

#### Horizontal Scaling
```yaml
spec:
  scaling:
    replicas: 3
    minReplicas: 1
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70
    targetMemoryUtilizationPercentage: 80
```

#### Vertical Scaling
```yaml
spec:
  scaling:
    vertical:
      enabled: true
      minCPU: "100m"
      maxCPU: "2"
      minMemory: "128Mi"
      maxMemory: "4Gi"
```

## Configuration Validation

### Required Fields

The parser enforces these required fields:

1. **apiVersion**: Must be present
2. **kind**: Must be exactly "Agent"
3. **metadata.name**: Agent name is required
4. **spec.runtime**: Runtime environment is required
5. **spec.model.provider**: Model provider is required
6. **spec.model.name**: Model name is required

### Validation Rules

#### Runtime Validation
- Runtime must be one of: `python`, `nodejs`, `go`, `rust`, `java`
- Case-sensitive matching
- No version suffixes in runtime field

#### Port Validation
- Container ports: 1-65535
- Host ports: 1-65535 (if specified)
- Port numbers must be valid integers
- No duplicate host ports

#### Resource Validation
- CPU values: Positive numbers or millicores
- Memory values: Positive numbers with units
- Storage values: Positive numbers with units
- Requests cannot exceed limits

### Validation Errors

Common validation error messages:

| Error | Cause | Solution |
|-------|-------|----------|
| `apiVersion is required` | Missing apiVersion field | Add `apiVersion: agent.dev/v1` |
| `kind must be 'Agent'` | Incorrect kind value | Set `kind: Agent` |
| `metadata.name is required` | Missing agent name | Add `metadata.name: "my-agent"` |
| `spec.runtime is required` | Missing runtime | Add `spec.runtime: "python"` |
| `invalid runtime 'invalid'` | Unsupported runtime | Use supported runtime value |
| `invalid container port 0` | Invalid port number | Use port 1-65535 |

## Configuration Examples

### Complete Examples

#### Basic Chatbot Agent
```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: simple-chatbot
  version: 1.0.0
  description: "Simple chatbot using OpenAI GPT-4"
  author: "ai-team@example.com"
  tags:
    - chatbot
    - ai
    - gpt-4
spec:
  runtime: python
  model:
    provider: openai
    name: gpt-4
    config:
      temperature: 0.7
      max_tokens: 200
  capabilities:
    - chat
    - text-generation
  dependencies:
    - openai==1.0.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
  environment:
    - name: OPENAI_API_KEY
      from: secret
    - name: LOG_LEVEL
      value: INFO
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
  healthCheck:
    command: ["curl", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
  resources:
    limits:
      cpu: "1"
      memory: "1Gi"
    requests:
      cpu: "500m"
      memory: "512Mi"
```

#### Advanced Data Analysis Agent
```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: data-analyzer
  version: 2.1.0
  description: "Advanced data analysis agent with visualization"
  author: "data-team@example.com"
  tags:
    - data-analysis
    - visualization
    - pandas
    - matplotlib
spec:
  runtime: python
  model:
    provider: anthropic
    name: claude-3-opus
    config:
      temperature: 0.3
      max_tokens: 1000
  capabilities:
    - data-analysis
    - visualization
    - text-generation
  dependencies:
    - pandas==2.0.3
    - numpy==1.24.3
    - matplotlib==3.7.2
    - seaborn==0.12.2
    - scikit-learn==1.3.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
  environment:
    - name: ANTHROPIC_API_KEY
      from: secret
    - name: LOG_LEVEL
      value: INFO
    - name: DEBUG
      value: false
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
    - container: 9090
      host: 9090
      protocol: tcp
  volumes:
    - source: ./data
      target: /app/data
      type: bind
    - source: model-cache
      target: /app/models
      type: volume
  healthCheck:
    type: http
    path: "/health"
    port: 8080
    interval: 30s
    timeout: 10s
    retries: 3
    startPeriod: 5s
  resources:
    limits:
      cpu: "2"
      memory: "4Gi"
      ephemeral-storage: "2Gi"
    requests:
      cpu: "1"
      memory: "2Gi"
      ephemeral-storage: "1Gi"
  config:
    data_processing:
      max_file_size: "100MB"
      supported_formats: ["csv", "json", "parquet"]
      batch_size: 1000
    visualization:
      default_theme: "dark"
      max_plot_elements: 10000
      export_formats: ["png", "svg", "pdf"]
```

#### Local LLM Agent
```yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: local-llm-agent
  version: 1.0.0
  description: "Local LLM agent using Ollama"
  author: "local-ai@example.com"
  tags:
    - local-llm
    - ollama
    - llama2
spec:
  runtime: python
  model:
    provider: local
    name: llama2
    config:
      base_url: "http://localhost:11434"
      temperature: 0.7
      top_k: 40
      top_p: 0.9
  capabilities:
    - text-generation
    - chat
  dependencies:
    - requests==2.31.0
    - fastapi==0.104.0
    - uvicorn==0.24.0
  environment:
    - name: OLLAMA_BASE_URL
      value: "http://localhost:11434"
    - name: LOG_LEVEL
      value: INFO
    - name: MODEL_NAME
      value: "llama2"
  ports:
    - container: 8080
      host: 8080
      protocol: tcp
  healthCheck:
    command: ["curl", "http://localhost:11434/api/tags"]
    interval: 60s
    timeout: 30s
    retries: 2
  resources:
    limits:
      cpu: "1"
      memory: "2Gi"
    requests:
      cpu: "500m"
      memory: "1Gi"
  config:
    ollama:
      model_pull_on_start: true
      model_cleanup: false
      max_concurrent_requests: 5
```

## Best Practices

### 1. Configuration Structure
- Use consistent indentation (2 spaces)
- Group related configurations together
- Use meaningful names and descriptions
- Add comprehensive tags for discovery

### 2. Security
- Never hardcode secrets in configuration
- Use environment variables for sensitive data
- Implement proper access controls
- Regular security audits

### 3. Resource Management
- Set appropriate resource limits
- Monitor resource usage
- Implement auto-scaling when needed
- Use resource quotas

### 4. Validation
- Validate configurations before deployment
- Use configuration schemas
- Implement automated testing
- Monitor configuration drift

### 5. Documentation
- Document all configuration options
- Provide usage examples
- Maintain change logs
- Include troubleshooting guides

## Troubleshooting

### Common Configuration Issues

#### YAML Syntax Errors
```bash
# Validate YAML syntax
yamllint agent.yaml

# Check with Python
python -c "import yaml; yaml.safe_load(open('agent.yaml'))"
```

#### Validation Errors
```bash
# Validate agent configuration
agent build --validate-only .

# Check specific fields
agent inspect agent.yaml
```

#### Runtime Issues
```bash
# Check runtime compatibility
agent build --dry-run .

# Verify dependencies
agent build --check-deps .
```

### Debug Configuration

#### Enable Debug Mode
```bash
# Set debug environment variable
export AAC_LOG_LEVEL=debug

# Run with verbose output
agent build -v -t my-agent .
```

#### Configuration Inspection
```bash
# Inspect parsed configuration
agent inspect agent.yaml --format json

# Show configuration summary
agent inspect agent.yaml --summary
```

## See Also

- [Parser](./parser.md) - Configuration parsing and validation
- [Builder](./builder.md) - Building agents from configurations
- [Runtime](./runtime.md) - Executing agent configurations
- [Registry](./registry.md) - Distributing agent configurations
- [CLI Overview](./cli-overview.md) - Command-line configuration tools