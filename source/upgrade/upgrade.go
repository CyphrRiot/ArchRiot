package upgrade

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/logger"
	"archriot-installer/tui"
)

// Global variable to track if kernel was upgraded
var kernelUpgraded bool

// Program holds the TUI program reference
var Program *tea.Program

// SetProgram sets the TUI program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// Global channel for upgrade confirmation result
var upgradeConfirmationDone chan bool

// IsKernelUpgraded returns whether kernel was upgraded in last system upgrade
func IsKernelUpgraded() bool {
	return kernelUpgraded
}

// PromptAndRun prompts the user for system upgrade and executes if confirmed
func PromptAndRun() error {
	logger.Log("Progress", "System", "Upgrade Prompt", "Asking user about system upgrade")

	// Initialize upgrade confirmation channel
	upgradeConfirmationDone = make(chan bool, 1)

	// Set up upgrade callback
	tui.SetUpgradeCallback(func(confirmed bool) {
		upgradeConfirmationDone <- confirmed
	})

	// Send upgrade confirmation through existing TUI system
	Program.Send(tui.UpgradeMsg{})

	// Wait for user response through callback
	confirmed := <-upgradeConfirmationDone

	if !confirmed {
		logger.Log("Info", "System", "Upgrade Choice", "User selected NO")
		return nil
	}

	logger.Log("Info", "System", "Upgrade Choice", "User selected YES")

	// Execute the upgrade within TUI
	if err := runSystemUpgrade(); err != nil {
		logger.Log("Warning", "System", "Package Upgrade", "Failed: "+err.Error())
		return fmt.Errorf("system upgrade failed: %v", err)
	}

	logger.Log("Success", "System", "Package Upgrade", "System upgrade completed successfully")
	return nil
}

