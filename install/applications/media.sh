#!/bin/bash

# ==============================================================================
# OhmArchy Media Applications Setup
# ==============================================================================
# Simple media application installation
# ==============================================================================

# Install audio applications
yay -S --noconfirm --needed \
    audacious \
    audacity \
    pavucontrol \
    pulseeffects

# Install video applications
yay -S --noconfirm --needed \
    vlc \
    mpv \
    obs-studio \
    ffmpeg

# Install graphics applications
yay -S --noconfirm --needed \
    gimp \
    inkscape \
    krita \
    blender

# Install media downloaders
yay -S --noconfirm --needed \
    yt-dlp \
    gallery-dl

echo "✅ Media applications setup complete!"
