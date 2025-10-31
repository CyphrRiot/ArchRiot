package session

// Inhibit starts a background systemd-inhibit session to keep the system awake.
// If cmdArgs is non-empty, it inhibits while that command runs; otherwise it
// runs a simple sleep loop. The process is detached into its own session so
// the caller is not blocked.
//
// Behavior:
// - If systemd-inhibit is not available, this is a no-op and returns nil.
// - Detaches stdio and session (new session via Setsid).
// - Returns any error from starting the inhibitor process.
import (
	"os/exec"
	"syscall"
)

// Inhibit launches a detached inhibitor and returns immediately.
func Inhibit(cmdArgs []string) error {
	// Prefer systemd-inhibit; if unavailable, silently no-op
	if _, err := exec.LookPath("systemd-inhibit"); err != nil {
		return nil
	}

	var cmd *exec.Cmd
	if len(cmdArgs) > 0 {
		// Inhibit while the provided command runs
		args := append([]string{"--what=sleep", "--why=ArchRiot Stay Awake"}, cmdArgs...)
		cmd = exec.Command("systemd-inhibit", args...)
	} else {
		// Background inhibitor loop
		cmd = exec.Command(
			"systemd-inhibit",
			"--what=sleep",
			"--why=ArchRiot Stay Awake",
			"bash", "-lc", "while :; do sleep 300; done",
		)
	}

	// Detach from the parent session and stdio
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Start()
}
