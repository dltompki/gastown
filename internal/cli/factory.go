package cli

import (
	"fmt"
	"sync"
)

// DefaultCLIFactory implements the CLIFactory interface.
type DefaultCLIFactory struct {
	mu       sync.RWMutex
	registry map[string]func() CLI
}

// NewCLIFactory creates a new CLI factory with built-in CLI types.
func NewCLIFactory() CLIFactory {
	factory := &DefaultCLIFactory{
		registry: make(map[string]func() CLI),
	}

	// Register built-in CLI types
	factory.RegisterCLI("claude", NewClaudeCodeCLI)
	factory.RegisterCLI("kiro", NewKiroCLI)

	return factory
}

// CreateCLI creates a CLI instance of the specified type.
func (f *DefaultCLIFactory) CreateCLI(cliType string) (CLI, error) {
	f.mu.RLock()
	constructor, exists := f.registry[cliType]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported CLI type: %s (supported: %v)", cliType, f.GetSupportedTypes())
	}

	return constructor(), nil
}

// RegisterCLI registers a new CLI type with its constructor.
func (f *DefaultCLIFactory) RegisterCLI(cliType string, constructor func() CLI) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.registry[cliType] = constructor
}

// GetSupportedTypes returns a list of supported CLI types.
func (f *DefaultCLIFactory) GetSupportedTypes() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	types := make([]string, 0, len(f.registry))
	for cliType := range f.registry {
		types = append(types, cliType)
	}
	return types
}

// Global factory instance
var (
	globalFactory CLIFactory
	factoryOnce   sync.Once
)

// GetGlobalFactory returns the global CLI factory instance.
func GetGlobalFactory() CLIFactory {
	factoryOnce.Do(func() {
		globalFactory = NewCLIFactory()
	})
	return globalFactory
}

// CreateCLIFromConfig creates a CLI instance based on configuration.
func CreateCLIFromConfig(cliType string) (CLI, error) {
	if cliType == "" {
		cliType = "claude" // Default to Claude for backwards compatibility
	}

	factory := GetGlobalFactory()
	return factory.CreateCLI(cliType)
}
