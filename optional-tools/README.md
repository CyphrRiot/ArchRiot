# 🔧 ArchRiot Optional Tools

**Advanced tools for power users - NOT part of the standard ArchRiot installation**

These tools are provided for advanced users who need specific functionality beyond the core ArchRiot experience. They are completely optional and should only be used if you understand the risks and requirements.

## ⚠️ Important Safety Notice

- **These tools can modify critical system components**
- **Always have a backup and recovery plan**
- **Test thoroughly in a virtual machine first**
- **Only use if you understand UEFI/BIOS concepts**

## 🛡️ Available Tools

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
2. **UEFI boot mode** (required for Secure Boot)
3. **Internet connection** (for package downloads)
4. **Basic UEFI/BIOS knowledge**
5. **Backup/recovery capability**

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
- **Pre-flight checks:** Verifies UEFI mode, internet, etc.
- **User confirmation:** Requires explicit approval for risky operations
- **Comprehensive logging:** All actions logged to `/var/log/archriot-secureboot.log`
- **Verification steps:** Checks setup integrity before completion

### Recovery Options
- **Disable Secure Boot:** Enter UEFI and turn off Secure Boot if issues occur
- **Verify files:** Use `sudo sbctl verify` to check signed files
- **Re-sign files:** Use `sudo sbctl sign -s <file>` to sign missing files

## 📋 System Requirements

### For Secure Boot Setup
- **UEFI firmware** (not Legacy BIOS)
- **Secure Boot capable system**
- **Setup Mode access** in UEFI/BIOS
- **Internet connection** for package installation

### Tested Hardware
- ✅ AMD Ryzen systems
- ✅ Intel Core/Xeon systems
- ✅ Standard UEFI motherboards
- ✅ Microsoft Surface devices
- ✅ Most gaming laptops/desktops

## 🚨 Important Notes

### Before You Start
1. **Test in VM first** if possible
2. **Create system backup** with ArchRiot's Migrate tool
3. **Understand your hardware** - know how to access UEFI setup
4. **Have recovery media ready** - ArchRiot USB or Arch Linux ISO

### During Setup
- **Follow prompts carefully** - the script guides you through each step
- **Don't skip UEFI steps** - Setup Mode must be enabled first
- **Don't interrupt the process** - let it complete fully

### After Setup
- **Test thoroughly** - ensure all hardware works correctly
- **Monitor updates** - kernels are automatically signed but verify occasionally
- **Keep recovery plan** - know how to disable Secure Boot if needed

## 🔄 Maintenance

### Automatic Maintenance
- **Kernel signing** happens automatically on pacman updates
- **No manual intervention** needed for normal updates

### Manual Checks
```bash
# Check Secure Boot status
sudo sbctl status

# Verify all files are signed
sudo sbctl verify

# List signed files
sudo sbctl list-files
```

## 🤝 Support

### Getting Help
1. **Check logs first:** `/var/log/archriot-secureboot.log`
2. **Verify setup:** Run `sudo sbctl status`
3. **Arch Wiki:** Comprehensive Secure Boot documentation
4. **ArchRiot Issues:** GitHub issues for ArchRiot-specific problems

### Common Issues
- **Boot failure:** Disable Secure Boot in UEFI, then troubleshoot
- **Unsigned files:** Use `sudo sbctl sign -s <file>` to sign missing files
- **Hardware incompatibility:** Some older systems may not support modern Secure Boot

## 📄 License

These optional tools are part of ArchRiot and licensed under the same terms as the main project.

---

**Remember:** These tools modify critical system components. Always have a backup plan and test thoroughly!

🛡️⚔️🪐 **Use at your own risk - Hack responsibly** 🪐⚔️🛡️
