#!/bin/bash

# ==============================================================================
# OhmArchy Installation Validation Script
# ==============================================================================
# Comprehensive validation to ensure OhmArchy will install successfully
# and deliver the expected CypherRiot Wayland + Hyprland experience
# ==============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Validation counters
TESTS_PASSED=0
TESTS_FAILED=0
WARNINGS=0

# Print colored output
print_status() {
    local status="$1"
    local message="$2"
    case "$status" in
        "PASS") echo -e "${GREEN}✓${NC} $message"; ((TESTS_PASSED++)) ;;
        "FAIL") echo -e "${RED}❌${NC} $message"; ((TESTS_FAILED++)) ;;
        "WARN") echo -e "${YELLOW}⚠${NC} $message"; ((WARNINGS++)) ;;
        "INFO") echo -e "${BLUE}ℹ${NC} $message" ;;
        "TEST") echo -e "${PURPLE}🧪${NC} $message" ;;
    esac
}

# Header
print_header() {
    # Read version
    local version="unknown"
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [[ -f "$script_dir/VERSION" ]]; then
        version=$(cat "$script_dir/VERSION" 2>/dev/null || echo "unknown")
    elif [[ -f "$HOME/.local/share/omarchy/VERSION" ]]; then
        version=$(cat "$HOME/.local/share/omarchy/VERSION" 2>/dev/null || echo "unknown")
    fi

    echo -e "${PURPLE}"
    echo "============================================================"
    echo "           OhmArchy Installation Validation"
    echo "                     Version: $version"
    echo "============================================================"
    echo -e "${NC}"
    echo "This script validates that OhmArchy will install successfully"
    echo "and deliver the expected CypherRiot Wayland + Hyprland experience."
    echo ""
}

# Test basic system requirements
test_system_requirements() {
    print_status "TEST" "Testing system requirements..."

    # Check if running on Arch
    if [[ -f /etc/arch-release ]]; then
        print_status "PASS" "Running on Arch Linux"
    else
        print_status "FAIL" "Not running on Arch Linux - OhmArchy requires Arch"
        return 1
    fi

    # Check internet connectivity (using IPv4 and HTTPS)
    if ping -c 1 -4 google.com &>/dev/null || curl -s --max-time 5 https://google.com >/dev/null 2>&1; then
        print_status "PASS" "Internet connectivity available"
    else
        print_status "FAIL" "No internet connectivity - required for package downloads"
        return 1
    fi

    # Check available disk space (need at least 5GB)
    local available_gb=$(df / | awk 'NR==2 {printf "%.0f", $4/1024/1024}')
    if [[ $available_gb -ge 5 ]]; then
        print_status "PASS" "Sufficient disk space available: ${available_gb}GB"
    else
        print_status "FAIL" "Insufficient disk space: ${available_gb}GB (need 5GB+)"
        return 1
    fi

    # Check if running in TTY (optimal for first install)
    if [[ $(tty) == /dev/tty* ]]; then
        print_status "PASS" "Running in TTY (optimal for installation)"
    else
        print_status "WARN" "Not running in TTY - installation will work but TTY is preferred"
    fi
}

# Test package manager availability
test_package_manager() {
    print_status "TEST" "Testing package manager..."

    # Check pacman
    if command -v pacman &>/dev/null; then
        print_status "PASS" "Pacman available"
    else
        print_status "FAIL" "Pacman not available"
        return 1
    fi

    # Check yay or ability to install it
    if command -v yay &>/dev/null; then
        print_status "PASS" "Yay AUR helper already available"
    elif command -v git &>/dev/null && command -v base-devel &>/dev/null; then
        print_status "PASS" "Git and base-devel available for yay installation"
    else
        print_status "WARN" "Yay not installed, git/base-devel may need to be installed first"
    fi

    # Test sudo access
    if sudo -n true 2>/dev/null; then
        print_status "PASS" "Sudo access available (passwordless)"
    elif groups | grep -q wheel || id -nG | grep -q wheel; then
        print_status "PASS" "User in wheel group - sudo access available (will prompt for password)"
    elif sudo -v 2>/dev/null; then
        print_status "PASS" "Sudo access confirmed (will prompt for password during install)"
    else
        print_status "FAIL" "No sudo access - required for package installation"
        return 1
    fi
}

