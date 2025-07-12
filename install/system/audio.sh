#!/bin/bash

# ==============================================================================
# OhmArchy Audio System Setup
# ==============================================================================
# Simple audio setup - PipeWire with essential utilities
# ==============================================================================

# Install PipeWire audio system
yay -S --noconfirm --needed \
    pipewire \
    pipewire-alsa \
    pipewire-pulse \
    wireplumber \
    pavucontrol \
    pamixer \
    playerctl

# Enable PipeWire services
systemctl --user enable --now pipewire.service pipewire-pulse.service wireplumber.service

echo "✅ Audio setup complete!"
