# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## ✅ COMPLETED TASKS

### ✅ PHASE 1 - TASK 1: Waybar Visual Bar Indicators - COMPLETED

**Problem**

Current waybar modules (temperature, CPU, memory) display numeric values which are less visually intuitive than graphical indicators. Users want immediate visual feedback through bar-style indicators with color coding.

**Requirements**

Replace numeric displays with visual bar indicators using Unicode block progression characters:

- `▁` = minimal level (10-25%)
- `▂` = low level (26-40%)
- `▃` = moderate level (41-55%)
- `▄` = medium level (56-70%)
- `▅` = high level (71-80%)
- `▆` = very high level (81-90%)
- `▇` = maximum level (91-95%)
- `█` = critical level (96-100%)

**Technical Specifications**

**TEMPERATURE MODULE:**

- Range: 30°C to 100°C
- Bar progression: No bar (≤30°C) → `▁` (31-40°C) → `▂` (41-50°C) → `▃` (51-60°C) → `▄` (61-70°C) → `▅` (71-80°C) → `▆` (81-90°C) → `▇` (91-95°C) → `█` (≥96°C)
- Color coding:
    - Normal: Default waybar color (≤70°C)
    - Warning: Reddish-purple (71-80°C)
    - Critical: Bright red (81-90°C)
    - Danger: Very bright red (≥91°C)

**CPU MODULE:**

- Range: 0% to 100%
- Bar progression: No bar (≤10%) → `▁` (11-25%) → `▂` (26-40%) → `▃` (41-55%) → `▄` (56-70%) → `▅` (71-80%) → `▆` (81-90%) → `▇` (91-95%) → `█` (≥96%)
- Color changes: Only above 90% (red warning)

**MEMORY MODULE:**

- Range: 0% to 100% of total RAM
- Bar progression: No bar (≤10%) → `▁` (11-25%) → `▂` (26-40%) → `▃` (41-55%) → `▄` (56-70%) → `▅` (71-80%) → `▆` (81-90%) → `▇` (91-95%) → `█` (≥96%)
- Color changes: Only above 90% (red warning)

**Implementation Approach**

**REUSE EXISTING ARCHITECTURE:**

- Modify existing Python scripts: `waybar-cpu-aggregate.py`, `waybar-memory-accurate.py`
- Use existing waybar temperature module (already configured)
- REUSE existing CSS classes and styling system in `config/waybar/style.css`

**Technical Steps:**

1. Update `waybar-cpu-aggregate.py` to output bar indicators instead of percentages
2. Update `waybar-memory-accurate.py` to output bar indicators instead of GB values
3. Configure temperature module to use bar format with color coding
4. Add new CSS classes for warning/critical states in `style.css`
5. Test all modules integrate properly with existing waybar layout

**Priority**: Medium - Visual improvement that enhances user experience

**Implementation Summary:**

- ✅ Created modular `get_visual_bar()` function for reuse across all modules
- ✅ Updated `waybar-cpu-aggregate.py` with realistic thresholds for modern multi-core systems
- ✅ Created `waybar-temp-bar.py` custom module with temperature-appropriate bar ranges
- ✅ Updated `waybar-memory-accurate.py` to use visual bar indicators
- ✅ Configured waybar to use grouped system metrics for tight visual clustering
- ✅ Maintained original CypherRiot color scheme and module icons
- ✅ Adjusted spacing for optimal visual balance

**Files Modified:**

- `config/bin/scripts/waybar-cpu-aggregate.py` - Added visual bars with realistic CPU thresholds
- `config/bin/scripts/waybar-temp-bar.py` - NEW: Custom temperature module with visual bars
- `config/bin/scripts/waybar-memory-accurate.py` - Added visual bars for memory usage
- `config/waybar/ModulesCustom` - Added custom/temp-bar module definition
- `config/waybar/ModulesGroups` - Added group/system-metrics for tight clustering
- `config/waybar/config` - Updated to use grouped system metrics
- `config/waybar/style.css` - Added temperature bar styling with color coding

**Result:** Waybar now displays intuitive visual bar indicators (`▁ ▂ ▃ ▄ ▅ ▆ ▇ █`) for temperature, CPU, and memory instead of raw numbers, with proper spacing and color schemes.

**Priority**: Medium - Visual improvement that enhances user experience ✅ COMPLETED

## 🎯 NEXT IMMEDIATE TASK

### PHASE 1 - TASK 2: TBD

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

## NEXT STEPS

**Remaining Tasks:**

- **TASK 1: Secure Boot Enablement** - Complex security enhancement requiring UEFI firmware interaction, key management, and bootloader signing

**Secure Boot Phased Implementation Plan:**

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
Define Phase 1 Task 2 or begin Secure Boot Phase 1 (Detection & User Interface)

## VERSION HISTORY

- **v2.9.8**: Waybar visual bar indicators - intuitive system metrics display
- **v2.9.7**: Kernel upgrade reboot detection - intelligent reboot prompting
- **v2.9.6**: Blue light persistence, Plymouth progress fixes, control panel --reapply
- **v2.9.5**: System upgrade integration, AUR resilience enhancements
- **v2.9.4**: Core installer stability improvements
