package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pxkundu/agent-as-code/internal/runtime"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [OPTIONS] IMAGE",
	Short: "Run an agent container",
	Long: `Run an agent container from a built image.

This command starts an agent container and manages its lifecycle.
The agent will be accessible on the specified port (default: 8080).

Examples:
  agent run my-agent:latest
  agent run -p 9000:8080 my-agent:latest
  agent run --env OPENAI_API_KEY=sk-... my-agent:latest
  agent run -d my-agent:latest`,
	Args: cobra.ExactArgs(1),
	RunE: runRun,
}

var (
	runPort        []string
	runEnv         []string
	runDetach      bool
	runName        string
	runVolume      []string
	runInteractive bool
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceVarP(&runPort, "port", "p", []string{}, "publish a container's port(s) to the host")
	runCmd.Flags().StringSliceVarP(&runEnv, "env", "e", []string{}, "set environment variables")
	runCmd.Flags().BoolVarP(&runDetach, "detach", "d", false, "run container in background")
	runCmd.Flags().StringVar(&runName, "name", "", "assign a name to the container")
	runCmd.Flags().StringSliceVarP(&runVolume, "volume", "v", []string{}, "bind mount a volume")
	runCmd.Flags().BoolVarP(&runInteractive, "interactive", "i", false, "run in interactive mode")
}

func runRun(cmd *cobra.Command, args []string) error {
	imageName := args[0]

	// Initialize runtime
	agentRuntime := runtime.New()

	// Run options
	options := &runtime.RunOptions{
		Image:       imageName,
		Ports:       runPort,
		Environment: runEnv,
		Detach:      runDetach,
		Name:        runName,
		Volumes:     runVolume,
		Interactive: runInteractive,
	}

	// Validate image exists
	if err := agentRuntime.ValidateImage(imageName); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	fmt.Printf("🚀 Starting agent: %s\n", imageName)

	// Start the agent
	container, err := agentRuntime.Run(options)
	if err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	if runDetach {
		fmt.Printf("✅ Agent started in background\n")
		fmt.Printf("   Container ID: %s\n", container.ID[:12])
		fmt.Printf("   Name: %s\n", container.Name)

		// Show port mappings
		if len(container.Ports) > 0 {
			fmt.Printf("   Ports:\n")
			for _, port := range container.Ports {
				fmt.Printf("     %s -> %s\n", port.Host, port.Container)
			}
		}

		fmt.Printf("\n💡 Use 'agent logs %s' to view logs\n", container.Name)
		fmt.Printf("💡 Use 'agent stop %s' to stop the agent\n", container.Name)
	} else {
		fmt.Printf("✅ Agent started successfully\n")
		fmt.Printf("   Container: %s\n", container.Name)

		// Show access information
		if len(container.Ports) > 0 {
			for _, port := range container.Ports {
				if port.Host != "" {
					fmt.Printf("   Access: http://localhost:%s\n", port.Host)
					break
				}
			}
		}

		fmt.Printf("\n📋 Press Ctrl+C to stop the agent\n\n")

		// Wait for interrupt signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		// Stream logs in foreground mode
		if !runDetach {
			go func() {
				if err := agentRuntime.StreamLogs(container.ID); err != nil {
					fmt.Printf("Error streaming logs: %v\n", err)
				}
			}()
		}

		// Wait for signal
		<-c
		fmt.Printf("\n🛑 Stopping agent...\n")

		// Stop the container
		if err := agentRuntime.Stop(container.ID); err != nil {
			return fmt.Errorf("failed to stop agent: %w", err)
		}

		fmt.Printf("✅ Agent stopped\n")
	}

	return nil
}
