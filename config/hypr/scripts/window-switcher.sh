#!/bin/bash
# ==============================================================================
# ArchRiot Window Switcher Script
# ==============================================================================
# Shows window list with app names using fuzzel for easy selection
# ==============================================================================

# Show all windows across all workspaces with clean format
windows=$(hyprctl clients -j | jq -r '.[] | "[\(.workspace.id)] \(.title) (\(.address))"')

# Check if we have any windows
if [[ -z "$windows" ]]; then
    notify-send "Window Switcher" "No windows found"
    exit 0
fi

# Use fuzzel to show selectable list with window names
selected=$(echo "$windows" | fuzzel --dmenu --prompt="Switch to: " --width=60 --lines=10)

if [[ -n "$selected" ]]; then
    # Extract window address from selection
    address=$(echo "$selected" | grep -o '(0x[^)]*)')
    address=${address//[()]/}

    # Focus the selected window and bring it to current workspace if needed
    hyprctl dispatch focuswindow "address:$address"

    # If it's a floating window, make sure it's visible
    hyprctl dispatch bringactivetotop
fi
