# Changelog

All notable changes to ArchRiot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.3] - 2025-07-21

### 🔍 Installation Reliability Improvements

- **Visible Failure Tracking**: Added comprehensive failure reporting system to installer
    - **NEW**: Installation failures are now impossible to miss with big warning messages
    - **NEW**: Detailed failure log with specific component names and errors
    - **NEW**: Automatic fix commands provided for failed components
    - **SOLUTION**: Prevents silent failures that leave systems partially broken
    - **IMPACT**: Users can immediately identify and fix installation issues

### 📚 Documentation

- **Installer Architecture Documentation**: Added comprehensive INSTALLER.md
    - **PREVENTS**: Future accidental breaking of installer functionality
    - **DOCUMENTS**: Complete installation flow and dependencies
    - **GUIDES**: Safe modification practices for developers

### 🔧 Technical Improvements

- **Enhanced Error Visibility**: Failures are logged and displayed prominently
- **Recovery Commands**: Automatic generation of fix commands for specific failures
- **Maintains Compatibility**: All improvements preserve existing installation behavior

---

## [1.2.1] - 2025-07-21

### 🚨 CRITICAL HOTFIX

- **Fixed Broken Systemd Timer**: Resolved conflicting output directives in version-check service
    - **ROOT CAUSE**: Debugging session left conflicting StandardOutput and ExecStart redirection
    - **IMPACT**: All new installations had broken automatic update timers
    - **SOLUTION**: Removed conflicting logging directives from systemd service configuration
    - **URGENCY**: Critical - affects all fresh installations and upgrades

---

## [1.2.0] - 2025-07-21

### 🎯 MAJOR: Revolutionary Waybar Update System

- **Complete Waybar Notification Redesign**: Three-state update system with always-visible status
    - **NEW**: `󰚰` (pulsing red) for new updates available
    - **NEW**: `󱧘` (steady purple) for updates seen but not upgraded
    - **NEW**: `-` (dark purple) for no updates available (always visible)
    - **SOLUTION**: Fast-loading upgrade dialog via cached version data (3s vs 10s)
    - **SOLUTION**: Eliminates systemd GTK environment issues completely
    - **IMPACT**: Reliable, always-visible update notifications that just work

### 🎨 UI/UX Improvements

- **Waybar Refinements**: Cleaner, more focused interface
    - **REMOVED**: Microphone module (reduced clutter)
    - **REMOVED**: IP addresses from network display (privacy & cleanliness)
    - **IMPROVED**: Tighter spacing between Mullvad VPN icon and location
    - **IMPACT**: Streamlined waybar with essential information only

### 🔧 Technical Improvements

- **Seamless Backward Compatibility**: Existing systemd timers automatically benefit
    - **PRESERVED**: Same `version-check` binary name and systemd service
    - **AUTOMATIC**: Timer switches from broken GTK-in-systemd to working waybar approach
    - **NO ACTION REQUIRED**: Users get improved notifications on next system update

### 🐛 Critical Bug Fixes

- **Missing Mako Configuration**: Fixed missing notification theme in installer
    - **ROOT CAUSE**: Installer lacked `config/mako/` directory for themed notifications
    - **SOLUTION**: Added proper CypherRiot-themed Mako config for all installations
    - **IMPACT**: New installations now have properly themed notifications from day one

---

## [1.1.80] - 2025-07-21

### 🎯 Major System Improvements

- **Waybar Update Notifications**: Complete redesign of upgrade notification system
    - **NEW**: Three-state waybar icon system for update notifications
    - **NEW**: `󰚰` (pulsing red) for new updates available
    - **NEW**: `󱧘` (steady purple) for updates seen but not upgraded
    - **NEW**: `󰸾` (muted blue) for no updates available (always visible)
    - **SOLUTION**: Fast-loading upgrade dialog via cached version data
    - **SOLUTION**: Eliminates systemd GTK environment issues completely
    - **IMPACT**: Reliable, always-visible update status in waybar

### 🔧 Critical Bug Fixes

- **Missing Mako Configuration**: Added missing notification theme configuration to installer
    - **ROOT CAUSE**: Installer was missing `config/mako/` directory entirely
    - **SOLUTION**: Added proper CypherRiot-themed Mako config for installation
    - **IMPACT**: New installations now have proper themed notifications

