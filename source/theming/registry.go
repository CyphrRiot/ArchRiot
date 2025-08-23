package theming

import (
	"archriot-installer/theming/applications"
	"fmt"
)

// Registry manages all theme appliers and coordinates batch operations
type Registry struct {
	appliers []applications.ThemeApplier
}

// NewRegistry creates a new theme registry with all supported applications
func NewRegistry() *Registry {
	registry := &Registry{
		appliers: make([]applications.ThemeApplier, 0),
	}

	// Register all supported applications
	registry.Register(&applications.WaybarApplier{})
	registry.Register(&applications.ZedApplier{})
	registry.Register(&applications.GhosttyApplier{})
	registry.Register(&applications.HyprlandApplier{})
	registry.Register(&applications.TextEditorApplier{})
	registry.Register(&applications.FuzzelApplier{})

	return registry
}

// Register adds a theme applier to the registry
func (r *Registry) Register(applier applications.ThemeApplier) {
	r.appliers = append(r.appliers, applier)
}

// ApplyAll applies theme to all registered appliers
func (r *Registry) ApplyAll(colors *applications.MatugenColors, dynamicEnabled bool) []error {
	var errors []error

	for _, applier := range r.appliers {
		if err := applier.ApplyTheme(colors, dynamicEnabled); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", applier.Name(), err))
		}
	}

	return errors
}

// GetAppliers returns all registered appliers
func (r *Registry) GetAppliers() []applications.ThemeApplier {
	return r.appliers
}

// GetApplierByName returns a specific applier by name
func (r *Registry) GetApplierByName(name string) applications.ThemeApplier {
	for _, applier := range r.appliers {
		if applier.Name() == name {
			return applier
		}
	}
	return nil
}
