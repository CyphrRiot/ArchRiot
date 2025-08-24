# ArchRiot Development Plan

## 🎉 COMPLETED TASKS

### ✅ TASK 1: Dynamic Color Theming System - COMPLETED (v2.18)

- Matugen color extraction from wallpapers
- Complete waybar theming with real-time reload (SIGUSR2)
- Control panel toggle for enable/disable dynamic theming
- CLI commands: `--apply-wallpaper-theme` and `--toggle-dynamic-theming`
- Fallback to CypherRiot static colors when disabled

### ✅ TASK 2: Terminal Emoji Fallback System - COMPLETED

- Automatic terminal emoji capability detection
- ASCII alternatives with exact visual spacing

### ✅ TASK 3: Kernel Upgrade Reboot Detection - COMPLETED

- Fixed reboot handling and continuation flow

### ✅ TASK 4: Text Editor Dynamic Theming - COMPLETED (v2.21)

- Complete text editor applier for GNOME Text Editor
- Dynamic XML theme generation from matugen colors
- Proper gtksourceview-5/styles integration with gsettings theme switching
- Dynamic theme: `cypherriot-dynamic.xml` ↔ Static theme: `cypherriot.xml`
- Improved contrast mapping with proper dark backgrounds
- Core ArchRiot application now fully themed

### ✅ TASK 5: Background/Theme Persistence - COMPLETED (v2.22)

- Fixed critical issue where settings reverted after reboot
- Smart startup-background script replaces hardcoded hyprland.conf
- Reads saved configuration and applies correct background/theme on startup
- All dynamic theming settings now persist across system restarts

### ✅ TASK 6: Fuzzel Launcher Dynamic Theming - COMPLETED (v2.22)

- Complete fuzzel applier for application launcher
- RGBA hex color format support (8-character hex)
- Dynamic color mapping with proper dark background contrast
- INI file parsing that preserves non-color settings
- Real-time color updates when switching themes

### ✅ TASK 7: Control Panel Dynamic Theming Fixes - COMPLETED (v2.23)

- Fixed UI state persistence - Control Panel now shows saved dynamic setting
- Fixed immediate color application when enabling dynamic theming
- Resolved config conflicts between Control Panel and ArchRiot binary
- Proper GTK Switch callback handling with `switch.get_active()`
- Control Panel no longer overwrites properly managed theme settings

## 🚧 CURRENT TASKS

### TASK 8: System-Wide Dynamic Theming Extension

**PRIORITY: MEDIUM**

**Current Status:**

- ✅ **Waybar** - Complete integration, real-time updates via SIGUSR2
- ✅ **Zed Editor** - Full theme override system with fallback
- ✅ **Ghostty Terminal** - Color palette updates (static file updates only)
- ✅ **Hyprland Window Manager** - Border color coordination
- ✅ **Text Editor** - XML theme generation with gsettings switching
- ✅ **Fuzzel Launcher** - Application launcher dynamic colors
- ✅ **Architecture** - Complete modular `ThemeRegistry` system
- ✅ **Persistence** - All settings survive reboot via startup-background script

**Current Clean Architecture (v2.23):**

```
source/theming/
├── theming.go              // Core orchestration + matugen integration
├── registry.go             // Coordinates all theme appliers
└── applications/
    ├── types.go           // Shared MatugenColors + ThemeApplier interface
    ├── waybar.go          // Waybar colors.css + SIGUSR2 reload
    ├── zed.go             // Zed theme_overrides
    ├── ghostty.go         // Ghostty terminal palette
    ├── hyprland.go        // Hyprland window borders
    ├── texteditor.go      // Text Editor XML themes + gsettings
    └── fuzzel.go          // Fuzzel launcher RGBA colors
```

## ⏭️ IMMEDIATE NEXT STEPS

### 🔄 Phase 1: Core ArchRiot Applications (Remaining)

**KNOWN ISSUE: Ghostty Real-time Reload**

- Current: Static file updates (users must manually reload with Ctrl+Shift+,)
- Research status: SIGUSR2 not implemented in Ghostty 1.0/1.1
- Alternative approaches being investigated
- Priority: LOW (existing functionality works, just not real-time)

### 🔄 Phase 2: Additional ArchRiot Applications

**Potential candidates for theming integration:**

- **btop** (`config/btop/`) - System monitor with color themes
- **mako** (`config/mako/`) - Notification daemon styling
- **fastfetch** (`config/fastfetch/`) - System info display colors
- **fish** (`config/fish/`) - Shell prompt and syntax highlighting

**Research needed for each:**

- How they're configured in `packages.yaml`
- Configuration file format and color properties
- Real-time reload capabilities
- Integration complexity vs. user value

### 🔄 Phase 3: System-wide Integration

- **gtk-3.0/gtk-4.0** (`config/gtk-*/`) - System-wide application theming
- **Thunar** (`config/Thunar/`) - File manager theming
- **Neovim** - Dynamic colorscheme generation

## 🚧 OTHER TASKS

### TASK 9: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

- **Problem:** Current implementation disabled due to boot failures
- **Solution:** Replace manual OpenSSL approach with reliable `sbctl`
- **Status:** Deferred until theming system complete

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming

**Core System (COMPLETED v2.23):**

- [✅] Modular architecture with ThemeRegistry
- [✅] Waybar real-time dynamic theming
- [✅] Zed Editor theme override system
- [✅] Ghostty terminal color integration
- [✅] Hyprland window manager border coordination
- [✅] Text Editor XML theme generation and switching
- [✅] Fuzzel launcher dynamic color integration
- [✅] All components respect dynamic/static toggle
- [✅] Settings persistence across reboots
- [✅] Control Panel full functionality

**Optional Extensions:**

- [ ] **Ghostty real-time reload** - Config refresh without restart
- [ ] **btop system monitor** - Color theme coordination
- [ ] **mako notification daemon** - Notification styling
- [ ] **System-wide GTK theming** - Application-wide consistency

### Key Principles

- **Matugen does color science** - appliers just update config files
- **Single responsibility** - each applier handles one application
- **Graceful failure** - individual applier errors don't break system
- **Research-driven** - understand existing mechanisms before implementing
- **Proper contrast** - maintain readability with dark backgrounds and dynamic accents

## 🔧 TECHNICAL DEBT

- Control Panel JSON boolean serialization issues - **RESOLVED v2.23**
- Installation process optimization
- Background/theme reversion after reboot - **RESOLVED v2.22**

## 🏆 ACHIEVEMENTS

**Dynamic Theming System v2.23:**

- **6 applications** fully integrated with dynamic theming
- **100% persistence** - all settings survive reboot
- **Unified control** - CLI commands and Control Panel interface
- **Professional architecture** - clean, maintainable, extensible
- **Proven stability** - handles errors gracefully, no breaking changes

---

**Current Status:** Core dynamic theming system is complete and production-ready. All critical ArchRiot applications support dynamic theming with proper persistence and user controls.

**Next Focus:** Optional extensions for additional applications based on user feedback and value assessment.
