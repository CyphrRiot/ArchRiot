# ArchRiot Version Check System - Implementation Plan

## Overview

Create a periodic version check system that compares the locally installed ArchRiot version with the latest version available in the repository. When a newer version is detected, show a centered popup window with update options.

## Current System Analysis

- **Current Version**: 1.1.59 (stored in `VERSION` file)
- **GUI Framework**: GTK3 (Python) - same as existing welcome window
- **Config System**: JSON files in `~/.config/archriot/`
- **Systemd Integration**: Already has timer-based services (battery-monitor)
- **Existing Infrastructure**: Welcome window provides perfect template

## Implementation Status: ✅ COMPLETE

### Phase 1: Create Version Check Script ✅

**File**: `bin/version-check` ✅ IMPLEMENTED

- ✅ Python script with proper error handling and network timeouts
- ✅ Fetches remote version from `https://cyphrriot.github.io/ArchRiot/VERSION`
- ✅ Semantic version comparison (handles major.minor.patch correctly)
- ✅ JSON configuration system for timing and ignore settings
- ✅ Command line options: `--test`, `--force`, `--reset`
- ✅ Launches update dialog when updates available

### Phase 2: Create Update Notification Window ✅

**File**: `bin/version-update-dialog` ✅ IMPLEMENTED

- ✅ **Window Properties**:
    - Title: "ArchRiot Update Available"
    - Size: 600x400 (smaller than welcome window's 800x600)
    - Centered positioning with proper transparency
    - No image at top (clean layout)
    - Same styling/theme as welcome window with GTK3 CSS

- ✅ **Window Content**:
    - Current version vs. available version display
    - Non-scrolling instruction text with proper word wrapping
    - Clean typography using Hack Nerd Font Mono
    - Three buttons with proper sizing:
        1. **"📥 Install"** (120px)
        2. **"🔕 Ignore Notifications"** (160px)
        3. **"❌ Close"** (100px)

### Phase 3: Button Actions ✅

- ✅ **Install**:
    - Opens terminal (ghostty, kitty, alacritty, or xterm fallback)
    - Executes: `curl https://ArchRiot.org/setup.sh | bash`
    - Proper error handling with dialog fallback

- ✅ **Ignore Future Notifications**:
    - Creates `~/.config/archriot/versions.cfg` with ignore flag
    - Shows confirmation dialog with re-enable instructions
    - Persists setting properly

- ✅ **Close Window**:
    - Simple dialog close with no config changes

### Phase 4: Systemd Timer Service ✅

**Files**: ✅ IMPLEMENTED

- ✅ `config/systemd/user/version-check.service`
- ✅ `config/systemd/user/version-check.timer`

**Timer Configuration**: ✅ WORKING

- ✅ Checks every 4 hours after 5-minute boot delay
- ✅ Proper graphical session dependency
- ✅ Environment variables (DISPLAY, XDG_RUNTIME_DIR) configured
- ✅ Service tested and active in systemd

### Phase 5: Integration Points ✅

- ✅ **Hyprland Config**: Window rules added for version dialog positioning
- ✅ **Install Script**: Integrated into `install/core/03-config.sh`
- ✅ **Service Installation**: Automatic systemd timer enable during setup
- ✅ **Version Config**: Proper config file structure and management

### Phase 6: Configuration File Structure ✅

**File**: `~/.config/archriot/versions.cfg` ✅ IMPLEMENTED

```json
{
    "ignore_notifications": false,
    "last_check": "2024-01-01T00:00:00Z",
    "last_notified_version": "1.1.59",
    "check_interval_hours": 4
}
```

## Technical Details

### Version Comparison Logic

- Parse semantic versions (major.minor.patch)
- Compare numerically, not lexically
- Handle edge cases (missing remote version, network errors)

### Network Check Strategy

- Primary: GitHub API (`https://api.github.com/repos/CyphrRiot/ArchRiot/releases/latest`)
- Fallback: Direct VERSION file fetch from ArchRiot.org
- Timeout and error handling

### Terminal Opening Implementation

- Use Hyprland IPC or keybinding simulation
- Alternative: Direct ghostty execution with specific geometry
- Center terminal window for update process

### Window Rules Addition

Add to `config/hypr/hyprland.conf`:

```
# ArchRiot Version Check window rules
windowrulev2 = float, title:^(ArchRiot Update Available)$
windowrulev2 = center, title:^(ArchRiot Update Available)$
windowrulev2 = size 600 400, title:^(ArchRiot Update Available)$
```

## File Structure

```
ArchRiot/
├── bin/
│   ├── version-check              # Main version check script
│   └── version-update-dialog      # Update notification window
├── config/systemd/user/
│   ├── version-check.service      # Systemd service definition
│   └── version-check.timer        # Systemd timer configuration
└── install/core/03-config.sh      # Updated to install version scripts
```

## Implementation Order ✅ COMPLETED

1. ✅ Create version-check script (background logic)
2. ✅ Create version-update-dialog (GUI window)
3. ✅ Test both scripts manually
4. ✅ Create systemd service files
5. ✅ Update installation scripts
6. ✅ Add Hyprland window rules
7. ✅ Test complete integration
8. ✅ Update VERSION and README
9. ✅ Commit and push changes

## Testing Results ✅

- ✅ Version comparison works correctly (1.1.57 vs 1.1.58 detected properly)
- ✅ Network handling robust with proper timeouts and fallbacks
- ✅ Systemd timer active and scheduled (next run in ~4 hours)
- ✅ All three buttons function correctly:
    - Install button launches terminal with update command
    - Ignore button creates config and shows confirmation
    - Close button exits cleanly
- ✅ Configuration persistence works across restarts
- ✅ Terminal opening supports multiple terminal emulators
- ✅ Dialog displays properly with correct transparency and styling

## Final Implementation Notes ✅

- ✅ Follows existing ArchRiot patterns (GTK3, JSON config, systemd timers)
- ✅ Maintains consistent styling with welcome window
- ✅ Respects user's choice to ignore notifications
- ✅ No disruption to running processes - uses oneshot service type
- ✅ Handles offline scenarios gracefully with proper error messages
- ✅ Integrated into standard ArchRiot installation process
- ✅ Automatic installation and configuration during setup
- ✅ Complete user documentation added to README.md

## Project Status: 🎉 COMPLETE

The ArchRiot automatic version check system is fully implemented and ready for production use. All phases completed successfully with comprehensive testing and integration.
