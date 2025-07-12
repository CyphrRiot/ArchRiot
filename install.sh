#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Simple error handler that gives retry instructions
trap 'echo "❌ OhmArchy installation failed! You can retry by running: source $HOME/.local/share/omarchy/install.sh"' ERR

# Load shared installation helpers
if [ -f "$HOME/.local/share/omarchy/install/lib/install-helpers.sh" ]; then
    source "$HOME/.local/share/omarchy/install/lib/install-helpers.sh"
fi

# Define installation order for modular structure
declare -a install_modules=(
    "core"           # Essential system components (base, identity, config, shell)
    "system"         # System-level functionality (audio, networking, bluetooth, etc.)
    "desktop"        # Desktop environment (hyprland, apps, theming, fonts)
    "development"    # Development tools (editors, tools, containers)
    "applications"   # User applications (media, productivity, communication, utilities)
    "optional"       # Optional components (specialty apps)
)

# Function to get all installer files in proper order
get_installer_files() {
    local install_dir="$HOME/.local/share/omarchy/install"
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
            "xtras.sh"|"webapps.sh")
                skip_file=true ;;
        esac

        if [ "$skip_file" = false ]; then
            files+=("$file")
        fi
    done

    printf '%s\n' "${files[@]}"
}

# Get all installer files
readarray -t installers < <(get_installer_files)
total=${#installers[@]}
current=0

echo "🚀 Starting OhmArchy Installation (Modular Structure)"
echo "===================================================="
echo "Total installers: $total"
echo "Start time: $(date)"
echo

# Process each installer
for installer_file in "${installers[@]}"; do
    # Skip lib directory files
    if [[ "$installer_file" == *"/lib/"* ]]; then
        continue
    fi

    current=$((current + 1))
    installer_name=$(basename "$installer_file" .sh)
    module_name=$(basename "$(dirname "$installer_file")")

    # Show module context for organized files
    if [[ "$module_name" != "install" ]]; then
        echo -e "\n[$current/$total] 📦 $module_name/$installer_name"
    else
        echo -e "\n[$current/$total] 📦 $installer_name"
    fi
    echo "================================================"
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
echo "Completed at: $(date)"

# Ensure gum is available for final prompt
if ! command -v gum &>/dev/null; then
    echo "Installing gum for final prompt..."
    yay -S --noconfirm --needed gum
fi

gum confirm "Reboot to apply all settings?" && reboot
