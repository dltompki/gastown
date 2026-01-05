// Package cli provides abstraction for different AI CLI implementations.
package cli

import (
	"fmt"
)

// RoleType indicates whether a role is autonomous or interactive.
type RoleType string

const (
	// Autonomous roles (polecat, witness, refinery) need special handling
	// because they may be triggered externally without user input.
	Autonomous RoleType = "autonomous"

	// Interactive roles (mayor, crew) wait for user input.
	Interactive RoleType = "interactive"
)

// CLI defines the interface all AI CLIs must implement.
type CLI interface {
	// Command building
	BuildStartupCommand(role, actor, rigPath, prompt string) string
	BuildResumeCommand(sessionID string) string

	// Configuration management
	CreateConfiguration(workDir string, roleType RoleType) error

	// Session management
	IsProcessRunning(sessionName string) bool
	GetSessionIDEnvVar() string
	SupportsSessionResume() bool

	// Process monitoring
	GetProcessNames() []string

	// Feature support
	SupportsHooks() bool
	SupportsForkSession() bool

	// CLI identification
	GetType() string
}

// CLIError represents an error from CLI operations.
type CLIError struct {
	CLI     string
	Command string
	Err     error
}

func (e *CLIError) Error() string {
	return fmt.Sprintf("CLI %s command failed: %s: %v", e.CLI, e.Command, e.Err)
}

func (e *CLIError) Unwrap() error {
	return e.Err
}

// NewCLIError creates a new CLI error.
func NewCLIError(cli, command string, err error) *CLIError {
	return &CLIError{
		CLI:     cli,
		Command: command,
		Err:     err,
	}
}

// CLIFactory creates CLI instances based on type.
type CLIFactory interface {
	CreateCLI(cliType string) (CLI, error)
	RegisterCLI(cliType string, constructor func() CLI)
	GetSupportedTypes() []string
}
