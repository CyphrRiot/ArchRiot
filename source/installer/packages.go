package installer

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"archriot-installer/logger"
)

var program *tea.Program

// SetProgram sets the TUI program reference for sending messages
func SetProgram(p *tea.Program) {
	program = p
}

// InstallPackages installs packages in a single batch to avoid database locks
func InstallPackages(packages []string) error {
	if len(packages) == 0 {
		logger.Log("Info", "Package", "Packages", "None to install")
		return nil
	}

	logger.LogMessage("INFO", fmt.Sprintf("Installing %d packages", len(packages)))

	// Check which packages are already installed
	var toInstall []string
	for _, pkg := range packages {
		if !isPackageInstalled(pkg) {
			toInstall = append(toInstall, pkg)
		}
	}

	if len(toInstall) == 0 {
		logger.LogMessage("INFO", "All packages already installed")
		return nil
	}

	logger.LogMessage("INFO", fmt.Sprintf("Installing %d new packages", len(toInstall)))

	// Log which packages we're about to install
	for i, pkg := range toInstall {
		logger.LogMessage("INFO", fmt.Sprintf("📦 Package %d/%d: %s", i+1, len(toInstall), pkg))
	}

	// Try to install all packages at once with pacman
	start := time.Now()
	logger.LogMessage("INFO", "🔄 Attempting installation with pacman...")
	cmd := exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm", "--needed"}, toInstall...)...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Show detailed pacman failure info
		outputStr := string(output)
		if len(outputStr) > 300 {
			outputStr = outputStr[:300] + "... (truncated)"
		}
		logger.LogMessage("WARNING", fmt.Sprintf("Pacman failed for packages: %v", toInstall))
		logger.LogMessage("WARNING", fmt.Sprintf("Pacman error: %s", outputStr))

		// Try ONE retry with fresh database sync before falling back to yay
		logger.LogMessage("INFO", "🔄 Refreshing databases and retrying pacman once...")
		if syncErr := SyncPackageDatabases(); syncErr != nil {
			logger.LogMessage("WARNING", fmt.Sprintf("Database refresh failed: %v", syncErr))
		}

		// Single retry attempt
		retryCmd := exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm", "--needed"}, toInstall...)...)
		retryOutput, retryErr := retryCmd.CombinedOutput()

		if retryErr == nil {
			logger.LogMessage("SUCCESS", "✅ Retry successful! Packages installed with pacman")
			duration := time.Since(start)
			logger.LogMessage("SUCCESS", fmt.Sprintf("✅ Package installation completed in %v", duration))
			for _, pkg := range toInstall {
				logger.LogMessage("SUCCESS", fmt.Sprintf("✓ Installed: %s", pkg))
			}
			return nil
		}

		// Retry also failed, log it and continue to yay fallback
		retryOutputStr := string(retryOutput)
		if len(retryOutputStr) > 300 {
			retryOutputStr = retryOutputStr[:300] + "... (truncated)"
		}
		logger.LogMessage("WARNING", fmt.Sprintf("Pacman retry also failed: %s", retryOutputStr))
		logger.LogMessage("INFO", "🔄 Switching to yay for AUR package support...")

		// Identify likely AUR packages
		logger.LogMessage("INFO", "🔍 These packages may be from AUR (requiring compilation):")
		for _, pkg := range toInstall {
			if containsAURIndicators(pkg) {
				logger.LogMessage("INFO", fmt.Sprintf("   📦 %s (likely AUR - may take several minutes to build)", pkg))
			} else {
				logger.LogMessage("INFO", fmt.Sprintf("   📦 %s", pkg))
			}
		}

		// Check if yay is available
		if !CommandExists("yay") {
			logger.LogMessage("ERROR", "Yay not available for AUR packages")
			return fmt.Errorf("package installation failed: pacman failed and yay not available")
		}

		logger.LogMessage("INFO", "🔨 Starting yay installation (AUR packages compile from source - please wait)...")
		logger.LogMessage("INFO", "⏳ This may take 5-30 minutes depending on packages and system speed")
		cmd = exec.Command("yay", append([]string{"-S", "--noconfirm", "--needed"}, toInstall...)...)
		output, err = cmd.CombinedOutput()

		if err != nil {
			// Check if it's a network error and retry once
			outputStr := string(output)
			if strings.Contains(outputStr, "failure in name resolution") || strings.Contains(outputStr, "dial tcp") || strings.Contains(outputStr, "network") {
				logger.LogMessage("WARNING", "Network error detected, retrying yay once...")
				time.Sleep(3 * time.Second)

				retryCmd := exec.Command("yay", append([]string{"-S", "--noconfirm", "--needed"}, toInstall...)...)
				retryOutput, retryErr := retryCmd.CombinedOutput()

				if retryErr == nil {
					logger.LogMessage("SUCCESS", "✅ Yay retry successful!")
					output = retryOutput
					err = nil
				} else {
					retryOutputStr := string(retryOutput)
					if len(retryOutputStr) > 300 {
						retryOutputStr = retryOutputStr[:300] + "... (truncated)"
					}
					logger.LogMessage("WARNING", fmt.Sprintf("Yay retry also failed: %s", retryOutputStr))
				}
			}
		}

		if err != nil {
			// Log the error and fail for critical packages
			outputStr := string(output)
			if len(outputStr) > 500 {
				outputStr = outputStr[:500] + "... (truncated)"
			}
			logger.LogMessage("ERROR", fmt.Sprintf("Package installation failed: %s", outputStr))
			logger.Log("Error", "Package", "Installation", "Critical packages failed to install")
			return fmt.Errorf("package installation failed: %w", err)
		}

		// CRITICAL: Verify packages actually got installed
		// yay can return exit code 0 even when packages don't exist ("nothing to do")
		var failedPackages []string
		for _, pkg := range toInstall {
			if !isPackageInstalled(pkg) {
				failedPackages = append(failedPackages, pkg)
			}
		}

		if len(failedPackages) > 0 {
			outputStr := string(output)
			if len(outputStr) > 500 {
				outputStr = outputStr[:500] + "... (truncated)"
			}
			logger.LogMessage("ERROR", fmt.Sprintf("Package installation failed - packages not found: %v", failedPackages))
			logger.LogMessage("ERROR", fmt.Sprintf("yay output: %s", outputStr))
			logger.Log("Error", "Package", "Installation", fmt.Sprintf("Packages not found: %v", failedPackages))
			return fmt.Errorf("package installation failed - packages not found: %v", failedPackages)
		}
	}

	duration := time.Since(start)
	logger.LogMessage("SUCCESS", fmt.Sprintf("✅ Package installation completed in %v", duration))

	// Log successful installations
	for _, pkg := range toInstall {
		logger.LogMessage("SUCCESS", fmt.Sprintf("✓ Installed: %s", pkg))
	}

	return nil
}

