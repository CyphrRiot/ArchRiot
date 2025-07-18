# =============================================================================
# Dependency-Free Fish Configuration with Git Integration
# =============================================================================

# Greeting with fastfetch if available - with delay for terminal sizing
function fish_greeting
    if command -v fastfetch >/dev/null 2>&1
        # Small delay to ensure terminal is fully sized
        sleep 0.1
        command fastfetch --logo-width 20 --logo arch
    end
end

# =============================================================================
# Path Configuration
# =============================================================================

# Add local bin directories to PATH (last prepend gets highest priority)
fish_add_path --prepend $HOME/.local/share/archriot/bin
fish_add_path --prepend $HOME/.local/bin

# Git Prompt Configuration
set -g __fish_git_prompt_showdirtystate yes
set -g __fish_git_prompt_showstashstate yes
set -g __fish_git_prompt_showuntrackedfiles yes
set -g __fish_git_prompt_showupstream yes
set -g __fish_git_prompt_color_branch yellow
set -g __fish_git_prompt_color_upstream_ahead purple
set -g __fish_git_prompt_color_upstream_behind red
set -g __fish_git_prompt_color_dirtystate 'ff8c00'
set -g __fish_git_prompt_char_dirtystate "●"
set -g __fish_git_prompt_char_stagedstate "→"
set -g __fish_git_prompt_char_untrackedfiles "☡"
set -g __fish_git_prompt_char_stashstate "↩"
set -g __fish_git_prompt_char_upstream_ahead "+"
set -g __fish_git_prompt_char_upstream_behind "-"
set -g __fish_git_prompt_char_upstream_equal ""

# Enhanced Prompt Function
function fish_prompt
    set -l last_status $status
    echo -n "Ω "  # 🐧 (penguin - commented for easy revert)
    set_color blue
    printf "%s" (string replace $HOME "~" (pwd))
    set_color normal
    printf "%s" (__fish_git_prompt)
    if test $last_status -ne 0
        set_color red
        printf " [%d]" $last_status
        set_color normal
    end
    set_color cyan
    printf " ➤ "
    set_color normal
end

# =============================================================================
# Aliases & Functions
# =============================================================================

# Vim alias to nvim
alias vim='nvim'

# Fastfetch with correct logo width
alias fastfetch='command fastfetch --logo-width 20'

# Disk usage - show top 10 largest items by size
function dum
    du -sm * | sort -nr | head -10
end

# =============================================================================
# Hyprland Auto-start
# =============================================================================

# Auto-start Hyprland on TTY1
if status is-login && test (tty) = /dev/tty1
    exec Hyprland
end
