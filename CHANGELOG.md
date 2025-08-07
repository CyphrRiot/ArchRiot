# Changelog

All notable changes to ArchRiot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.6.8] - 2025-01-08

### 📋 ENHANCED UPGRADE GUI FORMATTING

#### Multi-Line Commit Message Display

- **ENHANCED**: Upgrade GUI now shows first 2-3 lines of commit messages instead of just titles
- **IMPROVED**: Proper line breaks and spacing between commit entries for better readability
- **REFINED**: Minimal left margin usage with hanging indentation for wrapped text
- **OPTIMIZED**: Scrollable content area handles longer commit descriptions gracefully

#### User Experience Improvements

- **DETAILED**: Users now see comprehensive commit information including:
    - Full commit titles with version numbers
    - Key feature descriptions and fixes
    - Critical technical details about changes
- **READABLE**: Clean formatting with proper spacing between each commit entry
- **INFORMATIVE**: No more guessing what's in an update - full context provided

#### Technical Implementation

- **SMART**: Automatic line wrapping with proper hanging indentation
- **EFFICIENT**: Preserves original commit message formatting and structure
- **RESPONSIVE**: Content adapts to dialog width while maintaining readability
- **SCALABLE**: Handles both short and long commit messages elegantly

## [2.6.7] - 2025-01-08

### 🔧 CRITICAL FIX: Graphics Configuration Installation

#### Immediate Environment Variable Loading

- **FIXED**: Graphics environment variables not loading without session restart
- **RESOLVED**: AMD users experiencing Thunar artifacts after installation/upgrade
- **IMPROVED**: Graphics configuration now applies immediately during installation
- **ADDED**: `systemctl --user set-environment` commands for instant variable loading
- **ELIMINATED**: Need to logout/login for graphics fixes to take effect

#### Installation Process Improvements

- **FIXED**: File conflict where both source and target graphics configs were copied
- **REFINED**: Base installation now excludes GPU-specific graphics files
- **ENHANCED**: Hardware detection section handles all GPU-specific configuration
- **OPTIMIZED**: Clean separation between generic and hardware-specific configs

#### Technical Details

- **AMD Systems**: `GSK_RENDERER=gl` and `LIBGL_ALWAYS_SOFTWARE=1` load immediately
- **NVIDIA Systems**: Hardware acceleration settings apply without restart
- **Intel Systems**: Minimal configuration loads instantly
- **All Systems**: Graphics fixes now work immediately after installation

## [2.6.1] - 2025-01-08

### 🔧 Optimization

#### Build System

- **Reduced binary size by 30%**: Optimized Makefile build flags
    - Added `-ldflags="-s -w"` to strip debug symbols
    - Disabled CGO with `CGO_ENABLED=0` for better portability
    - Added `-trimpath` to remove build paths from binary
    - Binary size reduced from 5.2MB to 3.6MB
    - Added `ultra` target for maximum compression with UPX support

## [2.6.5] - 2025-01-10

### 🔧 Critical Fixes

#### Fixed Issues

- **CRITICAL**: Fixed system-wide text selection failure in all GTK applications
    - **Root Cause**: wl-clip-persist package was intercepting and breaking text selection events
    - **Impact**: Users could not select text with mouse or Ctrl+A in GNOME Text Editor, Control Panel, and other GTK apps
    - **Solution**: Removed wl-clip-persist from installer and added cleanup commands
    - **Affected**: All systems with wl-clip-persist installed (introduced ~5 days ago)

#### New Features

- **Proton Docs Integration** - Added Fuzzel launcher entry for quick access to Proton documentation
- **Text Editor Font** - Changed default font from iA Writer Mono to Hack Nerd Font Mono
- **btop Theme** - Set CypherRiot theme as default for all users instead of tokyo-night

#### Removed

- **AbiWord** - Removed package and all references (desktop entries, waybar icons, README mentions)

## [2.6.4] - 2024-12-16

### 🎨 Background Management System

#### New Features

- **Control Panel Background Widget** - Added comprehensive background management to ArchRiot Control Panel
- **Custom Background Support** - Users can now add personal wallpapers via file chooser dialog
- **Background Slider** - Visual slider to select from all available backgrounds with immediate preview
- **Visual Indicators** - System backgrounds numbered (1-20), user backgrounds labeled (U1, U2, etc.)
- **Upgrade Preservation** - User backgrounds and preferences survive ArchRiot upgrades
- **Background Persistence** - Selected backgrounds persist across sessions and reboots

#### Technical Improvements

- **Unified Slider System** - Eliminated code duplication with modular slider creation function
- **Natural Sorting** - Proper numerical ordering for riot_01.jpg, riot_02.jpg, etc.
- **GTK4 File Chooser** - Modern file dialog with image filtering and proper sizing
- **Process Detachment** - Improved background service management for reliability

#### User Interface

- **Workspace Colors** - Refined workspace indicator colors for better visibility
- **Control Panel Optimization** - Improved layout, spacing, and removed redundant elements
- **Header Styling** - Better version display and font sizing

