#!/bin/bash
# Custom Waybar Mullvad VPN Status Module
# This script outputs Mullvad VPN status in JSON format for Waybar

# Check if mullvad command is available
if ! command -v mullvad &> /dev/null; then
    echo '{"text": "Mullvad not found", "tooltip": "Mullvad CLI not found", "class": "error"}'
    exit 1
fi

# Get Mullvad status
MULLVAD_STATUS=$(mullvad status)
CONNECTED=$(echo "$MULLVAD_STATUS" | grep -i "connected" | head -1)
DISCONNECTED=$(echo "$MULLVAD_STATUS" | grep -i "disconnected")

if [[ -n "$CONNECTED" ]]; then
    # Get the region code (WAS for Washington DC, etc.)
    RELAY_INFO=$(echo "$MULLVAD_STATUS" | grep -i "Relay:" | head -1)
    # Try to extract location info (e.g., "us-was-wg-002" -> "WAS")
    if [[ "$RELAY_INFO" =~ us-([a-z]+) ]]; then
        REGION="${BASH_REMATCH[1]^^}"  # Convert to uppercase
    else
        # Fallback to parsing the visible location
        LOCATION_INFO=$(echo "$MULLVAD_STATUS" | grep -i "Visible location:" | head -1)
        if [[ "$LOCATION_INFO" =~ ([A-Z][a-z]+)[,\.] ]]; then
            REGION="${BASH_REMATCH[1]}"
        else
            REGION="VPN"  # Default if we can't parse the location
        fi
    fi
    
    # Use a lock emoji + location
    echo "{\"text\": \"🔒 ${REGION}\", \"tooltip\": \"$MULLVAD_STATUS\", \"class\": \"connected\"}"
elif [[ -n "$DISCONNECTED" ]]; then
    echo "{\"text\": \"🔓\", \"tooltip\": \"Mullvad VPN: Disconnected\", \"class\": \"disconnected\"}"
else
    echo "{\"text\": \"🔌\", \"tooltip\": \"$MULLVAD_STATUS\", \"class\": \"unknown\"}"
fi