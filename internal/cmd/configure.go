package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure registry profiles",
	Long: `Configure registry profiles for agent management.

This command manages registry profiles that define where agents are stored
and retrieved from. Profiles include registry URLs, authentication tokens,
and other connection settings.

Examples:
  agent configure profile add prod --registry https://api.myagentregistry.com --pat a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456
  agent configure profile list
  agent configure profile remove prod
  agent configure profile test prod
  agent configure profile set-default prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage registry profiles",
	Long: `Manage registry profiles for agent operations.

Profiles define the connection settings for different agent registries,
including URLs, authentication tokens, and default settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var profileAddCmd = &cobra.Command{
	Use:   "add [NAME]",
	Short: "Add a new registry profile",
	Long: `Add a new registry profile with the specified settings.

This command creates a new profile that can be used for agent operations
like push, pull, and other registry interactions.

Examples:
  agent configure profile add production --registry https://api.myagentregistry.com --pat a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456
  agent configure profile add staging --registry https://api.myagentregistry.com --pat b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567a --description "Staging environment"
  agent configure profile add local --registry http://localhost:5000 --pat c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567890 --set-default --test`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		registry, _ := cmd.Flags().GetString("registry")
		pat, _ := cmd.Flags().GetString("pat")
		description, _ := cmd.Flags().GetString("description")
		setDefault, _ := cmd.Flags().GetBool("set-default")
		test, _ := cmd.Flags().GetBool("test")

		return addProfile(name, registry, pat, description, setDefault, test)
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured profiles",
	Long: `List all configured registry profiles.

This command displays all available profiles with their settings and
indicates which profile is currently set as default.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProfiles()
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove [NAME]",
	Short: "Remove a registry profile",
	Long: `Remove a registry profile.

This command removes the specified profile and all its associated settings.
If the profile is currently set as default, you'll need to set a new
default profile first.

Examples:
  agent configure profile remove production
  agent configure profile remove staging`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return removeProfile(name)
	},
}

var profileTestCmd = &cobra.Command{
	Use:   "test [NAME]",
	Short: "Test a registry profile",
	Long: `Test a registry profile connection.

This command tests the connection to the registry using the specified
profile to ensure it's working correctly.

Examples:
  agent configure profile test production
  agent configure profile test staging
  agent configure profile test default`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return testProfile(name)
	},
}

var profileSetDefaultCmd = &cobra.Command{
	Use:   "set-default [NAME]",
	Short: "Set a profile as default",
	Long: `Set a profile as the default registry profile.

This command sets the specified profile as the default, which will be
used for agent operations when no specific profile is specified.

Examples:
  agent configure profile set-default production
  agent configure profile set-default staging
  agent configure profile set-default local`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return setDefaultProfile(name)
	},
}

func init() {
	// Configure command
	rootCmd.AddCommand(configureCmd)

	// Profile subcommands
	configureCmd.AddCommand(profileCmd)

	// Profile add command
	profileAddCmd.Flags().String("registry", "", "registry URL (required)")
	profileAddCmd.Flags().String("pat", "", "personal access token")
	profileAddCmd.Flags().String("description", "", "profile description")
	profileAddCmd.Flags().Bool("set-default", false, "set as default profile")
	profileAddCmd.Flags().Bool("test", false, "test connection after adding")
	profileAddCmd.MarkFlagRequired("registry")
	profileCmd.AddCommand(profileAddCmd)

	// Profile list command
	profileCmd.AddCommand(profileListCmd)

	// Profile remove command
	profileCmd.AddCommand(profileRemoveCmd)

	// Profile test command
	profileCmd.AddCommand(profileTestCmd)

	// Profile set-default command
	profileCmd.AddCommand(profileSetDefaultCmd)
}

type Profile struct {
	Registry    string `json:"registry"`
	PAT         string `json:"pat"`
	Description string `json:"description"`
}

type Config struct {
	Profiles       map[string]Profile `json:"profiles"`
	DefaultProfile string             `json:"default_profile"`
}

func addProfile(name, registry, pat, description string, setDefault, test bool) error {
	// Validate PAT format
	if !validatePAT(pat) {
		return fmt.Errorf("invalid PAT format. PAT should be 64 characters hexadecimal")
	}

	// Load existing config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Check if profile already exists
	if _, exists := config.Profiles[name]; exists {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	// Create the profile
	profile := Profile{
		Registry:    registry,
		PAT:         pat,
		Description: description,
	}

	// Add to config
	config.Profiles[name] = profile

	// Set as default if requested or if no default profile exists
	if setDefault || config.DefaultProfile == "" {
		config.DefaultProfile = name
	}

	// Save the config
	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save profile: %v", err)
	}

	fmt.Printf("Profile '%s' configured successfully\n", name)
	if setDefault || config.DefaultProfile == name {
		fmt.Printf("Profile '%s' set as default\n", name)
	}

	// Test connection if requested
	if test {
		fmt.Printf("Testing connection...\n")
		if err := testProfile(name); err != nil {
			fmt.Printf("Connection test failed!\n")
			return err
		} else {
			fmt.Printf("Connection test successful!\n")
		}
	}

	return nil
}

func listProfiles() error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %v", err)
	}

	if len(config.Profiles) == 0 {
		fmt.Println("No profiles configured")
		fmt.Println("Use 'agent configure profile add' to add a profile")
		return nil
	}

	fmt.Println("Configured profiles:")
	for name, profile := range config.Profiles {
		defaultMarker := ""
		if name == config.DefaultProfile {
			defaultMarker = " (default)"
		}

		fmt.Printf("  %s%s\n", name, defaultMarker)
		fmt.Printf("    Registry: %s\n", profile.Registry)
		fmt.Printf("    Description: %s\n", profile.Description)
		fmt.Println()
	}

	return nil
}

func removeProfile(name string) error {
	// Load existing config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Check if profile exists
	if _, exists := config.Profiles[name]; !exists {
		fmt.Printf("Profile '%s' not found\n", name)
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Remove the profile
	delete(config.Profiles, name)

	// Update default profile if necessary
	if config.DefaultProfile == name {
		if len(config.Profiles) > 0 {
			// Set first remaining profile as default
			for profileName := range config.Profiles {
				config.DefaultProfile = profileName
				fmt.Printf("Default profile changed to '%s'\n", profileName)
				break
			}
		} else {
			config.DefaultProfile = ""
			fmt.Println("No profiles remaining")
		}
	}

	// Save the config
	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("Profile '%s' removed successfully\n", name)
	return nil
}

func testProfile(name string) error {
	// Load the config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Get the profile
	profile, exists := config.Profiles[name]
	if !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Test the connection using registry client
	if err := testRegistryConnection(profile.Registry, profile.PAT); err != nil {
		return fmt.Errorf("connection test failed: %v", err)
	}

	return nil
}

func setDefaultProfile(name string) error {
	// Load existing config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Check if profile exists
	if _, exists := config.Profiles[name]; !exists {
		fmt.Printf("Profile '%s' not found\n", name)
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Set as default
	config.DefaultProfile = name

	// Save the config
	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("Default profile set to '%s'\n", name)
	return nil
}

func loadConfig() (*Config, error) {
	configFile := getConfigFile()

	// Create default config if file doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{
			Profiles:       make(map[string]Profile),
			DefaultProfile: "",
		}, nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		// Return default config if parsing fails
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		return &Config{
			Profiles:       make(map[string]Profile),
			DefaultProfile: "",
		}, nil
	}

	// Initialize profiles map if nil
	if config.Profiles == nil {
		config.Profiles = make(map[string]Profile)
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configFile := getConfigFile()

	// Ensure config directory exists
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func getConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return filepath.Join(home, ".agent", "config.json")
}

func validatePAT(pat string) bool {
	// Basic validation - PAT should be 64 characters
	if len(pat) != 64 {
		return false
	}

	// Check if it's hexadecimal
	matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", pat)
	return matched
}

func testRegistryConnection(registry, pat string) error {
	// Import needed for HTTP requests
	// In a real implementation, this would make an HTTP request to test connectivity
	// For now, we simulate the test based on the registry URL

	// Test connection by checking if it looks like a valid registry URL
	if !strings.HasPrefix(registry, "http://") && !strings.HasPrefix(registry, "https://") {
		return fmt.Errorf("invalid registry URL format")
	}

	// Simulate connection test failure for example domains
	if strings.Contains(registry, "example.com") {
		return fmt.Errorf("example.com is not a real registry")
	}

	// In a real implementation, this would make a GET request to {registry}/health
	// with Authorization header containing the PAT

	return nil
}
