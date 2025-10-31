package diagnostics

// WifiPowerSaveCheck performs a read-only check of NetworkManager's Wi‑Fi
// power-save configuration and prints the current runtime power-save state
// for the first detected wireless interface (if any).
//
// It mirrors the legacy `--wifi-powersave-check` CLI behavior and returns
// an exit code (0 = success). The function never exits the process.
import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// WifiPowerSaveCheck runs the Wi‑Fi power-save diagnostics and prints results.
func WifiPowerSaveCheck() int {
	fmt.Println("WiFi power-save diagnostics")

	// Check NetworkManager drop-in for wifi.powersave
	psConf := "/etc/NetworkManager/conf.d/40-wifi-powersave.conf"
	if b, err := os.ReadFile(psConf); err == nil {
		val := "(not set)"
		for _, ln := range strings.Split(string(b), "\n") {
			s := strings.TrimSpace(ln)
			if strings.HasPrefix(s, "wifi.powersave=") {
				val = strings.TrimSpace(strings.TrimPrefix(s, "wifi.powersave="))
				break
			}
		}
		fmt.Printf("%-24s %s\n", "drop-in:", psConf)
		fmt.Printf("%-24s %s\n", "wifi.powersave:", val)
		if val != "2" && val != "(not set)" {
			fmt.Println("Tip: set wifi.powersave=2 to avoid Wi‑Fi power saving")
			fmt.Println("e.g., echo -e \"[connection]\\nwifi.powersave=2\" | sudo tee /etc/NetworkManager/conf.d/40-wifi-powersave.conf && sudo systemctl reload NetworkManager")
		}
	} else {
		fmt.Printf("%-24s %s\n", "drop-in:", "missing")
		fmt.Println("Tip: create /etc/NetworkManager/conf.d/40-wifi-powersave.conf with wifi.powersave=2")
		fmt.Println("e.g., echo -e \"[connection]\\nwifi.powersave=2\" | sudo tee /etc/NetworkManager/conf.d/40-wifi-powersave.conf && sudo systemctl reload NetworkManager")
	}

	// Runtime power save state on a wireless iface (if any)
	iface := ""
	if out, err := exec.Command("sh", "-lc", "iw dev | awk '/Interface/ {print $2; exit}'").CombinedOutput(); err == nil {
		iface = strings.TrimSpace(string(out))
	}
	if iface != "" {
		if out, err := exec.Command("sh", "-lc", "iw dev "+iface+" get power_save 2>/dev/null | awk '{print tolower($0)}'").CombinedOutput(); err == nil {
			state := strings.TrimSpace(string(out))
			if state == "" {
				state = "unknown"
			}
			fmt.Printf("%-24s %s (%s)\n", "runtime power_save:", state, iface)
			if strings.Contains(state, "on") {
				fmt.Println("Tip: runtime power save is ON; temporarily disable with:")
				fmt.Println("sudo iw dev " + iface + " set power_save off")
			}
		}
	} else {
		fmt.Println("No wireless interface detected")
	}

	return 0
}
