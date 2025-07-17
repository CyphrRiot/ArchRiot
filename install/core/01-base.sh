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
    echo "📦 Installing base development tools (required for yay)..."

    # Critical packages needed for yay installation
    local base_packages="base-devel git rsync bc"

    if command -v install_packages_clean &>/dev/null; then
        install_packages_clean "$base_packages" "Installing base development tools" "BLUE"
    else
        echo "Installing: $base_packages"
        sudo pacman -S --needed --noconfirm $base_packages || {
            echo "❌ CRITICAL: Failed to install base development tools"
            echo "   These are required for yay (AUR helper) installation"
            echo "   Please ensure your system has internet access and pacman is working"
            return 1
        }
        echo "✓ Base development tools installed"
    fi
}

install_yay() {
    # Check if yay is already installed
    if command -v yay &>/dev/null; then
        echo "✓ yay AUR helper already installed"
        return 0
    fi

    echo "📦 Installing yay AUR helper (required for OhmArchy packages)..."

    # Verify prerequisites
    if ! command -v git &>/dev/null; then
        echo "❌ CRITICAL: git not found - required for yay installation"
        return 1
    fi

    if ! command -v makepkg &>/dev/null; then
        echo "❌ CRITICAL: makepkg not found - base-devel package required"
        return 1
    fi

    # Use clean UI if available, otherwise manual installation
    if command -v run_command_clean &>/dev/null; then
        if ! run_command_clean "cd /tmp && git clone https://aur.archlinux.org/yay-bin.git && cd yay-bin && makepkg -si --noconfirm && cd / && rm -rf /tmp/yay-bin" "Installing yay AUR helper" "BLUE"; then
            echo "❌ CRITICAL: yay installation failed (clean method)"
            return 1
        fi
    else
        # Manual installation with better error handling
        local temp_dir="/tmp/yay-install-$$"
        local original_dir="$(pwd)"

        echo "Creating temporary directory: $temp_dir"
        if ! mkdir -p "$temp_dir"; then
            echo "❌ CRITICAL: Cannot create temporary directory for yay installation"
            return 1
        fi

        echo "Cloning yay repository..."
        if ! cd "$temp_dir"; then
            echo "❌ CRITICAL: Cannot enter temporary directory"
            return 1
        fi

        if ! git clone https://aur.archlinux.org/yay-bin.git; then
            echo "❌ CRITICAL: Failed to clone yay repository"
            echo "   Check internet connection and try again"
            cd "$original_dir"
            rm -rf "$temp_dir"
            return 1
        fi

        echo "Building and installing yay..."
        if ! cd yay-bin; then
            echo "❌ CRITICAL: Cannot enter yay directory"
            cd "$original_dir"
            rm -rf "$temp_dir"
            return 1
        fi

        if ! makepkg -si --noconfirm; then
            echo "❌ CRITICAL: yay build/installation failed"
            echo "   This usually means:"
            echo "   1. Missing base-devel package"
            echo "   2. Permission issues"
            echo "   3. Network connectivity problems"
            cd "$original_dir"
            rm -rf "$temp_dir"
            return 1
        fi

        # Cleanup
        cd "$original_dir"
        rm -rf "$temp_dir"

        echo "✓ yay installed successfully"
    fi

    # Refresh PATH to make yay immediately available
    export PATH="/usr/bin:$PATH"
    hash -r 2>/dev/null || true

    # Final verification
    if command -v yay &>/dev/null; then
        echo "✓ yay AUR helper installation verified"
        yay_version=$(yay --version | head -1)
        echo "  $yay_version"
        return 0
    else
        echo "❌ CRITICAL: yay installation completed but command not found"
        echo "   Trying to locate yay binary..."
        if [ -f "/usr/bin/yay" ]; then
            echo "   Found at /usr/bin/yay, adding to PATH"
            export PATH="/usr/bin:$PATH"
            hash -r 2>/dev/null || true
            if command -v yay &>/dev/null; then
                echo "✓ yay now available after PATH refresh"
                return 0
            fi
        fi
        echo "   You may need to restart your shell or check PATH manually"
        return 1
    fi
}

# Main execution with proper error handling
echo "🚀 Starting base system setup..."

if ! fix_mirrors; then
    echo "❌ CRITICAL: Mirror configuration failed"
    exit 1
fi

if ! install_base_devel; then
    echo "❌ CRITICAL: Base development tools installation failed"
    echo "   Cannot proceed without these packages"
    exit 1
fi

if ! install_yay; then
    echo "❌ CRITICAL: yay AUR helper installation failed"
    echo "   OhmArchy requires yay to install AUR packages"
    echo "   Please install yay manually and try again"
    exit 1
fi

echo "✅ Base system setup completed successfully"
echo "   - Mirrors configured"
echo "   - Base development tools installed"
echo "   - yay AUR helper installed and verified"
