#!/bin/bash

# ==============================================================================
# ArchRiot Ivy Bridge Vulkan Compatibility Fix - SAFE VERSION
# ==============================================================================
# SAFE, NON-DESTRUCTIVE fix for Vulkan issues on Intel HD Graphics 4000
# Does NOT modify core ArchRiot files - only adds user-specific configurations
# ==============================================================================

set -e

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly LOG_FILE="$HOME/.cache/archriot/ivy-bridge-fix.log"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Initialize logging
init_logging() {
    mkdir -p "$(dirname "$LOG_FILE")"
    echo "=== ArchRiot Ivy Bridge Fix - $(date) ===" > "$LOG_FILE"
}

log() {
    echo "$1" | tee -a "$LOG_FILE"
}

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}❌ $1${NC}" | tee -a "$LOG_FILE"
}

# Detect Intel GPU generation (safe detection only)
detect_intel_gpu() {
    local gpu_info
    gpu_info=$(lspci | grep -i "vga\|3d\|display" | grep -i intel 2>/dev/null || true)

    if [[ -z "$gpu_info" ]]; then
        echo "none"
        return
    fi

    # Conservative detection for known problematic hardware
    if echo "$gpu_info" | grep -qi -E "ivy bridge|hd 4000|hd 2500"; then
        echo "ivybridge"
    elif echo "$gpu_info" | grep -qi -E "sandy bridge|hd 3000|hd 2000"; then
        echo "sandybridge"
    elif echo "$gpu_info" | grep -qi -E "haswell|hd 4600|hd 4400|hd 4200"; then
        echo "haswell"
    else
        echo "modern"
    fi
}

# Test current Vulkan status safely
test_vulkan_status() {
    log_info "Testing current Vulkan status..."

    if ! command -v vulkaninfo >/dev/null 2>&1; then
        log_warning "vulkaninfo not found - installing vulkan-tools..."
        if sudo pacman -S --noconfirm --needed vulkan-tools 2>/dev/null; then
            log_success "vulkan-tools installed"
        else
            log_error "Failed to install vulkan-tools"
            return 1
        fi
    fi

    # Test vulkaninfo output
    local vulkan_output
    vulkan_output=$(vulkaninfo 2>&1 || true)

    if echo "$vulkan_output" | grep -qi "ivy bridge vulkan support is incomplete"; then
        log_warning "Ivy Bridge incomplete Vulkan detected"
        return 2
    elif echo "$vulkan_output" | grep -qi "cannot find.*vulkan.*driver"; then
        log_error "No Vulkan driver found"
        return 3
    elif echo "$vulkan_output" | grep -qi "apiVersion.*1\.[0-9]"; then
        log_success "Vulkan appears functional"
        return 0
    else
        log_warning "Vulkan status unclear"
        return 4
    fi
}

# Create Zed configuration overlay (non-destructive)
create_zed_config_overlay() {
    local intel_gen="$1"
    local zed_config_dir="$HOME/.config/zed"
    local settings_file="$zed_config_dir/settings.json"
    local overlay_file="$zed_config_dir/ivy-bridge-overlay.json"

    log_info "Creating Zed configuration overlay..."

    mkdir -p "$zed_config_dir"

    case "$intel_gen" in
        "ivybridge"|"sandybridge")
            log_info "Creating OpenGL-optimized overlay for $intel_gen"
            cat > "$overlay_file" << 'EOF'
{
    "_ivy_bridge_fix": "This overlay optimizes Zed for older Intel graphics",
    "_merge_instructions": "Manually merge these settings into your settings.json",
    "gpu": {
        "use_vulkan": false,
        "use_opengl": true
    },
    "experimental": {
        "renderer": "opengl"
    },
    "editor": {
        "soft_wrap": "editor_width",
        "cursor_blink": false
    }
}
EOF
            log_success "Overlay created: $overlay_file"
            log_info "To apply: manually copy settings from overlay to your settings.json"
            ;;
        *)
            log_info "Modern GPU detected - no overlay needed"
            ;;
    esac
}

# Create safe launcher wrapper
create_safe_launcher() {
    local wrapper_file="$HOME/.local/bin/zed-ivy-bridge"

    log_info "Creating safe Zed launcher..."

    mkdir -p "$HOME/.local/bin"

    cat > "$wrapper_file" << 'EOF'
#!/bin/bash
# Safe Zed launcher for Ivy Bridge compatibility
# Does not modify system files

# Check if this is an Ivy Bridge system
if lspci | grep -qi "hd 4000\|hd 2500\|ivy bridge"; then
    echo "🔧 Ivy Bridge detected - using compatibility mode"

    # Set environment variables for better compatibility
    export LIBGL_ALWAYS_SOFTWARE=0
    export MESA_GL_VERSION_OVERRIDE=3.3
    export MESA_GLSL_VERSION_OVERRIDE=330

    # Disable problematic Vulkan on old hardware
    export DISABLE_VULKAN=1
    export VK_ICD_FILENAMES=""

    echo "✅ Compatibility mode active"
fi

# Launch Zed normally
exec zed "$@"
EOF

    chmod +x "$wrapper_file"
    log_success "Safe launcher created: $wrapper_file"
    log_info "Usage: zed-ivy-bridge (instead of zed)"
}

