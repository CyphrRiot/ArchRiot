#!/bin/bash

# ==============================================================================
# OhmArchy Power Management Setup
# ==============================================================================
# Simple power management - basic monitoring tools only
# Modern Linux power management is sufficient for most users
# ==============================================================================

# Install basic power monitoring tools
if ls /sys/class/power_supply/BAT* &>/dev/null; then
    # Laptop with battery - install monitoring tools
    echo "🔋 Detected laptop - installing power monitoring tools..."
    yay -S --noconfirm --needed \
        powertop \
        acpi
    echo "✓ Power monitoring tools installed"
else
    # Desktop system - just powertop for monitoring
    echo "🖥️  Desktop detected - installing power monitoring..."
    yay -S --noconfirm --needed powertop
    echo "✓ Power monitoring tools installed"
fi

echo "✅ Power management setup complete!"
echo "ℹ️  Modern Linux power management is used (no TLP conflicts)"
