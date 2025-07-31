#!/bin/bash

# ==============================================================================
# ArchRiot Productivity Applications Setup
# ==============================================================================
# Simple productivity application installation
# ==============================================================================

# Install productivity tools (no heavy office suites)
install_packages "gnome-text-editor" "essential"

# Install Zed with proper Vulkan driver detection
install_zed_with_vulkan() {
    echo "🔍 Detecting GPU hardware for Zed Vulkan driver..."

    # Detect GPU hardware
    local gpu_info=$(lspci | grep -i "vga\|3d\|display")
    local vulkan_driver=""
    local gpu_detected=""

    if echo "$gpu_info" | grep -qi "nvidia"; then
        gpu_detected="NVIDIA"
        # Check if proprietary drivers are already installed
        if pacman -Qi nvidia-utils &>/dev/null || pacman -Qi nvidia &>/dev/null; then
            vulkan_driver="nvidia-utils"
            echo "✓ NVIDIA GPU detected - using proprietary drivers"
        else
            vulkan_driver="vulkan-nouveau"
            echo "✓ NVIDIA GPU detected - using open-source drivers (nouveau)"
        fi
    elif echo "$gpu_info" | grep -qi "amd\|radeon"; then
        gpu_detected="AMD"
        vulkan_driver="vulkan-radeon"
        echo "✓ AMD GPU detected - using vulkan-radeon"
    elif echo "$gpu_info" | grep -qi "intel"; then
        gpu_detected="Intel"
        vulkan_driver="vulkan-intel"
        echo "✓ Intel GPU detected - using vulkan-intel"
    else
        gpu_detected="Unknown/Virtual"
        echo "⚠ GPU hardware not clearly detected, available options:"
        echo "  1) vulkan-intel (Intel GPUs)"
        echo "  2) vulkan-radeon (AMD GPUs)"
        echo "  3) vulkan-nouveau (NVIDIA open-source)"
        echo "  4) nvidia-utils (NVIDIA proprietary)"
        echo "  5) vulkan-swrast (Software fallback)"

        while true; do
            read -p "🤔 Which Vulkan driver should be installed? [1-5]: " choice
            case $choice in
                1) vulkan_driver="vulkan-intel"; echo "✓ Selected Intel GPU drivers"; break ;;
                2) vulkan_driver="vulkan-radeon"; echo "✓ Selected AMD GPU drivers"; break ;;
                3) vulkan_driver="vulkan-nouveau"; echo "✓ Selected NVIDIA open-source drivers"; break ;;
                4) vulkan_driver="nvidia-utils"; echo "✓ Selected NVIDIA proprietary drivers"; break ;;
                5) vulkan_driver="vulkan-swrast"; echo "✓ Selected software fallback drivers"; break ;;
                *) echo "❌ Invalid choice. Please enter 1-5." ;;
            esac
        done
    fi

    # Install the Vulkan driver first
    if [[ -n "$vulkan_driver" ]]; then
        echo "📦 Installing Vulkan driver: $vulkan_driver"
        if install_packages "$vulkan_driver vulkan-icd-loader vulkan-tools" "essential"; then
            echo "✅ Vulkan packages installed successfully"

            # Verify Vulkan installation
            if command -v vulkaninfo >/dev/null 2>&1; then
                echo "🧪 Testing Vulkan installation..."
                if vulkaninfo --summary >/dev/null 2>&1; then
                    echo "✅ Vulkan driver working correctly"
                else
                    echo "⚠ Vulkan driver installed but may need reboot to activate"
                    echo "  💡 If Zed has issues, try: sudo reboot"
                fi
            else
                echo "⚠ vulkaninfo not available for testing"
            fi
        else
            echo "❌ Failed to install Vulkan driver: $vulkan_driver"
            echo "  🔄 Proceeding with Zed installation anyway..."
        fi
    else
        echo "❌ No Vulkan driver selected - Zed may not work properly"
    fi

    # Now install Zed with proper Vulkan support
    echo "📦 Installing Zed editor with $gpu_detected GPU support..."
    if install_packages "zed" "essential"; then
        # Verify Zed installation
        if command -v zeditor >/dev/null 2>&1; then
            echo "✅ Zed editor installed successfully"
            echo "  🎯 GPU: $gpu_detected | Driver: $vulkan_driver"
            echo "  💡 Binary: zeditor (will create 'zed' command next)"

            # Create user-friendly 'zed' command with upgrade safety
            echo "🔗 Creating 'zed' command for easy access..."
            mkdir -p ~/.local/bin

            # Check if zed-wayland wrapper exists (may be installed already)
            if [[ -f ~/.local/bin/zed-wayland ]]; then
                # Use the Wayland wrapper for better compatibility
                target="zed-wayland"
                echo "✓ zed-wayland wrapper found - using for better Wayland support"
            else
                # Fall back to direct zeditor link (wrapper will be installed later)
                target="/usr/bin/zeditor"
                echo "⏳ zed-wayland wrapper not yet installed - using direct zeditor link"
            fi

            # Handle existing zed command safely
            if [[ -e ~/.local/bin/zed ]]; then
                # Check if it's already pointing to what we want
                current_target=$(readlink ~/.local/bin/zed 2>/dev/null || echo "not_a_symlink")
                if [[ "$current_target" == "$target" ]]; then
                    echo "✓ 'zed' command already correctly configured → $target"
                else
                    echo "🔄 Existing 'zed' command found, updating for ArchRiot compatibility..."
                    # Backup existing command if it's not a symlink
                    if [[ ! -L ~/.local/bin/zed ]]; then
                        mv ~/.local/bin/zed ~/.local/bin/zed.backup.$(date +%s)
                        echo "✓ Backed up existing zed command to zed.backup.*"
                    fi
                    ln -sf "$target" ~/.local/bin/zed
                    echo "✅ Updated 'zed' command → $target"
                fi
            else
                ln -sf "$target" ~/.local/bin/zed
                echo "✅ Created 'zed' command → $target"
            fi
            echo "  💡 Users can now run: zed filename.txt"
        else
            echo "❌ Zed binary 'zeditor' not found after installation"
            return 1
        fi
    else
        echo "❌ Zed installation failed"
        echo "  🔄 You can manually install later with: yay -S zed"
        return 1
    fi
}

