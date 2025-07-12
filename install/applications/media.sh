#!/bin/bash

# ==============================================================================
# OhmArchy Media Applications Setup
# ==============================================================================
# Simple media application installation
# ==============================================================================

# Install essential media tools
yay -S --noconfirm --needed \
    audacious \
    totem \
    mpv \
    pavucontrol \
    ffmpeg \
    yt-dlp

echo "✅ Media applications setup complete!"
