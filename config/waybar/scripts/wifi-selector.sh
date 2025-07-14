#!/bin/bash

# Simple WiFi selector using fuzzel
wifi_list=$(nmcli -t -f SSID,SIGNAL,SECURITY device wifi list | grep -v '^--' | sort -t: -k2 -nr)

if [ -z "$wifi_list" ]; then
    notify-send "No WiFi networks found"
    exit 1
fi

# Format for fuzzel display
formatted_list=$(echo "$wifi_list" | while IFS=: read -r ssid signal security; do
    if [ -n "$ssid" ]; then
        # Add security icon
        if [[ "$security" == *"WPA"* ]] || [[ "$security" == *"WEP"* ]]; then
            echo "🔒 $ssid ($signal%)"
        else
            echo "📶 $ssid ($signal%)"
        fi
    fi
done)

# Show fuzzel selector
selected=$(echo "$formatted_list" | fuzzel --dmenu --prompt "Select WiFi Network:")

if [ -n "$selected" ]; then
    # Extract SSID from selection
    ssid=$(echo "$selected" | sed 's/^[🔒📶] //' | sed 's/ ([0-9]*%)$//')

    # Check if network requires password
    security=$(echo "$wifi_list" | grep "^$ssid:" | cut -d: -f3)

    if [[ "$security" == *"WPA"* ]] || [[ "$security" == *"WEP"* ]]; then
        # Prompt for password
        password=$(fuzzel --dmenu --password --prompt "Password for $ssid:")
        if [ -n "$password" ]; then
            nmcli device wifi connect "$ssid" password "$password"
        fi
    else
        # Connect without password
        nmcli device wifi connect "$ssid"
    fi

    # Refresh waybar
    pkill -SIGUSR2 waybar
fi
