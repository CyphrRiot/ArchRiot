package tui

// LogMsg represents a log message
type LogMsg string

// ProgressMsg represents progress update
type ProgressMsg float64

// StepMsg represents a step update
type StepMsg string

// DoneMsg indicates completion
type DoneMsg struct{}

// GitUsernameMsg carries git username input
type GitUsernameMsg string

// GitEmailMsg carries git email input
type GitEmailMsg string

// RebootMsg carries reboot decision
type RebootMsg bool

// InputRequestMsg requests user input
type InputRequestMsg struct {
	Mode   string
	Prompt string
}

// GitConfirmMsg carries git credential confirmation
type GitConfirmMsg bool

// Helper functions for external packages
var versionGetter func() string
var logPathGetter func() string

// SetVersionGetter sets the function to get version
func SetVersionGetter(fn func() string) {
	versionGetter = fn
}

// SetLogPathGetter sets the function to get log path
func SetLogPathGetter(fn func() string) {
	logPathGetter = fn
}

// GetVersion returns the current version
func GetVersion() string {
	if versionGetter != nil {
		return versionGetter()
	}
	return "2.5.0"
}

// GetLogPath returns the current log path
func GetLogPath() string {
	if logPathGetter != nil {
		return logPathGetter()
	}
	return "/tmp/archriot-install.log"
}

// Git callback functions
var gitCompletionCallback func(bool)
var gitUsernameCallback func(string)
var gitEmailCallback func(string)

// SetGitCallbacks sets the callback functions for git credential handling
func SetGitCallbacks(completion func(bool), username func(string), email func(string)) {
	gitCompletionCallback = completion
	gitUsernameCallback = username
	gitEmailCallback = email
}
