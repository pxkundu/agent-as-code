package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pxkundu/agent-as-code/internal/parser"
)

// Builder handles agent building
type Builder struct {
	parser       *parser.Parser
	dockerClient *client.Client
}

// BuildOptions represents build options
type BuildOptions struct {
	Path     string
	Tag      string
	NoCache  bool
	Push     bool
	Platform string
}

// BuildResult represents build result
type BuildResult struct {
	ImageID string
	Size    string
	Tags    []string
}

// New creates a new builder instance
func New() *Builder {
	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		// If Docker is not available, continue without it (will show appropriate error later)
		dockerClient = nil
	}

	return &Builder{
		parser:       parser.New(),
		dockerClient: dockerClient,
	}
}

// ValidateContext validates the build context
func (b *Builder) ValidateContext(path string) error {
	// Check if agent.yaml exists
	agentFile, err := b.parser.FindAgentFile(path)
	if err != nil {
		return fmt.Errorf("no agent.yaml found: %w", err)
	}

	// Parse and validate agent.yaml
	_, err = b.parser.ParseFile(agentFile)
	if err != nil {
		return fmt.Errorf("invalid agent.yaml: %w", err)
	}

	return nil
}

// Build builds an agent from the given options
func (b *Builder) Build(options *BuildOptions) (*BuildResult, error) {
	// Find and parse agent.yaml
	agentFile, err := b.parser.FindAgentFile(options.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to find agent.yaml: %w", err)
	}

	spec, err := b.parser.ParseFile(agentFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent.yaml: %w", err)
	}

	// Generate Dockerfile
	dockerfile, err := b.generateDockerfile(spec, options.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	// Write Dockerfile to build context
	dockerfilePath := filepath.Join(options.Path, "Dockerfile.agent")
	if err := b.writeDockerfile(dockerfilePath, dockerfile); err != nil {
		return nil, fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	// Build Docker image
	imageID, err := b.buildDockerImage(options, dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("docker build failed: %w", err)
	}

	// Get image size
	size, err := b.getImageSize(imageID)
	if err != nil {
		size = "unknown"
	}

	// Prepare result
	result := &BuildResult{
		ImageID: imageID,
		Size:    size,
		Tags:    []string{},
	}

	if options.Tag != "" {
		result.Tags = append(result.Tags, options.Tag)
	}

	return result, nil
}

// generateDockerfile generates a Dockerfile from agent spec
func (b *Builder) generateDockerfile(spec *parser.AgentSpec, contextPath string) (string, error) {
	dockerfile := ""

	// Base image based on runtime
	switch spec.Spec.Runtime {
	case "python":
		dockerfile += "FROM python:3.11-slim\n\n"
	case "nodejs":
		dockerfile += "FROM node:18-slim\n\n"
	case "go":
		dockerfile += "FROM golang:1.21-alpine AS builder\n"
		dockerfile += "FROM alpine:latest\n\n"
	default:
		return "", fmt.Errorf("unsupported runtime: %s", spec.Spec.Runtime)
	}

	// Set working directory
	dockerfile += "WORKDIR /app\n\n"

	// Install dependencies
	if len(spec.Spec.Dependencies) > 0 {
		switch spec.Spec.Runtime {
		case "python":
			dockerfile += "# Install Python dependencies\n"
			dockerfile += "COPY requirements.txt .\n"
			dockerfile += "RUN pip install --no-cache-dir -r requirements.txt\n\n"
		case "nodejs":
			dockerfile += "# Install Node.js dependencies\n"
			dockerfile += "COPY package*.json .\n"
			dockerfile += "RUN npm ci --only=production\n\n"
		}
	}

	// Copy application code
	dockerfile += "# Copy application code\n"
	dockerfile += "COPY . .\n\n"

	// Set environment variables
	if len(spec.Spec.Environment) > 0 {
		dockerfile += "# Environment variables\n"
		for _, env := range spec.Spec.Environment {
			if env.Value != "" {
				dockerfile += fmt.Sprintf("ENV %s=%s\n", env.Name, env.Value)
			}
		}
		dockerfile += "\n"
	}

	// Expose ports
	if len(spec.Spec.Ports) > 0 {
		dockerfile += "# Expose ports\n"
		for _, port := range spec.Spec.Ports {
			dockerfile += fmt.Sprintf("EXPOSE %d\n", port.Container)
		}
		dockerfile += "\n"
	}

	// Health check
	if spec.Spec.HealthCheck != nil {
		dockerfile += "# Health check\n"
		dockerfile += "HEALTHCHECK "
		if spec.Spec.HealthCheck.Interval != "" {
			dockerfile += fmt.Sprintf("--interval=%s ", spec.Spec.HealthCheck.Interval)
		}
		if spec.Spec.HealthCheck.Timeout != "" {
			dockerfile += fmt.Sprintf("--timeout=%s ", spec.Spec.HealthCheck.Timeout)
		}
		if spec.Spec.HealthCheck.Retries > 0 {
			dockerfile += fmt.Sprintf("--retries=%d ", spec.Spec.HealthCheck.Retries)
		}
		if spec.Spec.HealthCheck.StartPeriod != "" {
			dockerfile += fmt.Sprintf("--start-period=%s ", spec.Spec.HealthCheck.StartPeriod)
		}
		dockerfile += "CMD " + joinCommand(spec.Spec.HealthCheck.Command) + "\n\n"
	}

	// Default command
	switch spec.Spec.Runtime {
	case "python":
		dockerfile += "# Run the application\n"
		dockerfile += "CMD [\"python\", \"main.py\"]\n"
	case "nodejs":
		dockerfile += "# Run the application\n"
		dockerfile += "CMD [\"node\", \"index.js\"]\n"
	case "go":
		dockerfile += "# Run the application\n"
		dockerfile += "CMD [\"./app\"]\n"
	}

	return dockerfile, nil
}

// writeDockerfile writes the Dockerfile to disk
func (b *Builder) writeDockerfile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// buildDockerImage builds the Docker image
func (b *Builder) buildDockerImage(options *BuildOptions, dockerfilePath string) (string, error) {
	if b.dockerClient == nil {
		return "", fmt.Errorf("Docker client not available. Please ensure Docker is running")
	}

	ctx := context.Background()

	// Create build context
	buildContext, err := b.createBuildContext(options.Path, dockerfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create build context: %w", err)
	}

	// Build options
	buildOpts := types.ImageBuildOptions{
		Dockerfile: filepath.Base(dockerfilePath),
		Tags:       []string{},
		Remove:     true,
		NoCache:    options.NoCache,
	}

	if options.Tag != "" {
		buildOpts.Tags = append(buildOpts.Tags, options.Tag)
	}

	// Build the image
	fmt.Printf("Building Docker image...\n")
	resp, err := b.dockerClient.ImageBuild(ctx, buildContext, buildOpts)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Stream build output
	var imageID string
	decoder := json.NewDecoder(resp.Body)
	for {
		var buildLine struct {
			Stream string `json:"stream"`
			Aux    struct {
				ID string `json:"ID"`
			} `json:"aux"`
			Error string `json:"error"`
		}

		if err := decoder.Decode(&buildLine); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to decode build output: %w", err)
		}

		if buildLine.Error != "" {
			return "", fmt.Errorf("build error: %s", buildLine.Error)
		}

		if buildLine.Stream != "" {
			fmt.Print(buildLine.Stream)
		}

		if buildLine.Aux.ID != "" {
			imageID = buildLine.Aux.ID
		}
	}

	if imageID == "" {
		return "", fmt.Errorf("failed to get image ID from build output")
	}

	fmt.Printf("Successfully built %s\n", imageID[:12])
	if options.Tag != "" {
		fmt.Printf("Successfully tagged %s\n", options.Tag)
	}

	return imageID, nil
}

// getImageSize gets the size of a Docker image
func (b *Builder) getImageSize(imageID string) (string, error) {
	if b.dockerClient == nil {
		return "unknown", nil
	}

	ctx := context.Background()
	imageInspect, _, err := b.dockerClient.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return "unknown", err
	}

	size := imageInspect.Size
	return formatSize(size), nil
}

