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

	logger.LogMessage("INFO", fmt.Sprintf("Installing %d new packages: %v", len(toInstall), toInstall))

	// Install in batches to avoid overwhelming the system
	batchSize := 10
	for i := 0; i < len(toInstall); i += batchSize {
		end := i + batchSize
		if end > len(toInstall) {
			end = len(toInstall)
		}
		batch := toInstall[i:end]

		if err := installPackageBatch(batch); err != nil {
			return fmt.Errorf("batch installation failed: %w", err)
		}
	}

	return nil
}

// installPackageBatch installs a batch of packages
func installPackageBatch(packages []string) error {
	start := time.Now()

	// Try pacman first, then yay if needed
	cmd := exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm", "--needed"}, packages...)...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If pacman fails, try with yay for AUR packages
		logger.LogMessage("WARNING", "Pacman failed, trying yay for AUR packages")

		cmd = exec.Command("yay", append([]string{"-S", "--noconfirm", "--needed"}, packages...)...)
		output, err = cmd.CombinedOutput()

		if err != nil {
			// Limit output to first 200 characters to prevent TUI spam
			outputStr := string(output)
			if len(outputStr) > 200 {
				outputStr = outputStr[:200] + "... (truncated)"
			}
			logger.Log("Error", "Package", "Package Error", "Failed: "+outputStr)
			return fmt.Errorf("batch installation failed: %w", err)
		}
	}

	duration := time.Since(start)
	logger.LogMessage("SUCCESS", fmt.Sprintf("Installed %d packages in %v", len(packages), duration))
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
