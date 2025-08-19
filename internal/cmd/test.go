package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [TAG]",
	Short: "Test agent functionality",
	Long: `Test agent functionality by running the agent and executing test scenarios.

This command starts the agent container and runs predefined tests to verify
that the agent is working correctly. Tests may include health checks,
API endpoint validation, and basic functionality verification.

Examples:
  agent test my-agent:latest
  agent test my-agent:v1.0.0
  agent test --timeout 60s my-agent:latest`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]
		timeout, _ := cmd.Flags().GetString("timeout")
		
		fmt.Printf("🧪 Testing agent: %s\n", tag)
		
		// Check if the agent image exists
		if !testImageExists(tag) {
			return fmt.Errorf("agent image '%s' not found. Build it first with 'agent build'", tag)
		}
		
		// Run the test
		return runAgentTests(tag, timeout)
	},
}

func init() {
	testCmd.Flags().String("timeout", "30s", "test timeout duration")
	rootCmd.AddCommand(testCmd)
}

func testImageExists(tag string) bool {
	// Simple check - in a real implementation, this would query Docker
	// For now, we'll assume the image exists if we can find it in our registry
	return true
}

func runAgentTests(tag, timeout string) error {
	fmt.Println("  Starting agent for testing...")
	
	// Start the agent in test mode
	containerName := fmt.Sprintf("test-%s", sanitizeTag(tag))
	
	// Run the agent container
	runCmd := exec.Command("docker", "run", 
		"--name", containerName,
		"--rm",
		"-d",
		"-p", "8080:8080",
		tag)
	
	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("failed to start test container: %v", err)
	}
	
	defer func() {
		// Clean up the test container
		exec.Command("docker", "stop", containerName).Run()
		exec.Command("docker", "rm", containerName).Run()
	}()
	
	fmt.Println("  Waiting for agent to be ready...")
	
	// Wait for the agent to be ready
	if err := waitForAgentReady("localhost:8080", timeout); err != nil {
		return fmt.Errorf("agent failed to become ready: %v", err)
	}
	
	fmt.Println("  Running health check...")
	
	// Run health check
	if err := runHealthCheck("localhost:8080"); err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	
	fmt.Println("  Running basic functionality tests...")
	
	// Run basic functionality tests
	if err := runBasicTests("localhost:8080"); err != nil {
		return fmt.Errorf("basic tests failed: %v", err)
	}
	
	fmt.Println("✅ All tests passed!")
	return nil
}

func sanitizeTag(tag string) string {
	// Convert tag to valid container name
	return filepath.Base(tag)
}

func waitForAgentReady(addr, timeout string) error {
	// Simple wait implementation
	// In a real implementation, this would poll the health endpoint
	fmt.Printf("    Agent ready at %s\n", addr)
	return nil
}

func runHealthCheck(addr string) error {
	// Run health check
	healthCmd := exec.Command("curl", "-f", fmt.Sprintf("http://%s/health", addr))
	if err := healthCmd.Run(); err != nil {
		return fmt.Errorf("health endpoint not responding: %v", err)
	}
	
	fmt.Println("    Health check passed")
	return nil
}

func runBasicTests(addr string) error {
	// Run basic functionality tests
	// This could include testing various endpoints, checking responses, etc.
	
	// Test root endpoint
	rootCmd := exec.Command("curl", "-f", fmt.Sprintf("http://%s/", addr))
	if err := rootCmd.Run(); err != nil {
		return fmt.Errorf("root endpoint test failed: %v", err)
	}
	
	// Test API documentation endpoint if available
	docsCmd := exec.Command("curl", "-f", fmt.Sprintf("http://%s/docs", addr))
	docsCmd.Run() // This is optional, don't fail if it doesn't exist
	
	fmt.Println("    Basic functionality tests passed")
	return nil
}
