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

### Implemented Features ✅

1. **🍅 Pomodoro Timer** - Toggle + Duration slider (5-60min, 5min increments)
2. **💡 Blue Light Filter** - Toggle + Temperature slider (2500K-5000K, 500K increments)
3. **🔒 Mullvad VPN** - Toggle + Account number field (auto-formatted: "1234 5678 9012 3456") + "Get Mullvad VPN" link
4. **🔊 Audio System** - Toggle switch with on/off status
5. **📷 Camera System** - Toggle switch with on/off status

### Stubbed Widgets ✅

6. **🖥️ Display Settings** - Monitor configuration and scaling
7. **⌨️ Input Devices** - Keyboard/mouse/touchpad settings
8. **🔋 Power Management** - Battery and performance profiles
9. **🛡️ Security** - Firewall and security settings

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

## 🎯 Current Status

### Working Features

- ✅ **GTK 4 Application** with ArchRiot theming + dark entry field styling
- ✅ **Modular widget architecture** (BaseControlWidget pattern)
- ✅ **Configuration management** with real-time saving
- ✅ **Debounced slider updates** (500ms delay)
- ✅ **Compact UI layout** (toggle rightmost, labels in title row)
- ✅ **Five functional widgets** with proper validation
- ✅ **Auto-formatted input fields** (Mullvad account spacing)
- ✅ **External link integration** (xdg-open for website links)
- ✅ **Background styling solved** - Consistent solid black backgrounds (no transparency issues)

### Code Architecture

```python
BaseControlWidget          # Reusable base class
├── PomodoroWidget        # Duration slider + toggle ✅
├── BlueLightWidget       # Temperature slider + toggle ✅
├── MullvadWidget         # Account field + toggle + link ✅
├── AudioWidget           # Toggle switch ✅
├── CameraWidget          # Toggle switch ✅
├── DisplayWidget         # (stubbed)
├── InputWidget           # (stubbed)
├── PowerWidget           # (stubbed)
└── SecurityWidget        # (stubbed)

ControlPanelWindow        # Main window assembly
ControlPanelApplication   # GTK 4 application wrapper
```

## ✅ Phase 1 Complete: All UI Controls Working

### Background Styling Solution

- **Issue resolved**: Eliminated opacity/transparency conflicts entirely
- **Solution**: Both window and frame backgrounds use solid black `rgba(0, 0, 0, 1.0)`
- **Consistency**: Welcome script updated to match control panel styling
- **Result**: Perfect visual consistency with no background mismatches

### Modularization Status

- ✅ **BaseControlWidget**: Reusable toggle/slider/entry patterns
- ✅ **Auto-formatting**: Could be extracted for other numeric inputs
- ✅ **External links**: Reusable button pattern for website links
- ✅ **Layout patterns**: Title row with right-aligned toggles standardized

### Potential Reusable Components

- `create_formatted_entry()` - For account numbers, phone numbers, etc.
- `create_external_link_button()` - For website links (implemented in MullvadWidget)
- `create_auto_spacing_field()` - For formatted numeric input (implemented in MullvadWidget)
- `create_simple_toggle()` - For basic on/off controls (implemented in AudioWidget, CameraWidget)

## 🚀 NEXT: Phase 2 - System Integration

### Priority: Connect UI to Live Systems

**CRITICAL**: UI controls must modify actual running systems, not just config files.

### Integration Requirements

#### 1. Pomodoro Timer → Waybar Integration

- **Current**: Updates `archriot.conf` only
- **Needed**: Modify `~/.local/bin/waybar-tomato-timer.py` duration
- **Method**: Update script variables + restart waybar module
- **Files**:
    - `waybar-tomato-timer.py` (duration variable)
    - `waybar` reload command

#### 2. Blue Light Filter → Hyprsunset Integration

- **Current**: Updates `archriot.conf` only
- **Needed**: Control live hyprsunset process
- **Method**: Kill/restart hyprsunset with new temperature
- **Files**:
    - `hyprland.conf` (exec-once line)
    - Live process management (`pkill hyprsunset`, restart)

#### 3. Mullvad VPN → System Integration

- **Current**: Updates `archriot.conf` only
- **Needed**: Configure actual Mullvad client
- **Method**: `mullvad account login` + auto-connect settings
- **Files**:
    - Mullvad CLI integration
    - Waybar VPN status updates

### Integration Implementation Plan

```python
# Add to each widget class:
def apply_to_system(self):
    """Apply configuration changes to running system"""
    pass

# Examples:
class PomodoroWidget:
    def apply_to_system(self):
        # Update waybar-tomato-timer.py duration variable
        # Reload waybar module

class BlueLightWidget:
    def apply_to_system(self):
        # Kill current hyprsunset process
        # Start new process with updated temperature

class MullvadWidget:
    def apply_to_system(self):
        # Run: mullvad account login {account_number}
        # Set: mullvad auto-connect {enabled}
```

## 🎯 Immediate Next Steps

### Phase 2A: Core System Integration

1. **Implement PomodoroWidget.apply_to_system()**
    - Update waybar timer script duration
    - Reload waybar to apply changes

2. **Implement BlueLightWidget.apply_to_system()**
    - Control live hyprsunset process
    - Handle enable/disable + temperature changes

3. **Implement MullvadWidget.apply_to_system()**
    - Integrate with Mullvad CLI
    - Handle account login + auto-connect

### Phase 2B: Boot Integration

1. **Startup configuration loading**
    - Apply settings on Hyprland startup
    - Integrate with existing autostart

2. **Service management**
    - Systemd service integration where needed
    - Proper process lifecycle management

### Phase 2C: Access Points

1. **Keyboard shortcut**: `SUPER+C` → Control Panel
2. **Fuzzel integration**: Desktop entry for app launcher
3. **Power menu integration**: Add to `SUPER+ESCAPE` menu

## 🔮 Future Phases

### Phase 3: Complete Widget Implementation

- Implement all stubbed widgets (Audio, Camera, Display, Input, Power, Security)
- Add widget-specific system integrations

### Phase 4: Advanced Features

- Import/Export configuration presets
- Real-time system state monitoring
- Advanced validation and error handling

### Phase 5: Polish & Distribution

- Comprehensive testing across different systems
- Documentation and user guides
- Integration into main ArchRiot installation

---

**Last Updated**: Phase 1 Complete - All UI Controls Fully Functional
**Status**: Perfect UI with solid black backgrounds, all 5 widgets working
**Background Issue**: RESOLVED - Consistent solid backgrounds eliminate opacity conflicts
**Next**: Commit Phase 1 work, then implement `apply_to_system()` methods for live system control
**Modularization**: Complete - standardized patterns for toggles, sliders, formatted entries, external links
