package theming

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ThemeConfig represents the theming configuration
type ThemeConfig struct {
	DynamicThemingEnabled bool              `json:"dynamic_theming"`
	CurrentBackground     string            `json:"current_background"`
	GeneratedPalette      map[string]string `json:"generated_palette,omitempty"`
}

// GetConfigPath returns the path to the theme config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "archriot", "background-prefs.json"), nil
}

// LoadThemeConfig loads the theme configuration
func LoadThemeConfig() (*ThemeConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &ThemeConfig{
			DynamicThemingEnabled: false,
			CurrentBackground:     "riot_01.jpg",
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config ThemeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &config, nil
}

// SaveThemeConfig saves the theme configuration
func SaveThemeConfig(config *ThemeConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// ExtractColorsFromWallpaper uses matugen to extract colors from wallpaper
func ExtractColorsFromWallpaper(wallpaperPath string) (*MatugenColors, error) {
	// Check if matugen is installed
	if _, err := exec.LookPath("matugen"); err != nil {
		return nil, fmt.Errorf("matugen not found: %w", err)
	}

	// Run matugen to extract colors
	cmd := exec.Command("matugen", "image", wallpaperPath, "--json", "hex")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running matugen: %w", err)
	}

	var colors MatugenColors
	if err := json.Unmarshal(output, &colors); err != nil {
		return nil, fmt.Errorf("parsing matugen output: %w", err)
	}

	return &colors, nil
}

// ApplyWallpaperTheme applies theming based on wallpaper
func ApplyWallpaperTheme(wallpaperPath string) error {
	// Load current config
	config, err := LoadThemeConfig()
	if err != nil {
		return fmt.Errorf("loading theme config: %w", err)
	}

	// Update current background
	config.CurrentBackground = filepath.Base(wallpaperPath)

	var colors *MatugenColors
	if config.DynamicThemingEnabled {
		// Extract colors from wallpaper
		colors, err = ExtractColorsFromWallpaper(wallpaperPath)
		if err != nil {
			// Fallback to static colors if extraction fails
			fmt.Printf("Warning: Color extraction failed, using static colors: %v\n", err)
			colors = nil
		} else {
			// Store extracted palette in config
			config.GeneratedPalette = map[string]string{
				"primary":   colors.Colors.Dark.Primary,
				"secondary": colors.Colors.Dark.Secondary,
				"surface":   colors.Colors.Dark.Surface,
			}
		}
	}

	// Save config
	if err := SaveThemeConfig(config); err != nil {
		return fmt.Errorf("saving theme config: %w", err)
	}

	// Apply theme to all registered applications
	registry := NewRegistry()
	errors := registry.ApplyAll(colors, config.DynamicThemingEnabled)

	// Log any application-specific errors but don't fail overall
	for _, err := range errors {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

// ToggleDynamicTheming enables or disables dynamic theming
func ToggleDynamicTheming(enabled bool) error {
	// Load current config
	config, err := LoadThemeConfig()
	if err != nil {
		return fmt.Errorf("loading theme config: %w", err)
	}

	// Update setting
	config.DynamicThemingEnabled = enabled

	// If we're enabling dynamic theming and have a current background, apply it
	if enabled && config.CurrentBackground != "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}

		// Try system backgrounds first
		wallpaperPath := filepath.Join(homeDir, ".local", "share", "archriot", "backgrounds", config.CurrentBackground)
		if _, err := os.Stat(wallpaperPath); os.IsNotExist(err) {
			// Try user backgrounds
			wallpaperPath = filepath.Join(homeDir, ".config", "archriot", "backgrounds", config.CurrentBackground)
		}

		if _, err := os.Stat(wallpaperPath); err == nil {
			return ApplyWallpaperTheme(wallpaperPath)
		}
	}

	// Save config
	if err := SaveThemeConfig(config); err != nil {
		return fmt.Errorf("saving theme config: %w", err)
	}

	// Apply static theme to all registered applications
	registry := NewRegistry()
	errors := registry.ApplyAll(nil, enabled)

	// Log any application-specific errors but don't fail overall
	for _, err := range errors {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

// IsDynamicThemingEnabled returns whether dynamic theming is currently enabled
func IsDynamicThemingEnabled() (bool, error) {
	config, err := LoadThemeConfig()
	if err != nil {
		return false, err
	}
	return config.DynamicThemingEnabled, nil
}
