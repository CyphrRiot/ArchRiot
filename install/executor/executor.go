package executor

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/git"
	"archriot-installer/handlers"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
	handlers.SetProgram(p)
}

// sendFormattedLog sends a properly formatted log message to TUI
func sendFormattedLog(status, emoji, name, description string) {
	if Program != nil {
		Program.Send(tui.LogMsg(fmt.Sprintf("%s %s %-15s %s", status, emoji, name, description)))
	}
}

// ExecuteModulesInOrder executes all modules according to priority order
func ExecuteModulesInOrder(cfg *config.Config) error {
	logger.LogMessage("INFO", "Starting module execution in priority order")
	if Program != nil {
		sendFormattedLog("🔄", "📦", "Module Exec", "Starting modules")
	}

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
		if Program != nil {
			sendFormattedLog("📋", "📦", strings.Title(category), "No modules")
		}
		return nil
	}

	priority := config.ModuleOrder[category]
	logger.LogMessage("INFO", fmt.Sprintf("Executing %s modules (priority %d)", category, priority))
	if Program != nil {
		sendFormattedLog("🔄", "📦", strings.Title(category), "Starting "+category+" modules")
	}

	for name, module := range modules {
		fullName := fmt.Sprintf("%s.%s", category, name)
		logger.LogMessage("INFO", fmt.Sprintf("Starting module: %s - %s", fullName, module.Description))
		if Program != nil {
			sendFormattedLog("🔄", "📦", fullName, module.Description)
		}

		// Install packages
		if err := installer.InstallPackages(module.Packages); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Package installation had issues for %s: %v", fullName, err))
		}

		// Handle Git configuration for identity module
		if category == "core" && name == "identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Git configuration had issues: %v", err))
				if Program != nil {
					sendFormattedLog("⚠️", "📦", "Git Setup", "Issues: "+err.Error())
				}
			}
		}

		// Copy configs
		if err := installer.CopyConfigs(module.Configs); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", fullName, err))
			if Program != nil {
				sendFormattedLog("⚠️", "📁", fullName, "Config issues: "+err.Error())
			}
		}

		// Execute handler if specified
		if module.Handler != "" {
			if err := handlers.ExecuteHandler(module.Handler); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Handler execution had issues for %s: %v", fullName, err))
				if Program != nil {
					sendFormattedLog("⚠️", "🔧", fullName, "Handler issues: "+err.Error())
				}
			}
		}

		logger.LogMessage("SUCCESS", fmt.Sprintf("Module %s completed", fullName))
		if Program != nil {
			sendFormattedLog("✅", "📦", fullName, "Complete")
		}
	}

	logger.LogMessage("SUCCESS", fmt.Sprintf("All %s modules completed", category))
	if Program != nil {
		sendFormattedLog("✅", "📦", strings.Title(category), "All done")
	}
	return nil
}
