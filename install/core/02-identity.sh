#!/bin/bash

# Install gum for interactive input
install_gum() {
    command -v gum &>/dev/null && return 0
    echo "📦 Installing gum for user input..."
    yay -S --noconfirm --needed gum || return 1
}

# Simplified input function
get_input() {
    local prompt="$1" placeholder="$2" validation="$3"
    local value

    while true; do
        if command -v gum &>/dev/null; then
            value=$(gum input --placeholder "$placeholder" --prompt "$prompt> ")
        else
            echo -n "$prompt> "
            read -r value
        fi

        [[ $value =~ $validation ]] && { echo "$value"; return; }
        echo "❌ Invalid input, please try again"
    done
}

# Get user identity with validation
get_user_identity() {
    echo -e "\n🔍 Enter identification for git and autocomplete..."

    OMARCHY_USER_NAME=$(get_input "Name" "Enter full name" "^[a-zA-Z].*")
    OMARCHY_USER_EMAIL=$(get_input "Email" "Enter email address" "^[^@]+@[^@]+\.[^@]+$")

    # Export and persist
    export OMARCHY_USER_NAME OMARCHY_USER_EMAIL
    local env_file="$HOME/.config/omarchy/user.env"
    mkdir -p "$(dirname "$env_file")"
    {
        echo "OMARCHY_USER_NAME='$OMARCHY_USER_NAME'"
        echo "OMARCHY_USER_EMAIL='$OMARCHY_USER_EMAIL'"
    } > "$env_file"

    echo "✓ User identity configured: $OMARCHY_USER_NAME <$OMARCHY_USER_EMAIL>"
}

# Main execution
install_gum && get_user_identity
