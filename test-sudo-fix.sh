#!/bin/bash

# ==============================================================================
# OhmArchy Sudo Fix Test Script
# ==============================================================================
# Test the sudo passwordless setup to validate the fix
# ==============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    local status="$1"
    local message="$2"
    case "$status" in
        "INFO") echo -e "${BLUE}ℹ ${NC} $message" ;;
        "PASS") echo -e "${GREEN}✓${NC} $message" ;;
        "FAIL") echo -e "${RED}❌${NC} $message" ;;
        "WARN") echo -e "${YELLOW}⚠${NC} $message" ;;
    esac
}

echo "🔒 Testing OhmArchy Sudo Passwordless Setup"
echo "=========================================="

# Load the sudo helper
if [ -f "install/lib/sudo-helper.sh" ]; then
    source install/lib/sudo-helper.sh
    print_status "PASS" "Sudo helper loaded successfully"
else
    print_status "FAIL" "Could not find install/lib/sudo-helper.sh"
    exit 1
fi

echo ""
print_status "INFO" "Current sudo configuration:"
show_sudo_status

echo ""
print_status "INFO" "Testing passwordless sudo setup..."

# Test setup
if setup_passwordless_sudo; then
    print_status "PASS" "Setup completed successfully"
else
    print_status "FAIL" "Setup failed"
    exit 1
fi

echo ""
print_status "INFO" "Testing passwordless commands..."

# Test each command
test_commands=("pacman --version" "yay --version" "systemctl --version")
all_passed=true

for cmd in "${test_commands[@]}"; do
    if sudo -n $cmd >/dev/null 2>&1; then
        print_status "PASS" "Passwordless: $cmd"
    else
        print_status "FAIL" "Still requires password: $cmd"
        all_passed=false
    fi
done

echo ""
if $all_passed; then
    print_status "PASS" "All passwordless commands working!"
else
    print_status "FAIL" "Some commands still require passwords"
fi

echo ""
print_status "INFO" "Current sudoers rule:"
if has_omarchy_rule; then
    sudo grep "$OMARCHY_SUDO_MARKER" /etc/sudoers 2>/dev/null | while read -r line; do
        echo "  $line"
    done
else
    print_status "WARN" "No OhmArchy rule found in sudoers"
fi

echo ""
print_status "INFO" "Testing cleanup..."
if cleanup_passwordless_sudo; then
    print_status "PASS" "Cleanup completed successfully"
else
    print_status "FAIL" "Cleanup failed"
fi

echo ""
print_status "INFO" "Final verification - should require passwords again:"
for cmd in "${test_commands[@]}"; do
    if sudo -n $cmd >/dev/null 2>&1; then
        print_status "WARN" "Still passwordless after cleanup: $cmd"
    else
        print_status "PASS" "Correctly requires password: $cmd"
    fi
done

echo ""
print_status "INFO" "Test complete!"
