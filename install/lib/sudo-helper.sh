#!/bin/bash

# ==============================================================================
# ArchRiot Sudo Helper - Simple and Safe
# ==============================================================================
# Ensures user is in wheel group and enables passwordless sudo
# Uses proper wheel group configuration instead of individual user rules
# ==============================================================================

CURRENT_USER="$(whoami)"
WHEEL_SUDO_RULE="%wheel ALL=(ALL) NOPASSWD: ALL"
ARCHRIOT_SUDO_MARKER="# ArchRiot Auto-Generated"
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

# Check if passwordless sudo works (general test)
is_passwordless() {
    sudo -n whoami >/dev/null 2>&1
}

# Check if wheel NOPASSWD rule is already active (any variant)
has_wheel_rule() {
    # Check for uncommented wheel NOPASSWD rule (exact or similar)
    sudo grep -q "^%wheel.*ALL.*NOPASSWD.*ALL" /etc/sudoers 2>/dev/null || \
    sudo grep -q "^%wheel.*ALL=(ALL).*NOPASSWD:.*ALL" /etc/sudoers 2>/dev/null
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
add_wheel_rule() {
    print_status "INFO" "Enabling passwordless sudo for wheel group..."

    # Double-check if rule already exists (any variant)
    if has_wheel_rule; then
        print_status "INFO" "✓ Wheel NOPASSWD rule already exists"
        return 0
    fi

    # Check if wheel rule already exists but is commented
    if sudo grep -q "^# %wheel ALL=(ALL) NOPASSWD: ALL" /etc/sudoers; then
        # Uncomment the existing rule
        if sudo sed -i 's/^# %wheel ALL=(ALL) NOPASSWD: ALL/%wheel ALL=(ALL) NOPASSWD: ALL/' /etc/sudoers; then
            print_status "INFO" "✓ Uncommented existing wheel NOPASSWD rule"
            return 0
        else
            print_status "ERROR" "Failed to uncomment wheel rule"
            return 1
        fi
    else
        # Add new wheel rule
        local temp_file=$(mktemp)
        sudo cp /etc/sudoers "$temp_file"

        # Check if we're about to create a duplicate
        if grep -q "^%wheel.*NOPASSWD" "$temp_file"; then
            rm -f "$temp_file"
            print_status "INFO" "✓ Wheel NOPASSWD rule already exists (variant found)"
            return 0
        fi

        # Add wheel rule after the wheel group section
        if grep -q "^%wheel ALL=(ALL) ALL" "$temp_file"; then
            # Add after existing wheel rule
            sed -i "/^%wheel ALL=(ALL) ALL/a\\$WHEEL_SUDO_RULE $ARCHRIOT_SUDO_MARKER" "$temp_file"
        else
            # Add at end of file
            echo "$WHEEL_SUDO_RULE $ARCHRIOT_SUDO_MARKER" >> "$temp_file"
        fi

        # Validate and apply
        if sudo visudo -c -f "$temp_file" >/dev/null 2>&1; then
            sudo cp "$temp_file" /etc/sudoers
            rm -f "$temp_file"
            print_status "INFO" "✓ Wheel group passwordless sudo enabled"
            return 0
        else
            rm -f "$temp_file"
            print_status "ERROR" "Failed to add wheel rule (syntax error)"
            return 1
        fi
    fi
}

# Remove ArchRiot sudo rule and old individual user files
remove_passwordless_rule() {
    print_status "INFO" "Removing ArchRiot sudo rule..."

    # Remove rule from main sudoers file
    if sudo grep -q "$ARCHRIOT_SUDO_MARKER" /etc/sudoers; then
        # Use sed to remove only our marked line
        if sudo sed -i "/$ARCHRIOT_SUDO_MARKER/d" /etc/sudoers; then
            print_status "INFO" "✓ ArchRiot sudo rule removed"
        else
            print_status "ERROR" "Failed to remove sudo rule"
            return 1
        fi
    else
        print_status "INFO" "No ArchRiot rule found to remove"
    fi

    # Remove old individual user file if it exists
    local user_sudo_file="/etc/sudoers.d/00-$CURRENT_USER"
    if [[ -f "$user_sudo_file" ]]; then
        if sudo rm -f "$user_sudo_file"; then
            print_status "INFO" "✓ Removed old individual user sudo file: 00-$CURRENT_USER"
        else
            print_status "WARN" "Failed to remove old user sudo file: 00-$CURRENT_USER"
        fi
    fi

    # Also remove any OhmArchy legacy files
    if sudo grep -q "$OMARCHY_SUDO_MARKER" /etc/sudoers 2>/dev/null; then
        if sudo sed -i "/$OMARCHY_SUDO_MARKER/d" /etc/sudoers; then
            print_status "INFO" "✓ Removed legacy OhmArchy sudo rule"
        fi
    fi

    return 0
}

# Show current sudo configuration state
show_current_sudo_state() {
    print_status "INFO" "Current sudo configuration:"
    print_status "INFO" "  • User: $CURRENT_USER"
    print_status "INFO" "  • In wheel group: $(is_in_wheel && echo "YES" || echo "NO")"
    print_status "INFO" "  • Passwordless sudo working: $(is_passwordless && echo "YES" || echo "NO")"
    print_status "INFO" "  • Wheel NOPASSWD rule exists: $(has_wheel_rule && echo "YES" || echo "NO")"

    # Check for individual user files
    local user_sudo_file="/etc/sudoers.d/00-$CURRENT_USER"
    if [[ -f "$user_sudo_file" ]]; then
        print_status "INFO" "  • Individual user file exists: YES (will be cleaned up)"
    else
        print_status "INFO" "  • Individual user file exists: NO"
    fi

    # Check for legacy rules
    if sudo grep -q "$OMARCHY_SUDO_MARKER\|$ARCHRIOT_SUDO_MARKER" /etc/sudoers 2>/dev/null; then
        print_status "INFO" "  • Legacy sudo rules found: YES (will be cleaned up)"
    else
        print_status "INFO" "  • Legacy sudo rules found: NO"
    fi
}

# Setup passwordless sudo (main function)
setup_passwordless_sudo() {
    print_status "INFO" "Setting up passwordless sudo using wheel group..."

    # Show current state first
    show_current_sudo_state

    # Check if already working
    if is_passwordless; then
        print_status "INFO" "Passwordless sudo already working - no changes needed"
        return 0
    fi

    # Ensure user is in wheel group (REQUIRED for wheel-based sudo)
    if ! is_in_wheel; then
        add_to_wheel || return 1
        print_status "INFO" "User added to wheel group - you may need to re-login for full effect"
    fi

    # Clean up any old individual user files first
    local user_sudo_file="/etc/sudoers.d/00-$CURRENT_USER"
    if [[ -f "$user_sudo_file" ]]; then
        print_status "INFO" "Cleaning up old individual user sudo file..."
        if sudo rm -f "$user_sudo_file"; then
            print_status "INFO" "✓ Removed old user sudo file: 00-$CURRENT_USER"
        else
            print_status "WARN" "Failed to remove old user sudo file"
        fi
    fi

    # Add wheel rule if not already present
    if ! has_wheel_rule; then
        add_wheel_rule || return 1
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

    if sudo -n whoami >/dev/null 2>&1; then
        print_status "INFO" "✓ Passwordless sudo working"
        return 0
    else
        print_status "WARN" "✗ Passwordless sudo failed"
        return 1
    fi
}

# Show current sudo status
show_sudo_status() {
    print_status "INFO" "Current sudo configuration:"
    print_status "INFO" "  • User: $CURRENT_USER"
    print_status "INFO" "  • In wheel group: $(is_in_wheel && echo "YES" || echo "NO")"
    print_status "INFO" "  • Passwordless sudo: $(is_passwordless && echo "YES" || echo "NO")"
    print_status "INFO" "  • Wheel NOPASSWD rule: $(has_wheel_rule && echo "YES" || echo "NO")"
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
