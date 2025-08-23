package theming

import "fmt"

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

// ThemeApplier defines the interface that all application theme implementations must satisfy
type ThemeApplier interface {
	// Name returns the human-readable name of this theme applier
	Name() string

	// ApplyTheme applies the theme colors to the application
	// If colors is nil or dynamicEnabled is false, should apply static/original theme
	ApplyTheme(colors *MatugenColors, dynamicEnabled bool) error

	// GetConfigPath returns the path to the application's config file
	GetConfigPath() (string, error)
}

// ColorTransformer provides utilities for color format conversion
type ColorTransformer interface {
	// HexToRGB converts hex color to RGB values
	HexToRGB(hex string) (r, g, b int)

	// HexToRGBA converts hex color to RGBA string with alpha
	HexToRGBA(hex string, alpha string) string

	// StripHex removes # prefix from hex color
	StripHex(hex string) string
}

// ConfigManager provides utilities for config file operations
type ConfigManager interface {
	// ReplaceProperty replaces a property value in config content
	ReplaceProperty(content, property, newValue string) string

	// ExtractProperty extracts a property value from config content
	ExtractProperty(content, property string) (string, error)

	// LoadTemplate loads the original config template
	LoadTemplate(templatePath string) (string, error)
}

// ThemeRegistry manages all registered theme appliers
type ThemeRegistry struct {
	appliers []ThemeApplier
}

// NewThemeRegistry creates a new theme registry
func NewThemeRegistry() *ThemeRegistry {
	return &ThemeRegistry{
		appliers: make([]ThemeApplier, 0),
	}
}

// Register adds a theme applier to the registry
func (r *ThemeRegistry) Register(applier ThemeApplier) {
	r.appliers = append(r.appliers, applier)
}

// ApplyAll applies theme to all registered appliers
func (r *ThemeRegistry) ApplyAll(colors *MatugenColors, dynamicEnabled bool) []error {
	var errors []error

	for _, applier := range r.appliers {
		if err := applier.ApplyTheme(colors, dynamicEnabled); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", applier.Name(), err))
		}
	}

	return errors
}

// GetAppliers returns all registered appliers
func (r *ThemeRegistry) GetAppliers() []ThemeApplier {
	return r.appliers
}
