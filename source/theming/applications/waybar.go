package applications

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WaybarApplier handles waybar theming
type WaybarApplier struct{}

// Name returns the human-readable name of this applier
func (w *WaybarApplier) Name() string {
	return "Waybar"
}

// GetConfigPath returns the path to waybar's colors.css file
func (w *WaybarApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "waybar", "colors.css"), nil
}

// ApplyTheme applies colors to waybar
func (w *WaybarApplier) ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error {
	colorsPath, err := w.GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(colorsPath), 0755); err != nil {
		return fmt.Errorf("creating colors directory: %w", err)
	}

	// Generate colors.css content
	content, err := w.generateColorsCSS(colors, dynamicEnabled)
	if err != nil {
		return fmt.Errorf("generating colors CSS: %w", err)
	}

	// Write colors.css
	if err := os.WriteFile(colorsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing colors CSS: %w", err)
	}

	// Reload waybar using SIGUSR2 (reload signal)
	if err := exec.Command("pkill", "-SIGUSR2", "waybar").Run(); err != nil {
		// Don't fail if waybar reload fails - just warn
		fmt.Printf("Warning: Failed to reload waybar: %v\n", err)
	}

	return nil
}

// generateColorsCSS creates the colors.css content
func (w *WaybarApplier) generateColorsCSS(colors *MatugenColors, dynamicEnabled bool) (string, error) {
	var content strings.Builder

	content.WriteString("/* ArchRiot Central Color Definitions */\n")
	content.WriteString("/* Waybar @define-color syntax for Dynamic Theming System */\n")
	content.WriteString("/* This file is the single source of truth for all ArchRiot colors */\n\n")

	if dynamicEnabled && colors != nil {
		// Use dynamic colors from wallpaper
		content.WriteString("/* Dynamic Colors from Wallpaper */\n")
		content.WriteString(fmt.Sprintf("@define-color primary_color %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color accent_color %s;\n", colors.Colors.Dark.Secondary))
		content.WriteString(fmt.Sprintf("@define-color secondary_color %s;\n", colors.Colors.Dark.Tertiary))
		content.WriteString(fmt.Sprintf("@define-color background_primary %s;\n", colors.Colors.Dark.Surface))
		content.WriteString(fmt.Sprintf("@define-color background_secondary %s;\n", colors.Colors.Dark.SurfaceVariant))
		content.WriteString(fmt.Sprintf("@define-color background_tertiary %s;\n", colors.Colors.Dark.SurfaceContainer))
		content.WriteString(fmt.Sprintf("@define-color background_sidebar %s;\n", colors.Colors.Dark.Background))
		content.WriteString(fmt.Sprintf("@define-color background_surface %s;\n", colors.Colors.Dark.Surface))
		content.WriteString(fmt.Sprintf("@define-color foreground_primary %s;\n", colors.Colors.Dark.OnSurface))
		content.WriteString(fmt.Sprintf("@define-color foreground_secondary %s;\n", colors.Colors.Dark.OnPrimary))
		content.WriteString(fmt.Sprintf("@define-color foreground_dim %s;\n", colors.Colors.Dark.OnSurface))
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

	content.WriteString("\n/* End of ArchRiot Color Definitions */\n")
	return content.String(), nil
}
