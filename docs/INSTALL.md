# 🚀 Installing Agent-as-Code

This guide covers all methods to install the Agent-as-Code binary on your system.

## 📋 Prerequisites

- **macOS**: 10.15+ (Catalina) or later
- **Linux**: Ubuntu 18.04+, CentOS 7+, or equivalent
- **Windows**: Windows 10+ (64-bit)
- **Docker**: Required for building and running agents

## 🔧 Installation Methods

### Method 1: Direct Binary Download (Recommended)

#### macOS (Intel)
```bash
# Download the binary
curl -L -o agent https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-darwin-amd64

# Make it executable
chmod +x agent

# Move to system PATH
sudo mv agent /usr/local/bin/agent

# Verify installation
agent version
```

#### macOS (Apple Silicon / ARM64)
```bash
# Download the binary
curl -L -o agent https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-darwin-arm64

# Make it executable
chmod +x agent

# Move to system PATH
sudo mv agent /usr/local/bin/agent

# Verify installation
agent version
```

#### Linux (AMD64)
```bash
# Download the binary
curl -L -o agent https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-linux-amd64

# Make it executable
chmod +x agent

# Move to system PATH
sudo mv agent /usr/local/bin/agent

# Verify installation
agent version
```

#### Linux (ARM64)
```bash
# Download the binary
curl -L -o agent https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-linux-arm64

# Make it executable
chmod +x agent

# Move to system PATH
sudo mv agent /usr/local/bin/agent

# Verify installation
agent version
```

#### Windows (AMD64)
```powershell
# Download the binary
Invoke-WebRequest -Uri "https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-windows-amd64.exe" -OutFile "agent.exe"

# Move to system PATH (run as Administrator)
Move-Item "agent.exe" "C:\Windows\System32\agent.exe"

# Verify installation
agent version
```

#### Windows (ARM64)
```powershell
# Download the binary
Invoke-WebRequest -Uri "https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-windows-arm64.exe" -OutFile "agent.exe"

# Move to system PATH (run as Administrator)
Move-Item "agent.exe" "C:\Windows\System32\agent.exe"

# Verify installation
agent version
```

### Method 2: Using wget (Alternative)

```bash
# For Linux/macOS users who prefer wget
wget -O agent https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)

chmod +x agent
sudo mv agent /usr/local/bin/agent
```

### Method 3: Manual Platform Detection

```bash
# Detect your platform automatically
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64) ARCH="arm64" ;;
esac

# Download appropriate binary
curl -L -o agent "https://github.com/pxkundu/agent-as-code/releases/latest/download/agent-${OS}-${ARCH}"

chmod +x agent
sudo mv agent /usr/local/bin/agent
```

## 🐳 Docker Installation

If you prefer to use Docker instead of installing the binary:

```bash
# Pull the latest image
docker pull pxkundu/agent-as-code:latest

# Run commands using Docker
docker run --rm pxkundu/agent-as-code:latest version
docker run --rm pxkundu/agent-as-code:latest init my-agent
```

## 📦 Python Package Installation

For Python users, you can also install via pip:

```bash
# Install from PyPI
pip install agent-as-code

# Verify installation
agent version
```

## ✅ Verification

After installation, verify everything works:

```bash
# Check version and ASCII art
agent version

# Check help
agent --help

# Test initialization
agent init test-agent --template chatbot
```

## 🔧 Troubleshooting

### Permission Denied
```bash
# Make sure the binary is executable
chmod +x agent

# Check file permissions
ls -la agent
```

### Command Not Found
```bash
# Verify the binary is in PATH
which agent

# Check if /usr/local/bin is in your PATH
echo $PATH | grep /usr/local/bin

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

### Binary Not Working
```bash
# Check if it's the right architecture
file $(which agent)

# Verify it's not corrupted
agent version
```

### Docker Issues
```bash
# Check Docker is running
docker --version
docker ps

# Restart Docker if needed
sudo systemctl restart docker  # Linux
# Or restart Docker Desktop on macOS/Windows
```

## 🌟 Next Steps

After successful installation:

1. **Read the [Quickstart Guide](./quickstart.md)** to create your first agent
2. **Explore [Templates](./templates.md)** for pre-built agent configurations
3. **Check [CLI Overview](./cli-overview.md)** for all available commands
4. **Review [Examples](./examples.md)** for practical use cases

## 📚 Additional Resources

- [CLI Overview](./cli-overview.md) - Complete command reference
- [Quickstart Guide](./quickstart.md) - Get started in minutes
- [Architecture](./architecture.md) - Understand the system design
- [Examples](./examples.md) - Real-world usage examples

## 🆘 Need Help?

- **GitHub Issues**: [Report bugs or request features](https://github.com/pxkundu/agent-as-code/issues)
- **Documentation**: Browse the `/docs` folder for comprehensive guides
- **Examples**: Check the `/examples` folder for working code samples

---

**Happy coding with Agent-as-Code! 🚀**
