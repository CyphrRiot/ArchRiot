# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## 🚧 OUTSTANDING TASKS

### TASK 1: Secure Boot Enablement

### Problem

Users without Secure Boot enabled are vulnerable to memory hijacking attacks on LUKS encrypted drives. The installer should offer to enable Secure Boot during installation/upgrade to improve system security.

### Requirements

- Detect if Secure Boot is currently disabled during installation/upgrade
- Prompt user with clear explanation of Secure Boot benefits for LUKS protection
- If user selects "YES", guide them through Secure Boot enablement process
- Handle the complexity of UEFI setup, key management, and bootloader signing
- Ensure process works across different hardware vendors (Dell, Lenovo, etc.)
- Provide fallback/recovery options if Secure Boot setup fails

### Implementation Challenges

- **UEFI Firmware Interaction**: Different vendors have different UEFI interfaces
- **Key Management**: Generating and managing Secure Boot keys (PK, KEK, db, dbx)
- **Bootloader Signing**: Signing GRUB/systemd-boot with custom keys
- **Kernel Signing**: Signing Linux kernel and modules for Secure Boot validation
- **User Guidance**: Walking users through BIOS/UEFI settings safely
- **Recovery Planning**: Ensuring users can disable Secure Boot if needed

### Technical Approach

**CRITICAL: REUSE EXISTING TUI ARCHITECTURE**

ArchRiot already has a sophisticated TUI system that MUST be reused:

- Existing message types in `source/tui/messages.go` (LogMsg, ProgressMsg, etc.)
- Existing input modes in `source/tui/model.go` (git-username, git-email, reboot, etc.)
- Existing confirmation prompt system with YES/NO options
- Existing callback pattern for user decisions

**Implementation Using EXISTING TUI:**

1. **Detection Phase**: Check `bootctl status` and `/sys/firmware/efi/efivars` for Secure Boot state
2. **TUI Integration**: Add new message types (`SecureBootPromptMsg`, `SecureBootStatusMsg`) to existing `messages.go`
3. **User Education**: REUSE existing confirmation prompt system to explain LUKS memory attack protection benefits
4. **Decision Flow**: EXTEND existing `inputMode` system for "secure-boot-confirm" mode
5. **Key Generation**: Create custom Secure Boot key hierarchy (PK → KEK → db)
6. **Bootloader Setup**: Configure and sign bootloader with custom keys
7. **Kernel Setup**: Sign kernel and modules for Secure Boot validation
8. **UEFI Guidance**: Provide vendor-specific instructions for enabling Secure Boot
9. **Validation**: Verify Secure Boot is working correctly after setup

**Priority**: Medium - Important security enhancement but complex implementation

**ARCHITECTURE REQUIREMENTS:**

- NO new TUI interfaces - extend existing system only
- Follow existing message/callback patterns in `tui/messages.go` and `tui/model.go`
- Integrate into existing installation flow, don't create separate flows
- Use existing `tools.go` framework for optional Secure Boot tool access

**CODEBASE ANALYSIS COMPLETE:**

Based on examination of the existing ArchRiot TUI architecture, the following patterns are established:

**TUI Message System (`source/tui/messages.go`):**

- Message types: `LogMsg`, `ProgressMsg`, `StepMsg`, `DoneMsg`, `FailureMsg`
- Input types: `GitUsernameMsg`, `GitEmailMsg`, `RebootMsg`, `UpgradeMsg`, `KernelUpgradeMsg`
- Callback system: External packages set callbacks via `SetGitCallbacks()`, `SetUpgradeCallback()`
- Helper functions: `SetVersionGetter()`, `SetLogPathGetter()` for external data access

**TUI Model System (`source/tui/model.go`):**

- Input modes: `"git-username"`, `"git-email"`, `"reboot"`, `""` (no input)
- Confirmation system: `showConfirm`, `confirmPrompt`, `cursor` (0=YES, 1=NO)
- Flow: InputRequestMsg → input mode → callback → confirmation → action
- External integration: Callbacks trigger actions in main packages

**Logging Integration (`source/logger/logger.go`):**

