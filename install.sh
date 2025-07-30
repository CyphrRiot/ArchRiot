#!/bin/bash

# ================================================================================
# ArchRiot Installation System v2.0 - Reliable Module Execution
# ================================================================================
# Eliminates set -e conflicts, provides proper logging, ensures reliable execution
# All output is visible and logged, no silent failures
# ================================================================================

# DO NOT use set -e - causes conflicts with module scripts
# Instead, handle errors explicitly where needed

# Installation configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly INSTALL_DIR="$HOME/.local/share/archriot/install"
readonly LOG_FILE="$HOME/.cache/archriot/install.log"
readonly ERROR_LOG="$HOME/.cache/archriot/install-errors.log"
readonly BACKUP_DIR="$HOME/.config-backup-$(date +%Y%m%d-%H%M%S)"

# Read version from VERSION file
if [ -f "$SCRIPT_DIR/VERSION" ]; then
    readonly ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
elif [ -f "$HOME/.local/share/archriot/VERSION" ]; then
    readonly ARCHRIOT_VERSION=$(cat "$HOME/.local/share/archriot/VERSION" 2>/dev/null || echo "unknown")
else
    readonly ARCHRIOT_VERSION=$(curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION 2>/dev/null || echo "unknown")
fi

# Module ordering priorities (lower number = earlier execution)
declare -A MODULE_PRIORITIES=(
    ["core"]=10
    ["system"]=20
    ["development"]=30
    ["desktop"]=40
    ["post-desktop"]=45
    ["applications"]=50
    ["optional"]=60
)

