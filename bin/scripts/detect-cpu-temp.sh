#!/bin/bash

# Simple CPU Temperature Detection for Waybar
# Finds the best CPU temperature sensor automatically

find_temp_sensor() {
    # Try coretemp first (most reliable)
    for hwmon in /sys/class/hwmon/hwmon*/name; do
        [[ -f "$hwmon" ]] && [[ "$(cat "$hwmon")" == "coretemp" ]] && {
            echo "${hwmon%/name}/temp1_input"
            return 0
        }
    done

    # Try thermal zones with temperature validation
    for zone in /sys/class/thermal/thermal_zone*/temp; do
        [[ -f "$zone" ]] && {
            temp=$(cat "$zone" 2>/dev/null)
            [[ $temp -gt 30000 && $temp -lt 100000 ]] && {
                echo "$zone"
                return 0
            }
        }
    done

    # Fallback
    echo "/sys/class/thermal/thermal_zone0/temp"
}

# Main execution
temp_path=$(find_temp_sensor)
temp=$(cat "$temp_path" 2>/dev/null || echo "0")
temp_c=$((temp / 1000))

case "${1:-detect}" in
    "waybar") echo "\"hwmon-path\": [\"$temp_path\"]," ;;
    *) echo "\"$temp_path\",  // ${temp_c}°C" ;;
esac
