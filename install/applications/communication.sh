#!/bin/bash

# ==============================================================================
# ArchRiot Communication Applications Setup
# ==============================================================================
# Simple communication application installation
# ==============================================================================

# Install browsers and communication apps
yay -S --noconfirm --needed \
    brave-bin \
    signal-desktop



# Source web2app function for creating web application launchers
source ~/.local/share/archriot/default/bash/functions

# Install essential web applications
echo "🌐 Installing essential web applications..."

# Copy desktop files for web apps
local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "$script_dir/../../applications/Proton Mail.desktop" ]]; then
    cp "$script_dir/../../applications/Proton Mail.desktop" ~/.local/share/applications/
    echo "✓ Proton Mail desktop file installed"
fi

if [[ -f "$script_dir/../../applications/Google Messages.desktop" ]]; then
    cp "$script_dir/../../applications/Google Messages.desktop" ~/.local/share/applications/
    echo "✓ Google Messages desktop file installed"
fi

if [[ -f "$script_dir/../../applications/X.desktop" ]]; then
    cp "$script_dir/../../applications/X.desktop" ~/.local/share/applications/
    echo "✓ X (Twitter) desktop file installed"
fi

# Copy icons
if [[ -d "$script_dir/../../applications/icons" ]]; then
    mkdir -p ~/.local/share/icons
    cp -r "$script_dir/../../applications/icons"/* ~/.local/share/icons/ 2>/dev/null || true
    echo "✓ Web app icons installed"
fi

echo "✅ Communication applications and web apps setup complete!"
