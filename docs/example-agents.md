# Example Agents

This guide provides a collection of example agents demonstrating different use cases and capabilities of the Agent-as-Code framework.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Advanced Examples](#advanced-examples)
- [Integration Examples](#integration-examples)
- [Production Examples](#production-examples)

## Basic Examples

### 1. Simple Chatbot

A basic chatbot that can engage in conversations using GPT-4.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: simple-chatbot
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

```python
# main.py
from agent_as_code import Agent

class ChatbotAgent(Agent):
    async def chat(self, message: str) -> str:
        response = await self.model.generate(message)
        return response

if __name__ == "__main__":
    agent = ChatbotAgent()
    agent.run()
```

### 2. Sentiment Analyzer

An agent that analyzes text sentiment using AI.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: sentiment-analyzer
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  capabilities:
    - sentiment-analysis
```

```python
# main.py
from agent_as_code import Agent
from typing import Dict

class SentimentAgent(Agent):
    async def analyze_sentiment(self, text: str) -> Dict[str, float]:
        prompt = f"Analyze the sentiment of this text: {text}"
        response = await self.model.generate(prompt)
        
        # Parse sentiment scores
        return {
            "positive": 0.8,  # Example scores
            "negative": 0.1,
            "neutral": 0.1
        }

if __name__ == "__main__":
    agent = SentimentAgent()
    agent.run()
```

## Advanced Examples

### 1. Multi-Model Agent

An agent that uses multiple AI models for different tasks.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: multi-model-agent
  version: 1.0.0
spec:
  runtime: python:3.11
  models:
    - name: gpt-4
      provider: openai
      use: text-generation
    - name: llama2
      provider: ollama
      use: classification
    - name: claude-3
      provider: anthropic
      use: analysis
```

```python
# main.py
from agent_as_code import Agent, Model

class MultiModelAgent(Agent):
    async def process(self, task: str, data: str) -> str:
        if task == "generate":
            model = self.get_model("gpt-4")
            return await model.generate(data)
        elif task == "classify":
            model = self.get_model("llama2")
            return await model.classify(data)
        elif task == "analyze":
            model = self.get_model("claude-3")
            return await model.analyze(data)
```

### 2. Data Processing Agent

An agent that processes and analyzes data using AI.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: data-processor
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  dependencies:
    - pandas
    - numpy
    - matplotlib
```

```python
# main.py
import pandas as pd
import numpy as np
from agent_as_code import Agent

class DataProcessorAgent(Agent):
    async def process_data(self, data: pd.DataFrame) -> dict:
        # Analyze data
        stats = data.describe()
        
        # Generate insights using AI
        insights = await self.model.generate(
            f"Analyze this data: {stats.to_string()}"
        )
        
        return {
            "stats": stats.to_dict(),
            "insights": insights,
            "recommendations": await self.generate_recommendations(data)
        }
```

## Integration Examples

### 1. Slack Integration

An agent that integrates with Slack for team communication.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: slack-bot
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  environment:
    - name: SLACK_TOKEN
      from: secret
```

```python
# main.py
from slack_bolt import App
from agent_as_code import Agent

class SlackAgent(Agent):
    def __init__(self):
        super().__init__()
        self.app = App(token=self.env.SLACK_TOKEN)
        
    async def handle_message(self, message: str, channel: str):
        response = await self.model.generate(message)
        await self.app.client.chat_postMessage(
            channel=channel,
            text=response
        )
```

### 2. Database Integration

An agent that works with databases and AI.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: db-agent
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  environment:
    - name: DATABASE_URL
      from: secret
```

```python
# main.py
from sqlalchemy import create_engine
from agent_as_code import Agent

class DatabaseAgent(Agent):
    def __init__(self):
        super().__init__()
        self.engine = create_engine(self.env.DATABASE_URL)
        
    async def query_and_analyze(self, query: str) -> dict:
        # Execute query
        results = self.engine.execute(query)
        
        # Analyze results with AI
        analysis = await self.model.generate(
            f"Analyze these results: {results}"
        )
        
        return {
            "results": results,
            "analysis": analysis
        }
```

## Production Examples

### 1. High-Availability Agent

An agent designed for high availability and scalability.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: ha-agent
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  replicas: 3
  resources:
    requests:
      memory: "512Mi"
      cpu: "250m"
    limits:
      memory: "1Gi"
      cpu: "500m"
  healthcheck:
    path: /health
    interval: 30s
```

### 2. Monitoring-Enabled Agent

An agent with comprehensive monitoring and logging.

```yaml
# agent.yaml
apiVersion: agent.dev/v1
kind: Agent
metadata:
  name: monitored-agent
  version: 1.0.0
spec:
  runtime: python:3.11
  model:
    provider: openai
    name: gpt-4
  monitoring:
    prometheus: true
    logging:
      level: INFO
      format: json
    tracing:
      enabled: true
      provider: jaeger
```

```python
# main.py
from prometheus_client import Counter, Histogram
from agent_as_code import Agent
import logging

class MonitoredAgent(Agent):
    def __init__(self):
        super().__init__()
        self.requests = Counter('requests_total', 'Total requests')
        self.latency = Histogram('request_latency_seconds', 'Request latency')
        
    async def process(self, request: dict) -> dict:
        self.requests.inc()
        with self.latency.time():
            response = await self.model.generate(request["prompt"])
            logging.info("Request processed", extra={
                "request_id": request["id"],
                "latency": self.latency
            })
            return {"response": response}
```

## Running the Examples

1. **Clone the Example**
   ```bash
   agent init my-agent --template example-name
   cd my-agent
   ```

2. **Configure Environment**
   ```bash
   # Create .env file
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Build and Run**
   ```bash
   agent build -t my-agent:latest .
   agent run my-agent:latest
   ```

## Best Practices

1. **Error Handling**
   - Implement proper error handling
   - Use try-except blocks
   - Log errors appropriately

2. **Security**
   - Never hardcode secrets
   - Use environment variables
   - Implement proper authentication

3. **Performance**
   - Cache responses when possible
   - Implement rate limiting
   - Monitor resource usage

4. **Testing**
   - Write unit tests
   - Implement integration tests
   - Use mock services for testing

## Next Steps

- Explore the [API Reference](./api-reference.md)
- Learn about [deployment strategies](./deployment-strategies.md)
- Read the [security guide](./security.md)
- Join our [community](https://github.com/pxkundu/agent-as-code/discussions)