- Semantic logging: `Log(status, logType, name, description)`
- TUI integration: `program.Send(tui.LogMsg(...))` for real-time display
- File logging: Separate detailed logs in `~/.cache/archriot/`

**Tools Framework (`source/tools/tools.go`):**

- Tool structure: `ID`, `Name`, `Description`, `Category`, `ExecuteFunc`, `Advanced`, `Available`
- Availability checking: Dynamic `checkSecureBootAvailable()` functions
- Execution pattern: System checks → Package installation → Configuration → Verification

**System Integration Patterns:**

- System calls: `exec.Command()` with proper error handling
- Package checks: `exec.LookPath()` for availability, `pacman -Q` for installation
- File system checks: `os.Stat()` for directories/files (e.g., `/sys/firmware/efi`)
- State management: Global variables with TUI callbacks for user decisions

## NEXT STEPS

**Remaining Tasks:**

- **TASK 1: Secure Boot Enablement** - Complex security enhancement requiring UEFI firmware interaction, key management, and bootloader signing

**Secure Boot Phased Implementation Plan:**

**PHASE 1: Detection & User Interface (Low Risk) - DETAILED IMPLEMENTATION**

**1.1 Secure Boot Detection System**

```go
// Add to source/installer/secureboot.go (NEW FILE)
func DetectSecureBootStatus() (enabled bool, supported bool, err error) {
    // Check UEFI support: /sys/firmware/efi directory exists
    if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
        return false, false, nil // Legacy BIOS system
    }

    // Check current Secure Boot state: bootctl status
    cmd := exec.Command("bootctl", "status")
    output, err := cmd.Output()
    if err != nil {
        return false, true, fmt.Errorf("failed to check bootctl status: %w", err)
    }

    // Parse bootctl output for Secure Boot state
    enabled = strings.Contains(string(output), "Secure Boot: enabled")
    supported = true // UEFI system with bootctl means Secure Boot is supported

    return enabled, supported, nil
}

func DetectLuksEncryption() (bool, []string, error) {
    var luksDevices []string

    // Method 1: Check /proc/cmdline for cryptdevice= parameters
    if cmdline, err := os.ReadFile("/proc/cmdline"); err == nil {
        if strings.Contains(string(cmdline), "cryptdevice=") ||
           strings.Contains(string(cmdline), "rd.luks.uuid=") {
            luksDevices = append(luksDevices, "cmdline")
        }
    }

    // Method 2: Check /etc/crypttab for LUKS entries
    if crypttab, err := os.ReadFile("/etc/crypttab"); err == nil {
        lines := strings.Split(string(crypttab), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line != "" && !strings.HasPrefix(line, "#") {
                parts := strings.Fields(line)
                if len(parts) >= 2 {
                    luksDevices = append(luksDevices, parts[0]) // mapper name
                }
            }
        }
    }

    // Method 3: Check active device mapper for crypt devices
    cmd := exec.Command("dmsetup", "ls", "--target", "crypt")
    if output, err := cmd.Output(); err == nil {
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line != "" && line != "No devices found" {
                parts := strings.Fields(line)
                if len(parts) > 0 {
                    luksDevices = append(luksDevices, parts[0])
                }
            }
        }
    }

    // Method 4: Use lsblk to find encrypted partitions
    cmd = exec.Command("lsblk", "-f", "-o", "NAME,FSTYPE")
    if output, err := cmd.Output(); err == nil {
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if strings.Contains(line, "crypto_LUKS") {
                parts := strings.Fields(line)
                if len(parts) > 0 {
                    deviceName := strings.TrimPrefix(parts[0], "├─")
                    deviceName = strings.TrimPrefix(deviceName, "└─")
                    luksDevices = append(luksDevices, deviceName)
                }
            }
        }
    }

    // Remove duplicates
    uniqueDevices := make(map[string]bool)
    var result []string
    for _, device := range luksDevices {
        if !uniqueDevices[device] {
            uniqueDevices[device] = true
            result = append(result, device)
        }
    }

    return len(result) > 0, result, nil
}
```

**1.2 TUI Message Extension**