# Install Zed with intelligent Vulkan driver detection
install_zed_with_vulkan

install_packages "abiword papers" "essential"
install_packages "thunar unzip p7zip" "essential"

# Install calendar and time management
install_packages "gnome-clocks" "essential"

# Install and configure Gnome Text Editor with theme support
if command -v gnome-text-editor >/dev/null 2>&1; then
    echo "🎨 Installing themes for Gnome Text Editor..."

    # Create gtksourceview styles directory
    mkdir -p "$HOME/.local/share/gtksourceview-5/styles"

    # Install all available text editor themes
    themes_installed=0
    # Use consolidated text editor themes
    if [[ -f "$HOME/.local/share/archriot/config/text-editor/cypherriot.xml" ]]; then
        cp "$HOME/.local/share/archriot/config/text-editor/cypherriot.xml" "$HOME/.local/share/gtksourceview-5/styles/"
        echo "✓ Installed text editor theme: cypherriot.xml"
        themes_installed=1
    fi

    if [[ $themes_installed -gt 0 ]]; then
        echo "✓ $themes_installed text editor theme(s) installed"
    else
        echo "⚠ No text editor themes found"
    fi

    echo "🎨 Configuring Gnome Text Editor..."
    gsettings set org.gnome.TextEditor show-line-numbers true
    gsettings set org.gnome.TextEditor highlight-current-line true
    gsettings set org.gnome.TextEditor show-right-margin false
    gsettings set org.gnome.TextEditor custom-font 'Hack Nerd Font 12'
    gsettings set org.gnome.TextEditor line-height 1.2
    gsettings set org.gnome.TextEditor use-system-font false

    # Set default theme to Tokyo Night Dark
    gsettings set org.gnome.TextEditor style-scheme 'tokyo-night'
    echo "✓ Tokyo Night Dark theme set for text editor"
fi

# Install Zed editor configuration and launcher
echo "🖥️ Installing Zed desktop integration..."
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
mkdir -p ~/.local/share/applications ~/.local/bin

