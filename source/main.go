package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/executor"
	"archriot-installer/git"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/orchestrator"
	"archriot-installer/tools"
	"archriot-installer/tui"
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

		case "--help", "-h":
			showHelp()
			return

		default:
			fmt.Printf("Unknown option: %s\n\n", os.Args[1])
			showHelp()
			os.Exit(1)
		}
	}

	// STEP 1: Setup passwordless sudo (critical for installation)
	if err := setupSudo(); err != nil {
		log.Fatalf("❌ Sudo setup failed: %v", err)
	}

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
		// User chose to retry - they will reboot to UEFI
		program.Send(tui.LogMsg("🔄 You chose to continue - reboot and access UEFI settings"))
		program.Send(tui.LogMsg("⏳ This program will run again after you enable Secure Boot"))
		program.Send(tui.ProgressMsg(1.0))
		program.Send(tui.DoneMsg{})
	} else {
		// User chose to cancel - restore system
		program.Send(tui.LogMsg("❌ User cancelled Secure Boot setup"))
		program.Send(tui.LogMsg("🔄 Restoring normal system behavior..."))
		if err := restoreHyprlandConfig(); err != nil {
			program.Send(tui.LogMsg("⚠️ Failed to restore hyprland.conf: " + err.Error()))
			program.Send(tui.FailureMsg{Error: "Failed to restore system configuration"})
		} else {
			program.Send(tui.LogMsg("✅ System restored to normal startup"))
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
  archriot --tools      Launch optional tools interface
  archriot --validate   Validate packages.yaml configuration
  archriot --version    Show version information
  archriot --help       Show this help message

Options:
  -t, --tools          Access optional advanced tools (Secure Boot, etc.)
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
