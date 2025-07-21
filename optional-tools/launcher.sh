#!/bin/bash

# ArchRiot Optional Tools Launcher
# Provides access to advanced system tools
# Version: 1.0.0

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

print_header() {
    clear
    echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${PURPLE}║${NC}  ${CYAN}🔧 ArchRiot Optional Tools Launcher 🔧${NC}                  ${PURPLE}║${NC}"
    echo -e "${PURPLE}║${NC}  Advanced tools for power users                            ${PURPLE}║${NC}"
    echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

show_warning() {
    echo -e "${RED}⚠️  IMPORTANT WARNING ⚠️${NC}"
    echo
    echo -e "${YELLOW}These are ADVANCED tools that can modify critical system components!${NC}"
    echo
    echo "• Tools are NOT part of standard ArchRiot installation"
    echo "• They can potentially break your system if used incorrectly"
    echo "• Always have a backup and recovery plan ready"
    echo "• Test in virtual machines when possible"
    echo
    print_warning "Only proceed if you understand the risks and requirements"
    echo
    echo -n "Do you understand and accept these risks? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Launcher cancelled for safety."
        exit 0
    fi
    echo
}

check_tool_availability() {
    local tool_path="$1"
    local tool_name="$2"

    if [[ -x "$tool_path" ]]; then
        print_success "$tool_name available"
        return 0
    elif [[ -f "$tool_path" ]]; then
        print_warning "$tool_name found but not executable"
        print_info "Making executable..."
        chmod +x "$tool_path"
        print_success "$tool_name is now executable"
        return 0
    else
        print_error "$tool_name not found at $tool_path"
        return 1
    fi
}

show_tool_info() {
    echo -e "${CYAN}Available Tools:${NC}"
    echo

    # Dell XPS Sleep Fix
    echo -e "${GREEN}1. Dell XPS Sleep Crash Fix${NC}"
    echo "   Purpose: Fix sleep/suspend crashes on Dell XPS with Intel Arc Graphics"
    echo "   Risk Level: ⚠️  MODERATE - Modifies kernel parameters and power management"
    echo "   Compatibility: Dell XPS laptops with Intel Lunar Lake Arc Graphics (130V/140V)"
    echo "   Requirements: Dell XPS hardware, XE graphics driver"

    if check_tool_availability "$SCRIPT_DIR/dell-sleep-fix/setup-dell-sleep-fix.sh" "Dell Sleep Fix" >/dev/null 2>&1; then
        echo -e "   Status: ${GREEN}✓ Available${NC}"
    else
        echo -e "   Status: ${RED}✗ Not Available${NC}"
    fi
    echo

    # Secure Boot Setup
    echo -e "${GREEN}2. Secure Boot Setup${NC}"
    echo "   Purpose: Implement UEFI Secure Boot using standard Arch methods"
    echo "   Risk Level: ⚠️  MODERATE - Uses tested Arch packages"
    echo "   Compatibility: AMD, Intel, Any UEFI system"
    echo "   Requirements: UEFI boot mode, internet connection"

    if check_tool_availability "$SCRIPT_DIR/secure-boot/setup-secure-boot.sh" "Secure Boot Setup" >/dev/null 2>&1; then
        echo -e "   Status: ${GREEN}✓ Available${NC}"
    else
        echo -e "   Status: ${RED}✗ Not Available${NC}"
    fi
    echo

    # Future tools can be added here
    echo -e "${BLUE}3. More tools coming soon...${NC}"
    echo "   Additional optional tools will be added in future releases"
    echo
}

run_dell_sleep_fix() {
    local tool_path="$SCRIPT_DIR/dell-sleep-fix/setup-dell-sleep-fix.sh"

    if check_tool_availability "$tool_path" "Dell Sleep Fix"; then
        echo
        print_info "Starting Dell XPS Sleep Crash Fix..."
        echo
        exec sudo "$tool_path"
    else
        print_error "Dell Sleep Fix tool is not available"
        return 1
    fi
}

run_secure_boot_setup() {
    local tool_path="$SCRIPT_DIR/secure-boot/setup-secure-boot.sh"

    if check_tool_availability "$tool_path" "Secure Boot Setup"; then
        echo
        print_info "Starting Secure Boot Setup..."
        echo
        exec "$tool_path"
    else
        print_error "Secure Boot Setup tool is not available"
        return 1
    fi
}

show_main_menu() {
    while true; do
        print_header
        show_tool_info

        echo -e "${CYAN}Options:${NC}"
        echo "1. Launch Dell XPS Sleep Crash Fix"
        echo "2. Launch Secure Boot Setup"
        echo "3. View Documentation"
        echo "4. Check System Requirements"
        echo "5. Exit"
        echo
        echo -n "Enter choice [1-5]: "
        read -r choice

        case $choice in
            1)
                run_dell_sleep_fix
                break
                ;;
            2)
                run_secure_boot_setup
                break
                ;;
            3)
                show_documentation
                ;;
            4)
                check_system_requirements
                ;;
            5)
                echo "Goodbye!"
                exit 0
                ;;
            *)
                print_error "Invalid choice. Please try again."
                sleep 2
                ;;
        esac
    done
}

