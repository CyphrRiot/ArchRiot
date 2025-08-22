# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## 🚧 OUTSTANDING TASKS

### TASK 1: Matugen Dynamic Color Theming Integration - IN PROGRESS

### Current Status: PARTIAL IMPLEMENTATION

**✅ COMPLETED:**

- Central color system using waybar's `@define-color` syntax
- Waybar CSS integration with colors.css import working
- Go theming package with matugen integration
- CLI commands: `--apply-wallpaper-theme` and `--toggle-dynamic-theming`
- Control panel toggle for dynamic theming
- Wallpaper change hooks in swaybg-next and control panel
- Package dependency: matugen added to packages.yaml

**❌ REMAINING WORK:**

- Go theming system needs to generate waybar `@define-color` syntax (not CSS variables)
- Hyprland config integration (string replacement approach)
- GTK theme integration
- Testing and validation of complete end-to-end flow

### Critical Lessons Learned

**WAYBAR CSS SYNTAX:**

- Waybar does NOT support CSS variables (`:root { --var: value }`)
- Waybar uses `@define-color colorname #value;` syntax only
- Import works with local files: `@import url("colors.css");`
- References use `@colorname` syntax (not `var(--colorname)`)

**IMPLEMENTATION APPROACH:**

- Colors.css lives in waybar directory, gets installed with waybar configs
- No separate config copying needed - waybar/\* pattern handles it
- Go theming system writes to `~/.config/waybar/colors.css`
- Must generate `@define-color` format, not CSS variables

## Next Implementation Steps

### IMMEDIATE PRIORITY: Fix Go Theming System

**Problem:** Current Go code generates CSS variables, but waybar requires `@define-color` syntax.

**Solution:** Update `GenerateColorsCSS()` function in `source/theming/theming.go` to output:

```
@define-color primary_color #7aa2f7;
@define-color accent_color #bb9af7;
```

Instead of:

```
:root {
  --primary-color: #7aa2f7;
  --accent-color: #bb9af7;
}
```

### REMAINING COMPONENTS TO INTEGRATE

1. **Hyprland Configuration** - String replacement approach for:
    - `col.active_border = rgba(89b4fa88)` → extracted primary color
    - `col.inactive_border`, shadow colors, group colors

2. **GTK Theme Integration** - Apply colors to:
    - `config/gtk-3.0/gtk.css` background/text colors
    - Selection and interactive element colors

3. **End-to-End Testing** - Verify complete workflow:
    - Toggle dynamic theming ON/OFF via control panel
    - Change wallpaper → colors update automatically
    - Fallback to CypherRiot when disabled

### TECHNICAL IMPLEMENTATION NOTES

**File Structure:**

- `colors.css` lives in waybar directory (`~/.config/waybar/colors.css`)
- Go theming system writes to waybar location
- Waybar imports with `@import url("colors.css");`

**Color Format Conversion:**

- Matugen outputs: `"primary": "#7aa2f7"`
- Must generate: `@define-color primary_color #7aa2f7;`
- Waybar references: `color: @primary_color;`

**Integration Points Working:**

- ✅ swaybg-next calls `--apply-wallpaper-theme`
- ✅ Control panel calls theming system
- ✅ CLI commands implemented
- ✅ Dynamic theming toggle saves to `background-prefs.json`

### TASK 2: Secure Boot Implementation Overhaul

### Problem Analysis

The current Secure Boot implementation has critical flaws that
cause boot failures:

1. **Disabled due to boot failures** - Currently disabled with `false &&` in orchestrator.go
2. **Complex manual approach** - Uses error-prone `openssl`/`efitools` instead of reliable `sbctl`
3. **Setup Mode issues** - Keys fail to install because UEFI Setup Mode isn't properly verified
4. **Bootloader signing failures** - Manual `sbsign` approach is unreliable
5. **No recovery mechanisms** - Silent failures with no user guidance

### Root Cause

Current implementation uses Method 3 (Manual) from Secure_Boot.md documentation, which is the most complex and error-prone approach. The documentation recommends Method 1 (`sbctl`) for simplicity and reliability.

### Implementation Strategy (TESTED APPROACH)

Replace the entire implementation with the proven `sbctl` method from our documentation:

## PHASE 1: Replace Key Generation System

### Objective

Replace complex `openssl`/`efitools` with simple `sbctl create-keys`

### Implementation

- Remove `generateSecureBootKeys()`, `generateKey()`, `generateUUID()` functions
- Replace with single `sbctl create-keys` command
- Remove `prepareKeyInstallation()` complexity
- Use `sbctl` default key locations (`/usr/share/secureboot/keys/`)

### Testing Requirements

1. **Unit Test**: Verify `sbctl create-keys` executes successfully
2. **Integration Test**: Check key files are created with proper permissions
3. **Verification Test**: Use `sbctl status` to confirm keys are valid
4. **Error Test**: Verify graceful handling if `sbctl` package missing

## PHASE 2: Fix Setup Mode Detection and Enrollment

### Objective

Implement reliable key enrollment using `sbctl enroll-keys` (custom keys only)

### Implementation

