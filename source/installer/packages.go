package installer

import (
	"fmt"
	"os/exec"
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

	// Try to install all packages at once with pacman
	start := time.Now()
	cmd := exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm", "--needed"}, toInstall...)...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If pacman fails, try with yay (handles both AUR and regular packages)
		logger.LogMessage("WARNING", "Pacman failed, trying yay for all packages")

		// Check if yay is available
		if !CommandExists("yay") {
			logger.LogMessage("ERROR", "Yay not available for AUR packages")
			return fmt.Errorf("package installation failed: pacman failed and yay not available")
		}

		cmd = exec.Command("yay", append([]string{"-S", "--noconfirm", "--needed"}, toInstall...)...)
		output, err = cmd.CombinedOutput()

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
	logger.LogMessage("SUCCESS", fmt.Sprintf("Package installation completed in %v", duration))
	return nil
}

// SyncPackageDatabases syncs pacman and yay databases
func SyncPackageDatabases() error {
	logger.LogMessage("INFO", "🔄 Syncing package databases...")
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