#### Background Management

- **System Backgrounds**: `~/.local/share/archriot/backgrounds/` (managed by installer)
- **User Backgrounds**: `~/.config/archriot/backgrounds/` (preserved during upgrades)
- **Remove Custom Backgrounds**: `rm ~/.config/archriot/backgrounds/*` or `rm -rf ~/.config/archriot/backgrounds/`

### Fixed

- **Audio Widget Removal** - Removed non-functional audio management widget
- **Slider Snapping** - All sliders now snap properly to discrete values
- **Background Ordering** - System backgrounds appear first, user backgrounds at end

---

## [2.6.2] - 2024-12-09

### Fixed

- Fixed installer error when copying single files to directory targets
- Resolved "is a directory" error for bin/upgrade-system installation
- Improved path handling for custom target configurations

## [2.6.1] - 2024-12-09

### 📋 DOCUMENTATION & FORMATTING FIXES

#### README.md Improvements

- **FIXED**: Removed 8 instances of non-functional `{:target="_blank"}` markdown syntax
- **IMPROVED**: GitHub rendering by cleaning up all target attribute remnants
- **CORRECTED**: Header hierarchy - Plymouth section now properly nested under VM guidance
- **ENHANCED**: Overall markdown formatting for better GitHub compatibility

#### Technical Notes

- **GitHub Compatibility**: Target attributes are stripped by GitHub for security, making the syntax non-functional
- **Formatting Consistency**: All section headers now follow proper markdown hierarchy
- **User Experience**: Links work correctly without broken syntax cluttering the README

## [2.6.0] - 2025-01-08

### 🚀 MAJOR RELEASE: INTELLIGENT GRAPHICS & DYNAMIC UPGRADE SYSTEM

#### GPU-Specific Graphics Configuration System

- **NEW**: Automatic GPU detection and optimal configuration during installation
- **CREATED**: `graphics-amd.conf` - Forces software rendering to eliminate AMD artifacts
- **CREATED**: `graphics-nvidia.conf` - Hardware acceleration with NVIDIA optimizations
- **CREATED**: `graphics-intel.conf` - Minimal configuration for Intel integrated graphics
- **FIXED**: AMD Barcelo GPU flashing and artifacts in Thunar and GTK4 applications
- **IMPROVED**: Installation process automatically selects correct graphics configuration
- **ELIMINATED**: Non-functional shell conditionals in environment.d files

#### Dynamic Upgrade GUI with Real-Time Commit Information

- **REVOLUTIONARY**: Upgrade dialog now shows actual commit messages instead of static text
- **ADDED**: GitHub API integration to fetch recent changes between versions
- **IMPROVED**: Compact single-line version display (`Update: v2.5.x → v2.6.0`)
- **ENHANCED**: Proper text wrapping with hanging indentation for bullet points
- **ADDED**: Scrollable content area to handle long commit lists
- **FIXED**: Deprecated GTK font styling with modern CSS approach
- **OPTIMIZED**: Intelligent fallback system when version tags unavailable
- **REMOVED**: Auto-timeout behavior for better user experience

#### Technical Excellence Improvements

- **MODERNIZED**: GTK styling system eliminates deprecation warnings
- **ROBUST**: Multiple fallback mechanisms for network/API failures
- **INTELLIGENT**: Hardware detection works across AMD/Intel/NVIDIA systems
- **RELIABLE**: Graceful degradation when remote services unavailable
- **USER-FRIENDLY**: No more static upgrade messages, real information every time

## [2.5.19] - 2025-01-08

### 🎮 GPU-SPECIFIC GRAPHICS CONFIGURATION

#### Intelligent Graphics Driver Detection and Configuration

- **NEW**: GPU-specific graphics configuration system for optimal compatibility
- **ADDED**: Automatic GPU detection during installation (AMD/Intel/NVIDIA)
- **CREATED**: `graphics-amd.conf` - Forces software rendering to fix AMD artifacts and flashing
- **CREATED**: `graphics-nvidia.conf` - Hardware acceleration with NVIDIA optimizations
- **CREATED**: `graphics-intel.conf` - Minimal configuration for Intel integrated graphics
- **FIXED**: AMD Barcelo GPU flashing/artifacts in Thunar and GTK4 applications
- **IMPROVED**: Installation process automatically selects correct graphics config
- **REMOVED**: Static graphics.conf with non-functional shell conditionals
- **IMPACT**: Eliminates graphics issues across different GPU vendors automatically

#### Technical Details

- **AMD Systems**: Uses `LIBGL_ALWAYS_SOFTWARE=1` to force stable software rendering
- **NVIDIA Systems**: Maintains hardware acceleration with threading optimizations
- **Intel Systems**: Conservative approach with minimal changes to preserve performance
- **Detection**: Uses `lspci` GPU vendor detection during hardware setup phase

## [2.5.7] - 2025-01-08

### 🚨 CRITICAL SUDO CHECK FIX

#### Emergency Patch for False Positive Sudo Detection