### ⚡ Performance Improvements

- **Fast Upgrade Dialog Launch**: New `--gui` flag uses cached version data
    - **IMPROVEMENT**: Launch time reduced from ~10 seconds to ~3 seconds
    - **BENEFIT**: Instant response when clicking waybar update icon

### 🔄 Backward Compatibility

- **Seamless Upgrade Path**: Existing systemd timers automatically benefit from new system
    - **PRESERVED**: Same `version-check` binary name and systemd service
    - **AUTOMATIC**: Timer switches from broken GTK-in-systemd to working waybar approach
    - **NO ACTION REQUIRED**: Users get improved notifications on next system update

---

## [1.1.79] - 2025-07-21

### 🔧 Critical Bug Fixes

- **Upgrade Dialog Visibility**: Fixed critical issue where upgrade notifications were marked as "shown" even when invisible to users
    - **ROOT CAUSE**: GTK dialog marked versions as "notified" regardless of user visibility in Hyprland/Wayland environments
    - **ROOT CAUSE**: Window manager focus issues caused dialogs to be created but not visible to users
    - **SOLUTION**: Added user interaction tracking - only mark as notified when user actually clicks buttons
    - **SOLUTION**: Added 30-second timeout with fallback system notifications when GTK dialog fails
    - **SOLUTION**: Enhanced dialog urgency hints and presentation for better visibility
    - **IMPACT**: Users will now always receive upgrade notifications, even with window manager issues

### 🧹 Code Cleanup

- **Consolidated Dialog System**: Removed obsolete `version-update-dialog` standalone script
    - **BENEFIT**: Single consolidated script reduces maintenance complexity
    - **BENEFIT**: All dialog functionality now integrated into `version-check` script
    - **IMPACT**: Existing systemd timers automatically benefit from fixes

---

## [1.1.9] - 2025-07-16

### 🔧 Critical Bug Fixes

- **Sudo Passwordless Setup**: Fixed password prompts during post-install operations
    - **ROOT CAUSE**: User-specific sudo rule needed instead of wheel group dependency
    - **ROOT CAUSE**: Gum installation happened after sudo cleanup, causing password prompt
    - **SOLUTION**: Changed to user-specific sudo rule and moved gum install before cleanup
    - **IMPACT**: Installation now truly passwordless from start to finish

---

## [1.1.8] - 2025-07-16

### 🔧 Critical Bug Fixes

- **Signal Keybinding**: Fixed SUPER+G not launching Signal on fresh installs
    - **ROOT CAUSE**: signal-desktop package was never being installed during setup
    - **SOLUTION**: Added signal-desktop to communication.sh yay installation list
    - **IMPACT**: SUPER+G keybinding now works immediately after installation

---

## [1.0.2] - 2025-01-15

### 🔧 Bug Fixes

- **Waybar Configuration**: Fixed duplicate JSON key in 4workspaces.json
    - Removed duplicate `"title<.*reddit.*>"` entry causing validation warnings
    - All diagnostics now pass clean without warnings

### ✨ Enhancements

- **Base System**: Added `bc` (basic calculator) to core installation
    - Fixes script failures when mathematical calculations are needed
    - Resolves `sync-media` and other script compatibility issues on minimal installations

### 🔍 Quality Improvements

- **Configuration Validation**: Enhanced JSON structure integrity
- **Code Quality**: Eliminated duplicate entries in waybar window-rewrite rules

---

## [1.0.1] - 2025-01-14

### 🔧 Minor Fixes

- **Initramfs errors**: Fixed in v1.0.1 - update if experiencing build issues

---

## [1.0.0] - 2025-01-14

### 🚀 Major Changes

#### Terminal Replacement: Ghostty → Kitty

- **BREAKING**: Replaced Kitty terminal with Ghostty as the default terminal emulator
- Added Ghostty shell integration for enhanced functionality
- Implemented proper theme system integration with `config-file` imports
- Added floating terminal support with `SUPER+SHIFT+RETURN` keybinding

#### Fastfetch Display Fix

