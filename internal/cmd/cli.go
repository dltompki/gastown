package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/steveyegge/gastown/internal/cli"
	"github.com/steveyegge/gastown/internal/config"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Manage CLI configuration",
	Long:  `Configure which AI CLI to use (claude or kiro).`,
}

var cliSetCmd = &cobra.Command{
	Use:   "set [claude|kiro]",
	Short: "Set the default CLI type",
	Long:  `Set the default AI CLI to use for all agents (claude or kiro).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCLISet,
}

var cliGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current CLI type",
	Long:  `Show the currently configured AI CLI type.`,
	RunE:  runCLIGet,
}

var cliListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available CLI types",
	Long:  `List all supported AI CLI types.`,
	RunE:  runCLIList,
}

func init() {
	cliCmd.AddCommand(cliSetCmd)
	cliCmd.AddCommand(cliGetCmd)
	cliCmd.AddCommand(cliListCmd)
	rootCmd.AddCommand(cliCmd)
}

func runCLISet(cmd *cobra.Command, args []string) error {
	cliType := args[0]

	// Validate CLI type
	factory := cli.GetGlobalFactory()
	supportedTypes := factory.GetSupportedTypes()
	
	valid := false
	for _, supported := range supportedTypes {
		if cliType == supported {
			valid = true
			break
		}
	}
	
	if !valid {
		return fmt.Errorf("unsupported CLI type: %s (supported: %v)", cliType, supportedTypes)
	}

	// Get town root
	townRoot, err := findTownRoot()
	if err != nil {
		return fmt.Errorf("finding town root: %w", err)
	}

	// Load or create town settings
	settingsPath := config.TownSettingsPath(townRoot)
	settings, err := config.LoadOrCreateTownSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("loading town settings: %w", err)
	}

	// Update CLI type
	settings.DefaultCLI = cliType

	// Save settings
	if err := saveTownSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("saving town settings: %w", err)
	}

	fmt.Printf("Set default CLI to: %s\n", cliType)
	return nil
}

func runCLIGet(cmd *cobra.Command, args []string) error {
	// Get town root
	townRoot, err := findTownRoot()
	if err != nil {
		return fmt.Errorf("finding town root: %w", err)
	}

	// Resolve CLI type
	cliType := config.ResolveCLIType(townRoot)
	fmt.Printf("Current CLI: %s\n", cliType)
	return nil
}

func runCLIList(cmd *cobra.Command, args []string) error {
	factory := cli.GetGlobalFactory()
	supportedTypes := factory.GetSupportedTypes()
	
	fmt.Println("Supported CLI types:")
	for _, cliType := range supportedTypes {
		fmt.Printf("  %s\n", cliType)
	}
	return nil
}

// saveTownSettings saves town settings to a file.
func saveTownSettings(path string, settings *config.TownSettings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding settings: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}

	return nil
}
