#!/bin/bash
# generate_keybindings_help - Simple 2-column HTML help page from Hyprland keybindings
# Tarso Galvão Sat Oct 11 10:56:11 PM -03 2025

DEFAULT_FILE="$HOME/.config/hypr/hyprland.conf"
USER_FILE="$HOME/.config/hypr/keybindings.conf"
OUTPUT_FILE="$HOME/.cache/archriot/help.html"
MOD_SYMBOL="⌘"

# Ensure default config exists
if [[ ! -f "$DEFAULT_FILE" ]]; then
    notify-send "Keymap Generator" "hyprland.conf not found, can't continue!"
    echo "File not found: $DEFAULT_FILE, can't continue."
    exit 1
fi

# --- Collect user unbinds (if any) ---
declare -A UNBINDS
if [[ -f "$USER_FILE" ]]; then
    while IFS= read -r line; do
        [[ "$line" =~ ^unbind ]] || continue
        keypair=$(echo "$line" | sed -E 's/^unbind *= *//;s/ *#.*//' | xargs)
        [[ "$keypair" =~ ^, ]] && keypair=$(echo "$keypair" | sed -E 's/^, *//')
        keypair=$(echo "$keypair" | sed 's/,/ + /g')
        keypair=$(echo "$keypair" | sed "s/\\\$mod/${MOD_SYMBOL}/g")
        UNBINDS["$keypair"]=1
    done < "$USER_FILE"
fi

# Write HTML header
cat > "$OUTPUT_FILE" <<'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>ArchRiot Keybind Mapping</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
html, body {
    font-family: monospace;
    background: #111;
    color: #eee;
    margin: 0;
    padding: 1em;
    overflow: auto;
    scrollbar-width: none;
    -ms-overflow-style: none;
}
h1 {
    text-align: center;
    color: #00e0ff;
    padding: 0;
}
table {
    width: 100%;
    border-collapse: collapse;
    margin: 1em auto;
    max-width: 600px;
}
th, td {
    padding: 8px 10px;
    border-bottom: 1px solid #333;
}
th {
    background: #222;
    color: #0ff;
}
tr:hover {
    background: #222;
}
.desc {
    color: #aaa;
}
</style>
</head>
<body>
<h1>ArchRiot Keybind Mapping</h1>
<table>
<thead><tr><th>Bind</th><th>Description</th></tr></thead>
<tbody>
EOF

# --- Parse DEFAULT_FILE ---
if [[ -f "$DEFAULT_FILE" ]]; then
    grep -E '^bind' "$DEFAULT_FILE" | while IFS= read -r line; do
        [[ "$line" != *"#"* ]] && continue
        desc="${line#*#}"
        desc=$(echo "$desc" | xargs)
        bind_keys=$(echo "$line" | sed -E 's/^bind *= *([^,]+,[^,]+).*/\1/' | xargs)
        if [[ "$bind_keys" =~ ^, ]]; then
            bind_keys=$(echo "$bind_keys" | sed -E 's/^, *//')
        fi
        bind_keys=$(echo "$bind_keys" | sed 's/,/ + /g')
        bind_keys=$(echo "$bind_keys" | sed "s/\\\$mod/${MOD_SYMBOL}/g")

        # Skip if this key was unbound by user
        if [[ -n "${UNBINDS[$bind_keys]}" ]]; then
            continue
        fi

        echo "<tr><td>${bind_keys}</td><td class=\"desc\">${desc}</td></tr>" >> "$OUTPUT_FILE"
    done
fi

# --- Parse USER_FILE ---
if [[ -f "$USER_FILE" ]]; then
    grep -E '^bind' "$USER_FILE" | while IFS= read -r line; do
        [[ "$line" != *"#"* ]] && continue
        desc="${line#*#}"
        desc=$(echo "$desc" | xargs)
        bind_keys=$(echo "$line" | sed -E 's/^bind *= *([^,]+,[^,]+).*/\1/' | xargs)
        if [[ "$bind_keys" =~ ^, ]]; then
            bind_keys=$(echo "$bind_keys" | sed -E 's/^, *//')
        fi
        bind_keys=$(echo "$bind_keys" | sed 's/,/ + /g')
        bind_keys=$(echo "$bind_keys" | sed "s/\\\$mod/${MOD_SYMBOL}/g")
        echo "<tr><td>${bind_keys}</td><td class=\"desc\">${desc}</td></tr>" >> "$OUTPUT_FILE"
    done
fi

# --- Close HTML ---
cat >> "$OUTPUT_FILE" <<'EOF'
</tbody>
</table>
</body>
</html>
EOF

echo "Help page generated: $(realpath "$OUTPUT_FILE")"
