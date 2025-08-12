# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## COMPLETED TASKS

### ✅ **System Upgrade Integration (v2.9.5)**

- Integrated full system upgrade functionality into TUI
- Added proper progress bar integration (90% → 95% → 98% → 100%)
- Implemented default "NO" for upgrade prompt
- Added detailed error logging for AUR and orphan cleanup
- Enhanced Plymouth progress bar timing to show intermediate steps

### ✅ **Blue Light Filter Persistence**

- **Problem**: Control panel blue light changes not surviving system upgrades
- **Root Cause**: `desktop.hyprland` module overwrites `~/.config/hypr/hyprland.conf` during upgrades
- **Solution**: Added `--reapply` flag to control panel that restores user settings from external config
- **Implementation**: Installer calls `archriot-control-panel --reapply` after `desktop.hyprland` module execution
- **Status**: Tested and working - user customizations now persist through upgrades

### ✅ **yay Installation Resilience Enhancement**

- Added retry logic and user choice prompts for AUR installation failures
- Up to 3 retry attempts with user choice between failures
- Option to continue without AUR packages
- Proper cleanup on failed attempts
