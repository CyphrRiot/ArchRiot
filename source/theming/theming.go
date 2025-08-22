package theming

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ThemeConfig represents the theming configuration
type ThemeConfig struct {
	DynamicThemingEnabled bool              `json:"dynamic_theming"`
	CurrentBackground     string            `json:"current_background"`
	GeneratedPalette      map[string]string `json:"generated_palette,omitempty"`
}

// MatugenColors represents the color palette from matugen
type MatugenColors struct {
	Colors struct {
		Light struct {
			Primary          string `json:"primary"`
			Secondary        string `json:"secondary"`
			Tertiary         string `json:"tertiary"`
			Surface          string `json:"surface"`
			SurfaceVariant   string `json:"surface_variant"`
			SurfaceContainer string `json:"surface_container"`
			OnPrimary        string `json:"on_primary"`
			OnSecondary      string `json:"on_secondary"`
			OnSurface        string `json:"on_surface"`
			Background       string `json:"background"`
		} `json:"light"`
		Dark struct {
			Primary          string `json:"primary"`
			Secondary        string `json:"secondary"`
			Tertiary         string `json:"tertiary"`
			Surface          string `json:"surface"`
			SurfaceVariant   string `json:"surface_variant"`
			SurfaceContainer string `json:"surface_container"`
			OnPrimary        string `json:"on_primary"`
			OnSecondary      string `json:"on_secondary"`
			OnSurface        string `json:"on_surface"`
			Background       string `json:"background"`
		} `json:"dark"`
	} `json:"colors"`
}

// ColorMapping maps matugen colors to ArchRiot CSS variables
type ColorMapping struct {
	MatugenColor string
	CSSVariable  string
}

// GetConfigPath returns the path to the theming config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "archriot", "background-prefs.json"), nil
}

// GetColorsPath returns the path to the central colors.css file
func GetColorsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "waybar", "colors.css"), nil
}

// GetZedConfigPath returns the path to the Zed settings.json file
func GetZedConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "zed", "settings.json"), nil
}

// GetZedTemplatePath returns the path to the original Zed config template
func GetZedTemplatePath() string {
	// This is relative to the ArchRiot source directory
	return "config/zed/settings.json"
}

// GetHyprlandConfigPath returns the path to the Hyprland config file
func GetHyprlandConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "hypr", "hyprland.conf"), nil
}

// LoadThemeConfig loads the theming configuration
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

// SaveThemeConfig saves the theming configuration
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

// ExtractColorsFromWallpaper uses matugen to extract colors from a wallpaper
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

	// Parse the JSON output
	var colors MatugenColors
	if err := json.Unmarshal(output, &colors); err != nil {
		return nil, fmt.Errorf("parsing matugen output: %w", err)
	}

	return &colors, nil
}

// GetColorMappings returns the mapping from matugen colors to CSS variables
func GetColorMappings() []ColorMapping {
	return []ColorMapping{
		// Primary colors from matugen dark theme
		{MatugenColor: "primary", CSSVariable: "--primary-color"},
		{MatugenColor: "secondary", CSSVariable: "--accent-color"},
		{MatugenColor: "tertiary", CSSVariable: "--secondary-color"},

		// Background colors
		{MatugenColor: "surface", CSSVariable: "--background-primary"},
		{MatugenColor: "surface_variant", CSSVariable: "--background-secondary"},
		{MatugenColor: "surface_container", CSSVariable: "--background-tertiary"},
		{MatugenColor: "background", CSSVariable: "--background-sidebar"},

		// Text colors
		{MatugenColor: "on_surface", CSSVariable: "--foreground-primary"},
		{MatugenColor: "on_primary", CSSVariable: "--foreground-secondary"},

		// Interactive elements
		{MatugenColor: "primary", CSSVariable: "--border-active"},
		{MatugenColor: "surface_variant", CSSVariable: "--border-inactive"},
	}
}

