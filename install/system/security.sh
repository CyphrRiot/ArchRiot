#!/bin/bash

# ==============================================================================
# ArchRiot Security Configuration
# ==============================================================================
# Configures system security components including PAM integration for
# gnome-keyring to enable automatic keyring unlock on login
# ==============================================================================

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/../lib/display.sh"

# Display section header
display_header "Security Configuration"

# ==============================================================================
# Configure PAM for Gnome Keyring Integration
# ==============================================================================

configure_gnome_keyring_pam() {
    display_info "Configuring PAM integration for gnome-keyring..."

    # Backup the original system-auth file
    if [[ ! -f /etc/pam.d/system-auth.backup ]]; then
        sudo cp /etc/pam.d/system-auth /etc/pam.d/system-auth.backup
        echo "✓ Created backup of /etc/pam.d/system-auth"
    fi

    # Check if gnome-keyring PAM modules are already configured
    if grep -q "pam_gnome_keyring.so" /etc/pam.d/system-auth; then
        echo "✓ Gnome keyring PAM modules already configured"
        return 0
    fi

    # Add gnome-keyring PAM modules to system-auth
    echo "⚙ Adding gnome-keyring PAM integration..."

    # Add auth module after pam_unix.so auth line
    sudo sed -i '/auth.*pam_unix.so.*try_first_pass nullok/a auth       optional                    pam_gnome_keyring.so' /etc/pam.d/system-auth

    # Add session module after pam_unix.so session line
    sudo sed -i '/session.*pam_unix.so/a session    optional                    pam_gnome_keyring.so auto_start' /etc/pam.d/system-auth

    # Add password module after pam_unix.so password line
    sudo sed -i '/password.*pam_unix.so.*try_first_pass nullok shadow/a password   optional                    pam_gnome_keyring.so' /etc/pam.d/system-auth

    echo "✓ Gnome keyring PAM integration configured"
}

# ==============================================================================
# Clean up duplicate PAM configurations
# ==============================================================================

cleanup_duplicate_pam_configs() {
    display_info "Cleaning up duplicate PAM configurations..."

    # Remove any duplicate gnome-keyring entries from login and system-login
    # since they should be handled by system-auth
    if grep -q "pam_gnome_keyring.so" /etc/pam.d/login; then
        sudo sed -i '/pam_gnome_keyring.so/d' /etc/pam.d/login
        echo "✓ Removed duplicate keyring config from /etc/pam.d/login"
    fi

    if grep -q "pam_gnome_keyring.so" /etc/pam.d/system-login; then
        sudo sed -i '/pam_gnome_keyring.so/d' /etc/pam.d/system-login
        echo "✓ Removed duplicate keyring config from /etc/pam.d/system-login"
    fi
}

# ==============================================================================
# Validate Security Configuration
# ==============================================================================

validate_security_config() {
    display_info "Validating security configuration..."

    local validation_passed=true

    # Check if gnome-keyring is installed
    if ! command -v gnome-keyring-daemon &>/dev/null; then
        echo "⚠ gnome-keyring not found - ensure desktop apps are installed"
        validation_passed=false
    else
        echo "✓ gnome-keyring daemon available"
    fi

    # Check PAM configuration
    if grep -q "auth.*optional.*pam_gnome_keyring.so" /etc/pam.d/system-auth; then
        echo "✓ PAM auth module configured"
    else
        echo "⚠ PAM auth module missing"
        validation_passed=false
    fi

    if grep -q "session.*optional.*pam_gnome_keyring.so.*auto_start" /etc/pam.d/system-auth; then
        echo "✓ PAM session module configured"
    else
        echo "⚠ PAM session module missing"
        validation_passed=false
    fi

    if grep -q "password.*optional.*pam_gnome_keyring.so" /etc/pam.d/system-auth; then
        echo "✓ PAM password module configured"
    else
        echo "⚠ PAM password module missing"
        validation_passed=false
    fi

    if [[ "$validation_passed" == true ]]; then
        echo "✅ Security configuration validation passed"
        return 0
    else
        echo "❌ Security configuration validation failed"
        return 1
    fi
}

# ==============================================================================
# Display Summary
# ==============================================================================

display_security_summary() {
    echo ""
    echo "🔒 Security configuration complete!"
    echo ""
    echo "📋 Configured components:"
    echo "  • Gnome keyring PAM integration"
    echo "  • Automatic keyring unlock on login"
    echo "  • Keyring password sync with login password"
    echo ""
    echo "ℹ️  Changes take effect on next login/reboot"
    echo "   Your keyring will automatically unlock with your login password"
}

# ==============================================================================
# Main Installation Function
# ==============================================================================

install_security() {
    display_info "Starting security configuration..."

    # Configure gnome-keyring PAM integration
    configure_gnome_keyring_pam

    # Clean up any duplicate configurations
    cleanup_duplicate_pam_configs

    # Validate the configuration
    if validate_security_config; then
        display_security_summary
        return 0
    else
        display_error "Security configuration failed validation"
        return 1
    fi
}

# ==============================================================================
# Script Execution
# ==============================================================================

# Run the installation if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_security
fi

echo "✅ Security configuration setup complete!"
