package executor

import (
	"fmt"
	"os/exec"
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
			logger.LogMessage("WARNING", fmt.Sprintf("Package installation had issues for %s: %v", moduleName, err))
		}

		// Handle git configuration for core.identity
		if moduleName == "core.identity" {
			if err := git.HandleGitConfiguration(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Git configuration failed for %s: %v", moduleName, err))
			}
		}

		// Copy config files
		if err := installer.CopyConfigs(module.Configs); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", moduleName, err))
			logger.Log("Warning", "File", moduleName, "Config issues: "+err.Error())
		}

		// Execute commands
		if err := executeCommands(module.Commands, moduleName); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Command execution had issues for %s: %v", moduleName, err))
			logger.Log("Warning", "Command", moduleName, "Command issues: "+err.Error())
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
