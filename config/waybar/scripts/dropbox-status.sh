#!/bin/bash

# Check if dropbox is installed
if ! command -v dropbox &> /dev/null; then
    # Dropbox not installed
    echo '{"text": "Not installed", "tooltip": "Dropbox is not installed", "class": "stopped", "alt": "stopped"}'
    exit 0
fi

# Get dropbox status
STATUS=$(dropbox status 2>&1)

# Handle different states
if [[ $STATUS == *"isn't running"* ]]; then
    # Dropbox is not running
    echo '{"text": "Stopped", "tooltip": "Dropbox is not running", "class": "stopped", "alt": "stopped"}'
elif [[ $STATUS == *"Another instance"* ]]; then
    # Another instance is running - use filestatus to get actual status
    FILE_STATUS=$(dropbox filestatus 2>&1)
    
    if [[ $FILE_STATUS == *"Up to date"* || $FILE_STATUS == *"not connected"* ]]; then
        echo '{"text": "Synced", "tooltip": "Dropbox is up to date", "class": "running", "alt": "running"}'
    elif [[ $FILE_STATUS == *"Syncing"* || $FILE_STATUS == *"Indexing"* ]]; then
        echo '{"text": "Syncing", "tooltip": "Dropbox is syncing files", "class": "syncing", "alt": "syncing"}'
    else
        # Running but status unknown
        echo '{"text": "Running", "tooltip": "Dropbox is running", "class": "running", "alt": "running"}'
    fi
elif [[ $STATUS == *"Syncing"* || $STATUS == *"Indexing"* || $STATUS == *"Uploading"* || $STATUS == *"Downloading"* ]]; then
    # Dropbox is syncing
    DETAILS=$(echo "$STATUS" | grep -v "Dropbox" | sed 's/^[[:space:]]*//')
    echo "{\"text\": \"Syncing\", \"tooltip\": \"$DETAILS\", \"class\": \"syncing\", \"alt\": \"syncing\"}"
elif [[ $STATUS == *"Up to date"* ]]; then
    # Dropbox is up to date
    echo '{"text": "Synced", "tooltip": "Dropbox is up to date", "class": "running", "alt": "running"}'
else
    # Unknown status - return raw status with status class
    STATUS_CLEAN=$(echo "$STATUS" | tr -d '\n')
    echo "{\"text\": \"Running\", \"tooltip\": \"${STATUS_CLEAN}\", \"class\": \"running\", \"alt\": \"running\"}"
fi