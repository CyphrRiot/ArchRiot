#!/usr/bin/env bash
# ==============================================================================
# ArchRiot Boot Logo Generator
# ==============================================================================
# Converts ArchRiot ASCII art to Plymouth boot screen logo
# ==============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Always use logo.png for LUKS boot screen - simple and reliable

# Configuration
LOGO_WIDTH=650
LOGO_HEIGHT=150
BACKGROUND_COLOR='#191a25'  # Exact Plymouth background color
TEXT_COLOR='#c0caf5'        # Tokyo Night foreground
PLYMOUTH_THEME_DIR="/usr/share/plymouth/themes/archriot"
TEMP_LOGO="/tmp/archriot_logo_temp.png"

echo -e "${BLUE}🎨 ArchRiot Boot Logo Generator${NC}"
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

# Automated change detection for logo updates
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_LOGO="$SCRIPT_DIR/../config/images/logo.png"
LOGO_SCRIPT_PATH="$SCRIPT_DIR/../bin/generate-boot-logo.sh"
LOGO_CHECKSUM_FILE="$HOME/.config/archriot/logo_checksums.txt"

# Create config directory if needed
mkdir -p "$HOME/.config/archriot"

if [[ ! -f "$SOURCE_LOGO" ]]; then
    echo -e "${RED}❌ Logo not found at: $SOURCE_LOGO${NC}"
    exit 1
fi

# Check if logo regeneration is needed (automated detection)
needs_logo_update() {
    # Check if Plymouth theme directory exists
    if [[ ! -f "$PLYMOUTH_THEME_DIR/logo.png" ]]; then
        echo "Plymouth logo missing - generation needed"
        return 0
    fi

    # Check if this script changed
    local script_hash=$(md5sum "$LOGO_SCRIPT_PATH" | cut -d' ' -f1)
    local stored_script_hash=$(grep "script_hash" "$LOGO_CHECKSUM_FILE" 2>/dev/null | cut -d' ' -f1 || echo "none")

    if [[ "$script_hash" != "$stored_script_hash" ]]; then
        echo "Logo generator script changed - regeneration needed"
        return 0
    fi

    # Check if source logo changed
    local source_hash=$(md5sum "$SOURCE_LOGO" | cut -d' ' -f1)
    local stored_source_hash=$(grep "source_logo" "$LOGO_CHECKSUM_FILE" 2>/dev/null | cut -d' ' -f1 || echo "none")

    if [[ "$source_hash" != "$stored_source_hash" ]]; then
        echo "Source logo changed - regeneration needed"
        return 0
    fi

    # Check if logo configuration changed (size, colors, etc.)
    local config_hash=$(echo "${LOGO_WIDTH}${LOGO_HEIGHT}${BACKGROUND_COLOR}${TEXT_COLOR}" | md5sum | cut -d' ' -f1)
    local stored_config_hash=$(grep "logo_config" "$LOGO_CHECKSUM_FILE" 2>/dev/null | cut -d' ' -f1 || echo "none")

    if [[ "$config_hash" != "$stored_config_hash" ]]; then
        echo "Logo configuration changed - regeneration needed"
        return 0
    fi

    echo "Logo already up-to-date - skipping generation"
    return 1
}

if ! needs_logo_update; then
    echo -e "${GREEN}✅ Boot logo already up-to-date${NC}"
    exit 0
fi

echo -e "${YELLOW}Generating updated boot logo...${NC}"
echo -e "${GREEN}✓ Using logo: $SOURCE_LOGO${NC}"

