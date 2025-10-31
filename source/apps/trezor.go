package apps

// RunTrezor focuses an existing Trezor Suite window or launches it using the best
// available method (native > Flatpak > AppImage), then focuses it once ready.
// Mirrors the original --trezor behavior without exiting the process.
import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// RunTrezor executes the focus-or-launch flow for Trezor Suite.
// Always returns 0 (best-effort).
func RunTrezor(args []string) int {
	home := os.Getenv("HOME")
	appImage := filepath.Join(home, ".local", "bin", "trezor-suite.AppImage")

	have := func(name string) bool {
		_, err := exec.LookPath(name)
		return err == nil
	}

	trezorWindowPresent := func() bool {
		// Check via hyprctl clients (no jq dependency)
		if err := exec.Command("sh", "-lc", "hyprctl clients 2>/dev/null | grep -qE 'class:\\s*Trezor Suite\\b'").Run(); err == nil {
			return true
		}
		return false
	}

	focusTrezor := func() bool {
		// Focus by class match
		if err := exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(Trezor Suite)$").Run(); err != nil {
			return false
		}
		_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
		return true
	}

	// If window exists, focus and exit
	if trezorWindowPresent() && focusTrezor() {
		return 0
	}

	// Otherwise launch best available target
	if have("trezor-suite") {
		_ = exec.Command("trezor-suite").Start()
	} else if have("flatpak") && exec.Command("flatpak", "info", "--show-commit", "com.trezor.TrezorSuite").Run() == nil {
		_ = exec.Command("flatpak", "run", "com.trezor.TrezorSuite").Start()
	} else if st, err := os.Stat(appImage); err == nil && !st.IsDir() {
		_ = exec.Command(appImage).Start()
	}

	// Small async wait to focus when it appears (best effort)
	go func() {
		for i := 0; i < 30; i++ {
			time.Sleep(250 * time.Millisecond)
			if trezorWindowPresent() {
				_ = exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(Trezor Suite)$").Run()
				_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
				break
			}
		}
	}()

	return 0
}