// GenerateColorsCSS generates the colors.css content with dynamic or static colors
func GenerateColorsCSS(colors *MatugenColors, dynamicEnabled bool) (string, error) {
	var content strings.Builder

	content.WriteString("/* ArchRiot Central Color Definitions */\n")
	content.WriteString("/* Waybar @define-color syntax for Dynamic Theming System */\n")
	content.WriteString("/* This file is the single source of truth for all ArchRiot colors */\n\n")

	if dynamicEnabled && colors != nil {
		// Use extracted colors
		content.WriteString("/* Dynamic Colors from Wallpaper */\n")
		content.WriteString(fmt.Sprintf("@define-color primary_color %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color accent_color %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color secondary_color %s;\n", colors.Colors.Dark.Tertiary))
		content.WriteString(fmt.Sprintf("@define-color background_primary %s;\n", colors.Colors.Dark.Surface))
		content.WriteString(fmt.Sprintf("@define-color background_secondary %s;\n", colors.Colors.Dark.SurfaceVariant))
		content.WriteString(fmt.Sprintf("@define-color background_tertiary %s;\n", colors.Colors.Dark.SurfaceContainer))
		content.WriteString(fmt.Sprintf("@define-color background_sidebar %s;\n", colors.Colors.Dark.Background))
		content.WriteString(fmt.Sprintf("@define-color foreground_primary %s;\n", colors.Colors.Dark.OnSurface))
		content.WriteString(fmt.Sprintf("@define-color foreground_secondary %s;\n", colors.Colors.Dark.OnPrimary))
		content.WriteString(fmt.Sprintf("@define-color border_active %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color border_inactive %s;\n", colors.Colors.Dark.SurfaceVariant))
	} else {
		// Use static CypherRiot colors
		content.WriteString("/* Static CypherRiot Colors */\n")
		content.WriteString("@define-color primary_color #7aa2f7;\n")
		content.WriteString("@define-color accent_color #bb9af7;\n")
		content.WriteString("@define-color secondary_color #9d7bd8;\n")
		content.WriteString("@define-color background_primary #0a0b10;\n")
		content.WriteString("@define-color background_secondary #0f1016;\n")
		content.WriteString("@define-color background_tertiary #1a1b26;\n")
		content.WriteString("@define-color background_sidebar #16161e;\n")
		content.WriteString("@define-color background_surface #292e42;\n")
		content.WriteString("@define-color foreground_primary #c0caf5;\n")
		content.WriteString("@define-color foreground_secondary #a9b1d6;\n")
		content.WriteString("@define-color foreground_dim #565f89;\n")
		content.WriteString("@define-color border_active #89b4fa;\n")
		content.WriteString("@define-color border_inactive #414868;\n")
	}

	// Always include semantic colors (keep static)
	content.WriteString("\n/* Semantic Colors (Always Static) */\n")
	content.WriteString("@define-color success_color #9ece6a;\n")
	content.WriteString("@define-color warning_color #e0af68;\n")
	content.WriteString("@define-color error_color #f7768e;\n")
	content.WriteString("@define-color info_color #0db9d7;\n")

	// Waybar specific colors
	content.WriteString("\n/* Waybar Specific Colors */\n")
	content.WriteString("@define-color waybar_bg_alpha rgba(0, 0, 0, 0.3);\n")
	content.WriteString("@define-color waybar_tooltip_bg #1e1e2e;\n")
	if dynamicEnabled && colors != nil {
		content.WriteString(fmt.Sprintf("@define-color workspace_active %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color workspace_inactive %s;\n", colors.Colors.Dark.SurfaceVariant))
		content.WriteString(fmt.Sprintf("@define-color workspace_hover %s;\n", colors.Colors.Dark.Secondary))
	} else {
		content.WriteString("@define-color workspace_active #3584e8;\n")
		content.WriteString("@define-color workspace_inactive #3a3a4a;\n")
		content.WriteString("@define-color workspace_hover #6b8ba6;\n")
	}

	// Component specific colors
	content.WriteString("\n/* Component Specific Colors */\n")
	if dynamicEnabled && colors != nil {
		content.WriteString(fmt.Sprintf("@define-color cpu_color %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color memory_color %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color temp_color %s;\n", colors.Colors.Dark.Tertiary))
		content.WriteString(fmt.Sprintf("@define-color power_color %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color lock_color %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color window_color %s;\n", colors.Colors.Dark.Tertiary))
	} else {
		content.WriteString("@define-color cpu_color #6a7de8;\n")
		content.WriteString("@define-color memory_color #8a95e8;\n")
		content.WriteString("@define-color temp_color #547ae0;\n")
		content.WriteString("@define-color power_color #a45ad0;\n")
		content.WriteString("@define-color lock_color #9c7ce8;\n")
		content.WriteString("@define-color window_color #4a2b7a;\n")
	}

	// Accent palette
	content.WriteString("\n/* Accent Palette */\n")
	if dynamicEnabled && colors != nil {
		content.WriteString(fmt.Sprintf("@define-color accent1 %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color accent2 %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color accent3 %s;\n", colors.Colors.Dark.Tertiary))
		content.WriteString(fmt.Sprintf("@define-color accent4 %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color accent5 %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color accent6 %s;\n", colors.Colors.Dark.Tertiary))
	} else {
		content.WriteString("@define-color accent1 #ff7a93;\n")
		content.WriteString("@define-color accent2 #0db9d7;\n")
		content.WriteString("@define-color accent3 #ff9e64;\n")
		content.WriteString("@define-color accent4 #bb9af7;\n")
		content.WriteString("@define-color accent5 #7da6ff;\n")
		content.WriteString("@define-color accent6 #0db9d7;\n")
	}

	// Selection/Highlight Colors
	content.WriteString("\n/* Selection/Highlight Colors */\n")
	content.WriteString("@define-color selection_bg rgba(122, 162, 247, 0.2);\n")
	content.WriteString("@define-color hover_bg rgba(122, 162, 247, 0.1);\n")

	// Shadow Colors
	content.WriteString("\n/* Shadow Colors */\n")
	content.WriteString("@define-color shadow_color rgba(26, 27, 38, 0.8);\n")
	content.WriteString("@define-color shadow_inactive rgba(26, 27, 38, 0.67);\n")

	content.WriteString("\n/* Dynamic Theming Flag */\n")
	content.WriteString("/* This will be modified by the theming system */\n")
	content.WriteString(fmt.Sprintf("/* dynamic_theming_enabled: %t */\n", dynamicEnabled))

	return content.String(), nil
}

