#!/bin/bash

# Load user environment and install packages
setup_packages() {
    echo "📦 Installing terminal tools and shell..."
    local env_file="$HOME/.config/omarchy/user.env"
    [[ -f "$env_file" ]] && source "$env_file"

    # Install essentials (critical)
    local essentials="wget curl unzip inetutils git neovim ghostty ghostty-shell-integration"
    yay -S --noconfirm --needed $essentials || return 1

    # Install shell tools (best effort)
    local shell_tools="fish fd eza fzf ripgrep zoxide bat lsd fastfetch btop"
    yay -S --noconfirm --needed $shell_tools || echo "⚠ Some shell tools may have failed"

    # Install utilities (nice-to-have)
    local utilities="wl-clipboard man tldr less whois plocate bash-completion"
    yay -S --noconfirm --needed $utilities || true

    echo "✓ Package installation completed"
}

# Setup and validate fish shell
setup_fish() {
    echo "🐠 Setting up Fish shell..."

    command -v fish &>/dev/null || return 1

    # Setup fish config
    mkdir -p ~/.config/fish
    local fish_config="$HOME/.local/share/omarchy/config/fish/config.fish"
    [[ -f "$fish_config" ]] || return 1
    cp "$fish_config" ~/.config/fish/config.fish || return 1

    # Validate configuration
    fish -c "source ~/.config/fish/config.fish" 2>/dev/null || {
        echo "❌ Fish configuration has syntax errors"
        return 1
    }

    echo "✓ Fish configuration installed and validated"
}

# Set fish as default shell
set_default_shell() {
    echo "🔧 Setting fish as default shell..."

    local fish_path="/usr/bin/fish"
    [[ -f "$fish_path" ]] || return 1
    [[ "$SHELL" == "$fish_path" ]] && return 0

    if sudo chsh -s "$fish_path" "$USER"; then
        echo "✓ Fish set as default shell (takes effect on next login)"
    else
        echo "⚠ Failed to set fish as default - run: sudo chsh -s /usr/bin/fish $USER"
        return 1
    fi
}

# Display setup summary
show_summary() {
    echo ""
    echo "🎉 Terminal and shell setup complete!"
    echo ""
    echo "📋 Installed: Fish shell, modern CLI tools, Ghostty terminal, Neovim"
    echo "✨ Features: Ω prompt, Git integration, Fastfetch greeting, vim→nvim alias"
    echo "🚀 Quick start: Type 'fish' to test or log out/in for default shell"

    # Show prompt preview
    command -v fish &>/dev/null && {
        echo ""
        echo "🐠 Prompt preview:"
        fish -c "echo -n 'Example: '; fish_prompt; echo 'your_command_here'" 2>/dev/null || true
    }
}

# Main execution
main() {
    echo "🚀 Starting terminal and shell setup..."

    setup_packages || return 1
    setup_fish || return 1
    set_default_shell || echo "⚠ Could not set fish as default shell"
    show_summary

    echo "✅ Terminal and shell setup completed!"
}

main "$@"
