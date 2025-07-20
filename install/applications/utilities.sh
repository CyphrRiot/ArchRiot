#!/bin/bash

# ==============================================================================
# ArchRiot Utilities Setup
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
# Feather Wallet
local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "$script_dir/../../applications/feather-wallet.desktop" ]]; then
  cp "$script_dir/../../applications/feather-wallet.desktop" ~/.local/share/applications/
  echo "✓ Feather Wallet desktop file installed"
else
  echo "⚠ Feather Wallet desktop file not found in repository applications"
fi

# Download Feather Wallet icon
if wget -q -O ~/.local/share/icons/hicolor/256x256/apps/feather-wallet.png "https://raw.githubusercontent.com/feather-wallet/feather/master/src/assets/images/feather.png"; then
  echo "✓ Feather Wallet icon downloaded"
else
  echo "⚠ Failed to download Feather Wallet icon"
fi

# Install Signal Wayland launcher and desktop file
# Note: signal-desktop package is installed by communication.sh module
if [[ -f "$script_dir/../../bin/signal-wayland" ]]; then
  cp "$script_dir/../../bin/signal-wayland" ~/.local/bin/
  chmod +x ~/.local/bin/signal-wayland
  echo "✓ Signal Wayland launcher installed"
else
  echo "⚠ Signal Wayland launcher not found in repository bin"
fi

if [[ -f "$script_dir/../../applications/signal-desktop.desktop" ]]; then
  cp "$script_dir/../../applications/signal-desktop.desktop" ~/.local/share/applications/
  echo "✓ Signal desktop file installed with Wayland support"
else
  echo "⚠ Signal desktop file not found in repository applications"
fi

# Install Brave Private desktop file
if [[ -f "$script_dir/../../applications/brave-private.desktop" ]]; then
  cp "$script_dir/../../applications/brave-private.desktop" ~/.local/share/applications/
  echo "✓ Brave Private desktop file installed"
else
  echo "⚠ Brave Private desktop file not found in repository applications"
fi

# Install hidden applications to suppress unwanted launchers
echo "🙈 Installing hidden applications..."
if [[ -d "$script_dir/../../applications/hidden" ]]; then
  cp "$script_dir/../../applications/hidden"/*.desktop ~/.local/share/applications/ 2>/dev/null || true
  hidden_count=$(find "$script_dir/../../applications/hidden" -name "*.desktop" 2>/dev/null | wc -l)
  echo "✓ $hidden_count hidden applications installed"
else
  echo "⚠ Hidden applications directory not found"
fi

# Install iwgtk desktop file with better name and icon
echo "📶 Installing WiFi Manager desktop file..."
if [[ -f "$script_dir/../../applications/iwgtk.desktop" ]]; then
  cp "$script_dir/../../applications/iwgtk.desktop" ~/.local/share/applications/
  echo "✓ WiFi Manager desktop file installed"
else
  echo "⚠ WiFi Manager desktop file not found in repository"
fi

# Install custom renamed desktop files with shorter names
echo "✂️ Installing custom desktop files with cleaner names..."

# Clean up any conflicting desktop files from previous installs
echo "🧹 Cleaning up conflicting desktop files from previous installs..."
rm -f ~/.local/share/applications/gnome-system-monitor-kde.desktop
rm -f ~/.local/share/applications/mpv.desktop
rm -f ~/.local/share/applications/thunar.desktop
rm -f ~/.local/share/applications/thunar-volman-settings.desktop
echo "✓ Old conflicting desktop files removed"

# System Monitor (renamed from GNOME System Monitor)
if [[ -f "$script_dir/../../applications/gnome-system-monitor-kde.desktop" ]]; then
  cp "$script_dir/../../applications/gnome-system-monitor-kde.desktop" ~/.local/share/applications/system-monitor.desktop
  echo "✓ System Monitor desktop file installed"
fi

# Media Player (renamed from mpv Media Player)
if [[ -f "$script_dir/../../applications/mpv.desktop" ]]; then
  cp "$script_dir/../../applications/mpv.desktop" ~/.local/share/applications/media-player.desktop
  echo "✓ Media Player desktop file installed"
fi

# File Manager (renamed from Thunar File Manager)
if [[ -f "$script_dir/../../applications/thunar.desktop" ]]; then
  cp "$script_dir/../../applications/thunar.desktop" ~/.local/share/applications/file-manager.desktop
  echo "✓ File Manager desktop file installed"
fi

# Removable Drives (shortened from Removable Drives and Media)
if [[ -f "$script_dir/../../applications/thunar-volman-settings.desktop" ]]; then
  cp "$script_dir/../../applications/thunar-volman-settings.desktop" ~/.local/share/applications/
  echo "✓ Removable Drives desktop file installed"
fi

# Install ArchRiot upgrade-system script
echo "🚀 Installing ArchRiot upgrade-system utility..."
local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "$script_dir/../../bin/upgrade-system" ]]; then
  cp "$script_dir/../../bin/upgrade-system" ~/.local/bin/
  chmod +x ~/.local/bin/upgrade-system
  echo "✓ ArchRiot upgrade-system installed to ~/.local/bin/"
  echo "  Usage: upgrade-system --help"
else
  echo "⚠ ArchRiot upgrade-system script not found in repository bin"
fi

# Update desktop database
update-desktop-database ~/.local/share/applications/

echo "✅ Utilities setup complete!"
