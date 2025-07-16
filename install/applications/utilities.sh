#!/bin/bash

# ==============================================================================
# OhmArchy Utilities Setup
# ==============================================================================
# Simple utility applications installation
# ==============================================================================

# Install system utilities
if command -v install_packages_clean &>/dev/null; then
    install_packages_clean "btop neofetch fastfetch tree wget curl unzip p7zip" "Installing system utilities" "GREEN"
    install_packages_clean "thunar thunar-volman thunar-archive-plugin gvfs gvfs-mtp" "Installing file management" "GREEN"
    install_packages_clean "gnome-system-monitor" "Installing system monitoring" "GREEN"
    install_packages_clean "iwgtk" "Installing network utilities" "GREEN"
    install_packages_clean "gnome-calculator file-roller" "Installing essential tools" "GREEN"
    install_packages_clean "featherwallet-bin" "Installing financial tools" "GREEN"
    install_packages_clean "fragments" "Installing torrent client" "GREEN"
else
    # Fallback to direct yay commands
    yay -S --noconfirm --needed \
      btop \
      neofetch \
      fastfetch \
      tree \
      wget \
      curl \
      unzip \
      p7zip

    yay -S --noconfirm --needed \
      thunar \
      thunar-volman \
      thunar-archive-plugin \
      gvfs \
      gvfs-mtp

    yay -S --noconfirm --needed \
      gnome-system-monitor

    yay -S --noconfirm --needed \
      iwgtk

    yay -S --noconfirm --needed \
      gnome-calculator \
      file-roller

    yay -S --noconfirm --needed \
      featherwallet-bin

    yay -S --noconfirm --needed \
      fragments
fi

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

# Install hidden applications to suppress unwanted launchers
echo "🙈 Installing hidden applications..."
if [[ -d "$HOME/.local/share/omarchy/applications/hidden" ]]; then
  cp "$HOME/.local/share/omarchy/applications/hidden"/*.desktop ~/.local/share/applications/ 2>/dev/null || true
  hidden_count=$(find "$HOME/.local/share/omarchy/applications/hidden" -name "*.desktop" 2>/dev/null | wc -l)
  echo "✓ $hidden_count hidden applications installed"
else
  echo "⚠ Hidden applications directory not found"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Utilities setup complete!"
