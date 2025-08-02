package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"gopkg.in/yaml.v3"

	"archriot-installer/config"
	"archriot-installer/git"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

const (
	CONFIG_FILE    = "packages.yaml"
	MAX_CONCURRENT = 3 // Limit concurrent package installations
)

// Global version variable read from VERSION file
var VERSION string

// Global program reference for TUI
var program *tea.Program

// Global git credentials handling
var (
	gitInputDone chan bool
)

// Global model instance
var model *tui.InstallModel

// sendFormattedLog sends a properly formatted log message to TUI
func sendFormattedLog(status, emoji, name, description string) {
	if program != nil {
		program.Send(tui.LogMsg(fmt.Sprintf("%s %s %-15s %s", status, emoji, name, description)))
	}
}

// PackageResult represents the result of a package installation
type PackageResult struct {
	Package  string
	Success  bool
	Error    error
	Duration time.Duration
}

func main() {
	// Read version from VERSION file first
	if err := readVersion(); err != nil {
		log.Fatalf("❌ Failed to read version: %v", err)
	}

	// Initialize logging first
	if err := logger.InitLogging(); err != nil {
		log.Fatalf("❌ Failed to initialize logging: %v", err)
	}
	defer logger.CloseLogging()

	// Set up TUI helper functions
	tui.SetVersionGetter(func() string { return VERSION })
	tui.SetLogPathGetter(func() string { return logger.GetLogPath() })

	// Initialize git input channel
	gitInputDone = make(chan bool, 1)

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

	// Set up git package (after program is created)
	git.SetProgram(program)
	git.SetGitInputChannel(gitInputDone)

	// Set up installer package
	installer.SetProgram(program)

	// Start installation in background
	go func() {
		// Small delay to let TUI initialize
		time.Sleep(100 * time.Millisecond)

		// Run installation in goroutine
		runInstallation()

		// Signal completion
		program.Send(tui.DoneMsg{})
	}()

	// Run TUI in main thread
	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}

func runInstallation() {
	// Send log messages to TUI
	sendLog := func(msg string) {
		program.Send(tui.LogMsg(msg))
	}

	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		program.Send(tui.ProgressMsg(progress))
	}

	// Send step updates to TUI
	sendStep := func(step string) {
		program.Send(tui.StepMsg(step))
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
	configPath := findConfigFile()
	if configPath == "" {
		sendLog("❌ packages.yaml not found")
		return
	}

	sendLog(fmt.Sprintf("📄 Loading config: %s", configPath))

	// Load and validate YAML
	config, err := loadConfig(configPath)
	if err != nil {
		sendLog(fmt.Sprintf("❌ Failed to load config: %v", err))
		return
	}

	sendLog("✅ 📋 YAML Config      Config loaded")
	sendStep("Installing modules...")
	sendProgress(0.3)

	// Execute modules in proper order
	if err := executeModulesInOrder(config); err != nil {
		sendFormattedLog("❌", "📦", "Module Exec", "Failed: "+err.Error())
		return
	}

	sendStep("Installation complete!")
	sendProgress(1.0)
	sendFormattedLog("✅", "📦", "Installation", "Complete!")
	sendFormattedLog("✅", "📦", "Module Exec", "All modules done")
	sendFormattedLog("📝", "📋", "Log File", "Available at: "+logger.GetLogPath())
}

// findConfigFile looks for packages.yaml in common locations
func findConfigFile() string {
	locations := []string{
		"packages.yaml",
		"install/packages.yaml",
		filepath.Join(os.Getenv("HOME"), ".local/share/archriot/install/packages.yaml"),
	}

	for _, path := range locations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// loadConfig reads and parses the YAML configuration
func loadConfig(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &cfg, nil
}



// readVersion reads version from VERSION file
func readVersion() error {
	// Get script directory
	scriptDir := filepath.Dir(os.Args[0])
	if scriptDir == "." {
		if wd, err := os.Getwd(); err == nil {
			scriptDir = wd
		}
	}

	// Look for VERSION file in parent directory (since we're in install/)
	versionFile := filepath.Join(filepath.Dir(scriptDir), "VERSION")

	if data, err := os.ReadFile(versionFile); err == nil {
		VERSION = strings.TrimSpace(string(data))
		return nil
	}

	// Fallback to home directory ArchRiot installation
	homeDir, err := os.UserHomeDir()
	if err != nil {
		VERSION = "unknown"
		return nil
	}

	versionFile = filepath.Join(homeDir, ".local", "share", "archriot", "VERSION")
	if data, err := os.ReadFile(versionFile); err == nil {
		VERSION = strings.TrimSpace(string(data))
		return nil
	}

	VERSION = "unknown"
	return nil
}

// executeModulesInOrder executes all modules according to priority order
func executeModulesInOrder(config *config.Config) error {
	logger.LogMessage("INFO", "Starting module execution in priority order")
	if program != nil {
		sendFormattedLog("🔄", "📦", "Module Exec", "Starting modules")
	}

	// Core modules (priority 10)
	if err := executeModuleCategory("core", config.Core); err != nil {
		return fmt.Errorf("core modules failed: %w", err)
	}

	// System modules (priority 20) - TODO: implement when we add system category
	// Development modules (priority 30)
	if err := executeModuleCategory("development", config.Development); err != nil {
		return fmt.Errorf("development modules failed: %w", err)
	}

	// Desktop modules (priority 40)
	if err := executeModuleCategory("desktop", config.Desktop); err != nil {
		return fmt.Errorf("desktop modules failed: %w", err)
	}

	// Media modules (priority 60 - treating as optional for now)
	if err := executeModuleCategory("media", config.Media); err != nil {
		return fmt.Errorf("media modules failed: %w", err)
	}

	logger.LogMessage("SUCCESS", "All module categories completed")
	return nil
}

// executeModuleCategory executes all modules in a category
func executeModuleCategory(category string, modules map[string]config.Module) error {
	if len(modules) == 0 {
		logger.LogMessage("INFO", fmt.Sprintf("No %s modules to execute", category))
		if program != nil {
			sendFormattedLog("📋", "📦", strings.Title(category), "No modules")
		}
		return nil
	}

	priority := config.ModuleOrder[category]
	logger.LogMessage("INFO", fmt.Sprintf("Executing %s modules (priority %d)", category, priority))
	if program != nil {
		sendFormattedLog("🔄", "📦", strings.Title(category), "Starting "+category+" modules")
	}

	for name, module := range modules {
		fullName := fmt.Sprintf("%s.%s", category, name)
		logger.LogMessage("INFO", fmt.Sprintf("Starting module: %s - %s", fullName, module.Description))
		if program != nil {
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
				if program != nil {
					sendFormattedLog("⚠️", "📦", "Git Setup", "Issues: "+err.Error())
				}
			}
		}

		// Copy configs
		if err := installer.CopyConfigs(module.Configs); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", fullName, err))
			if program != nil {
				sendFormattedLog("⚠️", "📁", fullName, "Config issues: "+err.Error())
			}
		}

		logger.LogMessage("SUCCESS", fmt.Sprintf("Module %s completed", fullName))
		if program != nil {
			sendFormattedLog("✅", "📦", fullName, "Complete")
		}
	}

	logger.LogMessage("SUCCESS", fmt.Sprintf("All %s modules completed", category))
	if program != nil {
		sendFormattedLog("✅", "📦", strings.Title(category), "All done")
	}
	return nil
}
