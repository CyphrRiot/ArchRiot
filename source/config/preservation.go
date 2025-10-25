package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"archriot-installer/logger"
)

const (
	// CacheDir is the directory where backups are stored
	CacheDir = ".cache/archriot"
)

// PreservationConfig represents the structure of preserve.yaml
type PreservationConfig struct {
	Version                  string        `yaml:"version"`
	UserCustomizableSettings []UserSetting `yaml:"user_customizable_settings"`
	Backup                   BackupConfig  `yaml:"backup"`
	ManagedFiles             []ManagedFile `yaml:"managed_files"`
	Prompts                  PromptsConfig `yaml:"prompts"`
}

// UserSetting represents a single preservable setting
type UserSetting struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Pattern     string `yaml:"pattern"`
}

// BackupConfig represents backup configuration
type BackupConfig struct {
	CacheDir       string `yaml:"cache_dir"`
	FilenameFormat string `yaml:"filename_format"`
	AlwaysBackup   bool   `yaml:"always_backup"`
}

// PromptsConfig represents user prompt configuration
type ManagedFile struct {
	Path   string `yaml:"path"`
	Backup string `yaml:"backup"` // e.g., "side_by_side"
	When   string `yaml:"when"`   // e.g., "on_overwrite_diff"
	Keep   int    `yaml:"keep"`   // how many old copies to keep (currently informational)
}

type PromptsConfig struct {
	RestorePrompt string            `yaml:"restore_prompt"`
	DefaultChoice int               `yaml:"default_choice"`
	Messages      map[string]string `yaml:"messages"`
}

// Global preservation configuration
var preservationConfig *PreservationConfig

// PreservationResult represents the outcome of config preservation
type PreservationResult struct {
	BackupPath      string
	PreservedValues map[string]string
	Message         string
}

// ExtractUserSettings extracts user-customized settings from existing config
func ExtractUserSettings(existingConfigPath string) (*PreservationResult, error) {
	result := &PreservationResult{
		BackupPath:      "",
		PreservedValues: make(map[string]string),
		Message:         "",
	}

	// Check if existing config exists
	if _, err := os.Stat(existingConfigPath); os.IsNotExist(err) {
		result.Message = "No existing config found - fresh install"
		logger.LogMessage("INFO", result.Message)
		return result, nil
	}

	// Read existing config
	existingContent, err := os.ReadFile(existingConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read existing config: %w", err)
	}

	// Early comparison - check if files are identical
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}
	newConfigPath := filepath.Join(homeDir, ".local", "share", "archriot", "config", "hypr", "hyprland.conf")

	newContent, err := os.ReadFile(newConfigPath)
	if err == nil {
		// If files are byte-for-byte identical, skip preservation entirely
		if string(existingContent) == string(newContent) {
			result.Message = "Configs are identical - no preservation needed"
			logger.LogMessage("INFO", result.Message)
			return result, nil
		}
	}

	// Extract user settings
	extractedSettings := extractSettings(string(existingContent))

	// Check if any extracted settings are actually different from new config
	meaningfulSettings := make(map[string]string)
	if len(extractedSettings) > 0 {
		// Read new config to compare
		newContent, err := os.ReadFile(newConfigPath)
		if err == nil {
			newSettings := extractSettings(string(newContent))

			// Only preserve settings that are actually different
			for key, userValue := range extractedSettings {
				newValue, exists := newSettings[key]
				if !exists || userValue != newValue {
					meaningfulSettings[key] = userValue
				}
			}
		} else {
			// If we can't read new config, preserve all extracted settings
			meaningfulSettings = extractedSettings
		}
	}

	result.PreservedValues = meaningfulSettings

	if len(meaningfulSettings) > 0 {
		result.Message = fmt.Sprintf("Extracted %d user settings for preservation", len(meaningfulSettings))
		logger.LogMessage("INFO", result.Message)
		for key, value := range meaningfulSettings {
			logger.LogMessage("INFO", fmt.Sprintf("Preserving: %s = %s", key, value))
		}
	} else {
		result.Message = "No user customizations found in existing config"
		logger.LogMessage("INFO", result.Message)
	}

	return result, nil
}

// CreateBackup creates a timestamped backup of the existing config
func CreateBackup(configPath, configName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	// Create cache directory
	cacheDir := filepath.Join(homeDir, CacheDir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("creating cache directory: %w", err)
	}

	// Generate timestamped backup filename
	timestamp := time.Now().Format("20060102")
	backupName := fmt.Sprintf("%s_%s", timestamp, configName)
	backupPath := filepath.Join(cacheDir, backupName)

	// Copy existing config to backup location
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("reading config for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", fmt.Errorf("writing backup file: %w", err)
	}

	logger.LogMessage("INFO", fmt.Sprintf("Config backed up to: %s", backupPath))
	return backupPath, nil
}

