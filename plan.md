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

## 🚧 CURRENT TASKS

### TASK 4: System-Wide Dynamic Theming Extension

**PRIORITY: HIGH**

**Current Status:**

- ✅ **Waybar** - Complete integration, real-time updates
- ✅ **Zed Editor** - Full theme override system with fallback
- ✅ **Ghostty Terminal** - Color palette updates (no real-time reload)
- ✅ **Hyprland Window Manager** - Border color coordination
- ✅ **Architecture** - `ThemeApplier` interface and `ThemeRegistry` system
- ✅ **v2.20 Production Fixes** - Template paths, config saving, complete color definitions

**Current Modular Architecture:**

```
source/theming/
├── theming.go          // Core matugen integration + config
├── interfaces.go       // ThemeApplier interface + ThemeRegistry
└── [monolithic code]   // Needs refactoring into appliers/
```

**Current Clean Architecture:**

```
source/theming/
├── theming.go              // Core orchestration + matugen
├── interfaces.go           // Legacy (mostly unused)
├── registry.go             // Coordinates all applications
└── applications/
    ├── types.go           // Shared MatugenColors + ThemeApplier interface
    ├── waybar.go          // Waybar colors.css + SIGUSR2 reload (~130 lines)
    ├── zed.go             // Zed theme_overrides (~170 lines)
    ├── ghostty.go         // Ghostty terminal palette (~160 lines)
    └── hyprland.go        // Hyprland window borders (~100 lines)
```

## ⏭️ IMMEDIATE NEXT STEPS

### 🔄 Phase 1: Missing ArchRiot Applications

**CRITICAL - Our Own Applications Missing from Theming:**

1. **Text Editor** (`config/text-editor/`) - Our default markdown editor
    - Has `cypherriot.xml` and `tokyo-night.xml` theme files
    - Should generate dynamic XML themes from matugen colors
    - Core ArchRiot application that users see daily

2. **Real-time Reload Issues:**
    - Ghostty needs config reload signal (like waybar's SIGUSR2)
    - Users should see theme changes immediately without restarting terminals

### 🔄 Phase 2: Config Directory Audit

**Potential Theming Candidates from `config/` directories:**

- **btop** (`config/btop/`) - System monitor with color themes
- **mako** (`config/mako/`) - Notification daemon styling
- **fuzzel** (`config/fuzzel/`) - Application launcher colors
- **fastfetch** (`config/fastfetch/`) - System info display colors
- **fish** (`config/fish/`) - Shell prompt and syntax highlighting
- **gtk-3.0/gtk-4.0** (`config/gtk-*/`) - System-wide application theming
- **Thunar** (`config/Thunar/`) - File manager theming

### 🔄 Phase 3: Extended System Integration

4. **Neovim** - Dynamic colorscheme generation
5. **System-wide GTK theming** - Consistent application colors

## 🚧 OTHER TASKS

### TASK 5: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

- **Problem:** Current implementation disabled due to boot failures
- **Solution:** Replace manual OpenSSL approach with reliable `sbctl`
- **Status:** Deferred until theming system complete

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming

- [✅] Core modular architecture with ThemeRegistry
- [✅] Waybar real-time dynamic theming with complete color definitions
- [✅] Zed Editor theme override system with proper syntax colors
- [✅] Ghostty terminal color integration (static file updates)
- [✅] Hyprland window manager border color coordination
- [✅] Clean applier refactor complete (all applications modularized)
- [✅] All components respect dynamic/static toggle
- [✅] Critical v2.20 production fixes (template paths, config saving)
- [ ] **Text Editor** - ArchRiot's markdown editor XML theme generation
- [ ] **Ghostty real-time reload** - Config refresh without terminal restart
- [ ] **btop system monitor** - Color theme coordination
- [ ] **mako notification daemon** - Notification styling
- [ ] **System-wide consistency** - All ArchRiot applications themed

### Key Principles

- **Matugen does color science** - appliers just update config files
- **Single responsibility** - each applier handles one application
- **Graceful failure** - individual applier errors don't break system
- **~20-30 lines per applier** - focused, maintainable code

## 🔧 TECHNICAL DEBT

- Control Panel JSON boolean serialization issues
- CLI `--toggle-dynamic-theming` save functionality
- Installation process optimization

### 🎯 PRIORITY ORDER

1. **URGENT:** Text Editor applier - Core ArchRiot application missing theming
2. **HIGH:** Ghostty real-time reload - Improve user experience
3. **MEDIUM:** btop, mako, fuzzel - Complete ArchRiot ecosystem
4. **LOW:** GTK system-wide, Neovim - External application support

---

**Next Action:** Add text-editor applier for ArchRiot's markdown editor
