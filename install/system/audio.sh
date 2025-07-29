#!/bin/bash

# ==============================================================================
# ArchRiot Audio System Setup
# ==============================================================================
# Simple audio setup - PipeWire with essential utilities
# ==============================================================================

# Install PipeWire audio system with conflict detection
echo "📦 Installing PipeWire audio system..."

# Check for conflicting packages first
echo "🔍 Checking for audio system conflicts..."
conflicting_packages=""
if pacman -Qs "^pulseaudio$" &>/dev/null; then
    conflicting_packages="$conflicting_packages pulseaudio"
fi
if pacman -Qs "^jack2$" &>/dev/null; then
    conflicting_packages="$conflicting_packages jack2"
fi

# Remove conflicting packages if found
if [[ -n "$conflicting_packages" ]]; then
    echo "⚠ Found conflicting audio packages:$conflicting_packages"
    echo "🗑️ Removing conflicting packages..."
    sudo pacman -Rdd --noconfirm $conflicting_packages || echo "⚠ Some packages couldn't be removed"
fi

# Install PipeWire with proper conflict resolution
echo "🔄 Installing PipeWire audio system..."
if sudo pacman -S --needed --noconfirm pipewire pipewire-alsa pipewire-pulse pipewire-audio pipewire-jack gst-plugin-pipewire libpipewire wireplumber; then
    echo "✓ PipeWire audio system installed successfully"
elif sudo pacman -Syu --noconfirm --needed pipewire pipewire-alsa pipewire-pulse wireplumber; then
    echo "✓ PipeWire installed via system upgrade"
else
    echo "❌ PipeWire installation failed"
    echo "ℹ️ Manual intervention may be required"
    echo "ℹ️ Try running: sudo pacman -Syu && sudo pacman -S pipewire-pulse"
    # Don't exit - continue with what we have
fi
install_packages "pavucontrol pamixer playerctl" "essential"

# Enable PipeWire services
systemctl --user enable --now pipewire.service pipewire-pulse.service wireplumber.service

echo "✅ Audio setup complete!"
