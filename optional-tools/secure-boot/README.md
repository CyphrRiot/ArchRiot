# 🛡️ ArchRiot Secure Boot Setup

**Clean, safe, and simple UEFI Secure Boot implementation for ArchRiot**

This tool implements Secure Boot using the **standard Arch Linux method** with `sbctl` and `shim-signed`. It's designed to work cleanly on any firmware (Intel, AMD) while maintaining maximum compatibility and safety.

## ⚠️ Important Safety Information

**This tool modifies critical boot components. Always have a backup plan!**

- ✅ **Safe**: Uses official Arch packages and follows Arch Wiki exactly
- ✅ **Compatible**: Works with Microsoft hardware requirements
- ✅ **Automatic**: Handles kernel signing on updates
- ⚠️ **Advanced**: Requires basic UEFI/BIOS knowledge
- ⚠️ **Risky**: Can prevent boot if misconfigured

## 🚀 Quick Start

### Prerequisites Checklist

Before you start, ensure you have:

- [ ] **Fresh ArchRiot installation** (recommended)
- [ ] **UEFI boot mode** (not Legacy BIOS)
- [ ] **Secure Boot capable system**
- [ ] **Access to UEFI/BIOS setup** (know your key: F2/F12/Delete)
- [ ] **Internet connection** (for package downloads)
- [ ] **Backup or recovery plan** (ArchRiot Migrate tool recommended)

### Quick Setup (Recommended)

```bash
# Navigate to ArchRiot directory
cd /path/to/ArchRiot

# Run the setup script
./optional-tools/secure-boot/setup-secure-boot.sh
```

**The script will guide you through:**
1. System compatibility checks
2. Package installation (`sbctl`, `shim-signed`)
3. UEFI setup instructions
4. Key creation and enrollment
5. Bootloader and kernel signing
6. Verification and next steps

## 📋 Detailed Setup Process

### Step 1: Preparation

1. **Boot into ArchRiot** normally
2. **Open terminal** as regular user (not root)
3. **Ensure internet** connection is available
4. **Close unnecessary applications**

### Step 2: Run Setup Script

```bash
./optional-tools/secure-boot/setup-secure-boot.sh
```

