# 🛠️ Dell XPS Sleep Crash Fix

**Purpose:** Fix sleep/suspend crashes on Dell XPS laptops with Intel Arc Graphics
**Risk Level:** ⚠️ MODERATE - Modifies kernel parameters and power management
**Compatibility:** Dell XPS laptops with Intel Lunar Lake Arc Graphics (130V/140V)

## 🚨 Problem Description

Dell XPS laptops with the new Intel Lunar Lake Arc Graphics (130V/140V) using the XE driver can experience system crashes when:
- Closing the laptop lid
- Entering sleep/suspend mode
- Waking from sleep

This typically manifests as:
- System completely freezes/crashes instead of sleeping
- Black screen on wake that requires hard reset
- System reboots unexpectedly after sleep
- Complete system lock-up requiring power button reset

## 🔧 What This Tool Fixes

### Root Causes Addressed:
1. **Intel XE Graphics Driver Issues** - New bleeding-edge driver with sleep bugs
2. **S2idle Power Management** - Modern standby compatibility problems
3. **Dell-Specific Hardware Quirks** - Vendor-specific ACPI/power management issues
4. **Kernel Parameter Conflicts** - Default settings incompatible with new hardware

### Technical Solutions Applied:
- **XE Driver Stabilization** - Disables problematic features (PSR, FBC, RC6)
- **Sleep Mode Optimization** - Forces stable s2idle-only mode
- **Pre/Post Sleep Hooks** - Manages graphics state transitions
- **Kernel Parameter Tuning** - Adds Intel Arc specific boot parameters
- **Power Management Override** - Dell-specific workarounds

## 📋 Prerequisites

### Hardware Requirements:
- **Dell XPS laptop** with Intel Lunar Lake processor
- **Intel Arc Graphics 130V or 140V**
- **UEFI boot mode** (systemd-boot)

### Software Requirements:
- **ArchRiot installed** and running
- **Kernel 6.15+** (included in ArchRiot)
- **XE graphics driver** (automatically detected)
- **systemd-boot** (ArchRiot default)

### Before Running:
```bash
# Verify you have the right hardware
lspci | grep VGA
# Should show: Intel Corporation Lunar Lake [Intel Arc Graphics 130V / 140V]

# Check if XE driver is loaded
lsmod | grep xe
# Should show xe driver loaded

# Verify systemd-boot
ls /boot/loader/entries/
# Should show .conf files
```

## 🚀 Installation

### Step 1: Navigate to Tool
```bash
cd /path/to/ArchRiot/optional-tools/dell-sleep-fix
```

### Step 2: Make Executable
```bash
chmod +x setup-dell-sleep-fix.sh
```

### Step 3: Run Installation
```bash
sudo ./setup-dell-sleep-fix.sh
```

### Step 4: Reboot
**IMPORTANT:** You MUST reboot for changes to take effect
```bash
sudo reboot
```

## 🔍 What Gets Modified

### System Files Created/Modified:
```
/etc/modprobe.d/xe-graphics-fix.conf          # XE driver parameters
/etc/systemd/sleep.conf.d/dell-sleep-fix.conf # Sleep configuration
/lib/systemd/system-sleep/dell-sleep-hook.sh  # Pre/post sleep scripts
/boot/loader/entries/*.conf                   # Kernel parameters
/etc/logrotate.d/dell-sleep                   # Log rotation
```

### Utilities Installed:
```
/usr/local/bin/check-sleep-blockers  # Diagnostic tool
/usr/local/bin/force-sleep-fix       # Emergency fix tool
```

### Kernel Parameters Added:
```
intel_iommu=off          # Disable IOMMU for stability
xe.enable_guc=3          # Enable GuC for power management
xe.enable_psr=0          # Disable Panel Self Refresh (causes crashes)
xe.enable_fbc=0          # Disable Frame Buffer Compression
mem_sleep_default=s2idle # Force s2idle sleep mode
```

## 🧪 Testing

### After Installation:
1. **Reboot the system**
2. **Test lid close/open**:
   - Close laptop lid
   - Wait 30 seconds
   - Open lid
   - System should wake properly

3. **Test manual suspend**:
   ```bash
   systemctl suspend
   ```

### If Problems Persist:
```bash
# Run diagnostic tool
check-sleep-blockers

# Check what's preventing sleep
journalctl -f &
# Then try to sleep and watch logs

# Emergency fix if stuck
force-sleep-fix
```

## 📊 Diagnostic Tools

### Check Sleep Blockers
```bash
check-sleep-blockers
```
**Shows:**
- Current sleep state capabilities
- Wake-enabled devices
- Recent sleep/wake events
- XE driver status
- Graphics power states

