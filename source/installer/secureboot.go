package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"archriot-installer/logger"
)

// SecureBootStatus represents the current Secure Boot state
type SecureBootStatus struct {
	Enabled   bool
	Supported bool
	SetupMode bool
}

// LuksInfo represents detected LUKS encryption information
type LuksInfo struct {
	Detected bool
	Devices  []string
	Methods  []string // Which detection methods found LUKS
}

// DetectSecureBootStatus checks the current Secure Boot state and Setup Mode
func DetectSecureBootStatus() (bool, bool, error) {
	logger.LogMessage("INFO", "🔍 Detecting Secure Boot status...")

	// Check if system supports UEFI
	if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
		logger.LogMessage("INFO", "Legacy BIOS system detected - Secure Boot not supported")
		return false, false, nil
	}

	// Check current Secure Boot state using bootctl
	cmd := exec.Command("bootctl", "status")
	output, err := cmd.Output()
	if err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to run bootctl status: %v", err))
		// Fallback to efivars method
		return detectSecureBootFromEfivars()
	}

	outputStr := string(output)
	enabled := strings.Contains(outputStr, "Secure Boot: enabled")
	supported := true // UEFI system with bootctl means Secure Boot is supported

	// Check if system is in Setup Mode (required for key installation)
	setupMode := detectSetupMode()
	if setupMode {
		logger.LogMessage("INFO", "UEFI is in Setup Mode - ready for custom key installation")
	} else {
		logger.LogMessage("WARNING", "UEFI is NOT in Setup Mode - custom keys cannot be installed")
	}

	logger.LogMessage("INFO", fmt.Sprintf("Secure Boot status: enabled=%v, supported=%v, setup_mode=%v", enabled, supported, setupMode))
	return enabled, supported, nil
}

// detectSecureBootFromEfivars fallback method using EFI variables
func detectSecureBootFromEfivars() (bool, bool, error) {
	logger.LogMessage("INFO", "Using efivars fallback for Secure Boot detection...")

	// Check if SecureBoot variable exists and is enabled
	sbPath := "/sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c"
	if sbData, err := os.ReadFile(sbPath); err == nil {
		// EFI variable format: first 4 bytes are attributes, then data
		if len(sbData) > 4 {
			enabled := sbData[4] == 1
			logger.LogMessage("INFO", fmt.Sprintf("SecureBoot efivar: enabled=%v", enabled))
			return enabled, true, nil
		}
	}

	// If we can't read the variable, assume supported but not enabled
	logger.LogMessage("WARNING", "Could not read SecureBoot efivar, assuming supported but disabled")
	return false, true, nil
}

// DetectLuksEncryption checks for LUKS encryption using multiple methods
func DetectLuksEncryption() (bool, []string, error) {
	logger.LogMessage("INFO", "🔍 Detecting LUKS encryption...")

	var luksDevices []string
	var detectionMethods []string

	// Method 1: Check /proc/cmdline for cryptdevice= parameters
	if devices := detectLuksFromCmdline(); len(devices) > 0 {
		luksDevices = append(luksDevices, devices...)
		detectionMethods = append(detectionMethods, "cmdline")
		logger.LogMessage("INFO", fmt.Sprintf("LUKS detected in cmdline: %v", devices))
	}

	// Method 2: Check /etc/crypttab for LUKS entries
	if devices := detectLuksFromCrypttab(); len(devices) > 0 {
		luksDevices = append(luksDevices, devices...)
		detectionMethods = append(detectionMethods, "crypttab")
		logger.LogMessage("INFO", fmt.Sprintf("LUKS detected in crypttab: %v", devices))
	}

	// Method 3: Check active device mapper for crypt devices
	if devices := detectLuksFromDmsetup(); len(devices) > 0 {
		luksDevices = append(luksDevices, devices...)
		detectionMethods = append(detectionMethods, "dmsetup")
		logger.LogMessage("INFO", fmt.Sprintf("LUKS detected via dmsetup: %v", devices))
	}

	// Method 4: Use lsblk to find encrypted partitions
	if devices := detectLuksFromLsblk(); len(devices) > 0 {
		luksDevices = append(luksDevices, devices...)
		detectionMethods = append(detectionMethods, "lsblk")
		logger.LogMessage("INFO", fmt.Sprintf("LUKS detected via lsblk: %v", devices))
	}

	// Remove duplicates
	uniqueDevices := removeDuplicates(luksDevices)
	detected := len(uniqueDevices) > 0

	if detected {
		logger.LogMessage("SUCCESS", fmt.Sprintf("LUKS encryption detected: %d devices found via %v",
			len(uniqueDevices), detectionMethods))
	} else {
		logger.LogMessage("INFO", "No LUKS encryption detected")
	}

	return detected, uniqueDevices, nil
}

