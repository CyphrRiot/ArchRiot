#!/bin/bash

# ==============================================================================
# ArchRiot Hardware Setup
# ==============================================================================
# Simple hardware driver installation
# ==============================================================================

# Install hardware detection utilities
install_packages "usbutils pciutils" "essential"

# Enable multilib for 32-bit support
sudo sed -i '/^#\[multilib\]/,/^#Include/ s/^#//' /etc/pacman.conf
# Note: No manual pacman -Sy needed - yay will sync when installing packages

# Install GPU drivers based on hardware
if lspci | grep -qi nvidia; then
    install_packages "nvidia-dkms nvidia-utils lib32-nvidia-utils" "essential"
fi

if lspci | grep -qi -E 'amd|radeon'; then
    install_packages "mesa lib32-mesa vulkan-radeon lib32-vulkan-radeon" "essential"
fi

if lspci | grep -i intel | grep -qi -E 'vga|3d|display|graphics'; then
    echo "🎮 Installing Intel graphics drivers..."
    install_packages "mesa lib32-mesa vulkan-intel lib32-vulkan-intel intel-media-driver intel-gmmlib libva-intel-driver" "essential"
    echo "✓ Intel graphics drivers installed"
fi

# Fix Apple keyboard if detected
if command -v lsusb >/dev/null 2>&1; then
    if lsusb | grep -i apple | grep -qi keyboard; then
        echo "🍎 Apple keyboard detected - configuring function keys..."
        echo "options hid_apple fnmode=2" | sudo tee /etc/modprobe.d/hid_apple.conf >/dev/null
        echo "✓ Apple keyboard configured (function keys work normally)"
    fi
else
    echo "⚠ lsusb not available - skipping Apple keyboard detection"
fi

echo "✅ Hardware setup complete!"