// Push pushes the image to a registry
func (b *Builder) Push(tag string) error {
	if b.dockerClient == nil {
		return fmt.Errorf("Docker client not available")
	}

	ctx := context.Background()

	// Push the image
	fmt.Printf("Pushing %s...\n", tag)
	resp, err := b.dockerClient.ImagePush(ctx, tag, types.ImagePushOptions{})
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	defer resp.Close()

	// Stream push output
	decoder := json.NewDecoder(resp)
	for {
		var pushLine struct {
			Status   string `json:"status"`
			Progress string `json:"progress"`
			Error    string `json:"error"`
		}

		if err := decoder.Decode(&pushLine); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode push output: %w", err)
		}

		if pushLine.Error != "" {
			return fmt.Errorf("push error: %s", pushLine.Error)
		}

		if pushLine.Status != "" {
			fmt.Printf("%s\n", pushLine.Status)
			if pushLine.Progress != "" {
				fmt.Printf(" %s", pushLine.Progress)
			}
		}
	}

	fmt.Printf("Push completed successfully\n")
	return nil
}

// createBuildContext creates a tar archive of the build context
func (b *Builder) createBuildContext(buildPath, dockerfilePath string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Walk through the build directory
	err := filepath.Walk(buildPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") && path != buildPath {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(buildPath, path)
		if err != nil {
			return err
		}

		// Skip the build path itself
		if relPath == "." {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// formatSize formats bytes to human readable string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper functions
func joinCommand(cmd []string) string {
	if len(cmd) == 0 {
		return ""
	}

	result := "["
	for i, part := range cmd {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", part)
	}
	result += "]"
	return result
}
