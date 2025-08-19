package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/pxkundu/agent-as-code/internal/builder"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [OPTIONS] PATH",
	Short: "Build an agent from agent.yaml",
	Long: `Build an agent container from the agent.yaml configuration file.

This command reads the agent.yaml file in the specified directory,
validates the configuration, and builds a container image that can
be run locally or pushed to a registry.

Examples:
  agent build .
  agent build -t my-agent:latest .
  agent build -t my-agent:v1.0.0 ./my-agent-dir
  agent build --no-cache -t my-agent .`,
	Args: cobra.ExactArgs(1),
	RunE: runBuild,
}

var (
	buildTag      string
	buildNoCache  bool
	buildPush     bool
	buildPlatform string
)

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&buildTag, "tag", "t", "", "name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().BoolVar(&buildNoCache, "no-cache", false, "do not use cache when building the image")
	buildCmd.Flags().BoolVar(&buildPush, "push", false, "push the image to registry after building")
	buildCmd.Flags().StringVar(&buildPlatform, "platform", "", "set platform if server is multi-platform capable")
}

func runBuild(cmd *cobra.Command, args []string) error {
	buildPath := args[0]

	// Convert to absolute path
	absPath, err := filepath.Abs(buildPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Initialize builder
	agentBuilder := builder.New()

	// Build options
	options := &builder.BuildOptions{
		Path:     absPath,
		Tag:      buildTag,
		NoCache:  buildNoCache,
		Push:     buildPush,
		Platform: buildPlatform,
	}

	// Validate build context
	if err := agentBuilder.ValidateContext(absPath); err != nil {
		return fmt.Errorf("invalid build context: %w", err)
	}

	fmt.Printf("🔨 Building agent from %s\n", absPath)

	// Build the agent
	result, err := agentBuilder.Build(options)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Success message
	fmt.Printf("✅ Agent built successfully!\n")
	fmt.Printf("   Image: %s\n", result.ImageID)
	fmt.Printf("   Size: %s\n", result.Size)

	if buildTag != "" {
		fmt.Printf("   Tag: %s\n", buildTag)
	}

	if buildPush {
		fmt.Printf("📤 Pushing to registry...\n")
		if err := agentBuilder.Push(buildTag); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
		fmt.Printf("✅ Push completed!\n")
	}

	return nil
}
