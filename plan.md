# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## 🚧 OUTSTANDING TASKS

### TASK 1: Matugen Dynamic Color Theming Integration (TOP PRIORITY)

### Problem Analysis

The current ArchRiot theming system has limitations that reduce user customization and visual coherence:

1. **Static color scheme** - Fixed CypherRiot colors regardless of wallpaper
2. **Hardcoded theme values** - Colors defined statically across multiple files (Go TUI, GTK, Waybar, Hyprland)
3. **Disconnected wallpaper system** - Wallpaper changes don't influence overall system aesthetics
4. **Limited personalization** - Users cannot have themes that match their chosen wallpapers
5. **Manual theme maintenance** - Any color changes require editing multiple configuration files

### Root Cause

Current implementation uses hardcoded CypherRiot color values throughout the system without any dynamic color extraction or template-based theming approach.

### Implementation Strategy

Integrate `matugen` (Material Design 3 color extraction tool) to create **optional** dynamic themes based on wallpaper colors. The system defaults to the carefully crafted CypherRiot theme with riot_01.jpg wallpaper, with dynamic theming available as an opt-in toggle in the archriot-control-panel.

## PHASE 1: Color Analysis and Template System

### Objective

Analyze current color usage and create template system for dynamic theming

### Implementation

- **Audit color usage**: Catalog all hardcoded colors across system components
- **Create color templates**: Convert static configs to templated versions for:
    - Hyprland configuration (`hyprland.conf`)
    - Waybar styling (`style.css`)
    - GTK 3/4 themes (`gtk.css`)
    - Terminal colors (Ghostty, Fish shell)
    - Python control panel colors
- **Design color mapping**: Map matugen's Material Design 3 palette to ArchRiot component needs
- **Template engine**: Implement templating system (likely using Go `text/template` or simple string replacement)

### Testing Requirements

1. **Color Audit Test**: Verify all hardcoded colors are identified and cataloged
2. **Template Generation Test**: Ensure templates generate valid configuration files
3. **Fallback Test**: Verify CypherRiot colors work when matugen unavailable
4. **Syntax Test**: Validate all templated configs have correct syntax

## PHASE 2: Matugen Integration

### Objective

Integrate matugen binary and color extraction workflow

### Implementation

- **Package installation**: Add `matugen` to ArchRiot package dependencies
- **Color extraction function**: Create Go function to run matugen on wallpaper files
- **JSON processing**: Parse matugen's JSON output for color values
- **Default behavior**: System starts with CypherRiot theme + riot_01.jpg wallpaper (dynamic theming disabled)
- **Error handling**: Graceful fallback to CypherRiot theme if extraction fails
- **Color validation**: Ensure extracted colors meet accessibility/contrast requirements

### Testing Requirements

1. **Package Test**: Verify matugen installs correctly on target systems
2. **Extraction Test**: Test color extraction on various image formats (jpg, png, webp)
3. **JSON Parse Test**: Validate parsing of matugen output across different images
4. **Fallback Test**: Ensure system works when matugen binary missing or fails
5. **Edge Case Test**: Test with problematic images (very dark, very bright, monochrome)

## PHASE 3: Template Application System

### Objective

Implement system to apply extracted colors to configuration templates

### Implementation

- **Template processor**: Function to replace template variables with extracted colors
- **File generation**: Write templated configs to appropriate locations
- **Backup system**: Preserve original configs before applying dynamic themes
- **Live reload integration**: Trigger config reloads for affected applications
- **Color scheme caching**: Cache generated themes to avoid regeneration on startup
- **Settings preservation**: Save dynamic theming toggle state to existing `background-prefs.json` for upgrade preservation

### Testing Requirements

1. **Template Processing Test**: Verify color substitution works correctly
2. **File Generation Test**: Confirm configs written to correct locations with proper permissions
3. **Backup Test**: Verify original configs are preserved and restorable
4. **Live Reload Test**: Test that applications pick up new colors without restart
5. **Cache Test**: Validate color scheme caching and invalidation logic
6. **Upgrade Preservation Test**: Verify dynamic theming setting survives ArchRiot upgrades

## PHASE 4: Wallpaper System Integration

### Objective

Integrate dynamic theming with existing wallpaper management system

### Implementation

- **Hook integration**: Modify wallpaper change system to trigger theme generation
- **Control panel integration**: Add "Dynamic Theming" toggle (YES/NO) to `archriot-control-panel`
- **CLI integration**: Add theme commands to main ArchRiot binary
- **Background scanning**: Extend background scanner to cache color palettes
- **Theme preview**: Allow users to preview themes before applying

### Testing Requirements