# Copy and resize the logo for LUKS boot screen
magick "$SOURCE_LOGO" \
    -resize ${LOGO_WIDTH}x${LOGO_HEIGHT} \
    -background "$BACKGROUND_COLOR" \
    -gravity center \
    -extent ${LOGO_WIDTH}x${LOGO_HEIGHT} \
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
    if [ -d "$HOME/.local/share/archriot/config/default/plymouth" ]; then
        sudo cp -r "$HOME/.local/share/archriot/config/default/plymouth/"* "$PLYMOUTH_THEME_DIR/"
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
    sudo plymouth-set-default-theme archriot
    echo -e "${GREEN}✓ Plymouth theme set to archriot${NC}"

    # Regenerate initramfs to apply changes
    echo -e "${YELLOW}Regenerating initramfs (this may take a moment)...${NC}"

    # Change to safe directory before rebuilding initramfs to avoid getcwd issues
    ORIGINAL_DIR="$(pwd)"
    cd /tmp 2>/dev/null || cd /

    sudo mkinitcpio -P

    # Return to original directory
    cd "$ORIGINAL_DIR" 2>/dev/null || true

    echo -e "${GREEN}✓ Initramfs regenerated${NC}"
else
    echo -e "${YELLOW}⚠ Plymouth not available, logo installed but theme not activated${NC}"
fi

# Update logo tracking checksums after successful generation
echo "📝 Updating logo tracking files..."
script_hash=$(md5sum "$LOGO_SCRIPT_PATH" | cut -d' ' -f1)
source_hash=$(md5sum "$SOURCE_LOGO" | cut -d' ' -f1)
config_hash=$(echo "${LOGO_WIDTH}${LOGO_HEIGHT}${BACKGROUND_COLOR}${TEXT_COLOR}" | md5sum | cut -d' ' -f1)

echo "$script_hash script_hash" > "$LOGO_CHECKSUM_FILE"
echo "$source_hash source_logo" >> "$LOGO_CHECKSUM_FILE"
echo "$config_hash logo_config" >> "$LOGO_CHECKSUM_FILE"
echo "✓ Script hash recorded: $script_hash"
echo "✓ Source logo hash recorded: $source_hash"
echo "✓ Config hash recorded: $config_hash"

# Create marker file to indicate custom logo is installed
echo "custom_ascii_logo_installed=$(date)" | sudo tee "$PLYMOUTH_THEME_DIR/.custom_logo_marker" > /dev/null
sudo chmod 644 "$PLYMOUTH_THEME_DIR/.custom_logo_marker"

# Create persistent backup for re-installations
PERSISTENT_BACKUP_DIR="$HOME/.config/archriot/plymouth-backup"
mkdir -p "$PERSISTENT_BACKUP_DIR"
cp "$PLYMOUTH_THEME_DIR/logo.png" "$PERSISTENT_BACKUP_DIR/custom_logo.png"
echo "custom_logo_backup_created=$(date)" > "$PERSISTENT_BACKUP_DIR/backup_info.txt"

# Cleanup
rm -f "$TEMP_LOGO"

echo -e "${GREEN}🎉 ArchRiot boot logo generation complete!${NC}"
echo ""
echo -e "${BLUE}The ArchRiot logo has been installed for your LUKS/boot screen.${NC}"
echo -e "${BLUE}You'll see it on next reboot during disk decryption.${NC}"
echo ""
echo -e "${GREEN}✓ Custom logo backup created at: $PERSISTENT_BACKUP_DIR/custom_logo.png${NC}"
echo -e "${GREEN}✓ This logo will be preserved during ArchRiot re-installations${NC}"
echo ""
echo -e "${YELLOW}To test the Plymouth theme without rebooting:${NC}"
echo -e "${YELLOW}  sudo plymouthd --debug --debug-file=/tmp/plymouth.log${NC}"
echo -e "${YELLOW}  sudo plymouth --show-splash${NC}"
echo -e "${YELLOW}  # Press Ctrl+Alt+F2 to see it, then:${NC}"
echo -e "${YELLOW}  sudo plymouth --quit${NC}"
echo ""
echo -e "${BLUE}NOTE: Your custom ASCII logo will survive ArchRiot re-installations!${NC}"
echo -e "${GREEN}🚀 Next run will skip logo generation (unless script or source image changes)${NC}"
