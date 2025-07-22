# ArchRiot Installer Architecture Documentation

## 🎯 Purpose
This document exists to prevent the installer from being broken by developers who don't understand how it works. **READ THIS BEFORE MODIFYING INSTALLER CODE.**

## 📁 File Structure

```
install/
├── core/                   # Core system setup (run in order)
│   ├── 01-base.sh         # Base packages, AUR helper
│   ├── 02-identity.sh     # User identity setup
│   ├── 03-config.sh       # Config installation (AFTER desktop)
│   └── 04-shell.sh        # Shell configuration
├── desktop/               # Desktop environment (ALL .sh files run)
│   ├── hyprland.sh        # Hyprland WM + utilities
│   ├── theming.sh         # CRITICAL: Theme symlinks & setup
│   ├── fonts.sh           # Font installation
│   └── apps.sh            # Desktop applications
├── system/                # System services
├── development/           # Dev tools
├── applications/          # User applications
└── optional/              # Optional components
```

## 🔄 Installation Flow

### 1. Entry Point
- `setup.sh` downloads/updates ArchRiot to `~/.local/share/archriot/`
- Calls `install.sh` to start installation

### 2. Module Execution Order (install.sh)
```bash
declare -a install_modules=(
    "core/02-identity.sh"   # User identity setup
    "desktop"               # Desktop environment (ALL .sh files)
    "core/03-config.sh"     # Config installation (AFTER desktop)
    "core/04-shell.sh"      # Shell configuration
    "system"                # System services
    "development"           # Development tools
    "applications"          # User applications
    "optional"              # Optional components
)
```

### 3. Desktop Module Execution
The "desktop" module runs **ALL .sh files** in `install/desktop/` in alphabetical order:
1. `apps.sh` - Desktop applications
2. `fonts.sh` - Font installation
3. `hyprland.sh` - Hyprland WM + utilities
4. `theming.sh` - **CRITICAL** theme setup

## 🎨 Theme System (theming.sh)

### Critical Functions
- `setup_archriot_theme_system()` - Creates `~/.config/archriot/` structure
- `set_default_theme()` - Creates `~/.config/archriot/current/theme` symlink
- `setup_theme_backgrounds()` - Creates background symlinks
- `link_theme_configs()` - Links theme configs to applications

### What Gets Created
```
~/.config/archriot/
├── current/
│   ├── theme -> ~/.local/share/archriot/themes/cypherriot
│   ├── background -> theme/backgrounds/riot_01.jpg
│   └── backgrounds -> theme/backgrounds/
├── themes/           # Symlinks to available themes
└── backgrounds/      # Background organization
```

### Applications That Depend On This
- **hyprlock** - Sources `~/.config/archriot/current/theme/hyprlock.conf`
- **waybar** - May use theme-specific configs
- **mako** - Uses `~/.config/archriot/current/theme/mako.ini`
- **Background systems** - Use `~/.config/archriot/current/background`

## ⚠️ Critical Dependencies

### Installation Order Matters
1. **Desktop components MUST install first** (hyprland, apps)
2. **theming.sh MUST run** to create theme symlinks
3. **config installation happens AFTER** desktop setup
4. **Theme configs depend on symlinks existing**

### Symlink Dependencies
- `hyprlock.conf` sources `~/.config/archriot/current/theme/hyprlock.conf`
- If `~/.config/archriot/current/` doesn't exist → **BROKEN LOCK SCREEN**
- If theme symlinks are broken → **BLACK SCREEN, NO UI**

## 🚨 What NOT To Do

### DON'T Remove theming.sh
- **NEVER** remove or skip `install/desktop/theming.sh`
- It creates critical symlinks that apps depend on
- Without it: broken lock screen, broken themes, broken backgrounds

### DON'T Change Module Order
- Desktop MUST run before config installation
- theming.sh MUST run before configs that depend on themes
- Changing order breaks dependencies

### DON'T Assume Files Don't Run
- If a .sh file exists in a module directory, **IT RUNS**
- Don't assume something is broken just because you can't see it
- **READ THE ACTUAL CODE** before making changes

## 🔍 Debugging Installation Issues

### Check If Modules Ran
```bash
# Check if theme system was set up
ls -la ~/.config/archriot/current/

# Check if symlinks exist
readlink ~/.config/archriot/current/theme
readlink ~/.config/archriot/current/background

# Check if hyprlock config sources theme
cat ~/.config/hypr/hyprlock.conf
```

### Common Failure Points
1. **Theme symlinks missing** → `setup_archriot_theme_system()` didn't run
2. **Hyprlock black screen** → Theme config sourcing fails
3. **Background not working** → Background symlinks broken
4. **Apps can't find theme** → Config linking failed

## 📝 Making Changes Safely

### Before Modifying Installer
1. **Read this document completely**
2. **Understand the dependency chain**
3. **Test on a fresh VM/system**
4. **Don't assume anything**

### Testing Checklist
- [ ] Fresh installation completes without errors
- [ ] `~/.config/archriot/current/` directory exists
- [ ] Theme and background symlinks work
- [ ] Hyprlock shows UI (not black screen)
- [ ] All dependent applications work

### When Adding New Components
- Consider where in the installation order it belongs
- Identify what it depends on (themes, configs, etc.)
- Test with fresh installation, not existing systems
- Document new dependencies here

## 🎯 Key Principles

1. **Dependencies flow downward** - later modules depend on earlier ones
2. **Theme system is critical** - many components depend on it
3. **Test with fresh installs** - existing systems hide dependency issues
4. **Document changes** - update this file when adding dependencies
5. **Don't break working code** - understand before modifying

## 📚 References

- Main installer: `install.sh`
- Theme system: `install/desktop/theming.sh`
- Config installation: `install/core/03-config.sh`
- Entry point: `setup.sh`

---

**Remember: The installer works. Don't break it by making assumptions.**
