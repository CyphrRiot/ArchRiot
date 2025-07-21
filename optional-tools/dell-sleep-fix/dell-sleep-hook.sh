#!/bin/bash
# Dell XPS Sleep Hook Script
# Prevents crashes when Dell XPS with Intel Arc Graphics goes to sleep
# Place in /lib/systemd/system-sleep/

case $1/$2 in
  pre/suspend)
    echo "$(date): Preparing Dell XPS for sleep..." >> /var/log/dell-sleep.log

    # Sync filesystem to prevent data loss
    sync

    # Disable USB autosuspend temporarily
    echo on > /sys/bus/usb/devices/*/power/control 2>/dev/null || true

    # Force Intel Arc graphics to idle state
    echo auto > /sys/class/drm/card*/device/power/control 2>/dev/null || true

    # Disable problematic devices that can prevent proper sleep
    echo 0 > /sys/bus/pci/devices/*/power/wakeup 2>/dev/null || true

    # Ensure all processes are ready for sleep
    killall -STOP firefox 2>/dev/null || true
    killall -STOP chrome 2>/dev/null || true
    killall -STOP discord 2>/dev/null || true

    # Wait a moment for processes to settle
    sleep 1
    ;;

  post/suspend)
    echo "$(date): Dell XPS waking from sleep..." >> /var/log/dell-sleep.log

    # Resume stopped processes
    killall -CONT firefox 2>/dev/null || true
    killall -CONT chrome 2>/dev/null || true
    killall -CONT discord 2>/dev/null || true

    # Re-enable USB autosuspend
    echo auto > /sys/bus/usb/devices/*/power/control 2>/dev/null || true

    # Force Intel Arc graphics to active state
    echo on > /sys/class/drm/card*/device/power/control 2>/dev/null || true

    # Reload XE driver if it appears to be stuck
    if ! lsmod | grep -q xe; then
        modprobe xe 2>/dev/null || true
    fi

    # Restart display manager if needed (last resort)
    # systemctl restart display-manager 2>/dev/null || true
    ;;
esac

exit 0
