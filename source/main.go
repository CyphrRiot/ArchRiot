package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/executor"
	"archriot-installer/git"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/orchestrator"
	"archriot-installer/theming"
	"archriot-installer/tools"
	"archriot-installer/tui"
	"archriot-installer/upgrade"
	"archriot-installer/version"
)

// Global program reference for TUI
var program *tea.Program

// Global git credentials handling
var (
	gitInputDone chan bool
)

// Global Secure Boot setup handling
var (
	secureBootSetupDone        chan bool
	secureBootContinuationDone chan bool
)

// Global model instance
var model *tui.InstallModel

// Global reboot flag
var shouldReboot bool

// Strict ABI guard flag
var strictABI bool

// confirmInstallation shows initial installation confirmation using TUI
func confirmInstallation() bool {
	// Set up TUI helper functions before creating model
	tui.SetVersionGetter(func() string { return version.Get() })

	// Create confirmation model
	model := tui.NewInstallModel()
	model.SetConfirmationMode("install", fmt.Sprintf("λ Install ArchRiot v%s?", version.Get()))

	program := tea.NewProgram(model)

	// Run the confirmation TUI
	finalModel, err := program.Run()
	if err != nil {
		log.Fatalf("❌ Confirmation failed: %v", err)
	}

	// Extract the result
	if m, ok := finalModel.(*tui.InstallModel); ok {
		return m.GetConfirmationResult()
	}

	return false
}

// setupSudo ensures passwordless sudo is configured
func setupSudo() error {
	log.Printf("🔐 Checking sudo configuration...")

	// Test if passwordless sudo already works
	if testSudo() {
		log.Printf("✅ Passwordless sudo is already working")
		return nil
	}

	log.Printf("❌ Passwordless sudo is required for ArchRiot installation")
	log.Printf("")
	log.Printf("🔧 Please configure passwordless sudo by running these commands:")
	log.Printf("")
	log.Printf("   1. Add your user to the wheel group:")
	log.Printf("      sudo usermod -aG wheel $USER")
	log.Printf("")
	log.Printf("   2. Enable passwordless sudo for wheel group:")
	log.Printf("      echo '%%wheel ALL=(ALL) NOPASSWD: ALL' | sudo tee -a /etc/sudoers")
	log.Printf("")
	log.Printf("   3. Verify the rule was added:")
	log.Printf("      sudo grep 'wheel.*NOPASSWD' /etc/sudoers")
	log.Printf("")
	log.Printf("   4. Log out and log back in, then run the installer again")
	log.Printf("")
	log.Printf("💡 This is required because the installer needs to install packages")
	log.Printf("   and configure system services without password prompts.")

	return fmt.Errorf("passwordless sudo not configured")
}

// testSudo tests if passwordless sudo is working
// testSudo checks if passwordless sudo is properly configured
func testSudo() bool {
	// Clear sudo timestamp cache to avoid false positives from cached credentials
	// Ignore error as this might fail if user never used sudo
	_ = exec.Command("sudo", "-k").Run()

	// Now test for actual passwordless sudo configuration
	cmd := exec.Command("sudo", "-n", "true")
	return cmd.Run() == nil
}

// installYay installs yay AUR helper with retry logic for AUR unreliability
func installYay() error {
	log.Printf("🔍 Checking for yay AUR helper...")

	// Check if yay is already installed
	if _, err := exec.LookPath("yay"); err == nil {
		log.Printf("✅ yay is already installed")
		return nil
	}

	log.Printf("📦 Installing yay AUR helper...")

	// Install prerequisites first (this is reliable)
	log.Printf("🔧 Installing prerequisites...")
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", "base-devel", "git")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install yay prerequisites: %v", err)
	}

	// Retry yay installation up to 3 times with user prompts
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 yay installation attempt %d/%d...", attempt, maxRetries)

		if err := attemptYayInstall(); err != nil {
			log.Printf("❌ Attempt %d failed: %v", attempt, err)

			if attempt < maxRetries {
				// Offer retry with different methods
				log.Printf("")
				log.Printf("🌐 AUR servers may be experiencing issues.")
				log.Printf("💡 This is common and usually resolves quickly.")
				log.Printf("")
				log.Printf("Options:")
				log.Printf("  1. Retry yay installation")
				log.Printf("  2. Continue without AUR packages (limited functionality)")
				log.Printf("  3. Exit and try again later")
				log.Printf("")
				log.Printf("Choice [1/2/3]: ")

				var choice string
				fmt.Scanln(&choice)

				switch choice {
				case "1", "":
					log.Printf("🔄 Retrying yay installation...")
					continue
				case "2":
					log.Printf("⚠️ Continuing without AUR support - some packages may be missing")
					return nil
				case "3":
					log.Printf("👋 Exiting installer - try again when AUR is stable")
					os.Exit(0)
				default:
					log.Printf("🔄 Invalid choice, retrying...")
					continue
				}
			} else {
				// Final attempt failed
				log.Printf("")
				log.Printf("❌ All yay installation attempts failed.")
				log.Printf("🌐 AUR servers appear to be down or unreachable.")
				log.Printf("")
				log.Printf("Options:")
				log.Printf("  1. Continue without AUR packages (limited functionality)")
				log.Printf("  2. Exit and try again later")
				log.Printf("")
				log.Printf("Choice [1/2]: ")

				var choice string
				fmt.Scanln(&choice)

				switch choice {
				case "1", "":
					log.Printf("⚠️ Continuing without AUR support - some packages may be missing")
					return nil
				default:
					log.Printf("👋 Exiting installer - try again when AUR is stable")
					os.Exit(0)
				}
			}
		} else {
			log.Printf("✅ yay installed successfully on attempt %d", attempt)
			return nil
		}
	}

	return nil
}

