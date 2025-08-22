# 🗑️ Uninstalling Agent-as-Code

This guide covers all methods to completely remove Agent-as-Code from your system.

## ⚠️ Before You Begin

- **Backup any important data** from your agent projects
- **Note your configuration** if you plan to reinstall later
- **Ensure you have admin/sudo access** for system-wide removal

## 🔧 Uninstallation Methods

### Method 1: Remove Binary Installation

#### macOS / Linux
```bash
# Remove the binary from system PATH
sudo rm /usr/local/bin/agent

# Verify removal
which agent
# Should return: agent not found
```

#### Windows
```powershell
# Remove the binary (run as Administrator)
Remove-Item "C:\Windows\System32\agent.exe"

# Verify removal
Get-Command agent
# Should return: Get-Command : The term 'agent' is not recognized
```

### Method 2: Remove Python Package

If you installed via pip:

```bash
# Uninstall the Python package
pip uninstall agent-as-code -y

# Verify removal
pip list | grep agent-as-code
# Should return nothing
```

### Method 3: Remove Docker Image

If you used Docker:

```bash
# Remove the Docker image
docker rmi pxkundu/agent-as-code:latest

# Verify removal
docker images | grep agent-as-code
# Should return nothing
```

## 🧹 Cleanup Configuration Files

### Remove User Configuration
```bash
# Remove user config file
rm ~/.agent-as-code.yaml

# Verify removal
ls -la ~/.agent-as-code.yaml
# Should return: No such file or directory
```

### Remove Global Configuration
```bash
# Remove global config (if exists)
sudo rm /etc/agent-as-code.yaml

# Verify removal
ls -la /etc/agent-as-code.yaml
# Should return: No such file or directory
```

## 🗂️ Cleanup Project Files

### Remove Agent Projects
```bash
# List your agent projects
find ~ -name "agent.yaml" -type f 2>/dev/null

# Remove specific projects (replace with your project names)
rm -rf ~/my-chatbot
rm -rf ~/sentiment-analyzer
rm -rf ~/my-agent
```

### Remove Build Artifacts
```bash
# Remove any Docker images created by agent-as-code
docker images | grep -E "(my-agent|test-agent|chatbot)" | awk '{print $3}' | xargs -r docker rmi

# Remove any containers
docker ps -a | grep -E "(my-agent|test-agent|chatbot)" | awk '{print $1}' | xargs -r docker rm
```

## 🔍 Verification Steps

Run these commands to ensure complete removal:

```bash
# 1. Check binary in PATH
which agent
# Expected: agent not found

# 2. Check system binary
ls -la /usr/local/bin/agent 2>/dev/null || echo "✅ System binary removed"

# 3. Check config files
ls -la ~/.agent-as-code.yaml 2>/dev/null || echo "✅ Config file removed"

# 4. Check Python package
pip list | grep agent-as-code 2>/dev/null || echo "✅ Python package uninstalled"

# 5. Check Docker image
docker images | grep agent-as-code 2>/dev/null || echo "✅ Docker image removed"
```

## 🚨 Troubleshooting

### Binary Still Found
```bash
# Check all possible locations
find /usr/local/bin /usr/bin /opt -name "agent" 2>/dev/null

# Check your shell's PATH
echo $PATH | tr ':' '\n' | xargs -I {} find {} -name "agent" 2>/dev/null

# Remove from all found locations
sudo rm /path/to/found/agent
```

### Permission Denied
```bash
# Check file permissions
ls -la /usr/local/bin/agent

# Use sudo for removal
sudo rm /usr/local/bin/agent
```

### Configuration File Not Found
```bash
# Check for hidden files
ls -la ~/.agent-as-code*

# Check for alternative locations
find ~ -name "*agent*" -type f 2>/dev/null
```

### Python Package Issues
```bash
# Force uninstall
pip uninstall agent-as-code -y --force

# Check for multiple installations
pip list | grep -i agent

# Remove from specific Python environments
python3 -m pip uninstall agent-as-code -y
```

## 🔄 Complete System Reset

If you want to completely reset your system:

```bash
# Remove all agent-related files
sudo find /usr/local -name "*agent*" -delete 2>/dev/null
sudo find /usr/bin -name "*agent*" -delete 2>/dev/null

# Clean user directory
rm -rf ~/.agent-as-code*

# Clean Docker completely
docker system prune -a --volumes

# Clean pip cache
pip cache purge
```

## 📋 Uninstallation Checklist

- [ ] **Binary removed** from system PATH
- [ ] **Python package** uninstalled (if applicable)
- [ ] **Docker image** removed (if applicable)
- [ ] **Configuration files** deleted
- [ ] **Project directories** cleaned up
- [ ] **Build artifacts** removed
- [ ] **Verification** completed successfully

## 🌟 Reinstallation

If you want to reinstall later:

1. **Follow the [Installation Guide](./INSTALL.md)**
2. **Download the latest binary** from GitHub releases
3. **Set up fresh configuration** as needed

## 🆘 Need Help?

- **GitHub Issues**: [Report uninstallation problems](https://github.com/pxkundu/agent-as-code/issues)
- **Documentation**: Check other guides in the `/docs` folder
- **Community**: Ask questions in GitHub discussions

## ⚡ Quick Uninstall Script

For advanced users, here's a one-liner to remove everything:

```bash
# macOS/Linux (use with caution!)
sudo rm -f /usr/local/bin/agent && \
rm -f ~/.agent-as-code.yaml && \
pip uninstall agent-as-code -y 2>/dev/null && \
docker rmi pxkundu/agent-as-code:latest 2>/dev/null && \
echo "✅ Agent-as-Code completely removed!"
```

---

**Agent-as-Code has been successfully removed from your system! 🎉**
