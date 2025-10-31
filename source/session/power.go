package session

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// PowerMenu shows a simple system power menu.
// Prefers fuzzel if available; falls back to a TTY prompt.
// Returns 0 always (best-effort).
func PowerMenu() int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	choose := func() string {
		if have("fuzzel") {
			items := []string{"Lock", "Suspend", "Reboot", "Power Off", "Logout", "Cancel", "Control Panel"}
			input := strings.Join(items, "\n")
			linesArg := fmt.Sprintf("--lines=%d", len(items))
			cmd := exec.Command("fuzzel", "--dmenu", "--prompt=Power: ", "--width=30", linesArg)
			cmd.Stdin = strings.NewReader(input)
			if out, err := cmd.Output(); err == nil {
				return strings.TrimSpace(string(out))
			}
		}
		// Fallback TTY prompt
		fmt.Println("Power menu:")
		fmt.Println("1) Lock")
		fmt.Println("2) Suspend")
		fmt.Println("3) Reboot")
		fmt.Println("4) Power Off")
		fmt.Println("5) Logout")
		fmt.Println("6) Cancel")
		fmt.Println("7) Control Panel")
		fmt.Print("Select [1-7]: ")
		reader := bufio.NewReader(os.Stdin)
		s, _ := reader.ReadString('\n')
		return strings.TrimSpace(s)
	}

	sel := choose()
	act := sel
	switch sel {
	case "1", "Lock":
		act = "Lock"
	case "2", "Suspend":
		act = "Suspend"
	case "3", "Reboot":
		act = "Reboot"
	case "4", "Power Off":
		act = "Power Off"
	case "5", "Logout":
		act = "Logout"
	case "7", "Control Panel":
		act = "Control Panel"
	default:
		return 0
	}

	switch act {
	case "Lock":
		if have("hyprlock") {
			cmd := exec.Command("hyprlock")
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			_ = cmd.Start()
		} else {
			_ = exec.Command("loginctl", "lock-session").Run()
		}
	case "Suspend":
		_ = exec.Command("systemctl", "suspend").Start()
	case "Reboot":
		_ = exec.Command("systemctl", "reboot").Start()
	case "Power Off":
		_ = exec.Command("systemctl", "poweroff").Start()
	case "Logout":
		{
			if have("hyprctl") {
				_ = exec.Command("hyprctl", "dispatch", "exit").Start()
			} else {
				user := os.Getenv("USER")
				if user == "" {
					if out, err := exec.Command("sh", "-lc", "id -un").Output(); err == nil {
						user = strings.TrimSpace(string(out))
					}
				}
				if user != "" {
					_ = exec.Command("loginctl", "terminate-user", user).Start()
				}
			}
		}
	case "Control Panel":
		{
			cp := os.Getenv("HOME") + "/.local/share/archriot/config/bin/archriot-control-panel"
			if _, err := os.Stat(cp); err == nil {
				cmd := exec.Command(cp)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
				_ = cmd.Start()
			}
		}
	}
	return 0
}
