package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/apps"
	"archriot-installer/audio"
	"archriot-installer/backgrounds"
	"archriot-installer/cleanup"
	"archriot-installer/cli"
	"archriot-installer/diagnostics"
	"archriot-installer/display"
	"archriot-installer/displays"
	"archriot-installer/executor"
	"archriot-installer/git"
	"archriot-installer/help"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/orchestrator"
	"archriot-installer/secboot"
	"archriot-installer/session"
	"archriot-installer/theming"
	"archriot-installer/tools"
	"archriot-installer/tui"
	"archriot-installer/upgradeguard"
	"archriot-installer/upgrunner"
	"archriot-installer/version"
	"archriot-installer/waybar"
	"archriot-installer/waybartools"
	"archriot-installer/windows"
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
			if err := cli.ValidateConfig(); err != nil {
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
			cli.ShowHelp()
			return

		case "--preflight":
			_ = diagnostics.Preflight()
			return

		case "--idle-diagnostics":
			_ = diagnostics.IdleDiagnostics()
			return

		case "--install":
			// Explicit install mode: continue to normal installation flow (no return)
		case "--upgrade":
			if err := upgrunner.Bootstrap(); err != nil {
				log.Fatalf("❌ Upgrade bootstrap failed: %v", err)
			}
			return
		case "--signal":
			_ = apps.RunSignal(os.Args[2:])
			return

		case "--telegram":
			_ = apps.RunTelegram(os.Args[2:])
			return

		case "--waybar-status":
			if status, _ := session.WaybarStatus(); status != "" {
				fmt.Println(status)
			} else {
				fmt.Println("stopped")
			}
			return

		case "--waybar-diagnostics":
			session.DiagnoseWaybar()
			return
		case "--waybar-sweep":
			session.SweepWaybar()
			return
		case "--waybar-restart":
			session.RestartWaybar()
			return

		case "--waybar-launch":
			_ = session.LaunchWaybar()
			return

		case "--trezor":
			_ = apps.RunTrezor(os.Args[2:])
			return

		case "--wallet":
			_ = apps.RunWallet(os.Args[2:])
			return

		case "--pomodoro-click":
			session.PomodoroClick("")
			return

		case "--pomodoro-delay-toggle":
			session.PomodoroDelayToggle()
			return

		case "--swaybg-next":
			_ = backgrounds.Next()
			return

		case "--waybar-workspace-click":
			if len(os.Args) < 3 {
				return
			}
			session.WorkspaceClick(os.Args[2])
			return

		case "--startup-background":
			session.StartupBackground()
			return

		case "--stabilize-session":
			withInhibit := false
			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] == "--inhibit" {
					withInhibit = true
					break
				}
			}
			session.StabilizeSession(withInhibit)
			return

		case "--zed":
			_ = apps.RunZed(os.Args[2:])
			return

			return

		case "--welcome":
			session.WelcomeLaunch()
			return

		case "--waybar-pomodoro":
			_ = waybar.EmitPomodoro(os.Args[2:])
			return

		case "--waybar-memory":
			_ = waybar.EmitMemory(os.Args[2:])
			return

		case "--waybar-cpu":
			_ = waybar.EmitCPU(os.Args[2:])
			return

		case "--waybar-temp":
			_ = waybar.EmitTemp(os.Args[2:])
			return

		case "--waybar-volume":
			_ = waybar.EmitVolume(os.Args[2:])
			return

		case "--waybar-reload":
			if self, err := os.Executable(); err == nil {
				_ = session.ReloadWaybar(self)
			} else {
				_ = session.ReloadWaybar("")
			}
			return

		case "--upgrade-smoketest":
			_ = tools.UpgradeSmokeTest(os.Args[2:])
			return

		case "--stay-awake":
			_ = session.Inhibit(os.Args[2:])
			return

		case "--brightness":
			_ = display.Run(os.Args[2:])
			return

		case "--volume":
			_ = audio.Run(os.Args[2:])
			return

		case "--wifi-powersave-check":
			_ = diagnostics.WifiPowerSaveCheck()
			return

		case "--help-binds":
			// Delegated to help.PrintBinds (SUPER-only)
			{
				filter := ""
				if len(os.Args) > 2 {
					filter = os.Args[2]
				}
				_ = help.PrintBinds(filter)
			}
			return

		case "--help-binds-gtk":
			// Delegated to help.OpenWeb (single canonical path)
			_ = help.OpenWeb()
			return

		case "--suspend-if-undocked":
			session.SuspendGuard()
			return

		case "--help-binds-web":
			// Delegated to help.OpenWeb (single canonical path)
			_ = help.OpenWeb()
			return

		case "--help-binds-generate":
			// Delegated to help.GenerateHTMLAndPrintPath (single canonical path)
			_ = help.GenerateHTMLAndPrintPath()
			return

		case "--help-binds-html":
			// Delegated to help.GenerateHTMLAndPrintPath (single canonical path)
			_ = help.GenerateHTMLAndPrintPath()
			return

		case "--fix-offscreen-windows":
			_ = windows.FixOffscreen()
			return

		case "--switch-window":
			_ = windows.Switcher()
			return

		case "--mullvad-startup":
			session.MullvadStartup()
			return

		case "--power-menu":
			_ = session.PowerMenu()
			return

		case "--setup-temperature":
			_ = waybartools.SetupTemperature()
			return

		case "--kanshi-autogen":
			_ = displays.Autogen()
			return

		default:
			fmt.Printf("Unknown option: %s\n\n", os.Args[1])
			cli.ShowHelp()
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
	upgradeguard.PreInstall(strictABI)

	// STEP 2: Install yay AUR helper (critical for AUR packages)
	if err := installer.EnsureYay(); err != nil {
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

	// Remove legacy shell helpers/scripts (best-effort; no-op if not present)
	cleanup.RemoveLegacyFiles()

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
		session.RebootNow()
	}
}

// Optional pre-install system upgrade helpers (ABI mismatch guard)

// validateConfig validates the packages.yaml configuration

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
		secboot.RunContinuation(program, model, secureBootContinuationDone)
	}()

	// Run TUI in main thread
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %v", err)
	}

	// Handle reboot if requested (same as main installer)
	if shouldReboot {
		session.RebootNow()
	}

	return nil
}

// restoreHyprlandConfig restores the original hyprland.conf with welcome
// Moved to secboot.RestoreHyprlandConfig()

// showHelp displays the help message
