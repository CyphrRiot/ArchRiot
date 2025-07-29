#!/usr/bin/env bash

# ==============================================================================
# ArchRiot Plymouth Manual Version Control
# ==============================================================================
# Simple version-based system - bump version to force reinstall
# ==============================================================================

# Manual version control - increment this to force LUKS reinstall
PLYMOUTH_VERSION="1.0"
VERSION_FILE="$HOME/.config/archriot/plymouth_version.txt"

# Create config directory if needed
mkdir -p "$HOME/.config/archriot"

# Check if reinstall is needed
needs_reinstall() {
    # Check 1: Initial install (no Plymouth installed)
    if ! command -v plymouth &>/dev/null; then
        echo "Initial Plymouth installation needed"
        return 0
    fi

    # Check 2: ArchRiot theme missing or broken
    local current_theme=$(sudo plymouth-set-default-theme 2>/dev/null || echo "none")
    if [[ "$current_theme" != "archriot" ]] || [[ ! -f "/usr/share/plymouth/themes/archriot/archriot.plymouth" ]]; then
        echo "ArchRiot Plymouth theme missing - reinstall needed"
        return 0
    fi

    # Check 3: LUKS header content comparison (most important check)
    echo "🔍 Checking if LUKS header content changed..."
    local new_logo="$HOME/.local/share/archriot/default/plymouth/logo.png"
    local current_logo="/usr/share/plymouth/themes/archriot/logo.png"

    if [[ -f "$new_logo" && -f "$current_logo" ]]; then
        if ! cmp -s "$new_logo" "$current_logo"; then
            echo "LUKS header logo content changed - reinstall needed"
            return 0
        else
            echo "✓ LUKS header logo unchanged - skipping reinstall"
        fi
    fi

    # Check 4: Version changed (manual trigger - only after content check)
    local stored_version=$(cat "$VERSION_FILE" 2>/dev/null || echo "0.0")
    if [[ "$stored_version" != "$PLYMOUTH_VERSION" ]]; then
        echo "Plymouth version updated ($stored_version → $PLYMOUTH_VERSION) but content unchanged - updating version file only"
        echo "$PLYMOUTH_VERSION" > "$VERSION_FILE"
    fi

    echo "Plymouth v$PLYMOUTH_VERSION already installed - skipping"
    exit 0
}

# Check if reinstall needed
if ! needs_reinstall; then
    echo "✅ Plymouth/LUKS v$PLYMOUTH_VERSION already configured"
    exit 0
fi

echo "🔄 Installing Plymouth v$PLYMOUTH_VERSION..."

# Install Plymouth if not already installed
if ! command -v plymouth &>/dev/null; then
  install_packages "plymouth" "essential"
fi

# Backup original mkinitcpio.conf (only when making changes)
backup_timestamp=$(date +"%Y%m%d%H%M%S")
sudo cp /etc/mkinitcpio.conf "/etc/mkinitcpio.conf.bak.${backup_timestamp}"

