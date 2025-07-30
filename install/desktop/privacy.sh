#!/bin/bash

# ==============================================================================
# ArchRiot Desktop Privacy Configuration
# ==============================================================================
# Disables recent files tracking across all desktop applications and file managers
# Non-optional privacy enhancement for cleaner, more private desktop experience
# ==============================================================================

# Load user environment if available
load_user_environment() {
    local env_file="$HOME/.config/archriot/user.env"
    if [[ -f "$env_file" ]]; then
        source "$env_file"
    fi
}

# Disable GNOME/desktop recent files tracking
disable_gnome_recents() {
    echo "🔒 Disabling GNOME recent files tracking..."

    # Set GNOME privacy setting to not remember recent files
    if command -v gsettings &>/dev/null; then
        gsettings set org.gnome.desktop.privacy remember-recent-files false
        echo "  ✓ GNOME recent files disabled via gsettings"
    else
        echo "  ⚠ gsettings not available, skipping GNOME setting"
    fi
}

# Remove existing recent files history
clear_existing_recents() {
    echo "🧹 Clearing existing recent files history..."

    local recent_file="$HOME/.local/share/recently-used.xbel"
    if [[ -f "$recent_file" ]]; then
        rm "$recent_file"
        echo "  ✓ Removed existing recent files database"
    else
        echo "  ✓ No existing recent files database found"
    fi
}

# Configure GTK settings to disable recent files
configure_gtk_recents() {
    echo "⚙️  Configuring GTK recent files settings..."

    # Configure GTK-3 settings
    local gtk3_dir="$HOME/.config/gtk-3.0"
    local gtk3_settings="$gtk3_dir/settings.ini"

    mkdir -p "$gtk3_dir"

    # Check if settings.ini exists, create or update it
    if [[ -f "$gtk3_settings" ]]; then
        echo "  📝 Updating existing GTK-3 settings..."

        # Remove any existing recent files settings
        sed -i '/gtk-recent-files-enabled/d' "$gtk3_settings"
        sed -i '/gtk-recent-files-max-age/d' "$gtk3_settings"
        sed -i '/gtk-recent-files-limit/d' "$gtk3_settings"

        # Add our settings to the [Settings] section
        if grep -q "^\[Settings\]" "$gtk3_settings"; then
            # Add after [Settings] line
            sed -i '/^\[Settings\]/a gtk-recent-files-enabled=0\ngtk-recent-files-max-age=0\ngtk-recent-files-limit=0' "$gtk3_settings"
        else
            # No [Settings] section, add everything
            echo "" >> "$gtk3_settings"
            echo "[Settings]" >> "$gtk3_settings"
            echo "gtk-recent-files-enabled=0" >> "$gtk3_settings"
            echo "gtk-recent-files-max-age=0" >> "$gtk3_settings"
            echo "gtk-recent-files-limit=0" >> "$gtk3_settings"
        fi
    else
        echo "  📄 Creating GTK-3 settings file..."
        cat > "$gtk3_settings" << 'EOF'
[Settings]
gtk-recent-files-enabled=0
gtk-recent-files-max-age=0
gtk-recent-files-limit=0
EOF
    fi

    echo "  ✓ GTK-3 recent files disabled"

    # Configure GTK-4 settings
    local gtk4_dir="$HOME/.config/gtk-4.0"
    local gtk4_settings="$gtk4_dir/settings.ini"

    mkdir -p "$gtk4_dir"

    if [[ -f "$gtk4_settings" ]]; then
        echo "  📝 Updating existing GTK-4 settings..."

        # Remove any existing recent files settings
        sed -i '/gtk-recent-files-enabled/d' "$gtk4_settings"
        sed -i '/gtk-recent-files-max-age/d' "$gtk4_settings"
        sed -i '/gtk-recent-files-limit/d' "$gtk4_settings"

        # Add our settings to the [Settings] section
        if grep -q "^\[Settings\]" "$gtk4_settings"; then
            sed -i '/^\[Settings\]/a gtk-recent-files-enabled=0\ngtk-recent-files-max-age=0\ngtk-recent-files-limit=0' "$gtk4_settings"
        else
            echo "" >> "$gtk4_settings"
            echo "[Settings]" >> "$gtk4_settings"
            echo "gtk-recent-files-enabled=0" >> "$gtk4_settings"
            echo "gtk-recent-files-max-age=0" >> "$gtk4_settings"
            echo "gtk-recent-files-limit=0" >> "$gtk4_settings"
        fi
    else
        echo "  📄 Creating GTK-4 settings file..."
        cat > "$gtk4_settings" << 'EOF'
[Settings]
gtk-recent-files-enabled=0
gtk-recent-files-max-age=0
gtk-recent-files-limit=0
EOF
    fi

    echo "  ✓ GTK-4 recent files disabled"
}

