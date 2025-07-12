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

echo "✅ Development editors setup complete!"
