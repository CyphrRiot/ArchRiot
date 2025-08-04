#!/usr/bin/env bash
# ==============================================================================
# Smart Signal Launcher for ArchRiot
# ==============================================================================
# Either focuses existing Signal window or launches new instance if not running
# ==============================================================================

# Check if Signal is already running
if pgrep -x "signal-desktop" > /dev/null; then
    # Signal is running, try to focus its window
    # Try different approaches to find and focus the window

    # Method 1: Try to focus by class name
    if hyprctl dispatch focuswindow class:Signal 2>/dev/null; then
        exit 0
    fi

    # Method 2: Try lowercase class name
    if hyprctl dispatch focuswindow class:signal 2>/dev/null; then
        exit 0
    fi

    # Method 3: Try to find window by title containing "Signal"
    if hyprctl clients -j | jq -r '.[] | select(.title | test("Signal"; "i")) | .address' | head -1 | xargs -I {} hyprctl dispatch focuswindow address:{} 2>/dev/null; then
        exit 0
    fi

    # Method 4: If all else fails, try to bring any Signal window to current workspace
    signal_address=$(hyprctl clients -j 2>/dev/null | jq -r '.[] | select(.class == "Signal" or .class == "signal" or (.title | test("Signal"; "i"))) | .address' | head -1)

    if [ -n "$signal_address" ]; then
        hyprctl dispatch focuswindow address:$signal_address 2>/dev/null
    else
        # Signal is running but window not found, might be on another workspace
        # Try to move to Signal's workspace
        signal_workspace=$(hyprctl clients -j 2>/dev/null | jq -r '.[] | select(.class == "Signal" or .class == "signal" or (.title | test("Signal"; "i"))) | .workspace.id' | head -1)
        if [ -n "$signal_workspace" ] && [ "$signal_workspace" != "null" ]; then
            hyprctl dispatch workspace $signal_workspace
        fi
    fi
else
    # Signal is not running, launch it
    env GDK_SCALE=1 signal-desktop --ozone-platform=wayland --enable-features=UseOzonePlatform > /dev/null 2>&1 &

    # Wait for Signal to start and then try to focus it
    sleep 2

    # Try to focus the new window
    for i in {1..10}; do
        if hyprctl dispatch focuswindow class:Signal 2>/dev/null || hyprctl dispatch focuswindow class:signal 2>/dev/null; then
            break
        fi
        sleep 0.5
    done
fi
