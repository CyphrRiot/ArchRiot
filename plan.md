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

### ✅ TASK 2: Kernel Upgrade Reboot Detection - COMPLETED

**Implementation Summary:**

- ✅ Added kernel upgrade detection in `upgrade.go` using `pacman -Qu` to check for "linux" package upgrades
- ✅ Created `KernelUpgradeMsg` message type to communicate kernel upgrade status to TUI
- ✅ Modified reboot prompt to default to "YES" with special message when kernel upgraded
- ✅ Maintained existing "NO" default for regular upgrades

**Files Modified:**

- `ArchRiot/source/upgrade/upgrade.go` - Added `detectKernelUpgrade()` function
- `ArchRiot/source/tui/model.go` - Added kernel upgrade handling in reboot prompt
- `ArchRiot/source/tui/messages.go` - Added `KernelUpgradeMsg` type

**Result:** System now automatically recommends reboot when kernel is upgraded, improving user experience and system stability.

## NEXT STEPS

**Remaining Task:**

- **TASK 1: Secure Boot Enablement** - Complex security enhancement requiring UEFI firmware interaction, key management, and bootloader signing

**Phased Implementation Plan:**

**PHASE 1: Detection & User Interface (Low Risk)**

- Add Secure Boot status detection using existing system calls
- Extend existing TUI message types in `source/tui/messages.go`
- Add new input mode to existing `source/tui/model.go` system
- Integrate Secure Boot prompt into existing installation flow
- **Deliverable**: User sees Secure Boot status and can choose to enable it

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

**Recommended Next Action:**
Start with Phase 1 - low risk TUI integration to establish foundation

## VERSION HISTORY

- **v2.9.7**: Kernel upgrade reboot detection - intelligent reboot prompting
- **v2.9.6**: Blue light persistence, Plymouth progress fixes, control panel --reapply
- **v2.9.5**: System upgrade integration, AUR resilience enhancements
- **v2.9.4**: Core installer stability improvements
