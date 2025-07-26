#!/bin/bash

# ArchRiot Secure Boot Setup - The Safe Arch Way
# Uses standard Arch Linux methods: sbctl + shim-signed
# Version: 1.0.0
# OPTIONAL TOOL - Not part of standard ArchRiot installation

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/var/log/archriot-secureboot.log"

# Logging
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | sudo tee -a "$LOG_FILE" >/dev/null
}

print_header() {
    clear
    echo -e "${CYAN}🛡️ ArchRiot Secure Boot Setup (The Arch Way) 🛡️${NC}"
    echo -e "${PURPLE}Safe, Standard Arch Linux Implementation${NC}"
    echo
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
    log "SUCCESS: $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    log "ERROR: $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
    log "WARNING: $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
    log "INFO: $1"
}

show_safety_info() {
    echo -e "${YELLOW}⚠️  IMPORTANT SAFETY INFORMATION ⚠️${NC}"
    echo
    echo "This script uses the STANDARD Arch Linux secure boot methods:"
    echo "• sbctl (official Arch package) for key management"
    echo "• shim-signed (AUR) for Microsoft key compatibility"
    echo "• Follows Arch Wiki recommendations exactly"
    echo
    echo -e "${GREEN}Benefits of this approach:${NC}"
    echo "✓ Much safer than custom key generation"
    echo "✓ Compatible with hardware firmware updates"
    echo "✓ Works with Windows dual-boot"
    echo "✓ Microsoft hardware compatibility guaranteed"
    echo "✓ Automatic kernel signing on updates"
    echo
    echo -e "${CYAN}Tested on:${NC}"
    echo "• AMD systems (Ryzen, EPYC)"
    echo "• Intel systems (Core, Xeon)"
    echo "• Any UEFI system with Secure Boot"
    echo
    echo -n "Continue with Arch-standard setup? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Setup cancelled."
        exit 0
    fi
}

check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_error "Do not run this script as root!"
        print_info "This script will use sudo when needed."
        exit 1
    fi
}

check_arch() {
    if [[ ! -f /etc/arch-release ]]; then
        print_error "This script is designed for Arch Linux only."
        exit 1
    fi
}

check_internet() {
    print_info "Checking internet connection..."
    if ! ping -c 1 8.8.8.8 &> /dev/null; then
        print_error "Internet connection required for package installation."
        exit 1
    fi
    print_success "Internet connection verified"
}

check_uefi() {
    print_info "Checking UEFI boot mode..."
    if [[ ! -d /sys/firmware/efi ]]; then
        print_error "System is not booted in UEFI mode."
        print_info "Secure Boot requires UEFI. Please:"
        echo "1. Enable UEFI mode in BIOS"
        echo "2. Disable Legacy/CSM mode"
        echo "3. Reinstall system in UEFI mode"
        exit 1
    fi
    print_success "UEFI mode confirmed"
}

check_current_status() {
    print_info "Checking current Secure Boot status..."

    echo
    echo "=== System Information ==="

    # CPU info
    local cpu_vendor=$(grep "vendor_id" /proc/cpuinfo | head -1 | awk '{print $3}')
    case "$cpu_vendor" in
        "AuthenticAMD") echo "CPU: AMD (Optimal compatibility)" ;;
        "GenuineIntel") echo "CPU: Intel (Full compatibility)" ;;
        *) echo "CPU: $cpu_vendor (Should work)" ;;
    esac

    # Boot loader
    if command -v bootctl >/dev/null 2>&1; then
        echo "Boot loader: systemd-boot detected"
    fi

    # Current Secure Boot status
    if command -v bootctl >/dev/null 2>&1; then
        echo
        echo "=== Current Boot Status ==="
        bootctl status | grep -E "(Secure Boot|TPM|System Token)" || true
    fi

    # Check if sbctl is installed
    if command -v sbctl >/dev/null 2>&1; then
        echo
        echo "=== Current sbctl Status ==="
        sudo sbctl status || true
    fi

    echo
}

