package displays

// Enforce ensures the correct monitor policy without restarting Hyprland:
// - When an external output (HDMI/DP/USB-C) is present: all internal panels (eDP/LVDS/DSI) are disabled.
// - When no external output is present: at least one internal panel is enabled using preferred mode.
// This uses `hyprctl keyword monitor` so changes apply live.
//
// Exit codes:
//   0 -> OK / best-effort completed
//   1 -> hyprctl unavailable or monitor enumeration failed (printed to stderr)
//
// Notes:
// - We intentionally keep the configuration simple: "preferred,auto,1" for the enabled output.
// - We do not restart the compositor; this is live-only.
// - This function is safe to call multiple times (idempotent in practice).
// - If kanshi is also managing outputs, ensure its configuration expresses the same policy to avoid "fighting".
//
// Example effects:
//   - External HDMI-A-1 present -> eDP-* disabled, HDMI-A-1 enabled.
//   - No external present -> internal eDP-* enabled (first found), externals disabled implicitly by absence.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// monitor describes the minimal fields we care about from `hyprctl -j monitors`
type monitor struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
}

// isInternal reports whether a connector name is a laptop/internal panel.
func isInternal(name string) bool {
	up := strings.ToUpper(name)
	return strings.HasPrefix(up, "EDP-") || strings.Contains(up, "EDP") ||
		strings.HasPrefix(up, "LVDS") || strings.HasPrefix(up, "DSI")
}

// hyprctlJSON runs `hyprctl -j <args...>` and returns stdout bytes.
func hyprctlJSON(args ...string) ([]byte, error) {
	argv := append([]string{"-j"}, args...)
	cmd := exec.Command("hyprctl", argv...)
	return cmd.CombinedOutput()
}

// hyprKeyword sets a Hyprland keyword live.
func hyprKeyword(key, value string) error {
	return exec.Command("hyprctl", "keyword", key, value).Run()
}

// pickPreferredExternal returns the focused external if any, else the first.
func pickPreferredExternal(externals []monitor) (monitor, bool) {
	for _, m := range externals {
		if m.Focused {
			return m, true
		}
	}
	if len(externals) > 0 {
		return externals[0], true
	}
	return monitor{}, false
}

// enablePreferred enables the monitor with a default layout.
func enablePreferred(name string) {
	// preferred mode, auto position, scale 1.0
	// Idempotence: if Hyprland already lists this monitor by name, assume enabled and skip.
	if out, err := hyprctlJSON("monitors"); err == nil {
		var mons []monitor
		if json.Unmarshal(out, &mons) == nil {
			for _, m := range mons {
				if m.Name == name {
					// Already present (enabled); no-op
					return
				}
			}
		}
	}
	_ = hyprKeyword("monitor", fmt.Sprintf("%s,preferred,auto,1", name))
}

// disableOutput disables a monitor by name.
func disableOutput(name string) {
	// Idempotence: if Hyprland does NOT list this monitor by name, assume it's already disabled.
	if out, err := hyprctlJSON("monitors"); err == nil {
		var mons []monitor
		if json.Unmarshal(out, &mons) == nil {
			found := false
			for _, m := range mons {
				if m.Name == name {
					found = true
					break
				}
			}
			if !found {
				// Not present (already disabled); no-op
				return
			}
		}
	}
	_ = hyprKeyword("monitor", fmt.Sprintf("%s,disable", name))
}

// detectConnectedConnectors reads sysfs DRM connector status to determine physical connection state.
// Returns connector names without the "cardX-" prefix (e.g., "HDMI-A-1", "eDP-1").
func detectConnectedConnectors() (externals, internals []string) {
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		// Expect entries like: card0-HDMI-A-1, card0-eDP-1, etc.
		if !strings.Contains(name, "-") {
			continue
		}
		// Read connection status
		stB, err := os.ReadFile(filepath.Join("/sys/class/drm", name, "status"))
		if err != nil {
			continue
		}
		status := strings.TrimSpace(strings.ToLower(string(stB)))
		if status != "connected" {
			continue
		}
		// Strip "cardN-" prefix to get connector short name
		short := name
		if idx := strings.Index(short, "-"); idx >= 0 {
			short = short[idx+1:]
		}
		if isInternal(short) {
			internals = append(internals, short)
		} else {
			externals = append(externals, short)
		}
	}
	return
}

