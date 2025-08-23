package applications

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"archriot-installer/theming"
)

// ZedApplier handles Zed editor theming
type ZedApplier struct{}

// Name returns the human-readable name of this applier
func (z *ZedApplier) Name() string {
	return "Zed Editor"
}

// GetConfigPath returns the path to Zed's settings.json file
func (z *ZedApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "zed", "settings.json"), nil
}

// getTemplatePath returns the path to the original Zed template
func (z *ZedApplier) getTemplatePath() string {
	return "config/zed/settings.json"
}

// ApplyTheme applies colors to Zed editor
func (z *ZedApplier) ApplyTheme(colors *theming.MatugenColors, dynamicEnabled bool) error {
	configPath, err := z.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting Zed config path: %w", err)
	}

	// Read existing settings
	var existingSettings map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &existingSettings); err != nil {
			return fmt.Errorf("parsing existing Zed settings: %w", err)
		}
	} else {
		existingSettings = make(map[string]interface{})
	}

	// Generate theme overrides
	themeOverrides := z.generateThemeOverrides(colors, dynamicEnabled)

	// Update the experimental.theme_overrides section
	existingSettings["experimental.theme_overrides"] = themeOverrides

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating Zed config directory: %w", err)
	}

	// Write updated settings
	data, err := json.MarshalIndent(existingSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling Zed settings: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing Zed settings: %w", err)
	}

	return nil
}

// ZedThemeOverrides represents Zed's theme override structure
type ZedThemeOverrides struct {
	EditorGutterBackground string `json:"editor.gutter.background"`
	PanelBackground        string `json:"panel.background"`
	BackgroundAppearance   string `json:"background.appearance"`
	ToolbarBackground      string `json:"toolbar.background"`
	EditorIndentGuide      string `json:"editor.indent_guide"`
	TitleBarBackground     string `json:"title_bar.background"`
	StatusBarBackground    string `json:"status_bar.background"`
	Background             string `json:"background"`
	EditorBackground       string `json:"editor.background"`
	TerminalBackground     string `json:"terminal.background"`
	SyntaxComment          string `json:"syntax.comment"`
	SyntaxString           string `json:"syntax.string"`
	SyntaxProperty         string `json:"syntax.property"`
	SyntaxTitle            string `json:"syntax.title"`
}

// generateThemeOverrides creates theme overrides based on colors and dynamic state
func (z *ZedApplier) generateThemeOverrides(colors *theming.MatugenColors, dynamicEnabled bool) ZedThemeOverrides {
	if dynamicEnabled && colors != nil {
		// Use dynamic colors from wallpaper
		return ZedThemeOverrides{
			EditorGutterBackground: colors.Colors.Dark.Background,
			PanelBackground:        colors.Colors.Dark.Background,
			BackgroundAppearance:   "opaque",
			ToolbarBackground:      colors.Colors.Dark.Background,
			EditorIndentGuide:      colors.Colors.Dark.SurfaceVariant,
			TitleBarBackground:     colors.Colors.Dark.Background,
			StatusBarBackground:    colors.Colors.Dark.Background,
			Background:             colors.Colors.Dark.Background,
			EditorBackground:       colors.Colors.Dark.Background,
			TerminalBackground:     colors.Colors.Dark.Background,
			SyntaxComment:          colors.Colors.Dark.OnSurface,
			SyntaxString:           colors.Colors.Dark.Secondary,
			SyntaxProperty:         colors.Colors.Dark.Tertiary,
			SyntaxTitle:            colors.Colors.Dark.Primary,
		}
	}

	// Load original ArchRiot theme from template
	originalOverrides, err := z.loadOriginalThemeOverrides()
	if err != nil {
		// Fallback to hardcoded values if template loading fails
		return ZedThemeOverrides{
			EditorGutterBackground: "#000000",
			PanelBackground:        "#000000",
			BackgroundAppearance:   "opaque",
			ToolbarBackground:      "#000000",
			EditorIndentGuide:      "#2d2d2d",
			TitleBarBackground:     "#000000",
			StatusBarBackground:    "#000000",
			Background:             "#000000",
			EditorBackground:       "#000000",
			TerminalBackground:     "#000000",
			SyntaxComment:          "#666666",
			SyntaxString:           "#a855f7",
			SyntaxProperty:         "#c084fc",
			SyntaxTitle:            "#8b5cf6",
		}
	}

	return *originalOverrides
}

// loadOriginalThemeOverrides loads the original theme from the template file
func (z *ZedApplier) loadOriginalThemeOverrides() (*ZedThemeOverrides, error) {
	templatePath := z.getTemplatePath()

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("reading Zed template: %w", err)
	}

	var templateSettings map[string]interface{}
	if err := json.Unmarshal(data, &templateSettings); err != nil {
		return nil, fmt.Errorf("parsing Zed template: %w", err)
	}

	// Extract experimental.theme_overrides from template
	overrides, exists := templateSettings["experimental.theme_overrides"]
	if !exists {
		return nil, fmt.Errorf("no theme_overrides found in template")
	}

	// Convert to our struct
	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		return nil, fmt.Errorf("marshaling template overrides: %w", err)
	}

	var themeOverrides ZedThemeOverrides
	if err := json.Unmarshal(overridesJSON, &themeOverrides); err != nil {
		return nil, fmt.Errorf("unmarshaling theme overrides: %w", err)
	}

	return &themeOverrides, nil
}
