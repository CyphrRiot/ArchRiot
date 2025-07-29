#!/bin/bash

# ================================================================================
# ArchRiot Installation Verification Script
# ================================================================================
# Comprehensive verification of ArchRiot installation
# Tests all critical components and reports status
# ================================================================================

# Configuration
readonly SCRIPT_NAME="ArchRiot Verification"
readonly LOG_FILE="$HOME/.cache/archriot/verify.log"
readonly CONFIG_DIR="$HOME/.config"
readonly ARCHRIOT_DIR="$HOME/.local/share/archriot"

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
WARNING_TESTS=0
declare -a FAILED_ITEMS=()
declare -a WARNING_ITEMS=()

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly WHITE='\033[1;37m'
readonly GRAY='\033[0;37m'
readonly NC='\033[0m' # No Color

# ================================================================================
# Utility Functions
# ================================================================================

print_header() {
    echo -e "${PURPLE}${WHITE}=================================="
    echo -e "🔍 $SCRIPT_NAME"
    echo -e "==================================${NC}"
    echo ""
}

print_section() {
    local section="$1"
    echo -e "${CYAN}📋 Testing: $section${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

test_result() {
    local test_name="$1"
    local status="$2"
    local details="$3"

    ((TOTAL_TESTS++))

    case "$status" in
        "PASS")
            echo -e "  ${GREEN}✅ $test_name${NC}"
            ((PASSED_TESTS++))
            ;;
        "FAIL")
            echo -e "  ${RED}❌ $test_name${NC}"
            if [[ -n "$details" ]]; then
                echo -e "     ${GRAY}$details${NC}"
            fi
            FAILED_ITEMS+=("$test_name")
            ((FAILED_TESTS++))
            ;;
        "WARN")
            echo -e "  ${YELLOW}⚠️  $test_name${NC}"
            if [[ -n "$details" ]]; then
                echo -e "     ${GRAY}$details${NC}"
            fi
            WARNING_ITEMS+=("$test_name")
            ((WARNING_TESTS++))
            ;;
    esac
}

# ================================================================================
# Test Functions
# ================================================================================

