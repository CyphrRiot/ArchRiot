#!/bin/bash

# ArchRiot Setup - Bulletproof Minimal Version
# Purpose: Download repository and run installer (nothing else!)

set -euo pipefail  # Exit on any error, undefined vars, or pipe failures

# Colors for output
readonly PURPLE='\033[0;35m'
readonly BLUE='\033[0;34m'
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly NC='\033[0m'

# Constants with safety checks
readonly HOME="${HOME:-$(getent passwd "$(whoami)" | cut -d: -f6)}"
readonly REPO_URL="https://github.com/CyphrRiot/ArchRiot/archive/refs/heads/master.tar.gz"
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

# Warning message
warn_msg() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

# Enhanced WiFi and network connectivity detection
check_network_connectivity() {
    info_msg "Checking network connectivity..."

    local connection_type=""
    local interface_found=false

    # Check for active network interfaces
    while read -r interface; do
        [[ -n "$interface" ]] || continue
        interface_found=true

        if [[ "$interface" =~ ^wl ]]; then
            connection_type="WiFi"
            info_msg "Found WiFi interface: $interface"

            # Check WiFi connection status
            if command -v iwconfig >/dev/null 2>&1; then
                local essid=$(iwconfig "$interface" 2>/dev/null | grep -o 'ESSID:"[^"]*"' | cut -d'"' -f2)
                if [[ -n "$essid" && "$essid" != "off/any" ]]; then
                    info_msg "Connected to WiFi network: $essid"
                else
                    warn_msg "WiFi interface found but not connected to a network"
                fi
            fi
        elif [[ "$interface" =~ ^en|^eth ]]; then
            connection_type="Ethernet"
            info_msg "Found Ethernet interface: $interface"
        fi
    done < <(ip route show default 2>/dev/null | awk '{print $5}' | sort -u)

    # If no interfaces found via routing, check all available interfaces
    if ! $interface_found; then
        info_msg "No default route found, checking all network interfaces..."

        while read -r line; do
            local iface=$(echo "$line" | awk '{print $2}' | tr -d ':')
            local status=$(echo "$line" | grep -o 'state [A-Z]*' | awk '{print $2}')

            if [[ "$iface" != "lo" && "$status" == "UP" ]]; then
                interface_found=true
                if [[ "$iface" =~ ^wl ]]; then
                    connection_type="WiFi"
                    info_msg "Found active WiFi interface: $iface"
                elif [[ "$iface" =~ ^en|^eth ]]; then
                    connection_type="Ethernet"
                    info_msg "Found active Ethernet interface: $iface"
                fi
            fi
        done < <(ip link show 2>/dev/null | grep -E "^[0-9]+:")
    fi

    # Test actual internet connectivity with multiple fallbacks
    info_msg "Testing internet connectivity..."

    local test_urls=(
        "8.8.8.8"           # Google DNS
        "1.1.1.1"           # Cloudflare DNS
        "github.com"        # GitHub
        "archlinux.org"     # Arch Linux
    )

    local connectivity_test_passed=false

    for url in "${test_urls[@]}"; do
        if ping -c 1 -W 3 "$url" >/dev/null 2>&1; then
            connectivity_test_passed=true
            break
        elif command -v curl >/dev/null 2>&1 && curl -s --connect-timeout 5 --max-time 10 "http://$url" >/dev/null 2>&1; then
            connectivity_test_passed=true
            break
        elif command -v wget >/dev/null 2>&1 && wget -q --timeout=5 --tries=1 --spider "http://$url" >/dev/null 2>&1; then
            connectivity_test_passed=true
            break
        fi
    done

    if ! $connectivity_test_passed; then
        error_exit "No internet connectivity detected. Please check your network connection and try again."
    fi

    # Additional WiFi-specific checks
    if [[ "$connection_type" == "WiFi" ]]; then
        # Check signal strength if iwconfig is available
        if command -v iwconfig >/dev/null 2>&1; then
            local wifi_interface=$(ip route show default 2>/dev/null | awk '/^default.*wl/ {print $5}' | head -1)
            if [[ -n "$wifi_interface" ]]; then
                local signal_info=$(iwconfig "$wifi_interface" 2>/dev/null | grep -E "(Signal level|Link Quality)")
                if [[ -n "$signal_info" ]]; then
                    info_msg "WiFi signal info: $(echo "$signal_info" | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g')"
                fi
            fi
        fi

        # Check for common WiFi issues
        if command -v networkctl >/dev/null 2>&1; then
            local wifi_status=$(networkctl status 2>/dev/null | grep -E "(State:|WiFi|Wireless)")
            if [[ -n "$wifi_status" ]]; then
                info_msg "Network status: $(echo "$wifi_status" | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g')"
            fi
        fi
    fi

    success_msg "Network connectivity verified ($connection_type connection)"
}

# Verify prerequisites
check_prerequisites() {
    info_msg "Checking prerequisites..."

    # Verify we have a home directory
    [[ -n "$HOME" && -d "$HOME" ]] || error_exit "No valid home directory found"

    # Check network connectivity first
    check_network_connectivity

    # Check if we have git
    if ! command -v git >/dev/null 2>&1; then
        info_msg "Installing git..."

        # Check for package manager and sudo
        command -v pacman >/dev/null 2>&1 || error_exit "pacman not found - are you on Arch Linux?"
        command -v sudo >/dev/null 2>&1 || error_exit "sudo not found - please install git manually"

        sudo pacman -Sy --noconfirm --needed git || error_exit "Failed to install git"
    fi

    # Test git functionality with actual repository
    info_msg "Testing repository connectivity..."
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

        # Remove any existing directory and create parent
        [[ -d "$INSTALL_DIR" ]] && rm -rf "$INSTALL_DIR"
        mkdir -p "$(dirname "$INSTALL_DIR")" || error_exit "Cannot create directory structure"

        # Download and extract release tarball (much faster than git clone)
        info_msg "Downloading latest release..."
        local temp_file="/tmp/archriot-latest.tar.gz"

        # Download tarball
        curl -fSL "$REPO_URL" -o "$temp_file" || error_exit "Failed to download repository"

        # Extract to temporary directory first
        local temp_dir="/tmp/archriot-extract"
        mkdir -p "$temp_dir" || error_exit "Cannot create temp directory"

        # Extract tarball (GitHub creates a subdirectory named ArchRiot-master)
        tar -xzf "$temp_file" -C "$temp_dir" || error_exit "Failed to extract repository"

        # Move from temp location to final location
        mv "$temp_dir/ArchRiot-master" "$INSTALL_DIR" || error_exit "Failed to move extracted files"

        # Cleanup
        rm -f "$temp_file"
        rm -rf "$temp_dir"

        success_msg "Repository downloaded and extracted (no git history)"
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
    echo -e "${BLUE}"
    echo '▄  ▄▀█ █▀█ █▀▀ █ █ █▀█ █ █▀█ ▀█▀  ▄'
    echo '▄  █▀█ █▀▄ █▄▄ █▀█ █▀▄ █ █▄█  █   ▄'
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
