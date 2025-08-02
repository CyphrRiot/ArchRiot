package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"archriot-installer/logger"
	"archriot-installer/tui"
)

var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// sendLog sends a formatted log message to TUI
func sendLog(emoji, name, description string) {
	if Program != nil {
		Program.Send(tui.LogMsg(fmt.Sprintf("✅ %s %-15s %s", emoji, name, description)))
	}
}

// enableService enables and starts a systemd service
func enableService(serviceName string) error {
	logger.LogMessage("INFO", fmt.Sprintf("Enabling %s service", serviceName))
	cmd := exec.Command("sudo", "systemctl", "enable", "--now", serviceName+".service")
	if err := cmd.Run(); err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to enable %s: %v", serviceName, err))
		return err
	}
	logger.LogMessage("SUCCESS", fmt.Sprintf("%s service enabled", serviceName))
	return nil
}

// runCommand executes a shell command
func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Run()
}

// setupAudioSystem handles PipeWire installation with conflict resolution
func setupAudioSystem() error {
	logger.LogMessage("INFO", "Setting up PipeWire audio system...")
	sendLog("🔊", "Audio", "Configuring PipeWire")

	// Check for conflicting packages and remove them
	conflictingPackages := []string{"pulseaudio", "jack2"}
	for _, pkg := range conflictingPackages {
		cmd := exec.Command("pacman", "-Qs", "^"+pkg+"$")
		if cmd.Run() == nil {
			logger.LogMessage("INFO", fmt.Sprintf("Removing conflicting package: %s", pkg))
			removeCmd := exec.Command("sudo", "pacman", "-Rdd", "--noconfirm", pkg)
			if err := removeCmd.Run(); err != nil {
				logger.LogMessage("WARNING", fmt.Sprintf("Could not remove %s: %v", pkg, err))
			}
		}
	}

	// Enable PipeWire services for user
	services := []string{"pipewire", "pipewire-pulse", "wireplumber"}
	for _, service := range services {
		cmd := exec.Command("systemctl", "--user", "enable", "--now", service+".service")
		if err := cmd.Run(); err != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Failed to enable %s: %v", service, err))
		}
	}

	logger.LogMessage("SUCCESS", "PipeWire audio system configured")
	return nil
}

// applyMemoryOptimization applies the copied sysctl memory configuration
func applyMemoryOptimization() error {
	logger.LogMessage("INFO", "Applying memory optimization...")
	sendLog("🧠", "Memory", "Applying sysctl config")

	// Apply the copied sysctl configuration
	configPath := "/etc/sysctl.d/99-memory-optimization.conf"
	applyCmd := fmt.Sprintf("sudo sysctl -p %s", configPath)
	if err := runCommand(applyCmd); err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to apply sysctl config: %v", err))
		return err
	}

	logger.LogMessage("SUCCESS", "Memory optimization applied")
	return nil
}

// configureMimeTypes sets up default application associations for file types
func configureMimeTypes() error {
	logger.LogMessage("INFO", "Configuring MIME type associations...")
	sendLog("📄", "MIME Types", "Setting default applications")

	// Update desktop database
	runCommand("update-desktop-database ~/.local/share/applications")

	// Set image viewer (imv) for all image formats
	imageTypes := []string{
		"image/png", "image/jpeg", "image/gif", "image/webp",
		"image/bmp", "image/tiff",
	}
	for _, mimeType := range imageTypes {
		runCommand(fmt.Sprintf("xdg-mime default imv.desktop %s", mimeType))
	}

	// Set PDF viewer
	runCommand("xdg-mime default org.gnome.Papers.desktop application/pdf")

	// Set default browser
	runCommand("xdg-settings set default-web-browser brave-browser.desktop")
	browserSchemes := []string{"x-scheme-handler/http", "x-scheme-handler/https"}
	for _, scheme := range browserSchemes {
		runCommand(fmt.Sprintf("xdg-mime default brave-browser.desktop %s", scheme))
	}

	// Set text editor for text and markdown
	textTypes := []string{
		"text/plain", "text/markdown", "text/x-markdown", "application/x-markdown",
	}
	for _, mimeType := range textTypes {
		runCommand(fmt.Sprintf("xdg-mime default org.gnome.TextEditor.desktop %s", mimeType))
	}

	// Set video player (mpv) for all video formats
	videoTypes := []string{
		"video/mp4", "video/x-msvideo", "video/x-matroska", "video/x-flv",
		"video/x-ms-wmv", "video/mpeg", "video/ogg", "video/webm",
		"video/quicktime", "video/3gpp", "video/3gpp2", "video/x-ms-asf",
		"video/x-ogm+ogg", "video/x-theora+ogg", "application/ogg",
	}
	for _, mimeType := range videoTypes {
		runCommand(fmt.Sprintf("xdg-mime default mpv.desktop %s", mimeType))
	}

	logger.LogMessage("SUCCESS", "MIME type associations configured")
	return nil
}

