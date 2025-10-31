package apps

// RunZed focuses an existing Zed window or launches Zed with Wayland-friendly
// settings and focuses it when ready. This mirrors the original --zed behavior
// from main.go without exiting the process. It performs best-effort actions and
// returns 0 regardless of launch/focus outcome.
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RunZed executes the focus-or-launch flow for Zed.
func RunZed(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	focusZed := func() bool {
		if exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(dev\\.zed\\.Zed|dev\\.zed\\.Zed-Preview)$").Run() == nil {
			_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
			return true
		}
		return false
	}

	// If a Zed window exists, focus and return
	if focusZed() {
		return 0
	}

	// Minimal runtime logging
	home := os.Getenv("HOME")
	logDir := filepath.Join(home, ".cache", "archriot")
	_ = os.MkdirAll(logDir, 0o755)
	logFile := filepath.Join(logDir, "runtime.log")
	logAppend := func(msg string) {
		if f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			defer f.Close()
			ts := time.Now().Format("2006-01-02 15:04:05")
			_, _ = f.WriteString("[" + ts + "] zed: " + msg + "\n")
		}
	}

	// Detect Intel GPU (ANV) to avoid Vulkan crashes; prefer WGPU_BACKEND=gl
	detectIntel := func() bool {
		// Prefer lspci if available (broadest signal)
		if have("lspci") {
			if out, err := exec.Command("sh", "-lc", "lspci -nnk | grep -iE 'vga|3d|display' | head -n1").CombinedOutput(); err == nil {
				s := strings.ToLower(string(out))
				if strings.Contains(s, "intel") {
					return true
				}
			}
		}
		// Fallback: presence of Mesa Intel DRI driver
		if _, err := os.Stat("/usr/lib/dri/iris_dri.so"); err == nil {
			return true
		}
		return false
	}
	intel := detectIntel()
	if intel {
		logAppend("Intel GPU detected; preferring WGPU_BACKEND=gl for Zed to avoid Vulkan/ANV instability")
	}

	// helper: short synchronous focus wait (reused across attempts)
	attemptFocus := func(loops int) bool {
		for i := 0; i < loops; i++ {
			time.Sleep(250 * time.Millisecond)
			if focusZed() {
				return true
			}
		}
		return false
	}

	launched := false

	// 1) zeditor first (more reliable on some systems)
	if !launched && have("zeditor") {
		// Build env for child only (no global environment changes)
		argsEnv := []string{}
		if intel {
			argsEnv = append(argsEnv, "WGPU_BACKEND=gl")
			logAppend("launch: env WGPU_BACKEND=gl zeditor (Wayland env applied)")
		} else {
			logAppend("launch: zeditor (Wayland env applied)")
		}
		argsEnv = append(argsEnv,
			"WAYLAND_DISPLAY="+os.Getenv("WAYLAND_DISPLAY"),
			"GDK_BACKEND=wayland",
			"QT_QPA_PLATFORM=wayland",
			"SDL_VIDEODRIVER=wayland",
			"zeditor",
		)
		_ = exec.Command("env", argsEnv...).Start()
		if attemptFocus(60) {
			launched = true
		} else {
			// Delayed retry
			time.Sleep(1200 * time.Millisecond)
			_ = exec.Command("env", argsEnv...).Start()
			if attemptFocus(60) {
				launched = true
			} else {
				logAppend("zeditor did not realize a window in time")
			}
		}
	}

	// 2) zed (native)
	if !launched && have("zed") {
		if intel {
			logAppend("launch: env WGPU_BACKEND=gl zed")
			_ = exec.Command("env", "WGPU_BACKEND=gl", "zed").Start()
		} else {
			logAppend("launch: zed")
			_ = exec.Command("zed").Start()
		}
		if attemptFocus(60) {
			launched = true
		} else {
			// Delayed retry
			time.Sleep(1200 * time.Millisecond)
			if intel {
				_ = exec.Command("env", "WGPU_BACKEND=gl", "zed").Start()
			} else {
				_ = exec.Command("zed").Start()
			}
			if attemptFocus(60) {
				launched = true
			} else {
				logAppend("zed did not realize a window in time")
			}
		}
	}

	// 3) Flatpak (pass GL backend when Intel is detected)
	if !launched && have("flatpak") {
		if intel {
			logAppend("launch: flatpak run --env=WGPU_BACKEND=gl dev.zed.Zed")
			_ = exec.Command("flatpak", "run", "--env=WGPU_BACKEND=gl", "dev.zed.Zed").Start()
		} else {
			logAppend("launch: flatpak run dev.zed.Zed")
			_ = exec.Command("flatpak", "run", "dev.zed.Zed").Start()
		}
		if attemptFocus(60) {
			launched = true
		} else {
			// Delayed retry
			time.Sleep(1200 * time.Millisecond)
			if intel {
				_ = exec.Command("flatpak", "run", "--env=WGPU_BACKEND=gl", "dev.zed.Zed").Start()
			} else {
				_ = exec.Command("flatpak", "run", "dev.zed.Zed").Start()
			}
			if attemptFocus(60) {
				launched = true
			} else {
				logAppend("flatpak zed did not realize a window in time")
			}
		}
	}

	// Async focus when window appears
	if launched {
		go func() {
			for i := 0; i < 40; i++ {
				time.Sleep(250 * time.Millisecond)
				if focusZed() {
					return
				}
			}
		}()
	}

	return 0
}