- **FIXED**: Fastfetch content truncation in Ghostty terminal windows
- Added 0.1s timing delay in `fish_greeting` to allow terminal sizing
- Implemented `--logo-width 20` parameter for consistent layout
- Ensures beautiful fastfetch display on all terminal launch methods

### ✨ New Features

- **Versioning System**: Added VERSION file and version display in installation
- **Version Command**: New `version` command to check installed version and system status
- **Enhanced Validation**: Updated pre-install validation to check Ghostty components
- **Floating Terminal**: Centered floating terminal window (1200x800) with proper opacity
- **Theme Conversion**: CypherRiot and Tokyo Night themes converted to Ghostty format

### 🔧 Fixed Issues

#### Critical Installation Breaks

- **Desktop Applications**: Fixed `.desktop` files to use `ghostty` instead of `kitty`
    - `About.desktop`, `Activity.desktop`, `nvim.desktop`
- **Management Scripts**: Updated bin scripts to reference Ghostty configs
    - `theme-next`, `validate-system`
- **Waybar Integration**: Fixed right-click terminal commands in Waybar configs
- **Theme System**: Fixed main Ghostty config to import themes instead of hardcoded colors

#### Configuration Updates

- **Hyprland**: Updated window rules for Ghostty class names (`com.mitchellh.ghostty`)
- **Fish Shell**: Added fastfetch timing fix with sleep delay in greeting function
- **Theme Files**: Converted theme files from Kitty syntax to Ghostty syntax
- **Package Installation**: Updated core installation to include `ghostty` and `ghostty-shell-integration`

### 🗑️ Removed

- **Legacy Configs**: Removed `/config/kitty/` directory and configurations
- **Outdated References**: Eliminated all Kitty references from documentation and configs

### 📚 Documentation

- **README**: Updated to reflect Ghostty as default terminal
- **Keybindings**: Added documentation for floating terminal keybinding
- **Theme Docs**: Updated CypherRiot theme documentation
- **Installation Guide**: Enhanced with version information display

### 🔍 Validation & Quality

- **Pre-Install Checks**: Enhanced validation script to verify Ghostty components
- **Component Testing**: Added checks for shell integration and theme files
- **Installation Verification**: Improved post-install validation for Ghostty functionality
- **Error Prevention**: Fixed multiple breaking issues that would prevent smooth installation

### 🎨 Visual Improvements

- **Terminal Theming**: Proper theme integration with dynamic color switching
- **Window Management**: Enhanced floating window rules and positioning
- **Opacity Settings**: Consistent transparency across terminal instances
- **Font Configuration**: Maintained CaskaydiaMono Nerd Font with proper sizing

### ⚡ Performance

- **Startup Speed**: Optimized terminal launch with proper initialization timing
- **Theme Switching**: Improved theme reload functionality for Ghostty
- **Resource Usage**: Maintained lightweight footprint while adding features

---

## Technical Details

### Package Changes

- **Added**: `ghostty`, `ghostty-shell-integration`
- **Removed**: `kitty` (replaced by Ghostty)

### Configuration Structure

```
~/.config/ghostty/
├── config                    # Main Ghostty configuration
└── current-theme.conf        # Symlinked theme file

~/.config/archriot/current/
└── theme/
    └── ghostty.conf          # Theme-specific colors
```

### New Commands

- `version` - Display version and system information
- `SUPER+SHIFT+RETURN` - Launch floating centered terminal

### Breaking Changes

- All terminal-related keybindings now use Ghostty instead of Kitty
- Theme files use Ghostty syntax (palette = N=#rrggbb)
- Desktop application launchers updated to use Ghostty

---

## Migration Notes

Users upgrading from a Kitty-based installation should:

1. Run the installation normally - Ghostty will be installed automatically
2. Themes will be converted to Ghostty format during installation
3. All keybindings remain the same (`SUPER+RETURN` for terminal)
4. Floating terminal is a new feature accessible via `SUPER+SHIFT+RETURN`

For a complete fresh installation, run:

```bash
curl -fsSL https://archriot.org/setup.sh | bash
```

---

## Verification

After installation, verify everything is working:

```bash
version                           # Check installation status
fastfetch                         # Verify display works correctly
ghostty --version                 # Confirm Ghostty is installed
```
