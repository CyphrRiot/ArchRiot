package applications

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
