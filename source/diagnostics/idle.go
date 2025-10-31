package diagnostics

// IdleDiagnostics performs a read-only inspection of Hyprland idle/lock setup,
// suspend guards, and related power state indicators. It mirrors the legacy
// `--idle-diagnostics` behavior and prints a concise report. It never exits
// the process; it returns an exit code (0 = success).
//
// Checks performed:
// - Presence of hypridle and hyprlock in PATH
// - Whether hypridle is running (and its command line)
// - hypridle.conf lock_cmd and a common "10m lock" listener
// - Legacy suspend-if-undocked.sh script existence/executable bit
// - Hyprland monitors (if hyprctl is available)
// - Power supply "online" state (AC/USB-PD)
// - systemd-logind drop-ins relevant to lid/idle behavior
//
// All steps are best-effort. Any failures are shown as warnings/notes.
// No changes are made to the system.
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IdleDiagnostics runs the idle diagnostics and prints results to stdout.
func IdleDiagnostics() int {
	home := os.Getenv("HOME")
	cfg := filepath.Join(home, ".config", "hypr", "hypridle.conf")
	fmt.Println("🔎 ArchRiot Idle Diagnostics")

	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
	runOut := func(cmd string) string {
		out, err := exec.Command("sh", "-lc", cmd).CombinedOutput()
		if err != nil {
			return strings.TrimSpace(string(out))
		}
		return strings.TrimSpace(string(out))
	}
	printKV := func(k, v string) { fmt.Printf("%-24s %s\n", k+":", v) }

	// Binaries
	printKV("hypridle in PATH", map[bool]string{true: "yes", false: "no"}[have("hypridle")])
	printKV("hyprlock in PATH", map[bool]string{true: "yes", false: "no"}[have("hyprlock")])

	// Processes
	printKV("hypridle running", map[bool]string{true: "yes", false: "no"}[exec.Command("pgrep", "-x", "hypridle").Run() == nil])
	if p := runOut("pgrep -a hypridle || true"); p != "" {
		printKV("hypridle pgrep", p)
	}

	// Config
	fmt.Println("Config:", cfg)
	if b, err := os.ReadFile(cfg); err != nil {
		fmt.Println("⚠️  Cannot read hypridle.conf:", err)
	} else {
		txt := string(b)
		// lock_cmd presence
		lockCmd := "(not found)"
		for _, line := range strings.Split(txt, "\n") {
			s := strings.TrimSpace(line)
			if strings.HasPrefix(s, "lock_cmd") {
				lockCmd = s
				break
			}
		}
		printKV("lock_cmd", lockCmd)

		// Look for a 10-minute lock listener
		has10 := strings.Contains(txt, "timeout = 600") && strings.Contains(txt, "on-timeout = lock")
		printKV("10m lock listener", map[bool]string{true: "present", false: "missing"}[has10])
	}

	// Suspend guard script (legacy)
	script := filepath.Join(home, ".local", "bin", "suspend-if-undocked.sh")
	if st, err := os.Stat(script); err == nil && !st.IsDir() {
		printKV("suspend-if-undocked.sh", "present (legacy)")
		printKV("executable", map[bool]string{true: "yes", false: "no"}[st.Mode().Perm()&0o111 != 0])
	} else {
		printKV("suspend-if-undocked.sh", "not needed (using CLI)")
	}
	// Preferred: native CLI guard
	printKV("suspend guard (CLI)", "available: archriot --suspend-if-undocked")

	// Dock/AC state (informational)
	if have("hyprctl") {
		ms := runOut("hyprctl monitors | sed -n '1,60p'")
		if ms != "" {
			fmt.Println("--- hyprctl monitors ---")
			fmt.Println(ms)
		}
	}

	// AC / USB-PD online
	ac := runOut(`for f in /sys/class/power_supply/*/online; do [ -f "$f" ] && n="$(basename "$(dirname "$f")")"; v="$(cat "$f" 2>/dev/null || true)"; [ -n "$v" ] && echo "$n: $v"; done`)
	if ac != "" {
		fmt.Println("--- power_supply online ---")
		fmt.Println(ac)
	}

	// Logind drop-ins
	ten := "/etc/systemd/logind.conf.d/10-docked-ignore-lid.conf"
	twenty := "/etc/systemd/logind.conf.d/20-idle-ignore.conf"
	fmt.Println("--- logind drop-ins ---")
	if out := runOut("sudo sed -n '1,50p' " + ten + " 2>/dev/null || true"); out != "" {
		fmt.Println(out)
	}
	if out := runOut("sudo sed -n '1,50p' " + twenty + " 2>/dev/null || true"); out != "" {
		fmt.Println(out)
	}

	fmt.Println("✅ Idle diagnostics complete")
	return 0
}