# Test critical installation files
test_installation_files() {
    print_status "TEST" "Testing installation files..."

    local install_dir="$HOME/.local/share/omarchy"

    # Test if we can clone/access the repo
    if [[ -d "$install_dir" ]]; then
        print_status "INFO" "OhmArchy already cloned"
    else
        print_status "INFO" "Testing repository clone..."
        if git clone --depth 1 https://github.com/CyphrRiot/OhmArchy.git /tmp/omarchy-test 2>/dev/null; then
            print_status "PASS" "Repository accessible and cloneable"
            rm -rf /tmp/omarchy-test
        else
            print_status "FAIL" "Cannot clone OhmArchy repository"
            return 1
        fi
    fi

    # Check critical install scripts exist (using current directory if available)
    local check_dir="."
    if [[ ! -f "./install.sh" ]]; then
        check_dir="$install_dir"
    fi

    local critical_files=(
        "install.sh"
        "install/core/01-base.sh"
        "install/desktop/hyprland.sh"
        "install/desktop/theming.sh"
        "config/ghostty/config"
        "config/fish/config.fish"
        "themes/cypherriot/config"
        "themes/cypherriot/backgrounds.sh"
        "config/waybar/style.css"
        "validate.sh"
    )

    for file in "${critical_files[@]}"; do
        if [[ -f "$check_dir/$file" ]]; then
            print_status "PASS" "Critical file exists: $file"
        else
            print_status "FAIL" "Missing critical file: $file"
            return 1
        fi
    done
}

# Test CypherRiot theme integrity
test_cypherriot_theme() {
    print_status "TEST" "Testing CypherRiot theme integrity..."

    local check_dir="."
    if [[ ! -d "./themes/cypherriot" ]]; then
        check_dir="$HOME/.local/share/omarchy"
    fi

    local theme_dir="$check_dir/themes/cypherriot"

    if [[ ! -d "$theme_dir" ]]; then
        print_status "FAIL" "CypherRiot theme directory not found"
        return 1
    fi

    # Check essential theme files
    local theme_files=(
        "config"
        "hyprland.conf"
        "hyprlock.conf"
        "fuzzel.ini"
        "neovim.lua"
        "ghostty.conf"
        "backgrounds/escape_velocity.jpg"
        "backgrounds.sh"
    )

    for file in "${theme_files[@]}"; do
        if [[ -f "$theme_dir/$file" ]]; then
            print_status "PASS" "Theme file exists: $file"
        else
            print_status "FAIL" "Missing theme file: $file"
            return 1
        fi
    done

    # Check if waybar config has theme specifics
    if grep -q "cypherriot\|CypherRiot\|Hack Nerd Font" "$theme_dir/config" 2>/dev/null; then
        print_status "PASS" "CypherRiot waybar config contains theme-specific settings"
    else
        print_status "WARN" "CypherRiot waybar config may be generic"
    fi
}

# Test Wayland/Hyprland compatibility
test_wayland_compatibility() {
    print_status "TEST" "Testing Wayland/Hyprland compatibility..."

    # Check if we're already in Wayland (would indicate compatibility)
    if [[ -n "$WAYLAND_DISPLAY" ]]; then
        print_status "PASS" "Already running in Wayland session"
    else
        print_status "INFO" "Not currently in Wayland (expected for TTY install)"
    fi

    # Check graphics driver compatibility
    if lspci | grep -i "vga\|3d\|display" | grep -qi "nvidia"; then
        print_status "WARN" "NVIDIA GPU detected - may need additional driver setup"
    elif lspci | grep -i "vga\|3d\|display" | grep -qi "amd\|radeon"; then
        print_status "PASS" "AMD GPU detected - good Wayland compatibility"
    elif lspci | grep -i "vga\|3d\|display" | grep -qi "intel"; then
        print_status "PASS" "Intel GPU detected - excellent Wayland compatibility"
    else
        print_status "WARN" "Unknown GPU type - may need manual driver configuration"
    fi

    # Check kernel version (newer kernels have better Wayland support)
    local kernel_version=$(uname -r | cut -d. -f1-2)
    local major=$(echo $kernel_version | cut -d. -f1)
    local minor=$(echo $kernel_version | cut -d. -f2)

    if [[ $major -gt 5 ]] || [[ $major -eq 5 && $minor -ge 15 ]]; then
        print_status "PASS" "Kernel version $kernel_version supports modern Wayland features"
    else
        print_status "WARN" "Kernel version $kernel_version may have limited Wayland support"
    fi
}

