package version

import (
	"os"
	"path/filepath"
	"strings"
)

// Version holds the current version string
var Version string

// ReadVersion reads version from VERSION file
func ReadVersion() error {
	// Get script directory
	scriptDir := filepath.Dir(os.Args[0])
	if scriptDir == "." {
		if wd, err := os.Getwd(); err == nil {
			scriptDir = wd
		}
	}

	// Look for VERSION file in parent directory (since we're in install/)
	versionFile := filepath.Join(filepath.Dir(scriptDir), "VERSION")

	if data, err := os.ReadFile(versionFile); err == nil {
		Version = strings.TrimSpace(string(data))
		return nil
	}

	// Fallback to home directory ArchRiot installation
	homeDir, err := os.UserHomeDir()
	if err != nil {
		Version = "unknown"
		return nil
	}

	versionFile = filepath.Join(homeDir, ".local", "share", "archriot", "VERSION")
	if data, err := os.ReadFile(versionFile); err == nil {
		Version = strings.TrimSpace(string(data))
		return nil
	}

	Version = "unknown"
	return nil
}

// Get returns the current version string
func Get() string {
	return Version
}