```go
// Add to source/tui/messages.go
type SecureBootStatusMsg struct {
    Enabled     bool
    Supported   bool
    LuksUsed    bool
    LuksDevices []string  // List of detected LUKS devices
}

type SecureBootPromptMsg struct{}

// Add callback support
var secureBootCallback func(bool)
func SetSecureBootCallback(callback func(bool)) {
    secureBootCallback = callback
}
```

**1.3 TUI Model Integration**

```go
// Add to source/tui/model.go InstallModel struct
type InstallModel struct {
    // ... existing fields ...
    secureBootEnabled   bool
    secureBootSupported bool
    luksDetected        bool
    luksDevices         []string
}

// Add to Update() method
case SecureBootStatusMsg:
    // Store secure boot status in model
    m.secureBootEnabled = msg.Enabled
    m.secureBootSupported = msg.Supported
    m.luksDetected = msg.LuksUsed
    m.luksDevices = msg.LuksDevices
    return m, nil

case SecureBootPromptMsg:
    if !m.secureBootEnabled && m.secureBootSupported && m.luksDetected {
        m.showConfirm = true
        deviceList := strings.Join(m.luksDevices, ", ")
        m.confirmPrompt = fmt.Sprintf("🛡️ Enable Secure Boot? (%s)", deviceList)
        m.cursor = 1 // Default to NO (conservative)
    }
    return m, nil

// Add to handleConfirmSelection()
} else if strings.HasPrefix(m.confirmPrompt, "🛡️ Enable Secure Boot protection?") {
    m.showConfirm = false
    m.confirmPrompt = ""
    if secureBootCallback != nil {
        secureBootCallback(m.cursor == 0) // YES = 0
    }
    return m, nil
}
```

**1.4 Installation Flow Integration**

```go
// Add to source/orchestrator/orchestrator.go in RunInstallation()
// After package installation, before final steps:

func checkSecureBootRecommendation() {
    logger.LogMessage("INFO", "🔍 Checking Secure Boot and LUKS status...")

    // Detect Secure Boot status
    sbEnabled, sbSupported, err := installer.DetectSecureBootStatus()
    if err != nil {
        logger.LogMessage("WARNING", fmt.Sprintf("Failed to detect Secure Boot status: %v", err))
        return
    }

    // Detect LUKS encryption
    luksUsed, luksDevices, err := installer.DetectLuksEncryption()
    if err != nil {
        logger.LogMessage("WARNING", fmt.Sprintf("Failed to detect LUKS encryption: %v", err))
        luksUsed = false
        luksDevices = []string{}
    }

    logger.LogMessage("INFO", fmt.Sprintf("Secure Boot: enabled=%v, supported=%v", sbEnabled, sbSupported))
    logger.LogMessage("INFO", fmt.Sprintf("LUKS: detected=%v, devices=%v", luksUsed, luksDevices))

    // Send status to TUI
    program.Send(tui.SecureBootStatusMsg{
        Enabled:     sbEnabled,
        Supported:   sbSupported,
        LuksUsed:    luksUsed,
        LuksDevices: luksDevices,
    })

    // Prompt user if Secure Boot should be enabled for LUKS protection
    if !sbEnabled && sbSupported && luksUsed {
        logger.LogMessage("INFO", "Prompting user for Secure Boot enablement...")
        program.Send(tui.SecureBootPromptMsg{})
        // Wait for user response via callback
    }
}
```

**1.5 Tools Integration**

```go
// Update source/tools/tools.go GetAvailableTools()
{
    ID:          "secure-boot-setup",
    Name:        "🛡️ Secure Boot Setup",
    Description: "Enable Secure Boot protection + encryption",
    Category:    "Security",
    ExecuteFunc: executeSecureBootSetup,
    Advanced:    false, // Make accessible to all users
    Available:   checkSecureBootSetupAvailable(),
}

func checkSecureBootSetupAvailable() bool {
    status, err := installer.DetectSecureBootStatus()
    return err == nil && status.supported && !status.enabled
}
```

**1.6 User Education Integration**

```go
func getSecureBootEducation(luksDevices []string) string {
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
```

**PHASE 1 DELIVERABLES:**

