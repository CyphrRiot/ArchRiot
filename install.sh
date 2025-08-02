#!/bin/bash

# ================================================================================
# ArchRiot Installation System v2.5.0 - YAML Processing Engine
# ================================================================================
# Clean implementation focused on YAML-driven package installation
# Replaces 30+ shell scripts with unified YAML configuration
# ================================================================================

# Installation configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly INSTALL_DIR="$HOME/.local/share/archriot/install"
readonly LOG_FILE="$HOME/.cache/archriot/install.log"
readonly ERROR_LOG="$HOME/.cache/archriot/install-errors.log"

# YAML processing configuration
readonly YAML_CONFIG_FILE="$INSTALL_DIR/packages.yaml"

# Read version from VERSION file
if [ -f "$SCRIPT_DIR/VERSION" ]; then
    readonly ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
else
    readonly ARCHRIOT_VERSION="2.5.0-dev"
fi

# ================================================================================
# Logging Functions
# ================================================================================

# Initialize logging
init_logging() {
    mkdir -p "$(dirname "$LOG_FILE")"
    mkdir -p "$(dirname "$ERROR_LOG")"

    echo "=== ArchRiot YAML Installation v$ARCHRIOT_VERSION - $(date) ===" > "$LOG_FILE"
    echo "=== ArchRiot YAML Installation Errors v$ARCHRIOT_VERSION - $(date) ===" > "$ERROR_LOG"
}

# Log messages
log_message() {
    local level="$1"
    local message="$2"
    local timestamp="[$(date '+%H:%M:%S')]"

    case "$level" in
        "INFO")
            echo "$timestamp $message" | tee -a "$LOG_FILE"
            ;;
        "SUCCESS")
            echo "$timestamp ✅ $message" | tee -a "$LOG_FILE"
            ;;
        "WARNING")
            echo "$timestamp ⚠️  $message" | tee -a "$LOG_FILE"
            ;;
        "ERROR")
            echo "$timestamp ❌ $message" | tee -a "$LOG_FILE" "$ERROR_LOG"
            ;;
    esac
}

# ================================================================================
# Dependency Management
# ================================================================================

# Install correct yq version (go-yq)
install_go_yq() {
    log_message "INFO" "Checking yq dependency..."

    # Check if we have the correct yq version
    if command -v yq >/dev/null 2>&1; then
        local yq_version=$(yq --version 2>/dev/null)
        if [[ "$yq_version" == *"mikefarah"* ]]; then
            log_message "SUCCESS" "Correct yq version already installed: $yq_version"
            return 0
        else
            log_message "WARNING" "Wrong yq version detected, replacing with go-yq..."
            # Remove old python-based yq
            if yay -R yq --noconfirm >> "$LOG_FILE" 2>&1; then
                log_message "SUCCESS" "Old yq removed"
            else
                log_message "WARNING" "Could not remove old yq, continuing anyway"
            fi
        fi
    fi

    log_message "INFO" "Installing go-yq..."
    if yay -S go-yq --noconfirm >> "$LOG_FILE" 2>&1; then
        log_message "SUCCESS" "go-yq installed successfully"
        # Refresh PATH
        export PATH="/usr/bin:$PATH"
        hash -r 2>/dev/null || true
        return 0
    else
        log_message "ERROR" "Failed to install go-yq"
        return 1
    fi
}

# Validate YAML file syntax
validate_yaml() {
    local yaml_file="$1"

    log_message "INFO" "Validating YAML syntax: $yaml_file"

    if [[ ! -f "$yaml_file" ]]; then
        log_message "ERROR" "YAML file not found: $yaml_file"
        return 1
    fi

    # Test YAML syntax
    if yq eval '.' "$yaml_file" >/dev/null 2>&1; then
        log_message "SUCCESS" "YAML syntax is valid"
    else
        log_message "ERROR" "YAML syntax error in: $yaml_file"
        log_message "ERROR" "Please fix YAML syntax before continuing"
        return 1
    fi

    # Test if we can read core modules
    local core_modules=$(yq eval '.core | keys' "$yaml_file" 2>/dev/null)
    if [[ -n "$core_modules" ]]; then
        log_message "SUCCESS" "YAML structure validated - core modules found"
        return 0
    else
        log_message "ERROR" "YAML structure invalid - no core modules found"
        return 1
    fi
}

# ================================================================================
# YAML Processing Functions
# ================================================================================

# Parse YAML packages for a given module
parse_yaml_packages() {
    local module_key="$1"

    if [[ ! -f "$YAML_CONFIG_FILE" ]]; then
        log_message "ERROR" "YAML config file not found: $YAML_CONFIG_FILE"
        return 1
    fi

    # Extract packages for the specified module
    yq eval ".${module_key}.packages[]?" "$YAML_CONFIG_FILE" 2>/dev/null
}

# Parse YAML configs for a given module
parse_yaml_configs() {
    local module_key="$1"

    if [[ ! -f "$YAML_CONFIG_FILE" ]]; then
        log_message "ERROR" "YAML config file not found: $YAML_CONFIG_FILE"
        return 1
    fi

    # Extract configs for the specified module
    yq eval ".${module_key}.configs[]?" "$YAML_CONFIG_FILE" 2>/dev/null
}

