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

# Get user identity with validation and smart credential detection
get_user_identity() {
    echo -e "\n🔍 Git Configuration (Optional)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Check for existing git credentials
    local existing_name=$(git config --global user.name 2>/dev/null || echo "")
    local existing_email=$(git config --global user.email 2>/dev/null || echo "")

    if [[ -n "$existing_name" || -n "$existing_email" ]]; then
        echo "🎉 GitHub credentials found!"
        echo ""

        # Calculate box width and format entries properly
        local box_width=59
        local name_display="${existing_name:-"(not set)"}"
        local email_display="${existing_email:-"(not set)"}"

        # Format lines with proper spacing
        local name_line=$(printf "│ Username: %-*s │" $((box_width-14)) "$name_display")
        local email_line=$(printf "│ Email:    %-*s │" $((box_width-14)) "$email_display")

        echo "┌─────────────────────────────────────────────────────────┐"
        echo "│                 📋 Current Git Config                   │"
        echo "├─────────────────────────────────────────────────────────┤"
        echo "$name_line"
        echo "$email_line"
        echo "└─────────────────────────────────────────────────────────┘"
        echo ""

        echo ""
        echo -n "Would you like to use these credentials? [Y/n] (auto-yes in 10s): "

        # Use timeout to prevent hanging - default to YES
        if read -t 10 -r response </dev/tty 2>/dev/null; then
            echo ""  # Move to next line after user input
        else
            response="y"  # Default to yes on timeout
            echo "y"      # Show the default choice
            echo ""
        fi

        case "$response" in
            [nN][oO]|[nN])
                echo ""
                echo "💬 No problem! Let's set up new credentials..."
                echo ""
                ARCHRIOT_USER_NAME=$(get_input "Name" "Your full name for Git commits" "^[a-zA-Z].*" "true")
                ARCHRIOT_USER_EMAIL=$(get_input "Email" "Your email for Git commits" "^[^@]+@[^@]+\.[^@]+$" "true")
                ;;
            *)
                echo ""
                echo "✅ Using existing credentials!"
                ARCHRIOT_USER_NAME="$existing_name"
                ARCHRIOT_USER_EMAIL="$existing_email"
                ;;
        esac
    else
        echo "Configure Git with your name and email for commits and development."
        echo "This is optional - you can skip by pressing Enter or configure later with:"
        echo "  git config --global user.name \"Your Name\""
        echo "  git config --global user.email \"your@email.com\""
        echo ""

        ARCHRIOT_USER_NAME=$(get_input "Name" "Your full name for Git commits" "^[a-zA-Z].*" "true")
        ARCHRIOT_USER_EMAIL=$(get_input "Email" "Your email for Git commits" "^[^@]+@[^@]+\.[^@]+$" "true")
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
