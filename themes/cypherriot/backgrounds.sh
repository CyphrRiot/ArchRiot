#!/bin/bash
# CypherRiot theme backgrounds
# Copy existing wallpapers from repo (NO DOWNLOADS NEEDED)

# Define backgrounds directory
BACKGROUNDS_DIR="$HOME/.config/archriot/backgrounds"

mkdir -p "$BACKGROUNDS_DIR/cypherriot"

# Clear existing backgrounds to prevent duplicates
rm -f "$BACKGROUNDS_DIR/cypherriot"/*.jpg "$BACKGROUNDS_DIR/cypherriot"/*.png 2>/dev/null || true

# Copy City Rainy Night wallpaper as default (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/City-Rainy-Night.png "$BACKGROUNDS_DIR/cypherriot/1-City-Rainy-Night.png"

# Copy tokyo_cat wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/tokyo_cat.jpeg "$BACKGROUNDS_DIR/cypherriot/2-tokyo_cat.jpeg"

# Copy escape velocity wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/escape_velocity.jpg "$BACKGROUNDS_DIR/cypherriot/3-escape_velocity.jpg"

# Copy blue night moon over lake wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/blue_night_moon_over_lake.jpg "$BACKGROUNDS_DIR/cypherriot/4-blue_night_moon_over_lake.jpg"

# Copy Staircase wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/Staircase.png "$BACKGROUNDS_DIR/cypherriot/5-Staircase.png"

# Copy Anime Purple Eyes wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/Anime-Purple-eyes.png "$BACKGROUNDS_DIR/cypherriot/6-Anime-Purple-eyes.png"

# Copy purple sky wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/purple_sky.jpeg "$BACKGROUNDS_DIR/cypherriot/7-purple_sky.jpeg"

# Copy cyber wallpaper (already in repo)
cp ~/.local/share/archriot/themes/cypherriot/backgrounds/cyber.jpg "$BACKGROUNDS_DIR/cypherriot/8-cyber.jpg"
