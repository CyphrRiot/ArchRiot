#!/usr/bin/env bash

# Install Plymouth if not already installed
if ! command -v plymouth &>/dev/null; then
  yay -S --noconfirm --needed plymouth
fi

# Always update Plymouth configuration (even on re-installs)
# Backup original mkinitcpio.conf just in case
backup_timestamp=$(date +"%Y%m%d%H%M%S")
sudo cp /etc/mkinitcpio.conf "/etc/mkinitcpio.conf.bak.${backup_timestamp}"

# Add plymouth to HOOKS array after 'base udev' or 'base systemd' if not already present
if ! grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "plymouth"; then
  if grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "base systemd"; then
    sudo sed -i '/^HOOKS=/s/base systemd/base systemd plymouth/' /etc/mkinitcpio.conf
  elif grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "base udev"; then
    sudo sed -i '/^HOOKS=/s/base udev/base udev plymouth/' /etc/mkinitcpio.conf
  else
    echo "Couldn't add the Plymouth hook"
  fi

  # Regenerate initramfs only if we modified hooks
  sudo mkinitcpio -P
fi

# Add kernel parameters for Plymouth (systemd-boot only)
if [ -d "/boot/loader/entries" ]; then
  echo "Detected systemd-boot"

  for entry in /boot/loader/entries/*.conf; do
    if [ -f "$entry" ]; then
      # Skip fallback entries
      if [[ "$(basename "$entry")" == *"fallback"* ]]; then
        echo "Skipped: $(basename "$entry") (fallback entry)"
        continue
      fi

      # Skip if splash it already present for some reason
      if ! grep -q "splash" "$entry"; then
        sudo sed -i '/^options/ s/$/ splash quiet/' "$entry"
      else
        echo "Skipped: $(basename "$entry") (splash already present)"
      fi
    fi
  done
elif [ -f "/etc/default/grub" ]; then
  echo "Detected grub"
  # Backup GRUB config before modifying
  backup_timestamp=$(date +"%Y%m%d%H%M%S")
  sudo cp /etc/default/grub "/etc/default/grub.bak.${backup_timestamp}"

  # Check if splash is already in GRUB_CMDLINE_LINUX_DEFAULT
  if ! grep -q "GRUB_CMDLINE_LINUX_DEFAULT.*splash" /etc/default/grub; then
    # Get current GRUB_CMDLINE_LINUX_DEFAULT value
    current_cmdline=$(grep "^GRUB_CMDLINE_LINUX_DEFAULT=" /etc/default/grub | cut -d'"' -f2)

    # Add splash and quiet if not present
    new_cmdline="$current_cmdline"
    if [[ ! "$current_cmdline" =~ splash ]]; then
      new_cmdline="$new_cmdline splash"
    fi
    if [[ ! "$current_cmdline" =~ quiet ]]; then
      new_cmdline="$new_cmdline quiet"
    fi

    # Trim any leading/trailing spaces
    new_cmdline=$(echo "$new_cmdline" | xargs)

    sudo sed -i "s/^GRUB_CMDLINE_LINUX_DEFAULT=\".*\"/GRUB_CMDLINE_LINUX_DEFAULT=\"$new_cmdline\"/" /etc/default/grub

    # Regenerate grub config
    sudo grub-mkconfig -o /boot/grub/grub.cfg
  else
    echo "GRUB already configured with splash kernel parameters"
  fi
elif [ -d "/etc/cmdline.d" ]; then
  echo "Detected a UKI setup"
  # Relying on mkinitcpio to assemble a UKI
  # https://wiki.archlinux.org/title/Unified_kernel_image
  if ! grep -q splash /etc/cmdline.d/*.conf; then
      # Need splash, create the ohmarchy file
      echo "splash" | sudo tee -a /etc/cmdline.d/ohmarchy.conf
  fi
  if ! grep -q quiet /etc/cmdline.d/*.conf; then
      # Need quiet, create or append the ohmarchy file
      echo "quiet" | sudo tee -a /etc/cmdline.d/ohmarchy.conf
  fi
elif [ -f "/etc/kernel/cmdline" ]; then
  # Alternate UKI kernel cmdline location
  echo "Detected a UKI setup"

  # Backup kernel cmdline config before modifying
  backup_timestamp=$(date +"%Y%m%d%H%M%S")
  sudo cp /etc/kernel/cmdline "/etc/kernel/cmdline.bak.${backup_timestamp}"

  current_cmdline=$(cat /etc/kernel/cmdline)

  # Add splash and quiet if not present
  new_cmdline="$current_cmdline"
  if [[ ! "$current_cmdline" =~ splash ]]; then
      new_cmdline="$new_cmdline splash"
  fi
  if [[ ! "$current_cmdline" =~ quiet ]]; then
      new_cmdline="$new_cmdline quiet"
  fi

  # Trim any leading/trailing spaces
  new_cmdline=$(echo "$new_cmdline" | xargs)

  # Write new file
  echo $new_cmdline | sudo tee /etc/kernel/cmdline
else
  echo ""
  echo "Neither systemd-boot nor GRUB detected. Please manually add these kernel parameters:"
  echo "  - splash (to see the graphical splash screen)"
  echo "  - quiet (for silent boot)"
  echo ""
fi

# AGGRESSIVE CLEANUP - Remove ALL old Plymouth themes and cached files
echo "🧹 Cleaning up old Plymouth themes..."
sudo rm -rf /usr/share/plymouth/themes/omarchy
sudo rm -rf /usr/share/plymouth/themes/ohmarchy

# Remove any cached theme references
sudo rm -f /var/lib/plymouth/themes/default.plymouth
sudo rm -f /etc/plymouth/plymouthd.conf

# FORCE UPDATE: Fetch latest Plymouth files from GitHub for re-installs
echo "🔄 Fetching latest Plymouth theme files..."
TEMP_PLYMOUTH_DIR="/tmp/ohmarchy-plymouth-$$"
mkdir -p "$TEMP_PLYMOUTH_DIR"

# Download latest Plymouth files directly from GitHub
if curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/ohmarchy.plymouth" -o "$TEMP_PLYMOUTH_DIR/ohmarchy.plymouth" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/ohmarchy.script" -o "$TEMP_PLYMOUTH_DIR/ohmarchy.script" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/logo.png" -o "$TEMP_PLYMOUTH_DIR/logo.png" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/entry.png" -o "$TEMP_PLYMOUTH_DIR/entry.png" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/lock.png" -o "$TEMP_PLYMOUTH_DIR/lock.png" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/bullet.png" -o "$TEMP_PLYMOUTH_DIR/bullet.png" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/progress_bar.png" -o "$TEMP_PLYMOUTH_DIR/progress_bar.png" &&
   curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/default/plymouth/progress_box.png" -o "$TEMP_PLYMOUTH_DIR/progress_box.png"; then

    echo "✓ Latest Plymouth files downloaded from GitHub"

    # Install from downloaded files
    echo "📦 Installing OhmArchy Plymouth theme..."
    sudo mkdir -p /usr/share/plymouth/themes/ohmarchy/
    sudo cp -r "$TEMP_PLYMOUTH_DIR"/* /usr/share/plymouth/themes/ohmarchy/

    # Update local installation too
    mkdir -p "$HOME/.local/share/omarchy/default/plymouth"
    cp -r "$TEMP_PLYMOUTH_DIR"/* "$HOME/.local/share/omarchy/default/plymouth/"
    echo "✓ Updated local Plymouth files"

else
    echo "⚠ Failed to download latest files, using local files as fallback..."
    sudo mkdir -p /usr/share/plymouth/themes/ohmarchy/
    sudo cp -r "$HOME/.local/share/omarchy/default/plymouth"/* /usr/share/plymouth/themes/ohmarchy/
fi

# Cleanup temp directory
rm -rf "$TEMP_PLYMOUTH_DIR"

# Verify theme files exist
if [[ ! -f /usr/share/plymouth/themes/ohmarchy/ohmarchy.plymouth ]]; then
    echo "❌ Failed to install Plymouth theme files!"
    exit 1
fi

echo "🎨 Setting OhmArchy as default Plymouth theme..."
sudo plymouth-set-default-theme -R ohmarchy

# Verify theme is set correctly
CURRENT_THEME=$(sudo plymouth-set-default-theme)
if [[ "$CURRENT_THEME" != "ohmarchy" ]]; then
    echo "❌ Failed to set OhmArchy theme! Current: $CURRENT_THEME"
    echo "🔧 Available themes:"
    sudo plymouth-set-default-theme --list
    exit 1
fi

echo "✅ Plymouth theme verified: $CURRENT_THEME (latest version from GitHub)"
