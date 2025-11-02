package session

// Waybar sweep and restart helpers.
//
// These functions handle cases where a single Waybar process leaves behind multiple
// bar layer surfaces after display hotplug (opening laptop lid or adding/removing
// external monitors). They provide quick actions that can be wired to kanshi,
// hypridle, or other hooks to keep exactly one bar per monitor.
//
// - SweepWaybar: Compares number of Waybar windows to number of monitors and
//   restarts Waybar if the window count is clearly excessive.
// - RestartWaybar: Forcibly restarts Waybar (kill then LaunchWaybar).
//
// Both are best-effort and print a brief summary.

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// countWaybarWindows returns the number of Waybar windows detected by hyprctl.
func countWaybarWindows() int {
	out, _ := exec.Command(
		"sh", "-lc",
		`hyprctl clients 2>/dev/null | grep -i 'class: waybar' | wc -l`,
	).CombinedOutput()
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// countMonitors returns the number of monitors reported by Hyprland.
// It uses text output for maximum compatibility and zero dependencies.
func countMonitors() int {
	// Prefer plain text to avoid strict JSON parsing and extra deps
	out, _ := exec.Command(
		"sh", "-lc",
		`hyprctl monitors 2>/dev/null | grep -c '^Monitor' || true`,
	).CombinedOutput()
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// SweepWaybar restarts Waybar only when the number of Waybar windows
// exceeds the number of monitors (with a small tolerance); otherwise it skips.
// It prints a short summary to stdout.
func SweepWaybar() {
	windows := countWaybarWindows()
	monitors := countMonitors()

	if monitors <= 0 {
		fmt.Println("Waybar sweep: monitors unknown; skipping")
		return
	}

	// Tolerance: allow up to monitors + 1 windows before considering it broken
	if windows > monitors+1 {
		fmt.Printf("Waybar sweep: windows=%d, monitors=%d -> restarting Waybar\n", windows, monitors)
		_ = exec.Command("pkill", "-x", "waybar").Run()
		_ = LaunchWaybar()
		return
	}

	fmt.Printf("Waybar sweep: windows=%d, monitors=%d -> skipping restart\n", windows, monitors)
}

// RestartWaybar forcibly restarts Waybar (kill + LaunchWaybar).
func RestartWaybar() {
	_ = exec.Command("pkill", "-x", "waybar").Run()
	_ = LaunchWaybar()
}
