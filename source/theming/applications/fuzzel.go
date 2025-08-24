package applications

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FuzzelApplier handles Fuzzel application launcher theming
type FuzzelApplier struct{}

// Name returns the human-readable name of this applier
func (f *FuzzelApplier) Name() string {
	return "Fuzzel Launcher"
}

// GetConfigPath returns the path to fuzzel's config file
func (f *FuzzelApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "fuzzel", "fuzzel.ini"), nil
}

// getTemplatePath returns the path to the original fuzzel template
func (f *FuzzelApplier) getTemplatePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "share", "archriot", "config", "fuzzel", "fuzzel.ini"), nil
}

// ApplyTheme applies colors to fuzzel launcher
func (f *FuzzelApplier) ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error {
	configPath, err := f.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting fuzzel config path: %w", err)
	}

	// Read current config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading fuzzel config: %w", err)
	}

	content := string(data)

	if dynamicEnabled && colors != nil {
		// Apply dynamic colors from wallpaper
		content = f.applyDynamicColors(content, colors)
	} else {
		// Restore static colors from template
		content, err = f.applyStaticColors(content)
		if err != nil {
			return fmt.Errorf("applying static colors: %w", err)
		}
	}

	// Write updated config
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing fuzzel config: %w", err)
	}

	return nil
}

// applyDynamicColors applies dynamic colors from matugen palette
func (f *FuzzelApplier) applyDynamicColors(content string, colors *MatugenColors) string {
	// Use dark background with good contrast, dynamic colors for accents
	background := f.toRGBA("#1a1b26", 0.95) // Original dark background with transparency
	text := f.toRGBA(colors.Colors.Dark.OnSurface, 1.0)
	border := f.toRGBA(colors.Colors.Dark.Primary, 1.0)
	selection := f.toRGBA(colors.Colors.Dark.Primary, 0.3)
	selectionText := f.toRGBA(colors.Colors.Dark.OnSurface, 1.0)
	selectionMatch := f.toRGBA(colors.Colors.Dark.Secondary, 1.0)

	// Replace color values in [colors] section
	content = f.replaceColor(content, "background", background)
	content = f.replaceColor(content, "text", text)
	content = f.replaceColor(content, "prompt", text)
	content = f.replaceColor(content, "input", text)
	content = f.replaceColor(content, "match", selectionMatch)
	content = f.replaceColor(content, "selection", selection)
	content = f.replaceColor(content, "selection-text", selectionText)
	content = f.replaceColor(content, "selection-match", selectionMatch)
	content = f.replaceColor(content, "border", border)

	return content
}

// applyStaticColors restores original colors from template
func (f *FuzzelApplier) applyStaticColors(content string) (string, error) {
	templatePath, err := f.getTemplatePath()
	if err != nil {
		return content, fmt.Errorf("getting template path: %w", err)
	}

	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return content, fmt.Errorf("reading fuzzel template: %w", err)
	}

	templateContent := string(templateData)

	// Restore colors from template
	colorKeys := []string{"background", "text", "prompt", "input", "match",
		"selection", "selection-text", "selection-match", "border"}

	for _, key := range colorKeys {
		content = f.restoreColor(content, templateContent, key)
	}

	return content, nil
}

// replaceColor replaces a color value in the [colors] section only
func (f *FuzzelApplier) replaceColor(content, key, newValue string) string {
	lines := strings.Split(content, "\n")
	inColorsSection := false

	for i, line := range lines {
		if strings.TrimSpace(line) == "[colors]" {
			inColorsSection = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "[") && strings.HasSuffix(strings.TrimSpace(line), "]") {
			inColorsSection = false
			continue
		}
		if inColorsSection && strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			// Replace this line
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			lines[i] = fmt.Sprintf("%s%s=%s", indent, key, newValue)
			break
		}
	}

	return strings.Join(lines, "\n")
}

// restoreColor restores a color value from template [colors] section only
func (f *FuzzelApplier) restoreColor(content, templateContent, key string) string {
	// Extract value from template [colors] section
	templateLines := strings.Split(templateContent, "\n")
	inColorsSection := false
	var templateValue string

	for _, line := range templateLines {
		if strings.TrimSpace(line) == "[colors]" {
			inColorsSection = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "[") && strings.HasSuffix(strings.TrimSpace(line), "]") {
			inColorsSection = false
			continue
		}
		if inColorsSection && strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
			if len(parts) == 2 {
				templateValue = parts[1]
				break
			}
		}
	}

	if templateValue != "" {
		return f.replaceColor(content, key, templateValue)
	}

	return content
}

// toRGBA converts hex color to RGBA hex format (RRGGBBAA)
func (f *FuzzelApplier) toRGBA(hexColor string, alpha float64) string {
	// Remove # if present
	hex := strings.TrimPrefix(hexColor, "#")

	// Ensure alpha is between 0 and 1
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}

	// Convert alpha to hex (0-255 -> 00-ff)
	alphaHex := fmt.Sprintf("%02x", int(alpha*255))

	return hex + alphaHex
}
