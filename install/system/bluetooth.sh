#!/bin/bash

# ==============================================================================
# OhmArchy Bluetooth System Setup
# ==============================================================================
# Simple bluetooth setup - install blueberry and enable service
# ==============================================================================

# Install bluetooth controls
yay -S --noconfirm --needed blueberry

# Turn on bluetooth by default
sudo systemctl enable --now bluetooth.service

echo "✅ Bluetooth setup complete!"
