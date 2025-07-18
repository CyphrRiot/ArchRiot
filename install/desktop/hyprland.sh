#!/bin/bash

# Install Hyprland configurations first
install_hyprland_configs() {
    echo "📁 Installing Hyprland configurations..."

    # Install OhmArchy configs (SAFELY - preserve user configs)
    local source_config="$HOME/.local/share/archriot/config"
    [[ -d "$source_config" ]] || return 1
    mkdir -p ~/.config

    # Install hypr configs - REPLACE existing configs (this is what users want!)
    for item in "$source_config"/*; do
        local basename=$(basename "$item")
        local target="$HOME/.config/$basename"

        # Focus on hypr configs needed for desktop setup
        if [[ "$basename" == "hypr" ]]; then
            if [[ -e "$target" ]]; then
                # Back up existing config before replacing
                echo "📦 Backing up existing config: $basename"
                cp -R "$target" "$target.backup-$(date +%Y%m%d-%H%M%S)" 2>/dev/null || true
                rm -rf "$target"
            fi
            # Install fresh ArchRiot config
            cp -R "$item" "$target" || return 1
            echo "✓ Installed ArchRiot config: $basename"
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

# Touchpad configuration is now handled in the base hyprland.conf
# No post-processing needed to avoid config corruption

# Main execution
main() {
    echo "🚀 Starting Hyprland desktop environment setup..."

    install_hyprland_configs || return 1
    setup_hyprland_packages || return 1
    validate_installation || return 1
    configure_hyprland

    # Reload Hyprland config if it's running
    if pgrep -x "Hyprland" >/dev/null; then
        echo "🔄 Reloading Hyprland configuration..."
        if hyprctl reload 2>/dev/null; then
            echo "✓ Hyprland configuration reloaded successfully"
        else
            echo "⚠ Failed to reload Hyprland - please restart Hyprland manually"
        fi
    else
        echo "ℹ Hyprland not running - configuration will be applied on next start"
    fi

    show_summary

    echo "✅ Hyprland desktop environment setup completed!"
}

main "$@"
