package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/steveyegge/gastown/internal/beads"
	"github.com/steveyegge/gastown/internal/claude"
	"github.com/steveyegge/gastown/internal/config"
	"github.com/steveyegge/gastown/internal/deps"
	"github.com/steveyegge/gastown/internal/formula"
	"github.com/steveyegge/gastown/internal/session"
	"github.com/steveyegge/gastown/internal/style"
	"github.com/steveyegge/gastown/internal/templates"
	"github.com/steveyegge/gastown/internal/workspace"
)

var (
	installForce      bool
	installName       string
	installOwner      string
	installPublicName string
	installNoBeads    bool
	installGit        bool
	installGitHub     string
	installPublic     bool
	installCLI        string
)

var installCmd = &cobra.Command{
	Use:     "install [path]",
	GroupID: GroupWorkspace,
	Short:   "Create a new Gas Town HQ (workspace)",
	Long: `Create a new Gas Town HQ at the specified path.

The HQ (headquarters) is the top-level directory where Gas Town is installed -
the root of your workspace where all rigs and agents live. It contains:
  - CLAUDE.md            Mayor role context (Mayor runs from HQ root)
  - mayor/               Mayor config, state, and rig registry
  - .beads/              Town-level beads DB (hq-* prefix for mayor mail)

If path is omitted, uses the current directory.

See docs/hq.md for advanced HQ configurations including beads
redirects, multi-system setups, and HQ templates.

Examples:
  gt install ~/gt                              # Create HQ at ~/gt
  gt install . --name my-workspace             # Initialize current dir
  gt install ~/gt --no-beads                   # Skip .beads/ initialization
  gt install ~/gt --git                        # Also init git with .gitignore
  gt install ~/gt --github=user/repo           # Create private GitHub repo (default)
  gt install ~/gt --github=user/repo --public  # Create public GitHub repo`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInstall,
}

