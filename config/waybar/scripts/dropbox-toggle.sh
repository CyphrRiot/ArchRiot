#!/bin/bash

# Check if dropbox is installed
if ! command -v dropbox &> /dev/null; then
    notify-send "Dropbox" "Dropbox is not installed"
    exit 1
fi

# Directly try to stop Dropbox first
dropbox stop
notify-send "Dropbox" "Stopping Dropbox..."

# Check if it's still running after a brief pause
sleep 2
STATUS=$(dropbox status 2>&1)

if [[ $STATUS == *"isn't running"* ]]; then
    # Dropbox is now stopped, try to start it
    dropbox start
    notify-send "Dropbox" "Starting Dropbox..."
fi