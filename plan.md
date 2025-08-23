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
- Core ArchRiot application now fully themed

## 🚧 CURRENT TASKS

### TASK 5: System-Wide Dynamic Theming Extension

**PRIORITY: HIGH**

**Current Status:**

- ✅ **Waybar** - Complete integration, real-time updates via SIGUSR2
- ✅ **Zed Editor** - Full theme override system with fallback
- ✅ **Ghostty Terminal** - Color palette updates (static file updates only)
- ✅ **Hyprland Window Manager** - Border color coordination
- ✅ **Text Editor** - XML theme generation with gsettings switching
- ✅ **Architecture** - Complete modular `ThemeRegistry` system

**Clean Modular Architecture (v2.21):**

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
    └── texteditor.go      // Text Editor XML themes + gsettings
```

## ⏭️ IMMEDIATE NEXT STEPS

### 🔄 Phase 1: Remaining ArchRiot Applications

**NEXT TARGET: Fuzzel Application Launcher**

- **fuzzel** (`config/fuzzel/`) - Application launcher colors
- Simple configuration, likely easy implementation
- Part of core ArchRiot user experience

**Research needed:**

- How fuzzel themes are applied in `packages.yaml`
- Configuration file format and color properties
- Real-time reload capabilities

### 🔄 Phase 2: System Enhancement

**KNOWN ISSUE: Ghostty Real-time Reload**

- Current: Static file updates (users must manually reload)
- Research status: SIGUSR2 not implemented in Ghostty 1.0/1.1
- Alternative approaches being investigated

**Additional Applications:**

- **btop** (`config/btop/`) - System monitor with color themes
- **mako** (`config/mako/`) - Notification daemon styling
- **fastfetch** (`config/fastfetch/`) - System info display colors

### 🔄 Phase 3: Advanced Integration

- **fish** (`config/fish/`) - Shell prompt and syntax highlighting
- **gtk-3.0/gtk-4.0** (`config/gtk-*/`) - System-wide application theming
- **Thunar** (`config/Thunar/`) - File manager theming
- **Neovim** - Dynamic colorscheme generation

## 🚧 OTHER TASKS

### TASK 6: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

- **Problem:** Current implementation disabled due to boot failures
- **Solution:** Replace manual OpenSSL approach with reliable `sbctl`
- **Status:** Deferred until theming system complete

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming

**Core System (COMPLETED):**

- [✅] Modular architecture with ThemeRegistry
- [✅] Waybar real-time dynamic theming
- [✅] Zed Editor theme override system
- [✅] Ghostty terminal color integration
- [✅] Hyprland window manager border coordination
- [✅] Text Editor XML theme generation and switching
- [✅] All components respect dynamic/static toggle

**Next Targets:**

- [ ] **Fuzzel launcher** - Application launcher theming
- [ ] **Ghostty real-time reload** - Config refresh without restart
- [ ] **btop system monitor** - Color theme coordination
- [ ] **mako notification daemon** - Notification styling

**Long-term Goals:**

- [ ] **System-wide consistency** - All ArchRiot applications themed
- [ ] **GTK integration** - System-wide application theming

### Key Principles

- **Matugen does color science** - appliers just update config files
- **Single responsibility** - each applier handles one application
- **Graceful failure** - individual applier errors don't break system
- **Research-driven** - understand existing mechanisms before implementing

## 🔧 TECHNICAL DEBT

- Control Panel JSON boolean serialization issues
- Installation process optimization

---

**Next Action:** Research and implement fuzzel application launcher theming