// extractSettings extracts user-customizable settings from config content
func extractSettings(content string) map[string]string {
	settings := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Load config if not already loaded
		if preservationConfig == nil {
			if err := loadPreservationConfig(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Failed to load preservation config: %v", err))
				return settings
			}
		}

		// Check each user-customizable setting
		for _, setting := range preservationConfig.UserCustomizableSettings {
			if strings.Contains(trimmed, setting.Pattern) && strings.Contains(trimmed, "=") {
				// Extract the value after the equals sign
				parts := strings.SplitN(trimmed, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])

					// Match the setting name and check if value is meaningful
					if key == setting.Pattern || strings.HasPrefix(key, setting.Pattern+" ") {
						// Only preserve if value is not empty or whitespace-only
						if value != "" && !isDefaultOrEmpty(value) {
							settings[setting.Name] = value
						}
						break
					}
				}
			}
		}
	}

	return settings
}

// ApplyUserSettings applies preserved user settings to new config content
func ApplyUserSettings(newConfigContent string, userSettings map[string]string) string {
	if len(userSettings) == 0 {
		return newConfigContent
	}

	lines := strings.Split(newConfigContent, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Load config if not already loaded
		if preservationConfig == nil {
			if err := loadPreservationConfig(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Failed to load preservation config: %v", err))
				continue
			}
		}

		// Check if this line contains a setting we want to preserve
		for settingName, userValue := range userSettings {
			// Find the corresponding setting in config
			var pattern string
			for _, setting := range preservationConfig.UserCustomizableSettings {
				if setting.Name == settingName {
					pattern = setting.Pattern
					break
				}
			}

			if pattern != "" && strings.Contains(trimmed, pattern) && strings.Contains(trimmed, "=") {
				parts := strings.SplitN(trimmed, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])

					// Match the setting name and replace with user value
					if key == pattern || strings.HasPrefix(key, pattern+" ") {
						// Preserve original indentation
						indent := ""
						for _, char := range line {
							if char == ' ' || char == '\t' {
								indent += string(char)
							} else {
								break
							}
						}
						lines[i] = indent + key + " = " + userValue
						logger.LogMessage("INFO", fmt.Sprintf("Applied user setting: %s = %s", settingName, userValue))
						break
					}
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

// loadPreservationConfig loads the preservation configuration from preserve.yaml
func loadPreservationConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".local", "share", "archriot", "install", "preserve.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("preserve.yaml not found at %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading preserve.yaml: %w", err)
	}

	// Parse YAML
	config := &PreservationConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("parsing preserve.yaml: %w", err)
	}

	preservationConfig = config
	logger.LogMessage("INFO", fmt.Sprintf("Loaded preservation config v%s with %d settings", config.Version, len(config.UserCustomizableSettings)))

	return nil
}

// GetPreservationConfig returns the loaded preservation configuration
func expandUserPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

// MatchManagedFile returns the managed_files entry matching the absolute target path, if any.
func MatchManagedFile(targetPath string) (*ManagedFile, bool) {
	cfg := GetPreservationConfig()
	if cfg == nil || len(cfg.ManagedFiles) == 0 {
		return nil, false
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		absTarget = filepath.Clean(targetPath)
	}
	absTarget = filepath.Clean(absTarget)

	for i := range cfg.ManagedFiles {
		mp := expandUserPath(cfg.ManagedFiles[i].Path)
		absMP, err := filepath.Abs(mp)
		if err != nil {
			absMP = filepath.Clean(mp)
		}
		absMP = filepath.Clean(absMP)
		if absMP == absTarget {
			return &cfg.ManagedFiles[i], true
		}
	}
	return nil, false
}

// ShouldSideBySideBackupOnOverwriteDiff reports whether we should create/refresh a .old
// side-by-side backup for this path when we are about to overwrite it with different content.
func ShouldSideBySideBackupOnOverwriteDiff(targetPath string) bool {
	mf, ok := MatchManagedFile(targetPath)
	if !ok {
		return false
	}
	if strings.EqualFold(mf.Backup, "side_by_side") {
		// Default policy if unspecified is to back up on overwrite diff
		if mf.When == "" || strings.EqualFold(mf.When, "on_overwrite_diff") {
			return true
		}
	}
	return false
}

func GetPreservationConfig() *PreservationConfig {
	if preservationConfig == nil {
		if err := loadPreservationConfig(); err != nil {
			logger.LogMessage("ERROR", fmt.Sprintf("Failed to load preservation config: %v", err))
			return nil
		}
	}
	return preservationConfig
}

// isDefaultOrEmpty checks if a value is empty, whitespace-only, or a common default
func isDefaultOrEmpty(value string) bool {
	trimmed := strings.TrimSpace(value)

	// Empty or whitespace-only
	if trimmed == "" {
		return true
	}

	// Common default values that shouldn't be preserved
	defaults := []string{
		"us",      // default keyboard layout
		"ghostty", // default terminal
		"brave",   // default browser
		"Thunar",  // default file manager
	}

	for _, defaultVal := range defaults {
		if trimmed == defaultVal {
			return true
		}
	}

	return false
}
