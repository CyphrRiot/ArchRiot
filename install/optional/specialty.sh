#!/bin/bash

# ==============================================================================
# ArchRiot Specialty Applications Setup
# ==============================================================================
# Optional specialty applications installation supporting:
# - Financial tools (cryptocurrency wallets, trading)
# - Creative writing tools (advanced editors, publishers)
# - Research and academic tools (reference managers, note-taking)
# - Media downloaders and converters
# - Experimental and niche applications
# ==============================================================================

# Load shared library functions
source "$(dirname "${BASH_SOURCE[0]}")/../lib/shared.sh" || {
    echo "❌ Failed to load shared library functions"
    exit 1
}

# Initialize installer
init_installer "specialty"

# Detect existing specialty applications
detect_existing_specialty() {
    echo "🔍 Detecting existing specialty applications..."

    local specialty_apps=(
        "feather-wallet-bin:Feather Monero Wallet"
        "gnome-text-editor:Gnome Text Editor"
        "papers:Papers Reference Manager"
        "yt-dlp:YouTube Downloader"
        "spotdl:Spotify Downloader"
    )

    local found_apps=()
    echo "Found specialty applications:"

    for app_info in "${specialty_apps[@]}"; do
        IFS=':' read -r cmd name <<< "$app_info"
        if command -v "$cmd" &>/dev/null; then
            echo "  ✓ $name"
            found_apps+=("$cmd")
        fi
    done

    DETECTED_SPECIALTY_APPS=("${found_apps[@]}")
    return ${#found_apps[@]}
}

# Install financial applications
install_financial_apps() {
    echo "💰 Installing financial applications..."

    local financial_apps=(
        "feather-wallet-bin"    # Feather Monero wallet
    )

    install_packages "Financial Applications" "${financial_apps[@]}"
}

# Install creative writing tools
install_writing_tools() {
    echo "✍️ Installing advanced writing tools..."

    local writing_apps=(
        "gnome-text-editor"     # Modern text editor with Tokyo Night theme
        "abiword"               # Lightweight word processor
    )

    install_packages "Writing Tools" "${writing_apps[@]}"
}

# Install research and academic tools
install_research_tools() {
    echo "📚 Installing research and academic tools..."

    local research_apps=(
        "papers"                # Papers reference manager
    )

    install_packages "Research Tools" "${research_apps[@]}"
}

# Install media downloaders and converters
install_media_tools() {
    echo "📥 Installing media downloaders and converters..."

    local media_apps=(
        "yt-dlp"                # YouTube and media downloader
        "ffmpeg"                # Media converter
    )

    # Install regular media packages
    install_packages "Media Tools" "${media_apps[@]}"

    # Install spotdl separately with --nocheck flag (API test failures)
    print_status "INFO" "Installing spotdl (requires --nocheck due to API test failures)"
    install_aur_nocheck "spotdl" "optional"
}

# Install experimental and niche applications
install_experimental_apps() {
    echo "🧪 Installing experimental applications..."

    local experimental_apps=(
        "lollypop"              # Modern GTK music player
    )

    install_packages "Experimental Tools" "${experimental_apps[@]}"
}

# Configure specialty applications
configure_specialty_apps() {
    echo "⚙️ Configuring specialty applications..."

    # Set up file associations for writing tools
    if command -v xdg-mime >/dev/null; then
        if command -v gnome-text-editor >/dev/null; then
            xdg-mime default org.gnome.TextEditor.desktop text/markdown 2>/dev/null || true
        fi
    fi

    # Create desktop shortcuts for research tools
    local shortcuts_dir="$HOME/.local/share/applications"
    mkdir -p "$shortcuts_dir"

    # Configure financial app security settings
    if command -v feather-wallet-bin >/dev/null; then
        echo "💡 Remember to backup your wallet seeds securely"
    fi
}

# Validate specialty applications installation
validate_specialty_apps() {
    echo "🧪 Validating specialty applications..."

    local validation_errors=0

    # Test financial apps
    local financial_apps=("feather-wallet-bin" "electrum")
    for app in "${financial_apps[@]}"; do
        if command -v "$app" >/dev/null; then
            echo "✓ $app available"
        else
            echo "⚠ $app not found"
            ((validation_errors++))
        fi
    done

    # Test writing tools
    local writing_apps=("gnome-text-editor" "typora")
    for app in "${writing_apps[@]}"; do
        if command -v "$app" >/dev/null; then
            echo "✓ $app available"
        else
            echo "⚠ $app not found"
            ((validation_errors++))
        fi
    done

    # Test media tools
    local media_apps=("yt-dlp" "ffmpeg")
    for app in "${media_apps[@]}"; do
        if command -v "$app" >/dev/null; then
            echo "✓ $app available"
        else
            echo "⚠ $app not found"
            ((validation_errors++))
        fi
    done

    if [[ $validation_errors -eq 0 ]]; then
        echo "✅ All specialty applications validated successfully"
        return 0
    else
        echo "⚠ $validation_errors validation errors found"
        return 1
    fi
}

# Display installation summary
display_specialty_summary() {
    echo ""
    echo "🎉 Specialty applications setup complete!"
    echo ""
    echo "💰 Financial Tools:"
    echo "  • Feather Wallet - Lightweight Monero wallet"
    echo ""
    echo "✍️ Creative Writing:"
    echo "  • Gnome Text Editor - Modern text editor with Tokyo Night theme"
    echo "  • AbiWord - Lightweight word processor"
    echo ""
    echo "📚 Research & Academic:"
    echo "  • Papers - Academic reference manager"
    echo ""
    echo "📥 Media Tools:"
    echo "  • yt-dlp - Universal media downloader"
    echo "  • spotdl - Spotify music downloader"
    echo "  • FFmpeg - Media conversion Swiss Army knife"
    echo ""
    echo "🧪 Experimental Tools:"
    echo "  • Lollypop - Modern GTK music player"
    echo ""
    echo "💡 Quick Start Tips:"
    echo "  • Launch writing tools from applications menu"
    echo "  • Financial apps: Always backup wallet seeds!"
    echo "  • Use 'yt-dlp --help' for download options"
    echo ""
}

# Main installation function
main() {
    echo "🚀 Starting specialty applications setup..."

    load_user_environment

    # Show current state
    detect_existing_specialty

    # Install specialty applications by category
    install_financial_apps
    install_writing_tools
    install_research_tools
    install_media_tools
    install_experimental_apps

    # Configure applications
    configure_specialty_apps

    # Validate installation
    if validate_specialty_apps; then
        display_specialty_summary
    else
        echo "⚠ Some specialty applications may not be functioning correctly"
        echo "Check the installation logs above for details"
    fi

    echo "✅ Specialty applications setup complete!"
    show_install_summary
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
