#!/bin/bash

# ==============================================================================
# ArchRiot Hyprlock Wrapper with Automatic Window Fix
# ==============================================================================
# Runs hyprlock and automatically fixes off-screen windows after unlock
# Solves the AMD GPU bug where windows get positioned off-screen after lock/unlock
# ==============================================================================

# Path to the window fix script
FIX_SCRIPT="$HOME/.local/bin/scripts/fix-offscreen-windows.sh"

# Run hyprlock (this blocks until unlock)
hyprlock

# After hyprlock exits (user unlocked), check if we need to fix windows
# Only run the fix on AMD systems since this is an AMD-specific bug
if lspci | grep -i "amd.*vga\|amd.*display\|vga.*amd\|display.*amd" >/dev/null 2>&1; then
    # Small delay to let windows settle after unlock
    sleep 0.5

    # Run the window fix script if it exists
    if [[ -x "$FIX_SCRIPT" ]]; then
        "$FIX_SCRIPT"
    else
        # Fallback notification if script is missing
        notify-send "Window Fix Error" "Fix script not found at $FIX_SCRIPT" --urgency=critical
    fi
fi