func init() {
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "Overwrite existing HQ")
	installCmd.Flags().StringVarP(&installName, "name", "n", "", "Town name (defaults to directory name)")
	installCmd.Flags().StringVar(&installOwner, "owner", "", "Owner email for entity identity (defaults to git config user.email)")
	installCmd.Flags().StringVar(&installPublicName, "public-name", "", "Public display name (defaults to town name)")
	installCmd.Flags().BoolVar(&installNoBeads, "no-beads", false, "Skip town beads initialization")
	installCmd.Flags().BoolVar(&installGit, "git", false, "Initialize git with .gitignore")
	installCmd.Flags().StringVar(&installGitHub, "github", "", "Create GitHub repo (format: owner/repo, private by default)")
	installCmd.Flags().BoolVar(&installPublic, "public", false, "Make GitHub repo public (use with --github)")
	installCmd.Flags().StringVar(&installCLI, "cli", "claude", "AI CLI to use (claude or kiro)")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Determine target path
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	// Expand ~ and resolve to absolute path
	if targetPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}
		targetPath = filepath.Join(home, targetPath[1:])
	}

	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	// Determine town name
	townName := installName
	if townName == "" {
		townName = filepath.Base(absPath)
	}

	// Check if already a workspace
	if isWS, _ := workspace.IsWorkspace(absPath); isWS && !installForce {
		return fmt.Errorf("directory is already a Gas Town HQ (use --force to reinitialize)")
	}

	// Check if inside an existing workspace
	if existingRoot, _ := workspace.Find(absPath); existingRoot != "" && existingRoot != absPath {
		style.PrintWarning("Creating HQ inside existing workspace at %s", existingRoot)
	}

	// Ensure beads (bd) is available before proceeding
	if !installNoBeads {
		if err := deps.EnsureBeads(true); err != nil {
			return fmt.Errorf("beads dependency check failed: %w", err)
		}
	}

	fmt.Printf("%s Creating Gas Town HQ at %s\n\n",
		style.Bold.Render("🏭"), style.Dim.Render(absPath))

	// Create directory structure
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Create mayor directory (holds config, state, and mail)
	mayorDir := filepath.Join(absPath, "mayor")
	if err := os.MkdirAll(mayorDir, 0755); err != nil {
		return fmt.Errorf("creating mayor directory: %w", err)
	}
	fmt.Printf("   ✓ Created mayor/\n")

	// Determine owner (defaults to git user.email)
	owner := installOwner
	if owner == "" {
		out, err := exec.Command("git", "config", "user.email").Output()
		if err == nil {
			owner = strings.TrimSpace(string(out))
		}
	}

	// Determine public name (defaults to town name)
	publicName := installPublicName
	if publicName == "" {
		publicName = townName
	}

	// Create town.json in mayor/
	townConfig := &config.TownConfig{
		Type:       "town",
		Version:    config.CurrentTownVersion,
		Name:       townName,
		Owner:      owner,
		PublicName: publicName,
		CreatedAt:  time.Now(),
	}
	townPath := filepath.Join(mayorDir, "town.json")
	if err := config.SaveTownConfig(townPath, townConfig); err != nil {
		return fmt.Errorf("writing town.json: %w", err)
	}
	fmt.Printf("   ✓ Created mayor/town.json\n")

	// Create rigs.json in mayor/
	rigsConfig := &config.RigsConfig{
		Version: config.CurrentRigsVersion,
		Rigs:    make(map[string]config.RigEntry),
	}
	rigsPath := filepath.Join(mayorDir, "rigs.json")
	if err := config.SaveRigsConfig(rigsPath, rigsConfig); err != nil {
		return fmt.Errorf("writing rigs.json: %w", err)
	}
	fmt.Printf("   ✓ Created mayor/rigs.json\n")

	// Create town settings with CLI configuration
	townSettings := config.NewTownSettings()
	townSettings.DefaultCLI = installCLI
	
	settingsDir := filepath.Join(absPath, "settings")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return fmt.Errorf("creating settings directory: %w", err)
	}
	
	settingsPath := config.TownSettingsPath(absPath)
	if err := saveTownSettings(settingsPath, townSettings); err != nil {
		return fmt.Errorf("saving town settings: %w", err)
	}
	fmt.Printf("   ✓ Created settings/config.json (CLI: %s)\n", installCLI)

	// Create Mayor CLI configuration at HQ root (Mayor runs from there)
	if err := createMayorCLIConfig(absPath, absPath); err != nil {
		fmt.Printf("   %s Could not create CLI configuration: %v\n", style.Dim.Render("⚠"), err)
	} else {
		fmt.Printf("   ✓ Created CLI configuration\n")
	}

	// Ensure Mayor has CLI settings with appropriate hooks.
	// This ensures gt prime runs on CLI startup, which outputs the Mayor
	// delegation protocol - critical for preventing direct implementation.
	if err := ensureMayorCLISettings(absPath); err != nil {
		fmt.Printf("   %s Could not create CLI settings: %v\n", style.Dim.Render("⚠"), err)
	} else {
		fmt.Printf("   ✓ Created CLI settings\n")
	}

	// Initialize town-level beads database (optional)
	// Town beads (hq- prefix) stores mayor mail, cross-rig coordination, and handoffs.
	// Rig beads are separate and have their own prefixes.
	if !installNoBeads {
		if err := initTownBeads(absPath); err != nil {
			fmt.Printf("   %s Could not initialize town beads: %v\n", style.Dim.Render("⚠"), err)
		} else {
			fmt.Printf("   ✓ Initialized .beads/ (town-level beads with hq- prefix)\n")

			// Provision embedded formulas to .beads/formulas/
			if count, err := formula.ProvisionFormulas(absPath); err != nil {
				// Non-fatal: formulas are optional, just convenience
				fmt.Printf("   %s Could not provision formulas: %v\n", style.Dim.Render("⚠"), err)
			} else if count > 0 {
				fmt.Printf("   ✓ Provisioned %d formulas\n", count)
			}
		}

		// Create town-level agent beads (Mayor, Deacon) and role beads.
		// These use hq- prefix and are stored in town beads for cross-rig coordination.
		if err := initTownAgentBeads(absPath); err != nil {
			fmt.Printf("   %s Could not create town-level agent beads: %v\n", style.Dim.Render("⚠"), err)
		}
	}

	// Detect and save overseer identity
	overseer, err := config.DetectOverseer(absPath)
	if err != nil {
		fmt.Printf("   %s Could not detect overseer identity: %v\n", style.Dim.Render("⚠"), err)
	} else {
		overseerPath := config.OverseerConfigPath(absPath)
		if err := config.SaveOverseerConfig(overseerPath, overseer); err != nil {
			fmt.Printf("   %s Could not save overseer config: %v\n", style.Dim.Render("⚠"), err)
		} else {
			fmt.Printf("   ✓ Detected overseer: %s (via %s)\n", overseer.FormatOverseerIdentity(), overseer.Source)
		}
	}

	// Provision CLI-specific commands
	cliType := config.ResolveCLIType(absPath)
	if cliType == "claude" {
		// Provision town-level slash commands (.claude/commands/)
		// All agents inherit these via Claude's directory traversal - no per-workspace copies needed.
		if err := templates.ProvisionCommands(absPath); err != nil {
			fmt.Printf("   %s Could not provision slash commands: %v\n", style.Dim.Render("⚠"), err)
		} else {
			fmt.Printf("   ✓ Created .claude/commands/ (slash commands for all agents)\n")
		}
	} else if cliType == "kiro" {
		// Provision Kiro prompts (.kiro/prompts/)
		// Provides equivalent functionality to Claude's slash commands via @ prompts
		if err := templates.ProvisionPrompts(absPath); err != nil {
			fmt.Printf("   %s Could not provision Kiro prompts: %v\n", style.Dim.Render("⚠"), err)
		} else {
			fmt.Printf("   ✓ Created .kiro/prompts/ (@ prompts for Gastown workflows)\n")
		}
	}

	// Initialize git if requested (--git or --github implies --git)
	if installGit || installGitHub != "" {
		fmt.Println()
		if err := InitGitForHarness(absPath, installGitHub, !installPublic); err != nil {
			return fmt.Errorf("git initialization failed: %w", err)
		}
	}

	fmt.Printf("\n%s HQ created successfully!\n", style.Bold.Render("✓"))
	fmt.Println()
	fmt.Println("Next steps:")
	step := 1
	if !installGit && installGitHub == "" {
		fmt.Printf("  %d. Initialize git: %s\n", step, style.Dim.Render("gt git-init"))
		step++
	}
	fmt.Printf("  %d. Add a rig: %s\n", step, style.Dim.Render("gt rig add <name> <git-url>"))
	step++
	fmt.Printf("  %d. Enter the Mayor's office: %s\n", step, style.Dim.Render("gt mayor attach"))

	return nil
}

