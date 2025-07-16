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

# Setup LazyVim configuration
setup_lazyvim() {
    local nvim_config="$HOME/.config/nvim"

    echo "🚀 Setting up LazyVim..."

    if [[ -d "$nvim_config" ]]; then
        # Backup existing config
        local backup_dir="$nvim_config.backup-$(date +%s)"
        echo "📦 Backing up existing nvim config to: $backup_dir"
        mv "$nvim_config" "$backup_dir"
    fi

    # Clone LazyVim starter
    if git clone https://github.com/LazyVim/starter "$nvim_config"; then
        rm -rf "$nvim_config/.git"
        echo "✓ LazyVim starter configuration installed"
        echo "💡 First nvim launch will install plugins automatically"
    else
        echo "❌ Failed to clone LazyVim starter"
        return 1
    fi
}

# Install LazyVim
setup_lazyvim

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
