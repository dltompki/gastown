package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// KiroCLI implements the CLI interface for Kiro CLI.
type KiroCLI struct {
	command string
}

// NewKiroCLI creates a new Kiro CLI instance.
func NewKiroCLI() CLI {
	return &KiroCLI{
		command: "kiro-cli",
	}
}

// BuildStartupCommand builds a startup command for Kiro CLI.
func (k *KiroCLI) BuildStartupCommand(role, actor, rigPath, prompt string) string {
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

	cmd := "export " + strings.Join(exports, " ") + " && " + k.command + " chat"

	// Add prompt if provided
	if prompt != "" {
		cmd += " " + shellescape(prompt)
	}

	return cmd
}

// BuildResumeCommand builds a resume command for Kiro CLI.
// Kiro CLI uses automatic session persistence, so this returns empty.
func (k *KiroCLI) BuildResumeCommand(sessionID string) string {
	// Kiro CLI handles session persistence automatically
	return ""
}

// CreateConfiguration creates Kiro CLI agent configuration.
func (k *KiroCLI) CreateConfiguration(workDir string, roleType RoleType) error {
	return k.createKiroAgentConfig(workDir, roleType)
}

// createKiroAgentConfig creates a Kiro CLI agent configuration file.
func (k *KiroCLI) createKiroAgentConfig(workDir string, roleType RoleType) error {
	kiroDir := filepath.Join(workDir, ".kiro")
	agentsDir := filepath.Join(kiroDir, "agents")
	configPath := filepath.Join(agentsDir, "gastown.json")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	// Create directories
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("creating .kiro/agents directory: %w", err)
	}

	// Create agent configuration based on role type
	config := k.buildAgentConfig(roleType)

	// Write configuration file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling agent config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("writing agent config: %w", err)
	}

	return nil
}

// buildAgentConfig builds the agent configuration for Kiro CLI.
func (k *KiroCLI) buildAgentConfig(roleType RoleType) map[string]interface{} {
	config := map[string]interface{}{
		"name":        "mayor",
		"description": "Gastown agent configuration",
		"tools":       []string{"*"}, // Allow all tools
		"allowedTools": []string{
			"fs_read",
			"fs_write", 
			"execute_bash",
			"todo_list",
		},
	}

	// Add role-specific configuration
	if roleType == Autonomous {
		// Autonomous roles need additional tool permissions
		config["toolsSettings"] = map[string]interface{}{
			"execute_bash": map[string]interface{}{
				"autoAllowReadonly": true,
			},
		}
	}

	return config
}

// IsProcessRunning checks if Kiro CLI is running.
func (k *KiroCLI) IsProcessRunning(sessionName string) bool {
	// Use tmux integration to check if Kiro CLI process is running
	// This would check for "kiro-cli" process in the session
	return false // Placeholder - will be integrated with tmux package
}

// GetSessionIDEnvVar returns empty as Kiro CLI doesn't use session ID env vars.
func (k *KiroCLI) GetSessionIDEnvVar() string {
	return ""
}

// SupportsSessionResume returns false as Kiro CLI uses automatic persistence.
func (k *KiroCLI) SupportsSessionResume() bool {
	return false
}

// GetProcessNames returns the process names to monitor for Kiro CLI.
func (k *KiroCLI) GetProcessNames() []string {
	return []string{"kiro-cli"}
}

// SupportsHooks returns true as Kiro CLI supports hooks via agent configuration.
func (k *KiroCLI) SupportsHooks() bool {
	return true
}

// SupportsForkSession returns false as Kiro CLI doesn't support fork sessions.
func (k *KiroCLI) SupportsForkSession() bool {
	return false
}

// GetType returns the CLI type identifier.
func (k *KiroCLI) GetType() string {
	return "kiro"
}
