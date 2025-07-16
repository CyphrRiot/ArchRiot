#!/bin/bash

# ==============================================================================
# OhmArchy Productivity Applications Setup
# ==============================================================================
# Simple productivity application installation
# ==============================================================================

# Install productivity tools (no heavy office suites)
yay -S --noconfirm --needed \
    gnome-text-editor \
    abiword \
    papers \
    thunar \
    unzip \
    p7zip

# Install calendar and time management
yay -S --noconfirm --needed \
    gnome-clocks

# Configure Gnome Text Editor with Tokyo Night theme
if command -v gnome-text-editor >/dev/null 2>&1; then
    echo "🎨 Configuring Gnome Text Editor..."
    gsettings set org.gnome.TextEditor show-line-numbers true
    gsettings set org.gnome.TextEditor highlight-current-line true
    gsettings set org.gnome.TextEditor show-right-margin false
    gsettings set org.gnome.TextEditor custom-font 'Hack Nerd Font 12'
    gsettings set org.gnome.TextEditor line-height 1.2
    gsettings set org.gnome.TextEditor use-system-font false
    echo "✓ Gnome Text Editor configured"
fi

echo "✅ Productivity applications setup complete!"
