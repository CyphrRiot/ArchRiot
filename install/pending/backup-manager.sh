#!/bin/bash

# ================================================================================
# ArchRiot Centralized Backup Manager
# ================================================================================
# Single backup system for all ArchRiot scripts
# Location: ~/.archriot/backups/
# Keeps only 3 most recent backups, auto-cleanup older ones
# ================================================================================

# Centralized backup configuration
readonly BACKUP_ROOT="$HOME/.archriot/backups"
readonly MAX_BACKUPS=3

# Colors for output (when needed)
BACKUP_RED='\033[0;31m'
BACKUP_GREEN='\033[0;32m'
BACKUP_YELLOW='\033[1;33m'
BACKUP_BLUE='\033[0;34m'
BACKUP_NC='\033[0m' # No Color

# Initialize backup system
init_backup_system() {
    mkdir -p "$BACKUP_ROOT"

    # Ensure proper permissions
    chmod 755 "$BACKUP_ROOT"

    # Log initialization
    echo "INFO: Backup system initialized at $BACKUP_ROOT" >> "${LOG_FILE:-/dev/null}" 2>&1
}

# Create a new backup with descriptive name
create_backup() {
    local backup_type="$1"
    local description="$2"

    # Validate input
    if [[ -z "$backup_type" ]]; then
        echo "ERROR: Backup type required" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi

    # Initialize backup system if needed
    init_backup_system

    # Generate backup directory name
    local timestamp=$(date +%Y-%m-%d-%H%M%S)
    local backup_dir="$BACKUP_ROOT/${backup_type}-${timestamp}"

    mkdir -p "$backup_dir"

    # Save backup metadata
    cat > "$backup_dir/BACKUP_INFO" << EOF
backup_type=$backup_type
description=${description:-"$backup_type backup"}
timestamp=$timestamp
created_by=${0##*/}
archriot_version=${ARCHRIOT_VERSION:-"unknown"}
created=$(date)
EOF

    echo "$backup_dir"
}

# Backup specific config directories
backup_configs() {
    local backup_type="$1"
    local configs_to_backup=("${@:2}")

    if [[ ${#configs_to_backup[@]} -eq 0 ]]; then
        echo "ERROR: No configs specified for backup" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi

    local backup_dir=$(create_backup "$backup_type" "Configuration backup")
    local backed_up=0

    for config in "${configs_to_backup[@]}"; do
        local source="$HOME/.config/$config"
        local target="$backup_dir/$config"

        if [[ -e "$source" ]]; then
            mkdir -p "$(dirname "$target")"
            if cp -R "$source" "$target" 2>/dev/null; then
                ((backed_up++))
                echo "INFO: Backed up: $config" >> "${LOG_FILE:-/dev/null}" 2>&1
            fi
        fi
    done

    # Save manifest of what was backed up
    printf '%s\n' "${configs_to_backup[@]}" > "$backup_dir/MANIFEST"

    if [[ $backed_up -gt 0 ]]; then
        echo "INFO: ✓ Configuration backup created: $backup_dir ($backed_up items)" >> "${LOG_FILE:-/dev/null}" 2>&1
        cleanup_old_backups "$backup_type"
        echo "$backup_dir"
    else
        # Remove empty backup
        rm -rf "$backup_dir"
        echo "INFO: No configurations found to backup" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi
}

# Backup entire ArchRiot installation
backup_archriot_install() {
    local reason="$1"
    local source_dir="$HOME/.local/share/archriot"

    if [[ ! -d "$source_dir" ]]; then
        echo "INFO: No existing ArchRiot installation to backup" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 0
    fi

    local backup_dir=$(create_backup "archriot" "${reason:-"ArchRiot installation backup"}")

    if cp -R "$source_dir" "$backup_dir/archriot" 2>/dev/null; then
        echo "INFO: ✓ ArchRiot installation backed up: $backup_dir" >> "${LOG_FILE:-/dev/null}" 2>&1
        cleanup_old_backups "archriot"
        echo "$backup_dir"
    else
        rm -rf "$backup_dir"
        echo "ERROR: Failed to backup ArchRiot installation" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi
}

# Restore from backup
restore_backup() {
    local backup_dir="$1"
    local target_base="${2:-$HOME/.config}"

    if [[ ! -d "$backup_dir" || ! -f "$backup_dir/MANIFEST" ]]; then
        echo "ERROR: Invalid backup directory or missing manifest: $backup_dir" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi

    echo "INFO: 🔄 Restoring from backup: $backup_dir" >> "${LOG_FILE:-/dev/null}" 2>&1

    local restored=0
    while IFS= read -r config; do
        local source="$backup_dir/$config"
        local target="$target_base/$config"

        if [[ -e "$source" ]]; then
            rm -rf "$target"
            if cp -R "$source" "$target" 2>/dev/null; then
                ((restored++))
                echo "INFO: Restored: $config" >> "${LOG_FILE:-/dev/null}" 2>&1
            fi
        fi
    done < "$backup_dir/MANIFEST"

    echo "INFO: ✓ Restored $restored items from backup" >> "${LOG_FILE:-/dev/null}" 2>&1
    return 0
}

# Clean up old backups, keeping only MAX_BACKUPS most recent
cleanup_old_backups() {
    local backup_type="$1"

    if [[ ! -d "$BACKUP_ROOT" ]]; then
        return 0
    fi

    # Get all backups of this type, sorted by date (newest first)
    local backups=($(find "$BACKUP_ROOT" -maxdepth 1 -type d -name "${backup_type}-*" | sort -r))

    # Remove excess backups
    if [[ ${#backups[@]} -gt $MAX_BACKUPS ]]; then
        local removed=0
        for ((i=MAX_BACKUPS; i<${#backups[@]}; i++)); do
            rm -rf "${backups[$i]}" 2>/dev/null && ((removed++))
        done

        if [[ $removed -gt 0 ]]; then
            echo "INFO: ✓ Cleaned up $removed old $backup_type backups (keeping $MAX_BACKUPS most recent)" >> "${LOG_FILE:-/dev/null}" 2>&1
        fi
    fi
}

# List all backups
list_backups() {
    local backup_type="$1"

    if [[ ! -d "$BACKUP_ROOT" ]]; then
        echo "No backups found"
        return 0
    fi

    local pattern="*"
    if [[ -n "$backup_type" ]]; then
        pattern="${backup_type}-*"
    fi

    echo "Available backups in $BACKUP_ROOT:"
    for backup_dir in "$BACKUP_ROOT"/$pattern; do
        if [[ -d "$backup_dir" && -f "$backup_dir/BACKUP_INFO" ]]; then
            local info_file="$backup_dir/BACKUP_INFO"
            local description=$(grep "^description=" "$info_file" | cut -d'=' -f2-)
            local created=$(grep "^created=" "$info_file" | cut -d'=' -f2-)

            echo "  $(basename "$backup_dir"): $description ($created)"
        fi
    done
}

# Get most recent backup of type
get_latest_backup() {
    local backup_type="$1"

    if [[ ! -d "$BACKUP_ROOT" ]]; then
        return 1
    fi

    # Find most recent backup of this type
    local latest=$(find "$BACKUP_ROOT" -maxdepth 1 -type d -name "${backup_type}-*" | sort -r | head -1)

    if [[ -n "$latest" && -d "$latest" ]]; then
        echo "$latest"
        return 0
    else
        return 1
    fi
}

# Emergency backup function - quick and dirty backup before risky operations
emergency_backup() {
    local reason="$1"

    echo "INFO: 🚨 Creating emergency backup: ${reason:-"Emergency backup"}" >> "${LOG_FILE:-/dev/null}" 2>&1

    # Backup critical configs
    local critical_configs=("hypr" "waybar" "archriot")
    backup_configs "emergency" "${critical_configs[@]}"
}

# Check if backup system is healthy
check_backup_health() {
    if [[ ! -d "$BACKUP_ROOT" ]]; then
        echo "WARN: Backup system not initialized" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi

    # Check permissions
    if [[ ! -w "$BACKUP_ROOT" ]]; then
        echo "ERROR: Backup directory not writable: $BACKUP_ROOT" >> "${LOG_FILE:-/dev/null}" 2>&1
        return 1
    fi

    # Check disk space (warn if less than 100MB available)
    local available=$(df "$BACKUP_ROOT" | awk 'NR==2 {print $4}')
    if [[ $available -lt 102400 ]]; then  # 100MB in KB
        echo "WARN: Low disk space for backups: ${available}KB available" >> "${LOG_FILE:-/dev/null}" 2>&1
    fi

    echo "INFO: ✓ Backup system healthy" >> "${LOG_FILE:-/dev/null}" 2>&1
    return 0
}

# Export functions for use by other scripts
export -f init_backup_system
export -f create_backup
export -f backup_configs
export -f backup_archriot_install
export -f restore_backup
export -f cleanup_old_backups
export -f list_backups
export -f get_latest_backup
export -f emergency_backup
export -f check_backup_health
