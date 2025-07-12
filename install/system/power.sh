#!/bin/bash

# ==============================================================================
# OhmArchy Power Management Setup
# ==============================================================================
# Simple power management - laptop battery support and CPU scaling
# ==============================================================================

# Install power management tools
if ls /sys/class/power_supply/BAT* &>/dev/null; then
    # Laptop with battery - install power management
    yay -S --noconfirm --needed \
        powertop \
        tlp \
        acpi

    # Enable TLP service
    sudo systemctl enable --now tlp.service
else
    # Desktop system - just basic CPU scaling
    yay -S --noconfirm --needed powertop
fi

echo "✅ Power management setup complete!"
