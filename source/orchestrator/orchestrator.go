package orchestrator

import (
	"crypto/rand"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/config"
	"archriot-installer/executor"
	"archriot-installer/installer"
	"archriot-installer/logger"
	"archriot-installer/plymouth"
	"archriot-installer/tui"
	"archriot-installer/upgrade"
)

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// countTotalModules counts all modules across all categories
func countTotalModules(cfg *config.Config) int {
	total := 0
	total += len(cfg.Core)
	total += len(cfg.System)
	total += len(cfg.Development)
	total += len(cfg.Desktop)
	total += len(cfg.Media)
	total += len(cfg.Utilities)
	total += len(cfg.Productivity)
	total += len(cfg.Specialty)
	total += len(cfg.Theming)
	return total
}

// roundToNearest5 rounds progress to nearest 5%
func roundToNearest5(progress float64) float64 {
	return math.Round(progress*20) / 20
}

// RunInstallation runs the main installation process
func RunInstallation() {
	// Send progress updates to TUI
	sendProgress := func(progress float64) {
		Program.Send(tui.ProgressMsg(roundToNearest5(progress)))
	}

	// Send step updates to TUI
	sendStep := func(step string) {
		Program.Send(tui.StepMsg(step))
	}

	sendStep("Preparing system...")
	logger.Log("Progress", "System", "System Prep", "Preparing system...")
	sendProgress(0.1)

	// Sync package databases first
	if err := installer.SyncPackageDatabases(); err != nil {
		logger.Log("Error", "Database", "Database Sync", "Failed: "+err.Error())
		logger.Log("Info", "System", "Manual Fix", "Please run 'sudo pacman -Sy' manually and try again")
		Program.Send(tui.FailureMsg{Error: "Failed to sync package databases: " + err.Error()})
		return
	}

	sendStep("Loading configuration...")
	sendProgress(0.2)

	// Find config file
	configPath := config.FindConfigFile()
	if configPath == "" {
		logger.Log("Error", "File", "Config", "packages.yaml not found")
		Program.Send(tui.FailureMsg{Error: "Configuration file packages.yaml not found"})
		return
	}

	logger.Log("Progress", "File", "Config Load", "Loading: "+configPath)

	// Load and validate YAML
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log("Error", "File", "Config Load", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Failed to load configuration: " + err.Error()})
		return
	}

	logger.Log("Success", "File", "YAML Config", "Config loaded")

	// Validate YAML configuration
	logger.Log("Progress", "File", "YAML Validation", "Validating configuration...")
	if err := config.ValidateConfig(cfg); err != nil {
		logger.Log("Error", "File", "YAML Validation", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Configuration validation failed: " + err.Error()})
		return
	}
	logger.Log("Success", "File", "YAML Validation", "Configuration validated")

	sendStep("Installing modules...")
	sendProgress(0.3)

	// Set up installer program reference for preservation prompts
	installer.SetProgram(Program)

	// Calculate dynamic progress
	totalModules := countTotalModules(cfg)
	moduleProgressRange := 0.7 // 70% of progress is for modules (30% already used for prep)
	progressPerModule := moduleProgressRange / float64(totalModules)

	// Create progress callback for executor
	completedModules := 0
	progressCallback := func() {
		completedModules++
		currentProgress := 0.3 + (float64(completedModules) * progressPerModule)
		sendProgress(currentProgress)
	}

	// Execute modules in proper order with progress tracking
	if err := executor.ExecuteModulesInOrderWithProgress(cfg, progressCallback); err != nil {
		logger.Log("Error", "System", "Module Exec", "Failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Module execution failed: " + err.Error()})
		return
	}

	// Optional system upgrade before Plymouth installation
	sendStep("Module installation complete!")
	sendProgress(0.90)

	upgrade.SetProgram(Program)
	if err := upgrade.PromptAndRun(); err != nil {
		logger.Log("Warning", "System", "Package Upgrade", "Failed: "+err.Error())
		// Continue anyway - upgrade failure should not stop installation
	}

	// Install Plymouth boot screen (critical system component)
	sendStep("Configuring boot screen...")
	sendProgress(0.95)

	plymouthManager, err := plymouth.NewPlymouthManager()
	if err != nil {
		logger.Log("Error", "System", "Plymouth", "Failed to initialize: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth initialization failed: " + err.Error()})
		return
	}

	sendProgress(0.96)
	if err := plymouthManager.InstallPlymouth(); err != nil {
		logger.Log("Error", "System", "Plymouth", "Installation failed: "+err.Error())
		Program.Send(tui.FailureMsg{Error: "Plymouth installation failed: " + err.Error()})
		return
	}

	sendProgress(0.98)
	sendStep("Checking Secure Boot recommendation...")

	// Check Secure Boot recommendation before completion
	checkSecureBootRecommendation()

	sendStep("Installation complete!")
	sendProgress(1.0)
	logger.Log("Success", "System", "Installation", "Complete!")
	logger.Log("Success", "System", "Module Exec", "All modules done")
	logger.Log("Info", "System", "Log File", "Available at: "+logger.GetLogPath())

	// Ensure Mullvad VPN connectivity (if installed and currently disconnected)
	refreshMullvadIfActive()

	// Explicit success banner in install.log
	logger.LogMessage("SUCCESS", "Installation completed")

	// Send success completion message
	Program.Send(tui.DoneMsg{})
}

