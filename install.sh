#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Log file will be set after progress system loads
INSTALL_LOG_FILE=""

echo ""
echo "🚀 Installing ArchRiot... "
echo ""


# Optional installation logging for troubleshooting
if [[ "${ARCHRIOT_DEBUG:-}" == "1" ]]; then
    mkdir -p "$(dirname "$INSTALL_LOG_FILE")"
    echo "🗂️  Debug mode enabled - logging to: $INSTALL_LOG_FILE"
    exec > >(tee -a "$INSTALL_LOG_FILE") 2>&1
fi

# Function to log failures visibly
log_failure() {
    local module="$1"
    local error="$2"
    echo "[$(date)] FAILURE: $module - $error" >> "$ARCHRIOT_LOG_FILE"
    echo "❌ $module failed but installation continuing..."
}

# Simple error handler that gives retry instructions and cleans up sudo
cleanup_on_exit() {
    echo "❌ ArchRiot installation failed! You can retry by running: source $HOME/.local/share/archriot/install.sh"
    echo "💡 Most components are idempotent - re-running will skip already installed items"
    # Clean up sudo if helper is available
    if [ -f "$HOME/.local/share/archriot/install/lib/sudo-helper.sh" ]; then
        source "$HOME/.local/share/archriot/install/lib/sudo-helper.sh" 2>/dev/null || true
        if command -v remove_passwordless_rule &>/dev/null; then
            echo "🔒 Cleaning up sudo configuration..."
            remove_passwordless_rule || true
        fi
    fi
}
trap cleanup_on_exit ERR

# Load shared installation helpers
if [ -f "$HOME/.local/share/archriot/install/lib/install-helpers.sh" ]; then
    source "$HOME/.local/share/archriot/install/lib/install-helpers.sh"
fi

# CRITICAL: Install yay FIRST before anything else that might need it
# Many ArchRiot packages come from AUR and require yay to install
#
# echo "🚀 CRITICAL: Installing yay AUR helper befoCRITICAL: Installre anything else..."

# Install base development tools and yay immediately
sudo pacman -Sy --noconfirm --needed base-devel git rsync bc || {
    echo "❌ CRITICAL: Failed to install base development tools"
    exit 1
}

# Install yay immediately
if ! command -v yay &>/dev/null; then
    echo "📦 Installing yay AUR helper..."
    cd /tmp
    git clone https://aur.archlinux.org/yay-bin.git || {
        echo "❌ CRITICAL: Failed to clone yay-bin repository"
        exit 1
    }
    cd yay-bin
    makepkg -si --noconfirm || {
        echo "❌ CRITICAL: yay installation failed"
        exit 1
    }
    cd /
    rm -rf /tmp/yay-bin

    # Refresh PATH
    export PATH="/usr/bin:$PATH"
    hash -r 2>/dev/null || true

    # Verify yay is now available after PATH refresh
    if ! command -v yay &>/dev/null; then
        echo "❌ CRITICAL: yay not found in PATH after installation"
        echo "   Try opening a new terminal and running the installer again"
        exit 1
    fi

    # Verify yay is now available
    if ! command -v yay &>/dev/null; then
        echo "❌ CRITICAL: yay installation failed"
        exit 1
    fi
    echo "✅ yay AUR helper installed successfully"
else
    echo "✅ yay AUR helper already available"
fi

# Load clean progress system for better user feedback during installation
if [ -f "$HOME/.local/share/archriot/install/lib/simple-progress.sh" ]; then
    source "$HOME/.local/share/archriot/install/lib/simple-progress.sh"
    # Initialize single error log file
    init_error_log
    # Set legacy variable for backward compatibility
    INSTALL_LOG_FILE="$ARCHRIOT_LOG_FILE"
fi



