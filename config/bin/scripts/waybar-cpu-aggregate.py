#!/usr/bin/env python3
"""
CPU usage meter for Waybar
Shows total CPU usage (averaged across all cores)
"""
import json
import psutil
import time

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
        return ""
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

def get_cpu_aggregate():
    # Get overall CPU usage (already averaged across all cores)
    total_usage = psutil.cpu_percent(interval=1)

    # Get number of cores for tooltip info
    core_count = psutil.cpu_count()

    # Show the proper averaged usage percentage
    display_percent = int(total_usage)

    # Get load average
    load_avg = psutil.getloadavg()[0]

    # Convert CPU usage to visual bar indicator using reusable function
    bar = get_visual_bar(total_usage, show_empty=False)

    # Determine color class based on usage
    if total_usage >= 90:
        css_class = "critical"
    elif total_usage >= 70:
        css_class = "warning"
    else:
        css_class = "normal"

    # Create output showing visual bar with CPU icon
    output = {
        "text": f"{bar} 󰍛",
        "tooltip": f"CPU Usage: {total_usage:.1f}% | {core_count} cores | Load: {load_avg:.2f}",
        "class": css_class,
        "percentage": display_percent
    }

    return json.dumps(output)

if __name__ == "__main__":
    try:
        print(get_cpu_aggregate())
    except Exception as e:
        # Fallback output if something goes wrong
        print(json.dumps({
            "text": "-- 󰍛",
            "tooltip": f"CPU Error: {str(e)}",
            "class": "critical",
            "percentage": 0
        }))
