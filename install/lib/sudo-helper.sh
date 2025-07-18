#!/bin/bash

# ==============================================================================
# OhmArchy Sudo Helper - Simple and Safe
# ==============================================================================
# Ensures user is in wheel group and enables passwordless package management
# Uses surgical changes instead of dangerous file backups
# ==============================================================================

CURRENT_USER="$(whoami)"
OMARCHY_SUDO_RULE="$CURRENT_USER ALL=(ALL) NOPASSWD: /usr/bin/pacman, /usr/bin/yay, /usr/bin/systemctl"
OMARCHY_SUDO_MARKER="# OhmArchy temporary rule"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    local status="$1"
    local message="$2"
    case "$status" in
        "INFO") echo -e "${GREEN}ℹ ${NC} $message" ;;
        "WARN") echo -e "${YELLOW}⚠${NC} $message" ;;
        "ERROR") echo -e "${RED}❌${NC} $message" ;;
    esac
}

# Check if user is in wheel group
is_in_wheel() {
    groups "$CURRENT_USER" | grep -q '\bwheel\b'
}

# Check if passwordless sudo works for our commands
is_passwordless() {
    sudo -n pacman --version >/dev/null 2>&1 && \
    sudo -n yay --version >/dev/null 2>&1 && \
    sudo -n systemctl --version >/dev/null 2>&1
}

# Check if our rule is already active
has_archriot_rule() {
    sudo grep -q "$OMARCHY_SUDO_MARKER" /etc/sudoers 2>/dev/null && \
    sudo grep -q "$CURRENT_USER ALL=(ALL) NOPASSWD:" /etc/sudoers 2>/dev/null
}

# Add user to wheel group (for system consistency)
add_to_wheel() {
    print_status "INFO" "Adding $CURRENT_USER to wheel group..."
    if sudo usermod -aG wheel "$CURRENT_USER"; then
        print_status "INFO" "✓ User added to wheel group"
        return 0
    else
        print_status "WARN" "Failed to add user to wheel group (continuing with user-specific rule)"
        return 0
    fi
}

# Add passwordless rule for package management
add_passwordless_rule() {
    print_status "INFO" "Adding passwordless sudo rule for package management..."

    # Create the rule with marker
    local rule_line="$OMARCHY_SUDO_RULE $OMARCHY_SUDO_MARKER"

    # Create a temporary file for visudo
    local temp_file=$(mktemp)
    sudo cp /etc/sudoers "$temp_file"

    # Add our rule to the temp file
    echo "$rule_line" >> "$temp_file"

    # Use visudo to validate and install
    if sudo visudo -c -f "$temp_file" >/dev/null 2>&1; then
        sudo cp "$temp_file" /etc/sudoers
        rm -f "$temp_file"
        print_status "INFO" "✓ Passwordless sudo enabled for: pacman, yay, systemctl"

        # Force reload sudo rules for current session
        sudo -k

        return 0
    else
        rm -f "$temp_file"
        print_status "ERROR" "Failed to add sudo rule - syntax validation failed"
        return 1
    fi
}

# Remove our passwordless rule
remove_passwordless_rule() {
    print_status "INFO" "Removing temporary passwordless sudo rule..."

    if has_archriot_rule; then
        # Use sed to remove only our marked line
        if sudo sed -i "/$OMARCHY_SUDO_MARKER/d" /etc/sudoers; then
            print_status "INFO" "✓ Temporary sudo rule removed"
            return 0
        else
            print_status "ERROR" "Failed to remove sudo rule"
            return 1
        fi
    else
        print_status "INFO" "No temporary rule found to remove"
        return 0
    fi
}

# Setup passwordless sudo (main function)
setup_passwordless_sudo() {
    print_status "INFO" "Setting up passwordless sudo for package management..."

    # Check if already working
    if is_passwordless; then
        print_status "INFO" "Passwordless sudo already working - no changes needed"
        return 0
    fi

    # Ensure user is in wheel group (for system consistency, but not required)
    if ! is_in_wheel; then
        add_to_wheel
    fi

    # Add passwordless rule if not already present
    if ! has_archriot_rule; then
        add_passwordless_rule || return 1
    fi

    # Force reload sudo rules and test
    sudo -k
    sleep 1

    # Test if it works
    if is_passwordless; then
        print_status "INFO" "✓ Passwordless sudo setup complete"

        # Validate it's actually working
        if validate_passwordless_sudo; then
            print_status "INFO" "✓ Passwordless sudo validation passed"
            return 0
        else
            print_status "WARN" "Setup completed but validation failed"
            return 1
        fi
    else
        print_status "ERROR" "Passwordless sudo setup failed - you may need to re-login"
        print_status "INFO" "Manual test: sudo -n pacman --version"
        return 1
    fi
}

# Cleanup passwordless sudo (main function)
cleanup_passwordless_sudo() {
    print_status "INFO" "Cleaning up temporary passwordless sudo..."
    remove_passwordless_rule
}

# Test passwordless sudo functionality
test_passwordless_sudo() {
    print_status "INFO" "Testing passwordless sudo..."

    local test_commands=("pacman --version" "yay --version" "systemctl --version")
    local success=true

    for cmd in "${test_commands[@]}"; do
        if sudo -n $cmd >/dev/null 2>&1; then
            print_status "INFO" "✓ Working: $cmd"
        else
            print_status "WARN" "✗ Failed: $cmd"
            success=false
        fi
    done

    if $success; then
        print_status "INFO" "✓ All package management commands working"
        return 0
    else
        return 1
    fi
}

# Show current sudo status
show_sudo_status() {
    print_status "INFO" "Current sudo configuration:"
    print_status "INFO" "  • User: $CURRENT_USER"
    print_status "INFO" "  • In wheel group: $(is_in_wheel && echo "YES" || echo "NO")"
    print_status "INFO" "  • Passwordless for packages: $(is_passwordless && echo "YES" || echo "NO")"
    print_status "INFO" "  • OhmArchy rule active: $(has_archriot_rule && echo "YES" || echo "NO")"
}

# Safe sudo wrapper - ensures passwordless operation
safe_sudo() {
    local cmd="$*"

    # Test if passwordless sudo works for this command
    if sudo -n true 2>/dev/null; then
        # Passwordless sudo is working, execute the command
        sudo -n "$@"
    else
        print_status "ERROR" "Passwordless sudo not working for: $cmd"
        print_status "ERROR" "Installation cannot continue without passwordless sudo"
        print_status "INFO" "Please run: sudo -v (enter password once)"
        print_status "INFO" "Or setup passwordless sudo: setup_passwordless_sudo"
        return 1
    fi
}

# Test and validate passwordless sudo is working
validate_passwordless_sudo() {
    print_status "INFO" "Validating passwordless sudo configuration..."

    local test_commands=("true" "pacman --version" "systemctl --version")
    local all_working=true

    for cmd in "${test_commands[@]}"; do
        if sudo -n $cmd >/dev/null 2>&1; then
            print_status "INFO" "✓ Passwordless sudo working for: $cmd"
        else
            print_status "ERROR" "✗ Passwordless sudo failed for: $cmd"
            all_working=false
        fi
    done

    if [[ "$all_working" == "true" ]]; then
        print_status "INFO" "✅ Passwordless sudo validation passed"
        return 0
    else
        print_status "ERROR" "❌ Passwordless sudo validation failed"
        print_status "INFO" "Some commands will prompt for password during installation"
        return 1
    fi
}

# Emergency cleanup (removes any OhmArchy rules)
emergency_cleanup() {
    print_status "WARN" "Emergency cleanup - removing all OhmArchy sudo rules..."
    remove_passwordless_rule
}