# Install packages from YAML for a given module
install_yaml_packages() {
    local module_key="$1"

    log_message "INFO" "Installing packages for module: $module_key"

    local packages
    packages=$(parse_yaml_packages "$module_key")

    if [[ -z "$packages" ]]; then
        log_message "WARNING" "No packages found for module: $module_key"
        return 0
    fi

    # Convert to array and install
    local package_array=($packages)
    if [[ ${#package_array[@]} -gt 0 ]]; then
        log_message "INFO" "Installing ${#package_array[@]} packages: ${package_array[*]}"

        if command -v yay >/dev/null 2>&1; then
            yay -S --noconfirm "${package_array[@]}" >> "$LOG_FILE" 2>&1
        elif command -v pacman >/dev/null 2>&1; then
            sudo pacman -S --noconfirm "${package_array[@]}" >> "$LOG_FILE" 2>&1
        else
            log_message "ERROR" "No package manager found"
            return 1
        fi

        if [[ $? -eq 0 ]]; then
            log_message "SUCCESS" "Packages installed successfully for $module_key"
        else
            log_message "ERROR" "Failed to install packages for $module_key"
            return 1
        fi
    fi
}

# Copy config files from YAML for a given module
copy_yaml_configs() {
    local module_key="$1"

    log_message "INFO" "Copying configs for module: $module_key"

    local configs
    configs=$(parse_yaml_configs "$module_key")

    if [[ -z "$configs" ]]; then
        log_message "WARNING" "No configs found for module: $module_key"
        return 0
    fi

    # Process each config pattern
    while IFS= read -r config_pattern; do
        [[ -z "$config_pattern" ]] && continue

        local source_path="$SCRIPT_DIR/dotfiles/$config_pattern"
        local dest_path="$HOME/.config/$config_pattern"

        # Create destination directory
        mkdir -p "$(dirname "$dest_path")"

        # Copy files matching pattern
        if [[ -d "$source_path" ]]; then
            cp -r "$source_path"/* "$(dirname "$dest_path")/" >> "$LOG_FILE" 2>&1
            log_message "INFO" "Copied directory: $config_pattern"
        elif [[ -f "$source_path" ]]; then
            cp "$source_path" "$dest_path" >> "$LOG_FILE" 2>&1
            log_message "INFO" "Copied file: $config_pattern"
        else
            log_message "WARNING" "Config not found: $source_path"
        fi
    done <<< "$configs"
}

# ================================================================================
# Testing Functions
# ================================================================================

# Test YAML engine with real packages
test_yaml_engine() {
    local test_module="core.base"

    log_message "INFO" "🧪 Testing YAML engine with module: $test_module"

    # Test package parsing
    log_message "INFO" "Testing package parsing..."
    local packages
    packages=$(parse_yaml_packages "$test_module")
    if [[ -n "$packages" ]]; then
        log_message "SUCCESS" "Packages found: $packages"
    else
        log_message "ERROR" "No packages found for $test_module"
        return 1
    fi

    # Test config parsing
    log_message "INFO" "Testing config parsing..."
    local configs
    configs=$(parse_yaml_configs "$test_module")
    if [[ -n "$configs" ]]; then
        log_message "SUCCESS" "Configs found: $configs"
    else
        log_message "WARNING" "No configs found for $test_module"
    fi

    # Test actual package installation (safe on existing system)
    log_message "INFO" "Testing package installation..."
    if install_yaml_packages "$test_module"; then
        log_message "SUCCESS" "Package installation test completed"
    else
        log_message "ERROR" "Package installation test failed"
        return 1
    fi

    log_message "SUCCESS" "🎉 YAML engine test completed successfully!"
    return 0
}

# ================================================================================
# Main Function
# ================================================================================

main() {
    init_logging

    log_message "INFO" "🚀 ArchRiot YAML Processing Engine v$ARCHRIOT_VERSION"

    # STEP 1: Install correct yq dependency
    if ! install_go_yq; then
        log_message "ERROR" "Failed to install yq dependency"
        exit 1
    fi

    # STEP 2: Validate YAML file before doing anything else
    if [[ -f "$YAML_CONFIG_FILE" ]]; then
        if ! validate_yaml "$YAML_CONFIG_FILE"; then
            log_message "ERROR" "YAML validation failed - cannot continue"
            exit 1
        fi
    else
        log_message "ERROR" "YAML config file not found: $YAML_CONFIG_FILE"
        log_message "ERROR" "Please ensure packages.yaml exists before running installer"
        exit 1
    fi

    # STEP 3: Test YAML processing engine
    log_message "INFO" "Testing YAML processing functionality..."
    if test_yaml_engine; then
        log_message "SUCCESS" "YAML processing engine fully validated"
    else
        log_message "ERROR" "YAML engine test failed"
        exit 1
    fi

    log_message "SUCCESS" "🎉 ArchRiot YAML Processing Engine ready for full deployment"
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
