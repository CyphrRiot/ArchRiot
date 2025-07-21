# 🔧 ArchRiot Optional Tools

**Advanced tools for power users - NOT part of the standard ArchRiot installation**

These tools are provided for advanced users who need specific functionality beyond the core ArchRiot experience. They are completely optional and should only be used if you understand the risks and requirements.

## ⚠️ Important Safety Notice

- **These tools can modify critical system components**
- **Always have a backup and recovery plan**
- **Test thoroughly in a virtual machine first**
- **Only use if you understand the risks involved**

## 🛠️ Available Tools

### Dell XPS Sleep Crash Fix

**Purpose:** Fix sleep/suspend crashes on Dell XPS laptops with Intel Arc Graphics
**Location:** `dell-sleep-fix/setup-dell-sleep-fix.sh`
**Risk Level:** ⚠️ MODERATE - Modifies kernel parameters and power management
**Compatibility:** Dell XPS laptops with Intel Lunar Lake Arc Graphics (130V/140V)

**What it does:**

- Fixes crashes when closing laptop lid or entering sleep mode
- Stabilizes Intel XE graphics driver for sleep/wake cycles
- Configures s2idle sleep mode for maximum compatibility
- Adds Dell-specific power management workarounds
- Creates diagnostic and emergency fix tools

**Features:**

- ✅ Fixes Dell XPS lid-close crashes
- ✅ Stabilizes Intel Arc Graphics sleep/wake
- ✅ Automatic pre/post sleep state management
- ✅ Comprehensive diagnostic tools
- ✅ Emergency recovery utilities

### Secure Boot Setup

**Purpose:** Implement UEFI Secure Boot using standard Arch Linux methods
**Location:** `secure-boot/setup-secure-boot.sh`
**Risk Level:** ⚠️ MODERATE - Uses tested Arch methods
**Compatibility:** AMD, Intel, Any UEFI system with Secure Boot

**What it does:**

- Uses `sbctl` (official Arch package) for key management
- Uses `shim-signed` (AUR) for Microsoft hardware compatibility
- Follows Arch Wiki recommendations exactly
- Automatically signs kernels on updates
- Supports Windows dual-boot scenarios

**Features:**

- ✅ Microsoft hardware compatibility guaranteed
- ✅ Automatic kernel signing with pacman hooks
- ✅ Windows dual-boot support
- ✅ Hardware firmware update compatibility
- ✅ Well-tested by Arch community

## 🚀 Usage

### Prerequisites

Before using any optional tool:

1. **Fresh ArchRiot installation** (recommended)
2. **System backup** with ArchRiot's Migrate tool
3. **Internet connection** (for package downloads)
4. **Understanding of the risks involved**
5. **Backup/recovery capability**

### Running Dell XPS Sleep Fix

```bash
# Navigate to ArchRiot directory
cd /path/to/ArchRiot

# Make executable
chmod +x optional-tools/dell-sleep-fix/setup-dell-sleep-fix.sh

# Run the installer
sudo ./optional-tools/dell-sleep-fix/setup-dell-sleep-fix.sh

# IMPORTANT: Reboot after installation
sudo reboot
```

**The script will:**

1. Check hardware compatibility (Dell XPS + Intel Arc Graphics)
2. Configure Intel XE graphics driver parameters
3. Set up systemd sleep configuration
4. Install pre/post sleep hooks
5. Add kernel parameters for stability
6. Create diagnostic and recovery tools

### Running Secure Boot Setup

```bash
# Navigate to ArchRiot directory
cd /path/to/ArchRiot

# Make executable
chmod +x optional-tools/secure-boot/setup-secure-boot.sh

# Run the installer
./optional-tools/secure-boot/setup-secure-boot.sh
```

**The script will:**

1. Check system compatibility
2. Install required packages (`sbctl`, `shim-signed`)
3. Guide you through UEFI setup
4. Create and enroll Secure Boot keys
5. Sign bootloader and kernels
6. Verify the setup

## 🔍 Safety Features

### Built-in Protections

- **Pre-flight checks:** Verifies system compatibility
- **User confirmation:** Requires explicit approval for risky operations
- **Comprehensive logging:** All actions logged for troubleshooting
- **Verification steps:** Checks setup integrity before completion

