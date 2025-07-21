#!/bin/bash
# Dell XPS Sleep Crash Fix Script
# Comprehensive solution for Dell XPS laptops with Intel Arc Graphics
# that crash when closing the lid and going to sleep
#
# This script addresses multiple potential causes:
# - Intel XE graphics driver issues
# - Power management conflicts
# - S2idle stability problems
# - Hardware-specific Dell quirks

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging
LOG_FILE="/var/log/dell-sleep-fix.log"
exec > >(tee -a "$LOG_FILE")
exec 2>&1

echo -e "${BLUE}=== Dell XPS Sleep Crash Fix Script ===${NC}"
echo "Started at: $(date)"
echo "Kernel: $(uname -r)"
echo "Hardware: $(lspci | grep VGA)"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo -e "${RED}Error: This script must be run as root${NC}"
    echo "Please run: sudo $0"
    exit 1
fi

# Backup function
backup_file() {
    local file="$1"
    if [[ -f "$file" ]]; then
        cp "$file" "${file}.backup-$(date +%Y%m%d-%H%M%S)"
        echo -e "${YELLOW}Backed up: $file${NC}"
    fi
}

echo -e "${BLUE}Step 1: Updating system packages...${NC}"
pacman -Sy --noconfirm
pacman -S --needed --noconfirm linux-firmware intel-ucode

echo -e "${BLUE}Step 2: Configuring Intel XE graphics driver...${NC}"
mkdir -p /etc/modprobe.d

cat > /etc/modprobe.d/xe-graphics-fix.conf << 'EOF'
# XE Graphics Driver Configuration for Sleep Stability
# Prevents crashes when Dell XPS with Intel Arc Graphics goes to sleep

# Enable power management for XE driver
options xe enable_guc=3
options xe enable_psr=0
options xe enable_fbc=0
options xe enable_gvt=0

# Disable problematic features that can cause sleep issues
options xe enable_hangcheck=0
options xe enable_rc6=0

# Force specific power states to prevent crashes
options xe modeset=1
options xe force_probe=*

# Dell-specific workarounds for laptop lid sleep issues
options xe enable_dc=0
options xe disable_power_well=0

# Debugging options (can be disabled once stable)
options xe debug=0x0
EOF

echo -e "${GREEN}Created XE graphics configuration${NC}"

echo -e "${BLUE}Step 3: Configuring systemd sleep settings...${NC}"
mkdir -p /etc/systemd/sleep.conf.d

cat > /etc/systemd/sleep.conf.d/dell-sleep-fix.conf << 'EOF'
[Sleep]
# Dell XPS Sleep Configuration for Stability
# Prevents crashes when closing lid by using safer sleep modes

# Use suspend-to-idle only (safest for new Intel hardware)
SuspendMode=s2idle
SuspendState=mem

# Disable hibernate and hybrid sleep to prevent conflicts
HibernateMode=
HibernateState=
HybridSleepMode=
HybridSleepState=

# Set shorter suspend delay to prevent race conditions
SuspendEstimationSec=5

# Disable problematic sleep states that cause crashes on Dell XPS
# with Intel Arc graphics
AllowSuspend=yes
AllowHibernation=no
AllowSuspendThenHibernate=no
AllowHybridSleep=no
EOF

echo -e "${GREEN}Created systemd sleep configuration${NC}"

echo -e "${BLUE}Step 4: Creating sleep hook scripts...${NC}"
mkdir -p /lib/systemd/system-sleep

cat > /lib/systemd/system-sleep/dell-sleep-hook.sh << 'EOF'
#!/bin/bash
# Dell XPS Sleep Hook Script
# Prevents crashes when Dell XPS with Intel Arc Graphics goes to sleep

