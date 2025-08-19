package cmd

import (
	"fmt"

	"github.com/pxkundu/agent-as-code/internal/registry"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [OPTIONS] NAME[:TAG]",
	Short: "Push an agent to a registry",
	Long: `Push an agent image to a registry.

This command uploads the specified agent image to a container registry
or the Agent as Code registry, making it available for others to pull
and use.

Examples:
  agent push my-agent:latest
  agent push registry.example.com/my-agent:v1.0.0
  agent push my-agent --registry myagentregistry.com`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

var (
	pushRegistry string
	pushAll      bool
)

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&pushRegistry, "registry", "", "registry to push to")
	pushCmd.Flags().BoolVarP(&pushAll, "all-tags", "a", false, "push all tagged images in the repository")
}

func runPush(cmd *cobra.Command, args []string) error {
	imageName := args[0]

	// Initialize registry client
	registryClient := registry.New()

	// Push options
	options := &registry.PushOptions{
		Image:    imageName,
		Registry: pushRegistry,
		AllTags:  pushAll,
	}

	// Validate image exists locally
	if err := registryClient.ValidateLocalImage(imageName); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	fmt.Printf("📤 Pushing %s\n", imageName)

	// Push the image
	result, err := registryClient.Push(options)
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	// Success message
	fmt.Printf("✅ Push completed successfully!\n")
	fmt.Printf("   Repository: %s\n", result.Repository)
	fmt.Printf("   Tag: %s\n", result.Tag)
	fmt.Printf("   Digest: %s\n", result.Digest)
	fmt.Printf("   Size: %s\n", result.Size)

	// Show registry URL if available
	if result.RegistryURL != "" {
		fmt.Printf("   Registry: %s\n", result.RegistryURL)
		fmt.Printf("\n💡 Others can now pull with: agent pull %s\n", imageName)
	}

	return nil
}
