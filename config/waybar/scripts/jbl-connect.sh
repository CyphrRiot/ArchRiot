#!/bin/bash
# Script to toggle JBL headphones connection with proper feedback

DEVICE_MAC="40:72:18:7D:0A:69"

# Create a file to indicate user action (prevents duplicate notifications from the status script)
touch /tmp/jbl_user_action

# Check current connection status
is_connected() {
    bluetoothctl info "$DEVICE_MAC" 2>/dev/null | grep -q "Connected: yes"
}

# Main script
if is_connected; then
    # Currently connected, so disconnect
    if bluetoothctl disconnect "$DEVICE_MAC"; then
        notify-send --app-name="JBL Headphones" --urgency="normal" "Disconnected" "JBL headphones have been disconnected"
        exit 0
    else
        notify-send --app-name="JBL Headphones" --urgency="critical" "Disconnection Failed" "Could not disconnect from JBL headphones"
        exit 1
    fi
else
    # Currently disconnected, so connect
    if bluetoothctl connect "$DEVICE_MAC"; then
        # Get battery level if available
        sleep 1  # Brief pause to allow connection to stabilize
        battery_info=$(bluetoothctl info "$DEVICE_MAC" 2>/dev/null | grep "Battery Percentage" | grep -o '[0-9]*' | head -1)
        
        if [ -n "$battery_info" ]; then
            notify-send --app-name="JBL Headphones" --urgency="normal" "Connected" "JBL headphones connected. Battery: ${battery_info}%"
        else
            notify-send --app-name="JBL Headphones" --urgency="normal" "Connected" "JBL headphones have been connected"
        fi
        exit 0
    else
        notify-send --app-name="JBL Headphones" --urgency="critical" "Connection Failed" "Could not connect to JBL headphones. Check if they are powered on and in range."
        exit 1
    fi
fi