#!/bin/bash

# ==============================================================================
# ArchRiot Audio System Setup
# ==============================================================================
# Simple audio setup - PipeWire with essential utilities
# ==============================================================================

# Install PipeWire audio system
install_packages "pipewire pipewire-alsa pipewire-pulse" "essential"
install_packages "wireplumber" "essential"
install_packages "pavucontrol pamixer playerctl" "essential"

# Enable PipeWire services
systemctl --user enable --now pipewire.service pipewire-pulse.service wireplumber.service

echo "✅ Audio setup complete!"