# Test if user setup will work
test_user_environment() {
    print_status "TEST" "Testing user environment..."

    # Check shell compatibility
    if [[ "$SHELL" == */fish ]]; then
        print_status "PASS" "Fish shell detected - optimal for OhmArchy"
    elif [[ "$SHELL" == */bash ]]; then
        print_status "PASS" "Bash shell detected - compatible with OhmArchy"
    else
        print_status "WARN" "Unusual shell detected: $SHELL - may need manual configuration"
    fi

    # Check home directory structure
    if [[ -w "$HOME" ]]; then
        print_status "PASS" "Home directory is writable"
    else
        print_status "FAIL" "Home directory is not writable"
        return 1
    fi

    # Check if .config exists or can be created
    if [[ -d "$HOME/.config" ]] || mkdir -p "$HOME/.config" 2>/dev/null; then
        print_status "PASS" "Config directory accessible"
    else
        print_status "FAIL" "Cannot create config directory"
        return 1
    fi

    # Test if we can create the omarchy directories
    if mkdir -p "$HOME/.config/omarchy/test" 2>/dev/null; then
        rmdir "$HOME/.config/omarchy/test" 2>/dev/null
        print_status "PASS" "Can create OhmArchy config directories"
    else
        print_status "FAIL" "Cannot create OhmArchy config directories"
        return 1
    fi
}

# Test package availability
test_package_availability() {
    print_status "TEST" "Testing critical package availability..."

    local critical_packages=(
        "hyprland"
        "waybar"
        "ghostty"
        "fish"
        "git"
        "curl"
        "gum"
    )

    for package in "${critical_packages[@]}"; do
        if pacman -Ss "^$package$" &>/dev/null; then
            print_status "PASS" "Package available: $package"
        else
            print_status "FAIL" "Package not available: $package"
            return 1
        fi
    done

    # Test AUR packages (if yay is available)
    if command -v yay &>/dev/null; then
        local aur_packages=("bibata-cursor-theme" "kora-icon-theme" "ghostty-shell-integration")
        for package in "${aur_packages[@]}"; do
            if yay -Ss "^$package$" &>/dev/null; then
                print_status "PASS" "AUR package available: $package"
            else
                print_status "WARN" "AUR package not found: $package"
            fi
        done
    fi
}

# Simulate critical installation steps
test_installation_simulation() {
    print_status "TEST" "Simulating critical installation steps..."

    # Test using local repo if available, otherwise clone from GitHub
    local test_dir="/tmp/omarchy-validation-$$"
    local using_local=false

    # Check if we're running from within the OhmArchy repo
    if [[ -f "./config/ghostty/config" && -f "./config/fish/config.fish" && -f "./install.sh" ]]; then
        print_status "PASS" "Using local OhmArchy repository for validation"
        test_dir="."
        using_local=true
    elif git clone --depth 1 https://github.com/CyphrRiot/OhmArchy.git "$test_dir" 2>/dev/null; then
        print_status "PASS" "Repository clone simulation successful"
    else
        print_status "FAIL" "Repository clone simulation failed"
        return 1
    fi

    if [[ "$using_local" == "true" || -d "$test_dir" ]]; then

        # Test if install.sh exists and is executable
        if [[ -x "$test_dir/install.sh" ]]; then
            print_status "PASS" "install.sh is executable"
        else
            print_status "WARN" "install.sh may not be executable"
        fi

        # Test theme files
        if [[ -f "$test_dir/themes/cypherriot/config" ]]; then
            print_status "PASS" "CypherRiot theme files present in clone"
        else
            print_status "FAIL" "CypherRiot theme files missing in clone"
        fi

        # Test Ghostty config
        if [[ -f "$test_dir/config/ghostty/config" ]]; then
            print_status "PASS" "Ghostty config present in clone"
        else
            print_status "FAIL" "Ghostty config missing in clone"
        fi

        # Test Fish config with fastfetch
        if [[ -f "$test_dir/config/fish/config.fish" ]] && grep -q "fastfetch" "$test_dir/config/fish/config.fish"; then
            print_status "PASS" "Fish config with fastfetch greeting present"
        else
            print_status "FAIL" "Fish config with fastfetch greeting missing"
        fi

        if [[ "$using_local" == "false" ]]; then
            rm -rf "$test_dir"
        fi
    fi

    # Test directory creation
    local test_config="/tmp/omarchy-config-test-$$"
    if mkdir -p "$test_config"/{themes,current,backgrounds} 2>/dev/null; then
        print_status "PASS" "Config directory structure creation works"
        rm -rf "$test_config"
    else
        print_status "FAIL" "Cannot create config directory structure"
        return 1
    fi
}

