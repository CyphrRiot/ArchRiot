#!/bin/bash

# ==============================================================================
# OhmArchy Sudo Helper - Safe Passwordless Setup
# ==============================================================================
# Temporarily enables passwordless sudo for package management during install
# Automatically reverts after installation for security
# ==============================================================================

SUDOERS_BACKUP="/tmp/omarchy-sudoers-backup-$$"
SUDOERS_TEMP="/tmp/omarchy-sudoers-temp-$$"
CURRENT_USER="$(whoami)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    local status="$1"
    local message="$2"
    case "$status" in
        "INFO") echo -e "${GREEN}ℹ${NC} $message" ;;
        "WARN") echo -e "${YELLOW}⚠${NC} $message" ;;
        "ERROR") echo -e "${RED}❌${NC} $message" ;;
    esac
}

# Check if user is already passwordless
is_passwordless() {
    sudo -n true 2>/dev/null
}

# Setup temporary passwordless sudo for package management only
setup_passwordless_sudo() {
    print_status "INFO" "Setting up temporary passwordless sudo for package management..."

    # Check if already passwordless
    if is_passwordless; then
        print_status "INFO" "Sudo is already passwordless - no changes needed"
        return 0
    fi

    # Backup current sudoers
    if ! sudo cp /etc/sudoers "$SUDOERS_BACKUP" 2>/dev/null; then
        print_status "ERROR" "Failed to backup sudoers file"
        return 1
    fi

    # Create temporary sudoers rule for package management
    cat > "$SUDOERS_TEMP" <<EOF
# Temporary OhmArchy installation rule - will be auto-removed
$CURRENT_USER ALL=(ALL) NOPASSWD: /usr/bin/pacman, /usr/bin/yay, /usr/bin/systemctl
EOF

    # Validate the temporary rule
    if ! sudo visudo -c -f "$SUDOERS_TEMP" >/dev/null 2>&1; then
        print_status "ERROR" "Invalid sudoers rule generated"
        return 1
    fi

    # Append the rule to sudoers
    if sudo sh -c "cat '$SUDOERS_TEMP' >> /etc/sudoers" 2>/dev/null; then
        print_status "INFO" "Temporary passwordless sudo enabled for package management"
        print_status "INFO" "Commands: pacman, yay, systemctl"
        rm -f "$SUDOERS_TEMP"
        return 0
    else
        print_status "ERROR" "Failed to update sudoers file"
        return 1
    fi
}

# Remove temporary passwordless sudo rules
cleanup_passwordless_sudo() {
    print_status "INFO" "Cleaning up temporary passwordless sudo..."

    # Check if backup exists
    if [[ ! -f "$SUDOERS_BACKUP" ]]; then
        print_status "WARN" "No sudoers backup found - manual cleanup may be needed"
        return 1
    fi

    # Restore original sudoers
    if sudo cp "$SUDOERS_BACKUP" /etc/sudoers 2>/dev/null; then
        print_status "INFO" "Original sudoers configuration restored"
        rm -f "$SUDOERS_BACKUP"
        return 0
    else
        print_status "ERROR" "Failed to restore sudoers - check /etc/sudoers manually"
        return 1
    fi
}

# Emergency cleanup function
emergency_cleanup() {
    print_status "WARN" "Emergency cleanup triggered"

    if [[ -f "$SUDOERS_BACKUP" ]]; then
        print_status "INFO" "Restoring sudoers from backup..."
        sudo cp "$SUDOERS_BACKUP" /etc/sudoers 2>/dev/null || {
            print_status "ERROR" "Emergency restore failed - manual intervention required"
            print_status "ERROR" "Backup location: $SUDOERS_BACKUP"
        }
    fi

    # Clean up temp files
    rm -f "$SUDOERS_TEMP" "$SUDOERS_BACKUP" 2>/dev/null
}

# Test if passwordless sudo is working for our commands
test_passwordless_sudo() {
    print_status "INFO" "Testing passwordless sudo setup..."

    local test_commands=("pacman --version" "systemctl --version")

    for cmd in "${test_commands[@]}"; do
        if sudo -n $cmd >/dev/null 2>&1; then
            print_status "INFO" "✓ Passwordless sudo working for: $cmd"
        else
            print_status "WARN" "✗ Passwordless sudo failed for: $cmd"
            return 1
        fi
    done

    return 0
}

# Show current sudo status
show_sudo_status() {
    print_status "INFO" "Current sudo status:"

    if is_passwordless; then
        print_status "INFO" "  • Passwordless: YES"
    else
        print_status "INFO" "  • Passwordless: NO"
    fi

    if groups | grep -q wheel; then
        print_status "INFO" "  • In wheel group: YES"
    else
        print_status "INFO" "  • In wheel group: NO"
    fi

    if sudo -l 2>/dev/null | grep -q NOPASSWD; then
        print_status "INFO" "  • Has NOPASSWD rules: YES"
    else
        print_status "INFO" "  • Has NOPASSWD rules: NO"
    fi
}

# Main function based on command
main() {
    local command="${1:-status}"

    case "$command" in
        "setup"|"enable")
            show_sudo_status
            setup_passwordless_sudo
            test_passwordless_sudo
            ;;
        "cleanup"|"disable"|"restore")
            cleanup_passwordless_sudo
            show_sudo_status
            ;;
        "test")
            test_passwordless_sudo
            ;;
        "status")
            show_sudo_status
            ;;
        "emergency")
            emergency_cleanup
            ;;
        *)
            echo "Usage: $0 {setup|cleanup|test|status|emergency}"
            echo ""
            echo "Commands:"
            echo "  setup     - Enable temporary passwordless sudo for package management"
            echo "  cleanup   - Restore original sudo configuration"
            echo "  test      - Test if passwordless sudo is working"
            echo "  status    - Show current sudo configuration"
            echo "  emergency - Emergency restore (use if installation fails)"
            echo ""
            echo "This tool temporarily enables passwordless sudo for pacman/yay/systemctl"
            echo "during OhmArchy installation, then automatically restores security."
            exit 1
            ;;
    esac
}

# Set up trap for emergency cleanup on script exit/failure
trap emergency_cleanup EXIT ERR

# Execute main function
main "$@"
