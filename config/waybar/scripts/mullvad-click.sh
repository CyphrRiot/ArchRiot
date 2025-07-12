#!/bin/bash

# Smart Mullvad click handler
# $1 = "left" or "right"

CLICK_TYPE="$1"
STATUS=$(mullvad status)

if echo "$STATUS" | grep -q "Connected"; then
    # Currently connected
    if [ "$CLICK_TYPE" = "right" ]; then
        # Right-click when connected: disconnect
        mullvad disconnect
    fi
    # Left-click when connected: do nothing
else
    # Currently disconnected  
    if [ "$CLICK_TYPE" = "left" ]; then
        # Left-click when disconnected: connect
        mullvad connect
    fi
    # Right-click when disconnected: do nothing
fi