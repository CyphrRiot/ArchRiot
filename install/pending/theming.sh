#!/bin/bash

# ==============================================================================
# ArchRiot Desktop Theming System
# ==============================================================================
# Installs and configures desktop theming with flexible theme support
# Handles cursors, icons, GTK themes, and ArchRiot-specific theming
# ==============================================================================

# Load user environment
load_user_environment() {
    local env_file="$HOME/.config/archriot/user.env"
    if [[ -f "$env_file" ]]; then
        source "$env_file"
    fi
}

# Installation state management
STATE_FILE="$HOME/.config/archriot/.install-state"

init_state_tracking() {
    mkdir -p "$(dirname "$STATE_FILE")"
    touch "$STATE_FILE"
}

is_component_installed() {
    local component="$1"
    grep -q "^${component}=completed$" "$STATE_FILE" 2>/dev/null
}

mark_component_completed() {
    local component="$1"
    init_state_tracking
    # Remove any existing entry for this component
    sed -i "/^${component}=/d" "$STATE_FILE" 2>/dev/null || true
    echo "${component}=completed" >> "$STATE_FILE"
    echo "✓ Marked $component as completed"
}

check_component_status() {
    local component="$1"
    if is_component_installed "$component"; then
        echo "✓ $component already completed, skipping"
        return 0
    else
        echo "🔄 Installing $component..."
        return 1
    fi
}

