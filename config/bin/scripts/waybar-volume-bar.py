#!/usr/bin/env python3
"""
Volume meter for Waybar with visual bar indicators
Shows speaker and microphone volume with bar progression
"""

import json
import subprocess
import shutil


def have_pamixer():
    return shutil.which("pamixer") is not None


def pipewire_ready():
    try:
        subprocess.run(
            ["pactl", "info"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            check=True,
            timeout=0.8,
        )
        return True
    except Exception:
        return False


def get_visual_bar(percentage, show_empty=True):
    """
    Convert percentage to visual bar indicator
    Args:
        percentage: Value from 0-100
        show_empty: Whether to show empty bar for values <= 10%
    Returns:
        Unicode block character representing the level
    """
    if percentage <= 0:
        return "▁"
    elif percentage <= 2:
        return "▁"
    elif percentage <= 5:
        return "▂"
    elif percentage <= 10:
        return "▃"
    elif percentage <= 20:
        return "▄"
    elif percentage <= 35:
        return "▅"
    elif percentage <= 50:
        return "▆"
    elif percentage <= 75:
        return "▇"
    else:
        return "█"


def get_volume_info():
    """Get speaker volume and mute status using pamixer"""
    if not have_pamixer() or not pipewire_ready():
        return None, None
    try:
        mute_result = subprocess.run(
            ["pamixer", "--get-mute"],
            capture_output=True,
            text=True,
            check=True,
            timeout=0.8,
        )
        is_muted = mute_result.stdout.strip() == "true"
        volume_result = subprocess.run(
            ["pamixer", "--get-volume"],
            capture_output=True,
            text=True,
            check=True,
            timeout=0.8,
        )
        volume = int(volume_result.stdout.strip())
        return is_muted, volume
    except Exception:
        return None, None


def get_volume_bar():
    """Get volume with visual bar indicator"""
    is_muted, volume = get_volume_info()

    if is_muted is None or volume is None:
        return {
            "text": "▁ 󰖁",
            "tooltip": "Audio not ready",
            "class": "muted",
            "percentage": 0,
        }
    elif is_muted:
        # Muted - show flat line bar with muted icon
        output = {
            "text": "▁ 󰖁",
            "tooltip": f"Speaker: Muted (was {volume}%)",
            "class": "muted",
            "percentage": 0,
        }
    else:
        # Get visual bar using reusable function
        bar = get_visual_bar(volume, show_empty=False)

        # Choose appropriate volume icon based on level
        if volume == 0:
            icon = "󰕿"  # No volume
        elif volume <= 33:
            icon = "󰖀"  # Low volume
        elif volume <= 66:
            icon = "󰕾"  # Medium volume
        else:
            icon = "󰕾"  # High volume (same as medium)

        # Determine color class based on volume
        if volume >= 100:
            css_class = "critical"
        elif volume >= 85:
            css_class = "warning"
        else:
            css_class = "normal"

        output = {
            "text": f"{bar} {icon}",
            "tooltip": f"Speaker Volume: {volume}%",
            "class": css_class,
            "percentage": volume,
        }

    return output


if __name__ == "__main__":
    try:
        result = get_volume_bar()
        print(json.dumps(result))
    except Exception as e:
        # Fallback output if something goes wrong
        print(
            json.dumps(
                {
                    "text": "-- 󰕿",
                    "tooltip": f"Volume Error: {str(e)}",
                    "class": "critical",
                    "percentage": 0,
                }
            )
        )
