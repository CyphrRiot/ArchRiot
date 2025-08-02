package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile      *os.File
	errorLogFile *os.File
	logPath      string
	errorLogPath string
)

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

// LogMessage logs a message with the specified level
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

// GetLogPath returns the current log file path
func GetLogPath() string {
	return logPath
}

// GetLevelIcon returns the appropriate icon for a log level
func GetLevelIcon(level string) string {
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
