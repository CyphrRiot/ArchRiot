#!/bin/bash

# ==============================================================================
# OhmArchy Communication Applications Setup
# ==============================================================================
# Simple communication application installation
# ==============================================================================

# Install browsers (no corporate messaging apps)
yay -S --noconfirm --needed \
    brave-bin \
    firefox \
    chromium

# Terminal already installed in core, but ensure kitty available
yay -S --noconfirm --needed \
    kitty

echo "✅ Communication applications setup complete!"
