# Implementation Plan - Gastown Kiro CLI Compatibility

## Implementation Checklist

- [ ] Step 1: Create CLI abstraction interface
- [ ] Step 2: Implement Claude Code CLI adapter  
- [ ] Step 3: Implement Kiro CLI adapter
- [ ] Step 4: Create CLI factory and configuration system
- [ ] Step 5: Integrate abstraction layer into Gastown core
- [ ] Step 6: Add global CLI configuration
- [ ] Step 7: Update session management
- [ ] Step 8: Update process monitoring
- [ ] Step 9: Create parallel testing framework
- [ ] Step 10: Add configuration migration support
- [ ] Step 11: Update documentation
- [ ] Step 12: End-to-end testing and validation

## Detailed Implementation Steps

### Step 1: Create CLI Abstraction Interface

**Objective**: Define the core interface that both Claude Code and Kiro CLI implementations will satisfy.

**Implementation Guidance**:
- Create `internal/cli/interface.go` with the main `CLI` interface
- Define supporting types like `RoleType`, `CLIError`
- Keep interface minimal and focused on essential operations
- Include comprehensive documentation for each method

**Tests**: Unit tests for interface contracts and error types

**Demo**: Interface compiles and basic types are defined correctly

### Step 2: Implement Claude Code CLI Adapter

**Objective**: Wrap existing Claude Code integration in the new CLI interface.

**Implementation Guidance**:
- Create `internal/cli/claude.go` implementing the `CLI` interface
- Extract existing Claude Code logic from various files into this adapter
- Maintain exact same behavior as current implementation
- Ensure all existing functionality continues to work

**Tests**: Verify Claude Code adapter produces identical commands to current implementation

**Demo**: Claude Code adapter generates correct startup and resume commands

### Step 3: Implement Kiro CLI Adapter

**Objective**: Create Kiro CLI implementation that provides equivalent functionality.

**Implementation Guidance**:
- Create `internal/cli/kiro.go` implementing the `CLI` interface
- Map Gastown concepts to Kiro CLI equivalents
- Handle session management differences (no manual session IDs)
- Create agent configuration files instead of `.claude/settings.json`

**Tests**: Unit tests for command generation and configuration creation

**Demo**: Kiro CLI adapter generates appropriate `kiro-cli chat` commands and agent configs

### Step 4: Create CLI Factory and Configuration System

**Objective**: Implement factory pattern for CLI creation and global configuration management.

**Implementation Guidance**:
- Create `internal/cli/factory.go` with CLI factory implementation
- Add global configuration support in town settings
- Implement CLI type detection and validation
- Add configuration migration for existing setups

**Tests**: Factory creates correct CLI instances based on configuration

**Demo**: Factory correctly instantiates Claude Code or Kiro CLI based on global setting

### Step 5: Integrate Abstraction Layer into Gastown Core

**Objective**: Replace direct Claude Code calls with CLI abstraction throughout Gastown.

**Implementation Guidance**:
- Update `internal/config/loader.go` to use CLI factory
- Modify `internal/session/manager.go` to use CLI interface
- Update all command building functions to use CLI abstraction
- Ensure backwards compatibility with existing configurations

**Tests**: Integration tests verify all core functionality works with abstraction

**Demo**: Mayor, Witness, Refinery, and Polecat all work with CLI abstraction

### Step 6: Add Global CLI Configuration

**Objective**: Implement global CLI type selection in town configuration.

**Implementation Guidance**:
- Extend `TownSettings` with CLI configuration field
- Add CLI type validation and defaults
- Create configuration migration for existing towns
- Add CLI configuration commands to `gt` CLI

**Tests**: Configuration loading, validation, and migration tests

**Demo**: Users can set global CLI preference and it affects all agents

### Step 7: Update Session Management

**Objective**: Abstract session management differences between CLIs.

**Implementation Guidance**:
- Update `internal/tmux/tmux.go` to use CLI interface for process detection
- Modify session health checks to work with both CLIs
- Handle session resume differences appropriately
- Maintain existing session semantics for Gastown workflows

**Tests**: Session management tests for both CLI types

**Demo**: Sessions start, resume, and health monitoring works with both CLIs

### Step 8: Update Process Monitoring

**Objective**: Abstract process monitoring to work with different CLI process names.

**Implementation Guidance**:
- Update `internal/doctor/orphan_check.go` to use CLI interface
- Modify process detection to handle both `claude` and `kiro-cli` processes
- Update health check logic to work with different process patterns
- Ensure orphan process cleanup works for both CLIs

**Tests**: Process monitoring tests for both CLI types

**Demo**: Process health monitoring and cleanup works correctly for both CLIs

### Step 9: Create Parallel Testing Framework

**Objective**: Implement testing framework that runs same tests against both CLIs.

**Implementation Guidance**:
- Create test utilities for running parallel CLI tests
- Implement test fixtures for both CLI types
- Add integration tests that verify identical behavior
- Create performance comparison tests

**Tests**: Framework itself has unit tests, plus parallel integration tests

**Demo**: Same workflow tests pass for both Claude Code and Kiro CLI

### Step 10: Add Configuration Migration Support

**Objective**: Provide smooth migration path for existing Gastown installations.

**Implementation Guidance**:
- Detect existing Claude Code configurations
- Provide migration commands for switching CLI types
- Handle edge cases and configuration conflicts
- Add validation for mixed configurations

**Tests**: Migration tests for various existing configuration scenarios

**Demo**: Existing Gastown installation can be migrated to use Kiro CLI

### Step 11: Update Documentation

**Objective**: Document the new CLI compatibility features and configuration options.

**Implementation Guidance**:
- Update README with CLI configuration instructions
- Add troubleshooting guide for CLI compatibility issues
- Document CLI-specific limitations and differences
- Update installation instructions for both CLIs

**Tests**: Documentation accuracy verification

**Demo**: Clear documentation explains how to configure and use both CLIs

### Step 12: End-to-End Testing and Validation

**Objective**: Comprehensive testing of all Gastown features with both CLIs.

**Implementation Guidance**:
- Test all major workflows (convoys, beads, formulas, molecules)
- Verify agent coordination works with both CLIs
- Test error scenarios and recovery
- Performance testing and comparison
- User acceptance testing

**Tests**: Comprehensive end-to-end test suite

**Demo**: Complete Gastown functionality working identically with both Claude Code and Kiro CLI
