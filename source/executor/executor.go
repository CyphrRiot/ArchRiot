package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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

// executeCommands runs a list of shell commands with critical/non-critical handling
func executeCommands(commands []string, moduleName string, isCritical bool) error {
	if len(commands) == 0 {
		return nil
	}

	logger.LogMessage("INFO", fmt.Sprintf("Executing %d commands for %s", len(commands), moduleName))

	for i, command := range commands {
		logger.LogMessage("INFO", fmt.Sprintf("Running command %d/%d: %s", i+1, len(commands), command))

		cmd := exec.Command("sh", "-c", command)
		if err := cmd.Run(); err != nil {
			if isCritical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL command failed for %s: %s (error: %v)", moduleName, command, err))
				return fmt.Errorf("critical command failed: %s (error: %v)", command, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical command failed for %s: %s (error: %v) - continuing installation", moduleName, command, err))
				continue
			}
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
	logger.LogMessage("INFO", "Starting module execution with dependency resolution")
	logger.Log("Progress", "System", "Module Exec", "Resolving dependencies")

	// Collect all modules with their full names
	allModules := make(map[string]config.Module)

	cfgValue := reflect.ValueOf(cfg).Elem()
	cfgType := cfgValue.Type()

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		categoryName := field.Tag.Get("yaml")
		fieldValue := cfgValue.Field(i)

		if fieldValue.Kind() == reflect.Map {
			modules := fieldValue.Interface().(map[string]config.Module)
			for moduleName, module := range modules {
				fullName := categoryName + "." + moduleName
				allModules[fullName] = module
			}
		}
	}

	// Resolve execution order based on dependencies
	executionOrder, err := resolveDependencies(allModules)
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	logger.LogMessage("INFO", fmt.Sprintf("Executing %d modules in dependency order", len(executionOrder)))
	logger.Log("Progress", "System", "Module Exec", "Starting modules")

	// Execute modules in resolved order
	for _, moduleName := range executionOrder {
		module := allModules[moduleName]

		logger.LogMessage("INFO", fmt.Sprintf("Starting module: %s - %s", moduleName, module.Start))
		logger.Log("Progress", module.Type, moduleName, module.Start)

		// Install packages
		if err := installer.InstallPackages(module.Packages); err != nil {
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL package installation FAILED for %s: %v", moduleName, err))
				logger.Log("Error", "Package", moduleName, "Critical installation failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (packages): %w", moduleName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical package installation FAILED for %s: %v - continuing installation", moduleName, err))
				logger.Log("Warning", "Package", moduleName, "Non-critical installation failed: "+err.Error())
			}
		}

		// Handle git configuration for core.identity
		if moduleName == "core.identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				if module.Critical {
					logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL git configuration FAILED for %s: %v", moduleName, err))
					logger.Log("Error", "Git", moduleName, "Critical configuration failed: "+err.Error())
					return fmt.Errorf("critical installation failed at module %s (git config): %w", moduleName, err)
				} else {
					logger.LogMessage("WARNING", fmt.Sprintf("Non-critical git configuration FAILED for %s: %v - continuing installation", moduleName, err))
					logger.Log("Warning", "Git", moduleName, "Non-critical configuration failed: "+err.Error())
				}
			}
		}

		// Copy config files
		if err := installer.CopyConfigs(module.Configs); err != nil {
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL config copying FAILED for %s: %v", moduleName, err))
				logger.Log("Error", "File", moduleName, "Critical config copy failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (config copy): %w", moduleName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical config copying FAILED for %s: %v - continuing installation", moduleName, err))
				logger.Log("Warning", "File", moduleName, "Non-critical config copy failed: "+err.Error())
			}
		}

		// Execute commands
		if err := executeCommands(module.Commands, moduleName, module.Critical); err != nil {
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL command execution FAILED for %s: %v", moduleName, err))
				logger.Log("Error", "Command", moduleName, "Critical command failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (commands): %w", moduleName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical command execution FAILED for %s: %v - continuing installation", moduleName, err))
				logger.Log("Warning", "Command", moduleName, "Non-critical command failed: "+err.Error())
			}
		}

		logger.LogMessage("SUCCESS", fmt.Sprintf("Module completed: %s - %s", moduleName, module.End))
		logger.Log("Success", module.Type, moduleName, module.End)

		// Call progress callback after each module completes
		progressCallback()
	}

	logger.LogMessage("SUCCESS", "All modules executed successfully")
	return nil
}

