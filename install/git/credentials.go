package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

var (
	GitUsername   string
	GitEmail      string
	GitConfirmUse bool
	GitInputDone  chan bool
	Program       *tea.Program
)

// SetProgram sets the TUI program reference for sending messages
func SetProgram(p *tea.Program) {
	Program = p
}

// SetGitInputChannel sets the channel for git input completion
func SetGitInputChannel(ch chan bool) {
	GitInputDone = ch
}

// SetGitCredentials sets the git credentials from TUI callbacks
func SetGitCredentials(username, email string, confirmed bool) {
	if username != "" {
		GitUsername = username
	}
	if email != "" {
		GitEmail = email
	}
	GitConfirmUse = confirmed
}

// SetGitConfirm sets only the confirmation status
func SetGitConfirm(confirmed bool) {
	GitConfirmUse = confirmed
}

// SetGitUsername sets only the username
func SetGitUsername(username string) {
	GitUsername = username
}

// SetGitEmail sets only the email
func SetGitEmail(email string) {
	GitEmail = email
}

// sendFormattedLog sends a properly formatted log message to TUI
func sendFormattedLog(status, emoji, name, description string) {
	if Program != nil {
		Program.Send(tui.LogMsg(fmt.Sprintf("%s %s %-15s %s", status, emoji, name, description)))
	}
}

// HandleGitConfiguration applies Git configuration with beautiful styling
func HandleGitConfiguration() error {
	logger.LogMessage("INFO", "🔧 Applying Git configuration...")

	if Program != nil {
		sendFormattedLog("🔄", "📦", "Git Setup", "Checking credentials")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Check for existing git credentials in git config first
	existingName, _ := runGitConfigGet("user.name")
	existingEmail, _ := runGitConfigGet("user.email")

	var userName, userEmail string

	// If we have existing git credentials, ask to use them
	if strings.TrimSpace(existingName) != "" || strings.TrimSpace(existingEmail) != "" {
		if Program != nil {
			sendFormattedLog("🎉", "📦", "Git Found", "Found existing git credentials")
			sendFormattedLog("📋", "📦", "Name", existingName)
			sendFormattedLog("📋", "📦", "Email", existingEmail)
			Program.Send(tui.InputRequestMsg{Mode: "git-confirm", Prompt: ""})
		}

		// Wait for confirmation
		<-GitInputDone

		if GitConfirmUse {
			userName = existingName
			userEmail = existingEmail
			if Program != nil {
				sendFormattedLog("✅", "📦", "Git Setup", "Using existing credentials")
			}
		} else {
			// User said no, prompt for new credentials
			if Program != nil {
				sendFormattedLog("💬", "📦", "Git Setup", "Setting up new credentials")
				Program.Send(tui.InputRequestMsg{Mode: "git-username", Prompt: "Git Username: "})
			}

			// Wait for new credentials
			<-GitInputDone

			userName = GitUsername
			userEmail = GitEmail
		}
	} else {
		// No existing credentials, prompt for new ones
		if Program != nil {
			sendFormattedLog("💬", "📦", "Git Setup", "No credentials found, setting up")
			Program.Send(tui.InputRequestMsg{Mode: "git-username", Prompt: "Git Username: "})
		}

		// Wait for credentials to be entered
		<-GitInputDone

		userName = GitUsername
		userEmail = GitEmail
	}

	// Save credentials to user.env
	if err := saveGitCredentials(homeDir, userName, userEmail); err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to save git credentials: %v", err))
	}

	// Apply Git user configuration
	if strings.TrimSpace(userName) != "" {
		if err := runGitConfig("user.name", userName); err != nil {
			return fmt.Errorf("setting git user.name: %w", err)
		}
		if Program != nil {
			sendFormattedLog("✅", "📦", "Git Identity", "User name set to: "+userName)
		}
		logger.LogMessage("SUCCESS", fmt.Sprintf("Git user.name set to: %s", userName))
	}

	if strings.TrimSpace(userEmail) != "" {
		if err := runGitConfig("user.email", userEmail); err != nil {
			return fmt.Errorf("setting git user.email: %w", err)
		}
		if Program != nil {
			sendFormattedLog("✅", "📦", "Git Identity", "User email set to: "+userEmail)
		}
		logger.LogMessage("SUCCESS", fmt.Sprintf("Git user.email set to: %s", userEmail))
	}

	// Apply Git aliases and defaults
	gitConfigs := map[string]string{
		"alias.co":           "checkout",
		"alias.br":           "branch",
		"alias.ci":           "commit",
		"alias.st":           "status",
		"pull.rebase":        "true",
		"init.defaultBranch": "master",
	}

	if Program != nil {
		sendFormattedLog("🔄", "📦", "Git Aliases", "Setting aliases")
	}

	for key, value := range gitConfigs {
		if err := runGitConfig(key, value); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Failed to set %s: %v", key, err))
		}
	}

	if Program != nil {
		sendFormattedLog("✅", "📦", "Git Setup", "Complete")
	}
	logger.LogMessage("SUCCESS", "Git configuration applied")

	return nil
}

// saveGitCredentials saves git credentials to user.env file
func saveGitCredentials(homeDir, userName, userEmail string) error {
	configDir := filepath.Join(homeDir, ".config", "archriot")
	envFile := filepath.Join(configDir, "user.env")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Read existing file if it exists
	var lines []string
	if content, err := os.ReadFile(envFile); err == nil {
		lines = strings.Split(string(content), "\n")
	}

	// Update or add git credentials
	updated := make(map[string]bool)
	for i, line := range lines {
		if strings.HasPrefix(line, "GIT_USERNAME=") {
			lines[i] = fmt.Sprintf("GIT_USERNAME=%s", userName)
			updated["username"] = true
		} else if strings.HasPrefix(line, "GIT_EMAIL=") {
			lines[i] = fmt.Sprintf("GIT_EMAIL=%s", userEmail)
			updated["email"] = true
		}
	}

	// Add missing credentials
	if !updated["username"] {
		lines = append(lines, fmt.Sprintf("GIT_USERNAME=%s", userName))
	}
	if !updated["email"] {
		lines = append(lines, fmt.Sprintf("GIT_EMAIL=%s", userEmail))
	}

	// Write back to file
	content := strings.Join(lines, "\n")
	return os.WriteFile(envFile, []byte(content), 0644)
}

// runGitConfig sets a git configuration value
func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	return cmd.Run()
}

// runGitConfigGet gets a git configuration value
func runGitConfigGet(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
