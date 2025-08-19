#!/usr/bin/env python3
"""
Basic Chatbot Agent - Real-world Example
A production-ready customer support chatbot with conversation memory
"""

import os
import json
import logging
from typing import Dict, Any, List, Optional
from datetime import datetime

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import openai

# Configure logging
logging.basicConfig(level=os.getenv("LOG_LEVEL", "INFO"))
logger = logging.getLogger(__name__)

# Configuration
openai.api_key = os.getenv("OPENAI_API_KEY")
MODEL_NAME = os.getenv("MODEL_NAME", "gpt-4")
MAX_CONVERSATION_HISTORY = int(os.getenv("MAX_CONVERSATION_HISTORY", "10"))
ESCALATION_KEYWORDS = [kw.strip() for kw in os.getenv("ESCALATION_KEYWORDS", "human,manager,supervisor,escalate").split(',')]

# FastAPI app
app = FastAPI(
    title="Basic Chatbot Agent",
    description="A customer support chatbot agent with conversation memory and escalation handling.",
    version="1.0.0",
)

# Data models
class ChatMessage(BaseModel):
    role: str
    content: str

class ChatRequest(BaseModel):
    message: str
    conversation_id: Optional[str] = None
    customer_id: Optional[str] = None
    context: Optional[Dict[str, Any]] = None

class ChatResponse(BaseModel):
    response: str
    conversation_id: str
    escalation_triggered: bool = False
    model_used: str
    timestamp: str

# In-memory conversation history (replace with Redis/DB in production)
conversation_history: Dict[str, List[ChatMessage]] = {}
customer_context: Dict[str, Dict[str, Any]] = {}

# System prompt for the chatbot
SYSTEM_PROMPT = """You are a helpful customer support agent for a technology company. 
You should be friendly, professional, and helpful. If you cannot help with something or 
the customer explicitly asks for a human agent, politely acknowledge the request.

Guidelines:
- Be concise but thorough
- Ask clarifying questions when needed
- Provide step-by-step instructions when appropriate
- Acknowledge when escalation to human agents is needed
- Remember context from previous messages in the conversation
"""

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "ok", "timestamp": datetime.now().isoformat()}

@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """Main chat endpoint"""
    try:
        # Generate conversation ID if not provided
        conv_id = request.conversation_id or f"conv_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        
        # Store customer context
        if request.customer_id and request.context:
            customer_context[request.customer_id] = request.context
        
        # Retrieve conversation history
        history = conversation_history.get(conv_id, [])
        
        # Add system message if this is a new conversation
        if not history:
            history.append(ChatMessage(role="system", content=SYSTEM_PROMPT))
        
        # Add user message
        history.append(ChatMessage(role="user", content=request.message))
        
        # Trim history to max length (keep system message)
        if len(history) > MAX_CONVERSATION_HISTORY + 1:  # +1 for system message
            history = [history[0]] + history[-(MAX_CONVERSATION_HISTORY):]
        
        logger.info(f"Processing message for conversation {conv_id}: {request.message}")
        
        # Check for escalation keywords
        escalation_triggered = any(keyword.lower() in request.message.lower() for keyword in ESCALATION_KEYWORDS)
        
        if escalation_triggered:
            response_content = "I understand you'd like to speak with a human agent. I'm connecting you with our support team now. Please hold on while I transfer your conversation."
            logger.info(f"Escalation triggered for conversation {conv_id}")
        else:
            # Generate AI response
            messages = [{"role": msg.role, "content": msg.content} for msg in history]
            
            response = openai.chat.completions.create(
                model=MODEL_NAME,
                messages=messages,
                temperature=0.7,
                max_tokens=500
            )
            
            response_content = response.choices[0].message.content
            logger.info(f"AI response generated for conversation {conv_id}")
        
        # Add assistant response to history
        history.append(ChatMessage(role="assistant", content=response_content))
        conversation_history[conv_id] = history
        
        return ChatResponse(
            response=response_content,
            conversation_id=conv_id,
            escalation_triggered=escalation_triggered,
            model_used=MODEL_NAME,
            timestamp=datetime.now().isoformat()
        )
        
    except Exception as e:
        logger.error(f"Error processing chat request: {e}")
        raise HTTPException(status_code=500, detail=f"Error processing request: {str(e)}")

@app.get("/conversations/{conversation_id}")
async def get_conversation(conversation_id: str):
    """Get conversation history"""
    history = conversation_history.get(conversation_id, [])
    return {
        "conversation_id": conversation_id,
        "messages": [{"role": msg.role, "content": msg.content} for msg in history],
        "message_count": len(history)
    }

@app.delete("/conversations/{conversation_id}")
async def clear_conversation(conversation_id: str):
    """Clear conversation history"""
    if conversation_id in conversation_history:
        del conversation_history[conversation_id]
        logger.info(f"Cleared conversation {conversation_id}")
        return {"message": f"Conversation {conversation_id} cleared"}
    else:
        raise HTTPException(status_code=404, detail="Conversation not found")

@app.get("/stats")
async def get_stats():
    """Get chatbot statistics"""
    total_conversations = len(conversation_history)
    total_messages = sum(len(history) for history in conversation_history.values())
    
    return {
        "total_conversations": total_conversations,
        "total_messages": total_messages,
        "active_conversations": total_conversations,
        "model_used": MODEL_NAME,
        "uptime": datetime.now().isoformat()
    }

if __name__ == "__main__":
    import uvicorn
    
    logger.info(f"Starting Basic Chatbot Agent with model: {MODEL_NAME}")
    uvicorn.run(app, host="0.0.0.0", port=8080)
