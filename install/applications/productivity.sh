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

# Install and configure Gnome Text Editor with theme support
if command -v gnome-text-editor >/dev/null 2>&1; then
    echo "🎨 Installing themes for Gnome Text Editor..."

    # Create gtksourceview styles directory
    mkdir -p "$HOME/.local/share/gtksourceview-5/styles"

    # Install all available text editor themes
    local themes_installed=0
    for theme_dir in "$HOME/.local/share/omarchy/themes"/*; do
        if [[ -d "$theme_dir/text-editor" ]]; then
            theme_name=$(basename "$theme_dir")
            for theme_file in "$theme_dir/text-editor"/*.xml; do
                if [[ -f "$theme_file" ]]; then
                    cp "$theme_file" "$HOME/.local/share/gtksourceview-5/styles/"
                    echo "✓ Installed $(basename "$theme_file") theme"
                    ((themes_installed++))
                fi
            done
        fi
    done

    if [[ $themes_installed -gt 0 ]]; then
        echo "✓ $themes_installed text editor theme(s) installed"
    else
        echo "⚠ No text editor themes found"
    fi

    echo "🎨 Configuring Gnome Text Editor..."
    gsettings set org.gnome.TextEditor show-line-numbers true
    gsettings set org.gnome.TextEditor highlight-current-line true
    gsettings set org.gnome.TextEditor show-right-margin false
    gsettings set org.gnome.TextEditor custom-font 'Hack Nerd Font 12'
    gsettings set org.gnome.TextEditor line-height 1.2
    gsettings set org.gnome.TextEditor use-system-font false

    # Set default theme (prefer Tokyo Night, fallback to first available)
    if [[ -f "$HOME/.local/share/gtksourceview-5/styles/tokyo-night.xml" ]]; then
        gsettings set org.gnome.TextEditor style-scheme 'tokyo-night'
        echo "✓ Gnome Text Editor configured with Tokyo Night theme"
    elif [[ -f "$HOME/.local/share/gtksourceview-5/styles/cypherriot.xml" ]]; then
        gsettings set org.gnome.TextEditor style-scheme 'cypherriot'
        echo "✓ Gnome Text Editor configured with CypherRiot theme"
    else
        echo "ℹ Using default text editor theme"
    fi
fi

echo "✅ Productivity applications setup complete!"
