# Changelog

All notable changes to OhmArchy will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.8] - 2025-07-16

### 🔧 Critical Bug Fixes

- **Signal Keybinding**: Fixed SUPER+G not launching Signal on fresh installs
    - **ROOT CAUSE**: signal-desktop package was never being installed during setup
    - **SOLUTION**: Added signal-desktop to communication.sh yay installation list
    - **IMPACT**: SUPER+G keybinding now works immediately after installation

---

## [1.0.2] - 2025-01-15

### 🔧 Bug Fixes

- **Waybar Configuration**: Fixed duplicate JSON key in 4workspaces.json
    - Removed duplicate `"title<.*reddit.*>"` entry causing validation warnings
    - All diagnostics now pass clean without warnings

### ✨ Enhancements

- **Base System**: Added `bc` (basic calculator) to core installation
    - Fixes script failures when mathematical calculations are needed
    - Resolves `sync-media` and other script compatibility issues on minimal installations

### 🔍 Quality Improvements

- **Configuration Validation**: Enhanced JSON structure integrity
- **Code Quality**: Eliminated duplicate entries in waybar window-rewrite rules

---

## [1.0.1] - 2025-01-14

### 🔧 Minor Fixes

- **Initramfs errors**: Fixed in v1.0.1 - update if experiencing build issues

---

## [1.0.0] - 2025-01-14

### 🚀 Major Changes

#### Terminal Replacement: Ghostty → Kitty

- **BREAKING**: Replaced Kitty terminal with Ghostty as the default terminal emulator
- Added Ghostty shell integration for enhanced functionality
- Implemented proper theme system integration with `config-file` imports
- Added floating terminal support with `SUPER+SHIFT+RETURN` keybinding

#### Fastfetch Display Fix

- **FIXED**: Fastfetch content truncation in Ghostty terminal windows
- Added 0.1s timing delay in `fish_greeting` to allow terminal sizing
- Implemented `--logo-width 20` parameter for consistent layout
- Ensures beautiful fastfetch display on all terminal launch methods

### ✨ New Features

- **Versioning System**: Added VERSION file and version display in installation
- **Version Command**: New `ohmarchy-version` command to check installed version and system status
- **Enhanced Validation**: Updated pre-install validation to check Ghostty components
- **Floating Terminal**: Centered floating terminal window (1200x800) with proper opacity
- **Theme Conversion**: CypherRiot and Tokyo Night themes converted to Ghostty format

### 🔧 Fixed Issues

#### Critical Installation Breaks

- **Desktop Applications**: Fixed `.desktop` files to use `ghostty` instead of `kitty`
    - `About.desktop`, `Activity.desktop`, `nvim.desktop`
- **Management Scripts**: Updated bin scripts to reference Ghostty configs
    - `ohmarchy-theme-next`, `ohmarchy-validate-system`
- **Waybar Integration**: Fixed right-click terminal commands in Waybar configs
- **Theme System**: Fixed main Ghostty config to import themes instead of hardcoded colors

#### Configuration Updates

- **Hyprland**: Updated window rules for Ghostty class names (`com.mitchellh.ghostty`)
- **Fish Shell**: Added fastfetch timing fix with sleep delay in greeting function
- **Theme Files**: Converted theme files from Kitty syntax to Ghostty syntax
- **Package Installation**: Updated core installation to include `ghostty` and `ghostty-shell-integration`

### 🗑️ Removed

- **Legacy Configs**: Removed `/config/kitty/` directory and configurations
- **Outdated References**: Eliminated all Kitty references from documentation and configs

### 📚 Documentation

- **README**: Updated to reflect Ghostty as default terminal
- **Keybindings**: Added documentation for floating terminal keybinding
- **Theme Docs**: Updated CypherRiot theme documentation
- **Installation Guide**: Enhanced with version information display

### 🔍 Validation & Quality

- **Pre-Install Checks**: Enhanced validation script to verify Ghostty components
- **Component Testing**: Added checks for shell integration and theme files
- **Installation Verification**: Improved post-install validation for Ghostty functionality
- **Error Prevention**: Fixed multiple breaking issues that would prevent smooth installation

### 🎨 Visual Improvements

- **Terminal Theming**: Proper theme integration with dynamic color switching
- **Window Management**: Enhanced floating window rules and positioning
- **Opacity Settings**: Consistent transparency across terminal instances
- **Font Configuration**: Maintained CaskaydiaMono Nerd Font with proper sizing

### ⚡ Performance

- **Startup Speed**: Optimized terminal launch with proper initialization timing
- **Theme Switching**: Improved theme reload functionality for Ghostty
- **Resource Usage**: Maintained lightweight footprint while adding features

---

## Technical Details

### Package Changes

- **Added**: `ghostty`, `ghostty-shell-integration`
- **Removed**: `kitty` (replaced by Ghostty)

### Configuration Structure

```
~/.config/ghostty/
├── config                    # Main Ghostty configuration
└── current-theme.conf        # Symlinked theme file

~/.config/omarchy/current/
└── theme/
    └── ghostty.conf          # Theme-specific colors
```

### New Commands

- `ohmarchy-version` - Display version and system information
- `SUPER+SHIFT+RETURN` - Launch floating centered terminal

### Breaking Changes

- All terminal-related keybindings now use Ghostty instead of Kitty
- Theme files use Ghostty syntax (palette = N=#rrggbb)
- Desktop application launchers updated to use Ghostty

---

## Migration Notes

Users upgrading from a Kitty-based installation should:

1. Run the installation normally - Ghostty will be installed automatically
2. Themes will be converted to Ghostty format during installation
3. All keybindings remain the same (`SUPER+RETURN` for terminal)
4. Floating terminal is a new feature accessible via `SUPER+SHIFT+RETURN`

For a complete fresh installation, run:

```bash
curl -fsSL https://ohmarchy.org/setup.sh | bash
```

---

## Verification

After installation, verify everything is working:

```bash
ohmarchy-version                   # Check installation status
fastfetch                         # Verify display works correctly
ghostty --version                 # Confirm Ghostty is installed
```
