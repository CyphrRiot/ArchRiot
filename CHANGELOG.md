# Changelog

All notable changes to ArchRiot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.15] - 2025-07-30

### 🚨 Critical Bug Fix

#### Media Applications Installation Error

- **Fixed "print_status: command not found" error** in media.sh during installation
- **Root cause**: Added function calls to script that doesn't source helper libraries
- **Solution**: Replaced unavailable functions with direct commands and echo statements
- **Impact**: Prevents installation failure at media applications step
- **Added print_status to exported functions** for scripts that do source helpers
- **Maintained spotdl --nocheck functionality** with proper availability

### ⚡ Performance Improvements

#### Installation Speed Optimization

- **Eliminated excessive temp file operations** during module execution
- **Root cause**: Every module created/wrote/deleted temp files unnecessarily
- **Solution**: Direct logging via `tee -a` instead of temp file + cat + rm operations
- **Impact**: Significantly faster installations, especially for already-installed packages
- **Reduced I/O operations** from 4 file operations to 1 per module
- **Maintains all functionality** (output visibility, error handling, complete logging)

This optimization resolves slow package installations (e.g., brave-bin taking 30+ seconds when already installed).

---

## [2.0.14] - 2025-07-30

### 🚀 Major Application Improvements

#### Zed Editor Complete Installation Overhaul

- **Intelligent Vulkan driver detection**: Automatically detects AMD/Intel/NVIDIA GPUs and installs correct Vulkan drivers
- **User-friendly command creation**: Users can now run `zed filename.txt` directly (not just `zed-wayland`)
- **Perfect Wayland integration**: Optimized environment variables and native rendering
- **Upgrade-safe installation**: Handles existing conflicting installations with automatic backups
- **Desktop integration**: Clean application menu entry with no duplicates
- **Functional testing**: Verifies Zed actually launches and works with GPU acceleration
- **Consolidated installation**: Removed duplicate installations, single source of truth

#### Screen Recording (Kooha) Enhancement

- **Comprehensive codec support**: Added x264, x265, libvpx, aom, gifski, gifsicle for complete format coverage
- **Video format support**: MP4 H.265), WebM (VP8/VP9), AV1, and all GStreamer formats
- **GIF creation tools**: High-quality GIF encoding via gifski and optimization via gifsicle
- **Hardware acceleration**: VA-API support for Intel/AMD GPU encoding
- **Professional-grade recording**: All codecs verified working through GStreamer pipeline

#### Spotify Downloader (spotdl) Fix

- **Root cause resolution**: Fixed python-syncedlyrics dependency failing due to API test failures
- **Specialized installation function**: Added `install_aur_nocheck()` for packages with failing tests
- **Bypassed broken tests**: Uses `--mflags "--nocheck"` to skip Musixmatch/Genius API tests
- **Future-proof solution**: Can handle other AUR packages with similar API test issues
- **Comprehensive documentation**: Explains why this approach is needed

### 🛠️ System Stability & Installation Improvements

#### Theme System Verification Enhancement

- **Eliminated false positives**: Fixed "Theme system not configured" appearing when theme IS configured
- **Detailed verification logging**: Shows exactly what's being checked (archriot.conf, backgrounds)
- **Individual component checks**: Verifies each theme component separately with clear feedback
- **Background file validation**: Ensures images actually exist, not just directories
- **Race condition prevention**: Added 2-second delay to prevent upgrade timing issues
- **Smart fix triggering**: Only runs theme setup when actually needed

#### Hardware Detection Improvements

- **Fixed missing lsusb command**: Installs usbutils package to provide lsusb for hardware detection
- **Enhanced Apple keyboard support**: Proper detection and configuration with clear user feedback
- **Better error handling**: Graceful fallback when hardware detection tools unavailable
- **Comprehensive hardware utilities**: Added both usbutils and pciutils for complete detection

#### Process Management Fix

- **Background service persistence**: Fixed services disappearing when installer terminal closes
- **Proper process detachment**: Added nohup and disown to swaybg and other background services
- **Upgrade-safe background services**: Services now survive terminal closure during upgrades
- **Comprehensive process audit**: Systematic review of all background process starts

