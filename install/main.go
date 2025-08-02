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
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
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

// Global logging infrastructure
var (
	logFile      *os.File
	errorLogFile *os.File
	logPath      string
	errorLogPath string
)

// Tokyo Night inspired Cypher Riot theme
var (
	primaryColor = lipgloss.Color("#7aa2f7") // Tokyo Night blue
	accentColor  = lipgloss.Color("#bb9af7") // Tokyo Night purple
	successColor = lipgloss.Color("#9ece6a") // Tokyo Night green
	warningColor = lipgloss.Color("#e0af68") // Tokyo Night yellow
	errorColor   = lipgloss.Color("#f7768e") // Tokyo Night red
	bgColor      = lipgloss.Color("#1a1b26") // Tokyo Night background
	fgColor      = lipgloss.Color("#c0caf5") // Tokyo Night foreground
	dimColor     = lipgloss.Color("#565f89") // Tokyo Night comment/dim - secondary text
)

// ASCII art for ArchRiot installer
const ArchRiotASCII = ` █████╗ ██████╗  ██████╗██╗  ██╗██████╗ ██╗ ██████╗ ████████╗
██╔══██╗██╔══██╗██╔════╝██║  ██║██╔══██╗██║██╔═══██╗╚══██╔══╝
███████║██████╔╝██║     ███████║██████╔╝██║██║   ██║   ██║
██╔══██║██╔══██╗██║     ██╔══██║██╔══██╗██║██║   ██║   ██║
██║  ██║██║  ██║╚██████╗██║  ██║██║  ██║██║╚██████╔╝   ██║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝ ╚═════╝    ╚═╝`

// TUI Model for installation progress
type InstallModel struct {
	progress    float64
	message     string
	logs        []string
	maxLogs     int
	width       int
	height      int
	done        bool
	operation   string
	currentStep string
	inputMode   string // "git-username", "git-email", "reboot", ""
	inputValue  string // current typed input
	inputPrompt string // what we're asking for
	showRebootButtons bool // show YES/NO buttons in scroll window
	selectedButton    int  // 0 = YES, 1 = NO
}

// NewInstallModel creates a new installation model
func NewInstallModel() *InstallModel {
	return &InstallModel{
		logs:        make([]string, 0),
		maxLogs:     12,
		width:       80,
		height:      24,
		operation:   "ArchRiot Installation",
		currentStep: "Initializing...",
		inputMode:   "",
		inputValue:  "",
		inputPrompt: "",
		showRebootButtons: false,
		selectedButton:    1, // Default to NO
	}
}

// Init implements tea.Model
func (m *InstallModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.showRebootButtons {
			// Handle button selection
			switch msg.Type {
			case tea.KeyLeft:
				m.selectedButton = 0 // YES
				return m, nil
			case tea.KeyRight:
				m.selectedButton = 1 // NO
				return m, nil
			case tea.KeyEnter:
				// Just quit directly - don't send messages
				return m, tea.Quit
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			}
			return m, nil
		}

		if m.inputMode != "" {
			// Handle input mode
			switch msg.Type {
			case tea.KeyEnter:
				logMessage("DEBUG", fmt.Sprintf("ENTER pressed with input: '%s'", m.inputValue))
				return m.handleInputSubmit()
			case tea.KeyBackspace:
				if len(m.inputValue) > 0 {
					m.inputValue = m.inputValue[:len(m.inputValue)-1]
				}
				return m, nil
			case tea.KeyRunes:
				m.inputValue += string(msg.Runes)
				// Test: write ANY keypress to log to see if input is working
				logMessage("DEBUG", fmt.Sprintf("Key received: '%s', input now: '%s'", string(msg.Runes), m.inputValue))
				return m, nil
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			}
			return m, nil
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// If installation is complete and no input mode, any key press exits
		if m.done && m.inputMode == "" {
			return m, tea.Quit
		}
	case LogMsg:
		m.addLog(string(msg))
		return m, nil
	case ProgressMsg:
		m.setProgress(float64(msg))
		return m, nil
	case StepMsg:
		m.setCurrentStep(string(msg))
		return m, nil
	case DoneMsg:
		m.done = true
		m.showRebootButtons = true
		m.addLog("")
		m.addLog("✅ 🎉 Installation     Complete!")
		m.addLog("")
		return m, nil
	case InputRequestMsg:
		m.setInputMode(msg.Mode, msg.Prompt)
		return m, nil
	case GitUsernameMsg:
		// Git username received, handled by installation process
		return m, nil
	case GitEmailMsg:
		// Git email received, handled by installation process
		return m, nil
	case RebootMsg:
		// Reboot decision received, handled by installation process
		return m, nil
	}
	return m, nil
}