# Load and setup sudo helper for passwordless installation
if [ -f "$HOME/.local/share/archriot/install/lib/sudo-helper.sh" ]; then
    source "$HOME/.local/share/archriot/install/lib/sudo-helper.sh"
    echo "🔒 Setting up temporary passwordless sudo for installation..."
    if setup_passwordless_sudo; then
        echo "🔍 Validating passwordless sudo configuration..."
        if validate_passwordless_sudo; then
            echo "✅ Passwordless sudo is working correctly"
        else
            echo "❌ Passwordless sudo validation failed!"
            echo "📋 Installation will continue but may prompt for passwords"
            echo "💡 You can enter your password when prompted during installation"
            sleep 3
        fi
    else
        echo "❌ Failed to setup passwordless sudo!"
        echo "📋 Installation will continue but will prompt for passwords"
        echo "💡 You will need to enter your password when prompted during installation"
        sleep 3
    fi
else
    echo "⚠ Sudo helper not found - installation may prompt for passwords"
fi

# Define installation order for modular structure
# CRITICAL: Order matters! Desktop must come before config validation
# This ensures waybar and other desktop components exist before validation
declare -a install_modules=(
    "core/02-identity.sh"   # User identity setup
    "desktop"               # Desktop environment (hyprland, waybar, apps, theming, fonts)
    "core/03-config.sh"     # Config installation and validation (after desktop components exist)
    "core/04-shell.sh"      # Shell configuration
    "system"                # System-level functionality (audio, networking, bluetooth, etc.)
    "development"           # Development tools (editors, tools)
    "applications"          # User applications (media, productivity, communication, utilities)
    "optional"              # Optional components (specialty apps)
)

# Additional standalone installers to run after modules
declare -a standalone_installers=(
    "plymouth.sh"    # Boot splash screen with ArchRiot branding
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
            "development.sh"|"nvim.sh")
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

