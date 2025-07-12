#!/bin/bash
# Enhanced hyprlock wrapper to handle monitor changes and ensure reliable locking
# Created to fix issues with external monitor disconnection

# Log function for debugging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> /tmp/hyprlock-wrapper.log
}

log "Lock requested"

# Kill any existing hyprlock instances first
if pgrep -x "hyprlock" > /dev/null; then
    log "Killing existing hyprlock instance"
    killall -q hyprlock
    sleep 0.5
fi

# Reset Hyprland lock permission if needed
hyprctl keyword misc:allow_session_lock_restore 1
log "Reset lock permission"

# Get current monitor setup
connected_monitors=$(hyprctl monitors -j | jq length)
log "Detected $connected_monitors connected monitors"

# Choose the appropriate config based on monitor resolution
primary_res=$(hyprctl monitors -j | jq -r '.[0].width')
log "Primary monitor width: $primary_res"

if [ "$primary_res" -lt 1920 ]; then
    log "Using 1080p config"
    config_file="$HOME/.config/hypr/hyprlock-1080p.conf"
else
    log "Using standard config"
    config_file="$HOME/.config/hypr/hyprlock.conf"
fi

# Execute hyprlock with reliable options
log "Launching hyprlock with config: $config_file"
/usr/bin/hyprlock --config $config_file --immediate --no-fade-in

# Record the exit status
exit_status=$?
log "hyprlock exited with status: $exit_status"

# Handle potential failures - if hyprlock fails, try the manual recovery command
if [ $exit_status -ne 0 ]; then
    log "hyprlock failed, attempting recovery"
    sleep 1
    hyprctl --instance 0 'keyword misc:allow_session_lock_restore 1'
    hyprctl --instance 0 'dispatch exec hyprlock'
fi

# Log success
log "Lock completed"
exit 0