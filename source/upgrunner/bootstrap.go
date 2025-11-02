package upgrunner

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/logger"
	"archriot-installer/tui"
	"archriot-installer/upgrade"
	"archriot-installer/version"
)

// Bootstrap initializes logging and the TUI, runs the upgrade flow,
// and performs post-upgrade stabilization. It encapsulates the upgrade
// orchestration to keep the entrypoint minimal and delegation-only.
//
// Behavior:
// - Initializes version and logging.
// - Wires the TUI program to logger and upgrade packages.
// - Starts a Mullvad connection guard if connected at start.
// - Runs the interactive upgrade (upgrade.PromptAndRun) in a background goroutine.
// - On success: restarts hypridle and signals DoneMsg.
// - Ensures the Mullvad guard is stopped after the TUI exits.
//
// Returns an error only when initialization or TUI run fails.
// Upgrade failures are surfaced via TUI messages and do not return as an error here.
func Bootstrap() error {
	// Read version first (required by UI)
	if err := version.ReadVersion(); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Initialize logging (paired with defer)
	if err := logger.InitLogging(); err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}
	defer logger.CloseLogging()

	// Wire TUI helpers
	tui.SetVersionGetter(func() string { return version.Get() })
	tui.SetLogPathGetter(func() string { return logger.GetLogPath() })

	// Create and wire TUI program
	upModel := tui.NewInstallModel()
	upProgram := tea.NewProgram(upModel)
	logger.SetProgram(upProgram)
	upgrade.SetProgram(upProgram)

	// Optional Mullvad connection guard:
	// If the user starts the upgrade connected, auto-reconnect on drops.
	var stopMullvadGuard chan struct{}
	if _, err := exec.LookPath("mullvad"); err == nil {
		connectedAtStart := false
		if out, err := exec.Command("mullvad", "status").CombinedOutput(); err == nil {
			s := strings.ToLower(string(out))
			if strings.Contains(s, "connected") {
				connectedAtStart = true
			}
		}
		if connectedAtStart {
			stopMullvadGuard = make(chan struct{})
			upProgram.Send(tui.LogMsg("🛡️ Mullvad: connection guard active during upgrade"))

			go func(stop <-chan struct{}) {
				ticker := time.NewTicker(5 * time.Second)
				defer ticker.Stop()
				lastAttempt := time.Now().Add(-20 * time.Second)
				for {
					select {
					case <-stop:
						return
					case <-ticker.C:
						out, err := exec.Command("mullvad", "status").CombinedOutput()
						if err != nil {
							continue
						}
						ss := strings.ToLower(string(out))
						// If still connected, nothing to do
						if strings.Contains(ss, "connected") {
							continue
						}
						// Only attempt reconnect if an account is present
						if acc, err := exec.Command("mullvad", "account", "get").CombinedOutput(); err == nil &&
							strings.Contains(strings.ToLower(string(acc)), "mullvad account") {
							// Throttle reconnect attempts
							if time.Since(lastAttempt) > 10*time.Second {
								_ = exec.Command("mullvad", "connect").Start()
								lastAttempt = time.Now()
							}
						}
					}
				}
			}(stopMullvadGuard)
		}
	}

	// Start upgrade flow (non-blocking worker)
	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := upgrade.PromptAndRun(); err != nil {
			upProgram.Send(tui.LogMsg("❌ Upgrade failed: " + err.Error()))
			upProgram.Send(tui.FailureMsg{Error: "Upgrade failed"})
			// Stop Mullvad guard on failure
			if stopMullvadGuard != nil {
				close(stopMullvadGuard)
				stopMullvadGuard = nil
			}
			return
		}

		// Success: post-upgrade stabilization
		upProgram.Send(tui.LogMsg("🎉 Upgrade completed!"))
		logger.LogMessage("SUCCESS", "Upgrade completed")

		// Refresh idle manager only if running and present (avoid disrupting apps)
		if _, err := exec.LookPath("hypridle"); err == nil {
			if exec.Command("pgrep", "-x", "hypridle").Run() == nil {
				upProgram.Send(tui.LogMsg("🔄 Refreshing idle manager (hypridle)…"))
				_ = exec.Command("pkill", "-x", "hypridle").Run()
				cmd := exec.Command("hypridle")
				cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
				_ = cmd.Start()
			}
		}

		// If Brave is running, skip display enforcement and Waybar reload to avoid crashes; rely on login-time enforcement
		if exec.Command("pgrep", "-x", "brave").Run() == nil || exec.Command("pgrep", "-x", "brave-browser").Run() == nil {
			upProgram.Send(tui.LogMsg("🖥️ Skipping display enforcement and Waybar reload (Brave detected); relying on login-time enforcement"))
		} else {
			// Enforce display policy, then safely reload Waybar (sweep only if reload fails)
			upProgram.Send(tui.LogMsg("🖥️ Enforcing display policy (external preferred)…"))
			_ = exec.Command("sh", "-lc", "$HOME/.local/share/archriot/install/archriot --displays-enforce").Run()

			// Try a safe Waybar reload first; only sweep if no instance is running afterwards
			upProgram.Send(tui.LogMsg("🔄 Reloading Waybar…"))
			_ = exec.Command("sh", "-lc", "$HOME/.local/share/archriot/install/archriot --waybar-reload").Run()
			if exec.Command("pgrep", "-x", "waybar").Run() != nil {
				upProgram.Send(tui.LogMsg("🧹 Sweeping Waybar surfaces…"))
				_ = exec.Command("sh", "-lc", "$HOME/.local/share/archriot/install/archriot --waybar-sweep").Run()
			}
		}

		// Signal TUI completion
		upProgram.Send(tui.DoneMsg{})
	}()

	// Run TUI (blocking)
	if _, err := upProgram.Run(); err != nil {
		// Ensure Mullvad guard is stopped even on TUI error
		if stopMullvadGuard != nil {
			close(stopMullvadGuard)
			stopMullvadGuard = nil
		}
		return fmt.Errorf("TUI error: %w", err)
	}

	// Stop Mullvad guard after TUI exits
	if stopMullvadGuard != nil {
		close(stopMullvadGuard)
		stopMullvadGuard = nil
	}

	return nil
}
