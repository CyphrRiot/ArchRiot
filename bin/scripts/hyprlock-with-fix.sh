#!/bin/bash

# ==============================================================================
# ArchRiot Hyprlock Wrapper with Automatic Window Fix
# ==============================================================================
# Runs hyprlock and automatically fixes off-screen windows after unlock
# Solves the AMD GPU bug where windows get positioned off-screen after lock/unlock
# ==============================================================================

# Debug log file
DEBUG_LOG="/tmp/hyprlock-wrapper-debug.log"

# Log function
log_debug() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$DEBUG_LOG"
}

# Path to the window fix script
FIX_SCRIPT="$HOME/.local/bin/scripts/fix-offscreen-windows.sh"

log_debug "=== HYPRLOCK WRAPPER STARTED ==="
log_debug "Running hyprlock..."

# Run hyprlock (this blocks until unlock)
hyprlock

log_debug "Hyprlock exited, checking for window fixes needed..."

# After hyprlock exits (user unlocked), check if we need to fix windows
# Only run the fix on AMD systems since this is an AMD-specific bug
if lspci | grep -i "amd.*vga\|amd.*display\|vga.*amd\|display.*amd" >/dev/null 2>&1; then
    log_debug "AMD GPU detected, proceeding with window fix"

    # Small delay to let windows settle after unlock
    log_debug "Waiting 0.5s for windows to settle..."
    sleep 0.5

    # Run the window fix script if it exists
    if [[ -x "$FIX_SCRIPT" ]]; then
        log_debug "Running fix script: $FIX_SCRIPT"
        "$FIX_SCRIPT" 2>&1 | while read line; do log_debug "FIX: $line"; done
        log_debug "Fix script completed"
    else
        log_debug "ERROR: Fix script not found or not executable at $FIX_SCRIPT"
        # Fallback notification if script is missing
        notify-send "Window Fix Error" "Fix script not found at $FIX_SCRIPT" --urgency=critical
    fi
else
    log_debug "No AMD GPU detected, skipping window fix"
fi

log_debug "=== HYPRLOCK WRAPPER FINISHED ==="
