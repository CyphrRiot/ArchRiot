#!/bin/bash

# =============================================================================
# ArchRiot Clean Progress System
# Shows clean progress bars while capturing package manager output to files
# =============================================================================

# ANSI Color Codes (ArchRiot Theme)
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
# Global log file location
export ARCHRIOT_LOG_FILE="$HOME/.cache/archriot/install.log"

# Disable progress in non-interactive terminals
if [[ ! -t 1 ]] || [[ -z "$TERM" ]] || [[ "$TERM" == "dumb" ]]; then
    PROGRESS_ENABLED=false
fi

# Initialize single log file
init_error_log() {
    mkdir -p "$(dirname "$ARCHRIOT_LOG_FILE")"
    echo "=== ArchRiot Installation Log - $(date) ===" > "$ARCHRIOT_LOG_FILE"
}

# Initialize progress tracking
init_clean_progress() {
    local total_phases=${1:-8}
    TOTAL_PHASES=$total_phases
    PHASE_COUNT=0
    START_TIME=$(date +%s)

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo
        echo -e "${COLORS[PURPLE]}${COLORS[BOLD]}🚀 ArchRiot Installation${COLORS[RESET]}"
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
        # Module name already shown in progress bar - no need to repeat
        true
    fi
}

# Install packages with clean progress and error handling
install_packages_clean() {
    local packages="$1"
    local phase_name="$2"
    local color="${3:-BLUE}"

    echo -e "${COLORS[BLUE]}📦 Installing: $packages${COLORS[RESET]}"

    # Install with output captured - use yay if available, otherwise fallback to pacman
    local install_cmd
    if command -v yay &>/dev/null; then
        install_cmd="yay -S --noconfirm --needed $packages"
    else
        install_cmd="sudo pacman -S --noconfirm --needed $packages"
    fi

    # Capture output to check for warnings/errors
    local temp_output=$(mktemp)
    if $install_cmd > "$temp_output" 2>&1; then
        # Check for warnings and log them
        if grep -q "warning:" "$temp_output"; then
            local warnings=$(grep "warning:" "$temp_output" | wc -l)
            echo -e "${COLORS[YELLOW]}  ⚠ $warnings warnings${COLORS[RESET]}"
            echo "[$(date)] WARNINGS in $phase_name:" >> "$ARCHRIOT_LOG_FILE"
            grep "warning:" "$temp_output" >> "$ARCHRIOT_LOG_FILE"
            echo "" >> "$ARCHRIOT_LOG_FILE"
        fi

        echo -e "${COLORS[PURPLE]}✓ Successfully installed${COLORS[RESET]}"
        rm -f "$temp_output"
        echo
        return 0
    else
        local exit_code=$?
        echo -e "${COLORS[RED]}❌ Installation failed (exit code: $exit_code)${COLORS[RESET]}"
        echo -e "${COLORS[YELLOW]}Error details logged to: $ARCHRIOT_LOG_FILE${COLORS[RESET]}"

        # Log error details
        echo "[$(date)] ERROR in $phase_name ($packages):" >> "$ARCHRIOT_LOG_FILE"
        cat "$temp_output" >> "$ARCHRIOT_LOG_FILE"
        echo "" >> "$ARCHRIOT_LOG_FILE"

        rm -f "$temp_output"
        echo
        return $exit_code
    fi
}

# Run a command with clean progress
run_command_clean() {
    local command="$1"
    local phase_name="$2"
    local color="${3:-BLUE}"

    show_phase_progress "$phase_name" "$color"

    # Module name already shown in progress bar - no redundant display needed

    # Run command with output captured
    local temp_output=$(mktemp)
    if eval "$command" > "$temp_output" 2>&1; then
        echo -e "${COLORS[PURPLE]}✓ Successfully completed${COLORS[RESET]}"
        rm -f "$temp_output"
        echo
        return 0
    else
        local exit_code=$?
        echo -e "${COLORS[RED]}❌ Command failed (exit code: $exit_code)${COLORS[RESET]}"
        echo -e "${COLORS[YELLOW]}Error details logged to: $ARCHRIOT_LOG_FILE${COLORS[RESET]}"

        # Log error details
        echo "[$(date)] ERROR in $phase_name:" >> "$ARCHRIOT_LOG_FILE"
        cat "$temp_output" >> "$ARCHRIOT_LOG_FILE"
        echo "" >> "$ARCHRIOT_LOG_FILE"

        rm -f "$temp_output"
        echo
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

    echo -e "${COLORS[GRAY]}Installation log: $ARCHRIOT_LOG_FILE${COLORS[RESET]}"
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
        echo -e "${COLORS[GRAY]}Installation log: $ARCHRIOT_LOG_FILE${COLORS[RESET]}"
        echo
        echo -e "${COLORS[YELLOW]}To retry installation:${COLORS[RESET]}"
        echo -e "${COLORS[WHITE]}  source ~/.local/share/archriot/install.sh${COLORS[RESET]}"
    fi
}

# Clean up old log files (optional)
cleanup_old_logs() {
    # Single log file - no cleanup needed
    # User can manually delete $ARCHRIOT_LOG_FILE if desired
    return 0
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

# Progress pause/resume functions for interactive input
PROGRESS_PAUSED=false
SAVED_PROGRESS_STATE=""

pause_progress() {
    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        PROGRESS_PAUSED=true
        SAVED_PROGRESS_STATE="$PROGRESS_ENABLED"
        PROGRESS_ENABLED=false
        echo # Add blank line for clean input
    fi
}

resume_progress() {
    if [[ "$PROGRESS_PAUSED" == "true" ]]; then
        PROGRESS_ENABLED="$SAVED_PROGRESS_STATE"
        PROGRESS_PAUSED=false
        echo # Add blank line before resuming
    fi
}

# Export functions for use in installers
export -f init_error_log init_clean_progress start_module show_phase_progress install_packages_clean
export -f run_command_clean complete_clean_installation show_clean_failure
export -f install_packages install_essential install_optional cleanup_old_logs
