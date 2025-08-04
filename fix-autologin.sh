#!/bin/bash
# Emergency Auto-Login Fix Script for ArchRiot
# Fixes the broken auto-login caused by $SUDO_USER issue

echo "🚨 ArchRiot Auto-Login Emergency Fix"
echo "======================================"

# Get the actual username (not root)
if [ "$EUID" -eq 0 ]; then
    echo "❌ Don't run this as root! Run as your regular user."
    exit 1
fi

USERNAME=$(whoami)
echo "🔧 Fixing auto-login for user: $USERNAME"

# Create the correct auto-login configuration
sudo mkdir -p /etc/systemd/system/getty@tty1.service.d

echo "📝 Creating correct getty override configuration..."
echo -e "[Service]\nExecStart=\nExecStart=-/usr/bin/agetty --autologin $USERNAME --noclear %I \$TERM" | sudo tee /etc/systemd/system/getty@tty1.service.d/override.conf >/dev/null

echo "🔄 Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "✅ Auto-login fix applied!"
echo "💡 Reboot your system to test the fix."
echo ""
echo "If you're SSH'd in, run: sudo reboot"
