package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Runtime handles agent execution
type Runtime struct {
	dockerClient *client.Client
}

// RunOptions represents runtime options
type RunOptions struct {
	Image       string
	Ports       []string
	Environment []string
	Detach      bool
	Name        string
	Volumes     []string
	Interactive bool
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID    string
	Name  string
	Ports []PortMapping
}

// PortMapping represents port mapping
type PortMapping struct {
	Host      string
	Container string
	Protocol  string
}

// New creates a new runtime instance
func New() *Runtime {
	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		// If Docker is not available, continue without it (will show appropriate error later)
		dockerClient = nil
	}

	return &Runtime{
		dockerClient: dockerClient,
	}
}

// ValidateImage validates that an image exists
func (r *Runtime) ValidateImage(imageName string) error {
	if r.dockerClient == nil {
		return fmt.Errorf("Docker client not available. Please ensure Docker is running")
	}

	ctx := context.Background()
	_, _, err := r.dockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return fmt.Errorf("image '%s' not found locally. Try 'agent pull %s' first", imageName, imageName)
	}

	fmt.Printf("✓ Image found: %s\n", imageName)
	return nil
}

// Run starts an agent container
func (r *Runtime) Run(options *RunOptions) (*ContainerInfo, error) {
	if r.dockerClient == nil {
		return nil, fmt.Errorf("Docker client not available. Please ensure Docker is running")
	}

	ctx := context.Background()

	// Generate container name if not provided
	containerName := options.Name
	if containerName == "" {
		containerName = generateContainerName(options.Image)
	}

	// Parse port mappings
	ports := parsePortMappings(options.Ports)
	portBindings := make(nat.PortMap)
	exposedPorts := make(nat.PortSet)

	for _, port := range ports {
		containerPort := nat.Port(fmt.Sprintf("%s/%s", port.Container, port.Protocol))
		exposedPorts[containerPort] = struct{}{}
		if port.Host != "" {
			portBindings[containerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port.Host,
				},
			}
		}
	}

	// Container configuration
	containerConfig := &container.Config{
		Image:        options.Image,
		Env:          options.Environment,
		ExposedPorts: exposedPorts,
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
	}

	if options.Interactive {
		containerConfig.Tty = true
		containerConfig.OpenStdin = true
		hostConfig.AutoRemove = true
	}

	// Add volume mounts
	if len(options.Volumes) > 0 {
		hostConfig.Binds = options.Volumes
	}

	fmt.Printf("Creating container: %s\n", containerName)

	// Create container
	resp, err := r.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	containerID := resp.ID
	fmt.Printf("Container ID: %s\n", containerID[:12])

	// Show port mappings
	if len(ports) > 0 {
		fmt.Printf("Port mappings:\n")
		for _, port := range ports {
			fmt.Printf("  %s:%s -> %s/%s\n", port.Host, port.Container, port.Container, port.Protocol)
		}
	}

	// Show environment variables
	if len(options.Environment) > 0 {
		fmt.Printf("Environment variables:\n")
		for _, env := range options.Environment {
			fmt.Printf("  %s\n", env)
		}
	}

	// Start the container
	fmt.Printf("Starting container...\n")
	err = r.dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("✅ Container started successfully\n")

	return &ContainerInfo{
		ID:    containerID,
		Name:  containerName,
		Ports: ports,
	}, nil
}

// Stop stops a running container
func (r *Runtime) Stop(containerID string) error {
	if r.dockerClient == nil {
		return fmt.Errorf("Docker client not available")
	}

	ctx := context.Background()
	timeout := int(30) // 30 second timeout

	fmt.Printf("Stopping container %s...\n", containerID[:12])

	err := r.dockerClient.ContainerStop(ctx, containerID, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Printf("✅ Container stopped\n")
	return nil
}

// StreamLogs streams container logs
func (r *Runtime) StreamLogs(containerID string) error {
	if r.dockerClient == nil {
		return fmt.Errorf("Docker client not available")
	}

	ctx := context.Background()

	fmt.Printf("Streaming logs for container %s...\n", containerID[:12])

	// Get container logs
	reader, err := r.dockerClient.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get container logs: %w", err)
	}
	defer reader.Close()

	// Stream logs to stdout
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("failed to stream logs: %w", err)
	}

	return nil
}

// List lists running containers
func (r *Runtime) List() ([]ContainerInfo, error) {
	// In a real implementation, this would list actual containers
	return []ContainerInfo{
		{
			ID:   "abcdef123456",
			Name: "my-agent",
			Ports: []PortMapping{
				{Host: "8080", Container: "8080", Protocol: "tcp"},
			},
		},
	}, nil
}

// Helper functions
func generateContainerName(imageName string) string {
	// Generate a unique container name based on image
	timestamp := time.Now().Unix()
	return fmt.Sprintf("agent-%d", timestamp)
}

func parsePortMappings(ports []string) []PortMapping {
	var mappings []PortMapping

	for _, portStr := range ports {
		// Parse port strings like "8080:8080", "80:8080/tcp", "8080"
		mapping := PortMapping{
			Protocol: "tcp", // Default protocol
		}

		// Split by protocol if specified
		parts := strings.Split(portStr, "/")
		if len(parts) == 2 {
			mapping.Protocol = parts[1]
			portStr = parts[0]
		}

		// Split host:container ports
		portParts := strings.Split(portStr, ":")
		switch len(portParts) {
		case 1:
			// Only container port specified (e.g., "8080")
			mapping.Container = strings.TrimSpace(portParts[0])
			mapping.Host = strings.TrimSpace(portParts[0]) // Same as container
		case 2:
			// Both host and container ports (e.g., "80:8080")
			mapping.Host = strings.TrimSpace(portParts[0])
			mapping.Container = strings.TrimSpace(portParts[1])
		default:
			// Invalid format, skip
			continue
		}

		// Validate port numbers
		if isValidPort(mapping.Host) && isValidPort(mapping.Container) {
			mappings = append(mappings, mapping)
		}
	}

	// Default port mapping if none specified
	if len(mappings) == 0 {
		mappings = append(mappings, PortMapping{
			Host:      "8080",
			Container: "8080",
			Protocol:  "tcp",
		})
	}

	return mappings
}

// isValidPort checks if a port string is a valid port number
func isValidPort(portStr string) bool {
	if portStr == "" {
		return false
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}

	return port > 0 && port <= 65535
}
