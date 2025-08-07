#!/bin/bash
# Simple ArchRiot Update Check for Waybar - No Cache Nonsense

# Get versions directly
local_version=$(cat ~/.local/share/archriot/VERSION 2>/dev/null || echo "unknown")
remote_version=$(timeout 10 curl -s https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION 2>/dev/null || echo "unknown")

# Function to compare semantic versions (returns 0 if remote is newer)
is_newer_version() {
    local local_ver="$1"
    local remote_ver="$2"

    [[ "$local_ver" == "unknown" || "$remote_ver" == "unknown" ]] && return 1

    # Split versions into arrays
    IFS='.' read -ra LOCAL <<< "$local_ver"
    IFS='.' read -ra REMOTE <<< "$remote_ver"

    # Compare major.minor.patch
    for i in {0..2}; do
        local local_part=${LOCAL[$i]:-0}
        local remote_part=${REMOTE[$i]:-0}

        if (( remote_part > local_part )); then
            return 0  # Remote is newer
        elif (( local_part > remote_part )); then
            return 1  # Local is newer
        fi
    done

    return 1  # Versions are equal
}

# Handle click events
if [[ "$1" == "--click" ]]; then
    # Show update dialog when clicked (using the correct command)
    if [[ "$remote_version" != "unknown" && "$local_version" != "unknown" ]] && is_newer_version "$local_version" "$remote_version"; then
        # Show launching notification only when there's an actual update
        notify-send -t 5000 "Launching Upgrade..." "Starting ArchRiot upgrade process..." &
        $HOME/.local/share/archriot/config/bin/version-check --gui 2>/dev/null &
    else
        # Already up to date
        notify-send -t 2000 "ArchRiot Up to Date" "Version $local_version is the latest" &
    fi
    exit 0
fi

# Compare versions and output for waybar
if [[ "$remote_version" == "unknown" || "$local_version" == "unknown" ]]; then
    # Network/file error - show nothing
    echo '{"text":"-","tooltip":"Update check unavailable","class":"update-none"}'
elif is_newer_version "$local_version" "$remote_version"; then
    # Update available - show icon
    echo '{"text":"󰚰","tooltip":"ArchRiot update available!\nCurrent: '$local_version'\nAvailable: '$remote_version'","class":"update-available"}'
else
    # Up to date - show nothing
    echo '{"text":"-","tooltip":"ArchRiot is up to date\nCurrent: '$local_version'","class":"update-none"}'
fi
