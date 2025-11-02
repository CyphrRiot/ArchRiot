package session

import (
	"os"
	"os/exec"
	"syscall"

	"archriot-installer/displays"
)

// StabilizeSession dedupes and relaunches session components for a healthy desktop state.
// Behavior (best-effort; no errors returned):
//   - Kills any existing Waybar instances, then launches a single managed instance.
//     Prefer launching via the current ArchRiot binary (--waybar-launch) when resolvable;
//     otherwise fallback to launching "waybar" directly.
//   - Restarts hypridle if present (kill then start detached).
//   - Optionally starts a background sleep inhibitor when withInhibit is true.
func StabilizeSession(withInhibit bool) {
	// Enforce external-preferred display policy (no compositor restart)
	_ = displays.Enforce()

	// Kill existing Waybar instances (best-effort)
	_ = exec.Command("pkill", "-x", "waybar").Run()

	// Relaunch Waybar via ArchRiot (preferred) or directly
	if self, err := os.Executable(); err == nil && self != "" {
		_ = exec.Command(self, "--waybar-launch").Start()
	} else {
		_ = exec.Command("waybar").Start()
	}

	// Restart hypridle (best-effort)
	_ = exec.Command("pkill", "hypridle").Run()
	if _, err := exec.LookPath("hypridle"); err == nil {
		cmd := exec.Command("hypridle")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		_ = cmd.Start()
	}

	// Optional sleep inhibitor
	if withInhibit {
		// Prefer native inhibitor helper in this package
		_ = Inhibit(nil)
	}
}
