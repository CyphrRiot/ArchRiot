#!/bin/bash

# ==============================================================================
# ArchRiot Power Management Setup
# ==============================================================================
# Simple power management - install power-profiles-daemon and monitoring tools
# Modern Linux power management is sufficient for most users
# ==============================================================================

# Install power management controls
install_packages "power-profiles-daemon powertop" "essential"

# Install battery monitoring for laptops
if ls /sys/class/power_supply/BAT* &>/dev/null; then
    install_packages "acpi" "optional"
fi

# Enable power-profiles-daemon service
sudo systemctl enable --now power-profiles-daemon.service

echo "✅ Power management setup complete!"
echo "ℹ️  Modern Linux power management is used (no TLP conflicts)"
