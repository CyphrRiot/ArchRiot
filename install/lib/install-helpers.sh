#!/bin/bash

# =============================================================================
# OhmArchy Installation Helper Library
# Centralized error handling and package installation functions
# =============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Installation context tracking
CURRENT_INSTALLER=""
INSTALL_START_TIME=""
FAILED_PACKAGES=()

# Initialize installer context
init_installer() {
    local installer_name="$1"
    CURRENT_INSTALLER="$installer_name"
    INSTALL_START_TIME=$(date +%s)
    echo -e "${BLUE}🔧 Initializing: $installer_name${NC}"
}

# Enhanced package installation with error handling
install_packages() {
    local packages="$1"
    local package_type="${2:-essential}"

    echo -e "${CYAN}📦 Installing $package_type packages: $packages${NC}"

    if yay -S --noconfirm --needed $packages; then
        echo -e "${GREEN}✓ Successfully installed: $packages${NC}"
        return 0
    else
        local exit_code=$?
        echo -e "${RED}❌ Failed to install: $packages${NC}"
        FAILED_PACKAGES+=("$packages")

        case $package_type in
            "essential"|"critical")
                handle_critical_failure "$packages" "$exit_code"
                ;;
            "optional")
                handle_optional_failure "$packages" "$exit_code"
                return 0
                ;;
            *)
                handle_critical_failure "$packages" "$exit_code"
                ;;
        esac

        return $exit_code
    fi
}

# Handle critical package installation failures
handle_critical_failure() {
    local failed_packages="$1"
    local exit_code="$2"

    echo
    echo -e "${RED}❌ CRITICAL INSTALLATION FAILURE${NC}"
    echo "=================================="
    echo "Installer: $CURRENT_INSTALLER"
    echo "Failed packages: $failed_packages"
    echo "Exit code: $exit_code"
    echo "Time: $(date)"
    echo

    # Log failure details
    log_failure_details "$failed_packages" "$exit_code"

    # Show context-specific troubleshooting
    show_troubleshooting "$CURRENT_INSTALLER" "$failed_packages"

    echo -e "${YELLOW}💡 You can retry the installation after fixing the issue:${NC}"
    echo "   source ~/.local/share/archriot/install.sh"
    echo

    exit $exit_code
}

# Handle optional package failures (warn but continue)
handle_optional_failure() {
    local failed_packages="$1"
    local exit_code="$2"

    echo -e "${YELLOW}⚠ Optional package installation failed: $failed_packages${NC}"
    echo -e "${YELLOW}  This may reduce functionality but installation will continue${NC}"

    # Still log for debugging
    echo "OPTIONAL FAILURE: $CURRENT_INSTALLER - $failed_packages (exit: $exit_code)" >> /tmp/archriot-warnings.log
}

# Log detailed failure information
log_failure_details() {
    local failed_packages="$1"
    local exit_code="$2"
    local log_file="/tmp/archriot-install-failure.log"

    {
        echo "===== FAILURE DETAILS ====="
        echo "Timestamp: $(date)"
        echo "Installer: $CURRENT_INSTALLER"
        echo "Failed packages: $failed_packages"
        echo "Exit code: $exit_code"
        echo "Duration: $(($(date +%s) - INSTALL_START_TIME))s"
        echo "System info:"
        echo "  - Architecture: $(uname -m)"
        echo "  - Kernel: $(uname -r)"
        echo "  - Available space: $(df -h / | tail -1 | awk '{print $4}')"
        echo "  - Memory: $(free -h | head -2 | tail -1 | awk '{print $7}')"
        echo "Failed packages list:"
        printf '  - %s\n' "${FAILED_PACKAGES[@]}"
        echo "=============================="
        echo
    } >> "$log_file"
}