### 🐛 Bug Fixes

- **Package installation robustness**: Better handling of AUR packages with test failures
- **Installation verification accuracy**: More precise detection of installed components
- **Desktop integration cleanup**: Removed duplicate installations and conflicting files
- **Hardware setup reliability**: Eliminated command not found errors

### 🧹 Code Quality Improvements

- **Specialized installation functions**: Added targeted functions for problematic packages
- **Better error messages**: More descriptive feedback when installations fail
- **Upgrade handling**: Safe file replacement with automatic backups
- **Documentation**: Clear explanations of why special handling is needed

## [2.0.4] - 2025-07-29

### 📚 Documentation Improvements

- **Added CHANGELOG links**: Added prominent changelog badge and table of contents link to README
- **Improved navigation**: Users can now easily find version history in proper location
- **Clean separation**: README stays focused, version history properly organized in CHANGELOG.md
- **Open source best practices**: Follows standard project documentation structure

---

## [1.8.2] - 2025-07-28

### 🔧 Critical Fixes

- **Theming Installation Order**: Renamed theming.sh to 1-theming.sh to fix execution order
- **Installation Reliability**: Fixed theming component loading during setup process

---

## [1.8.1] - 2025-07-28

### 🔧 Critical Fixes

- **Theming System**: Fixed theming system installation permanently
- **Component Loading**: Resolved theming component initialization issues

---

## [1.7.9] - 2025-07-28

### ✨ Enhancements

- **Theming Installation**: Fixed theming installation process
- **Control Panel**: Added pomodoro timer image for enhanced UI
- **Visual Improvements**: Enhanced Control Panel theming components

---

## [1.7.8] - 2025-07-28

### 🔧 Critical Fixes

- **Theming System**: Fixed theming system reliability issues
- **Control Panel Enhancement**: Improved Control Panel functionality and integration

---

## [1.7.7] - 2025-07-28

### 🚨 CRITICAL FIX

- **Control Panel Universal Support**: Fixed Control Panel to work on all systems including laptops
- **Cross-Platform Compatibility**: Resolved system-specific Control Panel launch issues

---

## [1.7.6] - 2025-07-28

### 🚨 CRITICAL FIX

- **Background Installation**: Fixed background installation to finally work on fresh installs
- **Installation Reliability**: Resolved background system setup failures

---

## [1.7.5] - 2025-07-28

### 🔧 Critical Fixes

- **Power Menu**: Fixed Power Menu functionality and integration
- **Script Consistency**: Improved script location consistency across components

---

## [1.7.4] - 2025-07-28

### 🚨 CRITICAL FIX

- **Background Installation**: Fixed background installation system
- **Control Panel Integration**: Enhanced Control Panel functionality and reliability

---

## [1.7.2] - 2025-07-28

### ✨ Major Enhancement

- **Control Panel Integration**: Complete Control Panel integration per specification
- **System Management**: Enhanced system configuration management interface

---

## [1.7.1] - 2025-07-28

### 🚨 CRITICAL FIX

- **Background Installation**: Force background installation during theming setup
- **Theming Reliability**: Ensured background system works on all installations

---

## [1.6.5] - 2025-07-28

### ✨ New Features

- **Intelligent Window Switcher**: Added smart window switcher for lost floating windows
- **Window Management**: Enhanced window recovery and navigation system

---

## [1.6.4] - 2025-07-28

### 🔧 Bug Fixes

- **Brightness Control**: Fixed brightness increments to proper 5% steps
- **System Controls**: Improved brightness adjustment precision

---

## [1.6.3] - 2025-07-28

### 🚨 CRITICAL FIX

- **Mako Configuration**: Removed invalid Mako config options breaking notifications
- **Notification System**: Fixed notification daemon configuration issues

---

## [1.6.2] - 2025-07-28

### 🔧 Critical Fixes

- **Waybar Stability**: Fixed waybar disappearing when exiting installer
- **Installation Process**: Improved installer exit handling

---

## [1.6.1] - 2025-07-28

