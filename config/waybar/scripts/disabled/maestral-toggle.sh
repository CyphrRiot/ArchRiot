#!/bin/bash

# Check if maestral is installed
if ! command -v maestral &> /dev/null; then
    notify-send "Maestral" "Maestral is not installed"
    exit 1
fi

# Get current status
STATUS=$(maestral status 2>&1)

if [[ $STATUS == *"Stopped"* ]]; then
    # Maestral is stopped, start it
    maestral start
    notify-send "Maestral" "Starting Maestral..."
elif [[ $STATUS == *"Paused"* ]]; then
    # Maestral is paused, resume it
    maestral resume
    notify-send "Maestral" "Resuming Maestral..."
elif [[ $STATUS == *"Up to date"* || $STATUS == *"Syncing"* || $STATUS == *"Indexing"* || $STATUS == *"Starting"* || $STATUS == *"Idle"* ]]; then
    # Maestral is running, pause it
    maestral pause
    notify-send "Maestral" "Pausing Maestral..."
else
    # Unknown state, try to stop it first then start it
    maestral stop
    sleep 1
    maestral start
    notify-send "Maestral" "Restarting Maestral..."
fi