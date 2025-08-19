# API Reference

This reference is aligned with the actual Go code under `internal/api`, `internal/runtime`, `internal/cmd`, and `internal/llm`. It covers:
- Binary API (version listing, file listing, download, upload)
- Profile configuration and PAT usage
- Runtime API (container lifecycle)
- Local LLM management (Ollama)

Import paths used in examples are based on this repository’s modules.

## Binary API (internal/api)

The Binary API client is implemented in `internal/api/client.go`, with helpers in `downloader.go` and `uploader.go`.

### Client

```go
import (
    "github.com/pxkundu/agent-as-code/internal/api"
)

// Create API client
client := api.NewClient("https://api.myagentregistry.com")

// Set authentication token for upload endpoints (required for UploadBinary)
client.SetAuthToken(os.Getenv("AGENT_REGISTRY_TOKEN"))
```

### Endpoints (used by the client)
- GET `{base}/binary/releases/agent-as-code/versions`
- GET `{base}/binary/releases/agent-as-code/{major}/{minor}/`
- GET `{base}/binary/releases/agent-as-code/{major}/{minor}/{filename}`
- POST `{base}/binary/releases/agent-as-code/{major}/{minor}/upload` (Authorization: `Bearer <token>`)

### Version Management

List available versions:
```go
versionsResp, err := client.ListVersions()
if err != nil { /* handle */ }
fmt.Println(versionsResp.Versions)
```

List files for a version (major.minor derived from semantic version string):
```go
filesResp, err := client.ListFiles(1, 2) // lists for 1.2.x
if err != nil { /* handle */ }
for _, f := range filesResp.Files {
    fmt.Printf("%s %s/%s %dB\n", f.Filename, f.Platform, f.Architecture, f.Size)
}
```

### Download Binary

Download a specific binary:
```go
bytes, err := client.DownloadBinary("1.2.0", "linux", "amd64")
if err != nil { /* handle */ }
err = api.SaveBinaryToFile(bytes, "./agent_as_code_1.2.0_linux_amd64.zip")
```

Get latest for current platform:
```go
info, err := client.GetLatestBinary()
if err != nil { /* handle */ }
fmt.Println("latest:", info.Version, info.DownloadURL)
```

### Upload Binary

Upload requires an auth token (Bearer) and reads the file from disk. The client computes checksum and base64 content.
```go
resp, err := client.UploadBinary(
    "./dist/agent-linux-amd64.zip", // filePath
    "1.2.0",                        // version
    "linux",                        // platform
    "amd64",                        // arch
)
if err != nil { /* handle */ }
fmt.Println("uploaded:", resp.Success, resp.Release.DownloadURL)
```

### Downloader helper (download & install)
```go
// Create downloader
loader := api.NewDownloader("https://api.myagentregistry.com")

// Download specific
res := loader.DownloadBinary(api.DownloadOptions{
    Version:      "1.2.0",
    Platform:     "darwin",
    Architecture: "arm64",
    OutputDir:    "./artifacts",
})

// Download latest for current OS/ARCH
res = loader.DownloadLatest("./artifacts")

// Install (extract zip and install executable)
res = loader.InstallBinary("1.2.0", "/usr/local/bin")
```

### Uploader helper (multi-platform)
```go
up := api.NewUploader(
    "https://api.myagentregistry.com",
    os.Getenv("AGENT_REGISTRY_TOKEN"),
    "1.2.0",
)

// Upload a single binary
r := up.UploadBinary(api.UploadOptions{
    Platform:     "linux",
    Architecture: "amd64",
    FilePath:     "./dist/agent-linux-amd64.zip",
})

// Upload all platforms (expects files in binDir per platform naming)
results := up.UploadAllPlatforms("./dist")
fmt.Println(api.GetUploadSummary(results))
```

## Profile Configuration & PAT (agent configure)

Profiles are managed by `agent configure profile ...` (implemented in `internal/cmd/configure.go`). Profiles define registry URL and Personal Access Token (PAT). Config is stored at:
- `~/.agent/config.json`

### PAT Format
- 64-character hexadecimal string (validated in code)

### Commands
```bash
# Add a profile
agent configure profile add production \
  --registry https://api.myagentregistry.com \
  --pat a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456

# List profiles
agent configure profile list

# Test a profile
agent configure profile test production

# Set default profile
agent configure profile set-default production

# Remove a profile
agent configure profile remove production
```

### Config file schema (`~/.agent/config.json`)
```json
{
  "profiles": {
    "production": {
      "registry": "https://api.myagentregistry.com",
      "pat": "<64-hex>",
      "description": "Production registry"
    }
  },
  "default_profile": "production"
}
```

## Runtime API (internal/runtime)

Container lifecycle is managed via Docker. See `internal/runtime/runtime.go`.

### Types
```go
type RunOptions struct {
    Image       string
    Ports       []string
    Environment []string
    Detach      bool
    Name        string
    Volumes     []string
    Interactive bool
}

type ContainerInfo struct {
    ID    string
    Name  string
    Ports []PortMapping
}

type PortMapping struct {
    Host      string
    Container string
    Protocol  string
}
```

### Start / Stop / Logs
```go
rt := runtime.New()

info, err := rt.Run(&runtime.RunOptions{
    Image:       "my-agent:latest",
    Ports:       []string{"8080:8080"},
    Environment: []string{"OPENAI_API_KEY=..."},
    Volumes:     []string{"./data:/app/data"},
})

err = rt.StreamLogs(info.ID)
err = rt.Stop(info.ID)
```

## Local LLM Management (agent llm)

Local model integration is implemented in `internal/llm/local_manager.go` and surfaced via `agent llm` in `internal/cmd/llm.go`. Backend: Ollama (`http://localhost:11434`).

### Commands
```bash
# Setup guidance
agent llm setup

# List local models
agent llm list

# Pull a model
agent llm pull llama2
agent llm pull mistral:7b

# Test a model
agent llm test llama2

# Remove a model
agent llm remove llama2

# Recommend models by use case (chatbot, code, general, fast)
agent llm recommend chatbot

# Model info
a gent llm info llama2
```

### Programmatic (from local_manager)
```go
mgr := llm.NewLocalLLMManager()
if err := mgr.CheckOllamaAvailability(); err != nil { /* guide user */ }

models, _ := mgr.ListLocalModels()
_ = mgr.PullModel("llama2:7b")
_ = mgr.TestModel("llama2")
_ = mgr.RemoveModel("llama2")
info, _ := mgr.GetModelInfo("llama2")
```

## Environment Variables (observed by CLI)

The CLI reads certain environment variables (see `internal/cmd/version.go`):
- `AGENT_REGISTRY_TOKEN`: registry authentication token (used for uploads)
- `AGENT_LLM_PROVIDER`, `AGENT_LLM_MODEL`: informational display for current LLM
- `AGENT_RUNTIME_PROVIDER`, `AGENT_RUNTIME_VERSION`: informational display for runtime
- `OPENAI_API_KEY`: commonly passed to agents via environment when running

## See Also
- `internal/api/client.go`, `downloader.go`, `uploader.go`
- `internal/cmd/configure.go` (profiles/PAT)
- `internal/runtime/runtime.go` (container lifecycle)
- `internal/cmd/llm.go`, `internal/llm/local_manager.go` (local LLM)
