package waybar

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// EmitMullvadStatus returns VPN status JSON for waybar.
// Outputs JSON with text, tooltip, and class fields.
func EmitMullvadStatus(args []string) int {
	// Run mullvad status with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "mullvad", "status")
	out, err := cmd.Output()

	connected := false
	location := ""

	if err == nil {
		output := strings.ToLower(string(out))
		if strings.Contains(output, "connected") {
			connected = true
			// Try to extract location from Relay: line
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "Relay:") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						relay := parts[1] // e.g., us-sjc-wg-501
						// Extract middle token (sjc)
						relayParts := strings.Split(relay, "-")
						if len(relayParts) >= 2 {
							location = strings.ToUpper(relayParts[1])
						}
					}
					break
				}
			}
		}
	}

	var text, tooltip, class string

	if connected {
		if location == "" {
			location = "VPN"
		}
		text = "󰌆 " + location
		tooltip = "Mullvad VPN Connected: " + location + "\nRight-click to disconnect"
		class = "mullvad-connected"
	} else {
		text = "󰌉"
		tooltip = "Mullvad VPN Disconnected\nLeft-click to connect"
		class = "mullvad-disconnected"
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

// EmitMullvadClick handles click events for Mullvad VPN toggle.
// Takes click type as argument: "left" or "right"
func EmitMullvadClick(args []string) int {
	clickType := "left"
	if len(args) > 0 {
		clickType = args[0]
	}
	// Get current status
	cmd := exec.Command("mullvad", "status")
	out, err := cmd.Output()

	isConnected := err == nil && strings.Contains(strings.ToLower(string(out)), "connected")

	var actionCmd *exec.Cmd

	if isConnected {
		// Currently connected
		if clickType == "right" {
			// Right-click: disconnect
			actionCmd = exec.Command("mullvad", "disconnect")
		}
		// Left-click when connected: do nothing
	} else {
		// Currently disconnected
		if clickType == "left" {
			// Left-click: connect
			actionCmd = exec.Command("mullvad", "connect")
		}
		// Right-click when disconnected: do nothing
	}

	if actionCmd != nil {
		_ = actionCmd.Run()
	}

	return 0
}
