package standalone

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Runtime represents a local standalone runtime configuration
type Runtime struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Registry    string `json:"registry" yaml:"registry"`
	Author      string `json:"author,omitempty" yaml:"author,omitempty"`
}

// RuntimeManager manages local runtime configurations
type RuntimeManager struct {
	configDir string
}

// NewRuntimeManager creates a new runtime manager
func NewRuntimeManager() *RuntimeManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".synadia", "connect", "standalone")
	return &RuntimeManager{
		configDir: configDir,
	}
}

// GetConfigDir returns the configuration directory
func (rm *RuntimeManager) GetConfigDir() string {
	return rm.configDir
}

// ensureConfigDir creates the config directory if it doesn't exist
func (rm *RuntimeManager) ensureConfigDir() error {
	return os.MkdirAll(rm.configDir, 0755)
}

// runtimesFile returns the path to the runtimes configuration file
func (rm *RuntimeManager) runtimesFile() string {
	return filepath.Join(rm.configDir, "runtimes.json")
}

// LoadRuntimes loads all configured runtimes
func (rm *RuntimeManager) LoadRuntimes() ([]Runtime, error) {
	if err := rm.ensureConfigDir(); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// If file doesn't exist, create it with default runtimes
	if _, err := os.Stat(rm.runtimesFile()); os.IsNotExist(err) {
		if err := rm.initializeDefaultRuntimes(); err != nil {
			return nil, fmt.Errorf("failed to initialize default runtimes: %w", err)
		}
	}

	data, err := os.ReadFile(rm.runtimesFile())
	if err != nil {
		return nil, fmt.Errorf("failed to read runtimes file: %w", err)
	}

	var runtimes []Runtime
	if err := json.Unmarshal(data, &runtimes); err != nil {
		return nil, fmt.Errorf("failed to parse runtimes file: %w", err)
	}

	return runtimes, nil
}

// SaveRuntimes saves the runtime configurations
func (rm *RuntimeManager) SaveRuntimes(runtimes []Runtime) error {
	if err := rm.ensureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(runtimes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal runtimes: %w", err)
	}

	if err := os.WriteFile(rm.runtimesFile(), data, 0644); err != nil {
		return fmt.Errorf("failed to write runtimes file: %w", err)
	}

	return nil
}

// GetRuntime retrieves a specific runtime by ID
func (rm *RuntimeManager) GetRuntime(id string) (*Runtime, error) {
	runtimes, err := rm.LoadRuntimes()
	if err != nil {
		return nil, err
	}

	for _, runtime := range runtimes {
		if runtime.ID == id {
			return &runtime, nil
		}
	}

	return nil, fmt.Errorf("runtime '%s' not found", id)
}

// AddRuntime adds a new runtime configuration
func (rm *RuntimeManager) AddRuntime(runtime Runtime) error {
	// Validate required fields
	if runtime.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if runtime.Registry == "" {
		return fmt.Errorf("Registry is required")
	}

	runtimes, err := rm.LoadRuntimes()
	if err != nil {
		return err
	}

	// Check if runtime already exists
	for _, existing := range runtimes {
		if existing.ID == runtime.ID {
			return fmt.Errorf("runtime '%s' already exists", runtime.ID)
		}
	}

	runtimes = append(runtimes, runtime)
	return rm.SaveRuntimes(runtimes)
}

// RemoveRuntime removes a runtime configuration
func (rm *RuntimeManager) RemoveRuntime(id string) error {
	// Protect default wombat runtime
	if id == "wombat" {
		return fmt.Errorf("cannot remove default runtime 'wombat'")
	}

	runtimes, err := rm.LoadRuntimes()
	if err != nil {
		return err
	}

	found := false
	newRuntimes := make([]Runtime, 0, len(runtimes))
	for _, runtime := range runtimes {
		if runtime.ID != id {
			newRuntimes = append(newRuntimes, runtime)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("runtime '%s' not found", id)
	}

	return rm.SaveRuntimes(newRuntimes)
}

// UpdateRuntime updates an existing runtime configuration
func (rm *RuntimeManager) UpdateRuntime(runtime Runtime) error {
	runtimes, err := rm.LoadRuntimes()
	if err != nil {
		return err
	}

	found := false
	for i, existing := range runtimes {
		if existing.ID == runtime.ID {
			runtimes[i] = runtime
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("runtime '%s' not found", runtime.ID)
	}

	return rm.SaveRuntimes(runtimes)
}

// initializeDefaultRuntimes creates the initial runtime configurations
func (rm *RuntimeManager) initializeDefaultRuntimes() error {
	defaultRuntimes := []Runtime{
		{
			ID:          "wombat",
			Name:        "Wombat Runtime",
			Description: "Default Wombat-based connector runtime for Synadia Connect",
			Registry:    "registry.synadia.io/connect-runtime-wombat",
			Author:      "Synadia",
		},
	}

	return rm.SaveRuntimes(defaultRuntimes)
}

// ResolveRuntimeImage resolves a runtime reference to a full Docker image name
// Examples:
//
//	"wombat" -> "registry.synadia.io/connect-runtime-wombat:latest"
//	"wombat:v1.0.3" -> "registry.synadia.io/connect-runtime-wombat:v1.0.3"
func (rm *RuntimeManager) ResolveRuntimeImage(runtimeRef string) (string, error) {
	// Parse runtime reference (id:version)
	parts := strings.SplitN(runtimeRef, ":", 2)
	runtimeID := parts[0]
	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	// Get runtime configuration
	runtime, err := rm.GetRuntime(runtimeID)
	if err != nil {
		return "", fmt.Errorf("runtime '%s' not found: %w", runtimeID, err)
	}

	// Build full image name
	return fmt.Sprintf("%s:%s", runtime.Registry, version), nil
}
