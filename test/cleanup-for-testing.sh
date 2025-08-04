#!/bin/bash

# ================================================================================
# ArchRiot Fresh Install Testing Cleanup Script
# ================================================================================
# Removes all ArchRiot-managed configurations for clean testing
# Run this before testing the installer to simulate a fresh system
# ================================================================================

echo "🧹 ArchRiot Fresh Install Testing Cleanup"
echo "=========================================="
echo ""
echo "⚠️  WARNING: This will remove ALL ArchRiot configurations!"
echo "   Only run this if you want to test a fresh installation."
echo ""
read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

echo ""
echo "💾 Creating backup of existing configs..."
BACKUP_DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="$HOME/config_$BACKUP_DATE.tgz"

# Create list of directories that exist (with proper path expansion)
DIRS_TO_BACKUP=()
DIRS_TO_CHECK=(
    "$HOME/.local/share/archriot"
    "$HOME/.config/fish"
    "$HOME/.config/nvim"
    "$HOME/.config/tmux"
    "$HOME/.config/environment.d"
    "$HOME/.config/fastfetch"
    "$HOME/.config/hypr"
    "$HOME/.config/waybar"
    "$HOME/.config/fuzzel"
    "$HOME/.config/mako"
    "$HOME/.config/gtk-3.0"
    "$HOME/.config/gtk-4.0"
    "$HOME/.config/ghostty"
    "$HOME/.config/Thunar"
    "$HOME/.config/btop"
    "$HOME/.config/text-editor"
    "$HOME/.config/xdg-desktop-portal"
    "$HOME/.config/zed"
    "$HOME/.config/archriot"
    "$HOME/.local/share/applications"
    "$HOME/.icons"
    "$HOME/.local/bin/upgrade-system"
)

# Check systemd files separately (glob patterns)
for file in "$HOME/.config/systemd/user/battery-monitor."* "$HOME/.config/systemd/user/version-check."*; do
    if [[ -e "$file" ]]; then
        DIRS_TO_CHECK+=("$file")
    fi
done

# Add existing directories to backup list
for dir in "${DIRS_TO_CHECK[@]}"; do
    if [[ -e "$dir" ]]; then
        DIRS_TO_BACKUP+=("$dir")
    fi
done

if [[ ${#DIRS_TO_BACKUP[@]} -gt 0 ]]; then
    tar -czf "$BACKUP_FILE" "${DIRS_TO_BACKUP[@]}" 2>/dev/null
    echo "✓ Backup created: $BACKUP_FILE"
    echo "  To restore: cd / && tar -xzf $BACKUP_FILE"
else
    echo "✓ No existing configs to backup"
fi

echo ""
echo "🗑️  Removing ArchRiot installation..."
rm -rf ~/.local/share/archriot
echo "✓ Removed ~/.local/share/archriot"

echo ""
echo "🗑️  Removing ArchRiot-managed config directories..."

# Core configs
rm -rf ~/.config/fish
rm -rf ~/.config/nvim
rm -rf ~/.config/tmux
rm -rf ~/.config/environment.d
rm -rf ~/.config/fastfetch

# Desktop environment
rm -rf ~/.config/hypr
rm -rf ~/.config/waybar
rm -rf ~/.config/fuzzel
rm -rf ~/.config/mako
rm -rf ~/.config/gtk-3.0
rm -rf ~/.config/gtk-4.0

# Applications
rm -rf ~/.config/ghostty
rm -rf ~/.config/Thunar
rm -rf ~/.config/btop
rm -rf ~/.config/text-editor
rm -rf ~/.config/xdg-desktop-portal
rm -rf ~/.config/zed

# ArchRiot runtime configs (created by scripts)
rm -rf ~/.config/archriot

# User directories that get desktop files
rm -rf ~/.local/share/applications

# User directories that get icons
rm -rf ~/.icons

# Systemd user services
rm -rf ~/.config/systemd/user/battery-monitor.*
rm -rf ~/.config/systemd/user/version-check.*

# Local bin directory (where upgrade-system gets copied)
rm -f ~/.local/bin/upgrade-system

echo ""
echo "✅ Fresh install cleanup complete!"
echo ""
echo "📋 Next steps for testing:"
echo "1. Copy repo to ~/.local/share/archriot/"
echo "2. Run: ~/.local/share/archriot/install/archriot"
echo "3. Test complete installation"
echo ""
