# ArchRiot Installer Architecture

## 🎯 Purpose

This document exists to prevent the installer from being broken by developers who think they know better. **READ THIS BEFORE TOUCHING INSTALLER CODE.**

The installer has been through hell and back. It works now. Don't break it.

## 📁 File Structure

```
install/
├── helpers/          # Shared utilities (common.sh, packages.sh, ui.sh)
├── modules/          # Installation modules
│   ├──
```

system/ # Base system (first)
│ ├── development/ # Dev tools (second)
│ ├── desktop/ # Desktop environment (third)
│ └── applications/ # User apps (last)
├── install.sh # Main orchestrator
└── setup.sh # Entry point

```

## 🔄 Installation Flow

### Entry Point
- `setup.sh` downloads/updates ArchRiot to `~/.local/share/archriot/`
- Calls `install.sh` to start the real work

### Module Execution
- Dynamic module discovery (maintains flexibility)
- Runs in order: system → development → desktop → applications
- Each module scans for `.sh` files and runs them alphabetically
- **If a .sh file exists in a module directory, IT RUNS** - don't assume it doesn't

## 🚨 Critical Architecture (DON'T BREAK THIS)
```

Rules

### Error Handling Revolution (v1.9.0+)

**CRITICAL CHANGE**: We eliminated the `set -e` nightmare that was causing silent failures.

- ❌ **OLD**: Global `set -e` + output capture = hidden failures
- ✅ **NEW**: Explicit error handling + visible output = reliable installation

**DO NOT BRING BACK `set -e`** - It breaks everything in subtle ways.

### Theme System (SOLVED)

**SOLUTION IMPLEMENTED**: Theme system has been consolidated to eliminate config override issues.

Changes made:

- CypherRiot theme integrated directly into main configs
- All theme directories removed (tokyo-night, cypherriot)
- Config linking logic eliminated from installer
- Single unified theme prevents maintenance issues

### Module Dependencies (DON'T REORDER)

- **Order matters**: system → development → desktop → applications
- Desktop components need system packages first
- Applications need desktop environment
- **Don't change module order** without understanding dependencies

## 🔍 Debugging

```bash
# Check module execution
ls ~/.local/share/archriot/install/modules/

# Verify consolidated backgrounds
ls ~/.config/archriot/backgrounds/

# Check waybar config integrity
grep -i "custom/updates" ~/.config/waybar/config
```

## ⚠️ Development Rules (MANDATORY)

### Before Touching Installer Code

1. **Read this document completely**
2. **Test on fresh VM** (not your customized development system)
3. **Check theme impact** (will this break in themed systems?)
4. **Verify error visibility** (don't hide failures behind pretty progress bars)

### Adding New Features

1. **Add to main config first**
2. **Test with ALL themes** (or you'll create silent failures)
3. **Document theme dependencies**
4. **Consider sudo/privilege requirements**
5. **Test with fresh installation** - existing systems hide dependency issues

### What NOT To Do

- **NEVER use global `set -e`** - It causes more problems than it solves
- **Don't hide command output** - Hidden output = hidden failures
- **Don't assume files don't run** - If a .sh file exists, it executes
- **Don't change module order** - Dependencies flow downward
- **Add features to main configs** - No more theme overrides to worry about
- **Don't break working code** - Understand before modifying

## 🎯 Key Files

- **Main installer**: `install.sh`
- **Theming system**: `modules/desktop/theming.sh`
- **Helper functions**: `helpers/common.sh`

---

---

**Final Warning**: This installer has survived multiple rewrites, silent failure epidemics, and the theme config disaster (now solved). It works now. Respect that.

**Translation**: Don't be the person who breaks the installer because you didn't read this document.
