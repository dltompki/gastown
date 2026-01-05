// Package clitest provides testing utilities for CLI abstraction.
package clitest

import (
	"testing"

	"github.com/steveyegge/gastown/internal/cli"
)

// CLITestSuite provides parallel testing for both CLI implementations.
type CLITestSuite struct {
	claudeCLI cli.CLI
	kiroCLI   cli.CLI
}

// NewCLITestSuite creates a new test suite with both CLI implementations.
func NewCLITestSuite() *CLITestSuite {
	return &CLITestSuite{
		claudeCLI: cli.NewClaudeCodeCLI(),
		kiroCLI:   cli.NewKiroCLI(),
	}
}

// TestWorkflow runs a test workflow against both CLI implementations.
func (s *CLITestSuite) TestWorkflow(t *testing.T, name string, workflow func(t *testing.T, cli cli.CLI)) {
	t.Run(name+"/Claude", func(t *testing.T) {
		workflow(t, s.claudeCLI)
	})
	
	t.Run(name+"/Kiro", func(t *testing.T) {
		workflow(t, s.kiroCLI)
	})
}

// TestCommandBuilding tests command building for both CLIs.
func (s *CLITestSuite) TestCommandBuilding(t *testing.T) {
	s.TestWorkflow(t, "CommandBuilding", func(t *testing.T, cli cli.CLI) {
		// Test startup command
		cmd := cli.BuildStartupCommand("crew", "test/crew/alice", "/tmp/test", "")
		if cmd == "" {
			t.Error("BuildStartupCommand returned empty string")
		}

		// Test resume command (if supported)
		if cli.SupportsSessionResume() {
			resumeCmd := cli.BuildResumeCommand("test-session")
			if resumeCmd == "" {
				t.Error("BuildResumeCommand returned empty string for CLI that supports resume")
			}
		}
	})
}

// TestConfiguration tests configuration creation for both CLIs.
func (s *CLITestSuite) TestConfiguration(t *testing.T) {
	s.TestWorkflow(t, "Configuration", func(t *testing.T, cliInstance cli.CLI) {
		// Test configuration creation
		tempDir := t.TempDir()
		err := cliInstance.CreateConfiguration(tempDir, cli.Autonomous)
		if err != nil {
			t.Errorf("CreateConfiguration failed: %v", err)
		}
	})
}

// TestFeatureSupport tests feature support flags for both CLIs.
func (s *CLITestSuite) TestFeatureSupport(t *testing.T) {
	s.TestWorkflow(t, "FeatureSupport", func(t *testing.T, cli cli.CLI) {
		// Test feature support methods
		_ = cli.SupportsHooks()
		_ = cli.SupportsForkSession()
		_ = cli.SupportsSessionResume()
		
		// Test process names
		processNames := cli.GetProcessNames()
		if len(processNames) == 0 {
			t.Error("GetProcessNames returned empty slice")
		}
		
		// Test CLI type
		cliType := cli.GetType()
		if cliType == "" {
			t.Error("GetType returned empty string")
		}
	})
}

// CompareResults compares results between CLI implementations.
func (s *CLITestSuite) CompareResults(t *testing.T, name string, testFunc func(cli cli.CLI) interface{}) {
	claudeResult := testFunc(s.claudeCLI)
	kiroResult := testFunc(s.kiroCLI)
	
	// Results don't need to be identical, but both should be valid
	if claudeResult == nil && kiroResult == nil {
		t.Errorf("%s: both CLIs returned nil", name)
	}
}
