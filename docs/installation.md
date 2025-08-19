# Installation Guide

This guide covers the different methods to install Agent-as-Code.

## Prerequisites
- Python 3.8 or later (for Python package installation)
- Docker (optional, for container-based deployment)

## Installation Methods

### 1. Python Package (Recommended for Developers)
```bash
pip install agent-as-code
```

### 2. Direct Binary Download (Fastest)
For Linux/macOS:
```bash
curl -L https://api.myagentregistry.com/install.sh | sh
```

For Windows (PowerShell):
```powershell
iwr -useb https://api.myagentregistry.com/install.ps1 | iex
```

### 3. Homebrew (macOS/Linux)
```bash
brew install agent-as-code
```

### 4. Build from Source
```bash
git clone https://github.com/pxkundu/agent-as-code
cd agent-as-code
make install
```

## Verification
After installation, verify by running:
```bash
agent --version
```

If you encounter issues, visit our [Community Forum](https://github.com/pxkundu/agent-as-code/discussions).

## Next Steps
- Proceed to the [Quick Start Guide](./quickstart.md) to create your first agent.
