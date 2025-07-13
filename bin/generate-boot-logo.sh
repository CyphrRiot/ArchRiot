#!/usr/bin/env bash
# ==============================================================================
# OhmArchy Boot Logo Generator
# ==============================================================================
# Converts OhmArchy ASCII art to Plymouth boot screen logo
# ==============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ASCII art from setup.sh
ASCII_ART=' ██████╗ ██╗  ██╗███╗   ███╗ █████╗ ██████╗  ██████╗██╗  ██╗██╗   ██╗
██╔═══██╗██║  ██║████╗ ████║██╔══██╗██╔══██╗██╔════╝██║  ██║╚██╗ ██╔╝
██║   ██║███████║██╔████╔██║███████║██████╔╝██║     ███████║ ╚████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══██║██╔══██╗██║     ██╔══██║  ╚██╔╝
╚██████╔╝██║  ██║██║ ╚═╝ ██║██║  ██║██║  ██║╚██████╗██║  ██║   ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝   ╚═╝'

# Configuration
LOGO_WIDTH=800
LOGO_HEIGHT=168
BACKGROUND_COLOR='#1a1b26'  # Tokyo Night background
TEXT_COLOR='#c0caf5'        # Tokyo Night foreground
PLYMOUTH_THEME_DIR="/usr/share/plymouth/themes/omarchy"
TEMP_LOGO="/tmp/ohmarchy_logo_temp.png"

echo -e "${BLUE}🎨 OhmArchy Boot Logo Generator${NC}"
echo -e "${BLUE}=================================${NC}"

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"
if ! command -v convert &> /dev/null; then
    echo -e "${RED}❌ ImageMagick not found. Installing...${NC}"
    if command -v yay &> /dev/null; then
        yay -S --noconfirm --needed imagemagick
    elif command -v pacman &> /dev/null; then
        sudo pacman -S --noconfirm --needed imagemagick
    else
        echo -e "${RED}❌ Cannot install ImageMagick. Please install manually.${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✓ ImageMagick available${NC}"

# Create temporary ASCII art file
echo -e "${YELLOW}Creating ASCII art image...${NC}"
echo "$ASCII_ART" > /tmp/ascii_art.txt

# Try different fonts in order of preference
FONTS=(
    "Cascadia-Mono"
    "CaskaydiaCove-Nerd-Font-Mono"
    "CaskaydiaMono-Nerd-Font"
    "JetBrainsMono-Nerd-Font"
    "FiraCode-Nerd-Font-Mono"
    "DejaVu-Sans-Mono"
    "Liberation-Mono"
    "monospace"
)

FONT_FOUND=""
for font in "${FONTS[@]}"; do
    if convert -list font | grep -qi "$font"; then
        FONT_FOUND="$font"
        echo -e "${GREEN}✓ Using font: $font${NC}"
        break
    fi
done

if [ -z "$FONT_FOUND" ]; then
    FONT_FOUND="monospace"
    echo -e "${YELLOW}⚠ Using fallback font: monospace${NC}"
fi

# Generate the logo
echo -e "${YELLOW}Generating boot logo...${NC}"
convert -size ${LOGO_WIDTH}x${LOGO_HEIGHT} \
    -background "$BACKGROUND_COLOR" \
    -fill "$TEXT_COLOR" \
    -font "$FONT_FOUND" \
    -pointsize 10 \
    -gravity center \
    label:@/tmp/ascii_art.txt \
    "$TEMP_LOGO"

if [ ! -f "$TEMP_LOGO" ]; then
    echo -e "${RED}❌ Failed to generate logo image${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Logo generated successfully${NC}"

# Check if Plymouth theme directory exists
if [ ! -d "$PLYMOUTH_THEME_DIR" ]; then
    echo -e "${YELLOW}⚠ Plymouth theme directory not found. Creating...${NC}"
    sudo mkdir -p "$PLYMOUTH_THEME_DIR"

    # Copy the entire Plymouth theme if it doesn't exist
    if [ -d "$HOME/.local/share/omarchy/default/plymouth" ]; then
        sudo cp -r "$HOME/.local/share/omarchy/default/plymouth/"* "$PLYMOUTH_THEME_DIR/"
        echo -e "${GREEN}✓ Plymouth theme installed${NC}"
    else
        echo -e "${RED}❌ Plymouth theme source not found${NC}"
        exit 1
    fi
