# Simple Chatbot Agent Example

This is a simple chatbot agent that demonstrates how to use the Agent-as-Code framework to create an AI-powered chatbot service.

## Features

- REST API using FastAPI
- OpenAI GPT integration
- Environment variable configuration
- Health check endpoint
- Docker containerization

## Setup

1. Create a `.env` file with your OpenAI API key:
   ```
   OPENAI_API_KEY=your_api_key_here
   MODEL_NAME=gpt-3.5-turbo  # optional, defaults to gpt-3.5-turbo
   ```

2. Build the agent:
   ```bash
   agent build .
   ```

3. Run the agent:
   ```bash
   agent run chatbot:latest
   ```

## API Usage

### Health Check
```bash
curl http://localhost:8080/health
```

### Chat
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello! How are you?"}
    ]
  }'
```

## Configuration

The agent can be configured through the following environment variables:
- `OPENAI_API_KEY`: Your OpenAI API key (required)
- `MODEL_NAME`: The OpenAI model to use (optional, defaults to gpt-3.5-turbo)

## Development

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Run locally:
   ```bash
   python main.py
   ```

3. Access the API documentation:
   ```
   http://localhost:8080/docs
   ```
