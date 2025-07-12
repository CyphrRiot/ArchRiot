#!/bin/bash

# ==============================================================================
# OhmArchy Communication Applications Setup
# ==============================================================================
# Simple communication application installation
# ==============================================================================

# Install messaging applications
yay -S --noconfirm --needed \
    signal-desktop \
    discord \
    element-desktop \
    telegram-desktop

# Install video conferencing
yay -S --noconfirm --needed \
    zoom \
    teams

# Install email clients
yay -S --noconfirm --needed \
    thunderbird

# Install browsers for web-based communication
yay -S --noconfirm --needed \
    firefox \
    chromium

echo "✅ Communication applications setup complete!"