fi

# Backup existing logo
if [ -f "$PLYMOUTH_THEME_DIR/logo.png" ]; then
    BACKUP_FILE="$PLYMOUTH_THEME_DIR/logo.png.backup.$(date +%Y%m%d%H%M%S)"
    sudo cp "$PLYMOUTH_THEME_DIR/logo.png" "$BACKUP_FILE"
    echo -e "${GREEN}✓ Backed up existing logo to: $(basename "$BACKUP_FILE")${NC}"
fi

# Install new logo
echo -e "${YELLOW}Installing new boot logo...${NC}"
sudo cp "$TEMP_LOGO" "$PLYMOUTH_THEME_DIR/logo.png"
sudo chmod 644 "$PLYMOUTH_THEME_DIR/logo.png"

# Set ownership
sudo chown root:root "$PLYMOUTH_THEME_DIR/logo.png"

echo -e "${GREEN}✓ Boot logo installed${NC}"

# Update Plymouth theme
echo -e "${YELLOW}Updating Plymouth configuration...${NC}"
if command -v plymouth-set-default-theme &> /dev/null; then
    sudo plymouth-set-default-theme omarchy
    echo -e "${GREEN}✓ Plymouth theme set to omarchy${NC}"

    # Regenerate initramfs to apply changes
    echo -e "${YELLOW}Regenerating initramfs (this may take a moment)...${NC}"
    sudo mkinitcpio -P
    echo -e "${GREEN}✓ Initramfs regenerated${NC}"
else
    echo -e "${YELLOW}⚠ Plymouth not available, logo installed but theme not activated${NC}"
fi

# Create marker file to indicate custom logo is installed
echo "custom_ascii_logo_installed=$(date)" | sudo tee "$PLYMOUTH_THEME_DIR/.custom_logo_marker" > /dev/null
sudo chmod 644 "$PLYMOUTH_THEME_DIR/.custom_logo_marker"

# Create persistent backup for re-installations
PERSISTENT_BACKUP_DIR="$HOME/.config/omarchy/plymouth-backup"
mkdir -p "$PERSISTENT_BACKUP_DIR"
cp "$PLYMOUTH_THEME_DIR/logo.png" "$PERSISTENT_BACKUP_DIR/custom_logo.png"
echo "custom_logo_backup_created=$(date)" > "$PERSISTENT_BACKUP_DIR/backup_info.txt"

# Cleanup
rm -f /tmp/ascii_art.txt "$TEMP_LOGO"

echo -e "${GREEN}🎉 OhmArchy boot logo generation complete!${NC}"
echo ""
echo -e "${BLUE}The ASCII art logo has been installed for your LUKS/boot screen.${NC}"
echo -e "${BLUE}You'll see it on next reboot during disk decryption.${NC}"
echo ""
echo -e "${GREEN}✓ Custom logo backup created at: $PERSISTENT_BACKUP_DIR/custom_logo.png${NC}"
echo -e "${GREEN}✓ This logo will be preserved during OhmArchy re-installations${NC}"
echo ""
echo -e "${YELLOW}To test the Plymouth theme without rebooting:${NC}"
echo -e "${YELLOW}  sudo plymouthd --debug --debug-file=/tmp/plymouth.log${NC}"
echo -e "${YELLOW}  sudo plymouth --show-splash${NC}"
echo -e "${YELLOW}  # Press Ctrl+Alt+F2 to see it, then:${NC}"
echo -e "${YELLOW}  sudo plymouth --quit${NC}"
echo ""
echo -e "${BLUE}NOTE: Your custom ASCII logo will survive OhmArchy re-installations!${NC}"
