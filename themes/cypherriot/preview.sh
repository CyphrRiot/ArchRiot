#!/bin/bash
# CypherRiot Theme Color Preview
# Shows the key colors used in the theme

echo -e "\n🎨 CypherRiot Theme Color Palette"
echo -e "═══════════════════════════════════════"

# Function to show color
show_color() {
    local name="$1"
    local hex="$2" 
    local rgb="${3:-}"
    printf "%-20s \033[48;2;%sm   \033[0m %s %s\n" "$name" "$rgb" "$hex" 
}

echo -e "\n🌟 Primary Colors:"
show_color "Background" "#1a1b26" "26;27;38"
show_color "Foreground" "#ffffff" "255;255;255"
show_color "Secondary BG" "#222436" "34;36;54"

echo -e "\n💜 Accent Colors:" 
show_color "Purple (Active)" "#9d7bd8" "157;123;216"
show_color "Light Purple" "#bb9af7" "187;154;247"
show_color "Blue" "#7da6ff" "125;166;255"
show_color "Cyan" "#0db9d7" "13;185;215"
show_color "Red" "#ff7a93" "255;122;147"
show_color "Orange" "#ff9e64" "255;158;100"

echo -e "\n🖥️  UI Colors:"
show_color "Window Border" "#4a2b7a" "74;43;122"
show_color "CPU Color" "#6a7de8" "106;125;232"
show_color "Memory Color" "#8a95e8" "138;149;232"
show_color "Green (Active)" "#9ece6a" "158;206;106"

echo -e "\n✨ Based on your current Waybar configuration"
echo -e "   Colors extracted from your personalized setup\n"