// resolveDependencies performs topological sort to resolve module execution order
func resolveDependencies(modules map[string]config.Module) ([]string, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize
	for name := range modules {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// Build edges
	for name, module := range modules {
		for _, dep := range module.Depends {
			if _, exists := modules[dep]; !exists {
				return nil, fmt.Errorf("module %s depends on non-existent module %s", name, dep)
			}
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	// Kahn's algorithm for topological sorting
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	var result []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for circular dependencies
	if len(result) != len(modules) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
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
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL package installation FAILED for %s: %v", fullName, err))
				logger.Log("Error", "Package", fullName, "Critical installation failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (packages): %w", fullName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical package installation FAILED for %s: %v - continuing installation", fullName, err))
				logger.Log("Warning", "Package", fullName, "Non-critical installation failed: "+err.Error())
			}
		}

		// Handle Git configuration for identity module
		if fullName == "core.identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				if module.Critical {
					logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL git configuration FAILED for %s: %v", fullName, err))
					logger.Log("Error", "Git", fullName, "Critical configuration failed: "+err.Error())
					return fmt.Errorf("critical installation failed at module %s (git config): %w", fullName, err)
				} else {
					logger.LogMessage("WARNING", fmt.Sprintf("Non-critical git configuration FAILED for %s: %v - continuing installation", fullName, err))
					logger.Log("Warning", "Git", fullName, "Non-critical configuration failed: "+err.Error())
				}
			}
		}

		// Copy configs
		if err := installer.CopyConfigs(module.Configs); err != nil {
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL config copying FAILED for %s: %v", fullName, err))
				logger.Log("Error", "File", fullName, "Critical config copy failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (config copy): %w", fullName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical config copying FAILED for %s: %v - continuing installation", fullName, err))
				logger.Log("Warning", "File", fullName, "Non-critical config copy failed: "+err.Error())
			}
		}

		// Execute commands
		if err := executeCommands(module.Commands, fullName, module.Critical); err != nil {
			if module.Critical {
				logger.LogMessage("ERROR", fmt.Sprintf("CRITICAL command execution FAILED for %s: %v", fullName, err))
				logger.Log("Error", "Command", fullName, "Critical command failed: "+err.Error())
				return fmt.Errorf("critical installation failed at module %s (commands): %w", fullName, err)
			} else {
				logger.LogMessage("WARNING", fmt.Sprintf("Non-critical command execution FAILED for %s: %v - continuing installation", fullName, err))
				logger.Log("Warning", "Command", fullName, "Non-critical command failed: "+err.Error())
			}
		}

		// Special post-hook for desktop.hyprland to re-apply control panel settings
		if fullName == "desktop.hyprland" {
			logger.Log("Progress", "Config", "Settings Restore", "Re-applying control panel settings")
			if err := reapplyControlPanelSettings(); err != nil {
				logger.Log("Warning", "Config", "Settings Restore", "Failed: "+err.Error())
			}
			// Ensure idle/lock setup
			logger.Log("Progress", "System", "Idle/Lock", "Ensuring hypridle/hyprlock and autostart")
			if err := ensureIdleLockSetup(); err != nil {
				logger.Log("Warning", "System", "Idle/Lock", "Failed: "+err.Error())
			} else {
				logger.Log("Success", "System", "Idle/Lock", "Configured")
				// Ensure Waybar single-instance launcher autostart
				logger.Log("Progress", "System", "Waybar", "Ensuring single-instance launcher autostart")
				if err := ensureWaybarLaunchSetup(); err != nil {
					logger.Log("Warning", "System", "Waybar", "Failed: "+err.Error())
				} else {
					logger.Log("Success", "System", "Waybar", "Launcher configured")
				}
			}
		}

		logger.LogMessage("INFO", fmt.Sprintf("Completed module: %s - %s", fullName, module.End))
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

