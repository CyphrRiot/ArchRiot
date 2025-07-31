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

# Move window to center by focusing first, then centering
center_window() {
    local address="$1"
    # First focus the window to bring it into view
    hyprctl dispatch focuswindow "address:$address"
    sleep 0.1
    # Then center it properly
    hyprctl dispatch centerwindow
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

    # Get all windows and check positions
    local fixed_count=0

    while read -r address x y width height class title; do
        # Skip if any values are empty or null
        [[ "$address" == "null" || -z "$address" ]] && continue
        [[ "$x" == "null" || -z "$x" ]] && continue
        [[ "$y" == "null" || -z "$y" ]] && continue

        # Check if window is off-screen
        if (( x < 0 || y < 0 || x > screen_width || y > screen_height )); then
            echo "Found off-screen window: $class ($title) at position $x,$y"
            center_window "$address"
            ((fixed_count++))
            # Small delay to avoid overwhelming Hyprland
            sleep 0.1
        fi
    done <<< "$(get_windows | jq -r '.[] | "\(.address) \(.at[0]) \(.at[1]) \(.size[0]) \(.size[1]) \(.class) \(.title)"')"

    if (( fixed_count > 0 )); then
        echo "Fixed $fixed_count off-screen windows"
        notify-send "Window Fix Complete" "Fixed $fixed_count off-screen windows" --urgency=normal
    else
        echo "No off-screen windows found"
        notify-send "Window Fix Complete" "No off-screen windows found - all good!" --urgency=low
    fi
}

# Run the fix
fix_offscreen_windows
