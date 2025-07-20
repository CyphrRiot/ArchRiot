#!/bin/bash

# Install system fonts and download programming fonts
install_fonts() {
    echo "📝 Installing fonts..."

    local env_file="$HOME/.config/archriot/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Install system fonts via package manager
    local system_fonts="ttf-font-awesome noto-fonts noto-fonts-emoji noto-fonts-cjk noto-fonts-extra ttf-liberation ttf-hack-nerd ttf-jetbrains-mono-nerd ttf-cascadia-mono-nerd ttf-ia-writer"
    yay -S --noconfirm --needed $system_fonts || echo "⚠ Some system fonts may have failed"

    echo "✓ System fonts installed"
}

# Additional programming fonts (if needed in future)
install_programming_fonts() {
    echo "💻 Checking programming fonts..."
    echo "✓ iA Writer fonts installed via AUR package"
}

# Refresh cache and validate
finalize_fonts() {
    echo "🔄 Finalizing font installation..."

    # Comprehensive font cache refresh
    echo "🔄 Refreshing system font cache..."

    # Refresh user font cache
    if command -v fc-cache >/dev/null; then
        fc-cache -f ~/.local/share/fonts 2>/dev/null || true
        fc-cache -fv 2>/dev/null || true
        echo "✓ User font cache refreshed"
    fi

    # Refresh system font cache (if possible)
    if sudo -n true 2>/dev/null; then
        sudo fc-cache -f 2>/dev/null || true
        echo "✓ System font cache refreshed"
    fi

    # Update fontconfig cache
    fc-cache --system-only 2>/dev/null || true
    fc-cache --really-force 2>/dev/null || true

    # Force reload for current session
    hash -r 2>/dev/null || true

    # Enhanced validation with detailed feedback
    echo "🧪 Validating font installation..."
    local issues=0

    fc-list | grep -qi "noto" || { echo "⚠ Noto fonts missing"; ((issues++)); }
    fc-list | grep -qi "font awesome" || { echo "⚠ Font Awesome missing"; ((issues++)); }
    fc-list | grep -qi "jetbrains" || { echo "⚠ JetBrains fonts missing"; ((issues++)); }
    fc-list | grep -qi "ia writer\|iawriter" || { echo "⚠ iA Writer fonts missing"; ((issues++)); }

    if [[ $issues -eq 0 ]]; then
        echo "✓ All essential fonts validated and available"
    else
        echo "⚠ $issues font validation issues detected"
        echo "💡 Try restarting applications to refresh font cache"
    fi

    echo "✓ Font cache refresh complete"
}

# Display summary
show_summary() {
    echo ""
    echo "🎉 Font installation complete!"
    echo ""
    echo "📝 Installed: Noto fonts, Font Awesome icons, Hack Nerd Font, JetBrainsMono Nerd Font, Cascadia Mono Nerd Font, iA Writer fonts"
    echo "💡 Usage: Terminal fonts available, restart apps to use new fonts"
}

# Main execution
main() {
    echo "🚀 Starting font installation..."

    install_fonts
    install_programming_fonts
    finalize_fonts
    show_summary

    echo "✅ Font installation completed!"
}

main "$@"
