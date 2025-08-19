// Binary Uploader Tool
// Uploads agent CLI binaries to the registry for distribution
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pxkundu/agent-as-code/internal/api"
)

func main() {
	var (
		version      = flag.String("version", "", "Version to upload (required)")
		registry     = flag.String("registry", "https://api.myagentregistry.com", "Registry URL")
		token        = flag.String("token", "", "Auth token (or use AGENT_REGISTRY_TOKEN env)")
		binDir       = flag.String("bin-dir", "bin", "Directory containing binaries")
		allPlatforms = flag.Bool("all-platforms", false, "Upload all platform binaries")
		platform     = flag.String("platform", "", "Specific platform to upload")
		arch         = flag.String("arch", "", "Specific architecture to upload")
		dryRun       = flag.Bool("dry-run", false, "Show what would be uploaded")
	)

	flag.Parse()

	if *version == "" {
		fmt.Println("Error: version is required")
		flag.Usage()
		os.Exit(1)
	}

	authToken := *token
	if authToken == "" {
		authToken = os.Getenv("AGENT_REGISTRY_TOKEN")
		if authToken == "" {
			fmt.Println("Error: auth token required (use -token or AGENT_REGISTRY_TOKEN env)")
			os.Exit(1)
		}
	}

	fmt.Printf("🚀 Agent CLI Binary Uploader\n")
	fmt.Printf("Version: %s\n", *version)
	fmt.Printf("Registry: %s\n", *registry)

	if *dryRun {
		fmt.Println("🔍 DRY RUN - No actual uploads will be performed")
	}

	uploader := api.NewUploader(*registry, authToken, *version)

	var results []*api.UploadResult

	if *allPlatforms {
		fmt.Printf("📦 Uploading agent CLI binaries for all platforms from %s...\n", *binDir)
		if !*dryRun {
			results = uploader.UploadAllPlatforms(*binDir)
		} else {
			fmt.Println("Would upload all platform binaries")
			results = []*api.UploadResult{
				{Platform: "linux", Architecture: "amd64", Success: true},
				{Platform: "linux", Architecture: "arm64", Success: true},
				{Platform: "darwin", Architecture: "amd64", Success: true},
				{Platform: "darwin", Architecture: "arm64", Success: true},
				{Platform: "windows", Architecture: "amd64", Success: true},
				{Platform: "windows", Architecture: "arm64", Success: true},
			}
		}
	} else if *platform != "" && *arch != "" {
		binaryPath := filepath.Join(*binDir, fmt.Sprintf("agent-%s-%s", *platform, *arch))
		if *platform == "windows" {
			binaryPath += ".exe"
		}

		fmt.Printf("📦 Uploading agent CLI binary for %s/%s...\n", *platform, *arch)

		if !*dryRun {
			opts := api.UploadOptions{
				Platform:     *platform,
				Architecture: *arch,
				FilePath:     binaryPath,
			}
			result := uploader.UploadBinary(opts)
			results = []*api.UploadResult{result}
		} else {
			fmt.Printf("Would upload: %s\n", binaryPath)
			results = []*api.UploadResult{
				{Platform: *platform, Architecture: *arch, Success: true},
			}
		}
	} else {
		fmt.Println("Error: specify either --all-platforms or both --platform and --arch")
		os.Exit(1)
	}

	// Display results
	summary := api.GetUploadSummary(results)
	fmt.Print(summary)

	// Check for failures
	for _, result := range results {
		if !result.Success {
			log.Fatal("Some uploads failed")
		}
	}

	fmt.Println("\n✅ Agent CLI binaries are now available for installation!")
	fmt.Printf("Users can install via:\n")
	fmt.Printf("  pip install agent-as-code==%s\n", *version)
	fmt.Printf("  curl -L %s/install.sh | sh\n", *registry)
}