// detectLuksFromCmdline checks /proc/cmdline for LUKS parameters
func detectLuksFromCmdline() []string {
	var devices []string

	cmdline, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return devices
	}

	cmdlineStr := string(cmdline)

	// Look for cryptdevice= parameters
	if strings.Contains(cmdlineStr, "cryptdevice=") {
		// Parse cryptdevice=UUID=xxx:name or cryptdevice=/dev/xxx:name
		for _, param := range strings.Fields(cmdlineStr) {
			if strings.HasPrefix(param, "cryptdevice=") {
				parts := strings.Split(param, ":")
				if len(parts) >= 2 {
					devices = append(devices, parts[1]) // mapper name
				}
			}
		}
	}

	// Look for rd.luks.uuid= parameters (dracut style)
	if strings.Contains(cmdlineStr, "rd.luks.uuid=") {
		for _, param := range strings.Fields(cmdlineStr) {
			if strings.HasPrefix(param, "rd.luks.uuid=") {
				uuid := strings.TrimPrefix(param, "rd.luks.uuid=")
				devices = append(devices, "luks-"+uuid)
			}
		}
	}

	return devices
}

// detectLuksFromCrypttab checks /etc/crypttab for LUKS entries
func detectLuksFromCrypttab() []string {
	var devices []string

	crypttab, err := os.ReadFile("/etc/crypttab")
	if err != nil {
		return devices
	}

	lines := strings.Split(string(crypttab), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				devices = append(devices, parts[0]) // mapper name
			}
		}
	}

	return devices
}

// detectLuksFromDmsetup checks active device mapper for crypt devices
func detectLuksFromDmsetup() []string {
	var devices []string

	cmd := exec.Command("dmsetup", "ls", "--target", "crypt")
	output, err := cmd.Output()
	if err != nil {
		return devices
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "No devices found" {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				// Remove device mapper formatting like (253:0)
				deviceName := parts[0]
				devices = append(devices, deviceName)
			}
		}
	}

	return devices
}

// detectLuksFromLsblk uses lsblk to find encrypted partitions
func detectLuksFromLsblk() []string {
	var devices []string

	cmd := exec.Command("lsblk", "-f", "-o", "NAME,FSTYPE", "--noheadings")
	output, err := cmd.Output()
	if err != nil {
		return devices
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "crypto_LUKS") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				deviceName := strings.TrimPrefix(parts[0], "├─")
				deviceName = strings.TrimPrefix(deviceName, "└─")
				deviceName = strings.TrimPrefix(deviceName, "│ ")
				deviceName = strings.TrimSpace(deviceName)
				if deviceName != "" {
					devices = append(devices, deviceName)
				}
			}
		}
	}

	return devices
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(items []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			result = append(result, item)
		}
	}

	return result
}

// CheckSecureBootRecommendation determines if Secure Boot should be recommended
func CheckSecureBootRecommendation() (bool, string) {
	logger.LogMessage("INFO", "🔍 Checking Secure Boot recommendation...")

	// Check Secure Boot status
	sbEnabled, sbSupported, err := DetectSecureBootStatus()
	if err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to detect Secure Boot status: %v", err))
		return false, "Could not detect Secure Boot status"
	}

	// Check LUKS encryption
	luksDetected, luksDevices, err := DetectLuksEncryption()
	if err != nil {
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to detect LUKS encryption: %v", err))
		return false, "Could not detect LUKS encryption status"
	}

	// Recommendation logic
	if !sbSupported {
		return false, "Secure Boot not supported (Legacy BIOS system)"
	}

	if sbEnabled {
		return false, "Secure Boot already enabled"
	}

	if !luksDetected {
		return false, "No LUKS encryption detected - Secure Boot not critical"
	}

	// Secure Boot recommended: UEFI system with LUKS but no Secure Boot
	deviceList := strings.Join(luksDevices, ", ")
	reason := fmt.Sprintf("LUKS encryption detected (%s) but Secure Boot disabled - recommended for memory attack protection", deviceList)

	logger.LogMessage("INFO", "Secure Boot enablement recommended")
	return true, reason
}