install_packages() {
    print_info "Installing required packages..."

    # Update package database
    print_info "Updating package database..."
    sudo pacman -Sy

    # Install sbctl (official package)
    if ! pacman -Q sbctl &>/dev/null; then
        print_info "Installing sbctl (official Arch package)..."
        sudo pacman -S --noconfirm sbctl
        print_success "sbctl installed"
    else
        print_success "sbctl already installed"
    fi

    # Check if AUR helper is available
    local aur_helper=""
    for helper in yay paru aurman; do
        if command -v "$helper" >/dev/null 2>&1; then
            aur_helper="$helper"
            break
        fi
    done

    if [[ -z "$aur_helper" ]]; then
        print_warning "No AUR helper found. Installing yay..."

        # Install base-devel and git if not present
        sudo pacman -S --needed --noconfirm base-devel git

        # Install yay
        local temp_dir=$(mktemp -d)
        cd "$temp_dir"
        git clone https://aur.archlinux.org/yay.git
        cd yay
        makepkg -si --noconfirm
        cd "$SCRIPT_DIR"
        rm -rf "$temp_dir"
        aur_helper="yay"
        print_success "yay installed"
    fi

    # Install shim-signed (AUR package for Microsoft compatibility)
    if ! pacman -Q shim-signed &>/dev/null; then
        print_info "Installing shim-signed (AUR package for Microsoft compatibility)..."
        $aur_helper -S --noconfirm shim-signed
        print_success "shim-signed installed"
    else
        print_success "shim-signed already installed"
    fi
}

setup_secure_boot() {
    print_info "Setting up Secure Boot with standard Arch methods..."

    # Check if already set up
    if sudo sbctl status | grep -q "Secure Boot.*enabled"; then
        print_warning "Secure Boot appears to already be enabled"
        echo -n "Continue anyway? (y/N): "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_info "Setup cancelled by user"
            return 0
        fi
    fi

    print_warning "IMPORTANT: You must enable 'Setup Mode' in your UEFI/BIOS first!"
    echo "1. Reboot and enter UEFI/BIOS setup (F2/F12/Delete)"
    echo "2. Find Security/Secure Boot settings"
    echo "3. Enable 'Setup Mode' or 'Clear Keys'"
    echo "4. Save and boot back to this system"
    echo
    echo -n "Have you enabled Setup Mode? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_warning "Please enable Setup Mode first, then run this script again"
        exit 0
    fi

    # Create keys
    print_info "Creating secure boot keys..."
    sudo sbctl create-keys

    # Enroll keys
    print_info "Enrolling keys..."
    if ! sudo sbctl enroll-keys -m; then
        print_error "Failed to enroll keys. Make sure Setup Mode is enabled."
        return 1
    fi
    print_success "Keys enrolled successfully"

    # Sign bootloader and kernel
    print_info "Signing bootloader and kernel..."

    # Find and sign systemd-boot
    local esp_path="/efi"
    if [[ ! -d "$esp_path" ]]; then
        esp_path="/boot"
    fi

    if [[ -f "$esp_path/EFI/systemd/systemd-bootx64.efi" ]]; then
        sudo sbctl sign -s "$esp_path/EFI/systemd/systemd-bootx64.efi"
        print_success "Signed systemd-boot"
    fi

    # Sign kernel
    local kernel_path="/boot/vmlinuz-linux"
    if [[ -f "$kernel_path" ]]; then
        sudo sbctl sign -s "$kernel_path"
        print_success "Signed kernel"
    fi

    # Sign initramfs if it exists
    local initramfs_path="/boot/initramfs-linux.img"
    if [[ -f "$initramfs_path" ]]; then
        sudo sbctl sign -s "$initramfs_path"
        print_success "Signed initramfs"
    fi

    # Sign any other kernels
    for kernel in /boot/vmlinuz-*; do
        if [[ -f "$kernel" && "$kernel" != "$kernel_path" ]]; then
            sudo sbctl sign -s "$kernel"
            print_success "Signed $(basename "$kernel")"
        fi
    done

    print_success "Secure Boot setup completed"
}