// refreshMullvadIfActive checks Mullvad status and reconnects if not connected.
func refreshMullvadIfActive() {
	if _, err := exec.LookPath("mullvad"); err != nil {
		return
	}
	out, err := exec.Command("mullvad", "status").CombinedOutput()
	if err != nil {
		logger.Log("Warning", "VPN", "Mullvad", "status failed: "+err.Error())
		return
	}
	s := strings.ToLower(string(out))
	if strings.Contains(s, "connected") {
		// Already connected; avoid disrupting an active session.
		logger.Log("Info", "VPN", "Mullvad", "Already connected; no refresh needed")
		return
	}
	// Attempt to reconnect when disconnected
	Program.Send(tui.LogMsg("🔐 Refreshing Mullvad VPN connection..."))
	_ = exec.Command("mullvad", "reconnect").Run()
	if out2, err2 := exec.Command("mullvad", "status").CombinedOutput(); err2 == nil && strings.Contains(strings.ToLower(string(out2)), "connected") {
		logger.Log("Success", "VPN", "Mullvad", "Connected")
	} else {
		logger.Log("Warning", "VPN", "Mullvad", "Still not connected after refresh")
	}
}

// secureBootSetupDone channel for synchronization
var secureBootSetupDone chan bool

// SetSecureBootSetupChannel sets the channel for Secure Boot setup synchronization
func SetSecureBootSetupChannel(ch chan bool) {
	secureBootSetupDone = ch
}

// checkSecureBootRecommendation checks if Secure Boot should be recommended and waits for user decision
func checkSecureBootRecommendation() {
	logger.Log("Info", "System", "SecureBoot", "Checking Secure Boot and LUKS status...")

	// Detect Secure Boot status
	sbEnabled, sbSupported, err := installer.DetectSecureBootStatus()
	if err != nil {
		logger.Log("Warning", "System", "SecureBoot", "Failed to detect Secure Boot status: "+err.Error())
		return
	}

	// Detect LUKS encryption
	luksUsed, luksDevices, err := installer.DetectLuksEncryption()
	if err != nil {
		logger.Log("Warning", "System", "SecureBoot", "Failed to detect LUKS encryption: "+err.Error())
		luksUsed = false
		luksDevices = []string{}
	}

	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("Secure Boot: enabled=%v, supported=%v", sbEnabled, sbSupported))
	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("LUKS: detected=%v, devices=%v", luksUsed, luksDevices))

	// Send status to TUI
	Program.Send(tui.SecureBootStatusMsg{
		Enabled:     sbEnabled,
		Supported:   sbSupported,
		LuksUsed:    luksUsed,
		LuksDevices: luksDevices,
	})

	// TODO: Secure Boot prompting disabled - bootloader signing issues cause boot failures
	if false && !sbEnabled && sbSupported && luksUsed {
		logger.Log("Info", "System", "SecureBoot", "Prompting user for Secure Boot enablement...")
		Program.Send(tui.SecureBootPromptMsg{})

		// Wait for user decision
		if secureBootSetupDone != nil {
			userWantsSecureBoot := <-secureBootSetupDone
			if userWantsSecureBoot {
				Program.Send(tui.StepMsg("Setting up Secure Boot..."))
				if err := performSecureBootSetup(); err != nil {
					logger.Log("Error", "System", "SecureBoot", "Setup failed: "+err.Error())
					Program.Send(tui.FailureMsg{Error: "Secure Boot setup failed: " + err.Error()})
					return
				}
				Program.Send(tui.LogMsg("✅ Secure Boot preparation complete"))
				logger.Log("Success", "System", "SecureBoot", "Setup completed successfully")
			} else {
				logger.Log("Info", "System", "SecureBoot", "User declined Secure Boot setup")
			}
		}
	}
}

