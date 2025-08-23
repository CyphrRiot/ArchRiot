package applications

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"archriot-installer/theming"
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
func (w *WaybarApplier) ApplyTheme(colors *theming.MatugenColors, dynamicEnabled bool) error {
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
func (w *WaybarApplier) generateColorsCSS(colors *theming.MatugenColors, dynamicEnabled bool) (string, error) {
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
		content.WriteString(fmt.Sprintf("@define-color foreground_primary %s;\n", colors.Colors.Dark.OnSurface))
		content.WriteString(fmt.Sprintf("@define-color foreground_secondary %s;\n", colors.Colors.Dark.OnPrimary))
		content.WriteString(fmt.Sprintf("@define-color border_active %s;\n", colors.Colors.Dark.Primary))
		content.WriteString(fmt.Sprintf("@define-color border_inactive %s;\n", colors.Colors.Dark.SurfaceVariant))
	} else {
		// Use static CypherRiot colors
		content.WriteString("/* Static CypherRiot Colors */\n")
		content.WriteString("@define-color primary_color #8b5cf6;\n")
		content.WriteString("@define-color accent_color #a855f7;\n")
		content.WriteString("@define-color secondary_color #c084fc;\n")
		content.WriteString("@define-color background_primary #1a1a1a;\n")
		content.WriteString("@define-color background_secondary #2d2d2d;\n")
		content.WriteString("@define-color background_tertiary #404040;\n")
		content.WriteString("@define-color background_sidebar #000000;\n")
		content.WriteString("@define-color foreground_primary #ffffff;\n")
		content.WriteString("@define-color foreground_secondary #e5e5e5;\n")
		content.WriteString("@define-color border_active #8b5cf6;\n")
		content.WriteString("@define-color border_inactive #404040;\n")
	}

	content.WriteString("\n/* End of ArchRiot Color Definitions */\n")
	return content.String(), nil
}
