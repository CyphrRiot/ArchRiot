# Install asdcontrol for controlling brightness on Apple Displays
if ! command -v asdcontrol &>/dev/null; then
  echo "📱 Installing asdcontrol for Apple Display brightness control..."

  temp_dir="/tmp/asdcontrol-$$"

  if git clone https://github.com/nikosdion/asdcontrol.git "$temp_dir" &&
     cd "$temp_dir" && make && sudo make install; then

    # Setup sudo-less controls
    echo "$USER ALL=(ALL) NOPASSWD: /usr/local/bin/asdcontrol" | sudo tee /etc/sudoers.d/asdcontrol >/dev/null
    sudo chmod 440 /etc/sudoers.d/asdcontrol

    echo "✓ asdcontrol installed and configured"
  else
    echo "⚠ asdcontrol installation failed"
  fi

  cd - >/dev/null 2>&1
  rm -rf "$temp_dir"
fi