// runSystemUpgrade executes the full system upgrade process within the TUI
func runSystemUpgrade() error {
	logger.Log("Progress", "System", "Package Upgrade", "Starting comprehensive system upgrade")

	// Reset kernel upgrade flag
	kernelUpgraded = false

	// Step 1: Update package databases
	Program.Send(tui.StepMsg("Updating package databases..."))
	logger.Log("Progress", "System", "Database Update", "Syncing package databases")

	cmd := exec.Command("sudo", "pacman", "-Sy", "--noconfirm")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("package database update failed: %v", err)
	}

	// Check what packages will be upgraded before performing upgrade
	Program.Send(tui.StepMsg("Checking for kernel upgrades..."))
	kernelUpgraded = detectKernelUpgrade()

	// Send kernel upgrade status to TUI
	Program.Send(tui.KernelUpgradeMsg(kernelUpgraded))

	// Step 2: Upgrade official packages with retry logic
	if kernelUpgraded {
		Program.Send(tui.StepMsg("Upgrading official packages (kernel upgrade detected)..."))
		logger.Log("Progress", "System", "Kernel Upgrade", "Linux kernel will be upgraded - reboot will be recommended")
	} else {
		Program.Send(tui.StepMsg("Upgrading official packages..."))
	}
	logger.Log("Progress", "System", "Pacman Upgrade", "Upgrading official packages")

	cmd = exec.Command("sudo", "pacman", "-Su", "--noconfirm")
	if err := cmd.Run(); err != nil {
		// Retry with cache cleaning
		logger.Log("Warning", "System", "Pacman Retry", "First attempt failed, cleaning cache and retrying")
		Program.Send(tui.StepMsg("Cleaning cache and retrying..."))

		// Clean partial downloads
		exec.Command("sudo", "rm", "-f", "/var/cache/pacman/pkg/*.part").Run()

		// Retry upgrade
		cmd = exec.Command("sudo", "pacman", "-Su", "--noconfirm")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("official package upgrade failed after retry: %v", err)
		}
	}

	logger.Log("Success", "System", "Pacman Upgrade", "Official packages upgraded successfully")

	// Step 3: Upgrade AUR packages if yay is available
	if _, err := exec.LookPath("yay"); err == nil {
		Program.Send(tui.StepMsg("Upgrading AUR packages..."))
		logger.Log("Progress", "System", "AUR Upgrade", "Upgrading AUR packages")

		cmd = exec.Command("yay", "-Su", "--noconfirm")

		// Capture both stdout and stderr for detailed error reporting
		combinedOutput, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log("Warning", "System", "AUR Upgrade", "Failed: "+err.Error())
			logger.Log("Warning", "System", "AUR Output", "Command output: "+string(combinedOutput))
			logger.Log("Info", "System", "AUR Manual", "You can manually upgrade AUR packages later with: yay -Su")
			// Continue anyway - AUR failures shouldn't stop the upgrade
		} else {
			logger.Log("Success", "System", "AUR Upgrade", "AUR packages upgraded successfully")
			if len(combinedOutput) > 0 {
				logger.Log("Info", "System", "AUR Output", "Upgrade output: "+string(combinedOutput))
			}
		}
	} else {
		logger.Log("Info", "System", "AUR Upgrade", "yay not available, skipping AUR upgrades")
	}

	// Step 4: Clean orphaned packages
	Program.Send(tui.StepMsg("Cleaning orphaned packages..."))
	logger.Log("Progress", "System", "Orphan Cleanup", "Checking for orphaned packages")

	cmd = exec.Command("pacman", "-Qtdq")
	output, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		orphanList := strings.TrimSpace(string(output))
		logger.Log("Info", "System", "Orphan Cleanup", "Found orphaned packages: "+orphanList)

		// Split package names and pass as arguments instead of stdin
		orphanPackages := strings.Fields(orphanList)
		args := append([]string{"pacman", "-Rns", "--noconfirm"}, orphanPackages...)
		cmd = exec.Command("sudo", args...)

		// Capture both stdout and stderr for detailed error reporting
		combinedOutput, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log("Warning", "System", "Orphan Cleanup", "Failed: "+err.Error())
			logger.Log("Warning", "System", "Orphan Output", "Command output: "+string(combinedOutput))
			logger.Log("Info", "System", "Orphan Manual", "You can manually remove orphans later with: sudo pacman -Rns $(pacman -Qtdq)")
			// Continue anyway - orphan removal failure shouldn't stop upgrade
		} else {
			logger.Log("Success", "System", "Orphan Cleanup", "Orphaned packages removed successfully")
			if len(combinedOutput) > 0 {
				logger.Log("Info", "System", "Orphan Output", "Removal output: "+string(combinedOutput))
			}
		}
	} else {
		if err != nil {
			logger.Log("Info", "System", "Orphan Cleanup", "Cannot check orphans: "+err.Error())
		} else {
			logger.Log("Info", "System", "Orphan Cleanup", "No orphaned packages found")
		}
	}

	// Log final kernel upgrade status
	if kernelUpgraded {
		logger.Log("Success", "System", "Kernel Upgrade", "System upgrade completed - kernel was upgraded, reboot recommended")
	} else {
		logger.Log("Success", "System", "Package Upgrade", "System upgrade completed - no kernel upgrade detected")
	}

	logger.Log("Success", "System", "Package Upgrade", "Complete system upgrade finished successfully")
	return nil
}

// detectKernelUpgrade checks if kernel packages will be upgraded
func detectKernelUpgrade() bool {
	// Get list of packages that would be upgraded
	cmd := exec.Command("pacman", "-Qu")
	output, err := cmd.Output()
	if err != nil {
		logger.Log("Warning", "System", "Kernel Detection", "Could not check for kernel upgrades: "+err.Error())
		return false
	}

	upgradeList := strings.TrimSpace(string(output))
	if upgradeList == "" {
		logger.Log("Info", "System", "Kernel Detection", "No packages to upgrade")
		return false
	}

	lines := strings.Split(upgradeList, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract package name (first field before space)
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		packageName := parts[0]

		// Check for exact "linux" package match
		if packageName == "linux" {
			logger.Log("Info", "System", "Kernel Detection", "Kernel package upgrade detected: "+packageName)
			return true
		}
	}

	logger.Log("Info", "System", "Kernel Detection", "No kernel upgrades detected")
	return false
}
