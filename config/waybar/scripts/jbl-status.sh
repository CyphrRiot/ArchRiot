#!/bin/bash
# JBL LIVE660NC Waybar Module Script
# Shows Headphone icon with status colors

DEVICE_MAC="40:72:18:7D:0A:69"

# Check if device is connected
is_connected() {
    bluetoothctl info "$DEVICE_MAC" 2>/dev/null | grep -q "Connected: yes"
}

# Check if device exists (is paired)
device_exists() {
    bluetoothctl devices 2>/dev/null | grep -q "$DEVICE_MAC"
}

# Main function
main() {
    if ! device_exists; then
        # Device not paired - show grey icon
        echo '{"text": "󰋋", "tooltip": "JBL LIVE660NC: Not Paired", "class": "disconnected"}'
    elif is_connected; then
        # Device connected - show green icon
        # Get battery level if available
        battery_info=$(bluetoothctl info "$DEVICE_MAC" 2>/dev/null | grep "Battery Percentage" | grep -o '[0-9]*' | head -1)
        if [ -n "$battery_info" ]; then
            tooltip="JBL LIVE660NC: Connected (${battery_info}%)"
        else
            tooltip="JBL LIVE660NC: Connected"
        fi
        echo '{"text": "󰋋", "tooltip": "'"$tooltip"'", "class": "connected"}'
    else
        # Device paired but not connected - show grey icon
        echo '{"text": "󰋋", "tooltip": "JBL LIVE660NC: Disconnected", "class": "disconnected"}'
    fi
}

main