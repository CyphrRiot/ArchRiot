package apps

// RunTelegram focuses an existing Telegram window or launches Telegram Desktop and
// focuses it when ready. It mirrors the original --telegram behavior without exiting
// the process. Returns 0 on success path (focused/launched best-effort), 1 when no
// launch path succeeds.
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RunTelegram executes the focus-or-launch flow for Telegram Desktop.
func RunTelegram(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	notify := func(title, msg string, ms int) {
		if have("notify-send") {
			_ = exec.Command("notify-send", "-t", fmt.Sprintf("%d", ms), title, msg).Start()
		}
	}

	// Minimal runtime logging to assist debugging
	home := os.Getenv("HOME")
	logDir := filepath.Join(home, ".cache", "archriot")
	_ = os.MkdirAll(logDir, 0o755)
	logFile := filepath.Join(logDir, "runtime.log")
	logAppend := func(msg string) {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return
		}
		defer f.Close()
		ts := time.Now().Format("2006-01-02 15:04:05")
		_, _ = f.WriteString(fmt.Sprintf("[%s] telegram: %s\n", ts, msg))
	}

	// Focus without scanning clients to avoid brittle parsing (use regex class match)
	focusTelegram := func() bool {
		// Broad match set: org.telegram.desktop (Flatpak/native), telegram-desktop (native),
		// TelegramDesktop (legacy), telegramdesktop/Telegram (rare)
		if exec.Command(
			"hyprctl", "dispatch", "focuswindow",
			"class:^(org\\.telegram\\.desktop|org\\.telegram\\.desktop\\.TelegramDesktop|org\\.telegram\\..*desktop.*|telegram-desktop|TelegramDesktop|telegramdesktop|Telegram)$",
		).Run() == nil {
			_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
			return true
		}
		return false
	}

	// Focus only if a Telegram window is present; otherwise proceed to launch
	if exec.Command(
		"sh", "-lc",
		"hyprctl clients 2>/dev/null | grep -qE 'class:\\s*(org\\.telegram\\.desktop(\\.TelegramDesktop)?|org\\.telegram\\..*desktop.*|telegram-desktop|TelegramDesktop|telegramdesktop|Telegram)\\b'",
	).Run() == nil {
		logAppend("window present; focusing")
		if focusTelegram() {
			return 0
		}
	}

	// Between attempts, wait briefly and check for a realized window to avoid duplicate spawns.
	waitFocus := func(loops int) bool {
		for i := 0; i < loops; i++ {
			time.Sleep(250 * time.Millisecond)
			if focusTelegram() {
				return true
			}
		}
		return false
	}

	// 1) Native binary
	if have("telegram-desktop") {
		logAppend("launching native telegram-desktop")
		notify("Telegram", "Launching Telegram Desktop…", 3000)
		_ = exec.Command("telegram-desktop").Start()
		if waitFocus(60) {
			return 0
		}
		time.Sleep(1500 * time.Millisecond)
		_ = exec.Command("telegram-desktop").Start()
		if waitFocus(60) {
			return 0
		}
		logAppend("native launch did not realize a window in time")
	}

	// 2) Desktop entry (gtk-launch), try common IDs and dynamically discovered ID
	if have("gtk-launch") {
		logAppend("trying gtk-launch candidates")
		notify("Telegram", "Launching Telegram (desktop)…", 3000)
		candidates := []string{"org.telegram.desktop", "telegram-desktop"}
		for _, id := range candidates {
			logAppend("gtk-launch " + id)
			_ = exec.Command("gtk-launch", id).Start()
			if waitFocus(60) {
				return 0
			}
			// Delayed retry if the window didn't realize on the first attempt
			time.Sleep(1500 * time.Millisecond)
			_ = exec.Command("gtk-launch", id).Start()
			if waitFocus(60) {
				return 0
			}
		}
		// Try to discover a Telegram desktop file dynamically
		out, _ := exec.Command(
			"sh", "-lc",
			"for d in ~/.local/share/applications /usr/local/share/applications /usr/share/applications; do "+
				"for f in \"$d\"/*[Tt]elegram*.desktop; do [ -f \"$f\" ] && basename \"${f%.desktop}\"; done; "+
				"done | head -n 1",
		).CombinedOutput()
		dyn := strings.TrimSpace(string(out))
		if dyn != "" {
			logAppend("gtk-launch discovered desktop id: " + dyn)
			_ = exec.Command("gtk-launch", dyn).Start()
			if waitFocus(60) {
				return 0
			}
			// One more attempt after a short delay for slower/cold systems
			time.Sleep(1500 * time.Millisecond)
			_ = exec.Command("gtk-launch", dyn).Start()
			if waitFocus(60) {
				return 0
			}
		}
	}

	// 3) Flatpak (attempt without info check; slower systems may need a retry)
	if have("flatpak") {
		logAppend("launching Flatpak org.telegram.desktop (no info probe)")
		notify("Telegram", "Launching Telegram Desktop (Flatpak)…", 3000)
		_ = exec.Command("flatpak", "run", "org.telegram.desktop").Start()
		if waitFocus(60) {
			return 0
		}
		time.Sleep(1500 * time.Millisecond)
		_ = exec.Command("flatpak", "run", "org.telegram.desktop").Start()
		if waitFocus(60) {
			return 0
		}
		logAppend("flatpak launch did not realize a window in time")
	}

	// 4) Parse .desktop Exec fallback
	// Last-resort: parse Exec lines from all Telegram*.desktop entries and try them in order
	if out, _ := exec.Command(
		"sh", "-lc",
		"for d in ~/.local/share/applications /usr/local/share/applications /usr/share/applications; do "+
			"for f in \"$d\"/*[Tt]elegram*.desktop; do "+
			"[ -f \"$f\" ] && grep -m1 '^Exec=' \"$f\" | sed -E 's/^Exec=//; s/%[fFuUdDnNickvm]//g'; "+
			"done; done",
	).CombinedOutput(); len(out) > 0 {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			cmdline := strings.TrimSpace(line)
			if cmdline == "" {
				continue
			}
			logAppend("exec from .desktop: " + cmdline)
			fields := strings.Fields(cmdline)
			if len(fields) == 0 {
				continue
			}
			_ = exec.Command(fields[0], fields[1:]...).Start()
			if waitFocus(90) {
				return 0
			}
			// Slow systems: retry the same command once after a short delay
			time.Sleep(1500 * time.Millisecond)
			_ = exec.Command(fields[0], fields[1:]...).Start()
			if waitFocus(90) {
				return 0
			}
		}
	}

	logAppend("all launch paths failed")
	notify("Telegram", "Unable to launch Telegram", 2000)
	return 1
}
