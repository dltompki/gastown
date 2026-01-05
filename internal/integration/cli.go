// Package integration provides integration between config and CLI packages.
package integration

import (
	"path/filepath"

	"github.com/steveyegge/gastown/internal/cli"
	"github.com/steveyegge/gastown/internal/config"
)

// BuildCrewStartupCommandWithCLI builds a crew startup command using CLI abstraction.
func BuildCrewStartupCommandWithCLI(rigName, crewName, rigPath, prompt string) string {
	// Derive town root from rig path
	townRoot := filepath.Dir(rigPath)
	
	// Resolve CLI type and create CLI instance
	cliType := config.ResolveCLIType(townRoot)
	cliInstance, err := cli.CreateCLIFromConfig(cliType)
	if err != nil {
		// Fallback to Claude Code for backwards compatibility
		cliInstance = cli.NewClaudeCodeCLI()
	}

	// Build actor string
	bdActor := rigName + "/crew/" + crewName
	
	// Use CLI to build startup command
	return cliInstance.BuildStartupCommand("crew", bdActor, rigPath, prompt)
}

// BuildPolecatStartupCommandWithCLI builds a polecat startup command using CLI abstraction.
func BuildPolecatStartupCommandWithCLI(rigName, polecatName, rigPath, prompt string) string {
	// Derive town root from rig path
	townRoot := filepath.Dir(rigPath)
	
	// Resolve CLI type and create CLI instance
	cliType := config.ResolveCLIType(townRoot)
	cliInstance, err := cli.CreateCLIFromConfig(cliType)
	if err != nil {
		// Fallback to Claude Code for backwards compatibility
		cliInstance = cli.NewClaudeCodeCLI()
	}

	// Build actor string
	bdActor := rigName + "/polecats/" + polecatName
	
	// Use CLI to build startup command
	return cliInstance.BuildStartupCommand("polecat", bdActor, rigPath, prompt)
}

// BuildAgentStartupCommandWithCLI builds an agent startup command using CLI abstraction.
func BuildAgentStartupCommandWithCLI(role, bdActor, rigPath, prompt string) string {
	// Derive town root from rig path
	townRoot := filepath.Dir(rigPath)
	
	// Resolve CLI type and create CLI instance
	cliType := config.ResolveCLIType(townRoot)
	cliInstance, err := cli.CreateCLIFromConfig(cliType)
	if err != nil {
		// Fallback to Claude Code for backwards compatibility
		cliInstance = cli.NewClaudeCodeCLI()
	}

	// Use CLI to build startup command
	return cliInstance.BuildStartupCommand(role, bdActor, rigPath, prompt)
}

// CreateCLIConfiguration creates CLI-specific configuration for a role.
func CreateCLIConfiguration(workDir string, roleType string, townRoot string) error {
	// Resolve CLI type
	cliType := config.ResolveCLIType(townRoot)
	cliInstance, err := cli.CreateCLIFromConfig(cliType)
	if err != nil {
		// Fallback to Claude Code
		cliInstance = cli.NewClaudeCodeCLI()
	}

	// Convert string to RoleType
	var rt cli.RoleType
	if roleType == "autonomous" {
		rt = cli.Autonomous
	} else {
		rt = cli.Interactive
	}

	return cliInstance.CreateConfiguration(workDir, rt)
}

// GetCLIProcessNames returns process names for the configured CLI.
func GetCLIProcessNames(townRoot string) []string {
	cliType := config.ResolveCLIType(townRoot)
	cliInstance, err := cli.CreateCLIFromConfig(cliType)
	if err != nil {
		// Fallback to Claude Code
		cliInstance = cli.NewClaudeCodeCLI()
	}

	return cliInstance.GetProcessNames()
}
