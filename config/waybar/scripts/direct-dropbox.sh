#!/bin/bash

# Simple, direct approach to check Dropbox status
# This script prioritizes reliability over features

if pgrep -f "dropbox-dist" > /dev/null; then
    # Dropbox is running
    echo '{"text": "󰇚", "tooltip": "Dropbox is running", "class": "running", "alt": "running"}'
else
    # Dropbox is not running
    echo '{"text": "󰇙", "tooltip": "Dropbox is not running", "class": "stopped", "alt": "stopped"}'
fi