### 🎯 Major Improvements

- **Package Installation**: Major package installation standardization
- **Plymouth Fixes**: Fixed Plymouth crash issues during installation
- **System Stability**: Enhanced installation reliability

---

## [1.6.0] - 2025-07-28

### 🎯 MAJOR: Installation System Overhaul

- **Installation Output**: Major cleanup of installation output and messaging
- **Session Stability**: Enhanced session stability during installation process
- **User Experience**: Dramatically improved installation feedback and reliability

---

## [1.5.5] - 2025-07-28

### 🔧 Critical Fixes

- **Package Cleanup**: Removed deprecated neofetch causing installation failures
- **Installation Reliability**: Fixed package-related installation issues

---

## [1.5.4] - 2025-07-28

### 🚨 CRITICAL FIX

- **Waybar Validation**: Removed invalid waybar --dry-run that broke upgrades
- **Upgrade System**: Fixed upgrade process validation issues

---

## [1.5.3] - 2025-07-28

### 🚨 CRITICAL HOTFIX

- **Upgrade Dialog**: Fixed broken Install button in upgrade dialog
- **User Interface**: Restored upgrade dialog functionality

---

## [1.5.2] - 2025-07-28

### ✨ Enhancements

- **Brightness Notifications**: Added brightness change notifications
- **Notification System**: Fixed notification stacking issues

---

## [1.5.0] - 2025-07-28

### 🎯 MAJOR: Installation System Improvements

- **Installation Reliability**: Major installation system improvements
- **Critical Fixes**: Multiple critical installation fixes
- **System Stability**: Enhanced overall system installation process

---

## [1.1.78] - 2025-07-21

### 🔧 Bug Fixes

- **Editor Aliases**: Fixed vi/vim editor aliases for system commands
- **Command Line**: Improved editor command integration

---

## [2.0.3] - 2025-07-29

### 🔧 Critical Fixes

- **Restored beautiful window borders**: Fixed borders broken during v2.0.0 theme consolidation
    - Restored original light blue gradient active borders `rgba(89b4fa88) 45deg`
    - Restored subtle dark inactive borders `rgba(1a1a1a60)`
    - Corrected border size back to 1 pixel for elegant appearance
- **Enhanced upgrade notifications**:
    - Added immediate "Launching Upgrade..." feedback for waybar upgrade button
    - Added consistent Control Panel launch notifications from Power Menu
    - Reduced notification timeout to 2 seconds for better UX
- **Fixed user experience inconsistencies**: Unified feedback across all launcher methods

---

## [2.0.2] - 2025-07-29

### 🎨 Boot Logo Enhancement

- **Fixed LUKS boot logo**: Replaced old "Ohmarchy" branding with proper "ArchRiot" logo
- **Enhanced logo design**:
    - Uses Hack Nerd Font Mono Bold for system consistency
    - Larger, more readable fonts (100px + 34px)
    - Perfect background color matching (#191a25)
    - Transparent background for seamless integration
- **Consolidated logo management**: Single source logo file eliminates duplicates
- **Fixed background/lock screen issues**: Corrected theme consolidation path problems
- **Improved boot logo script**: Removed syntax errors and optimized generation logic

---

## [2.0.1] - 2025-07-29

### 🚨 CRITICAL HOTFIX

- **Fixed migrate tool download**: Resolved "text file busy" errors during installation
- **Resolved PipeWire conflicts**: Fixed audio system installation issues
- **Fixed theme consolidation**: Corrected upgrade paths from v2.0.0
- **Background system repair**: Fixed missing backgrounds and lock screen issues
- **Simplified update notifications**: Removed problematic caching system

---

## [2.0.0] - 2025-07-29

### 🎯 MAJOR: Unified Theme System

- **Eliminated theme override system**: Removed maintenance-heavy theme switching
- **Consolidated CypherRiot theme**: Integrated directly into main configurations
- **Removed tokyo-night theme**: Simplified to single, beautiful theme
- **Fixed installation paths**: Streamlined upgrade and installation processes
- **Enhanced system reliability**: Eliminated theme-related conflicts and issues

---

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