// reapplyControlPanelSettings calls archriot-control-panel --reapply to restore user settings
func reapplyControlPanelSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	controlPanelPath := filepath.Join(homeDir, ".local", "share", "archriot", "config", "bin", "archriot-control-panel")
	logger.LogMessage("INFO", fmt.Sprintf("Checking for control panel at: %s", controlPanelPath))

	// Check if control panel exists
	if _, err := os.Stat(controlPanelPath); err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Control panel not found at %s - skipping settings restore", controlPanelPath))
		return nil // Not an error - control panel not installed yet
	}

	logger.LogMessage("INFO", "Control panel found, executing --reapply to restore blue light settings")

	// Execute control panel with --reapply flag
	cmd := exec.Command(controlPanelPath, "--reapply")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogMessage("ERROR", fmt.Sprintf("Control panel reapply failed: %v, output: %s", err, string(output)))
		return fmt.Errorf("control panel reapply failed: %w", err)
	}

	logger.LogMessage("SUCCESS", fmt.Sprintf("Control panel reapply completed successfully: %s", string(output)))
	return nil
}

// ensureWaybarLaunchSetup ensures Waybar is launched via the single-instance ArchRiot launcher
func ensureWaybarLaunchSetup() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprDir := filepath.Join(homeDir, ".config", "hypr")
	if err := os.MkdirAll(hyprDir, 0755); err != nil {
		return fmt.Errorf("ensuring hypr config dir: %w", err)
	}

	hyprlandConf := filepath.Join(hyprDir, "hyprland.conf")
	b, err := os.ReadFile(hyprlandConf)
	if err != nil {
		// Not fatal; user may not have a hyprland.conf yet
		logger.LogMessage("WARNING", fmt.Sprintf("Hyprland config not readable: %v", err))
		return nil
	}
	txt := string(b)

	// If already configured (either managed launcher or direct waybar), do nothing
	if strings.Contains(txt, "--waybar-launch") || strings.Contains(txt, "waybar") {
		return nil
	}

	// Append managed launcher exec-once
	if !strings.HasSuffix(txt, "\n") {
		txt += "\n"
	}
	txt += "# Autostart Waybar via ArchRiot single-instance launcher\n"
	txt += "exec-once = $HOME/.local/share/archriot/install/archriot --waybar-launch\n"

	if err := os.WriteFile(hyprlandConf, []byte(txt), 0644); err != nil {
		return fmt.Errorf("writing hyprland.conf: %w", err)
	}

	logger.LogMessage("INFO", "Added Waybar autostart via archriot --waybar-launch")
	return nil
}

// ensureIdleLockSetup ensures hypridle/hyprlock are installed and config/autostart are set
func ensureIdleLockSetup() error {
	// Ensure required packages
	if err := installer.InstallPackages([]string{"hypridle", "hyprlock"}); err != nil {
		return fmt.Errorf("installing hypridle/hyprlock: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprDir := filepath.Join(homeDir, ".config", "hypr")
	if err := os.MkdirAll(hyprDir, 0755); err != nil {
		return fmt.Errorf("ensuring hypr config dir: %w", err)
	}

	// Ensure hypridle autostart in hyprland.conf
	hyprlandConf := filepath.Join(hyprDir, "hyprland.conf")
	if b, err := os.ReadFile(hyprlandConf); err == nil {
		txt := string(b)
		if !strings.Contains(txt, "hypridle") {
			txt = txt + "\n# Autostart hypridle to enable idle lock\nexec-once = hypridle\n"
			if writeErr := os.WriteFile(hyprlandConf, []byte(txt), 0644); writeErr != nil {
				return fmt.Errorf("writing hyprland.conf: %w", writeErr)
			}
			logger.LogMessage("INFO", "Added hypridle autostart to hyprland.conf")
		}
	} else {
		// Not fatal; user may not have a hyprland.conf yet
		logger.LogMessage("WARNING", fmt.Sprintf("Hyprland config not readable: %v", err))
	}

	// Ensure minimal hypridle.conf exists
	hypridleConf := filepath.Join(hyprDir, "hypridle.conf")
	if _, err := os.Stat(hypridleConf); os.IsNotExist(err) {
		content := "general {\n  lock_cmd = hyprlock\n}\n\nlistener {\n  timeout = 600\n  on-timeout = lock\n  on-resume = unlock\n}\n"
		if writeErr := os.WriteFile(hypridleConf, []byte(content), 0644); writeErr != nil {
			return fmt.Errorf("writing hypridle.conf: %w", writeErr)
		}
		logger.LogMessage("INFO", "Created default hypridle.conf (10-minute lock)")
	}

	return nil
}
