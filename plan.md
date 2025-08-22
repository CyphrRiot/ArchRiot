# ArchRiot Development Plan

## OVERVIEW

This document tracks development priorities and outstanding tasks for the ArchRiot installer and system components.

## 🎉 COMPLETED TASKS

### ✅ TASK 1: Dynamic Color Theming System - COMPLETED (v2.18)

**Full matugen-based dynamic theming system with complete waybar integration.**

**What Works:**

- ✅ Matugen color extraction from any wallpaper image
- ✅ Complete waybar theming (workspace colors, CPU indicators, accents, etc.)
- ✅ Control panel toggle for enable/disable dynamic theming
- ✅ SUPER+CTRL+SPACE wallpaper cycling with automatic theme updates
- ✅ CLI commands: `--apply-wallpaper-theme` and `--toggle-dynamic-theming`
- ✅ Proper fallback to CypherRiot static colors when disabled
- ✅ Real-time waybar reload without killing the process (SIGUSR2)
- ✅ JSON config preservation and error handling

**Implementation Details:**

- Uses `@define-color` syntax for waybar compatibility
- Colors.css lives in `~/.config/waybar/` directory
- Theming system respects dynamic_theming setting in background-prefs.json
- All component colors (workspace, CPU, memory, accents) update dynamically
- System gracefully handles matugen failures with static color fallback

**User Experience:**

- Toggle dynamic theming in Control Panel
- Change wallpapers via Control Panel or keyboard shortcut
- Colors update instantly across waybar
- Upgrades reset to defaults (documented behavior)

### ✅ TASK 2: Terminal Emoji Fallback System - COMPLETED

**Automatic detection and ASCII fallbacks for emoji support.**

- Automatic terminal emoji capability detection
- ASCII alternatives maintain exact visual spacing
- Test mode with `--ascii-only` flag
- Both logger system and TUI components supported

### ✅ TASK 3: Kernel Upgrade Reboot Detection - COMPLETED

**Fixed reboot handling and continuation flow.**

- Improved reboot detection using prefix matching
- Enhanced continuation workflow for system updates
- Released as v2.10.8

## 🚧 OUTSTANDING TASKS

### TASK 4: System-Wide Dynamic Theming Extension

**PRIORITY: HIGH**

**Problem:** Current dynamic theming only affects waybar. Text editors, GTK applications, terminal themes, and other system components still use static CypherRiot colors.

**Goal:** Extend matugen-based theming to create a cohesive system-wide color experience.

#### Phase 1: Text Editor Integration

**Target Applications:**

- **Zed Editor** - Primary development environment
- **Neovim** - Terminal-based editing
- **VS Code** (if installed) - Backup editor support

**Implementation Strategy:**

- Research Zed theme format and configuration location
- Create dynamic Zed theme templates with matugen color mapping
- Implement Neovim colorscheme generation (likely via lua config)
- Add editor theming to Go theming system

**Expected Files to Modify:**

- `source/theming/theming.go` - Add editor theme generation functions
- New: `config/zed/` - Zed theme templates
- New: `config/nvim/` - Neovim colorscheme templates

#### Phase 2: GTK Application Theming

**Target Applications:**

- **File Manager (Thunar/Nautilus)**
- **System dialogs and notifications**
- **Control panel and GUI applications**

**Implementation Strategy:**

- Generate GTK3/GTK4 CSS themes from matugen colors
- Update `~/.config/gtk-3.0/gtk.css` and `~/.config/gtk-4.0/gtk.css`
- Handle both light and dark theme variants
- Ensure proper contrast ratios for accessibility

#### Phase 3: Terminal and Shell Integration

**Target Components:**

- **Ghostty terminal** - Color palette and background
- **Fish shell** - Syntax highlighting and prompt colors
- **Prompt themes** - Starship or custom prompt integration

**Implementation Strategy:**

- Generate terminal color profiles (16-color palette from matugen)
- Update shell syntax highlighting themes
- Integrate with existing prompt customizations

#### Phase 4: Hyprland Window Manager Integration

**Target Components:**

- **Window borders and decorations**
- **Workspace indicators**
- **Notification styling**

**Implementation Strategy:**

- String replacement in hyprland.conf for border colors
- Update notification daemon themes
- Coordinate with waybar theming for consistent experience

### TASK 5: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

**Problem:** Current Secure Boot implementation disabled due to boot failures and complexity.

**Root Cause:** Uses error-prone manual `openssl`/`efitools` approach instead of reliable `sbctl`.

**Solution Strategy:**

1. **Replace Key Generation** - Use `sbctl create-keys` instead of manual OpenSSL
2. **Fix Enrollment** - Use `sbctl enroll-keys` with proper Setup Mode detection
3. **Reliable Signing** - Replace manual `sbsign` with `sbctl sign` commands
4. **Add Recovery** - Implement `sbctl reset` and failure recovery mechanisms
5. **Re-enable** - Remove disable flag after thorough testing

**Implementation Phases:**

- Phase 1: Key generation replacement and testing
- Phase 2: Setup Mode detection and key enrollment
- Phase 3: Signing system overhaul with pacman hooks
- Phase 4: Recovery mechanisms and error handling
- Phase 5: End-to-end integration testing

## 🔧 TECHNICAL DEBT AND IMPROVEMENTS

### Control Panel JSON Serialization

- Fix GTK boolean serialization issues in background preferences
- Improve error handling for corrupted config files

### Installation Process Optimization

- Review file copying patterns for efficiency
- Consider preservation system expansion for waybar configs

### Testing and Quality Assurance

- Automated testing for theming system across different wallpapers
- VM-based testing for Secure Boot implementation
- Integration testing for upgrade scenarios

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming (Task 4)

- [ ] Zed editor themes update with wallpaper changes
- [ ] Neovim colorschemes follow dynamic colors
- [ ] GTK applications use consistent color palette
- [ ] Terminal color scheme matches system theme
- [ ] Hyprland decorations coordinate with waybar
- [ ] All components respect dynamic theming toggle
- [ ] Fallback to CypherRiot themes when disabled

### Secure Boot Implementation (Task 5)

- [ ] `sbctl status` shows "Secure Boot: Enabled" after setup
- [ ] System boots successfully with custom keys
- [ ] Kernel updates maintain Secure Boot signatures automatically
- [ ] Clear recovery instructions for failed setups
- [ ] Multi-vendor UEFI hardware compatibility

## 📋 DEVELOPMENT PRINCIPLES

**Code Quality:**

- One change at a time, wait for feedback
- Comprehensive testing before claiming completion
- Clear error handling and user guidance
- Modular design with separation of concerns

**User Experience:**

- Consistent theming across all applications
- Simple toggle controls (enable/disable)
- Graceful fallbacks for any failures
- Educational content about security benefits

**Security Focus:**

- Custom Secure Boot keys for maximum control
- No reliance on third-party certificate authorities
- Clear chain of trust documentation
- Recovery mechanisms for all failure scenarios
