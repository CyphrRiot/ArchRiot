package session

// SuspendGuard prevents suspend when the system appears docked or on AC power.
// If an external display is detected (via Hyprland or DRM) or power is online,
// the function exits without suspending. Otherwise, it triggers a suspend.
//
// Behavior:
// - Hyprland detection: checks hyprctl monitors output for external connectors.
// - DRM detection: reads /sys/class/drm/*/status for "connected" connectors and
//   classifies internal (eDP/LVDS/DSI) vs external (HDMI/DP/USB-C/etc).
// - Power detection: reads /sys/class/power_supply for "Mains"/USB online, or
//   any battery whose status is not "Discharging" (treated as AC).
//
// Notes:
// - This function is best-effort and does not return an error; callers can ignore.
// - When suspend is appropriate (undocked, on battery), it invokes:
//   systemctl suspend
//
// Example usage:
//   session.SuspendGuard()
//
// This implementation is extracted from the previous inline logic in main.go to
// keep the entrypoint as delegation-only and reduce complexity.
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SuspendGuard performs the undocked-on-battery suspend guard.
func SuspendGuard() {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
	trim := func(b []byte) string { return strings.TrimSpace(string(b)) }

	// Detect external displays via Hyprland (best-effort)
	isDockedHyprctl := func() bool {
		if !have("hyprctl") {
			return false
		}
		out, err := exec.Command("hyprctl", "monitors").CombinedOutput()
		if err != nil {
			return false
		}
		s := strings.ToUpper(string(out))
		// Common external connector hints in Hyprland monitor names
		return strings.Contains(s, " HDMI") ||
			strings.Contains(s, " DP-") ||
			strings.Contains(s, " DISPLAYPORT") ||
			strings.Contains(s, " DVI") ||
			strings.Contains(s, " VGA") ||
			strings.Contains(s, " USB-C")
	}

	// Detect docked state via DRM connector status
	isDockedDRM := func() bool {
		entries, err := os.ReadDir("/sys/class/drm")
		if err != nil {
			return false
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			connDir := filepath.Join("/sys/class/drm", e.Name())
			stPath := filepath.Join(connDir, "status")
			if fi, err := os.Stat(stPath); err == nil && !fi.IsDir() {
				state, _ := os.ReadFile(stPath)
				if !strings.EqualFold(trim(state), "connected") {
					continue
				}
				// classify connector, e.g., "card0-HDMI-A-1" -> "HDMI-A-1"
				name := strings.ToUpper(e.Name())
				if idx := strings.Index(name, "-"); idx > 0 {
					name = name[idx+1:]
				}
				isInternal := strings.HasPrefix(name, "EDP") ||
					strings.HasPrefix(name, "LVDS") ||
					strings.HasPrefix(name, "DSI")
				isExternal := strings.HasPrefix(name, "HDMI") ||
					strings.HasPrefix(name, "DP") ||
					strings.HasPrefix(name, "DISPLAYPORT") ||
					strings.HasPrefix(name, "DVI") ||
					strings.HasPrefix(name, "VGA") ||
					strings.HasPrefix(name, "USB") ||
					strings.HasPrefix(name, "USB-C")
				if isExternal && !isInternal {
					return true
				}
			}
		}
		return false
	}

	// Consider AC/USB-PD online as "docked" (skip suspend)
	isOnAC := func() bool {
		psEntries, err := os.ReadDir("/sys/class/power_supply")
		if err == nil {
			// Primary pass: Mains/USB online
			for _, e := range psEntries {
				base := filepath.Join("/sys/class/power_supply", e.Name())
				typB, terr := os.ReadFile(filepath.Join(base, "type"))
				onlB, oerr := os.ReadFile(filepath.Join(base, "online"))
				if terr == nil && oerr == nil {
					typ := strings.TrimSpace(string(typB))
					on := strings.TrimSpace(string(onlB))
					if (strings.EqualFold(typ, "Mains") || strings.HasPrefix(strings.ToUpper(typ), "USB")) && on == "1" {
						return true
					}
				}
			}
			// Fallback: any battery not Discharging -> treat as on AC
			for _, e := range psEntries {
				if !strings.HasPrefix(strings.ToUpper(e.Name()), "BAT") {
					continue
				}
				stB, err := os.ReadFile(filepath.Join("/sys/class/power_supply", e.Name(), "status"))
				if err == nil {
					st := strings.TrimSpace(string(stB))
					if st != "" && !strings.EqualFold(st, "Discharging") {
						return true
					}
				}
			}
		}
		return false
	}

	// If docked or on AC → do nothing
	if isDockedHyprctl() || isDockedDRM() || isOnAC() {
		return
	}

	// Undocked on battery → suspend
	_ = exec.Command("systemctl", "suspend").Run()
}