The script provides a menu with options:
- **Quick Setup** (recommended for most users)
- **Check Current Status** (see what's already configured)
- **Verify Existing Setup** (check if already set up)
- **Advanced Manual Setup** (step-by-step control)

### Step 3: UEFI Configuration

When prompted, you'll need to:

1. **Reboot** your system
2. **Enter UEFI/BIOS setup** (press F2/F12/Delete during boot)
3. **Navigate** to Security → Secure Boot settings
4. **Enable "Setup Mode"** or **"Clear Keys"**
5. **Save and exit** to boot back to ArchRiot
6. **Run the script again** to continue

### Step 4: Complete Setup

After UEFI configuration:

1. Script will **create Secure Boot keys**
2. **Enroll keys** into firmware
3. **Sign bootloader** and kernels
4. **Verify** everything is working
5. **Provide next steps**

### Step 5: Final UEFI Configuration

Final step requires another UEFI visit:

1. **Reboot** again
2. **Enter UEFI/BIOS setup**
3. **Disable "Setup Mode"**
4. **Enable "Secure Boot"**
5. **Save and exit**

## 🔧 Technical Details

### What the Script Does

1. **System Checks**: Verifies UEFI mode, internet, Arch Linux
2. **Package Installation**:
   - `sbctl` (official Arch package for key management)
   - `shim-signed` (AUR package for Microsoft compatibility)
3. **Key Management**:
   - Creates custom Secure Boot keys
   - Enrolls keys in firmware
4. **File Signing**:
   - Signs systemd-boot bootloader
   - Signs Linux kernel(s)
   - Signs initramfs
5. **Automatic Updates**: Sets up pacman hooks for future kernel signing

### Files That Get Signed

- `/boot/vmlinuz-linux` (main kernel)
- `/boot/vmlinuz-*` (other kernels)
- `/boot/initramfs-linux.img` (initramfs)
- `/efi/EFI/systemd/systemd-bootx64.efi` (bootloader)

### Logging

All operations are logged to: `/var/log/archriot-secureboot.log`

## 🔍 Verification & Maintenance

### Check Secure Boot Status

```bash
# Check overall status
sudo sbctl status

# Verify all files are signed
sudo sbctl verify

# List signed files
sudo sbctl list-files

# Check boot status
bootctl status
```

### Automatic Maintenance

- **Kernel updates**: Automatically signed during pacman updates
- **No manual intervention** needed for normal operations
- **Monitor occasionally** with `sudo sbctl status`

## 🚨 Troubleshooting

### Common Issues

**Boot Failure After Enabling Secure Boot**
1. Enter UEFI/BIOS setup
2. Disable Secure Boot temporarily
3. Boot normally
4. Run `sudo sbctl verify` to check signed files
5. Sign any missing files: `sudo sbctl sign -s <file_path>`

**Keys Not Enrolling**
- Ensure "Setup Mode" is enabled in UEFI
- Some systems call it "Clear Keys" or "Factory Reset Keys"
- Try different UEFI menu locations (Security/Boot/Advanced)

**Hardware Compatibility Issues**
- Some older systems may not support modern Secure Boot
- Try updating UEFI/BIOS firmware first
- Check manufacturer documentation

**Package Installation Fails**
- Ensure internet connection is stable
- Update package database: `sudo pacman -Sy`
- For AUR issues, check if base-devel is installed

### Recovery Options

**Complete Recovery (if system won't boot)**
1. Boot from ArchRiot USB/Arch Linux ISO
2. Mount your root filesystem
3. Chroot into system
4. Run `sudo sbctl status` to diagnose
5. Re-sign files or disable Secure Boot

**Quick Recovery (system boots but Secure Boot disabled)**
1. Boot normally
2. Run verification: `sudo sbctl verify`
3. Sign missing files: `sudo sbctl sign -s <file>`
4. Re-enable Secure Boot in UEFI

## 📊 Compatibility

### Tested Hardware

✅ **AMD Systems**
- Ryzen 3000/5000/7000 series
- EPYC server processors
- Standard AM4/AM5 motherboards

✅ **Intel Systems**
- Core 8th gen and newer
- Xeon server processors
- Standard LGA1200/LGA1700 motherboards

✅ **Laptop/Mobile**
- Gaming laptops (ASUS, MSI, etc.)
- Business laptops (ThinkPad, Dell, HP)
- Microsoft Surface devices

✅ **Special Cases**
- Windows dual-boot configurations
- Multiple Linux distributions
- Server hardware

### Known Limitations

⚠️ **Legacy Systems**
- Pre-2012 hardware may lack proper UEFI support
- Some older motherboards have incomplete Secure Boot implementation

⚠️ **Custom Firmware**
- Coreboot/Libreboot systems may need special configuration
- Some gaming motherboards with custom UEFI may behave differently

## 🔒 Security Benefits

### What Secure Boot Provides

- **Boot integrity**: Prevents malicious bootloaders
- **Kernel verification**: Ensures kernel hasn't been tampered with
- **Chain of trust**: Maintains security from firmware to OS
- **Hardware compliance**: Required for some enterprise environments

### What It Doesn't Protect Against

- **Runtime attacks**: Only secures boot process
- **User-space malware**: Doesn't protect running system
- **Hardware attacks**: Physical access can bypass
- **Social engineering**: User-level compromise still possible

## 📚 Additional Resources

### Official Documentation
- [Arch Wiki - Secure Boot](https://wiki.archlinux.org/title/Unified_Extensible_Firmware_Interface/Secure_Boot)
- [sbctl GitHub](https://github.com/Foxboron/sbctl)

### ArchRiot Resources
- [Main ArchRiot Documentation](../../README.md)
- [ArchRiot Migrate Tool](https://github.com/CyphrRiot/Migrate)

### UEFI/Hardware Resources
- Your motherboard manufacturer's documentation
- System UEFI/BIOS manual

## 🆘 Getting Help

### Before Asking for Help

1. **Check logs**: `/var/log/archriot-secureboot.log`
2. **Run diagnostics**: `sudo sbctl status` and `sudo sbctl verify`
3. **Try recovery steps** listed above
4. **Document your hardware** and exact error messages

### Where to Get Help

1. **Arch Wiki**: Most comprehensive Secure Boot documentation
2. **ArchRiot GitHub Issues**: For ArchRiot-specific problems
3. **Arch Linux Forums**: Community support
4. **Reddit r/archlinux**: Community discussions

---

**Remember**: Secure Boot is a powerful security feature, but it adds complexity. Only enable it if you understand the trade-offs and have a solid backup plan.

🛡️⚔️🪐 **Secure your boot, hack responsibly** 🪐⚔️🛡️
