#!/bin/bash

# ==============================================================================
# OhmArchy Utilities Setup
# ==============================================================================
# Simple utility applications installation
# ==============================================================================

# Install system utilities
yay -S --noconfirm --needed \
    btop \
    neofetch \
    fastfetch \
    tree \
    wget \
    curl \
    unzip \
    p7zip

# Install file management utilities
yay -S --noconfirm --needed \
    thunar \
    thunar-volman \
    thunar-archive-plugin \
    gvfs \
    gvfs-mtp

# Install system monitoring
yay -S --noconfirm --needed \
    gnome-system-monitor

# Install network utilities
yay -S --noconfirm --needed \
    nm-connection-editor \
    networkmanager-applet

# Install essential tools
yay -S --noconfirm --needed \
    gnome-calculator \
    file-roller

echo "✅ Utilities setup complete!"
