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

cleanup_old_icon_themes() {
    echo "🧹 Cleaning up old icon themes..."

    # Remove Obsidian icon themes if installed
    if pacman -Qi obsidian-icon-theme &>/dev/null; then
        echo "🗑️  Removing obsidian-icon-theme..."
        yay -Rs --noconfirm obsidian-icon-theme 2>/dev/null || sudo pacman -Rs --noconfirm obsidian-icon-theme 2>/dev/null || true
    fi

    # Clean up any remaining Obsidian icon directories
    if [ -d "/usr/share/icons/Obsidian" ] || [ -d "/usr/share/icons/Obsidian-Purple" ]; then
        echo "🗑️  Removing leftover Obsidian icon directories..."
        sudo rm -rf /usr/share/icons/Obsidian* 2>/dev/null || true
    fi

    echo "✓ Old icon themes cleaned up"
}

install_icon_theme() {
    echo "🎨 Installing icon theme..."

    # Clean up old themes first
    cleanup_old_icon_themes

    # Install Kora icon theme (required by Fuzzel config)
    if yay -S --noconfirm --needed kora-icon-theme; then
        echo "✓ Kora icon theme installed (Fuzzel dependency)"
    else
        echo "⚠ Failed to install kora-icon-theme"
    fi

    if yay -S --noconfirm --needed tela-icon-theme-purple-git; then
        echo "✓ Tela purple icon theme installed"

        # Reset XDG user directories and icons
        xdg-user-dirs-update
        xdg-user-dirs-gtk-update

        # Force refresh folder icons
        gtk-update-icon-cache -f -t ~/.local/share/icons/ 2>/dev/null || true
        gtk-update-icon-cache -f -t /usr/share/icons/ 2>/dev/null || true

        echo "✓ Icon theme configuration complete"
    else
        echo "⚠ Failed to install tela-icon-theme-purple-git (using fallback)"
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

    local archriot_cursor_index="$HOME/.local/share/archriot/default/icons/default/index.theme"
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

# Setup ArchRiot theme system
setup_archriot_theme_system() {
    echo "🎨 Setting up ArchRiot theme system..."

    # Create theme directories
    mkdir -p ~/.config/archriot/{themes,current,backgrounds}

    # Link available themes
    local themes_source="$HOME/.local/share/archriot/themes"
    if [[ -d "$themes_source" ]]; then
        for theme_dir in "$themes_source"/*; do
            if [[ -d "$theme_dir" ]]; then
                local theme_name=$(basename "$theme_dir")
                ln -sf "$theme_dir" ~/.config/archriot/themes/
                echo "✓ Linked theme: $theme_name"
            fi
        done
    else
        echo "⚠ ArchRiot themes directory not found at: $themes_source"
        return 1
    fi

    echo "✓ ArchRiot theme system initialized"
}

# Set default theme (CypherRiot or first available)
# The theme system works by creating symlinks in ~/.config/archriot/current/
# that point to the active theme's configuration files
set_default_theme() {
    echo "🎯 Setting default theme..."

    # CypherRiot is the flagship theme for ArchRiot
    local preferred_theme="cypherriot"
    local theme_path="$HOME/.config/archriot/themes/$preferred_theme"

    # Check if preferred theme exists and is properly structured
    if [[ -d "$theme_path" ]]; then
        echo "✓ Using preferred theme: $preferred_theme"
    else
        # Find first available theme
        local first_theme=$(find ~/.config/archriot/themes -maxdepth 1 -type d -not -name themes | head -1)
        if [[ -n "$first_theme" ]]; then
            preferred_theme=$(basename "$first_theme")
            theme_path="$first_theme"
            echo "✓ Using fallback theme: $preferred_theme"
        else
            echo "❌ No themes available"
            return 1
        fi
    fi

    # Create symlink: ~/.config/archriot/current/theme -> selected theme directory
    # This allows other components (hyprlock, waybar, etc.) to source theme configs
    # Set current theme link
    ln -snf "$theme_path" ~/.config/archriot/current/theme
    echo "✓ Active theme set to: $preferred_theme"

    # Initialize the background cycling system for this theme
    # Force error propagation when sourced
    set -e
    if setup_theme_backgrounds "$preferred_theme"; then
        echo "✅ Background system setup completed successfully"
        set +e
    else
        echo "🚨 CRITICAL: Background system setup FAILED!"
        echo "🚨 Background cycling will not work!"
        echo "🚨 Please check the installation logs above for details"
        set +e
        exit 1  # Force failure even when sourced
    fi
}

# Setup background system for selected theme
# Each theme can have its own collection of background images
# The system creates numbered copies for easy cycling (01-name.png, 02-name.png, etc.)
setup_theme_backgrounds() {
    local theme_name="$1"
    echo "🖼️  Setting up backgrounds for theme: $theme_name"

    # Ensure required directories exist
    export BACKGROUNDS_DIR="$HOME/.config/archriot/backgrounds"
    mkdir -p "$BACKGROUNDS_DIR/$theme_name"
    mkdir -p "$HOME/.config/archriot/current"

    # Copy background files directly from repo
    # Use absolute path resolution to work in any execution context
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    # Resolve the absolute path to avoid any relative path issues when sourced
    local archriot_root="$(cd "$script_dir/../.." && pwd)"
    local source_bg_dir="$archriot_root/themes/$theme_name/backgrounds"

    echo "🔍 Checking source directory: $source_bg_dir"
    echo "🔍 Script directory: $script_dir"
    echo "🔍 ArchRiot root: $archriot_root"
    echo "🔍 Current working directory: $(pwd)"

    if [[ ! -d "$source_bg_dir" ]]; then
        echo "🚨 CRITICAL: Source background directory not found: $source_bg_dir"
        echo "🚨 This will cause background cycling to fail!"
        return 1
    fi

    echo "📦 Installing backgrounds from: $source_bg_dir"

    # Clear existing numbered backgrounds
    find "$BACKGROUNDS_DIR/$theme_name" -name "[0-9][0-9]-*" -type f -delete 2>/dev/null || true

    # Find all background files
    mapfile -t bg_files < <(find "$source_bg_dir" -maxdepth 1 -type f \( -iname "*.png" -o -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.webp" \) | sort)

    if [[ ${#bg_files[@]} -eq 0 ]]; then
        echo "🚨 CRITICAL: No background files found in $source_bg_dir"
        echo "🚨 This will cause background cycling to fail!"
        return 1
    fi

    # Copy all backgrounds with numbered prefixes
    local counter=1
    local failed_copies=0
    for bg in "${bg_files[@]}"; do
        filename=$(basename "$bg")
        local dest_file="$BACKGROUNDS_DIR/$theme_name/$(printf "%02d" $counter)-$filename"

        if cp "$bg" "$dest_file"; then
            echo "✓ Installed: $(printf "%02d" $counter)-$filename"
        else
            echo "❌ Failed to install: $filename"
            ((failed_copies++))
        fi
        ((counter++))
    done

    if [[ $failed_copies -gt 0 ]]; then
        echo "🚨 WARNING: $failed_copies background files failed to copy"
    fi

    local total_installed=$((counter-1))
    echo "✓ Installed $total_installed backgrounds for $theme_name"

    # Verify backgrounds were actually installed
    local actual_count=$(find "$BACKGROUNDS_DIR/$theme_name" -name "[0-9][0-9]-*" -type f | wc -l)
    if [[ $actual_count -eq 0 ]]; then
        echo "🚨 CRITICAL: No backgrounds were successfully installed!"
        echo "🚨 Background cycling will not work!"
        return 1
    fi

    echo "✓ Verified: $actual_count background files in destination"

    # Create directory structure and symlinks
    echo "🔗 Setting up directory structure..."

    # Link background directory
    if ln -snf "$BACKGROUNDS_DIR/$theme_name" ~/.config/archriot/current/backgrounds; then
        echo "✓ Background directory linked"
    else
        echo "❌ Failed to link background directory"
        return 1
    fi

    # Set default background (riot_01.jpg preferred, or first available)
    local riot_01_bg=$(find "$BACKGROUNDS_DIR/$theme_name" -name "*riot_01*" | head -1)

    if [[ -n "$riot_01_bg" && -f "$riot_01_bg" ]]; then
        if ln -snf "$riot_01_bg" ~/.config/archriot/current/background; then
            echo "✓ Default background set: $(basename "$riot_01_bg")"
        else
            echo "❌ Failed to set default background"
            return 1
        fi
    else
        # Fallback to first numbered background (should be 01-)
        local first_bg=$(find "$BACKGROUNDS_DIR/$theme_name" -type f \( -name "*.jpg" -o -name "*.png" -o -name "*.jpeg" -o -name "*.webp" \) | sort | head -1)
        if [[ -n "$first_bg" ]]; then
            if ln -snf "$first_bg" ~/.config/archriot/current/background; then
                echo "✓ Default background set: $(basename "$first_bg")"
            else
                echo "❌ Failed to set default background"
                return 1
            fi
        else
            echo "🚨 CRITICAL: No background files found after installation"
            return 1
        fi
    fi

    # Final verification
    echo "🧪 Verifying background setup..."
    if [[ -L ~/.config/archriot/current/background && -L ~/.config/archriot/current/backgrounds ]]; then
        local bg_target=$(readlink ~/.config/archriot/current/background)
        local bg_dir_target=$(readlink ~/.config/archriot/current/backgrounds)

        if [[ -f "$bg_target" && -d "$bg_dir_target" ]]; then
            echo "✅ Background setup verification passed"
            echo "   • Background: $(basename "$bg_target")"
            echo "   • Directory: $actual_count files"
            return 0
        else
            echo "🚨 CRITICAL: Background setup verification failed!"
            echo "   • Background exists: $([[ -f "$bg_target" ]] && echo "YES" || echo "NO")"
            echo "   • Directory exists: $([[ -d "$bg_dir_target" ]] && echo "YES" || echo "NO")"
            return 1
        fi
    else
        echo "🚨 CRITICAL: Required symlinks not created!"
        return 1
    fi
}

# Link theme configurations to application configs
link_theme_configs() {
    echo "🔗 Linking theme configurations..."

    local theme_dir="$HOME/.config/archriot/current/theme"
    if [[ ! -d "$theme_dir" ]]; then
        echo "❌ No active theme found"
        return 1
    fi

    # Backup existing configs
    backup_existing_configs

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

    # Handle hyprlock config specially (link main config that sources theme)
    local main_hyprlock="$HOME/.local/share/archriot/config/hypr/hyprlock.conf"
    local target_hyprlock="$HOME/.config/hypr/hyprlock.conf"

    if [[ -f "$main_hyprlock" ]]; then
        mkdir -p "$HOME/.config/hypr"

        # Remove existing file/symlink to prevent "same file" errors
        [[ -e "$target_hyprlock" ]] && rm -f "$target_hyprlock"

        ln -snf "$main_hyprlock" "$target_hyprlock"
        echo "✓ Linked: hyprlock.conf (sources theme)"
    else
        echo "⚠ Main hyprlock config not found"
    fi

    # Link other application configs
    local config_links=(
        "fuzzel.ini:$HOME/.config/fuzzel/fuzzel.ini"
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

    local backup_dir="$HOME/.config/archriot-backups/$(date +%Y-%m-%d-%H%M%S)"
    mkdir -p "$backup_dir"

    local configs_to_backup=(
        "$HOME/.config/waybar/config"
        "$HOME/.config/fuzzel/fuzzel.ini"
        "$HOME/.config/hypr/hyprlock.conf"
        "$HOME/.config/mako/config"
        "$HOME/.config/gtk-3.0/gtk.css"
        "$HOME/.config/gtk-3.0/settings.ini"
        "$HOME/.config/gtk-4.0/settings.ini"
    )

    local backed_up=0
    for config in "${configs_to_backup[@]}"; do
        if [[ -f "$config" && ! -L "$config" ]]; then
            local relative_path="${config#$HOME/.config/}"
            local backup_path="$backup_dir/$relative_path"

            mkdir -p "$(dirname "$backup_path")"
            cp "$config" "$backup_path" 2>/dev/null && ((backed_up++))
        fi
    done

    if [[ $backed_up -gt 0 ]]; then
        echo "✓ Configuration backups created at: $backup_dir"
    else
        rmdir "$backup_dir" 2>/dev/null || true
        echo "✓ No configurations to backup"
    fi
}

# Link waybar configuration with special handling
link_waybar_config() {
    echo "📊 Configuring waybar theme..."

    local theme_config="$HOME/.config/archriot/current/theme/config"
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

    # Services will be restarted at end of installation

    # Start background service
    local bg_file="$HOME/.config/archriot/current/background"
    if [[ -f "$bg_file" ]] || [[ -L "$bg_file" ]]; then
        # Resolve symlink if needed
        local actual_bg=$(readlink -f "$bg_file" 2>/dev/null || echo "$bg_file")
        if [[ -f "$actual_bg" ]]; then
            swaybg -i "$actual_bg" -m fill >/dev/null 2>&1 &
            echo "✓ Background service started with: $(basename "$actual_bg")"
        else
            echo "⚠ Background file not found: $actual_bg"
        fi
    else
        echo "⚠ No background configured"
    fi

    # Waybar will be started at the end of installation
    echo "✓ Waybar configuration ready"

    # Start blue light filter by default (no user prompt)
    start_blue_light_filter
    echo "✓ Blue light filter enabled by default (hyprsunset configured in base hyprland.conf)"
}

# Stop theme-related services safely
stop_theme_services() {
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
        hyprsunset -t 3500 >/dev/null 2>&1 &
        sleep 2

        if pgrep hyprsunset >/dev/null; then
            echo "✓ Blue light filter started (3500K)"
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
    if [[ -L ~/.config/archriot/current/theme ]]; then
        local active_theme=$(basename "$(readlink ~/.config/archriot/current/theme)")
        echo "✓ Active theme: $active_theme"
    else
        echo "❌ No active theme set"
        ((validation_errors++))
    fi

    # Check background
    if [[ -f ~/.config/archriot/current/background ]]; then
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
    echo "  • Icons: Tela-purple-dark"
    echo "  • GTK: Adwaita Dark"
    echo "  • ArchRiot: $(basename "$(readlink ~/.config/archriot/current/theme 2>/dev/null)" 2>/dev/null || echo "Not set")"
    echo ""
    echo "🎯 Active features:"
    echo "  • Dark mode preference"
    echo "  • Blue light filter (4000K)"
    echo "  • Dynamic wallpapers"
    echo "  • Consistent theming across applications"
    echo ""
    echo "🔧 Theme management:"
    echo "  • Use theme-next to cycle themes"
    echo "  • Themes located in ~/.config/archriot/themes/"
    echo "  • Current theme: ~/.config/archriot/current/theme"
}

# Main execution with comprehensive error handling
main() {
    # Force error propagation when sourced by installer
    set -e

    echo "🚀 Starting desktop theming setup..."

    load_user_environment

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

    # Setup ArchRiot theme system
    setup_archriot_theme_system || {
        echo "❌ Failed to setup ArchRiot theme system"
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

    # Setup fuzzel cache directory
    setup_fuzzel_cache || {
        echo "⚠ Fuzzel cache setup had issues"
    }

    # Start services
    start_theme_services

    # Validate setup
    validate_theme_setup || {
        echo "⚠ Theme validation had issues"
    }

    display_theming_summary
    echo "✅ Desktop theming setup completed!"
    set +e  # Reset error handling
}

# Execute main function only if not sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
else
    # When sourced by installer, still run main but handle errors differently
    main "$@" || exit 1
fi
