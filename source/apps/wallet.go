package apps

// RunWallet focuses an existing crypto wallet window (Trezor Suite or Ledger Live)
// or launches whichever is installed, then focuses it once the window appears.
// This mirrors the original --wallet behavior from main.go without exiting the process.
// Always returns 0 (best-effort).
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RunWallet executes the focus-or-launch flow for Trezor Suite or Ledger Live.
func RunWallet(args []string) int {
	home := os.Getenv("HOME")
	trezorAppImage := filepath.Join(home, ".local", "bin", "trezor-suite.AppImage")
	ledgerAppImage := filepath.Join(home, ".local", "bin", "ledger-live.AppImage")

	// Wallet logging (helpful for runtime diagnostics)
	trezorLogDir := filepath.Join(home, ".cache", "archriot")
	_ = os.MkdirAll(trezorLogDir, 0o755)
	trezorLogFile := filepath.Join(trezorLogDir, "runtime.log")
	logAppend := func(msg string) {
		f, err := os.OpenFile(trezorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err == nil {
			defer f.Close()
			ts := time.Now().Format("2006-01-02 15:04:05")
			_, _ = f.WriteString(fmt.Sprintf("[%s] %s\n", ts, msg))
		}
	}

	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	notify := func(title, msg string) {
		if have("notify-send") {
			_ = exec.Command("notify-send", "-t", "1500", title, msg).Start()
		}
	}
	notifyDur := func(title, msg string, ms int) {
		if have("notify-send") {
			_ = exec.Command("notify-send", "-t", fmt.Sprintf("%d", ms), title, msg).Start()
		}
	}

	hyprClientsContains := func(substr string) bool {
		out, err := exec.Command("hyprctl", "clients").Output()
		return err == nil && strings.Contains(string(out), substr)
	}

	focusClass := func(classRegex string) bool {
		if err := exec.Command("hyprctl", "dispatch", "focuswindow", "class:^("+classRegex+")$").Run(); err != nil {
			return false
		}
		_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
		return true
	}

	// Focus if a wallet window exists (Trezor first, then Ledger)
	if hyprClientsContains("class: Trezor Suite") && focusClass("Trezor Suite") {
		logAppend("Focusing Trezor Suite")
		notify("Wallet", "Focusing Trezor Suite…")
		return 0
	}
	if (hyprClientsContains("class: Ledger Live") || hyprClientsContains("class: ledger live")) && focusClass("Ledger Live|ledger live") {
		logAppend("Focusing Ledger Live")
		notify("Wallet", "Focusing Ledger Live…")
		return 0
	}

	// Launch whichever is installed, preferring Trezor
	launched := false

	// Trezor: native > Flatpak > AppImage
	switch {
	case have("trezor-suite"):
		logAppend("Launching Trezor Suite (native)")
		notify("Wallet", "Opening Trezor Suite…")
		_ = exec.Command("trezor-suite").Start()
		launched = true
	case have("flatpak") && exec.Command("flatpak", "info", "--show-commit", "com.trezor.TrezorSuite").Run() == nil:
		logAppend("Launching Trezor Suite (Flatpak)")
		notify("Wallet", "Opening Trezor Suite (Flatpak)…")
		_ = exec.Command("flatpak", "run", "com.trezor.TrezorSuite").Start()
		launched = true
	default:
		if st, err := os.Stat(trezorAppImage); err == nil && !st.IsDir() {
			logAppend("Launching Trezor Suite (AppImage)")
			notify("Wallet", "Opening Trezor Suite (AppImage)…")
			_ = exec.Command(trezorAppImage).Start()
			launched = true
		}
	}

	// If Trezor wasn’t launched, try Ledger: native > Flatpak > AppImage
	if !launched {
		switch {
		case have("ledger-live") || have("ledger-live-desktop"):
			logAppend("Launching Ledger Live (native)")
			notifyDur("Wallet", "Opening Ledger Live…", 6000)
			if have("ledger-live") {
				_ = exec.Command("ledger-live").Start()
			} else {
				_ = exec.Command("ledger-live-desktop").Start()
			}
			launched = true
		case have("flatpak") && exec.Command("flatpak", "info", "--show-commit", "com.ledgerhq.LedgerLive").Run() == nil:
			logAppend("Launching Ledger Live (Flatpak)")
			notifyDur("Wallet", "Opening Ledger Live (Flatpak)…", 6000)
			_ = exec.Command("flatpak", "run", "com.ledgerhq.LedgerLive").Start()
			launched = true
		default:
			if st, err := os.Stat(ledgerAppImage); err == nil && !st.IsDir() {
				logAppend("Launching Ledger Live (AppImage)")
				notifyDur("Wallet", "Opening Ledger Live (AppImage)…", 6000)
				_ = exec.Command(ledgerAppImage).Start()
				launched = true
			}
		}
	}

	// Async focus once spawned (best-effort)
	if launched {
		go func() {
			for i := 0; i < 40; i++ {
				time.Sleep(250 * time.Millisecond)
				if hyprClientsContains("class: Trezor Suite") {
					_ = exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(Trezor Suite)$").Run()
					_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
					logAppend("Focused Trezor Suite after launch")
					return
				}
				if hyprClientsContains("class: Ledger Live") || hyprClientsContains("class: ledger live") {
					_ = exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(Ledger Live|ledger live)$").Run()
					_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
					logAppend("Focused Ledger Live after launch")
					return
				}
			}
		}()
	}

	return 0
}
