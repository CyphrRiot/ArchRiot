#!/usr/bin/env python3
"""
Temperature meter for Waybar with visual bar indicators
Shows CPU temperature with bar progression and color coding
"""
import json
import glob
import os

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
    elif percentage <= 10:
        return "▁"
    elif percentage <= 25:
        return "▂"
    elif percentage <= 40:
        return "▃"
    elif percentage <= 55:
        return "▄"
    elif percentage <= 70:
        return "▅"
    elif percentage <= 80:
        return "▆"
    elif percentage <= 90:
        return "▇"
    else:
        return "█"

def find_temp_sensor():
    """Find the best CPU temperature sensor"""
    # Priority order for sensor types
    sensor_paths = []

    # 1. Look for coretemp hwmon (Intel/AMD)
    for hwmon_path in glob.glob('/sys/class/hwmon/hwmon*'):
        name_file = os.path.join(hwmon_path, 'name')
        if os.path.exists(name_file):
            try:
                with open(name_file, 'r') as f:
                    name = f.read().strip()
                if name == 'coretemp':
                    temp_file = os.path.join(hwmon_path, 'temp1_input')
                    if os.path.exists(temp_file):
                        sensor_paths.append(temp_file)
            except:
                continue

    # 2. Look for x86_pkg_temp thermal zone
    for zone_path in glob.glob('/sys/class/thermal/thermal_zone*'):
        type_file = os.path.join(zone_path, 'type')
        if os.path.exists(type_file):
            try:
                with open(type_file, 'r') as f:
                    zone_type = f.read().strip()
                if zone_type == 'x86_pkg_temp':
                    temp_file = os.path.join(zone_path, 'temp')
                    if os.path.exists(temp_file):
                        sensor_paths.append(temp_file)
            except:
                continue

    # 3. Fallback to thermal_zone0
    fallback = '/sys/class/thermal/thermal_zone0/temp'
    if os.path.exists(fallback):
        sensor_paths.append(fallback)

    return sensor_paths[0] if sensor_paths else None

def read_temperature():
    """Read CPU temperature in Celsius"""
    temp_sensor = find_temp_sensor()

    if not temp_sensor:
        return None

    try:
        with open(temp_sensor, 'r') as f:
            temp_raw = int(f.read().strip())

        # Convert from millidegrees to degrees
        temp_celsius = temp_raw / 1000
        return temp_celsius
    except:
        return None

def get_temp_bar():
    """Get temperature with visual bar indicator"""
    temp = read_temperature()

    if temp is None:
        return {
            "text": "-- 󰈸",
            "tooltip": "Temperature sensor not available",
            "class": "critical",
            "percentage": 0
        }

    temp_int = int(temp)

    # Convert temperature to percentage for bar calculation (55°C = 0%, 80°C = 100%)
    temp_percentage = max(0, min(100, (temp - 55) * 100 / 25))

    # Get visual bar using reusable function
    bar = get_visual_bar(temp_percentage, show_empty=True)

    # Determine color class based on temperature
    if temp >= 90:
        css_class = "critical"
    elif temp >= 80:
        css_class = "warning"
    elif temp >= 70:
        css_class = "hot"
    else:
        css_class = "normal"

    # Create output with bar and temperature icon
    output = {
        "text": f"{bar} 󰈸",
        "tooltip": f"CPU Temperature: {temp:.1f}°C",
        "class": css_class,
        "percentage": temp_int
    }

    return output

if __name__ == "__main__":
    try:
        result = get_temp_bar()
        print(json.dumps(result))
    except Exception as e:
        # Fallback output if something goes wrong
        print(json.dumps({
            "text": "-- 󰈸",
            "tooltip": f"Temperature Error: {str(e)}",
            "class": "critical",
            "percentage": 0
        }))
