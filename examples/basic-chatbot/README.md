# Basic Chatbot Agent

A production-ready customer support chatbot with conversation memory and escalation handling.

## Features

- **Conversation Memory**: Maintains context across multiple messages
- **Escalation Handling**: Automatically detects when customers need human support
- **Customer Context**: Stores and uses customer information for personalized responses
- **Health Monitoring**: Built-in health checks and statistics
- **Production Ready**: Comprehensive error handling and logging

## Quick Start

### 1. Build and Run

```bash
# Using Agent-as-Code CLI
agent build -t basic-chatbot:latest .
agent run basic-chatbot:latest

# Or manually
pip install -r requirements.txt
export OPENAI_API_KEY="your-api-key-here"
python main.py
```

### 2. Test the Chatbot

```bash
# Start a conversation
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hi, I need help with my account",
    "customer_id": "customer_123",
    "context": {"plan": "premium", "status": "active"}
  }'

# Continue the conversation using the returned conversation_id
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "I forgot my password",
    "conversation_id": "conv_20241201_143022"
  }'

# Trigger escalation
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "I want to speak to a manager",
    "conversation_id": "conv_20241201_143022"
  }'
```

### 3. Monitor and Manage

```bash
# Check health
curl http://localhost:8080/health

# Get statistics
curl http://localhost:8080/stats

# View conversation history
curl http://localhost:8080/conversations/conv_20241201_143022

# Clear a conversation
curl -X DELETE http://localhost:8080/conversations/conv_20241201_143022
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API key (required) | - |
| `MODEL_NAME` | OpenAI model to use | `gpt-4` |
| `LOG_LEVEL` | Logging level | `INFO` |
| `MAX_CONVERSATION_HISTORY` | Max messages to keep in memory | `10` |
| `ESCALATION_KEYWORDS` | Comma-separated keywords that trigger escalation | `human,manager,supervisor,escalate` |

### Agent Configuration

The `agent.yaml` file contains the complete agent configuration including:

- **Model Settings**: OpenAI GPT-4 with customizable temperature and token limits
- **Dependencies**: All required Python packages
- **Environment**: Configuration for secrets and environment variables
- **Health Checks**: Automated health monitoring
- **Ports**: Container port mapping

## API Reference

### POST /chat

Send a message to the chatbot.

**Request:**
```json
{
  "message": "Hello, I need help",
  "conversation_id": "optional-conversation-id",
  "customer_id": "optional-customer-id",
  "context": {
    "plan": "premium",
    "status": "active"
  }
}
```

**Response:**
```json
{
  "response": "Hello! I'd be happy to help you. What can I assist you with today?",
  "conversation_id": "conv_20241201_143022",
  "escalation_triggered": false,
  "model_used": "gpt-4",
  "timestamp": "2024-12-01T14:30:22.123456"
}
```

### GET /conversations/{conversation_id}

Retrieve conversation history.

### DELETE /conversations/{conversation_id}

Clear conversation history.

### GET /health

Health check endpoint.

### GET /stats

Get chatbot usage statistics.

## Production Deployment

### Docker Deployment

```bash
# Build using Agent-as-Code
agent build -t basic-chatbot:latest .

# Run with environment variables
docker run -d \
  -p 8080:8080 \
  -e OPENAI_API_KEY="your-api-key" \
  -e LOG_LEVEL="INFO" \
  basic-chatbot:latest
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: basic-chatbot
spec:
  replicas: 3
  selector:
    matchLabels:
      app: basic-chatbot
  template:
    metadata:
      labels:
        app: basic-chatbot
    spec:
      containers:
      - name: chatbot
        image: basic-chatbot:latest
        ports:
        - containerPort: 8080
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: openai-secret
              key: api-key
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

## Production Considerations

### Scalability

- **Conversation Storage**: Replace in-memory storage with Redis or database
- **Load Balancing**: Use sticky sessions or external session storage
- **Caching**: Cache frequent responses to reduce API costs

### Security

- **API Key Management**: Use secure secret management (Kubernetes secrets, AWS Secrets Manager)
- **Rate Limiting**: Implement rate limiting to prevent abuse
- **Input Validation**: Add comprehensive input validation and sanitization

### Monitoring

- **Metrics**: Add Prometheus metrics for conversation counts, response times, escalations
- **Logging**: Structured logging with correlation IDs
- **Alerting**: Set up alerts for escalation spikes or error rates

### Cost Optimization

- **Model Selection**: Use cheaper models (GPT-3.5) for simple queries
- **Response Caching**: Cache common responses
- **Token Management**: Optimize prompts to reduce token usage

## Customization

### Custom Escalation Logic

```python
def should_escalate(message: str, context: Dict) -> bool:
    # Custom escalation logic
    if context.get("previous_escalations", 0) > 2:
        return True
    
    sentiment_score = analyze_sentiment(message)
    if sentiment_score < -0.5:  # Very negative
        return True
    
    return any(keyword in message.lower() for keyword in ESCALATION_KEYWORDS)
```

### Custom System Prompts

Modify the `SYSTEM_PROMPT` variable to customize the chatbot's behavior:

```python
SYSTEM_PROMPT = f"""You are {COMPANY_NAME}'s AI customer support assistant.
Company information: {COMPANY_INFO}
Available services: {SERVICES}

Guidelines:
- Always be helpful and professional
- Provide accurate information about our services
- Escalate complex technical issues to human agents
"""
```

## Troubleshooting

### Common Issues

1. **OpenAI API Key Not Working**
   - Verify the API key is correctly set
   - Check API key permissions and billing

2. **High Memory Usage**
   - Reduce `MAX_CONVERSATION_HISTORY`
   - Implement conversation cleanup

3. **Slow Response Times**
   - Monitor OpenAI API latency
   - Consider using faster models for simple queries

### Debugging

Enable debug logging:
```bash
export LOG_LEVEL=DEBUG
python main.py
```

## License

This example is part of the Agent-as-Code framework and is licensed under the MIT License.
