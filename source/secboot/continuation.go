package secboot

import (
	"fmt"
	"os"
	"path/filepath"
)

// Logger is a minimal logger function signature for dependency injection.
// Pass a function that records log messages with a level and message.
// Example injection:
//
//	secboot.SetLogger(func(level, msg string) { logger.LogMessage(level, msg) })
type LoggerFunc func(level, msg string)

var logFn LoggerFunc

// SetLogger configures an optional logger used by this package.
// When unset, logging calls are ignored (no-op).
func SetLogger(fn LoggerFunc) {
	logFn = fn
}

func log(level, msg string) {
	if logFn != nil {
		logFn(level, msg)
	}
}

// RestoreHyprlandConfig restores the user's Hyprland configuration from a backup
// created by the Secure Boot continuation flow.
//
// Behavior:
//   - Looks for ~/.config/hypr/hyprland.conf.archriot-backup
//   - Restores it to ~/.config/hypr/hyprland.conf
//   - Removes the backup file upon success
//   - Uses injected logger when provided (via SetLogger)
//
// Returns:
//   - error when the backup file is missing or I/O operations fail
func RestoreHyprlandConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprlandConfigPath := filepath.Join(homeDir, ".config", "hypr", "hyprland.conf")
	backupPath := hyprlandConfigPath + ".archriot-backup"

	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup hyprland.conf not found at %s", backupPath)
	}

	// Read backup
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("reading backup: %w", err)
	}

	// Ensure target directory exists (defensive)
	if err := os.MkdirAll(filepath.Dir(hyprlandConfigPath), 0o755); err != nil {
		return fmt.Errorf("ensuring target dir: %w", err)
	}

	// Write restored config
	if err := os.WriteFile(hyprlandConfigPath, backupData, 0o644); err != nil {
		return fmt.Errorf("writing restored config: %w", err)
	}

	// Remove backup file (best-effort; return error if removal fails to avoid confusion)
	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("removing backup: %w", err)
	}

	log("SUCCESS", "Restored original hyprland.conf")
	return nil
}
