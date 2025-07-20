#!/bin/bash

# Install Hyprland configurations first
install_hyprland_configs() {
    echo "📁 Installing Hyprland configurations..."

    # Install ArchRiot configs (SAFELY - preserve user configs)
    local source_config="$HOME/.local/share/archriot/config"
    [[ -d "$source_config" ]] || return 1
    mkdir -p ~/.config

    # PRESERVE USER MONITORS.CONF BEFORE NUKING CONFIGS
    local monitors_conf="$HOME/.config/hypr/monitors.conf"
    local monitors_backup=""
    if [[ -f "$monitors_conf" ]]; then
        # Check if user has customized monitors.conf
        if grep -q "# User customized" "$monitors_conf" ||
           ! grep -q "# ArchRiot auto-generated" "$monitors_conf" ||
           grep -q -v "^#\|^$\|env = GDK_SCALE\|monitor=,preferred,auto" "$monitors_conf"; then
            monitors_backup="/tmp/archriot_monitors_backup.conf"
            cp "$monitors_conf" "$monitors_backup"
            echo "🖥️ Backing up user-customized monitors.conf"
        fi
    fi

    # Install hypr configs - REPLACE existing configs (this is what users want!)
    for item in "$source_config"/*; do
        local basename=$(basename "$item")
        local target="$HOME/.config/$basename"

        # Focus on hypr configs needed for desktop setup
        if [[ "$basename" == "hypr" ]]; then
            if [[ -e "$target" ]]; then
                # Remove existing config (backup already created by consolidated system)
                echo "📦 Replacing existing config: $basename"
                rm -rf "$target"
            fi
            # Install fresh ArchRiot config
            cp -R "$item" "$target" || return 1
            echo "✓ Installed ArchRiot config: $basename"
        fi
    done

    # RESTORE USER MONITORS.CONF IF IT WAS BACKED UP
    if [[ -n "$monitors_backup" && -f "$monitors_backup" ]]; then
        cp "$monitors_backup" "$monitors_conf"
        rm "$monitors_backup"
        echo "✓ Restored user-customized monitors.conf"
    fi

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
    local utilities="waybar fuzzel mako swaybg hyprlock hypridle swayosd grim slurp hyprshot hyprpicker hyprland-qtutils kooha gst-libav gst-plugins-ugly x264"
    yay -S --noconfirm --needed $utilities || echo "⚠ Some Hyprland utilities may have failed"

    # Display management GUI
    yay -S --noconfirm --needed nwg-displays || echo "⚠ nwg-displays installation failed"

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

# Restart waybar with new configuration if in graphical session
restart_waybar() {
    echo "🔄 Restarting waybar with new configuration..."

    # Kill existing waybar processes more thoroughly
    pkill -f waybar 2>/dev/null || true
    killall waybar 2>/dev/null || true
    sleep 2

    # Ensure waybar is completely stopped
    local attempts=0
    while pgrep -f waybar >/dev/null && [ $attempts -lt 5 ]; do
        echo "⏳ Waiting for waybar to stop..."
        pkill -9 waybar 2>/dev/null || true
        sleep 1
        ((attempts++))
    done

    # Start waybar if we're in a graphical session
    if [[ -n "$WAYLAND_DISPLAY" ]] || [[ -n "$DISPLAY" ]]; then
        echo "🚀 Starting waybar with new configuration..."
        # Use nohup to properly detach from installation process
        nohup waybar </dev/null >/dev/null 2>&1 &

        # Wait longer for waybar to initialize
        sleep 3

        # Verify waybar started successfully
        if pgrep -f waybar >/dev/null; then
            echo "✓ Waybar is running successfully"
        else
            echo "⚠ Waybar failed to start - try: waybar &"
            echo "🔍 Check waybar config: waybar --log-level debug"
        fi
    else
        echo "ℹ No graphical session detected - waybar will start on next login"
    fi
}

# Touchpad configuration is now handled in the base hyprland.conf
# Setup VM-specific scaling configuration
setup_vm_scaling() {
    echo "🖥️ Configuring display scaling for current environment..."

    local monitors_conf="$HOME/.config/hypr/monitors.conf"

    # Check if user has already customized monitors.conf (should have been preserved above)
    if [[ -f "$monitors_conf" ]]; then
        # Check if it contains user customizations (not just default content)
        if grep -q "# User customized" "$monitors_conf" ||
           ! grep -q "# ArchRiot auto-generated" "$monitors_conf" ||
           grep -q -v "^#\|^$\|env = GDK_SCALE\|monitor=,preferred,auto" "$monitors_conf"; then
            echo "🖥️ User-customized monitors.conf preserved - skipping auto-scaling"
            echo "ℹ️ To reset auto-scaling: rm ~/.config/hypr/monitors.conf && re-run installer"
            return 0
        fi
    fi

    # Detect if running in a virtual machine
    local virt_type="none"
    if command -v systemd-detect-virt >/dev/null 2>&1; then
        if systemd-detect-virt >/dev/null 2>&1; then
            virt_type=$(systemd-detect-virt)
        else
            virt_type="none"
        fi
    fi

    # Generate appropriate scaling based on environment
    if [[ "$virt_type" != "none" ]]; then
        echo "🖥️ Virtual machine detected ($virt_type) - applying VM-optimized scaling"
        cat > "$monitors_conf" << 'EOF'
# See https://wiki.hyprland.org/Configuring/Monitors/

# VM-optimized scaling for better visibility
# Change to 1 if display appears too large, or 1.5 for even larger scaling
env = GDK_SCALE,1.25

# Use single default monitor with VM-optimized scaling
# Format: monitor = [port], resolution, position, scale
monitor=,preferred,auto,1.25

# Alternative scaling options for VMs:
# monitor=,preferred,auto,1.0    # Normal scaling (may be small)
# monitor=,preferred,auto,1.5    # Larger scaling for high-DPI VM displays

# Example for specific VM resolutions:
# monitor=,1920x1080@60.00, auto, 1.25
# monitor=,2560x1440@60.00, auto, 1.5

# ArchRiot auto-generated scaling config
# To prevent overwriting: add "# User customized" anywhere in this file
# or modify any monitor/scaling settings below
EOF
        echo "✓ VM-optimized display scaling configured (1.25x)"
        echo "ℹ️ To customize scaling: edit ~/.config/hypr/monitors.conf"
    else
        echo "🖥️ Physical hardware detected - applying standard scaling"
        cat > "$monitors_conf" << 'EOF'
# See https://wiki.hyprland.org/Configuring/Monitors/

# Physical hardware scaling
# Change to 1.25 or 1.5 for fractional scaling on high-DPI displays
env = GDK_SCALE,1

# Use single default monitor (see all monitors with: hyprctl monitors)
# Format: monitor = [port], resolution, position, scale
monitor=,preferred,auto,1.0

# Example for fractional scaling on high-DPI displays:
# env = GDK_SCALE,1.75
# monitor=,preferred,auto,1.666667

# Example multi-monitor setup:
# monitor = DP-5, 6016x3384@60.00, auto, 2
# monitor = eDP-1, 2880x1920@120.00, auto, 2

# ArchRiot auto-generated scaling config
# To prevent overwriting: add "# User customized" anywhere in this file
# or modify any monitor/scaling settings below
EOF
        echo "✓ Standard display scaling configured (1.0x)"
        echo "ℹ️ To customize scaling: edit ~/.config/hypr/monitors.conf"
    fi

    echo "✓ Display scaling configuration complete"
}

# No post-processing needed to avoid config corruption

# Main execution
main() {
    echo "🚀 Starting Hyprland desktop environment setup..."

    install_hyprland_configs || return 1
    setup_hyprland_packages || return 1
    setup_vm_scaling || return 1
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

    # Restart waybar after configuration changes
    restart_waybar

    show_summary

    echo "✅ Hyprland desktop environment setup completed!"
}

main "$@"