- **CRITICAL FIX**: Passwordless sudo check giving false positives due to cached credentials
- **ROOT CAUSE**: `testSudo()` function succeeded when sudo timestamp cache was active, even without passwordless sudo configured
- **SOLUTION**: Clear sudo timestamp cache with `sudo -k` before testing actual passwordless sudo configuration
- **IMPACT**: Prevents installation failures when sudo cache expires during installation process
- **AFFECTED**: All users who recently ran sudo commands before installing ArchRiot
- **SYMPTOM**: Installer showed "✅ Passwordless sudo is already working" but failed later during installation

#### Technical Details

- Added `exec.Command("sudo", "-k").Run()` to clear timestamp cache before testing
- Ensures `sudo -n true` tests actual sudoers configuration, not cached credentials
- Eliminates false positives that caused delayed installation failures

## [2.5.4] - 2025-01-13

### 🚨 CRITICAL AUTO-LOGIN FIX

#### Emergency Patch for Broken Auto-Login

- **CRITICAL FIX**: Auto-login broken during installation causing endless `--noclear (automatic login)` loop
- **ROOT CAUSE**: Installer using `$SUDO_USER` variable which is unreliable during installation context
- **SOLUTION**: Changed to `$(whoami)` for proper user detection in getty configuration
- **IMPACT**: Prevents systems from being locked out at login prompt after fresh installation
- **AFFECTED**: All v2.5.0+ installations where auto-login was configured
- **REGRESSION**: Introduced in v2.5.0 during shell→Go installer transition

This was a **HIGH SEVERITY** regression that could render systems unusable after installation.

## [2.5.1] - 2025-01-08

### 🌡️ Critical Hardware Compatibility Fix

#### Temperature Sensor Auto-Detection System

- **FIXED**: Missing temperature sensor configuration during installation (HIGH PRIORITY from ROADMAP)
- **ADDED**: `setup-temperature` script execution to hyprland installation phase
- **RESOLVED**: Waybar temperature module failures on different hardware configurations
- **IMPACT**: Temperature monitoring now works correctly across Intel/AMD CPUs and different motherboards
- **TECHNICAL**: Auto-detects coretemp, x86_pkg_temp, or thermal_zone fallback paths instead of hardcoded hwmon6/hwmon1
- **RESULT**: 100% hardware compatibility for temperature monitoring in waybar status bar

### 🔧 Installation System Improvements

- **ENHANCED**: Installation process now includes automatic hardware-specific temperature sensor detection
- **VERIFIED**: Tested and confirmed working on fresh installations with config removal

## [2.5.0] - 2025-08-03

### 🧹 Script Cleanup and System Optimization

- **REMOVED**: 8 unused/orphaned scripts cluttering the bin directory
    - `apple-display-brightness` - Niche utility with no keybindings or usage
    - `asciinema-to-gif` - Terminal recording converter, unused
    - `fingerprint-setup` - Problematic fingerprint authentication script
    - `fix-touchpad-config` - Touchpad configuration utility, no active usage
    - `generate-boot-logo.sh` - Redundant Plymouth logo generator (static PNG preferred)
    - `signal-wayland` - Comprehensive Wayland script superseded by signal-launcher.sh
    - `update` - Outdated git-based update script (replaced by binary installer)
    - `volume-osd` - Volume notification script (superseded by Volume.sh + mako)
- **FIXED**: Waybar hypridle module by moving `toggle-idle` to `scripts/Hypridle.sh`
    - **ADDED**: Proper JSON status output for waybar integration
    - **ADDED**: Support for both `status` and `toggle` arguments
    - **RESULT**: Functional idle lock toggle in waybar status bar
- **CLEANED**: Help system references to deleted commands
- **REMOVED**: 16 unused waybar custom modules with broken script references
    - `custom/weather`, `custom/file_manager`, `custom/tty`, `custom/browser`
    - `custom/settings`, `custom/cycle_wall`, `custom/hint`, `custom/hypridle`
    - `custom/light_dark`, `custom/cava_mviz`, `custom/playerctl`, `custom/reboot`
    - `custom/quit`, `custom/updater`, `custom/tray-open`, `custom/tray-close`
    - **ROOT CAUSE**: Copy-pasted waybar config from different rice with non-existent scripts
    - **RESULT**: All modules referenced scripts in `/config/hypr/scripts/` that never existed
- **FIXED**: `custom/power` module removed broken right-click reference to `ChangeBlur.sh`
- **IMPACT**: Cleaner waybar config with only functional, actively used modules

### 🎉 MAJOR: Complete YAML Migration and Architectural Overhaul

#### 🏗️ Go Binary Installer Revolution

- **REPLACED**: Fragile shell script installer with compiled Go binary for 100% reliable installations
- **ELIMINATED**: Shell injection vulnerabilities, buffer overflows, and undefined behavior
- **ADDED**: Proper error handling, structured logging, and predictable behavior across environments
- **IMPLEMENTED**: Atomic operations with rollback capabilities - no more broken partial installs
- **RESULT**: Architectural superiority over traditional shell script installations

