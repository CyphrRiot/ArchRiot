#!/bin/bash

# ==============================================================================
# ArchRiot Utilities Setup
# ==============================================================================
# Simple utility applications installation
# ==============================================================================

# Install system utilities
# Use centralized package installation for consistency
install_packages "btop fastfetch tree wget curl unzip p7zip" "essential"
install_packages "thunar thunar-volman thunar-archive-plugin gvfs gvfs-mtp" "essential"
install_packages "gnome-system-monitor" "essential"
install_packages "iwgtk" "essential"
install_packages "gnome-calculator file-roller secrets" "essential"
install_packages "featherwallet-bin" "optional"
install_packages "fragments" "optional"

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

# Install off-screen window fix script (AMD DPMS bug recovery)
mkdir -p ~/.local/bin/scripts
if [[ -f "$script_dir/../../bin/scripts/fix-offscreen-windows.sh" ]]; then
  cp "$script_dir/../../bin/scripts/fix-offscreen-windows.sh" ~/.local/bin/scripts/
  chmod +x ~/.local/bin/scripts/fix-offscreen-windows.sh
  echo "✓ Off-screen window fix script installed (AMD DPMS recovery)"
else
  echo "⚠ Off-screen window fix script not found in repository bin/scripts"
fi

# Install Brave Private desktop file
if [[ -f "$script_dir/../../applications/brave-private.desktop" ]]; then
  cp "$script_dir/../../applications/brave-private.desktop" ~/.local/share/applications/
  echo "✓ Brave Private desktop file installed"
else
  echo "⚠ Brave Private desktop file not found in repository applications"
fi

# Install web application and custom icon desktop files
echo "🌐 Installing web applications and custom icon apps..."
custom_apps=("Google Messages.desktop" "Proton Mail.desktop" "X.desktop" "Activity.desktop" "zed.desktop")
for app in "${custom_apps[@]}"; do
  if [[ -f "$script_dir/../../applications/$app" ]]; then
    cp "$script_dir/../../applications/$app" ~/.local/share/applications/
    echo "✓ Custom app installed: $(basename "$app" .desktop)"
  else
    echo "⚠ Custom app not found: $app"
  fi
done

# Replicate exact working local system setup
echo "🎯 Installing applications exactly like working local system..."

# Hide ALL system monitor duplicates with NoDisplay=true (matches local working system)
echo '[Desktop Entry]
NoDisplay=true' > ~/.local/share/applications/gnome-system-monitor-kde.desktop

echo '[Desktop Entry]
NoDisplay=true' > ~/.local/share/applications/org.gnome.SystemMonitor.desktop

echo '[Desktop Entry]
NoDisplay=true' > ~/.local/share/applications/btop.desktop

# Install our clean renamed System Monitor (matches local working system)
if [[ -f "$script_dir/../../applications/system-monitor.desktop" ]]; then
  cp "$script_dir/../../applications/system-monitor.desktop" ~/.local/share/applications/
  echo "✓ System Monitor (clean renamed version) installed"
else
  echo "⚠ System Monitor desktop file not found"
fi

# Fix disk space issues that cause 0-byte icon files
echo "💾 Checking disk space before copying icons..."
available_space=$(df ~/.local/share/icons/ | tail -1 | awk '{print $4}')
if [[ $available_space -lt 100000 ]]; then
  echo "⚠ Low disk space detected. Cleaning up..."
  sudo pacman -Scc --noconfirm 2>/dev/null || true
  sudo journalctl --vacuum-time=7d 2>/dev/null || true
  echo "✓ Disk cleanup completed"
fi

# Install iwgtk desktop file with better name and icon
echo "📶 Installing WiFi Manager desktop file..."
if [[ -f "$script_dir/../../applications/iwgtk.desktop" ]]; then
  cp "$script_dir/../../applications/iwgtk.desktop" ~/.local/share/applications/
  echo "✓ WiFi Manager desktop file installed"
else
  echo "⚠ WiFi Manager desktop file not found in repository"
fi

# Copy icons for applications - ensure they're not 0-byte files
echo "🎨 Installing application icons (preventing 0-byte corruption)..."
icons_dir="$script_dir/../../applications/icons"
if [[ -d "$icons_dir" ]]; then
    mkdir -p "$HOME/.local/share/icons/hicolor/256x256/apps"
    for icon_file in "$icons_dir"/*.png; do
        if [[ -f "$icon_file" && -s "$icon_file" ]]; then
            # Only copy non-empty icon files
            cp "$icon_file" "$HOME/.local/share/icons/hicolor/256x256/apps/"
            echo "✓ Copied $(basename "$icon_file")"
        elif [[ -f "$icon_file" ]]; then
            echo "⚠ Skipping empty icon file: $(basename "$icon_file")"
        fi
    done
    echo "✓ Application icons installed (0-byte files prevented)"
fi

# Install custom renamed desktop files with shorter names
echo "✂️ Installing custom desktop files with cleaner names..."

# Clean up old duplicate desktop files
echo "🧹 Cleaning up old duplicate desktop files..."
rm -f ~/.local/share/applications/media-player.desktop
rm -f ~/.local/share/applications/file-manager.desktop
echo "✓ Old duplicate desktop files removed"

# Media Player (renamed from mpv Media Player)
if [[ -f "$script_dir/../../applications/mpv.desktop" ]]; then
  cp "$script_dir/../../applications/mpv.desktop" ~/.local/share/applications/media-player.desktop
  echo "✓ Media Player desktop file installed"
fi

# File Manager (renamed from Thunar File Manager)
if [[ -f "$script_dir/../../applications/thunar.desktop" ]]; then
  cp "$script_dir/../../applications/thunar.desktop" ~/.local/share/applications/thunar.desktop
  echo "✓ File Manager desktop file installed (overwriting original)"
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

# FINAL FIX: Ensure all files are owned by user and icons work
echo "🔧 Final ownership and icon cache fixes..."
sudo chown -R "$USER:$USER" "$HOME/.local/share/applications/" 2>/dev/null || true
sudo chown -R "$USER:$USER" "$HOME/.local/share/icons/" 2>/dev/null || true

# Critical: Update icon cache so custom icons are found
gtk-update-icon-cache -f "$HOME/.local/share/icons/hicolor/" 2>/dev/null || true

echo "✓ Duplicate applications hidden, clean versions installed"
echo "✓ Icon cache updated - custom icons should now work"

# Force update desktop database multiple times to ensure changes take effect
echo "🔄 Updating desktop database..."
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
sleep 1
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
echo "✓ Desktop database forcefully updated"

echo "✅ Utilities setup complete!"
echo "🔄 You may need to restart fuzzel/rofi for changes to take full effect"
