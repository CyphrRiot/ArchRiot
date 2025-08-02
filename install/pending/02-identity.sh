#!/bin/bash

# Apply Git configuration from environment variables
# This runs after the interactive prompt in install.sh
# NO INTERACTIVE CODE - just applies config

echo "🔧 Applying Git configuration..."

# Load user environment if it exists
env_file="$HOME/.config/archriot/user.env"
[[ -f "$env_file" ]] && source "$env_file"

# Apply Git configuration if variables are set
if [[ -n "${ARCHRIOT_USER_NAME// /}" ]]; then
    git config --global user.name "$ARCHRIOT_USER_NAME"
    echo "✓ Git user.name set to: $ARCHRIOT_USER_NAME"
fi

if [[ -n "${ARCHRIOT_USER_EMAIL// /}" ]]; then
    git config --global user.email "$ARCHRIOT_USER_EMAIL"
    echo "✓ Git user.email set to: $ARCHRIOT_USER_EMAIL"
fi

# Set useful Git aliases and defaults
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global pull.rebase true
git config --global init.defaultBranch master

echo "✓ Git configuration applied"
