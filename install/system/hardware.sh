#!/bin/bash

# ==============================================================================
# ArchRiot Hardware Setup
# ==============================================================================
# Simple hardware driver installation
# ==============================================================================

# Enable multilib for 32-bit support
sudo sed -i '/^#\[multilib\]/,/^#Include/ s/^#//' /etc/pacman.conf
sudo pacman -Sy

# Install GPU drivers based on hardware
if lspci | grep -qi nvidia; then
    yay -S --noconfirm --needed nvidia-dkms nvidia-utils lib32-nvidia-utils
fi

if lspci | grep -qi -E 'amd|radeon'; then
    yay -S --noconfirm --needed mesa lib32-mesa vulkan-radeon lib32-vulkan-radeon
fi

if lspci | grep -qi intel | grep -qi -E 'vga|3d|display|graphics'; then
    echo "🎮 Installing Intel graphics drivers..."
    yay -S --noconfirm --needed mesa lib32-mesa vulkan-intel lib32-vulkan-intel intel-media-driver intel-gmmlib libva-intel-driver
    echo "✓ Intel graphics drivers installed"
fi

# Fix Apple keyboard if detected
if lsusb | grep -qi apple | grep -qi keyboard; then
    echo "options hid_apple fnmode=2" | sudo tee /etc/modprobe.d/hid_apple.conf
fi

echo "✅ Hardware setup complete!"
