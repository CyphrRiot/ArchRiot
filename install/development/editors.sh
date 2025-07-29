#!/bin/bash

# ==============================================================================
# ArchRiot Development Editors Setup
# ==============================================================================
# Simple editor installation - Neovim with LazyVim and alternatives
# ==============================================================================

# Install core editors and tools
yay -S --noconfirm --needed \
    neovim \
    tree-sitter-cli \
    ripgrep \
    fd \
    fzf \
    zed

# Install common LSP servers
yay -S --noconfirm --needed \
    lua-language-server \
    pyright \
    typescript-language-server \
    bash-language-server

# Note: Neovim configuration with TokyoNight theme is installed by the main config installer
# This preserves the ArchRiot dark theme setup and prevents overwriting user configs
echo "✓ Neovim will be configured with ArchRiot TokyoNight theme via main config installer"

# Install Zed desktop file and Wayland launcher
echo "🎯 Installing Zed desktop integration..."
mkdir -p ~/.local/share/applications ~/.local/bin

# Install Wayland launcher script
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ -f "$script_dir/../../bin/zed-wayland" ]]; then
    cp "$script_dir/../../bin/zed-wayland" ~/.local/bin/
    chmod +x ~/.local/bin/zed-wayland
    echo "✓ Zed Wayland launcher installed"
else
    echo "⚠ Zed Wayland launcher not found in repository bin"
fi

if [[ -f "$script_dir/../../applications/zed.desktop" ]]; then
    # Hide system zed desktop files to prevent duplicates
    mkdir -p ~/.local/share/applications/

    # Hide dev.zed.Zed.desktop (actual system file)
    if [[ -f "/usr/share/applications/dev.zed.Zed.desktop" ]]; then
        echo "[Desktop Entry]
NoDisplay=true" > ~/.local/share/applications/dev.zed.Zed.desktop
        echo "✓ System dev.zed.Zed.desktop hidden"
    fi

    # Hide zed.desktop (legacy system file)
    if [[ -f "/usr/share/applications/zed.desktop" ]]; then
        echo "[Desktop Entry]
NoDisplay=true" > ~/.local/share/applications/zed.desktop.bak
        echo "✓ System zed.desktop hidden"
    fi

    # Install our Wayland-fixed version with proper HOME expansion
    sed "s|\$HOME|$HOME|g" "$script_dir/../../applications/zed.desktop" > ~/.local/share/applications/zed.desktop
    echo "✓ Zed desktop file installed with Wayland support"
else
    echo "⚠ Zed desktop file not found in repository applications"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Development editors setup complete!"
