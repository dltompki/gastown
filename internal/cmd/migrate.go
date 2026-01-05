package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steveyegge/gastown/internal/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate configuration for CLI compatibility",
	Long:  `Migrate existing Gastown configuration to support CLI compatibility.`,
}

var migrateCliCmd = &cobra.Command{
	Use:   "cli [claude|kiro]",
	Short: "Migrate to a different CLI type",
	Long:  `Migrate existing configuration to use a different AI CLI.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateCLI,
}

func init() {
	migrateCmd.AddCommand(migrateCliCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runMigrateCLI(cmd *cobra.Command, args []string) error {
	targetCLI := args[0]
	
	// Validate target CLI
	if targetCLI != "claude" && targetCLI != "kiro" {
		return fmt.Errorf("unsupported CLI type: %s (supported: claude, kiro)", targetCLI)
	}

	// Get town root
	townRoot, err := findTownRoot()
	if err != nil {
		return fmt.Errorf("finding town root: %w", err)
	}

	// Load current settings
	settingsPath := config.TownSettingsPath(townRoot)
	settings, err := config.LoadOrCreateTownSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("loading town settings: %w", err)
	}

	currentCLI := settings.DefaultCLI
	if currentCLI == "" {
		currentCLI = "claude" // default
	}

	if currentCLI == targetCLI {
		fmt.Printf("Already using %s CLI\n", targetCLI)
		return nil
	}

	fmt.Printf("Migrating from %s to %s CLI...\n", currentCLI, targetCLI)

	// Update CLI type
	settings.DefaultCLI = targetCLI

	// Save updated settings
	if err := saveTownSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("saving town settings: %w", err)
	}

	// Check for existing sessions that might need restart
	if err := checkActiveSessions(townRoot, targetCLI); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	fmt.Printf("Successfully migrated to %s CLI\n", targetCLI)
	fmt.Println("Note: Restart active sessions to use the new CLI")
	
	return nil
}

// checkActiveSessions warns about active sessions that need restart.
func checkActiveSessions(townRoot, targetCLI string) error {
	// This is a simplified check - in practice would use tmux integration
	// to detect active Gastown sessions
	
	// Check if there are any .claude directories that might conflict with Kiro
	if targetCLI == "kiro" {
		claudeDirs := []string{}
		
		// Walk through rigs looking for .claude directories
		entries, err := os.ReadDir(townRoot)
		if err != nil {
			return nil // non-fatal
		}
		
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			
			claudeDir := filepath.Join(townRoot, entry.Name(), ".claude")
			if _, err := os.Stat(claudeDir); err == nil {
				claudeDirs = append(claudeDirs, entry.Name())
			}
		}
		
		if len(claudeDirs) > 0 {
			return fmt.Errorf("found Claude configuration directories in rigs: %v", claudeDirs)
		}
	}
	
	return nil
}
