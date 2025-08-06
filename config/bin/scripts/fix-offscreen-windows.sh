#!/bin/bash

# ==============================================================================
# Fix Off-Screen Windows Script for ArchRiot
# ==============================================================================
# Moves windows that are positioned beyond screen boundaries back to center
# Fixes the AMD DPMS wake bug where windows get offset by thousands of pixels
# ==============================================================================

# Get screen resolution
get_screen_resolution() {
    hyprctl monitors -j | jq -r '.[0] | "\(.width) \(.height)"'
}

# Get all windows and their positions
get_windows() {
    hyprctl clients -j
}

# Move window to center by switching workspace, focusing, then centering
center_window() {
    local address="$1"
    local workspace="$2"
    local current_workspace="$3"

    # Switch to the window's workspace first
    hyprctl dispatch workspace "$workspace"
    sleep 0.1

    # Focus the window to bring it into view
    hyprctl dispatch focuswindow "address:$address"
    sleep 0.1

    # Center it properly
    hyprctl dispatch centerwindow
    sleep 0.1

    # Switch back to original workspace
    hyprctl dispatch workspace "$current_workspace"
}

# Main function
fix_offscreen_windows() {
    # Get screen dimensions
    read screen_width screen_height <<< "$(get_screen_resolution)"

    if [[ -z "$screen_width" || -z "$screen_height" ]]; then
        echo "Error: Could not get screen resolution"
        notify-send "Window Fix Error" "Could not get screen resolution" --urgency=critical
        exit 1
    fi

    echo "Screen resolution: ${screen_width}x${screen_height}"

    # Get current workspace to return to later
    local current_workspace=$(hyprctl activewindow -j | jq -r '.workspace.id // 1')

    # Get all windows and check positions
    local fixed_count=0

    while read -r address x y width height floating workspace_id class title; do
        # Skip if any values are empty or null
        [[ "$address" == "null" || -z "$address" ]] && continue
        [[ "$x" == "null" || -z "$x" ]] && continue
        [[ "$y" == "null" || -z "$y" ]] && continue
        [[ "$floating" == "null" || -z "$floating" ]] && continue
        [[ "$workspace_id" == "null" || -z "$workspace_id" ]] && continue

        # Only check floating windows (tiled windows don't get thrown off-screen)
        if [[ "$floating" != "true" ]]; then
            continue
        fi

        echo "Checking floating window: $class ($title) at position $x,$y"

        # Check if window is off-screen (expanded detection)
        if (( x < -100 || y < -100 || x > (screen_width + 100) || y > (screen_height + 100) )); then
            echo "Found off-screen floating window: $class ($title) at position $x,$y on workspace $workspace_id"
            center_window "$address" "$workspace_id" "$current_workspace"
            ((fixed_count++))
            # Small delay to avoid overwhelming Hyprland
            sleep 0.2
        fi
    done <<< "$(get_windows | jq -r '.[] | "\(.address) \(.at[0]) \(.at[1]) \(.size[0]) \(.size[1]) \(.floating) \(.workspace.id) \(.class) \(.title)"')"

    if (( fixed_count > 0 )); then
        echo "Fixed $fixed_count off-screen floating windows"
        notify-send "Window Fix Complete" "Fixed $fixed_count off-screen floating windows" --urgency=normal
    else
        echo "No off-screen floating windows found"
    fi
}

# Run the fix
fix_offscreen_windows
