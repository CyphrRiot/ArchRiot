#!/usr/bin/env bash
# suspend-if-undocked.sh
#
# Purpose:
#   Suspend the system on idle ONLY when undocked (i.e., no external display connected).
#   If any external connector (HDMI/DP/DisplayPort/USB-C/DVI/VGA) is connected, exit without suspending.
#
# Usage:
#   Intended to be called by hypridle (or any idle manager) as an on-timeout action.
#   Example hypridle block:
#       listener {
#           timeout = 1800
#           on-timeout = ~/.local/bin/suspend-if-undocked.sh
#       }
#
# Behavior:
#   - Detects docked state using kernel DRM connectors in /sys/class/drm/*/status
#   - Exits immediately (no suspend) if an external connector is "connected"
#   - Otherwise, runs "systemctl suspend"
#
# Environment overrides:
#   - DOCK_INHIBIT_EXTERNAL_REGEX  (default: ^(HDMI|DP|DisplayPort|DVI|VGA|USB-?C)-)
#   - DOCK_INHIBIT_INTERNAL_REGEX  (default: ^(eDP|LVDS|DSI))
#   - DOCK_SUSPEND_CMD             (default: systemctl suspend)
#   - DOCK_DEBUG=1                 (enable debug logs to stderr)
#   - DRY_RUN=1                    (show actions but do not suspend)
#
# Exit codes:
#   0 on success or when staying awake (docked).
#   Non-zero only on unexpected errors.

set -euo pipefail

EXTERNAL_REGEX="${DOCK_INHIBIT_EXTERNAL_REGEX:-^(HDMI|DP|DisplayPort|DVI|VGA|USB-?C)-}"
INTERNAL_REGEX="${DOCK_INHIBIT_INTERNAL_REGEX:-^(eDP|LVDS|DSI)}"
SUSPEND_CMD="${DOCK_SUSPEND_CMD:-systemctl suspend}"

log() {
  if [[ "${DOCK_DEBUG:-0}" == "1" ]]; then
    printf '[suspend-if-undocked] %s\n' "$*" >&2 || true
  fi
}

is_docked_drm() {
  local s conn name state
  for s in /sys/class/drm/*/status; do
    [[ -e "$s" ]] || continue
    state="$(cat "$s" 2>/dev/null || true)"
    conn="$(basename "$(dirname "$s")")"  # e.g., card0-HDMI-A-1
    name="${conn#*-}"                     # e.g., HDMI-A-1

    # Treat any connected external connector as "docked"
    if printf '%s' "$name" | grep -Eiq "$EXTERNAL_REGEX" && \
       ! printf '%s' "$name" | grep -Eiq "$INTERNAL_REGEX"; then
      log "Connector $name state=$state (external)"
      if [[ "$state" == "connected" ]]; then
        return 0
      fi
    else
      log "Connector $name state=$state (internal/ignored)"
    fi
  done
  return 1
}

is_docked_hyprctl() {
  # Detect external outputs via Hyprland (if available)
  if ! command -v hyprctl >/dev/null 2>&1; then
    return 1
  fi
  local out
  out="$(hyprctl monitors 2>/dev/null || true)"
  if [[ -z "$out" ]]; then
    return 1
  fi
  # Match common external connector names seen in Hyprland monitor listings
  if printf '%s\n' "$out" | grep -Eiq '\b(HDMI|DP-|DisplayPort|DVI|VGA|USB-?C)\b'; then
    return 0
  fi
  return 1
}

is_on_ac() {
  # Prefer explicit mains/USB-PD online indicators
  local ps type online
  for ps in /sys/class/power_supply/*; do
    [[ -d "$ps" ]] || continue
    if [[ -f "$ps/type" && -f "$ps/online" ]]; then
      type="$(cat "$ps/type" 2>/dev/null || true)"
      online="$(cat "$ps/online" 2>/dev/null || true)"
      if printf '%s' "$type" | grep -Eiq '^(Mains|USB|USB_PD|USB-C)$' && [[ "$online" == "1" ]]; then
        log "Power supply $(basename "$ps") type=$type online=$online"
        return 0
      fi
    fi
  done

  # Fallback: if any battery is not Discharging, consider on AC
  local b
  for b in /sys/class/power_supply/BAT*; do
    [[ -d "$b" ]] || continue
    if [[ -f "$b/status" ]]; then
      local status
      status="$(cat "$b/status" 2>/dev/null || true)"
      if [[ "$status" != "Discharging" && -n "$status" ]]; then
        log "Battery status=$status (treat as on AC)"
        return 0
      fi
    fi
  done
  return 1
}

main() {
  if is_docked_hyprctl || is_docked_drm || is_on_ac; then
    log "External display or AC power detected → staying awake (no suspend)."
    exit 0
  fi

  log "No external display connected → suspending."
  if [[ "${DRY_RUN:-0}" == "1" ]]; then
    log "(dry-run) Would execute: $SUSPEND_CMD"
    exit 0
  fi

  # Best-effort suspend; do not crash the caller if suspend command fails.
  if ! eval "$SUSPEND_CMD"; then
    log "Warning: Suspend command failed: $SUSPEND_CMD"
    exit 1
  fi

  exit 0
}

main "$@"
