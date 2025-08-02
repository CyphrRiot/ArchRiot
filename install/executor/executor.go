package executor

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/git"
	"archriot-installer/installer"
	"archriot-installer/logger"


)

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}



// executeCommands runs a list of shell commands
func executeCommands(commands []string, moduleName string) error {
	if len(commands) == 0 {
		return nil
	}

	logger.LogMessage("INFO", fmt.Sprintf("Executing %d commands for %s", len(commands), moduleName))

	for i, command := range commands {
		logger.LogMessage("INFO", fmt.Sprintf("Running command %d/%d: %s", i+1, len(commands), command))

		cmd := exec.Command("sh", "-c", command)
		if err := cmd.Run(); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Command failed for %s: %s (error: %v)", moduleName, command, err))
			return fmt.Errorf("command failed: %s", command)
		}

		logger.LogMessage("SUCCESS", fmt.Sprintf("Command completed: %s", command))
	}

	return nil
}

// ExecuteModulesInOrder executes all modules according to priority order
func ExecuteModulesInOrder(cfg *config.Config) error {
	logger.LogMessage("INFO", "Starting module execution in priority order")
	logger.Log("Progress", "System", "Module Exec", "Starting modules")

	// Core modules (priority 10)
	if err := executeModuleCategory("core", cfg.Core); err != nil {
		return fmt.Errorf("core modules failed: %w", err)
	}

	// System modules (priority 20)
	if err := executeModuleCategory("system", cfg.System); err != nil {
		return fmt.Errorf("system modules failed: %w", err)
	}

	// Development modules (priority 30)
	if err := executeModuleCategory("development", cfg.Development); err != nil {
		return fmt.Errorf("development modules failed: %w", err)
	}

	// Desktop modules (priority 40)
	if err := executeModuleCategory("desktop", cfg.Desktop); err != nil {
		return fmt.Errorf("desktop modules failed: %w", err)
	}

	// Media modules (priority 60 - treating as optional for now)
	if err := executeModuleCategory("media", cfg.Media); err != nil {
		return fmt.Errorf("media modules failed: %w", err)
	}

	logger.LogMessage("SUCCESS", "All module categories completed")
	return nil
}

// executeModuleCategory executes all modules in a category
func executeModuleCategory(category string, modules map[string]config.Module) error {
	if len(modules) == 0 {
		logger.LogMessage("INFO", fmt.Sprintf("No %s modules to execute", category))
		logger.Log("Info", "Module", strings.Title(category), "No modules")
		return nil
	}

	priority := config.ModuleOrder[category]
	logger.LogMessage("INFO", fmt.Sprintf("Executing %s modules (priority %d)", category, priority))
	logger.Log("Progress", "Module", strings.Title(category), "Starting "+category+" modules")

	for name, module := range modules {
		fullName := fmt.Sprintf("%s.%s", category, name)
		logger.LogMessage("INFO", fmt.Sprintf("Starting module: %s - %s", fullName, module.Start))
		logger.Log("Progress", module.Type, fullName, module.Start)

		// Install packages
		if err := installer.InstallPackages(module.Packages); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Package installation had issues for %s: %v", fullName, err))
		}

		// Handle Git configuration for identity module
		if category == "core" && name == "identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Git configuration had issues: %v", err))
				logger.Log("Warning", "Git", "Git Setup", "Issues: "+err.Error())
			}
		}

		// Copy configs
		if err := installer.CopyConfigs(module.Configs); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", fullName, err))
			logger.Log("Warning", "File", fullName, "Config issues: "+err.Error())
		}

		// Execute commands if specified
		if len(module.Commands) > 0 {
			if err := executeCommands(module.Commands, fullName); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Command execution had issues for %s: %v", fullName, err))
				logger.Log("Warning", "System", fullName, "Command issues: "+err.Error())
			}
		}



		logger.LogMessage("SUCCESS", fmt.Sprintf("Module %s completed", fullName))
		logger.Log("Complete", module.Type, fullName, module.End)
	}

	logger.LogMessage("SUCCESS", fmt.Sprintf("All %s modules completed", category))
	logger.Log("Complete", "Module", strings.Title(category), "All done")
	return nil
}
