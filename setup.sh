# ArchRiot theme deep purple color
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

ascii_art=' █████╗ ██████╗  ██████╗██╗  ██╗██████╗ ██╗ ██████╗ ████████╗
██╔══██╗██╔══██╗██╔════╝██║  ██║██╔══██╗██║██╔═══██╗╚══██╔══╝
███████║██████╔╝██║     ███████║██████╔╝██║██║   ██║   ██║
██╔══██║██╔══██╗██║     ██╔══██║██╔══██╗██║██║   ██║   ██║
██║  ██║██║  ██║╚██████╗██║  ██║██║  ██║██║╚██████╔╝   ██║
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝ ╚═════╝    ╚═╝   '

echo -e "\n${PURPLE}$ascii_art${NC}\n"

# Read and display version
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Read version from VERSION file (single source of truth)
if [[ -f "$SCRIPT_DIR/VERSION" ]]; then
    ARCHRIOT_VERSION=$(cat "$SCRIPT_DIR/VERSION" 2>/dev/null || echo "unknown")
else
    # Fetch version from GitHub when running via curl - WITH CACHE BUSTING
    # CRITICAL: Never use cached content for installer scripts!
    # CDN caching causes stale versions and prevents immediate upgrades
    CACHE_BUSTER=$(date +%s)
    ARCHRIOT_VERSION=$(curl -fsSL -H "Cache-Control: no-cache" -H "Pragma: no-cache" "https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/VERSION?v=$CACHE_BUSTER" 2>/dev/null || echo "unknown")
fi

echo -e "🎭 ArchRiot Setup - Version: $ARCHRIOT_VERSION"
echo -e "================================================\n"

# Check if ArchRiot is already installed with the same version
LOCAL_VERSION=""
if [[ -f "$HOME/.local/share/archriot/VERSION" ]]; then
    LOCAL_VERSION=$(cat "$HOME/.local/share/archriot/VERSION" 2>/dev/null || echo "")
fi

# Compare versions and exit if they match
if [[ -n "$LOCAL_VERSION" && "$LOCAL_VERSION" == "$ARCHRIOT_VERSION" ]]; then
    echo -e "✨ ${PURPLE}ArchRiot v$ARCHRIOT_VERSION is already installed!${NC}"
    echo -e "📦 Your system is up to date - no upgrade needed."
    echo -e "🔄 To force reinstall: rm -rf ~/.local/share/archriot && curl -fsSL https://ArchRiot.org/setup.sh | bash"
    echo -e "🆔 To check version: cat ~/.local/share/archriot/VERSION"
    echo -e "\nHave a great day! 🎉"
    exit 0
fi

# Install git if missing
pacman -Q git &>/dev/null || sudo pacman -Sy --noconfirm --needed git

# NUCLEAR cleanup of ALL ArchRiot/OhmArchy directories before fresh install
echo -e "\n🧹 NUCLEAR cleanup of old installations..."
echo "   Removing ALL ArchRiot and OhmArchy directories..."
rm -rf ~/.local/share/archriot/
rm -rf ~/.local/share/omarchy/

# Force fresh download with cache busting
echo -e "📥 Downloading ArchRiot with cache-busting headers..."
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
