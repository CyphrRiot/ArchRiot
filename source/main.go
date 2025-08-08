package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"

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

// installYay installs yay AUR helper if not already present
func installYay() error {
	log.Printf("🔍 Checking for yay AUR helper...")

	// Check if yay is already installed
	if _, err := exec.LookPath("yay"); err == nil {
		log.Printf("✅ yay is already installed")
		return nil
	}

	log.Printf("📦 Installing yay AUR helper...")

	// Install prerequisites first
	log.Printf("🔧 Installing prerequisites...")
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", "base-devel", "git")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install yay prerequisites: %v", err)
	}

	// Create temporary directory for yay installation
	tempDir := "/tmp/yay-bin-install"
	if err := exec.Command("rm", "-rf", tempDir).Run(); err != nil {
		log.Printf("⚠️ Could not clean temp directory: %v", err)
	}

	// Clone yay-bin repository (precompiled binary version)
	log.Printf("📥 Downloading yay-bin from AUR...")
	cmd = exec.Command("git", "clone", "https://aur.archlinux.org/yay-bin.git", tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone yay-bin repository: %v", err)
	}

	// Build and install yay-bin (this downloads the precompiled binary)
	log.Printf("🔨 Installing yay-bin...")
	cmd = exec.Command("makepkg", "-si", "--noconfirm")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build yay-bin: %v", err)
	}

	// Clean up temp directory
	exec.Command("rm", "-rf", tempDir).Run()

	// Verify yay installation
	if _, err := exec.LookPath("yay"); err != nil {
		return fmt.Errorf("yay installation failed - not found in PATH")
	}

	log.Printf("✅ yay installed successfully")
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

		case "--version", "-v":
			if err := version.ReadVersion(); err != nil {
				log.Fatalf("❌ Failed to read version: %v", err)
			}
			fmt.Printf("ArchRiot version %s\n", version.Get())
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

// showHelp displays the help message
func showHelp() {
	fmt.Printf(`ArchRiot - The (Arch) Linux System You've Always Wanted

Usage:
  archriot              Run the main installer
  archriot --tools      Launch optional tools interface
  archriot --version    Show version information
  archriot --help       Show this help message

Options:
  -t, --tools          Access optional advanced tools (Secure Boot, etc.)
  -v, --version        Display version information
  -h, --help           Display this help message

Examples:
  archriot             # Start installation
  archriot --tools     # Open tools menu

For more information, visit: https://github.com/CyphrRiot/ArchRiot
`)
}