# Get all installation modules dynamically with proper ordering
get_installation_modules() {
    local modules=()

    # Scan for individual core files first (numbered for ordering)
    for file in "$INSTALL_DIR/core"/*.sh; do
        [[ -f "$file" ]] && modules+=("core/$(basename "$file")")
    done

    # Then scan directory modules by priority
    for dir in "$INSTALL_DIR"/*; do
        if [[ -d "$dir" && "$(basename "$dir")" != "core" && "$(basename "$dir")" != "lib" ]]; then
            modules+=("$(basename "$dir")")
        fi
    done

    printf '%s\n' "${modules[@]}"
}

# Get installation modules in correct order
declare -a INSTALL_MODULES=($(get_installation_modules))

# Standalone installers
declare -a STANDALONE_INSTALLERS=(
    "plymouth.sh"            # Boot splash screen
)

# Performance tracking
INSTALL_START_TIME=$(date +%s)
INSTALL_START_DATE=$(date)
FAILED_MODULES=()
SKIPPED_MODULES=()

# ================================================================================
# Logging and Output Functions
# ================================================================================

# Initialize comprehensive logging
init_logging() {
    mkdir -p "$(dirname "$LOG_FILE")"
    mkdir -p "$(dirname "$ERROR_LOG")"

    # Clear previous logs
    echo "=== ArchRiot Installation v$ARCHRIOT_VERSION - $(date) ===" > "$LOG_FILE"
    echo "=== ArchRiot Installation Errors v$ARCHRIOT_VERSION - $(date) ===" > "$ERROR_LOG"

    echo ""
    echo "🚀 ArchRiot Installation System v2.0"
    echo "====================================="
    echo "Version: $ARCHRIOT_VERSION"
    echo "Start time: $INSTALL_START_DATE"
    echo "📝 Installation log: $LOG_FILE"
    echo "❌ Error log: $ERROR_LOG"
    echo ""
}

# Log function with both display and file output
log_message() {
    local level="$1"
    local message="$2"
    local timestamp="[$(date '+%H:%M:%S')]"

    case "$level" in
        "INFO")
            echo "$timestamp $message" | tee -a "$LOG_FILE"
            ;;
        "SUCCESS")
            echo "$timestamp ✅ $message" | tee -a "$LOG_FILE"
            ;;
        "WARNING")
            echo "$timestamp ⚠️  $message" | tee -a "$LOG_FILE"
            ;;
        "ERROR")
            echo "$timestamp ❌ $message" | tee -a "$LOG_FILE" | tee -a "$ERROR_LOG"
            ;;
        "CRITICAL")
            echo "$timestamp 🚨 CRITICAL: $message" | tee -a "$LOG_FILE" | tee -a "$ERROR_LOG"
            ;;
    esac
}

# ================================================================================
# System Preparation Functions
# ================================================================================

# Install essential tools immediately
install_essential_tools() {
    log_message "INFO" "Installing essential tools..."

    # Install base development tools
    if ! sudo pacman -Sy --noconfirm --needed base-devel git rsync bc; then
        log_message "CRITICAL" "Failed to install base development tools"
        return 1
    fi

    # Install yay if not present
    if ! command -v yay &>/dev/null; then
        log_message "INFO" "Installing yay AUR helper..."

        cd /tmp || return 1
        if ! git clone https://aur.archlinux.org/yay-bin.git; then
            log_message "CRITICAL" "Failed to clone yay-bin repository"
            return 1
        fi

        cd yay-bin || return 1
        if ! makepkg -si --noconfirm; then
            log_message "CRITICAL" "yay installation failed"
            return 1
        fi

        cd /
        rm -rf /tmp/yay-bin

        # Verify yay installation
        if ! command -v yay &>/dev/null; then
            log_message "CRITICAL" "yay not available after installation"
            return 1
        fi

        log_message "SUCCESS" "yay AUR helper installed"
    else
        log_message "SUCCESS" "yay AUR helper already available"
    fi

    return 0
}

# Setup sudo configuration
setup_sudo() {
    log_message "INFO" "Setting up sudo configuration..."

    if [ -f "$INSTALL_DIR/lib/sudo-helper.sh" ]; then
        source "$INSTALL_DIR/lib/sudo-helper.sh"

        if setup_passwordless_sudo; then
            if validate_passwordless_sudo; then
                log_message "SUCCESS" "Passwordless sudo configured"
                return 0
            else
                log_message "WARNING" "Passwordless sudo validation failed"
            fi
        else
            log_message "WARNING" "Failed to setup passwordless sudo"
        fi
    else
        log_message "WARNING" "Sudo helper not found"
    fi

    log_message "INFO" "Installation will prompt for passwords when needed"
    return 0
}

# Handle Git credentials with gum (before module processing)
handle_git_credentials() {
    echo ""
    echo "🔍 Git Configuration (Optional)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Check for existing git credentials
    local existing_name=$(git config --global user.name 2>/dev/null || echo "")
    local existing_email=$(git config --global user.email 2>/dev/null || echo "")

    if [[ -n "$existing_name" || -n "$existing_email" ]]; then
        echo "🎉 GitHub credentials found!"
        echo ""

        # Calculate box width and format entries properly
        local box_width=60
        local name_display="${existing_name:-"(not set)"}"
        local email_display="${existing_email:-"(not set)"}"

        # Format lines with proper spacing
        local name_line=$(printf "│ Username: %-*s │" $((box_width-15)) "$name_display")
        local email_line=$(printf "│ Email:    %-*s │" $((box_width-15)) "$email_display")

        echo "┌─────────────────────────────────────────────────────────┐"
        echo "│                 📋 Current Git Config                   │"
        echo "├─────────────────────────────────────────────────────────┤"
        echo "$name_line"
        echo "$email_line"
        echo "└─────────────────────────────────────────────────────────┘"
        echo ""

        # Use gum confirm (works here since not piped)
        if command -v gum &>/dev/null; then
            if gum confirm "Would you like to use these credentials?"; then
                echo "✅ Using existing credentials!"
                export ARCHRIOT_USER_NAME="$existing_name"
                export ARCHRIOT_USER_EMAIL="$existing_email"
            else
                echo ""
                echo "💬 No problem! Let's set up new credentials..."
                echo ""
                export ARCHRIOT_USER_NAME=$(gum input --placeholder "Your full name for Git commits" --prompt "Name: ")
                export ARCHRIOT_USER_EMAIL=$(gum input --placeholder "Your email for Git commits" --prompt "Email: ")
            fi
        else
            echo -n "Would you like to use these credentials? [Y/n]: "
            read -r response
            case "$response" in
                [nN][oO]|[nN])
                    echo ""
                    echo "💬 No problem! Let's set up new credentials..."
                    echo ""
                    echo -n "Name (Your full name for Git commits): "
                    read -r ARCHRIOT_USER_NAME
                    echo -n "Email (Your email for Git commits): "
                    read -r ARCHRIOT_USER_EMAIL
                    export ARCHRIOT_USER_NAME ARCHRIOT_USER_EMAIL
                    ;;
                *)
                    echo ""
                    echo "✅ Using existing credentials!"
                    export ARCHRIOT_USER_NAME="$existing_name"
                    export ARCHRIOT_USER_EMAIL="$existing_email"
                    ;;
            esac
        fi
    else
        echo "Configure Git with your name and email for commits and development."
        echo "This is optional - you can skip by pressing Enter or configure later."
        echo ""

        if command -v gum &>/dev/null; then
            export ARCHRIOT_USER_NAME=$(gum input --placeholder "Your full name for Git commits (optional)" --prompt "Name: ")
            export ARCHRIOT_USER_EMAIL=$(gum input --placeholder "Your email for Git commits (optional)" --prompt "Email: ")
        else
            echo -n "Name (optional): "
            read -r ARCHRIOT_USER_NAME
            echo -n "Email (optional): "
            read -r ARCHRIOT_USER_EMAIL
            export ARCHRIOT_USER_NAME ARCHRIOT_USER_EMAIL
        fi
    fi

    # Persist for modules to use
    local env_file="$HOME/.config/archriot/user.env"
    mkdir -p "$(dirname "$env_file")"
    {
        echo "ARCHRIOT_USER_NAME='$ARCHRIOT_USER_NAME'"
        echo "ARCHRIOT_USER_EMAIL='$ARCHRIOT_USER_EMAIL'"
    } > "$env_file"

    if [[ -n "$ARCHRIOT_USER_NAME" || -n "$ARCHRIOT_USER_EMAIL" ]]; then
        echo "✓ Git identity configured: ${ARCHRIOT_USER_NAME:-"(no name)"} <${ARCHRIOT_USER_EMAIL:-"(no email)"}>"
    else
        echo "⚠ Git identity skipped - you can configure later if needed"
    fi
    echo ""
}

# ================================================================================
# Module Execution Functions
# ================================================================================

# Execute a single module with comprehensive error handling
execute_module() {
    local module_path="$1"
    local module_name="$(basename "$module_path" .sh)"
    local start_time=$(date +%s)

    if [[ ! -f "$module_path" ]]; then
        log_message "ERROR" "Module not found: $module_path"
        FAILED_MODULES+=("$module_name (not found)")
        return 1
    fi

    log_message "INFO" "🔄 Executing module: $module_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Execute module directly with full output visibility
    local temp_log=$(mktemp)
    if bash "$module_path" 2>&1 | tee "$temp_log"; then
        # Append module output to main log
        cat "$temp_log" >> "$LOG_FILE"
        rm -f "$temp_log"

        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        log_message "SUCCESS" "$module_name completed (${duration}s)"
        echo ""
        return 0
    else
        local exit_code=$?

        # Log the failure details
        echo "Module output:" >> "$ERROR_LOG"
        cat "$temp_log" >> "$ERROR_LOG"
        cat "$temp_log" >> "$LOG_FILE"
        rm -f "$temp_log"

        log_message "ERROR" "$module_name failed (exit code: $exit_code)"
        FAILED_MODULES+=("$module_name")
        echo ""

        # For critical modules, ask user how to proceed
        if [[ "$module_name" == *"identity"* ]] || [[ "$module_name" == *"hyprland"* ]] || [[ "$module_name" == *"theming"* ]]; then
            echo "🚨 CRITICAL MODULE FAILED: $module_name"
            echo "This module is essential for ArchRiot functionality."
            echo ""
            echo "Options:"
            echo "1. Continue installation (may result in broken system)"
            echo "2. Retry this module"
            echo "3. Abort installation"
            echo ""

            read -p "Choose [1/2/3]: " choice
            case "$choice" in
                2)
                    log_message "INFO" "Retrying $module_name..."
                    return $(execute_module "$module_path")
                    ;;
                3)
                    log_message "CRITICAL" "Installation aborted by user"
                    exit 1
                    ;;
                *)
                    log_message "WARNING" "Continuing despite critical module failure"
                    ;;
            esac
        fi

        return $exit_code
    fi
}

# Execute all modules in a directory
execute_module_directory() {
    local module_dir="$1"
    local directory_name="$(basename "$module_dir")"

    if [[ ! -d "$module_dir" ]]; then
        log_message "WARNING" "Module directory not found: $module_dir"
        SKIPPED_MODULES+=("$directory_name (directory not found)")
        return 0
    fi

    log_message "INFO" "📁 Processing module directory: $directory_name"

    local module_count=0
    local failed_count=0

    # Execute all .sh files in the directory
    for module_file in "$module_dir"/*.sh; do
        if [[ -f "$module_file" ]]; then
            ((module_count++))
            if ! execute_module "$module_file"; then
                ((failed_count++))
            fi
        fi
    done

    if [[ $module_count -eq 0 ]]; then
        log_message "WARNING" "No modules found in directory: $directory_name"
        SKIPPED_MODULES+=("$directory_name (no modules)")
    elif [[ $failed_count -eq 0 ]]; then
        log_message "SUCCESS" "All modules in $directory_name completed successfully"
    else
        log_message "WARNING" "$failed_count of $module_count modules failed in $directory_name"
    fi

    return 0
}

# Process all installation modules
process_installation_modules() {
    log_message "INFO" "🚀 Starting module processing..."

    for module in "${INSTALL_MODULES[@]}"; do
        if [[ "$module" == *".sh" ]]; then
            # Individual module file
            local module_path="$INSTALL_DIR/$module"
            execute_module "$module_path"
        else
            # Module directory
            local module_dir="$INSTALL_DIR/$module"
            execute_module_directory "$module_dir"
        fi
    done

    # Execute standalone installers
    for standalone in "${STANDALONE_INSTALLERS[@]}"; do
        local standalone_path="$INSTALL_DIR/$standalone"
        if [[ -f "$standalone_path" ]]; then
            execute_module "$standalone_path"
        else
            log_message "WARNING" "Standalone installer not found: $standalone"
            SKIPPED_MODULES+=("$standalone")
        fi
    done
}

# ================================================================================
# Verification Functions
# ================================================================================

# Verify critical installations
verify_installation() {
    log_message "INFO" "🔍 Verifying critical installations..."

    # Small delay to ensure all file operations complete during upgrades
    sleep 2

    local failures=0

    # Verify yay
    if ! command -v yay &>/dev/null; then
        log_message "ERROR" "yay AUR helper not found"
        ((failures++))
    fi

    # Verify Hyprland
    if ! command -v hyprland &>/dev/null; then
        log_message "ERROR" "Hyprland not installed"
        ((failures++))
    fi

    # Verify theme system (consolidated structure) with detailed logging
    log_message "INFO" "Checking theme system configuration..."
    local theme_issues=0

    # Check archriot.conf
    if [[ ! -f "$HOME/.config/archriot/archriot.conf" ]]; then
        log_message "ERROR" "Missing: ~/.config/archriot/archriot.conf"
        ((theme_issues++))
    else
        log_message "INFO" "✓ archriot.conf found"
    fi

    # Check backgrounds directory
    if [[ ! -d "$HOME/.config/archriot/backgrounds" ]]; then
        log_message "ERROR" "Missing: ~/.config/archriot/backgrounds/"
        ((theme_issues++))
    else
        local bg_count=$(find "$HOME/.config/archriot/backgrounds" -type f \( -name "*.jpg" -o -name "*.png" \) 2>/dev/null | wc -l)
        if [[ $bg_count -gt 0 ]]; then
            log_message "INFO" "✓ backgrounds directory found ($bg_count files)"
        else
            log_message "ERROR" "backgrounds directory exists but contains no image files"
            ((theme_issues++))
        fi
    fi

    # Only attempt fix if there are actual issues
    if [[ $theme_issues -gt 0 ]]; then
        log_message "ERROR" "Theme system has $theme_issues issue(s)"
        ((failures++))

        # Attempt to fix theme system
        log_message "INFO" "Attempting to fix theme system..."
        local theming_script="$INSTALL_DIR/desktop/theming.sh"
        if [[ -f "$theming_script" ]]; then
            if bash "$theming_script" 2>&1 | tee -a "$LOG_FILE"; then
                log_message "SUCCESS" "Theme system fixed"
                ((failures--))
            else
                log_message "ERROR" "Failed to fix theme system"
            fi
        fi
    else
        log_message "INFO" "✓ Theme system properly configured"
    fi

    # Verify gum (required for UI)
    if ! command -v gum &>/dev/null; then
        log_message "WARNING" "gum not found, installing..."
        if yay -S --noconfirm --needed gum; then
            log_message "SUCCESS" "gum installed"
        else
            log_message "ERROR" "Failed to install gum"
            ((failures++))
        fi
    fi

    # Verify waybar
    if ! command -v waybar &>/dev/null; then
        log_message "WARNING" "waybar not found"
    fi

    if [[ $failures -eq 0 ]]; then
        log_message "SUCCESS" "All critical components verified"
        return 0
    else
        log_message "WARNING" "$failures critical verification failures detected"
        return 1
    fi
}

# ================================================================================
# Cleanup and Finalization Functions
# ================================================================================

# System updates and cleanup
finalize_installation() {
    log_message "INFO" "🔧 Finalizing installation..."

    # Update locate database
    log_message "INFO" "Updating locate database..."
    sudo updatedb || log_message "WARNING" "Failed to update locate database"

    # Update font cache
    log_message "INFO" "Updating font cache..."
    fc-cache -fv >/dev/null 2>&1 || log_message "WARNING" "Failed to update font cache"

    # Update icon cache
    log_message "INFO" "Updating icon cache..."
    gtk-update-icon-cache -f ~/.local/share/icons/hicolor/ 2>/dev/null || true
    update-desktop-database ~/.local/share/applications/ 2>/dev/null || true

    # Reload running services
    log_message "INFO" "Reloading system services..."
    if pgrep -x "Hyprland" >/dev/null; then
        hyprctl reload 2>/dev/null || true
    fi

    # Configure thumbnails (disable PDF thumbnails while keeping others)
    log_message "INFO" "Configuring thumbnails (disable PDF thumbnails)..."
    if [[ -f "$INSTALL_DIR/bin/fix-thunar-thumbnails" ]]; then
        if bash "$INSTALL_DIR/bin/fix-thunar-thumbnails" 2>&1 | tee -a "$LOG_FILE"; then
            log_message "SUCCESS" "Thumbnail configuration completed"
        else
            log_message "WARNING" "Thumbnail configuration had issues but continuing"
        fi
    else
        log_message "WARNING" "Thumbnail fix script not found"
    fi

    # Ensure waybar is running after installation (only in Wayland environment)
    if [[ -n "$WAYLAND_DISPLAY" ]] || pgrep -x "Hyprland" >/dev/null; then
        if pgrep -x "waybar" >/dev/null; then
            log_message "INFO" "Restarting waybar..."
            pkill waybar 2>/dev/null || true
            sleep 1
        else
            log_message "INFO" "Starting waybar..."
        fi
        nohup waybar &>/dev/null & disown
    else
        log_message "INFO" "Skipping waybar start (no Wayland environment detected)"
    fi

    # Update version file
    if [[ -n "$ARCHRIOT_VERSION" && "$ARCHRIOT_VERSION" != "unknown" ]]; then
        mkdir -p "$HOME/.local/share/archriot"
        echo "$ARCHRIOT_VERSION" > "$HOME/.local/share/archriot/VERSION"
    fi

    log_message "SUCCESS" "Installation finalized"
}

# Clean up sudo configuration
cleanup_sudo() {
    if [ -f "$INSTALL_DIR/lib/sudo-helper.sh" ]; then
        source "$INSTALL_DIR/lib/sudo-helper.sh"
        if command -v remove_passwordless_rule &>/dev/null; then
            log_message "INFO" "Cleaning up sudo configuration..."
            remove_passwordless_rule || true
        fi
    fi
}

# Show installation summary
show_installation_summary() {
    local end_time=$(date +%s)
    local duration=$((end_time - INSTALL_START_TIME))
    local duration_min=$((duration / 60))
    local duration_sec=$((duration % 60))

    echo ""
    echo "==============================================="
    echo "🎉 ArchRiot v$ARCHRIOT_VERSION Installation Complete!"
    echo "==============================================="
    echo "⏱️  Total time: ${duration_min}m ${duration_sec}s"
    echo "📁 Installation log: $LOG_FILE"

    if [[ ${#FAILED_MODULES[@]} -gt 0 ]]; then
        echo ""
        echo "⚠️  Failed modules (${#FAILED_MODULES[@]}):"
        for module in "${FAILED_MODULES[@]}"; do
            echo "   ❌ $module"
        done
        echo ""
        echo "📝 Check error log for details: $ERROR_LOG"
        echo "🔄 To retry failed modules: source ~/.local/share/archriot/install.sh"
    fi

    if [[ ${#SKIPPED_MODULES[@]} -gt 0 ]]; then
        echo ""
        echo "ℹ️  Skipped modules (${#SKIPPED_MODULES[@]}):"
        for module in "${SKIPPED_MODULES[@]}"; do
            echo "   ⏭️  $module"
        done
    fi

    echo ""
    echo "🎯 Quick Commands:"
    echo "  • Launch Apps: Super + D"
    echo "  • Backgrounds: Super + Ctrl + Space"
    echo "  • View Help: Super + H"
    echo ""

    # Show backup location if created
    if [[ -d "$BACKUP_DIR" ]]; then
        echo "📦 Previous configuration backed up to:"
        echo "   $BACKUP_DIR"
        echo ""
    fi
}

# Error handler for cleanup
handle_installation_error() {
    log_message "CRITICAL" "Installation failed unexpectedly!"
    log_message "INFO" "Performing cleanup..."

    cleanup_sudo

    echo ""
    echo "❌ ArchRiot installation failed!"
    echo "📝 Check logs for details:"
    echo "   Main log: $LOG_FILE"
    echo "   Error log: $ERROR_LOG"
    echo ""
    echo "🔄 To retry: source $HOME/.local/share/archriot/install.sh"
    echo "💡 Most components are idempotent - re-running will skip installed items"

    exit 1
}

# ================================================================================
# Main Installation Flow
# ================================================================================

main() {
    # Set up error handling
    trap handle_installation_error ERR

    # Initialize logging
    init_logging

    # Essential system preparation
    log_message "INFO" "🔧 Preparing system..."
    install_essential_tools || {
        log_message "CRITICAL" "Failed to install essential tools"
        exit 1
    }

    setup_sudo

    # Load specific helper libraries if available (avoid shared.sh which has dependencies)
    for helper_name in "install-helpers.sh" "simple-progress.sh"; do
        local helper_file="$INSTALL_DIR/lib/$helper_name"
        if [[ -f "$helper_file" ]]; then
            source "$helper_file" 2>/dev/null || {
                log_message "WARNING" "Failed to load $helper_name helper"
            }
        fi
    done

    # Handle Git credentials (before modules so gum works properly)
    handle_git_credentials

    # Process all installation modules
    process_installation_modules

    # Verify critical installations
    verify_installation

    # Finalize installation
    finalize_installation

    # Clean up
    cleanup_sudo

    # Show summary
    show_installation_summary

    # Final check for major failures
    if [[ ${#FAILED_MODULES[@]} -gt 0 ]]; then
        log_message "WARNING" "Installation completed with some failures"

        # Check if any critical modules failed
        local critical_failed=false
        for module in "${FAILED_MODULES[@]}"; do
            if [[ "$module" == *"identity"* ]] || [[ "$module" == *"hyprland"* ]] || [[ "$module" == *"theming"* ]]; then
                critical_failed=true
                break
            fi
        done

        if [[ "$critical_failed" == "true" ]]; then
            echo ""
            echo "🚨 CRITICAL: Essential modules failed!"
            echo "Your ArchRiot installation may not function properly."
            echo "Please review the error log and retry installation."
            echo ""
            read -p "Press Enter to continue..."
            return 1
        fi
    else
        log_message "SUCCESS" "All modules completed successfully!"
    fi

    echo "✅ System ready! New terminals will have updated configs."
    echo ""

    # Optional reboot prompt
    if command -v gum &>/dev/null; then
        if gum confirm "Reboot to ensure all settings are fully applied?"; then
            reboot
        fi
    else
        read -p "Reboot to ensure all settings are fully applied? [y/N]: " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            reboot
        fi
    fi
}

# Execute main function
main "$@"
