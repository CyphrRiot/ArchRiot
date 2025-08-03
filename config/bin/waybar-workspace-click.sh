#!/bin/bash
# Safe workspace click handler for waybar
# Only switches workspace if a valid workspace name is provided

if [[ -n "$1" && "$1" =~ ^[0-9]+$ ]]; then
    # Valid workspace number provided
    hyprctl dispatch workspace "$1"
else
    # No workspace or invalid workspace - do nothing
    exit 0
fi