// View implements tea.Model - exact structure from Migrate
func (m *InstallModel) View() string {
	var s strings.Builder

	// Header - ASCII + title + version (like Migrate) with spacing
	s.WriteString("\n\n") // Blank lines before ASCII logo
	asciiStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	ascii := asciiStyle.Render(ArchRiotASCII)
	s.WriteString(ascii + "\n")

	titleStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	title := titleStyle.Render(fmt.Sprintf("ArchRiot Installer v%s", VERSION))
	s.WriteString(title + "\n")

	versionStyle := lipgloss.NewStyle().Foreground(dimColor)
	subtitle := versionStyle.Render("Tokyo Night Inspired • Cypher Riot Themed")
	s.WriteString(subtitle + "\n\n")

	// Operation title
	operationStyle := lipgloss.NewStyle().Foreground(successColor).Bold(true)
	s.WriteString(operationStyle.Render("📦 "+m.operation) + "\n")

	// Info section - operation details
	infoStyle := lipgloss.NewStyle().Foreground(fgColor)
	logStyle := lipgloss.NewStyle().Foreground(dimColor)

	s.WriteString(infoStyle.Render("📋 Current Step:   "+m.currentStep) + "\n")
	s.WriteString(logStyle.Render("📝 Log File:       "+logPath) + "\n")

	// Progress bar
	s.WriteString(m.renderProgressBar() + "\n\n")

	// Scroll window - bordered content area
	s.WriteString(m.renderScrollWindow())

	// Buttons below scroll window if shown
	if m.showRebootButtons {
		var yesButton, noButton string

		if m.selectedButton == 0 {
			yesButton = "[►YES◄]"
			noButton = "[ NO ]"
		} else {
			yesButton = "[ YES ]"
			noButton = "[►NO◄]"
		}

		s.WriteString(fmt.Sprintf("\n\nReboot now?  %s  %s  (← → to select, Enter to confirm)", yesButton, noButton))
	} else if m.inputMode != "" {
		s.WriteString("\n\n" + m.inputPrompt + m.inputValue + "_")
	} else {
		s.WriteString("\n\nPress 'q' to quit or 'ctrl+c' to exit")
	}

	return s.String()
}

// renderProgressBar creates a progress bar with percentage
func (m *InstallModel) renderProgressBar() string {
	progressStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)

	width := 50
	filled := int(m.progress * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := fmt.Sprintf("%.1f%%", m.progress*100)

	return progressStyle.Render(fmt.Sprintf("Progress: [%s] %s", bar, percentage))
}

