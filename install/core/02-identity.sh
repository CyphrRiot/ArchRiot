#!/bin/bash

# Install gum for interactive input
install_gum() {
    command -v gum &>/dev/null && return 0
    echo "📦 Installing gum for user input..."
    yay -S --noconfirm --needed gum || return 1
}

# Simplified input function with optional support
get_input() {
    local prompt="$1" placeholder="$2" validation="$3" optional="$4"
    local value

    while true; do
        if command -v gum &>/dev/null; then
            if [[ "$optional" == "true" ]]; then
                value=$(gum input --placeholder "$placeholder (optional - press Enter to skip)" --prompt "$prompt> ")
            else
                value=$(gum input --placeholder "$placeholder" --prompt "$prompt> ")
            fi
        else
            if [[ "$optional" == "true" ]]; then
                echo -n "$prompt (optional - press Enter to skip)> "
            else
                echo -n "$prompt> "
            fi
            read -r value
        fi

        # If optional and empty, return empty
        if [[ "$optional" == "true" && -z "$value" ]]; then
            echo ""
            return
        fi

        # If not empty, validate
        if [[ -n "$value" && $value =~ $validation ]]; then
            echo "$value"
            return
        fi

        if [[ -z "$value" ]]; then
            echo "❌ This field is required, please enter a value"
        else
            echo "❌ Invalid input format, please try again"
        fi
    done
}

# Get user identity with validation
get_user_identity() {
    echo -e "\n🔍 Git Configuration (Automated - Skipped)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Git configuration skipped for automated installation."
    echo "Configure manually later with:"
    echo "  git config --global user.name \"Your Name\""
    echo "  git config --global user.email \"your@email.com\""
    echo ""

    # No interactive prompts - just set empty values
    ARCHRIOT_USER_NAME=""
    ARCHRIOT_USER_EMAIL=""

    # Export and persist empty values
    export ARCHRIOT_USER_NAME ARCHRIOT_USER_EMAIL
    local env_file="$HOME/.config/archriot/user.env"
    mkdir -p "$(dirname "$env_file")"
    {
        echo "ARCHRIOT_USER_NAME=''"
        echo "ARCHRIOT_USER_EMAIL=''"
    } > "$env_file"

    echo "⚠ Git identity skipped - configure manually when needed"
}

# Main execution
install_gum && get_user_identity