#### 📄 YAML Configuration System

- **MIGRATED**: All installation logic from shell scripts to declarative YAML configuration
- **CENTRALIZED**: Single `install/packages.yaml` file defines entire system
- **ADDED**: Intelligent module dependency resolution that shell scripts couldn't provide
- **IMPLEMENTED**: Type-safe configuration eliminating shell script parsing nightmares
- **STRUCTURED**: Clean separation of packages, configurations, and commands with proper dependency ordering

#### 🗂️ Repository Structure Reorganization

- **SEPARATED**: Source code (`source/`) from installation files (`install/`)
- **INCLUDED**: Pre-built installer binary (`install/archriot`) for immediate use
- **CONSOLIDATED**: All scripts moved to single location (`~/.local/share/archriot/config/bin/`)
- **REMOVED**: 6,092 lines of obsolete shell script code
- **CLEANED**: Proper build system with Makefile for development workflow

#### 🔧 Critical Bug Fixes Resolved

- **FIXED**: Waybar JSON syntax errors (duplicate keys, trailing commas) that broke status bar
- **FIXED**: All broken script paths in waybar modules now reference correct locations
- **FIXED**: Hidden desktop files properly managed - no more unwanted menu entries (btop++, About Xfce)
- **FIXED**: Duplicate application entries (Zed override, system monitor hiding)
- **RESTORED**: Missing Media Player (mpv.desktop) and proper application visibility
- **ADDED**: Missing packages for screen recording (kooha + GStreamer codecs)
- **RESOLVED**: XCompose file path issues and waybar script permissions

#### 📦 Package System Improvements

- **COMPLETED**: All 23 config directories properly handled in YAML configuration
- **ADDED**: Missing system packages (btop, hyprshot, ufw, linux-firmware, efibootmgr)
- **IMPLEMENTED**: Video group permissions for proper screen recording functionality
- **ORGANIZED**: Clean module structure with proper dependency management
- **OPTIMIZED**: Intelligent package installation with conflict resolution

#### 🧪 Validation and Testing System

- **BUILT-IN**: Comprehensive validation directly in Go binary (`--validate` flag)
- **ENHANCED**: YAML configuration integrity and module dependency validation
- **ADDED**: Configuration conflict detection and rollback verification
- **IMPLEMENTED**: Performance checks for memory, disk space, and network connectivity
- **CREATED**: New test system (`test/test-riot`) with proper backup and cleanup

#### 🛠️ Development and Maintenance

- **ADDED**: Proper build system with `make`, `make install`, `make test` targets
- **IMPLEMENTED**: Clean development workflow with source/install separation
- **CREATED**: Idempotent installation system - safe to re-run installer
- **ESTABLISHED**: Single source of truth for all system configuration
- **FUTURE-PROOFED**: Maintainable architecture for ongoing development

### 📊 Migration Impact

- **Files Changed**: 97 files modified
- **Code Removed**: 6,092 lines of obsolete shell scripts
- **Code Added**: 661 lines of clean Go code and YAML configuration
- **Architecture**: Complete migration from shell scripts to Go/YAML system
- **Reliability**: 100% reliable installations with complete system state management
- **Maintainability**: Clean, type-safe configuration system for future development

### ⚠️ Breaking Changes

- **Removed**: Legacy shell script installer (`install.sh`)
- **Removed**: Old validation script (`validate.sh`)
- **Changed**: All management commands now use absolute paths to new binary
- **Updated**: Installation method remains the same but uses new architecture internally

### 🚀 Upgrade Instructions

No changes required for users - same installation command works with new architecture:

```bash
curl -fsSL https://archriot.org/setup.sh | bash
```

The new system automatically handles upgrades and maintains backward compatibility.

## [2.2.4] - 2025-08-01

### 🔐 Privacy & System Cleanup

#### Privacy Violations Removed

- **REMOVED**: Microsoft Visual Studio Code (415 MB) - Telemetry violation not part of ArchRiot
- **REMOVED**: crush-bin AI tool - Making unauthorized API calls to x.ai
- **REMOVED**: 2+ GB of debugging artifacts from development sessions
- **CLEANED**: Docker, VLC, 7zip, npm, qemu, audacity, and other unwanted packages

#### Document Viewer Optimization

- **REMOVED**: Evince and Sushi (102 MB) - Redundant with Papers reference manager
- **STREAMLINED**: Single document solution using Papers for all PDF viewing
- **FIXED**: Thunar-based file management (sushi was Nautilus-specific)

#### Waybar Interface Fixes

- **FIXED**: Audio volume icons - Now shows proper speaker icons at all volume levels
- **FIXED**: Large gap between package update and power menu modules
- **IMPROVED**: Consistent waybar spacing and visual alignment

#### Extended Lock Screen Fixes

