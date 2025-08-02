package orchestrator

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/executor"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// sendFormattedLog sends a properly formatted log message to TUI
func sendFormattedLog(status, emoji, name, description string) {
	if Program != nil {
		Program.Send(tui.LogMsg(fmt.Sprintf("%s %s %-15s %s", status, emoji, name, description)))
	}
}

// RunInstallation runs the main installation process
func RunInstallation() {
	// Send log messages to TUI
	sendLog := func(msg string) {
		Program.Send(tui.LogMsg(msg))
	}

	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		Program.Send(tui.ProgressMsg(progress))
	}

	// Send step updates to TUI
	sendStep := func(step string) {
		Program.Send(tui.StepMsg(step))
	}

	sendStep("Preparing system...")
	sendLog("🔧 Preparing system...")
	sendProgress(0.1)

	// Sync package databases first
	if err := installer.SyncPackageDatabases(); err != nil {
		sendLog(fmt.Sprintf("❌ Failed to sync package databases: %v", err))
		sendLog("💡 Please run 'sudo pacman -Sy' manually and try again")
		return
	}

	sendStep("Loading configuration...")
	sendProgress(0.2)

	// Find config file
	configPath := config.FindConfigFile()
	if configPath == "" {
		sendLog("❌ packages.yaml not found")
		return
	}

	sendLog(fmt.Sprintf("📄 Loading config: %s", configPath))

	// Load and validate YAML
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		sendLog(fmt.Sprintf("❌ Failed to load config: %v", err))
		return
	}

	sendLog("✅ 📋 YAML Config      Config loaded")
	sendStep("Installing modules...")
	sendProgress(0.3)

	// Execute modules in proper order
	if err := executor.ExecuteModulesInOrder(cfg); err != nil {
		sendFormattedLog("❌", "📦", "Module Exec", "Failed: "+err.Error())
		return
	}

	sendStep("Installation complete!")
	sendProgress(1.0)
	sendFormattedLog("✅", "📦", "Installation", "Complete!")
	sendFormattedLog("✅", "📦", "Module Exec", "All modules done")
	sendFormattedLog("📝", "📋", "Log File", "Available at: "+logger.GetLogPath())
}
