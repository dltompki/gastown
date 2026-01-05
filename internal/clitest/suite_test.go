package clitest

import (
	"testing"

	"github.com/steveyegge/gastown/internal/cli"
)

func TestCLIAbstraction(t *testing.T) {
	suite := NewCLITestSuite()
	
	// Test command building
	suite.TestCommandBuilding(t)
	
	// Test configuration
	suite.TestConfiguration(t)
	
	// Test feature support
	suite.TestFeatureSupport(t)
}

func TestCLIFactory(t *testing.T) {
	factory := cli.GetGlobalFactory()
	
	// Test Claude CLI creation
	claudeCLI, err := factory.CreateCLI("claude")
	if err != nil {
		t.Fatalf("Failed to create Claude CLI: %v", err)
	}
	if claudeCLI.GetType() != "claude" {
		t.Errorf("Expected claude, got %s", claudeCLI.GetType())
	}
	
	// Test Kiro CLI creation
	kiroCLI, err := factory.CreateCLI("kiro")
	if err != nil {
		t.Fatalf("Failed to create Kiro CLI: %v", err)
	}
	if kiroCLI.GetType() != "kiro" {
		t.Errorf("Expected kiro, got %s", kiroCLI.GetType())
	}
	
	// Test unsupported CLI type
	_, err = factory.CreateCLI("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported CLI type")
	}
}

func TestCommandEquivalence(t *testing.T) {
	suite := NewCLITestSuite()
	
	// Test that both CLIs can build valid commands
	suite.CompareResults(t, "StartupCommand", func(cli cli.CLI) interface{} {
		return cli.BuildStartupCommand("crew", "test/crew/alice", "/tmp/test", "hello")
	})
	
	// Test process names
	suite.CompareResults(t, "ProcessNames", func(cli cli.CLI) interface{} {
		return cli.GetProcessNames()
	})
}