// addAlpha adds alpha channel to a hex color
func addAlpha(hexColor, alpha string) string {
	// Remove # if present
	color := strings.TrimPrefix(hexColor, "#")

	// Convert hex to rgba
	if len(color) == 6 {
		return fmt.Sprintf("rgba(%s, %s)", hexToRGB(color), alpha)
	}

	return hexColor // Return original if can't parse
}

// hexToRGB converts hex color to RGB values
func hexToRGB(hex string) string {
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d, %d, %d", r, g, b)
}

// UpdateColorsFile updates the central colors.css file
func UpdateColorsFile(colors *MatugenColors, dynamicEnabled bool) error {
	colorsPath, err := GetColorsPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(colorsPath), 0755); err != nil {
		return fmt.Errorf("creating colors directory: %w", err)
	}

	// Generate new colors.css content
	content, err := GenerateColorsCSS(colors, dynamicEnabled)
	if err != nil {
		return fmt.Errorf("generating colors CSS: %w", err)
	}

	// Write to file
	if err := os.WriteFile(colorsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing colors file: %w", err)
	}

	return nil
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

	// Update colors file
	if err := UpdateColorsFile(colors, config.DynamicThemingEnabled); err != nil {
		return fmt.Errorf("updating colors file: %w", err)
	}

	// Save config
	if err := SaveThemeConfig(config); err != nil {
		return fmt.Errorf("saving theme config: %w", err)
	}

	// Update Zed editor theme
	if err := UpdateZedTheme(colors, config.DynamicThemingEnabled); err != nil {
		fmt.Printf("Warning: Failed to update Zed theme: %v\n", err)
	}

	// Update Hyprland window border colors
	if err := UpdateHyprlandColors(colors, config.DynamicThemingEnabled); err != nil {
		fmt.Printf("Warning: Failed to update Hyprland colors: %v\n", err)
	}

	// Reload waybar to pick up new colors using SIGUSR2 (reload signal)
	if err := exec.Command("pkill", "-SIGUSR2", "waybar").Run(); err != nil {
		// Don't fail if waybar reload fails
		fmt.Printf("Warning: Failed to reload waybar: %v\n", err)
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

	// Update colors file with static colors
	if err := UpdateColorsFile(nil, enabled); err != nil {
		return fmt.Errorf("updating colors file: %w", err)
	}

	// Update Zed editor theme
	if err := UpdateZedTheme(nil, enabled); err != nil {
		fmt.Printf("Warning: Failed to update Zed theme: %v\n", err)
	}

	// Update Hyprland window border colors
	if err := UpdateHyprlandColors(nil, enabled); err != nil {
		fmt.Printf("Warning: Failed to update Hyprland colors: %v\n", err)
	}

	// Save config
	if err := SaveThemeConfig(config); err != nil {
		return fmt.Errorf("saving theme config: %w", err)
	}

	return nil
}

