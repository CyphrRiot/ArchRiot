#!/bin/bash

# ==============================================================================
# ArchRiot Media Applications Setup
# ==============================================================================
# Simple media application installation
# ==============================================================================

# Install essential media tools
install_packages "mpv pavucontrol ffmpeg" "essential"
install_packages "totem" "essential"
install_packages "yt-dlp" "essential"
install_packages "lollypop" "optional"
install_packages "spotdl" "optional"

echo "✅ Media applications setup complete!"
