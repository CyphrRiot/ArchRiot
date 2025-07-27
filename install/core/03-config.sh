#!/bin/bash

# CRITICAL BUG FIX: Waybar Binary Corruption Prevention + Installation Order Fix
#
# This script previously had a critical bug where the waybar binary at ~/.local/bin/waybar
# could be overwritten with a CSS file during installation, causing waybar to fail with
# errors like "@define-color: command not found".
#
# Root cause: The original script used:
#   find "$script_source" -type f \( -name "*.py" -o -name "*.sh" \) -exec cp {} "$script_dest/" \;
#   chmod +x "$script_dest"/*
#
# This could overwrite system binaries if there were naming conflicts or path issues.
#
# INSTALLATION ORDER FIX:
# This script now runs AFTER the desktop module installation, ensuring waybar and related
# components are properly installed before configuration and validation occurs.
# Previous order: core -> system -> desktop (waybar validation failed)
# Fixed order: core/01-02 -> desktop -> core/03-04 -> system (waybar exists for validation)
#
# Fixes implemented:
# 1. Individual script copying with safety checks
# 2. System binary protection (waybar, hyprland, etc.)
# 3. Pre-installation corruption detection and cleanup
# 4. Post-installation verification of waybar functionality
# 5. Enhanced error reporting and validation
# 6. Fixed installation module ordering
#
# This prevents the critical "curl -fsSL https://OhmArchy.org/setup.sh | bash" bug
# that was corrupting users' waybar installations and the "No waybar binary found" error.

# Load user environment and create backup
setup_environment() {
    local env_file="$HOME/.config/archriot/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Create surgical backup system
    create_surgical_backup
}

# Surgical backup system - only backup ArchRiot-managed configs
create_surgical_backup() {
    local backup_dir="$HOME/.archriot/backups/$(date +%Y-%m-%d-%H%M%S)"
    local configs_to_backup=(
        "btop" "environment.d" "fastfetch" "fish" "fuzzel"
        "ghostty" "gtk-3.0" "gtk-4.0" "hypr" "nvim"
        "systemd" "waybar" "xdg-desktop-portal" "zed"
    )

    echo "📦 Creating surgical backup of ArchRiot-managed configs only..."

    for config in "${configs_to_backup[@]}"; do
        local source="$HOME/.config/$config"
        local target="$backup_dir/$config"

        if [[ -e "$source" ]]; then
            echo "  → Backing up: $config"
            mkdir -p "$(dirname "$target")"
            cp -R "$source" "$target"
        fi
    done

    # Save backup manifest
    printf '%s\n' "${configs_to_backup[@]}" > "$backup_dir/MANIFEST"
    echo "$backup_dir" > /tmp/archriot-config-backup
    echo "✓ Surgical backup created at: $backup_dir"
}

# Smart restoration from surgical backup
restore_from_backup() {
    local backup_dir="$1"

    if [[ ! -d "$backup_dir" || ! -f "$backup_dir/MANIFEST" ]]; then
        echo "⚠️ Invalid backup directory or missing manifest"
        return 1
    fi

    echo "🔄 Restoring from surgical backup: $backup_dir"

    while IFS= read -r config; do
        local source="$backup_dir/$config"
        local target="$HOME/.config/$config"

        if [[ -e "$source" ]]; then
            echo "  → Restoring: $config"
            rm -rf "$target"
            cp -R "$source" "$target"
        fi
    done < "$backup_dir/MANIFEST"

    echo "✓ Configuration restored from backup"
}

