# Rough Idea: Make Gastown Compatible with Kiro CLI

## Original Request
"help me make gastown compatible with kiro-cli, atm it only works with claude code"

## Context
- Gastown is a multi-agent orchestrator for Claude Code
- It currently only works with Claude Code CLI (Anthropic's agentic coding tool)
- Need to add compatibility with Kiro CLI (AWS's AI assistant)
- Both are terminal-based AI coding assistants with similar capabilities

## Current State
- Gastown has roles: Mayor, Deacon, Witness, Refinery, Polecat
- Uses tmux sessions for agent coordination
- Tracks work with "convoys" and "beads" (git-backed issue tracker)
- Spawns Claude Code instances to execute work
- Has structured workflows via "formulas" and "molecules"

## Goal
Make Gastown work with both Claude Code and Kiro CLI, allowing users to choose their preferred AI assistant while maintaining all orchestration capabilities.
