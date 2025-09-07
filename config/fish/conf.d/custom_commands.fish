# =============================================================================
# Custom Commands Loader
# =============================================================================
# This file automatically loads custom commands from ~/.config/fish/custom_commands.fish
# if it exists. This follows Fish shell best practices for modular configuration.

# Define the custom commands file path
set -l custom_commands_file "$HOME/.config/fish/custom_commands.fish"

# Load custom commands if the file exists
if test -f "$custom_commands_file"
    source "$custom_commands_file"
    echo "✓ Loaded custom commands from $custom_commands_file"
end