// containsAURIndicators checks if a package name suggests it's from AUR
func containsAURIndicators(pkg string) bool {
	indicators := []string{"-git", "-bin", "-devel", "-beta", "-alpha", "-rc", "-nightly"}
	for _, indicator := range indicators {
		if len(pkg) > len(indicator) && pkg[len(pkg)-len(indicator):] == indicator {
			return true
		}
	}
	return false
}

// SyncPackageDatabases syncs pacman and yay databases
func SyncPackageDatabases() error {
	logger.LogMessage("INFO", "🔄 Syncing package databases... (can take a while)")
	logger.Log("Progress", "Database", "Database Sync", "Syncing databases")

	start := time.Now()

	// Sync pacman database
	if err := syncPackmanDatabase(); err != nil {
		logger.LogMessage("ERROR", fmt.Sprintf("Failed to sync pacman database: %v", err))
		return fmt.Errorf("pacman sync failed: %w", err)
	}

	// Sync yay database (non-critical)
	cmd := exec.Command("yay", "-Sy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Yay sync failure is not critical, just log it
		outputStr := string(output)
		if len(outputStr) > 200 {
			outputStr = outputStr[:200] + "... (truncated)"
		}
		logger.LogMessage("WARNING", fmt.Sprintf("Failed to sync yay database: %s", outputStr))
		logger.Log("Warning", "Database", "Database Sync", "Yay sync failed, continuing")
	} else {
		logger.LogMessage("SUCCESS", "Yay database synced")
	}

	duration := time.Since(start)
	logger.LogMessage("SUCCESS", fmt.Sprintf("Database sync completed in %v", duration))
	return nil
}

// syncPackmanDatabase syncs the pacman database
func syncPackmanDatabase() error {
	cmd := exec.Command("sudo", "pacman", "-Sy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pacman -Sy failed: %w, output: %s", err, string(output))
	}
	return nil
}

// isPackageInstalled checks if a package is already installed
func isPackageInstalled(packageName string) bool {
	cmd := exec.Command("pacman", "-Q", packageName)
	err := cmd.Run()
	return err == nil
}

// CommandExists checks if a command is available in PATH
func CommandExists(name string) bool {
	cmd := exec.Command("which", name)
	err := cmd.Run()
	return err == nil
}