func createMayorCLIConfig(hqRoot, townRoot string) error {
	// Resolve CLI type for this town
	cliType := config.ResolveCLIType(townRoot)
	
	// Create CLI-specific configuration
	switch cliType {
	case "kiro":
		return createMayorKiroConfig(hqRoot, townRoot)
	default:
		return createMayorCLAUDEmd(hqRoot, townRoot)
	}
}

func createMayorKiroConfig(hqRoot, townRoot string) error {
	// Create .kiro/agents/ directory
	kiroDir := filepath.Join(hqRoot, ".kiro", "agents")
	if err := os.MkdirAll(kiroDir, 0755); err != nil {
		return fmt.Errorf("creating .kiro/agents directory: %w", err)
	}

	// Create mayor agent configuration
	config := map[string]interface{}{
		"name":        "mayor",
		"description": "Gastown Mayor agent - global coordinator",
		"tools":       []string{"*"},
		"allowedTools": []string{
			"fs_read",
			"fs_write",
			"execute_bash",
		},
		"resources": []string{
			"file://MAYOR.md", // Kiro uses MAYOR.md instead of CLAUDE.md
		},
	}

	configPath := filepath.Join(kiroDir, "mayor.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling mayor config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing mayor config: %w", err)
	}

	// Create MAYOR.md (Kiro equivalent of CLAUDE.md)
	return createMayorMD(hqRoot, townRoot, "MAYOR.md")
}