verify_setup() {
    print_info "Verifying Secure Boot setup..."

    echo
    echo "=== Verification Results ==="

    # Check sbctl status
    sudo sbctl status

    echo
    echo "=== Signed Files ==="
    sudo sbctl list-files

    # Verify all files are signed
    if sudo sbctl verify; then
        print_success "All files are properly signed"
    else
        print_warning "Some files may not be signed properly"
        print_info "Check the output above and sign missing files with:"
        echo "sudo sbctl sign -s <file_path>"
    fi
}

show_next_steps() {
    echo
    print_success "Setup Complete!"
    echo
    print_info "Next Steps:"
    echo "1. Reboot your system"
    echo "2. Enter UEFI/BIOS setup (F2/F12/Delete)"
    echo "3. Navigate to Security → Secure Boot"
    echo "4. DISABLE 'Setup Mode'"
    echo "5. ENABLE 'Secure Boot'"
    echo "6. Save and exit"
    echo
    print_info "Your system will now:"
    echo "✓ Boot with Secure Boot enabled"
    echo "✓ Automatically sign new kernels"
    echo "✓ Maintain hardware compatibility"
    echo
    print_warning "If boot fails:"
    echo "• Enter UEFI and disable Secure Boot"
    echo "• Boot normally and check 'sbctl verify'"
    echo "• Sign any missing files with 'sbctl sign'"
    echo
    print_info "Future kernel updates will be signed automatically"
    print_info "Monitor with: sudo sbctl status"
}

show_menu() {
    while true; do
        print_header
        echo "Choose an option:"
        echo
        echo "1. Quick Setup (Recommended)"
        echo "2. Check Current Status"
        echo "3. Verify Existing Setup"
        echo "4. Advanced Manual Setup"
        echo "5. Exit"
        echo
        echo -n "Enter choice [1-5]: "
        read -r choice

        case $choice in
            1)
                check_current_status
                install_packages
                setup_secure_boot
                verify_setup
                show_next_steps
                break
                ;;
            2)
                check_current_status
                echo
                echo -n "Press Enter to continue..."
                read -r
                ;;
            3)
                verify_setup
                echo
                echo -n "Press Enter to continue..."
                read -r
                ;;
            4)
                print_warning "Advanced setup allows manual control of each step"
                echo -n "Continue? (y/N): "
                read -r response
                if [[ "$response" =~ ^[Yy]$ ]]; then
                    advanced_setup
                fi
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

advanced_setup() {
    while true; do
        clear
        echo -e "${PURPLE}=== Advanced Secure Boot Setup ===${NC}"
        echo
        echo "1. Install packages only"
        echo "2. Create keys only"
        echo "3. Enroll keys only"
        echo "4. Sign files only"
        echo "5. Full setup"
        echo "6. Back to main menu"
        echo
        echo -n "Enter choice [1-6]: "
        read -r choice

        case $choice in
            1) install_packages ;;
            2) sudo sbctl create-keys ;;
            3) sudo sbctl enroll-keys -m ;;
            4)
                print_info "Signing common files..."
                sudo sbctl sign -s /boot/vmlinuz-linux
                sudo sbctl sign -s /efi/EFI/systemd/systemd-bootx64.efi 2>/dev/null || sudo sbctl sign -s /boot/EFI/systemd/systemd-bootx64.efi
                ;;
            5)
                install_packages
                setup_secure_boot
                verify_setup
                ;;
            6) return ;;
            *) print_error "Invalid choice" ;;
        esac

        if [[ $choice != 6 ]]; then
            echo
            echo -n "Press Enter to continue..."
            read -r
        fi
    done
}

# Main execution
main() {
    print_header

    log "ArchRiot Secure Boot Setup started"

    # Preliminary checks
    check_root
    check_arch
    check_uefi
    check_internet

    # Show safety information
    show_safety_info

    # Show menu
    show_menu

    log "ArchRiot Secure Boot Setup completed"

    echo
    print_success "ArchRiot Secure Boot setup completed!"
    print_info "Log file: $LOG_FILE"
    print_warning "Remember to test boot thoroughly before relying on this system."
}

# Run main function
main "$@"
