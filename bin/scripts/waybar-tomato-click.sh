#!/bin/bash
# Simple file-based tomato timer click handler
# Single click: start/pause
# Double click: reset to 25:00

STATE_FILE="/tmp/waybar-tomato-timer.state"
CLICK_FILE="/tmp/waybar-tomato-click"

# Check for double click
if [[ -f "$CLICK_FILE" ]]; then
    # Double click detected - reset timer
    rm -f "$CLICK_FILE"
    
    # Create a reset command state
    cat > "$STATE_FILE" << EOF
{
    "action": "reset",
    "timestamp": $(date +%s)
}
EOF
else
    # First click - set marker for double-click detection
    echo "$(date +%s%N)" > "$CLICK_FILE"
    (
        sleep 0.4
        if [[ -f "$CLICK_FILE" ]]; then
            # Single click confirmed - toggle start/pause
            rm -f "$CLICK_FILE"
            
            # Create a toggle command state
            cat > "$STATE_FILE" << EOF
{
    "action": "toggle",
    "timestamp": $(date +%s)
}
EOF
        fi
    ) &
fi