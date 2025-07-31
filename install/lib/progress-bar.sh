#!/bin/bash

# ================================================================================
# Modern Progress Bar System for ArchRiot v2.1
# ================================================================================
# Beautiful static progress bar with friendly module names
# All verbose output goes to background logs, only errors surface to console
# ================================================================================

# Progress tracking globals
declare -g TOTAL_MODULES=0
declare -g CURRENT_MODULE=0
declare -g PROGRESS_BAR_WIDTH=50
declare -g PROGRESS_INITIALIZED=false

# Module friendly names mapping
declare -A FRIENDLY_NAMES=(
    # Core modules
    ["01-base"]="Base System Setup"
    ["02-identity"]="User Identity Configuration"
    ["04-shell"]="Shell Environment Setup"

    # Directory modules
    ["core"]="Core System Components"
    ["system"]="System Configuration"
    ["development"]="Development Tools"
    ["desktop"]="Desktop Environment"
    ["post-desktop"]="Desktop Finalization"
    ["applications"]="Essential Applications"
    ["optional"]="Optional Components"

    # Standalone installers
    ["asdcontrol"]="System Control Panel"
    ["mimetypes"]="File Type Associations"
    ["plymouth"]="Boot Screen Setup"
    ["printer"]="Printer Support"
)

# Colors for beautiful output
PROGRESS_RED='\033[0;31m'
PROGRESS_GREEN='\033[0;32m'
PROGRESS_YELLOW='\033[1;33m'
PROGRESS_BLUE='\033[0;34m'
PROGRESS_PURPLE='\033[0;35m'
PROGRESS_CYAN='\033[0;36m'
PROGRESS_WHITE='\033[1;37m'
PROGRESS_NC='\033[0m' # No Color

# Get friendly name for a module
get_friendly_name() {
    local module_path="$1"
    local module_name=$(basename "$module_path" .sh)
    local dir_name=$(basename "$(dirname "$module_path")")

    # Check for specific module name first
    if [[ -n "${FRIENDLY_NAMES[$module_name]}" ]]; then
        echo "${FRIENDLY_NAMES[$module_name]}"
    elif [[ -n "${FRIENDLY_NAMES[$dir_name]}" ]]; then
        echo "${FRIENDLY_NAMES[$dir_name]}"
    else
        # Fallback: prettify the name
        echo "$module_name" | sed 's/-/ /g' | sed 's/\b\w/\U&/g'
    fi
}

# Count total modules for accurate percentage (must match installer logic exactly)
count_total_modules() {
    # Simply count all install scripts
    find "$INSTALL_DIR" -name "*.sh" | wc -l
}

# Draw the static progress bar (stays in same position)
draw_progress_bar() {
    local percentage=$1
    local current_task="$2"
    local status="$3"  # SUCCESS, ERROR, or empty

    # Only draw if progress system is initialized
    [[ "$PROGRESS_INITIALIZED" == "true" ]] || return 0

    # Calculate filled portion
    local filled=$((percentage * PROGRESS_BAR_WIDTH / 100))
    local empty=$((PROGRESS_BAR_WIDTH - filled))

    # Move cursor to progress bar line (line 6)
    printf '\033[6;1H'

    # Clear the line and draw progress bar
    printf '\033[K'  # Clear line

    # Choose color based on status
    local bar_color="$PROGRESS_CYAN"
    local status_icon="⚡"

    case "$status" in
        "SUCCESS")
            bar_color="$PROGRESS_GREEN"
            status_icon="✅"
            ;;
        "ERROR")
            bar_color="$PROGRESS_RED"
            status_icon="❌"
            ;;
    esac

    # Draw the progress bar with simple ASCII characters
    printf "${PROGRESS_WHITE}Progress: ${bar_color}["
    for ((i=1; i<=filled; i++)); do
        printf "="
    done
    for ((i=1; i<=empty; i++)); do
        printf "-"
    done
    printf "]${PROGRESS_NC} ${PROGRESS_WHITE}%3d%%${PROGRESS_NC}\n" $percentage

    # Show current task (line 7)
    printf '\033[7;1H'
    printf '\033[K'  # Clear line
    printf "${status_icon} ${PROGRESS_WHITE}%s${PROGRESS_NC}\n" "$current_task"

    # Show module progress (line 8)
    printf '\033[8;1H'
    printf '\033[K'  # Clear line
    printf "${PROGRESS_PURPLE}Module ${CURRENT_MODULE} of ${TOTAL_MODULES}${PROGRESS_NC}\n"
}

# Initialize progress bar system
init_progress_system() {
    # Count total modules
    TOTAL_MODULES=$(count_total_modules)
    CURRENT_MODULE=0
    PROGRESS_INITIALIZED=true

    # Clear screen and show header
    clear
    echo ""
    echo "🚀 ArchRiot Installation v$ARCHRIOT_VERSION"
    echo "══════════════════════════════════════════════════════════════════════"
    echo ""
    echo "Installing to: $HOME/.local/share/archriot"
    echo ""
    echo ""  # Space for progress bar (line 6)
    echo ""  # Space for current task (line 7)
    echo ""  # Space for module counter (line 8)
    echo ""
    echo ""  # Space for errors (line 10+)

    # Show initial progress bar
    draw_progress_bar 0 "Initializing installation system..."
}

# Update progress for a module start
progress_module_start() {
    local module_path="$1"
    ((CURRENT_MODULE++))

    local friendly_name=$(get_friendly_name "$module_path")
    local percentage=$((CURRENT_MODULE * 100 / TOTAL_MODULES))

    draw_progress_bar $percentage "$friendly_name"
}

