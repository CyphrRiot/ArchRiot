#!/bin/bash
# Simple ArchRiot Update Check for Waybar - No Cache Nonsense

# Get versions directly
local_version=$(cat ~/.local/share/archriot/VERSION 2>/dev/null || echo "unknown")
remote_version=$(timeout 10 curl -s https://ArchRiot.org/VERSION 2>/dev/null || echo "unknown")

# Handle click events
if [[ "$1" == "--click" ]]; then
    # Show update dialog when clicked (using the correct command)
    if [[ "$remote_version" != "unknown" && "$local_version" != "unknown" && "$remote_version" != "$local_version" ]]; then
        ~/.local/bin/version-check --gui 2>/dev/null &
    fi
    exit 0
fi

# Compare versions and output for waybar
if [[ "$remote_version" == "unknown" || "$local_version" == "unknown" ]]; then
    # Network/file error - show nothing
    echo '{"text":"-","tooltip":"Update check unavailable","class":"update-none"}'
elif [[ "$remote_version" != "$local_version" ]]; then
    # Update available - show icon
    echo '{"text":"󰚰","tooltip":"ArchRiot update available!\nCurrent: '$local_version'\nAvailable: '$remote_version'","class":"update-available"}'
else
    # Up to date - show nothing
    echo '{"text":"-","tooltip":"ArchRiot is up to date\nCurrent: '$local_version'","class":"update-none"}'
fi
