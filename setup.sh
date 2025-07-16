ascii_art=' ██████╗ ██╗  ██╗███╗   ███╗ █████╗ ██████╗  ██████╗██╗  ██╗██╗   ██╗
██╔═══██╗██║  ██║████╗ ████║██╔══██╗██╔══██╗██╔════╝██║  ██║╚██╗ ██╔╝
██║   ██║███████║██╔████╔██║███████║██████╔╝██║     ███████║ ╚████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══██║██╔══██╗██║     ██╔══██║  ╚██╔╝
╚██████╔╝██║  ██║██║ ╚═╝ ██║██║  ██║██║  ██║╚██████╗██║  ██║   ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝   ╚═╝'

echo -e "\n$ascii_art\n"

# Read and display version
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OMARCHY_VERSION="1.0.12"
if [[ -f "$SCRIPT_DIR/VERSION" ]]; then
    OMARCHY_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "1.0.12")
else
    # Fetch version from GitHub when running via curl
    OMARCHY_VERSION=$(curl -fsSL https://raw.githubusercontent.com/CyphrRiot/OhmArchy/master/VERSION 2>/dev/null || echo "1.0.12")
fi

echo -e "🚀 OhmArchy Setup - Version: $OMARCHY_VERSION"
echo -e "================================================\n"

# Install git if missing
pacman -Q git &>/dev/null || sudo pacman -Sy --noconfirm --needed git

# Clone OhmArchy repository
echo -e "\nCloning OhmArchy..."
rm -rf ~/.local/share/omarchy/
git clone https://github.com/CyphrRiot/OhmArchy.git ~/.local/share/omarchy || {
    echo "Error: Failed to clone OhmArchy repository. Check your internet connection."
    exit 1
}

# Switch to custom branch if specified
if [[ -n "$OMARCHY_REF" ]]; then
    echo -e "\nUsing branch: $OMARCHY_REF"
    if cd ~/.local/share/omarchy &&
       git fetch origin "${OMARCHY_REF}" &&
       git checkout "${OMARCHY_REF}"; then
        echo "✓ Switched to branch $OMARCHY_REF"
    else
        echo "⚠ Failed to switch to branch $OMARCHY_REF, using default"
    fi
    cd - >/dev/null
fi

# Start installation
echo -e "\nOhmArchy installation starting..."
source ~/.local/share/omarchy/install.sh
