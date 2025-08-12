package config

import (
	"strings"
	"testing"
)

func TestValidateDependencies_ValidConfig(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base": {
				Start:   "Installing base",
				End:     "Base installed",
				Type:    "Package",
				Depends: []string{},
			},
			"shell": {
				Start:   "Installing shell",
				End:     "Shell installed",
				Type:    "Package",
				Depends: []string{"core.base"},
			},
		},
		Desktop: map[string]Module{
			"hyprland": {
				Start:   "Installing Hyprland",
				End:     "Hyprland installed",
				Type:    "Package",
				Depends: []string{"core.base", "core.shell"},
			},
		},
	}

	err := ValidateDependencies(cfg)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestValidateDependencies_MissingDependency(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base": {
				Start:   "Installing base",
				End:     "Base installed",
				Type:    "Package",
				Depends: []string{"core.nonexistent"}, // Missing dependency
			},
		},
	}

	err := ValidateDependencies(cfg)
	if err == nil {
		t.Error("Expected error for missing dependency, but validation passed")
	}

	if !strings.Contains(err.Error(), "non-existent module core.nonexistent") {
		t.Errorf("Expected error about missing dependency, got: %v", err)
	}
}

func TestValidateDependencies_CircularDependency(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base": {
				Start:   "Installing base",
				End:     "Base installed",
				Type:    "Package",
				Depends: []string{"core.shell"}, // A depends on B
			},
			"shell": {
				Start:   "Installing shell",
				End:     "Shell installed",
				Type:    "Package",
				Depends: []string{"core.identity"}, // B depends on C
			},
			"identity": {
				Start:   "Setting up identity",
				End:     "Identity configured",
				Type:    "Git",
				Depends: []string{"core.base"}, // C depends on A - CYCLE!
			},
		},
	}

	err := ValidateDependencies(cfg)
	if err == nil {
		t.Error("Expected error for circular dependency, but validation passed")
	}

	if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("Expected circular dependency error, got: %v", err)
	}
}

func TestValidateDependencies_ComplexValidChain(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base": {
				Start:   "Installing base",
				End:     "Base installed",
				Type:    "Package",
				Depends: []string{},
			},
			"shell": {
				Start:   "Installing shell",
				End:     "Shell installed",
				Type:    "Package",
				Depends: []string{"core.base"},
			},
		},
		System: map[string]Module{
			"fonts": {
				Start:   "Installing fonts",
				End:     "Fonts installed",
				Type:    "System",
				Depends: []string{"core.base"},
			},
			"audio": {
				Start:   "Setting up audio",
				End:     "Audio configured",
				Type:    "System",
				Depends: []string{"system.fonts", "core.shell"},
			},
		},
		Desktop: map[string]Module{
			"hyprland": {
				Start:   "Installing Hyprland",
				End:     "Hyprland installed",
				Type:    "Package",
				Depends: []string{"core.base", "core.shell", "system.fonts", "system.audio"},
			},
		},
	}

	err := ValidateDependencies(cfg)
	if err != nil {
		t.Errorf("Expected complex valid chain to pass validation, got error: %v", err)
	}
}

func TestValidateDependencies_SelfDependency(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base": {
				Start:   "Installing base",
				End:     "Base installed",
				Type:    "Package",
				Depends: []string{"core.base"}, // Self-dependency
			},
		},
	}

	err := ValidateDependencies(cfg)
	if err == nil {
		t.Error("Expected error for self-dependency, but validation passed")
	}

	if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("Expected circular dependency error for self-dependency, got: %v", err)
	}
}

func TestCollectAllModules(t *testing.T) {
	cfg := &Config{
		Core: map[string]Module{
			"base":  {Start: "Base", End: "Base done", Type: "Package"},
			"shell": {Start: "Shell", End: "Shell done", Type: "Package"},
		},
		Desktop: map[string]Module{
			"hyprland": {Start: "Hyprland", End: "Hyprland done", Type: "Package"},
		},
		Development: map[string]Module{
			"tools": {Start: "Tools", End: "Tools done", Type: "Package"},
		},
	}

	allModules := collectAllModules(cfg)

	expectedModules := []string{
		"core.base",
		"core.shell",
		"desktop.hyprland",
		"development.tools",
	}

	if len(allModules) != len(expectedModules) {
		t.Errorf("Expected %d modules, got %d", len(expectedModules), len(allModules))
	}

	for _, expectedModule := range expectedModules {
		if _, exists := allModules[expectedModule]; !exists {
			t.Errorf("Expected module %s not found in collected modules", expectedModule)
		}
	}
}