// Enforce applies the external-preferred display policy with verification/retry to handle hotplug races.
func Enforce() int {
	// User opt-out: skip enforcement when marker file is present
	home := os.Getenv("HOME")
	if home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			home = h
		}
	}
	if home != "" {
		if st, err := os.Stat(filepath.Join(home, ".config", "archriot", "disable-display-enforcement")); err == nil && !st.IsDir() {
			fmt.Println("Displays enforcement: disabled by ~/.config/archriot/disable-display-enforcement")
			return 0
		}
	}

	if _, err := exec.LookPath("hyprctl"); err != nil {
		fmt.Fprintln(os.Stderr, "hyprctl not found in PATH")
		return 1
	}

	apply := func() (externals, internals []monitor, ok bool) {
		out, err := hyprctlJSON("monitors")
		if err != nil || len(strings.TrimSpace(string(out))) == 0 {
			// Fall back to sysfs-only decision if hyprctl JSON not ready
			sysExt, sysInt := detectConnectedConnectors()
			// If sysfs indicates an external, force-apply it
			if len(sysExt) > 0 {
				for _, in := range sysInt {
					disableOutput(in)
				}
				enablePreferred(sysExt[0])
				// Synthesize monitor slices for verification
				for _, n := range sysExt {
					externals = append(externals, monitor{Name: n})
				}
				for _, n := range sysInt {
					internals = append(internals, monitor{Name: n})
				}
				return externals, internals, true
			}
			// If only internals are connected, enable first internal
			if len(sysInt) > 0 {
				enablePreferred(sysInt[0])
				for _, n := range sysInt {
					internals = append(internals, monitor{Name: n})
				}
				return externals, internals, true
			}
			return nil, nil, false
		}
		var mons []monitor
		if err := json.Unmarshal(out, &mons); err != nil {
			return nil, nil, false
		}

		// Classify by Hyprland view
		for _, m := range mons {
			if strings.TrimSpace(m.Name) == "" {
				continue
			}
			if isInternal(m.Name) {
				internals = append(internals, m)
			} else {
				externals = append(externals, m)
			}
		}

		// Cross-check with sysfs; if Hyprland reports no externals but sysfs has one, force-enable it.
		sysExt, sysInt := detectConnectedConnectors()
		if len(externals) == 0 && len(sysExt) > 0 {
			for _, in := range sysInt {
				disableOutput(in)
			}
			enablePreferred(sysExt[0])
			// Update slices to reflect desired state
			externals = []monitor{{Name: sysExt[0]}}
			internals = nil
			for _, n := range sysInt {
				internals = append(internals, monitor{Name: n})
			}
			return externals, internals, true
		}

		// Apply desired state (Hyprland view)
		if len(externals) > 0 {
			for _, in := range internals {
				disableOutput(in.Name)
			}
			if sel, ok := pickPreferredExternal(externals); ok {
				enablePreferred(sel.Name)
			}
		} else if len(internals) > 0 {
			// If no externals according to Hyprland and sysfs, keep internal on
			enablePreferred(internals[0].Name)
		}
		return externals, internals, true
	}

	// Attempt, verify, and retry up to 5 times with small delays to handle races.
	for i := 0; i < 5; i++ {
		_, _, ok := apply()
		if !ok {
			time.Sleep(250 * time.Millisecond)
			continue
		}
		// Allow Hyprland/kanshi to settle
		time.Sleep(250 * time.Millisecond)

		// Verify state
		out, err := hyprctlJSON("monitors")
		if err == nil {
			var mons []monitor
			if json.Unmarshal(out, &mons) == nil {
				haveInternal := false
				haveExternal := false
				for _, m := range mons {
					if strings.TrimSpace(m.Name) == "" {
						continue
					}
					if isInternal(m.Name) {
						haveInternal = true
					} else {
						haveExternal = true
					}
				}
				// Success conditions
				if haveExternal && !haveInternal {
					{
						ext, in := detectConnectedConnectors()
						fmt.Printf("Displays enforced: external present -> internal panels disabled [%s]; externals [%s]\n", strings.Join(in, ","), strings.Join(ext, ","))
					}
					return 0
				}
				if !haveExternal && haveInternal {
					{
						_, in := detectConnectedConnectors()
						fmt.Printf("Displays enforced: no external -> internal panel enabled [%s]\n", strings.Join(in, ","))
					}
					return 0
				}
			}
		}
	}

	// If we reach here, we applied with retries; final state may depend on hotplug timing.
	fmt.Println("Displays enforcement: applied with retries; final state may depend on hotplug timing")
	return 0
}
