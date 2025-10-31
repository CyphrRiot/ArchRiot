package session

// Waybar diagnostics helper. This inspects common multi-launch sources and
// prints a structured report to stdout so users can remove redundant starters
// and rely on a single, race-proof launcher (archriot --waybar-launch).
//
// It is safe to call at any time. Errors are printed as best-effort notes.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DiagnoseWaybar prints a structured diagnostic report to stdout:
// - Current Waybar processes (PID, PPID, cmd)
// - Parent process command for each Waybar (who launched it)
// - Hyprland exec-once lines that may launch Waybar
// - systemd --user unit files that reference Waybar
// - Autostart .desktop entries referencing Waybar
func DiagnoseWaybar() {
	fmt.Println("Waybar Diagnostics — launch sources and processes")

	// 1) Processes: list all Waybar instances with PPID and cmd
	fmt.Println("\n[Processes]")
	if out, err := exec.Command("sh", "-lc", "ps -o pid=,ppid=,cmd= -C waybar 2>/dev/null || true").CombinedOutput(); err == nil {
		s := strings.TrimSpace(string(out))
		if s == "" {
			fmt.Println("  No waybar processes found")
		} else {
			lines := strings.Split(s, "\n")
			for _, ln := range lines {
				fmt.Println(" ", ln)
			}
			// For each PID, print parent process cmdline
			for _, ln := range lines {
				fields := strings.Fields(ln)
				if len(fields) < 2 {
					continue
				}
				ppid := fields[1]
				cmdline := fmt.Sprintf("tr '\\0' ' ' </proc/%s/cmdline 2>/dev/null | sed 's/[[:space:]]\\+$//' || true", ppid)
				if pout, err := exec.Command("sh", "-lc", cmdline).CombinedOutput(); err == nil {
					p := strings.TrimSpace(string(pout))
					if p != "" {
						fmt.Printf("  PPID %s cmd: %s\n", ppid, p)
					}
				}
			}
		}
	} else {
		fmt.Println("  Unable to query processes")
	}

	// 2) Hyprland exec-once lines
	fmt.Println("\n[Hyprland exec-once]")
	home := os.Getenv("HOME")
	hyprCfgs := []string{
		filepath.Join(home, ".config", "hypr", "hyprland.conf"),
		filepath.Join(home, ".config", "hypr", "keybindings.conf"),
		filepath.Join(home, ".config", "hypr", "autostart.conf"),
	}
	found := false
	for _, p := range hyprCfgs {
		if b, err := os.ReadFile(p); err == nil {
			lines := strings.Split(string(b), "\n")
			for _, ln := range lines {
				s := strings.TrimSpace(ln)
				if strings.HasPrefix(s, "exec-once") &&
					(strings.Contains(s, "waybar") || strings.Contains(s, "archriot --waybar-launch")) {
					if !found {
						found = true
					}
					fmt.Printf("  %s: %s\n", p, s)
				}
			}
		}
	}
	if !found {
		fmt.Println("  (none detected)")
	}

	// 3) systemd --user units referencing Waybar
	fmt.Println("\n[systemd --user]")
	if out, err := exec.Command("sh", "-lc", "systemctl --user list-units --type=service --all 2>/dev/null | grep -i waybar || true").CombinedOutput(); err == nil {
		s := strings.TrimSpace(string(out))
		if s == "" {
			fmt.Println("  (no systemd --user services matching waybar)")
		} else {
			fmt.Println(s)
			// Try to cat unit files for detail
			if out2, err2 := exec.Command("sh", "-lc", "systemctl --user list-unit-files 2>/dev/null | grep -i waybar | awk '{print $1}'").CombinedOutput(); err2 == nil {
				units := strings.Fields(strings.TrimSpace(string(out2)))
				for _, u := range units {
					if out3, err3 := exec.Command("sh", "-lc", fmt.Sprintf("systemctl --user cat %s 2>/dev/null || true", u)).CombinedOutput(); err3 == nil {
						trim := strings.TrimSpace(string(out3))
						if trim != "" {
							fmt.Printf("\n-- %s --\n%s\n", u, trim)
						}
					}
				}
			}
		}
	} else {
		fmt.Println("  (systemd --user not available)")
	}

	// 4) Autostart .desktop entries with waybar in Exec=
	fmt.Println("\n[Autostart .desktop]")
	autostartDirs := []string{
		filepath.Join(home, ".config", "autostart"),
		"/etc/xdg/autostart",
	}
	aFound := false
	for _, d := range autostartDirs {
		entries, _ := os.ReadDir(d)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".desktop") {
				continue
			}
			p := filepath.Join(d, e.Name())
			if b, err := os.ReadFile(p); err == nil {
				txt := strings.ToLower(string(b))
				// match any Exec line that includes 'waybar'
				if strings.Contains(txt, "exec=") && strings.Contains(txt, "waybar") {
					if !aFound {
						aFound = true
					}
					fmt.Printf("  %s\n", p)
				}
			}
		}
	}
	if !aFound {
		fmt.Println("  (none detected)")
	}

	fmt.Println("\nHint:")
	fmt.Println("- Keep exactly one Waybar launcher. Prefer Hyprland exec-once with `archriot --waybar-launch`.")
	fmt.Println("- Disable other launch sources (systemd --user service, autostart .desktop) if present.")
	fmt.Println("- The ArchRiot launcher now uses a non-blocking lock to avoid races on resume.")
}
