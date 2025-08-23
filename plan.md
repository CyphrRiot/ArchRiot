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
- ✅ **Architecture** - `ThemeApplier` interface and `ThemeRegistry` system

**Current Modular Architecture:**

```
source/theming/
├── theming.go          // Core matugen integration + config
├── interfaces.go       // ThemeApplier interface + ThemeRegistry
└── [monolithic code]   // Needs refactoring into appliers/
```

**Target Clean Architecture:**

```
source/theming/
├── theming.go          // Core orchestration
├── interfaces.go       // Simplified ColorApplier interface
├── registry.go         // ApplierRegistry coordination
└── appliers/          // Focused single-purpose appliers (~20-30 lines each)
    ├── waybar.go      // Move existing waybar logic here
    ├── zed.go         // Move existing Zed logic here
    ├── ghostty.go     // NEW: Terminal color palette
    └── hyprland.go    // NEW: Window border colors
```

## ⏭️ IMMEDIATE NEXT STEPS

### 🔄 Phase 1: Clean Modular Refactor

1. **Extract existing logic** - Move waybar/Zed code to focused appliers
2. **Simplify interfaces** - Leverage matugen's format capabilities
3. **Implement registry** - Clean batch coordination

### 🔄 Phase 2: New Application Support

1. **Ghostty Terminal** - Color palette from matugen output
2. **Hyprland** - Window border color string replacement
3. **btop** - System monitor theming
4. **Neovim** - Dynamic colorscheme generation

## 🚧 OTHER TASKS

### TASK 5: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

- **Problem:** Current implementation disabled due to boot failures
- **Solution:** Replace manual OpenSSL approach with reliable `sbctl`
- **Status:** Deferred until theming system complete

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming

- [✅] Core modular architecture with ThemeRegistry
- [✅] Waybar real-time dynamic theming
- [✅] Zed Editor theme override system
- [ ] Clean applier refactor (waybar.go, zed.go extracted)
- [ ] Ghostty terminal color integration
- [ ] Hyprland window manager coordination
- [ ] btop system monitor theming
- [ ] Neovim editor colorscheme support
- [ ] All components respect dynamic/static toggle

### Key Principles

- **Matugen does color science** - appliers just update config files
- **Single responsibility** - each applier handles one application
- **Graceful failure** - individual applier errors don't break system
- **~20-30 lines per applier** - focused, maintainable code

## 🔧 TECHNICAL DEBT

- Control Panel JSON boolean serialization issues
- CLI `--toggle-dynamic-theming` save functionality
- Installation process optimization

---

**Next Action:** Refactor existing theming code into clean modular appliers
