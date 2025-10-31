package installer

import (
	"fmt"
	"os/exec"
)

// EnsureYay verifies that the AUR helper "yay" is installed.
// If it's not available, it attempts to install it (with minimal interaction).
// Returns an error only when prerequisites fail outright; otherwise this function
// prefers to be non-fatal (installation attempts may be skipped or fail quietly,
// allowing the rest of the system to continue without AUR support).
func EnsureYay() error {
	// Already present? Nothing to do.
	if _, err := exec.LookPath("yay"); err == nil {
		return nil
	}

	// Install prerequisites first (base-devel, git). If this fails, propagate the error
	// so the caller can present a meaningful message.
	if err := exec.Command("sudo", "pacman", "-S", "--noconfirm", "--needed", "base-devel", "git").Run(); err != nil {
		return fmt.Errorf("failed to install yay prerequisites (base-devel, git): %w", err)
	}

	// Retry yay installation a few times; AUR can be flaky.
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := attemptYayInstall(); err != nil {
			// On the final attempt, give up quietly (non-fatal) and continue without yay.
			// Earlier failures simply retry.
			if attempt == maxRetries {
				return nil
			}
			continue
		}
		// Installed successfully
		return nil
	}

	return nil
}

// attemptYayInstall performs a single attempt to install yay from the AUR using yay-bin.
// It returns an error on failure; the caller may choose to retry.
func attemptYayInstall() error {
	tempDir := "/tmp/yay-bin-install"

	// Clean up any previous attempts (best-effort)
	_ = exec.Command("rm", "-rf", tempDir).Run()

	// Clone yay-bin repository (precompiled binary variant)
	if err := exec.Command("git", "clone", "https://aur.archlinux.org/yay-bin.git", tempDir).Run(); err != nil {
		_ = exec.Command("rm", "-rf", tempDir).Run()
		return fmt.Errorf("failed to clone yay-bin repository: %w", err)
	}

	// Build and install yay-bin (this downloads and installs the binary via makepkg)
	cmd := exec.Command("makepkg", "-si", "--noconfirm")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		_ = exec.Command("rm", "-rf", tempDir).Run()
		return fmt.Errorf("failed to build/install yay-bin: %w", err)
	}

	// Clean up temp directory
	_ = exec.Command("rm", "-rf", tempDir).Run()

	// Verify yay installation
	if _, err := exec.LookPath("yay"); err != nil {
		return fmt.Errorf("yay installation appears to have failed (not found in PATH)")
	}

	return nil
}
