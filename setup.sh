ascii_art=' █████╗ ██████╗  ██████╗██╗  ██╗██████╗ ██╗ ██████╗ ████████╗
██╔══██╗██╔══██╗██╔════╝██║  ██║██╔══██╗██║██╔═══██╗╚══██╔══╝
███████║██████╔╝██║     ███████║██████╔╝██║██║   ██║   ██║
██╔══██║██╔══██╗██║     ██╔══██║██╔══██╗██║██║   ██║   ██║
██║  ██║██║  ██║╚██████╗██║  ██║██║  ██║██║╚██████╔╝   ██║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝ ╚═════╝    ╚═╝   '

echo -e "\n$ascii_art\n"

# Read and display version
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Read version from VERSION file (single source of truth)
if [[ -f "$SCRIPT_DIR/VERSION" ]]; then
    ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
else
    # Fetch version from GitHub when running via curl
    ARCHRIOT_VERSION=$(curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION 2>/dev/null || echo "unknown")
fi

echo -e "🎭 ArchRiot Setup - Version: $ARCHRIOT_VERSION"
echo -e "================================================\n"

# Install git if missing
pacman -Q git &>/dev/null || sudo pacman -Sy --noconfirm --needed git

# NUCLEAR cleanup of ALL ArchRiot/OhmArchy directories before fresh install
echo -e "\n🧹 NUCLEAR cleanup of old installations..."
echo "   Removing ALL ArchRiot and OhmArchy directories..."
rm -rf ~/.local/share/archriot/
rm -rf ~/.local/share/omarchy/
rm -rf ~/.config/archriot/
rm -rf ~/.config/omarchy/
echo "✓ Cleanup complete - fresh install guaranteed"

# Clone ArchRiot repository
echo -e "\nCloning ArchRiot..."
git clone https://github.com/CyphrRiot/ArchRiot.git ~/.local/share/archriot || {
    echo "Error: Failed to clone ArchRiot repository. Check your internet connection."
    exit 1
}

# Switch to custom branch if specified
if [[ -n "$OMARCHY_REF" ]]; then
    echo -e "\nUsing branch: $OMARCHY_REF"
    if cd ~/.local/share/archriot &&
       git fetch origin "${OMARCHY_REF}" &&
       git checkout "${OMARCHY_REF}"; then
        echo "✓ Switched to branch $OMARCHY_REF"
    else
        echo "⚠ Failed to switch to branch $OMARCHY_REF, using default"
    fi
    cd - >/dev/null
fi

# Start installation
echo -e "\nArchRiot installation starting..."
source ~/.local/share/archriot/install.sh
