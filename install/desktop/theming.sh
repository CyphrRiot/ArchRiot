#!/bin/bash

# ==============================================================================
# OhmArchy Desktop Theming System
# ==============================================================================
# Installs and configures desktop theming with flexible theme support
# Handles cursors, icons, GTK themes, and OhmArchy-specific theming
# ==============================================================================

# Load user environment
load_user_environment() {
    local env_file="$HOME/.config/omarchy/user.env"
    [[ -f "$env_file" ]] && source "$env_file"
}

# Install cursor and icon themes
install_cursor_theme() {
    echo "🎯 Installing cursor theme..."

    # Install Bibata Modern Ice cursor theme
    if yay -S --noconfirm --needed bibata-cursor-theme; then
        echo "✓ Bibata cursor theme installed"
    else
        echo "❌ Failed to install bibata-cursor-theme"
        return 1
    fi

    # Verify cursor theme is accessible
    local cursor_locations=(
        "/usr/share/icons/Bibata-Modern-Ice"
        "$HOME/.local/share/icons/Bibata-Modern-Ice"
        "$HOME/.icons/Bibata-Modern-Ice"
    )

    local cursor_found=false
    for location in "${cursor_locations[@]}"; do
        if [[ -d "$location" ]]; then
            cursor_found=true
            echo "✓ Cursor theme found at: $location"
            break
        fi
    done

    if [[ "$cursor_found" != "true" ]]; then
        echo "❌ Bibata-Modern-Ice cursor theme not accessible after installation"
        return 1
    fi

    echo "✓ Cursor theme installation verified"
}

install_icon_theme() {
    echo "🎨 Installing icon theme..."

    if yay -S --noconfirm --needed kora-icon-theme; then
        echo "✓ Kora icon theme installed"
    else
        echo "⚠ Failed to install kora-icon-theme (using fallback)"
        return 1
    fi
}

install_gtk_themes() {
    echo "🖼️  Installing GTK themes..."

    # Install base GTK themes
    if sudo pacman -S --noconfirm gnome-themes-extra; then
        echo "✓ GNOME themes installed"
    else
        echo "⚠ Failed to install GNOME themes"
    fi

    # Install Qt theming support
    if sudo pacman -S --noconfirm kvantum-qt5; then
        echo "✓ Qt theming support installed"
    else
        echo "⚠ Failed to install Qt theming support"
    fi
}

# Configure system-wide theme settings
configure_gtk_settings() {
    echo "⚙️  Configuring GTK settings..."

    if command -v gsettings >/dev/null 2>&1; then
        # Set dark theme preference
        gsettings set org.gnome.desktop.interface gtk-theme "Adwaita-dark" 2>/dev/null || true
        gsettings set org.gnome.desktop.interface color-scheme "prefer-dark" 2>/dev/null || true

        # Set icon theme
        gsettings set org.gnome.desktop.interface icon-theme "kora" 2>/dev/null || true

        # Set cursor theme
        gsettings set org.gnome.desktop.interface cursor-theme "Bibata-Modern-Ice" 2>/dev/null || true
        gsettings set org.gnome.desktop.interface cursor-size 24 2>/dev/null || true

        echo "✓ GTK settings configured"
    else
        echo "⚠ gsettings not available, will use environment variables"
    fi
}

# Setup cursor theme links for applications that don't use gsettings
setup_cursor_links() {
    echo "🔗 Setting up cursor theme links..."

    mkdir -p ~/.icons/default

    local omarchy_cursor_index="$HOME/.local/share/omarchy/default/icons/default/index.theme"
    if [[ -f "$omarchy_cursor_index" ]]; then
        cp "$omarchy_cursor_index" ~/.icons/default/index.theme
        echo "✓ Default cursor theme links created"
    else
        # Create basic index.theme if OhmArchy version not found
        cat > ~/.icons/default/index.theme <<EOF
[Icon Theme]
Name=Default
Comment=Default cursor theme
Inherits=Bibata-Modern-Ice
EOF
        echo "✓ Fallback cursor theme links created"
    fi
}

