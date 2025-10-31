package secboot

import (
	"archriot-installer/installer"
	"archriot-installer/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// RunContinuation handles the Secure Boot continuation flow within the TUI.
// It mirrors the previous inline logic in main.go, but lives in the secboot
// package to keep the entrypoint small and delegation-only.
//
// The caller should inject a logger via SetLogger(logger.LogMessage) before
// invoking this function, so RestoreHyprlandConfig() can report via logs.
//
// Parameters:
//   - program: the active Bubble Tea program
//   - model: the active TUI model (unused but kept for parity/extension)
//   - secureBootContinuationDone: channel receiving user's YES/NO decision on continuation
func RunContinuation(program *tea.Program, model *tui.InstallModel, secureBootContinuationDone <-chan bool) {
	_ = model // reserved for future use

	program.Send(tui.StepMsg("Validating Secure Boot setup..."))
	program.Send(tui.ProgressMsg(0.1))

	// Check current Secure Boot status
	sbEnabled, sbSupported, err := installer.DetectSecureBootStatus()
	if err != nil {
		program.Send(tui.LogMsg("❌ Failed to detect Secure Boot status: " + err.Error()))
		program.Send(tui.FailureMsg{Error: "Could not validate Secure Boot status"})
		return
	}

	program.Send(tui.ProgressMsg(0.3))

	if !sbSupported {
		program.Send(tui.LogMsg("❌ System does not support Secure Boot (Legacy BIOS detected)"))
		program.Send(tui.FailureMsg{Error: "Secure Boot not supported on this system"})
		return
	}

	if sbEnabled {
		// Success! Secure Boot is enabled
		program.Send(tui.StepMsg("Secure Boot successfully enabled!"))
		program.Send(tui.LogMsg("✅ Secure Boot validation successful"))
		program.Send(tui.LogMsg("🔒 LUKS encryption is now protected against memory attacks"))
		program.Send(tui.LogMsg("🎉 Setup complete - restoring normal system behavior"))
		program.Send(tui.ProgressMsg(0.8))

		// Restore hyprland.conf to normal welcome
		if err := RestoreHyprlandConfig(); err != nil {
			program.Send(tui.LogMsg("⚠️ Failed to restore hyprland.conf: " + err.Error()))
		} else {
			program.Send(tui.LogMsg("✅ Restored normal startup configuration"))
		}

		program.Send(tui.ProgressMsg(1.0))
		program.Send(tui.DoneMsg{})
		return
	}

	// Secure Boot not enabled - provide detailed guidance
	program.Send(tui.StepMsg("Secure Boot setup incomplete"))
	program.Send(tui.LogMsg("⚠️ Secure Boot is supported but not yet enabled"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("📋 UEFI SETUP INSTRUCTIONS:"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("1. Restart your computer"))
	program.Send(tui.LogMsg("2. Press the UEFI/BIOS key during boot:"))
	program.Send(tui.LogMsg("   • Dell: F2 or F12"))
	program.Send(tui.LogMsg("   • HP: F10 or ESC"))
	program.Send(tui.LogMsg("   • Lenovo: F1, F2, or Enter"))
	program.Send(tui.LogMsg("   • ASUS: F2 or DEL"))
	program.Send(tui.LogMsg("   • MSI: DEL or F2"))
	program.Send(tui.LogMsg("   • Acer: F2 or DEL"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("3. Navigate to Security or Boot settings"))
	program.Send(tui.LogMsg("4. Find 'Secure Boot' option"))
	program.Send(tui.LogMsg("5. Enable Secure Boot"))
	program.Send(tui.LogMsg("6. Save settings and exit (usually F10)"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.LogMsg("⚠️ IMPORTANT: If system fails to boot after enabling:"))
	program.Send(tui.LogMsg("   • Return to UEFI settings"))
	program.Send(tui.LogMsg("   • Disable Secure Boot"))
	program.Send(tui.LogMsg("   • System will boot normally"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.ProgressMsg(0.7))

	// Show retry/cancel options
	program.Send(tui.LogMsg("Choose an option:"))
	program.Send(tui.LogMsg("• YES: I will reboot to UEFI settings now"))
	program.Send(tui.LogMsg("• NO: Cancel setup and restore normal system"))
	program.Send(tui.LogMsg(""))
	program.Send(tui.SecureBootContinuationPromptMsg{})

	// Wait for user decision
	userWantsRetry := <-secureBootContinuationDone
	if userWantsRetry {
		// User chose to retry - show reboot prompt to go to UEFI
		program.Send(tui.LogMsg("🔄 You chose to continue - reboot to access UEFI settings"))
		program.Send(tui.LogMsg("⏳ This program will run again after you enable Secure Boot"))
		program.Send(tui.LogMsg(""))
		program.Send(tui.StepMsg("Reboot to enable Secure Boot in UEFI"))
		program.Send(tui.ProgressMsg(1.0))

		// Trigger reboot prompt (same as main installer)
		program.Send(tui.DoneMsg{})
	} else {
		// User chose to cancel - restore welcome and exit
		program.Send(tui.LogMsg("❌ User cancelled Secure Boot setup"))
		program.Send(tui.LogMsg("🔄 Restoring normal system behavior..."))
		if err := RestoreHyprlandConfig(); err != nil {
			program.Send(tui.LogMsg("⚠️ Failed to restore hyprland.conf: " + err.Error()))
			program.Send(tui.FailureMsg{Error: "Failed to restore system configuration"})
		} else {
			program.Send(tui.LogMsg("✅ System restored - welcome will show on next login"))
			program.Send(tui.LogMsg(""))
			program.Send(tui.LogMsg("Press any key to exit"))
			program.Send(tui.ProgressMsg(1.0))
			program.Send(tui.DoneMsg{})
		}
	}
}