// attemptYayInstall performs a single yay installation attempt
func attemptYayInstall() error {
	tempDir := "/tmp/yay-bin-install"

	// Clean up any previous attempts
	if err := exec.Command("rm", "-rf", tempDir).Run(); err != nil {
		log.Printf("⚠️ Could not clean temp directory: %v", err)
	}

	// Clone yay-bin repository (precompiled binary version)
	log.Printf("📥 Downloading yay-bin from AUR...")
	cmd := exec.Command("git", "clone", "https://aur.archlinux.org/yay-bin.git", tempDir)
	if err := cmd.Run(); err != nil {
		exec.Command("rm", "-rf", tempDir).Run() // Clean up on failure
		return fmt.Errorf("failed to clone yay-bin repository: %v", err)
	}

	// Build and install yay-bin (this downloads the precompiled binary)
	log.Printf("🔨 Installing yay-bin...")
	cmd = exec.Command("makepkg", "-si", "--noconfirm")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		exec.Command("rm", "-rf", tempDir).Run() // Clean up on failure
		return fmt.Errorf("failed to build yay-bin: %v", err)
	}

	// Clean up temp directory
	exec.Command("rm", "-rf", tempDir).Run()

	// Verify yay installation
	if _, err := exec.LookPath("yay"); err != nil {
		return fmt.Errorf("yay installation failed - not found in PATH")
	}

	return nil
}

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--ascii-only":
			// Force ASCII mode for testing
			logger.SetForceAsciiMode(true)
			// Continue with normal installation (fall through to main logic)

		case "--strict-abi":
			// Enable strict ABI guard: block install when compositor/Wayland upgrades are pending
			strictABI = true
			// Do not return; allow normal installation or other flags to proceed

		case "--tools", "-t":
			// Initialize basic logging for tools
			if err := logger.InitLogging(); err != nil {
				log.Fatalf("❌ Failed to initialize logging: %v", err)
			}
			defer logger.CloseLogging()

			// Run tools interface
			if err := tools.RunToolsInterface(); err != nil {
				log.Fatalf("❌ Tools interface failed: %v", err)
			}
			return

		case "--validate":
			// Validate packages.yaml configuration
			if err := validateConfig(); err != nil {
				log.Fatalf("❌ Configuration validation failed: %v", err)
			}
			fmt.Println("✅ Configuration is valid")
			return

		case "--version", "-v":
			if err := version.ReadVersion(); err != nil {
				log.Fatalf("❌ Failed to read version: %v", err)
			}
			fmt.Printf("ArchRiot version %s\n", version.Get())
			return

		case "--secure_boot_stage":
			// Run Secure Boot continuation after reboot
			if err := runSecureBootContinuation(); err != nil {
				log.Fatalf("❌ Secure Boot continuation failed: %v", err)
			}
			return

		case "--apply-wallpaper-theme":
			// Apply wallpaper-based theming
			if len(os.Args) < 3 {
				log.Fatalf("❌ --apply-wallpaper-theme requires wallpaper path argument")
			}
			wallpaperPath := os.Args[2]
			if err := theming.ApplyWallpaperTheme(wallpaperPath); err != nil {
				log.Fatalf("❌ Failed to apply wallpaper theme: %v", err)
			}
			return

		case "--toggle-dynamic-theming":
			// Toggle dynamic theming on/off
			if len(os.Args) < 3 {
				log.Fatalf("❌ --toggle-dynamic-theming requires true/false argument")
			}
			enabled := os.Args[2] == "true"
			if err := theming.ToggleDynamicTheming(enabled); err != nil {
				log.Fatalf("❌ Failed to toggle dynamic theming: %v", err)
			}
			return

		case "--help", "-h":
			showHelp()
			return

		case "--preflight":
			// Read-only system audit: validate config, paths, binds, exec-once, memory opt-in, and Waybar status
			home := os.Getenv("HOME")
			fmt.Println("🔎 ArchRiot Preflight (read-only)")

			// 1) Config validation
			if err := validateConfig(); err != nil {
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

				if strings.Contains(txt, `bind = $mod, G, exec,`) &&
					(strings.Contains(txt, "--telegram") ||
						strings.Contains(txt, `org\.telegram\.desktop`) ||
						strings.Contains(txt, "gtk-launch org.telegram.desktop") ||
						strings.Contains(txt, "telegram-desktop")) {
					fmt.Println("✅ Bind(G): Telegram mapping present")
				} else {
					fmt.Println("⚠️  Bind(G): Telegram mapping not found")
				}

				if strings.Contains(txt, `bind = $mod, S, exec,`) && strings.Contains(txt, "--signal") {
					fmt.Println("✅ Bind(S): Signal mapping present")
				} else {
					fmt.Println("⚠️  Bind(S): Signal mapping not found")
				}

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

			// Portals stack checks
			fmt.Println("Portals:")
			havePortal := func(bin, libPath, proc string) bool {
				// 1) PATH check
				if _, err := exec.LookPath(bin); err == nil {
					return true
				}
				// 2) Well-known lib path (most distros place portals under /usr/lib)
				if st, err := os.Stat(libPath); err == nil && !st.IsDir() {
					return true
				}
				// 3) Running process check (best-effort)
				if err := exec.Command("pgrep", "-f", proc).Run(); err == nil {
					return true
				}
				return false
			}
			printKV := func(k, v string) { fmt.Printf("%-24s %s\n", k+":", v) }

			// Detect presence via PATH, lib path, or running process
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
							def = strings.TrimPrefix(s, "default=")
							break
						}
					}
					if def != "" {
						printKV("portals.conf", confPath)
						printKV("default chain", def)
						if !strings.HasPrefix(def, "hyprland;") && strings.ToLower(def) != "hyprland" {
							fmt.Println("⚠️  Tip: Put 'hyprland' first in the default chain for ScreenCast/Screenshot on Wayland.")
						}
					}
				}
			} else {
				fmt.Println("ℹ️  No portals.conf found; system defaults apply (hyprland portal should be active).")
			}

			fmt.Println("WiFi power-save:")

			// Check NetworkManager drop-in for wifi.powersave
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
				}
			} else {
				fmt.Printf("%-24s %s\n", "drop-in:", "missing")
				fmt.Println("⚠️  Tip: create /etc/NetworkManager/conf.d/40-wifi-powersave.conf with wifi.powersave=2")
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
						fmt.Println("⚠️  Tip: runtime power save is ON; consider turning it off to avoid drops while locked")
					}
				}
			} else {
				fmt.Println("ℹ️  No wireless interface detected")
			}

			fmt.Println("✅ Preflight complete (warnings shown above if any)")
			return

		case "--idle-diagnostics":
			// Inspect hypridle/hyprlock presence, running status, config, and suspend guard wiring
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

				// look for a 10-minute lock listener
				has10 := strings.Contains(txt, "timeout = 600") && strings.Contains(txt, "on-timeout = lock")
				printKV("10m lock listener", map[bool]string{true: "present", false: "missing"}[has10])
			}

			// Suspend guard script
			script := filepath.Join(home, ".local", "bin", "suspend-if-undocked.sh")
			if st, err := os.Stat(script); err == nil && !st.IsDir() {
				printKV("suspend-if-undocked.sh", "present")
				printKV("executable", map[bool]string{true: "yes", false: "no"}[st.Mode().Perm()&0o111 != 0])
			} else {
				printKV("suspend-if-undocked.sh", "missing")
			}

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
			return

		case "--install":
			// Explicit install mode: continue to normal installation flow (no return)
		case "--upgrade":
			// Run upgrade flow in a standalone TUI and exit
			if err := version.ReadVersion(); err != nil {
				log.Fatalf("❌ Failed to read version: %v", err)
			}
			if err := logger.InitLogging(); err != nil {
				log.Fatalf("❌ Failed to initialize logging: %v", err)
			}
			defer logger.CloseLogging()
			tui.SetVersionGetter(func() string { return version.Get() })
			tui.SetLogPathGetter(func() string { return logger.GetLogPath() })
			upModel := tui.NewInstallModel()
			upProgram := tea.NewProgram(upModel)
			logger.SetProgram(upProgram)
			upgrade.SetProgram(upProgram)

			// 🛡️ Mullvad VPN connection guard (auto-reconnect during upgrade if it drops)
			var stopMullvadGuard chan struct{}
			if _, err := exec.LookPath("mullvad"); err == nil {
				connectedAtStart := false
				if out, err := exec.Command("mullvad", "status").CombinedOutput(); err == nil {
					s := strings.ToLower(string(out))
					if strings.Contains(s, "connected") {
						connectedAtStart = true
					}
				}
				// Only guard if user was connected when starting the upgrade
				if connectedAtStart {
					stopMullvadGuard = make(chan struct{})
					// Informational log in the TUI
					upProgram.Send(tui.LogMsg("🛡️ Mullvad: connection guard active during upgrade"))

					go func(stop <-chan struct{}) {
						ticker := time.NewTicker(5 * time.Second)
						defer ticker.Stop()
						lastAttempt := time.Now().Add(-20 * time.Second)
						for {
							select {
							case <-stop:
								return
							case <-ticker.C:
								out, err := exec.Command("mullvad", "status").CombinedOutput()
								if err != nil {
									continue
								}
								ss := strings.ToLower(string(out))
								// If still connected, nothing to do
								if strings.Contains(ss, "connected") {
									continue
								}
								// Only attempt reconnect if an account is present
								if acc, err := exec.Command("mullvad", "account", "get").CombinedOutput(); err == nil && strings.Contains(strings.ToLower(string(acc)), "mullvad account") {
									// Throttle reconnect attempts
									if time.Since(lastAttempt) > 10*time.Second {
										_ = exec.Command("mullvad", "connect").Start()
										lastAttempt = time.Now()
									}
								}
							}
						}
					}(stopMullvadGuard)
				}
			}

			// Start upgrade prompt in background
			go func() {
				time.Sleep(100 * time.Millisecond)
				if err := upgrade.PromptAndRun(); err != nil {
					upProgram.Send(tui.LogMsg("❌ Upgrade failed: " + err.Error()))
					upProgram.Send(tui.FailureMsg{Error: "Upgrade failed"})
					// Stop Mullvad guard on failure
					if stopMullvadGuard != nil {
						close(stopMullvadGuard)
						stopMullvadGuard = nil
					}
					return
				}
				// Success banner and finalization
				upProgram.Send(tui.LogMsg("🎉 Upgrade completed!"))
				logger.LogMessage("SUCCESS", "Upgrade completed")

				// Refresh idle manager so 10m lock works immediately post-upgrade
				upProgram.Send(tui.LogMsg("🔄 Refreshing idle manager (hypridle)…"))
				_ = exec.Command("pkill", "hypridle").Run()
				if _, err := exec.LookPath("hypridle"); err == nil {
					cmd := exec.Command("hypridle")
					cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
					_ = cmd.Start()
				}

				installer.FinalizePackageManagers()
				upProgram.Send(tui.DoneMsg{})
			}()
			// Run TUI and exit
			if _, err := upProgram.Run(); err != nil {
				log.Fatalf("TUI error: %v", err)
			}
			// Stop Mullvad guard after TUI exits
			if stopMullvadGuard != nil {
				close(stopMullvadGuard)
				stopMullvadGuard = nil
			}
			return
		case "--signal":
			// Focus-or-launch Signal Desktop with Wayland flags and async focus
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
			// If a Signal window exists, focus it
			if hyprClientsContains("class: Signal") || hyprClientsContains("class: signal") {
				if focusSignal() {
					return
				}
			}
			// Otherwise launch Signal (Wayland/Ozone) and focus when ready
			notify("Signal", "Launching Signal Desktop…", 3000)
			_ = exec.Command("env", "GDK_SCALE=1", "signal-desktop", "--ozone-platform=wayland", "--enable-features=UseOzonePlatform").Start()
			go func() {
				for i := 0; i < 20; i++ {
					time.Sleep(250 * time.Millisecond)
					if focusSignal() {
						return
					}
				}
			}()
			return

		case "--telegram":
			// Focus-or-launch Telegram Desktop with async focus; prefer native first; add runtime logging
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

			// Focus without scanning clients to avoid brittle parsing
			focusTelegram := func() bool {
				// Broad match set: org.telegram.desktop (Flatpak/native), telegram-desktop (native),
				// TelegramDesktop (legacy), telegramdesktop/Telegram (rare)
				if exec.Command("hyprctl", "dispatch", "focuswindow", "class:^(org\\.telegram\\.desktop|org\\.telegram\\.desktop\\.TelegramDesktop|org\\.telegram\\..*desktop.*|telegram-desktop|TelegramDesktop|telegramdesktop|Telegram)$").Run() == nil {
					_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
					return true
				}
				return false
			}

			// Focus only if a Telegram window is present; otherwise proceed to launch
			if exec.Command("sh", "-lc", "hyprctl clients 2>/dev/null | grep -qE 'class:\\s*(org\\.telegram\\.desktop(\\.TelegramDesktop)?|org\\.telegram\\..*desktop.*|telegram-desktop|TelegramDesktop|telegramdesktop|Telegram)\\b'").Run() == nil {
				logAppend("window present; focusing")
				if focusTelegram() {
					return
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
					return
				}
				time.Sleep(1500 * time.Millisecond)
				_ = exec.Command("telegram-desktop").Start()
				if waitFocus(60) {
					return
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
						return
					}
					// Delayed retry if the window didn't realize on the first attempt
					time.Sleep(1500 * time.Millisecond)
					_ = exec.Command("gtk-launch", id).Start()
					if waitFocus(60) {
						return
					}
				}
				// Try to discover a Telegram desktop file dynamically
				out, _ := exec.Command("sh", "-lc", "for d in ~/.local/share/applications /usr/local/share/applications /usr/share/applications; do for f in \"$d\"/*[Tt]elegram*.desktop; do [ -f \"$f\" ] && basename \"${f%.desktop}\"; done; done | head -n 1").CombinedOutput()
				dyn := strings.TrimSpace(string(out))
				if dyn != "" {
					logAppend("gtk-launch discovered desktop id: " + dyn)
					_ = exec.Command("gtk-launch", dyn).Start()
					if waitFocus(60) {
						return
					}
					// One more attempt after a short delay for slower/cold systems
					time.Sleep(1500 * time.Millisecond)
					_ = exec.Command("gtk-launch", dyn).Start()
					if waitFocus(60) {
						return
					}
				}
			}

			// 3) Flatpak (attempt without info check; slower systems may need a retry)
			if have("flatpak") {
				logAppend("launching Flatpak org.telegram.desktop (no info probe)")
				notify("Telegram", "Launching Telegram Desktop (Flatpak)…", 3000)
				_ = exec.Command("flatpak", "run", "org.telegram.desktop").Start()
				if waitFocus(60) {
					return
				}
				time.Sleep(1500 * time.Millisecond)
				_ = exec.Command("flatpak", "run", "org.telegram.desktop").Start()
				if waitFocus(60) {
					return
				}
				logAppend("flatpak launch did not realize a window in time")
			}

			// 4) Parse .desktop Exec fallback
			// Last-resort: parse Exec lines from all Telegram*.desktop entries and try them in order
			if out, _ := exec.Command("sh", "-lc", "for d in ~/.local/share/applications /usr/local/share/applications /usr/share/applications; do for f in \"$d\"/*[Tt]elegram*.desktop; do [ -f \"$f\" ] && grep -m1 '^Exec=' \"$f\" | sed -E 's/^Exec=//; s/%[fFuUdDnNickvm]//g'; done; done").CombinedOutput(); len(out) > 0 {
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
						return
					}
					// Slow systems: retry the same command once after a short delay
					time.Sleep(1500 * time.Millisecond)
					_ = exec.Command(fields[0], fields[1:]...).Start()
					if waitFocus(90) {
						return
					}
				}
			}

			logAppend("all launch paths failed")
			notify("Telegram", "Unable to launch Telegram", 2000)
			os.Exit(1)

		case "--waybar-status":
			// Print Waybar status: running or stopped
			if exec.Command("pgrep", "-x", "waybar").Run() == nil {
				fmt.Println("running")
			} else {
				fmt.Println("stopped")
			}
			return

		case "--waybar-launch":
			// Single-instance Waybar launcher with non-blocking lock and logging
			home := os.Getenv("HOME")
			logDir := filepath.Join(home, ".cache", "archriot")
			_ = os.MkdirAll(logDir, 0o755)
			lockPath := filepath.Join(logDir, "waybar-launch.lock")
			logPath := filepath.Join(logDir, "runtime.log")
			// Cap size at ~5MB: if exceeded, truncate to start a fresh session
			if fi, err := os.Stat(logPath); err == nil && fi.Size() > 5*1024*1024 {
				_ = os.Truncate(logPath, 0)
			}

			// Dedupe if multiple Waybar PIDs exist, then check running status
			if out, err := exec.Command("sh", "-lc", "pgrep -x waybar | wc -l").Output(); err == nil {
				c := strings.TrimSpace(string(out))
				if c != "" && c != "0" && c != "1" {
					_ = exec.Command("pkill", "-x", "waybar").Run()
				}
			}
			// If Waybar already running (single), exit quietly
			if exec.Command("pgrep", "-x", "waybar").Run() == nil {
				return
			}

			// Acquire non-blocking lock
			lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
					_ = lockFile.Close()
					// Another launcher holds the lock; exit quietly
					return
				}
				defer func() {
					// Best-effort unlock on exit (closing fd releases lock)
					_ = lockFile.Close()
				}()
			}

			// Re-check after lock to avoid races
			if exec.Command("pgrep", "-x", "waybar").Run() == nil {
				return
			}

			// Open log file and start Waybar
			logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				// Fallback: start without logging
				_ = exec.Command("waybar").Start()
				return
			}
			defer logFile.Close()

			// Session header to delineate launches
			_, _ = fmt.Fprintf(logFile, "\n==== Waybar session start %s ====\n", time.Now().Format(time.RFC3339))

			cmd := exec.Command("waybar")
			cmd.Stdout = logFile
			cmd.Stderr = logFile

			// Start Waybar in its own session (detached); do not wait, so closing the
			// launching session/terminal won't kill Waybar.
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			cmd.Stdin = nil
			_ = cmd.Start()
			return

		case "--trezor":
			// Focus-or-launch Trezor Suite (native > Flatpak > AppImage), then async focus when ready
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
				return
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
			return

		case "--wallet":
			// Focus-or-launch crypto wallet: prefer existing window; otherwise open whichever is installed
			// Targets:
			// - Trezor Suite (native/Flatpak/AppImage)
			// - Ledger Live (native/Flatpak/AppImage)
			home := os.Getenv("HOME")
			trezorAppImage := filepath.Join(home, ".local", "bin", "trezor-suite.AppImage")
			ledgerAppImage := filepath.Join(home, ".local", "bin", "ledger-live.AppImage")

			// Wallet logging (Task 8 acceptance: logs helpful)
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
				return
			}
			if (hyprClientsContains("class: Ledger Live") || hyprClientsContains("class: ledger live")) && focusClass("Ledger Live|ledger live") {
				logAppend("Focusing Ledger Live")
				notify("Wallet", "Focusing Ledger Live…")
				return
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
			return

		case "--pomodoro-click":
			// File-based click handler: double-click = reset, single-click = toggle (after short delay)
			clickFile := "/tmp/waybar-tomato-click"
			stateFile := "/tmp/waybar-tomato-timer.state"

			// If a click marker exists, treat as double-click (reset)
			if _, err := os.Stat(clickFile); err == nil {
				_ = os.Remove(clickFile)
				ts := time.Now().Unix()
				_ = os.WriteFile(stateFile, []byte(fmt.Sprintf("{\n    \"action\": \"reset\",\n    \"timestamp\": %d\n}\n", ts)), 0644)
				return
			}

			// First click: create marker and spawn delayed toggle worker
			_ = os.WriteFile(clickFile, []byte(fmt.Sprintf("%d", time.Now().UnixNano())), 0644)

			// Re-exec self with a delayed toggle worker (non-blocking)
			if self, err := os.Executable(); err == nil {
				_ = exec.Command(self, "--pomodoro-delay-toggle").Start()
			}
			return

		case "--pomodoro-delay-toggle":
			// After a brief delay, if click marker still exists, it's a single-click -> toggle
			clickFile := "/tmp/waybar-tomato-click"
			stateFile := "/tmp/waybar-tomato-timer.state"

			time.Sleep(500 * time.Millisecond)
			if _, err := os.Stat(clickFile); err == nil {
				_ = os.Remove(clickFile)
				ts := time.Now().Unix()
				_ = os.WriteFile(stateFile, []byte(fmt.Sprintf("{\n    \"action\": \"toggle\",\n    \"timestamp\": %d\n}\n", ts)), 0644)
			}
			return

		case "--swaybg-next":
			home := os.Getenv("HOME")
			bgsDir := filepath.Join(home, ".local", "share", "archriot", "backgrounds")
			stateFile := filepath.Join(home, ".config", "archriot", ".current-background")

			entries, err := os.ReadDir(bgsDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Backgrounds directory not found: %s\n", bgsDir)
				os.Exit(1)
			}

			var files []string
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				name := e.Name()
				lower := strings.ToLower(name)
				if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
					strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") {
					files = append(files, filepath.Join(bgsDir, name))
				}
			}
			if len(files) == 0 {
				fmt.Fprintf(os.Stderr, "No background images found in %s\n", bgsDir)
				os.Exit(1)
			}
			sort.Strings(files)

			current := ""
			if b, err := os.ReadFile(stateFile); err == nil {
				current = strings.TrimSpace(string(b))
			}
			idx := -1
			for i, f := range files {
				if f == current {
					idx = i
					break
				}
			}
			next := files[(idx+1)%len(files)]

			_ = os.MkdirAll(filepath.Dir(stateFile), 0o755)
			_ = os.WriteFile(stateFile, []byte(next+"\n"), 0o644)

			_ = exec.Command("pkill", "-x", "swaybg").Run()
			time.Sleep(500 * time.Millisecond)

			cmd := exec.Command("swaybg", "-i", next, "-m", "fill")
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			cmd.Stdin = nil
			cmd.Stdout = nil
			cmd.Stderr = nil
			_ = cmd.Start()

			// Best-effort theme refresh
			if self, err := os.Executable(); err == nil {
				_ = exec.Command(self, "--apply-wallpaper-theme", next).Start()
			}

			time.Sleep(1 * time.Second)
			if exec.Command("pgrep", "-x", "swaybg").Run() == nil {
				fmt.Printf("🖼️  Switched to background: %s\n", filepath.Base(next))
				fmt.Println("✓ Background service restarted")
			} else {
				fmt.Println("⚠ Background service may not have started properly")
			}
			return

		case "--waybar-workspace-click":
			// Usage: archriot --waybar-workspace-click <workspaceNumber>
			if len(os.Args) < 3 {
				return
			}
			ws := strings.TrimSpace(os.Args[2])

			// Ensure hyprctl is available
			if _, err := exec.LookPath("hyprctl"); err != nil {
				return
			}
			// Validate numeric workspace name
			if _, err := strconv.Atoi(ws); err != nil {
				return
			}
			// Switch workspace
			_ = exec.Command("hyprctl", "dispatch", "workspace", ws).Run()
			return

		case "--startup-background":
			// Start wallpaper at login from saved preferences, without applying theme (to avoid startup races)
			home := os.Getenv("HOME")
			configFile := filepath.Join(home, ".config", "archriot", "background-prefs.json")
			bgsDir := filepath.Join(home, ".local", "share", "archriot", "backgrounds")
			stateFile := filepath.Join(home, ".config", "archriot", ".current-background")
			defaultName := "riot_01.jpg"

			// Read desired background name from JSON (current_background)
			bgName := ""
			if b, err := os.ReadFile(configFile); err == nil {
				var obj map[string]interface{}
				if json.Unmarshal(b, &obj) == nil {
					if v, ok := obj["current_background"]; ok {
						if s, ok2 := v.(string); ok2 {
							bgName = s
						}
					}
				}
			}
			if strings.TrimSpace(bgName) == "" || strings.EqualFold(bgName, "null") {
				bgName = defaultName
			}

			// Resolve file path with fallbacks
			pick := filepath.Join(bgsDir, bgName)
			if st, err := os.Stat(pick); err != nil || st.IsDir() {
				pick = filepath.Join(bgsDir, defaultName)
			}
			if st, err := os.Stat(pick); err != nil || st.IsDir() {
				entries, _ := os.ReadDir(bgsDir)
				for _, e := range entries {
					if e.IsDir() {
						continue
					}
					name := e.Name()
					lower := strings.ToLower(name)
					if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
						strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") {
						pick = filepath.Join(bgsDir, name)
						break
					}
				}
			}

			if pick == "" {
				// No backgrounds found; nothing to do
				return
			}

			// Update state file for runtime cycling compatibility
			_ = os.MkdirAll(filepath.Dir(stateFile), 0o755)
			_ = os.WriteFile(stateFile, []byte(pick+"\n"), 0o644)

			// Relaunch swaybg (detached); do not apply theme on startup to avoid conflicts
			_ = exec.Command("pkill", "-x", "swaybg").Run()
			time.Sleep(300 * time.Millisecond)
			cmd := exec.Command("swaybg", "-i", pick, "-m", "fill")
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			cmd.Stdin = nil
			cmd.Stdout = nil
			cmd.Stderr = nil
			_ = cmd.Start()
			time.Sleep(500 * time.Millisecond)
			return

		case "--stabilize-session":
			// Dedupe Waybar and restart hypridle; optional inhibitor via --inhibit
			_ = exec.Command("pkill", "-x", "waybar").Run()
			if self, err := os.Executable(); err == nil {
				_ = exec.Command(self, "--waybar-launch").Start()
			} else {
				_ = exec.Command("waybar").Start()
			}
			_ = exec.Command("pkill", "hypridle").Run()
			if _, err := exec.LookPath("hypridle"); err == nil {
				cmd := exec.Command("hypridle")
				cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
				cmd.Stdin = nil
				cmd.Stdout = nil
				cmd.Stderr = nil
				_ = cmd.Start()
			}
			// Optional inhibitor if provided
			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] == "--inhibit" {
					if self, err := os.Executable(); err == nil {
						_ = exec.Command(self, "--stay-awake").Start()
					} else {
						cmd := exec.Command("systemd-inhibit", "--what=sleep", "--why=ArchRiot Stay Awake", "bash", "-lc", "while :; do sleep 300; done")
						cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
						cmd.Stdin = nil
						cmd.Stdout = nil
						cmd.Stderr = nil
						_ = cmd.Start()
					}
					break
				}
			}
			return

		case "--zed":
			// Focus-or-launch Zed (native > Flatpak) with Wayland-friendly env; async focus
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
				return
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
				if _, err := exec.LookPath("lspci"); err == nil {
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

			// Launch: prefer native zed/zeditor (with GL on Intel), then Flatpak
			launched := false
			if have("zed") {
				if intel {
					logAppend("launch: env WGPU_BACKEND=gl zed")
					_ = exec.Command("env", "WGPU_BACKEND=gl", "zed").Start()
				} else {
					logAppend("launch: zed")
					_ = exec.Command("zed").Start()
				}
				launched = true
			} else if have("zeditor") {
				// Build env for child only (no global environment changes)
				args := []string{}
				if intel {
					args = append(args, "WGPU_BACKEND=gl")
					logAppend("launch: env WGPU_BACKEND=gl zeditor (Wayland env applied)")
				} else {
					logAppend("launch: zeditor (Wayland env applied)")
				}
				args = append(args,
					"WAYLAND_DISPLAY="+os.Getenv("WAYLAND_DISPLAY"),
					"GDK_BACKEND=wayland",
					"QT_QPA_PLATFORM=wayland",
					"SDL_VIDEODRIVER=wayland",
					"zeditor",
				)
				_ = exec.Command("env", args...).Start()
				launched = true
			} else if have("flatpak") && exec.Command("flatpak", "info", "--show-commit", "dev.zed.Zed").Run() == nil {
				// Flatpak env overrides for WGPU won't apply; still log for traceability
				if intel {
					logAppend("launch: flatpak run dev.zed.Zed (Intel detected; GL override not applied to Flatpak)")
				} else {
					logAppend("launch: flatpak run dev.zed.Zed")
				}
				_ = exec.Command("flatpak", "run", "dev.zed.Zed").Start()
				launched = true
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
			return

		case "--welcome":
			// Launch the existing welcome app non-blocking (if present)
			// Path: $HOME/.local/share/archriot/config/bin/welcome
			{
				home := os.Getenv("HOME")
				welcome := filepath.Join(home, ".local", "share", "archriot", "config", "bin", "welcome")
				if _, err := os.Stat(welcome); err == nil {
					cmd := exec.Command(welcome)
					cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
					cmd.Stdin = nil
					cmd.Stdout = nil
					cmd.Stderr = nil
					_ = cmd.Start()
				}
			}
			return

		case "--waybar-pomodoro":
			// Waybar Pomodoro JSON emitter (native replacement for waybar-tomato-timer.py)
			type pomoState struct {
				Mode            string   `json:"mode"` // "idle", "work", "break", "break_complete"
				Running         bool     `json:"running"`
				EndTimeISO      string   `json:"end_time"`
				PausedRemaining *float64 `json:"paused_remaining"`
			}
			type pomoOut struct {
				Text       string `json:"text"`
				Tooltip    string `json:"tooltip"`
				Class      string `json:"class"`
				Percentage int    `json:"percentage"`
			}

			now := time.Now()
			home := os.Getenv("HOME")
			configPath := filepath.Join(home, ".config", "archriot", "archriot.conf")
			timerFile := "/tmp/waybar-tomato.json"
			stateFile := "/tmp/waybar-tomato-timer.state"

			enabled := true
			workMinutes := 25
			breakMinutes := 5

			// Minimal INI parse for [pomodoro]
			if b, err := os.ReadFile(configPath); err == nil {
				sc := bufio.NewScanner(strings.NewReader(string(b)))
				inSection := false
				for sc.Scan() {
					line := strings.TrimSpace(sc.Text())
					if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
						continue
					}
					if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
						sect := strings.ToLower(strings.Trim(line, "[]"))
						inSection = sect == "pomodoro"
						continue
					}
					if !inSection {
						continue
					}
					if strings.HasPrefix(strings.ToLower(line), "enabled") && strings.Contains(line, "=") {
						val := strings.ToLower(strings.TrimSpace(strings.SplitN(line, "=", 2)[1]))
						if val == "false" || val == "0" || val == "no" || val == "off" {
							enabled = false
						}
					}
					if strings.HasPrefix(strings.ToLower(line), "duration") && strings.Contains(line, "=") {
						val := strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
						if n, err := strconv.Atoi(val); err == nil && n > 0 && n <= 120 {
							workMinutes = n
						}
					}
				}
			}

			// Helper: JSON marshal or fallback
			mustJSON := func(v interface{}) []byte {
				j, err := json.Marshal(v)
				if err != nil {
					return []byte("{}")
				}
				return j
			}

			// Load current timer state
			state := pomoState{Mode: "idle", Running: false, EndTimeISO: "", PausedRemaining: nil}
			if b, err := os.ReadFile(timerFile); err == nil {
				_ = json.Unmarshal(b, &state)
			}

			// Helper to persist state
			save := func() {
				_ = os.WriteFile(timerFile, mustJSON(state), 0644)
			}

			// Process click action if present
			if _, err := os.Stat(stateFile); err == nil {
				// Read and remove click file
				var act struct {
					Action string `json:"action"`
				}
				if b, err := os.ReadFile(stateFile); err == nil {
					_ = json.Unmarshal(b, &act)
				}
				_ = os.Remove(stateFile)

				if enabled {
					switch act.Action {
					case "reset":
						state.Mode = "idle"
						state.Running = false
						state.EndTimeISO = ""
						state.PausedRemaining = nil
						save()
					case "toggle":
						switch {
						case state.Mode == "idle" || state.Mode == "break_complete":
							state.Mode = "work"
							state.Running = true
							state.EndTimeISO = now.Add(time.Duration(workMinutes) * time.Minute).Format(time.RFC3339)
							state.PausedRemaining = nil
							save()
						case state.Running:
							// Pause
							if state.EndTimeISO != "" {
								if t, err := time.Parse(time.RFC3339, state.EndTimeISO); err == nil {
									rem := t.Sub(now).Seconds()
									if rem < 0 {
										rem = 0
									}
									r := rem
									state.PausedRemaining = &r
								}
							}
							state.Running = false
							save()
						default:
							// Resume
							if state.PausedRemaining != nil && *state.PausedRemaining > 0 {
								end := now.Add(time.Duration(*state.PausedRemaining) * time.Second)
								state.EndTimeISO = end.Format(time.RFC3339)
								state.PausedRemaining = nil
							}
							state.Running = true
							save()
						}
					}
				}
			}

			// Compute remaining and handle transitions
			remaining := 0.0
			total := 0.0
			switch state.Mode {
			case "work":
				total = float64(workMinutes * 60)
			case "break":
				total = float64(breakMinutes * 60)
			}

			if state.Running && state.EndTimeISO != "" {
				if t, err := time.Parse(time.RFC3339, state.EndTimeISO); err == nil {
					remaining = t.Sub(now).Seconds()
				}
			} else if !state.Running && state.PausedRemaining != nil {
				remaining = *state.PausedRemaining
			}

			if remaining <= 0 {
				switch state.Mode {
				case "work":
					// Transition to break
					state.Mode = "break"
					state.Running = true
					state.EndTimeISO = now.Add(time.Duration(breakMinutes) * time.Minute).Format(time.RFC3339)
					state.PausedRemaining = nil
					save()
					remaining = float64(breakMinutes * 60)
					total = remaining
				case "break":
					// Break complete
					state.Mode = "break_complete"
					state.Running = false
					state.EndTimeISO = ""
					state.PausedRemaining = nil
					save()
				default:
					remaining = 0
				}
			}

			// Build output
			out := pomoOut{Text: "", Tooltip: "", Class: "idle", Percentage: 0}
			if !enabled {
				out.Text = "󰌾 --:--"
				out.Tooltip = "Pomodoro Timer - Disabled"
				out.Class = "disabled"
				fmt.Println(string(mustJSON(out)))
				return
			}

			switch state.Mode {
			case "idle":
				out.Text = fmt.Sprintf("󰌾 %02d:00", workMinutes)
				out.Tooltip = "Pomodoro Timer - Click to start"
				out.Class = "idle"
			case "break_complete":
				out.Text = "󰌾 Ready"
				out.Tooltip = "Break over! Click to start next session"
				out.Class = "break_complete"
			case "work", "break":
				rem := int(remaining + 0.5)
				if rem < 0 {
					rem = 0
				}
				mins := rem / 60
				secs := rem % 60
				icon := "󰔛"
				if state.Mode == "break" {
					icon = "☕"
				}
				if !state.Running {
					icon = "󰏤"
				}
				status := "Work"
				if state.Mode == "break" {
					status = "Break"
				}
				if !state.Running {
					status = "Paused"
				}
				out.Text = fmt.Sprintf("%s %02d:%02d", icon, mins, secs)
				out.Tooltip = fmt.Sprintf("%s - %02d:%02d remaining", status, mins, secs)
				if !state.Running {
					out.Class = "paused"
				} else {
					out.Class = state.Mode
				}
				if total > 0 {
					progress := 100.0 - ((float64(rem) / total) * 100.0)
					if progress < 0 {
						progress = 0
					}
					if progress > 100 {
						progress = 100
					}
					out.Percentage = int(progress + 0.5)
				}
			}

			fmt.Println(string(mustJSON(out)))
			return

		case "--waybar-memory":
			// Waybar memory module JSON: traditional percent, modern/traditional in tooltip
			type memOut struct {
				Text       string `json:"text"`
				Tooltip    string `json:"tooltip"`
				Class      string `json:"class"`
				Percentage int    `json:"percentage"`
			}
			data, err := os.ReadFile("/proc/meminfo")
			if err != nil {
				fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: cannot read /proc/meminfo","class":"critical","percentage":0}`)
				return
			}
			parse := func(key string) int64 {
				for _, ln := range strings.Split(string(data), "\n") {
					fields := strings.Fields(ln)
					if len(fields) >= 2 && strings.TrimSuffix(fields[0], ":") == key {
						val, _ := strconv.ParseInt(fields[1], 10, 64)
						return val
					}
				}
				return 0
			}
			memTotal := parse("MemTotal")
			memAvailable := parse("MemAvailable")
			memFree := parse("MemFree")
			buffers := parse("Buffers")
			cached := parse("Cached")

			if memTotal <= 0 {
				fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: invalid totals","class":"critical","percentage":0}`)
				return
			}

			usedModernKB := memTotal - memAvailable
			usedTraditionalKB := memTotal - memFree - buffers - cached

			totalGB := float64(memTotal) / (1024.0 * 1024.0)
			availGB := float64(memAvailable) / (1024.0 * 1024.0)
			usedModernGB := float64(usedModernKB) / (1024.0 * 1024.0)
			usedTraditionalGB := float64(usedTraditionalKB) / (1024.0 * 1024.0)

			percent := (float64(usedTraditionalKB) / float64(memTotal)) * 100.0

			bar := func(p float64) string {
				switch {
				case p <= 0:
					return ""
				case p <= 15:
					return "▁"
				case p <= 30:
					return "▂"
				case p <= 45:
					return "▃"
				case p <= 60:
					return "▄"
				case p <= 75:
					return "▅"
				case p <= 85:
					return "▆"
				case p <= 95:
					return "▇"
				default:
					return "█"
				}
			}(percent)

			class := "normal"
			if percent >= 90 {
				class = "critical"
			} else if percent >= 75 {
				class = "warning"
			}

			out := memOut{
				Text:       fmt.Sprintf("%s 󰾆", bar),
				Tooltip:    fmt.Sprintf("Used (Modern): %.1fGB\nUsed (Traditional): %.1fGB\nAvailable: %.1fGB\nTotal: %.1fGB (%.1f%%)", usedModernGB, usedTraditionalGB, availGB, totalGB, percent),
				Class:      class,
				Percentage: int(percent + 0.5),
			}
			if js, err := json.Marshal(out); err == nil {
				fmt.Println(string(js))
			} else {
				fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: marshal","class":"critical","percentage":0}`)
			}
			return

		case "--waybar-cpu":
			// Waybar CPU aggregate usage JSON using /proc/stat deltas
			type cpuOut struct {
				Text       string `json:"text"`
				Tooltip    string `json:"tooltip"`
				Class      string `json:"class"`
				Percentage int    `json:"percentage"`
			}
			readCPU := func() (idle, total uint64, ok bool) {
				data, err := os.ReadFile("/proc/stat")
				if err != nil {
					return 0, 0, false
				}
				for _, ln := range strings.Split(string(data), "\n") {
					if strings.HasPrefix(ln, "cpu ") {
						fields := strings.Fields(ln)
						if len(fields) < 8 {
							return 0, 0, false
						}
						parse := func(s string) uint64 {
							v, _ := strconv.ParseUint(s, 10, 64)
							return v
						}
						user := parse(fields[1])
						nice := parse(fields[2])
						system := parse(fields[3])
						idleVal := parse(fields[4])
						iowait := parse(fields[5])
						irq := parse(fields[6])
						softirq := parse(fields[7])
						steal := uint64(0)
						if len(fields) > 8 {
							steal = parse(fields[8])
						}
						idle = idleVal + iowait
						total = user + nice + system + idleVal + iowait + irq + softirq + steal
						return idle, total, true
					}
				}
				return 0, 0, false
			}
			i1, t1, ok := readCPU()
			if !ok {
				fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: cannot read /proc/stat","class":"critical","percentage":0}`)
				return
			}
			time.Sleep(120 * time.Millisecond)
			i2, t2, ok := readCPU()
			if !ok || t2 <= t1 || i2 < i1 {
				fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: invalid delta","class":"critical","percentage":0}`)
				return
			}
			dIdle := i2 - i1
			dTotal := t2 - t1
			usage := (float64(dTotal-dIdle) / float64(dTotal)) * 100.0

			bar := func(p float64) string {
				switch {
				case p <= 0:
					return ""
				case p <= 15:
					return "▁"
				case p <= 30:
					return "▂"
				case p <= 45:
					return "▃"
				case p <= 60:
					return "▄"
				case p <= 75:
					return "▅"
				case p <= 85:
					return "▆"
				case p <= 95:
					return "▇"
				default:
					return "█"
				}
			}(usage)

			class := "normal"
			if usage >= 90 {
				class = "critical"
			} else if usage >= 75 {
				class = "warning"
			}

			out := cpuOut{
				Text:       fmt.Sprintf("%s 󰍛", bar),
				Tooltip:    fmt.Sprintf("CPU Usage: %.1f%%", usage),
				Class:      class,
				Percentage: int(usage + 0.5),
			}
			if js, err := json.Marshal(out); err == nil {
				fmt.Println(string(js))
			} else {
				fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: marshal","class":"critical","percentage":0}`)
			}
			return

		case "--waybar-temp":
			// Waybar CPU temperature JSON with visual bar; sensor autodetect
			type tempOut struct {
				Text       string `json:"text"`
				Tooltip    string `json:"tooltip"`
				Class      string `json:"class"`
				Percentage int    `json:"percentage"`
			}

			// Find best sensor path
			findSensor := func() string {
				// 1) hwmon: coretemp/k10temp/zenpower -> temp1_input
				if entries, err := os.ReadDir("/sys/class/hwmon"); err == nil {
					for _, e := range entries {
						hp := filepath.Join("/sys/class/hwmon", e.Name())
						nb, err := os.ReadFile(filepath.Join(hp, "name"))
						if err != nil {
							continue
						}
						name := strings.TrimSpace(string(nb))
						switch name {
						case "coretemp", "k10temp", "zenpower":
							tf := filepath.Join(hp, "temp1_input")
							if st, err := os.Stat(tf); err == nil && !st.IsDir() {
								return tf
							}
						}
					}
				}
				// 2) thermal zones: x86_pkg_temp -> temp
				if entries, err := os.ReadDir("/sys/class/thermal"); err == nil {
					for _, e := range entries {
						if !strings.HasPrefix(e.Name(), "thermal_zone") {
							continue
						}
						zp := filepath.Join("/sys/class/thermal", e.Name())
						tb, err := os.ReadFile(filepath.Join(zp, "type"))
						if err != nil {
							continue
						}
						if strings.TrimSpace(string(tb)) == "x86_pkg_temp" {
							tf := filepath.Join(zp, "temp")
							if st, err := os.Stat(tf); err == nil && !st.IsDir() {
								return tf
							}
						}
					}
				}
				// 3) Fallback: thermal_zone0/temp
				tf := "/sys/class/thermal/thermal_zone0/temp"
				if st, err := os.Stat(tf); err == nil && !st.IsDir() {
					return tf
				}
				return ""
			}

			readTempC := func() (float64, bool) {
				s := findSensor()
				if s == "" {
					return 0, false
				}
				b, err := os.ReadFile(s)
				if err != nil {
					return 0, false
				}
				raw := strings.TrimSpace(string(b))
				// Some sensors return millidegrees, some degrees
				val, _ := strconv.ParseFloat(raw, 64)
				if val == 0 {
					// Try to parse when files contain e.g., "42000"
					if iv, e := strconv.ParseInt(raw, 10, 64); e == nil {
						if iv > 1000 {
							return float64(iv) / 1000.0, true
						}
						return float64(iv), true
					}
					return 0, false
				}
				if val > 1000 {
					return val / 1000.0, true
				}
				return val, true
			}

			tempC, ok := readTempC()
			if !ok {
				fmt.Println(`{"text":"-- 󰈸","tooltip":"Temperature sensor not available","class":"critical","percentage":0}`)
				return
			}

			// Map temperature into a bar (like the Python helper)
			// Percentage baseline: 60°C -> 0%, 95°C -> 100%
			tempPct := (tempC - 60.0) * 100.0 / 35.0
			if tempPct < 0 {
				tempPct = 0
			}
			if tempPct > 100 {
				tempPct = 100
			}
			bar := func(p float64) string {
				switch {
				case p <= 10:
					return "▁"
				case p <= 25:
					return "▂"
				case p <= 40:
					return "▃"
				case p <= 55:
					return "▄"
				case p <= 70:
					return "▅"
				case p <= 80:
					return "▆"
				case p <= 90:
					return "▇"
				default:
					return "█"
				}
			}(tempPct)

			class := "normal"
			switch {
			case tempC >= 90:
				class = "critical"
			case tempC >= 80:
				class = "warning"
			}

			out := tempOut{
				Text:       fmt.Sprintf("%s 󰈸", bar),
				Tooltip:    fmt.Sprintf("CPU Temperature: %.1f°C", tempC),
				Class:      class,
				Percentage: int(tempC + 0.5),
			}
			if js, err := json.Marshal(out); err == nil {
				fmt.Println(string(js))
			} else {
				fmt.Println(`{"text":"-- 󰈸","tooltip":"Temperature Error: marshal","class":"critical","percentage":0}`)
			}
			return

		case "--waybar-volume":
			// Waybar speaker volume JSON with visual bar (wpctl > pamixer > pactl)
			{
				have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

				vol := func() (int, bool) {
					if have("wpctl") {
						out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
						if err == nil {
							v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
							return v, true
						}
					}
					if have("pamixer") {
						// Ensure PipeWire/Pulse is ready
						if have("pactl") && exec.Command("pactl", "info").Run() != nil {
							return 0, false
						}
						out, err := exec.Command("pamixer", "--get-volume").Output()
						if err == nil {
							v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
							return v, true
						}
					}
					if have("pactl") {
						out, err := exec.Command("sh", "-lc", "pactl get-sink-volume @DEFAULT_SINK@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
						if err == nil {
							v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
							return v, true
						}
					}
					return 0, false
				}

				muted := func() (bool, bool) {
					if have("wpctl") {
						err := exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run()
						return err == nil, true
					}
					if have("pamixer") {
						// Ensure PipeWire/Pulse is ready
						if have("pactl") && exec.Command("pactl", "info").Run() != nil {
							return false, false
						}
						out, err := exec.Command("pamixer", "--get-mute").Output()
						if err == nil {
							return strings.TrimSpace(string(out)) == "true", true
						}
					}
					if have("pactl") {
						err := exec.Command("sh", "-lc", "pactl get-sink-mute @DEFAULT_SINK@ | grep -qi yes").Run()
						return err == nil, true
					}
					return false, false
				}

				v, vok := vol()
				m, mok := muted()
				if !vok || !mok {
					fmt.Println(`{"text":"▁ 󰖁","tooltip":"Audio not ready","class":"muted","percentage":0}`)
					return
				}
				if m {
					fmt.Println(`{"text":"▁ 󰖁","tooltip":"Speaker: Muted","class":"muted","percentage":0}`)
					return
				}

				var bar string
				switch {
				case v <= 2:
					bar = "▁"
				case v <= 5:
					bar = "▂"
				case v <= 10:
					bar = "▃"
				case v <= 20:
					bar = "▄"
				case v <= 35:
					bar = "▅"
				case v <= 50:
					bar = "▆"
				case v <= 75:
					bar = "▇"
				default:
					bar = "█"
				}

				icon := "󰕾"
				if v == 0 {
					icon = "󰕿"
				} else if v <= 33 {
					icon = "󰖀"
				} else if v <= 66 {
					icon = "󰕾"
				}

				class := "normal"
				if v >= 100 {
					class = "critical"
				} else if v >= 85 {
					class = "warning"
				}

				fmt.Println(fmt.Sprintf(`{"text":"%s %s","tooltip":"Speaker Volume: %d%%","class":"%s","percentage":%d}`, bar, icon, v, class, v))
			}
			return

		case "--waybar-reload":
			// Robust Waybar reload: dedupe PIDs; prefer internal --waybar-launch instead of external script
			// Dedupe if multiple Waybar PIDs exist
			if out, err := exec.Command("sh", "-lc", "pgrep -x waybar | wc -l").Output(); err == nil {
				c := strings.TrimSpace(string(out))
				if c != "" && c != "0" && c != "1" {
					_ = exec.Command("pkill", "-x", "waybar").Run()
				}
			}

			// If Waybar isn't running, start via internal launcher
			if err := exec.Command("pgrep", "-x", "waybar").Run(); err != nil {
				if self, err2 := os.Executable(); err2 == nil {
					_ = exec.Command(self, "--waybar-launch").Start()
				} else {
					_ = exec.Command("waybar").Start()
				}
				return
			}

			// Send SIGUSR2 to reload
			_ = exec.Command("pkill", "-SIGUSR2", "waybar").Run()

			// Short window to detect crash and auto-restart via internal launcher
			crashed := false
			for i := 0; i < 10; i++ {
				time.Sleep(100 * time.Millisecond)
				if err := exec.Command("pgrep", "-x", "waybar").Run(); err != nil {
					crashed = true
					break
				}
			}
			if crashed {
				if self, err := os.Executable(); err == nil {
					_ = exec.Command(self, "--waybar-launch").Start()
				} else {
					_ = exec.Command("waybar").Start()
				}
			}
			return

		case "--upgrade-smoketest":
			// Local upgrade smoke test:
			// - Parses packages.yaml and collects all packages across modules
			// - Compares with pacman -Qq
			// - Exit codes: 0 OK; 2 potential reintroductions; 3 unavailable (missing tools/files)
			// Flags: --config PATH, --json, --quiet
			home := os.Getenv("HOME")
			configPath := filepath.Join(home, ".local", "share", "archriot", "install", "packages.yaml")
			allowlistPath := filepath.Join(home, ".config", "archriot", "upgrade-allowlist.txt")
			outputJSON := false
			quiet := false

			// Parse simple flags
			for i := 2; i < len(os.Args); i++ {
				arg := os.Args[i]
				switch arg {
				case "--config":
					if i+1 < len(os.Args) {
						configPath = os.Args[i+1]
						i++
					} else {
						fmt.Fprintln(os.Stderr, "Missing value for --config")
						os.Exit(3)
					}
				case "--json":
					outputJSON = true
				case "--quiet":
					quiet = true
				case "-h", "--help":
					if !outputJSON && !quiet {
						fmt.Println("ArchRiot Local Upgrade Smoke Test")
						fmt.Println("Usage: archriot --upgrade-smoketest [--config PATH] [--json] [--quiet]")
					}
					os.Exit(0)
				default:
					if !quiet {
						fmt.Fprintf(os.Stderr, "Unknown argument: %s\n", arg)
					}
					os.Exit(3)
				}
			}

			// Ensure pacman exists
			if _, err := exec.LookPath("pacman"); err != nil {
				if outputJSON {
					payload := map[string]interface{}{
						"status":    "unavailable",
						"message":   "pacman not found in PATH; cannot assess installed packages",
						"config":    configPath,
						"allowlist": allowlistPath,
						"missing":   []string{},
					}
					if b, e := json.Marshal(payload); e == nil {
						fmt.Println(string(b))
					} else {
						fmt.Println(`{"status":"unavailable","message":"pacman not found in PATH","missing":[]}`)
					}
				} else if !quiet {
					fmt.Fprintln(os.Stderr, "pacman not found; cannot assess installed packages.")
				}
				os.Exit(3)
			}

			// Ensure config exists
			if _, err := os.Stat(configPath); err != nil {
				if outputJSON {
					payload := map[string]interface{}{
						"status":    "unavailable",
						"message":   fmt.Sprintf("packages.yaml not found at: %s", configPath),
						"config":    configPath,
						"allowlist": allowlistPath,
						"missing":   []string{},
					}
					if b, e := json.Marshal(payload); e == nil {
						fmt.Println(string(b))
					} else {
						fmt.Printf("{\"status\":\"unavailable\",\"message\":\"packages.yaml not found at: %s\",\"missing\":[]}\n", configPath)
					}
				} else if !quiet {
					fmt.Fprintf(os.Stderr, "packages.yaml not found at: %s\n", configPath)
				}
				os.Exit(3)
			}

			// Load configuration
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				if outputJSON {
					payload := map[string]interface{}{
						"status":    "unavailable",
						"message":   fmt.Sprintf("failed to load packages.yaml: %v", err),
						"config":    configPath,
						"allowlist": allowlistPath,
						"missing":   []string{},
					}
					if b, e := json.Marshal(payload); e == nil {
						fmt.Println(string(b))
					} else {
						fmt.Printf("{\"status\":\"unavailable\",\"message\":\"failed to load packages.yaml\",\"missing\":[]}\n")
					}
				} else if !quiet {
					fmt.Fprintf(os.Stderr, "failed to load packages.yaml: %v\n", err)
				}
				os.Exit(3)
			}

			// Collect desired packages from all module maps via reflection
			desired := make(map[string]struct{})
			cv := reflect.ValueOf(cfg).Elem()
			ct := cv.Type()
			for i := 0; i < cv.NumField(); i++ {
				field := cv.Field(i)
				if field.Kind() != reflect.Map || field.Type().String() != "map[string]config.Module" {
					continue
				}
				iter := field.MapRange()
				for iter.Next() {
					modVal := iter.Value()
					mod, ok := modVal.Interface().(config.Module)
					if !ok {
						continue
					}
					for _, p := range mod.Packages {
						if p = strings.TrimSpace(p); p != "" {
							desired[p] = struct{}{}
						}
					}
				}
			}
			_ = ct // silence unused in case of build tag differences

			// Load allowlist (optional)
			allow := make(map[string]struct{})
			if f, err := os.Open(allowlistPath); err == nil {
				sc := bufio.NewScanner(f)
				for sc.Scan() {
					line := sc.Text()
					if idx := strings.IndexByte(line, '#'); idx >= 0 {
						line = line[:idx]
					}
					line = strings.TrimSpace(line)
					if line != "" {
						allow[line] = struct{}{}
					}
				}
				_ = f.Close()
			}

			// Get installed packages
			out, err := exec.Command("pacman", "-Qq").Output()
			if err != nil {
				if outputJSON {
					payload := map[string]interface{}{
						"status":    "unavailable",
						"message":   "failed to query installed packages with pacman -Qq",
						"config":    configPath,
						"allowlist": allowlistPath,
						"missing":   []string{},
					}
					if b, e := json.Marshal(payload); e == nil {
						fmt.Println(string(b))
					} else {
						fmt.Println(`{"status":"unavailable","message":"failed to query installed packages","missing":[]}`)
					}
				} else if !quiet {
					fmt.Fprintln(os.Stderr, "failed to query installed packages with pacman -Qq")
				}
				os.Exit(3)
			}
			installed := make(map[string]struct{})
			{
				sc := bufio.NewScanner(strings.NewReader(string(out)))
				for sc.Scan() {
					p := strings.TrimSpace(sc.Text())
					if p != "" {
						installed[p] = struct{}{}
					}
				}
			}

			// Compute missing (present in YAML but not installed, and not allowlisted)
			var missing []string
			for p := range desired {
				if _, ok := allow[p]; ok {
					continue
				}
				if _, ok := installed[p]; !ok {
					missing = append(missing, p)
				}
			}

			// Output
			if outputJSON {
				status := "ok"
				message := "No potential reintroductions detected"
				if len(missing) > 0 {
					status = "warn"
					message = "Potential reintroductions detected"
				}
				payload := map[string]interface{}{
					"status":    status,
					"message":   message,
					"config":    configPath,
					"allowlist": allowlistPath,
					"missing":   missing,
				}
				if b, e := json.Marshal(payload); e == nil {
					fmt.Println(string(b))
				} else {
					// Minimal fallback
					fmt.Printf("{\"status\":\"%s\",\"message\":\"%s\",\"missing_count\":%d}\n", status, message, len(missing))
				}
			} else if !quiet {
				fmt.Println("ArchRiot Local Upgrade Smoke Test")
				fmt.Printf("Config:    %s\n", configPath)
				if _, err := os.Stat(allowlistPath); err == nil {
					fmt.Printf("Allowlist: %s\n", allowlistPath)
				}
				fmt.Println()
				if len(missing) == 0 {
					fmt.Println("✅ No potential reintroductions detected.")
					fmt.Println("Local upgrade appears safe with respect to previously removed packages.")
				} else {
					fmt.Println("⚠️  Potential reintroductions detected (present in packages.yaml but not installed):")
					for _, p := range missing {
						fmt.Printf("  - %s\n", p)
					}
					fmt.Println()
					fmt.Println("This suggests these packages were removed (or never installed) and would be installed by a Local upgrade.")
					if _, err := os.Stat(allowlistPath); err == nil {
						fmt.Printf("You can add specific packages to %s (one per line) to suppress this warning.\n", allowlistPath)
					}
				}
			}

			// Exit codes
			if len(missing) == 0 {
				os.Exit(0)
			} else {
				os.Exit(2)
			}
			return

		case "--stay-awake":
			// Always detach: run systemd-inhibit in its own session
			if len(os.Args) > 2 {
				// Pass the remainder as the inhibited command
				args := append([]string{"--what=sleep", "--why=ArchRiot Stay Awake"}, os.Args[2:]...)
				cmd := exec.Command("systemd-inhibit", args...)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
				// Detach stdio so parent exit doesn't affect child
				cmd.Stdin = nil
				cmd.Stdout = nil
				cmd.Stderr = nil
				_ = cmd.Start()
			} else {
				// Background inhibitor loop
				cmd := exec.Command("systemd-inhibit", "--what=sleep", "--why=ArchRiot Stay Awake", "bash", "-lc", "while :; do sleep 300; done")
				cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
				cmd.Stdin = nil
				cmd.Stdout = nil
				cmd.Stderr = nil
				_ = cmd.Start()
			}
			return

		case "--brightness":
			// Usage:
			//   archriot --brightness up
			//   archriot --brightness down
			//   archriot --brightness set <0-100>
			//   archriot --brightness get
			have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
			notify := func(title, msg, icon string) {
				// Best-effort notifications; safe if not installed
				if have("makoctl") {
					_ = exec.Command("makoctl", "dismiss", "--all").Run()
				}
				if have("notify-send") {
					_ = exec.Command("notify-send", "--replace-id=9999", "--app-name=Brightness Control", "--urgency=normal", "--icon", icon, title, msg).Start()
				}
			}
			readPct := func() string {
				out, err := exec.Command("sh", "-lc", "brightnessctl -m | cut -d, -f4 | tr -d '%'").Output()
				if err != nil {
					return "0"
				}
				return strings.TrimSpace(string(out))
			}
			iconFor := func(pct string) string {
				// Thresholds: >=75 high, >=50 medium, >=25 low, else min
				p := pct
				if len(p) == 1 {
					p = "0" + p
				}
				switch {
				case p >= "75":
					return "brightness-high"
				case p >= "50":
					return "brightness-medium"
				case p >= "25":
					return "brightness-low"
				default:
					return "brightness-min"
				}
			}
			if !have("brightnessctl") {
				notify("Brightness Error", "brightnessctl not found", "dialog-error")
				fmt.Fprintln(os.Stderr, "Error: brightnessctl is not installed")
				os.Exit(1)
			}
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: archriot --brightness [up|down|set <0-100>|get]")
				os.Exit(1)
			}
			switch os.Args[2] {
			case "up":
				_ = exec.Command("brightnessctl", "set", "5%+").Run()
				pct := readPct()
				notify("Brightness", pct+"%", iconFor(pct))
			case "down":
				_ = exec.Command("brightnessctl", "set", "5%-").Run()
				pct := readPct()
				notify("Brightness", pct+"%", iconFor(pct))
			case "set":
				if len(os.Args) < 4 {
					fmt.Fprintln(os.Stderr, "Usage: archriot --brightness set <0-100>")
					os.Exit(1)
				}
				val := strings.TrimSpace(os.Args[3])
				if val == "" {
					fmt.Fprintln(os.Stderr, "Usage: archriot --brightness set <0-100>")
					os.Exit(1)
				}
				_ = exec.Command("brightnessctl", "set", val+"%").Run()
				pct := readPct()
				notify("Brightness", pct+"%", iconFor(pct))
			case "get":
				fmt.Println(readPct())
			default:
				fmt.Fprintln(os.Stderr, "Usage: archriot --brightness [up|down|set <0-100>|get]")
				os.Exit(1)
			}
			return

		case "--volume":
			// Usage:
			//   archriot --volume toggle|inc|dec|get
			//   archriot --volume mic-toggle|mic-inc|mic-dec|mic-get
			have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
			notify := func(title, msg, icon string) {
				if have("makoctl") {
					_ = exec.Command("makoctl", "dismiss", "--all").Run()
				}
				if have("notify-send") {
					_ = exec.Command("notify-send", "--replace-id=8888", "--app-name=Volume Control", "--urgency=normal", "--icon", icon, title, msg).Start()
				}
			}
			usePamixer := have("pamixer")
			useWpctl := have("wpctl")
			usePactl := have("pactl")
			if !usePamixer && !useWpctl && !usePactl {
				notify("Volume Error", "pamixer/wpctl/pactl not found", "dialog-error")
				fmt.Fprintln(os.Stderr, "Error: pamixer/wpctl/pactl is not installed")
				os.Exit(1)
			}
			vol := func() string {
				// Prefer PipeWire (wpctl) first, then pamixer, then pactl
				if useWpctl {
					out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
					if err != nil {
						return "0"
					}
					return strings.TrimSpace(string(out))
				}
				if usePamixer {
					out, err := exec.Command("pamixer", "--get-volume").Output()
					if err != nil {
						return "0"
					}
					return strings.TrimSpace(string(out))
				}
				// pactl fallback
				out, err := exec.Command("sh", "-lc", "pactl get-sink-volume @DEFAULT_SINK@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
				if err != nil {
					return "0"
				}
				return strings.TrimSpace(string(out))
			}
			isMuted := func() bool {
				// Prefer wpctl (PipeWire), then pamixer, then pactl
				if useWpctl {
					return exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run() == nil
				}
				if usePamixer {
					return exec.Command("sh", "-lc", "pamixer --get-mute | grep -q true").Run() == nil
				}
				return exec.Command("sh", "-lc", "pactl get-sink-mute @DEFAULT_SINK@ | grep -qi yes").Run() == nil
			}
			micVol := func() string {
				// Prefer wpctl first, then pamixer, then pactl
				if useWpctl {
					out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
					if err != nil {
						return "0"
					}
					return strings.TrimSpace(string(out))
				}
				if usePamixer {
					out, err := exec.Command("pamixer", "--default-source", "--get-volume").Output()
					if err != nil {
						return "0"
					}
					return strings.TrimSpace(string(out))
				}
				// pactl fallback
				out, err := exec.Command("sh", "-lc", "pactl get-source-volume @DEFAULT_SOURCE@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
				if err != nil {
					return "0"
				}
				return strings.TrimSpace(string(out))
			}
			micMuted := func() bool {
				// Prefer wpctl (PipeWire), then pamixer, then pactl
				if useWpctl {
					return exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run() == nil
				}
				if usePamixer {
					return exec.Command("sh", "-lc", "pamixer --default-source --get-mute | grep -q true").Run() == nil
				}
				return exec.Command("sh", "-lc", "pactl get-source-mute @DEFAULT_SOURCE@ | grep -qi yes").Run() == nil
			}
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: archriot --volume [toggle|inc|dec|get|mic-toggle|mic-inc|mic-dec|mic-get]")
				os.Exit(1)
			}
			switch os.Args[2] {
			case "toggle":
				// Prefer PipeWire native control first
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") toggle").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--toggle-mute").Run()
				} else {
					_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle").Run()
				}
				if isMuted() {
					notify("Audio Muted", "Speaker: Muted", "audio-volume-muted")
				} else {
					notify("Audio Unmuted", "Speaker: "+vol()+"%", "audio-volume-high")
				}
			case "inc":
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
					_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%+").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--increase", "5", "--unmute").Run()
				} else {
					_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "0").Run()
					_ = exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "+5%").Run()
				}
				notify("Volume Up", "Speaker: "+vol()+"%", "audio-volume-high")
			case "dec":
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
					_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%-").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--decrease", "5", "--unmute").Run()
				} else {
					_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "0").Run()
					_ = exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "-5%").Run()
				}
				notify("Volume Down", "Speaker: "+vol()+"%", "audio-volume-low")
			case "get":
				fmt.Println(vol())
			case "mic-toggle":
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") toggle").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--default-source", "--toggle-mute").Run()
				} else {
					_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "toggle").Run()
				}
				if micMuted() {
					notify("Microphone Muted", "Microphone: Muted", "microphone-sensitivity-muted")
				} else {
					notify("Microphone Unmuted", "Microphone: "+micVol()+"%", "microphone-sensitivity-high")
				}
			case "mic-inc":
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
					_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%+").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--default-source", "--increase", "5", "--unmute").Run()
				} else {
					_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "0").Run()
					_ = exec.Command("pactl", "set-source-volume", "@DEFAULT_SOURCE@", "+5%").Run()
				}
				notify("Microphone Volume", "Microphone: "+micVol()+"%", "microphone-sensitivity-high")
			case "mic-dec":
				if useWpctl {
					_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
					_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%-").Run()
				} else if usePamixer {
					_ = exec.Command("pamixer", "--default-source", "--decrease", "5", "--unmute").Run()
				} else {
					_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "0").Run()
					_ = exec.Command("pactl", "set-source-volume", "@DEFAULT_SOURCE@", "-5%").Run()
				}
				notify("Microphone Volume", "Microphone: "+micVol()+"%", "microphone-sensitivity-low")
			case "mic-get":
				fmt.Println(micVol())
			default:
				fmt.Fprintln(os.Stderr, "Usage: archriot --volume [toggle|inc|dec|get|mic-toggle|mic-inc|mic-dec|mic-get]")
				os.Exit(1)
			}
			return

		case "--wifi-powersave-check":
			// Targeted WiFi power-save diagnostics (read-only)
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
			return

		case "--help-binds":
			// Print Hyprland keybindings from hyprland.conf (optional filter substring)
			home := os.Getenv("HOME")
			path := filepath.Join(home, ".config", "hypr", "hyprland.conf")
			b, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read %s: %v\n", path, err)
				os.Exit(1)
			}
			filter := ""
			if len(os.Args) > 2 {
				filter = strings.ToLower(os.Args[2])
			}
			sc := bufio.NewScanner(strings.NewReader(string(b)))
			fmt.Println("Hyprland keybindings (bind/bindm):")
			for sc.Scan() {
				line := strings.TrimSpace(sc.Text())
				if strings.HasPrefix(line, "bind") {
					if filter == "" || strings.Contains(strings.ToLower(line), filter) {
						fmt.Println(line)
					}
				}
			}
			return

		case "--help-binds-gtk":
			// Native GTK window for keybindings using yad or zenity; fallback to HTML
			self, _ := os.Executable()
			have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
			home := os.Getenv("HOME")
			outDir := filepath.Join(home, ".cache", "archriot", "help")
			_ = os.MkdirAll(outDir, 0o755)
			outFile := filepath.Join(outDir, "keybindings.txt")

			// Generate text via existing printer
			if self != "" {
				if b, err := exec.Command(self, "--help-binds").Output(); err == nil && len(b) > 0 {
					_ = os.WriteFile(outFile, b, 0o644)
				} else {
					_ = os.WriteFile(outFile, []byte("Hyprland keybindings (no data)\n"), 0o644)
				}
			}

			if have("yad") {
				_ = exec.Command("yad",
					"--title=ArchRiot — Keybindings",
					"--width=900", "--height=700", "--center",
					"--fontname=monospace 11",
					"--text-info",
					"--wrap",
					"--filename", outFile,
				).Start()
				return
			}
			if have("zenity") {
				_ = exec.Command("zenity", "--text-info",
					"--title=ArchRiot — Keybindings",
					"--width=900", "--height=700",
					"--filename", outFile,
				).Start()
				return
			}
			// Fallback HTML
			if self != "" {
				_ = exec.Command(self, "--help-binds-html").Start()
			}
			return

		case "--help-binds-web":
			// Generate HTML via script and open Brave app window with a stable class
			{
				home := os.Getenv("HOME")
				script := filepath.Join(home, ".local", "share", "archriot", "config", "bin", "scripts", "generate-keybindings-help.sh")
				outDir := filepath.Join(home, ".cache", "archriot", "help")
				outFile := filepath.Join(outDir, "keybindings.html")
				_ = os.MkdirAll(outDir, 0o755)

				// If the generator script exists, use it; otherwise fall back to internal HTML generator
				if st, err := os.Stat(script); err == nil && !st.IsDir() {
					// Use the script to generate without auto-opening
					_ = exec.Command("bash", "-lc", fmt.Sprintf("%q --output %q --no-open", script, outFile)).Run()
				} else {
					// Fallback to internal generator
					if self, err := os.Executable(); err == nil {
						_ = exec.Command(self, "--help-binds-html").Start()
					}
					return
				}

				// Build file:// URL
				url := "file://" + outFile

				// Prefer brave, then brave-browser; fallback to xdg-open if neither present
				class := "brave-archriot-keybinds"
				if _, err := exec.LookPath("brave"); err == nil {
					_ = exec.Command("brave", "--app="+url, "--class="+class).Start()
				} else if _, err := exec.LookPath("brave-browser"); err == nil {
					_ = exec.Command("brave-browser", "--app="+url, "--class="+class).Start()
				} else if _, err := exec.LookPath("xdg-open"); err == nil {
					_ = exec.Command("xdg-open", url).Start()
				}
			}
			return

		case "--help-binds-html":
			// Generate an HTML keybindings page from Hyprland config and open it
			home := os.Getenv("HOME")
			candidates := []string{
				filepath.Join(home, ".config", "hypr", "keybindings.conf"),
				filepath.Join(home, ".config", "hypr", "hyprland.conf"),
			}
			configPath := ""
			for _, p := range candidates {
				if _, err := os.Stat(p); err == nil {
					configPath = p
					break
				}
			}

			outDir := filepath.Join(home, ".cache", "archriot", "help")
			outFile := filepath.Join(outDir, "keybindings.html")
			_ = os.MkdirAll(outDir, 0o755)

			var html strings.Builder
			html.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\">")
			html.WriteString("<title>ArchRiot — Keybindings Help</title>")
			html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
			html.WriteString("<style>:root{--bg:#0f172a;--fg:#cbd5e1;--muted:#94a3b8;--accent:#7aa2f7;--line:#22263a}html,body{background:var(--bg);color:var(--fg);margin:0;padding:0;font-family:ui-monospace,monospace} .wrap{max-width:1100px;margin:0 auto;padding:28px 20px 40px} h1{color:var(--accent);font-size:28px;margin:0 0 14px;line-height:1.15} table{width:100%;border-collapse:collapse;margin:12px 0} th,td{padding:8px 10px;border-bottom:1px solid var(--line);vertical-align:top} th{text-align:left} .bind{white-space:pre-wrap}</style>")
			html.WriteString("</head><body><div class=\"wrap\">")
			html.WriteString("<h1>ArchRiot — Keybindings Help</h1>")

			type row struct {
				Bind string
				Desc string
			}
			var rows []row

			if configPath == "" {
				// No config; render a helpful page
				html.WriteString("<p>No Hyprland config found. Checked:</p><ul>")
				for _, p := range candidates {
					html.WriteString("<li>" + p + "</li>")
				}
				html.WriteString("</ul>")
			} else {
				// Load and parse binds
				data, err := os.ReadFile(configPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", configPath, err)
					os.Exit(1)
				}
				sc := bufio.NewScanner(strings.NewReader(string(data)))
				for sc.Scan() {
					line := sc.Text()
					trim := strings.TrimSpace(line)
					if !strings.HasPrefix(trim, "bind") {
						continue
					}
					// Extract description after '#', if any
					desc := ""
					if idx := strings.Index(trim, "#"); idx >= 0 {
						desc = strings.TrimSpace(trim[idx+1:])
						trim = strings.TrimSpace(trim[:idx])
					}
					// Normalize "bind=" to "bind ="
					if strings.HasPrefix(trim, "bind=") {
						trim = "bind = " + strings.TrimPrefix(trim, "bind=")
					}
					// Extract first two comma-separated fields after "bind ="
					bindStr := trim
					if i := strings.Index(bindStr, "bind ="); i >= 0 {
						bindStr = strings.TrimSpace(bindStr[i+len("bind ="):])
					}
					parts := strings.Split(bindStr, ",")
					if len(parts) < 2 {
						continue
					}
					mk := strings.TrimSpace(parts[0]) + ", " + strings.TrimSpace(parts[1])
					// $mod → SUPER, commas → " + "
					mk = strings.ReplaceAll(mk, "$mod", "SUPER")
					mk = strings.ReplaceAll(mk, ",", " + ")
					rows = append(rows, row{Bind: mk, Desc: desc})
				}
				_ = sc.Err()

				// Render table
				html.WriteString("<table><thead><tr><th>Bind</th><th>Description</th></tr></thead><tbody>")
				for _, r := range rows {
					// Simple HTML escape for &, <, >
					bindEsc := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(r.Bind)
					descEsc := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(r.Desc)
					html.WriteString("<tr><td class=\"bind\">" + bindEsc + "</td><td>" + descEsc + "</td></tr>")
				}
				html.WriteString("</tbody></table>")
				html.WriteString("<p>Source: " + configPath + "</p>")
			}
			html.WriteString("</div></body></html>")

			if err := os.WriteFile(outFile, []byte(html.String()), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", outFile, err)
				os.Exit(1)
			}

			// Best-effort open
			if _, err := exec.LookPath("xdg-open"); err == nil {
				_ = exec.Command("xdg-open", outFile).Start()
			}
			fmt.Println(outFile)
			return

		case "--fix-offscreen-windows":
			// Center off-screen floating windows (multi‑monitor aware) using hyprctl JSON
			have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
			if !have("hyprctl") {
				fmt.Fprintln(os.Stderr, "hyprctl not found in PATH")
				os.Exit(1)
			}

			// Read monitors (gather absolute geometry for all monitors)
			type mon struct {
				X      int `json:"x"`
				Y      int `json:"y"`
				Width  int `json:"width"`
				Height int `json:"height"`
			}
			monOut, err := exec.Command("hyprctl", "monitors", "-j").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read monitors: %v\n", err)
				os.Exit(1)
			}
			var mons []mon
			if err := json.Unmarshal(monOut, &mons); err != nil || len(mons) == 0 {
				fmt.Fprintln(os.Stderr, "Failed to parse monitors JSON")
				os.Exit(1)
			}
			type rect struct{ x0, y0, x1, y1 int }
			var screens []rect
			for _, m := range mons {
				screens = append(screens, rect{m.X, m.Y, m.X + m.Width, m.Y + m.Height})
			}

			// Read current workspace from activewindow
			type active struct {
				Workspace struct {
					ID int `json:"id"`
				} `json:"workspace"`
			}
			awOut, _ := exec.Command("hyprctl", "activewindow", "-j").Output()
			curWS := 1
			if len(awOut) > 0 {
				var aw active
				if json.Unmarshal(awOut, &aw) == nil && aw.Workspace.ID != 0 {
					curWS = aw.Workspace.ID
				}
			}

			// Read clients
			type client struct {
				Address   string `json:"address"`
				At        []int  `json:"at"`
				Size      []int  `json:"size"`
				Floating  bool   `json:"floating"`
				Workspace struct {
					ID int `json:"id"`
				} `json:"workspace"`
				Class string `json:"class"`
				Title string `json:"title"`
			}
			clOut, err := exec.Command("hyprctl", "clients", "-j").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read clients: %v\n", err)
				os.Exit(1)
			}
			var cls []client
			if err := json.Unmarshal(clOut, &cls); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse clients JSON")
				os.Exit(1)
			}

			// Consider a window "off-screen" when its rectangle doesn't intersect any monitor rect (with tolerance)
			overlaps := func(a rect, bx0, by0, bx1, by1 int) bool {
				// axis-aligned intersection test
				return a.x0 < bx1 && a.x1 > bx0 && a.y0 < by1 && a.y1 > by0
			}
			isOffscreen := func(x, y, w, h int) bool {
				// add a small tolerance to catch nearly invisible positions
				tol := 8
				wx0, wy0 := x, y
				wx1, wy1 := x+w, y+h
				wx0 -= tol
				wy0 -= tol
				wx1 += tol
				wy1 += tol
				for _, s := range screens {
					if overlaps(s, wx0, wy0, wx1, wy1) {
						return false
					}
				}
				return true
			}

			fixed := 0
			for _, c := range cls {
				// Only floating windows with geometry
				if !c.Floating || len(c.At) < 2 || len(c.Size) < 2 {
					continue
				}
				x, y := c.At[0], c.At[1]
				w, h := c.Size[0], c.Size[1]

				if !isOffscreen(x, y, w, h) {
					continue
				}
				ws := c.Workspace.ID
				if ws == 0 {
					ws = curWS
				}

				// Switch to window's workspace, focus, center, return
				_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", ws)).Run()
				time.Sleep(100 * time.Millisecond)
				_ = exec.Command("hyprctl", "dispatch", "focuswindow", "address:"+c.Address).Run()
				time.Sleep(120 * time.Millisecond)
				_ = exec.Command("hyprctl", "dispatch", "centerwindow").Run()
				time.Sleep(120 * time.Millisecond)
				_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", curWS)).Run()

				fmt.Printf("Centered off-screen window: %s (%s) from %d,%d size %dx%d on ws %d\n", c.Class, c.Title, x, y, w, h, ws)
				fixed++
			}

			// Desktop notifications (best-effort)
			if have("notify-send") {
				if fixed > 0 {
					_ = exec.Command("notify-send", "-t", "2500", "ArchRiot", fmt.Sprintf("Centered %d off-screen floating window(s)", fixed)).Start()
				} else {
					_ = exec.Command("notify-send", "-t", "2000", "ArchRiot", "No off-screen floating windows").Start()
				}
			}

			if fixed > 0 {
				fmt.Printf("Fixed %d off-screen floating windows\n", fixed)
			} else {
				fmt.Println("No off-screen floating windows found")
			}
			return

		case "--switch-window":
			// Window switcher using hyprctl JSON and fuzzel (no external jq)
			have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
			if !have("hyprctl") {
				fmt.Fprintln(os.Stderr, "hyprctl not found in PATH")
				os.Exit(1)
			}
			if !have("fuzzel") {
				fmt.Fprintln(os.Stderr, "fuzzel not found in PATH")
				os.Exit(1)
			}

			// Read clients
			type client struct {
				Address   string `json:"address"`
				Workspace struct {
					ID int `json:"id"`
				} `json:"workspace"`
				Class string `json:"class"`
				Title string `json:"title"`
			}
			clOut, err := exec.Command("hyprctl", "clients", "-j").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read clients: %v\n", err)
				os.Exit(1)
			}
			var cls []client
			if err := json.Unmarshal(clOut, &cls); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse clients JSON")
				os.Exit(1)
			}

			// Build selection list
			type option struct {
				Display string
				Address string
				WS      int
			}
			var opts []option
			for _, c := range cls {
				display := fmt.Sprintf("[%d] %s — %s", c.Workspace.ID, c.Title, c.Class)
				opts = append(opts, option{Display: display, Address: c.Address, WS: c.Workspace.ID})
			}
			if len(opts) == 0 {
				fmt.Println("No windows found")
				return
			}

			var b strings.Builder
			for i, o := range opts {
				if i > 0 {
					b.WriteByte('\n')
				}
				b.WriteString(o.Display)
			}

			cmd := exec.Command("fuzzel", "--dmenu", "--prompt=Switch to: ", "--width=60", "--lines=10")
			cmd.Stdin = strings.NewReader(b.String())
			sel, err := cmd.Output()
			if err != nil {
				return
			}
			choice := strings.TrimSpace(string(sel))
			if choice == "" {
				return
			}

			var chosen *option
			for i := range opts {
				if opts[i].Display == choice {
					chosen = &opts[i]
					break
				}
			}
			if chosen == nil {
				return
			}

			_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", chosen.WS)).Run()
			_ = exec.Command("hyprctl", "dispatch", "focuswindow", "address:"+chosen.Address).Run()
			_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
			return

		case "--mullvad-startup":
			// Safe, conditional Mullvad GUI startup (no-ops if not installed or not configured)
			{
				home := os.Getenv("HOME")
				logAppend := func(msg string) {
					// Best-effort append to runtime log
					logDir := filepath.Join(home, ".cache", "archriot")
					_ = os.MkdirAll(logDir, 0o755)
					logFile := filepath.Join(logDir, "runtime.log")
					if f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
						defer f.Close()
						ts := time.Now().Format("15:04:05")
						_, _ = f.WriteString("[" + ts + "] " + msg + "\n")
					}
				}

				// Basic presence checks (CLI and GUI)
				if _, err := exec.LookPath("mullvad"); err != nil {
					// Not installed; silently exit
					return
				}
				guiPath := "/opt/Mullvad VPN/mullvad-gui"
				if _, err := os.Stat(guiPath); err != nil {
					// GUI not installed; silently exit
					return
				}

				// Already running? Exit quietly
				if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
					logAppend("Mullvad GUI already running, skipping startup")
					fmt.Println("Mullvad GUI already running")
					return
				}

				// Account check: consider "status" Connected/Disconnected/Blocked as account-present,
				// or "account get" containing an account number.
				accountPresent := false
				if out, err := exec.Command("mullvad", "status").CombinedOutput(); err == nil {
					s := strings.ToLower(string(out))
					if strings.Contains(s, "connected") || strings.Contains(s, "disconnected") || strings.Contains(s, "blocked") {
						accountPresent = true
					}
				}
				if !accountPresent {
					if out, err := exec.Command("mullvad", "account", "get").CombinedOutput(); err == nil {
						if strings.Contains(string(out), "Mullvad account:") {
							accountPresent = true
						}
					}
				}
				if !accountPresent {
					logAppend("No Mullvad account configured, skipping GUI startup")
					return
				}

				// Respect Mullvad auto-connect: skip GUI when auto-connect is OFF/disabled
				if out, err := exec.Command("mullvad", "auto-connect", "get").CombinedOutput(); err == nil {
					s := strings.ToLower(string(out))
					if !(strings.Contains(s, "on") || strings.Contains(s, "enabled")) {
						logAppend("Mullvad auto-connect is OFF; skipping GUI startup")
						return
					}
				}

				// Ensure startMinimized=true in GUI settings if present
				settingsPath := filepath.Join(home, ".config", "Mullvad VPN", "gui_settings.json")
				if b, err := os.ReadFile(settingsPath); err == nil {
					type guiSettings struct {
						StartMinimized bool `json:"startMinimized"`
					}
					var gs guiSettings
					// If parsing fails, ignore and continue with defaults
					if json.Unmarshal(b, &gs) == nil {
						if !gs.StartMinimized {
							gs.StartMinimized = true
							if nb, err := json.MarshalIndent(gs, "", "  "); err == nil {
								// Preserve permissions best-effort
								_ = os.WriteFile(settingsPath, nb, 0o644)
							}
						}
					}
				}

				// Delay to allow tray to be ready
				time.Sleep(10 * time.Second)

				// Final guard: if started in the meantime, exit
				if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
					return
				}

				// Start minimized; detach
				logAppend("Starting Mullvad GUI minimized to tray")
				cmd := exec.Command(guiPath, "--minimize-to-tray")
				_ = cmd.Start()

				// Brief verify window
				time.Sleep(2 * time.Second)
				if err := exec.Command("pgrep", "-x", "mullvad-gui").Run(); err == nil {
					logAppend("✓ Mullvad GUI started successfully")
				} else {
					logAppend("⚠ Mullvad GUI did not appear to start")
				}
			}
			return

		case "--power-menu":
			// System power menu: prefer fuzzel; fallback to stdin prompt
			{
				have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
				choose := func() string {
					if have("fuzzel") {
						input := "Lock\nSuspend\nReboot\nPower Off\nCancel\n"
						cmd := exec.Command("fuzzel", "--dmenu", "--prompt=Power: ", "--width=30", "--lines=5")
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
					fmt.Println("5) Cancel")
					fmt.Print("Select [1-5]: ")
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
				default:
					return
				}
				switch act {
				case "Lock":
					if _, err := exec.LookPath("hyprlock"); err == nil {
						_ = exec.Command("hyprlock").Start()
					} else {
						_ = exec.Command("loginctl", "lock-session").Run()
					}
				case "Suspend":
					_ = exec.Command("systemctl", "suspend").Start()
				case "Reboot":
					_ = exec.Command("systemctl", "reboot").Start()
				case "Power Off":
					_ = exec.Command("systemctl", "poweroff").Start()
				}
			}
			return

		case "--setup-temperature":
			// Detect CPU temperature sensors and update Waybar Modules hwmon-path
			{
				home := os.Getenv("HOME")
				modPath := filepath.Join(home, ".config", "waybar", "Modules")

				// Detect coretemp hwmon (Intel) at /sys/class/hwmon/hwmon*/name == coretemp and temp1_input
				findCoretemp := func() string {
					entries, _ := os.ReadDir("/sys/class/hwmon")
					for _, e := range entries {
						hp := filepath.Join("/sys/class/hwmon", e.Name())
						nameB, err := os.ReadFile(filepath.Join(hp, "name"))
						if err != nil {
							continue
						}
						if strings.TrimSpace(string(nameB)) == "coretemp" {
							if _, err := os.Stat(filepath.Join(hp, "temp1_input")); err == nil {
								return filepath.Join(hp, "temp1_input")
							}
						}
					}
					return ""
				}

				// Detect x86_pkg_temp thermal zone
				findPkgTemp := func() string {
					entries, _ := os.ReadDir("/sys/class/thermal")
					for _, e := range entries {
						if !strings.HasPrefix(e.Name(), "thermal_zone") {
							continue
						}
						zp := filepath.Join("/sys/class/thermal", e.Name())
						tb, err := os.ReadFile(filepath.Join(zp, "type"))
						if err != nil {
							continue
						}
						if strings.TrimSpace(string(tb)) == "x86_pkg_temp" {
							return filepath.Join(zp, "temp")
						}
					}
					return ""
				}

				core := findCoretemp()
				pkg := findPkgTemp()

				// Build array: include detected + thermal_zone0 as fallback
				var paths []string
				if core != "" {
					paths = append(paths, fmt.Sprintf("\"%s\"", core))
				}
				if pkg != "" {
					paths = append(paths, fmt.Sprintf("\"%s\"", pkg))
				}
				paths = append(paths, "\"/sys/class/thermal/thermal_zone0/temp\"")
				hwmonArray := strings.Join(paths, ",")

				// Read Modules file and replace the hwmon-path array if present
				data, err := os.ReadFile(modPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Waybar Modules not found: %s\n", modPath)
					return
				}
				lines := strings.Split(string(data), "\n")
				var out []string
				replaced := false
				inHwmonBlock := false

				for i := 0; i < len(lines); i++ {
					line := lines[i]
					trim := strings.TrimSpace(line)

					// Match multi-line array starting at `"hwmon-path": [`
					if !replaced && !inHwmonBlock &&
						strings.HasPrefix(trim, "\"hwmon-path\"") &&
						strings.Contains(trim, "[") && !strings.Contains(trim, "]") {

						inHwmonBlock = true
						// Keep the line that opens the array
						out = append(out, line)

						// Skip until closing bracket
						i++
						for i < len(lines) && !strings.Contains(lines[i], "]") {
							i++
						}
						// Determine trailing comma based on original line
						comma := ""
						if i < len(lines) && strings.HasSuffix(strings.TrimSpace(lines[i]), "],") {
							comma = ","
						}

						// Insert our array content (indent with tabs for consistency)
						out = append(out, "\t\t\""+`hwmon-path`+"\": [")
						out = append(out, "\t\t\t"+hwmonArray)
						out = append(out, "\t\t]"+comma)

						replaced = true
						inHwmonBlock = false
						continue
					}

					// One-line array case: "hwmon-path": [ ... ],
					if !replaced && strings.HasPrefix(trim, "\"hwmon-path\"") &&
						strings.Contains(trim, "[") && strings.Contains(trim, "]") {

						indent := line[:len(line)-len(strings.TrimLeft(line, "\t "))]
						out = append(out, indent+"\"hwmon-path\": [")
						out = append(out, "\t\t"+hwmonArray)
						out = append(out, indent+"],")
						replaced = true
						continue
					}

					out = append(out, line)
				}

				if !replaced {
					fmt.Println("Note: \"hwmon-path\" not found in Waybar Modules; no changes made")
					return
				}

				if err := os.WriteFile(modPath, []byte(strings.Join(out, "\n")), 0o644); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update Waybar Modules: %v\n", err)
					return
				}

				fmt.Printf("✓ Updated Waybar temperature configuration\nPaths: %s\n", hwmonArray)
			}
			return

		default:
			fmt.Printf("Unknown option: %s\n\n", os.Args[1])
			showHelp()
			os.Exit(1)
		}
	}

	// Skip normal installation if we handled a special flag above (allow --install and --strict-abi to proceed)
	if len(os.Args) > 1 && os.Args[1] != "--ascii-only" && os.Args[1] != "--install" && os.Args[1] != "--strict-abi" {
		return
	}

	// STEP 1: Setup passwordless sudo (critical for installation)
	if err := setupSudo(); err != nil {
		log.Fatalf("❌ Sudo setup failed: %v", err)
	}

	// Optional: pre-install upgrade guard to avoid partial-upgrade ABI mismatches (e.g., wlroots/Hyprland)
	// This will update the keyring and advise or block (when --strict-abi) if risky packages are pending.
	preInstallUpgradeGuard()

	// STEP 2: Install yay AUR helper (critical for AUR packages)
	if err := installYay(); err != nil {
		log.Fatalf("❌ yay installation failed: %v", err)
	}

	// Read version from VERSION file first
	if err := version.ReadVersion(); err != nil {
		log.Fatalf("❌ Failed to read version: %v", err)
	}

	// STEP 3: Initial installation confirmation
	if !confirmInstallation() {
		fmt.Println("❌ Installation cancelled by user")
		os.Exit(0)
	}

	// STEP 4: Initialize logging
	if err := logger.InitLogging(); err != nil {
		log.Fatalf("❌ Failed to initialize logging: %v", err)
	}
	defer logger.CloseLogging()

	// Set up TUI helper functions
	tui.SetVersionGetter(func() string { return version.Get() })
	tui.SetLogPathGetter(func() string { return logger.GetLogPath() })

	// Initialize git input channel
	gitInputDone = make(chan bool, 1)

	// Initialize Secure Boot setup channel
	secureBootSetupDone = make(chan bool, 1)
	secureBootContinuationDone = make(chan bool, 1)

	// Set up reboot flag callback
	tui.SetRebootFlag = func(reboot bool) {
		shouldReboot = reboot
	}

	// Set up git credential callbacks
	tui.SetGitCallbacks(
		func(confirmed bool) {
			git.SetGitConfirm(confirmed)
			gitInputDone <- true
		},
		func(username string) {
			git.SetGitUsername(username)
		},
		func(email string) {
			git.SetGitEmail(email)
			gitInputDone <- true
		},
	)

	// Initialize TUI model
	model = tui.NewInstallModel()
	program = tea.NewProgram(model)

	// Set up unified logger with TUI program (must be first)
	logger.SetProgram(program)

	// Set up TUI emoji callback
	logger.SetTuiEmojiCallback(tui.SetEmojiMode)

	// Set up git package (after program is created)
	git.SetProgram(program)
	git.SetGitInputChannel(gitInputDone)

	// Set up Secure Boot callback
	tui.SetSecureBootCallback(func(confirmed bool) {
		secureBootSetupDone <- confirmed
	})

	// Set up Secure Boot continuation callback
	tui.SetSecureBootContinuationCallback(func(retry bool) {
		secureBootContinuationDone <- retry
	})

	// Set up failure retry callback (re-run orchestrator without exiting TUI)
	tui.SetRetryCallback(func(retry bool) {
		if retry {
			// Minimal UI feedback before retry
			program.Send(tui.LogMsg("↻ Retrying installation..."))
			program.Send(tui.StepMsg("Retrying installation..."))

			// Re-run the installation in background (same TUI/program)
			go func() {
				time.Sleep(100 * time.Millisecond)
				orchestrator.RunInstallation()
			}()
		}
	})

	// Set up orchestrator Secure Boot channel
	orchestrator.SetSecureBootSetupChannel(secureBootSetupDone)

	// Set up installer package
	installer.SetProgram(program)

	// Set up executor package
	executor.SetProgram(program)

	// Set up orchestrator package
	orchestrator.SetProgram(program)

	// Start installation in background
	go func() {
		// Small delay to let TUI initialize
		time.Sleep(100 * time.Millisecond)

		// Run installation in goroutine
		// NOTE: orchestrator.RunInstallation() will send either DoneMsg{} or FailureMsg{} itself
		// Do NOT send DoneMsg here as it overrides failure messages
		orchestrator.RunInstallation()
	}()

	// Run TUI in main thread
	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}

	// Handle reboot if requested
	if shouldReboot {
		log.Println("🔄 Preparing for system reboot...")
		log.Println("💾 Syncing filesystems...")

		// Sync filesystems
		if err := exec.Command("sync").Run(); err != nil {
			log.Printf("⚠️ Failed to sync filesystems: %v", err)
		}

		// Give time for any background processes to finish
		log.Println("⏳ Waiting for processes to complete...")
		time.Sleep(2 * time.Second)

		// Clean shutdown and reboot
		log.Println("🔄 Initiating system reboot...")
		if err := exec.Command("sudo", "shutdown", "-r", "now").Run(); err != nil {
			log.Printf("❌ Failed to reboot: %v", err)
			// Fallback to systemctl if shutdown fails
			log.Println("🔄 Trying fallback reboot method...")
			exec.Command("sudo", "systemctl", "reboot").Run()
		}
	}
}

