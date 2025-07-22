#!/bin/bash

# Get the status
STATUS=$(mullvad status)

if echo "$STATUS" | grep -q "Connected"; then
    # Extract the location - get the relay line (us-was-wg-002)
    LOCATION=$(echo "$STATUS" | grep "Relay:" | awk '{print $2}' | cut -d'-' -f2 | tr '[:lower:]' '[:upper:]')

    # Return JSON with purple class for connected state
    echo "{\"text\": \"󰌆 $LOCATION\", \"class\": \"mullvad-connected\", \"tooltip\": \"Mullvad VPN Connected: $LOCATION\\nRight-click to disconnect\"}"
else
    # Return JSON with blue class for disconnected state
    echo "{\"text\": \"󰌉\", \"class\": \"mullvad-disconnected\", \"tooltip\": \"Mullvad VPN Disconnected\\nLeft-click to connect\"}"
fi
