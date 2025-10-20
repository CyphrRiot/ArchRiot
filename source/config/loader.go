package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FindConfigFile looks for packages.yaml in common locations
func FindConfigFile() string {
	locations := []string{
		filepath.Join(os.Getenv("HOME"), ".local/share/archriot/install/packages.yaml"),
		filepath.Join("install", "packages.yaml"),
	}

	for _, path := range locations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// LoadConfig reads and parses the YAML configuration
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &cfg, nil
}

// ValidateConfig validates that all required fields are present in the YAML
func ValidateConfig(cfg *Config) error {
	// Validate core modules
	if err := validateModuleCategory("core", cfg.Core); err != nil {
		return err
	}

	// Validate system modules
	if err := validateModuleCategory("system", cfg.System); err != nil {
		return err
	}

	// Validate desktop modules
	if err := validateModuleCategory("desktop", cfg.Desktop); err != nil {
		return err
	}

	// Validate development modules
	if err := validateModuleCategory("development", cfg.Development); err != nil {
		return err
	}

	// Validate media modules
	if err := validateModuleCategory("media", cfg.Media); err != nil {
		return err
	}

	return nil
}

// validateModuleCategory validates all modules in a category
func validateModuleCategory(category string, modules map[string]Module) error {
	for name, module := range modules {
		fullName := fmt.Sprintf("%s.%s", category, name)

		if module.Start == "" {
			return fmt.Errorf("module %s missing required 'start' field", fullName)
		}

		if module.End == "" {
			return fmt.Errorf("module %s missing required 'end' field", fullName)
		}

		if module.Type == "" {
			return fmt.Errorf("module %s missing required 'type' field", fullName)
		}

		// Validate type is one of the allowed values
		validTypes := []string{"Package", "Git", "System", "File", "Module"}
		typeValid := false
		for _, validType := range validTypes {
			if module.Type == validType {
				typeValid = true
				break
			}
		}
		if !typeValid {
			return fmt.Errorf("module %s has invalid type '%s', must be one of: %v", fullName, module.Type, validTypes)
		}
	}

	return nil
}