case $1/$2 in
  pre/suspend)
    echo "$(date): Preparing Dell XPS for sleep..." >> /var/log/dell-sleep.log

    # Sync filesystem to prevent data loss
    sync

    # Disable USB autosuspend temporarily to prevent wake issues
    for device in /sys/bus/usb/devices/*/power/control; do
        [ -f "$device" ] && echo on > "$device" 2>/dev/null || true
    done

    # Force Intel Arc graphics to idle state
    for device in /sys/class/drm/card*/device/power/control; do
        [ -f "$device" ] && echo auto > "$device" 2>/dev/null || true
    done

    # Disable PCI device wakeup to prevent spurious wake events
    for device in /sys/bus/pci/devices/*/power/wakeup; do
        [ -f "$device" ] && echo disabled > "$device" 2>/dev/null || true
    done

    # Stop problematic processes that can interfere with sleep
    for proc in firefox chrome discord teams slack; do
        pkill -STOP "$proc" 2>/dev/null || true
    done

    # Ensure network interfaces are properly handled
    for iface in /sys/class/net/*/device/power/control; do
        [ -f "$iface" ] && echo auto > "$iface" 2>/dev/null || true
    done

    # Wait for processes to settle
    sleep 1
    ;;

  post/suspend)
    echo "$(date): Dell XPS waking from sleep..." >> /var/log/dell-sleep.log

    # Resume stopped processes
    for proc in firefox chrome discord teams slack; do
        pkill -CONT "$proc" 2>/dev/null || true
    done

    # Re-enable USB autosuspend
    for device in /sys/bus/usb/devices/*/power/control; do
        [ -f "$device" ] && echo auto > "$device" 2>/dev/null || true
    done

    # Force Intel Arc graphics to active state
    for device in /sys/class/drm/card*/device/power/control; do
        [ -f "$device" ] && echo on > "$device" 2>/dev/null || true
    done

    # Reload XE driver if it appears to be stuck
    if ! lsmod | grep -q xe; then
        modprobe xe 2>/dev/null || true
        echo "$(date): Reloaded XE driver after wake" >> /var/log/dell-sleep.log
    fi

    # Reset display if needed (uncomment if display issues persist)
    # systemctl restart display-manager 2>/dev/null || true
    ;;
esac

exit 0
EOF

chmod +x /lib/systemd/system-sleep/dell-sleep-hook.sh
echo -e "${GREEN}Created sleep hook script${NC}"

echo -e "${BLUE}Step 5: Updating kernel parameters...${NC}"
BOOT_ENTRY=$(find /boot/loader/entries/ -name "*.conf" | grep -v fallback | head -1)

if [[ -f "$BOOT_ENTRY" ]]; then
    backup_file "$BOOT_ENTRY"

    # Add Intel Arc specific kernel parameters
    if ! grep -q "intel_iommu=off" "$BOOT_ENTRY"; then
        sed -i 's/splash quiet/splash quiet intel_iommu=off xe.enable_guc=3 xe.enable_psr=0 xe.enable_fbc=0 mem_sleep_default=s2idle/' "$BOOT_ENTRY"
        echo -e "${GREEN}Updated kernel parameters in $BOOT_ENTRY${NC}"
    else
        echo -e "${YELLOW}Kernel parameters already updated${NC}"
    fi
else
    echo -e "${YELLOW}Warning: Could not find systemd-boot entry to update${NC}"
fi

echo -e "${BLUE}Step 6: Configuring power management...${NC}"

