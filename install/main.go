package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"gopkg.in/yaml.v3"

	"archriot-installer/logger"
	"archriot-installer/tui"
)

// Config represents the YAML structure
type Config struct {
	Core        map[string]Module `yaml:"core"`
	Desktop     map[string]Module `yaml:"desktop"`
	Development map[string]Module `yaml:"development"`
	Media       map[string]Module `yaml:"media"`
}

// Module represents a single installation module
type Module struct {
	Packages    []string     `yaml:"packages"`
	Configs     []ConfigRule `yaml:"configs"`
	Handler     string       `yaml:"handler,omitempty"`
	Depends     []string     `yaml:"depends,omitempty"`
	Description string       `yaml:"description"`
}

// ConfigRule represents a configuration copying rule
type ConfigRule struct {
	Pattern          string   `yaml:"pattern"`
	PreserveIfExists []string `yaml:"preserve_if_exists,omitempty"`
}

// ModuleOrder defines the execution order for different module categories
var ModuleOrder = map[string]int{
	"core":         10,
	"system":       20,
	"development":  30,
	"desktop":      40,
	"post-desktop": 45,
	"applications": 50,
	"optional":     60,
}

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
	gitUsername   string
	gitEmail      string
	gitConfirmUse bool
	gitInputDone  chan bool
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
			gitConfirmUse = confirmed
			gitInputDone <- true
		},
		func(username string) {
			gitUsername = username
		},
		func(email string) {
			gitEmail = email
			gitInputDone <- true
		},
	)

	// Initialize TUI model
	model = tui.NewInstallModel()
	program = tea.NewProgram(model)

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
	if err := syncPackageDatabases(); err != nil {
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
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &config, nil
}

// installPackages installs packages in a single batch to avoid database locks
func installPackages(packages []string) error {
	if len(packages) == 0 {
		if program != nil {
			sendFormattedLog("📋", "📦", "Packages", "None to install")
		}
		return nil
	}

	// Log package installation to file only, not TUI

	// Filter out already installed packages for cleaner output
	needed, alreadyInstalled := checkPackageStatus(packages)

	if len(alreadyInstalled) > 0 {
		// Already installed packages logged to file only
	}

	if len(needed) == 0 {
		// All packages already installed - no TUI spam
		return nil
	}

	// Install all needed packages in one batch
	if err := installPackageBatch(needed); err != nil {
		return fmt.Errorf("batch installation failed: %w", err)
	}

	// Package installation success logged to file only
	return nil
}

// checkPackageStatus separates already installed packages from needed ones
func checkPackageStatus(packages []string) (needed []string, alreadyInstalled []string) {
	for _, pkg := range packages {
		if isPackageInstalled(pkg) {
			alreadyInstalled = append(alreadyInstalled, pkg)
		} else {
			needed = append(needed, pkg)
		}
	}
	return needed, alreadyInstalled
}

// isPackageInstalled checks if a package is already installed
func isPackageInstalled(pkg string) bool {
	cmd := exec.Command("pacman", "-Q", pkg)
	err := cmd.Run()
	return err == nil
}

// installPackageBatch installs multiple packages in a single transaction
func installPackageBatch(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	var cmd *exec.Cmd

	if commandExists("yay") {
		args := append([]string{"-S", "--noconfirm", "--needed"}, packages...)
		cmd = exec.Command("yay", args...)
	} else if commandExists("pacman") {
		args := append([]string{"pacman", "-S", "--noconfirm", "--needed"}, packages...)
		cmd = exec.Command("sudo", args...)
	} else {
		return fmt.Errorf("no package manager found (yay or pacman)")
	}

	// Command execution logged to file only, not TUI

	output, err := cmd.CombinedOutput()
	if err != nil {
		if program != nil {
			// Limit output to first 200 characters to prevent TUI spam
			outputStr := string(output)
			if len(outputStr) > 200 {
				outputStr = outputStr[:200] + "... (truncated)"
			}
			sendFormattedLog("❌", "📦", "Package Error", "Failed: "+outputStr)
		}
		return fmt.Errorf("batch installation failed: %w", err)
	}

	return nil
}

