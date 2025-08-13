#!/usr/bin/env python3
"""
Accurate memory usage for Waybar
Uses the same calculation as modern tools like btop: Total - Available
"""
import json

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
    elif percentage <= 15:
        return "▁"
    elif percentage <= 30:
        return "▂"
    elif percentage <= 45:
        return "▃"
    elif percentage <= 60:
        return "▄"
    elif percentage <= 75:
        return "▅"
    elif percentage <= 85:
        return "▆"
    elif percentage <= 95:
        return "▇"
    else:
        return "█"

def get_memory_info():
    # Read /proc/meminfo
    meminfo = {}
    with open('/proc/meminfo', 'r') as f:
        for line in f:
            if line.strip():
                key, value = line.split()[:2]
                meminfo[key.rstrip(':')] = int(value)

    # Convert from KB to GB
    total_gb = meminfo['MemTotal'] / (1024 * 1024)
    avail_gb = meminfo['MemAvailable'] / (1024 * 1024)
    free_gb = meminfo['MemFree'] / (1024 * 1024)

    # Calculate used memory - two different methods:

    # Method 1: Total - Available (modern tools like btop)
    used_modern_gb = total_gb - avail_gb

    # Method 2: Traditional (used = total - free - buffers - cached)
    buffers_gb = meminfo.get('Buffers', 0) / (1024 * 1024)
    cached_gb = meminfo.get('Cached', 0) / (1024 * 1024)
    used_traditional_gb = total_gb - free_gb - buffers_gb - cached_gb

    # Choose which method to use (change this line to switch)
    used_gb = used_traditional_gb  # Change to used_traditional_gb for old-school display

    # Calculate percentage
    percentage = (used_gb / total_gb) * 100

    # Get visual bar using reusable function
    bar = get_visual_bar(percentage, show_empty=False)

    # Determine status class
    if percentage >= 90:
        css_class = "critical"
    elif percentage >= 75:
        css_class = "warning"
    else:
        css_class = "normal"

    output = {
        "text": f"{bar} 󰾆",
        "tooltip": f"Used (Modern): {used_modern_gb:.1f}GB | Used (Traditional): {used_traditional_gb:.1f}GB | Available: {avail_gb:.1f}GB | Total: {total_gb:.1f}GB ({percentage:.1f}%)",
        "class": css_class,
        "percentage": round(percentage)
    }

    return json.dumps(output)

if __name__ == "__main__":
    try:
        print(get_memory_info())
    except Exception as e:
        # Fallback output if something goes wrong
        print(json.dumps({
            "text": "-- 󰾆",
            "tooltip": f"Memory Error: {str(e)}",
            "class": "critical",
            "percentage": 0
        }))