# Update progress for a module completion
progress_module_complete() {
    local module_path="$1"
    local success="$2"  # true/false

    local friendly_name=$(get_friendly_name "$module_path")
    local percentage=$((CURRENT_MODULE * 100 / TOTAL_MODULES))

    if [[ "$success" == "true" ]]; then
        draw_progress_bar $percentage "$friendly_name" "SUCCESS"
    else
        draw_progress_bar $percentage "$friendly_name" "ERROR"
    fi

    sleep 0.3  # Brief pause to show status
}

# Show error on console (surfaces from background)
progress_show_error() {
    local error_message="$1"
    local module_name="$2"

    # Show error below progress area (line 10+)
    printf '\033[10;1H'
    printf "${PROGRESS_RED}⚠️  Error in %s: %s${PROGRESS_NC}\n" "$module_name" "$error_message"
}

# Clear error messages
progress_clear_errors() {
    # Clear error area (lines 10-15)
    for i in {10..15}; do
        printf '\033[%d;1H\033[K' $i
    done
}

# Show beautiful completion summary
progress_show_completion() {
    local end_time=$(date +%s)

    # Get start time from main installer or fallback to current time
    local start_time="${INSTALL_START_TIME:-$end_time}"

    # Debug logging
    echo "DEBUG: end_time=$end_time, start_time=$start_time, INSTALL_START_TIME=$INSTALL_START_TIME" >> "${LOG_FILE:-/dev/null}" 2>&1

    local duration=$((end_time - start_time))
    local duration_min=$((duration / 60))
    local duration_sec=$((duration % 60))

    # Ensure we have a reasonable duration (not negative or impossibly large)
    if [[ $duration -lt 0 ]] || [[ $duration -gt 7200 ]]; then
        echo "DEBUG: Invalid duration $duration, resetting to 0" >> "${LOG_FILE:-/dev/null}" 2>&1
        duration_min=0
        duration_sec=0
    fi

    local success_count=$((TOTAL_MODULES - ${#FAILED_MODULES[@]}))
    local success_rate=$((success_count * 100 / TOTAL_MODULES))

    # Move below progress area
    printf '\033[12;1H'

    echo ""
    printf "${PROGRESS_GREEN}══════════════════════════════════════════════════════════════════════${PROGRESS_NC}\n"
    printf "${PROGRESS_GREEN}🎉 ArchRiot Installation Complete!${PROGRESS_NC}\n"
    printf "${PROGRESS_GREEN}══════════════════════════════════════════════════════════════════════${PROGRESS_NC}\n"
    echo ""
    printf "⏱️  ${PROGRESS_WHITE}Total time:${PROGRESS_NC} ${duration_min}m ${duration_sec}s\n"
    printf "📦 ${PROGRESS_WHITE}Modules processed:${PROGRESS_NC} ${success_count}/${TOTAL_MODULES}\n"
    printf "✅ ${PROGRESS_WHITE}Success rate:${PROGRESS_NC} ${success_rate}%%\n"
    echo ""
    printf "📁 ${PROGRESS_WHITE}Installation log:${PROGRESS_NC} $LOG_FILE\n"
    printf "🏠 ${PROGRESS_WHITE}ArchRiot location:${PROGRESS_NC} ~/.local/share/archriot\n"

    if [[ ${#FAILED_MODULES[@]} -gt 0 ]]; then
        echo ""
        printf "${PROGRESS_YELLOW}⚠️  Failed modules (${#FAILED_MODULES[@]}):${PROGRESS_NC}\n"
        for module in "${FAILED_MODULES[@]}"; do
            printf "   ${PROGRESS_RED}❌ $module${PROGRESS_NC}\n"
        done
        echo ""
        printf "📝 ${PROGRESS_WHITE}Check error log:${PROGRESS_NC} $ERROR_LOG\n"
    fi

    if [[ ${#SKIPPED_MODULES[@]} -gt 0 ]]; then
        echo ""
        printf "${PROGRESS_BLUE}ℹ️  Skipped modules (${#SKIPPED_MODULES[@]}):${PROGRESS_NC}\n"
        for module in "${SKIPPED_MODULES[@]}"; do
            printf "   ${PROGRESS_BLUE}⏭️  $module${PROGRESS_NC}\n"
        done
    fi

    echo ""
    if [[ ${#FAILED_MODULES[@]} -eq 0 ]]; then
        printf "${PROGRESS_CYAN}🚀 Perfect! You have updated to the latest ArchRiot.${PROGRESS_NC}\n"
    else
        printf "${PROGRESS_YELLOW}🔄 Some modules failed. Check logs and retry if needed.${PROGRESS_NC}\n"
    fi
    echo ""
}

# Temporarily pause progress display for user interaction
progress_pause_for_input() {
    local message="$1"

    # Move below progress area for user interaction
    printf '\033[12;1H'
    echo ""
    echo "$message"
    echo ""
}

# Resume progress display after user interaction
progress_resume_after_input() {
    # Clear the interaction area and redraw progress
    clear
    echo ""
    echo "🚀 ArchRiot Installation v$ARCHRIOT_VERSION"
    echo "══════════════════════════════════════════════════════════════════════"
    echo ""
    echo "Installing to: $HOME/.local/share/archriot"
    echo ""
    echo ""  # Space for progress bar (line 6)
    echo ""  # Space for current task (line 7)
    echo ""  # Space for module counter (line 8)
    echo ""
    echo ""  # Space for errors (line 10+)

    # Redraw current progress
    local percentage=$((CURRENT_MODULE * 100 / TOTAL_MODULES))
    draw_progress_bar $percentage "Resuming installation..."
}

# Disable progress system (for fallback)
disable_progress_system() {
    PROGRESS_INITIALIZED=false
}