// syncPackageDatabases ensures yay is installed then synchronizes pacman and yay package databases
func syncPackageDatabases() error {
	logger.LogMessage("INFO", "🔄 Syncing package databases...")
	if program != nil {
		sendFormattedLog("🔄", "📦", "Database Sync", "Syncing databases")
	}

	start := time.Now()

	// Install yay if not present
	if !commandExists("yay") {
		logger.LogMessage("INFO", "Installing yay AUR helper...")
		// Yay installation logged to file only

		cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", "yay")
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Limit output to prevent log spam
			outputStr := string(output)
			if len(outputStr) > 200 {
				outputStr = outputStr[:200] + "... (truncated)"
			}
			logger.LogMessage("CRITICAL", fmt.Sprintf("Failed to install yay: %s", outputStr))
			return fmt.Errorf("yay installation failed: %w", err)
		}

		logger.LogMessage("SUCCESS", "yay AUR helper installed")
		// Yay installation success logged to file only
	} else {
		logger.LogMessage("SUCCESS", "yay AUR helper already available")
		// Yay already available logged to file only
	}

	// Sync pacman database - fail gracefully if it doesn't work
	if err := syncPackmanDatabase(); err != nil {
		logger.LogMessage("CRITICAL", fmt.Sprintf("Pacman sync failed: %v", err))
		return fmt.Errorf("pacman sync failed - please run 'sudo pacman -Sy' manually: %w", err)
	}

	// Sync yay database
	cmd := exec.Command("yay", "-Sy", "--noconfirm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Limit output to prevent log spam
		outputStr := string(output)
		if len(outputStr) > 200 {
			outputStr = outputStr[:200] + "... (truncated)"
		}
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to sync yay database: %s", outputStr))
		if program != nil {
			sendFormattedLog("⚠️", "📦", "Database Sync", "Yay sync failed, continuing")
		}
	} else {
		logger.LogMessage("SUCCESS", "Yay database synced")
		// Yay sync success logged to file only
	}

	duration := time.Since(start)
	logger.LogMessage("SUCCESS", fmt.Sprintf("Package databases synced in %v", duration))
	// Database sync timing logged to file only

	return nil
}

// syncPackmanDatabase attempts to sync pacman database
func syncPackmanDatabase() error {
	cmd := exec.Command("sudo", "pacman", "-Sy", "--noconfirm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Limit output to prevent massive error dumps
		outputStr := string(output)
		if len(outputStr) > 200 {
			outputStr = outputStr[:200] + "... (truncated)"
		}
		return fmt.Errorf("pacman sync failed: %s", outputStr)
	}

	logger.LogMessage("SUCCESS", "Pacman database synced")
	// Pacman sync success logged to file only
	return nil
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



// copyConfigs copies configuration files with preservation logic
func copyConfigs(configs []ConfigRule) error {
	if len(configs) == 0 {
		if program != nil {
			sendFormattedLog("📋", "📁", "Config Copy", "None to copy")
		}
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Use proper ArchRiot installation directory for config source
	configSourceDir := filepath.Join(homeDir, ".local", "share", "archriot", "config")

	// logMessage("INFO", fmt.Sprintf("Copying configs from: %s", configSourceDir))
	if program != nil {
		sendFormattedLog("🔄", "📁", "Config Copy", "From: "+configSourceDir)
	}

	for _, configRule := range configs {
		// logMessage("INFO", fmt.Sprintf("Processing config pattern: %s", configRule.Pattern))

		if err := copyConfigPattern(configSourceDir, homeDir, configRule); err != nil {
			// logMessage("WARNING", fmt.Sprintf("Failed to copy config %s: %v", configRule.Pattern, err))
			if program != nil {
				sendFormattedLog("❌", "📄", configRule.Pattern, "Failed: "+err.Error())
			}
		} else {
			if program != nil {
				sendFormattedLog("✅", "📄", configRule.Pattern, "Copied successfully")
			}
		}
	}

	return nil
}

// copyConfigPattern copies files matching a config pattern with preservation
func copyConfigPattern(sourceDir, homeDir string, rule ConfigRule) error {
	// Parse pattern (e.g., "hypr/*" -> source: config/hypr, dest: ~/.config/hypr)
	pattern := rule.Pattern
	var sourcePath, destPath string

	if strings.HasSuffix(pattern, "/*") {
		// Directory pattern like "hypr/*"
		dirName := strings.TrimSuffix(pattern, "/*")
		sourcePath = filepath.Join(sourceDir, dirName)
		destPath = filepath.Join(homeDir, ".config", dirName)
	} else {
		// Single file pattern
		sourcePath = filepath.Join(sourceDir, pattern)
		destPath = filepath.Join(homeDir, ".config", pattern)
	}

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source not found: %s", sourcePath)
	}

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("creating dest directory: %w", err)
	}

	// Copy files
	return copyFileOrDirectory(sourcePath, destPath, rule.PreserveIfExists)
}

