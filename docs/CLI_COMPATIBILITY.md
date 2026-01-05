# CLI Compatibility Guide

This guide explains how to use Gastown with different AI CLI tools.

## Supported CLIs

Gastown supports two AI CLI tools:

- **Claude Code** - Anthropic's agentic coding tool (default)
- **Kiro CLI** - AWS's AI assistant

## Configuration

### Setting CLI Type

Configure which AI CLI to use globally:

```bash
# Use Claude Code (default)
gt cli set claude

# Use Kiro CLI  
gt cli set kiro

# Check current setting
gt cli get

# List available options
gt cli list
```

### Migration

Switch between CLIs for existing installations:

```bash
# Migrate to Kiro CLI
gt migrate cli kiro

# Migrate back to Claude Code
gt migrate cli claude
```

## Differences

While all Gastown features work identically, the underlying CLIs have different characteristics:

### Claude Code
- **Command**: `claude --dangerously-skip-permissions`
- **Configuration**: `.claude/settings.json` files
- **Session Management**: Manual session IDs via `CLAUDE_SESSION_ID`
- **Permissions**: Global bypass flag
- **Resume**: `claude --resume <session-id>`

### Kiro CLI
- **Command**: `kiro-cli chat`
- **Configuration**: `.kiro/agents/` JSON files
- **Session Management**: Automatic persistence
- **Permissions**: Granular tool permissions
- **Resume**: Automatic (no manual session IDs)

## Troubleshooting

### CLI Not Found
```bash
# Check if CLI is installed
which claude     # or which kiro-cli
which kiro-cli

# Install missing CLI
curl -fsSL https://claude.ai/install.sh | bash        # Claude Code
curl -fsSL https://cli.kiro.dev/install | bash        # Kiro CLI
```

### Configuration Issues
```bash
# Check current configuration
gt cli get

# Reset to default
gt cli set claude

# Check for conflicting configurations
ls -la .claude/    # Claude Code configs
ls -la .kiro/      # Kiro CLI configs
```

### Session Problems
```bash
# Check active sessions
gt status

# Restart sessions after CLI change
gt shutdown
gt start
```

### Process Monitoring
```bash
# Check for orphaned processes
gt doctor

# Fix orphaned processes
gt doctor --fix
```

## Best Practices

1. **Choose One CLI**: Don't mix CLIs in the same workspace
2. **Restart After Migration**: Restart active sessions when switching CLIs
3. **Check Compatibility**: Ensure your chosen CLI is installed and working
4. **Use Migration Commands**: Use `gt migrate cli` instead of manual configuration changes
5. **Test After Changes**: Run `gt doctor` to verify configuration

## Advanced Configuration

### Custom CLI Settings

For advanced users, CLI-specific settings can be customized in town configuration:

```json
{
  "type": "town-settings",
  "version": 1,
  "default_cli": "kiro",
  "agents": {
    "custom-agent": {
      "command": "custom-cli",
      "args": ["--special-flag"]
    }
  }
}
```

### Environment Variables

Some CLIs use environment variables:

- **Claude Code**: `CLAUDE_SESSION_ID`, `CLAUDE_CONFIG_DIR`
- **Kiro CLI**: Standard environment variables

Gastown handles these automatically based on your CLI configuration.
