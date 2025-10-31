package session

import (
	"os/exec"
	"strconv"
	"strings"
)

// WorkspaceClick handles Waybar workspace click routing.
// It validates that the provided name is a numeric workspace identifier and,
// if valid, dispatches a Hyprland workspace switch via `hyprctl`.
// The name should be the raw {name} from Waybar.
func WorkspaceClick(name string) {
	ws := strings.TrimSpace(name)
	if ws == "" {
		return
	}

	// Ensure hyprctl is available
	if _, err := exec.LookPath("hyprctl"); err != nil {
		return
	}

	// Validate numeric workspace name
	if _, err := strconv.Atoi(ws); err != nil {
		return
	}

	// Switch workspace
	_ = exec.Command("hyprctl", "dispatch", "workspace", ws).Run()
}
