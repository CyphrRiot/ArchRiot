#!/bin/bash
# ArchRiot Update Notification Script for Waybar
# Shows update icon when ArchRiot updates are available

set -e

# Configuration
CONFIG_DIR="$HOME/.config/archriot"
CONFIG_FILE="$CONFIG_DIR/versions.cfg"
UPDATE_FLAG_FILE="$CONFIG_DIR/update_available"
UPDATE_STATE_FILE="$CONFIG_DIR/update_state"
CACHE_FILE="$CONFIG_DIR/.update_cache"
CACHE_DURATION=3600  # 1 hour cache

# Ensure config directory exists
mkdir -p "$CONFIG_DIR"

# Function to get local version
get_local_version() {
    local version_file="$HOME/.local/share/archriot/VERSION"
    if [[ -f "$version_file" ]]; then
        cat "$version_file" 2>/dev/null || echo "unknown"
    else
        echo "unknown"
    fi
}

# Function to get remote version (cached)
get_remote_version() {
    local current_time=$(date +%s)

    # Check if cache exists and is fresh
    if [[ -f "$CACHE_FILE" ]]; then
        local cache_time=$(stat -c %Y "$CACHE_FILE" 2>/dev/null || echo 0)
        local age=$((current_time - cache_time))

        if [[ $age -lt $CACHE_DURATION ]]; then
            cat "$CACHE_FILE" 2>/dev/null && return 0
        fi
    fi

    # Fetch new version with timeout
    local remote_version
    if remote_version=$(timeout 10 curl -s "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION" 2>/dev/null); then
        # Validate version format (should be like 1.1.79)
        if [[ "$remote_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "$remote_version" > "$CACHE_FILE"
            echo "$remote_version"
            return 0
        fi
    fi

    # Fallback to cached version if available
    if [[ -f "$CACHE_FILE" ]]; then
        cat "$CACHE_FILE" 2>/dev/null || echo "unknown"
    else
        echo "unknown"
    fi
}

# Function to compare versions (returns 0 if remote is newer)
is_newer_version() {
    local local_ver="$1"
    local remote_ver="$2"

    [[ "$local_ver" == "unknown" || "$remote_ver" == "unknown" ]] && return 1

    # Split versions into components
    IFS='.' read -ra LOCAL <<< "$local_ver"
    IFS='.' read -ra REMOTE <<< "$remote_ver"

    # Compare each component
    for i in {0..2}; do
        local local_part=${LOCAL[$i]:-0}
        local remote_part=${REMOTE[$i]:-0}

        if (( remote_part > local_part )); then
            return 0  # Remote is newer
        elif (( local_part > remote_part )); then
            return 1  # Local is newer or equal
        fi
    done

    return 1  # Versions are equal
}

# Function to get current update state
get_update_state() {
    local remote_version="$1"

    # Check if user has seen this version (clicked and closed dialog)
    if [[ -f "$UPDATE_STATE_FILE" ]]; then
        local seen_version=$(cat "$UPDATE_STATE_FILE" 2>/dev/null)
        if [[ "$seen_version" == "$remote_version" ]]; then
            echo "seen"
            return
        fi
    fi

    # Check if notifications are ignored via config
    if [[ -f "$CONFIG_FILE" ]] && command -v jq >/dev/null 2>&1; then
        local ignore_notifications=$(jq -r '.ignore_notifications // false' "$CONFIG_FILE" 2>/dev/null)
        if [[ "$ignore_notifications" == "true" ]]; then
            echo "ignored"
            return
        fi
    fi

    echo "new"
}

# Function to handle click action
handle_click() {
    if [[ -f "$UPDATE_FLAG_FILE" ]]; then
        # Launch upgrade dialog with fast cached launch
        if command -v version-check >/dev/null 2>&1; then
            version-check --gui 2>/dev/null &

            # Mark this version as seen (user clicked the icon)
            if [[ -f "$UPDATE_FLAG_FILE" ]]; then
                local update_info=$(cat "$UPDATE_FLAG_FILE" 2>/dev/null)
                if [[ "$update_info" =~ UPDATE_AVAILABLE:.*-\>(.*)$ ]]; then
                    local remote_version="${BASH_REMATCH[1]}"
                    echo "$remote_version" > "$UPDATE_STATE_FILE"
                fi
            fi
        else
            # Fallback to notify
            notify-send --urgency=critical --icon=system-software-update \
                "ArchRiot Update" "Please run: version-check --force"
        fi
    fi
}

# Handle command line arguments
case "${1:-}" in
    --click)
        handle_click
        exit 0
        ;;
    --local-version)
        get_local_version
        exit 0
        ;;
    --remote-version)
        get_remote_version
        exit 0
        ;;
esac

# Main logic for waybar display
main() {
    local local_version=$(get_local_version)
    local remote_version=$(get_remote_version)

    # Check if update is available
    if is_newer_version "$local_version" "$remote_version"; then
        echo "UPDATE_AVAILABLE:$local_version->$remote_version" > "$UPDATE_FLAG_FILE"

        local update_state=$(get_update_state "$remote_version")

        case "$update_state" in
            "new")
                # New update - pulsing download icon
                echo "{\"text\":\"󰚰\",\"tooltip\":\"NEW ArchRiot Update Available\\n$local_version → $remote_version\\n\\nClick to upgrade\",\"class\":\"update-new\"}"
                ;;
            "seen")
                # User saw it but didn't upgrade - notification bell
                echo "{\"text\":\"󱧘\",\"tooltip\":\"ArchRiot Update Available\\n$local_version → $remote_version\\n\\nClick to upgrade\",\"class\":\"update-seen\"}"
                ;;
            "ignored")
                # User ignored notifications - no icon
                echo ""
                ;;
        esac
    else
        # No update available - solid circle, darker
        echo "{\"text\":\"-\",\"tooltip\":\"ArchRiot is up to date\\nCurrent: $local_version\",\"class\":\"update-none\"}"

        # Clean up flag files
        [[ -f "$UPDATE_FLAG_FILE" ]] && rm -f "$UPDATE_FLAG_FILE"
        [[ -f "$UPDATE_STATE_FILE" ]] && rm -f "$UPDATE_STATE_FILE"
    fi
}

# Run main function
main
