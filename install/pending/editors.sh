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
    fzf

# Install common LSP servers
yay -S --noconfirm --needed \
    lua-language-server \
    pyright \
    typescript-language-server \
    bash-language-server

# Note: Neovim configuration with TokyoNight theme is installed by the main config installer
# This preserves the ArchRiot dark theme setup and prevents overwriting user configs
echo "✓ Neovim will be configured with ArchRiot TokyoNight theme via main config installer"

# Create vi symlink to nvim for system tools (visudo, etc.)
echo "🔗 Creating vi -> nvim symlink for system compatibility..."
if sudo ln -sf /usr/bin/nvim /usr/bin/vi; then
    echo "✓ vi symlink created (system tools will use nvim)"
else
    echo "⚠ Failed to create vi symlink"
fi

# Note: Zed editor is now installed and configured via install/applications/productivity.sh
# This includes intelligent Vulkan driver detection, Wayland support, and desktop integration
echo "💡 Zed editor installation handled by productivity applications module"

echo "✅ Development editors setup complete!"
