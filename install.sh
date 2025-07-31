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
readonly ARCHRIOT_DATA_DIR="$HOME/.local/share/archriot"
readonly ARCHRIOT_ICONS_DIR="$HOME/.local/share/icons/hicolor"
readonly ARCHRIOT_APPS_DIR="$HOME/.local/share/applications"
readonly ARCHRIOT_BACKUP_PREFIX="$HOME/.local/share/archriot-backup-"
readonly TEMP_BUILD_DIR="/tmp"
# Backup directory now handled by centralized backup system

# Read version from VERSION file
if [ -f "$SCRIPT_DIR/VERSION" ]; then
    readonly ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
elif [ -f "$ARCHRIOT_DATA_DIR/VERSION" ]; then
    readonly ARCHRIOT_VERSION=$(cat "$ARCHRIOT_DATA_DIR/VERSION" 2>/dev/null || echo "unknown")
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

# Execute all installation modules with direct execution approach
execute_installation_modules() {
    echo "INFO: 🚀 Starting module processing..." >> "$LOG_FILE" 2>&1

    # Core modules first (numbered order)
    for file in "$INSTALL_DIR/core"/*.sh; do
        [[ -f "$file" ]] && execute_module "$file"
    done

    # Directory modules in specific order
    for dir in "$INSTALL_DIR"/{system,desktop,applications,development,post-desktop,optional}; do
        [[ -d "$dir" ]] && execute_module_directory "$dir"
    done

    # Execute standalone installers
    for standalone in "${STANDALONE_INSTALLERS[@]}"; do
        local standalone_path="$INSTALL_DIR/$standalone"
        if [[ -f "$standalone_path" ]]; then
            execute_module "$standalone_path"
        else
            echo "WARNING: Standalone installer not found: $standalone" >> "$LOG_FILE" 2>&1
            SKIPPED_MODULES+=("$standalone")
        fi
    done
}

# Standalone installers
declare -a STANDALONE_INSTALLERS=(
    "plymouth.sh"            # Boot splash screen
)

# Performance tracking
export INSTALL_START_TIME=$(date +%s)
export INSTALL_START_DATE=$(date)
FAILED_MODULES=()
SKIPPED_MODULES=()

# ================================================================================
# Logging and Output Functions
# ================================================================================

# Load progress bar system
if [[ -f "$INSTALL_DIR/lib/progress-bar.sh" ]]; then
    source "$INSTALL_DIR/lib/progress-bar.sh"
else
    # Fallback if progress system not available
    init_progress_system() { echo "🚀 ArchRiot Installation v$ARCHRIOT_VERSION"; }
    progress_module_start() { echo "Starting: $(basename "$1" .sh)"; }
    progress_module_complete() { echo "Completed: $(basename "$1" .sh)"; }
    progress_pause_for_input() { echo ""; }
    progress_resume_after_input() { echo ""; }
    progress_show_completion() { echo "Installation complete!"; }
fi

# Load centralized backup system
if [[ -f "$INSTALL_DIR/lib/backup-manager.sh" ]]; then
    source "$INSTALL_DIR/lib/backup-manager.sh"
else
    # Fallback backup functions
    backup_configs() { echo "INFO: Backup system not available" >> "${LOG_FILE:-/dev/null}" 2>&1; }
    emergency_backup() { echo "INFO: Emergency backup not available" >> "${LOG_FILE:-/dev/null}" 2>&1; }
fi

# Initialize comprehensive logging
init_logging() {
    mkdir -p "$(dirname "$LOG_FILE")"
    mkdir -p "$(dirname "$ERROR_LOG")"

    # Clear previous logs
    echo "=== ArchRiot Installation v$ARCHRIOT_VERSION - $(date) ===" > "$LOG_FILE"
    echo "=== ArchRiot Installation Errors v$ARCHRIOT_VERSION - $(date) ===" > "$ERROR_LOG"

    # Initialize modern progress bar system
    init_progress_system
}

# Log function with both display and file output
log_message() {
    local level="$1"
    local message="$2"
    local timestamp="[$(date '+%H:%M:%S')]"
    local formatted_message

    case "$level" in
        "INFO")
            formatted_message="$timestamp $message"
            ;;
        "SUCCESS")
            formatted_message="$timestamp ✅ $message"
            ;;
        "WARNING")
            formatted_message="$timestamp ⚠️  $message"
            ;;
        "ERROR")
            formatted_message="$timestamp ❌ $message"
            echo "$formatted_message" >> "$ERROR_LOG"
            ;;
        "CRITICAL")
            formatted_message="$timestamp 🚨 CRITICAL: $message"
            echo "$formatted_message" >> "$ERROR_LOG"
            ;;
    esac

    # Single output operation
    echo "$formatted_message" >> "$LOG_FILE"
}

# ================================================================================
# System Preparation Functions
# ================================================================================

# Install essential tools immediately
# Install base development tools and yay AUR helper
install_essential_tools() {
    echo "Installing essential tools..." >> "$LOG_FILE" 2>&1

    # Install base development tools (silent)
    if ! sudo pacman -Sy --noconfirm --needed base-devel git rsync bc >> "$LOG_FILE" 2>&1; then
        echo "CRITICAL: Failed to install base development tools" >> "$LOG_FILE" 2>&1
        return 1
    fi

    # Install yay if not present (silent)
    if ! command -v yay &>/dev/null; then
        echo "Installing yay AUR helper..." >> "$LOG_FILE" 2>&1

        cd "$TEMP_BUILD_DIR" || return 1
        if ! git clone https://aur.archlinux.org/yay-bin.git >> "$LOG_FILE" 2>&1; then
            echo "CRITICAL: Failed to clone yay-bin repository" >> "$LOG_FILE" 2>&1
            return 1
        fi

        cd yay-bin || return 1
        if ! makepkg -si --noconfirm >> "$LOG_FILE" 2>&1; then
            echo "CRITICAL: yay installation failed" >> "$LOG_FILE" 2>&1
            return 1
        fi

        cd /
        rm -rf "$TEMP_BUILD_DIR/yay-bin"

        # Verify yay installation
        if ! command -v yay &>/dev/null; then
            echo "CRITICAL: yay not available after installation" >> "$LOG_FILE" 2>&1
            return 1
        fi

        echo "SUCCESS: yay AUR helper installed" >> "$LOG_FILE" 2>&1
    else
        echo "SUCCESS: yay AUR helper already available" >> "$LOG_FILE" 2>&1
    fi

    return 0
}

# Setup sudo configuration
setup_sudo() {
    echo "INFO: Setting up sudo configuration..." >> "$LOG_FILE" 2>&1

    if [ -f "$INSTALL_DIR/lib/sudo-helper.sh" ]; then
        source "$INSTALL_DIR/lib/sudo-helper.sh"

        if setup_passwordless_sudo; then
            if validate_passwordless_sudo; then
                echo "SUCCESS: Passwordless sudo configured" >> "$LOG_FILE" 2>&1
                return 0
            else
                echo "WARNING: Passwordless sudo validation failed" >> "$LOG_FILE" 2>&1
            fi
        else
            echo "WARNING: Failed to setup passwordless sudo" >> "$LOG_FILE" 2>&1
        fi
    else
        echo "WARNING: Sudo helper not found" >> "$LOG_FILE" 2>&1
    fi

    echo "INFO: Installation will prompt for passwords when needed" >> "$LOG_FILE" 2>&1
    return 0
}

# Handle Git credentials with gum (before module processing)
handle_git_credentials() {
    # Pause progress display for user interaction
    progress_pause_for_input "🔍 Git Configuration (Optional)"
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
        echo "ERROR: Module not found: $module_path" >> "$LOG_FILE" 2>&1
        FAILED_MODULES+=("$module_name (not found)")
        return 1
    fi

    # Update progress display
    progress_module_start "$module_path"

    # Log to file only (verbose output hidden from console)
    echo "INFO: 🚀 Starting module: $module_name" >> "$LOG_FILE" 2>&1

    # Execute module with all output redirected to log file (silent execution)
    if bash "$module_path" >> "$LOG_FILE" 2>&1; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        echo "SUCCESS: $module_name completed in ${duration}s" >> "$LOG_FILE" 2>&1

        # Update progress display
        progress_module_complete "$module_path" "true"
        return 0
    else
        local exit_code=$?

        # Log failure details
        echo "Module: $module_name (exit code: $exit_code)" >> "$ERROR_LOG" 2>&1
        echo "See main log for details: $LOG_FILE" >> "$ERROR_LOG" 2>&1
        echo "ERROR: $module_name failed (exit code: $exit_code)" >> "$LOG_FILE" 2>&1

        # Update progress display and show error
        progress_module_complete "$module_path" "false"
        progress_show_error "Installation failed (exit code: $exit_code)" "$(get_friendly_name "$module_path")"

        FAILED_MODULES+=("$module_name")

        # For critical modules, fail the installation
        if [[ "$module_name" == *"identity"* ]] || [[ "$module_name" == *"hyprland"* ]] || [[ "$module_name" == *"theming"* ]]; then
            echo "🚨 CRITICAL MODULE FAILED: $module_name"
            echo "This module is essential for ArchRiot functionality."
            echo "Installation cannot continue with this failure."
            echo ""
            echo "Check the log file for details: $LOG_FILE"
            echo "Fix the issue and re-run the installation."
            echo ""
            echo "CRITICAL: Installation failed due to $module_name failure" >> "$LOG_FILE" 2>&1
            exit 1
        fi

        return $exit_code
    fi
}

# Execute all modules in a directory
execute_module_directory() {
    local module_dir="$1"
    local directory_name="$(basename "$module_dir")"

    if [[ ! -d "$module_dir" ]]; then
        echo "WARNING: Module directory not found: $module_dir" >> "$LOG_FILE" 2>&1
        SKIPPED_MODULES+=("$directory_name (directory not found)")
        return 0
    fi

    # Update progress display
    progress_module_start "$module_dir"
    echo "INFO: 📁 Processing module directory: $directory_name" >> "$LOG_FILE" 2>&1

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
        echo "WARNING: No modules found in directory: $directory_name" >> "$LOG_FILE" 2>&1
        SKIPPED_MODULES+=("$directory_name (no modules)")
        progress_module_complete "$module_dir" "false"
    elif [[ $failed_count -eq 0 ]]; then
        echo "SUCCESS: All modules in $directory_name completed successfully" >> "$LOG_FILE" 2>&1
        progress_module_complete "$module_dir" "true"
    else
        echo "WARNING: $failed_count of $module_count modules failed in $directory_name" >> "$LOG_FILE" 2>&1
        progress_module_complete "$module_dir" "false"
    fi

    return 0
}

# Removed - functionality moved to execute_installation_modules()

# ================================================================================
# Verification Functions
# ================================================================================

# Clean up old scattered backup directories from previous backup systems
cleanup_old_backups() {
    echo "INFO: 🧹 Cleaning up old backup directories..." >> "$LOG_FILE" 2>&1

    local total_removed=0

    # Remove old install.sh backups (~/.config-backup-*)
    for backup_dir in ~/.config-backup-*; do
        if [[ -d "$backup_dir" ]]; then
            rm -rf "$backup_dir" 2>/dev/null && ((total_removed++))
        fi
    done

    # Remove old setup.sh backups (archriot-backup-*)
    for backup_dir in "${ARCHRIOT_BACKUP_PREFIX}"*; do
        if [[ -d "$backup_dir" ]]; then
            rm -rf "$backup_dir" 2>/dev/null && ((total_removed++))
        fi
    done

    # Remove old theming.sh backups (archriot-backups/*)
    if [[ -d ~/.config/archriot-backups ]]; then
        for backup_dir in ~/.config/archriot-backups/*; do
            if [[ -d "$backup_dir" ]]; then
                rm -rf "$backup_dir" 2>/dev/null && ((total_removed++))
            fi
        done
        rmdir ~/.config/archriot-backups 2>/dev/null
    fi

    # Remove old bin backups (.local/bin.backup-*)
    for backup_dir in ~/.local/bin.backup-*; do
        if [[ -d "$backup_dir" ]]; then
            rm -rf "$backup_dir" 2>/dev/null && ((total_removed++))
        fi
    done

    if [[ $total_removed -gt 0 ]]; then
        echo "INFO: ✓ Cleaned up $total_removed old backup directories" >> "$LOG_FILE" 2>&1
    fi
}

# Verify critical installations
verify_installation() {
    echo "INFO: 🔍 Verifying critical installations..." >> "$LOG_FILE" 2>&1

    # Small delay to ensure all file operations complete during upgrades
    sleep 2

    local failures=0

    # Verify yay
    if ! command -v yay &>/dev/null; then
        echo "ERROR: yay AUR helper not found" >> "$LOG_FILE" 2>&1
        ((failures++))
    fi

    # Verify Hyprland
    if ! command -v hyprland &>/dev/null; then
        echo "ERROR: Hyprland not installed" >> "$LOG_FILE" 2>&1
        ((failures++))
    fi

    # Verify theme system (consolidated structure) with detailed logging
    echo "INFO: Checking theme system configuration..." >> "$LOG_FILE" 2>&1
    local theme_issues=0

    # Check archriot.conf
    if [[ ! -f "$HOME/.config/archriot/archriot.conf" ]]; then
        echo "ERROR: Missing: ~/.config/archriot/archriot.conf" >> "$LOG_FILE" 2>&1
        ((theme_issues++))
    else
        echo "INFO: ✓ archriot.conf found" >> "$LOG_FILE" 2>&1
    fi

    # Check backgrounds directory
    if [[ ! -d "$HOME/.config/archriot/backgrounds" ]]; then
        echo "ERROR: Missing: ~/.config/archriot/backgrounds/" >> "$LOG_FILE" 2>&1
        ((theme_issues++))
    else
        local bg_count=$(find "$HOME/.config/archriot/backgrounds" -type f \( -name "*.jpg" -o -name "*.png" \) 2>/dev/null | wc -l)
        if [[ $bg_count -gt 0 ]]; then
            echo "INFO: ✓ backgrounds directory found ($bg_count files)" >> "$LOG_FILE" 2>&1
        else
            echo "ERROR: backgrounds directory exists but contains no image files" >> "$LOG_FILE" 2>&1
            ((theme_issues++))
        fi
    fi

    # Only attempt fix if there are actual issues
    if [[ $theme_issues -gt 0 ]]; then
        echo "ERROR: Theme system has $theme_issues issue(s)" >> "$LOG_FILE" 2>&1
        ((failures++))

        # Attempt to fix theme system
        echo "INFO: Attempting to fix theme system..." >> "$LOG_FILE" 2>&1
        local theming_script="$INSTALL_DIR/desktop/theming.sh"
        if [[ -f "$theming_script" ]]; then
            if bash "$theming_script" >> "$LOG_FILE" 2>&1; then
                echo "SUCCESS: Theme system fixed" >> "$LOG_FILE" 2>&1
                ((failures--))
            else
                echo "ERROR: Failed to fix theme system" >> "$LOG_FILE" 2>&1
            fi
        fi
    else
        echo "INFO: ✓ Theme system properly configured" >> "$LOG_FILE" 2>&1
    fi

    # Verify gum (required for UI)
    if ! command -v gum &>/dev/null; then
        echo "WARNING: gum not found, installing..." >> "$LOG_FILE" 2>&1
        if yay -S --noconfirm --needed gum >> "$LOG_FILE" 2>&1; then
            echo "SUCCESS: gum installed" >> "$LOG_FILE" 2>&1
        else
            echo "ERROR: Failed to install gum" >> "$LOG_FILE" 2>&1
            ((failures++))
        fi
    fi

    # Verify waybar
    if ! command -v waybar &>/dev/null; then
        echo "WARNING: waybar not found" >> "$LOG_FILE" 2>&1
    fi

    if [[ $failures -eq 0 ]]; then
        echo "SUCCESS: All critical components verified" >> "$LOG_FILE" 2>&1
        return 0
    else
        echo "WARNING: $failures critical verification failures detected" >> "$LOG_FILE" 2>&1
        return 1
    fi
}

# ================================================================================
# Cleanup and Finalization Functions
# ================================================================================

# System updates and cleanup
finalize_installation() {
    echo "INFO: 🔧 Finalizing installation..." >> "$LOG_FILE" 2>&1

    # Update locate database
    echo "INFO: Updating locate database..." >> "$LOG_FILE" 2>&1
    sudo updatedb >> "$LOG_FILE" 2>&1 || echo "WARNING: Failed to update locate database" >> "$LOG_FILE" 2>&1

    # Update font cache
    echo "INFO: Updating font cache..." >> "$LOG_FILE" 2>&1
    fc-cache -fv >> "$LOG_FILE" 2>&1 || echo "WARNING: Failed to update font cache" >> "$LOG_FILE" 2>&1

    # Update icon cache
    echo "INFO: Updating icon cache..." >> "$LOG_FILE" 2>&1
    gtk-update-icon-cache -f "$ARCHRIOT_ICONS_DIR/" 2>/dev/null || true
    update-desktop-database "$ARCHRIOT_APPS_DIR/" 2>/dev/null || true

    # Services will be restarted at the end of installation

    # Configure thumbnails (disable PDF thumbnails while keeping others)
    echo "INFO: Configuring thumbnails (disable PDF thumbnails)..." >> "$LOG_FILE" 2>&1
    if [[ -f "$INSTALL_DIR/bin/fix-thunar-thumbnails" ]]; then
        if bash "$INSTALL_DIR/bin/fix-thunar-thumbnails" >> "$LOG_FILE" 2>&1; then
            echo "SUCCESS: Thumbnail configuration completed" >> "$LOG_FILE" 2>&1
        else
            echo "WARNING: Thumbnail configuration had issues but continuing" >> "$LOG_FILE" 2>&1
        fi
    else
        echo "WARNING: Thumbnail fix script not found" >> "$LOG_FILE" 2>&1
    fi

    # Services will be restarted in final restart function

    # Update version file
    if [[ -n "$ARCHRIOT_VERSION" && "$ARCHRIOT_VERSION" != "unknown" ]]; then
        mkdir -p "$ARCHRIOT_DATA_DIR"
        echo "$ARCHRIOT_VERSION" > "$ARCHRIOT_DATA_DIR/VERSION"
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

# Show installation summary (now handled by progress system)
show_installation_summary() {
    # Use the beautiful progress system completion summary
    progress_show_completion
}

# Final service restart - consolidate all restarts into one operation
final_service_restart() {
    echo "INFO: 🔄 Final service restart..." >> "$LOG_FILE" 2>&1

    # Only restart if we're in a Wayland environment
    if [[ -n "$WAYLAND_DISPLAY" ]] || pgrep -x "Hyprland" >/dev/null; then

        # Restart Mako notifications with new configs
        if pgrep -x "mako" >/dev/null; then
            echo "INFO: Restarting mako notifications..." >> "$LOG_FILE" 2>&1
            pkill mako 2>/dev/null || true
            sleep 1
        fi
        if command -v mako &>/dev/null; then
            nohup mako &>/dev/null & disown
        fi

        # Restart Waybar with new configs
        if pgrep -x "waybar" >/dev/null; then
            echo "INFO: Restarting waybar..." >> "$LOG_FILE" 2>&1
            pkill waybar 2>/dev/null || true
            sleep 1
            if command -v waybar &>/dev/null; then
                nohup waybar &>/dev/null & disown
            fi
        elif command -v waybar &>/dev/null; then
            echo "INFO: Starting waybar..." >> "$LOG_FILE" 2>&1
            nohup waybar &>/dev/null & disown
        fi

        # Reload Hyprland configs (no restart needed)
        if pgrep -x "Hyprland" >/dev/null; then
            echo "INFO: Reloading Hyprland configuration..." >> "$LOG_FILE" 2>&1
            hyprctl reload 2>/dev/null || true
        fi

        echo "INFO: ✓ Desktop services restarted with new configurations" >> "$LOG_FILE" 2>&1
    else
        echo "INFO: Skipping service restart (no Wayland environment detected)" >> "$LOG_FILE" 2>&1
    fi
}

# Show final tips and commands
show_final_tips() {
    echo "🎯 Quick Commands:"
    echo "  • Launch Apps: Super + D"
    echo "  • Backgrounds: Super + Ctrl + Space"
    echo "  • View Help: Super + H"
    echo ""

    # Show backup location if centralized backups exist
    if [[ -d "$HOME/.archriot/backups" ]] && [[ -n "$(ls -A "$HOME/.archriot/backups" 2>/dev/null)" ]]; then
        echo "📦 Configuration backups available at:"
        echo "   ~/.archriot/backups/ (keeping 3 most recent)"
        echo ""
    fi
}

# Error handler for cleanup
handle_installation_error() {
    echo "CRITICAL: Installation failed unexpectedly!" >> "$LOG_FILE" 2>&1
    echo "INFO: Performing cleanup..." >> "$LOG_FILE" 2>&1

    cleanup_sudo

    echo ""
    echo "❌ ArchRiot installation failed!"
    echo "📝 Check logs for details:"
    echo "   Main log: $LOG_FILE"
    echo "   Error log: $ERROR_LOG"
    echo ""
    echo "🔄 To retry: source $ARCHRIOT_DATA_DIR/install.sh"
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
    echo "INFO: 🔧 Preparing system..." >> "$LOG_FILE" 2>&1
    install_essential_tools || {
        echo "CRITICAL: Failed to install essential tools" >> "$LOG_FILE" 2>&1
        exit 1
    }

    # Clean up old backup directories from previous systems
    cleanup_old_backups

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

    # Resume progress display after user interaction
    progress_resume_after_input

    # Process all installation modules
    execute_installation_modules

    # Verify critical installations
    verify_installation

    # Finalize installation
    finalize_installation

    # Clean up
    cleanup_sudo

    # Final service restart - only once at the very end
    final_service_restart

    # Show summary
    show_installation_summary
    show_final_tips

    # Final check for major failures
    if [[ ${#FAILED_MODULES[@]} -gt 0 ]]; then
        echo "WARNING: Installation completed with some failures" >> "$LOG_FILE" 2>&1

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

    # Pause progress for reboot prompt
    progress_pause_for_input "🔄 System Reboot"

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
