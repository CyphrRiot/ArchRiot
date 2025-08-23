package applications

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GhosttyApplier handles Ghostty terminal theming
type GhosttyApplier struct{}

// Name returns the human-readable name of this applier
func (g *GhosttyApplier) Name() string {
	return "Ghostty Terminal"
}

// GetConfigPath returns the path to Ghostty's config file
func (g *GhosttyApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "ghostty", "config"), nil
}

// getTemplatePath returns the path to the original Ghostty template
func (g *GhosttyApplier) getTemplatePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "share", "archriot", "config", "ghostty", "config"), nil
}

// ApplyTheme applies colors to Ghostty terminal
func (g *GhosttyApplier) ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error {
	configPath, err := g.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting Ghostty config path: %w", err)
	}

	// Read current config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading Ghostty config: %w", err)
	}

	content := string(data)

	if dynamicEnabled && colors != nil {
		// Use dynamic colors from wallpaper
		content = g.applyDynamicColors(content, colors)
	} else {
		// Restore static colors from template
		content, err = g.applyStaticColors(content)
		if err != nil {
			return fmt.Errorf("applying static colors: %w", err)
		}
	}

	// Write updated config
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing Ghostty config: %w", err)
	}

	return nil
}

// applyDynamicColors applies dynamic colors from matugen palette
func (g *GhosttyApplier) applyDynamicColors(content string, colors *MatugenColors) string {
	// Extract colors and remove # prefix
	bg := strings.TrimPrefix(colors.Colors.Dark.Background, "#")
	fg := strings.TrimPrefix(colors.Colors.Dark.OnSurface, "#")
	primary := strings.TrimPrefix(colors.Colors.Dark.Primary, "#")
	secondary := strings.TrimPrefix(colors.Colors.Dark.Secondary, "#")
	tertiary := strings.TrimPrefix(colors.Colors.Dark.Tertiary, "#")
	surfaceVariant := strings.TrimPrefix(colors.Colors.Dark.SurfaceVariant, "#")

	// Replace basic colors
	content = g.replaceProperty(content, "background", bg)
	content = g.replaceProperty(content, "foreground", fg)
	content = g.replaceProperty(content, "cursor-color", fg)
	content = g.replaceProperty(content, "cursor-text", bg)
	content = g.replaceProperty(content, "selection-background", surfaceVariant)
	content = g.replaceProperty(content, "selection-foreground", fg)

	// Replace 16-color palette with dynamic colors
	content = g.replaceProperty(content, "palette = 0", bg)             // black
	content = g.replaceProperty(content, "palette = 8", surfaceVariant) // bright black
	content = g.replaceProperty(content, "palette = 4", primary)        // blue
	content = g.replaceProperty(content, "palette = 12", primary)       // bright blue
	content = g.replaceProperty(content, "palette = 5", secondary)      // magenta
	content = g.replaceProperty(content, "palette = 13", secondary)     // bright magenta
	content = g.replaceProperty(content, "palette = 6", tertiary)       // cyan
	content = g.replaceProperty(content, "palette = 14", tertiary)      // bright cyan
	content = g.replaceProperty(content, "palette = 7", fg)             // white
	content = g.replaceProperty(content, "palette = 15", fg)            // bright white

	return content
}

// applyStaticColors restores original colors from template
func (g *GhosttyApplier) applyStaticColors(content string) (string, error) {
	templatePath, err := g.getTemplatePath()
	if err != nil {
		return content, fmt.Errorf("getting template path: %w", err)
	}
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return content, fmt.Errorf("reading Ghostty template: %w", err)
	}

	templateContent := string(templateData)

	// Restore basic colors from template
	content = g.restoreProperty(content, templateContent, "background")
	content = g.restoreProperty(content, templateContent, "foreground")
	content = g.restoreProperty(content, templateContent, "cursor-color")
	content = g.restoreProperty(content, templateContent, "cursor-text")
	content = g.restoreProperty(content, templateContent, "selection-background")
	content = g.restoreProperty(content, templateContent, "selection-foreground")

	// Restore palette colors
	for i := 0; i <= 17; i++ {
		content = g.restoreProperty(content, templateContent, fmt.Sprintf("palette = %d", i))
	}

	return content, nil
}

// replaceProperty replaces a property value in Ghostty config
func (g *GhosttyApplier) replaceProperty(content, property, newValue string) string {
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

// restoreProperty restores a property value from template
func (g *GhosttyApplier) restoreProperty(content, templateContent, property string) string {
	// Extract value from template
	templateLines := strings.Split(templateContent, "\n")
	var templateValue string
	for _, line := range templateLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, property+" = ") {
			parts := strings.SplitN(trimmed, " = ", 2)
			if len(parts) == 2 {
				templateValue = parts[1]
				break
			}
		}
	}

	if templateValue != "" {
		return g.replaceProperty(content, property, templateValue)
	}
	return content
}