# Configure Thunar to disable recent files
configure_thunar_recents() {
    echo "📁 Configuring Thunar recent files settings..."

    local thunar_config_dir="$HOME/.config/Thunar"
    local thunar_config="$thunar_config_dir/thunarrc"

    mkdir -p "$thunar_config_dir"

    if [[ -f "$thunar_config" ]]; then
        echo "  📝 Updating existing Thunar configuration..."

        # Remove any existing recent files settings
        sed -i '/LastShowRecent=/d' "$thunar_config"
        sed -i '/MiscRememberGeometry=/d' "$thunar_config"

        # Add setting to disable recent files sidebar
        echo "LastShowRecent=FALSE" >> "$thunar_config"
        echo "  ✓ Thunar recent files disabled"
    else
        echo "  📄 Creating Thunar configuration..."
        cat > "$thunar_config" << 'EOF'
[Configuration]
LastShowRecent=FALSE
MiscRememberGeometry=TRUE
EOF
        echo "  ✓ Thunar recent files disabled"
    fi
}

# Create a script to maintain privacy settings
create_privacy_maintenance() {
    echo "🛠️  Creating privacy maintenance tools..."

    local privacy_script="$HOME/.local/bin/clear-recents"
    mkdir -p "$HOME/.local/bin"

    cat > "$privacy_script" << 'EOF'
#!/bin/bash
# ArchRiot Privacy Maintenance Script
# Clears any recent files that may have been created

echo "🧹 Clearing recent files..."

# Remove recent files database
if [[ -f "$HOME/.local/share/recently-used.xbel" ]]; then
    rm "$HOME/.local/share/recently-used.xbel"
    echo "✓ Cleared recent files database"
fi

# Clear any GTK recent files
for gtk_dir in "$HOME/.local/share/gtk-"*; do
    if [[ -d "$gtk_dir" ]]; then
        find "$gtk_dir" -name "*recent*" -type f -delete 2>/dev/null
    fi
done

# Clear thumbnail cache recent files
if [[ -d "$HOME/.cache/thumbnails" ]]; then
    find "$HOME/.cache/thumbnails" -name "*.xbel" -type f -delete 2>/dev/null
fi

echo "✓ Privacy maintenance complete"
EOF

    chmod +x "$privacy_script"
    echo "  ✓ Privacy maintenance script created at ~/.local/bin/clear-recents"
}

# Verify privacy configuration
verify_privacy_settings() {
    echo "🔍 Verifying privacy configuration..."

    local issues=0

    # Check GNOME setting
    if command -v gsettings &>/dev/null; then
        local gnome_setting=$(gsettings get org.gnome.desktop.privacy remember-recent-files 2>/dev/null)
        if [[ "$gnome_setting" == "false" ]]; then
            echo "  ✓ GNOME recent files disabled"
        else
            echo "  ❌ GNOME recent files still enabled"
            ((issues++))
        fi
    fi

    # Check GTK settings
    for gtk_version in "gtk-3.0" "gtk-4.0"; do
        local gtk_settings="$HOME/.config/$gtk_version/settings.ini"
        if [[ -f "$gtk_settings" ]] && grep -q "gtk-recent-files-enabled=0" "$gtk_settings"; then
            echo "  ✓ $gtk_version recent files disabled"
        else
            echo "  ❌ $gtk_version recent files not properly disabled"
            ((issues++))
        fi
    done

    # Check for recent files database
    if [[ ! -f "$HOME/.local/share/recently-used.xbel" ]]; then
        echo "  ✓ No recent files database present"
    else
        echo "  ⚠ Recent files database still exists (may be recreated by apps)"
    fi

    if [[ $issues -eq 0 ]]; then
        echo "  ✅ All privacy settings verified successfully"
    else
        echo "  ⚠ $issues privacy setting issues detected"
    fi
}

# Main execution function
main() {
    echo "🔒 Starting desktop privacy configuration..."

    load_user_environment

    # Apply all privacy configurations
    disable_gnome_recents
    clear_existing_recents
    configure_gtk_recents
    configure_thunar_recents
    create_privacy_maintenance

    # Verify everything worked
    verify_privacy_settings

    echo ""
    echo "✅ Desktop privacy configuration completed!"
    echo ""
    echo "📋 Privacy features configured:"
    echo "  • GNOME recent files tracking disabled"
    echo "  • GTK-3 and GTK-4 recent files disabled"
    echo "  • Thunar recent files sidebar disabled"
    echo "  • Existing recent files history cleared"
    echo "  • Privacy maintenance script installed"
    echo ""
    echo "🛠️  Maintenance: Run 'clear-recents' anytime to clear any new recent files"
    echo ""
}

# Execute main function if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