// IsDynamicThemingEnabled checks if dynamic theming is enabled
func IsDynamicThemingEnabled() (bool, error) {
	config, err := LoadThemeConfig()
	if err != nil {
		return false, err
	}
	return config.DynamicThemingEnabled, nil
}

// ZedThemeOverrides represents Zed's experimental.theme_overrides section
type ZedThemeOverrides struct {
	EditorGutterBackground string                 `json:"editor.gutter.background"`
	PanelBackground        string                 `json:"panel.background"`
	BackgroundAppearance   string                 `json:"background.appearance"`
	ToolbarBackground      string                 `json:"toolbar.background"`
	EditorIndentGuide      string                 `json:"editor.indent_guide"`
	TitleBarBackground     string                 `json:"title_bar.background"`
	StatusBarBackground    string                 `json:"status_bar.background"`
	Background             string                 `json:"background"`
	EditorBackground       string                 `json:"editor.background"`
	TerminalBackground     string                 `json:"terminal.background"`
	Syntax                 map[string]interface{} `json:"syntax"`
}

// ZedSettings represents the structure of Zed's settings.json
type ZedSettings struct {
	ExperimentalThemeOverrides ZedThemeOverrides `json:"experimental.theme_overrides"`
	// We'll preserve other settings as raw JSON
	Other map[string]interface{} `json:"-"`
}

// LoadOriginalZedThemeOverrides loads the original theme overrides from the ArchRiot config template
func LoadOriginalZedThemeOverrides() (*ZedThemeOverrides, error) {
	templatePath := GetZedTemplatePath()

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("reading Zed template: %w", err)
	}

	var templateSettings map[string]interface{}
	if err := json.Unmarshal(data, &templateSettings); err != nil {
		return nil, fmt.Errorf("parsing Zed template: %w", err)
	}

	// Extract theme overrides from template
	themeOverridesRaw, exists := templateSettings["experimental.theme_overrides"]
	if !exists {
		return nil, fmt.Errorf("no theme overrides found in template")
	}

	// Convert to JSON and back to get proper struct
	overridesJSON, err := json.Marshal(themeOverridesRaw)
	if err != nil {
		return nil, fmt.Errorf("marshaling theme overrides: %w", err)
	}

	var overrides ZedThemeOverrides
	if err := json.Unmarshal(overridesJSON, &overrides); err != nil {
		return nil, fmt.Errorf("unmarshaling theme overrides: %w", err)
	}

	return &overrides, nil
}

// GenerateZedThemeOverrides creates theme overrides for Zed editor
func GenerateZedThemeOverrides(colors *MatugenColors, dynamicEnabled bool) ZedThemeOverrides {
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
			Syntax: map[string]interface{}{
				"comment":         map[string]string{"color": colors.Colors.Dark.OnSurface},
				"string":          map[string]string{"color": colors.Colors.Dark.Secondary},
				"emphasis.strong": map[string]string{"color": colors.Colors.Dark.Primary},
				"title":           map[string]string{"color": colors.Colors.Dark.Primary},
				"property":        map[string]string{"color": colors.Colors.Dark.Tertiary},
				"variable":        map[string]string{"color": colors.Colors.Dark.Tertiary},
			},
		}
	} else {
		// Use original ArchRiot theme from template
		originalOverrides, err := LoadOriginalZedThemeOverrides()
		if err != nil {
			// Fallback to hardcoded values if template loading fails
			return ZedThemeOverrides{
				EditorGutterBackground: "#000000",
				PanelBackground:        "#000000",
				BackgroundAppearance:   "opaque",
				ToolbarBackground:      "#000000",
				EditorIndentGuide:      "#000002",
				TitleBarBackground:     "#000000",
				StatusBarBackground:    "#000000",
				Background:             "#000000",
				EditorBackground:       "#000000",
				TerminalBackground:     "#000000",
				Syntax: map[string]interface{}{
					"comment":         map[string]string{"color": "#4d6878"},
					"string":          map[string]string{"color": "#9d7dd8"},
					"emphasis.strong": map[string]string{"color": "#4a90e2"},
					"title":           map[string]string{"color": "#7dd3fc"},
					"property":        map[string]string{"color": "#60a5fa"},
					"variable":        map[string]string{"color": "#60a5fa"},
				},
			}
		}
		return *originalOverrides
	}
}

