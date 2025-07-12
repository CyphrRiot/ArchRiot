#!/bin/bash

# Install base development tools and yay AUR helper
install_base_devel() {
    echo "📦 Installing base development tools..."
    sudo pacman -S --needed --noconfirm base-devel || return 1
}

install_yay() {
    command -v yay &>/dev/null && return 0

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
}

# Main execution
install_base_devel && install_yay