- Enhance Setup Mode detection with `sbctl status` output parsing
- Replace manual EFI variable manipulation with `sbctl enroll-keys`
- Add clear user guidance for entering Setup Mode via UEFI settings
- Remove complex `.auth` file generation

### Testing Requirements

1. **Setup Mode Test**: Verify detection when UEFI is/isn't in Setup Mode
2. **Enrollment Test**: Test `sbctl enroll-keys` with custom keys only (NO --microsoft)
3. **Status Test**: Verify `sbctl status` shows enrolled keys correctly
4. **User Guidance Test**: Verify clear instructions for UEFI Setup Mode

## PHASE 3: Implement Reliable Signing

### Objective

Replace manual `sbsign` with proven `sbctl sign` commands

### Implementation

- Remove `signBootComponents()`, `signFile()` functions
- Use `sbctl sign -s` for individual files
- Use `sbctl sign-all` for batch operations
- Implement documented pacman hook from Secure_Boot.md

### Testing Requirements

1. **Individual Signing Test**: Verify bootloader and kernel signing works
2. **Batch Signing Test**: Test `sbctl sign-all` functionality
3. **Hook Test**: Verify pacman hook auto-signs on kernel updates
4. **Verification Test**: Use `sbverify` to confirm signatures are valid

## PHASE 4: Add Recovery Mechanisms

### Objective

Implement proper error handling and recovery options

### Implementation

- Add `sbctl verify` checks after signing
- Implement boot failure recovery instructions
- Add `sbctl reset` option for failed setups
- Enhance user guidance with vendor-specific UEFI instructions

### Testing Requirements

1. **Verification Test**: Test `sbctl verify` on signed components
2. **Recovery Test**: Verify reset functionality works
3. **Boot Test**: Actual reboot test with Secure Boot enabled
4. **Failure Test**: Test recovery when Secure Boot fails to boot

## PHASE 5: Re-enable and Integration Testing

### Objective

Remove disable flag and test complete end-to-end flow

### Implementation

- Remove `false &&` disable flag in `checkSecureBootRecommendation()`
- Update user prompts to reflect custom-key-only approach
- Add comprehensive logging throughout process
- Implement status reporting and progress indication

### Testing Requirements

1. **End-to-End Test**: Complete flow from detection to working Secure Boot
2. **Reboot Test**: Verify system boots successfully with Secure Boot enabled
3. **Multiple Hardware Test**: Test on different UEFI implementations
4. **Error Scenario Test**: Test all failure modes and recovery paths

## COMPREHENSIVE TESTING STRATEGY

### Pre-Implementation Testing

- [ ] Verify `sbctl` package is available in repositories
- [ ] Test `sbctl` commands manually on development system
- [ ] Document exact command sequences that work

### Per-Phase Testing

- [ ] Unit tests for each function replacement
- [ ] Integration tests for phase functionality
- [ ] Regression tests to ensure no functionality loss
- [ ] Manual verification of each step

### Final Integration Testing

- [ ] VM testing with different UEFI configurations
- [ ] Physical hardware testing on multiple vendors
- [ ] Dual-boot testing (ensure Windows unaffected)
- [ ] Recovery testing (intentional failure scenarios)

### Success Criteria

- [ ] `sbctl status` shows "Secure Boot: Enabled" after setup
- [ ] System boots successfully with Secure Boot enabled
- [ ] Kernel updates automatically maintain Secure Boot signatures
- [ ] Clear user guidance for any required manual steps
- [ ] Robust error handling with recovery options

## NEXT STEPS

**Current Priority:**

1. Phase 1: Key Generation Replacement
2. Comprehensive testing of each phase before proceeding
3. Documentation updates to reflect new approach

**Development Approach:**

- Make ONE change at a time
- Test thoroughly after each change
- Wait for feedback before proceeding
- Maintain backup and recovery options

**Technical Notes:**

- Use `sbctl` exclusively - no manual `openssl`/`efitools`
- Custom keys only - no Microsoft dual-boot support needed
- Follow Method 1 from Secure_Boot.md documentation exactly
- Prioritize reliability over complexity

## COMPLETED TASKS

### ✅ TASK 1: Terminal Emoji Fallback System - COMPLETED

- Automatic detection of terminal emoji support implemented
- ASCII fallbacks for all emojis working correctly
- Maintains exact visual width/spacing
- Test mode with `--ascii-only` flag functional
- Both logger system and TUI hardcoded strings handled properly

### ✅ TASK 2: Kernel Upgrade Reboot Detection - COMPLETED

- Fixed reboot handling using prefix matching
- Improved continuation flow for reboot requests
- Released as v2.10.8

## ARCHITECTURAL PRINCIPLES

**Code Quality:**

- Modular design with clear separation of concerns
- TUI-based user interaction with semantic logging
- Comprehensive error handling and user feedback
- Extensive testing at every level

**Security Focus:**

- Custom Secure Boot keys for maximum control
- No reliance on third-party certificate authorities
- Clear chain of trust documentation
- Recovery mechanisms for failed configurations

**User Experience:**

- Clear educational content about Secure Boot benefits
- Step-by-step guidance for UEFI configuration
- Robust error messages and recovery instructions
- Automated where possible, guided manual steps when necessary
