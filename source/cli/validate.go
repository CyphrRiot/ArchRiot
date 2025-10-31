package cli

import (
	"fmt"

	"archriot-installer/config"
)

// ValidateConfig validates the packages.yaml configuration and prints progress messages.
// It mirrors the previous inline validation logic while living in the cli package
// to keep main.go focused on delegation only.
func ValidateConfig() error {
	fmt.Println("🔍 Validating packages.yaml configuration...")

	// Find configuration file
	configPath := config.FindConfigFile()
	if configPath == "" {
		return fmt.Errorf("packages.yaml not found")
	}

	fmt.Printf("📁 Found config: %s\n", configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Println("✅ YAML structure is valid")

	// Validate basic config structure
	if err := config.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	fmt.Println("✅ Required fields are present")

	// Validate dependencies
	if err := config.ValidateDependencies(cfg); err != nil {
		return fmt.Errorf("dependency validation: %w", err)
	}

	fmt.Println("✅ Dependencies are valid (no cycles, all references exist)")

	// Validate command safety
	if err := config.ValidateAllCommands(cfg); err != nil {
		return fmt.Errorf("command safety validation: %w", err)
	}

	fmt.Println("✅ Commands are safe (no dangerous patterns detected)")

	return nil
}
