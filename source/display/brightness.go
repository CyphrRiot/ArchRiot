package display

// Package display provides brightness controls extracted from main.go.
//
// Usage (delegated from CLI):
//   archriot --brightness up
//   archriot --brightness down
//   archriot --brightness set <0-100>
//   archriot --brightness get
//
// Behavior:
// - Uses brightnessctl if available; otherwise prints an error and returns non-zero.
// - For up/down/set, shows a desktop notification (if notify-send is available).
// - Never calls os.Exit; returns an exit code for the caller to handle.

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Run executes the brightness subcommands and returns an exit code (0 = success).
func Run(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	notify := func(title, msg, icon string) {
		// Best-effort notifications; safe if not installed
		if have("makoctl") {
			_ = exec.Command("makoctl", "dismiss", "--all").Run()
		}
		if have("notify-send") {
			_ = exec.Command("notify-send",
				"--replace-id=9999",
				"--app-name=Brightness Control",
				"--urgency=normal",
				"--icon", icon,
				title, msg,
			).Start()
		}
	}

	readPct := func() string {
		out, err := exec.Command("sh", "-lc", "brightnessctl -m | cut -d, -f4 | tr -d '%'").Output()
		if err != nil {
			return "0"
		}
		return strings.TrimSpace(string(out))
	}

	iconFor := func(pct string) string {
		n, err := strconv.Atoi(strings.TrimSpace(pct))
		if err != nil {
			return "brightness-medium"
		}
		switch {
		case n >= 75:
			return "brightness-high"
		case n >= 50:
			return "brightness-medium"
		case n >= 25:
			return "brightness-low"
		default:
			return "brightness-min"
		}
	}

	usage := func() int {
		fmt.Fprintln(os.Stderr, "Usage: archriot --brightness [up|down|set <0-100>|get]")
		return 1
	}

	if !have("brightnessctl") {
		notify("Brightness Error", "brightnessctl not found", "dialog-error")
		fmt.Fprintln(os.Stderr, "Error: brightnessctl is not installed")
		return 1
	}

	if len(args) < 1 {
		return usage()
	}

	switch args[0] {
	case "up":
		_ = exec.Command("brightnessctl", "set", "5%+").Run()
		pct := readPct()
		notify("Brightness", pct+"%", iconFor(pct))
		return 0

	case "down":
		_ = exec.Command("brightnessctl", "set", "5%-").Run()
		pct := readPct()
		notify("Brightness", pct+"%", iconFor(pct))
		return 0

	case "set":
		if len(args) < 2 {
			return usage()
		}
		val := strings.TrimSpace(args[1])
		if val == "" {
			return usage()
		}
		_ = exec.Command("brightnessctl", "set", val+"%").Run()
		pct := readPct()
		notify("Brightness", pct+"%", iconFor(pct))
		return 0

	case "get":
		fmt.Println(readPct())
		return 0

	default:
		return usage()
	}
}
