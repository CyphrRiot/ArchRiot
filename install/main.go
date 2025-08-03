package main

import (
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/git"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/tui"
	"archriot-installer/version"
	"archriot-installer/executor"
	"archriot-installer/orchestrator"
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

func main() {
	// Read version from VERSION file first
	if err := version.ReadVersion(); err != nil {
		log.Fatalf("❌ Failed to read version: %v", err)
	}

	// Initialize logging first
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
		orchestrator.RunInstallation()

		// Signal completion
		program.Send(tui.DoneMsg{})
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