- Secure Boot status detection in installation flow
- User education about LUKS memory attack protection
- Optional Secure Boot enablement prompt during install
- Tools menu access for post-install setup
- Foundation for Phase 2 key management

**PHASE 2: Key Management Foundation (Medium Risk)**

- Implement secure key generation (PK, KEK, db hierarchy)
- Add key storage and backup mechanisms
- Create validation functions for key integrity
- **Deliverable**: Secure Boot keys can be generated and managed safely

**PHASE 3: Bootloader Integration (High Risk)**

- Sign GRUB/systemd-boot with generated keys
- Implement pacman hook for automatic kernel signing
- Add recovery mechanisms if signing fails
- **Deliverable**: System can boot with Secure Boot enabled

**PHASE 4: User Guidance & Recovery (Medium Risk)**

- Vendor-specific UEFI setup instructions
- Automated verification of Secure Boot status
- Recovery tools if Secure Boot breaks system
- **Deliverable**: Complete, production-ready Secure Boot enablement

**PHASE 1 IMPLEMENTATION PROGRESS:**

✅ **COMPLETED (v2.10.3 - Released):**

- Created `source/installer/secureboot.go` with comprehensive detection system
- Implemented `DetectSecureBootStatus()` with bootctl and efivars fallback
- Implemented `DetectLuksEncryption()` with 4 detection methods (cmdline, crypttab, dmsetup, lsblk)
- Added educational functions and prerequisite validation
- Extended `source/tui/messages.go` with SecureBootStatusMsg, SecureBootPromptMsg, SecureBootConfirmMsg
- Added Secure Boot callback system matching existing git/upgrade patterns
- Extended `source/tui/model.go` InstallModel struct with Secure Boot fields
- Added SecureBootStatusMsg and SecureBootPromptMsg handling in TUI Update method
- Added Secure Boot confirmation handling to handleConfirmSelection method
- Integrated Secure Boot checking into orchestrator RunInstallation() flow (at 98.5% completion)
- Updated tools framework to include Secure Boot setup option with dynamic availability
- Added `--secure_boot_stage` command line flag for post-reboot continuation
- Implemented `runSecureBootContinuation()` function in main.go
- Added hyprland.conf modification and restoration logic
- Connected callback channels between TUI and orchestrator
- **Build verification: Code compiles successfully with `make`**
- **Released as v2.10.3 with Secure Boot prompting temporarily disabled**

🔄 **CURRENT STATUS (Bootloader Signing Critical Failure):**

**CRITICAL ISSUE DISCOVERED**: Secure Boot implementation causes complete boot failure with "Operating System Loader has no signature" error. Users get locked out of their systems.

**Root Cause**: Bootloader signing process fails or signed bootloader is not being used properly, causing UEFI to reject boot with Secure Boot enabled.

**SAFETY MEASURE**: Secure Boot prompting disabled again with `if false &&` to protect users from boot failures.

**SECURE BOOT CONTINUATION ARCHITECTURE:**

**Problem Solved**: How to continue Secure Boot setup after mandatory reboot for UEFI settings.

**Solution**: Hyprland.conf modification approach

1. **During Setup**: Replace `exec-once = welcome` with `exec-once = /home/username/.local/share/archriot/install/archriot --secure_boot_stage` (using existing `os.UserHomeDir()` + `filepath.Join()` pattern from codebase)
2. **Post-Reboot**: User boots into Secure Boot continuation instead of normal welcome screen
3. **Continuation Flow**:
    - Validates current Secure Boot status
    - Guides user through UEFI settings if needed
    - Provides troubleshooting for boot issues
    - Handles completion/cancellation
    - Restores original hyprland.conf when done
4. **Self-Cleaning**: Automatic restoration of normal system behavior

**IMPLEMENTATION STAGES:**

**Stage 1**: Pre-Reboot Setup (98.5% of installation)

- Detect Secure Boot/LUKS status
- User confirmation prompt
- If YES: Generate keys, sign bootloader, modify hyprland.conf
- Complete installation and reboot

**Stage 2**: Post-Reboot Continuation (`--secure_boot_stage`)