// Optional pre-install system upgrade helpers (ABI mismatch guard)
func preInstallUpgradeGuard() {
	// If pacman is unavailable, nothing to do
	if _, err := exec.LookPath("pacman"); err != nil {
		return
	}

	// Best-effort: ensure keyring is current so subsequent queries/installs don't fail
	_ = exec.Command("sudo", "pacman", "-Sy", "--noconfirm", "archlinux-keyring").Run()

	// Check for pending upgrades; if core Wayland compositor/portal pieces are queued,
	// recommend a full system upgrade to avoid ABI mismatches (multi-monitor/wlroots issues).
	out, err := exec.Command("pacman", "-Qu").CombinedOutput()
	if err != nil {
		return
	}
	s := strings.ToLower(string(out))
	if strings.Contains(s, "hyprland") ||
		strings.Contains(s, "wlroots") ||
		strings.Contains(s, "xdg-desktop-portal-hyprland") ||
		strings.Contains(s, "wayland") {

		log.Println("⚠️  Detected pending compositor/portal updates (e.g., Hyprland/wlroots).")
		log.Println("    To avoid ABI mismatches on multi‑monitor systems, update your system before continuing:")
		log.Println("    sudo pacman -Sy archlinux-keyring && yay -Syu && yay -Yc && sudo paccache -r")

		// Enforce blocking when strict ABI mode is enabled
		if strictABI {
			log.Println("❌ Strict ABI mode: blocking installation until system is fully upgraded.")
			os.Exit(2)
		}
	}
}