test_essential_packages() {
    print_section "Essential Packages"

    # Test yay
    if command -v yay &>/dev/null; then
        test_result "yay AUR helper" "PASS"
    else
        test_result "yay AUR helper" "FAIL" "yay not found in PATH"
    fi

    # Test git
    if command -v git &>/dev/null; then
        test_result "git version control" "PASS"
    else
        test_result "git version control" "FAIL" "git not found"
    fi

    # Test essential base-devel packages (individually installed)
    local essential_devel_packages=("gcc" "make" "autoconf" "automake" "binutils" "fakeroot" "file" "findutils" "gawk" "gettext" "grep" "gzip" "libtool" "m4" "pacman" "patch" "pkgconf" "sed" "sudo" "texinfo" "which")
    local missing_packages=0
    local total_packages=${#essential_devel_packages[@]}

    for package in "${essential_devel_packages[@]}"; do
        if ! pacman -Q "$package" &>/dev/null; then
            ((missing_packages++))
        fi
    done

    if [[ $missing_packages -eq 0 ]]; then
        test_result "base-devel packages ($total_packages/$(pacman -Qg base-devel 2>/dev/null | wc -l || echo $total_packages))" "PASS"
    elif [[ $missing_packages -le 2 ]]; then
        test_result "base-devel packages ($((total_packages - missing_packages))/$total_packages)" "WARN" "$missing_packages essential packages missing"
    else
        test_result "base-devel packages ($((total_packages - missing_packages))/$total_packages)" "FAIL" "$missing_packages essential packages missing"
    fi

    echo ""
}

test_desktop_environment() {
    print_section "Desktop Environment"

    # Test Hyprland
    if command -v hyprland &>/dev/null; then
        test_result "Hyprland compositor" "PASS"
    else
        test_result "Hyprland compositor" "FAIL" "Hyprland not found"
    fi

    # Test Waybar
    if command -v waybar &>/dev/null; then
        test_result "Waybar status bar" "PASS"
    else
        test_result "Waybar status bar" "FAIL" "Waybar not found"
    fi

    # Test gum (UI component)
    if command -v gum &>/dev/null; then
        test_result "gum UI component" "PASS"
    else
        test_result "gum UI component" "FAIL" "gum not found"
    fi

    # Test fuzzel (app launcher)
    if command -v fuzzel &>/dev/null; then
        test_result "Fuzzel app launcher" "PASS"
    else
        test_result "Fuzzel app launcher" "WARN" "fuzzel not found"
    fi

    # Test notification daemon
    if command -v dunst &>/dev/null; then
        test_result "Dunst notifications" "PASS"
    elif command -v mako &>/dev/null; then
        test_result "Mako notifications" "PASS"
    else
        test_result "Notification daemon" "WARN" "No notification daemon found"
    fi

    echo ""
}

test_theming_system() {
    print_section "Theming System (CRITICAL)"

    # Test theme directory structure
    if [[ -d "$CONFIG_DIR/archriot/current" ]]; then
        test_result "Theme directory structure" "PASS"

        # Check for theme symlinks
        local theme_links=0
        for link in "$CONFIG_DIR/archriot/current"/*; do
            if [[ -L "$link" ]]; then
                ((theme_links++))
            fi
        done

        if [[ $theme_links -gt 0 ]]; then
            test_result "Theme symlinks ($theme_links found)" "PASS"
        else
            test_result "Theme symlinks" "FAIL" "No theme symlinks found"
        fi
    else
        test_result "Theme directory structure" "FAIL" "~/.config/archriot/current not found"
    fi

    # Test cursor theme
    local cursor_theme_found=false
    for cursor_dir in "/usr/share/icons/Bibata-Modern-Ice" "$HOME/.local/share/icons/Bibata-Modern-Ice" "$HOME/.icons/Bibata-Modern-Ice"; do
        if [[ -d "$cursor_dir" ]]; then
            cursor_theme_found=true
            break
        fi
    done

    if [[ "$cursor_theme_found" == "true" ]]; then
        test_result "Bibata cursor theme" "PASS"
    else
        test_result "Bibata cursor theme" "FAIL" "Bibata-Modern-Ice not found"
    fi

    # Test icon theme
    local icon_theme_found=false
    for icon_dir in "/usr/share/icons/Tela-purple" "/usr/share/icons/Tela-purple-dark" "/usr/share/icons/kora"; do
        if [[ -d "$icon_dir" ]]; then
            icon_theme_found=true
            break
        fi
    done

    if [[ "$icon_theme_found" == "true" ]]; then
        test_result "Icon themes" "PASS"
    else
        test_result "Icon themes" "WARN" "No expected icon themes found"
    fi

    # Test GTK theme configuration
    if [[ -f "$CONFIG_DIR/gtk-3.0/settings.ini" ]]; then
        test_result "GTK-3 theme config" "PASS"
    else
        test_result "GTK-3 theme config" "WARN" "GTK-3 settings.ini not found"
    fi

    echo ""
}

test_configuration_files() {
    print_section "Configuration Files"

    # Test Hyprland config
    if [[ -f "$CONFIG_DIR/hypr/hyprland.conf" ]]; then
        test_result "Hyprland configuration" "PASS"
    else
        test_result "Hyprland configuration" "FAIL" "hyprland.conf not found"
    fi

    # Test Waybar config
    if [[ -f "$CONFIG_DIR/waybar/config" ]]; then
        test_result "Waybar configuration" "PASS"
    else
        test_result "Waybar configuration" "FAIL" "waybar config not found"
    fi

    # Test shell configuration
    local shell_config_found=false
    for shell_config in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.config/fish/config.fish"; do
        if [[ -f "$shell_config" ]]; then
            shell_config_found=true
            break
        fi
    done

    if [[ "$shell_config_found" == "true" ]]; then
        test_result "Shell configuration" "PASS"
    else
        test_result "Shell configuration" "WARN" "No shell config found"
    fi

    # Test ArchRiot user environment
    if [[ -f "$CONFIG_DIR/archriot/user.env" ]]; then
        test_result "ArchRiot user environment" "PASS"
    else
        test_result "ArchRiot user environment" "WARN" "user.env not found"
    fi

    echo ""
}

test_applications() {
    print_section "Applications"

    # Test terminal
    local terminal_found=false
    for terminal in "kitty" "alacritty" "wezterm" "foot"; do
        if command -v "$terminal" &>/dev/null; then
            test_result "Terminal ($terminal)" "PASS"
            terminal_found=true
            break
        fi
    done

    if [[ "$terminal_found" == "false" ]]; then
        test_result "Terminal" "FAIL" "No supported terminal found"
    fi

    # Test file manager
    if command -v thunar &>/dev/null; then
        test_result "File manager (Thunar)" "PASS"
    elif command -v nautilus &>/dev/null; then
        test_result "File manager (Nautilus)" "PASS"
    else
        test_result "File manager" "WARN" "No GUI file manager found"
    fi

    # Test text editor
    if command -v nvim &>/dev/null; then
        test_result "Text editor (Neovim)" "PASS"
    elif command -v vim &>/dev/null; then
        test_result "Text editor (Vim)" "PASS"
    else
        test_result "Text editor" "WARN" "No vim/nvim found"
    fi

    # Test browser
    if command -v firefox &>/dev/null; then
        test_result "Web browser (Firefox)" "PASS"
    elif command -v chromium &>/dev/null; then
        test_result "Web browser (Chromium)" "PASS"
    else
        test_result "Web browser" "WARN" "No browser found"
    fi

    echo ""
}

test_system_services() {
    print_section "System Services"

    # Test audio system
    if command -v pulseaudio &>/dev/null || command -v pipewire &>/dev/null; then
        test_result "Audio system" "PASS"
    else
        test_result "Audio system" "WARN" "No audio system found"
    fi

    # Test network management
    if systemctl is-enabled NetworkManager &>/dev/null; then
        test_result "NetworkManager" "PASS"
    elif systemctl is-enabled systemd-networkd &>/dev/null; then
        test_result "systemd-networkd" "PASS"
    else
        test_result "Network management" "WARN" "No network manager enabled"
    fi

    # Test bluetooth
    if systemctl is-enabled bluetooth &>/dev/null; then
        test_result "Bluetooth service" "PASS"
    else
        test_result "Bluetooth service" "WARN" "Bluetooth not enabled"
    fi

    echo ""
}

test_fonts() {
    print_section "Fonts"

    # Test font cache
    if fc-list | grep -q .; then
        test_result "Font cache" "PASS"
    else
        test_result "Font cache" "FAIL" "Font cache empty"
    fi

    # Test for specific font families
    local fonts_found=0
    for font in "Noto Sans" "JetBrains Mono" "Font Awesome"; do
        if fc-list | grep -qi "$font"; then
            ((fonts_found++))
        fi
    done

    if [[ $fonts_found -gt 0 ]]; then
        test_result "Essential fonts ($fonts_found/3 found)" "PASS"
    else
        test_result "Essential fonts" "WARN" "No essential fonts found"
    fi

    echo ""
}

test_installation_integrity() {
    print_section "Installation Integrity"

    # Test ArchRiot installation directory
    if [[ -d "$ARCHRIOT_DIR" ]]; then
        test_result "ArchRiot installation directory" "PASS"
    else
        test_result "ArchRiot installation directory" "FAIL" "~/.local/share/archriot not found"
    fi

    # Test version file
    if [[ -f "$ARCHRIOT_DIR/VERSION" ]]; then
        local version=$(cat "$ARCHRIOT_DIR/VERSION" 2>/dev/null)
        test_result "ArchRiot version ($version)" "PASS"
    else
        test_result "ArchRiot version" "WARN" "VERSION file not found"
    fi

    # Test installation logs
    if [[ -f "$HOME/.cache/archriot/install.log" ]]; then
        test_result "Installation logs" "PASS"
    else
        test_result "Installation logs" "WARN" "No installation logs found"
    fi

    echo ""
}

# ================================================================================
# Special Tests for Known Issues
# ================================================================================

test_known_issues() {
    print_section "Known Issue Checks"

    # Test the specific theming issue that was problematic
    if [[ -f "$ARCHRIOT_DIR/install/desktop/theming.sh" ]]; then
        test_result "Theming script exists" "PASS"

        # Check if theming script can be executed
        if bash -n "$ARCHRIOT_DIR/install/desktop/theming.sh" 2>/dev/null; then
            test_result "Theming script syntax" "PASS"
        else
            test_result "Theming script syntax" "FAIL" "Syntax errors in theming.sh"
        fi
    else
        test_result "Theming script exists" "FAIL" "theming.sh not found"
    fi

    # Test for lock screen functionality (Super+L issue)
    if command -v hyprlock &>/dev/null; then
        test_result "Lock screen (hyprlock)" "PASS"
    elif command -v swaylock &>/dev/null; then
        test_result "Lock screen (swaylock)" "PASS"
    else
        test_result "Lock screen" "FAIL" "No lock screen program found"
    fi

    # Test for mysterious directory creation issue
    if [[ -d "$CONFIG_DIR/archriot" ]]; then
        local dir_permissions=$(stat -c "%a" "$CONFIG_DIR/archriot" 2>/dev/null)
        if [[ "$dir_permissions" == "755" ]] || [[ "$dir_permissions" == "700" ]]; then
            test_result "ArchRiot config directory permissions" "PASS"
        else
            test_result "ArchRiot config directory permissions" "WARN" "Unusual permissions: $dir_permissions"
        fi
    fi

    echo ""
}

# ================================================================================
# Summary and Recommendations
# ================================================================================

print_summary() {
    echo -e "${PURPLE}${WHITE}=================================="
    echo -e "📊 Verification Summary"
    echo -e "==================================${NC}"
    echo ""

    echo -e "${WHITE}Total Tests: $TOTAL_TESTS${NC}"
    echo -e "${GREEN}✅ Passed: $PASSED_TESTS${NC}"
    echo -e "${YELLOW}⚠️  Warnings: $WARNING_TESTS${NC}"
    echo -e "${RED}❌ Failed: $FAILED_TESTS${NC}"
    echo ""

    # Calculate success rate
    local success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))

    if [[ $FAILED_TESTS -eq 0 ]]; then
        echo -e "${GREEN}🎉 All critical tests passed! ArchRiot installation is healthy.${NC}"
    elif [[ $FAILED_TESTS -le 2 ]]; then
        echo -e "${YELLOW}⚠️  Minor issues detected. System should be mostly functional.${NC}"
    else
        echo -e "${RED}🚨 Significant issues detected. ArchRiot may not function properly.${NC}"
    fi

    echo -e "${WHITE}Success Rate: $success_rate%${NC}"
    echo ""

    # Show failed items
    if [[ ${#FAILED_ITEMS[@]} -gt 0 ]]; then
        echo -e "${RED}Failed Tests:${NC}"
        for item in "${FAILED_ITEMS[@]}"; do
            echo -e "  ${RED}❌ $item${NC}"
        done
        echo ""
    fi

    # Show warnings
    if [[ ${#WARNING_ITEMS[@]} -gt 0 ]]; then
        echo -e "${YELLOW}Warnings:${NC}"
        for item in "${WARNING_ITEMS[@]}"; do
            echo -e "  ${YELLOW}⚠️  $item${NC}"
        done
        echo ""
    fi

    # Recommendations
    echo -e "${CYAN}📋 Recommendations:${NC}"

    if [[ $FAILED_TESTS -gt 0 ]]; then
        echo -e "  ${WHITE}• Re-run installation to fix failed components:${NC}"
        echo -e "    ${GRAY}source ~/.local/share/archriot/install.sh${NC}"
        echo ""
    fi

    if [[ ${#FAILED_ITEMS[@]} -gt 0 ]]; then
        for item in "${FAILED_ITEMS[@]}"; do
            case "$item" in
                *"Theme directory"*|*"theming"*)
                    echo -e "  ${WHITE}• Fix theming system:${NC}"
                    echo -e "    ${GRAY}bash ~/.local/share/archriot/install/desktop/theming.sh${NC}"
                    ;;
                *"yay"*)
                    echo -e "  ${WHITE}• Install yay AUR helper manually${NC}"
                    ;;
                *"Hyprland"*)
                    echo -e "  ${WHITE}• Install Hyprland: yay -S hyprland${NC}"
                    ;;
            esac
        done
        echo ""
    fi

    echo -e "${GRAY}Log saved to: $LOG_FILE${NC}"
    echo ""
}

# ================================================================================
# Main Execution
# ================================================================================

main() {
    # Initialize logging
    mkdir -p "$(dirname "$LOG_FILE")"

    # Start verification
    print_header

    # Run all tests
    test_essential_packages
    test_desktop_environment
    test_theming_system
    test_configuration_files
    test_applications
    test_system_services
    test_fonts
    test_installation_integrity
    test_known_issues

    # Show summary
    print_summary

    # Save log
    {
        echo "ArchRiot Verification Results - $(date)"
        echo "=========================================="
        echo "Total: $TOTAL_TESTS, Passed: $PASSED_TESTS, Warnings: $WARNING_TESTS, Failed: $FAILED_TESTS"
        echo ""
        echo "Failed items:"
        printf '%s\n' "${FAILED_ITEMS[@]}"
        echo ""
        echo "Warning items:"
        printf '%s\n' "${WARNING_ITEMS[@]}"
    } > "$LOG_FILE"

    # Exit with appropriate code
    if [[ $FAILED_TESTS -eq 0 ]]; then
        exit 0
    elif [[ $FAILED_TESTS -le 2 ]]; then
        exit 1
    else
        exit 2
    fi
}

# Execute main function if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
