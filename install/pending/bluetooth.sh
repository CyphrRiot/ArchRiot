#!/bin/bash

# ==============================================================================
# ArchRiot Bluetooth System Setup
# ==============================================================================
# Simple bluetooth setup - install blueberry and enable service
# ==============================================================================

# Install bluetooth controls
install_packages "blueberry" "essential"

# Turn on bluetooth by default
sudo systemctl enable --now bluetooth.service

echo "✅ Bluetooth setup complete!"