// performSecureBootSetup executes the Secure Boot setup process
func performSecureBootSetup() error {
	Program.Send(tui.LogMsg("🛡️ Preparing Secure Boot configuration..."))
	Program.Send(tui.LogMsg("📋 Validating system prerequisites..."))

	// Validate prerequisites
	if err := installer.ValidateSecureBootPrerequisites(); err != nil {
		Program.Send(tui.LogMsg("❌ Prerequisites validation failed"))
		return fmt.Errorf("prerequisites not met: %w", err)
	}
	Program.Send(tui.LogMsg("✅ System prerequisites validated"))

	Program.Send(tui.LogMsg("⚙️ Configuring post-reboot continuation..."))
	// Modify hyprland.conf for continuation after reboot
	if err := modifyHyprlandForContinuation(); err != nil {
		Program.Send(tui.LogMsg("❌ Failed to configure continuation"))
		return fmt.Errorf("failed to setup continuation: %w", err)
	}
	Program.Send(tui.LogMsg("✅ Post-reboot continuation configured"))

	// Phase 2: Generate custom Secure Boot keys
	Program.Send(tui.LogMsg("🔑 Generating Secure Boot signing keys..."))
	if err := generateSecureBootKeys(); err != nil {
		Program.Send(tui.LogMsg("❌ Failed to generate signing keys"))
		return fmt.Errorf("key generation failed: %w", err)
	}
	Program.Send(tui.LogMsg("✅ Secure Boot keys generated"))

	// Phase 2: Sign bootloader and kernel
	Program.Send(tui.LogMsg("🖊️ Signing bootloader and kernel..."))
	if err := signBootComponents(); err != nil {
		Program.Send(tui.LogMsg("❌ Failed to sign boot components"))
		return fmt.Errorf("signing failed: %w", err)
	}
	Program.Send(tui.LogMsg("✅ Boot components signed"))

	// Phase 2: Set up automatic signing for future updates
	Program.Send(tui.LogMsg("⚙️ Setting up automatic kernel signing..."))
	if err := setupPackmanHooks(); err != nil {
		Program.Send(tui.LogMsg("⚠️ Warning: Failed to setup automatic signing"))
		logger.Log("Warning", "System", "SecureBoot", "Pacman hooks setup failed: "+err.Error())
	} else {
		Program.Send(tui.LogMsg("✅ Automatic signing configured"))
	}

	// Phase 2: Prepare key installation guidance
	Program.Send(tui.LogMsg("📋 Preparing UEFI key installation..."))
	if err := prepareKeyInstallation(); err != nil {
		Program.Send(tui.LogMsg("⚠️ Warning: Failed to prepare key installation"))
		logger.Log("Warning", "System", "SecureBoot", "Key preparation failed: "+err.Error())
	} else {
		Program.Send(tui.LogMsg("✅ Key installation prepared"))
	}

	Program.Send(tui.LogMsg(""))
	Program.Send(tui.LogMsg("🔄 WHAT HAPPENS NEXT:"))
	Program.Send(tui.LogMsg("• After reboot, you'll see setup continuation"))
	Program.Send(tui.LogMsg("• Install custom keys in UEFI settings"))
	Program.Send(tui.LogMsg("• Enable Secure Boot with custom keys"))
	Program.Send(tui.LogMsg("• System will validate and complete setup"))
	Program.Send(tui.LogMsg(""))

	logger.Log("Info", "System", "SecureBoot", "Secure Boot setup completed - install keys and enable in UEFI settings after reboot")
	return nil
}

