package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version string
	commit  string
	date    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "Docker-like CLI for AI agents",
	Long: `Agent as Code (AaC) - A Docker-like CLI for building, running, and managing AI agents.

Think of it as "Docker for AI agents" - create, package, and share AI agents
using simple, declarative configuration files.

Examples:
  agent init my-chatbot --template chatbot
  agent build -t my-chatbot:latest .
  agent run my-chatbot:latest
  agent push my-chatbot:latest`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = getVersionString()
}

func getVersionString() string {
	// ASCII art banner
	banner := `
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

`
	// Build version info
	info := fmt.Sprintf("\n🚀 AaC v%s - Agent as Code\n", version)
	info += "Declarative AI agent configuration framework\n"

	// About section
	info += "╭──────────────────  About  ───────────────────╮\n"
	info += "│ 🔗 https://agent-as-code.myagentregistry.com │\n"
	info += "│ 🐙 GitHub: @pxkundu/agent-as-code            │\n"
	info += "╰──────────────────────────────────────────────╯\n"

	// System info
	info += fmt.Sprintf("💻 System: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	info += fmt.Sprintf("🔨 Go Version: %s\n", runtime.Version())

	// Python info
	pythonVersion := getPythonVersion()
	if pythonVersion != "" {
		info += fmt.Sprintf("🐍 Python: %s\n", pythonVersion)
	}

	// Docker info
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err == nil {
		dockerInfo, err := dockerClient.Info(context.Background())
		if err == nil {
			info += fmt.Sprintf("🐳 Docker: %s\n", dockerInfo.ServerVersion)
			info += fmt.Sprintf("🏗️  Runtime: %s\n", dockerInfo.DefaultRuntime)
		}
	}

	// Build info
	if commit != "" && commit != "dev" {
		if len(commit) >= 8 {
			info += fmt.Sprintf("📝 Build: %s (%s)\n", commit[:8], date)
		} else {
			info += fmt.Sprintf("📝 Build: %s (%s)\n", commit, date)
		}
	}

	return banner + info
}

func getPythonVersion() string {
	// Try python3 first
	output, err := exec.Command("python3", "--version").Output()
	if err == nil {
		return strings.TrimSpace(strings.TrimPrefix(string(output), "Python "))
	}

	// Try python as fallback
	output, err = exec.Command("python", "--version").Output()
	if err == nil {
		return strings.TrimSpace(strings.TrimPrefix(string(output), "Python "))
	}

	return ""
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agent-as-code.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".agent-as-code" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".agent-as-code")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
