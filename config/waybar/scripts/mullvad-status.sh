#!/bin/bash
# Robust Mullvad status for Waybar (JSON)
# - Exits cleanly on SIGPIPE (Waybar reloads)
# - Suppresses stderr noise
# - Defensive parsing of mullvad status output

set -euo pipefail
trap 'exit 0' PIPE
exec 2>/dev/null

STATUS="$(mullvad status || true)"

connected=false
location=""

# Detect connection
if printf '%s' "$STATUS" | grep -qE '^[[:space:]]*Connected'; then
  connected=true
  # Try to extract location from Relay: line (format: Relay: us-sjc-wg-501)
  relay_line="$(printf '%s' "$STATUS" | grep -E '^\s*Relay:' | head -n1 || true)"
  if [ -n "$relay_line" ]; then
    # Get second field (us-sjc-wg-501) and take the middle token (sjc)
    second_field="$(printf '%s' "$relay_line" | awk '{print $2}' 2>/dev/null || true)"
    mid_token="$(printf '%s' "$second_field" | awk -F'-' '{print $2}' 2>/dev/null || true)"
    if [ -n "$mid_token" ]; then
      location="$(printf '%s' "$mid_token" | tr '[:lower:]' '[:upper:]')"
    fi
  fi
fi

if [ "$connected" = true ]; then
  [ -n "$location" ] || location="VPN"
  printf '{"text":"󰌆 %s","class":"mullvad-connected","tooltip":"Mullvad VPN Connected: %s\\nRight-click to disconnect"}\n' "$location" "$location"
else
  printf '{"text":"󰌉","class":"mullvad-disconnected","tooltip":"Mullvad VPN Disconnected\\nLeft-click to connect"}\n'
fi