// copyFileOrDirectory recursively copies files or directories with preservation
func copyFileOrDirectory(source, dest string, preserveFiles []string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return copyDirectory(source, dest, preserveFiles)
	}
	return copyFile(source, dest, preserveFiles)
}

// copyDirectory recursively copies a directory
func copyDirectory(source, dest string, preserveFiles []string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(sourcePath, destPath, preserveFiles); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourcePath, destPath, preserveFiles); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file with preservation logic
func copyFile(source, dest string, preserveFiles []string) error {
	fileName := filepath.Base(dest)

	// Check if this file should be preserved
	for _, preserveFile := range preserveFiles {
		if fileName == preserveFile {
			if _, err := os.Stat(dest); err == nil {
				logger.LogMessage("INFO", fmt.Sprintf("Preserving existing file: %s", dest))
				return nil // File exists and should be preserved
			}
		}
	}

	// Copy the file
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// executeModulesInOrder executes all modules according to priority order
func executeModulesInOrder(config *Config) error {
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

	return nil
}

// executeModuleCategory executes all modules in a category
func executeModuleCategory(category string, modules map[string]Module) error {
	if len(modules) == 0 {
		logger.LogMessage("INFO", fmt.Sprintf("No %s modules to execute", category))
		if program != nil {
			sendFormattedLog("📋", "📦", strings.Title(category), "No modules")
		}
		return nil
	}

	priority := ModuleOrder[category]
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
		if err := installPackages(module.Packages); err != nil {
			logger.LogMessage("ERROR", fmt.Sprintf("Package installation failed for %s: %v", fullName, err))
			return fmt.Errorf("package installation failed for %s: %w", fullName, err)
		}

		// Handle Git configuration for identity module
		if category == "core" && name == "identity" {
			if err := handleGitConfiguration(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Git configuration had issues: %v", err))
				if program != nil {
					sendFormattedLog("⚠️", "📦", "Git Setup", "Issues: "+err.Error())
				}
			}
		}

		// Copy configs
		if err := copyConfigs(module.Configs); err != nil {
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

// handleGitConfiguration applies Git configuration with beautiful styling

func handleGitConfiguration() error {
	logger.LogMessage("INFO", "🔧 Applying Git configuration...")

	if program != nil {
		sendFormattedLog("🔄", "📦", "Git Setup", "Checking credentials")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Check for existing git credentials in git config first
	existingName, _ := runGitConfigGet("user.name")
	existingEmail, _ := runGitConfigGet("user.email")

	var userName, userEmail string

	// If we have existing git credentials, ask to use them
	if strings.TrimSpace(existingName) != "" || strings.TrimSpace(existingEmail) != "" {
		if program != nil {
			sendFormattedLog("🎉", "📦", "Git Found", "Found existing git credentials")
			sendFormattedLog("📋", "📦", "Name", existingName)
			sendFormattedLog("📋", "📦", "Email", existingEmail)
			program.Send(tui.InputRequestMsg{Mode: "git-confirm", Prompt: ""})
		}

		// Wait for confirmation
		<-gitInputDone

		if gitConfirmUse {
			userName = existingName
			userEmail = existingEmail
			if program != nil {
				sendFormattedLog("✅", "📦", "Git Setup", "Using existing credentials")
			}
		} else {
			// User said no, prompt for new credentials
			if program != nil {
				sendFormattedLog("💬", "📦", "Git Setup", "Setting up new credentials")
				program.Send(tui.InputRequestMsg{Mode: "git-username", Prompt: "Git Username: "})
			}

			// Wait for new credentials
			<-gitInputDone

			userName = gitUsername
			userEmail = gitEmail
		}
	} else {
		// No existing credentials, prompt for new ones
		if program != nil {
			sendFormattedLog("💬", "📦", "Git Setup", "No credentials found, setting up")
			program.Send(tui.InputRequestMsg{Mode: "git-username", Prompt: "Git Username: "})
		}

		// Wait for credentials to be entered
		<-gitInputDone

		userName = gitUsername
		userEmail = gitEmail
	}

	// Save credentials to user.env
	if err := saveGitCredentials(homeDir, userName, userEmail); err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to save git credentials: %v", err))
	}

	// Apply Git user configuration
	if strings.TrimSpace(userName) != "" {
		if err := runGitConfig("user.name", userName); err != nil {
			return fmt.Errorf("setting git user.name: %w", err)
		}
		if program != nil {
			sendFormattedLog("✅", "📦", "Git Identity", "User name set to: "+userName)
		}
		logger.LogMessage("SUCCESS", fmt.Sprintf("Git user.name set to: %s", userName))
	}

	if strings.TrimSpace(userEmail) != "" {
		if err := runGitConfig("user.email", userEmail); err != nil {
			return fmt.Errorf("setting git user.email: %w", err)
		}
		if program != nil {
			sendFormattedLog("✅", "📦", "Git Identity", "User email set to: "+userEmail)
		}
		logger.LogMessage("SUCCESS", fmt.Sprintf("Git user.email set to: %s", userEmail))
	}

	// Apply Git aliases and defaults
	gitConfigs := map[string]string{
		"alias.co":           "checkout",
		"alias.br":           "branch",
		"alias.ci":           "commit",
		"alias.st":           "status",
		"pull.rebase":        "true",
		"init.defaultBranch": "master",
	}

	if program != nil {
		sendFormattedLog("🔄", "📦", "Git Aliases", "Setting aliases")
	}

	for key, value := range gitConfigs {
		if err := runGitConfig(key, value); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Failed to set %s: %v", key, err))
		}
	}

	if program != nil {
		sendFormattedLog("✅", "📦", "Git Setup", "Complete")
	}
	logger.LogMessage("SUCCESS", "Git configuration applied")

	return nil
}

// saveGitCredentials saves git credentials to user.env file
func saveGitCredentials(homeDir, userName, userEmail string) error {
	configDir := filepath.Join(homeDir, ".config", "archriot")
	envFile := filepath.Join(configDir, "user.env")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Create user.env content
	content := fmt.Sprintf("ARCHRIOT_USER_NAME=\"%s\"\nARCHRIOT_USER_EMAIL=\"%s\"\n", userName, userEmail)

	// Write to file
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing user.env file: %w", err)
	}

	logger.LogMessage("SUCCESS", "Git credentials saved to user.env")
	return nil
}

// runGitConfig sets a git configuration value
func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	return cmd.Run()
}

// runGitConfigGet gets a git configuration value
func runGitConfigGet(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// commandExists checks if a command is available in PATH
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
