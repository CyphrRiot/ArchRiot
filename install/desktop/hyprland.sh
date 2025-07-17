#!/bin/bash

# Install Hyprland configurations first
install_hyprland_configs() {
    echo "📁 Installing Hyprland configurations..."

    # Install OhmArchy configs (SAFELY - preserve user configs)
    local source_config="$HOME/.local/share/archriot/config"
    [[ -d "$source_config" ]] || return 1
    mkdir -p ~/.config

    # Safe config installation - focus on hypr configs needed for desktop setup
    for item in "$source_config"/*; do
        local basename=$(basename "$item")
        local target="$HOME/.config/$basename"

        # Install hypr configs if they don't exist or are OhmArchy-managed
        if [[ "$basename" == "hypr" ]]; then
            if [[ ! -e "$target" ]]; then
                # New installation - safe to copy
                cp -R "$item" "$target" || return 1
                echo "✓ Installed new config: $basename"
            elif [[ -L "$target" ]] && [[ "$(readlink "$target")" == *"archriot"* ]]; then
                # OhmArchy-managed symlink - safe to update
                rm -f "$target"
                cp -R "$item" "$target" || return 1
                echo "✓ Updated OhmArchy config: $basename"
            else
                # USER'S EXISTING CONFIG - create reference copy
                echo "⚠ Preserving existing user config: $basename"
                cp -R "$item" "$target.archriot-default" 2>/dev/null || true
                echo "  → Created reference copy: $basename.archriot-default"
            fi
        fi
    done

    echo "✓ Hyprland configurations installed"
}

# Load environment and install Hyprland packages
setup_hyprland_packages() {
    echo "🪟 Installing Hyprland desktop environment..."

    local env_file="$HOME/.config/archriot/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Core packages (critical)
    local core="hyprland xdg-desktop-portal-hyprland xdg-desktop-portal-gtk polkit-gnome"
    yay -S --noconfirm --needed $core || return 1

    # Utilities (best effort)
    local utilities="waybar fuzzel mako swaybg hyprlock hypridle swayosd grim slurp hyprshot hyprpicker hyprland-qtutils kooha gst-libav"
    yay -S --noconfirm --needed $utilities || echo "⚠ Some Hyprland utilities may have failed"

    # Python dependencies
    yay -S --noconfirm --needed python python-psutil || return 1
    python3 -c "import psutil" || return 1

    echo "✓ Hyprland packages installed"
}

# Validate critical components
validate_installation() {
    echo "🧪 Validating Hyprland installation..."

    local issues=0

    # Check critical binaries
    for cmd in Hyprland waybar hyprlock fuzzel; do
        command -v "$cmd" &>/dev/null || ((issues++))
    done

    if [[ $issues -eq 0 ]]; then
        echo "✓ Critical components validated"
        return 0
    else
        echo "❌ $issues critical components missing"
        return 1
    fi
}

# Setup autostart and verify config
configure_hyprland() {
    echo "🚀 Configuring Hyprland validation..."

    # Verify essential configs exist
    local hyprland_conf="$HOME/.config/hypr/hyprland.conf"
    if [[ -f "$hyprland_conf" ]] && grep -q "bind" "$hyprland_conf"; then
        echo "✓ Hyprland configuration validated"
    else
        echo "⚠ Hyprland configuration may need attention"
    fi

    echo "✓ Hyprland autostart is handled by shell configs (bash/fish)"
    echo "✓ Hyprland config setup complete"
}

# Display setup summary
show_summary() {
    echo ""
    echo "🎉 Hyprland desktop environment setup complete!"
    echo ""
    echo "📦 Installed: Hyprland WM, Waybar, Fuzzel launcher, Mako notifications, SwayOSD volume overlay, Hyprlock screen locker, screen tools, Kooha recorder"
    echo "🚀 Getting started: Log out/in to start Hyprland, or type 'Hyprland'"
    echo "⌨️  Key bindings: Super+Return (terminal), Super+D (launcher)"
}

# Fix touchpad configuration based on hardware detection
fix_touchpad_config() {
    echo "🔧 Configuring touchpad settings..."

    local touchpad_script="$HOME/.local/share/archriot/bin/fix-touchpad-config"
    if [[ -f "$touchpad_script" ]]; then
        chmod +x "$touchpad_script"
        "$touchpad_script"
    else
        echo "⚠ Touchpad config script not found - skipping"
    fi
}

# Main execution
main() {
    echo "🚀 Starting Hyprland desktop environment setup..."

    install_hyprland_configs || return 1
    setup_hyprland_packages || return 1
    validate_installation || return 1
    configure_hyprland
    fix_touchpad_config
    show_summary

    echo "✅ Hyprland desktop environment setup completed!"
}

main "$@"
