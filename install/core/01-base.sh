#!/bin/bash

# Fix mirrors if needed
fix_mirrors() {
    echo "🔧 Checking mirrors..."
    if ! sudo pacman -Sy 2>/dev/null; then
        echo "⚠ Mirror issues detected, fixing..."
        sudo pacman -S --noconfirm reflector 2>/dev/null || true
        if command -v reflector &>/dev/null; then
            sudo reflector --country US --age 12 --protocol https --sort rate --fastest 10 --save /etc/pacman.d/mirrorlist
            sudo pacman -Syy
            echo "✓ Mirrors fixed"
        else
            echo "⚠ Could not fix mirrors automatically"
        fi
    fi
}

# Install base development tools and yay AUR helper
install_base_devel() {
    if command -v install_packages_clean &>/dev/null; then
        install_packages_clean "base-devel rsync bc" "Installing base development tools" "BLUE"
    else
        echo "📦 Installing base development tools..."
        sudo pacman -S --needed --noconfirm base-devel rsync bc || return 1
    fi
}

install_yay() {
    command -v yay &>/dev/null && return 0

    if command -v run_command_clean &>/dev/null; then
        run_command_clean "cd /tmp && git clone https://aur.archlinux.org/yay-bin.git && cd yay-bin && makepkg -si --noconfirm && cd / && rm -rf /tmp/yay-bin" "Installing yay AUR helper" "BLUE"
    else
        echo "📦 Installing yay AUR helper..."
        local temp_dir="/tmp/yay-install-$$"

        mkdir -p "$temp_dir" || return 1

        cd "$temp_dir" &&
        git clone https://aur.archlinux.org/yay-bin.git &&
        cd yay-bin &&
        makepkg -si --noconfirm &&
        cd / &&
        rm -rf "$temp_dir" || {
            echo "Error: yay installation failed"
            return 1
        }

        echo "✓ yay installed successfully"
    fi
}

# Main execution
fix_mirrors && install_base_devel && install_yay
