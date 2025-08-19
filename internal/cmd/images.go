package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/pxkundu/agent-as-code/internal/registry"
	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:   "images [OPTIONS]",
	Short: "List agent images",
	Long: `List agent images available locally.

This command shows all agent images that have been built or pulled
to the local system, along with their tags, sizes, and creation dates.

Examples:
  agent images
  agent images --filter "name=my-agent"
  agent images --format json
  agent images -q`,
	RunE: runImages,
}

var (
	imagesFilter []string
	imagesFormat string
	imagesQuiet  bool
	imagesAll    bool
)

func init() {
	rootCmd.AddCommand(imagesCmd)

	imagesCmd.Flags().StringSliceVar(&imagesFilter, "filter", []string{}, "filter output based on conditions provided")
	imagesCmd.Flags().StringVar(&imagesFormat, "format", "table", "pretty-print images using a Go template")
	imagesCmd.Flags().BoolVarP(&imagesQuiet, "quiet", "q", false, "only show image IDs")
	imagesCmd.Flags().BoolVarP(&imagesAll, "all", "a", false, "show all images (default hides intermediate images)")
}

func runImages(cmd *cobra.Command, args []string) error {
	// Initialize registry client
	registryClient := registry.New()

	// List options
	options := &registry.ListOptions{
		Filter: imagesFilter,
		All:    imagesAll,
	}

	// Get images
	images, err := registryClient.ListLocal(options)
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(images) == 0 {
		fmt.Println("No agent images found")
		fmt.Println("\n💡 Build an agent with: agent build -t my-agent .")
		fmt.Println("💡 Or pull an agent with: agent pull my-agent:latest")
		return nil
	}

	// Handle different output formats
	switch {
	case imagesQuiet:
		for _, image := range images {
			fmt.Println(image.ID[:12])
		}
	case imagesFormat == "json":
		return printImagesJSON(images)
	default:
		return printImagesTable(images)
	}

	return nil
}

func printImagesTable(images []registry.ImageInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Header
	fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")

	// Rows
	for _, image := range images {
		repository := image.Repository
		if repository == "" {
			repository = "<none>"
		}

		tag := image.Tag
		if tag == "" {
			tag = "<none>"
		}

		created := formatTime(image.Created)
		size := formatSize(image.Size)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			repository, tag, image.ID[:12], created, size)
	}

	return nil
}

func printImagesJSON(images []registry.ImageInfo) error {
	// Simple JSON output (in a real implementation, use json.Marshal)
	fmt.Println("[")
	for i, image := range images {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", image.ID)
		fmt.Printf("    \"repository\": \"%s\",\n", image.Repository)
		fmt.Printf("    \"tag\": \"%s\",\n", image.Tag)
		fmt.Printf("    \"created\": \"%s\",\n", image.Created.Format(time.RFC3339))
		fmt.Printf("    \"size\": %d\n", image.Size)
		if i < len(images)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")
	return nil
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Just now"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	} else {
		return t.Format("2006-01-02")
	}
}

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