# Test expected post-install state
test_expected_outcome() {
    print_status "TEST" "Testing expected post-installation outcome..."

    print_status "INFO" "After successful installation, you should have:"
    print_status "INFO" "  • Hyprland Wayland compositor"
    print_status "INFO" "  • CypherRiot theme with purple/blue aesthetics"
    print_status "INFO" "  • Waybar with custom modules and Hack Nerd Font consistency"
    print_status "INFO" "  • Ghostty terminal with Fish shell and fastfetch greeting"
    print_status "INFO" "  • Floating terminal support (SUPER+SHIFT+RETURN)"
    print_status "INFO" "  • Multiple wallpapers including escape_velocity.jpg default"
    print_status "INFO" "  • Auto-login to Hyprland from TTY1"
    print_status "INFO" "  • Standardized font usage (Hack Nerd Font) across all components"

    # Check if system is already configured
    if [[ -f "$HOME/.config/hypr/hyprland.conf" ]]; then
        print_status "INFO" " Hyprland config already exists - this may be a re-install"
    fi

    if [[ -L "$HOME/.config/omarchy/current/theme" ]]; then
        local current_theme=$(basename "$(readlink "$HOME/.config/omarchy/current/theme")")
        print_status "INFO" " Current theme already set: $current_theme"
    fi
}

# Generate validation report
generate_report() {
    echo ""
    echo -e "${PURPLE}============================================================${NC}"
    echo -e "${PURPLE}                    VALIDATION REPORT${NC}"
    echo -e "${PURPLE}============================================================${NC}"
    echo ""

    echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
    echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
    echo ""

    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "${GREEN}🎉 VALIDATION SUCCESSFUL!${NC}"
        echo ""
        echo -e "${GREEN}✓ OhmArchy should install successfully${NC}"
        echo -e "${GREEN}✓ You should get a beautiful Wayland + Hyprland experience${NC}"
        echo -e "${GREEN}✓ CypherRiot theme should work properly${NC}"
        echo ""
        echo -e "${BLUE}Ready to install? Run:${NC}"
        echo -e "${BLUE}curl -fsSL https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/setup.sh | bash${NC}"
        echo ""

        if [[ $WARNINGS -gt 0 ]]; then
            echo -e "${YELLOW}Note: $WARNINGS warnings detected - installation should work but may need minor adjustments${NC}"
        fi

        return 0
    else
        echo -e "${RED}❌ VALIDATION FAILED!${NC}"
        echo ""
        echo -e "${RED}Issues found that may prevent successful installation:${NC}"
        echo -e "${RED}• $TESTS_FAILED critical tests failed${NC}"
        echo -e "${RED}• Please resolve the failed tests before installing${NC}"
        echo ""
        return 1
    fi
}

# Main execution
main() {
    print_header

    # Run all validation tests
    test_system_requirements || true
    test_package_manager || true
    test_installation_files || true
    test_cypherriot_theme || true
    test_wayland_compatibility || true
    test_user_environment || true
    test_package_availability || true
    test_installation_simulation || true
    test_expected_outcome || true

    # Generate final report
    generate_report
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