# Count total installer files for accurate progress tracking
count_total_installers() {
    local total=0
    for module in "${install_modules[@]}"; do
        if [[ "$module" == *".sh" ]]; then
            # Individual file
            total=$((total + 1))
        else
            # Module directory - count .sh files
            local module_dir="$HOME/.local/share/archriot/install/$module"
            if [[ -d "$module_dir" ]]; then
                local count=0
                for file in "$module_dir"/*.sh; do
                    if [[ -f "$file" ]]; then
                        count=$((count + 1))
                    fi
                done
                total=$((total + count))
            fi
        fi
    done

    # Add standalone installers
    total=$((total + ${#standalone_installers[@]}))
    echo $total
}

# Process installation modules in the correct order
process_installation_modules() {
    # Count total installers for accurate progress
    local total_installers=$(count_total_installers)

    # Initialize clean progress system
    if command -v init_clean_progress &>/dev/null; then
        init_clean_progress $total_installers
    fi

    for module in "${install_modules[@]}"; do
        # Handle individual files vs module directories
        if [[ "$module" == *".sh" ]]; then
            # Individual file (like core/01-base.sh)
            local installer_file="$HOME/.local/share/archriot/install/$module"
            local installer_name=$(basename "$module" .sh)

            if [[ -f "$installer_file" ]]; then
                process_installer_with_progress "$installer_file" "$installer_name"
            else
                if command -v fail_step &>/dev/null; then
                    fail_step "Installer not found: $module"
                fi
                exit 1
            fi
        else
            # Module directory (like desktop, system, etc.)
            local module_dir="$HOME/.local/share/archriot/install/$module"

            if [[ -d "$module_dir" ]]; then
                # Process all .sh files in the module directory
                for installer_file in "$module_dir"/*.sh; do
                    if [[ -f "$installer_file" ]]; then
                        local installer_name=$(basename "$installer_file" .sh)
                        process_installer_with_progress "$installer_file" "$installer_name"
                    fi
                done
            else
                if command -v fail_step &>/dev/null; then
                    fail_step "Module directory not found: $module"
                fi
                exit 1
            fi
        fi
    done
}

# Process individual installer with clean progress
process_installer_with_progress() {
    local installer_file="$1"
    local installer_name="$2"

    # Get color for installer type
    local color="BLUE"
    case "$installer_name" in
        *base*|*core*|*yay*) color="BLUE" ;;
        *identity*|*config*) color="CYAN" ;;
        *desktop*|*hypr*|*waybar*) color="PURPLE" ;;
        *theme*|*font*) color="YELLOW" ;;
        *shell*|*terminal*) color="GREEN" ;;
        *audio*|*network*|*bluetooth*|*power*) color="ORANGE" ;;
        *development*|*nvim*) color="CYAN" ;;
        *application*|*media*|*productivity*) color="GREEN" ;;
        *optional*|*xtras*) color="YELLOW" ;;
        *plymouth*|*final*) color="PURPLE" ;;
    esac

    # Check if this is an interactive installer that needs direct execution
    if [[ "$installer_name" == *"identity"* ]] || [[ "$installer_name" == *"config"* ]]; then
        # Interactive installers - run directly without output capture
        if command -v start_module &>/dev/null; then
            start_module "$installer_name" "$color"
        fi

        if source "$installer_file"; then
            echo "✓ Successfully completed"
            echo ""
        else
            echo "❌ Failed: $installer_name"
            echo "⚠ System may be in partially configured state"
            echo "💡 To rollback: restore from backup at ~/.config-backup-$(date +%Y%m%d) (if created)"
            echo "🔄 To retry: source ~/.local/share/archriot/install.sh"
            exit 1
        fi
    elif command -v run_command_clean &>/dev/null; then
        # Non-interactive installers - use clean progress with output capture
        if command -v start_module &>/dev/null; then
            start_module "$installer_name" "$color"
        fi

        # No special warning for Plymouth - let the script decide if it's actually installing

        run_command_clean "source '$installer_file'" "$installer_name" "$color"
    else
        # Fallback to original method

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
                echo "✓ Completed (${duration}s)"
            fi
        else
            echo "❌ Failed: $installer_name"
            echo "⚠ System may be in partially configured state"
            echo "💡 To rollback: restore from backup at ~/.config-backup-$(date +%Y%m%d) (if created)"
            echo "🔄 To retry: source ~/.local/share/archriot/install.sh"
            exit 1
        fi
    fi
}

# Read version from VERSION file (single source of truth)
# Priority: local repo (development) → installed location → GitHub fallback
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/VERSION" ]; then
    ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
elif [ -f "$HOME/.local/share/archriot/VERSION" ]; then
    ARCHRIOT_VERSION=$(cat "$HOME/.local/share/archriot/VERSION" 2>/dev/null || echo "unknown")
else
    # Fetch version from GitHub when running via curl
    ARCHRIOT_VERSION=$(curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION 2>/dev/null || echo "unknown")
fi

# Track installation performance
INSTALL_START_TIME=$(date +%s)
INSTALL_START_DATE=$(date)

echo "🚀 Starting ArchRiot Installation (Fixed Module Order)"
echo "====================================================="
echo "Version: $ARCHRIOT_VERSION"
echo "Total modules: ${#install_modules[@]}"
echo "Start time: $INSTALL_START_DATE"
echo "🔒 Sudo status: $(if sudo -n true 2>/dev/null; then echo "Passwordless ✓"; else echo "Will prompt for password"; fi)"
echo "🔧 Fix: Desktop environment installs before config validation"
echo



# Process installation modules in correct order
process_installation_modules

# Run standalone installers (included in progress tracking)
for standalone in "${standalone_installers[@]}"; do
    standalone_path="$HOME/.local/share/archriot/install/$standalone"
    standalone_name=$(basename "$standalone" .sh)

    if [[ -f "$standalone_path" ]]; then
        process_installer_with_progress "$standalone_path" "$standalone_name"
    else
        if command -v fail_step &>/dev/null; then
            fail_step "Standalone installer not found: $standalone"
        else
            echo "⚠ Standalone installer not found: $standalone"
        fi
    fi
done

# Clean up installation files if helpers are available
if command -v cleanup_install_files &>/dev/null; then
    cleanup_install_files
fi

# Ensure locate is up to date now that everything has been installed
sudo updatedb

# Final installation validation


# Show completion summary with progress system
if command -v complete_clean_installation &>/dev/null; then
    complete_clean_installation
fi

# Calculate installation performance
INSTALL_END_TIME=$(date +%s)
INSTALL_DURATION=$((INSTALL_END_TIME - INSTALL_START_TIME))
INSTALL_DURATION_MIN=$((INSTALL_DURATION / 60))
INSTALL_DURATION_SEC=$((INSTALL_DURATION % 60))

# Ensure gum is available for final prompt (BEFORE sudo cleanup)
if ! command -v gum &>/dev/null; then
    yay -S --noconfirm --needed gum || {
        echo "❌ CRITICAL: Failed to install gum"
        echo "   gum is required for ArchRiot user interface"
        exit 1
    }
fi

# Verify gum is actually working
if ! command -v gum &>/dev/null; then
    echo "❌ CRITICAL: gum not available after installation"
    exit 1
fi

# Update local version file after successful installation
if [[ -n "$ARCHRIOT_VERSION" && "$ARCHRIOT_VERSION" != "unknown" ]]; then
    mkdir -p "$HOME/.local/share/archriot"
    echo "$ARCHRIOT_VERSION" > "$HOME/.local/share/archriot/VERSION"
fi

# Clean up passwordless sudo after installation
if command -v remove_passwordless_rule &>/dev/null; then
    remove_passwordless_rule || true
fi

echo "================================="
echo "🎉 ArchRiot v$ARCHRIOT_VERSION installed! (${INSTALL_DURATION_MIN}m ${INSTALL_DURATION_SEC}s)"
echo ""
echo "🎯 Quick Commands:"
echo "  • Launch Apps: Super + D"
echo "  • Backgrounds: Super + Ctrl + Space"
echo "  • View Help: Super + H"
echo ""

# Show backup location if backup was created
if [[ -f /tmp/archriot-config-backup ]]; then
    local backup_location=$(cat /tmp/archriot-config-backup)
    echo "📦 Your previous configuration files were backed up to:"
    echo "   $backup_location"
    echo ""
fi

# Silently apply system updates
if pgrep -x "Hyprland" >/dev/null; then
    hyprctl reload 2>/dev/null || true
fi

if pgrep -x "waybar" >/dev/null; then
    pkill waybar 2>/dev/null || true
    sleep 1
    waybar &>/dev/null &
fi

fc-cache -fv >/dev/null 2>&1
gtk-update-icon-cache -f ~/.local/share/icons/hicolor/ 2>/dev/null || true
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
# Reload shell configuration
echo "🐚 Shell configuration will apply to new terminals"

# Check for installation failures and report them (ignore yay "up to date" warnings)
if [[ -f "$ARCHRIOT_LOG_FILE" ]] && grep -q "FAILURE:\|CRITICAL:\|❌" "$ARCHRIOT_LOG_FILE" 2>/dev/null; then
    # Double-check: exclude harmless "up to date" and "warning:" messages
    if grep -v "is up to date\|warning:" "$ARCHRIOT_LOG_FILE" | grep -q "FAILURE:\|CRITICAL:\|❌" 2>/dev/null; then
        echo ""
        echo "⚠️ ⚠️ ⚠️  INSTALLATION COMPLETED WITH FAILURES  ⚠️ ⚠️ ⚠️"
        echo ""
        echo "Some components failed during installation."
        echo "Check the log for details: $ARCHRIOT_LOG_FILE"
        echo ""
        echo "🔧 To fix these issues, run:"
        echo "   source ~/.local/share/archriot/install.sh"
        echo ""
        echo "⚠️  Your system may have missing functionality until these are resolved."
        echo ""
    else
        echo "✅ System ready! New terminals will have updated configs."
        echo "📝 Minor warnings logged (packages up-to-date) - system is healthy."
        echo ""
    fi
else
    echo "✅ System ready! New terminals will have updated configs."
    echo ""
fi

gum confirm "Reboot to ensure all settings are fully applied?" && reboot