# Context-aware troubleshooting suggestions
show_troubleshooting() {
    local installer="$1"
    local failed_packages="$2"

    echo -e "${CYAN}🔍 TROUBLESHOOTING STEPS:${NC}"
    echo "1. Check internet connection: ping -c3 archlinux.org"
    echo "2. Update system: sudo pacman -Syu"
    echo "3. Clear package cache: yay -Sc"
    echo "4. Check disk space: df -h"
    echo

    echo -e "${CYAN}📋 CONTEXT-SPECIFIC SOLUTIONS:${NC}"
    case "$installer" in
        *terminal*|*development*|*xtras*)
            echo "• AUR build failures:"
            echo "  - Wait 30 minutes and retry (packages may be updating)"
            echo "  - Install build dependencies: sudo pacman -S base-devel"
            echo "  - Check if package exists: yay -Ss <package-name>"
            ;;
        *desktop*|*hyprlandia*)
            echo "• Desktop environment conflicts:"
            echo "  - Remove conflicting WMs: sudo pacman -Rs gnome kde-plasma"
            echo "  - Install GPU drivers first"
            echo "  - Check graphics: lspci | grep -i vga"
            ;;
        *network*|*bluetooth*)
            echo "• Hardware/driver issues:"
            echo "  - Check hardware: lspci | lsusb"
            echo "  - Install firmware: sudo pacman -S linux-firmware"
            echo "  - Load modules: sudo modprobe <module-name>"
            ;;
        *fonts*|*mimetypes*)
            echo "• Configuration issues:"
            echo "  - Clear font cache: fc-cache -fv"
            echo "  - Update desktop database: update-desktop-database"
            ;;
        *)
            echo "• General package issues:"
            echo "  - Check package conflicts: yay -Si <package-name>"
            echo "  - Search alternatives: yay -Ss <similar-name>"
            echo "  - Check Arch forums/wiki"
            ;;
    esac

    echo
    echo -e "${CYAN}📝 Additional debugging:${NC}"
    echo "• Full logs: journalctl -xe"
    echo "• Package manager logs: tail -f /var/log/pacman.log"
    echo "• Failure details: cat /tmp/archriot-install-failure.log"
}

# Validate that packages were actually installed correctly
validate_packages() {
    local packages="$1"
    local package_type="${2:-essential}"
    local failed_validation=()

    echo -e "${BLUE}🔍 Validating installation: $packages${NC}"

    for package in $packages; do
        if ! pacman -Qi "$package" >/dev/null 2>&1; then
            failed_validation+=("$package")
        fi
    done

    if [ ${#failed_validation[@]} -gt 0 ]; then
        echo -e "${YELLOW}⚠ Validation failed for: ${failed_validation[*]}${NC}"
        if [ "$package_type" = "essential" ] || [ "$package_type" = "critical" ]; then
            echo -e "${RED}Critical packages missing after installation!${NC}"
            return 1
        fi
    else
        echo -e "${GREEN}✓ All packages validated successfully${NC}"
    fi

    return 0
}

# Quick wrapper for essential packages (most common use case)
install_essential() {
    install_packages "$1" "essential"
}

# Quick wrapper for optional packages
install_optional() {
    install_packages "$1" "optional"
}

# Clean up temporary files on successful completion
cleanup_install_files() {
    rm -f /tmp/archriot-install-failure.log /tmp/archriot-warnings.log 2>/dev/null || true
}

# Show installation summary
show_install_summary() {
    local duration=$(($(date +%s) - INSTALL_START_TIME))
    echo
    echo -e "${GREEN}✓ $CURRENT_INSTALLER completed successfully${NC}"
    echo -e "${BLUE}Duration: ${duration}s${NC}"

    if [ -f /tmp/archriot-warnings.log ]; then
        local warning_count=$(wc -l < /tmp/archriot-warnings.log)
        if [ $warning_count -gt 0 ]; then
            echo -e "${YELLOW}⚠ $warning_count warnings (check /tmp/archriot-warnings.log)${NC}"
        fi
    fi
}

# Export all functions for use in other scripts
export -f init_installer install_packages install_essential install_optional
export -f validate_packages cleanup_install_files show_install_summary
export -f handle_critical_failure handle_optional_failure
export -f log_failure_details show_troubleshooting
