# Claude Code Integration Research

This document analyzes how Gastown currently integrates with Claude Code CLI.

## Current Integration Points

### 1. Agent Configuration System

**Agent Presets (internal/config/agents.go):**
- Built-in presets: `claude`, `gemini`, `codex`
- Default preset: `AgentClaude`
- Command: `"claude"`
- Default args: `["--dangerously-skip-permissions"]`
- Session ID env: `"CLAUDE_SESSION_ID"`
- Resume flag: `"--resume"`
- Supports hooks: `true`
- Supports fork session: `true`

### 2. Command Building

**Startup Commands:**
- `BuildCrewStartupCommand(rigName, crewName, rigPath, prompt)` - Builds full startup command with environment variables
- `BuildResumeCommand(agentName, sessionID)` - Builds resume command for existing sessions
- Environment variables: `GT_ROLE`, `GT_RIG`, `GT_CREW`, `BD_ACTOR`, `GIT_AUTHOR_NAME`

**Command Structure:**
```bash
export GT_ROLE=crew GT_RIG=myproject GT_CREW=alice BD_ACTOR=myproject/crew/alice GIT_AUTHOR_NAME=alice && claude --dangerously-skip-permissions
```

### 3. Claude Code Settings Integration

**Settings Files (.claude/settings.json):**
- **Autonomous roles** (polecat, witness, refinery): Include mail injection in SessionStart
- **Interactive roles** (mayor, crew): Mail injection in UserPromptSubmit only

**Hook Configuration:**
- `SessionStart`: `gt prime && gt mail check --inject && gt nudge deacon session-started`
- `PreCompact`: `gt prime`
- `UserPromptSubmit`: `gt mail check --inject`
- `Stop`: `gt costs record`

### 4. Session Management

**Session Detection:**
- `IsClaudeRunning(sessionName)` - Checks if Claude process is active in tmux session
- Zombie detection: tmux session exists but Claude process is dead
- Session restart logic for crashed agents

**Session Lifecycle:**
- Creates `.claude/settings.json` with appropriate hooks
- Launches Claude in tmux sessions
- Monitors Claude process health
- Handles session recovery and restart

### 5. Process Management

**Claude Process Detection:**
- Searches for `claude` and `claude-code` processes
- Orphan process cleanup
- Process health monitoring
- Timeout handling for Claude startup

### 6. Integration Dependencies

**Hard Dependencies on Claude Code:**
- Hook system relies on Claude Code's `.claude/settings.json` format
- Session ID environment variable `CLAUDE_SESSION_ID`
- Command-line flags like `--dangerously-skip-permissions`, `--resume`
- Process name detection for health monitoring
- Fork session support for seance command

**Configuration Files:**
- `.claude/settings.json` - Claude Code configuration
- Hook commands executed by Claude Code
- Plugin configuration (beads marketplace)

## Key Integration Patterns

1. **Agent Abstraction**: Uses `AgentPresetInfo` to define CLI characteristics
2. **Command Building**: Constructs full command strings with environment setup
3. **Hook Integration**: Leverages Claude Code's hook system for Gastown coordination
4. **Session Management**: Manages Claude Code instances in tmux sessions
5. **Process Monitoring**: Tracks Claude Code process health and lifecycle
