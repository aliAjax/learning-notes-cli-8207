package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config contains runtime configuration shared by all commands.
type Config struct {
	DataDir string
}

// New returns a validated config with an absolute data directory.
func New(dataDir string) (Config, error) {
	dataDir = strings.TrimSpace(dataDir)
	if dataDir == "" {
		return Config{}, fmt.Errorf("data directory cannot be empty")
	}

	abs, err := filepath.Abs(dataDir)
	if err != nil {
		return Config{}, fmt.Errorf("resolve data directory %q: %v", dataDir, err)
	}

	return Config{DataDir: abs}, nil
}

// EnsureDataDir creates the configured data directory if it does not exist.
func (c Config) EnsureDataDir() error {
	if err := os.MkdirAll(c.DataDir, 0o755); err != nil {
		return fmt.Errorf("create data directory %q: %w", c.DataDir, err)
	}
	return nil
}

// DefaultDataDir returns the default notes directory under the user's home.
func DefaultDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".learning-notes"), nil
}
