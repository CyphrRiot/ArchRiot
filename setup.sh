#!/bin/bash

# ArchRiot Setup - Bulletproof Minimal Version
# Purpose: Download repository and run installer (nothing else!)

set -euo pipefail  # Exit on any error, undefined vars, or pipe failures

# Colors for output
readonly PURPLE='\033[0;35m'
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly NC='\033[0m'

# Constants with safety checks
readonly HOME="${HOME:-$(getent passwd "$(whoami)" | cut -d: -f6)}"
readonly REPO_URL="https://github.com/CyphrRiot/ArchRiot.git"
readonly INSTALL_DIR="$HOME/.local/share/archriot"
readonly INSTALLER_PATH="$INSTALL_DIR/install/archriot"

# Error handler
error_exit() {
    echo -e "${RED}❌ Error: $1${NC}" >&2
    exit 1
}

# Success message
success_msg() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Info message
info_msg() {
    echo -e "${PURPLE}🔄 $1${NC}"
}

# Verify prerequisites
check_prerequisites() {
    info_msg "Checking prerequisites..."

    # Verify we have a home directory
    [[ -n "$HOME" && -d "$HOME" ]] || error_exit "No valid home directory found"

    # Check if we have git
    if ! command -v git >/dev/null 2>&1; then
        info_msg "Installing git..."

        # Check for package manager and sudo
        command -v pacman >/dev/null 2>&1 || error_exit "pacman not found - are you on Arch Linux?"
        command -v sudo >/dev/null 2>&1 || error_exit "sudo not found - please install git manually"

        sudo pacman -Sy --noconfirm --needed git || error_exit "Failed to install git"
    fi

    # Test git functionality with actual repository
    if ! git ls-remote --exit-code --heads "$REPO_URL" >/dev/null 2>&1; then
        error_exit "Cannot connect to ArchRiot repository - check network connection"
    fi

    success_msg "Prerequisites verified"
}

# Download or update repository
setup_repository() {
    info_msg "Setting up ArchRiot repository..."

    if [[ -d "$INSTALL_DIR/.git" ]]; then
        info_msg "Updating existing installation..."

        # Safely update in subshell to avoid directory issues
        (
            cd "$INSTALL_DIR" || error_exit "Cannot access $INSTALL_DIR"
            git fetch origin || error_exit "Failed to fetch updates"
            git reset --hard origin/master || error_exit "Failed to update repository"
        )

        success_msg "Repository updated"
    else
        info_msg "Fresh installation..."

        # Remove any non-git directory and create parent
        [[ -d "$INSTALL_DIR" ]] && rm -rf "$INSTALL_DIR"
        mkdir -p "$(dirname "$INSTALL_DIR")" || error_exit "Cannot create directory structure"

        # Clone repository
        git clone --depth 1 "$REPO_URL" "$INSTALL_DIR" || error_exit "Failed to clone repository"

        success_msg "Repository cloned"
    fi
}

# Verify installer
verify_installer() {
    info_msg "Verifying installer..."

    [[ -f "$INSTALLER_PATH" ]] || error_exit "Installer binary not found at $INSTALLER_PATH"
    [[ -x "$INSTALLER_PATH" ]] || error_exit "Installer binary is not executable"

    # Test installer responds
    "$INSTALLER_PATH" --version >/dev/null 2>&1 || error_exit "Installer binary failed basic test"

    success_msg "Installer verified"
}

# Main execution
main() {
    echo -e "${PURPLE}"
    echo ' █████╗ ██████╗  ██████╗██╗  ██╗██████╗ ██╗ ██████╗ ████████╗'
    echo '██╔══██╗██╔══██╗██╔════╝██║  ██║██╔══██╗██║██╔═══██╗╚══██╔══╝'
    echo '███████║██████╔╝██║     ███████║██████╔╝██║██║   ██║   ██║   '
    echo '██╔══██║██╔══██╗██║     ██╔══██║██╔══██╗██║██║   ██║   ██║   '
    echo '██║  ██║██║  ██║╚██████╗██║  ██║██║  ██║██║╚██████╔╝   ██║   '
    echo '╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝ ╚═════╝    ╚═╝   '
    echo -e "${NC}"
    echo
    echo -e "${PURPLE}🎭 ArchRiot Setup${NC}"
    echo -e "${PURPLE}=====================${NC}"
    echo

    # Execute setup steps
    check_prerequisites
    setup_repository
    verify_installer

    echo
    info_msg "Starting ArchRiot installer..."
    echo

    # Hand off to the real installer
    exec "$INSTALLER_PATH"
}

# Run main function
main "$@"