// setupFishShell handles the complex fish shell configuration
func setupFishShell() error {
	logger.LogMessage("INFO", "Setting up Fish shell...")
	sendLog("🐠", "Fish Shell", "Configuring shell environment")

	homeDir, _ := os.UserHomeDir()

	// Setup fish config
	fishConfigDir := filepath.Join(homeDir, ".config", "fish")
	os.MkdirAll(fishConfigDir, 0755)

	// Copy config if available
	configSource := filepath.Join(homeDir, ".local/share/archriot/config/fish/config.fish")
	configDest := filepath.Join(fishConfigDir, "config.fish")
	if _, err := os.Stat(configSource); err == nil {
		runCommand(fmt.Sprintf("cp %s %s", configSource, configDest))
	}

	// Set as default shell
	user := os.Getenv("USER")
	if user != "" {
		exec.Command("sudo", "chsh", "-s", "/usr/bin/fish", user).Run()
	}

	logger.LogMessage("SUCCESS", "Fish shell configured")
	return nil
}

// Handler registry - simple and clean
var HandlerRegistry = map[string]func() error{
	"setup_base_system": func() error {
		logger.LogMessage("SUCCESS", "Base system ready")
		return nil
	},
	"setup_fish_shell": setupFishShell,
	"refresh_font_cache": func() error {
		runCommand("fc-cache -f")
		sendLog("🎨", "Fonts", "Cache refreshed")
		return nil
	},
	"enable_bluetooth_service": func() error {
		enableService("bluetooth")
		sendLog("📶", "Bluetooth", "Service enabled")
		return nil
	},
	"enable_cups_service": func() error {
		enableService("cups")
		sendLog("🖨️", "Printing", "CUPS enabled")
		return nil
	},
	"setup_power_management": func() error {
		enableService("power-profiles-daemon")
		sendLog("⚡", "Power", "Profiles enabled")
		return nil
	},
	"enable_udisks_service": func() error {
		enableService("udisks2")
		sendLog("💾", "Storage", "Auto-mount enabled")
		return nil
	},
	"setup_wireless_networking": func() error {
		enableService("iwd")
		sendLog("📡", "Wireless", "iwd enabled")
		return nil
	},
	"setup_media_tools": func() error {
		exec.Command("yay", "-S", "--noconfirm", "--needed", "--mflags", "--nocheck", "spotdl").Run()
		sendLog("🎵", "Media", "Tools configured")
		return nil
	},
	"install_migrate_tool": func() error {
		homeDir, _ := os.UserHomeDir()
		binDir := filepath.Join(homeDir, ".local/bin")
		os.MkdirAll(binDir, 0755)

		migrateURL := "https://raw.githubusercontent.com/CyphrRiot/Migrate/main/bin/migrate"
		migratePath := filepath.Join(binDir, "migrate")

		cmd := fmt.Sprintf("curl -L -o %s %s && chmod +x %s", migratePath, migrateURL, migratePath)
		runCommand(cmd)
		sendLog("📦", "Migrate", "Tool installed")
		return nil
	},
	"setup_audio_system": setupAudioSystem,
	"apply_memory_optimization": applyMemoryOptimization,
	"configure_mime_types": configureMimeTypes,
}

// ExecuteHandler executes a handler by name
func ExecuteHandler(handlerName string) error {
	if handler, exists := HandlerRegistry[handlerName]; exists {
		logger.LogMessage("INFO", fmt.Sprintf("Executing handler: %s", handlerName))
		return handler()
	}
	logger.LogMessage("WARNING", fmt.Sprintf("Handler not found: %s", handlerName))
	return fmt.Errorf("handler not found: %s", handlerName)
}
