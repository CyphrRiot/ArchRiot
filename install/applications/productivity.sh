#!/bin/bash

# ==============================================================================
# OhmArchy Productivity Applications Setup
# ==============================================================================
# Simple productivity application installation
# ==============================================================================

# Install productivity tools (no heavy office suites)
yay -S --noconfirm --needed \
    apostrophe \
    papers \
    thunar \
    unzip \
    p7zip \
    ark

# Install calendar and time management
yay -S --noconfirm --needed \
    gnome-clocks

echo "✅ Productivity applications setup complete!"