- **FIXED**: AMD GPU window positioning bug for extended lock sessions (10.5min+)
- **ADDED**: Automatic window recovery for DPMS wake events
- **ENHANCED**: Dedicated DPMS wake fix script with proper timing (2s delay vs 0.5s)
- **RESULT**: Complete fix for disappearing windows after any lock/unlock scenario

#### Installation System Improvements

- **ADDED**: Package database sync before installation begins
- **ENHANCED**: Prevents version conflicts and package installation issues
- **IMPROVED**: More reliable installation process

### 📦 Package Management

- **RESTORED**: lsof utility (accidentally removed with VS Code)
- **VERIFIED**: All remaining packages are official ArchRiot components
- **AUDITED**: System now adheres to "no telemetry, no data collection" promise

---

## [2.2.0] - 2025-01-31

### 🖥️ VM Environment Compatibility

#### Critical VM Installation Fixes

- **FIXED**: Missing multilib repository causing package installation failures in VMs
    - **AUTO-DETECTION**: Automatically enables multilib repository when missing
    - **DATABASE SYNC**: Added database synchronization before package installations
    - **RESULT**: Eliminates "database file for 'multilib' does not exist" errors

#### GPU Detection Improvements

- **ELIMINATED**: Interactive GPU selection prompts that hang in non-interactive environments
    - **SMART FALLBACK**: Uses vulkan-swrast software rendering for VMs/unknown hardware
    - **PRESERVED**: Existing GPU detection logic for NVIDIA/AMD/Intel hardware
    - **RESULT**: No more hanging prompts in VMs or automated installations

#### Installation Experience Enhancements

- **IMPROVED**: Reboot prompt defaults to "No" for safer automation
- **FIXED**: Text corruption after reboot prompt with proper screen cleanup
- **ENHANCED**: Package installation reliability with automatic repository configuration

---

## [2.1.9] - 2025-01-31

### 🐛 Progress Display Fixes

#### Installation Interface Improvements

- **REVERTED**: Broken reserved-line progress system that caused overlapping text
- **FIXED**: Completion screen clearing properly without artifacts
- **REMOVED**: Initial progress bar that created display mess before credentials
- **PRESERVED**: Startup messages by removing premature screen clearing

---

## [2.1.8] - 2025-01-31

### 🐛 Critical Bug Fixes

#### Installation System Stability

- **FIXED**: Infinite loop in critical module failure dialog causing 2.4GB log files
    - **ISSUE**: Interactive prompts in non-interactive environments (VMs, automated installs) caused endless loops
    - **SOLUTION**: Critical module failures now exit cleanly with clear error messages
    - **IMPACT**: Prevents installation hangs and massive log file generation

#### Path Resolution Improvements

- **ELIMINATED**: Remaining `../../` relative path navigation in theming system
    - **REPLACED**: Background path resolution now uses proper `$HOME/.local/share/archriot/backgrounds`
    - **RESULT**: More reliable path handling across different execution contexts

---

## [2.1.7] - 2025-01-31

### 🚀 Installation System Optimizations

#### Core Installer Performance Improvements

- **OPTIMIZED**: `install.sh` main installer with systematic efficiency improvements
    - **Eliminated multiple tee operations**: Replaced process spawning in `log_message()` with direct file operations
    - **Simplified module discovery**: Direct execution approach eliminating 49 lines of complex array building logic
    - **Centralized hardcoded paths**: All path constants defined in configuration section for better maintainability
    - **IMPACT**: Faster installation startup and reduced I/O overhead during module execution

#### Process Management Cleanup

- **REMOVED**: Orphaned `config/waybar/waybar.sh` script with improper process detachment
    - **ISSUE**: Script existed but was never called, contained background processes without `nohup & disown`
    - **AUDIT**: Verified all critical services (swaybg, waybar, mako) properly detached in installer
    - **RESULT**: Clean codebase with consistent background process management

### 🎨 AMD Graphics Compatibility

#### Progressive Graphics Rendering System

- **IMPLEMENTED**: AMD-specific graphics artifact fixes for GTK4 applications
    - **PRIMARY**: `GSK_RENDERER=gl` for stable GL rendering across all systems
    - **AMD OPTIMIZATIONS**: `mesa_glthread=true` and Mesa version overrides for better compatibility
    - **PROGRESSIVE FALLBACKS**: Cairo renderer and software rendering options available if needed
    - **IMPACT**: Eliminates Thunar and other GTK4 visual artifacts while preserving hardware acceleration

### 📋 Documentation & Planning Updates

- **UPDATED**: ROADMAP.md with completed optimizations and architectural proposal
- **PROPOSED**: YAML-based declarative package system for v2.2.0 development
- **CLEANED**: Redundant roadmap sections, focused on current development priorities

### 🔧 Code Quality Improvements

- **Consistent Installation Logic**: Unified approach to background service management
- **Maintainable Path Handling**: Centralized configuration for all installer paths
- **Reduced Code Duplication**: Eliminated redundant operations across core installer

## [2.1.6] - 2025-01-31

### 🚨 Critical System Fixes

#### AMD DPMS Window Management Bug - SOLVED

