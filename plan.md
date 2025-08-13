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

## 🚧 OUTSTANDING TASKS

### TASK 1: Kernel Upgrade Reboot Detection

### Problem

When system upgrade upgrades the Linux kernel, the system should default to "YES" for reboot prompt with clear messaging about why reboot is needed.

### Requirements

- Detect when kernel package is upgraded during system upgrade process
- Default reboot prompt to "YES" instead of "NO" when kernel is upgraded
- Show clear message in TUI scrolling window: "Linux Kernel upgraded, you really should reboot"
- Maintain existing reboot prompt behavior for non-kernel upgrades

### Implementation

- Add kernel upgrade detection during pacman upgrade process
- Modify reboot prompt to change default based on kernel upgrade status
- Add informative logging about kernel upgrade requirement

**Priority**: Medium - Improves user experience and system stability

### TASK 2: System Upgrade Display Bug

### Problem

During system upgrade, mysterious "Pack..." text appears in the TUI scrolling window without context or explanation.

### Example Output

```
│ ✅ 💫 Orphan Cleanup       Orphaned packages removed successfully                  │
│ 📋 💫 Orphan Output        Removal output: checking dependencies...                │
│                                                                                    │
│ Pack...                                                                            │
│ ✅ 💫 Package Upgrade      Complete system upgrade finished successfully           │
```

### Investigation Required

- Identify source of "Pack..." output in system upgrade process
- Determine if it's truncated package manager output
- Fix or remove mysterious display text
- Ensure all upgrade output is properly formatted and informative

**Priority**: Low - Cosmetic issue but affects user experience

## 📋 NEXT PRIORITIES

1. **Kernel Upgrade Reboot Detection** - Improve reboot prompting for kernel upgrades
2. **System Upgrade Display Bug** - Fix mysterious "Pack..." output

## VERSION HISTORY

- **v2.9.6**: Blue light persistence, Plymouth progress fixes, control panel --reapply
- **v2.9.5**: System upgrade integration, AUR resilience enhancements
- **v2.9.4**: Core installer stability improvements
