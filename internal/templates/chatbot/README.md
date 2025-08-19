# Chatbot Agent

AI-powered customer support chatbot with conversation memory and escalation handling.

## Features

- **Conversation Memory**: Maintains conversation context across messages
- **Escalation Detection**: Automatically detects when users need human assistance
- **Session Management**: Tracks user sessions and conversation history
- **RESTful API**: Simple HTTP API for chat interactions

## Configuration

Set the following environment variables:

- `OPENAI_API_KEY`: Your OpenAI API key
- `LOG_LEVEL`: Logging level (default: INFO)
- `MAX_CONVERSATION_HISTORY`: Max messages to keep in memory (default: 10)
- `ESCALATION_KEYWORDS`: Comma-separated keywords that trigger escalation (default: human,manager,supervisor,escalate)

## Usage

### Start the agent
```bash
python main.py
```

### Send a message
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, I need help with my order", "session_id": "user123"}'
```

### Check health
```bash
curl http://localhost:8080/health
```

## API Endpoints

- `POST /chat` - Send a chat message
- `GET /health` - Health check
- `GET /` - Root endpoint with basic info
