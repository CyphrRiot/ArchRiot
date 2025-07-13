#!/bin/bash

# ==============================================================================
# OhmArchy Communication Applications Setup
# ==============================================================================
# Simple communication application installation
# ==============================================================================

# Install browsers (no corporate messaging apps)
yay -S --noconfirm --needed \
    brave-bin

# Terminal already installed in core, but ensure kitty available
yay -S --noconfirm --needed \
    kitty

# Source web2app function for creating web application launchers
source ~/.local/share/omarchy/default/bash/functions

# Install essential web applications
echo "🌐 Installing essential web applications..."
web2app "Proton Mail" https://mail.proton.me/u/11/inbox https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/protonmail.png
web2app "Google Messages" https://messages.google.com/web/conversations https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/google-messages.png
web2app "X" https://x.com/ https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/x-light.png

echo "✅ Communication applications and web apps setup complete!"
