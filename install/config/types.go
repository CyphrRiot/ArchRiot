package config

// Config represents the YAML structure
type Config struct {
	Core        map[string]Module `yaml:"core"`
	System      map[string]Module `yaml:"system"`
	Desktop     map[string]Module `yaml:"desktop"`
	Development map[string]Module `yaml:"development"`
	Media       map[string]Module `yaml:"media"`
}

// Module represents a single installation module
type Module struct {
	Packages    []string     `yaml:"packages"`
	Configs     []ConfigRule `yaml:"configs"`
	Commands    []string     `yaml:"commands,omitempty"`
	Handler     string       `yaml:"handler,omitempty"`
	Depends     []string     `yaml:"depends,omitempty"`
	Description string       `yaml:"description"`
}

// ConfigRule represents a configuration copying rule
type ConfigRule struct {
	Pattern          string   `yaml:"pattern"`
	Target           string   `yaml:"target,omitempty"`
	PreserveIfExists []string `yaml:"preserve_if_exists,omitempty"`
}

// ModuleOrder defines the execution order for different module categories
var ModuleOrder = map[string]int{
	"core":         10,
	"system":       20,
	"development":  30,
	"desktop":      40,
	"post-desktop": 45,
	"applications": 50,
	"optional":     60,
}
