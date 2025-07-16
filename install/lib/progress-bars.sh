#!/bin/bash

# OhmArchy Progress Bar Library
# Safe integration into existing installer without breaking functionality

# Global progress bar configuration
PROGRESS_ENABLED=${PROGRESS_ENABLED:-true}
PROGRESS_WIDTH=${PROGRESS_WIDTH:-35}
PROGRESS_LABEL_WIDTH=${PROGRESS_LABEL_WIDTH:-25}

# OhmArchy Theme Colors
declare -A PROGRESS_COLORS=(
    ["RESET"]='\033[0m'
    ["BOLD"]='\033[1m'
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

# Progress tracking variables
CURRENT_TASK=""
CURRENT_PROGRESS=0
TOTAL_TASKS=0
COMPLETED_TASKS=0

# Initialize progress system
init_progress() {
    local total=${1:-1}
    TOTAL_TASKS=$total
    COMPLETED_TASKS=0

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo -e "${PROGRESS_COLORS[PURPLE]}${PROGRESS_COLORS[BOLD]}Installation Progress:${PROGRESS_COLORS[RESET]}"
        echo ""
    fi
}

# Show a progress bar for current task
show_progress() {
    local percent=$1
    local label="$2"
    local color="${3:-BLUE}"

    # Skip if progress is disabled
    if [[ "$PROGRESS_ENABLED" != "true" ]]; then
        return 0
    fi

    # Validate inputs
    if [[ ! "$percent" =~ ^[0-9]+$ ]] || [[ $percent -lt 0 ]] || [[ $percent -gt 100 ]]; then
        return 0
    fi

    # Calculate filled portion
    local filled=$((percent * PROGRESS_WIDTH / 100))
    local empty=$((PROGRESS_WIDTH - filled))

    # Build progress bar
    local bar=""
    for ((i=0; i<filled; i++)); do bar+="█"; done
    for ((i=0; i<empty; i++)); do bar+="░"; done

    # Print with perfect alignment
    printf "\r${PROGRESS_COLORS[WHITE]}%-${PROGRESS_LABEL_WIDTH}s ${PROGRESS_COLORS[GRAY]}▐${PROGRESS_COLORS[$color]}%s${PROGRESS_COLORS[GRAY]}▌ ${PROGRESS_COLORS[GREEN]}%3d%%${PROGRESS_COLORS[RESET]}" \
        "$label" "$bar" "$percent"
}

# Start a new task with progress
start_task() {
    local task_name="$1"
    local color="${2:-BLUE}"

    CURRENT_TASK="$task_name"
    CURRENT_PROGRESS=0

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        show_progress 0 "$task_name" "$color"
    fi
}

# Update current task progress
update_progress() {
    local percent=$1
    local color="${2:-BLUE}"

    CURRENT_PROGRESS=$percent

    if [[ "$PROGRESS_ENABLED" == "true" && -n "$CURRENT_TASK" ]]; then
        show_progress "$percent" "$CURRENT_TASK" "$color"
    fi
}

# Complete current task
complete_task() {
    local color="${1:-GREEN}"

    if [[ "$PROGRESS_ENABLED" == "true" && -n "$CURRENT_TASK" ]]; then
        show_progress 100 "$CURRENT_TASK" "$color"
        echo ""  # New line after completion
        COMPLETED_TASKS=$((COMPLETED_TASKS + 1))
    fi

    CURRENT_TASK=""
    CURRENT_PROGRESS=0
}

# Fail current task
fail_task() {
    local error_msg="${1:-Task failed}"

    if [[ "$PROGRESS_ENABLED" == "true" && -n "$CURRENT_TASK" ]]; then
        show_progress $CURRENT_PROGRESS "$CURRENT_TASK" "RED"
        echo ""  # New line after failure
        echo -e "${PROGRESS_COLORS[RED]}❌ $error_msg${PROGRESS_COLORS[RESET]}"
    fi

    CURRENT_TASK=""
    CURRENT_PROGRESS=0
}

# Show overall installation progress
show_overall_progress() {
    if [[ "$PROGRESS_ENABLED" == "true" && $TOTAL_TASKS -gt 0 ]]; then
        local overall_percent=$((COMPLETED_TASKS * 100 / TOTAL_TASKS))
        echo ""
        show_progress $overall_percent "Overall Progress" "PURPLE"
        echo ""
    fi
}

# Wrapper for existing installer functions to add progress
run_with_progress() {
    local task_name="$1"
    local command="$2"
    local color="${3:-BLUE}"

    # Start the task
    start_task "$task_name" "$color"

    # If progress is disabled, run normally with original output
    if [[ "$PROGRESS_ENABLED" != "true" ]]; then
        eval "$command"
        return $?
    fi

    # Run command and capture output
    local temp_file=$(mktemp)
    local exit_code=0

    # Simulate progress during command execution
    {
        eval "$command" > "$temp_file" 2>&1
        exit_code=$?
    } &
    local cmd_pid=$!

    # Animate progress while command runs
    local progress=0
    while kill -0 $cmd_pid 2>/dev/null; do
        update_progress $progress "$color"
        progress=$(( (progress + 5) % 95 ))  # Never reach 100% until done
        sleep 0.2
    done

    # Wait for command to complete
    wait $cmd_pid
    exit_code=$?

    # Show final result
    if [[ $exit_code -eq 0 ]]; then
        complete_task "GREEN"
    else
        fail_task "Command failed"
        # Show error output if progress is enabled
        echo ""
        echo -e "${PROGRESS_COLORS[RED]}Error output:${PROGRESS_COLORS[RESET]}"
        cat "$temp_file"
    fi

    # Cleanup
    rm -f "$temp_file"
    return $exit_code
}

# Safe fallback for missing commands
progress_echo() {
    local message="$1"
    local color="${2:-WHITE}"

    if [[ "$PROGRESS_ENABLED" == "true" ]]; then
        echo -e "${PROGRESS_COLORS[$color]}$message${PROGRESS_COLORS[RESET]}"
    else
        echo "$message"
    fi
}

# Disable progress bars (fallback to normal output)
disable_progress() {
    PROGRESS_ENABLED=false
}

# Enable progress bars
enable_progress() {
    PROGRESS_ENABLED=true
}

# Check if we're in a compatible terminal
check_progress_compatibility() {
    # Disable progress in non-interactive sessions
    if [[ ! -t 1 ]]; then
        PROGRESS_ENABLED=false
        return 1
    fi

    # Disable if terminal doesn't support colors
    if [[ -z "$TERM" || "$TERM" == "dumb" ]]; then
        PROGRESS_ENABLED=false
        return 1
    fi

    return 0
}

# Auto-detect and set progress compatibility
auto_configure_progress() {
    if ! check_progress_compatibility; then
        progress_echo "Progress bars disabled - using standard output" "YELLOW"
    fi
}

# Initialize on source
auto_configure_progress
