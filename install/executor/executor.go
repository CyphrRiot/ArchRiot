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

// ExecuteModulesInOrder executes all modules in the correct priority order
func ExecuteModulesInOrder(cfg *config.Config) error {
	return ExecuteModulesInOrderWithProgress(cfg, nil)
}

// ExecuteModulesInOrderWithProgress executes all modules with progress callback
func ExecuteModulesInOrderWithProgress(cfg *config.Config, progressCallback func()) error {
	logger.LogMessage("INFO", "Starting module execution in priority order")
	logger.Log("Progress", "System", "Module Exec", "Starting modules")

	// Core modules (priority 10)
	if err := executeModuleCategoryWithProgress("core", cfg.Core, progressCallback); err != nil {
		return fmt.Errorf("core modules failed: %w", err)
	}

	// System modules (priority 20)
	if err := executeModuleCategoryWithProgress("system", cfg.System, progressCallback); err != nil {
		return fmt.Errorf("system modules failed: %w", err)
	}

	// Development modules (priority 30)
	if err := executeModuleCategoryWithProgress("development", cfg.Development, progressCallback); err != nil {
		return fmt.Errorf("development modules failed: %w", err)
	}

	// Desktop modules (priority 40)
	if err := executeModuleCategoryWithProgress("desktop", cfg.Desktop, progressCallback); err != nil {
		return fmt.Errorf("desktop modules failed: %w", err)
	}

	// Media modules (priority 50)
	if err := executeModuleCategoryWithProgress("media", cfg.Media, progressCallback); err != nil {
		return fmt.Errorf("media modules failed: %w", err)
	}

	// Utilities modules
	if err := executeModuleCategoryWithProgress("utilities", cfg.Utilities, progressCallback); err != nil {
		return fmt.Errorf("utilities modules failed: %w", err)
	}

	// Productivity modules
	if err := executeModuleCategoryWithProgress("productivity", cfg.Productivity, progressCallback); err != nil {
		return fmt.Errorf("productivity modules failed: %w", err)
	}

	// Specialty modules
	if err := executeModuleCategoryWithProgress("specialty", cfg.Specialty, progressCallback); err != nil {
		return fmt.Errorf("specialty modules failed: %w", err)
	}

	// Theming modules
	if err := executeModuleCategoryWithProgress("theming", cfg.Theming, progressCallback); err != nil {
		return fmt.Errorf("theming modules failed: %w", err)
	}

	logger.LogMessage("SUCCESS", "All modules executed successfully")
	return nil
}

// executeModuleCategory executes all modules in a category
func executeModuleCategory(category string, modules map[string]config.Module) error {
	return executeModuleCategoryWithProgress(category, modules, nil)
}

// executeModuleCategoryWithProgress executes all modules in a category with progress callback
func executeModuleCategoryWithProgress(category string, modules map[string]config.Module, progressCallback func()) error {
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
		if fullName == "core.identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Git configuration failed for %s: %v", fullName, err))
			}
		}

		// Copy configs
		if err := installer.CopyConfigs(module.Configs); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", fullName, err))
			logger.Log("Warning", "File", fullName, "Config issues: "+err.Error())
		}

		// Execute commands
		if err := executeCommands(module.Commands, fullName); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Command execution had issues for %s: %v", fullName, err))
			logger.Log("Warning", "Command", fullName, "Command issues: "+err.Error())
		}

		logger.LogMessage("SUCCESS", fmt.Sprintf("Module completed: %s - %s", fullName, module.End))
		logger.Log("Success", module.Type, fullName, module.End)

		// Call progress callback after each module completes
		if progressCallback != nil {
			progressCallback()
		}
	}

	logger.LogMessage("SUCCESS", fmt.Sprintf("%s modules completed", strings.Title(category)))
	logger.Log("Success", "Module", strings.Title(category), "All "+category+" modules done")
	return nil
}
