#!/bin/bash

# ==============================================================================
# ArchRiot Productivity Applications Setup
# ==============================================================================
# Simple productivity application installation
# ==============================================================================

# Install productivity tools (no heavy office suites)
install_packages "gnome-text-editor" "essential"
install_packages "zed" "essential"
install_packages "abiword papers" "essential"
install_packages "thunar unzip p7zip" "essential"

# Install calendar and time management
install_packages "gnome-clocks" "essential"

# Install and configure Gnome Text Editor with theme support
if command -v gnome-text-editor >/dev/null 2>&1; then
    echo "🎨 Installing themes for Gnome Text Editor..."

    # Create gtksourceview styles directory
    mkdir -p "$HOME/.local/share/gtksourceview-5/styles"

    # Install all available text editor themes
    local themes_installed=0
    for theme_dir in "$HOME/.local/share/archriot/themes"/*; do
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

# Install Zed editor configuration and launcher
echo "🖥️ Installing Zed editor configuration..."
local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install Zed Wayland launcher
if [[ -f "$script_dir/../../bin/zed-wayland" ]]; then
  cp "$script_dir/../../bin/zed-wayland" ~/.local/bin/
  chmod +x ~/.local/bin/zed-wayland
  echo "✓ Zed Wayland launcher installed"
else
  echo "⚠ Zed Wayland launcher not found in repository"
fi

# Install Zed desktop file
if [[ -f "$script_dir/../../applications/zed.desktop" ]]; then
  cp "$script_dir/../../applications/zed.desktop" ~/.local/share/applications/
  echo "✓ Zed desktop file installed"
else
  echo "⚠ Zed desktop file not found in repository"
fi

# Install Zed configuration
if [[ -d "$script_dir/../../config/zed" ]]; then
  mkdir -p ~/.config/zed
  cp -r "$script_dir/../../config/zed"/* ~/.config/zed/
  echo "✓ Zed configuration installed"
else
  echo "⚠ Zed configuration not found in repository"
fi

echo "✅ Productivity applications setup complete!"
