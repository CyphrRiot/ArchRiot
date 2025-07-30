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
install_packages "python-opencv" "essential"
install_packages "lollypop" "optional"
# Install spotdl with --nocheck flag to bypass failing API tests
print_status "INFO" "Installing spotdl (requires --nocheck due to API test failures)"
install_aur_nocheck "spotdl" "optional"

echo "✅ Media applications setup complete!"