func createMayorMD(hqRoot, townRoot, filename string) error {
	tmpl, err := templates.New()
	if err != nil {
		return err
	}

	// Get town name for session names
	townName, _ := workspace.GetTownName(townRoot)

	data := templates.RoleData{
		Role:          "mayor",
		TownRoot:      townRoot,
		TownName:      townName,
		WorkDir:       hqRoot,
		MayorSession:  session.MayorSessionName(),
		DeaconSession: session.DeaconSessionName(),
	}

	content, err := tmpl.RenderRole("mayor", data)
	if err != nil {
		return err
	}

	mdPath := filepath.Join(hqRoot, filename)
	return os.WriteFile(mdPath, []byte(content), 0644)
}

func ensureMayorCLISettings(hqRoot string) error {
	// Resolve CLI type
	cliType := config.ResolveCLIType(hqRoot)
	
	switch cliType {
	case "kiro":
		// Kiro CLI settings are handled in agent configuration
		return nil
	default:
		// Use existing Claude settings
		return claude.EnsureSettingsForRole(hqRoot, "mayor")
	}
}

func createMayorCLAUDEmd(hqRoot, townRoot string) error {
	return createMayorMD(hqRoot, townRoot, "CLAUDE.md")
}

func writeJSON(path string, data interface{}) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// initTownBeads initializes town-level beads database using bd init.
// Town beads use the "hq-" prefix for mayor mail and cross-rig coordination.
func initTownBeads(townPath string) error {
	// Run: bd init --prefix hq
	cmd := exec.Command("bd", "init", "--prefix", "hq")
	cmd.Dir = townPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if beads is already initialized
		if strings.Contains(string(output), "already initialized") {
			// Already initialized - still need to ensure fingerprint exists
		} else {
			return fmt.Errorf("bd init failed: %s", strings.TrimSpace(string(output)))
		}
	}

	// Ensure database has repository fingerprint (GH #25).
	// This is idempotent - safe on both new and legacy (pre-0.17.5) databases.
	// Without fingerprint, the bd daemon fails to start silently.
	if err := ensureRepoFingerprint(townPath); err != nil {
		// Non-fatal: fingerprint is optional for functionality, just daemon optimization
		fmt.Printf("   %s Could not verify repo fingerprint: %v\n", style.Dim.Render("⚠"), err)
	}

	return nil
}

