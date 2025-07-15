#!/bin/bash

# Install system fonts and download programming fonts
install_fonts() {
    echo "📝 Installing fonts..."

    local env_file="$HOME/.config/omarchy/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Install system fonts via package manager
    local system_fonts="ttf-font-awesome noto-fonts noto-fonts-emoji noto-fonts-cjk noto-fonts-extra ttf-liberation ttf-jetbrains-mono-nerd ttf-cascadia-mono-nerd"
    yay -S --noconfirm --needed $system_fonts || echo "⚠ Some system fonts may have failed"

    echo "✓ System fonts installed"
}

# Download and install programming fonts
install_programming_fonts() {
    echo "💻 Installing programming fonts..."

    mkdir -p ~/.local/share/fonts

    # Install iA Writer Mono fonts
    if ! fc-list | grep -qi "iA Writer Mono"; then
        echo "⬇️  Downloading iA Writer fonts..."
        local temp_dir="/tmp/iawriter-$$"
        mkdir -p "$temp_dir"

        if wget -q "https://github.com/iaolo/iA-Fonts/archive/refs/heads/master.zip" -O "$temp_dir/iafonts.zip" &&
           cd "$temp_dir" && unzip -q iafonts.zip &&
           find . -name "iAWriterMonoS-*.ttf" -exec cp {} ~/.local/share/fonts/ \;; then
            echo "✓ iA Writer Mono fonts installed"
        else
            echo "⚠ iA Writer fonts installation failed"
        fi
        rm -rf "$temp_dir"
    fi
}

# Refresh cache and validate
finalize_fonts() {
    echo "🔄 Finalizing font installation..."

    # Refresh font cache
    command -v fc-cache >/dev/null && fc-cache -f ~/.local/share/fonts 2>/dev/null

    # Quick validation
    local issues=0
    fc-list | grep -qi "noto" || ((issues++))
    fc-list | grep -qi "font awesome" || ((issues++))

    if [[ $issues -eq 0 ]]; then
        echo "✓ Essential fonts validated"
    else
        echo "⚠ $issues font issues detected"
    fi

    echo "✓ Font cache refreshed"
}

# Display summary
show_summary() {
    echo ""
    echo "🎉 Font installation complete!"
    echo ""
    echo "📝 Installed: Noto fonts, Font Awesome icons, JetBrainsMono Nerd Font, Cascadia Mono Nerd Font, iA Writer Mono"
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
