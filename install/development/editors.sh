#!/bin/bash

# ==============================================================================
# OhmArchy Development Editors Setup
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

# Setup LazyVim if no nvim config exists
if [[ ! -d ~/.config/nvim ]]; then
    git clone https://github.com/LazyVim/starter ~/.config/nvim
    rm -rf ~/.config/nvim/.git
fi

# Install Zed desktop file and Wayland launcher
echo "🎯 Installing Zed desktop integration..."
mkdir -p ~/.local/share/applications ~/.local/bin

# Install Wayland launcher script
if [[ -f "$HOME/.local/share/omarchy/bin/zed-wayland" ]]; then
    cp "$HOME/.local/share/omarchy/bin/zed-wayland" ~/.local/bin/
    chmod +x ~/.local/bin/zed-wayland
    echo "✓ Zed Wayland launcher installed"
else
    echo "⚠ Zed Wayland launcher not found in OhmArchy bin"
fi

# Replace system desktop file with OhmArchy version
if [[ -f "$HOME/.local/share/omarchy/applications/zed.desktop" ]]; then
    sudo cp "$HOME/.local/share/omarchy/applications/zed.desktop" /usr/share/applications/dev.zed.Zed.desktop
    echo "✓ System Zed desktop file replaced with Wayland support"
else
    echo "⚠ Zed desktop file not found in OhmArchy applications"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Development editors setup complete!"
