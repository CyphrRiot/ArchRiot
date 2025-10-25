package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

// CopyConfigs copies configuration files with preservation logic
func CopyConfigs(configs []config.ConfigRule) error {
	if len(configs) == 0 {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	configSourceDir := filepath.Join(homeDir, ".local", "share", "archriot", "config")

	// logMessage("INFO", fmt.Sprintf("Copying configs from: %s", configSourceDir))
	logger.Log("Progress", "File", "Config Copy", "From: "+configSourceDir)

	for _, configRule := range configs {
		logger.LogMessage("INFO", fmt.Sprintf("Processing config pattern: %s", configRule.Pattern))

		if err := copyConfigPattern(configSourceDir, homeDir, configRule); err != nil {
			// logMessage("WARNING", fmt.Sprintf("Failed to copy config %s: %v", configRule.Pattern, err))
			logger.Log("Error", "File", configRule.Pattern, "Failed: "+err.Error())
		} else {
			logger.Log("Success", "File", configRule.Pattern, "Copied successfully")
		}
	}

	return nil
}

// expandTildePath expands ~ to the user's home directory
func expandTildePath(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// copyConfigPattern copies files matching a config pattern with preservation
func copyConfigPattern(sourceDir, homeDir string, configRule config.ConfigRule) error {
	// Parse pattern (e.g., "hypr/*" -> source: config/hypr, dest: ~/.config/hypr)
	pattern := configRule.Pattern
	var sourcePath, destPath string

	// Check if this is hyprland config that needs preservation
	if pattern == "hypr/*" {
		if err := handleHyprlandPreservation(sourceDir, homeDir, program); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Hyprland preservation failed: %v - continuing with normal copy", err))
		}
	}

	if configRule.Target != "" {
		// Custom target specified - expand ~ if present
		expandedTarget := expandTildePath(configRule.Target, homeDir)

		if strings.HasSuffix(pattern, "/*") {
			// Directory pattern with custom target
			dirName := strings.TrimSuffix(pattern, "/*")
			sourcePath = filepath.Join(sourceDir, dirName)
			destPath = expandedTarget
		} else {
			// File pattern with custom target
			sourcePath = filepath.Join(sourceDir, pattern)
			if strings.HasSuffix(configRule.Target, "/") {
				// Target is a directory, append filename
				destPath = filepath.Join(expandedTarget, filepath.Base(pattern))
			} else {
				// Target is a full file path
				destPath = expandedTarget
			}
		}
	} else if strings.HasSuffix(pattern, "/*") {
		// Directory pattern: "hypr/*" -> copy all files from hypr/ to ~/.config/hypr/
		dirName := strings.TrimSuffix(pattern, "/*")
		sourcePath = filepath.Join(sourceDir, dirName)
		destPath = filepath.Join(homeDir, ".config", dirName)
	} else {
		// File pattern: "hypr/hyprland.conf" -> copy specific file
		sourcePath = filepath.Join(sourceDir, pattern)
		destPath = filepath.Join(homeDir, ".config", pattern)
	}

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source not found: %s", sourcePath)
	}

	// Create destination directory
	// For file targets, create parent directory; for directory targets, create the directory itself
	targetInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("checking source: %w", err)
	}

	if targetInfo.IsDir() {
		// Source is a directory, ensure dest directory exists
		if err := os.MkdirAll(destPath, 0755); err != nil {
			return fmt.Errorf("creating dest directory: %w", err)
		}
	} else {
		// Source is a file, ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("creating dest directory: %w", err)
		}
	}

	// Copy files
	return copyFileOrDirectory(sourcePath, destPath, configRule.PreserveIfExists)
}

// copyFileOrDirectory recursively copies files or directories with preservation
func copyFileOrDirectory(source, dest string, preserveFiles []string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("reading source info: %w", err)
	}

	if sourceInfo.IsDir() {
		return copyDirectory(source, dest, preserveFiles)
	}
	return copyFile(source, dest, preserveFiles)
}

