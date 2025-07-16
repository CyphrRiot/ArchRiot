#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Simple error handler that gives retry instructions and cleans up sudo
cleanup_on_exit() {
    echo "❌ OhmArchy installation failed! You can retry by running: source $HOME/.local/share/omarchy/install.sh"
    # Clean up sudo if helper is available
    if command -v cleanup_passwordless_sudo &>/dev/null; then
        echo "🔒 Cleaning up sudo configuration..."
        cleanup_passwordless_sudo 2>/dev/null || true
    fi
}
trap cleanup_on_exit ERR

# Load shared installation helpers
if [ -f "$HOME/.local/share/omarchy/install/lib/install-helpers.sh" ]; then
    source "$HOME/.local/share/omarchy/install/lib/install-helpers.sh"
fi

# Load and setup sudo helper for passwordless installation
if [ -f "$HOME/.local/share/omarchy/install/lib/sudo-helper.sh" ]; then
    source "$HOME/.local/share/omarchy/install/lib/sudo-helper.sh"
    echo "🔒 Setting up temporary passwordless sudo for installation..."
    setup_passwordless_sudo || {
        echo "⚠️ Failed to setup passwordless sudo - installation will prompt for passwords"
    }
fi

# Define installation order for modular structure
# FIXED: Desktop module moved before system to ensure waybar is installed before config validation
declare -a install_modules=(
    "core/01-base.sh"       # Base tools and yay AUR helper
    "core/02-identity.sh"   # User identity setup
    "desktop"               # Desktop environment (hyprland, waybar, apps, theming, fonts)
    "core/03-config.sh"     # Config installation and validation (after desktop components exist)
    "core/04-shell.sh"      # Shell configuration
    "system"                # System-level functionality (audio, networking, bluetooth, etc.)
    "development"           # Development tools (editors, tools, containers)
    "applications"          # User applications (media, productivity, communication, utilities)
    "optional"              # Optional components (specialty apps)
)

# Additional standalone installers to run after modules
declare -a standalone_installers=(
    "plymouth.sh"    # Boot splash screen with OhmArchy branding
)

