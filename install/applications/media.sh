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
echo "ℹ Installing spotdl (requires --nocheck due to API test failures)"
if yay -S --noconfirm --needed --mflags "--nocheck" spotdl; then
    echo "✓ spotdl installed successfully"
else
    echo "⚠ spotdl installation failed - this may reduce functionality but installation will continue"
fi

echo "✅ Media applications setup complete!"
