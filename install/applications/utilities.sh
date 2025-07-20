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
    install_packages_clean "gnome-calculator file-roller secrets" "Installing essential tools" "GREEN"
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
      file-roller \
      secrets \

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
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
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
echo "DEBUG: script_dir = $script_dir"
echo "DEBUG: PWD = $(pwd)"
hidden_dir="$script_dir/../../applications/hidden"
echo "DEBUG: hidden_dir = $hidden_dir"
echo "DEBUG: hidden_dir resolved = $(readlink -f "$hidden_dir" 2>/dev/null || echo "FAILED TO RESOLVE")"
echo "DEBUG: Directory exists? $(test -d "$hidden_dir" && echo "YES" || echo "NO")"
if [[ -d "$hidden_dir" ]]; then
  echo "DEBUG: Files in hidden_dir:"
  ls -la "$hidden_dir"
  # Force copy all hidden desktop files to suppress system ones
  for hidden_file in "$hidden_dir"/*.desktop; do
    if [[ -f "$hidden_file" ]]; then
      echo "DEBUG: Copying $hidden_file to $HOME/.local/share/applications/"
      cp "$hidden_file" "$HOME/.local/share/applications/"
      echo "DEBUG: Copy result: $?"
    fi
  done
  hidden_count=$(find "$hidden_dir" -name "*.desktop" 2>/dev/null | wc -l)
  echo "✓ $hidden_count hidden applications installed"
else
  echo "⚠ Hidden applications directory not found at: $hidden_dir"
  echo "DEBUG: Let's see what's around here:"
  echo "DEBUG: Contents of $script_dir:"
  ls -la "$script_dir" 2>/dev/null || echo "script_dir doesn't exist"
  echo "DEBUG: Contents of $script_dir/../..:"
  ls -la "$script_dir/../.." 2>/dev/null || echo "parent dir doesn't exist"
fi

# Install iwgtk desktop file with better name and icon
echo "📶 Installing WiFi Manager desktop file..."
if [[ -f "$script_dir/../../applications/iwgtk.desktop" ]]; then
  cp "$script_dir/../../applications/iwgtk.desktop" ~/.local/share/applications/
  echo "✓ WiFi Manager desktop file installed"
else
  echo "⚠ WiFi Manager desktop file not found in repository"
fi

# Copy icons for applications (critical for missing icons issue)
echo "🎨 Installing application icons..."
icons_dir="$script_dir/../../applications/icons"
if [[ -d "$icons_dir" ]]; then
    mkdir -p "$HOME/.local/share/icons/hicolor/256x256/apps"
    for icon_file in "$icons_dir"/*.png; do
        if [[ -f "$icon_file" ]]; then
            cp "$icon_file" "$HOME/.local/share/icons/hicolor/256x256/apps/"
        fi
    done
    echo "✓ Application icons installed"
fi

# Install custom renamed desktop files with shorter names
echo "✂️ Installing custom desktop files with cleaner names..."

# Clean up ONLY old custom desktop files from previous installs (NOT hidden files)
echo "🧹 Cleaning up old custom desktop files from previous installs..."
rm -f ~/.local/share/applications/system-monitor.desktop
rm -f ~/.local/share/applications/media-player.desktop
rm -f ~/.local/share/applications/file-manager.desktop
echo "✓ Old custom desktop files removed (hidden files preserved)"

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
if [[ -f "$script_dir/../../bin/upgrade-system" ]]; then
  cp "$script_dir/../../bin/upgrade-system" ~/.local/bin/
  chmod +x ~/.local/bin/upgrade-system
  echo "✓ ArchRiot upgrade-system installed to ~/.local/bin/"
  echo "  Usage: upgrade-system --help"
else
  echo "⚠ ArchRiot upgrade-system script not found in repository bin"
fi

# Force update desktop database multiple times to ensure changes take effect
echo "🔄 Updating desktop database..."
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
sleep 1
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
echo "✓ Desktop database forcefully updated"

echo "✅ Utilities setup complete!"
echo "🔄 You may need to restart fuzzel/rofi for changes to take full effect"
