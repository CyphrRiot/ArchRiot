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
    zed \
    micro

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

# Install desktop file
if [[ -f "$HOME/.local/share/omarchy/applications/zed.desktop" ]]; then
    cp "$HOME/.local/share/omarchy/applications/zed.desktop" ~/.local/share/applications/
    echo "✓ Zed desktop file installed with Wayland support"
else
    echo "⚠ Zed desktop file not found in OhmArchy applications"
fi

# Override system Zed desktop file to prevent duplicates
if [[ -f /usr/share/applications/dev.zed.Zed.desktop ]]; then
    echo -e "[Desktop Entry]\nHidden=true" > ~/.local/share/applications/dev.zed.Zed.desktop
    echo "✓ System Zed desktop file overridden (hidden from launcher)"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Development editors setup complete!"
