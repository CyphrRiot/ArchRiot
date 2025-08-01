#!/bin/bash

# ==============================================================================
# ArchRiot DPMS Wake Window Fix Script
# ==============================================================================
# Fixes off-screen windows after DPMS wake events (extended lock sessions)
# Uses same timing and approach as hyprlock wrapper for consistency
# ==============================================================================

# Debug log file
DEBUG_LOG="/tmp/dpms-wake-fix-debug.log"

# Log function
log_debug() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] DPMS-WAKE: $1" >> "$DEBUG_LOG"
}

# Path to the window fix script
FIX_SCRIPT="$HOME/.local/bin/scripts/fix-offscreen-windows.sh"

log_debug "=== DPMS WAKE FIX STARTED ==="

# Turn screen back on
log_debug "Turning DPMS on..."
hyprctl dispatch dpms on

# Only run the fix on AMD systems since this is an AMD-specific bug
if lspci | grep -i "amd.*vga\|amd.*display\|vga.*amd\|display.*amd" >/dev/null 2>&1; then
    log_debug "AMD GPU detected, proceeding with window fix after DPMS wake"

    # Longer delay for DPMS wake - system needs more time to settle
    log_debug "Waiting 2s for system to settle after DPMS wake..."
    sleep 2

    # Run the window fix script if it exists
    if [[ -x "$FIX_SCRIPT" ]]; then
        log_debug "Running fix script: $FIX_SCRIPT"
        "$FIX_SCRIPT" 2>&1 | while read line; do log_debug "FIX: $line"; done
        log_debug "Fix script completed"

        # Send notification about auto-fix
        if command -v notify-send &>/dev/null; then
            notify-send "ArchRiot" "Auto-fixed windows after DPMS wake" --icon=preferences-desktop-display --urgency=low
        fi
    else
        log_debug "ERROR: Fix script not found or not executable at $FIX_SCRIPT"
        # Fallback notification if script is missing
        if command -v notify-send &>/dev/null; then
            notify-send "Window Fix Error" "Fix script not found at $FIX_SCRIPT" --urgency=critical
        fi
    fi
else
    log_debug "No AMD GPU detected, skipping window fix"
fi

log_debug "=== DPMS WAKE FIX FINISHED ==="
