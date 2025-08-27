# ArchRiot Development Plan

## 🚧 CURRENT TASKS

### TASK 1: Hyprlock Screen Dynamic Theming - DEFERRED

**PRIORITY: LOW**

**Current Status:**

- ❌ **Complex implementation challenges** - Multiple failed approaches
- ⚠️ **Template preservation issues** - Difficult to maintain all UI elements
- 🔄 **Needs simpler approach** - Current pattern doesn't fit hyprlock's complex config

**Lessons Learned:**

- Hyprlock has complex nested config with many label sections for UI elements
- Template-based replacement is error-prone due to formatting complexity
- String replacement approaches break the rich UI (time, date, system stats)
- Config generation approach loses important template elements
- Unlike other apps, hyprlock requires preserving extensive formatting and comments

**Proposed Solution for Future:**

- Study simpler approaches like CSS variable injection (if hyprlock supports it)
- Consider leaving hyprlock with static theming only (it's a lock screen, not primary UI)
- Focus resources on more impactful theming targets

## ⏭️ IMMEDIATE NEXT STEPS

### 🔄 Phase 1: Complete Hyprlock Integration

1. **Focus on remaining high-impact theming targets**
2. **Consider hyprlock static-only acceptable for lock screen**
3. **Research if hyprlock adds theme variable support in future versions**
4. **Prioritize applications users interact with more frequently**

### 🔄 Phase 2: Remaining ArchRiot Applications

**KNOWN ISSUE: Ghostty Real-time Reload**

- Current: Static file updates (users must manually reload with Ctrl+Shift+,)
- Research status: SIGUSR2 not implemented in Ghostty 1.0/1.1
- Priority: LOW (existing functionality works, just not real-time)

**Potential candidates for theming integration:**

- **btop** (`config/btop/`) - System monitor with color themes
- **mako** (`config/mako/`) - Notification daemon styling
- **fastfetch** (`config/fastfetch/`) - System info display colors
- **fish** (`config/fish/`) - Shell prompt and syntax highlighting

### 🔄 Phase 3: System-wide Integration

- **gtk-3.0/gtk-4.0** (`config/gtk-*/`) - System-wide application theming
- **Thunar** (`config/Thunar/`) - File manager theming
- **Neovim** - Dynamic colorscheme generation

## 🚧 OTHER TASKS

### TASK 2: Secure Boot Implementation Overhaul

**PRIORITY: MEDIUM**

- **Problem:** Current implementation disabled due to boot failures
- **Solution:** Replace manual OpenSSL approach with reliable `sbctl`
- **Status:** Deferred until theming system complete

## 🎯 SUCCESS CRITERIA

### System-Wide Dynamic Theming

**Core System (COMPLETED v2.25):**

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
- [✅] Fast, reliable setup and installation
- [🔄] **Hyprlock screen theming** - Research phase, needs simpler approach
- [❌] **Hyprlock screen theming** - Deferred due to implementation complexity

**Current Architecture (v2.25):**

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
- **One change at a time** - wait for feedback after each modification

## 🔧 TECHNICAL

DEBT

- Installation process optimization - **RESOLVED v2.25**
- Background/theme reversion after reboot - **RESOLVED v2.24**
- Control Panel state management issues - **RESOLVED v2.23**

## 🏆 ACHIEVEMENTS

**Dynamic Theming System v2.25:**

- **6 applications** integrated with dynamic theming (all tested and working)
- **100% persistence** - all settings survive reboot with bulletproof config reading
- **Unified control** - CLI commands and Control Panel interface working perfectly
- **Professional architecture** - clean, maintainable, extensible applier system
- **Proven stability** - handles errors gracefully, no breaking changes
- **Lightning-fast setup** - 26MB download vs 297MB, seconds vs minutes
- **Production ready** - reliable background persistence and theme management

---

**Current Status:** Core dynamic theming system is complete and production-ready. All critical ArchRiot applications support dynamic theming with proper persistence, fast installation, and reliable user controls.

**Next Action:** Research and design a simpler approach for hyprlock theming based on lessons learned from failed complex implementation.
