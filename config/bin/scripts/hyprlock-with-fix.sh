#!/bin/bash

# ==============================================================================
# Hyprlock Wrapper with Floating Window Fix
# ==============================================================================
# This wrapper runs hyprlock and automatically fixes floating windows that get
# repositioned off-screen when returning from the lock screen.
#
# This addresses a known Hyprland bug where floating windows can be moved to
# extreme coordinates (4000+ pixels) after monitor state changes.
# ==============================================================================

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Log function for debugging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] hyprlock-with-fix: $1" >> /tmp/hyprlock-with-fix.log
}

# Function to get all floating windows before locking
save_floating_window_states() {
    log "Saving floating window states before locking"
    # Save current floating window positions to temp file
    hyprctl clients -j | jq -r '.[] | select(.floating == true) | "\(.address)|\(.at[0])|\(.at[1])|\(.size[0])|\(.size[1])|\(.workspace.id)|\(.class)|\(.title)"' > /tmp/floating-windows-before-lock.txt

    if [ -s /tmp/floating-windows-before-lock.txt ]; then
        log "Saved $(wc -l < /tmp/floating-windows-before-lock.txt) floating windows"
    else
        log "No floating windows to save"
    fi
}

# Function to run the fix-offscreen-windows script
fix_floating_windows() {
    log "Running floating window fix"

    # Wait a moment for the desktop to stabilize after unlock
    sleep 0.5

    # Run the existing fix script
    if [ -x "$SCRIPT_DIR/fix-offscreen-windows.sh" ]; then
        "$SCRIPT_DIR/fix-offscreen-windows.sh"
    else
        log "ERROR: fix-offscreen-windows.sh not found or not executable"
    fi
}

# Main execution
main() {
    log "Starting hyprlock wrapper"

    # Save current workspace before locking
    CURRENT_WORKSPACE=$(hyprctl activewindow -j | jq -r '.workspace.id // 1')
    log "Current workspace before lock: $CURRENT_WORKSPACE"

    # Save floating window states before locking
    save_floating_window_states

    # Run hyprlock and wait for it to exit (when user unlocks)
    log "Launching hyprlock"
    hyprlock
    LOCK_EXIT_CODE=$?

    log "Hyprlock exited with code: $LOCK_EXIT_CODE"

    # Give the compositor a moment to stabilize
    sleep 0.3

    # Fix any floating windows that got repositioned
    fix_floating_windows

    # Return to original workspace to force visual refresh
    log "Returning to workspace $CURRENT_WORKSPACE to refresh display"
    hyprctl dispatch workspace "$CURRENT_WORKSPACE"
    sleep 0.2

    # Clean up temp file
    rm -f /tmp/floating-windows-before-lock.txt

    log "Hyprlock wrapper completed"

    exit $LOCK_EXIT_CODE
}

# Run main function
main
