package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/steveyegge/gastown/internal/claude"
)

// ClaudeCodeCLI implements the CLI interface for Claude Code.
type ClaudeCodeCLI struct {
	command string
	args    []string
}

// NewClaudeCodeCLI creates a new Claude Code CLI instance.
func NewClaudeCodeCLI() CLI {
	return &ClaudeCodeCLI{
		command: "claude",
		args:    []string{"--dangerously-skip-permissions"},
	}
}

// BuildStartupCommand builds a full startup command with environment exports.
func (c *ClaudeCodeCLI) BuildStartupCommand(role, actor, rigPath, prompt string) string {
	envVars := map[string]string{
		"GT_ROLE":         role,
		"BD_ACTOR":        actor,
		"GIT_AUTHOR_NAME": actor,
	}

	// Add role-specific environment variables
	switch role {
	case "crew":
		if parts := strings.Split(actor, "/"); len(parts) >= 3 {
			envVars["GT_RIG"] = parts[0]
			envVars["GT_CREW"] = parts[2]
		}
	case "polecat":
		if parts := strings.Split(actor, "/"); len(parts) >= 3 {
			envVars["GT_RIG"] = parts[0]
			envVars["GT_POLECAT"] = parts[2]
		}
	}

	// Build environment export prefix
	var exports []string
	for k, v := range envVars {
		exports = append(exports, fmt.Sprintf("%s=%s", k, v))
	}

	// Sort for deterministic output
	sort.Strings(exports)

	cmd := "export " + strings.Join(exports, " ") + " && " + c.command
	if len(c.args) > 0 {
		cmd += " " + strings.Join(c.args, " ")
	}

	// Add prompt if provided
	if prompt != "" {
		cmd += " " + shellescape(prompt)
	}

	return cmd
}

// BuildResumeCommand builds a command to resume a Claude session.
func (c *ClaudeCodeCLI) BuildResumeCommand(sessionID string) string {
	if sessionID == "" {
		return ""
	}

	cmd := c.command
	if len(c.args) > 0 {
		cmd += " " + strings.Join(c.args, " ")
	}
	cmd += " --resume " + sessionID

	return cmd
}

// CreateConfiguration creates Claude Code configuration files.
func (c *ClaudeCodeCLI) CreateConfiguration(workDir string, roleType RoleType) error {
	return claude.EnsureSettings(workDir, claude.RoleType(roleType))
}

// IsProcessRunning checks if Claude is running in the given session.
func (c *ClaudeCodeCLI) IsProcessRunning(sessionName string) bool {
	// Use tmux integration to check if Claude process is running
	// This integrates with existing tmux.IsClaudeRunning logic
	return false // Placeholder - will be integrated with tmux package
}

// GetSessionIDEnvVar returns the environment variable for Claude session IDs.
func (c *ClaudeCodeCLI) GetSessionIDEnvVar() string {
	return "CLAUDE_SESSION_ID"
}

// SupportsSessionResume returns true as Claude supports session resumption.
func (c *ClaudeCodeCLI) SupportsSessionResume() bool {
	return true
}

// GetProcessNames returns the process names to monitor for Claude.
func (c *ClaudeCodeCLI) GetProcessNames() []string {
	return []string{"claude", "claude-code"}
}

// SupportsHooks returns true as Claude supports the hooks system.
func (c *ClaudeCodeCLI) SupportsHooks() bool {
	return true
}

// SupportsForkSession returns true as Claude supports fork session.
func (c *ClaudeCodeCLI) SupportsForkSession() bool {
	return true
}

// GetType returns the CLI type identifier.
func (c *ClaudeCodeCLI) GetType() string {
	return "claude"
}

// shellescape provides basic shell escaping for arguments.
func shellescape(s string) string {
	if strings.ContainsAny(s, " \t\n\r\"'\\$`") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}
