# Zed Editor Wayland Integration

## Overview

ArchRiot includes native Wayland support for the Zed editor to ensure optimal rendering and performance on Wayland-based desktop environments like Hyprland.

## What's Included

### 1. Wayland Launcher Script (`zed-wayland`)

Located at `bin/zed-wayland`, this script ensures Zed runs with proper Wayland backend:

- Forces Wayland display server usage
- Disables X11 fallback for native rendering
- Sets proper scaling and DPI settings
- Optimizes for Wayland-specific features

### 2. Desktop Integration

The `applications/zed.desktop` file provides:

- Proper MIME type associations for code files
- Native Wayland execution through `zed-wayland` wrapper
- Application launcher integration
- File manager context menu support

### 3. Automatic Installation

The development editors installer (`install/development/editors.sh`) automatically:

- Installs the official Arch `zed` package
- Copies the Wayland launcher to `~/.local/bin/`
- Installs the desktop file to `~/.local/share/applications/`
- Updates the desktop database for immediate availability

## Manual Installation

If you need to install the Wayland integration manually:

```bash
# Copy the launcher script
cp ~/.local/share/omarchy/bin/zed-wayland ~/.local/bin/
chmod +x ~/.local/bin/zed-wayland

# Install desktop file
cp ~/.local/share/omarchy/applications/zed.desktop ~/.local/share/applications/
update-desktop-database ~/.local/share/applications/
```

## Usage

### Command Line

```bash
# Launch with Wayland support
zed-wayland

# Open specific files
zed-wayland file1.rs file2.py

# Create new window
zed-wayland --new
```

### Desktop Environment

- Use the application launcher (Super+A in Hyprland)
- Search for "Zed"
- Right-click files in file manager and "Open with Zed"

## Environment Variables Set

The launcher script sets these Wayland-specific variables:

- `WAYLAND_DISPLAY=wayland-1`
- `MOZ_ENABLE_WAYLAND=1`
- `GDK_BACKEND=wayland`
- `QT_QPA_PLATFORM=wayland`
- `SDL_VIDEODRIVER=wayland`
- `_JAVA_AWT_WM_NONREPARENTING=1`

## Benefits

- **Native Wayland rendering**: No XWayland compatibility layer
- **Better performance**: Direct GPU access and compositing
- **Proper scaling**: Correct HiDPI handling
- **Theme integration**: Native GTK/Qt theming support
- **Security**: Wayland's improved security model

## Troubleshooting

### If Zed appears "grey" or broken:

1. Verify you're using the correct package:

    ```bash
    pacman -Q zed
    which zeditor
    ```

2. Test the Wayland launcher:

    ```bash
    zed-wayland --help
    ```

3. Check environment variables:
    ```bash
    echo $WAYLAND_DISPLAY
    echo $GDK_BACKEND
    ```

### Fallback to X11:

If needed, you can still run Zed through XWayland:

```bash
GDK_BACKEND=x11 zeditor
```

## Package Information

- **Package**: `zed` (official Arch repository)
- **Binary**: `/usr/bin/zeditor` (CLI launcher)
- **Main executable**: `/usr/lib/zed/zed-editor`
- **Version**: Updates with Arch package system
