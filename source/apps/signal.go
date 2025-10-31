package apps

// RunSignal focuses an existing Signal window or launches Signal Desktop and
// focuses it when ready. It mirrors the original --signal behavior from main.go
// without exiting the process. Always returns 0.
import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RunSignal executes the focus-or-launch flow for Signal Desktop.
func RunSignal(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	notify := func(title, msg string, ms int) {
		if have("notify-send") {
			_ = exec.Command("notify-send", "-t", fmt.Sprintf("%d", ms), title, msg).Start()
		}
	}

	hyprClientsContains := func(substr string) bool {
		out, err := exec.Command("hyprctl", "clients").Output()
		return err == nil && strings.Contains(string(out), substr)
	}

	focusSignal := func() bool {
		if exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(Signal)$").Run() == nil {
			_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
			return true
		}
		if exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(signal)$").Run() == nil {
			_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
			return true
		}
		return false
	}

	// If a Signal window exists, focus it and return
	if hyprClientsContains("class: Signal") || hyprClientsContains("class: signal") {
		if focusSignal() {
			return 0
		}
	}

	// Otherwise launch Signal (Wayland/Ozone) and focus when ready
	notify("Signal", "Launching Signal Desktop…", 3000)
	_ = exec.Command(
		"env",
		"GDK_SCALE=1",
		"signal-desktop",
		"--ozone-platform=wayland",
		"--enable-features=UseOzonePlatform",
	).Start()

	// Async focus retries for a short period while the window appears
	go func() {
		for i := 0; i < 20; i++ {
			time.Sleep(250 * time.Millisecond)
			if focusSignal() {
				return
			}
		}
	}()

	return 0
}
