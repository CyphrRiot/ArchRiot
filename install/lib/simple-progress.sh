#!/bin/bash

# =============================================================================
# OhmArchy Clean Progress System
# Shows clean progress bars while capturing package manager output to files
# =============================================================================

# ANSI Color Codes (OhmArchy Theme)
declare -A COLORS=(
    ["RESET"]='\033[0m'
    ["BOLD"]='\033[1m'
    ["DIM"]='\033[2m'
    ["WHITE"]='\033[38;2;255;255;255m'
    ["GRAY"]='\033[38;2;104;104;104m'
    ["GREEN"]='\033[38;2;158;206;106m'
    ["BLUE"]='\033[38;2;125;166;255m'
    ["PURPLE"]='\033[38;2;157;123;216m'
    ["CYAN"]='\033[38;2;13;185;215m'
    ["ORANGE"]='\033[38;2;255;158;100m'
    ["YELLOW"]='\033[38;2;224;175;104m'
    ["RED"]='\033[38;2;255;122;147m'
)

# Progress configuration
PROGRESS_ENABLED=true
PROGRESS_WIDTH=40
ESTIMATED_TOTAL_MINUTES=3
START_TIME=""
CURRENT_PHASE=""
PHASE_COUNT=0
TOTAL_PHASES=0
LOG_DIR="/tmp/archriot-install"

# Disable progress in non-interactive terminals
if [[ ! -t 1 ]] || [[ -z "$TERM" ]] || [[ "$TERM" == "dumb" ]]; then
    PROGRESS_ENABLED=false
fi

# Create log directory
mkdir -p "$LOG_DIR"

# Initialize progress tracking
init_clean_progress() {
    local total_phases=${1:-8}
    TOTAL_PHASES=$total_phases
    PHASE_COUNT=0
    START_TIME=$(date +%s)

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo
        echo -e "${COLORS[PURPLE]}${COLORS[BOLD]}🚀 OhmArchy Installation${COLORS[RESET]}"
        echo -e "${COLORS[GRAY]}$(printf '━%.0s' {1..60})${COLORS[RESET]}"
        echo -e "${COLORS[CYAN]}Estimated time: ~$ESTIMATED_TOTAL_MINUTES minutes${COLORS[RESET]}"
        echo
    fi
}

# Start a new module
start_module() {
    local module_name="$1"
    local color="${2:-BLUE}"

    if [[ "$PROGRESS_ENABLED" != "true" ]]; then
        return 0
    fi

    PHASE_COUNT=$((PHASE_COUNT + 1))
    CURRENT_PHASE="$module_name"

    # Calculate progress based on modules completed
    local percent=$((PHASE_COUNT * 100 / TOTAL_PHASES))

    # Cap at 95% until completion
    if [[ $percent -gt 95 ]]; then
        percent=95
    fi

    # Build progress bar
    local filled=$((percent * PROGRESS_WIDTH / 100))
    local empty=$((PROGRESS_WIDTH - filled))
    local bar=""
    for ((i=0; i<filled; i++)); do bar+="█"; done
    for ((i=0; i<empty; i++)); do bar+="░"; done

    # Calculate time estimates
    local elapsed_seconds=$(($(date +%s) - START_TIME))
    local elapsed_minutes=$((elapsed_seconds / 60))
    local remaining_minutes=$((ESTIMATED_TOTAL_MINUTES - elapsed_minutes))
    if [[ $remaining_minutes -lt 0 ]]; then
        remaining_minutes=0
    fi

    # Show module progress
    echo -e "${COLORS[WHITE]}[${PHASE_COUNT}/${TOTAL_PHASES}] ${COLORS[$color]}${module_name}${COLORS[RESET]}"
    echo -e "${COLORS[PURPLE]}▐${bar}▌ ${COLORS[WHITE]}${percent}%${COLORS[RESET]} ${COLORS[GRAY]}(~${remaining_minutes}m remaining)${COLORS[RESET]}"
    echo
}

# Show progress for current phase (legacy compatibility)
show_phase_progress() {
    local phase_name="$1"
    local color="${2:-BLUE}"
    # Just show the phase name without incrementing counter
    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo -e "${COLORS[$color]}$phase_name${COLORS[RESET]}"
    fi
}

# Install packages with clean progress and error handling
install_packages_clean() {
    local packages="$1"
    local phase_name="$2"
    local color="${3:-BLUE}"
    local log_file="$LOG_DIR/$(date +%s)-$(echo "$phase_name" | tr ' ' '-' | tr '[:upper:]' '[:lower:]').log"

    echo -e "${COLORS[BLUE]}📦 Installing: $packages${COLORS[RESET]}"
    echo -e "${COLORS[GRAY]}   (Output logged to: $log_file)${COLORS[RESET]}"

    # Install with output captured - use yay if available, otherwise fallback to pacman
    local install_cmd
    if command -v yay &>/dev/null; then
        install_cmd="yay -S --noconfirm --needed $packages"
    else
        install_cmd="sudo pacman -S --noconfirm --needed $packages"
    fi

    if $install_cmd > "$log_file" 2>&1; then
        echo -e "${COLORS[PURPLE]}✓ Successfully installed${COLORS[RESET]}"

        # Show any warnings from the log
        if grep -q "warning:" "$log_file"; then
            local warnings=$(grep "warning:" "$log_file" | wc -l)
            echo -e "${COLORS[YELLOW]}  ⚠ $warnings warnings (see log for details)${COLORS[RESET]}"
        fi

        echo
        return 0
    else
        local exit_code=$?
        echo -e "${COLORS[RED]}❌ Installation failed (exit code: $exit_code)${COLORS[RESET]}"
        echo
        echo -e "${COLORS[RED]}Error details:${COLORS[RESET]}"
        echo -e "${COLORS[GRAY]}$(printf '─%.0s' {1..60})${COLORS[RESET]}"
        cat "$log_file"
        echo -e "${COLORS[GRAY]}$(printf '─%.0s' {1..60})${COLORS[RESET]}"
        echo
        echo -e "${COLORS[YELLOW]}Full log saved to: $log_file${COLORS[RESET]}"
        return $exit_code
    fi
}