- **FIXED**: Floating windows disappearing after DPMS wake events on AMD systems
    - **ROOT CAUSE**: Windows positioned beyond screen boundaries (coordinates >4000px) after monitor wake
    - **SOLUTION**: Automatic off-screen window detection and recovery system
    - **IMPLEMENTATION**: `fix-offscreen-windows.sh` script with focus-then-center approach
    - **INTEGRATION**: Automatic recovery on DPMS wake + manual SUPER+SHIFT+TAB keybinding
    - **IMPACT**: All floating windows (Signal, Calculator, System Monitor, etc.) now survive screen lock/wake cycles

#### Installation System Architecture Cleanup

- **ELIMINATED**: 23+ instances of confusing `"$script_dir/../../"` relative path navigation
- **REPLACED**: All installation scripts now use clean `~/.local/share/archriot/` absolute paths
- **IMPROVED**: Installation code readability and maintainability across entire system
- **FILES AFFECTED**: communication.sh, productivity.sh, utilities.sh, post-desktop/01-config.sh

### 🧹 Code Quality Improvements

- **Standardized installation path patterns** across all modules
- **Enhanced installation system reliability** with direct path references
- **Improved debugging experience** with clear, readable installation logic

---

## [2.1.5] - 2025-07-30

### 🔧 Hardware Detection Fixes

#### Critical Pipeline Bugs Fixed

- **Fixed Intel graphics detection** - `lspci | grep -qi intel | grep -qi -E 'vga|3d|display|graphics'` → `lspci | grep -i intel | grep -qi -E 'vga|3d|display|graphics'`
- **Fixed Apple keyboard detection** - `lsusb | grep -qi apple | grep -qi keyboard` → `lsusb | grep -i apple | grep -qi keyboard`
- **Root cause**: Silent first grep (`-q`) in pipeline broke second grep - no input to search

#### Impact

- **Intel graphics drivers** now install properly on Intel systems
- **Apple keyboard function keys** now configure correctly on Apple keyboards
- **Hardware detection** works as intended across all supported hardware

#### Technical Details

- **Pipeline issue**: Using `grep -q` (quiet) as first command in pipeline kills output for subsequent commands
- **Solution**: Remove `-q` from first grep, keep `-q` on final grep for boolean result
- **Affected systems**: Any ArchRiot installation on Intel graphics or with Apple keyboards

This release fixes critical hardware detection that was silently failing due to broken shell pipelines.

## [2.1.4] - 2025-07-30

### 🗂️ Backup System Consolidation & Performance Optimization

#### Centralized Backup System

- **Single backup location** - All backups now stored in `~/.archriot/backups/` instead of scattered directories
- **Smart cleanup** - Automatically keeps only 3 most recent backups, removes older ones
- **Unified API** - All installation scripts use same backup functions for consistency
- **Automatic legacy cleanup** - Removes old scattered backup directories on upgrade
- **Metadata tracking** - Each backup includes creation info, manifest, and restoration details
- **Space efficient** - Eliminates duplicate backups and wasted disk space

#### Installation Performance Improvements

- **Consolidated service restarts** - Single restart at end instead of multiple restarts during installation
- **Eliminated screen flashing** - No more Hyprland/Waybar restarts interrupting user workflow
- **Fixed timing calculation** - Installation duration now displays correctly instead of "0m 0s"
- **Reduced process conflicts** - Services no longer fight over resources during installation
- **Professional experience** - Smooth installation without desktop environment disruption

#### Legacy Cleanup

- **Removed Ohmarchy references** - Eliminated obsolete cleanup messages for non-existent legacy system
- **Backup directory cleanup** - Automatic removal of scattered backup directories from previous versions
- **Code consolidation** - Unified backup and restart logic across all installation modules

#### Technical Architecture

- **New backup manager** - `install/lib/backup-manager.sh` provides centralized backup functionality
- **Export fixes** - Proper variable export ensures timing calculations work correctly
- **Conflict resolution** - Separate timing variables prevent conflicts between different installer components
- **Silent operation** - All backup operations logged to files, no console spam

### Impact

This release significantly improves the installation experience by eliminating chaotic service restarts and consolidating the backup system. Users get a professional, smooth installation with intelligent backup management and accurate progress reporting.

## [2.1.2] - 2025-07-30

### 🔧 Minor Improvements

#### Notification Duration Enhancement

- **Extended upgrade notification** - "Starting ArchRiot Upgrade" notification now displays for 3 seconds instead of 2 seconds
- **Better user feedback** - Users have more time to see the upgrade launch confirmation
- **Improved UX** - More consistent with other system notifications timing

#### Setup Performance Optimization

- **Faster downloads** - Added `--depth 1` to git clone commands in setup.sh for shallow clones
- **Reduced download size** - Only downloads latest commit instead of full git history
- **Improved installation speed** - Significantly faster for users on slower connections
- **Same functionality** - No impact on installation reliability or features

These minor improvements enhance the user experience with better feedback timing and faster installation performance.

## [2.1.1] - 2025-07-30