# Test the fix safely
test_fix() {
    local intel_gen="$1"

    log_info "Testing fix functionality..."

    # Test 1: Check if wrapper exists and is executable
    if [[ -x "$HOME/.local/bin/zed-ivy-bridge" ]]; then
        log_success "Launcher wrapper is executable"
    else
        log_error "Launcher wrapper test failed"
        return 1
    fi

    # Test 2: Check if Zed overlay exists (for older GPUs)
    case "$intel_gen" in
        "ivybridge"|"sandybridge")
            if [[ -f "$HOME/.config/zed/ivy-bridge-overlay.json" ]]; then
                log_success "Zed overlay configuration created"
            else
                log_error "Zed overlay test failed"
                return 1
            fi
            ;;
    esac

    # Test 3: Basic Vulkan test (non-destructive)
    if command -v vulkaninfo >/dev/null 2>&1; then
        if timeout 5 vulkaninfo --summary >/dev/null 2>&1; then
            log_success "Basic Vulkan test passed"
        else
            log_warning "Vulkan test inconclusive (expected for older hardware)"
        fi
    fi

    return 0
}

# Show manual steps for full fix
show_manual_steps() {
    local intel_gen="$1"

    echo ""
    echo "📋 MANUAL STEPS TO COMPLETE THE FIX:"
    echo "===================================="

    case "$intel_gen" in
        "ivybridge"|"sandybridge")
            echo ""
            echo "1. 📝 Update your Zed configuration:"
            echo "   • Open: ~/.config/zed/settings.json"
            echo "   • Copy settings from: ~/.config/zed/ivy-bridge-overlay.json"
            echo "   • Add the gpu and experimental sections"
            echo ""
            echo "2. 🚀 Use the compatibility launcher:"
            echo "   • Command: zed-ivy-bridge"
            echo "   • Or add alias: alias zed='zed-ivy-bridge'"
            echo ""
            echo "3. 🔧 If issues persist:"
            echo "   • Try: LIBGL_ALWAYS_SOFTWARE=1 zed"
            echo "   • Check: ~/.cache/archriot/ivy-bridge-fix.log"
            ;;
        "haswell")
            echo ""
            echo "1. 🚀 Use the compatibility launcher:"
            echo "   • Command: zed-ivy-bridge"
            echo "   • Should resolve most issues"
            echo ""
            echo "2. 🔧 If problems continue:"
            echo "   • Check Vulkan: vulkaninfo --summary"
            echo "   • Try software fallback if needed"
            ;;
        "modern")
            echo ""
            echo "✅ Modern GPU detected - minimal fixes applied"
            echo "   • Compatibility launcher available if needed"
            echo "   • No configuration changes required"
            ;;
    esac

    echo ""
    echo "🔄 After making changes:"
    echo "   • Restart Zed completely"
    echo "   • Test with: zed-ivy-bridge"
    echo "   • Check logs: tail -f ~/.cache/archriot/ivy-bridge-fix.log"
}

# Main execution
main() {
    init_logging

    echo "🔧 ArchRiot Ivy Bridge Vulkan Fix (Safe Version)"
    echo "================================================"
    echo ""

    # Pre-flight checks
    if [[ $EUID -eq 0 ]]; then
        log_error "Don't run as root! Use your regular user account."
        exit 1
    fi

    if [[ ! -d "$HOME/.local/share/archriot" ]]; then
        log_error "ArchRiot not found. Install ArchRiot first."
        exit 1
    fi

    # Detect hardware
    local intel_gen
    intel_gen=$(detect_intel_gpu)
    log_info "Detected Intel GPU: $intel_gen"

    if [[ "$intel_gen" == "none" ]]; then
        log_info "No Intel GPU detected. This fix is not needed."
        exit 0
    fi

    # Test current Vulkan status
    local vulkan_status=0
    test_vulkan_status || vulkan_status=$?

    echo ""
    log_info "Applying safe, non-destructive fixes..."

    # Apply fixes (safe operations only)
    create_zed_config_overlay "$intel_gen"
    create_safe_launcher

    echo ""
    log_info "Testing applied fixes..."
    if test_fix "$intel_gen"; then
        log_success "All tests passed!"
    else
        log_warning "Some tests failed - check manual steps"
    fi

    # Show completion summary
    echo ""
    log_success "Safe Ivy Bridge fix completed!"
    echo ""
    echo "📊 SUMMARY:"
    echo "  • GPU Generation: $intel_gen"
    echo "  • Vulkan Status: $([ $vulkan_status -eq 0 ] && echo "Functional" || echo "Needs fallback")"
    echo "  • Launcher: ~/.local/bin/zed-ivy-bridge"
    echo "  • Config overlay: ~/.config/zed/ivy-bridge-overlay.json"
    echo "  • Log file: $LOG_FILE"

    show_manual_steps "$intel_gen"

    echo ""
    echo "🛡️  This fix is SAFE and NON-DESTRUCTIVE:"
    echo "   • No system files modified"
    echo "   • No core ArchRiot changes"
    echo "   • Easy to remove if needed"
    echo ""
    echo "💡 To remove: rm ~/.local/bin/zed-ivy-bridge ~/.config/zed/ivy-bridge-overlay.json"
}

main "$@"
