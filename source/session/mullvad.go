package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MullvadStartup safely starts the Mullvad GUI minimized to tray when appropriate.
// Behavior:
// - No-ops if Mullvad CLI or GUI is not installed.
// - Skips if GUI is already running.
// - Starts only when an account is present and auto-connect is enabled.
// - Ensures startMinimized=true in GUI settings when present.
// - Adds small delays to allow the tray to be ready and verifies startup.
func MullvadStartup() {
	home := os.Getenv("HOME")

	// Append to runtime log (best-effort)
	logAppend := func(msg string) {
		logDir := filepath.Join(home, ".cache", "archriot")
		_ = os.MkdirAll(logDir, 0o755)
		logFile := filepath.Join(logDir, "runtime.log")
		if f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			defer f.Close()
			ts := time.Now().Format("15:04:05")
			_, _ = f.WriteString("[" + ts + "] " + msg + "\n")
		}
	}

	// CLI presence check
	if _, err := exec.LookPath("mullvad"); err != nil {
		// Not installed; silently exit
		return
	}

	// GUI presence check
	guiPath := "/opt/Mullvad VPN/mullvad-gui"
	if _, err := os.Stat(guiPath); err != nil {
		// GUI not installed; silently exit
		return
	}

	// Already running? Exit quietly
	if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
		logAppend("Mullvad GUI already running, skipping startup")
		fmt.Println("Mullvad GUI already running")
		return
	}

	// Account check: consider "status" Connected/Disconnected/Blocked as account-present,
	// or "account get" containing an account number.
	accountPresent := false
	if out, err := exec.Command("mullvad", "status").CombinedOutput(); err == nil {
		s := strings.ToLower(string(out))
		if strings.Contains(s, "connected") || strings.Contains(s, "disconnected") || strings.Contains(s, "blocked") {
			accountPresent = true
		}
	}
	if !accountPresent {
		if out, err := exec.Command("mullvad", "account", "get").CombinedOutput(); err == nil {
			if strings.Contains(string(out), "Mullvad account:") {
				accountPresent = true
			}
		}
	}
	if !accountPresent {
		logAppend("No Mullvad account configured, skipping GUI startup")
		return
	}

	// Respect Mullvad auto-connect: skip GUI when auto-connect is OFF/disabled
	if out, err := exec.Command("mullvad", "auto-connect", "get").CombinedOutput(); err == nil {
		s := strings.ToLower(string(out))
		if !(strings.Contains(s, "on") || strings.Contains(s, "enabled")) {
			logAppend("Mullvad auto-connect is OFF; skipping GUI startup")
			return
		}
	}

	// Ensure startMinimized=true in GUI settings if present
	settingsPath := filepath.Join(home, ".config", "Mullvad VPN", "gui_settings.json")
	if b, err := os.ReadFile(settingsPath); err == nil {
		type guiSettings struct {
			StartMinimized bool `json:"startMinimized"`
		}
		var gs guiSettings
		// If parsing fails, ignore and continue with defaults
		if json.Unmarshal(b, &gs) == nil {
			if !gs.StartMinimized {
				gs.StartMinimized = true
				if nb, err := json.MarshalIndent(gs, "", "  "); err == nil {
					_ = os.WriteFile(settingsPath, nb, 0o644)
				}
			}
		}
	}

	// Delay to allow tray to be ready
	time.Sleep(10 * time.Second)

	// Final guard: if started in the meantime, exit
	if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
		return
	}

	// Start minimized; detach (no need to wait)
	logAppend("Starting Mullvad GUI minimized to tray")
	cmd := exec.Command(guiPath, "--minimize-to-tray")
	_ = cmd.Start()

	// Brief verify window
	time.Sleep(2 * time.Second)
	if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
		logAppend("✓ Mullvad GUI started successfully")
	} else {
		logAppend("⚠ Mullvad GUI did not appear to start")
	}
}
