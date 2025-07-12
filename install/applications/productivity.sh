#!/bin/bash

# ==============================================================================
# OhmArchy Productivity Applications Setup
# ==============================================================================
# Simple productivity application installation
# ==============================================================================

# Install office suite
yay -S --noconfirm --needed \
    libreoffice-fresh \
    hunspell \
    hunspell-en_us

# Install PDF tools
yay -S --noconfirm --needed \
    okular \
    evince

# Install text editors and note-taking
yay -S --noconfirm --needed \
    obsidian \
    typora \
    apostrophe

# Install calendar and time management
yay -S --noconfirm --needed \
    kalendar \
    gnome-clocks

# Install file management
yay -S --noconfirm --needed \
    ark \
    unzip \
    p7zip

echo "✅ Productivity applications setup complete!"