### Force Sleep Fix
```bash
force-sleep-fix
```
**Emergency Actions:**
- Reloads XE graphics driver
- Restarts systemd-logind
- Kills problematic processes
- Forces graphics reset

### Manual Diagnostics
```bash
# Check sleep capabilities
cat /sys/power/state
cat /sys/power/mem_sleep

# Monitor real-time sleep events
journalctl -f | grep -i sleep

# Check XE driver status
dmesg | grep xe

# View graphics power state
cat /sys/class/drm/card*/device/power/runtime_status
```

## 📝 Logs

### Automatic Logging:
- **Sleep Events:** `/var/log/dell-sleep.log`
- **Installation:** `/var/log/dell-sleep-fix.log`
- **System Events:** `journalctl | grep sleep`

### Log Rotation:
- Logs rotate weekly
- 4 weeks retention
- Compressed after rotation

## 🔄 Maintenance

### Automatic:
- **No maintenance required** - hooks run automatically
- **Kernel updates** - kernel parameters persist
- **Log rotation** - handled automatically

### Manual Checks:
```bash
# Verify configurations are still in place
ls -la /etc/modprobe.d/xe-graphics-fix.conf
ls -la /etc/systemd/sleep.conf.d/dell-sleep-fix.conf
ls -la /lib/systemd/system-sleep/dell-sleep-hook.sh

# Check if kernel parameters are active
cat /proc/cmdline | grep xe.enable_psr=0
```

## ⚠️ Known Issues & Limitations

### May Not Work For:
- **Non-Dell laptops** (Dell-specific workarounds)
- **Older Intel graphics** (XE driver specific)
- **AMD graphics** (Intel-specific fixes)
- **NVIDIA primary graphics** (different power management)

### Side Effects:
- **Slightly higher power consumption** (disabled power saving features)
- **No hibernation** (disabled for stability)
- **Longer boot times** (additional kernel parameters)

### If It Doesn't Work:
1. **Check hardware compatibility** - must be Intel Lunar Lake Arc Graphics
2. **Verify XE driver** - older systems use i915 driver
3. **Update BIOS** - Dell may have fixed issues in newer BIOS
4. **Check logs** - look for specific error patterns

## 🔙 Uninstall

### To Remove All Changes:
```bash
# Remove configuration files
sudo rm -f /etc/modprobe.d/xe-graphics-fix.conf
sudo rm -f /etc/systemd/sleep.conf.d/dell-sleep-fix.conf
sudo rm -f /lib/systemd/system-sleep/dell-sleep-hook.sh
sudo rm -f /etc/logrotate.d/dell-sleep

# Remove utilities
sudo rm -f /usr/local/bin/check-sleep-blockers
sudo rm -f /usr/local/bin/force-sleep-fix

# Restore original kernel parameters
sudo cp /boot/loader/entries/*.conf.backup-* /boot/loader/entries/
# (Rename to remove .backup-timestamp suffix)

# Reload configurations
sudo systemctl daemon-reload
sudo mkinitcpio -P

# Reboot
sudo reboot
```

## 🤝 Support

### Getting Help:
1. **Check logs first** - `/var/log/dell-sleep.log`
2. **Run diagnostics** - `check-sleep-blockers`
3. **Search issues** - ArchRiot GitHub issues
4. **Dell forums** - Dell-specific hardware issues

### Reporting Bugs:
Include this information:
```bash
# Hardware info
lspci | grep VGA
uname -r
dmesg | grep xe | head -10

# Configuration status
ls -la /etc/modprobe.d/xe-graphics-fix.conf
cat /proc/cmdline

# Recent logs
tail -20 /var/log/dell-sleep.log
```

## 📚 Technical Details

### Why This Happens:
- **Intel Arc Graphics** is bleeding-edge hardware (2024/2025)
- **XE driver** is newer than i915, still has bugs
- **Dell power management** has vendor-specific quirks
- **S2idle vs S3 sleep** compatibility issues
- **ACPI implementation** varies between manufacturers

### Why These Fixes Work:
- **Disable PSR/FBC** - Known to cause crashes on new Intel graphics
- **Force s2idle** - Most compatible sleep mode for modern hardware
- **Power state management** - Prevents graphics from getting stuck
- **Process handling** - Prevents applications from blocking sleep

---

**⚠️ Important:** This tool specifically targets Dell XPS laptops with Intel Lunar Lake Arc Graphics. Using it on incompatible hardware may cause issues.

🛡️⚔️🪐 **Hack the Planet (But Sleep Peacefully)** 🪐⚔️🛡️
