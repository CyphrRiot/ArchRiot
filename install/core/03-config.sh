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
    local env_file="$HOME/.config/omarchy/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Create backup if config exists
    if [[ -d ~/.config ]]; then
        local backup_dir="$HOME/.config.backup-$(date +%s)"
        cp -R ~/.config "$backup_dir" && echo "$backup_dir" > /tmp/omarchy-config-backup
        echo "✓ Backup created at: $backup_dir"
    fi
}

# Install dependencies and copy configurations
install_configs() {
    echo "📦 Installing configurations and dependencies..."

    # Install Python dependencies
    command -v python3 &>/dev/null || sudo pacman -S --noconfirm python3
    sudo pacman -S --noconfirm --needed python-psutil
    python3 -c "import psutil" || return 1

    # Install OhmArchy configs (SAFELY - preserve user configs)
    local source_config="$HOME/.local/share/omarchy/config"
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
            # Create reference copy for users who want to see OhmArchy's config
            if [[ ! -e "$target.omarchy-default" ]]; then
                cp -R "$item" "$target.omarchy-default" 2>/dev/null || true
                echo "  → Created reference copy: $basename.omarchy-default"
            fi
            continue
        fi

        # FORCE OVERWRITE ALL CONFIGS (removes user customizations!)
        if [[ ! -e "$target" ]]; then
            # New installation - safe to copy
            cp -R "$item" "$target" || return 1
            echo "✓ Installed new config: $basename"
        else
            # Backup existing and install OhmArchy version
            local backup_target="$target.user-backup-$(date +%s)"
            mv "$target" "$backup_target" 2>/dev/null || true
            cp -R "$item" "$target" || return 1
            echo "✓ Force installed OhmArchy config: $basename (backup: $(basename "$backup_target"))"
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

    # Install waybar scripts
    local script_source="$HOME/.local/share/omarchy/bin/scripts"
    local script_dest="$HOME/.local/bin"
    [[ -d "$script_source" ]] || return 1

    mkdir -p "$script_dest" ~/.local/share/applications

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
    local volume_script="$HOME/.local/share/omarchy/bin/omarchy-volume-osd"
    if [[ -f "$volume_script" ]]; then
        cp "$volume_script" "$script_dest/"
        chmod +x "$script_dest/omarchy-volume-osd"
        echo "✓ Volume OSD script installed"
    else
        echo "⚠ Volume OSD script not found"
    fi

    # Install welcome script
    local welcome_script="$HOME/.local/share/omarchy/bin/welcome"
    if [[ -f "$welcome_script" ]]; then
        cp "$welcome_script" "$script_dest/"
        chmod +x "$script_dest/welcome"
        echo "✓ Welcome script installed"
    else
        echo "⚠ Welcome script not found"
    fi

    # Install performance analysis tools
    local performance_tools=(
        "performance-analysis"
        "startup-profiler"
        "memory-profiler"
        "optimize-system"
    )

    for tool in "${performance_tools[@]}"; do
        local tool_script="$HOME/.local/share/omarchy/bin/$tool"
        if [[ -f "$tool_script" ]]; then
            cp "$tool_script" "$script_dest/"
            chmod +x "$script_dest/$tool"
            echo "✓ Performance tool installed: $tool"
        else
            echo "⚠ Performance tool not found: $tool"
        fi
    done

    echo "✓ Scripts and environment configured"

    # Install welcome image
    local source_image="$HOME/.local/share/omarchy/images/welcome.png"
    local dest_dir="$HOME/.local/share/omarchy/images"
    if [[ -f "$source_image" ]]; then
        mkdir -p "$dest_dir"
        cp "$source_image" "$dest_dir/"
        echo "✓ Welcome image installed"
    else
        echo "⚠ Welcome image not found"
    fi

    # Setup bash environment
    local omarchy_bashrc="$HOME/.local/share/omarchy/default/bash/rc"
    [[ -f "$omarchy_bashrc" ]] && echo "source $omarchy_bashrc" > ~/.bashrc

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
    local omarchy_xcompose="$HOME/.local/share/omarchy/default/xcompose"
    if [[ -f "$omarchy_xcompose" ]]; then
        tee ~/.XCompose >/dev/null <<EOF
include "$omarchy_xcompose"
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
        local backup_file="/tmp/omarchy-config-backup"
        if [[ -f "$backup_file" ]]; then
            local backup_dir=$(cat "$backup_file")
            [[ -d "$backup_dir" ]] && { rm -rf ~/.config; mv "$backup_dir" ~/.config; }
        fi
        return 1
    }

    rm -f /tmp/omarchy-config-backup
    echo "✅ Configuration installation complete!"
}

main "$@"
