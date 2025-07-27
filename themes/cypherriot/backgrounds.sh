#!/bin/bash
# CypherRiot theme backgrounds
# Automatically discover and copy all background files from repo
# Prevents duplicates and handles existing files intelligently

# Define backgrounds directory
BACKGROUNDS_DIR="$HOME/.config/archriot/backgrounds"

mkdir -p "$BACKGROUNDS_DIR/cypherriot"

# Source directory for backgrounds
SOURCE_DIR="$HOME/.local/share/archriot/themes/cypherriot/backgrounds"

# Check if source directory exists
if [[ ! -d "$SOURCE_DIR" ]]; then
    echo "❌ Source backgrounds directory not found: $SOURCE_DIR"
    exit 1
fi

echo "🖼️  Discovering background files..."

# Find all image files and sort them
mapfile -t ALL_BACKGROUNDS < <(find "$SOURCE_DIR" -maxdepth 1 -type f \( -iname "*.png" -o -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.webp" \) | sort)

if [[ ${#ALL_BACKGROUNDS[@]} -eq 0 ]]; then
    echo "❌ No background files found in $SOURCE_DIR"
    exit 1
fi

echo "📋 Found ${#ALL_BACKGROUNDS[@]} background files"

# Clear existing numbered backgrounds to prevent duplicates/conflicts
echo "🗑️  Cleaning existing numbered backgrounds..."
find "$BACKGROUNDS_DIR/cypherriot" -name "[0-9][0-9]-*" -type f -delete 2>/dev/null || true

# Separate riot_zero.png from other backgrounds for priority ordering
RIOT_ZERO=""
OTHER_BACKGROUNDS=()

for bg in "${ALL_BACKGROUNDS[@]}"; do
    filename=$(basename "$bg")
    if [[ "$filename" == "riot_zero.png" ]]; then
        RIOT_ZERO="$bg"
    else
        OTHER_BACKGROUNDS+=("$bg")
    fi
done

# Copy riot_zero.png as #1 if it exists (this is the new default)
counter=1
DEFAULT_BG=""

if [[ -n "$RIOT_ZERO" ]]; then
    filename=$(basename "$RIOT_ZERO")
    dest_file="$BACKGROUNDS_DIR/cypherriot/$(printf "%02d" $counter)-$filename"

    # Only copy if source is different from destination (avoid copying file to itself)
    if [[ "$RIOT_ZERO" != "$dest_file" ]]; then
        cp "$RIOT_ZERO" "$dest_file"
        echo "✓ Copied default: $(printf "%02d" $counter)-$filename"
        DEFAULT_BG="$dest_file"
    else
        echo "✓ Default already in place: $(printf "%02d" $counter)-$filename"
        DEFAULT_BG="$dest_file"
    fi
    ((counter++))
fi

# Copy all other backgrounds
for bg in "${OTHER_BACKGROUNDS[@]}"; do
    if [[ -f "$bg" ]]; then
        filename=$(basename "$bg")
        dest_file="$BACKGROUNDS_DIR/cypherriot/$(printf "%02d" $counter)-$filename"

        # Only copy if source is different from destination
        if [[ "$bg" != "$dest_file" ]]; then
            cp "$bg" "$dest_file"
            echo "✓ Copied: $(printf "%02d" $counter)-$filename"
        else
            echo "✓ Already in place: $(printf "%02d" $counter)-$filename"
        fi

        # If no riot_zero was found, use the first background as default
        if [[ -z "$DEFAULT_BG" ]]; then
            DEFAULT_BG="$dest_file"
        fi

        ((counter++))
    fi
done

echo ""
echo "🎉 Successfully processed $((counter-1)) backgrounds for CypherRiot theme"

# If this script is being run during installation, set the default background
if [[ -n "$DEFAULT_BG" && -d "$HOME/.config/archriot/current" ]]; then
    echo "🖼️  Setting default background..."
    ln -snf "$DEFAULT_BG" "$HOME/.config/archriot/current/background"
    echo "✓ Default background set: $(basename "$DEFAULT_BG")"
fi

echo ""
echo "💡 Available backgrounds (in order):"
for bg_file in "$BACKGROUNDS_DIR/cypherriot"/[0-9][0-9]-*; do
    if [[ -f "$bg_file" ]]; then
        echo "  • $(basename "$bg_file")"
    fi
done

echo ""
echo "🎮 Use SUPER+CTRL+SPACE to cycle through backgrounds"
