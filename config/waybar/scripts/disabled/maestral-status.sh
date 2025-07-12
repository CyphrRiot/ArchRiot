#!/bin/bash

# Check if maestral is installed
if ! command -v maestral &> /dev/null; then
    # Maestral not installed
    echo '{"text": "Not installed", "tooltip": "Maestral is not installed", "class": "stopped", "alt": "stopped"}'
    exit 0
fi

# Get maestral status
STATUS=$(maestral status 2>&1)

# Function to format tooltip nicely with emojis
format_tooltip() {
    # Extract information from status
    local account=$(echo "$1" | grep -o 'Account.*' | sed 's/Account\s*//' || echo "Unknown")
    local usage=$(echo "$1" | grep -o 'Usage.*' | sed 's/Usage\s*//' || echo "Unknown")
    local status=$(echo "$1" | grep -o 'Status.*' | sed 's/Status\s*//' || echo "Unknown")
    local errors=$(echo "$1" | grep -o 'Sync errors.*' | sed 's/Sync errors\s*//' || echo "0")
    
    # Create a nicely formatted tooltip
    echo "👤 $account | 💾 $usage | 🔄 $status | ⚠️ Errors: $errors"
}

# Handle different states
if [[ $STATUS == *"Stopped"* ]]; then
    # Maestral is not running
    echo '{"text": "Stopped", "tooltip": "Maestral is not running", "class": "stopped", "alt": "stopped"}'
elif [[ $STATUS == *"Paused"* ]]; then
    # Maestral is paused - format the tooltip
    DETAILS=$(format_tooltip "$STATUS")
    echo "{\"text\": \"Paused\", \"tooltip\": \"$DETAILS\", \"class\": \"paused\", \"alt\": \"paused\"}"
elif [[ $STATUS == *"Syncing"* || $STATUS == *"Indexing"* || $STATUS == *"Starting"* ]]; then
    # Maestral is syncing - format the tooltip
    DETAILS=$(format_tooltip "$STATUS")
    echo "{\"text\": \"Syncing\", \"tooltip\": \"$DETAILS\", \"class\": \"syncing\", \"alt\": \"syncing\"}"
elif [[ $STATUS == *"Up to date"* || $STATUS == *"Idle"* ]]; then
    # Maestral is up to date - format the tooltip
    DETAILS=$(format_tooltip "$STATUS")
    echo "{\"text\": \"Synced\", \"tooltip\": \"$DETAILS\", \"class\": \"running\", \"alt\": \"running\"}"
else
    # Unknown status - return raw status with a clean tooltip
    DETAILS=$(format_tooltip "$STATUS")
    echo "{\"text\": \"Running\", \"tooltip\": \"$DETAILS\", \"class\": \"running\", \"alt\": \"running\"}"
fi