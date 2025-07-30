#!/bin/bash

# Install gum for interactive input
install_gum() {
    command -v gum &>/dev/null && return 0
    echo "📦 Installing gum for user input..."
    yay -S --noconfirm --needed gum || return 1
}

# Simplified input function with gum support
get_input() {
    local prompt="$1" placeholder="$2" validation="$3" optional="$4"
    local value

    if command -v gum &>/dev/null; then
        if [[ "$optional" == "true" ]]; then
            value=$(gum input --placeholder "$placeholder (optional - press Enter to skip)" --prompt "$prompt> " || echo "")
        else
            value=$(gum input --placeholder "$placeholder" --prompt "$prompt> " || echo "")
        fi
    else
        if [[ "$optional" == "true" ]]; then
            echo -n "$prompt ($placeholder - optional, press Enter to skip)> " >&2
        else
            echo -n "$prompt ($placeholder)> " >&2
        fi
        read -r value || echo ""
    fi

    # If optional and empty, return empty
    if [[ "$optional" == "true" && -z "$value" ]]; then
        echo ""
        return
    fi

    # Basic validation - if it fails, just return what we got rather than loop
    if [[ -n "$value" && ! $value =~ $validation ]]; then
        echo "⚠ Warning: $value may not be a valid format" >&2
    fi

    echo "$value"
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
        local box_width=60
        local name_display="${existing_name:-"(not set)"}"
        local email_display="${existing_email:-"(not set)"}"

        # Format lines with proper spacing
        local name_line=$(printf "│ Username: %-*s │" $((box_width-15)) "$name_display")
        local email_line=$(printf "│ Email:    %-*s │" $((box_width-15)) "$email_display")

        echo "┌─────────────────────────────────────────────────────────┐"
        echo "│                 📋 Current Git Config                   │"
        echo "├─────────────────────────────────────────────────────────┤"
        echo "$name_line"
        echo "$email_line"
        echo "└─────────────────────────────────────────────────────────┘"
        echo ""

        # Use gum confirm if available, otherwise fallback
        if command -v gum &>/dev/null; then
            if gum confirm "Would you like to use these credentials?"; then
                echo "✅ Using existing credentials!"
                ARCHRIOT_USER_NAME="$existing_name"
                ARCHRIOT_USER_EMAIL="$existing_email"
            else
                echo ""
                echo "💬 No problem! Let's set up new credentials..."
                echo ""
                ARCHRIOT_USER_NAME=$(get_input "Name" "Your full name for Git commits" "^[a-zA-Z].*" "false")
                ARCHRIOT_USER_EMAIL=$(get_input "Email" "Your email for Git commits" "^[^@]+@[^@]+\.[^@]+$" "false")
            fi
        else
            echo -n "Would you like to use these credentials? [Y/n]: "
            read -r response
            case "$response" in
                [nN][oO]|[nN])
                    echo ""
                    echo "💬 No problem! Let's set up new credentials..."
                    echo ""
                    ARCHRIOT_USER_NAME=$(get_input "Name" "Your full name for Git commits" "^[a-zA-Z].*" "false")
                    ARCHRIOT_USER_EMAIL=$(get_input "Email" "Your email for Git commits" "^[^@]+@[^@]+\.[^@]+$" "false")
                    ;;
                *)
                    echo ""
                    echo "✅ Using existing credentials!"
                    ARCHRIOT_USER_NAME="$existing_name"
                    ARCHRIOT_USER_EMAIL="$existing_email"
                    ;;
            esac
        fi
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