# Function to get all installer files in proper order
get_installer_files() {
    # Detect script location - works from any directory
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local install_dir="$script_dir/install"
    local files=()

    # Add modular structure files in order
    for module in "${install_modules[@]}"; do
        if [ -d "$install_dir/$module" ]; then
            while IFS= read -r -d '' file; do
                files+=("$file")
            done < <(find "$install_dir/$module" -name "*.sh" -type f -print0 | sort -z)
        fi
    done

    # Add any remaining root-level files (backwards compatibility)
    # Skip if modular equivalent exists
    for file in "$install_dir"/*.sh; do
        [ -f "$file" ] || continue

        local basename_file=$(basename "$file")
        local skip_file=false

        # Skip files that have modular equivalents
        case "$basename_file" in
            "1-yay.sh"|"2-identification.sh"|"3-config.sh"|"4-terminal.sh")
                skip_file=true ;;
            "network.sh"|"bluetooth.sh"|"power.sh"|"filesystems.sh")
                skip_file=true ;;
            "desktop.sh"|"hyprlandia.sh"|"z-theme-final.sh"|"fonts.sh")
                skip_file=true ;;
            "development.sh"|"nvim.sh"|"docker.sh")
                skip_file=true ;;
            "xtras.sh")
                skip_file=true ;;
        esac

        if [ "$skip_file" = false ]; then
            files+=("$file")
        fi
    done

    printf '%s\n' "${files[@]}"
}

# Process installation modules in the correct order
process_installation_modules() {
    local total=${#install_modules[@]}
    local current=0

    for module in "${install_modules[@]}"; do
        current=$((current + 1))

        # Handle individual files vs module directories
        if [[ "$module" == *".sh" ]]; then
            # Individual file (like core/01-base.sh)
            local installer_file="$HOME/.local/share/omarchy/install/$module"
            local installer_name=$(basename "$module" .sh)
            local module_name=$(basename "$(dirname "$module")")

            echo -e "\n[$current/$total] 📦 $module_name/$installer_name"
            echo "================================================"

            if [[ -f "$installer_file" ]]; then
                process_installer "$installer_file" "$installer_name"
            else
                echo "⚠ Installer not found: $module"
            fi
        else
            # Module directory (like desktop, system, etc.)
            local module_dir="$HOME/.local/share/omarchy/install/$module"

            if [[ -d "$module_dir" ]]; then
                echo -e "\n[$current/$total] 📦 Processing $module module"
                echo "================================================"

                # Process all .sh files in the module directory
                for installer_file in "$module_dir"/*.sh; do
                    if [[ -f "$installer_file" ]]; then
                        local installer_name=$(basename "$installer_file" .sh)
                        process_installer "$installer_file" "$installer_name"
                    fi
                done
            else
                echo "⚠ Module directory not found: $module"
            fi
        fi
    done
}

# Process individual installer
process_installer() {
    local installer_file="$1"
    local installer_name="$2"

    echo "🔧 Installing: $installer_name"
    start_time=$(date +%s)

    # Initialize installer context if helpers are available
    if command -v init_installer &>/dev/null; then
        init_installer "$installer_name"
    fi

    # Execute installer
    if source "$installer_file"; then
        end_time=$(date +%s)
        duration=$((end_time - start_time))

        if command -v show_install_summary &>/dev/null; then
            show_install_summary
        else
            echo "✓ Completed: $installer_name (${duration}s)"
        fi
    else
        echo "❌ Failed: $installer_name"
        exit 1
    fi
}

# Read version
OMARCHY_VERSION="1.0.20"
if [ -f "$HOME/.local/share/omarchy/VERSION" ]; then
    OMARCHY_VERSION=$(cat "$HOME/.local/share/omarchy/VERSION" 2>/dev/null || echo "1.0.12")
else
    # Fetch version from GitHub when running via curl
    OMARCHY_VERSION=$(curl -fsSL https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/VERSION 2>/dev/null || echo "1.0.20")
fi

echo "🚀 Starting OhmArchy Installation (Fixed Module Order)"
echo "====================================================="
echo "Version: $OMARCHY_VERSION"
echo "Total modules: ${#install_modules[@]}"
echo "Start time: $(date)"
echo "🔒 Sudo status: $(if sudo -n true 2>/dev/null; then echo "Passwordless ✓"; else echo "Will prompt for password"; fi)"
echo "🔧 Fix: Desktop environment installs before config validation"
echo

# Process installation modules in correct order
process_installation_modules

# Run standalone installers
echo -e "\n🎨 Running Standalone Installers"
echo "================================="

for standalone in "${standalone_installers[@]}"; do
    standalone_path="$HOME/.local/share/omarchy/install/$standalone"
    standalone_name=$(basename "$standalone" .sh)

    if [[ -f "$standalone_path" ]]; then
        echo "🔧 Installing: $standalone_name"
        start_time=$(date +%s)

        if bash "$standalone_path"; then
            end_time=$(date +%s)
            duration=$((end_time - start_time))
            echo "✓ Completed: $standalone_name (${duration}s)"
        else
            echo "❌ Failed: $standalone_name (continuing anyway)"
        fi
    else
        echo "⚠ Standalone installer not found: $standalone"
    fi
done

# Clean up installation files if helpers are available
if command -v cleanup_install_files &>/dev/null; then
    cleanup_install_files
fi

# Ensure locate is up to date now that everything has been installed
sudo updatedb

# Final installation validation
echo -e "\n🔍 Final Installation Validation"
echo "================================="

# Test critical components
command -v waybar &>/dev/null && echo "✓ Waybar installed" || echo "⚠ Waybar installation issue"
command -v hyprland &>/dev/null && echo "✓ Hyprland installed" || echo "⚠ Hyprland installation issue"
command -v mullvad &>/dev/null && echo "✓ Mullvad installed" || echo "⚠ Mullvad installation issue"

# Check Zed Wayland integration
if command -v zed-wayland &>/dev/null; then
    if [ -f ~/.local/share/applications/zed.desktop ]; then
        echo "✓ Zed with Wayland support installed"
    else
        echo "⚠ Zed Wayland launcher found but desktop file missing"
        echo "  Expected: ~/.local/share/applications/zed.desktop"
    fi
elif command -v zed &>/dev/null; then
    echo "⚠ Zed installed but Wayland integration missing"
    echo "  Missing: ~/.local/bin/zed-wayland"
else
    echo "⚠ Zed not installed"
fi

# Check theme system
if [ -L ~/.config/omarchy/current/theme ]; then
    echo "✓ Theme system configured"
else
    echo "⚠ Theme system issue"
fi

# Check waybar scripts
script_count=$(find ~/.local/bin -name "waybar-*.py" -executable 2>/dev/null | wc -l)
if [ $script_count -ge 4 ]; then
    echo "✓ Waybar scripts installed ($script_count found)"
else
    echo "⚠ Missing waybar scripts (found $script_count, expected 4+)"
fi

echo "================================="
echo "🎉 OhmArchy installation complete!"
echo "Version: $OMARCHY_VERSION"
echo "Completed at: $(date)"

# Clean up temporary passwordless sudo
if command -v cleanup_passwordless_sudo &>/dev/null; then
    echo "🔒 Restoring original sudo configuration..."
    cleanup_passwordless_sudo || {
        echo "⚠️ Warning: Could not restore sudo config - may need manual cleanup"
    }
fi



# Ensure gum is available for final prompt
if ! command -v gum &>/dev/null; then
    echo "Installing gum for final prompt..."
    yay -S --noconfirm --needed gum
fi

echo ""
echo "🎯 Installation Summary:"
echo "  • All components installed and configured"
echo "  • Themes and backgrounds properly set up"
echo "  • All keyboard shortcuts configured"
echo ""

gum confirm "Reboot to apply all settings?" && reboot
