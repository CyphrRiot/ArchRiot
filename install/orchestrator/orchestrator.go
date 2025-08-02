package orchestrator

import (
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



// RunInstallation runs the main installation process
func RunInstallation() {
	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		Program.Send(tui.ProgressMsg(progress))
	}

	// Send step updates to TUI
	sendStep := func(step string) {
		Program.Send(tui.StepMsg(step))
	}

	sendStep("Preparing system...")
	logger.Log("Progress", "System", "System Prep", "Preparing system...")
	sendProgress(0.1)

	// Sync package databases first
	if err := installer.SyncPackageDatabases(); err != nil {
		logger.Log("Error", "Database", "Database Sync", "Failed: "+err.Error())
		logger.Log("Info", "System", "Manual Fix", "Please run 'sudo pacman -Sy' manually and try again")
		return
	}

	sendStep("Loading configuration...")
	sendProgress(0.2)

	// Find config file
	configPath := config.FindConfigFile()
	if configPath == "" {
		logger.Log("Error", "File", "Config", "packages.yaml not found")
		return
	}

	logger.Log("Progress", "File", "Config Load", "Loading: "+configPath)

	// Load and validate YAML
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log("Error", "File", "Config Load", "Failed: "+err.Error())
		return
	}

	logger.Log("Success", "File", "YAML Config", "Config loaded")

	// Validate YAML configuration
	logger.Log("Progress", "File", "YAML Validation", "Validating configuration...")
	if err := config.ValidateConfig(cfg); err != nil {
		logger.Log("Error", "File", "YAML Validation", "Failed: "+err.Error())
		return
	}
	logger.Log("Success", "File", "YAML Validation", "Configuration validated")

	sendStep("Installing modules...")
	sendProgress(0.3)

	// Execute modules in proper order
	if err := executor.ExecuteModulesInOrder(cfg); err != nil {
		logger.Log("Error", "System", "Module Exec", "Failed: "+err.Error())
		return
	}

	sendStep("Installation complete!")
	sendProgress(1.0)
	logger.Log("Success", "System", "Installation", "Complete!")
	logger.Log("Success", "System", "Module Exec", "All modules done")
	logger.Log("Info", "System", "Log File", "Available at: "+logger.GetLogPath())
}
