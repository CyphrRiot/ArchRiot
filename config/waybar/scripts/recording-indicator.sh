#!/usr/bin/env sh
# ArchRiot - Waybar Kooha recording indicator
# Lightweight screencast detector for Waybar.
# - Shows a red dot when a PipeWire screencast is active (typically via xdg-desktop-portal-wlr).
# - On click, attempts to stop Kooha gracefully.
#
# Usage in Waybar (custom module example):
# {
#   "custom/rec-dot": {
#     "format": "{}",
#     "return-type": "json",
#     "interval": 2,
#     "exec": "$HOME/.config/waybar/scripts/recording-indicator.sh",
#     "on-click": "$HOME/.config/waybar/scripts/recording-indicator.sh --click",
#     "tooltip": true
#   }
# }
#
# CSS (style.css):
# .recording { color: #ff4d4d; }

set -eu

# Attempt to gracefully stop Kooha when clicked
if [ "${1:-}" = "--click" ]; then
  # Try SIGINT first (let Kooha finalize the file cleanly), then TERM as fallback
  pkill -INT -x kooha 2>/dev/null || pkill -TERM -x kooha 2>/dev/null || true
  exit 0
fi

active=0

# Detect active screencast via PipeWire introspection.
# Prefer pw-dump (PipeWire >= 0.3.61 typically has -N), fallback to pw-cli.
if command -v pw-dump >/dev/null 2>&1; then
  # -N prints only Nodes; look for common screencast markers:
  #   - node.name starting with xdpw_ (xdg-desktop-portal-wlr)
  #   - application.name containing xdg-desktop-portal
  #   - generic "screencast" hints present in node metadata
  if pw-dump -N 2>/dev/null | grep -E -iq '"node.name":"xdpw_|"application.name":"xdg-desktop-portal|screencast|screen-?cast'; then
    active=1
  fi
elif command -v pw-cli >/dev/null 2>&1; then
  # Older systems: list Nodes and search similar hints
  if pw-cli ls Node 2>/dev/null | grep -E -iq 'node.name = "xdpw_|application.name = "xdg-desktop-portal|screencast|screen-?cast'; then
    active=1
  fi
fi

# Best-effort fallback: if Kooha is running, assume recording.
# (Kooha typically exits when not recording; if it stays resident, this may show false positive.)
if [ "$active" -eq 0 ] && command -v pgrep >/dev/null 2>&1; then
  if pgrep -x kooha >/dev/null 2>&1; then
    active=1
  fi
fi

# Output Waybar JSON
if [ "$active" -eq 1 ]; then
  printf '{"text":"●","class":"recording","tooltip":"Screen recording (click to stop)"}\n'
else
  printf '{"text":"","class":"","tooltip":"No screen recording"}\n'
fi
