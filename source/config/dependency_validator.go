package config

import (
	"fmt"
	"reflect"
)

// ValidateDependencies checks for circular dependencies and missing references
func ValidateDependencies(cfg *Config) error {
	allModules := collectAllModules(cfg)

	// Check for missing dependencies
	for moduleID, module := range allModules {
		for _, dep := range module.Depends {
			if _, exists := allModules[dep]; !exists {
				return fmt.Errorf("module %s depends on non-existent module %s", moduleID, dep)
			}
		}
	}

	// Check for circular dependencies
	for moduleID := range allModules {
		if hasCyclicDependency(moduleID, allModules, make(map[string]bool), make(map[string]bool)) {
			return fmt.Errorf("circular dependency detected involving module %s", moduleID)
		}
	}

	return nil
}

// collectAllModules gathers all modules from all categories into a single map
func collectAllModules(cfg *Config) map[string]Module {
	allModules := make(map[string]Module)

	// Use reflection to dynamically discover all module categories
	configValue := reflect.ValueOf(cfg).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		fieldType := configType.Field(i)

		// Skip non-map fields or maps that aren't map[string]Module
		if field.Kind() != reflect.Map || field.Type().String() != "map[string]config.Module" {
			continue
		}

		// Get category name from YAML tag or field name
		categoryName := fieldType.Tag.Get("yaml")
		if categoryName == "" {
			categoryName = fieldType.Name
		}

		// Iterate through modules in this category
		for _, key := range field.MapKeys() {
			moduleName := key.String()
			moduleValue := field.MapIndex(key)
			module := moduleValue.Interface().(Module)

			allModules[categoryName+"."+moduleName] = module
		}
	}

	return allModules
}

// hasCyclicDependency uses depth-first search to detect cycles
func hasCyclicDependency(moduleID string, modules map[string]Module, visiting, visited map[string]bool) bool {
	if visiting[moduleID] {
		return true // Found cycle - we're visiting a node we're already visiting
	}
	if visited[moduleID] {
		return false // Already processed this node completely
	}

	visiting[moduleID] = true
	module := modules[moduleID]

	for _, dep := range module.Depends {
		if hasCyclicDependency(dep, modules, visiting, visited) {
			return true
		}
	}

	visiting[moduleID] = false
	visited[moduleID] = true
	return false
}