# Add plymouth to HOOKS array in correct position for LUKS support
if ! grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "plymouth"; then
  echo "Adding Plymouth hook in correct position for LUKS encryption..."

  # For LUKS systems, plymouth must come AFTER consolefont but BEFORE encrypt
  if grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "encrypt"; then
    echo "✓ LUKS encryption detected - placing Plymouth before encrypt hook"
    sudo sed -i '/^HOOKS=/s/\(consolefont.*\)\(block\)/\1 plymouth \2/' /etc/mkinitcpio.conf
  elif grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "base systemd"; then
    echo "✓ Systemd initramfs detected"
    sudo sed -i '/^HOOKS=/s/base systemd/base systemd plymouth/' /etc/mkinitcpio.conf
  elif grep "^HOOKS=" /etc/mkinitcpio.conf | grep -q "base udev"; then
    echo "✓ Standard udev initramfs detected"
    sudo sed -i '/^HOOKS=/s/\(consolefont.*\)\(block\)/\1 plymouth \2/' /etc/mkinitcpio.conf
  else
    echo "⚠ Couldn't determine initramfs type - adding Plymouth after base"
    sudo sed -i '/^HOOKS=/s/base udev/base udev plymouth/' /etc/mkinitcpio.conf
  fi

  # Regenerate initramfs only if we modified hooks
  echo "🔄 Regenerating initramfs with Plymouth support..."
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
      # Need splash, create the archriot file
      echo "splash" | sudo tee -a /etc/cmdline.d/archriot.conf
  fi
  if ! grep -q quiet /etc/cmdline.d/*.conf; then
      # Need quiet, create or append the archriot file
      echo "quiet" | sudo tee -a /etc/cmdline.d/archriot.conf
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
sudo rm -rf /usr/share/plymouth/themes/archriot
sudo rm -rf /usr/share/plymouth/themes/ohmarchy
sudo rm -rf /usr/share/plymouth/themes/archriot

# Remove any cached theme references
sudo rm -f /var/lib/plymouth/themes/default.plymouth
sudo rm -f /etc/plymouth/plymouthd.conf

# PREFER LOCAL FILES: Use local Plymouth files to avoid old logo issues
echo "📦 Installing ArchRiot Plymouth theme from local files..."
sudo mkdir -p /usr/share/plymouth/themes/archriot/

# First try to use local files (preferred - ensures correct logo)
if [ -d "$HOME/.local/share/archriot/default/plymouth" ] && [ -f "$HOME/.local/share/archriot/default/plymouth/logo.png" ]; then
    echo "✓ Using local Plymouth files (ensures correct ArchRiot logo)"
    sudo cp -r "$HOME/.local/share/archriot/default/plymouth"/* /usr/share/plymouth/themes/archriot/

else
    # Fallback: Download from GitHub only if local files missing
    echo "⚠ Local files not found, downloading from GitHub as fallback..."
    TEMP_PLYMOUTH_DIR="/tmp/archriot-plymouth-$$"
    mkdir -p "$TEMP_PLYMOUTH_DIR"

    if curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/archriot.plymouth" -o "$TEMP_PLYMOUTH_DIR/archriot.plymouth" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/archriot.script" -o "$TEMP_PLYMOUTH_DIR/archriot.script" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/logo.png" -o "$TEMP_PLYMOUTH_DIR/logo.png" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/entry.png" -o "$TEMP_PLYMOUTH_DIR/entry.png" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/lock.png" -o "$TEMP_PLYMOUTH_DIR/lock.png" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/bullet.png" -o "$TEMP_PLYMOUTH_DIR/bullet.png" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/progress_bar.png" -o "$TEMP_PLYMOUTH_DIR/progress_bar.png" &&
       curl -fsSL "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/default/plymouth/progress_box.png" -o "$TEMP_PLYMOUTH_DIR/progress_box.png"; then

        echo "✓ Plymouth files downloaded from GitHub"
        sudo cp -r "$TEMP_PLYMOUTH_DIR"/* /usr/share/plymouth/themes/archriot/

        # Update local installation for future use
        mkdir -p "$HOME/.local/share/archriot/default/plymouth"
        cp -r "$TEMP_PLYMOUTH_DIR"/* "$HOME/.local/share/archriot/default/plymouth/"
        echo "✓ Updated local Plymouth files"


    else
        echo "❌ Failed to download Plymouth files from GitHub!"
        return 1
    fi
fi

# Cleanup temp directory (if it was created)
if [ -n "$TEMP_PLYMOUTH_DIR" ] && [ -d "$TEMP_PLYMOUTH_DIR" ]; then
    rm -rf "$TEMP_PLYMOUTH_DIR"
fi

# Verify theme files exist
if [[ ! -f /usr/share/plymouth/themes/archriot/archriot.plymouth ]]; then
    echo "❌ Failed to install Plymouth theme files!"
    return 1
fi

echo "🎨 Setting ArchRiot as default Plymouth theme..."
sudo plymouth-set-default-theme -R archriot

# Verify theme is set correctly
CURRENT_THEME=$(sudo plymouth-set-default-theme)
if [[ "$CURRENT_THEME" != "archriot" ]]; then
    echo "❌ Failed to set ArchRiot theme! Current: $CURRENT_THEME"
    echo "🔧 Available themes:"
    sudo plymouth-set-default-theme --list
    return 1
fi

echo "✅ Plymouth theme verified: $CURRENT_THEME"

# Use existing logo directly (no ASCII generation needed)
echo "✓ Using existing Plymouth logo - no regeneration needed"

# Record successful installation
echo "$PLYMOUTH_VERSION" > "$VERSION_FILE"
echo "✅ Plymouth v$PLYMOUTH_VERSION installation complete"
echo "🚀 To force reinstall, increment PLYMOUTH_VERSION in this script"
