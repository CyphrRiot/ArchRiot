#!/usr/bin/env bash
# generate-keybindings-help.sh — Dynamic 2-column HTML from Hyprland keybinds
# Outputs to ~/.cache/archriot/help/keybindings.html and opens it
# Inspired by community contribution (Issue #33)

set -euo pipefail

# Defaults
DEFAULT_OUT="$HOME/.cache/archriot/help/keybindings.html"
CANDIDATES=(
  "$HOME/.config/hypr/keybindings.conf"
  "$HOME/.config/hypr/hyprland.conf"
)

CONFIG_FILE=""
OUT_FILE="$DEFAULT_OUT"
OPEN_RESULT=1

# --- Args --------------------------------------------------------------------
# --config PATH     Use a specific Hyprland keybindings file
# --output PATH     Write HTML to PATH (default: ~/.cache/archriot/help/keybindings.html)
# --no-open         Do not auto-open the generated file
while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)
      shift
      [[ $# -gt 0 ]] || { echo "Missing value for --config" >&2; exit 2; }
      CONFIG_FILE="$1"
      ;;
    --output)
      shift
      [[ $# -gt 0 ]] || { echo "Missing value for --output" >&2; exit 2; }
      OUT_FILE="$1"
      ;;
    --no-open)
      OPEN_RESULT=0
      ;;
    -h|--help)
      cat <<EOF
Usage: $(basename "$0") [--config PATH] [--output PATH] [--no-open]

Generate a dynamic keybindings help page from your Hyprland config and open it.

Options:
  --config PATH   Read binds from PATH (default: first of:
                    ~/.config/hypr/keybindings.conf or ~/.config/hypr/hyprland.conf)
  --output PATH   Write HTML to PATH (default: $DEFAULT_OUT)
  --no-open       Do not auto-open the generated result
EOF
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 2
      ;;
  esac
  shift
done

# --- Resolve config ----------------------------------------------------------
if [[ -z "$CONFIG_FILE" ]]; then
  for c in "${CANDIDATES[@]}"; do
    if [[ -f "$c" ]]; then
      CONFIG_FILE="$c"
      break
    fi
  done
fi

# Ensure an output directory exists
mkdir -p "$(dirname "$OUT_FILE")"

# If no config found, generate a helpful page and open it
if [[ -z "$CONFIG_FILE" || ! -f "$CONFIG_FILE" ]]; then
  cat > "$OUT_FILE" <<'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>ArchRiot — Keybindings Help (No Config Found)</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { --bg:#0f172a; --fg:#cbd5e1; --muted:#94a3b8; --accent:#7aa2f7; }
  html,body{background:var(--bg);color:var(--fg);margin:0;padding:0;font-family:ui-monospace,Menlo,Consolas,"Liberation Mono",monospace}
  .wrap{max-width:900px;margin:0 auto;padding:28px 20px}
  h1{color:var(--accent);margin:0 0 14px}
  p{color:var(--muted)}
  code{background:#11162a;color:#e5e7eb;padding:2px 6px;border-radius:6px;border:1px solid #22263a}
</style>
</head>
<body>
  <div class="wrap">
    <h1>ArchRiot — Keybindings Help</h1>
    <p>No Hyprland config was found. Checked:</p>
    <ul>
      <li>~/.config/hypr/keybindings.conf</li>
      <li>~/.config/hypr/hyprland.conf</li>
    </ul>
    <p>Create <code>~/.config/hypr/keybindings.conf</code> with lines like<br>
    <code>bind = $mod, H, exec, ... # Website help</code><br>
    or ensure binds with inline comments exist in <code>~/.config/hypr/hyprland.conf</code>.</p>
    <p>Then run this generator again to build a dynamic help page.</p>
    <p>Close: <code>SUPER+Q</code> or <code>SUPER+W</code></p>
  </div>
</body>
</html>
EOF

  if [[ $OPEN_RESULT -eq 1 ]] && command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$OUT_FILE" >/dev/null 2>&1 &
  fi
  echo "$OUT_FILE"
  exit 0
fi

# --- Helpers -----------------------------------------------------------------
html_escape() {
  # Escapes &, <, > for HTML safety
  # shellcheck disable=SC2001
  local s
  s="$(printf '%s' "$1" | sed -e 's/&/\&amp;/g' -e 's/</\&lt;/g' -e 's/>/\&gt;/g')"
  printf '%s' "$s"
}

format_keys() {
  # Given a "mod,key" string, map $mod → SUPER and replace commas with ' + '
  local s="$1"
  s="$(echo "$s" | sed -E 's/^\s*,\s*//')"          # remove leading comma if no modifier
  s="$(echo "$s" | sed -E 's/\s+//g')"              # strip spaces around commas
  s="$(echo "$s" | sed -E 's/\$mod/SUPER/g')"       # replace $mod with SUPER
  s="$(echo "$s" | sed 's/,/ + /g')"                # commas to visual plus
  printf '%s' "$s"
}

timestamp() { date +"%Y-%m-%d %H:%M:%S %Z"; }

# --- Generate HTML -----------------------------------------------------------
cat > "$OUT_FILE" <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>ArchRiot — Keybindings Help</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { --bg:#0f172a; --bg2:#1a1b26; --fg:#cbd5e1; --muted:#94a3b8; --accent:#7aa2f7; --accent2:#bb9af7; --line:#22263a; }
  html,body{background:var(--bg);color:var(--fg);margin:0;padding:0;font-family:ui-monospace,Menlo,Consolas,"Liberation Mono",monospace;line-height:1.35}
  .wrap{max-width:1100px;margin:0 auto;padding:28px 20px 40px}
  h1{color:var(--accent);font-size:28px;margin:0 0 14px;line-height:1.15}
  p.note{color:var(--muted);margin:8px 0 18px}
  table{width:100%;border-collapse:collapse;margin:12px 0}
  th,td{padding:8px 10px;border-bottom:1px solid var(--line);vertical-align:top}
  th{color:var(--accent2);text-align:left;font-weight:700}
  tr:hover{background:#12172a}
  .bind{color:#c7d2fe;white-space:pre-wrap}
  .desc{color:var(--muted)}
  .meta{color:var(--muted);font-size:12px;margin-top:6px}
  .footer{color:var(--muted);font-size:12px;margin-top:16px}
  code,kbd{background:#11162a;color:#e5e7eb;padding:2px 6px;border-radius:6px;border:1px solid var(--line);font-size:.95em}
</style>
</head>
<body>
  <div class="wrap">
    <h1>ArchRiot — Keybindings Help</h1>
    <p class="note">
      Generated from <code>$(html_escape "$CONFIG_FILE")</code> on <code>$(timestamp)</code>.<br>
      Tip: Add inline comments after binds to describe them. Example:
      <code>bind = \$mod, H, exec, ... # Website help</code>
    </p>
    <table>
      <thead><tr><th>Bind</th><th>Description</th></tr></thead>
      <tbody>
EOF

# Pull only 'bind' lines with an inline comment (#)
# Normalize 'bind=' to 'bind =', then extract first two comma-separated fields after 'bind ='
while IFS= read -r line; do
  # Skip lines without a comment (no description)
  [[ "$line" == *"#"* ]] || continue

  # Normalize "bind=" to "bind ="
  norm="$(echo "$line" | sed -E 's/^[[:space:]]*bind[[:space:]]*=/bind =/')" || true

  # Extract "mod,key" (first two fields)
  keys="$(echo "$norm" | sed -E 's/^bind[[:space:]]*=[[:space:]]*([^,]+,[^,]+).*/\1/')" || true
  keys="$(echo "$keys" | xargs)" # trim

  # Extract description after "#"
  desc="${line#*#}"
  desc="$(echo "$desc" | xargs)"

  # Format keys and escape HTML
  keys_fmt="$(format_keys "$keys")"
  keys_fmt_esc="$(html_escape "$keys_fmt")"
  desc_esc="$(html_escape "$desc")"

  printf '        <tr><td class="bind">%s</td><td class="desc">%s</td></tr>\n' \
    "$keys_fmt_esc" "$desc_esc" >> "$OUT_FILE"
done < <(grep -E '^[[:space:]]*bind[[:space:]]*=' "$CONFIG_FILE" || true)

# Close HTML
cat >> "$OUT_FILE" <<'EOF'
      </tbody>
    </table>
    <p class="footer">
      Close: <kbd>SUPER</kbd>+<kbd>Q</kbd> or <kbd>SUPER</kbd>+<kbd>W</kbd> •
      Regenerate any time to reflect config changes.
    </p>
  </div>
</body>
</html>
EOF

# Auto-open result (best-effort)
if [[ $OPEN_RESULT -eq 1 ]] && command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$OUT_FILE" >/dev/null 2>&1 &
fi

echo "$OUT_FILE"
exit 0
