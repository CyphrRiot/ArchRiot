package upgradeguard

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// PreInstall performs a best‑effort upgrade guard before installation.
//
// Behavior:
// - If pacman is unavailable, it returns immediately (no-op).
// - Attempts to refresh archlinux-keyring to avoid signature/key issues during subsequent installs.
// - Checks for pending compositor/Wayland/portal upgrades (Hyprland/wlroots/etc).
//   - Logs a caution and a recommended upgrade sequence.
//   - When strictABI is true, exits the process with code 2 to prevent continuing
//     under potentially mismatched ABIs on multi-monitor systems.
//
// This mirrors the previous inline guard behavior that lived in main, but is now
// extracted to keep the entrypoint small and delegation-only.
func PreInstall(strictABI bool) {
	// If pacman is unavailable, nothing to do.
	if _, err := exec.LookPath("pacman"); err != nil {
		return
	}

	// Best-effort: ensure keyring is current so subsequent queries/installs don't fail.
	_ = exec.Command("sudo", "pacman", "-Sy", "--noconfirm", "archlinux-keyring").Run()

	// Check for pending upgrades; if core Wayland compositor/portal pieces are queued,
	// recommend a full system upgrade to avoid ABI mismatches (multi-monitor/wlroots issues).
	out, err := exec.Command("pacman", "-Qu").CombinedOutput()
	if err != nil {
		return
	}
	s := strings.ToLower(string(out))
	if strings.Contains(s, "hyprland") ||
		strings.Contains(s, "wlroots") ||
		strings.Contains(s, "xdg-desktop-portal-hyprland") ||
		strings.Contains(s, "wayland") {

		log.Println("⚠️  Detected pending compositor/portal updates (e.g., Hyprland/wlroots).")
		log.Println("    To avoid ABI mismatches on multi‑monitor systems, update your system before continuing:")
		log.Println("    sudo pacman -Sy archlinux-keyring && yay -Syu && yay -Yc && sudo paccache -r")

		// Enforce blocking when strict ABI mode is enabled.
		if strictABI {
			log.Println("❌ Strict ABI mode: blocking installation until system is fully upgraded.")
			os.Exit(2)
		}
	}
}