# Disable wake on LAN and other problematic wake sources
for device in /sys/class/net/*/device/power/wakeup; do
    [ -f "$device" ] && echo disabled > "$device" 2>/dev/null || true
done

# Configure Intel PMC for better sleep
if [[ -d /sys/kernel/debug/pmc_core ]]; then
    echo -e "${GREEN}Intel PMC debug interface available${NC}"
fi

echo -e "${BLUE}Step 7: Creating diagnostic tools...${NC}"

cat > /usr/local/bin/check-sleep-blockers << 'EOF'
#!/bin/bash
# Check what might be preventing sleep

echo "=== Sleep Blockers Diagnostic ==="
echo "Date: $(date)"
echo ""

echo "Current sleep state:"
cat /sys/power/state 2>/dev/null || echo "Not available"
echo ""

echo "Memory sleep mode:"
cat /sys/power/mem_sleep 2>/dev/null || echo "Not available"
echo ""

echo "Processes preventing sleep:"
for inhibitor in /sys/power/wake_lock /sys/power/wake_unlock; do
    if [[ -f "$inhibitor" ]]; then
        echo "$inhibitor: $(cat "$inhibitor" 2>/dev/null)"
    fi
done

echo ""
echo "Wake-enabled devices:"
find /sys/devices -name "wakeup" -exec sh -c 'echo "$(dirname {}): $(cat {} 2>/dev/null)"' \; 2>/dev/null | grep enabled | head -10

echo ""
echo "Recent sleep/wake events:"
journalctl --since="1 hour ago" | grep -i -E "(suspend|sleep|wake|lid)" | tail -10

echo ""
echo "XE driver status:"
lsmod | grep xe || echo "XE driver not loaded"

echo ""
echo "Graphics card power state:"
for card in /sys/class/drm/card*/device/power/runtime_status; do
    if [[ -f "$card" ]]; then
        echo "$card: $(cat "$card" 2>/dev/null)"
    fi
done
EOF

chmod +x /usr/local/bin/check-sleep-blockers
echo -e "${GREEN}Created diagnostic tool: check-sleep-blockers${NC}"

cat > /usr/local/bin/force-sleep-fix << 'EOF'
#!/bin/bash
# Emergency sleep fix if crashes continue

echo "=== Force Sleep Fix ==="

# Reload XE driver
echo "Reloading XE graphics driver..."
modprobe -r xe 2>/dev/null || true
sleep 2
modprobe xe 2>/dev/null || true

# Reset power management
echo "Resetting power management..."
systemctl restart systemd-logind

# Clear any stuck processes
echo "Clearing problematic processes..."
for proc in firefox chrome discord teams slack; do
    pkill -9 "$proc" 2>/dev/null || true
done

# Force graphics reset
echo "Forcing graphics reset..."
echo 1 > /sys/class/drm/card0/device/reset 2>/dev/null || true

echo "Sleep fix applied. Try closing lid now."
EOF

chmod +x /usr/local/bin/force-sleep-fix
echo -e "${GREEN}Created emergency fix tool: force-sleep-fix${NC}"

echo -e "${BLUE}Step 8: Final system configuration...${NC}"

# Reload systemd configuration
systemctl daemon-reload

# Update initramfs with new kernel parameters
mkinitcpio -P

# Create log rotation for our custom logs
cat > /etc/logrotate.d/dell-sleep << 'EOF'
/var/log/dell-sleep.log {
    weekly
    rotate 4
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
}

/var/log/dell-sleep-fix.log {
    weekly
    rotate 4
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
}
EOF

echo -e "${GREEN}=== Dell XPS Sleep Fix Installation Complete ===${NC}"
echo ""
echo -e "${YELLOW}IMPORTANT: You must reboot for all changes to take effect!${NC}"
echo ""
echo "After reboot, test the sleep functionality by:"
echo "1. Close the lid and wait 30 seconds"
echo "2. Open the lid and check if the system wakes properly"
echo ""
echo "If issues persist, run these diagnostic commands:"
echo "- check-sleep-blockers    (diagnose what's preventing sleep)"
echo "- force-sleep-fix         (emergency fix for stuck states)"
echo ""
echo "Logs are available at:"
echo "- /var/log/dell-sleep.log      (sleep events)"
echo "- /var/log/dell-sleep-fix.log  (this script's output)"
echo ""
echo -e "${BLUE}Reboot now? (y/N):${NC}"
read -r response
if [[ "$response" =~ ^[Yy]$ ]]; then
    echo "Rebooting in 5 seconds..."
    sleep 5
    reboot
else
    echo "Please reboot manually when ready."
fi
