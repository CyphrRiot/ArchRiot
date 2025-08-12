package orchestrator

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/executor"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/plymouth"
	"archriot-installer/tui"
	"archriot-installer/upgrade"
)

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// countTotalModules counts all modules across all categories
func countTotalModules(cfg *config.Config) int {
	total := 0
	total += len(cfg.Core)
	total += len(cfg.System)
	total += len(cfg.Development)
	total += len(cfg.Desktop)
	total += len(cfg.Media)
	total += len(cfg.Utilities)
	total += len(cfg.Productivity)
	total += len(cfg.Specialty)
	total += len(cfg.Theming)
	return total
}

// roundToNearest5 rounds progress to nearest 5%
func roundToNearest5(progress float64) float64 {
	return math.Round(progress*20) / 20
}

// RunInstallation runs the main installation process
func RunInstallation() {
	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		Program.Send(tui.ProgressMsg(roundToNearest5(progress)))
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
		Program.Send(tui.FailureMsg{Error: "Failed to sync package databases: " + err.Error()})
		return
	}

	sendStep("Loading configuration...")
	sendProgress(0.2)

	// Find config file
	configPath := config.FindConfigFile()
	if configPath == "" {
		logger.Log("Error", "File", "Config", "packages.yaml not found")
		Program.Send(tui.FailureMsg{Error: "Configuration file packages.yaml not found"})
		return
	}

	logger.Log("Progress", "File", "Config Load", "Loading: "+configPath)

	// Load and validate YAML
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log("Error", "File", "Config Load", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Failed to load configuration: " + err.Error()})
		return
	}

	logger.Log("Success", "File", "YAML Config", "Config loaded")

	// Validate YAML configuration
	logger.Log("Progress", "File", "YAML Validation", "Validating configuration...")
	if err := config.ValidateConfig(cfg); err != nil {
		logger.Log("Error", "File", "YAML Validation", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Configuration validation failed: " + err.Error()})
		return
	}
	logger.Log("Success", "File", "YAML Validation", "Configuration validated")

	sendStep("Installing modules...")
	sendProgress(0.3)

	// Calculate dynamic progress
	totalModules := countTotalModules(cfg)
	moduleProgressRange := 0.7 // 70% of progress is for modules (30% already used for prep)
	progressPerModule := moduleProgressRange / float64(totalModules)

	// Create progress callback for executor
	completedModules := 0
	progressCallback := func() {
		completedModules++
		currentProgress := 0.3 + (float64(completedModules) * progressPerModule)
		sendProgress(currentProgress)
	}

	// Execute modules in proper order with progress tracking
	if err := executor.ExecuteModulesInOrderWithProgress(cfg, progressCallback); err != nil {
		logger.Log("Error", "System", "Module Exec", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Module execution failed: " + err.Error()})
		return
	}

	// Optional system upgrade before Plymouth installation
	sendStep("Module installation complete!")
	sendProgress(0.90)

	upgrade.SetProgram(Program)
	if err := upgrade.PromptAndRun(); err != nil {
		logger.Log("Warning", "System", "Package Upgrade", "Failed: "+err.Error())
		// Continue anyway - upgrade failure should not stop installation
	}

	// Install Plymouth boot screen (critical system component)
	sendStep("Configuring boot screen...")
	sendProgress(0.98)

	plymouthManager, err := plymouth.NewPlymouthManager()
	if err != nil {
		logger.Log("Error", "System", "Plymouth", "Failed to initialize: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth initialization failed: " + err.Error()})
		return
	}

	if err := plymouthManager.InstallPlymouth(); err != nil {
		logger.Log("Error", "System", "Plymouth", "Installation failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth installation failed: " + err.Error()})
		return
	}

	sendStep("Installation complete!")
	sendProgress(1.0)
	logger.Log("Success", "System", "Installation", "Complete!")
	logger.Log("Success", "System", "Module Exec", "All modules done")
	logger.Log("Info", "System", "Log File", "Available at: "+logger.GetLogPath())

	// Send success completion message
	Program.Send(tui.DoneMsg{})
}