// ensureRepoFingerprint runs bd migrate --update-repo-id to ensure the database
// has a repository fingerprint. Legacy databases (pre-0.17.5) lack this, which
// prevents the daemon from starting properly.
func ensureRepoFingerprint(beadsPath string) error {
	cmd := exec.Command("bd", "migrate", "--update-repo-id")
	cmd.Dir = beadsPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bd migrate --update-repo-id: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// initTownAgentBeads creates town-level agent and role beads using hq- prefix.
// This creates:
// - hq-mayor, hq-deacon (agent beads for town-level agents)
// - hq-mayor-role, hq-deacon-role, hq-witness-role, hq-refinery-role,
//   hq-polecat-role, hq-crew-role (role definition beads)
//
// These beads are stored in town beads (~/gt/.beads/) and are shared across all rigs.
// Rig-level agent beads (witness, refinery) are created by gt rig add in rig beads.
//
// ERROR HANDLING ASYMMETRY:
// Agent beads (Mayor, Deacon) use hard fail - installation aborts if creation fails.
// Role beads use soft fail - logs warning and continues if creation fails.
//
// Rationale: Agent beads are identity beads that track agent state, hooks, and
// form the foundation of the CV/reputation ledger. Without them, agents cannot
// be properly tracked or coordinated. Role beads are documentation templates
// that define role characteristics but are not required for agent operation -
// agents can function without their role bead existing.
func initTownAgentBeads(townPath string) error {
	bd := beads.New(townPath)

	// Role beads (global templates)
	roleDefs := []struct {
		id    string
		title string
		desc  string
	}{
		{
			id:    beads.MayorRoleBeadIDTown(),
			title: "Mayor Role",
			desc:  "Role definition for Mayor agents. Global coordinator for cross-rig work.",
		},
		{
			id:    beads.DeaconRoleBeadIDTown(),
			title: "Deacon Role",
			desc:  "Role definition for Deacon agents. Daemon beacon for heartbeats and monitoring.",
		},
		{
			id:    beads.DogRoleBeadIDTown(),
			title: "Dog Role",
			desc:  "Role definition for Dog agents. Town-level workers for cross-rig tasks.",
		},
		{
			id:    beads.WitnessRoleBeadIDTown(),
			title: "Witness Role",
			desc:  "Role definition for Witness agents. Per-rig worker monitor with progressive nudging.",
		},
		{
			id:    beads.RefineryRoleBeadIDTown(),
			title: "Refinery Role",
			desc:  "Role definition for Refinery agents. Merge queue processor with verification gates.",
		},
		{
			id:    beads.PolecatRoleBeadIDTown(),
			title: "Polecat Role",
			desc:  "Role definition for Polecat agents. Ephemeral workers for batch work dispatch.",
		},
		{
			id:    beads.CrewRoleBeadIDTown(),
			title: "Crew Role",
			desc:  "Role definition for Crew agents. Persistent user-managed workspaces.",
		},
	}

	for _, role := range roleDefs {
		// Check if already exists
		if _, err := bd.Show(role.id); err == nil {
			continue // Already exists
		}

		// Create role bead using bd create --type=role
		cmd := exec.Command("bd", "create",
			"--type=role",
			"--id="+role.id,
			"--title="+role.title,
			"--description="+role.desc,
		)
		cmd.Dir = townPath
		if output, err := cmd.CombinedOutput(); err != nil {
			// Log but continue - role beads are optional
			fmt.Printf("   %s Could not create role bead %s: %s\n",
				style.Dim.Render("⚠"), role.id, strings.TrimSpace(string(output)))
			continue
		}
		fmt.Printf("   ✓ Created role bead: %s\n", role.id)
	}

	// Town-level agent beads
	agentDefs := []struct {
		id       string
		roleType string
		title    string
	}{
		{
			id:       beads.MayorBeadIDTown(),
			roleType: "mayor",
			title:    "Mayor - global coordinator, handles cross-rig communication and escalations.",
		},
		{
			id:       beads.DeaconBeadIDTown(),
			roleType: "deacon",
			title:    "Deacon (daemon beacon) - receives mechanical heartbeats, runs town plugins and monitoring.",
		},
	}

	existingAgents, err := bd.List(beads.ListOptions{
		Status:   "all",
		Type:     "agent",
		Priority: -1,
	})
	if err != nil {
		return fmt.Errorf("listing existing agent beads: %w", err)
	}
	existingAgentIDs := make(map[string]struct{}, len(existingAgents))
	for _, issue := range existingAgents {
		existingAgentIDs[issue.ID] = struct{}{}
	}

	for _, agent := range agentDefs {
		if _, ok := existingAgentIDs[agent.id]; ok {
			continue
		}

		fields := &beads.AgentFields{
			RoleType:   agent.roleType,
			Rig:        "", // Town-level agents have no rig
			AgentState: "idle",
			HookBead:   "",
			RoleBead:   beads.RoleBeadIDTown(agent.roleType),
		}

		if _, err := bd.CreateAgentBead(agent.id, agent.title, fields); err != nil {
			return fmt.Errorf("creating %s: %w", agent.id, err)
		}
		fmt.Printf("   ✓ Created agent bead: %s\n", agent.id)
	}

	return nil
}