// validateConfig validates the packages.yaml configuration
func validateConfig() error {
	fmt.Println("🔍 Validating packages.yaml configuration...")

	// Find configuration file
	configPath := config.FindConfigFile()
	if configPath == "" {
		return fmt.Errorf("packages.yaml not found")
	}

	fmt.Printf("📁 Found config: %s\n", configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Println("✅ YAML structure is valid")

	// Validate basic config structure
	if err := config.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	fmt.Println("✅ Required fields are present")

	// Validate dependencies
	if err := config.ValidateDependencies(cfg); err != nil {
		return fmt.Errorf("dependency validation: %w", err)
	}

	fmt.Println("✅ Dependencies are valid (no cycles, all references exist)")

	// Validate command safety
	if err := config.ValidateAllCommands(cfg); err != nil {
		return fmt.Errorf("command safety validation: %w", err)
	}

	fmt.Println("✅ Commands are safe (no dangerous patterns detected)")

	return nil
}

// runSecureBootContinuation handles the post-reboot Secure Boot setup continuation
func runSecureBootContinuation() error {
	// Read version from VERSION file first
	if err := version.ReadVersion(); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Initialize logging
	if err := logger.InitLogging(); err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}
	defer logger.CloseLogging()

	// Set up TUI helper functions
	tui.SetVersionGetter(func() string { return version.Get() })
	tui.SetLogPathGetter(func() string { return logger.GetLogPath() })

	// Initialize TUI model for Secure Boot continuation
	model := tui.NewInstallModel()
	program := tea.NewProgram(model)

	// Set up unified logger with TUI program
	logger.SetProgram(program)

	// Start Secure Boot continuation in background
	go func() {
		time.Sleep(100 * time.Millisecond)
		runSecureBootContinuationFlow(program, model)
	}()

	// Run TUI in main thread
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %v", err)
	}

	// Handle reboot if requested (same as main installer)
	if shouldReboot {
		log.Println("🔄 Preparing for system reboot...")
		log.Println("💾 Syncing filesystems...")

		// Sync filesystems
		if err := exec.Command("sync").Run(); err != nil {
			log.Printf("⚠️ Failed to sync filesystems: %v", err)
		}

		// Give time for any background processes to finish
		log.Println("⏳ Waiting for processes to complete...")
		time.Sleep(2 * time.Second)

		// Clean shutdown and reboot
		log.Println("🔄 Initiating system reboot...")
		if err := exec.Command("sudo", "shutdown", "-r", "now").Run(); err != nil {
			log.Printf("❌ Failed to reboot: %v", err)
			// Fallback to systemctl if shutdown fails
			log.Println("🔄 Trying fallback reboot method...")
			exec.Command("sudo", "systemctl", "reboot").Run()
		}
	}

	return nil
}

