#!/usr/bin/env bash

# Hypridle toggle and status script for waybar

case "$1" in
  "status")
    if pgrep -x hypridle > /dev/null; then
      echo '{"text": "󱫗", "tooltip": "Idle lock is ON - Click to disable", "class": "active"}'
    else
      echo '{"text": "󱫗", "tooltip": "Idle lock is OFF - Click to enable", "class": "notactive"}'
    fi
    ;;
  "toggle"|*)
    if pgrep -x hypridle > /dev/null; then
      pkill -x hypridle
      notify-send "Stop locking computer when idle"
    else
      setsid hypridle &> /dev/null &
      notify-send "Now locking computer when idle"
    fi
    ;;
esac