### Recovery Options

- **System restore:** Use ArchRiot's Migrate tool to restore backups
- **Disable features:** Instructions provided for disabling changes
- **Emergency tools:** Built-in recovery utilities
- **Manual rollback:** Detailed uninstall instructions in each tool's README

## 📋 System Requirements

### For Dell XPS Sleep Fix

- **Dell XPS laptop** with Intel Lunar Lake processor
- **Intel Arc Graphics 130V/140V** (check with `lspci | grep VGA`)
- **XE graphics driver** (automatically detected)
- **systemd-boot** (ArchRiot default)
- **Kernel 6.15+** (included in ArchRiot)

### For Secure Boot Setup

- **UEFI firmware** (not Legacy BIOS)
- **Secure Boot capable system**
- **Setup Mode access** in UEFI/BIOS
- **Internet connection** for package installation

### Tested Hardware

#### Dell Sleep Fix:

- ✅ Dell XPS 13 Plus (Intel Lunar Lake)
- ✅ Dell XPS 15 (Intel Arc Graphics 130V/140V)
- ✅ Other Dell XPS models with Intel Arc Graphics
- ❌ Non-Dell laptops (Dell-specific workarounds)
- ❌ AMD or NVIDIA graphics (Intel-specific fixes)

#### Secure Boot:

- ✅ AMD Ryzen systems
- ✅ Intel Core/Xeon systems
- ✅ Standard UEFI motherboards
- ✅ Microsoft Surface devices
- ✅ Most gaming laptops/desktops

## 🚨 Important Notes

### Before You Start

1. **Create system backup** with ArchRiot's Migrate tool
2. **Understand your hardware** - verify compatibility first
3. **Have recovery media ready** - ArchRiot USB or Arch Linux ISO
4. **Read the specific tool's README** for detailed information

### During Setup

- **Follow prompts carefully** - the scripts guide you through each step
- **Don't interrupt the process** - let installations complete fully
- **Pay attention to warnings** - some operations require manual steps

### After Setup

- **Test thoroughly** - ensure all hardware works correctly
- **Monitor logs** - check for any issues or warnings
- **Keep recovery plan** - know how to disable features if needed

## 🔄 Maintenance

### Dell Sleep Fix

- **Automatic:** No maintenance required - hooks run automatically
- **Manual checks:** Verify configurations persist after updates
- **Diagnostics:** Use `check-sleep-blockers` command

### Secure Boot

- **Automatic:** Kernel signing happens automatically on updates
- **Manual checks:** Use `sudo sbctl status` to verify
- **Maintenance:** Occasionally verify signed files

## 🤝 Support

### Getting Help

1. **Check tool-specific README** - Each tool has detailed documentation
2. **Check logs** - Located in `/var/log/` directory
3. **Arch Wiki** - Comprehensive documentation for underlying technologies
4. **ArchRiot Issues** - GitHub issues for ArchRiot-specific problems

### Common Issues

- **Hardware incompatibility:** Verify your hardware matches requirements
- **Boot failure:** Use recovery media to disable problematic features
- **Driver conflicts:** Check that you're using the correct drivers

## 📁 Tool Structure

```
optional-tools/
├── README.md                 # This file
├── launcher.sh              # Future: unified launcher
├── dell-sleep-fix/
│   ├── README.md            # Detailed Dell sleep fix documentation
│   ├── setup-dell-sleep-fix.sh  # Main installation script
│   ├── xe-graphics-fix.conf     # XE driver configuration
│   ├── dell-sleep-fix.conf      # systemd sleep configuration
│   └── dell-sleep-hook.sh       # Pre/post sleep scripts
└── secure-boot/
    ├── README.md            # Detailed Secure Boot documentation
    └── setup-secure-boot.sh # Main installation script
```

## 📄 License

These optional tools are part of ArchRiot and licensed under the same terms as the main project.

---

**Remember:** These tools modify critical system components. Always have a backup plan and test thoroughly!

🛡️⚔️🪐 **Use at your own risk - Hack responsibly** 🪐⚔️🛡️