// copyDirectory recursively copies a directory
func copyDirectory(source, dest string, preserveFiles []string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
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
	// Check if this file should be preserved
	fileName := filepath.Base(dest)
	for _, preserve := range preserveFiles {
		if fileName == preserve {
			if _, err := os.Stat(dest); err == nil {
				logger.LogMessage("INFO", fmt.Sprintf("Preserving existing file: %s", dest))
				return nil // Skip copying, preserve existing
			}
		}
	}

	sourceData, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	// Managed files backup: honor preserve.yaml managed_files (side-by-side .old) on overwrite-diff
	// Applies to files like ~/.config/hypr/hyprlock.conf and ~/.config/waybar/config
	if config.ShouldSideBySideBackupOnOverwriteDiff(dest) {
		if _, err := os.Stat(dest); err == nil {
			if existing, err2 := os.ReadFile(dest); err2 == nil {
				// Only back up when content differs
				if string(existing) != string(sourceData) {
					if err := os.WriteFile(dest+".old", existing, 0644); err == nil {
						logger.Log("Info", "File", "Backup", "Side-by-side backup created: "+dest+".old")
					} else {
						logger.Log("Warning", "File", "Backup", "Failed to create side-by-side backup: "+dest+".old")
					}
				}
			}
		}
	}

	if err := os.WriteFile(dest, sourceData, 0644); err != nil {
		return fmt.Errorf("writing dest file: %w", err)
	}

	return nil
}

// Global variables for preservation state
var (
	preservationResult     *config.PreservationResult
	preservationNewConfig  string
	preservationSourcePath string
	preservationDone       chan bool
)

// handleHyprlandPreservation handles the preservation logic for hyprland.conf
func handleHyprlandPreservation(sourceDir, homeDir string, prog *tea.Program) error {
	existingConfigPath := filepath.Join(homeDir, ".config", "hypr", "hyprland.conf")
	newConfigPath := filepath.Join(sourceDir, "hypr", "hyprland.conf")

	// Extract user settings from existing config
	result, err := config.ExtractUserSettings(existingConfigPath)
	if err != nil {
		return fmt.Errorf("extracting user settings: %w", err)
	}

	// If no existing config or no user settings found, skip preservation
	if len(result.PreservedValues) == 0 {
		logger.LogMessage("INFO", "No user settings to preserve for hyprland.conf")
		return nil
	}

	// Create backup of existing config
	backupPath, err := config.CreateBackup(existingConfigPath, "hyprland.conf")
	if err != nil {
		return fmt.Errorf("creating backup: %w", err)
	}
	result.BackupPath = backupPath

	// Read new config content
	newContent, err := os.ReadFile(newConfigPath)
	if err != nil {
		return fmt.Errorf("reading new config: %w", err)
	}

	// Store preservation state for callback
	preservationResult = result
	preservationNewConfig = string(newContent)
	preservationSourcePath = newConfigPath
	preservationDone = make(chan bool, 1)

	// Set up preservation callback
	tui.SetPreservationCallback(func(shouldRestore bool) {
		if shouldRestore {
			// Apply user settings to new config
			modifiedContent := config.ApplyUserSettings(preservationNewConfig, preservationResult.PreservedValues)

			// Write the modified config back to source location
			if err := os.WriteFile(preservationSourcePath, []byte(modifiedContent), 0644); err != nil {
				logger.LogMessage("ERROR", fmt.Sprintf("Failed to write modified config: %v", err))
			} else {
				logger.LogMessage("INFO", fmt.Sprintf("Applied %d user settings to new hyprland.conf", len(preservationResult.PreservedValues)))
			}
		} else {
			logger.LogMessage("INFO", "User chose not to restore hyprland modifications")
		}
		preservationDone <- true
	})

	// Send preservation prompt to TUI
	if prog != nil {
		prog.Send(tui.PreservationPromptMsg{})

		// Wait for user response with timeout
		select {
		case <-preservationDone:
			// User responded
		case <-time.After(30 * time.Second):
			// Timeout - default to NO restoration
			logger.LogMessage("WARNING", "Preservation prompt timed out - defaulting to no restoration")
		}
	}

	return nil
}