# Run a command with clean progress
run_command_clean() {
    local command="$1"
    local phase_name="$2"
    local color="${3:-BLUE}"
    local log_file="$LOG_DIR/$(date +%s)-$(echo "$phase_name" | tr ' ' '-' | tr '[:upper:]' '[:lower:]').log"

    show_phase_progress "$phase_name" "$color"

    echo -e "${COLORS[BLUE]}⚙ Running: $phase_name${COLORS[RESET]}"
    echo -e "${COLORS[GRAY]}   (Output logged to: $log_file)${COLORS[RESET]}"

    # Run command with output captured
    if eval "$command" > "$log_file" 2>&1; then
        echo -e "${COLORS[PURPLE]}✓ Successfully completed${COLORS[RESET]}"
        echo
        return 0
    else
        local exit_code=$?
        echo -e "${COLORS[RED]}❌ Command failed (exit code: $exit_code)${COLORS[RESET]}"
        echo
        echo -e "${COLORS[RED]}Error details:${COLORS[RESET]}"
        echo -e "${COLORS[GRAY]}$(printf '─%.0s' {1..60})${COLORS[RESET]}"
        cat "$log_file"
        echo -e "${COLORS[GRAY]}$(printf '─%.0s' {1..60})${COLORS[RESET]}"
        echo
        echo -e "${COLORS[YELLOW]}Full log saved to: $log_file${COLORS[RESET]}"
        return $exit_code
    fi
}

# Show installation completion
complete_clean_installation() {
    if [[ "$PROGRESS_ENABLED" != "true" ]]; then
        return 0
    fi

    local total_seconds=$(($(date +%s) - START_TIME))
    local total_minutes=$((total_seconds / 60))
    local remaining_seconds=$((total_seconds % 60))

    # Show 100% completion
    local bar=""
    for ((i=0; i<PROGRESS_WIDTH; i++)); do bar+="█"; done

    echo -e "${COLORS[GRAY]}$(printf '━%.0s' {1..60})${COLORS[RESET]}"
    echo -e "${COLORS[PURPLE]}${COLORS[BOLD]}🎉 Installation Complete!${COLORS[RESET]}"
    echo -e "${COLORS[PURPLE]}▐${bar}▌ ${COLORS[WHITE]}100%${COLORS[RESET]}"

    if [[ $total_minutes -gt 0 ]]; then
        echo -e "${COLORS[GRAY]}Total time: ${total_minutes}m ${remaining_seconds}s${COLORS[RESET]}"
    else
        echo -e "${COLORS[GRAY]}Total time: ${total_seconds}s${COLORS[RESET]}"
    fi

    echo -e "${COLORS[GRAY]}Installation logs: $LOG_DIR${COLORS[RESET]}"
    echo
}

# Show installation failure
show_clean_failure() {
    local error_msg="${1:-Installation failed}"

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo
        echo -e "${COLORS[RED]}${COLORS[BOLD]}❌ $error_msg${COLORS[RESET]}"
        echo -e "${COLORS[GRAY]}Current phase: $CURRENT_PHASE${COLORS[RESET]}"

        local elapsed_seconds=$(($(date +%s) - START_TIME))
        local elapsed_minutes=$((elapsed_seconds / 60))
        echo -e "${COLORS[GRAY]}Time elapsed: ${elapsed_minutes}m${COLORS[RESET]}"
        echo -e "${COLORS[GRAY]}Installation logs: $LOG_DIR${COLORS[RESET]}"
        echo
        echo -e "${COLORS[YELLOW]}To retry installation:${COLORS[RESET]}"
        echo -e "${COLORS[WHITE]}  source ~/.local/share/omarchy/install.sh${COLORS[RESET]}"
    fi
}

# Clean up old log files (optional)
cleanup_old_logs() {
    find "$LOG_DIR" -name "*.log" -mtime +7 -delete 2>/dev/null || true
}

# Enhanced wrapper functions for backward compatibility
install_packages() {
    local packages="$1"
    local package_type="${2:-packages}"
    install_packages_clean "$packages" "Installing $package_type" "BLUE"
}

install_essential() {
    local packages="$1"
    install_packages_clean "$packages" "Installing essential packages" "BLUE"
}

install_optional() {
    local packages="$1"
    install_packages_clean "$packages" "Installing optional packages" "YELLOW"
}

# Export functions for use in installers
export -f init_clean_progress start_module show_phase_progress install_packages_clean
export -f run_command_clean complete_clean_installation show_clean_failure
export -f install_packages install_essential install_optional cleanup_old_logs
