#!/usr/bin/env python3
"""
Sentiment Analysis Agent - Analyzes text sentiment using AI
"""

import os
import logging
from datetime import datetime
from typing import Dict, List, Optional

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import openai

# Configure logging
logging.basicConfig(level=os.getenv("LOG_LEVEL", "INFO"))
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Sentiment Analysis Agent",
    description="AI-powered sentiment analysis for text content",
    version="1.0.0"
)

# Request/Response models
class SentimentRequest(BaseModel):
    text: str
    include_confidence: Optional[bool] = True

class SentimentResponse(BaseModel):
    sentiment: str  # positive, negative, neutral
    confidence: Optional[float] = None
    timestamp: str

class HealthResponse(BaseModel):
    status: str
    uptime: str
    timestamp: str

class SentimentAgent:
    def __init__(self):
        self.client = openai.OpenAI(
            api_key=os.getenv("OPENAI_API_KEY")
        )
        
    async def analyze_sentiment(self, request: SentimentRequest) -> SentimentResponse:
        """Analyze sentiment of the provided text"""
        try:
            # Prepare prompt for sentiment analysis
            prompt = f"""Analyze the sentiment of the following text and respond with only one word: positive, negative, or neutral.

Text: {request.text}

Sentiment:"""

            # Generate sentiment using OpenAI
            response = self.client.chat.completions.create(
                model="gpt-4",
                messages=[
                    {"role": "system", "content": "You are a sentiment analysis expert. Respond with only: positive, negative, or neutral."},
                    {"role": "user", "content": prompt}
                ],
                max_tokens=10,
                temperature=0.1
            )
            
            sentiment = response.choices[0].message.content.strip().lower()
            
            # Validate sentiment
            if sentiment not in ["positive", "negative", "neutral"]:
                sentiment = "neutral"
            
            # Calculate confidence (simplified approach)
            confidence = None
            if request.include_confidence:
                confidence = 0.85  # Simplified confidence score
            
            return SentimentResponse(
                sentiment=sentiment,
                confidence=confidence,
                timestamp=datetime.now().isoformat()
            )
            
        except Exception as e:
            logger.error(f"Error analyzing sentiment: {e}")
            raise HTTPException(status_code=500, detail="Internal server error")

# Initialize sentiment agent
sentiment_agent = SentimentAgent()

@app.post("/analyze", response_model=SentimentResponse)
async def analyze_sentiment(request: SentimentRequest):
    """Analyze sentiment of provided text"""
    return await sentiment_agent.analyze_sentiment(request)

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
    return {"message": "Sentiment Analysis Agent API", "status": "running", "version": "1.0.0"}

if __name__ == "__main__":
    import uvicorn
    import time
    
    start_time = time.time()
    
    logger.info("Starting Sentiment Analysis Agent...")
    uvicorn.run(
        app, 
        host="0.0.0.0", 
        port=8080,
        log_level=os.getenv("LOG_LEVEL", "info").lower()
    )