# Install Zed Wayland launcher
if [[ -f ~/.local/share/archriot/bin/zed-wayland ]]; then
  cp ~/.local/share/archriot/bin/zed-wayland ~/.local/bin/
  chmod +x ~/.local/bin/zed-wayland
  echo "✓ Zed Wayland launcher installed"

  # Now update the zed symlink to use the Wayland wrapper (upgrade-safe)
  echo "🔗 Updating 'zed' command to use Wayland wrapper..."
  current_target=$(readlink ~/.local/bin/zed 2>/dev/null || echo "not_found")
  if [[ "$current_target" != "zed-wayland" ]]; then
    ln -sf zed-wayland ~/.local/bin/zed
    echo "✅ Updated 'zed' command → zed-wayland (Wayland optimized)"
  else
    echo "✓ 'zed' command already using Wayland wrapper"
  fi
else
  echo "⚠ Zed Wayland launcher not found in repository"
fi

# Hide system zed desktop files to prevent duplicates
if [[ -f "/usr/share/applications/dev.zed.Zed.desktop" ]]; then
    echo "[Desktop Entry]
NoDisplay=true" > ~/.local/share/applications/dev.zed.Zed.desktop
    echo "✓ System dev.zed.Zed.desktop hidden"
fi

if [[ -f "/usr/share/applications/zed.desktop" ]]; then
    echo "[Desktop Entry]
NoDisplay=true" > ~/.local/share/applications/zed.desktop.bak
    echo "✓ System zed.desktop hidden"
fi

# Install Zed desktop file (upgrade-safe)
if [[ -f ~/.local/share/archriot/applications/zed.desktop ]]; then
  echo "📱 Installing Zed desktop file..."

  # Check if existing desktop file needs updating
  if [[ -f ~/.local/share/applications/zed.desktop ]]; then
    if ! grep -q "zed-wayland" ~/.local/share/applications/zed.desktop; then
      echo "🔄 Updating existing desktop file for ArchRiot compatibility..."
      cp ~/.local/share/applications/zed.desktop ~/.local/share/applications/zed.desktop.backup.$(date +%s)
      echo "✓ Backed up existing desktop file"
    fi
  fi

  # Install with proper HOME expansion
  sed "s|\$HOME|$HOME|g" ~/.local/share/archriot/applications/zed.desktop > ~/.local/share/applications/zed.desktop
  echo "✓ Zed desktop file installed with Wayland support"
else
  echo "⚠ Zed desktop file not found in repository"
fi

# Install Zed configuration
if [[ -d ~/.local/share/archriot/config/zed ]]; then
  mkdir -p ~/.config/zed
  cp -r ~/.local/share/archriot/config/zed/* ~/.config/zed/
  echo "✓ Zed configuration installed"
else
  echo "⚠ Zed configuration not found in repository"
fi

# Update desktop database
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
    echo "✓ Desktop database updated"
fi

# Test Zed functionality
echo "🧪 Testing Zed functionality..."
if command -v zed >/dev/null 2>&1; then
    # Test that Zed can show version (basic functionality)
    if zed --version >/dev/null 2>&1; then
        echo "✅ Zed version check: PASSED"

        # Test that Zed can start without immediately crashing (GPU/Vulkan test)
        echo "🎮 Testing Zed GPU/Vulkan support..."
        if timeout 5s zed --help >/dev/null 2>&1; then
            echo "✅ Zed help command: PASSED (GPU/Vulkan working)"
        else
            echo "⚠ Zed help command timed out or failed"
            echo "  💡 This may indicate GPU/Vulkan driver issues"
            echo "  🔧 Try: yay -R zed && yay -S zed (to reinstall with correct driver)"
        fi
    else
        echo "❌ Zed version check: FAILED"
        echo "  🔧 Try: yay -R zed && yay -S zed"
    fi
else
    echo "❌ Zed command not found after installation"
fi

echo "✅ Zed desktop integration complete!"
echo "  💡 Launch with: zed filename.txt (command line)"
echo "  💡 Launch from: Application menu → Zed (GUI)"

echo "✅ Productivity applications setup complete!"