1. **Hook Test**: Verify theme regeneration on wallpaper changes
2. **Control Panel Test**: Test new theming controls in GUI
3. **CLI Test**: Validate command-line theme management
4. **Preview Test**: Ensure theme preview functionality works correctly
5. **Performance Test**: Measure impact on wallpaper switching speed

## PHASE 5: Advanced Features and Polish

### Objective

Implement advanced theming features and user experience improvements

### Implementation

- **Theme variants**: Support light/dark mode variants from same wallpaper
- **Color adjustment**: Allow user fine-tuning of extracted colors (saturation, brightness)
- **Theme persistence**: Save and restore user-customized themes
- **Upgrade integration**: Add dynamic theming toggle to existing `background-prefs.json` preservation system
- **Automatic scheduling**: Support time-based theme/wallpaper combinations
- **Theme sharing**: Export/import theme configurations

### Testing Requirements

1. **Variant Test**: Verify light/dark theme generation works correctly
2. **Adjustment Test**: Test color modification controls and real-time preview
3. **Persistence Test**: Verify themes survive system restarts and updates
4. **Upgrade Integration Test**: Verify dynamic theming setting preserved across ArchRiot upgrades
5. **Scheduling Test**: Test automatic theme changes based on time/conditions
6. **Import/Export Test**: Validate theme sharing functionality

## COMPREHENSIVE TESTING STRATEGY

### Pre-Implementation Testing

- [ ] Research matugen capabilities and limitations
- [ ] Test matugen manually on various wallpaper types
- [ ] Verify Material Design 3 color mapping appropriateness

### Per-Phase Testing

- [ ] Unit tests for each color extraction and template function
- [ ] Integration tests for wallpaper-to-theme pipeline
- [ ] Visual regression tests for theme consistency
- [ ] Performance tests for theme generation speed

### Final Integration Testing

- [ ] End-to-end testing: wallpaper change to full system theme update
- [ ] Multiple wallpaper testing: various image types, colors, compositions
- [ ] Accessibility testing: ensure contrast ratios meet WCAG guidelines
- [ ] Upgrade testing: verify themes persist across ArchRiot updates

### Success Criteria

- [ ] System theme automatically adapts to wallpaper changes within 2 seconds
- [ ] Generated themes maintain visual coherence across all applications
- [ ] Fallback to CypherRiot theme works seamlessly when needed
- [ ] User can fine-tune extracted colors through control panel
- [ ] Theme generation doesn't negatively impact system performance
- [ ] All existing ArchRiot functionality remains unchanged

## FUNCTIONAL SPECIFICATIONS

### Core Requirements

**FR1: Dynamic Color Extraction**

- System MUST extract color palette from current wallpaper using matugen
- System MUST generate Material Design 3 compliant color scheme
- System MUST complete color extraction within 5 seconds for typical wallpapers

**FR2: System-Wide Theme Application**

- System MUST apply extracted colors to all major components:
    - Hyprland window manager (borders, decorations)
    - Waybar status bar (background, text, modules)
    - GTK applications (buttons, headers, selections)
    - Terminal applications (background, foreground, ANSI colors)
    - Control panel interface
- System MUST maintain visual consistency across all components

**FR3: Automatic Theme Updates**

- System MUST automatically regenerate theme when wallpaper changes
- System MUST reload affected applications without requiring logout/restart
- System MUST preserve user workflow during theme transitions

**FR4: Fallback and Recovery**

- System MUST fall back to CypherRiot theme if matugen fails, unavailable, or dynamic theming disabled
- System MUST preserve original configurations for recovery
- System MUST handle corrupted or invalid wallpaper files gracefully

**FR5: User Control**

- System MUST provide "Dynamic Theming" toggle (YES/NO) in archriot-control-panel
- System MUST default to CypherRiot theme with riot_01.jpg wallpaper (dynamic theming OFF)
- System MUST preserve dynamic theming toggle state across ArchRiot upgrades
- System SHOULD allow manual theme regeneration when dynamic theming enabled
- System SHOULD support theme preview before application
- System MAY support user adjustment of extracted colors

### Non-Functional Requirements

**NFR1: Performance**

- Theme generation MUST complete within 5 seconds
- System MUST NOT block wallpaper changes during theme processing
- Memory usage SHOULD NOT increase by more than 50MB during theme generation

**NFR2: Reliability**

- System MUST work with all supported image formats (jpg, png, webp)
- System MUST handle edge cases (monochrome, very dark/bright images)
- Fallback mechanism MUST activate within 2 seconds of matugen failure

**NFR3: Compatibility**

- System MUST maintain backward compatibility with existing ArchRiot installations
- System MUST preserve user customizations not related to colors
- System MUST integrate with ArchRiot's existing upgrade preservation system
- System MUST work across all supported hardware configurations

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