// generateSecureBootKeys creates the Secure Boot key hierarchy using sbctl
func generateSecureBootKeys() error {
	logger.Log("Info", "System", "SecureBoot", "Creating Secure Boot key hierarchy with sbctl...")

	// Check if sbctl is available
	if _, err := exec.LookPath("sbctl"); err != nil {
		return fmt.Errorf("sbctl not found - install 'sbctl' package: %w", err)
	}

	// Use sbctl to create keys (stores in /usr/share/secureboot/keys/ by default)
	cmd := exec.Command("sbctl", "create-keys")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sbctl create-keys failed: %w - output: %s", err, string(output))
	}

	logger.Log("Success", "System", "SecureBoot", "Secure Boot key hierarchy created with sbctl")
	logger.Log("Info", "System", "SecureBoot", string(output))
	return nil
}

// signBootComponents signs the bootloader and kernel using sbctl
func signBootComponents() error {
	logger.Log("Info", "System", "SecureBoot", "Signing boot components with sbctl...")

	// Check if sbctl is available
	if _, err := exec.LookPath("sbctl"); err != nil {
		return fmt.Errorf("sbctl not found - install 'sbctl' package: %w", err)
	}

	// Sign systemd-boot bootloader
	bootloaderPaths := []string{
		"/boot/EFI/systemd/systemd-bootx64.efi",
		"/boot/EFI/BOOT/BOOTX64.EFI",
	}

	for _, bootloaderPath := range bootloaderPaths {
		if _, err := os.Stat(bootloaderPath); err == nil {
			if err := signFileWithSbctl(bootloaderPath); err != nil {
				logger.Log("Warning", "System", "SecureBoot", fmt.Sprintf("Failed to sign %s: %v", bootloaderPath, err))
			} else {
				logger.Log("Success", "System", "SecureBoot", fmt.Sprintf("Signed bootloader: %s", bootloaderPath))
			}
		}
	}

	// Sign kernel
	kernelPath := "/boot/vmlinuz-linux"
	if _, err := os.Stat(kernelPath); err == nil {
		if err := signFileWithSbctl(kernelPath); err != nil {
			logger.Log("Warning", "System", "SecureBoot", fmt.Sprintf("Failed to sign kernel: %v", err))
		} else {
			logger.Log("Success", "System", "SecureBoot", "Signed kernel")
		}
	}

	logger.Log("Success", "System", "SecureBoot", "Boot component signing completed")
	return nil
}

// signFileWithSbctl signs a single file using sbctl
func signFileWithSbctl(filePath string) error {
	// Use sbctl sign with save option to add to database for auto-signing
	cmd := exec.Command("sbctl", "sign", "-s", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sbctl sign failed for %s: %w - output: %s", filePath, err, string(output))
	}

	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("sbctl output: %s", string(output)))
	return nil
}

// setupPackmanHooks creates pacman hooks for automatic kernel signing using sbctl
func setupPackmanHooks() error {
	hookDir := "/etc/pacman.d/hooks"
	if err := os.MkdirAll(hookDir, 0755); err != nil {
		return fmt.Errorf("creating hooks directory: %w", err)
	}

	hookPath := filepath.Join(hookDir, "99-secureboot.hook")

	// Comprehensive hook for all components that need signing
	hookContent := `[Trigger]
Operation = Install
Operation = Upgrade
Type = Package
Target = linux
Target = linux-lts
Target = linux-zen
Target = linux-hardened
Target = systemd
Target = plymouth
Target = mkinitcpio

[Action]
Description = Signing kernel, initramfs and boot loader for Secure Boot
When = PostTransaction
Exec = /usr/bin/sbctl sign-all
`

	if err := os.WriteFile(hookPath, []byte(hookContent), 0644); err != nil {
		return fmt.Errorf("writing pacman hook: %w", err)
	}

	logger.Log("Success", "System", "SecureBoot", "Pacman hooks configured for automatic signing with sbctl")
	return nil
}

