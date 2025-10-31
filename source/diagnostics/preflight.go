package diagnostics

// Preflight performs a read-only system audit similar to the legacy `--preflight`
// flag. It prints a concise report and returns an exit code (0 = success).
//
// Checks performed:
// 1) Config validation for packages.yaml
// 2) ArchRiot binary path sanity
// 3) Hyprland binds/exec-once (Telegram, Signal, Waybar launcher)
// 4) Memory optimization opt-in status
// 5) Waybar status (process count)
// 6) Portals stack presence and defaults
// 7) Wi-Fi power-save configuration and runtime state
//
// This function never exits the process. It always returns an int exit code.
// All operations are best-effort; errors are surfaced as warnings.
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"archriot-installer/config"
)

// Preflight runs the read-only diagnostics and prints results to stdout.
func Preflight() int {
	home := os.Getenv("HOME")
	fmt.Println("🔎 ArchRiot Preflight (read-only)")

	// 1) Config validation
	if err := validatePackagesConfig(); err != nil {
		fmt.Printf("⚠️  Config: %v\n", err)
	} else {
		fmt.Println("✅ Config: packages.yaml is valid")
	}

	// 2) Binary path check
	if self, err := os.Executable(); err == nil {
		expected := filepath.Join(home, ".local", "share", "archriot", "install", "archriot")
		if self == expected {
			fmt.Println("✅ Binary path:", self)
		} else {
			fmt.Printf("⚠️  Binary path mismatch: using %s; expected %s\n", self, expected)
		}
	} else {
		fmt.Println("⚠️  Binary path: could not determine")
	}

	// 3) Hyprland binds and exec-once (user config)
	hyprCfg := filepath.Join(home, ".config", "hypr", "hyprland.conf")
	if b, err := os.ReadFile(hyprCfg); err == nil {
		txt := string(b)

		// Telegram bind on $mod+G
		if strings.Contains(txt, `bind = $mod, G, exec,`) &&
			(strings.Contains(txt, "--telegram") ||
				strings.Contains(txt, `org\.telegram\.desktop`) ||
				strings.Contains(txt, "gtk-launch org.telegram.desktop") ||
				strings.Contains(txt, "telegram-desktop")) {
			fmt.Println("✅ Bind(G): Telegram mapping present")
		} else {
			fmt.Println("⚠️  Bind(G): Telegram mapping not found")
		}

		// Signal bind on $mod+S
		if strings.Contains(txt, `bind = $mod, S, exec,`) && strings.Contains(txt, "--signal") {
			fmt.Println("✅ Bind(S): Signal mapping present")
		} else {
			fmt.Println("⚠️  Bind(S): Signal mapping not found")
		}

		// Waybar exec-once launcher
		if strings.Contains(txt, "$HOME/.local/share/archriot/install/archriot --waybar-launch") &&
			strings.Contains(txt, "exec-once") {
			fmt.Println("✅ Exec-once: Waybar uses archriot --waybar-launch")
		} else {
			fmt.Println("⚠️  Exec-once: Waybar launcher not found (expected archriot --waybar-launch)")
		}
	} else {
		fmt.Printf("⚠️  Hyprland config not readable: %v\n", err)
	}

	// 4) Memory optimization (opt-in status)
	if _, err := os.Stat(filepath.Join(home, ".config", "archriot", "enable-memory-optimizations")); err == nil {
		fmt.Println("ℹ️  Memory: opt-in file present (~/.config/archriot/enable-memory-optimizations)")
	} else {
		fmt.Println("✅ Memory: no system VM tweaks (opt-in disabled)")
	}

	// 5) Waybar status (no changes, informational only)
	if out, err := exec.Command("sh", "-lc", "pgrep -x waybar | wc -l").Output(); err == nil {
		c := strings.TrimSpace(string(out))
		switch c {
		case "0":
			fmt.Println("ℹ️  Waybar: not running")
		case "1":
			fmt.Println("✅ Waybar: single instance running")
		default:
			fmt.Printf("⚠️  Waybar: %s instances detected (consider archriot --waybar-launch)\n", c)
		}
	} else {
		fmt.Println("⚠️  Waybar: unable to query status")
	}

	// 6) Portals stack checks
	fmt.Println("Portals:")
	havePortal := func(bin, libPath, proc string) bool {
		// PATH check
		if _, err := exec.LookPath(bin); err == nil {
			return true
		}
		// Well-known lib path
		if st, err := os.Stat(libPath); err == nil && !st.IsDir() {
			return true
		}
		// Running process (best-effort)
		if err := exec.Command("pgrep", "-f", proc).Run(); err == nil {
			return true
		}
		return false
	}
	printKV := func(k, v string) { fmt.Printf("%-24s %s\n", k+":", v) }

	printKV("xdg-desktop-portal", map[bool]string{true: "present", false: "missing"}[havePortal("xdg-desktop-portal", "/usr/lib/xdg-desktop-portal", "xdg-desktop-portal")])
	printKV("portal-hyprland", map[bool]string{true: "present", false: "missing"}[havePortal("xdg-desktop-portal-hyprland", "/usr/lib/xdg-desktop-portal-hyprland", "xdg-desktop-portal-hyprland")])
	printKV("portal-gtk", map[bool]string{true: "present", false: "missing"}[havePortal("xdg-desktop-portal-gtk", "/usr/lib/xdg-desktop-portal-gtk", "xdg-desktop-portal-gtk")])

	// Running portal processes (brief)
	if out, err := exec.Command("sh", "-lc", "ps -C xdg-desktop-portal -o pid,cmd --no-headers; ps -C xdg-desktop-portal-hyprland -o pid,cmd --no-headers; ps -C xdg-desktop-portal-gtk -o pid,cmd --no-headers").CombinedOutput(); err == nil {
		if s := strings.TrimSpace(string(out)); s != "" {
			fmt.Println("--- running portals ---")
			fmt.Println(s)
		}
	}

	// portals.conf default chain
	confPath := ""
	for _, p := range []string{
		filepath.Join(home, ".config", "xdg-desktop-portal", "portals.conf"),
		"/etc/xdg-desktop-portal/portals.conf",
	} {
		if _, err := os.Stat(p); err == nil {
			confPath = p
			break
		}
	}
	if confPath != "" {
		if b, err := os.ReadFile(confPath); err == nil {
			def := ""
			for _, ln := range strings.Split(string(b), "\n") {
				s := strings.TrimSpace(ln)
				if strings.HasPrefix(strings.ToLower(s), "default=") {
					def = strings.TrimSpace(strings.TrimPrefix(s, "default="))
					break
				}
			}
			if def != "" {
				printKV("portals.conf", confPath)
				printKV("default chain", def)
				l := strings.ToLower(def)
				if !strings.HasPrefix(l, "hyprland;") && l != "hyprland" {
					fmt.Println("⚠️  Tip: Put 'hyprland' first in the default chain for ScreenCast/Screenshot on Wayland.")
				}
			}
		}
	} else {
		fmt.Println("ℹ️  No portals.conf found; system defaults apply (hyprland portal should be active).")
	}

	// 7) WiFi power-save
	fmt.Println("WiFi power-save:")
	psConf := "/etc/NetworkManager/conf.d/40-wifi-powersave.conf"
	if b, err := os.ReadFile(psConf); err == nil {
		line := string(b)
		fmt.Printf("%-24s %s\n", "drop-in:", psConf)
		val := "(not set)"
		for _, ln := range strings.Split(line, "\n") {
			s := strings.TrimSpace(ln)
			if strings.HasPrefix(s, "wifi.powersave=") {
				val = strings.TrimSpace(strings.TrimPrefix(s, "wifi.powersave="))
				break
			}
		}
		fmt.Printf("%-24s %s\n", "wifi.powersave:", val)
		if val != "2" && val != "(not set)" {
			fmt.Println("⚠️  Tip: set wifi.powersave=2 to avoid Wi‑Fi power saving")
			fmt.Println("    e.g., echo -e \"[connection]\\nwifi.powersave=2\" | sudo tee /etc/NetworkManager/conf.d/40-wifi-powersave.conf && sudo systemctl reload NetworkManager")
		}
	} else {
		fmt.Printf("%-24s %s\n", "drop-in:", "missing")
		fmt.Println("⚠️  Tip: create /etc/NetworkManager/conf.d/40-wifi-powersave.conf with wifi.powersave=2")
		fmt.Println("    e.g., echo -e \"[connection]\\nwifi.powersave=2\" | sudo tee /etc/NetworkManager/conf.d/40-wifi-powersave.conf && sudo systemctl reload NetworkManager")
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
				fmt.Println("⚠️  Tip: runtime power save is ON; temporarily disable with:")
				fmt.Println("    sudo iw dev " + iface + " set power_save off")
			}
		}
	} else {
		fmt.Println("ℹ️  No wireless interface detected")
	}

	fmt.Println("✅ Preflight complete (warnings shown above if any)")
	return 0
}

// validatePackagesConfig validates packages.yaml using the config package.
func validatePackagesConfig() error {
	// Find configuration file
	configPath := config.FindConfigFile()
	if configPath == "" {
		return fmt.Errorf("packages.yaml not found")
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Validate structure and dependencies, and command safety
	if err := config.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}
	if err := config.ValidateDependencies(cfg); err != nil {
		return fmt.Errorf("dependency validation: %w", err)
	}
	if err := config.ValidateAllCommands(cfg); err != nil {
		return fmt.Errorf("command safety validation: %w", err)
	}

	return nil
}
