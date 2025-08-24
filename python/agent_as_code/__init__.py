"""
Agent as Code (AaC) - Python Package
====================================

A Docker-like CLI for AI agents with Python ecosystem integration.

This package provides a Python wrapper around the Go-based Agent as Code CLI,
enabling seamless integration with Python development workflows while
maintaining the performance benefits of the Go implementation.

Usage:
    # Command line usage
    $ pip install agent-as-code
    $ agent init my-chatbot --template chatbot
    $ agent build -t my-chatbot:latest .
    $ agent run my-chatbot:latest

    # Enhanced LLM Commands
    $ agent llm create-agent chatbot
    $ agent llm optimize llama2 chatbot
    $ agent llm benchmark
    $ agent llm deploy-agent my-agent
    $ agent llm analyze llama2

    # Python API usage
    from agent_as_code import AgentCLI
    
    cli = AgentCLI()
    cli.execute(['init', 'my-agent', '--template', 'chatbot'])
    cli.execute(['build', '-t', 'my-agent:latest', '.'])
    cli.execute(['run', 'my-agent:latest'])
    
    # Enhanced LLM API usage
    cli.create_agent('sentiment-analyzer')
    cli.optimize_model('llama2', 'chatbot')
    cli.benchmark_models(['chatbot', 'code-generation'])
    cli.deploy_agent('my-agent', test_suite='comprehensive')
    cli.analyze_model('mistral:7b', detailed=True)

Features:
    - Docker-like CLI commands for AI agents
    - Enhanced LLM intelligence for agent creation and optimization
    - Multiple agent templates (chatbot, sentiment, etc.)
    - Local and cloud LLM support via Ollama
    - Python ecosystem integration
    - Cross-platform compatibility
    - Intelligent agent generation with automated testing
    - Model optimization and benchmarking
    - Comprehensive deployment and validation
"""

__version__ = "1.1.0"
__author__ = "Partha Sarathi Kundu"
__email__ = "inboxpartha@outlook.com"

from .cli import AgentCLI, main

__all__ = [
    "AgentCLI",
    "main",
    "__version__",
    "__author__",
    "__email__",
]

# Version compatibility check
def check_binary_version():
    """Check if the Go binary version matches the Python package version."""
    try:
        cli = AgentCLI()
        binary_version = cli.get_version()
        if binary_version and binary_version != __version__:
            import warnings
            warnings.warn(
                f"Version mismatch: Python package v{__version__} "
                f"but Go binary v{binary_version}. "
                f"Consider updating: pip install --upgrade agent-as-code",
                UserWarning
            )
    except Exception:
        # Ignore errors during version check
        pass

# Perform version check on import (optional)
# check_binary_version()