// GetSecureBootEducation returns educational information about Secure Boot + LUKS
func GetSecureBootEducation(luksDevices []string) string {
	deviceList := "None detected"
	if len(luksDevices) > 0 {
		deviceList = strings.Join(luksDevices, ", ")
	}

	return fmt.Sprintf(`🛡️ SECURE BOOT + LUKS PROTECTION

WHY ENABLE SECURE BOOT?
• Prevents cold boot attacks on LUKS encryption keys
• Validates bootloader and kernel signatures before execution
• Protects against evil maid attacks and firmware tampering
• Industry security best practice for encrypted systems

DETECTED LUKS DEVICES: %s

WHAT WILL HAPPEN?
• Generate custom Secure Boot signing keys (PK, KEK, db)
• Sign bootloader (GRUB/systemd-boot) and kernel components
• Set up automatic signing for future kernel updates
• Guide you through UEFI firmware configuration
• Verify complete chain of trust

REQUIREMENTS:
• UEFI firmware (detected: ✅)
• LUKS encryption (detected: ✅)
• Administrative access to UEFI settings
• System reboot required

⚠️  IMPORTANT: You'll need to:
1. Reboot after setup completion
2. Enter UEFI/BIOS settings
3. Enable Secure Boot with custom keys
4. Clear existing Secure Boot database (if any)

This process is REVERSIBLE - you can disable Secure Boot anytime.`, deviceList)
}

// ValidateSecureBootPrerequisites checks if system meets requirements
func ValidateSecureBootPrerequisites() error {
	logger.LogMessage("INFO", "🔍 Validating Secure Boot prerequisites...")

	// Check UEFI mode
	if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
		return fmt.Errorf("system is not running in UEFI mode")
	}
	logger.LogMessage("SUCCESS", "✓ UEFI mode detected")

	// Check if running as root or with sudo
	if os.Geteuid() != 0 {
		// Test sudo access
		cmd := exec.Command("sudo", "-n", "true")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("root privileges required for Secure Boot setup")
		}
	}
	logger.LogMessage("SUCCESS", "✓ Administrative privileges confirmed")

	// Check available disk space in /boot
	cmd := exec.Command("df", "-BM", "/boot")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check /boot disk space: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 && strings.Contains(lines[1], "0M") {
		return fmt.Errorf("/boot partition has insufficient free space")
	}
	logger.LogMessage("SUCCESS", "✓ Sufficient disk space available")

	// Check if required tools are available
	requiredTools := []string{"openssl", "mokutil", "bootctl"}
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("required tool not found: %s", tool)
		}
	}
	logger.LogMessage("SUCCESS", "✓ Required tools available")

	// CRITICAL: Check if UEFI is in Setup Mode for key installation
	if !detectSetupMode() {
		return fmt.Errorf("UEFI must be in Setup Mode to install custom keys - clear Secure Boot keys in BIOS first")
	}
	logger.LogMessage("SUCCESS", "✓ UEFI Setup Mode confirmed - ready for key installation")

	return nil
}

// detectSetupMode checks if UEFI is in Setup Mode (required for custom key installation)
func detectSetupMode() bool {
	// Check SetupMode EFI variable
	setupPath := "/sys/firmware/efi/efivars/SetupMode-8be4df61-93ca-11d2-aa0d-00e098032b8c"
	if setupData, err := os.ReadFile(setupPath); err == nil {
		// EFI variable format: first 4 bytes are attributes, then data
		if len(setupData) > 4 {
			setupMode := setupData[4] == 1
			return setupMode
		}
	}
	return false
}

// GetSetupModeInstructions returns user instructions for enabling Setup Mode
func GetSetupModeInstructions() string {
	return `🔧 UEFI SETUP MODE REQUIRED

Your UEFI firmware must be in "Setup Mode" to install custom Secure Boot keys.

STEPS TO ENABLE SETUP MODE:

1. Reboot and enter UEFI/BIOS settings:
   • Dell: Press F2 during boot
   • HP: Press F10 or ESC during boot
   • Lenovo: Press F1 or F2 during boot
   • ASUS: Press F2 or DEL during boot

2. Navigate to Security → Secure Boot settings

3. Look for one of these options:
   • "Clear Secure Boot Keys" → Execute
   • "Reset to Setup Mode" → Enable
   • "Secure Boot Mode" → Set to "Custom"
   • "Key Management" → "Clear Keys"

4. Save settings and exit UEFI

5. Boot back to Linux and run ArchRiot again

⚠️  WARNING: This will clear existing Secure Boot keys (including Microsoft keys).
Your system will boot normally, but Windows may show warnings until custom keys are installed.

✅ BENEFIT: After setup, your system will have stronger security than default Microsoft keys.`
}
