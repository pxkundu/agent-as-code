#!/usr/bin/env python3
"""
AI Agent - Chatbot Implementation
Supports both OpenAI and Ollama (local) models
"""

import os
import json
import logging
from typing import Dict, Any, List, Optional
from datetime import datetime

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import requests

# Configure logging
logging.basicConfig(level=os.getenv("LOG_LEVEL", "INFO"))
logger = logging.getLogger(__name__)

# FastAPI app
app = FastAPI(title="AI Agent Chatbot", version="1.0.0")

# Configuration
MODEL_PROVIDER = os.getenv("MODEL_PROVIDER", "openai")
MODEL_NAME = os.getenv("MODEL_NAME", "gpt-4")
OLLAMA_BASE_URL = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")

# Request/Response models
class ChatRequest(BaseModel):
    message: str
    conversation_id: Optional[str] = None

class ChatResponse(BaseModel):
    response: str
    conversation_id: str
    provider: str
    model: str
    timestamp: str

class HealthResponse(BaseModel):
    status: str
    provider: str
    model: str
    timestamp: str

# LLM Client classes
class LLMClient:
    """Base LLM client interface"""
    
    def chat(self, message: str, conversation_id: str = None) -> str:
        raise NotImplementedError

class OpenAIClient(LLMClient):
    """OpenAI API client"""
    
    def __init__(self, api_key: str, model: str):
        self.api_key = api_key
        self.model = model
        self.base_url = "https://api.openai.com/v1"
        
        if not self.api_key:
            raise ValueError("OPENAI_API_KEY environment variable is required for OpenAI models")
    
    def chat(self, message: str, conversation_id: str = None) -> str:
        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }
        
        payload = {
            "model": self.model,
            "messages": [
                {"role": "user", "content": message}
            ],
            "max_tokens": 500,
            "temperature": 0.7
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/chat/completions",
                headers=headers,
                json=payload,
                timeout=30
            )
            response.raise_for_status()
            
            result = response.json()
            return result["choices"][0]["message"]["content"]
            
        except Exception as e:
            logger.error(f"OpenAI API error: {e}")
            raise HTTPException(status_code=500, detail=f"OpenAI API error: {str(e)}")

class OllamaClient(LLMClient):
    """Ollama local LLM client"""
    
    def __init__(self, base_url: str, model: str):
        self.base_url = base_url.rstrip('/')
        self.model = model
    
    def chat(self, message: str, conversation_id: str = None) -> str:
        payload = {
            "model": self.model,
            "prompt": message,
            "stream": False,
            "options": {
                "temperature": 0.7,
                "num_predict": 500
            }
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/api/generate",
                json=payload,
                timeout=60  # Longer timeout for local models
            )
            response.raise_for_status()
            
            result = response.json()
            return result.get("response", "No response generated")
            
        except requests.exceptions.ConnectionError:
            logger.error("Cannot connect to Ollama. Make sure Ollama is running: ollama serve")
            raise HTTPException(
                status_code=503, 
                detail="Ollama is not running. Please start Ollama with 'ollama serve'"
            )
        except Exception as e:
            logger.error(f"Ollama API error: {e}")
            raise HTTPException(status_code=500, detail=f"Ollama API error: {str(e)}")

# Initialize LLM client based on provider
def create_llm_client() -> LLMClient:
    """Create appropriate LLM client based on configuration"""
    provider = MODEL_PROVIDER.lower()
    
    if provider == "ollama":
        logger.info(f"Initializing Ollama client with model: {MODEL_NAME}")
        return OllamaClient(OLLAMA_BASE_URL, MODEL_NAME)
    
    elif provider == "openai":
        logger.info(f"Initializing OpenAI client with model: {MODEL_NAME}")
        return OpenAIClient(OPENAI_API_KEY, MODEL_NAME)
    
    else:
        raise ValueError(f"Unsupported provider: {provider}. Use 'openai' or 'ollama'")

# Global LLM client
llm_client = create_llm_client()

# API Endpoints
@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "AI Agent Chatbot is running",
        "provider": MODEL_PROVIDER,
        "model": MODEL_NAME,
        "timestamp": datetime.now().isoformat()
    }

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        provider=MODEL_PROVIDER,
        model=MODEL_NAME,
        timestamp=datetime.now().isoformat()
    )

@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """Chat with the AI agent"""
    try:
        # Generate conversation ID if not provided
        conversation_id = request.conversation_id or f"conv_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        
        logger.info(f"Processing chat request: {request.message[:50]}...")
        
        # Get response from LLM
        response = llm_client.chat(request.message, conversation_id)
        
        logger.info(f"Generated response: {response[:50]}...")
        
        return ChatResponse(
            response=response,
            conversation_id=conversation_id,
            provider=MODEL_PROVIDER,
            model=MODEL_NAME,
            timestamp=datetime.now().isoformat()
        )
        
    except Exception as e:
        logger.error(f"Chat error: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/models")
async def list_models():
    """List available models"""
    if MODEL_PROVIDER.lower() == "ollama":
        try:
            response = requests.get(f"{OLLAMA_BASE_URL}/api/tags", timeout=10)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            raise HTTPException(status_code=503, detail=f"Cannot connect to Ollama: {e}")
    
    return {"provider": MODEL_PROVIDER, "current_model": MODEL_NAME}

# Startup event
@app.on_event("startup")
async def startup_event():
    """Application startup"""
    logger.info("🚀 Starting AI Agent Chatbot")
    logger.info(f"📡 Provider: {MODEL_PROVIDER}")
    logger.info(f"🤖 Model: {MODEL_NAME}")
    
    if MODEL_PROVIDER.lower() == "ollama":
        logger.info(f"🔗 Ollama URL: {OLLAMA_BASE_URL}")
        logger.info("💡 Make sure Ollama is running: ollama serve")
    
    logger.info("✅ Chatbot is ready!")

if __name__ == "__main__":
    import uvicorn
    
    port = int(os.getenv("PORT", 8080))
    host = os.getenv("HOST", "0.0.0.0")
    
    logger.info(f"🌐 Starting server on {host}:{port}")
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=False,
        log_level="info"
    )
