#!/bin/bash

# =============================================================================
# OhmArchy Shared Library
# Central library that loads all common functions for installation scripts
# =============================================================================

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the install helpers library
source "$SCRIPT_DIR/install-helpers.sh" || {
    echo "❌ Failed to load install helpers library"
    exit 1
}

# Source the sudo helper library
source "$SCRIPT_DIR/sudo-helper.sh" || {
    echo "❌ Failed to load sudo helper library"
    exit 1
}

# Common utility functions that might be used across scripts

# Load user environment configuration
load_user_environment() {
    local env_file="$HOME/.config/archriot/user.env"
    if [[ -f "$env_file" ]]; then
        source "$env_file"
        echo "✓ Loaded user environment from $env_file"
    else
        echo "⚠ User environment file not found: $env_file"
    fi
}

# Check if running on Arch Linux
check_arch_system() {
    if [[ ! -f /etc/arch-release ]]; then
        echo "❌ This script requires Arch Linux"
        exit 1
    fi
}

# Check if yay is available
check_yay_available() {
    if ! command -v yay >/dev/null 2>&1; then
        echo "❌ yay AUR helper is required but not found"
        echo "Please install yay first: https://github.com/Jguer/yay"
        exit 1
    fi
}

# Update system packages
update_system() {
    echo "🔄 Updating system packages..."
    if yay -Syu --noconfirm; then
        echo "✓ System updated successfully"
    else
        echo "⚠ System update failed, continuing anyway"
    fi
}

# Create directory if it doesn't exist
ensure_directory() {
    local dir="$1"
    if [[ ! -d "$dir" ]]; then
        mkdir -p "$dir"
        echo "✓ Created directory: $dir"
    fi
}

# Check available disk space (in GB)
check_disk_space() {
    local required_gb="${1:-5}"
    local available_gb=$(df / | tail -1 | awk '{print int($4/1024/1024)}')

    if [[ $available_gb -lt $required_gb ]]; then
        echo "❌ Insufficient disk space: ${available_gb}GB available, ${required_gb}GB required"
        return 1
    else
        echo "✓ Sufficient disk space: ${available_gb}GB available"
        return 0
    fi
}

# Export all functions for use in other scripts
export -f load_user_environment check_arch_system check_yay_available
export -f update_system ensure_directory check_disk_space
