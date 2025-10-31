package session

// Package session provides small lifecycle helpers for desktop session components.
// This scaffold focuses on Waybar management and is intentionally conservative:
// - Minimal behavior (start/reload/status) with no environment mutations.
// - Best-effort process checks using standard tools (pgrep/pkill).
// - No logs/locks yet; those can be added incrementally without breaking callers.
//
// Callers should treat these helpers as non-fatal: they try their best and
// return an error only when something clearly failed to execute.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// WaybarStatus returns "running" or "stopped" based on whether a Waybar process is active.
// It returns an error only when process detection itself fails unexpectedly.
func WaybarStatus() (string, error) {
	if err := exec.Command("pgrep", "-x", "waybar").Run(); err == nil {
		return "running", nil
	}
	// pgrep exits non-zero when no matches; treat as stopped without error
	return "stopped", nil
}

// LaunchWaybar starts Waybar if it is not already running.
// - If multiple Waybar instances are detected, it attempts to dedupe before starting.
// - Adds a non-blocking flock to prevent multi-spawn races.
// - Waybar is started detached (new session); the caller is not blocked.
func LaunchWaybar() error {
	// Dedupe if multiple Waybar PIDs exist
	if out, err := exec.Command("sh", "-lc", "pgrep -x waybar | wc -l").Output(); err == nil {
		c := strings.TrimSpace(string(out))
		if n, _ := strconv.Atoi(c); n > 1 {
			_ = exec.Command("pkill", "-x", "waybar").Run()
		}
	}

	// Acquire non-blocking single-instance lock
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		lockDir := filepath.Join(home, ".cache", "archriot")
		_ = os.MkdirAll(lockDir, 0o755)
		lockPath := filepath.Join(lockDir, "waybar-launch.lock")
		if lf, e := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0644); e == nil {
			// try to acquire non-blocking exclusive lock
			if e := syscall.Flock(int(lf.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
				_ = lf.Close()
				// Another launcher holds the lock; exit quietly
				return nil
			}
			defer lf.Close()
		}
	}

	// Re-check after lock to avoid races
	if err := exec.Command("pgrep", "-x", "waybar").Run(); err == nil {
		return nil
	}

	// Start Waybar detached
	cmd := exec.Command("waybar")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

// ReloadWaybar reloads Waybar configuration if running; if not running, it starts Waybar.
// If selfPath is provided, it prefers to ask the ArchRiot binary to launch Waybar via
// `--waybar-launch` for consistency; otherwise it invokes `waybar` directly.
//
// Behavior:
// - Dedupe multiple running instances before action.
// - If running, send SIGUSR2 and watch briefly for a crash; auto-restart on crash.
// - If not running, start (prefer selfPath --waybar-launch).
func ReloadWaybar(selfPath string) error {
	// Dedupe if multiple Waybar PIDs exist
	if out, err := exec.Command("sh", "-lc", "pgrep -x waybar | wc -l").Output(); err == nil {
		c := strings.TrimSpace(string(out))
		if n, _ := strconv.Atoi(c); n > 1 {
			_ = exec.Command("pkill", "-x", "waybar").Run()
		}
	}

	// If Waybar isn't running, start via internal launcher when available
	if err := exec.Command("pgrep", "-x", "waybar").Run(); err != nil {
		if strings.TrimSpace(selfPath) != "" {
			_ = exec.Command(selfPath, "--waybar-launch").Start()
			return nil
		}
		cmd := exec.Command("waybar")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		_ = cmd.Start()
		return nil
	}

	// Send SIGUSR2 to reload
	_ = exec.Command("pkill", "-SIGUSR2", "waybar").Run()

	// Short window to detect crash and auto-restart
	crashed := false
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		if err := exec.Command("pgrep", "-x", "waybar").Run(); err != nil {
			crashed = true
			break
		}
	}
	if crashed {
		if strings.TrimSpace(selfPath) != "" {
			_ = exec.Command(selfPath, "--waybar-launch").Start()
			return nil
		}
		_ = LaunchWaybar()
	}
	return nil
}
