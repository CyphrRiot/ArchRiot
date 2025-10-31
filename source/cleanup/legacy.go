package cleanup

// Package cleanup provides best-effort routines to remove legacy artifacts
// that were replaced by native CLI features. These helpers are safe to call
// during installation; they will never cause the installer to fail.
//
// Policy:
// - Only remove files known to be legacy helpers.
// - Operate best-effort: ignore errors, do not escalate permissions.
// - Never remove current, supported assets.

import (
	"os"
	"path/filepath"
)

// RemoveLegacyFiles deletes known legacy helper scripts (if present) from the
// user's home directory. All operations are best-effort and errors are ignored.
// This aligns with plan.md: "Script consolidation into first-class CLI".
func RemoveLegacyFiles() {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}

	// Known legacy helpers that must not exist anymore
	candidates := []string{
		// Legacy suspend guard script (replaced by --suspend-if-undocked)
		filepath.Join(home, ".local", "bin", "suspend-if-undocked.sh"),

		// Legacy Waybar helpers (replaced by native CLI emitters)
		filepath.Join(home, ".local", "bin", "waybar-tomato-timer.py"),
		filepath.Join(home, ".local", "bin", "waybar-memory-accurate.py"),
		filepath.Join(home, ".local", "share", "archriot", "install", "waybar-tomato-timer.py"),
		filepath.Join(home, ".local", "share", "archriot", "install", "waybar-memory-accurate.py"),

		// Keybindings helpers (remove all variants; rely on CLI --help-binds-web)
		// Hyphen variants
		filepath.Join(home, ".local", "share", "archriot", "config", "bin", "generate-keybindings-help.sh"),
		filepath.Join(home, ".config", "bin", "scripts", "generate-keybindings-help.sh"),
		filepath.Join(home, ".local", "share", "archriot", "config", "bin", "scripts", "generate-keybindings-help.sh"),
		// Underscore variants
		filepath.Join(home, ".local", "share", "archriot", "config", "bin", "generate_keybindings_help.sh"),
		filepath.Join(home, ".config", "bin", "scripts", "generate_keybindings_help.sh"),
		filepath.Join(home, ".local", "share", "archriot", "config", "bin", "scripts", "generate_keybindings_help.sh"),
	}

	for _, p := range candidates {
		// Ignore errors (file may not exist, or permissions may restrict removal)
		_ = os.Remove(p)
	}
}
