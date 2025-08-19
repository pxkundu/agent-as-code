package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rmiCmd = &cobra.Command{
	Use:   "rmi [TAG]",
	Short: "Remove agent image",
	Long: `Remove an agent image from the local system.

This command removes the specified agent image, freeing up disk space.
If the image is currently being used by running containers, the removal
will fail unless the --force flag is used.

Examples:
  agent rmi my-agent:latest
  agent rmi my-agent:v1.0.0
  agent rmi --force my-agent:latest
  agent rmi --all-tags my-agent`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]
		force, _ := cmd.Flags().GetBool("force")
		allTags, _ := cmd.Flags().GetBool("all-tags")
		
		if allTags {
			return removeAllTags(strings.Split(tag, ":")[0], force)
		}
		
		return removeImage(tag, force)
	},
}

func init() {
	rmiCmd.Flags().Bool("force", false, "force removal even if image is in use")
	rmiCmd.Flags().Bool("all-tags", false, "remove all tags for the specified image")
	rootCmd.AddCommand(rmiCmd)
}

func removeImage(tag string, force bool) error {
	fmt.Printf("🗑️  Removing agent image: %s\n", tag)
	
	// Check if the image exists
	if !imageExists(tag) {
		return fmt.Errorf("agent image '%s' not found", tag)
	}
	
	// Check if the image is being used by running containers
	if !force && imageInUse(tag) {
		return fmt.Errorf("cannot remove image '%s': image is in use by running containers. Use --force to override", tag)
	}
	
	// Remove the image
	args := []string{"rmi"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, tag)
	
	rmiCmd := exec.Command("docker", args...)
	if err := rmiCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove image '%s': %v", tag, err)
	}
	
	fmt.Printf("✅ Successfully removed agent image: %s\n", tag)
	return nil
}

func removeAllTags(imageName string, force bool) error {
	fmt.Printf("🗑️  Removing all tags for agent: %s\n", imageName)
	
	// Get all tags for the image
	tags, err := getImageTags(imageName)
	if err != nil {
		return fmt.Errorf("failed to get tags for image '%s': %v", imageName, err)
	}
	
	if len(tags) == 0 {
		fmt.Printf("ℹ️  No tags found for image: %s\n", imageName)
		return nil
	}
	
	fmt.Printf("Found %d tags: %s\n", len(tags), strings.Join(tags, ", "))
	
	// Remove each tag
	removedCount := 0
	for _, tag := range tags {
		if err := removeImage(tag, force); err != nil {
			fmt.Printf("⚠️  Warning: failed to remove tag '%s': %v\n", tag, err)
			continue
		}
		removedCount++
	}
	
	if removedCount > 0 {
		fmt.Printf("✅ Successfully removed %d/%d tags for agent: %s\n", removedCount, len(tags), imageName)
	}
	
	if removedCount < len(tags) {
		return fmt.Errorf("some tags could not be removed. Check warnings above")
	}
	
	return nil
}

func imageExists(tag string) bool {
	// Check if the image exists using docker images
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	images := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, image := range images {
		if image == tag {
			return true
		}
	}
	
	return false
}

func imageInUse(tag string) bool {
	// Check if the image is being used by running containers
	cmd := exec.Command("docker", "ps", "--format", "{{.Image}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	containers := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, container := range containers {
		if container == tag {
			return true
		}
	}
	
	return false
}

func getImageTags(imageName string) ([]string, error) {
	// Get all tags for the specified image
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	var tags []string
	images := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, image := range images {
		if strings.HasPrefix(image, imageName+":") {
			tags = append(tags, image)
		}
	}
	
	return tags, nil
}
