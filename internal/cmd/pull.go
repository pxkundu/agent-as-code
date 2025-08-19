package cmd

import (
	"fmt"

	"github.com/pxkundu/agent-as-code/internal/registry"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [OPTIONS] NAME[:TAG]",
	Short: "Pull an agent from a registry",
	Long: `Pull an agent image from a registry.

This command downloads the specified agent image from a container registry
or the Agent as Code registry, making it available to run locally.

Examples:
  agent pull my-agent:latest
  agent pull registry.example.com/my-agent:v1.0.0
  agent pull my-agent --registry myagentregistry.com`,
	Args: cobra.ExactArgs(1),
	RunE: runPull,
}

var (
	pullRegistry string
	pullQuiet    bool
)

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVar(&pullRegistry, "registry", "", "registry to pull from")
	pullCmd.Flags().BoolVarP(&pullQuiet, "quiet", "q", false, "suppress verbose output")
}

func runPull(cmd *cobra.Command, args []string) error {
	imageName := args[0]

	// Initialize registry client
	registryClient := registry.New()

	// Pull options
	options := &registry.PullOptions{
		Image:    imageName,
		Registry: pullRegistry,
		Quiet:    pullQuiet,
	}

	if !pullQuiet {
		fmt.Printf("📥 Pulling %s\n", imageName)
	}

	// Pull the image
	result, err := registryClient.Pull(options)
	if err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}

	// Success message
	if !pullQuiet {
		fmt.Printf("✅ Pull completed successfully!\n")
		fmt.Printf("   Image: %s\n", result.ImageID)
		fmt.Printf("   Size: %s\n", result.Size)
		fmt.Printf("   Digest: %s\n", result.Digest)

		if result.RegistryURL != "" {
			fmt.Printf("   Registry: %s\n", result.RegistryURL)
		}

		fmt.Printf("\n💡 Run with: agent run %s\n", imageName)
	}

	return nil
}