// renderScrollWindow creates the bordered scroll window
func (m *InstallModel) renderScrollWindow() string {
	// Calculate responsive dimensions based on terminal size
	contentWidth := m.width - 4 // Account for borders and padding
	if contentWidth < 40 {
		contentWidth = 40 // Minimum width
	}

	// Calculate available height for scroll window
	// Account for: ASCII (7) + Title/subtitle (4) + Operation (2) + Progress (2) + Footer (3) + Buffer (2)
	usedHeight := 20
	availableHeight := m.height - usedHeight
	if availableHeight < 5 {
		availableHeight = 5 // Minimum scroll window height
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(contentWidth).
		Height(availableHeight)

	content := strings.Builder{}
	content.WriteString("Installation Log\n")
	separatorWidth := contentWidth - 2 // Account for padding
	if separatorWidth < 10 {
		separatorWidth = 10
	}
	content.WriteString(strings.Repeat("─", separatorWidth) + "\n")

	// Calculate how many log lines we can display
	maxLogLines := availableHeight - 2 // Account for header and separator
	if maxLogLines < 1 {
		maxLogLines = 1
	}

	// Show recent logs - limit to calculated space
	start := 0
	if len(m.logs) > maxLogLines {
		start = len(m.logs) - maxLogLines
	}

	actualLogCount := len(m.logs) - start
	if actualLogCount > maxLogLines {
		actualLogCount = maxLogLines
	}

	for i := start; i < start + actualLogCount; i++ {
		line := m.logs[i]
		maxLineWidth := contentWidth - 4 // Account for padding
		if maxLineWidth < 10 {
			maxLineWidth = 10
		}
		if len(line) > maxLineWidth {
			line = line[:maxLineWidth-3] + "..."
		}
		content.WriteString(line + "\n")
	}

	// Buttons are now rendered outside scroll window

	// Fill remaining lines in log area only
	for i := actualLogCount; i < maxLogLines; i++ {
		content.WriteString("\n")
	}

	return boxStyle.Render(content.String())
}

// renderComplete - REMOVED, completion shows at bottom of normal view

// addLog adds a new log entry
func (m *InstallModel) addLog(message string) {
	m.logs = append(m.logs, message)
}

// setProgress updates the progress value
func (m *InstallModel) setProgress(progress float64) {
	m.progress = progress
}

// setCurrentStep updates the current step
func (m *InstallModel) setCurrentStep(step string) {
	m.currentStep = step
}

// setInputMode sets the input mode and prompt
func (m *InstallModel) setInputMode(mode, prompt string) {
	m.inputMode = mode
	m.inputPrompt = prompt
	m.inputValue = ""
}

// handleInputSubmit processes submitted input
func (m *InstallModel) handleInputSubmit() (tea.Model, tea.Cmd) {
	switch m.inputMode {
	case "git-username":
		// Store git username and ask for email
		program.Send(GitUsernameMsg(m.inputValue))
		m.setInputMode("git-email", "Git Email: ")
		return m, nil
	case "git-email":
		// Store git email and continue
		program.Send(GitEmailMsg(m.inputValue))
		m.inputMode = ""
		m.inputPrompt = ""
		m.inputValue = ""
		return m, nil
	case "reboot":
		value := strings.ToLower(strings.TrimSpace(m.inputValue))

		if value == "y" || value == "yes" {
			program.Send(RebootMsg(true))
		} else {
			program.Send(RebootMsg(false))
		}
		// Force exit with os.Exit instead of tea.Quit
		os.Exit(0)
	}
	return m, nil
}

// LogMsg represents a log message
type LogMsg string

// ProgressMsg represents progress update
type ProgressMsg float64

// StepMsg represents a step update
type StepMsg string

// DoneMsg indicates completion
type DoneMsg struct{}

// GitUsernameMsg carries git username input
type GitUsernameMsg string

// GitEmailMsg carries git email input
type GitEmailMsg string

// RebootMsg carries reboot decision
type RebootMsg bool

// InputRequestMsg requests user input
type InputRequestMsg struct {
	Mode   string
	Prompt string
}

// Global model instance
var model *InstallModel
var program *tea.Program

// Styled components
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Margin(1, 0)

	welcomeStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accentColor).
			Padding(1, 3).
			Margin(1, 0).
			Width(60).
			Align(lipgloss.Center)

	infoStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 1)

	progressStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	asciiStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Align(lipgloss.Center)
)

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
	if err := initLogging(); err != nil {
		log.Fatalf("❌ Failed to initialize logging: %v", err)
	}
	defer closeLogging()

	// Initialize TUI model
	model = NewInstallModel()
	program = tea.NewProgram(model)

	// Start installation in background
	go func() {
		// Small delay to let TUI initialize
		time.Sleep(100 * time.Millisecond)

		// Run installation in goroutine
		runInstallation()

		// Signal completion
		program.Send(DoneMsg{})
	}()

	// Run TUI in main thread
	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}

func runInstallation() {
	// Send log messages to TUI
	sendLog := func(msg string) {
		program.Send(LogMsg(msg))
	}

	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		program.Send(ProgressMsg(progress))
	}

	// Send step updates to TUI
	sendStep := func(step string) {
		program.Send(StepMsg(step))
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
		sendLog(fmt.Sprintf("❌ 🚀 Module Exec      Failed: %v", err))
		return
	}

	sendStep("Installation complete!")
	sendProgress(1.0)
	sendLog("✅ 🎉 Installation     Complete!")
	sendLog("✅ 🚀 Module Exec      All modules done")
	sendLog(fmt.Sprintf("📝 📋 Log File        Available at: %s", logPath))
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
			program.Send(LogMsg("ℹ️  📋 Packages        None to install"))
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
			program.Send(LogMsg(fmt.Sprintf("❌ 📦 Package Error    Failed: %s", outputStr)))
		}
		return fmt.Errorf("batch installation failed: %w", err)
	}

	return nil
}

