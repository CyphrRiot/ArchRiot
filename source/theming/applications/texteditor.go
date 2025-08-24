package applications

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TextEditorApplier handles GNOME Text Editor theming
type TextEditorApplier struct{}

// Name returns the human-readable name of this applier
func (t *TextEditorApplier) Name() string {
	return "Text Editor"
}

// GetConfigPath returns the path to the text editor's dynamic theme file
func (t *TextEditorApplier) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	stylesDir := filepath.Join(homeDir, ".local", "share", "gtksourceview-5", "styles")
	return filepath.Join(stylesDir, "cypherriot-dynamic.xml"), nil
}

// getTemplatePath returns the path to the original CypherRiot theme template
func (t *TextEditorApplier) getTemplatePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "share", "archriot", "config", "text-editor", "cypherriot.xml"), nil
}

// ApplyTheme applies colors to Text Editor
func (t *TextEditorApplier) ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error {
	configPath, err := t.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting text editor config path: %w", err)
	}

	// Ensure styles directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating styles directory: %w", err)
	}

	if dynamicEnabled && colors != nil {
		// Generate dynamic theme XML
		content, err := t.generateDynamicTheme(colors)
		if err != nil {
			return fmt.Errorf("generating dynamic theme: %w", err)
		}

		// Write dynamic theme file
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing dynamic theme: %w", err)
		}

		// Set gsettings to use dynamic theme
		if err := exec.Command("gsettings", "set", "org.gnome.TextEditor", "style-scheme", "cypherriot-dynamic").Run(); err != nil {
			return fmt.Errorf("setting dynamic theme via gsettings: %w", err)
		}
	} else {
		// Remove dynamic theme file to fall back to static
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing dynamic theme: %w", err)
		}

		// Set gsettings back to static theme
		if err := exec.Command("gsettings", "set", "org.gnome.TextEditor", "style-scheme", "cypherriot").Run(); err != nil {
			return fmt.Errorf("setting static theme via gsettings: %w", err)
		}
	}

	return nil
}

// generateDynamicTheme creates dynamic XML theme content from matugen colors
func (t *TextEditorApplier) generateDynamicTheme(colors *MatugenColors) (string, error) {
	// Use proper contrast mapping - keep original dark backgrounds, apply dynamic colors to accents
	bg := "#0a0b10"      // Original CypherRiot very dark background
	bgDark := "#16161e"  // Original darker background
	bgFloat := "#222436" // Original floating background with blue tint

	// Extract dynamic colors for accents and UI elements
	fg := t.ensureHexPrefix(colors.Colors.Dark.OnSurface)
	primary := t.ensureHexPrefix(colors.Colors.Dark.Primary)
	secondary := t.ensureHexPrefix(colors.Colors.Dark.Secondary)
	tertiary := t.ensureHexPrefix(colors.Colors.Dark.Tertiary)

	// Use dynamic surface colors for selections and highlights (better contrast)
	surface := t.ensureHexPrefix(colors.Colors.Dark.SurfaceVariant)

	var content strings.Builder
	content.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!--
 CypherRiot Dynamic Theme for GNOME Text Editor
 Generated automatically from wallpaper colors via ArchRiot theming system
-->
<style-scheme id="cypherriot-dynamic" _name="CypherRiot Dynamic" version="1.0">
  <author>ArchRiot Dynamic Theming</author>
  <_description>CypherRiot theme with dynamic colors from wallpaper</_description>

  <!-- Dynamic Color Palette -->
`)
	content.WriteString(fmt.Sprintf(`  <color name="bg"         value="%s"/>  <!-- Primary background (dark) -->
`, bg))
	content.WriteString(fmt.Sprintf(`  <color name="bg-dark"    value="%s"/>  <!-- Darker background -->
`, bgDark))
	content.WriteString(fmt.Sprintf(`  <color name="bg-float"   value="%s"/>  <!-- Floating background -->
`, bgFloat))
	content.WriteString(fmt.Sprintf(`  <color name="fg"         value="%s"/>  <!-- Foreground text -->
`, fg))
	content.WriteString(fmt.Sprintf(`  <color name="primary"    value="%s"/>  <!-- Primary accent -->
`, primary))
	content.WriteString(fmt.Sprintf(`  <color name="secondary"  value="%s"/>  <!-- Secondary accent -->
`, secondary))
	content.WriteString(fmt.Sprintf(`  <color name="tertiary"   value="%s"/>  <!-- Tertiary accent -->
`, tertiary))
	content.WriteString(fmt.Sprintf(`  <color name="surface"    value="%s"/>  <!-- Selection/highlight color -->
`, surface))

	content.WriteString(`
  <!-- Global Settings -->
  <style name="text"                        foreground="fg" background="bg"/>
  <style name="selection"                   foreground="fg" background="surface"/>
  <style name="cursor"                      foreground="fg"/>
  <style name="current-line"                background="bg-float"/>
  <style name="current-line-number"         foreground="primary" background="bg-float" bold="true"/>
  <style name="line-numbers"                foreground="secondary" background="bg-dark"/>

  <!-- Bracket Matching -->
  <style name="bracket-match"               foreground="primary" background="surface" bold="true"/>
  <style name="bracket-mismatch"            foreground="tertiary" background="bg-dark" bold="true"/>

  <!-- Search Matching -->
  <style name="search-match"                foreground="bg" background="primary" bold="true"/>

  <!-- Comments -->
  <style name="def:comment"                 foreground="secondary" italic="true"/>
  <style name="def:shebang"                 foreground="secondary" bold="true"/>

  <!-- Constants -->
  <style name="def:constant"                foreground="tertiary"/>
  <style name="def:string"                  foreground="primary"/>
  <style name="def:character"               foreground="primary"/>
  <style name="def:number"                  foreground="tertiary"/>
  <style name="def:boolean"                 foreground="tertiary"/>

  <!-- Identifiers -->
  <style name="def:identifier"              foreground="fg"/>
  <style name="def:function"                foreground="secondary"/>

  <!-- Statements -->
  <style name="def:statement"               foreground="primary" bold="true"/>
  <style name="def:keyword"                 foreground="primary" bold="true"/>
  <style name="def:operator"                foreground="secondary"/>

  <!-- Types -->
  <style name="def:type"                    foreground="tertiary" bold="true"/>

  <!-- Markup -->
  <style name="def:emphasis"                italic="true"/>
  <style name="def:strong-emphasis"         foreground="primary" bold="true"/>
  <style name="def:heading"                 foreground="secondary" bold="true"/>
  <style name="def:link-text"               foreground="tertiary" underline="single"/>

  <!-- Others -->
  <style name="def:preprocessor"            foreground="secondary"/>
  <style name="def:error"                   foreground="tertiary" background="bg-dark" bold="true"/>
  <style name="def:warning"                 foreground="primary" background="bg-dark"/>

</style-scheme>`)

	return content.String(), nil
}

// ensureHexPrefix ensures color has # prefix
func (t *TextEditorApplier) ensureHexPrefix(color string) string {
	if !strings.HasPrefix(color, "#") {
		return "#" + color
	}
	return color
}
