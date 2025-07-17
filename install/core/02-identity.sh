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
    echo -e "\n🔍 Git Configuration (Optional)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Configure Git with your name and email for commits and development."
    echo "This is optional - you can skip by pressing Enter or configure later with:"
    echo "  git config --global user.name \"Your Name\""
    echo "  git config --global user.email \"your@email.com\""
    echo ""

    # Check if we're in an interactive terminal
    if [[ -t 0 && -t 1 ]]; then
        # Interactive mode - no timeout
        echo -n "Git Name (optional - press Enter to skip): "
        read -r ARCHRIOT_USER_NAME

        echo -n "Git Email (optional - press Enter to skip): "
        read -r ARCHRIOT_USER_EMAIL
    else
        # Non-interactive mode - skip prompts
        echo "⚠ Non-interactive mode detected - skipping git configuration"
        ARCHRIOT_USER_NAME=""
        ARCHRIOT_USER_EMAIL=""
    fi

    # Export and persist
    export ARCHRIOT_USER_NAME ARCHRIOT_USER_EMAIL
    local env_file="$HOME/.config/archriot/user.env"
    mkdir -p "$(dirname "$env_file")"
    {
        echo "ARCHRIOT_USER_NAME='$ARCHRIOT_USER_NAME'"
        echo "ARCHRIOT_USER_EMAIL='$ARCHRIOT_USER_EMAIL'"
    } > "$env_file"

    if [[ -n "$ARCHRIOT_USER_NAME" || -n "$ARCHRIOT_USER_EMAIL" ]]; then
        echo "✓ Git identity configured: ${ARCHRIOT_USER_NAME:-"(no name)"} <${ARCHRIOT_USER_EMAIL:-"(no email)"}>"
    else
        echo "⚠ Git identity skipped - you can configure later if needed"
    fi
}

# Main execution
install_gum && get_user_identity
