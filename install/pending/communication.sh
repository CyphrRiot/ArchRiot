#!/bin/bash

# ==============================================================================
# ArchRiot Communication Applications Setup
# ==============================================================================
# Simple communication application installation
# ==============================================================================

# Install browsers and communication apps
install_packages "brave-bin" "essential"
install_packages "signal-desktop" "essential"



# Source web2app function for creating web application launchers
source ~/.local/share/archriot/default/bash/functions

# Install essential web applications
echo "🌐 Installing essential web applications..."

# Copy desktop files for web apps
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f ~/.local/share/archriot/applications/Proton\ Mail.desktop ]]; then
    cp ~/.local/share/archriot/applications/Proton\ Mail.desktop ~/.local/share/applications/
    echo "✓ Proton Mail desktop file installed"
fi

if [[ -f ~/.local/share/archriot/applications/Google\ Messages.desktop ]]; then
    cp ~/.local/share/archriot/applications/Google\ Messages.desktop ~/.local/share/applications/
    echo "✓ Google Messages desktop file installed"
fi

if [[ -f ~/.local/share/archriot/applications/X.desktop ]]; then
    cp ~/.local/share/archriot/applications/X.desktop ~/.local/share/applications/
    echo "✓ X (Twitter) desktop file installed"
fi

# Copy icons
if [[ -d ~/.local/share/archriot/applications/icons ]]; then
    mkdir -p ~/.local/share/icons
    cp -r ~/.local/share/archriot/applications/icons/* ~/.local/share/icons/ 2>/dev/null || true
    echo "✓ Web app icons installed"
fi

echo "✅ Communication applications and web apps setup complete!"
