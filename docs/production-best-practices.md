# Production Best Practices

This guide outlines best practices for running Agent-as-Code in production environments. Following these guidelines will help ensure your agents are secure, performant, and maintainable.

## Table of Contents

- [Security](#security)
- [Performance](#performance)
- [Monitoring](#monitoring)
- [Deployment](#deployment)
- [Development](#development)
- [Operations](#operations)
- [Compliance](#compliance)

## Security

### 1. Secret Management

✅ **Do**
```yaml
# agent.yaml
spec:
  environment:
    - name: API_KEY
      from: secret
    - name: DATABASE_URL
      from: secret
```

❌ **Don't**
```yaml
spec:
  environment:
    - name: API_KEY
      value: "sk-1234..." # Never hardcode secrets
```

### 2. Access Control

✅ **Do**
```yaml
spec:
  security:
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    securityContext:
      capabilities:
        drop: ["ALL"]
```

### 3. Network Security

✅ **Do**
```yaml
spec:
  network:
    ingress:
      - port: 8080
        tls: true
    egress:
      - to: ["api.openai.com"]
        ports: [443]
```

### 4. Authentication & Authorization

```python
from agent_as_code import Agent, auth

class SecureAgent(Agent):
    @auth.require_token
    async def sensitive_operation(self, data: dict):
        # Only authenticated requests can access this
        pass
```

## Performance

### 1. Resource Management

✅ **Do**
```yaml
spec:
  resources:
    requests:
      memory: "512Mi"
      cpu: "250m"
    limits:
      memory: "1Gi"
      cpu: "500m"
```

### 2. Caching

```python
from agent_as_code import Agent, cache

class CachedAgent(Agent):
    @cache.ttl(minutes=5)
    async def expensive_operation(self, input: str):
        return await self.model.generate(input)
```

### 3. Rate Limiting

```python
from agent_as_code import Agent, rate_limit

class RateLimitedAgent(Agent):
    @rate_limit(requests=100, period="1m")
    async def api_operation(self, request: dict):
        return await self.process(request)
```

### 4. Connection Pooling

```python
from agent_as_code import Agent, db

class DatabaseAgent(Agent):
    def __init__(self):
        self.db = db.Pool(
            min_connections=5,
            max_connections=20,
            idle_timeout="5m"
        )
```

## Monitoring

### 1. Metrics Collection

```yaml
spec:
  monitoring:
    metrics:
      prometheus: true
      path: /metrics
      scrape_interval: 15s
```

```python
from prometheus_client import Counter, Histogram
from agent_as_code import Agent

class MonitoredAgent(Agent):
    def __init__(self):
        self.requests = Counter('requests_total', 'Total requests')
        self.latency = Histogram('request_latency_seconds', 'Request latency')
```

### 2. Logging

```yaml
spec:
  logging:
    level: INFO
    format: json
    destination: stdout
```

```python
import logging
from agent_as_code import Agent

class LoggedAgent(Agent):
    def __init__(self):
        logging.basicConfig(
            format='%(asctime)s %(levelname)s %(message)s',
            level=logging.INFO
        )
```

### 3. Tracing

```yaml
spec:
  tracing:
    enabled: true
    provider: jaeger
    sampling_rate: 0.1
```

```python
from opentelemetry import trace
from agent_as_code import Agent

class TracedAgent(Agent):
    async def operation(self, request: dict):
        with trace.get_tracer(__name__).start_as_current_span("operation") as span:
            span.set_attribute("request.id", request["id"])
            return await self.process(request)
```

### 4. Alerting

```yaml
spec:
  alerts:
    - name: HighLatency
      condition: p99_latency > 1s
      for: 5m
      severity: warning
    - name: HighErrorRate
      condition: error_rate > 0.1
      for: 1m
      severity: critical
```

## Deployment

### 1. Rolling Updates

```yaml
spec:
  deployment:
    strategy: RollingUpdate
    maxUnavailable: 25%
    maxSurge: 25%
```

### 2. Health Checks

```yaml
spec:
  healthcheck:
    liveness:
      path: /health
      interval: 30s
      timeout: 10s
    readiness:
      path: /ready
      interval: 10s
```

### 3. Backup and Recovery

```yaml
spec:
  backup:
    schedule: "0 2 * * *"  # Daily at 2 AM
    retention: 7d
    destination: s3://backups
```

### 4. Scaling

```yaml
spec:
  scaling:
    min: 2
    max: 10
    metrics:
      - type: Resource
        resource:
          name: cpu
          target:
            type: Utilization
            averageUtilization: 70
```

## Development

### 1. Version Control

```yaml
metadata:
  name: my-agent
  version: 1.2.3
  labels:
    git-commit: ${GIT_COMMIT}
    build-time: ${BUILD_TIME}
```

### 2. Testing

```python
from agent_as_code import Agent, testing

class TestableAgent(Agent):
    async def operation(self, input: str) -> str:
        # Business logic
        pass

    @testing.mock_model
    async def test_operation(self):
        result = await self.operation("test input")
        assert result == "expected output"
```

### 3. Documentation

```python
class DocumentedAgent(Agent):
    """
    Agent for processing text data.
    
    Attributes:
        model: The AI model used for processing
        cache: Response cache
    
    Environment Variables:
        API_KEY: OpenAI API key
        CACHE_URL: Redis connection URL
    """
```

### 4. Code Quality

```yaml
spec:
  quality:
    linting:
      - pylint
      - mypy
    coverage:
      minimum: 80%
    security:
      - bandit
      - safety
```

## Operations

### 1. Disaster Recovery

```yaml
spec:
  dr:
    backup:
      enabled: true
      schedule: daily
    failover:
      region: us-west-2
      automatic: true
```

### 2. Capacity Planning

```yaml
spec:
  capacity:
    planning:
      growth_rate: 10%
      buffer: 30%
    alerts:
      - threshold: 80%
        action: notify
```

### 3. Cost Management

```yaml
spec:
  cost:
    budget:
      monthly: 1000
    optimization:
      enabled: true
    alerts:
      - threshold: 80%
        notification: email
```

### 4. Incident Response

```yaml
spec:
  incidents:
    notification:
      - slack
      - email
    runbooks:
      location: git://runbooks
    escalation:
      - team: oncall
        delay: 5m
      - team: manager
        delay: 15m
```

## Compliance

### 1. Audit Logging

```yaml
spec:
  audit:
    enabled: true
    retention: 90d
    events:
      - category: security
        level: INFO
      - category: data
        level: WARNING
```

### 2. Data Privacy

```yaml
spec:
  privacy:
    pii:
      detection: true
      redaction: true
    retention:
      policy: 30d
```

### 3. Compliance Reporting

```yaml
spec:
  compliance:
    frameworks:
      - soc2
      - gdpr
    reporting:
      schedule: monthly
      format: pdf
```

### 4. Access Reviews

```yaml
spec:
  access:
    review:
      schedule: quarterly
      approvers:
        - security-team
        - compliance-team
```

## Next Steps

- Review the [Security Guide](./security.md) for detailed security configurations
- Explore [Monitoring Setup](./monitoring.md) for advanced monitoring
- Learn about [Scaling Strategies](./scaling.md)
- Join our [Community Forum](https://github.com/pxkundu/agent-as-code/discussions) for more tips
