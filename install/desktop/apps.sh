#!/bin/bash

# ==============================================================================
# OhmArchy Desktop Applications Setup
# ==============================================================================
# Installs essential desktop applications organized by functionality
# with graceful degradation for non-critical components
# ==============================================================================

# Load user environment
load_user_environment() {
    local env_file="$HOME/.config/archriot/user.env"
    [[ -f "$env_file" ]] && source "$env_file"
}

# Install system control utilities (critical for desktop functionality)
install_system_controls() {
    echo "🎛️  Installing system control utilities..."

    local control_packages=(
        "brightnessctl"     # Screen brightness control
        "playerctl"         # Media player control
        "pamixer"          # Audio mixer control
        "pavucontrol"      # Audio control GUI
        "wireplumber"      # Audio session manager
        "xdg-user-dirs"    # XDG folder icons (Downloads, Documents, etc.)
    )

    for package in "${control_packages[@]}"; do
        yay -S --noconfirm --needed "$package" || {
            echo "❌ Failed to install critical control utility: $package"
            return 1
        }
    done

    echo "✓ System control utilities installed"
}

# Install input method support
install_input_methods() {
    echo "⌨️  Installing input method support..."

    local input_packages=(
        "fcitx5"
        "fcitx5-gtk"
        "fcitx5-qt"
        "fcitx5-configtool"
    )

    local failed_input=()
    for package in "${input_packages[@]}"; do
        if yay -S --noconfirm --needed "$package"; then
            echo "✓ Installed: $package"
        else
            failed_input+=("$package")
            echo "⚠ Failed to install: $package"
        fi
    done

    if [[ ${#failed_input[@]} -gt 0 ]]; then
        echo "⚠ Some input method components failed: ${failed_input[*]}"
        echo "  International input may not work properly"
    else
        echo "✓ Input method support installed"
    fi
}

# Install clipboard and session utilities
install_session_utilities() {
    echo "📋 Installing session utilities..."

    local session_packages=(
        "wl-clip-persist"   # Clipboard persistence
        "gnome-keyring"     # Credential storage
        "hyprsunset"        # Blue light filter (theming.sh controls usage)
    )

    for package in "${session_packages[@]}"; do
        if yay -S --noconfirm --needed "$package"; then
            echo "✓ Installed: $package"
        else
            echo "⚠ Failed to install: $package (continuing...)"
        fi
    done

    echo "✓ Session utilities setup complete"
}

# Install file management applications
install_file_management() {
    echo "📁 Installing file management applications..."

    # Critical file manager
    yay -S --noconfirm --needed "thunar" || {
        echo "❌ Failed to install file manager (thunar)"
        return 1
    }

    # File management enhancements (non-critical)
    local file_enhancements=(
        "sushi"                # File previewer
        "ffmpegthumbnailer"   # Video thumbnails
        "tumbler"             # Image/document thumbnails for Thunar
        "libwebp"             # WEBP image format support
        "libheif"             # HEIF/HEIC image format support (iPhone photos)
        "libavif"             # AVIF image format support
        "raw-thumbnailer"     # RAW camera file thumbnails
    )

    for enhancement in "${file_enhancements[@]}"; do
        if yay -S --noconfirm --needed "$enhancement"; then
            echo "✓ Installed: $enhancement"
        else
            echo "⚠ Failed to install: $enhancement (file previews may be limited)"
        fi
    done

    echo "✓ File management applications installed"
}

# Install productivity applications
install_productivity_apps() {
    echo "📊 Installing productivity applications..."

    local productivity_packages=(
        "gnome-calculator"  # Calculator
    )

    for package in "${productivity_packages[@]}"; do
        if yay -S --noconfirm --needed "$package"; then
            echo "✓ Installed: $package"
        else
            echo "⚠ Failed to install: $package"
        fi
    done

    echo "✓ Productivity applications installed"
}

# Install notification and drawer applications
install_ui_enhancements() {
    echo "🔔 Installing UI enhancement applications..."

    local ui_packages=(
        "nwg-drawer"    # Application drawer
        "swaync"        # Notification center
        "fuzzel"        # Application launcher (replaces fuzzel)
    )

    for package in "${ui_packages[@]}"; do
        if yay -S --noconfirm --needed "$package"; then
            echo "✓ Installed: $package"
        else
            echo "⚠ Failed to install: $package (UI features may be limited)"
        fi
    done

    echo "✓ UI enhancement applications installed"
}

# Install media applications
install_media_apps() {
    echo "🎬 Installing media applications..."

    local media_packages=(
        "mpv"       # Video player
        "imv"       # Image viewer
        "evince"    # PDF viewer
    )

    local failed_media=()
    for package in "${media_packages[@]}"; do
        if yay -S --noconfirm --needed "$package"; then
            echo "✓ Installed: $package"
        else
            failed_media+=("$package")
            echo "⚠ Failed to install: $package"
        fi
    done

    if [[ ${#failed_media[@]} -gt 0 ]]; then
        echo "⚠ Some media applications failed: ${failed_media[*]}"
        echo "  Media playback may be limited"
    fi

    echo "✓ Media applications setup complete"
}

# Install web browser (critical for modern desktop)
install_web_browser() {
    echo "🌐 Installing web browser..."

    # Try Brave first (preferred)
    if yay -S --noconfirm --needed "brave-bin"; then
        echo "✓ Brave browser installed"
    else
        echo "❌ Failed to install Brave browser"
        return 1
    fi

    echo "✓ Web browser installation complete"
}

# Install VPN client (optional but recommended)
install_vpn_client() {
    echo "🔒 Installing VPN client..."

    if yay -S --noconfirm --needed "mullvad-vpn-bin"; then
        echo "✓ Mullvad VPN client installed"

        # Check if Mullvad is already configured
        if mullvad status 2>/dev/null | grep -q "Logged in"; then
            echo "✓ Mullvad VPN is already configured and logged in"
            echo ""
        else
            # Always non-interactive for automated installs - no prompting
            echo ""
            echo "📋 Mullvad VPN installed but not activated automatically"
            echo "   To activate later:"
            echo "   1. Create account at https://mullvad.net"
            echo "   2. Run: mullvad account login YOUR_ACCOUNT_NUMBER"
            echo "   3. Run: mullvad auto-connect set on"
            echo ""
        fi
    else
        echo "⚠ Failed to install Mullvad VPN (privacy features limited)"
    fi

    echo "✅ VPN client setup completed!"
}

# Setup XDG user directories for proper folder icons
setup_xdg_directories() {
    echo "📁 Setting up XDG user directories..."

    # Create standard user directories and configure them
    if command -v xdg-user-dirs-update >/dev/null 2>&1; then
        xdg-user-dirs-update
        echo "✓ XDG user directories configured"

        # Verify the configuration was created
        if [[ -f "$HOME/.config/user-dirs.dirs" ]]; then
            echo "✓ XDG configuration file created"
        else
            echo "⚠ XDG configuration file not found"
        fi
    else
        echo "⚠ xdg-user-dirs-update not available"
    fi
}

# Validate critical applications are working
validate_desktop_apps() {
    echo "🧪 Validating desktop applications..."

    local validation_errors=0

    # Check critical applications
    local critical_apps=(
        "thunar:File manager"
        "brightnessctl:Brightness control"
        "playerctl:Media control"
        "pamixer:Audio control"
        "fuzzel:Application launcher"
    )

    for app_info in "${critical_apps[@]}"; do
        IFS=':' read -r app desc <<< "$app_info"
        if command -v "$app" &>/dev/null; then
            echo "✓ $desc ($app) available"
        else
            echo "❌ $desc ($app) missing"
            ((validation_errors++))
        fi
    done

    # Check optional applications
    local optional_apps=(
        "mpv:Video player"
        "imv:Image viewer"
        "brave:Web browser"
        "gnome-calculator:Calculator"
        "nwg-drawer:App drawer"
        "swaync:Notification center"
    )

    local browser_found=false
    for app_info in "${optional_apps[@]}"; do
        IFS=':' read -r app desc <<< "$app_info"
        if command -v "$app" &>/dev/null; then
            echo "✓ $desc ($app) available"
            [[ "$app" == "brave" ]] && browser_found=true
        fi
    done

    # Ensure at least one browser is available
    if [[ "$browser_found" != "true" ]]; then
        echo "❌ No web browser found"
        ((validation_errors++))
    fi

    if [[ $validation_errors -eq 0 ]]; then
        echo "✅ Desktop applications validation passed"
        return 0
    else
        echo "❌ Desktop applications validation failed with $validation_errors critical errors"
        return 1
    fi
}

# Test application functionality
test_application_functionality() {
    echo "🔧 Testing application functionality..."

    # Test file manager
    if command -v thunar &>/dev/null; then
        thunar --help >/dev/null 2>&1 && echo "✓ File manager functional"
    fi

    # Test brightness control
    if command -v brightnessctl &>/dev/null; then
        brightnessctl get >/dev/null 2>&1 && echo "✓ Brightness control functional"
    fi

    # Test audio control
    if command -v pamixer &>/dev/null; then
        pamixer --get-volume >/dev/null 2>&1 && echo "✓ Audio control functional"
    fi

    # Test media control
    if command -v playerctl &>/dev/null; then
        playerctl --help >/dev/null 2>&1 && echo "✓ Media control functional"
    fi

    echo "✓ Application functionality testing complete"
}

# Display applications setup summary
display_apps_summary() {
    echo ""
    echo "🎉 Desktop applications setup complete!"
    echo ""
    echo "📦 Installed categories:"
    echo "  • System controls (brightness, audio, media)"
    echo "  • File management (thunar + enhancements)"
    echo "  • Input methods (fcitx5 for international input)"
    echo "  • Session utilities (clipboard, keyring, blue light filter)"
    echo "  • Media applications (mpv, imv, evince)"
    echo "  • Web browser (Brave)"
    echo "  • VPN client (Mullvad - configure with your account)"
    echo "  • UI enhancements (fuzzel launcher, notification center, app drawer)"
    echo ""
    echo "🚀 Quick access:"
    echo "  • Super+E for file manager"
    echo "  • Super+Shift+S for screenshots"
    echo "  • Media keys for volume/brightness control"
    echo "  • Super+N for notification center"
    echo "  • Waybar VPN indicator (click to connect/disconnect)"
    echo ""
    echo "🔒 Mullvad VPN: To activate, run 'mullvad account login YOUR_ACCOUNT_NUMBER'"
    echo "💡 Tip: All applications are available via fuzzel launcher (Super+D) or app drawer (Super+A)"
}

# Main execution with comprehensive error handling
main() {
    echo "🚀 Starting desktop applications setup..."

    load_user_environment

    # Install applications in logical order
    install_system_controls || {
        echo "❌ Failed to install critical system controls"
        return 1
    }

    install_web_browser || {
        echo "❌ Failed to install web browser"
        return 1
    }

    install_file_management || {
        echo "❌ Failed to install file management"
        return 1
    }

    # Install optional components (can fail without breaking system)
    install_input_methods
    install_session_utilities
    install_productivity_apps
    install_ui_enhancements
    install_media_apps
    install_vpn_client

    # Setup XDG user directories for proper folder icons
    setup_xdg_directories

    # Validate installation
    validate_desktop_apps || {
        echo "❌ Desktop applications validation failed"
        return 1
    }

    # Test functionality
    test_application_functionality

    display_apps_summary
    echo "✅ Desktop applications setup completed!"
}

# Execute main function
main "$@"