# Install dependencies and copy configurations
install_configs() {
    echo "📦 Installing configurations and dependencies..."

    # Install Python dependencies
    command -v python3 &>/dev/null || sudo pacman -S --noconfirm python3
    sudo pacman -S --noconfirm --needed python-psutil
    python3 -c "import psutil" || return 1

    # Install OhmArchy configs (SAFELY - preserve user configs)
    local source_config="$HOME/.local/share/archriot/config"
    [[ -d "$source_config" ]] || return 1
    mkdir -p ~/.config

    # Safe config installation - preserve existing user configurations
    for item in "$source_config"/*; do
        local basename=$(basename "$item")
        local target="$HOME/.config/$basename"

        # Skip hypr configs - they're installed by desktop module
        if [[ "$basename" == "hypr" ]]; then
            echo "ℹ️ Skipping hypr config (installed by desktop module)"
            continue
        fi

        # Skip Zed configs - preserve user editor customizations
        if [[ "$basename" == "zed" ]]; then
            echo "ℹ️ Skipping zed config (preserving user editor settings)"
            # Create reference copy for users who want to see ArchRiot's config
            if [[ ! -e "$target.archriot-default" ]]; then
                cp -R "$item" "$target.archriot-default" 2>/dev/null || true
                echo "  → Created reference copy: $basename.archriot-default"
            fi
            continue
        fi

        # Smart GTK config handling - preserve user bookmarks
        if [[ "$basename" == "gtk-3.0" ]]; then
            echo "ℹ️ Smart GTK-3.0 config installation (preserving bookmarks)"

            # Ensure target directory exists
            mkdir -p "$target"

            # Copy our theme files
            cp -R "$item"/* "$target/" 2>/dev/null || true

            # Create proper bookmarks with expanded HOME paths
            local bookmarks_file="$target/bookmarks"

            # Check if bookmarks exist and contain literal $HOME (broken)
            if [[ -f "$bookmarks_file" ]] && grep -q '\$HOME' "$bookmarks_file"; then
                echo "🔧 Fixing broken Thunar bookmarks (literal \$HOME found)"
                rm "$bookmarks_file"
            fi

            if [[ ! -f "$bookmarks_file" ]]; then
                echo "📁 Creating default Thunar bookmarks with proper paths"
                echo "file://${HOME}/Downloads Downloads" > "$bookmarks_file"
                echo "file://${HOME}/Documents Documents" >> "$bookmarks_file"
                echo "file://${HOME}/Pictures Pictures" >> "$bookmarks_file"
                echo "file://${HOME}/Music Music" >> "$bookmarks_file"
                echo "file://${HOME}/Videos Videos" >> "$bookmarks_file"
                echo "✓ Thunar bookmarks created with expanded paths"
            else
                echo "✓ Preserved existing Thunar bookmarks"
            fi
            echo "✓ Smart installed GTK config with bookmark preservation"
            continue
        fi

        # FORCE OVERWRITE ALL CONFIGS (removes user customizations!)
        if [[ ! -e "$target" ]]; then
            # New installation - safe to copy
            cp -R "$item" "$target" || return 1
            echo "✓ Installed new config: $basename"
        else
            # Remove existing and install ArchRiot version (backup already created)
            rm -rf "$target" 2>/dev/null || true
            cp -R "$item" "$target" || return 1
            echo "✓ Force installed ArchRiot config: $basename"
        fi
    done

    echo "✓ Configurations installed"
}

# Pre-installation safety check
pre_installation_safety_check() {
    echo "🔍 Running pre-installation safety checks..."

    # Check if waybar binary exists and is not corrupted
    if [[ -f "$HOME/.local/bin/waybar" ]]; then
        # Check if it's a CSS file (corrupted)
        if head -1 "$HOME/.local/bin/waybar" 2>/dev/null | grep -q "define-color\|/\*.*\*/\|@import"; then
            echo "🚨 CRITICAL: Corrupted waybar binary detected! Removing before installation..."
            rm -f "$HOME/.local/bin/waybar"
            echo "✓ Removed corrupted waybar binary"
        fi
    fi

    # Backup existing ~/.local/bin if it has potential conflicts
    local backup_bin_dir="$HOME/.local/bin.backup-$(date +%s)"
    if [[ -d "$HOME/.local/bin" ]]; then
        # Check for any non-executable files that might cause issues
        local problematic_files=()
        while IFS= read -r -d '' file; do
            if [[ ! -x "$file" ]] && [[ -f "$file" ]]; then
                problematic_files+=("$file")
            fi
        done < <(find "$HOME/.local/bin" -type f -print0 2>/dev/null)

        if [[ ${#problematic_files[@]} -gt 0 ]]; then
            echo "⚠ Found ${#problematic_files[@]} non-executable files in ~/.local/bin"
            echo "Creating backup at: $backup_bin_dir"
            cp -R "$HOME/.local/bin" "$backup_bin_dir"
        fi
    fi

    echo "✓ Pre-installation safety checks completed"
}

# Setup scripts and environment
setup_scripts_and_env() {
    echo "📊 Setting up scripts and environment..."

    # Install ALL desktop applications (including hidden folder to remove menu bloat)
    local app_source="$HOME/.local/share/archriot/applications"
    if [[ -d "$app_source" ]]; then
        echo "📱 Installing desktop applications and menu cleanup..."
        mkdir -p ~/.local/share/applications

        # Copy ALL applications including the critical hidden folder
        cp -r "$app_source"/* ~/.local/share/applications/ 2>/dev/null || true
        echo "✓ Desktop applications installed (including hidden menu cleanup)"

        # Update desktop database to apply changes immediately
        if command -v update-desktop-database >/dev/null 2>&1; then
            update-desktop-database ~/.local/share/applications/
            echo "✓ Desktop database updated"
        fi
    else
        echo "⚠ Applications folder not found at $app_source"
    fi

    # Install waybar scripts
    local script_source="$HOME/.local/share/archriot/bin/scripts"
    local script_dest="$HOME/.local/bin"
    [[ -d "$script_source" ]] || return 1

    mkdir -p "$script_dest"

    # Copy scripts individually to avoid overwriting system binaries
    find "$script_source" -type f \( -name "*.py" -o -name "*.sh" \) -print0 | while IFS= read -r -d '' script_file; do
        script_name="$(basename "$script_file")"
        # Safety check: don't overwrite system binaries
        if [[ "$script_name" =~ ^(waybar|hyprland|fish|bash|zsh|nvim|vim)$ ]]; then
            echo "⚠ Skipping system binary: $script_name"
            continue
        fi
        cp "$script_file" "$script_dest/$script_name"
        chmod +x "$script_dest/$script_name"
    done

    # Install volume OSD script
    local volume_script="$HOME/.local/share/archriot/bin/volume-osd"
    if [[ -f "$volume_script" ]]; then
        cp "$volume_script" "$script_dest/"
        chmod +x "$script_dest/volume-osd"
        echo "✓ Volume OSD script installed"
    else
        echo "⚠ Volume OSD script not found"
    fi

    # Install welcome script (force overwrite to ensure latest version)
    local welcome_script="$HOME/.local/share/archriot/bin/welcome"
    if [[ -f "$welcome_script" ]]; then
        # Remove existing welcome script first to ensure fresh install
        rm -f "$script_dest/welcome" 2>/dev/null || true
        cp "$welcome_script" "$script_dest/"
        chmod +x "$script_dest/welcome"
        echo "✓ Welcome script installed (latest version)"
    else
        echo "⚠ Welcome script not found"
    fi

    # Install version check scripts
    local version_check_script="$HOME/.local/share/archriot/bin/version-check"
    if [[ -f "$version_check_script" ]]; then
        cp "$version_check_script" "$script_dest/"
        chmod +x "$script_dest/version-check"
        echo "✓ Version check script installed"
    else
        echo "⚠ Version check script not found"
    fi



    # Install performance analysis tools
    local performance_tools=(
        "performance-analysis"
        "startup-profiler"
        "memory-profiler"
        "optimize-system"
    )

    for tool in "${performance_tools[@]}"; do
        local tool_script="$HOME/.local/share/archriot/bin/$tool"
        if [[ -f "$tool_script" ]]; then
            cp "$tool_script" "$script_dest/"
            chmod +x "$script_dest/$tool"
            echo "✓ Performance tool installed: $tool"
        else
            echo "⚠ Performance tool not found: $tool"
        fi
    done

    echo "✓ Scripts and environment configured"

    # Install systemd version check service
    local systemd_user_dir="$HOME/.config/systemd/user"
    local service_source="$HOME/.local/share/archriot/config/systemd/user"

    mkdir -p "$systemd_user_dir"

    if [[ -f "$service_source/version-check.service" ]] && [[ -f "$service_source/version-check.timer" ]]; then
        cp "$service_source/version-check.service" "$systemd_user_dir/"
        cp "$service_source/version-check.timer" "$systemd_user_dir/"

        # Reload systemd and enable the timer
        systemctl --user daemon-reload
        systemctl --user enable version-check.timer
        systemctl --user start version-check.timer

        echo "✓ Version check systemd timer installed and enabled"
        echo "🔄 Automatic update notifications will check every 4 hours"
    else
        echo "⚠ Version check systemd files not found"
    fi

    # Install welcome image from repository
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local source_image="$script_dir/../../images/welcome.png"
    local dest_dir="$HOME/.local/share/archriot/images"
    local dest_image="$dest_dir/welcome.png"

    # Create destination directory
    mkdir -p "$dest_dir"

    if [[ -f "$source_image" ]]; then
        # Check if source and destination are the same file (avoid copying to self)
        if [[ "$(realpath "$source_image")" == "$(realpath "$dest_image" 2>/dev/null || echo "$dest_image")" ]]; then
            echo "✓ Welcome image already in correct location"
        else
            cp "$source_image" "$dest_dir/"
            echo "✓ Welcome image installed"
        fi
    else
        echo "⚠ Welcome image not found at: $source_image"
    fi

    # Setup bash environment
    local archriot_bashrc="$HOME/.local/share/archriot/default/bash/rc"
    [[ -f "$archriot_bashrc" ]] && echo "source $archriot_bashrc" > ~/.bashrc

    # Setup tmux configuration
    local archriot_tmux="$HOME/.local/share/archriot/default/tmux.conf"
    if [[ -f "$archriot_tmux" ]]; then
        cp "$archriot_tmux" ~/.tmux.conf
        echo "✓ tmux configuration installed"
    else
        echo "⚠ tmux configuration not found"
    fi

    echo "✓ Scripts and environment configured"
}

# Configure system services
configure_system() {
    echo "🔐 Configuring system services..."

    # Skip display manager removal if we're already running Hyprland
    if [[ "$XDG_CURRENT_DESKTOP" == "Hyprland" ]] || [[ -n "$HYPRLAND_INSTANCE_SIGNATURE" ]]; then
        echo "✓ Running in Hyprland session - skipping display manager changes"
    else
        # Remove existing display managers (we use LUKS + hyprlock instead)
        local display_managers=("sddm" "gdm" "lightdm" "lxdm")
        local removed_any=false

        for dm in "${display_managers[@]}"; do
            if pacman -Qi "$dm" &>/dev/null; then
                echo "🔍 Found display manager: $dm"
                echo "🗑️ Removing $dm (redundant with LUKS + hyprlock)..."
                sudo systemctl stop "$dm" 2>/dev/null || true
                sudo systemctl disable "$dm" 2>/dev/null || true
                sudo pacman -Rns --noconfirm "$dm" 2>/dev/null || true
                removed_any=true
            fi
        done

        if [[ "$removed_any" == true ]]; then
            echo "✓ Display managers removed - using LUKS + autologin + hyprlock"
        else
            echo "✓ No display managers found to remove"
        fi
    fi

    # Setup autologin
    [[ -n "$USER" ]] || return 1
    sudo mkdir -p /etc/systemd/system/getty@tty1.service.d
    sudo tee /etc/systemd/system/getty@tty1.service.d/override.conf >/dev/null <<EOF
[Service]
ExecStart=
ExecStart=-/usr/bin/agetty --autologin $USER --noclear %I \$TERM
EOF

    echo "✓ Autologin configured for: $USER"
}

# Configure git and XCompose
configure_user_tools() {
    echo "📝 Configuring user tools..."

    # Git configuration
    git config --global alias.co checkout
    git config --global alias.br branch
    git config --global alias.ci commit
    git config --global alias.st status
    git config --global pull.rebase true
    git config --global init.defaultBranch master

    [[ -n "${OMARCHY_USER_NAME// /}" ]] && git config --global user.name "$OMARCHY_USER_NAME"
    [[ -n "${OMARCHY_USER_EMAIL// /}" ]] && git config --global user.email "$OMARCHY_USER_EMAIL"

    # XCompose setup
    local archriot_xcompose="$HOME/.local/share/archriot/default/xcompose"
    if [[ -f "$archriot_xcompose" ]]; then
        tee ~/.XCompose >/dev/null <<EOF
include "$archriot_xcompose"
<Multi_key> <space> <n> : "${OMARCHY_USER_NAME:-}"
<Multi_key> <space> <e> : "${OMARCHY_USER_EMAIL:-}"
EOF
    fi

    echo "✓ Git and XCompose configured"
}

# Validate critical components
validate_installation() {
    local issues=0

    # Check Python
    python3 -c "import psutil" 2>/dev/null || ((issues++))

    # CRITICAL: Final verification that waybar binary is not corrupted
    if [[ -f "$HOME/.local/bin/waybar" ]]; then
        # Check if waybar in ~/.local/bin is a CSS file (corrupted)
        if head -1 "$HOME/.local/bin/waybar" 2>/dev/null | grep -q "define-color\|/\*.*\*/\|@import\|font-family"; then
            echo "🚨 CRITICAL: Corrupted waybar binary detected! Fixing..."
            rm -f "$HOME/.local/bin/waybar"
            echo "✓ Removed corrupted waybar binary"
            ((issues++))
        else
            # Check if it's executable and not the system waybar
            if [[ ! -x "$HOME/.local/bin/waybar" ]] || ! "$HOME/.local/bin/waybar" --version &>/dev/null; then
                echo "🚨 WARNING: Invalid waybar binary in ~/.local/bin! Removing..."
                rm -f "$HOME/.local/bin/waybar"
                echo "✓ Removed invalid waybar binary"
                ((issues++))
            fi
        fi
    fi

    # Ensure system waybar is accessible
    if ! command -v waybar &>/dev/null; then
        echo "🚨 CRITICAL: No waybar binary found!"
        ((issues++))
    else
        # Test waybar version to ensure it's working
        if waybar --version &>/dev/null; then
            echo "✓ Waybar binary verified working"
        else
            echo "🚨 WARNING: Waybar binary exists but not functioning properly"
            ((issues++))
        fi
    fi

    # Check essential scripts
    for script in waybar-tomato-timer.py waybar-cpu-aggregate.py waybar-memory-accurate.py waybar-mic-status.py volume-osd welcome; do
        [[ -x "$HOME/.local/bin/$script" ]] || ((issues++))
    done

    if [[ $issues -eq 0 ]]; then
        echo "✅ Validation passed"
        return 0
    else
        echo "⚠ $issues validation issues detected"
        return 1
    fi
}

# Main execution with rollback
main() {
    echo "🚀 Starting OhmArchy configuration setup..."

    {
        pre_installation_safety_check &&
        setup_environment &&
        install_configs &&
        setup_scripts_and_env &&
        configure_system &&
        configure_user_tools &&
        validate_installation
    } || {
        echo "❌ Setup failed, attempting rollback..."
        local backup_file="/tmp/archriot-config-backup"
        if [[ -f "$backup_file" ]]; then
            local backup_dir=$(cat "$backup_file")
            if [[ -d "$backup_dir" ]]; then
                echo "🔄 Restoring from backup: $backup_dir"
                rm -rf ~/.config
                cp -R "$backup_dir" ~/.config
                echo "✓ Configuration restored from backup"
            fi
        fi
        return 1
    }

    rm -f /tmp/archriot-config-backup
    echo "✅ Configuration installation complete!"
}

main "$@"