### 🚨 Critical Hotfix

#### Pomodoro Timer Notification Spam Fix

- **Fixed notification spam during installation** - Pomodoro timer was sending "Timer Reset" notifications continuously during system updates
- **Root cause**: `reset_state()` function always sent notifications, even during automatic resets from service restarts
- **Solution**: Added `notify` parameter to `reset_state()` - only sends notifications when explicitly reset by user
- **Impact**: Eliminates notification spam that made installations unusable

#### Technical Details

- **Modified**: `bin/scripts/waybar-tomato-timer.py`
- **Change**: `reset_state(notify=False)` for automatic resets, `reset_state(notify=True)` for user-initiated resets
- **Affected scenarios**: Installation, system service restarts, Waybar reloads
- **User experience**: Silent automatic resets, notifications only for intentional user actions

This critical hotfix resolves the installation experience regression introduced in v2.1.0.

## [2.1.0] - 2025-07-30

### 📊 Modern Progress Bar System

#### Revolutionary Installer Experience

- **Static Progress Bar** - Beautiful real-time progress display that updates in place, no more scrolling walls of text
- **Friendly Module Names** - Shows "Desktop Environment" instead of cryptic "install/desktop/apps.sh" paths
- **Background Logging** - All verbose output redirected to log files, only progress and errors visible on screen
- **Accurate Progress Tracking** - Real percentages based on actual module count (32 install scripts)
- **Professional Interface** - Clean, modern installer experience with proper visual consistency
- **User Interaction Support** - Seamlessly handles Git credentials and reboot prompts without breaking progress display
- **Error Surfacing** - Critical errors still appear on console when user attention is needed
- **Beautiful Completion Summary** - Rich statistics display with timing, success rates, and module status

#### Technical Architecture

- **Progress Bar Library** - New `install/lib/progress-bar.sh` system with modular design
- **Silent Execution** - All package installations and system commands run silently in background
- **Proper Time Tracking** - Fixed duration calculations and progress percentage accuracy
- **Fallback Support** - Graceful degradation if progress system unavailable

#### Visual Improvements

- **Consistent Separator Lines** - Top border matches progress bar width for perfect alignment
- **Color-coded Status** - Green for success, red for errors, cyan for active progress
- **Real-time Updates** - Module counter and percentage update as installation progresses
- **Clean Layout** - Organized display with proper spacing and visual hierarchy

This represents a complete overhaul of the installation experience, transforming it from a verbose, overwhelming process into a clean, professional, and user-friendly system.

## [2.0.19] - 2025-07-30

### 🍅 Pomodoro Timer & System Enhancements

#### Enhanced Pomodoro Timer System

- **Added comprehensive notifications** - Beautiful purple-bordered notifications for all timer state changes
- **Added hotkey support** - SUPER + comma to start/pause timer, double-tap to reset
- **Notification system optimized** - 8-second duration with normal urgency preserves CypherRiot styling
- **Complete state feedback** - Notifications for start, pause, resume, work complete, break complete, and reset
- **Real-time timer integration** - Seamless waybar integration with immediate state updates

#### Window Management Improvements

- **Lollypop music player rules** - Now opens as 50% width/height floating window, centered
- **Consistent floating behavior** - Matches pattern of other media applications like Signal and Feather

#### Application Launcher Cleanup

- **Removed migrate from Fuzzel** - Eliminated sudo integration issues by removing migrate.desktop
- **Migrate still available** - Command remains functional from terminal where sudo works properly
- **Cleaner application menu** - Fuzzel now shows only GUI-appropriate applications

#### Notification System Stabilization

- **Eliminated SwayNC conflicts** - Completely removed SwayNC modules and styling from waybar
- **Mako-only notifications** - Single notification system prevents conflicts and ensures consistency
- **Preserved CypherRiot styling** - Beautiful purple-bordered notifications with proper theming

#### Keybind Optimization

- **Pomodoro timer hotkey** - SUPER + comma for timer control
- **Reorganized notification keys** - Moved dismiss notifications to SUPER + period
- **Logical key grouping** - Related functions grouped for better user experience

#### Technical Improvements

- **Removed redundant waybar modules** - Cleaned up ModulesCustom, ModulesGroups, and CSS
- **Notification system reliability** - Single notification daemon eliminates race conditions
- **Timer state management** - Robust state tracking with proper notification timing
- **Waybar integration** - Seamless timer display with real-time updates

## [2.0.18] - 2025-07-30

### 🏗️ Major Control Panel Architecture Overhaul

#### Control Panel System Integration Revolution

- **Eliminated redundant config system** - Removed conflicting dual source of truth between Control Panel config and actual system configs
- **Direct system management implementation** - All widgets now read/write directly to system services and config files
- **Single source of truth architecture** - Each setting has one authoritative location (system configs, not Control Panel config)
- **Real-time system integration** - Changes apply immediately to actual system state without config conflicts

#### Component-by-Component Conversion

