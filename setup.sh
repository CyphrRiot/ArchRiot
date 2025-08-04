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
if [[ -n "$LOCAL_VERSION" && "$LOCAL_VERSION" == "$ARCHRIOT_VERSION" && "$ARCHRIOT_VERSION" != "unknown" ]]; then
    echo -e "✨ ${PURPLE}ArchRiot v$ARCHRIOT_VERSION is already installed!${NC}"
    echo -e "📦 Your system is up to date - no upgrade needed."
    echo -e "🔄 To force reinstall:"
    echo -e "   rm -rf ~/.local/share/archriot && curl -fsSL https://ArchRiot.org/setup.sh | bash"
    echo -e "🆔 To check version: cat ~/.local/share/archriot/VERSION"
    echo -e "\nHave a great day! 🎉"
    exit 0
fi

# Install git if missing
pacman -Q git &>/dev/null || sudo pacman -Sy --noconfirm --needed git



# Smart ArchRiot installation/update
if [[ -d ~/.local/share/archriot/.git ]]; then
    echo -e "\n🔄 Updating existing ArchRiot installation..."
    cd ~/.local/share/archriot

    # Backup any local changes before destructive reset
    if ! git diff --quiet || ! git diff --cached --quiet; then
        # Load centralized backup system
        if [[ -f ~/.local/share/archriot/install/lib/backup-manager.sh ]]; then
            source ~/.local/share/archriot/install/lib/backup-manager.sh
            backup_archriot_install "Local changes detected during upgrade"
            echo "📦 Local changes backed up to: ~/.archriot/backups/"
        else
            # Fallback to old method if backup system not available
            backup_dir="$HOME/.local/share/archriot-backup-$(date +%Y%m%d-%H%M%S)"
            cp -r ~/.local/share/archriot "$backup_dir"
            echo "📦 Local changes backed up to: $backup_dir"
        fi
    fi

    # Safe update with fallback to fresh clone
    if git fetch origin && git reset --hard origin/master; then
        echo "✓ ArchRiot updated successfully"
    else
        echo "⚠ Update failed, performing fresh installation..."
        cd - >/dev/null
        rm -rf ~/.local/share/archriot
        git clone --depth 1 https://github.com/CyphrRiot/ArchRiot.git ~/.local/share/archriot || {
            echo "Error: Failed to clone ArchRiot repository. Check your internet connection."
            echo "🧹 Cleaning up partial download..."
            rm -rf ~/.local/share/archriot
            exit 1
        }

        # Verify clone contains required files
        if [[ ! -f ~/.local/share/archriot/install.sh ]]; then
            echo "❌ CRITICAL: Clone incomplete - install.sh missing"
            echo "   Repository may be corrupted or incomplete"
            rm -rf ~/.local/share/archriot
            exit 1
        fi

        echo "✓ Fresh installation completed"
    fi
    cd - >/dev/null
else
    echo -e "\n📥 Fresh ArchRiot installation..."
    rm -rf ~/.local/share/archriot  # Remove any non-git directory
    git clone --depth 1 https://github.com/CyphrRiot/ArchRiot.git ~/.local/share/archriot || {
        echo "Error: Failed to clone ArchRiot repository. Check your internet connection."
        echo "🧹 Cleaning up partial download..."
        rm -rf ~/.local/share/archriot
        exit 1
    }

    # Verify clone contains required files
    if [[ ! -f ~/.local/share/archriot/install.sh ]]; then
        echo "❌ CRITICAL: Clone incomplete - install.sh missing"
        echo "   Repository may be corrupted or incomplete"
        rm -rf ~/.local/share/archriot
        exit 1
    fi

    echo "✓ ArchRiot cloned successfully"
fi

# Switch to custom branch if specified
if [[ -n "$ARCHRIOT_REF" ]]; then
    echo -e "\nUsing branch: $ARCHRIOT_REF"
    if cd ~/.local/share/archriot &&
       git fetch origin "${ARCHRIOT_REF}" &&
       git checkout "${ARCHRIOT_REF}"; then
        echo "✓ Switched to branch $ARCHRIOT_REF"
    else
        echo "⚠ Failed to switch to branch $ARCHRIOT_REF, using default"
    fi
    cd - >/dev/null
fi

# Start installation
echo -e "\nArchRiot installation starting..."

# Check if pre-built installer exists
if [[ ! -f ~/.local/share/archriot/install/archriot ]]; then
    echo -e "❌ Pre-built installer not found. This may be a development version."
    echo -e "   Expected: ~/.local/share/archriot/install/archriot"
    exit 1
fi

# Make sure installer is executable
chmod +x ~/.local/share/archriot/install/archriot

echo -e "🚀 Running ArchRiot installer..."
~/.local/share/archriot/install/archriot
