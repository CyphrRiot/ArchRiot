#!/bin/bash

# ==============================================================================
# OhmArchy Utilities Setup
# ==============================================================================
# Simple utility applications installation
# ==============================================================================

# Install system utilities
yay -S --noconfirm --needed \
    htop \
    btop \
    neofetch \
    tree \
    wget \
    curl \
    unzip \
    p7zip \
    ark

# Install file management utilities
yay -S --noconfirm --needed \
    thunar \
    thunar-volman \
    thunar-archive-plugin \
    gvfs \
    gvfs-mtp

# Install system monitoring
yay -S --noconfirm --needed \
    gnome-system-monitor \
    baobab \
    gnome-disk-utility

# Install network utilities
yay -S --noconfirm --needed \
    nm-connection-editor \
    networkmanager-applet

# Install calculator and simple tools
yay -S --noconfirm --needed \
    gnome-calculator \
    gnome-font-viewer \
    file-roller

echo "✅ Utilities setup complete!"