// runSecureBootContinuationFlow handles the actual Secure Boot continuation logic
func runSecureBootContinuationFlow(program *tea.Program, model *tui.InstallModel) {
	program.Send(tui.StepMsg("Validating Secure Boot setup..."))
	program.Send(tui.ProgressMsg(0.1))

	// Check current Secure Boot status
	sbEnabled, sbSupported, err := installer.DetectSecureBootStatus()
	if err != nil {
		program.Send(tui.LogMsg("❌ Failed to detect Secure Boot status: " + err.Error()))
		program.Send(tui.FailureMsg{Error: "Could not validate Secure Boot status"})
		return
	}

	program.Send(tui.ProgressMsg(0.3))

	if !sbSupported {
		program.Send(tui.LogMsg("❌ System does not support Secure Boot (Legacy BIOS detected)"))
		program.Send(tui.FailureMsg{Error: "Secure Boot not supported on this system"})
		return
	}

	if sbEnabled {
		// Success! Secure Boot is enabled
		program.Send(tui.StepMsg("Secure Boot successfully enabled!"))
		program.Send(tui.LogMsg("✅ Secure Boot validation successful"))
		program.Send(tui.LogMsg("🔒 LUKS encryption is now protected against memory attacks"))
		program.Send(tui.LogMsg("🎉 Setup complete - restoring normal system behavior"))
		program.Send(tui.ProgressMsg(0.8))

		// Restore hyprland.conf to normal welcome
		if err := restoreHyprlandConfig(); err != nil {
			program.Send(tui.LogMsg("⚠️ Failed to restore hyprland.conf: " + err.Error()))
		} else {
			program.Send(tui.LogMsg("✅ Restored normal startup configuration"))
		}

		program.Send(tui.ProgressMsg(1.0))
		program.Send(tui.DoneMsg{})
		return
	}

	// Secure Boot not enabled - provide detailed guidance
	program.Send(tui.StepMsg("Secure Boot setup incomplete"))
	program.Send(tui.LogMsg("⚠️ Secure Boot is supported but not yet enabled"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("📋 UEFI SETUP INSTRUCTIONS:"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("1. Restart your computer"))
	program.Send(tui.LogMsg("2. Press the UEFI/BIOS key during boot:"))
	program.Send(tui.LogMsg("   • Dell: F2 or F12"))
	program.Send(tui.LogMsg("   • HP: F10 or ESC"))
	program.Send(tui.LogMsg("   • Lenovo: F1, F2, or Enter"))
	program.Send(tui.LogMsg("   • ASUS: F2 or DEL"))
	program.Send(tui.LogMsg("   • MSI: DEL or F2"))
	program.Send(tui.LogMsg("   • Acer: F2 or DEL"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("3. Navigate to Security or Boot settings"))
	program.Send(tui.LogMsg("4. Find 'Secure Boot' option"))
	program.Send(tui.LogMsg("5. Enable Secure Boot"))
	program.Send(tui.LogMsg("6. Save settings and exit (usually F10)"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("⚠️ IMPORTANT: If system fails to boot after enabling:"))
	program.Send(tui.LogMsg("   • Return to UEFI settings"))
	program.Send(tui.LogMsg("   • Disable Secure Boot"))
	program.Send(tui.LogMsg("   • System will boot normally"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.ProgressMsg(0.7))

	// Show retry/cancel options
	program.Send(tui.LogMsg("Choose an option:"))
	program.Send(tui.LogMsg("• YES: I will reboot to UEFI settings now"))
	program.Send(tui.LogMsg("• NO: Cancel setup and restore normal system"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.SecureBootContinuationPromptMsg{})

	// Wait for user decision
	userWantsRetry := <-secureBootContinuationDone
	if userWantsRetry {
		// User chose to retry - show reboot prompt to go to UEFI
		program.Send(tui.LogMsg("🔄 You chose to continue - reboot to access UEFI settings"))
		program.Send(tui.LogMsg("⏳ This program will run again after you enable Secure Boot"))
		program.Send(tui.LogMsg(""))
		program.Send(tui.StepMsg("Reboot to enable Secure Boot in UEFI"))
		program.Send(tui.ProgressMsg(1.0))

		// Trigger reboot prompt (same as main installer)
		program.Send(tui.DoneMsg{})
	} else {
		// User chose to cancel - restore welcome and exit
		program.Send(tui.LogMsg("❌ User cancelled Secure Boot setup"))
		program.Send(tui.LogMsg("🔄 Restoring normal system behavior..."))
		if err := restoreHyprlandConfig(); err != nil {
			program.Send(tui.LogMsg("⚠️ Failed to restore hyprland.conf: " + err.Error()))
			program.Send(tui.FailureMsg{Error: "Failed to restore system configuration"})
		} else {
			program.Send(tui.LogMsg("✅ System restored - welcome will show on next login"))
			program.Send(tui.LogMsg(""))
			program.Send(tui.LogMsg("Press any key to exit"))
			program.Send(tui.ProgressMsg(1.0))
			program.Send(tui.DoneMsg{})
		}
	}
}

// restoreHyprlandConfig restores the original hyprland.conf with welcome
func restoreHyprlandConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprlandConfigPath := filepath.Join(homeDir, ".config", "hypr", "hyprland.conf")
	backupPath := hyprlandConfigPath + ".archriot-backup"

	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup hyprland.conf not found at %s", backupPath)
	}

	// Restore from backup
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("reading backup: %w", err)
	}

	if err := os.WriteFile(hyprlandConfigPath, backupData, 0644); err != nil {
		return fmt.Errorf("writing restored config: %w", err)
	}

	// Remove backup file
	os.Remove(backupPath)

	logger.LogMessage("SUCCESS", "Restored original hyprland.conf")
	return nil
}

// showHelp displays the help message
func showHelp() {
	fmt.Printf(`ArchRiot - The (Arch) Linux System You've Always Wanted

Usage:
  archriot              Run the main installer
  archriot --install    Run the main installer (explicit)
  archriot --upgrade    Launch the TUI upgrade flow
  archriot --tools      Launch optional tools interface
  archriot --validate   Validate packages.yaml configuration
  archriot --version    Show version information
  archriot --help       Show this help message

Options:
  -t, --tools                    Access optional advanced tools (Secure Boot, etc.)
  --apply-wallpaper-theme PATH   Apply dynamic theming based on wallpaper
  --toggle-dynamic-theming BOOL  Enable/disable dynamic theming (true/false)
  --strict-abi                   Block install if compositor/Wayland upgrades are pending
      --validate       Validate configuration without installing
  -v, --version        Display version information
  -h, --help           Display this help message

Examples:
  archriot             # Start installation
  archriot --tools     # Open tools menu
  archriot --validate  # Check config for errors

For more information, visit: https://github.com/CyphrRiot/ArchRiot
`)
}