# Install cursor and icon themes
install_cursor_theme() {
    if check_component_status "cursor_theme"; then
        return 0
    fi

    echo "🎯 Installing cursor theme..."

    # Install Bibata Modern Ice cursor theme
    if timeout 120 yay -S --noconfirm --needed bibata-cursor-theme; then
        echo "✓ Bibata cursor theme installed"
        mark_component_completed "cursor_theme"
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

cleanup_old_icon_themes() {
    echo "🧹 Cleaning up old icon themes..."

    # Remove Obsidian icon themes if installed
    if pacman -Qi obsidian-icon-theme &>/dev/null; then
        echo "🗑️  Removing obsidian-icon-theme..."
        timeout 60 yay -Rs --noconfirm obsidian-icon-theme 2>/dev/null || timeout 60 sudo pacman -Rs --noconfirm obsidian-icon-theme 2>/dev/null || true
    fi



    # Clean up any remaining Obsidian icon directories
    if [ -d "/usr/share/icons/Obsidian" ] || [ -d "/usr/share/icons/Obsidian-Purple" ]; then
        echo "🗑️  Removing leftover Obsidian icon directories..."
        sudo rm -rf /usr/share/icons/Obsidian* 2>/dev/null || true
    fi

    # Remove unwanted desktop entries that clutter Fuzzel
    echo "🧹 Cleaning up unwanted desktop entries..."
    local unwanted_desktop_entries=(
        "/usr/share/applications/lstopo.desktop"  # Hardware Locality utility
    )

    for entry in "${unwanted_desktop_entries[@]}"; do
        if [[ -f "$entry" ]]; then
            echo "🗑️  Removing $(basename "$entry")..."
            sudo rm -f "$entry" 2>/dev/null || true
        fi
    done

    echo "✓ Old icon themes cleaned up"
}

install_icon_theme() {
    if check_component_status "icon_theme"; then
        return 0
    fi

    echo "🎨 Installing icon theme..."

    # Clean up old themes first
    cleanup_old_icon_themes

    # Install Kora icon theme (required by Fuzzel config)
    if timeout 120 yay -S --noconfirm --needed kora-icon-theme; then
        echo "✓ Kora icon theme installed (Fuzzel dependency)"
    else
        echo "⚠ Failed to install kora-icon-theme"
    fi

    if timeout 120 yay -S --noconfirm --needed tela-icon-theme-purple-git; then
        echo "✓ Tela purple icon theme installed"

        # Reset XDG user directories and icons
        xdg-user-dirs-update
        # Only run gtk update if available (not always installed)
        command -v xdg-user-dirs-gtk-update &>/dev/null && xdg-user-dirs-gtk-update || true

        # Force refresh folder icons
        gtk-update-icon-cache -f -t ~/.local/share/icons/ 2>/dev/null || true
        gtk-update-icon-cache -f -t /usr/share/icons/ 2>/dev/null || true

        echo "✓ Icon theme configuration complete"
        mark_component_completed "icon_theme"
    else
        echo "⚠ Failed to install tela-icon-theme-purple-git (using fallback)"
        return 1
    fi
}

install_gtk_themes() {
    if check_component_status "gtk_themes"; then
        return 0
    fi

    echo "🖼️  Installing GTK themes..."

    # Install base GTK themes
    if timeout 120 sudo pacman -S --noconfirm gnome-themes-extra; then
        echo "✓ GNOME themes installed"
    else
        echo "⚠ Failed to install GNOME themes"
    fi

    # Install Qt theming support
    if timeout 120 sudo pacman -S --noconfirm kvantum-qt5; then
        echo "✓ Qt theming support installed"
        mark_component_completed "gtk_themes"
    else
        echo "⚠ Failed to install Qt theming support"
    fi
}

# Configure system-wide theme settings
configure_gtk_settings() {
    echo "⚙️  Configuring GTK settings..."

    if command -v gsettings >/dev/null 2>&1; then
        # Set dark theme preference for GTK3/4 applications
        gsettings set org.gnome.desktop.interface gtk-theme "Adwaita-dark" 2>/dev/null || true

        # Set color scheme for libadwaita applications (prevents GTK4 warnings)
        # This replaces the deprecated gtk-application-prefer-dark-theme setting
        gsettings set org.gnome.desktop.interface color-scheme "prefer-dark" 2>/dev/null || true

        # Set window manager theme to match GTK theme (prevents grey-brown borders)
        gsettings set org.gnome.desktop.wm.preferences theme "Adwaita-dark" 2>/dev/null || true

        # Set icon theme (check multiple possible names)
        if gsettings set org.gnome.desktop.interface icon-theme "Tela-purple-dark" 2>/dev/null; then
            echo "✓ Icon theme set to Tela-purple-dark"
        elif gsettings set org.gnome.desktop.interface icon-theme "Tela-dark" 2>/dev/null; then
            echo "✓ Icon theme set to Tela-dark"
        elif gsettings set org.gnome.desktop.interface icon-theme "Tela" 2>/dev/null; then
            echo "✓ Icon theme set to Tela"
        else
            echo "⚠ Could not set icon theme"
        fi

        # Enable symbolic folder icons
        gsettings set org.gnome.desktop.interface icon-theme-use-symbolic true 2>/dev/null || true

        # Set cursor theme
        gsettings set org.gnome.desktop.interface cursor-theme "Bibata-Modern-Ice" 2>/dev/null || true
        gsettings set org.gnome.desktop.interface cursor-size 24 2>/dev/null || true

        # Fix window button layout (ensure close button is on right side)
        gsettings set org.gnome.desktop.wm.preferences button-layout ":minimize,maximize,close" 2>/dev/null || true

        echo "✓ GTK settings configured"
    else
        echo "⚠ gsettings not available, will use environment variables"
    fi
}

# Setup cursor theme links for applications that don't use gsettings
setup_cursor_links() {
    echo "🔗 Setting up cursor theme links..."

    mkdir -p ~/.icons/default

    local archriot_cursor_index="$HOME/.local/share/archriot/config/default/icons/default/index.theme"
    if [[ -f "$archriot_cursor_index" ]]; then
        cp "$archriot_cursor_index" ~/.icons/default/index.theme
        echo "✓ Default cursor theme links created"
    else
        # Create basic index.theme if ArchRiot version not found
        cat > ~/.icons/default/index.theme <<EOF
[Icon Theme]
Name=Default
Comment=Default cursor theme
Inherits=Bibata-Modern-Ice
EOF
        echo "✓ Fallback cursor theme links created"
    fi
}

# Theme system removed - using consolidated configs

# Fix upgrade path - update hyprlock.conf if still using old theme sourcing
fix_hyprlock_upgrade() {
    echo "🔧 Checking hyprlock configuration for upgrade compatibility..."

    local hyprlock_config="$HOME/.config/hypr/hyprlock.conf"
    if [[ -f "$hyprlock_config" ]] && grep -q "source.*theme.*hyprlock.conf" "$hyprlock_config"; then
        echo "🔄 Updating hyprlock.conf from theme sourcing to consolidated config"

        # Backup old config
        cp "$hyprlock_config" "$hyprlock_config.backup-$(date +%s)"

        # Copy new consolidated config
        local new_config="$HOME/.local/share/archriot/config/hypr/hyprlock.conf"
        if [[ -f "$new_config" ]]; then
            cp "$new_config" "$hyprlock_config"
            echo "✓ Hyprlock configuration updated for v2.0.0+"
        else
            echo "❌ New hyprlock config not found"
            return 1
        fi
    else
        echo "✓ Hyprlock configuration already up to date"
    fi
}

# Setup background system for consolidated backgrounds
# Uses the unified backgrounds directory created during theme consolidation
setup_backgrounds() {
    if check_component_status "backgrounds"; then
        return 0
    fi

    echo "🖼️  Setting up consolidated backgrounds..."

    # Ensure required directories exist
    export BACKGROUNDS_DIR="$HOME/.config/archriot/backgrounds"
    mkdir -p "$BACKGROUNDS_DIR"

    # Copy background files from ArchRiot data directory
    local source_bg_dir="$HOME/.local/share/archriot/backgrounds"

    echo "🔍 Checking source directory: $source_bg_dir"

    if [[ ! -d "$source_bg_dir" ]]; then
        echo "🚨 CRITICAL: Source background directory not found: $source_bg_dir"
        echo "🚨 This will cause background cycling to fail!"
        return 1
    fi

    echo "📦 Installing backgrounds from: $source_bg_dir"

    # Copy backgrounds directly (no numbered prefixes needed)
    local failed_copies=0
    local copied_count=0

    for bg_file in "$source_bg_dir"/*.{jpg,jpeg,png,webp}; do
        [[ -f "$bg_file" ]] || continue

        local filename=$(basename "$bg_file")
        if cp "$bg_file" "$BACKGROUNDS_DIR/$filename"; then
            echo "✓ Installed: $filename"
            ((copied_count++))
        else
            echo "❌ Failed to install: $filename"
            ((failed_copies++))
        fi
    done

    if [[ $copied_count -eq 0 ]]; then
        echo "🚨 CRITICAL: No background files were successfully installed!"
        return 1
    fi

    echo "✓ Installed $copied_count backgrounds"

    mark_component_completed "backgrounds"
    echo "✅ Background setup completed successfully"
    return 0
}

# Link theme configurations to application configs
link_main_configs() {
    echo "🔗 Linking main configurations..."

    # Copy GTK CSS for Thunar dark theme
    if [[ -f "$HOME/.local/share/archriot/config/gtk-3.0/gtk.css" ]]; then
        mkdir -p "$HOME/.config/gtk-3.0"
        cp "$HOME/.local/share/archriot/config/gtk-3.0/gtk.css" "$HOME/.config/gtk-3.0/gtk.css"
        echo "✓ GTK dark theme CSS applied"
    fi

    # Copy GTK settings (includes recent files disabled)
    if [[ -f "$HOME/.local/share/archriot/config/gtk-3.0/settings.ini" ]]; then
        mkdir -p "$HOME/.config/gtk-3.0"
        cp "$HOME/.local/share/archriot/config/gtk-3.0/settings.ini" "$HOME/.config/gtk-3.0/settings.ini"
        echo "✓ GTK-3.0 settings applied (recent files disabled)"
    fi

    if [[ -f "$HOME/.local/share/archriot/config/gtk-4.0/settings.ini" ]]; then
        mkdir -p "$HOME/.config/gtk-4.0"
        cp "$HOME/.local/share/archriot/config/gtk-4.0/settings.ini" "$HOME/.config/gtk-4.0/settings.ini"
        echo "✓ GTK-4.0 settings applied (recent files disabled)"
    fi

    # Copy Fuzzel configuration (includes Kora icon theme)
    if [[ -f "$HOME/.local/share/archriot/config/fuzzel/fuzzel.ini" ]]; then
        mkdir -p "$HOME/.config/fuzzel"
        cp "$HOME/.local/share/archriot/config/fuzzel/fuzzel.ini" "$HOME/.config/fuzzel/fuzzel.ini"
        echo "✓ Fuzzel configuration applied (Kora icons)"
    fi

    echo "✓ Main configurations linked"
}

# Setup fuzzel cache directory to prevent ~/yes file creation
setup_fuzzel_cache() {
    echo "📁 Setting up fuzzel cache directory..."

    local cache_dir="$HOME/.cache/fuzzel"



    # Ensure .cache directory exists first
    if ! [[ -d "$HOME/.cache" ]]; then
        echo "📁 Creating .cache directory first..."
        if ! mkdir -p "$HOME/.cache"; then
            echo "❌ Failed to create .cache directory: $HOME/.cache"
            echo "   Check filesystem permissions and disk space"
            return 1
        fi
    fi

    # Create fuzzel cache directory
    if mkdir -p "$cache_dir" 2>/dev/null; then
        echo "✓ Fuzzel cache directory created: $cache_dir"
    else
        echo "❌ Failed to create fuzzel cache directory: $cache_dir"
        echo "   Error details: $(mkdir -p "$cache_dir" 2>&1 || true)"
        echo "   Continuing installation (fuzzel will still work)"
        return 0  # Don't fail the entire installation
    fi

    # Verify directory is writable
    if [[ -w "$cache_dir" ]]; then
        echo "✓ Fuzzel cache directory is writable"
    else
        echo "❌ Fuzzel cache directory is not writable: $cache_dir"
        echo "   This may cause fuzzel cache issues but won't break functionality"
        return 0  # Don't fail the entire installation
    fi

    echo "✓ Fuzzel cache setup completed"
}

# Backup existing configurations
backup_existing_configs() {
    echo "💾 Creating configuration backups..."

    # Load centralized backup system
    if [[ -f "$INSTALL_DIR/lib/backup-manager.sh" ]]; then
        source "$INSTALL_DIR/lib/backup-manager.sh"

        local configs_to_backup=(
            "waybar"
            "fuzzel"
            "hypr"
            "mako"
            "gtk-3.0"
            "gtk-4.0"
        )

        if backup_configs "theming" "${configs_to_backup[@]}"; then
            echo "✓ Configuration backups created at: ~/.archriot/backups/"
        else
            echo "✓ No configurations to backup"
        fi
    else
        echo "✓ Backup system not available, skipping backup"
    fi
}

# Link waybar configuration with special handling
# Waybar config linking removed - using consolidated config directly

# Start theme-related services safely
start_background_service() {
    echo "🚀 Starting background service..."

    # Find first available background from consolidated directory
    local bg_file=$(find "$HOME/.config/archriot/backgrounds" -name "*riot_01*" | head -1)
    if [[ -z "$bg_file" ]]; then
        bg_file=$(find "$HOME/.config/archriot/backgrounds" -type f \( -name "*.jpg" -o -name "*.png" -o -name "*.jpeg" -o -name "*.webp" \) | sort | head -1)
    fi

    if [[ -f "$bg_file" ]]; then
        nohup swaybg -i "$bg_file" -m fill >/dev/null 2>&1 & disown
        echo "✓ Background service started with: $(basename "$bg_file")"
    else
        echo "⚠ No background files found in ~/.config/archriot/backgrounds"
    fi

    # Start blue light filter by default
    start_blue_light_filter
}

# Stop theme-related services safely
stop_theme_services() {
    # Stop background service
    pkill swaybg 2>/dev/null || true
    sleep 1

    echo "✓ Services stopped"
}





# Setup blue light filtering (handled by hyprland.conf)
start_blue_light_filter() {
    echo "🌙 Blue light filter managed by Hyprland..."

    # Blue light filter is started automatically by hyprland.conf:
    # exec-once = hyprsunset -t 3500
    echo "✓ Blue light filter configured in Hyprland startup"
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

    # Check consolidated backgrounds
    if [[ -d ~/.config/archriot/backgrounds ]] && [[ $(find ~/.config/archriot/backgrounds -type f | wc -l) -gt 0 ]]; then
        echo "✓ Backgrounds configured"
    else
        echo "⚠ Backgrounds not configured"
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
    echo "  • Icons: Tela-purple-dark"
    echo "  • GTK: Adwaita Dark"
    echo "  • ArchRiot: CypherRiot (consolidated)"
    echo ""
    echo "🎯 Active features:"
    echo "  • Dark mode preference"
    echo "  • Blue light filter (4000K)"
    echo "  • Dynamic wallpapers"
    echo "  • Consistent theming across applications"
    echo ""
    echo "🔧 Background management:"
    echo "  • Use swaybg-next to cycle backgrounds"
    echo "  • Backgrounds located in ~/.config/archriot/backgrounds/"
    echo "  • CypherRiot theme integrated into main configs"
}

# Main execution with comprehensive error handling
main() {
    echo "🚀 Starting desktop theming setup..."

    load_user_environment
    init_state_tracking

    # Install theme components
    install_cursor_theme || {
        echo "❌ Failed to install cursor theme"
        set +e
        exit 1
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

    # Fix upgrade path for v2.0.0+ transition
    fix_hyprlock_upgrade || {
        echo "⚠ Hyprlock upgrade had issues"
    }

    # Setup consolidated backgrounds
    setup_backgrounds || {
        echo "❌ Failed to setup backgrounds"
        return 1
    }

    # Link main configurations
    link_main_configs || {
        echo "⚠ Some configuration linking had issues"
    }

    # Setup fuzzel cache directory
    setup_fuzzel_cache || {
        echo "⚠ Fuzzel cache setup had issues"
    }

    # Start services
    start_background_service

    # Validate setup
    validate_theme_setup || {
        echo "⚠ Theme validation had issues"
    }

    display_theming_summary
    echo "✅ Desktop theming setup completed!"
}

# Execute main function
main "$@"