- Launched instead of welcome screen
- Validate Secure Boot enablement
- Guide user through UEFI configuration
- Verify boot chain integrity
- Restore normal system on completion

✅ **PHASE 1 IMPLEMENTATION COMPLETE - UNTESTED:**

All core Phase 1 functionality implemented, compiles successfully, and deployed in v2.10.3 but completely untested.

**CURRENT STATE:**

- **Code Status**: Implemented but untested
- **Deployment Status**: Released in v2.10.3 but disabled via `if false &&` condition
- **Re-enablement**: Requires testing before activation

**TESTING REQUIREMENTS (For Re-enablement):**

1. **Detection Testing**: Verify Secure Boot and LUKS detection on various systems
2. **Flow Testing**: Test installation prompt appears only when appropriate (UEFI + LUKS + no Secure Boot)
3. **Continuation Testing**: Verify hyprland.conf modification and post-reboot continuation
4. **Restoration Testing**: Ensure hyprland.conf restores properly after completion
5. **Edge Case Testing**: Legacy BIOS systems, already-enabled Secure Boot, no LUKS encryption

**STAGE 4 IMPLEMENTATION COMPLETE:**

The post-reboot continuation (`--secure_boot_stage`) has been significantly enhanced:

**Implemented Features:**

✅ **Detailed UEFI instructions** - Vendor-specific key combinations (Dell, HP, Lenovo, ASUS, MSI, Acer)
✅ **Step-by-step guidance** - Clear instructions for enabling Secure Boot in UEFI settings
✅ **Retry/Cancel options** - Users can choose to continue or cancel setup gracefully
✅ **Recovery mechanism** - Automatically restores hyprland.conf if user cancels
✅ **Better user guidance** - All instructions displayed in scrollable log window
✅ **Proper confirmation prompts** - Uses existing TUI system for user choices

**Flow Implementation:**

1. **Status Detection** - Checks if Secure Boot is enabled after reboot
2. **Success Path** - If enabled, validates setup and restores normal system
3. **Guidance Path** - If not enabled, shows detailed UEFI instructions
4. **User Choice** - "Continue setup?" with clear YES/NO explanations
5. **Retry Handling** - YES keeps continuation active for next reboot
6. **Cancel Handling** - NO restores hyprland.conf and returns to normal system

**Code Status:** Implemented and compiles successfully but CAUSES BOOT FAILURES - DISABLED FOR SAFETY

**PHASE 2 IMPLEMENTATION STATUS:**

✅ **Code Complete**: All Phase 2 functions implemented:

- `generateSecureBootKeys()` - Creates PK, KEK, db key hierarchy ✅
- `signBootComponents()` - Signs bootloader and kernel with custom keys ❌ **BROKEN**
- `setupPackmanHooks()` - Creates automatic signing hooks for kernel updates ✅
- `installKeysIntoUEFI()` - Installs keys directly into UEFI firmware ✅
- Fixed storage paths to `~/.config/archriot/` to avoid installer overwrites ✅
- Fixed UUID generation for proper UEFI key identification ✅
- Fixed Setup Mode detection and user guidance ✅

❌ **CRITICAL FAILURE**: Bootloader signing causes boot failures

**DISCOVERED ISSUES:**

1. **Key Generation**: Works correctly with proper UUIDs
2. **Key Installation**: Works correctly in Setup Mode
3. **Bootloader Signing**: FAILS - causes "no signature" boot error
4. **User Experience**: Complete nightmare - users get locked out

**SPECIFIC FAILURES:**

- Bootloader signing doesn't work properly with custom keys
- Signed bootloader may not be in correct location
- UEFI rejects bootloader even with custom keys installed
- No recovery mechanism for boot failures
- Users must manually disable Secure Boot to recover

**ARCHITECTURAL PROBLEMS:**

- Missing bootloader validation step
- No verification that signed bootloader works before enabling Secure Boot
- No fallback unsigned bootloader for recovery
- Continuation flow doesn't handle "already enabled but failing" state

**IMMEDIATE NEXT ACTION:**
Fix bootloader signing process and add validation before allowing Secure Boot enablement. Must ensure signed bootloader actually works before exposing to users.
