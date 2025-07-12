#!/usr/bin/env bash

if ! command -v plymouth &>/dev/null; then
  yay -S --noconfirm --needed plymouth

  # Single backup timestamp for all operations
  BACKUP_TS=$(date +"%Y%m%d%H%M%S")

  # Backup and modify mkinitcpio.conf
  sudo cp /etc/mkinitcpio.conf "/etc/mkinitcpio.conf.bak.${BACKUP_TS}"

  if grep -q "^HOOKS=.*base systemd" /etc/mkinitcpio.conf; then
    sudo sed -i '/^HOOKS=/s/base systemd/base systemd plymouth/' /etc/mkinitcpio.conf
  elif grep -q "^HOOKS=.*base udev" /etc/mkinitcpio.conf; then
    sudo sed -i '/^HOOKS=/s/base udev/base udev plymouth/' /etc/mkinitcpio.conf
  else
    echo "Couldn't add Plymouth hook - manual intervention required"
    exit 1
  fi

  sudo mkinitcpio -P

  # Configure bootloader with splash parameters
  configure_bootloader() {
    if [ -d "/boot/loader/entries" ]; then
      # systemd-boot configuration
      for entry in /boot/loader/entries/*.conf; do
        [[ "$(basename "$entry")" == *"fallback"* ]] && continue
        grep -q "splash" "$entry" || sudo sed -i '/^options/ s/$/ splash quiet/' "$entry"
      done
    elif [ -f "/etc/default/grub" ]; then
      # GRUB configuration
      sudo cp /etc/default/grub "/etc/default/grub.bak.${BACKUP_TS}"
      if ! grep -q "splash" /etc/default/grub; then
        sudo sed -i '/^GRUB_CMDLINE_LINUX_DEFAULT=/ s/"$/ splash quiet"/' /etc/default/grub
        sudo grub-mkconfig -o /boot/grub/grub.cfg
      fi
    elif [ -d "/etc/cmdline.d" ]; then
      # UKI setup with cmdline.d
      grep -q splash /etc/cmdline.d/*.conf 2>/dev/null || echo "splash" | sudo tee -a /etc/cmdline.d/omarchy.conf
      grep -q quiet /etc/cmdline.d/*.conf 2>/dev/null || echo "quiet" | sudo tee -a /etc/cmdline.d/omarchy.conf
    elif [ -f "/etc/kernel/cmdline" ]; then
      # UKI setup with kernel cmdline
      sudo cp /etc/kernel/cmdline "/etc/kernel/cmdline.bak.${BACKUP_TS}"
      cmdline=$(cat /etc/kernel/cmdline)
      [[ "$cmdline" =~ splash ]] || cmdline="$cmdline splash"
      [[ "$cmdline" =~ quiet ]] || cmdline="$cmdline quiet"
      echo "${cmdline// / }" | sudo tee /etc/kernel/cmdline
    else
      echo "Manual kernel parameter setup required: splash quiet"
    fi
  }

  configure_bootloader

  # Install Plymouth theme with logo handling
  THEME_DIR="/usr/share/plymouth/themes/omarchy"
  CUSTOM_LOGO="$THEME_DIR/logo.png"

  # Backup existing custom logo if different from default
  [ -f "$CUSTOM_LOGO" ] && sudo cp "$CUSTOM_LOGO" "/tmp/custom_logo_backup.png"

  # Install theme
  sudo cp -r "$HOME/.local/share/omarchy/default/plymouth" "$THEME_DIR/"

  # Restore custom logo if backed up
  [ -f "/tmp/custom_logo_backup.png" ] && {
    sudo cp "/tmp/custom_logo_backup.png" "$CUSTOM_LOGO"
    rm -f "/tmp/custom_logo_backup.png"
  }

  # Generate OhmArchy ASCII logo
  generate_logo() {
    command -v magick >/dev/null || sudo pacman -S --noconfirm --needed imagemagick

    local ascii_art=' ▄██████▄   ██    ██  ████████████    ▄████████    ▄████████    ▄█    █▄    ▄██   ▄
███    ███  ██    ██  ██    ██   ██   ███    ███   ███    ███   ███    ███   ███   ██▄
███    ███  ██▀▀▀▀██  ██    ██   ██   ███    ███   ███    █▀    ███    ███   ███▄▄▄███
███    ███  ██    ██  ████████████   ███    ███  ▄███▄▄▄▄██▀  ▄███▄▄▄▄███▄▄ ▀▀▀▀▀▀███
███    ███  ██▀▀▀▀██  ██    ██   ██ ▀███████████ ▀▀███▀▀▀▀▀   ▀▀███▀▀▀▀███▀  ▄██   ███
███    ███  ██    ██  ██    ██   ██   ███    ███ ▀███████████   ███    ███   ███   ███
███    ███  ██    ██  ██    ██   ██   ███    ███   ███    ███   ███    ███   ███   ███
 ▀██████▀   ██    ██  ██    ██   ██   ███    █▀    ███    ███   ███    █▀     ▀█████▀
                                                  ███    ███                         '

    for font in DejaVu-Sans-Mono Liberation-Mono monospace; do
      if magick -size 800x168 -background '#1a1b26' -fill '#c0caf5' -font "$font" -pointsize 10 -gravity center label:"$ascii_art" "/tmp/ohmarchy_logo.png" 2>/dev/null; then
        sudo cp "/tmp/ohmarchy_logo.png" "$CUSTOM_LOGO"
        sudo chown root:root "$CUSTOM_LOGO"
        sudo chmod 644 "$CUSTOM_LOGO"
        rm -f "/tmp/ohmarchy_logo.png"
        echo "✓ Logo generated with font: $font"
        return 0
      fi
    done
    echo "⚠ Logo generation failed, using default"
  }

  generate_logo

  # Check for persistent backup
  BACKUP_LOGO="$HOME/.config/omarchy/plymouth-backup/custom_logo.png"
  [ -f "$BACKUP_LOGO" ] && {
    sudo cp "$BACKUP_LOGO" "$CUSTOM_LOGO"
    sudo chown root:root "$CUSTOM_LOGO"
    sudo chmod 644 "$CUSTOM_LOGO"
    echo "✓ Custom logo restored from backup"
  }

  sudo plymouth-set-default-theme -R omarchy
fi
