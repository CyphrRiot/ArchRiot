package tui

// LogMsg represents a log message
type LogMsg string

// ProgressMsg represents progress update
type ProgressMsg float64

// StepMsg represents a step update
type StepMsg string

// DoneMsg indicates successful completion
type DoneMsg struct{}

// FailureMsg indicates installation failure
type FailureMsg struct {
	Error string
}

// GitUsernameMsg carries git username input
type GitUsernameMsg string

// GitEmailMsg carries git email input
type GitEmailMsg string

// RebootMsg carries reboot decision
type RebootMsg bool

// UpgradeMsg triggers upgrade confirmation prompt
type UpgradeMsg struct{}

// KernelUpgradeMsg indicates kernel was upgraded and reboot is recommended
type KernelUpgradeMsg bool

// PreservationPromptMsg triggers preservation confirmation prompt
type PreservationPromptMsg struct{}

// InputRequestMsg requests user input
type InputRequestMsg struct {
	Mode   string // git-username, git-email, git-confirm, etc.
	Prompt string
}

// GitConfirmMsg carries git credential confirmation
type GitConfirmMsg bool

// SecureBootStatusMsg carries Secure Boot and LUKS detection results
type SecureBootStatusMsg struct {
	Enabled     bool
	Supported   bool
	LuksUsed    bool
	LuksDevices []string
}

// SecureBootPromptMsg triggers Secure Boot enablement prompt
type SecureBootPromptMsg struct{}

// SecureBootConfirmMsg carries Secure Boot enablement decision
type SecureBootConfirmMsg bool

// SecureBootContinuationPromptMsg triggers retry/cancel prompt for Secure Boot continuation
type SecureBootContinuationPromptMsg struct{}

// SecureBootContinuationMsg carries continuation decision (retry/cancel)
type SecureBootContinuationMsg bool

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
	return "unknown"
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

// Upgrade callback function
var upgradeCompletionCallback func(bool)

// SetUpgradeCallback sets the callback function for upgrade confirmation handling
func SetUpgradeCallback(callback func(bool)) {
	upgradeCompletionCallback = callback
}

// Preservation callback function
var preservationCompletionCallback func(bool)

// SetPreservationCallback sets the callback function for preservation confirmation handling
func SetPreservationCallback(callback func(bool)) {
	preservationCompletionCallback = callback
}

// Secure Boot callback function
var secureBootCompletionCallback func(bool)

// SetSecureBootCallback sets the callback function for Secure Boot confirmation handling
func SetSecureBootCallback(callback func(bool)) {
	secureBootCompletionCallback = callback
}

// Secure Boot continuation callback function
var secureBootContinuationCallback func(bool)

// SetSecureBootContinuationCallback sets the callback function for Secure Boot continuation handling
func SetSecureBootContinuationCallback(callback func(bool)) {
	secureBootContinuationCallback = callback
}
