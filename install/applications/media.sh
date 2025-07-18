#!/bin/bash

# ==============================================================================
# ArchRiot Media Applications Setup
# ==============================================================================
# Simple media application installation
# ==============================================================================

# Install essential media tools
yay -S --noconfirm --needed \
    lollypop \
    totem \
    mpv \
    pavucontrol \
    ffmpeg \
    yt-dlp \
    spotdl

echo "✅ Media applications setup complete!"
