#!/usr/bin/env python3
"""
Enhanced LLM Commands Demo
==========================

This script demonstrates the new enhanced LLM commands available in
Agent as Code Python package v1.1.0+.

Features demonstrated:
- Intelligent agent creation
- Model optimization
- Benchmarking
- Agent deployment
- Model analysis
"""

from agent_as_code import AgentCLI
import time

def main():
    """Demonstrate enhanced LLM commands."""
    print("🧠 Enhanced LLM Commands Demo")
    print("=" * 50)
    
    try:
        # Initialize the CLI
        print("🔧 Initializing Agent CLI...")
        cli = AgentCLI()
        print("✅ CLI initialized successfully")
        
        # List available models
        print("\n📋 Listing available models...")
        if cli.list_models():
            print("✅ Models listed successfully")
        else:
            print("⚠️  No models available or Ollama not running")
        
        # Create an intelligent agent
        print("\n🚀 Creating intelligent chatbot agent...")
        if cli.create_agent('chatbot'):
            print("✅ Chatbot agent created successfully")
            print("📁 Project directory: chatbot-agent")
        else:
            print("❌ Failed to create chatbot agent")
            return
        
        # Create another agent for demonstration
        print("\n🚀 Creating intelligent sentiment analyzer agent...")
        if cli.create_agent('sentiment-analyzer'):
            print("✅ Sentiment analyzer agent created successfully")
            print("📁 Project directory: sentiment-analyzer-agent")
        else:
            print("❌ Failed to create sentiment analyzer agent")
            return
        
        # Deploy and test the first agent
        print("\n🚀 Deploying and testing chatbot agent...")
        if cli.deploy_agent('chatbot-agent', test_suite='comprehensive'):
            print("✅ Chatbot agent deployed and tested successfully")
            print("🔗 Access: http://localhost:8080")
        else:
            print("❌ Failed to deploy chatbot agent")
        
        # Demonstrate model optimization (if models are available)
        print("\n⚡ Demonstrating model optimization...")
        print("Note: This requires models to be available via Ollama")
        if cli.optimize_model('llama2', 'chatbot'):
            print("✅ Model optimization completed")
        else:
            print("⚠️  Model optimization failed (model may not be available)")
        
        # Demonstrate benchmarking
        print("\n📊 Running model benchmarks...")
        if cli.benchmark_models(['chatbot', 'sentiment-analysis']):
            print("✅ Benchmarking completed")
        else:
            print("⚠️  Benchmarking failed (models may not be available)")
        
        # Demonstrate model analysis
        print("\n🔍 Analyzing model capabilities...")
        if cli.analyze_model('llama2', detailed=True):
            print("✅ Model analysis completed")
        else:
            print("⚠️  Model analysis failed (model may not be available)")
        
        print("\n🎉 Enhanced LLM Commands Demo Completed!")
        print("\n💡 Next steps:")
        print("   - Explore the generated agent projects")
        print("   - Customize the agents for your specific needs")
        print("   - Deploy to production environments")
        print("   - Integrate with your existing workflows")
        
    except Exception as e:
        print(f"❌ Error during demo: {e}")
        print("\n💡 Troubleshooting tips:")
        print("   - Ensure Ollama is running: ollama serve")
        print("   - Pull required models: ollama pull llama2")
        print("   - Check network connectivity")
        print("   - Verify Go binary is available")

if __name__ == "__main__":
    main()
