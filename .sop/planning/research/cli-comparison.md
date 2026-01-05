# CLI Comparison Analysis

This document compares Claude Code and Kiro CLI interfaces to identify compatibility requirements.

## Command Interface Comparison

| Feature | Claude Code | Kiro CLI | Compatibility Notes |
|---------|-------------|----------|-------------------|
| **Basic Invocation** | `claude` | `kiro-cli chat` | Different command structure |
| **Session Resume** | `claude --resume <id>` | Automatic session persistence | Different resume mechanisms |
| **Permissions** | `--dangerously-skip-permissions` | Tool-based permission system | Completely different approaches |
| **Configuration** | `.claude/settings.json` | `kiro-cli settings` + agent configs | Different config systems |
| **Hooks** | JSON-based hooks in settings | Agent-based hooks | Similar concept, different implementation |
| **Session Management** | Manual session IDs | Automatic conversation management | Different session models |

## Key Differences

### 1. Command Structure
- **Claude Code**: Single binary with flags (`claude --resume --dangerously-skip-permissions`)
- **Kiro CLI**: Subcommand structure (`kiro-cli chat --agent myagent`)

### 2. Session Management
- **Claude Code**: Manual session IDs via `CLAUDE_SESSION_ID` environment variable
- **Kiro CLI**: Automatic conversation persistence with `/save` and `/load`

### 3. Configuration System
- **Claude Code**: JSON configuration in `.claude/settings.json`
- **Kiro CLI**: Settings via `kiro-cli settings` command + agent JSON files

### 4. Permission Model
- **Claude Code**: Global `--dangerously-skip-permissions` flag
- **Kiro CLI**: Granular tool permissions in agent configuration

### 5. Hook System
- **Claude Code**: Hooks defined in settings.json, executed at specific lifecycle events
- **Kiro CLI**: Agent-based hooks with similar lifecycle events but different configuration

## Abstraction Requirements

To make Gastown work with both CLIs, we need to abstract:

### 1. Command Building
- **Current**: Hardcoded Claude Code command structure
- **Needed**: Pluggable command builders for different CLI types

### 2. Session Management
- **Current**: `CLAUDE_SESSION_ID` environment variable
- **Needed**: CLI-specific session handling (env vars vs conversation management)

### 3. Configuration Management
- **Current**: Creates `.claude/settings.json` files
- **Needed**: CLI-specific configuration file generation

### 4. Process Detection
- **Current**: Searches for `claude` and `claude-code` processes
- **Needed**: CLI-specific process name detection

### 5. Hook Integration
- **Current**: Claude Code hook format and lifecycle
- **Needed**: Translate Gastown hooks to CLI-specific formats

## Compatibility Challenges

### Major Challenges
1. **Different session models**: Claude Code uses manual session IDs, Kiro CLI uses automatic persistence
2. **Permission systems**: Completely different approaches to tool permissions
3. **Configuration formats**: JSON vs command-line settings
4. **Command structure**: Single binary vs subcommands

### Minor Challenges
1. **Process names**: Different binary names for health monitoring
2. **Environment variables**: Different session ID variable names
3. **Hook timing**: May need to map hook events between systems

## Implementation Strategy

### 1. CLI Abstraction Layer
Create interface that abstracts:
- Command building
- Session management
- Configuration generation
- Process monitoring
- Hook integration

### 2. CLI-Specific Implementations
- `ClaudeCodeCLI` - Current implementation
- `KiroCLI` - New implementation for Kiro CLI

### 3. Configuration Mapping
- Map Gastown concepts to CLI-specific configurations
- Handle permission model differences
- Translate hook events and commands

### 4. Session Bridging
- Abstract session management differences
- Handle resume/persistence model differences
- Maintain Gastown's session semantics regardless of underlying CLI
