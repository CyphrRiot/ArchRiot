package logger

import (
	"archriot-installer/tui"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	logFile          *os.File
	errorLogFile     *os.File
	logPath          string
	errorLogPath     string
	program          *tea.Program
	emojiSupport     bool
	forceAsciiMode   bool
	tuiEmojiCallback func(bool)
)

// Type constants mapping to emojis
var typeEmojis = map[string]string{
	"Package":  "📦",
	"Git":      "🔧",
	"Database": "🗄️",
	"Module":   "🏗️",
	"File":     "📁",
	"System":   "💫",
}

// Single-character ASCII fallbacks for status emojis
var statusFallbacks = map[string]string{
	"Progress": ".",
	"Complete": "!",
	"Success":  "+",
	"Warning":  "?",
	"Error":    "X",
	"Info":     "i",
}

// Single-character ASCII fallbacks for type emojis
var typeFallbacks = map[string]string{
	"Package":  "*",
	"Git":      "+",
	"Database": "#",
	"Module":   "=",
	"File":     "-",
	"System":   "~",
}

// detectTerminalEmojiSupport checks if terminal supports emoji display
func detectTerminalEmojiSupport() bool {
	if forceAsciiMode {
		return false
	}

	term := os.Getenv("TERM")
	lang := os.Getenv("LANG")

	// Basic terminals don't support emojis
	if term == "linux" || term == "console" || term == "dumb" {
		return false
	}

	// No UTF-8 = no emojis
	if !strings.Contains(lang, "UTF-8") {
		return false
	}

	return true
}

// SetForceAsciiMode allows forcing ASCII mode (for --ascii-only flag)
func SetForceAsciiMode(force bool) {
	forceAsciiMode = force
}

// SetTuiEmojiCallback sets callback function for TUI emoji mode updates
func SetTuiEmojiCallback(callback func(bool)) {
	tuiEmojiCallback = callback
}

// getStatusEmoji returns emoji or ASCII based on terminal support
func getStatusEmoji(status string) string {
	if !emojiSupport {
		if fallback, exists := statusFallbacks[status]; exists {
			return fallback
		}
		return "i" // default fallback
	}

	switch status {
	case "Progress":
		return "⏳"
	case "Complete":
		return "🎉"
	case "Success":
		return "✅"
	case "Warning":
		return "⚠️"
	case "Error":
		return "❌"
	case "Info":
		return "📋"
	default:
		return "📋"
	}
}

// getTypeEmoji returns type emoji or ASCII based on terminal support
func getTypeEmoji(logType string) string {
	if !emojiSupport {
		if fallback, exists := typeFallbacks[logType]; exists {
			return fallback
		}
		return "*" // default fallback
	}

	if emoji, exists := typeEmojis[logType]; exists {
		return emoji
	}
	return logType // fallback to original string if type not found
}

// getCurrentStepEmoji returns the current step indicator
func getCurrentStepEmoji() string {
	if !emojiSupport {
		return ">"
	}
	return "🎯"
}

// getLogFileEmoji returns the log file indicator
func getLogFileEmoji() string {
	if !emojiSupport {
		return "-"
	}
	return "📁"
}

// GetCurrentStepEmoji returns the current step indicator (public for TUI)
func GetCurrentStepEmoji() string {
	return getCurrentStepEmoji()
}

// GetLogFileEmoji returns the log file indicator (public for TUI)
func GetLogFileEmoji() string {
	return getLogFileEmoji()
}

// SetProgram sets the TUI program instance for logging
func SetProgram(p *tea.Program) {
	program = p
}

// InitLogging initializes the logging system with proper file paths
func InitLogging() error {
	// Detect terminal emoji support
	emojiSupport = detectTerminalEmojiSupport()

	// Notify TUI of emoji mode if callback is set
	if tuiEmojiCallback != nil {
		tuiEmojiCallback(emojiSupport)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	logPath = filepath.Join(homeDir, ".cache", "archriot", "install.log")
	errorLogPath = filepath.Join(homeDir, ".cache", "archriot", "install-errors.log")

	// Create directories
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	// Open log files - truncate at start, then append during session
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	errorLogFile, err = os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening error log file: %w", err)
	}

	return nil
}

// CloseLogging closes the logging files
func CloseLogging() {
	if logFile != nil {
		logFile.Close()
	}
	if errorLogFile != nil {
		errorLogFile.Close()
	}
}

// LogMessage logs a message with the specified level (file only)
func LogMessage(level, message string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)

	// Write to appropriate files based on level
	switch level {
	case "ERROR", "CRITICAL":
		if logFile != nil {
			logFile.WriteString(logEntry)
		}
		if errorLogFile != nil {
			errorLogFile.WriteString(logEntry)
		}
	default:
		if logFile != nil {
			logFile.WriteString(logEntry)
		}
	}
}

// Log logs to both file and TUI using semantic constants
func Log(status, logType, name, description string) {
	// Log to file with traditional format
	logLevel := mapStatusToLevel(status)
	fileMessage := fmt.Sprintf("%s.%s - %s: %s", logType, name, status, description)
	LogMessage(logLevel, fileMessage)

	// Send to TUI with formatted display
	if program != nil {
		// Get emojis/ASCII automatically through logger functions
		statusEmoji := getStatusEmoji(status)
		typeEmoji := getTypeEmoji(logType)

		// Truncate name if longer than 20 characters
		if len(name) > 20 {
			name = name[:17] + "..."
		}
		program.Send(tui.LogMsg(fmt.Sprintf("%s %s %-20s %s", statusEmoji, typeEmoji, name, description)))
	}
}

// GetLogPath returns the current log file path
func GetLogPath() string {
	return logPath
}

// mapStatusToLevel converts semantic status to log level
func mapStatusToLevel(status string) string {
	switch status {
	case "Progress", "Info":
		return "INFO"
	case "Success", "Complete":
		return "SUCCESS"
	case "Warning":
		return "WARNING"
	case "Error":
		return "ERROR"
	default:
		return "INFO"
	}
}
