package config

import (
	"fmt"
	"regexp"
	"strings"
)

// dangerousPatterns defines regex patterns for commands that should never be allowed
var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`rm\s+-rf\s+/[^h]`),     // Never rm -rf / (except /home/user subdirs)
	regexp.MustCompile(`rm\s+-rf\s+/\s*$`),     // Never rm -rf / at end of line
	regexp.MustCompile(`>\s*/dev/sd[a-z]`),     // Don't write to block devices
	regexp.MustCompile(`dd\s+.*of=/dev/`),      // Dangerous dd commands to devices
	regexp.MustCompile(`chmod\s+777`),          // Overly permissive permissions
	regexp.MustCompile(`\$\(.*curl.*\|.*sh\)`), // Pipe curl to shell
	regexp.MustCompile(`sudo\s+rm\s+-rf\s+/`),  // Sudo + dangerous rm
	regexp.MustCompile(`mkfs\s+`),              // Filesystem creation
	regexp.MustCompile(`fdisk\s+`),             // Disk partitioning
	regexp.MustCompile(`parted\s+`),            // Disk partitioning
	regexp.MustCompile(`wipefs\s+`),            // Wipe filesystem signatures
}

// suspiciousPatterns defines patterns that are risky but might be legitimate with proper safety
var suspiciousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`rm\s+-rf\s+~`),            // Removing from home directory
	regexp.MustCompile(`rm\s+-rf\s+/tmp`),         // Removing /tmp contents
	regexp.MustCompile(`rm\s+-rf\s+/var/cache`),   // Removing cache directories
	regexp.MustCompile(`sudo\s+systemctl\s+stop`), // Stopping services
}

// systemModifyingCommands are commands that modify system state and should have error handling
var systemModifyingCommands = []string{
	"systemctl",
	"usermod",
	"groupadd",
	"useradd",
	"pacman",
	"yay",
	"mount",
	"umount",
	"chown",
	"chmod",
}

// ValidateCommands validates all commands in a module for safety
func ValidateCommands(module Module, moduleID string) error {
	for i, cmd := range module.Commands {
		if err := validateSingleCommand(cmd, moduleID, i+1); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleCommand validates a single command for safety patterns
func validateSingleCommand(cmd, moduleID string, cmdIndex int) error {
	// Check for dangerous patterns
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(cmd) {
			return fmt.Errorf("module %s command %d contains dangerous pattern '%s': %s",
				moduleID, cmdIndex, pattern.String(), cmd)
		}
	}

	// Check for suspicious patterns only if they're genuinely risky
	for _, pattern := range suspiciousPatterns {
		if pattern.MatchString(cmd) && isGenuinelyRisky(cmd) {
			return fmt.Errorf("module %s command %d contains risky pattern '%s': %s",
				moduleID, cmdIndex, pattern.String(), cmd)
		}
	}

	// Check for commands that should never use certain flags
	if err := validateCommandFlags(cmd, moduleID, cmdIndex); err != nil {
		return err
	}

	return nil
}

// hasSafetyHandling checks if a command has proper error handling
func hasSafetyHandling(cmd string) bool {
	return strings.Contains(cmd, "|| true") ||
		strings.Contains(cmd, "2>/dev/null") ||
		strings.Contains(cmd, "&>/dev/null") ||
		strings.Contains(cmd, "|| echo") ||
		strings.Contains(cmd, "|| return")
}

// isSystemModifyingCommand checks if a command modifies system state
func isSystemModifyingCommand(cmd string) bool {
	cmdLower := strings.ToLower(cmd)

	for _, syscmd := range systemModifyingCommands {
		if strings.Contains(cmdLower, syscmd) {
			return true
		}
	}

	// Check for sudo usage (generally system-modifying)
	if strings.Contains(cmdLower, "sudo") {
		return true
	}

	return false
}

// isKnownSafeCommand checks if a command is known to be safe without error handling
func isKnownSafeCommand(cmd string) bool {
	safeCommands := []string{
		"pacman -S --noconfirm --needed", // Package installation with safe flags
		"systemctl --user enable",        // User service management
		"systemctl --user start",         // User service management
		"systemctl --user daemon-reload", // User daemon reload
	}

	for _, safe := range safeCommands {
		if strings.Contains(cmd, safe) {
			return true
		}
	}

	return false
}

// validateCommandFlags checks for problematic command flags
func validateCommandFlags(cmd, moduleID string, cmdIndex int) error {
	// Check for recursive removal of dangerous targets
	if strings.Contains(cmd, "rm -rf") {
		// Block dangerous targets regardless of error handling
		dangerousPaths := []string{
			"rm -rf /",
			"rm -rf /*",
			"rm -rf /bin",
			"rm -rf /usr",
			"rm -rf /etc",
			"rm -rf /lib",
			"rm -rf /boot",
		}

		for _, dangerous := range dangerousPaths {
			if strings.Contains(cmd, dangerous) {
				return fmt.Errorf("module %s command %d attempts to remove critical system directory: %s",
					moduleID, cmdIndex, cmd)
			}
		}
	}

	// Check for overly broad pacman operations
	if strings.Contains(cmd, "pacman -R") && strings.Contains(cmd, "--cascade") {
		return fmt.Errorf("module %s command %d uses dangerous pacman --cascade flag: %s",
			moduleID, cmdIndex, cmd)
	}

	return nil
}

// isGenuinelyRisky determines if a suspicious command is actually dangerous
func isGenuinelyRisky(cmd string) bool {
	// Only flag commands that could cause permanent damage
	genuinelyRiskyPatterns := []string{
		"rm -rf ~",        // Removing entire home directory
		"rm -rf /home",    // Removing all user data
		"rm -rf /var/lib", // Removing system data
	}

	for _, risky := range genuinelyRiskyPatterns {
		if strings.Contains(cmd, risky) {
			return true
		}
	}

	return false
}

// ValidateAllCommands validates commands across all modules in a config
func ValidateAllCommands(cfg *Config) error {
	allModules := collectAllModules(cfg)

	for moduleID, module := range allModules {
		if err := ValidateCommands(module, moduleID); err != nil {
			return err
		}
	}

	return nil
}
