package orchestrator

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

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
	sendProgress(0.95)

	plymouthManager, err := plymouth.NewPlymouthManager()
	if err != nil {
		logger.Log("Error", "System", "Plymouth", "Failed to initialize: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth initialization failed: " + err.Error()})
		return
	}

	sendProgress(0.96)
	if err := plymouthManager.InstallPlymouth(); err != nil {
		logger.Log("Error", "System", "Plymouth", "Installation failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth installation failed: " + err.Error()})
		return
	}

	sendProgress(0.98)
	sendStep("Checking Secure Boot recommendation...")

	// Check Secure Boot recommendation before completion
	checkSecureBootRecommendation()

	sendStep("Installation complete!")
	sendProgress(1.0)
	logger.Log("Success", "System", "Installation", "Complete!")
	logger.Log("Success", "System", "Module Exec", "All modules done")
	logger.Log("Info", "System", "Log File", "Available at: "+logger.GetLogPath())

	// Send success completion message
	Program.Send(tui.DoneMsg{})
}

// secureBootSetupDone channel for synchronization
var secureBootSetupDone chan bool

// SetSecureBootSetupChannel sets the channel for Secure Boot setup synchronization
func SetSecureBootSetupChannel(ch chan bool) {
	secureBootSetupDone = ch
}

// checkSecureBootRecommendation checks if Secure Boot should be recommended and waits for user decision
func checkSecureBootRecommendation() {
	logger.Log("Info", "System", "SecureBoot", "Checking Secure Boot and LUKS status...")

	// Detect Secure Boot status
	sbEnabled, sbSupported, err := installer.DetectSecureBootStatus()
	if err != nil {
		logger.Log("Warning", "System", "SecureBoot", "Failed to detect Secure Boot status: "+err.Error())
		return
	}

	// Detect LUKS encryption
	luksUsed, luksDevices, err := installer.DetectLuksEncryption()
	if err != nil {
		logger.Log("Warning", "System", "SecureBoot", "Failed to detect LUKS encryption: "+err.Error())
		luksUsed = false
		luksDevices = []string{}
	}

	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("Secure Boot: enabled=%v, supported=%v", sbEnabled, sbSupported))
	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("LUKS: detected=%v, devices=%v", luksUsed, luksDevices))

	// Send status to TUI
	Program.Send(tui.SecureBootStatusMsg{
		Enabled:     sbEnabled,
		Supported:   sbSupported,
		LuksUsed:    luksUsed,
		LuksDevices: luksDevices,
	})

	// Prompt user if Secure Boot should be enabled for LUKS protection
	if !sbEnabled && sbSupported && luksUsed {
		logger.Log("Info", "System", "SecureBoot", "Prompting user for Secure Boot enablement...")
		Program.Send(tui.SecureBootPromptMsg{})

		// Wait for user decision
		if secureBootSetupDone != nil {
			userWantsSecureBoot := <-secureBootSetupDone
			if userWantsSecureBoot {
				Program.Send(tui.StepMsg("Setting up Secure Boot..."))
				if err := performSecureBootSetup(); err != nil {
					logger.Log("Error", "System", "SecureBoot", "Setup failed: "+err.Error())
					Program.Send(tui.FailureMsg{Error: "Secure Boot setup failed: " + err.Error()})
					return
				}
				logger.Log("Success", "System", "SecureBoot", "Setup completed successfully")
			} else {
				logger.Log("Info", "System", "SecureBoot", "User declined Secure Boot setup")
			}
		}
	}
}

// performSecureBootSetup executes the Secure Boot setup process
func performSecureBootSetup() error {
	// Validate prerequisites
	if err := installer.ValidateSecureBootPrerequisites(); err != nil {
		return fmt.Errorf("prerequisites not met: %w", err)
	}

	// Modify hyprland.conf for continuation after reboot
	if err := modifyHyprlandForContinuation(); err != nil {
		return fmt.Errorf("failed to setup continuation: %w", err)
	}

	// TODO: Implement actual Secure Boot setup steps for Phase 2:
	// 1. Generate custom signing keys
	// 2. Sign bootloader and kernel
	// 3. Set up pacman hooks for future updates
	// 4. Prepare UEFI guidance for user

	logger.Log("Info", "System", "SecureBoot", "Secure Boot setup prepared - enable in UEFI settings after reboot")
	return nil
}

// modifyHyprlandForContinuation modifies hyprland.conf to launch Secure Boot continuation after reboot
func modifyHyprlandForContinuation() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprlandConfigPath := filepath.Join(homeDir, ".config", "hypr", "hyprland.conf")
	backupPath := hyprlandConfigPath + ".archriot-backup"

	// Read current hyprland.conf
	configData, err := os.ReadFile(hyprlandConfigPath)
	if err != nil {
		return fmt.Errorf("reading hyprland.conf: %w", err)
	}

	// Backup original
	if err := os.WriteFile(backupPath, configData, 0644); err != nil {
		return fmt.Errorf("creating backup: %w", err)
	}

	// Get ArchRiot binary path
	archRiotPath := filepath.Join(homeDir, ".local", "share", "archriot", "install", "archriot")

	// Modify config to replace welcome with Secure Boot continuation
	configStr := string(configData)
	lines := strings.Split(configStr, "\n")

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "exec-once") && strings.Contains(trimmedLine, "welcome") {
			// Replace with Secure Boot continuation
			lines[i] = fmt.Sprintf("exec-once = %s --secure_boot_stage", archRiotPath)
			logger.Log("Info", "System", "SecureBoot", "Modified hyprland.conf for Secure Boot continuation")
			break
		}
	}

	// Write modified config
	modifiedConfig := strings.Join(lines, "\n")
	if err := os.WriteFile(hyprlandConfigPath, []byte(modifiedConfig), 0644); err != nil {
		// Restore backup on failure
		os.WriteFile(hyprlandConfigPath, configData, 0644)
		return fmt.Errorf("writing modified config: %w", err)
	}

	logger.Log("Success", "System", "SecureBoot", "Hyprland configuration modified for post-reboot continuation")
	return nil
}