// UpdateZedTheme updates the Zed editor theme with dynamic or static colors
func UpdateZedTheme(colors *MatugenColors, dynamicEnabled bool) error {
	zedConfigPath, err := GetZedConfigPath()
	if err != nil {
		return fmt.Errorf("getting Zed config path: %w", err)
	}

	// Read existing settings
	var existingSettings map[string]interface{}
	if data, err := os.ReadFile(zedConfigPath); err == nil {
		if err := json.Unmarshal(data, &existingSettings); err != nil {
			return fmt.Errorf("parsing existing Zed settings: %w", err)
		}
	} else {
		// If file doesn't exist, start with empty settings
		existingSettings = make(map[string]interface{})
	}

	// Generate new theme overrides
	themeOverrides := GenerateZedThemeOverrides(colors, dynamicEnabled)

	// Update theme overrides in existing settings
	existingSettings["experimental.theme_overrides"] = themeOverrides

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(zedConfigPath), 0755); err != nil {
		return fmt.Errorf("creating Zed config directory: %w", err)
	}

	// Write updated settings
	data, err := json.MarshalIndent(existingSettings, "", "    ")
	if err != nil {
		return fmt.Errorf("marshaling Zed settings: %w", err)
	}

	if err := os.WriteFile(zedConfigPath, data, 0644); err != nil {
		return fmt.Errorf("writing Zed settings file: %w", err)
	}

	return nil
}

// UpdateHyprlandColors updates Hyprland window border colors with dynamic or static colors
func UpdateHyprlandColors(colors *MatugenColors, dynamicEnabled bool) error {
	hyprlandConfigPath, err := GetHyprlandConfigPath()
	if err != nil {
		return fmt.Errorf("getting Hyprland config path: %w", err)
	}

	// Read current config
	data, err := os.ReadFile(hyprlandConfigPath)
	if err != nil {
		return fmt.Errorf("reading Hyprland config: %w", err)
	}

	content := string(data)

	if dynamicEnabled && colors != nil {
		// Use dynamic colors from wallpaper
		activeColor := colors.Colors.Dark.Primary
		inactiveColor := colors.Colors.Dark.SurfaceVariant

		// Convert hex to 6-char format and add alpha/gradient
		activeHex := strings.TrimPrefix(activeColor, "#")
		inactiveHex := strings.TrimPrefix(inactiveColor, "#")

		activeBorder := fmt.Sprintf("rgba(%s88) 45deg", activeHex)
		inactiveBorder := fmt.Sprintf("rgba(%s60)", inactiveHex)

		// Replace border colors using line-by-line replacement
		content = replaceHyprlandProperty(content, "col.active_border", activeBorder)
		content = replaceHyprlandProperty(content, "col.inactive_border", inactiveBorder)
		content = replaceHyprlandProperty(content, "col.border_active", activeBorder)
	} else {
		// Use static CypherRiot colors
		content = replaceHyprlandProperty(content, "col.active_border", "rgba(89b4fa88) 45deg")
		content = replaceHyprlandProperty(content, "col.inactive_border", "rgba(1a1a1a60)")
		content = replaceHyprlandProperty(content, "col.border_active", "rgba(89b4fa88) 45deg")
	}

	// Write updated config
	if err := os.WriteFile(hyprlandConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing Hyprland config: %w", err)
	}

	return nil
}

// replaceHyprlandProperty replaces a Hyprland property value in the config
func replaceHyprlandProperty(content, property, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, property+" = ") {
			// Keep the original indentation
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = fmt.Sprintf("%s%s = %s", indent, property, newValue)
			break
		}
	}
	return strings.Join(lines, "\n")
}