show_documentation() {
    clear
    echo -e "${PURPLE}=== Documentation ===${NC}"
    echo

    if [[ -f "$SCRIPT_DIR/README.md" ]]; then
        print_info "README.md found - showing key information:"
        echo
        head -n 50 "$SCRIPT_DIR/README.md"
        echo
        print_info "For complete documentation, read: $SCRIPT_DIR/README.md"
    else
        print_warning "README.md not found"
        echo
        print_info "For documentation, visit:"
        echo "• Arch Wiki: https://wiki.archlinux.org/title/Unified_Extensible_Firmware_Interface/Secure_Boot"
        echo "• ArchRiot GitHub: https://github.com/CyphrRiot/ArchRiot"
    fi

    echo
    echo -n "Press Enter to continue..."
    read -r
}

check_system_requirements() {
    clear
    echo -e "${PURPLE}=== System Requirements Check ===${NC}"
    echo

    # Check if running Arch Linux
    if [[ -f /etc/arch-release ]]; then
        print_success "Running Arch Linux"
    else
        print_error "Not running Arch Linux - tools designed for Arch only"
    fi

    # Check UEFI mode
    if [[ -d /sys/firmware/efi ]]; then
        print_success "System booted in UEFI mode"
    else
        print_error "System not in UEFI mode - required for Secure Boot"
    fi

    # Check internet connectivity
    if ping -c 1 archlinux.org &> /dev/null; then
        print_success "Internet connection available"
    else
        print_warning "No internet connection - required for package installation"
    fi

    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        print_warning "Running as root - tools should be run as regular user"
    else
        print_success "Running as regular user (recommended)"
    fi

    # Check for Dell XPS hardware (for Dell sleep fix)
    if lspci | grep -q "Intel Corporation Lunar Lake.*Graphics"; then
        print_success "Intel Lunar Lake Arc Graphics detected (Dell sleep fix compatible)"
    else
        print_info "Intel Lunar Lake Arc Graphics not detected (Dell sleep fix not needed)"
    fi

    # Check sudo access
    if sudo -n true 2>/dev/null; then
        print_success "Sudo access available"
    else
        print_warning "Sudo access may be required for some operations"
    fi

    echo
    echo -n "Press Enter to continue..."
    read -r
}

check_prerequisites() {
    # Check if we're in the right directory
    if [[ ! -f "$SCRIPT_DIR/README.md" ]]; then
        print_warning "Optional tools directory structure may be incomplete"
    fi

    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        print_error "Do not run the launcher as root!"
        print_info "Tools will use sudo when needed."
        exit 1
    fi
}

main() {
    print_header

    # Preliminary checks
    check_prerequisites

    # Show safety warning
    show_warning

    # Show main menu
    show_main_menu
}

# Run main function
main "$@"
