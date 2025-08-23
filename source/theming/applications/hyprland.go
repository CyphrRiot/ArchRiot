package applications

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HyprlandApplier handles Hyprland window manager theming
type HyprlandApplier struct{}

// Name returns the human-readable name of this applier
func (h *HyprlandApplier) Name() string {
	return "Hyprland Window Manager"
}

// GetConfigPath returns the path to Hyprland's config file
func (h *HyprlandApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "hypr", "hyprland.conf"), nil
}

// ApplyTheme applies colors to Hyprland window manager
func (h *HyprlandApplier) ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error {
	configPath, err := h.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting Hyprland config path: %w", err)
	}

	// Read current config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading Hyprland config: %w", err)
	}

	content := string(data)

	if dynamicEnabled && colors != nil {
		// Use dynamic colors from wallpaper
		content = h.applyDynamicColors(content, colors)
	} else {
		// Use static CypherRiot colors
		content = h.applyStaticColors(content)
	}

	// Write updated config
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing Hyprland config: %w", err)
	}

	return nil
}

// applyDynamicColors applies dynamic border colors from matugen palette
func (h *HyprlandApplier) applyDynamicColors(content string, colors *MatugenColors) string {
	// Extract colors and convert to Hyprland format
	activeColor := strings.TrimPrefix(colors.Colors.Dark.Primary, "#")
	inactiveColor := strings.TrimPrefix(colors.Colors.Dark.SurfaceVariant, "#")

	// Convert to Hyprland RGBA format with alpha and gradients
	activeBorder := fmt.Sprintf("rgba(%s88) 45deg", activeColor)
	inactiveBorder := fmt.Sprintf("rgba(%s60)", inactiveColor)

	// Replace border colors
	content = h.replaceProperty(content, "col.active_border", activeBorder)
	content = h.replaceProperty(content, "col.inactive_border", inactiveBorder)
	content = h.replaceProperty(content, "col.border_active", activeBorder)

	return content
}

// applyStaticColors applies static CypherRiot border colors
func (h *HyprlandApplier) applyStaticColors(content string) string {
	// Use static CypherRiot colors
	content = h.replaceProperty(content, "col.active_border", "rgba(89b4fa88) 45deg")
	content = h.replaceProperty(content, "col.inactive_border", "rgba(1a1a1a60)")
	content = h.replaceProperty(content, "col.border_active", "rgba(89b4fa88) 45deg")

	return content
}

// replaceProperty replaces a Hyprland property value in the config
func (h *HyprlandApplier) replaceProperty(content, property, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, property+" = ") {
			// Keep original indentation
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = fmt.Sprintf("%s%s = %s", indent, property, newValue)
			break
		}
	}
	return strings.Join(lines, "\n")
}
