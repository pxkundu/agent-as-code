# Integration Patterns

This guide covers common integration patterns and extension points for Agent-as-Code.

## Model Integration

### Custom Model Provider

Implement a custom model provider:

```python
from agent_as_code.models import ModelProvider, ModelResponse

class CustomModelProvider(ModelProvider):
    def __init__(self, config: dict):
        super().__init__(config)
        self.api_key = config.get("api_key")
        self.base_url = config.get("base_url")
    
    async def generate(self, prompt: str) -> ModelResponse:
        # Implement model generation
        response = await self._call_api(prompt)
        return ModelResponse(
            text=response["text"],
            usage=response["usage"]
        )
    
    async def stream(self, prompt: str) -> AsyncIterator[str]:
        # Implement streaming
        async for chunk in self._stream_api(prompt):
            yield chunk["text"]
```

Configuration:
```yaml
spec:
  model:
    provider: custom
    name: my-model
    config:
      api_key: ${MY_MODEL_API_KEY}
      base_url: https://api.mymodel.com
```

### Local Model Integration

Integrate with local models:

```python
from agent_as_code.models import LocalModelProvider

class OllamaProvider(LocalModelProvider):
    async def setup(self):
        # Initialize Ollama client
        self.client = OllamaClient(self.config["base_url"])
    
    async def generate(self, prompt: str) -> str:
        return await self.client.generate(
            model=self.config["model_name"],
            prompt=prompt
        )
```

Configuration:
```yaml
spec:
  model:
    provider: local
    name: llama2
    config:
      base_url: http://localhost:11434
```

## Runtime Integration

### Custom Runtime Environment

Create a custom runtime:

```python
from agent_as_code.runtime import Runtime, RuntimeConfig

class CustomRuntime(Runtime):
    def __init__(self, config: RuntimeConfig):
        super().__init__(config)
        
    async def start(self):
        # Initialize runtime
        await self._setup_environment()
        
    async def stop(self):
        # Cleanup
        await self._cleanup()
    
    async def execute(self, input: dict) -> dict:
        # Execute agent logic
        return await self._process(input)
```

Configuration:
```yaml
spec:
  runtime: custom:1.0
  config:
    custom_setting: value
```

### Container Runtime Extension

Extend container runtime:

```go
type CustomRuntime struct {
    runtime.Runtime
    customConfig map[string]interface{}
}

func (r *CustomRuntime) Run(options *runtime.RunOptions) (*runtime.ContainerInfo, error) {
    // Add custom logic before running container
    if err := r.preRun(options); err != nil {
        return nil, err
    }
    
    // Run container
    info, err := r.Runtime.Run(options)
    if err != nil {
        return nil, err
    }
    
    // Add custom logic after running container
    if err := r.postRun(info); err != nil {
        return nil, err
    }
    
    return info, nil
}
```

## API Integration

### Custom API Endpoints

Add custom API endpoints:

```python
from fastapi import FastAPI, APIRouter
from agent_as_code.api import AgentAPI

router = APIRouter()

@router.post("/custom")
async def custom_endpoint(request: dict):
    # Custom endpoint logic
    return {"result": "success"}

class CustomAPI(AgentAPI):
    def __init__(self):
        super().__init__()
        self.app.include_router(router, prefix="/api/v1")
```

### Middleware Integration

Add custom middleware:

```python
from starlette.middleware.base import BaseHTTPMiddleware

class CustomMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        # Pre-processing
        response = await call_next(request)
        # Post-processing
        return response

app.add_middleware(CustomMiddleware)
```

## Storage Integration

### Custom Storage Backend

Implement custom storage:

```python
from agent_as_code.storage import StorageBackend

class CustomStorage(StorageBackend):
    async def save(self, key: str, data: bytes) -> str:
        # Save data
        location = await self._store(key, data)
        return location
    
    async def load(self, key: str) -> bytes:
        # Load data
        return await self._retrieve(key)
    
    async def delete(self, key: str):
        # Delete data
        await self._remove(key)
```

Configuration:
```yaml
spec:
  storage:
    provider: custom
    config:
      endpoint: https://storage.example.com
      credentials: ${STORAGE_CREDENTIALS}
```

## Monitoring Integration

### Custom Metrics

Add custom metrics:

```python
from prometheus_client import Counter, Histogram
from agent_as_code.monitoring import Metrics

class CustomMetrics(Metrics):
    def __init__(self):
        self.requests = Counter(
            'custom_requests_total',
            'Total custom requests'
        )
        self.latency = Histogram(
            'custom_latency_seconds',
            'Request latency'
        )
    
    def record_request(self):
        self.requests.inc()
    
    def record_latency(self, duration: float):
        self.latency.observe(duration)
```

### Health Checks

Implement custom health checks:

```python
from agent_as_code.health import HealthCheck

class CustomHealthCheck(HealthCheck):
    async def check(self) -> bool:
        # Perform health check
        status = await self._check_dependencies()
        return status.is_healthy()
    
    async def get_metrics(self) -> dict:
        return {
            "custom_metric": await self._get_metric()
        }
```

## Security Integration

### Custom Authentication

Implement custom authentication:

```python
from agent_as_code.security import AuthProvider

class CustomAuth(AuthProvider):
    async def authenticate(self, credentials: dict) -> bool:
        # Verify credentials
        return await self._verify(credentials)
    
    async def authorize(self, token: str, resource: str) -> bool:
        # Check authorization
        return await self._check_access(token, resource)
```

Configuration:
```yaml
spec:
  security:
    auth_provider: custom
    config:
      auth_url: https://auth.example.com
      client_id: ${CLIENT_ID}
```

## Event Integration

### Custom Event Handlers

Implement custom event handling:

```python
from agent_as_code.events import EventHandler

class CustomEventHandler(EventHandler):
    async def handle(self, event: dict):
        # Process event
        if event["type"] == "custom":
            await self._process_custom_event(event)
    
    async def emit(self, event: dict):
        # Emit event
        await self._publish(event)
```

Configuration:
```yaml
spec:
  events:
    handler: custom
    config:
      broker_url: amqp://localhost
```

## Best Practices

### 1. Error Handling
- Implement proper error handling
- Use custom error types
- Provide detailed error messages

### 2. Configuration
- Use environment variables
- Support configuration files
- Validate configuration

### 3. Testing
- Write unit tests
- Implement integration tests
- Use mocks for external services

### 4. Documentation
- Document integration points
- Provide examples
- Include configuration options

### 5. Security
- Implement authentication
- Use secure communication
- Handle sensitive data properly

## See Also

- [API Reference](./api-reference.md)
- [Runtime Guide](./runtime.md)
- [Security Guide](./security.md)
- [Configuration Guide](./agent-configuration.md)