# Setup OhmArchy theme system
setup_omarchy_theme_system() {
    echo "🎨 Setting up OhmArchy theme system..."

    # Create theme directories
    mkdir -p ~/.config/omarchy/{themes,current,backgrounds}

    # Link available themes
    local themes_source="$HOME/.local/share/omarchy/themes"
    if [[ -d "$themes_source" ]]; then
        for theme_dir in "$themes_source"/*; do
            if [[ -d "$theme_dir" ]]; then
                local theme_name=$(basename "$theme_dir")
                ln -sf "$theme_dir" ~/.config/omarchy/themes/
                echo "✓ Linked theme: $theme_name"
            fi
        done
    else
        echo "⚠ OhmArchy themes directory not found at: $themes_source"
        return 1
    fi

    echo "✓ OhmArchy theme system initialized"
}

# Set default theme (CypherRiot or first available)
set_default_theme() {
    echo "🎯 Setting default theme..."

    local preferred_theme="cypherriot"
    local theme_path="$HOME/.config/omarchy/themes/$preferred_theme"

    # Check if preferred theme exists
    if [[ -d "$theme_path" ]]; then
        echo "✓ Using preferred theme: $preferred_theme"
    else
        # Find first available theme
        local first_theme=$(find ~/.config/omarchy/themes -maxdepth 1 -type d -not -name themes | head -1)
        if [[ -n "$first_theme" ]]; then
            preferred_theme=$(basename "$first_theme")
            theme_path="$first_theme"
            echo "✓ Using fallback theme: $preferred_theme"
        else
            echo "❌ No themes available"
            return 1
        fi
    fi

    # Set current theme link
    ln -snf "$theme_path" ~/.config/omarchy/current/theme
    echo "✓ Active theme set to: $preferred_theme"

    # Setup background system if available
    setup_theme_backgrounds "$preferred_theme"
}

# Setup background system for selected theme
setup_theme_backgrounds() {
    local theme_name="$1"
    echo "🖼️  Setting up backgrounds for theme: $theme_name"

    # Source background script if available
    local bg_script="$HOME/.local/share/omarchy/themes/$theme_name/backgrounds.sh"
    if [[ -f "$bg_script" ]]; then
        source "$bg_script" 2>/dev/null || echo "⚠ Background script failed to execute"
    fi

    # Link background directory
    local bg_dir="$HOME/.config/omarchy/backgrounds/$theme_name"
    if [[ -d "$bg_dir" ]]; then
        ln -snf "$bg_dir" ~/.config/omarchy/current/backgrounds

        # Set default background (first available)
        local first_bg=$(find "$bg_dir" -name "*.jpg" -o -name "*.png" | head -1)
        if [[ -n "$first_bg" ]]; then
            ln -snf "$first_bg" ~/.config/omarchy/current/background
            echo "✓ Default background set: $(basename "$first_bg")"
        fi
    else
        echo "⚠ Background directory not found for theme: $theme_name"
    fi
}

# Link theme configurations to application configs
link_theme_configs() {
    echo "🔗 Linking theme configurations..."

    local theme_dir="$HOME/.config/omarchy/current/theme"
    if [[ ! -d "$theme_dir" ]]; then
        echo "❌ No active theme found"
        return 1
    fi

    # Backup existing configs
    backup_existing_configs

    # Link application configs
    local config_links=(
        "hyprlock.conf:$HOME/.config/hypr/hyprlock.conf"
        "wofi.css:$HOME/.config/wofi/style.css"
        "neovim.lua:$HOME/.config/nvim/lua/plugins/theme.lua"
        "btop.theme:$HOME/.config/btop/themes/current.theme"
        "mako.ini:$HOME/.config/mako/config"
    )

    for link_info in "${config_links[@]}"; do
        IFS=':' read -r source_file target_path <<< "$link_info"
        local source_path="$theme_dir/$source_file"

        if [[ -f "$source_path" ]]; then
            mkdir -p "$(dirname "$target_path")"
            ln -snf "$source_path" "$target_path"
            echo "✓ Linked: $(basename "$source_file")"
        else
            echo "⚠ Theme file not found: $source_file"
        fi
    done

    # Handle waybar config specially (with backup)
    link_waybar_config
}

# Backup existing configurations
backup_existing_configs() {
    echo "💾 Creating configuration backups..."

    local configs_to_backup=(
        "$HOME/.config/waybar/config"
        "$HOME/.config/wofi/style.css"
        "$HOME/.config/hypr/hyprlock.conf"
        "$HOME/.config/mako/config"
    )

    for config in "${configs_to_backup[@]}"; do
        if [[ -f "$config" && ! -L "$config" ]]; then
            cp "$config" "${config}.backup-$(date +%s)" 2>/dev/null || true
        fi
    done

    echo "✓ Configuration backups created"
}

# Link waybar configuration with special handling
link_waybar_config() {
    echo "📊 Configuring waybar theme..."

    local theme_config="$HOME/.config/omarchy/current/theme/config"
    local waybar_config="$HOME/.config/waybar/config"
    local default_config="$HOME/.config/waybar/config.default"

    # Create default backup if it doesn't exist
    if [[ -f "$waybar_config" && ! -f "$default_config" && ! -L "$waybar_config" ]]; then
        cp "$waybar_config" "$default_config"
        echo "✓ Created waybar default backup"
    fi

    # Link theme config or restore default
    if [[ -f "$theme_config" ]]; then
        ln -snf "$theme_config" "$waybar_config"
        echo "✓ Waybar theme config linked"
    elif [[ -f "$default_config" ]]; then
        ln -snf "$default_config" "$waybar_config"
        echo "✓ Waybar default config restored"
    else
        echo "⚠ No waybar config available"
    fi
}

# Start theme-related services safely
start_theme_services() {
    echo "🚀 Starting theme services..."

    # Stop existing services gracefully
    stop_theme_services

    # Start background service
    local bg_file="$HOME/.config/omarchy/current/background"
    if [[ -f "$bg_file" ]]; then
        swaybg -i "$bg_file" -m fill >/dev/null 2>&1 &
        echo "✓ Background service started"
    fi

    # Start waybar if config exists
    if [[ -f "$HOME/.config/waybar/config" ]]; then
        waybar >/dev/null 2>&1 &
        sleep 2
        if pgrep waybar >/dev/null; then
            echo "✓ Waybar started successfully"
        else
            echo "⚠ Waybar failed to start"
        fi
    fi

    # Start blue light filter
    start_blue_light_filter
}

# Stop theme-related services safely
stop_theme_services() {
    echo "🛑 Stopping existing theme services..."

    # Stop waybar
    pkill waybar 2>/dev/null || true
    sleep 1

    # Stop background service
    pkill swaybg 2>/dev/null || true
    sleep 1

    echo "✓ Services stopped"
}

# Setup blue light filtering
start_blue_light_filter() {
    echo "🌙 Setting up blue light filter..."

    # Check if already running
    if pgrep hyprsunset >/dev/null; then
        echo "✓ Blue light filter already running"
        return 0
    fi

    # Start hyprsunset
    if command -v hyprsunset >/dev/null 2>&1; then
        hyprsunset -t 4000 >/dev/null 2>&1 &
        sleep 2

        if pgrep hyprsunset >/dev/null; then
            echo "✓ Blue light filter started (4000K)"
        else
            echo "⚠ Blue light filter failed to start"
        fi
    else
        echo "⚠ hyprsunset not found, blue light filter not available"
    fi
}

# Validate theme installation
validate_theme_setup() {
    echo "🧪 Validating theme setup..."

    local validation_errors=0

    # Check cursor theme
    if [[ -f ~/.icons/default/index.theme ]]; then
        echo "✓ Cursor theme configuration found"
    else
        echo "❌ Cursor theme configuration missing"
        ((validation_errors++))
    fi

    # Check active theme
    if [[ -L ~/.config/omarchy/current/theme ]]; then
        local active_theme=$(basename "$(readlink ~/.config/omarchy/current/theme)")
        echo "✓ Active theme: $active_theme"
    else
        echo "❌ No active theme set"
        ((validation_errors++))
    fi

    # Check background
    if [[ -f ~/.config/omarchy/current/background ]]; then
        echo "✓ Background configured"
    else
        echo "⚠ Background not configured"
    fi

    # Check waybar
    if pgrep waybar >/dev/null; then
        echo "✓ Waybar running"
    else
        echo "⚠ Waybar not running"
    fi

    if [[ $validation_errors -eq 0 ]]; then
        echo "✅ Theme validation passed"
        return 0
    else
        echo "❌ Theme validation failed with $validation_errors errors"
        return 1
    fi
}

# Display theming setup summary
display_theming_summary() {
    echo ""
    echo "🎉 Desktop theming setup complete!"
    echo ""
    echo "🎨 Installed themes:"
    echo "  • Cursor: Bibata Modern Ice"
    echo "  • Icons: Kora"
    echo "  • GTK: Adwaita Dark"
    echo "  • OhmArchy: $(basename "$(readlink ~/.config/omarchy/current/theme 2>/dev/null)" 2>/dev/null || echo "Not set")"
    echo ""
    echo "🎯 Active features:"
    echo "  • Dark mode preference"
    echo "  • Blue light filter (4000K)"
    echo "  • Dynamic wallpapers"
    echo "  • Consistent theming across applications"
    echo ""
    echo "🔧 Theme management:"
    echo "  • Use omarchy-theme-next to cycle themes"
    echo "  • Themes located in ~/.config/omarchy/themes/"
    echo "  • Current theme: ~/.config/omarchy/current/theme"
}

# Main execution with comprehensive error handling
main() {
    echo "🚀 Starting desktop theming setup..."

    load_user_environment

    # Install theme components
    install_cursor_theme || {
        echo "❌ Failed to install cursor theme"
        return 1
    }

    install_icon_theme || {
        echo "⚠ Icon theme installation had issues"
    }

    install_gtk_themes || {
        echo "⚠ GTK theme installation had issues"
    }

    # Configure system settings
    configure_gtk_settings
    setup_cursor_links

    # Setup OhmArchy theme system
    setup_omarchy_theme_system || {
        echo "❌ Failed to setup OhmArchy theme system"
        return 1
    }

    set_default_theme || {
        echo "❌ Failed to set default theme"
        return 1
    }

    # Link configurations
    link_theme_configs || {
        echo "⚠ Some theme configurations failed to link"
    }

    # Start services
    start_theme_services

    # Validate setup
    validate_theme_setup || {
        echo "⚠ Theme validation had issues"
    }

    display_theming_summary
    echo "✅ Desktop theming setup completed!"
}

# Execute main function
main "$@"
