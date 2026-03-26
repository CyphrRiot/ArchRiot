package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// EmitRecordingStatus returns recording status JSON for waybar.
// Outputs JSON with text, tooltip, and class fields.
func EmitRecordingStatus(args []string) int {
	active := false

	// Try pw-dump first (PipeWire >= 0.3.61)
	cmd := exec.Command("pw-dump", "-N")
	out, err := cmd.Output()
	if err == nil {
		output := strings.ToLower(string(out))
		if strings.Contains(output, `"node.name":"xdpw_`) ||
			strings.Contains(output, `"application.name":"xdg-desktop-portal`) ||
			strings.Contains(output, "screencast") ||
			strings.Contains(output, "screen-cast") {
			active = true
		}
	}

	// Fallback: check if Kooha is running
	if !active {
		cmd = exec.Command("pgrep", "-x", "kooha")
		err = cmd.Run()
		active = err == nil
	}

	var text, tooltip, class string

	if active {
		text = "●"
		tooltip = "Screen recording (click to stop)"
		class = "recording"
	} else {
		text = ""
		tooltip = "No screen recording"
		class = ""
	}

	output := Output{
		Text:    text,
		Tooltip: tooltip,
		Class:   class,
	}

	jsonOut, err := json.Marshal(output)
	if err != nil {
		return 1
	}

	fmt.Println(string(jsonOut))
	return 0
}

// EmitRecordingClick handles click events to stop Kooha recording.
// Takes click type as argument: "left" or "right"
func EmitRecordingClick(args []string) int {
	// Try SIGINT first (let Kooha finalize cleanly), then TERM as fallback
	pid, err := findProcessPID("kooha")
	if err == nil {
		syscall.Kill(pid, syscall.SIGINT)
		// Give it a moment, then try SIGTERM if still running
		go func() {
			exec.Command("sleep", "2").Run()
			syscall.Kill(pid, syscall.SIGTERM)
		}()
	}

	return 0
}

// findProcessPID returns the PID of a process by name, or error if not found
func findProcessPID(name string) (int, error) {
	cmd := exec.Command("pidof", name)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	pidStr := strings.TrimSpace(string(out))
	if pidStr == "" {
		return 0, fmt.Errorf("process not found")
	}
	var pid int
	_, err = fmt.Sscanf(pidStr, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// EmitUpdateCheck returns update status JSON for waybar.
// Checks local vs remote version and outputs appropriate state.
func EmitUpdateCheck(args []string) int {
	// Read local version
	localVersion := readFileOrDefault(os.ExpandEnv("$HOME/.local/share/archriot/VERSION"), "unknown")

	// Read remote version with timeout
	remoteVersion := fetchRemoteVersion("https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION", 10)

	var text, tooltip, class string

	// Handle error states
	if localVersion == "unknown" || remoteVersion == "unknown" {
		text = "-"
		tooltip = "Update check unavailable"
		class = "update-none"
	} else if isNewerVersion(localVersion, remoteVersion) {
		// Update available - check if user has seen it
		stateFile := os.ExpandEnv("$HOME/.cache/archriot/update-state")
		seenVersion := readFileOrDefault(stateFile, "")

		if seenVersion == remoteVersion {
			// User has seen this version
			text = "󰏖"
			tooltip = fmt.Sprintf("ArchRiot update available (seen)\nCurrent: %s\nAvailable: %s", localVersion, remoteVersion)
			class = "update-seen"
		} else {
			// New version not seen yet
			text = "󰚰"
			tooltip = fmt.Sprintf("ArchRiot update available!\nCurrent: %s\nAvailable: %s", localVersion, remoteVersion)
			class = "update-available"
		}
	} else {
		// Up to date
		text = "-"
		tooltip = fmt.Sprintf("ArchRiot is up to date\nCurrent: %s", localVersion)
		class = "update-none"
	}

	output := Output{
		Text:    text,
		Tooltip: tooltip,
		Class:   class,
	}

	jsonOut, err := json.Marshal(output)
	if err != nil {
		return 1
	}

	fmt.Println(string(jsonOut))
	return 0
}

// EmitUpdateClick handles click events for update notification.
func EmitUpdateClick(args []string) int {
	localVersion := readFileOrDefault(os.ExpandEnv("$HOME/.local/share/archriot/VERSION"), "unknown")
	remoteVersion := fetchRemoteVersion("https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION", 10)

	if localVersion != "unknown" && remoteVersion != "unknown" && isNewerVersion(localVersion, remoteVersion) {
		// Mark as clicked/seen
		stateFile := os.ExpandEnv("$HOME/.cache/archriot/update-state")
		os.MkdirAll(os.ExpandEnv("$HOME/.cache/archriot"), 0755)
		os.WriteFile(stateFile, []byte(remoteVersion), 0644)

		// Launch upgrade
		go func() {
			exec.Command("notify-send", "-t", "5000", "Launching Upgrade...", "Starting ArchRiot upgrade process...").Run()
			exec.Command(os.ExpandEnv("$HOME/.local/share/archriot/config/bin/version-check"), "--gui").Run()
		}()
	} else {
		// Already up to date
		go func() {
			exec.Command("notify-send", "-t", "2000", "ArchRiot Up to Date", fmt.Sprintf("Version %s is the latest", localVersion)).Run()
		}()
	}

	return 0
}

// Helper functions

func readFileOrDefault(path, defaultVal string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultVal
	}
	return strings.TrimSpace(string(data))
}

func fetchRemoteVersion(url string, timeoutSec int) string {
	cmd := exec.Command("timeout", fmt.Sprintf("%d", timeoutSec), "curl", "-s", url)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func isNewerVersion(local, remote string) bool {
	if local == "unknown" || remote == "unknown" {
		return false
	}

	localParts := strings.Split(local, ".")
	remoteParts := strings.Split(remote, ".")

	for i := 0; i < 3; i++ {
		localPart := 0
		remotePart := 0
		if i < len(localParts) {
			fmt.Sscanf(localParts[i], "%d", &localPart)
		}
		if i < len(remoteParts) {
			fmt.Sscanf(remoteParts[i], "%d", &remotePart)
		}

		if remotePart > localPart {
			return true
		}
		if localPart > remotePart {
			return false
		}
	}

	return false // Versions are equal
}
