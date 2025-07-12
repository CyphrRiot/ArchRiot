#!/bin/bash

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

    # Install OhmArchy configs
    local source_config="$HOME/.local/share/omarchy/config"
    [[ -d "$source_config" ]] || return 1
    mkdir -p ~/.config
    cp -R "$source_config"/* ~/.config/ || return 1

    echo "✓ Configurations installed"
}

# Setup scripts and environment
setup_scripts_and_env() {
    echo "📊 Setting up scripts and environment..."

    # Install waybar scripts
    local script_source="$HOME/.local/share/omarchy/bin/scripts"
    local script_dest="$HOME/.local/bin"
    [[ -d "$script_source" ]] || return 1

    mkdir -p "$script_dest" ~/.local/share/applications
    find "$script_source" -type f \( -name "*.py" -o -name "*.sh" \) -exec cp {} "$script_dest/" \;
    chmod +x "$script_dest"/*

    # Setup bash environment
    local omarchy_bashrc="$HOME/.local/share/omarchy/default/bash/rc"
    [[ -f "$omarchy_bashrc" ]] && echo "source $omarchy_bashrc" > ~/.bashrc

    echo "✓ Scripts and environment configured"
}

# Configure system services
configure_system() {
    echo "🔐 Configuring system services..."

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

    # Check essential scripts
    for script in waybar-tomato-timer.py waybar-cpu-aggregate.py waybar-memory-accurate.py waybar-mic-status.py; do
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
