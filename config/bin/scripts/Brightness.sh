#!/usr/bin/env bash
# ==============================================================================
# Brightness Control Script for ArchRiot
# ==============================================================================
# Handles screen brightness controls with notifications for waybar integration
# ==============================================================================

# Function to send notifications
notify() {
    local title="$1"
    local message="$2"
    local icon="$3"

    # Dismiss existing notifications first (for mako)
    if command -v makoctl &> /dev/null; then
        makoctl dismiss --all
    fi

    if command -v notify-send &> /dev/null; then
        notify-send --replace-id=9999 --app-name="Brightness Control" --urgency="normal" --icon="$icon" "$title" "$message"
    fi
}

# Function to get current brightness percentage
get_brightness() {
    if command -v brightnessctl &> /dev/null; then
        brightnessctl -m | cut -d, -f4 | tr -d '%'
    else
        echo "0"
    fi
}

# Function to get brightness icon based on level
get_brightness_icon() {
    local brightness="$1"

    if [ "$brightness" -ge 75 ]; then
        echo "brightness-high"
    elif [ "$brightness" -ge 50 ]; then
        echo "brightness-medium"
    elif [ "$brightness" -ge 25 ]; then
        echo "brightness-low"
    else
        echo "brightness-min"
    fi
}

# Main brightness control function
control_brightness() {
    local action="$1"
    local current_brightness
    local new_brightness
    local icon

    case "$action" in
        --up)
            if command -v brightnessctl &> /dev/null; then
                brightnessctl set 5%+ > /dev/null
                new_brightness=$(get_brightness)
                icon=$(get_brightness_icon "$new_brightness")
                notify "Brightness" "${new_brightness}%" "$icon"
            fi
            ;;
        --down)
            if command -v brightnessctl &> /dev/null; then
                brightnessctl set 5%- > /dev/null
                new_brightness=$(get_brightness)
                icon=$(get_brightness_icon "$new_brightness")
                notify "Brightness" "${new_brightness}%" "$icon"
            fi
            ;;
        --set)
            local value="$2"
            if [[ -n "$value" && "$value" =~ ^[0-9]+$ ]]; then
                if command -v brightnessctl &> /dev/null; then
                    brightnessctl set "${value}%" > /dev/null
                    new_brightness=$(get_brightness)
                    icon=$(get_brightness_icon "$new_brightness")
                    notify "Brightness" "${new_brightness}%" "$icon"
                fi
            else
                echo "Error: Invalid brightness value. Use a number between 0-100."
                exit 1
            fi
            ;;
        --get)
            get_brightness
            ;;
        *)
            echo "Usage: $0 {--up|--down|--set <value>|--get}"
            echo "  --up      Increase brightness by 5%"
            echo "  --down    Decrease brightness by 5%"
            echo "  --set <n> Set brightness to n% (0-100)"
            echo "  --get     Get current brightness percentage"
            exit 1
            ;;
    esac
}

# Check if brightnessctl is installed
if ! command -v brightnessctl &> /dev/null; then
    notify "Brightness Error" "brightnessctl not found" "dialog-error"
    echo "Error: brightnessctl is not installed"
    exit 1
fi

# Main execution
control_brightness "$@"