- **BlueLightWidget**: Direct hyprland.conf management - reads/writes actual temperature settings
- **DisplayWidget**: Direct hyprctl integration with smart aspect ratio resolution detection
- **MullvadWidget**: Direct mullvad CLI integration - account numbers and auto-connect from actual VPN service
- **PowerWidget**: Direct powerprofilesctl management - tested balanced ↔ power-saver transitions
- **AudioWidget**: Already using direct pactl system management (verified working)
- **CameraWidget**: Already using direct v4l2/device system management (verified working)
- **PomodoroWidget**: Preserved config system for ArchRiot-specific functionality

#### Smart Resolution Detection

- **Intelligent aspect ratio detection** - Automatically detects display's native aspect ratio
- **Maintains proper ratios** - Only offers resolutions that prevent image distortion
- **Comprehensive ratio support**: 16:9, 16:10, 4:3, 21:9 ultrawide, 32:9 super ultrawide
- **Eliminates nonsensical options** - No more weird resolutions like 1920x1440
- **Filtered by capability** - Only shows resolutions ≤ current display maximum

#### User Experience Improvements

- **No more config conflicts** - Blue light filter showing 3500K in Control Panel while hyprland.conf had 3000K (FIXED)
- **Settings persistence** - Control Panel changes now persist correctly across sessions
- **Immediate application** - Changes apply to system in real-time without restart required
- **Accurate status display** - Control Panel shows actual current system state, not stale config values

#### Technical Architecture

- **Future-proofed design** - Preserved config system for ArchRiot-specific features like Pomodoro timer
- **Reduced complexity** - Eliminated unnecessary abstraction layer that caused conflicts
- **Improved reliability** - System state is always authoritative, preventing drift and inconsistencies
- **Better maintainability** - Single code path for each setting, easier to debug and extend

#### Success Criteria Achieved

- ✅ Control Panel changes persist across sessions/reboots
- ✅ No config conflicts between Control Panel and system
- ✅ User settings preserved during ArchRiot updates
- ✅ Simplified architecture with single source of truth per setting
- ✅ Smart resolution detection prevents display distortion
- ✅ All system integrations work in real-time

## [2.0.17] - 2025-07-30

### 🎯 Fuzzel Application Launcher Improvements

#### Icon Theme Optimization & Desktop Entry Cleanup

- **Maintained Kora icon theme for Fuzzel** - Kora provides superior application icon coverage compared to Tela-purple-dark
- **Added automatic desktop entry cleanup** - Removes unwanted/low-level system utilities from application launcher
- **Removed Hardware Locality (lstopo) entry** - Eliminates confusing "Hardware Locality..." item with missing icon
- **Improved user experience** - Clean application launcher with only useful, properly-themed applications
- **Future-proofed cleanup system** - Expandable array for removing additional unwanted desktop entries

#### Technical Implementation

- **Smart icon theme strategy** - Kora for comprehensive application icons, Tela-purple-dark for GTK/system theming
- **Automated cleanup during installation** - Unwanted desktop entries removed automatically in theming.sh
- **Preserved icon functionality** - All application icons display correctly in Fuzzel launcher
- **State-aware installation** - Installation state management prevents redundant operations

#### User Experience Improvements

- **Clean application launcher** - No more mysterious entries without icons
- **Consistent theming** - Proper icon display across all applications
- **Reduced confusion** - Only user-relevant applications appear in launcher
- **Maintained performance** - Fast icon loading and application discovery

## [2.0.16] - 2025-07-30

### 🎨 Major Theme System Optimization

#### Theme System Redundancy Elimination (ROADMAP Priority #1)

- **Eliminated redundant background service starts** - Removed duplicate `swaybg` startup calls in theming.sh
- **Fixed background cycling for consolidated theme system** - Updated `swaybg-next` to work with flat directory structure
- **Resolved "Cannot find theme directory" validation errors** - Updated `validate.sh` to check consolidated backgrounds directory
- **Fixed hyprlock background loading** - Updated `hyprlock-wrapper.sh` to use consolidated system without symlinks
- **Removed old theme structure dependencies** - All scripts now use `~/.config/archriot/backgrounds/` directly
- **Implemented simple state tracking** - Background cycling now uses `.current-background` state file instead of broken symlinks
- **Verified 23 background files working** - All riot_01.jpg through riot_23.png backgrounds accessible and cycling properly
- **Zero theme validation errors** - Complete theme system now passes all validation checks

#### Technical Improvements

- **Streamlined background service management** - Single background service initialization eliminates restart loops
- **Robust error handling** - Graceful fallbacks to default backgrounds when state files missing
- **Simplified directory structure** - No more complex symlink dependencies or `/current/` directory requirements
- **Performance optimization** - Reduced I/O operations and process spawning during background operations

#### Success Criteria Achieved

- ✅ Zero "Cannot find theme directory" errors during validation
- ✅ Background cycling works perfectly with consolidated system
- ✅ Background service starts exactly once and persists correctly
- ✅ No duplicate theme setup operations
- ✅ Consistent theme state across fresh installs and upgrades

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
