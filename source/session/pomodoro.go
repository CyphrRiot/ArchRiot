package session

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// PomodoroClick handles Waybar Pomodoro click behavior:
// - Double-click (two quick clicks): reset timer
// - Single-click: toggle timer (handled via delayed worker)
//
// Pass selfPath to explicitly control the binary used for re-exec.
// If selfPath is empty, PomodoroClick will attempt to resolve the current executable.
func PomodoroClick(selfPath string) {
	const clickFile = "/tmp/waybar-tomato-click"
	const stateFile = "/tmp/waybar-tomato-timer.state"

	// If a click marker exists, treat as double-click (reset)
	if _, err := os.Stat(clickFile); err == nil {
		_ = os.Remove(clickFile)
		ts := time.Now().Unix()
		_ = os.WriteFile(stateFile, []byte(fmt.Sprintf("{\n    \"action\": \"reset\",\n    \"timestamp\": %d\n}\n", ts)), 0o644)
		return
	}

	// First click: create marker and spawn delayed toggle worker
	_ = os.WriteFile(clickFile, []byte(fmt.Sprintf("%d", time.Now().UnixNano())), 0o644)

	// Determine binary to re-exec for the delayed worker
	if selfPath == "" {
		if resolved, err := os.Executable(); err == nil {
			selfPath = resolved
		}
	}
	if selfPath == "" {
		// Could not resolve a path for the current executable; nothing else to do
		return
	}

	// Non-blocking spawn of delayed toggle worker
	_ = exec.Command(selfPath, "--pomodoro-delay-toggle").Start()
}

// PomodoroDelayToggle performs the delayed single-click toggle.
// If the click marker persists after a short delay, it's a single click → toggle.
func PomodoroDelayToggle() {
	const clickFile = "/tmp/waybar-tomato-click"
	const stateFile = "/tmp/waybar-tomato-timer.state"

	time.Sleep(500 * time.Millisecond)
	if _, err := os.Stat(clickFile); err == nil {
		_ = os.Remove(clickFile)
		ts := time.Now().Unix()
		_ = os.WriteFile(stateFile, []byte(fmt.Sprintf("{\n    \"action\": \"toggle\",\n    \"timestamp\": %d\n}\n", ts)), 0o644)
	}
}
