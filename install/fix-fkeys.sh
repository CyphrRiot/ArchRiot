if [[ ! -f /etc/modprobe.d/hid_apple.conf ]]; then
  echo "options hid_apple fnmode=2" | sudo tee /etc/modprobe.d/hid_apple.conf

  # Change to safe directory before rebuilding initramfs to avoid getcwd issues
  ORIGINAL_DIR="$(pwd)"
  cd /tmp 2>/dev/null || cd /

  sudo mkinitcpio -P

  # Return to original directory
  cd "$ORIGINAL_DIR" 2>/dev/null || true
fi
