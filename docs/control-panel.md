# ArchRiot Control Panel Implementation Plan

## 📋 Development Guidelines

### Code Changes

1. **ALWAYS show proposal first**: Reasoning + code change description
2. **ALWAYS make the code change** after approval
3. **ALWAYS ask**: "Can I update your system and run the new code?"
4. **ONE THING AT A TIME**: Single focused modification per interaction
5. **NO assumptions**: Ask questions if anything is unclear

### Git Workflow

1. **NEVER commit** until very end when everything works
2. **Test first**: Verify all functionality before any commits
3. **Explicit file staging**: Never use `git add -A` or `git add .`

### UI Guidelines

1. **Toggle switches ALWAYS rightmost** in title rows
2. **Compact layout**: Minimal vertical spacing
3. **Purple section headers** (#bb9af7) for prominence
4. **Modular widgets**: Each feature as separate class
5. **Real-time updates**: Immediate visual feedback

## ✅ COMPLETED - Phase 1: Core Framework

### Configuration System ✅

- **Location**: `~/.config/archriot/archriot.conf`
- **Parser**: `ArchRiot/bin/archriot-config` (Python CLI tool)
- **Backup System**: Automatic backups with timestamp
- **Import/Export**: JSON export capability

### GUI Framework ✅

- **Base**: GTK 4.0 (first GTK 4 component in ArchRiot)
- **Style**: ArchRiot theme integration with Tokyo Night colors
- **Window
  **: Centered floating dialog (900x700, Hyprland window rules)
- **Modular Architecture**: BaseControlWidget + specific widget classes

### Fully Implemented & System Integrated ✅

1. **🍅 Pomodoro Timer** - Real-time waybar integration + Duration slider (5-60min, snaps to 5min) + "Learn More" dialog + Disabled state support
2. **💡 Blue Light Filter** - Real-time hyprsunset control + Temperature slider (2500K-5000K, snaps to 500K) + Persistent hyprland config updates
3. **🛡️ Mullvad VPN** - Real account detection + Auto-connect toggle + Account privacy (show/hide with eye) + Hyprland GUI integration
4. **🔊 Audio** - Real-time mute/unmute via PulseAudio + Safe system control (no service breaking)
5. **📷 Camera** - Real device access control + Resolution slider + Test Camera with live GTK preview + OpenCV integration
6. **🖥️ Display Settings** - Real monitor resolution detection + Live resolution/scaling changes via hyprctl + Evenly spaced slider
7. **🔋 Power Management** - Real powerprofilesctl integration + Icon slider (🐢🏎️⚖️⚡) + Live profile switching

### Widget Architecture Completed ✅

- **Real-time system integration** - All widgets control live system, not just config
- **Persistent changes** - Survive control panel exit and system reboot
- **Visual feedback** - Immediate updates when changes are made
- **Privacy protection** - Account numbers hidden by default with show/hide
- **Educational content** - "Learn More" dialogs with ArchRiot theming
- **System safety** - Mute instead of killing services, permissions instead of breaking camera

### Temporarily Removed (TODO: Re-implement Later)

**⌨️ Input Devices** - Keyboard/mouse/touchpad settings (commented out)
**🛡️ Security** - Firewall and security settings (commented out)

## 📁 File Structure

### Configuration Files

```
~/.config/archriot/
├── archriot.conf          # Main configuration file
├── control-panel/         # Control panel specific files
└── backups/               # Configuration backups
```

### Implementation Files

```
ArchRiot/
├── bin/
│   ├── archriot-config           # Configuration CLI tool
│   └── archriot-control-panel    # Main GUI application
├── config/
│   ├── archriot/archriot.conf    # Default configuration template
│   └── hypr/hyprland.conf        # Updated with window rules
└── plan.md                       # This file
```

### Window Rules Added

```ini
# ArchRiot Control Panel window rules
windowrulev2 = float, title:^(ArchRiot Control Panel)$
windowrulev2 = center, title:^(ArchRiot Control Panel)$
windowrulev2 = size 900 80%, title:^(ArchRiot Control Panel)$
```

## ✅ Current Status - PRODUCTION READY

### Advanced Features Completed

- ✅ **GTK 4 Application** with ArchRiot theming + dark entry field styling + modular ArchRiotDialog system
- ✅ **Real-time system control** - All widgets modify live running system immediately
- ✅ **"Exit without Saving" vs "Save and Exit"** workflow with original config restoration
- ✅ **Privacy protection** - Account numbers with show/hide eye toggle (👁️/🔒)
- ✅ **Live video preview** - GTK camera test window with OpenCV integration
- ✅ **Educational integration** - "Learn More" dialogs with comprehensive information
- ✅ **Cross-platform power management** - Works with both AMD and Intel systems
- ✅ **Persistent system changes** - Settings survive reboot via config file modifications
- ✅ **Evenly spaced sliders** - Perfect tick mark distribution regardless of underlying values
- ✅ **Icon-based interfaces** - Intuitive emoji sliders for power profiles
- ✅ **Real device detection** - Camera resolution, monitor capabilities, VPN status, audio state
- ✅ **Safe system control** - Mute instead of kill, permissions instead of break

### Advanced Code Architecture

```python
ArchRiotDialog            # Modular reusable dialog system with proper theming ✅
BaseControlWidget         # Reusable base class with change tracking ✅
├── PomodoroWidget       # Waybar integration + Learn More + snapping ✅
├── BlueLightWidget      # Hyprsunset real-time + persistent config ✅
├── MullvadWidget        # Account detection + privacy + GUI integration ✅
├── AudioWidget          # PulseAudio mute control + status detection ✅
├── CameraWidget         # Device permissions + resolution + video test ✅
├── DisplayWidget        # Hyprctl live changes + monitor detection ✅
├── PowerWidget          # Powerprofilesctl + icon slider + profile detection ✅
├── InputWidget          # (commented out - TODO)
└── SecurityWidget       # (commented out - TODO)

ControlPanelWindow       # Exit without saving + restore original config ✅
ControlPanelApplication  # GTK 4 application + global CSS management ✅
```

## ✅ Phase 1 Complete: Core UI + Apply Changes Workflow

### New Apply Changes System

- ✅ **No auto-save**: Changes only applied when "Apply Changes" button is clicked
- ✅ **Visual state tracking**: Apply button disabled when no changes, highlighted when enabled
- ✅ **Unsaved changes protection**: Modal dialog warns before closing with unsaved changes
- ✅ **Change tracking**: Widgets notify window when modifications are made
- ✅ **Batch application**: All widget changes applied together for consistency

### Display Settings Implementation

- ✅ **Real monitor detection**: Uses `hyprctl monitors` to get actual available resolutions
- ✅ **Discrete resolution slider**: Up to 5 valid widths only, no interpolation
- ✅ **Discrete scaling slider**: 100%, 125%, 150%, 175%, 200% only
- ✅ **Snap-to-nearest**: Sliders automatically snap to valid tick marks
- ✅ **Proper defaults**: Scaling defaults to 100% (not 125%)

### Power Management Implementation

- ✅ **Real power profile detection**: Uses `powerprofilesctl get` to detect current system profile
- ✅ **Standard profile options**: Performance, Balanced, Power-saver (filtered from system output)
- ✅ **System integration**: Actually sets power profiles via `powerprofilesctl set <profile>`
- ✅ **Cross-platform support**: Works with both AMD and Intel systems
- ✅ **Proper installation**: Added `power-profiles-daemon` to ArchRiot installer via `power.sh`

### Widget Status

- ✅ **7 widgets implemented**: Pomodoro, Blue Light, Mullvad, Audio, Camera, Display, Power
- ✅ **All core widgets complete**: Full UI and system integration working
- ✅ **Compact layout achieved**: No scrolling needed, proper footer positioning
- ✅ **Visual polish complete**: Consistent spacing, proper margins, clean button layout

## 🚀 NEXT: Phase 2 - Input Devices & Security Widgets

### Remaining Widget Implementation

**⌨️ Input Devices Widget** - Keyboard/mouse/touchpad settings

- Keyboard repeat rate and delay sliders
- Mouse sensitivity and acceleration controls
- Touchpad gesture enable/disable toggles
- Input method switching for international users

**🛡️ Security Widget** - Firewall and security settings

- UFW firewall enable/disable toggle
- Common port toggles (SSH, HTTP, HTTPS)
- Failed login attempt monitoring
- System security status overview

### Advanced Integration Opportunities

**Waybar Status Integration:**

- VPN connection status indicator
- Camera access indicator when active
- Power profile indicator in waybar
- Audio mute status synchronization

**Hyprland Workspace Integration:**

- Per-workspace power profiles
- Display settings per workspace
- Camera permissions per application

**System Monitoring Integration:**

- Real-time resource usage in Power Management
- Network traffic monitoring in VPN widget
- Camera usage tracking and privacy alerts

### Dependency Management Completed

**Added to ArchRiot Installer:**

- `power-profiles-daemon` → `install/system/power.sh`
- `python-opencv` → `install/applications/media.sh`

**Installation Integration:**

```bash
# Power management (follows bluetooth.sh pattern)
install_packages "power-profiles-daemon powertop" "essential"
sudo systemctl enable --now power-profiles-daemon.service

# Camera support
install_packages "python-opencv" "essential"
```

**Future Installer Additions Needed:**

- `v4l-utils` for camera resolution control
- Additional camera testing utilities
- Input device configuration tools

## 🎯 Next Implementation Priorities

## 🚨 CRITICAL CODE ARCHITECTURE ISSUES

### URGENT: Eliminate Code Duplication

**Problem**: Every slider implements the SAME snapping pattern with copy-paste code:

```python
# This exact pattern is duplicated 5+ times across widgets:
def on_value_changed(self, scale):
    value = scale.get_value()
    snapped_value = round(value)  # or custom snapping logic

    if abs(value - snapped_value) > 0.5:
        scale.set_value(snapped_value)
        return

    # Apply change immediately
    index = int(snapped_value)
    if 0 <= index < len(self.available_values):
        # Custom logic here
        pass

    # Mark changes pending
    self.mark_changes_pending()
```

**Solution Needed**: Create reusable slider components:

```python
class SnappingSlider(Gtk.Scale):
    """Reusable slider that snaps to discrete values"""
    def __init__(self, values, current_value, callback):
        super().__init__()
        self.values = values
        self.callback = callback
        self.setup_snapping()

    def setup_snapping(self):
        # Common snapping logic here
        pass

class IconSlider(SnappingSlider):
    """Slider with icon labels"""
    pass

class ResolutionSlider(SnappingSlider):
    """Slider for resolution selection"""
    pass
```

**Impact**: 80% code reduction in slider implementations, consistent behavior, easier maintenance.

### Phase 2A: Complete Widget Set

1. **Input Devices Widget Implementation**
    - Keyboard repeat rate control via `xset r rate`
    - Mouse sensitivity via `xinput` or similar
    - Touchpad gesture configuration
    - Input method switching support

2. **Security Widget Implementation**
    - UFW firewall status and control
    - Common security port management
    - System security overview dashboard
    - Privacy-focused security recommendations

### Phase 2B: Advanced Integration

1. **Waybar Status Indicators**
    - Real-time VPN connection status
    - Camera access privacy indicator
    - Current power profile display
    - Audio system status synchronization

2. **Application Integration**
    - Per-application camera permissions
    - Power profile automation rules
    - Display settings per workspace
    - VPN kill-switch integration

### Phase 2C: User Experience Polish

1. **Keyboard Shortcuts & Access**
    - `SUPER+C` → Control Panel shortcut
    - Fuzzel desktop entry integration
    - Power menu integration (`SUPER+ESCAPE`)
    - Quick settings overlay for common controls

2. **Advanced Configuration**
    - Settings import/export functionality
    - Configuration profiles and presets
    - System health monitoring integration
    - Backup and restore capabilities

## 🔮 Future Expansion Possibilities

### Phase 3: Advanced System Integration

- **Workspace-aware settings**: Different power/display profiles per Hyprland workspace
- **Application-specific controls**: Camera permissions, power profiles per application
- **Automation rules**: Automatic profile switching based on conditions
- **System health monitoring**: Resource usage integration into power management

### Phase 4: Enterprise & Advanced Features

- **Multi-user support**: Per-user configuration profiles
- **Remote management**: Control panel accessible via web interface
- **Backup/restore system**: Full configuration backup with versioning
- **Import/export**: Share configurations between systems
- **Plugin architecture**: Third-party widget development support

### Phase 5: Distribution & Ecosystem

- **Comprehensive testing**: Validation across different hardware configurations
- **Documentation suite**: Video tutorials, user guides, developer documentation
- **Community integration**: User-contributed widgets and themes
- **ArchRiot ecosystem**: Integration with other ArchRiot tools and workflows

---

**Current Status**: ⚠️ **MOSTLY FUNCTIONAL** - 7 widgets implemented but with critical issues
**Major Achievement**: Revolutionary "Exit without Saving" with live preview + restoration
**Technical Excellence**: Modular architecture, privacy protection, educational content, system safety
**Critical Issues**: Power slider not snapping, massive footer spacing problems, code duplication epidemic
**Code Quality**: 🚨 **URGENT REFACTORING NEEDED** - Copy-paste hell across all slider implementations
**UI Polish**: 🚨 **SPACING BROKEN** - 20% empty space in footer area, inconsistent layout
**Next Priority**: Fix fundamental slider snapping and layout issues BEFORE any new features
**Installation Ready**: All dependencies added to ArchRiot installer system

## 🚨 CRITICAL ISSUES BLOCKING PROGRESS

### BROKEN FUNCTIONALITY - MUST FIX FIRST

1. **Power Management Slider** - 🚨 **CRITICAL**: Still not snapping despite 5+ "fixes"
2. **Footer Layout** - 🚨 **CRITICAL**: 20% empty space, poor user experience
3. **Code Quality** - 🚨 **URGENT**: Slider snapping logic copied 7+ times with subtle differences

### ROOT CAUSE ANALYSIS

**Power Slider Issue**: Claiming to use "exact same pattern" but each implementation has different bugs
**Layout Issue**: Constantly adjusting margins/spacing without understanding GTK layout principles
**Architecture Issue**: No reusable components, everything copy-pasted and modified incorrectly

## 🔥 IMMEDIATE ACTION ITEMS

### Priority 1: Fix Broken Core Functionality

1. **Power slider snapping** - Copy EXACT working code from camera widget (proven to work)
2. **Footer spacing** - Remove excessive margins causing 20% empty space
3. **Layout consistency** - Establish standard spacing throughout UI

### Priority 2: Stop Copy-Paste Development

1. **Create SnappingSlider base class** - Single implementation for all sliders
2. **Extract layout helpers** - Standard spacing/margin functions
3. **Widget factory pattern** - Consistent widget creation

### Priority 3: Quality Control

1. **Test all sliders** - Verify snapping works on every widget
2. **Layout verification** - Ensure consistent spacing without empty areas
3. **Code review** - Eliminate duplicate implementations

## 🎯 HONEST ASSESSMENT

**What Works**: 6/7 widgets have functional sliders, real-time system integration, privacy features
**What's Broken**: Power slider, footer layout, massive code duplication
**What's Needed**: Focus on fixing existing issues instead of adding new features
**Time Estimate**: 2-3 focused sessions to resolve core issues before any new development
