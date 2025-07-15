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
  networkmanager-applet \
  iwgtk

# Install essential tools
yay -S --noconfirm --needed \
  gnome-calculator \
  file-roller

# Install financial tools
yay -S --noconfirm --needed \
  featherwallet-bin

# install fragments (torrent client)
yay -S --noconfirm --needed \
  fragments

# Install custom desktop files and launchers
echo "🎯 Installing custom desktop integrations..."
mkdir -p ~/.local/share/applications ~/.local/bin ~/.local/share/icons/hicolor/256x256/apps

# Install Feather Wallet desktop file and icon
if [[ -f "$HOME/.local/share/omarchy/applications/feather-wallet.desktop" ]]; then
  cp "$HOME/.local/share/omarchy/applications/feather-wallet.desktop" ~/.local/share/applications/
  echo "✓ Feather Wallet desktop file installed"
else
  echo "⚠ Feather Wallet desktop file not found in OhmArchy applications"
fi

# Download Feather Wallet icon
if wget -q -O ~/.local/share/icons/hicolor/256x256/apps/feather-wallet.png "https://raw.githubusercontent.com/feather-wallet/feather/master/src/assets/images/feather.png"; then
  echo "✓ Feather Wallet icon downloaded"
else
  echo "⚠ Failed to download Feather Wallet icon"
fi

# Install Signal Wayland launcher and desktop file
if [[ -f "$HOME/.local/share/omarchy/bin/signal-wayland" ]]; then
  cp "$HOME/.local/share/omarchy/bin/signal-wayland" ~/.local/bin/
  chmod +x ~/.local/bin/signal-wayland
  echo "✓ Signal Wayland launcher installed"
else
  echo "⚠ Signal Wayland launcher not found in OhmArchy bin"
fi

if [[ -f "$HOME/.local/share/omarchy/applications/signal-desktop.desktop" ]]; then
  cp "$HOME/.local/share/omarchy/applications/signal-desktop.desktop" ~/.local/share/applications/
  echo "✓ Signal desktop file installed with Wayland support"
else
  echo "⚠ Signal desktop file not found in OhmArchy applications"
fi

# Install Brave Private desktop file
if [[ -f "$HOME/.local/share/omarchy/applications/brave-private.desktop" ]]; then
  cp "$HOME/.local/share/omarchy/applications/brave-private.desktop" ~/.local/share/applications/
  echo "✓ Brave Private desktop file installed"
else
  echo "⚠ Brave Private desktop file not found in OhmArchy applications"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Utilities setup complete!"