// syncPackageDatabases ensures yay is installed then synchronizes pacman and yay package databases
func syncPackageDatabases() error {
	logMessage("INFO", "🔄 Syncing package databases...")
	if program != nil {
		program.Send(LogMsg("🔄 🗄️  Database Sync   Syncing databases"))
	}

	start := time.Now()

	// Install yay if not present
	if !commandExists("yay") {
		logMessage("INFO", "Installing yay AUR helper...")
		// Yay installation logged to file only

		cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", "yay")
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Limit output to prevent log spam
			outputStr := string(output)
			if len(outputStr) > 200 {
				outputStr = outputStr[:200] + "... (truncated)"
			}
			logMessage("CRITICAL", fmt.Sprintf("Failed to install yay: %s", outputStr))
			return fmt.Errorf("yay installation failed: %w", err)
		}

		logMessage("SUCCESS", "yay AUR helper installed")
		// Yay installation success logged to file only
	} else {
		logMessage("SUCCESS", "yay AUR helper already available")
		// Yay already available logged to file only
	}

	// Sync pacman database - fail gracefully if it doesn't work
	if err := syncPackmanDatabase(); err != nil {
		logMessage("CRITICAL", fmt.Sprintf("Pacman sync failed: %v", err))
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
		logMessage("WARNING", fmt.Sprintf("Failed to sync yay database: %s", outputStr))
		if program != nil {
			program.Send(LogMsg("⚠️  🗄️  Database Sync   Yay sync failed, continuing"))
		}
	} else {
		logMessage("SUCCESS", "Yay database synced")
		// Yay sync success logged to file only
	}

	duration := time.Since(start)
	logMessage("SUCCESS", fmt.Sprintf("Package databases synced in %v", duration))
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

	logMessage("SUCCESS", "Pacman database synced")
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

// initLogging initializes log files and creates necessary directories
func initLogging() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(homeDir, ".cache", "archriot")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	// Set log file paths
	logPath = filepath.Join(logDir, "install.log")
	errorLogPath = filepath.Join(logDir, "install-errors.log")

	// Create/truncate log files
	logFile, err = os.Create(logPath)
	if err != nil {
		return fmt.Errorf("creating log file: %w", err)
	}

	errorLogFile, err = os.Create(errorLogPath)
	if err != nil {
		return fmt.Errorf("creating error log file: %w", err)
	}

	// Write headers
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "=== ArchRiot Go Installer v%s - %s ===\n", VERSION, timestamp)
	fmt.Fprintf(errorLogFile, "=== ArchRiot Go Installer Errors v%s - %s ===\n", VERSION, timestamp)

	return nil
}

// closeLogging closes log files
func closeLogging() {
	if logFile != nil {
		logFile.Close()
	}
	if errorLogFile != nil {
		errorLogFile.Close()
	}
}

// logMessage writes structured log messages to files and optionally console
func logMessage(level string, message string) {
	timestamp := time.Now().Format("15:04:05")
	formattedMessage := fmt.Sprintf("[%s] %s %s", timestamp, getLevelIcon(level), message)

	// Always write to main log file
	if logFile != nil {
		fmt.Fprintln(logFile, formattedMessage)
		logFile.Sync() // Ensure immediate write
	}

	// Write errors to error log file too
	if level == "ERROR" || level == "CRITICAL" {
		if errorLogFile != nil {
			fmt.Fprintln(errorLogFile, formattedMessage)
			errorLogFile.Sync()
		}
	}
}

// getLevelIcon returns appropriate icon for log level
func getLevelIcon(level string) string {
	switch level {
	case "INFO":
		return "ℹ️ "
	case "SUCCESS":
		return "✅"
	case "WARNING":
		return "⚠️ "
	case "ERROR":
		return "❌"
	case "CRITICAL":
		return "🚨"
	default:
		return "📝"
	}
}

