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

# Install and configure Gnome Text Editor with Tokyo Night theme
if command -v gnome-text-editor >/dev/null 2>&1; then
    echo "🎨 Installing Tokyo Night theme for Gnome Text Editor..."

    # Create gtksourceview styles directory
    mkdir -p "$HOME/.local/share/gtksourceview-5/styles"

    # Install Tokyo Night theme
    if [[ -f "$HOME/.local/share/omarchy/themes/tokyo-night/text-editor/tokyo-night.xml" ]]; then
        cp "$HOME/.local/share/omarchy/themes/tokyo-night/text-editor/tokyo-night.xml" "$HOME/.local/share/gtksourceview-5/styles/"
        echo "✓ Tokyo Night theme installed"
    fi

    echo "🎨 Configuring Gnome Text Editor..."
    gsettings set org.gnome.TextEditor show-line-numbers true
    gsettings set org.gnome.TextEditor highlight-current-line true
    gsettings set org.gnome.TextEditor show-right-margin false
    gsettings set org.gnome.TextEditor custom-font 'Hack Nerd Font 12'
    gsettings set org.gnome.TextEditor line-height 1.2
    gsettings set org.gnome.TextEditor use-system-font false
    gsettings set org.gnome.TextEditor style-scheme 'tokyo-night'
    echo "✓ Gnome Text Editor configured with Tokyo Night theme"
fi

echo "✅ Productivity applications setup complete!"
