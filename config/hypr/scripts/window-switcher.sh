#!/bin/bash
# ==============================================================================
# ArchRiot Window Switcher Script
# ==============================================================================
# Shows window list with app names using fuzzel for easy selection
# ==============================================================================

# Show all windows across all workspaces with clean format (but keep address internally)
windows=$(hyprctl clients -j | jq -r '.[] | "[\(.workspace.id)] \(.title)||||\(.address)"')

# Check if we have any windows
if [[ -z "$windows" ]]; then
    notify-send "Window Switcher" "No windows found"
    exit 0
fi

# Show clean list to user (without addresses)
display_list=$(echo "$windows" | cut -d'|' -f1)
selected_display=$(echo "$display_list" | fuzzel --dmenu --prompt="Switch to: " --width=60 --lines=10)

if [[ -n "$selected_display" ]]; then
    # Find the corresponding address from the original list
    address=$(echo "$windows" | grep -F "$selected_display" | cut -d'|' -f5)

    # Focus the selected window and bring it to current workspace if needed
    hyprctl dispatch focuswindow "address:$address"

    # If it's a floating window, make sure it's visible
    hyprctl dispatch bringactivetotop
fi