// copyConfigs copies configuration files with preservation logic
func copyConfigs(configs []ConfigRule) error {
	if len(configs) == 0 {
		if program != nil {
			program.Send(LogMsg("ℹ️  📁 Config Copy      None to copy"))
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
		program.Send(LogMsg(fmt.Sprintf("🔄 📁 Config Copy      From: %s", configSourceDir)))
	}

	for _, configRule := range configs {
		// logMessage("INFO", fmt.Sprintf("Processing config pattern: %s", configRule.Pattern))

		if err := copyConfigPattern(configSourceDir, homeDir, configRule); err != nil {
			// logMessage("WARNING", fmt.Sprintf("Failed to copy config %s: %v", configRule.Pattern, err))
			if program != nil {
				program.Send(LogMsg(fmt.Sprintf("❌ 📄 %-15s Failed: %v", configRule.Pattern, err)))
			}
		} else {
			if program != nil {
				program.Send(LogMsg(fmt.Sprintf("✅ 📄 %-15s Copied successfully", configRule.Pattern)))
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
				logMessage("INFO", fmt.Sprintf("Preserving existing file: %s", dest))
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
	logMessage("INFO", "Starting module execution in priority order")
	if program != nil {
		program.Send(LogMsg("🔄 🚀 Module Exec      Starting modules"))
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
		logMessage("INFO", fmt.Sprintf("No %s modules to execute", category))
		if program != nil {
			program.Send(LogMsg(fmt.Sprintf("ℹ️  ⏭️  %-15s No modules", strings.Title(category))))
		}
		return nil
	}

	priority := ModuleOrder[category]
	logMessage("INFO", fmt.Sprintf("Executing %s modules (priority %d)", category, priority))
	if program != nil {
		program.Send(LogMsg(fmt.Sprintf("🔄 🔧 %-15s Starting %s modules", strings.Title(category), category)))
	}

	for name, module := range modules {
		fullName := fmt.Sprintf("%s.%s", category, name)
		logMessage("INFO", fmt.Sprintf("Starting module: %s - %s", fullName, module.Description))
		if program != nil {
			program.Send(LogMsg(fmt.Sprintf("🔄 📦 %-15s %s", fullName, module.Description)))
		}

		// Install packages
		if err := installPackages(module.Packages); err != nil {
			logMessage("ERROR", fmt.Sprintf("Package installation failed for %s: %v", fullName, err))
			return fmt.Errorf("package installation failed for %s: %w", fullName, err)
		}

		// Handle Git configuration for identity module
		if category == "core" && name == "identity" {
			if err := handleGitConfiguration(); err != nil {
				logMessage("WARNING", fmt.Sprintf("Git configuration had issues: %v", err))
				if program != nil {
					program.Send(LogMsg(fmt.Sprintf("⚠️  🔧 Git Setup        Issues: %v", err)))
				}
			}
		}

		// Copy configs
		if err := copyConfigs(module.Configs); err != nil {
			logMessage("WARNING", fmt.Sprintf("Config copying had issues for %s: %v", fullName, err))
			if program != nil {
				program.Send(LogMsg(fmt.Sprintf("⚠️  📁 %-15s Config issues: %v", fullName, err)))
			}
		}

		logMessage("SUCCESS", fmt.Sprintf("Module %s completed", fullName))
		if program != nil {
			program.Send(LogMsg(fmt.Sprintf("✅ ✓ %-15s Complete", fullName)))
		}
	}

	logMessage("SUCCESS", fmt.Sprintf("All %s modules completed", category))
	if program != nil {
		program.Send(LogMsg(fmt.Sprintf("✅ 🎉 %-15s All done", strings.Title(category))))
	}
	return nil
}

// handleGitConfiguration applies Git configuration with beautiful styling

func handleGitConfiguration() error {
	logMessage("INFO", "🔧 Applying Git configuration...")

	if program != nil {
		program.Send(LogMsg("🔄 🔧 Git Setup        Configuring identity"))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Load user environment if it exists
	envFile := filepath.Join(homeDir, ".config", "archriot", "user.env")
	var userName, userEmail string

	if _, err := os.Stat(envFile); err == nil {
		logMessage("INFO", "Loading user environment from user.env")
		if data, err := os.ReadFile(envFile); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ARCHRIOT_USER_NAME=") {
					userName = strings.Trim(strings.TrimPrefix(line, "ARCHRIOT_USER_NAME="), "\"'")
				}
				if strings.HasPrefix(line, "ARCHRIOT_USER_EMAIL=") {
					userEmail = strings.Trim(strings.TrimPrefix(line, "ARCHRIOT_USER_EMAIL="), "\"'")
				}
			}
		}
	}

	// Apply Git user configuration
	if strings.TrimSpace(userName) != "" {
		if err := runGitConfig("user.name", userName); err != nil {
			return fmt.Errorf("setting git user.name: %w", err)
		}
		if program != nil {
			program.Send(LogMsg(fmt.Sprintf("✅ 👤 Git Identity     User name set to: %s", userName)))
		}
		logMessage("SUCCESS", fmt.Sprintf("Git user.name set to: %s", userName))
	}

	if strings.TrimSpace(userEmail) != "" {
		if err := runGitConfig("user.email", userEmail); err != nil {
			return fmt.Errorf("setting git user.email: %w", err)
		}
		if program != nil {
			program.Send(LogMsg(fmt.Sprintf("✅ 📧 Git Identity     User email set to: %s", userEmail)))
		}
		logMessage("SUCCESS", fmt.Sprintf("Git user.email set to: %s", userEmail))
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
		program.Send(LogMsg("🔄 ⚙️  Git Aliases      Setting aliases"))
	}

	for key, value := range gitConfigs {
		if err := runGitConfig(key, value); err != nil {
			logMessage("WARNING", fmt.Sprintf("Failed to set %s: %v", key, err))
		}
	}

	if program != nil {
		program.Send(LogMsg("✅ 🔧 Git Setup        Complete"))
	}
	logMessage("SUCCESS", "Git configuration applied")

	return nil
}

// runGitConfig executes a git config command
func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	return cmd.Run()
}

// commandExists checks if a command is available in PATH
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
