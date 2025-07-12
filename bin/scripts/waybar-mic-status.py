#!/usr/bin/env python3
"""
Waybar microphone status module
Shows microphone icon based on mute status
"""

import json
import subprocess
import sys

def get_mic_status():
    """Get microphone mute status using pamixer"""
    try:
        # Check if microphone is muted
        result = subprocess.run(
            ['pamixer', '--default-source', '--get-mute'],
            capture_output=True,
            text=True,
            check=True
        )
        is_muted = result.stdout.strip() == 'true'
        
        # Get microphone volume
        volume_result = subprocess.run(
            ['pamixer', '--default-source', '--get-volume'],
            capture_output=True,
            text=True,
            check=True
        )
        volume = int(volume_result.stdout.strip())
        
        return is_muted, volume
    except (subprocess.CalledProcessError, FileNotFoundError, ValueError):
        return True, 0  # Assume muted if command fails

def main():
    is_muted, volume = get_mic_status()
    
    # Choose icon based on status - using the FILLED vertical icons
    if is_muted:
        icon = "󰍭"  # Muted microphone icon (FILLED - from vertical version)
        css_class = "muted"
        status_text = "Muted"
        tooltip = f"Microphone: {status_text}"
    else:
        icon = "󰍬"  # Active microphone icon (FILLED - from vertical version)
        css_class = "active"
        status_text = f"{volume}%"
        tooltip = f"Microphone: {status_text}"
    
    # Return JSON for waybar
    output = {
        "text": icon,
        "tooltip": tooltip,
        "class": css_class,
        "alt": status_text
    }
    
    print(json.dumps(output))

if __name__ == "__main__":
    main()