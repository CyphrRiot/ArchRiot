package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"archriot-installer/config"
	"archriot-installer/logger"
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
		// Store if original path was intended as directory (ends with /)
		wasDirectory := strings.HasSuffix(path, "/")
		expanded := filepath.Join(homeDir, path[2:])
		// Only restore trailing slash if original had one (indicating directory intent)
		if wasDirectory && !strings.HasSuffix(expanded, "/") {
			expanded += "/"
		}
		return expanded
	}
	return path
}

// copyConfigPattern copies files matching a config pattern with preservation
func copyConfigPattern(sourceDir, homeDir string, configRule config.ConfigRule) error {
	// Parse pattern (e.g., "hypr/*" -> source: config/hypr, dest: ~/.config/hypr)
	pattern := configRule.Pattern
	var sourcePath, destPath string

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
			if strings.HasSuffix(expandedTarget, "/") {
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
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("creating dest directory: %w", err)
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

	if err := os.WriteFile(dest, sourceData, 0644); err != nil {
		return fmt.Errorf("writing dest file: %w", err)
	}

	return nil
}
