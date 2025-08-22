package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	// ASCII art banner
	fmt.Print(`
╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  █████████  █████  █████████  █████                                                                                    ║
║ ███░░░░░███ ░░███ ███░░░░░███ ░░███                                                                                    ║
║░███    ░███ ███████ ███     ░░░ ███████   ██████  ███████  ██████                                                    ║
║░███████████ ░░░███░ ░███        ███░░███ ███░░███ ███░░███ ███░░███                                                   ║
║░███░░░░░███  ░███  ░███       ░███ ░███░███ ░███░███ ░███░███████                                                    ║
║░███    ░███  ███  ░░███     ███░███ ░███░███ ░███░███ ░███░███░░░                                                     ║
║ █████   █████░░██████░░██████ ░░██████ ░░████████░░██████░░██████                                                    ║
║ ░░░░░   ░░░░░  ░░░░░░  ░░░░░░   ░░░░░░   ░░░░░░░░  ░░░░░░  ░░░░░░                                                     ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`)

	// Version info
	fmt.Printf("🚀 AaC v%s - Agent as Code\n", version)
	fmt.Println("Declarative AI agent configuration framework")

	// About section
	fmt.Println("╭──────────────────  About  ───────────────────╮")
	fmt.Println("│ 🔗 https://agent-as-code.myagentregistry.com │")
	fmt.Println("│ 🐙 GitHub: @pxkundu/agent-as-code            │")
	fmt.Println("╰──────────────────────────────────────────────╯")

	// System info
	fmt.Printf("💻 System: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🔨 Go Version: %s\n", runtime.Version())

	// Docker info
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err == nil {
		info, err := dockerClient.Info(cmd.Context())
		if err == nil {
			fmt.Printf("🐳 Docker: %s\n", info.ServerVersion)
			fmt.Printf("🏗️  Runtime: %s\n", info.DefaultRuntime)
		}
	}

	// Python info
	pythonVersion := getPythonVersion()
	if pythonVersion != "" {
		fmt.Printf("🐍 Python: %s\n", pythonVersion)
	}

	// Environment info
	if os.Getenv("OPENAI_API_KEY") != "" {
		fmt.Println("🔑 OpenAI API: ✅ Configured")
	}
	if os.Getenv("AGENT_REGISTRY_TOKEN") != "" {
		fmt.Println("📦 Registry: ✅ Authenticated")
	}

	// LLM info
	llmProvider := os.Getenv("AGENT_LLM_PROVIDER")
	llmModel := os.Getenv("AGENT_LLM_MODEL")
	if llmProvider != "" && llmModel != "" {
		fmt.Printf("🤖 LLM: %s:%s\n", llmProvider, llmModel)
	} else if llmProvider != "" {
		fmt.Printf("🤖 LLM: %s\n", llmProvider)
	} else {
		fmt.Println("🤖 LLM: Not configured")
	}

	// Runtime info
	runtimeProvider := os.Getenv("AGENT_RUNTIME_PROVIDER")
	runtimeVersion := os.Getenv("AGENT_RUNTIME_VERSION")
	if runtimeProvider != "" && runtimeVersion != "" {
		fmt.Printf("🏗️  Runtime: %s:%s\n", runtimeProvider, runtimeVersion)
	} else if runtimeProvider != "" {
		fmt.Printf("🏗️  Runtime: %s\n", runtimeProvider)
	} else {
		fmt.Println("🏗️  Runtime: Not configured")
	}

	// Build info
	if commit != "" && commit != "dev" {
		if len(commit) >= 8 {
			fmt.Printf("📝 Build: %s (%s)\n", commit[:8], date)
		} else {
			fmt.Printf("📝 Build: %s (%s)\n", commit, date)
		}
	}
}
