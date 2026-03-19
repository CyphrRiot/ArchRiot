#!/usr/bin/env bash
# Refresh crypto cache for hyprlock - run this with a keybinding
rm -f ~/.cache/hyprlock-crypto.json ~/.cache/hyprlock-crypto-prev.json
~/.local/share/archriot/install/archriot --crypto ROWML > /dev/null 2>&1
