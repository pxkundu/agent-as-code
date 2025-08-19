#!/usr/bin/env python3
"""
Chatbot Agent - Customer Support Agent with Conversation Memory
"""

import os
import asyncio
import logging
from datetime import datetime
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import openai

# Configure logging
logging.basicConfig(level=os.getenv("LOG_LEVEL", "INFO"))
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Chatbot Agent",
    description="AI-powered customer support chatbot with conversation memory",
    version="1.0.0"
)

# Request/Response models
class ChatRequest(BaseModel):
    message: str
    user_id: Optional[str] = None
    session_id: Optional[str] = None

class ChatResponse(BaseModel):
    response: str
    session_id: str
    timestamp: str

class HealthResponse(BaseModel):
    status: str
    uptime: str
    timestamp: str

@dataclass
class ConversationHistory:
    messages: List[Dict[str, str]]
    created_at: datetime
    updated_at: datetime

class ChatbotAgent:
    def __init__(self):
        self.client = openai.OpenAI(
            api_key=os.getenv("OPENAI_API_KEY")
        )
        self.conversation_history: Dict[str, ConversationHistory] = {}
        self.max_history = int(os.getenv("MAX_CONVERSATION_HISTORY", "10"))
        self.escalation_keywords = os.getenv("ESCALATION_KEYWORDS", "human,manager,supervisor,escalate").split(",")
        
    async def process_message(self, request: ChatRequest) -> ChatResponse:
        """Process incoming chat message and generate response"""
        try:
            session_id = request.session_id or f"session_{datetime.now().timestamp()}"
            
            # Get or create conversation history
            if session_id not in self.conversation_history:
                self.conversation_history[session_id] = ConversationHistory(
                    messages=[],
                    created_at=datetime.now(),
                    updated_at=datetime.now()
                )
            
            history = self.conversation_history[session_id]
            
            # Check for escalation keywords
            if any(keyword.lower() in request.message.lower() for keyword in self.escalation_keywords):
                response_text = ("I understand you'd like to speak with a human representative. "
                               "I'm transferring you to our support team. Please hold while I connect you.")
                # In a real implementation, this would trigger escalation workflow
                logger.info(f"Escalation triggered for session {session_id}")
            else:
                # Prepare conversation context
                messages = [
                    {"role": "system", "content": "You are a helpful customer support assistant. Be friendly, professional, and helpful."}
                ]
                
                # Add conversation history
                for msg in history.messages[-self.max_history:]:
                    messages.append(msg)
                
                # Add current message
                messages.append({"role": "user", "content": request.message})
                
                # Generate response using OpenAI
                response = await asyncio.to_thread(
                    self.client.chat.completions.create,
                    model="gpt-4",
                    messages=messages,
                    max_tokens=500,
                    temperature=0.7
                )
                
                response_text = response.choices[0].message.content
            
            # Update conversation history
            history.messages.extend([
                {"role": "user", "content": request.message},
                {"role": "assistant", "content": response_text}
            ])
            history.updated_at = datetime.now()
            
            # Trim history if too long
            if len(history.messages) > self.max_history * 2:
                history.messages = history.messages[-self.max_history * 2:]
            
            return ChatResponse(
                response=response_text,
                session_id=session_id,
                timestamp=datetime.now().isoformat()
            )
            
        except Exception as e:
            logger.error(f"Error processing message: {e}")
            raise HTTPException(status_code=500, detail="Internal server error")

# Initialize chatbot
chatbot = ChatbotAgent()

@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """Chat endpoint for processing messages"""
    return await chatbot.process_message(request)

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    import time
    uptime = time.time() - start_time
    return HealthResponse(
        status="healthy",
        uptime=f"{uptime:.2f}s",
        timestamp=datetime.now().isoformat()
    )

@app.get("/")
async def root():
    """Root endpoint"""
    return {"message": "Chatbot Agent API", "status": "running", "version": "1.0.0"}

if __name__ == "__main__":
    import uvicorn
    import time
    
    start_time = time.time()
    
    logger.info("Starting Chatbot Agent...")
    uvicorn.run(
        app, 
        host="0.0.0.0", 
        port=8080,
        log_level=os.getenv("LOG_LEVEL", "info").lower()
    )
