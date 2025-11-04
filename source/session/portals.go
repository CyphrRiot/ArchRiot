package session

import (
	"fmt"
	"os/exec"
	"time"
)

// RestartPortals restarts xdg-desktop-portal providers safely in the user session.
//
// Behavior (best-effort; non-fatal):
// - Prefers systemd --user restarts for hyprland and desktop providers.
// - Falls back to killing stale processes and (re)starting units.
// - Waits briefly for providers to settle before returning.
// - Always returns 0 to avoid breaking callers.
//
// Intended use:
//
//	archriot --portals-restart
func RestartPortals() int {
	have := func(name string) bool {
		_, err := exec.LookPath(name)
		return err == nil
	}

	run := func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		return cmd.Run()
	}

	tryRestartUnit := func(unit string) bool {
		// Attempt a systemd --user restart; ignore errors.
		if !have("systemctl") {
			return false
		}
		if err := run("systemctl", "--user", "restart", unit); err != nil {
			return false
		}
		return true
	}

	tryStartUnit := func(unit string) {
		if have("systemctl") {
			_ = run("systemctl", "--user", "start", unit)
		}
	}

	// 1) Prefer systemd-controlled restarts (hyprland portal, then desktop portal)
	okHypr := tryRestartUnit("xdg-desktop-portal-hyprland.service")
	okDesk := tryRestartUnit("xdg-desktop-portal.service")

	// 2) If either restart failed (units missing or not managed), clean up stale processes
	if !okHypr || !okDesk {
		_ = run("pkill", "-x", "xdg-desktop-portal")
		_ = run("pkill", "-x", "xdg-desktop-portal-hyprland")
		// Try to (re)start the usual providers via systemd
		tryStartUnit("xdg-desktop-portal-hyprland.service")
		tryStartUnit("xdg-desktop-portal.service")
	}

	// 3) Let providers settle briefly to avoid immediate clipboard/screenshot stalls
	for i := 0; i < 8; i++ {
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Println("✓ Portals restarted (hyprland, desktop)")
	return 0
}
