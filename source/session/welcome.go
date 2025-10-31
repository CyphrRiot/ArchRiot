package session

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// WelcomeLaunch starts the ArchRiot Welcome app non-blocking (if present).
// It looks for the executable at:
//
//	$HOME/.local/share/archriot/config/bin/welcome
//
// Behavior:
//   - If the welcome binary exists, it is launched detached (new session) so it
//     is not tied to the parent process.
//   - If it does not exist, this function returns silently (no error).
func WelcomeLaunch() {
	home := os.Getenv("HOME")
	if home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			home = h
		}
	}
	if home == "" {
		// Could not resolve HOME; nothing to do
		return
	}

	welcome := filepath.Join(home, ".local", "share", "archriot", "config", "bin", "welcome")
	if _, err := os.Stat(welcome); err != nil {
		// Not present; exit quietly
		return
	}

	cmd := exec.Command(welcome)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Start()
}
