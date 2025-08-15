#!/bin/bash

# ================================================================================
# Mullvad GUI Startup Script for ArchRiot
# ================================================================================
# Conditionally starts Mullvad GUI minimized to tray if account is logged in
# Handles startup delays and ensures proper tray integration
# ================================================================================

# Colors for output (when needed)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log function for debugging
log_message() {
    local message="$1"
    echo "[$(date '+%H:%M:%S')] $message" >> "$HOME/.cache/mullvad-startup.log"
}

# Check if Mullvad is installed
check_mullvad_installed() {
    if ! command -v mullvad &>/dev/null; then
        log_message "Mullvad CLI not found, skipping GUI startup"
        return 1
    fi

    if ! [ -f "/opt/Mullvad VPN/mullvad-gui" ]; then
        log_message "Mullvad GUI not found, skipping startup"
        return 1
    fi

    return 0
}

# Check if Mullvad account is logged in
check_mullvad_account() {
    local status_output=$(mullvad status 2>/dev/null)
    local account_output=$(mullvad account get 2>/dev/null)

    # Check if connected (means account is active)
    if echo "$status_output" | grep -q "Connected\|Disconnected\|Blocked"; then
        log_message "Mullvad account is active, will start GUI"
        return 0
    # Check if account command succeeds AND contains an actual account number
    elif echo "$account_output" | grep -q "Mullvad account:" && echo "$account_output" | grep -E "Mullvad account:\s+[0-9]+"; then
        log_message "Mullvad account found, will start GUI"
        return 0
    else
        log_message "No Mullvad account configured, skipping GUI startup"
        return 1
    fi
}

# Check if Mullvad GUI is already running
check_mullvad_running() {
    if pgrep -x "mullvad-gui" >/dev/null; then
        log_message "Mullvad GUI already running, skipping startup"
        return 0
    fi
    return 1
}

# Start Mullvad GUI minimized to tray
start_mullvad_gui() {
    log_message "Starting Mullvad GUI minimized to tray"

    # Ensure GUI settings are configured to start minimized
    local settings_file="$HOME/.config/Mullvad VPN/gui_settings.json"
    if [ -f "$settings_file" ]; then
        # Update startMinimized to true in the JSON settings
        if command -v jq &>/dev/null; then
            log_message "Updating Mullvad GUI settings to start minimized"
            jq '.startMinimized = true' "$settings_file" > "${settings_file}.tmp" && mv "${settings_file}.tmp" "$settings_file"
        else
            # Fallback: use sed to update the setting
            log_message "Updating Mullvad GUI settings to start minimized (fallback method)"
            sed -i 's/"startMinimized":false/"startMinimized":true/g' "$settings_file"
        fi
    fi

    # Start GUI with minimize flag and redirect output to avoid console spam
    "/opt/Mullvad VPN/mullvad-gui" --minimize-to-tray &>/dev/null & disown

    # Give it a moment to start
    sleep 2

    # Verify it started
    if pgrep -x "mullvad-gui" >/dev/null; then
        log_message "✓ Mullvad GUI started successfully"
        return 0
    else
        log_message "⚠ Failed to start Mullvad GUI"
        return 1
    fi
}

# Main execution
main() {
    # Wait longer for desktop environment and system tray to be ready
    sleep 10

    log_message "=== Mullvad GUI Startup Check ==="

    # Check prerequisites
    if ! check_mullvad_installed; then
        exit 0
    fi

    # Check if already running
    if check_mullvad_running; then
        exit 0
    fi

    # Check if account is logged in
    if ! check_mullvad_account; then
        exit 0
    fi

    # Start the GUI
    start_mullvad_gui

    log_message "=== Startup check complete ==="
}

# Run the main function
main "$@"
