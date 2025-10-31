package session

import (
	"log"
	"os/exec"
	"time"
)

// RebootNow centralizes the system reboot sequence.
//
// Behavior:
// 1) Sync filesystems (best-effort)
// 2) Brief delay to allow background processes to quiesce
// 3) Attempt sudo shutdown -r now
// 4) Fallback to sudo systemctl reboot on failure
//
// Notes:
// - All steps are best-effort; errors are logged and the function returns.
// - Callers should ensure appropriate privileges/environment for reboot.
func RebootNow() {
	log.Println("🔄 Preparing for system reboot...")
	log.Println("💾 Syncing filesystems...")

	// Sync filesystems (best-effort)
	_ = exec.Command("sync").Run()

	// Give time for any background processes to finish
	log.Println("⏳ Waiting for processes to complete...")
	time.Sleep(2 * time.Second)

	// Clean shutdown and reboot
	log.Println("🔄 Initiating system reboot...")
	if err := exec.Command("sudo", "shutdown", "-r", "now").Run(); err != nil {
		log.Printf("❌ Failed to reboot via shutdown: %v", err)
		// Fallback to systemctl if shutdown fails
		log.Println("🔄 Trying fallback reboot method (systemctl)...")
		_ = exec.Command("sudo", "systemctl", "reboot").Run()
	}
}
