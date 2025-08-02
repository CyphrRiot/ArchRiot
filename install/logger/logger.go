package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"archriot-installer/tui"
)

var (
	logFile      *os.File
	errorLogFile *os.File
	logPath      string
	errorLogPath string
	program      *tea.Program
)

// Type constants mapping to emojis
var typeEmojis = map[string]string{
	"Package":  "📦",
	"Git":      "🔧",
	"Database": "🗄️",
	"Module":   "🏗️",
	"File":     "📁",
	"System":   "🚀",
}

// getStatusEmoji returns emoji based on context automatically
func getStatusEmoji(status string) string {
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

// SetProgram sets the TUI program instance for logging
func SetProgram(p *tea.Program) {
	program = p
}

// InitLogging initializes the logging system with proper file paths
func InitLogging() error {
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

	// Open log files
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	errorLogFile, err = os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
		// Get emojis automatically
		statusEmoji := getStatusEmoji(status)
		typeEmoji := typeEmojis[logType]

		// Fallback to original string if type not found
		if typeEmoji == "" {
			typeEmoji = logType
		}

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
