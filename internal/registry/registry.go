package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Registry handles registry operations
type Registry struct {
	dockerClient *client.Client
	registryURL  string
	authToken    string
}

// PushOptions represents push options
type PushOptions struct {
	Image    string
	Registry string
	AllTags  bool
}

// PullOptions represents pull options
type PullOptions struct {
	Image    string
	Registry string
	Quiet    bool
}

// ListOptions represents list options
type ListOptions struct {
	Filter []string
	All    bool
}

// PushResult represents push result
type PushResult struct {
	Repository  string
	Tag         string
	Digest      string
	Size        string
	RegistryURL string
}

// PullResult represents pull result
type PullResult struct {
	ImageID     string
	Size        string
	Digest      string
	RegistryURL string
}

// ImageInfo represents image information
type ImageInfo struct {
	ID         string
	Repository string
	Tag        string
	Created    time.Time
	Size       int64
}

// New creates a new registry instance
func New() *Registry {
	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		dockerClient = nil
	}

	return &Registry{
		dockerClient: dockerClient,
		registryURL:  os.Getenv("AGENT_REGISTRY_URL"),
		authToken:    os.Getenv("AGENT_REGISTRY_TOKEN"),
	}
}

// ValidateLocalImage validates that an image exists locally
func (r *Registry) ValidateLocalImage(imageName string) error {
	if r.dockerClient == nil {
		return fmt.Errorf("Docker client not available. Please ensure Docker is running")
	}

	ctx := context.Background()
	_, _, err := r.dockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return fmt.Errorf("image '%s' not found locally. Build it first with 'agent build'", imageName)
	}

	fmt.Printf("✓ Local image validated: %s\n", imageName)
	return nil
}

// Push pushes an image to a registry
func (r *Registry) Push(options *PushOptions) (*PushResult, error) {
	if r.dockerClient == nil {
		return nil, fmt.Errorf("Docker client not available")
	}

	fmt.Printf("Pushing image: %s\n", options.Image)

	// Use registry-specific logic or Docker Hub
	if r.isAgentRegistry(options.Registry) {
		return r.pushToAgentRegistry(options)
	}

	// Default Docker registry push
	return r.pushToDockerRegistry(options)
}

// Pull pulls an image from a registry
func (r *Registry) Pull(options *PullOptions) (*PullResult, error) {
	if r.dockerClient == nil {
		return nil, fmt.Errorf("Docker client not available")
	}

	if !options.Quiet {
		fmt.Printf("Pulling image: %s\n", options.Image)
	}

	// Use registry-specific logic or Docker Hub
	if r.isAgentRegistry(options.Registry) {
		return r.pullFromAgentRegistry(options)
	}

	// Default Docker registry pull
	return r.pullFromDockerRegistry(options)
}

// ListLocal lists local images
func (r *Registry) ListLocal(options *ListOptions) ([]ImageInfo, error) {
	if r.dockerClient == nil {
		return nil, fmt.Errorf("Docker client not available")
	}

	ctx := context.Background()

	// List Docker images
	dockerImages, err := r.dockerClient.ImageList(ctx, types.ImageListOptions{
		All: options.All,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list local images: %w", err)
	}

	var images []ImageInfo
	for _, img := range dockerImages {
		// Skip images without repository tags
		if len(img.RepoTags) == 0 {
			continue
		}

		for _, repoTag := range img.RepoTags {
			repository, tag := parseImageName(repoTag)

			imageInfo := ImageInfo{
				ID:         img.ID,
				Repository: repository,
				Tag:        tag,
				Created:    time.Unix(img.Created, 0),
				Size:       img.Size,
			}

			// Apply filters
			if r.matchesFilters(imageInfo, options.Filter) {
				images = append(images, imageInfo)
			}
		}
	}

	return images, nil
}

// isAgentRegistry checks if we're using the agent registry
func (r *Registry) isAgentRegistry(registryURL string) bool {
	if registryURL == "" {
		registryURL = r.registryURL
	}
	return strings.Contains(registryURL, "myagentregistry.com") || strings.Contains(registryURL, "agent-registry")
}

// pushToAgentRegistry pushes to the agent registry using the documented API
func (r *Registry) pushToAgentRegistry(options *PushOptions) (*PushResult, error) {
	// This would implement the actual agent registry push logic
	// For now, fall back to Docker registry
	return r.pushToDockerRegistry(options)
}

// pullFromAgentRegistry pulls from the agent registry
func (r *Registry) pullFromAgentRegistry(options *PullOptions) (*PullResult, error) {
	// This would implement the actual agent registry pull logic
	// For now, fall back to Docker registry
	return r.pullFromDockerRegistry(options)
}

// pushToDockerRegistry pushes to Docker registry
func (r *Registry) pushToDockerRegistry(options *PushOptions) (*PushResult, error) {
	ctx := context.Background()

	// Push the image
	resp, err := r.dockerClient.ImagePush(ctx, options.Image, types.ImagePushOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to push image: %w", err)
	}
	defer resp.Close()

	// Stream push output
	if _, err := io.Copy(os.Stdout, resp); err != nil {
		return nil, fmt.Errorf("failed to stream push output: %w", err)
	}

	// Parse image name
	repository, tag := parseImageName(options.Image)

	return &PushResult{
		Repository:  repository,
		Tag:         tag,
		Digest:      "sha256:unknown", // Would be extracted from response
		Size:        "unknown",
		RegistryURL: options.Registry,
	}, nil
}

// pullFromDockerRegistry pulls from Docker registry
func (r *Registry) pullFromDockerRegistry(options *PullOptions) (*PullResult, error) {
	ctx := context.Background()

	// Pull the image
	resp, err := r.dockerClient.ImagePull(ctx, options.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull image: %w", err)
	}
	defer resp.Close()

	// Stream pull output if not quiet
	if !options.Quiet {
		if _, err := io.Copy(os.Stdout, resp); err != nil {
			return nil, fmt.Errorf("failed to stream pull output: %w", err)
		}
	} else {
		// Still need to read the response to complete the pull
		if _, err := io.Copy(io.Discard, resp); err != nil {
			return nil, fmt.Errorf("failed to complete pull: %w", err)
		}
	}

	return &PullResult{
		ImageID:     "sha256:unknown", // Would be extracted from response
		Size:        "unknown",
		Digest:      "sha256:unknown",
		RegistryURL: options.Registry,
	}, nil
}

// matchesFilters checks if an image matches the given filters
func (r *Registry) matchesFilters(image ImageInfo, filters []string) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		// Simple filter matching
		if strings.Contains(image.Repository, filter) || strings.Contains(image.Tag, filter) {
			return true
		}
	}

	return false
}

// Helper functions
func parseImageName(imageName string) (repository, tag string) {
	// Split on the last ':' to handle registry URLs with ports
	lastColon := strings.LastIndex(imageName, ":")
	if lastColon == -1 {
		return imageName, "latest"
	}

	// Check if what's after the colon looks like a tag (not a port)
	potentialTag := imageName[lastColon+1:]
	if strings.Contains(potentialTag, "/") {
		// This is likely a port number, not a tag
		return imageName, "latest"
	}

	return imageName[:lastColon], potentialTag
}