// prepareKeyInstallation prepares keys for UEFI installation
func prepareKeyInstallation() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	keyDir := filepath.Join(homeDir, ".config", "archriot", "secureboot", "keys")
	eslDir := filepath.Join(homeDir, ".config", "archriot", "secureboot", "esl")

	if err := os.MkdirAll(eslDir, 0755); err != nil {
		return fmt.Errorf("creating ESL directory: %w", err)
	}

	logger.Log("Info", "System", "SecureBoot", "Preparing keys for UEFI installation...")

	// Convert certificates to EFI signature lists
	keys := []string{"PK", "KEK", "db"}
	for _, keyName := range keys {
		certPath := filepath.Join(keyDir, keyName+".crt")
		eslPath := filepath.Join(eslDir, keyName+".esl")
		authPath := filepath.Join(eslDir, keyName+".auth")

		// Generate proper UUID for this key
		uuid, err := generateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for %s: %w", keyName, err)
		}

		// Create ESL from certificate with proper UUID
		cmd := exec.Command("cert-to-efi-sig-list", "-g", uuid, certPath, eslPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("creating ESL for %s: %w", keyName, err)
		}

		// Create authenticated variable for UEFI
		if keyName == "PK" {
			// PK is self-signed
			cmd = exec.Command("sign-efi-sig-list", "-k", filepath.Join(keyDir, "PK.key"), "-c", filepath.Join(keyDir, "PK.crt"), keyName, eslPath, authPath)
		} else {
			// KEK and db are signed by PK
			cmd = exec.Command("sign-efi-sig-list", "-k", filepath.Join(keyDir, "PK.key"), "-c", filepath.Join(keyDir, "PK.crt"), keyName, eslPath, authPath)
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("signing %s: %w", keyName, err)
		}

		logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("Prepared %s for UEFI installation", keyName))
	}

	// Create installation script
	scriptPath := filepath.Join(homeDir, ".local", "share", "archriot", "secureboot", "install-keys.sh")
	scriptContent := fmt.Sprintf(`#!/bin/bash
# ArchRiot Secure Boot Key Installation Script

ESL_DIR="%s"

echo "Installing ArchRiot Secure Boot keys..."

# Install keys using efi-updatevar
sudo efi-updatevar -f "$ESL_DIR/db.auth" db
sudo efi-updatevar -f "$ESL_DIR/KEK.auth" KEK
sudo efi-updatevar -f "$ESL_DIR/PK.auth" PK

echo "Keys installed. You can now enable Secure Boot in UEFI settings."
`, eslDir)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("creating installation script: %w", err)
	}

	logger.Log("Success", "System", "SecureBoot", "Key installation prepared")
	logger.Log("Info", "System", "SecureBoot", fmt.Sprintf("Installation script: %s", scriptPath))
	return nil
}

// generateUUID creates a proper UUID for UEFI key identification
func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant bits according to RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// modifyHyprlandForContinuation modifies hyprland.conf to launch Secure Boot continuation after reboot
func modifyHyprlandForContinuation() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	hyprlandConfigPath := filepath.Join(homeDir, ".config", "hypr", "hyprland.conf")
	backupPath := hyprlandConfigPath + ".archriot-backup"

	// Read current hyprland.conf
	configData, err := os.ReadFile(hyprlandConfigPath)
	if err != nil {
		return fmt.Errorf("reading hyprland.conf: %w", err)
	}

	// Backup original
	if err := os.WriteFile(backupPath, configData, 0644); err != nil {
		return fmt.Errorf("creating backup: %w", err)
	}

	// Get ArchRiot binary path
	archRiotPath := filepath.Join(homeDir, ".local", "share", "archriot", "install", "archriot")

	// Modify config to replace welcome with Secure Boot continuation in Ghostty
	configStr := string(configData)
	lines := strings.Split(configStr, "\n")

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "exec-once") && strings.Contains(trimmedLine, "welcome") {
			// Replace with Secure Boot continuation in terminal
			terminalCmd := fmt.Sprintf("ghostty --class=com.mitchellh.ghostty-floating -e bash -c 'cd $HOME; %s --secure_boot_stage; echo \"Press Enter to close.\"; read'", archRiotPath)
			lines[i] = fmt.Sprintf("exec-once = sleep 2 && %s", terminalCmd)
			logger.Log("Info", "System", "SecureBoot", "Modified hyprland.conf for Secure Boot continuation")
			break
		}
	}

	// Write modified config
	modifiedConfig := strings.Join(lines, "\n")
	if err := os.WriteFile(hyprlandConfigPath, []byte(modifiedConfig), 0644); err != nil {
		// Restore backup on failure
		os.WriteFile(hyprlandConfigPath, configData, 0644)
		return fmt.Errorf("writing modified config: %w", err)
	}

	logger.Log("Success", "System", "SecureBoot", "Hyprland configuration modified for post-reboot continuation")
	return nil
}
