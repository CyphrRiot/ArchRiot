<div align="center">

# :: 𝔸𝕣𝕔𝕙ℝ𝕚𝕠𝕥 ::

![Version](https://img.shields.io/badge/version-4.2.8-blue?labelColor=0052cc)
![License](https://img.shields.io/github/license/CyphrRiot/ArchRiot?color=4338ca&labelColor=3730a3)
![Platform](https://img.shields.io/badge/platform-linux-4338ca?logo=linux&logoColor=white&labelColor=3730a3)
![Arch Linux](https://img.shields.io/badge/Arch_Linux-1e1b4b?logo=arch-linux&logoColor=8b5cf6&labelColor=0f172a)
![Wayland](https://img.shields.io/badge/Wayland-312e81?logo=wayland&logoColor=a855f7&labelColor=1e1b4b)

![Last Commit](https://img.shields.io/github/last-commit/CyphrRiot/ArchRiot?color=5b21b6&labelColor=4c1d95)
![Code Size](https://img.shields.io/github/languages/code-size/CyphrRiot/ArchRiot?color=4338ca&labelColor=3730a3)
![Code](https://img.shields.io/badge/human-coded-blue?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmZmZmYiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIiBjbGFzcz0ibHVjaWRlIGx1Y2lkZS1wZXJzb24tc3RhbmRpbmctaWNvbiBsdWNpZGUtcGVyc29uLXN0YW5kaW5nIj48Y2lyY2xlIGN4PSIxMiIgY3k9IjUiIHI9IjEiLz48cGF0aCBkPSJtOSAyMCAzLTYgMyA2Ii8+PHBhdGggZD0ibTYgOCA2IDIgNi0yIi8+PHBhdGggZD0iTTEyIDEwdjQiLz48L3N2Zz4=&logoColor=a855f7&labelColor=1e1b4b)

![Language](https://img.shields.io/badge/language-Go-4338ca?logo=go&logoColor=c7d2fe&labelColor=3730a3)
![Language](https://img.shields.io/badge/language-YAML-5b21b6?logo=yaml&logoColor=e0e7ff&labelColor=4c1d95)
![Language](https://img.shields.io/badge/Python-312e81?label=language&logo=python&logoColor=c7d2fe&labelColor=1e1b4b&color=312e81)

</div>

## **ArchRiot: The (Arch) Linux System You've Always Wanted (but never had)**

### **One Command. Complete Environment. Zero Compromises.**

ArchRiot is the answer to every time you've thought "why can't Linux just work correctly from the start?" We've spent hundreds of hours perfecting the details so you get a blazing-fast, secure, beautiful system that actually respects your time and intelligence.

**Curated to be correct:**

- **🪟 Hyprland Tiling** - Makes other Window Managers feel primitive
- **⚡ Go Installer** - Atomic operations, instant rollbacks, zero dependency hell
- **🛡️ Privacy** - Zero telemetry, zero tracking, zero corporate data harvesting
- **🎨 Aesthetics** - Carefully crafted dark themes that work at any hour
- **💻 Development** - Zed, Neovim, shell enhancements, and other upgrades

_Built on Arch Linux with Hyprland, because compromises are for other people. This isn't maintained by committee or corporate roadmap -- it's maintained by someone with an obsessive, singular focus on getting it right the first time, because crappy Linux environments are an insult to what computing should be._

![ArchRiot Screenshot](config/images/screenshot.png)

<iframe width="560" height="315" src="https://www.youtube.com/embed/qrraIOvAcdg?si=yooSbgD0FJKLuAkc" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>

## 📚 Navigate This Guide

- [🚀 Choose Your ArchRiot Experience](#choose-your-archriot-experience)
    - [🔥 Method 1: Install Script](#method-1-install-script)
    - [⚡ Method 2: ArchRiot ISO](#method-2-archriot-iso)
        - [ISO Install: Wi‑Fi abort fix (iwd PSKs)](#iso-wifi-abort)

- [⌨️ Master Your ArchRiot Desktop](#master-your-archriot-desktop)
- [🎛️ Control Panel](#archriot-control-panel)
- [🎯 Key Customizations](#key-customizations)
- [🔄 System Management](#system-management)
- [🧰 Advanced Usage: CLI Flags](#advanced-usage-archriot-cli-flags)
- [🔧 Troubleshooting](#troubleshooting)
- [📄 License](#license)

## ⚡ Quick Install

- 🔥 Method 1 — Install Script (existing Arch):

    ```bash
    curl -fsSL https://ArchRiot.org/setup.sh | bash
    ```

    See: [🔥 Method 1: Install Script](#method-1-install-script)

- ⚡ Method 2 — ArchRiot ISO (fresh install):
    - Download from Releases and boot the ISO:
      https://github.com/CyphrRiot/ArchRiot/releases
    - See: [⚡ Method 2 — ArchRiot ISO](#method-2-archriot-iso)

<a id="choose-your-archriot-experience"></a>

## 🚀 Choose Your ArchRiot Experience

<a id="method-1-install-script"></a>

### 🔥 Method 1: Install Script

#### You already have Arch Linux installed

**Transform your current Arch system into ArchRiot**

```bash
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

**Perfect for:**

- 🏠 **System preservation** - Keep your data, configs, and custom tweaks intact
- 🔧 **Arch variants** - CachyOS, Manjaro, EndeavourOS, or pure Arch installations
- 🎨 **Desktop upgrade** - Transform just your desktop environment, leave the rest alone
- ⚡ **Quick wins** - Get ArchRiot's best features without starting over

**What you get:**

- ArchRiot desktop environment and apps
- Hyprland tiling window manager
- CypherRiot themes and customizations
- Keeps your existing system intact

#### **Note:** For fresh installations, we also have an **ArchRiot ISO available** (See below "Method 2")

<a id="method-2-archriot-iso"></a>

### ⚡ Method 2: ArchRiot ISO

#### You do NOT have Arch Linux installed

![ArchRiot Installation Demo](config/images/riot.gif)

#### ⚠️ Warning: ISO will replace a drive with ArchRiot Linux. ⚠️

1. **📥 Download ArchRiot Linux ISO**
    - **[ArchRiot Linux ISO](https://github.com/CyphrRiot/ArchRiot/releases/download/v2.27/archriot.iso)**
    - **[SHA256 Checksum](https://github.com/CyphrRiot/ArchRiot/releases/download/v2.27/archriot.sha256)**

Note: _Do not use wget_ for this. Curl is your friend:

```bash
curl -L -o archriot.iso https://github.com/CyphrRiot/ArchRiot/releases/download/v2.27/archriot.iso
```

2. **💾 Download Ventoy**
    - Download [Ventoy](https://www.ventoy.net/) and create a bootable USB drive
    - _Pro tip: Ventoy lets you boot multiple ISOs from one USB - perfect for testing_

3. **📂 Copy ISO to Ventoy Drive**
    - Copy the ArchRiot ISO file directly to your Ventoy USB drive
    - No flashing needed - just copy the file

4. **⚡ Boot ArchRiot ISO**
    - Boot from USB (disable Secure Boot if your BIOS is being difficult)
    - Select ArchRiot ISO from Ventoy menu

_Note: In some weird, rare cases, dhcpd is not running. In this case, run:_

`sudo systemctl start dhcpcd`

Then run `curl -fsSL https://ArchRiot.org/setup.sh | bash` to continue the final step of the installation.

Note: If There's a wifi issue (after first boot), see: [ISO Install: Wi‑Fi abort fix (iwd PSKs)](#iso-wifi-abort)

```bash
iwctl
device list
station {device} scan
station {device} get-networks
station {device} connect {network}
station {device} show
exit
```

**Perfect for:**

- 🖥️ **Fresh hardware** - New builds, clean slates, virtual machines
- 🚀 **Instant gratification** - Boot → `riot` → Perfect desktop in minutes
- 💀 **System replacement** - When your current setup has disappointed you for the last time
- 🎯 **Zero configuration** - For people who have better things to do than tweak configs

**What you get:**

- Complete Arch Linux + ArchRiot system
- Boot ISO → Run `riot` → Complete guided setup
- Hyprland, themes, apps all pre-configured
- No manual setup required

<a id="iso-wifi-abort"></a>

### Troubleshooting (ISO Install): Wi‑Fi abort at “Wi‑Fi configuration/services not properly set on target”

Cause

- The live ISO uses iwd/iwctl, while the installer’s first‑boot handoff previously expected NetworkManager (NM) connections on the target. If no usable connection is set up on the target, the installer aborts to avoid rebooting without network.

Fix without modifying the ISO
Run these from the archiso prompt after the error (target root is mounted at /mnt):

1. Preferred: copy iwd PSKs and enable iwd + DHCP (no NetworkManager)

```bash
# Enable iwd on target; avoid conflicts
arch-chroot /mnt systemctl enable iwd
arch-chroot /mnt systemctl mask wpa_supplicant

# Ensure DHCP will bring up networking at boot (when not using NM)
arch-chroot /mnt pacman -Sy --noconfirm --needed dhcpcd
arch-chroot /mnt systemctl enable dhcpcd

# Copy the live iwd credentials (PSKs) to the target
arch-chroot /mnt mkdir -p /var/lib/iwd
cp -a /var/lib/iwd/*.psk /mnt/var/lib/iwd/ 2>/dev/null || true
arch-chroot /mnt chmod 600 /var/lib/iwd/*.psk 2>/dev/null || true
```

2. Optional alternative: use NetworkManager on the target (with iwd backend)

```bash
# Install and enable NetworkManager on target
arch-chroot /mnt pacman -Sy --noconfirm --needed networkmanager iwd
arch-chroot /mnt bash -lc 'mkdir -p /etc/NetworkManager/conf.d &&
  printf "[device]\nwifi.backend=iwd\n" > /etc/NetworkManager/conf.d/10-wifi-backend.conf &&
  printf "[connection]\nwifi.powersave=2\n" > /etc/NetworkManager/conf.d/40-wifi-powersave.conf'
arch-chroot /mnt systemctl enable NetworkManager

# Seed NM connections from iwd (best effort)
for f in /var/lib/iwd/*.psk; do
  [ -e "$f" ] || continue
  ssid="$(basename "$f" .psk)"
  pass="$(awk -F= '/^Passphrase=/{print $2}' "$f")"
  if [ -n "$pass" ]; then
    arch-chroot /mnt nmcli c add type wifi ifname "*" con-name "$ssid" ssid "$ssid"
    arch-chroot /mnt nmcli c modify "$ssid" wifi-sec.key-mgmt wpa-psk wifi-sec.psk "$pass"
  else
    echo "No Passphrase in $f; skipping NM import for $ssid (PSK-only)"
  fi
done
arch-chroot /mnt bash -lc 'chmod 600 /etc/NetworkManager/system-connections/*.nmconnection 2>/dev/null || true'
```

3. Reboot and connect

- iwd path: system should auto‑connect from PSKs; if not, run “iwctl” and connect to your SSID.
- NM path: use “nmtui” if auto‑connect doesn’t happen.

Quick workaround

- Use a wired (Ethernet) connection during install; the safety gate will pass with a detected wired link.

Note on Secure Boot

- This specific Wi‑Fi abort is not caused by Secure Boot. Some systems still require disabling Secure Boot for other reasons; consult your vendor’s guide if needed.

---

**Security Note:** Your system remains secure through LUKS disk encryption and screen lock. Passwordless sudo is standard for automated system installations and doesn't compromise security when disk encryption is properly configured.

#### 🚀 One-Line Install or Upgrade

> Warning: On multi‑monitor setups, upgrading may temporarily stop Brave while displays are being reconciled.
>
> Mitigation:
>
> - Prefer upgrading with Brave closed.
> - If Brave is affected, safely recover with:
>     - `~/.local/share/archriot/install/archriot --displays-enforce`
>     - `~/.local/share/archriot/install/archriot --waybar-reload`
> - Displays are also re‑applied automatically at next login (login‑time enforcement).

**The only command you need to remember:**

```bash
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

![ArchRiot Upgrade Demo](config/images/upgrade.gif)

This downloads and runs our bulletproof Go binary installer with intelligent YAML configuration. Upgrading is exactly the same command - because simplicity is the ultimate sophistication.

**What happens:** Automatic package installation, configuration deployment, and system setup with complete rollback capability if anything goes wrong.

<a id="master-your-archriot-desktop"></a>

## ⌨️ Master Your ArchRiot Desktop

_The keyboard shortcuts that'll make you wonder how you ever used a mouse_

### 🎯 Getting Started

_Your gateway to ArchRiot mastery - memorize these first_

| Keybinding       | Action                                     |
| ---------------- | ------------------------------------------ |
| `SUPER + H`      | Show HELP - Your lifeline when lost        |
| `SUPER + D`      | App launcher - Find anything instantly     |
| `SUPER + RETURN` | Terminal - Your command center (Ghostty)   |
| `SUPER + L`      | Lock screen - Beautiful CypherRiot styling |
| `SUPER + ESCAPE` | Power menu - Sleep, restart, or shutdown   |

### 🪟 Window Management (Most Used)

_Tiling mastery - where ArchRiot really shines_

| Keybinding                   | Action                                          |
| ---------------------------- | ----------------------------------------------- |
| `SUPER + W` or `SUPER + Q`   | Close window - Goodbye forever                  |
| `SUPER + V`                  | Toggle floating - Break free                    |
| `SUPER + J`                  | Toggle split - Reorganize space                 |
| `SUPER + Arrow Keys`         | Move focus - Navigate like a pro                |
| `SUPER + SHIFT + Arrow Keys` | Swap windows - Rearrange perfection             |
| `SUPER + CTRL + Arrow Keys`  | Smart movement - Let Hyprland think             |
| `SUPER + SHIFT + TAB`        | Fix off-screen windows - Manual rescue tool     |
| `SUPER + 1-4`                | Switch workspace - Your digital rooms           |
| `SUPER + SHIFT + 1-4`        | Move window to workspace - Relocate             |
| `SUPER + SHIFT + RETURN`     | Floating terminal - When you need overlay power |

### 💻 Core Applications

_The tools that matter, one keystroke away_

| Keybinding               | Action                                                |
| ------------------------ | ----------------------------------------------------- |
| `SUPER + F`              | File manager - Navigate your digital kingdom (Thunar) |
| `SUPER + B`              | Browser - Privacy-focused web (Brave)                 |
| `SUPER + Z`              | Code editor - Modern development (Zed)                |
| `SUPER + N`              | Text editor - Power user paradise (Neovim)            |
| `SUPER + O`              | Simple text editor - Clean and fast (GNOME)           |
| `SUPER + T`              | System monitor - See everything (btop)                |
| `SUPER + SHIFT + RETURN` | Floating terminal - Overlay power mode                |

### 💬 Communication & Social

_Stay connected without selling your soul to data miners_

| Keybinding  | Action                                                |
| ----------- | ----------------------------------------------------- |
| `SUPER + E` | Email - Encrypted and private (Proton Mail)           |
| `SUPER + G` | Telegram messenger - Focus-or-launch (native/Flatpak) |
| `SUPER + S` | Signal messenger - Smart Signal integration           |
| `SUPER + M` | Google Messages - When you must use the machine       |
| `SUPER + K` | Google Keep - Quick notes and lists                   |
| `SUPER + X` | X/Twitter                                             |

Note: `SUPER + SHIFT + K` — Stabilize session: relaunches Waybar, restarts hypridle, and enforces displays (no full reload).

### 📸 Screenshots & Recording

_Capture your ArchRiot greatness and share it with the world_

| Keybinding          | Action                                                  |
| ------------------- | ------------------------------------------------------- |
| `SUPER + SHIFT + S` | Region screenshot - Select exactly what matters         |
| `SUPER + SHIFT + W` | Window screenshot - Perfect app captures                |
| `SUPER + SHIFT + F` | Full screen - Show off your entire desktop              |
| `Kooha`             | Screen recorder - Make tutorials (launch via `SUPER+D`) |

### 🎨 Wallpaper Management & Dynamic Theming

_Keep your desktop fresh with ArchRiot's intelligent wallpaper system and automatic color theming_

| Keybinding             | Action                                    |
| ---------------------- | ----------------------------------------- |
| `SUPER + CTRL + SPACE` | Cycle backgrounds - Fresh vibes on demand |

**🌈 Dynamic Color Theming:**

- **Optional intelligent theming** - Enable in Control Panel to extract colors from your wallpaper
- **Real-time waybar updates** - When enabled, workspace colors, CPU indicators, and accents change instantly to match your wallpaper
- **Toggle control** - Enable/disable dynamic theming in the Control Panel or keep classic CypherRiot colors
- **Smart fallback** - Always reverts to beautiful CypherRiot theme when disabled
- **⚠️ Upgrade note** - System upgrades will reset dynamic theming to CypherRiot defaults (simply re-enable in Control Panel)

**🎛️ Control Panel Magic:**

- Launch ArchRiot Control Panel for drag-and-drop wallpaper management
- Dynamic theming toggle - Enable wallpaper-based colors or stick with CypherRiot classics
- System wallpapers (1-15) + your custom collection (U1, U2, etc.) perfectly organized
- Changes apply instantly - no restarts, no waiting, just beauty

**🖼️ Custom Wallpaper Power Moves:**

```bash
# Pro method: Use the Control Panel's file chooser
# Quick method: Drop files directly into ~/.config/archriot/backgrounds/

# Clean house - remove specific wallpaper
rm ~/.config/archriot/backgrounds/user_01.jpg

# Nuclear option - start fresh
rm ~/.config/archriot/backgrounds/*
```

### ⚙️ System Management

_Keep your ArchRiot system running like a well-oiled machine_

```bash
archriot                             # ArchRiot's intelligent installer/upgrade engine
                                     # Automatically detects if upgrade is needed
                                     # Prompts for confirmation before proceeding

migrate                              # Backup/restore wizard - your insurance policy
sudo systemctl reboot                # Fresh start - sometimes you need it
sudo systemctl poweroff              # Graceful shutdown - not a crash
```

### 🎵 Audio & Media

_Hardware media keys that actually work - imagine that!_

```bash
XF86AudioRaiseVolume                 # Volume up - with gorgeous overlay
XF86AudioLowerVolume                 # Volume down - smooth as silk
XF86AudioMute                        # Toggle mute - instant feedback
XF86AudioMicMute                     # Microphone toggle - privacy at a keystroke
XF86AudioPlay/Pause                  # Media control - works with everything

# Beautiful volume overlay appears instantly with progress bar
# These are your actual hardware keys working the way they should
```

### 📊 Waybar Controls (Status Bar)

_Your desktop's mission control - everything you need at a glance_

```bash
Click tomato timer                   # Pomodoro focus mode - stay productive
Double-click tomato timer            # Reset to 25:00 - fresh start
Click network icon                   # Network manager - connect anywhere
Click volume icon                    # Audio settings - fine-tune your sound
Click battery icon                   # Power management - stay charged
```

### Lid/Idle policy (laptops, docked vs undocked)

- Docked (external display connected): closing the lid does not suspend; system does not idle-suspend.
- Undocked: closing the lid suspends; idle suspend occurs at 30 minutes by default.

Implemented via:

- `/etc/systemd/logind.conf.d/10-docked-ignore-lid.conf`:
    - `HandleLidSwitchDocked=ignore`, `HandleLidSwitch=suspend`, `HandleLidSwitchExternalPower=suspend`
- `/etc/systemd/logind.conf.d/20-idle-ignore.conf`:
    - `IdleAction=ignore` (logind does not idle-suspend)
- hypridle:
    - 10 min: `on-timeout = lock` (reliable hyprlock trigger)
    - 30 min: `on-timeout = $HOME/.local/share/archriot/install/archriot --suspend-if-undocked || $HOME/.local/bin/suspend-if-undocked.sh` (suspends only when undocked)
- kanshi autostart (hotplug profiles) for dock/undock screen management

Note: We never restart systemd-logind during install/upgrade; drop-ins take effect after reboot or a manual restart you perform later.

## 📋 Evolution Log

_The relentless march toward Linux perfection_

**🔥 Current Release:** v4.2.1 - Docs & UX polish: Brave wrapper consolidation, Waybar logs/reload guidance, fractional scaling notes, Control Panel sizing, Thunar opacity

**🚀 Recent Milestones:**

- **v2.7.4:** DPMS wake fixes and system optimization
- **v2.7.3:** Enhanced theme consistency and bug fixes
- **v2.7.2:** Critical yay installation + dependency fixes
- **v2.7.1:** Hardware module path issue fixes
- **v2.7.0:** Bulletproof fail-fast error handling - because failure is not an option
- **v2.6.14:** Waybar 3-state + package safety restoration
- **v2.6.13:** Bulletproof setup.sh + Control Panel UX improvements

```bash
# View all version changes
git log --grep="FIX" --oneline
```

<a id="archriot-control-panel"></a>

## 🎛️ ArchRiot Control Panel

**The command center that makes Linux actually user-friendly**

ArchRiot's Control Panel isn't just another settings app - it's the missing piece that makes advanced Linux features accessible to humans. Built with modern GTK4, it's fast, beautiful, and actually makes sense.

**Launch it:** `SUPER+C` or run `archriot-control-panel` from anywhere

![ArchRiot Control Panel](config/images/control-panel.png)

### 🎛️ **Features**

- **🍅 Pomodoro Timer** - Waybar-integrated productivity timer with 5-60 minute intervals
- **💡 Blue Light Filter** - Real-time screen temperature control (2500K-5000K) via hyprsunset
- **🛡️ Mullvad VPN** - Account management with privacy controls and auto-connect
- **🔊 Audio System** - Safe mute/unmute controls without breaking services
- **📷 Camera Control** - Device permissions, resolution settings, and live preview testing
- **🖥️ Display Settings** - Monitor resolution and scaling with live preview
- **🔋 Power Management** - Battery profiles (Power Saver, Balanced, Performance)

### 🛡️ **Privacy & Safety**

- **Account Privacy** - Sensitive information hidden by default with show/hide toggle
- **Safe Controls** - Mute instead of killing services, permissions instead of breaking :
- **Live Preview** - Real-time system changes with "Exit without Saving" option
- **Educational Content** - "Learn More" dialogs with comprehensive feature explanations

### 🎨 **Technical Excellence**

_Because you deserve software that doesn't suck_

- **GTK4 Application** - Cutting-edge interface with gorgeous CypherRiot theming
- **Real-time Integration** - Watch changes happen instantly - no "apply" buttons needed
- **Bulletproof Persistence** - Your settings survive reboots, updates, and system chaos
- **Modular Architecture** - Built right from day one - extensible and maintainable

## 💾 **Built-in Backup & Recovery with Migrate**

**Your perfect ArchRiot setup is precious - protect it like the digital treasure it is**

ArchRiot automatically installs and integrates **[Migrate](https://github.com/CyphrRiot/Migrate)** - our battle-tested backup and recovery system. This isn't some afterthought tool - it's your insurance policy against hardware failures, user errors, and the inevitable "what did I just delete?" moments.

_Because spending weeks recreating your perfect setup is a special kind of hell._

### 🛡️ **Why Migrate Matters**

- **Complete System Backup** - Every dotfile, every tweak, every perfect configuration preserved
- **Interactive TUI** - Gorgeous terminal interface that makes complex operations simple
- **Live System Recovery** - Restore without nuking your system - keep working while fixing
- **Cross-Installation Migration** - Clone your setup to new machines in minutes
- **Zero Maintenance** - Updates automatically with ArchRiot - one less thing to worry about

### 🔥 **Quick Start**

```bash
migrate                              # Launch interactive backup/restore interface
```

**No flags, no complexity** - just run `migrate` and use the intuitive menu to backup or restore your entire ArchRiot setup in minutes!

[![Migrate Tool](https://img.shields.io/badge/Backup_Tool-Migrate-6B46C1?style=for-the-badge&logo=github&logoColor=white)](https://github.com/CyphrRiot/Migrate)

<a id="key-customizations"></a>

## 🎯 Key Customizations

### 🔧 **Core System Changes**

- **Terminal:** Ghostty (replaces Kitty) with 90% transparency and dark theme
- **Browser:** Brave (replaces Chromium) with native Wayland support
- **File Manager:** Thunar (replaces Nautilus) with comprehensive dark theming
- **Shell:** Fish as default (replaces Bash) with proper PATH configuration
- **Theme:** CypherRiot integrated as unified theme system
- **Code Editor:** Zed (Wayland) + Neovim with proper theme integration
- **Applications:** All major apps now run native Wayland (no more XWayland issues)
- **Migrate Backup Tool:** CyphrRiot's comprehensive system backup/restore solution (built-in)
- **Memory Optimization:** Intelligent memory management that actually works
- **Blue Light Filter:** hyprsunset at 3500K for reduced eye strain (configurable)
- **GTK Theming:** Dark theme everywhere - no more jarring white dialogs
- **DPI Scaling:** Fixed scaling issues for consistent UI across all applications

#### 🧠 **Memory Management Fix**

Linux's default memory management is **aggressively stupid** about caching. The kernel will happily consume 90%+ of your RAM for file caches, then struggle to free it when applications actually need memory.

**ArchRiot's Solution:** Comprehensive memory management tuning that provides:

**Core Improvements:**

- **Smart Caching** - Reserves 1GB RAM, reduces aggressive file system caching
- **Minimal Swapping** - 10% swappiness (vs 60% default) keeps everything in RAM
- **Lag-Free Writing** - 5% dirty page limit prevents massive write bursts
- **Background Cleanup** - 2% background writeback for smooth performance

**Advanced Protection:**

- **Memory Overcommit Control** - Prevents dangerous memory allocation that causes crashes
- **Proactive Defragmentation** - Reduces memory fragmentation for better allocation
- **Smart OOM Killer** - Kills problematic processes, not random system services
- **Enhanced Responsiveness** - Optimized dirty page intervals and memory bandwidth

**Real-World Impact:**

- **No more lag spikes** when opening applications or switching windows
- **Better responsiveness** under heavy memory pressure (tested with 75%+ RAM usage)
- **Reduced swap usage** with intelligent RAM utilization
- **System stability** under extreme load - no freezes or crashes

**Result:** Your system stays fast and responsive even when running multiple applications, compiling code, or under extreme stress testing.

### 📱 **Advanced Waybar Integration**

ArchRiot includes a highly customized Waybar (status bar) with comprehensive system integration:

**Built-in Modules:**

- **🍅 Tomato Timer** - Built-in Pomodoro timer (idle/running/break/finished)
- **🛡️ Mullvad VPN Status** - Real-time VPN connection status with location display
- **📊 System Monitoring** - CPU aggregate usage, accurate memory monitoring
- **📊 Visual System Metrics** - Temperature, CPU, memory, and volume
    - shown as intuitive bar indicators (▁ ▂ ▃ ▄ ▅ ▆ ▇ █)
- **🎤 Microphone Control** - Visual mic status with one-click toggle
- **📶 Network Management** - WiFi status with nmtui integration
- **🔊 Audio Controls** - Volume display with hardware key integration

**Technical Improvements:**

- **CSS Parser Fixed** - Eliminated all !important declarations causing waybar errors
- **Custom Separators** - Clean, organized module layout for better readability
- **Transparency System** - Consistent 90-98% opacity across all applications
- **Font Optimization** - Improved date format (Sunday • July 13 • 01:49 PM)
- **Error-free Operation** - All modules validated and tested for reliability

### 📱 **Clean Web Applications**

- **Proton Mail** (SUPER+E / XF86Mail) - Privacy-focused email in floating window
- **Telegram** (SUPER+G) - Private messaging app
- **Signal** (SUPER+S) - Private messaging app
- **Google Messages** (SUPER+M) - Web-based messaging in floating window
- **X/Twitter** (SUPER+X) - Social platform in floating window
- **GitHub** - Development platform with proper icons from homarr-labs

### 🪟 **Responsive Window Management**

- **Percentage-based sizing** - Windows scale properly across different screen resolutions (1080p, 1440p, 4K, ultrawide)
- **Smart centering** - All floating windows automatically center regardless of monitor size
- **Cross-resolution compatibility** - No hardcoded pixel positions, works on any display setup
- **Optimized app windows**:
    - **X/Twitter**: `40% x 90%` (mobile-style layout with responsive height)
    - **Proton Mail**: `45% x 80%` (perfect email reading dimensions)
    - **Google Messages**: `40% x 85%` (comfortable messaging interface)
    - **Signal**: `40% x 80%` (maintains native desktop experience)
- **Future-proof design** - Window rules adapt automatically to new monitor configurations

### 🎮 **GPU Support**

ArchRiot automatically detects and installs optimal drivers for all major GPUs:

- **NVIDIA**: Proprietary drivers with Wayland and hardware acceleration
- **AMD/Radeon**: Open-source Mesa drivers with Vulkan support
- **Intel**: Mesa drivers including Intel Arc support

All GPUs get proper Wayland integration and hardware video acceleration for optimal performance.

**Performance Features:**

- **Hardware acceleration** - Video playback, compositing effects, and application rendering
- **Wayland native support** - No XWayland compatibility issues
- **Automatic driver selection** - No manual configuration required
- **Vulkan support** - Modern graphics API for gaming and development
- **Multi-monitor optimization** - Proper scaling and display management

### ⌨️ **Enhanced Keybindings & Productivity**

- **SUPER+D** = **SUPER+SPACE** (Unified app launcher)
- **Left-click Arch icon** - nwg-drawer app grid
- **Right-click Arch icon** - fuzzel app launcher
- **XF86Mail** - Floating Proton Mail window
- **SUPER+SHIFT+S** - Region screenshot (primary)
- **SUPER+SHIFT+W** - Window screenshot
- **SUPER+SHIFT+F** - Full screen screenshot
- **Key repeat enabled** (40 rate, 600 delay for responsive typing)
- **All media keys** - Volume, brightness, playback controls

### Keybindings Help (SUPER+SHIFT+H)

- Local and dynamic: The help window is generated from your Hyprland config (inline comments on bind lines).
- Sources:
    - ~/.config/hypr/keybindings.conf (preferred)
    - ~/.config/hypr/hyprland.conf (fallback)
- How to add descriptions: append a comment after a bind line, e.g. bind = $mod, D, exec, fuzzel # App launcher
- Launch: SUPER+SHIFT+H opens a compact web app window (Brave app mode) showing current binds.
- Tip: Edit your binds and press SUPER+SHIFT+H again to regenerate.

### Workspace Styles (Waybar)

You can switch the workspace module style in your Waybar config. Edit ~/.config/waybar/config and change the module name from:

- hyprland/workspaces#rw (default: “Rewrite” style with per-window icons)

to any of these:

- hyprland/workspaces Circles (● ◉ ○)
- hyprland/workspaces#numbers Plain 1..4
- hyprland/workspaces#roman I II III IV
- hyprland/workspaces#kanji 一 二 三 四
- hyprland/workspaces#pacman Pac‑Man themed
- hyprland/workspaces#cam Uno/Due/Tre/Quattro
- hyprland/workspaces#4 Numbers + per‑workspace icons

Notes:

- Persistent workspaces default to 1–4; 5–6 appear automatically when occupied. To change persistence, adjust persistent-workspaces in config/waybar/ModulesWorkspaces.
- Reload the bar safely after changes:
  archriot --waybar-reload

### 🎨 **Document & Media Handling**

- **Gnome Text Editor** - Default for text/markdown files (clean, modern text editing with CypherRiot theme)
- **Papers** - Default PDF viewer (GNOME's modern document viewer)
- **MPV** - Video playback with optimal performance
- **Better waybar network** - nmtui instead of impala for reliable WiFi management
- **Screenshot tools** - grim/slurp/hyprshot integration for all capture needs
- **Screen recording** - Kooha for simple GUI-based screen recording

### 🚫 **Removed Bloat & Corporate Apps**

- **Removed 37signals/Basecamp tools** - Hey, Basecamp web apps
- **Removed corporate social** - Discord, proprietary messaging
- **Removed heavy productivity** - Obsidian, LibreOffice, OBS Studio, KDEnlive, Pinta
- **Removed proprietary services** - 1Password, Typora, Dropbox, Spotify, Zoom
- **Removed entertainment** - YouTube webapp, WhatsApp webapp

<a id="system-management"></a>

## 🔄 System Management

### Updates

**System Updates:**

```bash
sudo pacman -Syu                     # Standard Arch Linux system update
yay -Syu                             # Update AUR packages
```

**ArchRiot Updates:**

```bash
curl -fsSL https://ArchRiot.org/setup.sh | bash    # Update ArchRiot (same as install)
```

**Automatic Update Notifications**: ArchRiot automatically checks for updates every 4 hours and shows a notification dialog when newer versions are available. You can install updates, ignore notifications, or simply close the dialog.

Update dialog options:

- Install from cURL — uses the maintainer’s install/packages.yaml (may reintroduce packages you previously removed).
- Install from Local — uses your local ~/.local/share/archriot/install/packages.yaml and respects packages you’ve removed.

When to choose which:

- Local: you’ve removed packages you don’t want reintroduced, you’ve customized modules locally, or you want the safest upgrade that preserves your current set.
- cURL: you want to align to the maintainer’s latest package set, you’re repairing a broken local state, or you prefer the canonical defaults after larger changes.

Extras:

- View Diff — shows a unified diff between your local packages.yaml and the maintainer’s, so you can preview changes before upgrading.
- Local safety check — a smoke test runs before “Local” upgrades and will block if it detects packages that would be reintroduced. You can optionally allowlist specific packages via ~/.config/archriot/upgrade-allowlist.txt.
- Bluesky (optional) — The Bluesky PWA is not installed by default. To enable it:

    ```bash
    mkdir -p ~/.local/share/applications
    cp ~/.local/share/archriot/config/applications/xtras/Bluesky.desktop ~/.local/share/applications/
    update-desktop-database ~/.local/share/applications
    ```

    It will then appear in Fuzzel as "Bluesky."

Brave wrapper and per-user flags

- Launchers: Only two visible entries: Brave and Brave (Private). The internal archriot-brave.desktop is hidden and used for default browser mapping; it won’t appear in launchers.
- Wrapper: The archriot-brave launcher enforces Wayland when available and enables GPU rasterization by default for stability and performance.
- Per-user overrides: Create ~/.config/archriot/brave-flags.conf (one flag per line) to add or override flags. Examples:
  --disable-gpu
  --ozone-platform=x11
- Precedence: Command-line args > brave-flags.conf > built-in defaults. Duplicate flags are resolved by key (higher precedence wins).
- PATH: The wrapper is also available at ~/.local/bin/archriot-brave for consistent shell access.
- Defaults applied automatically (when on Wayland):
  --ozone-platform=wayland
  --enable-features=UseOzonePlatform
  --enable-gpu-rasterization

<div align="center">
<img src="config/images/upgrade.png" alt="ArchRiot Update Dialog" width="600">
<br><em>Waybar update notifications: 󰚰 (new), 󱧘 (seen), - (up-to-date) with one-click upgrade dialog</em>
</div>

The ArchRiot updater downloads the latest YAML configuration and pre-built binary, then intelligently applies only the changes needed. The YAML-based system ensures **atomic updates** with proper dependency resolution - no partial failures or broken states like traditional shell script updaters.

### Portals and Screencast Troubleshooting

- Confirm active portals:
    - systemctl --user status xdg-desktop-portal xdg-desktop-portal-hyprland
    - ps aux | grep xdg-desktop-portal
- Verify preferred portal selection:
    - Check ~/.config/xdg-desktop-portal/portals.conf (or /etc/xdg-desktop-portal/portals.conf)
    - Ensure default=hyprland;gtk; and hyprland is selected for ScreenCast/Screenshot/Clipboard
- Reset portals (fix stale processes):
    - systemctl --user restart xdg-desktop-portal xdg-desktop-portal-hyprland
    - If needed: pkill xdg-desktop-portal; pkill xdg-desktop-portal-hyprland; then restart the services
- Quick screencast tests:
    - Kooha: kooha (record a short clip). You should see a Wayland screencast prompt and successful capture on clean installs.
    - OBS Studio: In Sources, click + and ensure "Screen Capture (PipeWire)" is available. Select it and verify monitor/window capture works. If missing, restart portals (see above) or relaunch your Hyprland session.
- Common issues and fixes:
    - Make sure xdg-desktop-portal-hyprland is installed and up to date
    - Avoid conflicting portals; ensure hyprland is first in the default list
    - If prompts do not appear, restart the user session or relaunch Hyprland

### Waybar Portability Troubleshooting

- Reload without duplicates:
    - Send a config reload: archriot --waybar-reload (sends SIGUSR2; auto-restarts Waybar on crash)
    - If the bar is unresponsive, restart cleanly: killall waybar; archriot --waybar-launch
- Check single-instance status and logs:
    - Status: archriot --waybar-status
    - Logs: tail -n 200 ~/.cache/archriot/runtime.log
- Idle/Lock screen not triggering? See "Hypridle: Lock Not Triggering (exec syntax)" below
- Dynamic detection tips:
    - Network: modules/scripts should auto-detect active interfaces; ensure your interface is up (ip link), and NetworkManager is running if you rely on it
    - Temperature: scripts try coretemp/k10temp/zenpower first, then thermal zones; install lm_sensors and run sensors-detect if temps are missing
    - Battery: laptop-only; verify /sys/class/power_supply/BAT\* exists
- Common fixes:
    - After theme/config changes, prefer SIGUSR2 reloads (pkill -SIGUSR2 waybar) over full restarts for stability
    - If modules look stale, restart portals (see Portals Troubleshooting) and then Waybar
    - Validate your configs exist under ~/.config/waybar and that custom scripts are executable

### WiFi drops on lock/idle

If your WiFi drops when the screen locks or during idle:

1. Quick diagnostic (read-only)

- Run:

```/dev/null/commands.sh#L1-3
~/.local/share/archriot/install/archriot --wifi-powersave-check
```

- This prints:
    - NetworkManager drop-in status: expects wifi.powersave=2 at /etc/NetworkManager/conf.d/40-wifi-powersave.conf
    - Runtime power save on your wireless interface (e.g., wlan0): expects power_save off

2. Fix NetworkManager drop-in (persistent)

- If the drop-in is missing or not set to 2:

```/dev/null/commands.sh#L1-2
echo -e "[connection]\nwifi.powersave=2" | sudo tee /etc/NetworkManager/conf.d/40-wifi-powersave.conf
sudo systemctl reload NetworkManager
```

3. Turn off runtime WiFi power save (temporary)

- If the diagnostic shows “power_save on” for your interface:

```/dev/null/commands.sh#L1-1
sudo iw dev <iface> set power_save off
```

- Replace <iface> with your wireless interface (e.g., wlan0). This helps avoid WiFi power saving while the session is locked.

4. Re-run the check

- Verify both the drop-in and runtime state:

```/dev/null/commands.sh#L1-1
~/.local/share/archriot/install/archriot --wifi-powersave-check
```

### Backup & Restore

```bash
migrate
```

**🎯 Migrate** is a separate project by Cypher Riot that gets automatically installed during ArchRiot setup. It's a TUI (Text User Interface) with **no command-line options**. Simply run the command and use the interactive menu to:

- Create comprehensive system backups
- Restore from previous backups
- Migrate configurations between installations
- Preserve all your customizations

**Integration Details:** ArchRiot automatically downloads and installs the latest version of Migrate from [CypherRiot/Migrate](https://github.com/CyphrRiot/Migrate) during installation, ensuring you always have the most current backup capabilities without any manual setup.

## 🎨 CypherRiot Theme System

There is one theme: **CypherRiot**, a beautiful Neo Tokyo Dark inspired theme. If you don't like it, the theme files are at `~/.local/share/archriot/config/` and can be edited.

**Visual Design:**

- **Style:** Custom Neo Tokyo Dark aesthetic with dark elegance
- **Color Palette:** Deep purples, electric blues, and charcoal backgrounds
- **Integration:** Complete system theming (waybar, hyprlock, fuzzel, terminals, applications)
- **Backgrounds:** 23 riot-themed wallpapers for dynamic cycling

**System Integration:**

- **Window Manager:** CypherRiot colors in Hyprland decorations and borders
- **Status Bar:** Custom waybar with CypherRiot purple accents and consistent styling
- **Lock Screen:** Beautiful hyprlock with CypherRiot theme and system status
- **Applications:** Unified theming across GTK, terminal, and desktop applications

### Wallpaper Management

- **Instant application:** Background changes apply immediately
- **Persistent settings:** Background preferences survive reboots and updates

**Background System:**

- **CypherRiot collection:** 23 riot-themed wallpapers included
- **Easy cycling:** Use `SUPER + CTRL + SPACE` to cycle through backgrounds
- **Dynamic switching:** Script-based background rotation for variety
- **High quality:** Curated wallpapers optimized for the CypherRiot aesthetic

**Advanced Customization:**

- **Background location:** `~/.local/share/archriot/backgrounds/` (consolidated directory)
- **Custom backgrounds:** Add your own wallpapers to the backgrounds directory
- **Script integration:** Background cycling integrates with waybar and system status
- **Application consistency:** CypherRiot theme provides unified styling across all applications

## Understanding the YAML Configuration

**For developers and power users who want to customize the system before installation:**

```bash
git clone https://github.com/CyphrRiot/ArchRiot.git ~/.local/share/archriot
# Edit install/packages.yaml to customize
~/.local/share/archriot/install/archriot
```

### YAML Architecture

ArchRiot uses a modern YAML-based configuration system that replaces traditional shell scripts. The entire system is defined in `install/packages.yaml`, which contains:

**Structure:**

- **Categories** (core, desktop, development, system, media)
- **Modules** within each category (base, hyprland, tools, etc.)
- **Each module defines:**
    - `packages:` - List of packages to install
    - `configs:` - Configuration files to deploy with patterns and targets
    - `commands:` - Post-installation commands to run
    - `depends:` - Dependencies on other modules
    - `type:` - Installation type (pacman, yay, flatpak, etc.)

**Example module:**

```yaml
desktop:
    hyprland:
        packages: [hyprland, waybar, wofi, hyprpaper]
        configs:
            - pattern: "config/hypr/*"
              target: "~/.config/hypr/"
              preserve_if_exists: [monitors.conf]
            - pattern: "config/waybar/*"
              target: "~/.config/waybar/"
        commands: ["systemctl --user enable hyprland"]
        depends: [core.base]
        type: pacman
```

This YAML system provides clean separation of packages, configurations, and commands while maintaining full dependency resolution and proper installation ordering.

### Configuration Preservation

ArchRiot includes an intelligent preservation system via `install/preserve.yaml` that maintains your personal customizations during upgrades and reinstalls. This ensures your keyboard layout, application preferences, and other settings survive system updates.

**Preserved Settings:**

- **Keyboard Configuration** - Layout (us, fr, de), variant, and model
- **Default Applications** - Browser, terminal, file manager, and dock preferences
- **User Customizations** - Any settings you've modified in Hyprland configs

**How It Works:**

1. **Detection** - Scans existing configs for user-customizable settings
2. **Backup** - Creates timestamped backups in `~/.cache/archriot/`
3. **Extraction** - Pulls out your personal settings using defined patterns
4. **Restoration** - Optionally applies saved settings to new configs

**Example preserve.yaml entry:**

```yaml
user_customizable_settings:
    - name: "kb_layout"
      description: "Keyboard layout (us, fr, de, etc.)"
      pattern: "kb_layout"

    - name: "browser"
      description: "Default browser application"
      pattern: "$browser"
```

The system prompts during installation: **"Restore your hyprland modifications?"** allowing you to choose whether to apply your saved preferences to the fresh configuration.

## ⚡ ArchRiot At a Glance

_The elevator pitch - everything that makes ArchRiot special in one place_

**🪟 Wayland Excellence:** Hyprland compositor with smooth animations, intelligent tiling, and zero XWayland compromises

**💻 Developer Paradise:** Fish shell, Zed editor, Neovim, modern CLI tools (lsd, ripgrep), and containers that just work

**🛡️ Privacy Arsenal:** Brave browser, Proton Mail, Mullvad VPN, Signal messaging, and Feather Wallet - all with native Wayland support

**🎨 Visual Perfection:** CypherRiot dark themes, optional blue light filtering, and beautiful interfaces that don't hurt your eyes at 3AM

**⚡ Performance Tuned:** Intelligent memory management, hardware acceleration, and optimized audio stack for lag-free computing

## 🔀 Differences from Omarchy

ArchRiot was (once) a **heavily customized fork** with these key distinctions:

### Core Philosophy

- **Privacy-first approach** - Proton Mail, Brave browser, Signal messaging vs. corporate alternatives
- **Developer-focused** - Zed editor, modern CLI tools, Fish shell, comprehensive dev environment
- **Performance over bloat** - Lightweight applications, intelligent memory management, native Wayland
- **Clean aesthetics** - CypherRiot theme, consistent dark mode, minimal distractions

### Major Technical Differences

- **Built-in backup system** - Integrated Migrate tool for complete system backup/restore
- **Enhanced window management** - Responsive percentage-based sizing across all resolutions
- **Comprehensive GPU support** - Automatic detection and optimization for NVIDIA, AMD, Intel
- **Advanced Waybar integration** - Custom modules, Pomodoro timer, VPN status, system monitoring
- **Modern application stack** - Ghostty terminal, Brave browser, native Wayland applications
- **Intelligent system tuning** - Memory management fixes, blue light filtering, DPI scaling

ArchRiot transforms Omarchy from a general productivity setup into a specialized development and privacy-focused environment.

## 🔍 Installation Verification System

ArchRiot includes a comprehensive verification system to ensure everything is working correctly:

### What the Installer Actually Does

The ArchRiot installer validates everything automatically as it runs. You don't need separate validation because the installer IS the validation system:

**Real-time validation during installation:**

- YAML configuration integrity and module dependencies
- Essential packages (yay, git, base-devel)
- Desktop environment (Hyprland, Waybar, fuzzel, mako)
- Configuration file deployment verification
- Applications (terminal, file manager, browser, text editor)
- System services (audio, network, bluetooth)
- Network connectivity and repository accessibility
- Memory, disk space, and system requirements

**If anything fails:** The installer stops immediately with clear error messages and diagnostic information. Fix the issue and re-run - it's completely safe and idempotent.

- **Detailed failure analysis** with specific recommendations
- **Fix suggestions** for failed components

### Verification

Installation logs are in `~/.cache/archriot/install.log` and can be reviewed for the full installation details or any errors during the installation process.

### Expected Defaults

After fresh installation, you should see:

- **Default theme:** CypherRiot (Neo Tokyo Dark aesthetic)
- **Default background:** riot_zero.png (riot-themed wallpaper)
- **PDF files:** Show proper document icons (not thumbnails)
- **Image files:** Show thumbnail previews in Thunar
- **Waybar:** Running with tomato timer, system stats, and transparent microphone button

## ArchRiot CLI Flags

See Advanced Usage: CLI Flags for the complete list and details.

[Jump to Advanced Usage: CLI Flags](#advanced-usage-archriot-cli-flags)

## Hyprland binds and exec-once (examples)

- Exec-once (Waybar):
    - exec-once = $HOME/.local/share/archriot/install/archriot --waybar-launch
- Reload (Waybar):
    - bind = $mod SHIFT, SPACE, exec, $HOME/.local/share/archriot/install/archriot --waybar-reload
- Wallet:
    - bind = $mod, R, exec, $HOME/.local/share/archriot/install/archriot --wallet
- Pomodoro:
    - bind = $mod, comma, exec, $HOME/.local/share/archriot/install/archriot --pomodoro-click

## Bluesky (optional; opt-in only)

- Bluesky is not installed by default and does not appear in Fuzzel on a clean install.
- If a previous install left a desktop entry behind, the installer now removes it from:
    - `~/.local/share/applications/`, `/usr/local/share/applications/`, `/usr/share/applications/` (both Bluesky.desktop/bluesky.desktop)
    - Then refreshes the desktop database.

Enable Bluesky (opt-in):

```bash
cp ~/.local/share/archriot/config/applications/xtras/Bluesky.desktop ~/.local/share/applications/
update-desktop-database ~/.local/share/applications
```

<a id="advanced-usage-archriot-cli-flags"></a>

## 🧰 Advanced Usage: CLI Flags

These first-class CLI flags replace legacy helper scripts and are used directly in Hyprland keybinds. All of them execute the built binary at:
$HOME/.local/share/archriot/install/archriot

- --waybar-launch
    - Purpose: Single-instance launcher with non-blocking lock and robust logging.
    - Behavior: Ensures exactly one Waybar instance; writes logs to ~/.cache/archriot/runtime.log.
    - Example (Hyprland exec-once): exec-once = $HOME/.local/share/archriot/install/archriot --waybar-launch

- --waybar-reload
    - Purpose: Robust Waybar reload. Uses SIGUSR2 when possible; fallback to controlled restart.
    - Behavior: Dedupe PIDs and avoid duplicate Waybar after reloads/resume.
    - Example (default bind): SUPER+SHIFT+SPACE → archriot --waybar-reload

- --wallet
    - Purpose: Focus-or-launch the configured wallet; avoids duplicate instances.
    - Example (default bind): SUPER+R → archriot --wallet

- --pomodoro-click
    - Purpose: Simulate a click on the Waybar tomato timer (toggle/reset via module behavior).
    - Example (default bind): SUPER+, → archriot --pomodoro-click

- --upgrade-smoketest
    - Purpose: Local upgrade smoketest used by the update dialog.
    - Exit codes: 0 OK; 2 potential reintroductions; 3 unavailable. Honor allowlist at ~/.config/archriot/upgrade-allowlist.txt

- --stay-awake
    - Purpose: Prevent system suspend while a task runs; always detaches so your app isn't tied to the launcher.
    - Usage:
        - archriot --stay-awake curl -L -o file.iso https://example.com/file.iso
        - archriot --stay-awake
    - Notes: Uses systemd-inhibit to block sleep; does not affect screen lock.

- --volume
    - Purpose: Unified audio control for speakers and microphone across PipeWire/PulseAudio.
    - Backend priority: wpctl (PipeWire with starred sink/source) → pamixer → pactl
    - Usage:
        - archriot --volume toggle # Toggle speaker mute
        - archriot --volume inc # Increase volume 5%
        - archriot --volume dec # Decrease volume 5%
        - archriot --volume get # Get current volume percentage

- --suspend-if-undocked
    - Purpose: Suspend only when undocked (no external display connected) and not on AC/USB-PD; intended for idle managers.
    - Behavior: No-ops when an external display or AC/USB-PD is detected; otherwise calls systemctl suspend.
    - Example (Hypridle listener): on-timeout = $HOME/.local/share/archriot/install/archriot --suspend-if-undocked

### Hypridle: Lock Not Triggering (exec syntax)

Symptoms:

- Screen never locks after idle timeout, but `hyprlock` works when launched manually.
- Hypridle verbose logs show attempts to run `exec, /usr/bin/hyprlock` and errors like `exec,: command not found`.

Cause:

- Hyprland config style (`exec, ...`) was mistakenly used inside hypridle.conf.
- Hypridle expects either `lock` (uses `lock_cmd`) or a direct executable path. It does not interpret `exec,`.

Fix:

1. In `~/.config/hypr/hypridle.conf`, replace:

```ini
on-timeout = exec, /usr/bin/hyprlock
```

with either:

```ini
on-timeout = lock
```

or (explicit path):

```ini
on-timeout = /usr/bin/hyprlock
```

2. Use absolute paths in `general {}` and listeners to avoid PATH/env issues under Hyprland:

```ini
general {
  lock_cmd = /usr/bin/hyprlock
  before_sleep_cmd = /usr/bin/loginctl lock-session
  after_sleep_cmd = /usr/bin/hyprctl dispatch dpms on
}
```

3. Restart hypridle:

```bash
pkill hypridle
hyprctl dispatch exec hypridle
```

Notes:

- Waybar’s `idle_inhibitor` module can block idle. Ensure it’s not activated if locks don’t trigger.
- Brightness dim at 5 minutes won’t affect external HDMI/DP displays; consider adding DPMS off/on at lock for a visible cue on external monitors.

- Troubleshooting:
    - Check backend: wpctl status | head -30
    - Test manually: archriot --volume get
    - Waybar uses this for all volume controls (scroll, click, microphone)

- --brightness
    - Purpose: Backlight control for laptop displays.
    - Usage:
        - archriot --brightness up # Increase brightness
        - archriot --brightness down # Decrease brightness
    - Notes: Uses brightnessctl; integrates with Waybar backlight module scroll actions

- --startup-background
    - Purpose: Start wallpaper at login from saved preferences; avoids theme apply during boot to prevent races.
    - Behavior: Reads ~/.config/archriot/background-prefs.json (key: current_background), falls back to riot_01.jpg or the first available image under ~/.local/share/archriot/backgrounds, writes the state file at ~/.config/archriot/.current-background, and restarts swaybg detached.
    - Example (Hyprland exec-once): exec-once = $HOME/.local/share/archriot/install/archriot --startup-background

- --swaybg-next
    - Purpose: Cycle to the next wallpaper and refresh theming if dynamic theming is enabled.
    - Behavior: Iterates images in ~/.local/share/archriot/backgrounds, updates ~/.config/archriot/.current-background, restarts swaybg detached, and triggers best-effort theme refresh.
    - Example (default bind): SUPER+CTRL+SPACE → $HOME/.local/share/archriot/install/archriot --swaybg-next

- --waybar-workspace-click (optional; recommend native on-click "activate" for correct per‑monitor routing)
    - Purpose: Safe Waybar workspace click handler (numeric-only).
    - Usage: $HOME/.local/share/archriot/install/archriot --waybar-workspace-click {name}
    - Notes: Validates numeric workspace and dispatches: hyprctl dispatch workspace {name}
    - Example (custom modules): "on-click": "$HOME/.local/share/archriot/install/archriot --waybar-workspace-click {name} "

- --waybar-cpu
    - Purpose: Aggregate CPU usage meter for Waybar.
    - Behavior: Reads /proc/stat deltas; renders a bar, percentage, and class via JSON.
    - Example (Waybar): "exec": "$HOME/.local/share/archriot/install/archriot --waybar-cpu"

- --waybar-temp
    - Purpose: CPU temperature meter for Waybar.
    - Behavior: Autodetects sensor (hwmon coretemp/k10temp/zenpower → temp1_input; x86_pkg_temp thermal zone; thermal_zone0 fallback), renders bar and class via JSON.
    - Example (Waybar): "exec": "$HOME/.local/share/archriot/install/archriot --waybar-temp"

- --waybar-volume
    - Purpose: Speaker volume meter for Waybar.
    - Backend priority: wpctl (PipeWire) → pamixer → pactl. Renders bar and icon via JSON; shows “audio not ready” gracefully.
    - Example (Waybar): "exec": "$HOME/.local/share/archriot/install/archriot --waybar-volume"

- --waybar-memory
    - Purpose: Memory usage meter for Waybar.
    - Behavior: Computes traditional percentage, shows modern vs traditional in tooltip; renders bar and class via JSON.
    - Example (Waybar): "exec": "$HOME/.local/share/archriot/install/archriot --waybar-memory"

- --stabilize-session
    - Purpose: Session recovery utility. Dedupe Waybar and relaunch a single, managed instance; restart hypridle to ensure idle/lock policy is active.
    - Usage:
        - archriot --stabilize-session
        - archriot --stabilize-session --inhibit # starts a detached sleep inhibitor for long-running work

- --zed
    - Purpose: Focus-or-launch Zed (native > Flatpak) with Wayland-friendly environment; focuses an existing window if present.
    - Example (bind): SUPER+Z → $HOME/.local/share/archriot/install/archriot --zed
    - Example (desktop entry Exec): $HOME/.local/share/archriot/install/archriot --zed %U

- --welcome
    - Purpose: Launch the ArchRiot Welcome window (Python GTK) in a detached manner.
    - Example (Hyprland exec-once): sleep 2 && $HOME/.local/share/archriot/install/archriot --welcome

### Hyprland bind examples (copy/paste)

These are already present by default; use if you need to reapply or test live.

```bash
# Reload Waybar safely
hyprctl keyword bind "$mod SHIFT, SPACE, exec, $HOME/.local/share/archriot/install/archriot --waybar-reload"

# Wallet
hyprctl keyword bind "$mod, R, exec, $HOME/.local/share/archriot/install/archriot --wallet"

# Pomodoro
hyprctl keyword bind "$mod, comma, exec, $HOME/.local/share/archriot/install/archriot --pomodoro-click"

# Telegram (PATH-resolved; shows 2s 'Opening Telegram...' notification)
hyprctl keyword bind "$mod, G, exec, sh -lc 'notify-send -t 2000 "Telegram" "Opening Telegram..." >/dev/null 2>&1; Telegram'"
```

## Brave Wrapper and Handler Mapping

ArchRiot routes all browser launches through a PATH-resolved wrapper: `archriot-brave`.

- Executable resolution:
    - The installer ensures the wrapper is on PATH.
    - Any stale `/usr/local/bin/archriot-brave` is proactively removed during install/upgrade to avoid shadowing.

- Defaults and overrides:
    - Defaults enable Wayland/Ozone and GPU rasterization.
    - User flags file: `~/.config/archriot/brave-flags.conf`
    - Flag precedence: CLI flags > user flags file > wrapper defaults
    - Example (user flags file; one flag per line):
      --ozone-platform-hint=wayland
      --enable-gpu-rasterization

- Verify GPU acceleration:
    - Open: `brave://gpu`
    - Expectation: Most features show “Hardware accelerated.”
    - If debugging issues, you can temporarily launch with a safe mode:
      archriot-brave --disable-gpu
    - Crash on workspace/monitor switch: test with the safe mode above. If stable, keep `--disable-gpu` temporarily and verify GPU drivers (Mesa/NVIDIA) and Wayland flags; report your GPU/driver combo.

- Handler policy and verification:
    - HTTP/HTTPS should resolve to the wrapper via the Brave desktop entry.
    - Verify current defaults:
      xdg-settings get default-web-browser
      xdg-mime query default x-scheme-handler/http
      xdg-mime query default x-scheme-handler/https
    - Expected: ArchRiot’s wrapper desktop entry is the default (`archriot-brave.desktop`). Only “Brave” and “Brave (Private)” should appear as visible handlers in common menus.

Notes:

- This wrapper approach avoids `$HOME` expansion pitfalls and keeps a consistent Wayland configuration by default.
- Use CLI flags for one-off tests; prefer the user flags file for persistent changes.

### Multi‑monitor (optional flags)

Some multi‑monitor Wayland setups benefit from these flags. Add them to ~/.config/archriot/brave-flags.conf (one per line), then restart Brave:

```
--disable-features=WaylandWpColorManagerV1
--enable-features=VaapiVideoDecoder,VaapiIgnoreDriverChecks,Vulkan,DefaultANGLEVulkan,VulkanFromANGLE,UseOzonePlatform
--ozone-platform=wayland
--enable-gpu-rasterization
--enable-zero-copy
--ignore-gpu-blocklist
```

Tip:

- Test incrementally and keep only what you need.
- Flag precedence: command‑line args > brave‑flags.conf > defaults.

## Waybar: Logs, Reloads, and Debugging

See:

- 🧰 Advanced Usage: CLI Flags — waybar commands (launch/reload/status)
- 🔧 Troubleshooting — logs path, SIGUSR2 reloads, dedupe tips, and recovery steps

<a id="troubleshooting"></a>

## 🔧 Troubleshooting

### Display enforcement opt-out

If your display configuration should not be automatically enforced (for example, to prevent your laptop panel from being disabled), you can opt out safely:

- Disable enforcement (no monitor changes will be applied):
    - Create the marker file:
        - touch ~/.config/archriot/disable-display-enforcement

- Re-enable enforcement:
    - Remove the marker file:
        - rm ~/.config/archriot/disable-display-enforcement

Notes:

- The login-time service will no-op when this marker exists.
- You can still run manual commands when needed:
    - ~/.local/share/archriot/install/archriot --displays-enforce
    - ~/.local/share/archriot/install/archriot --waybar-reload

### Lock screen: Weather and Crypto

- Weather (enable/disable)
    - Disable: touch ~/.config/archriot/disable-weather
    - Enable: rm ~/.config/archriot/disable-weather
    - Optional location config (use either file):
        - ~/.config/waybar/weather.conf
        - ~/.config/archriot/weather.conf
        - Put a line: LOCATION="City, CC"

- Crypto (enable/disable)
    - Disable: touch ~/.config/archriot/disable-crypto
    - Enable: rm ~/.config/archriot/disable-crypto
    - Configure coins: ~/.config/archriot/crypto.cfg
        - Space- or comma-separated CoinGecko ids (up to 5), e.g.: bitcoin ethereum solana

- Notes
    - Weather uses stormy; if missing, it simply shows nothing (no errors).

- Hyprlock crypto (holdings/P&L and custom list)

  ArchRiot’s lock screen uses a more powerful crypto script that supports holdings, entry price (to show profit/loss), ordering, and blank-line separators — all defined from a single list.

  Where:
  - Script path: `$HOME/.local/share/archriot/config/bin/hyprlock-crypto.sh`
  - Wired in hyprlock: see `config/hypr/hyprlock.conf` (execs with `ROWML` mode)

  Edit the `PAIRS` list near the top of the script. Use 2‑tuples for price‑only, or 4‑tuples to enable P/L:

  ```python
  PAIRS = [
      ("ZEC", "zcash",   0,   0.0),  # 4‑tuple → shows P/L when entry > 0 (held, entry_price)
      ("XMR", "monero",  0,   0.0),
      ("LTC", "litecoin", 0,   0.0),
      ("",    ""),               # blank separator line (ROWML only)
      ("BTC", "bitcoin"),        # 2‑tuple → price only (no P/L)
      ("ETH", "ethereum"),
  ]
  ```

  Quick rules:
  - Add a token: append a tuple. Use CoinGecko’s id as the second field (from the coin’s URL, e.g., `https://www.coingecko.com/en/coins/zcash` → `zcash`).
  - Remove a token: delete its tuple.
  - Show P/L: switch to a 4‑tuple and set your holdings (`held`) and your average entry price (`entry`). P/L appears only when `entry > 0`.
  - Add a blank separator: insert `("", "")` exactly where you want the blank line (affects only ROWML mode).
  - Alignment: amounts and percents use fixed widths for clean columns on monospace fonts. When entry is 0, P/L is hidden.


    - Crypto uses curl + jq + CoinGecko; if dependencies are missing, it simply shows nothing.

### Multi‑monitor anomalies after install (ABI mismatch)

If you see display glitches or crashes on multi‑monitor setups immediately after an install or upgrade, it’s likely a partial‑upgrade ABI mismatch (e.g., Hyprland/wlroots/portals not in sync). Fix by fully upgrading first, then run the installer again:

```bash
sudo pacman -Sy archlinux-keyring && yay -Syu && yay -Yc && sudo paccache -r
```

Pro tip: run the installer with `--strict-abi` to block installation until compositor/Wayland updates are applied.

### Installer Sync Recovery (pacman db lock / mirrors)

If you see pacman errors like “could not lock database” or “failed to synchronize,” try these steps:

Symptoms:

- could not lock database (db.lck present)
- failed to synchronize packages / failed retrieving file
- temporary failure in name resolution

Steps:

1. Ensure no other package manager is running (pacman/yay/paru).
2. Clear stale lock (safe if no pacman is running):
   sudo rm -f /var/lib/pacman/db.lck
3. Refresh databases:
   sudo pacman -Sy
    # If issues persist, force a full refresh:
    sudo pacman -Syy
4. Retry your install/upgrade and watch for transient mirror/network hiccups.

### Preflight Audit (read-only)

- Run: `~/.local/share/archriot/install/archriot --preflight`
- What it checks (no changes made; safe to run anytime):
    - Config: validates `packages.yaml`
    - Binary path: confirms you’re using `$HOME/.local/share/archriot/install/archriot`
    - Hyprland binds: verifies `SUPER+G` (Telegram) and `SUPER+S` (Signal)
    - Exec-once: ensures Waybar uses `archriot --waybar-launch`
    - Memory tuning: shows opt-in status (does not modify kernel settings)
    - Waybar: shows instance count (does not kill any process); logs at ~/.cache/archriot/runtime.log

### Memory Tuning (Basics by default; Advanced Opt-in)

By default, ArchRiot disables ZRAM to prevent compressed-RAM thrash and lockups under pressure, and applies minimal, safe VM settings system‑wide (vm.zone_reclaim_mode=0 and vm.overcommit_memory=0). To re‑enable ZRAM, remove or edit /etc/systemd/zram-generator.conf (set zram-size to a non‑zero value), and if you use systemd-swap, enable it with: sudo systemctl enable --now systemd-swap.service; then reboot. For advanced, dynamic tuning (including vm.min_free_kbytes ≈ 1% of RAM capped at 256MB, swappiness/dirty ratios), opt‑in by creating the flag file and re‑running install/upgrade. Optional (opt‑in): touch ~/.config/archriot/enable-memory-optimizations and re‑run installer/upgrade; verify /etc/sysctl.d/99-memory-optimization.conf contains the computed vm.min_free_kbytes.

```bash
touch ~/.config/archriot/enable-memory-optimizations
```

Apply now (opt-in):

```bash
sudo cp ~/.local/share/archriot/config/system/99-memory-optimization.conf /etc/sysctl.d/99-memory-optimization.conf
total_kb=$(awk '/MemTotal/ {print $2}' /proc/meminfo); calc=$(awk -v t="$total_kb" 'BEGIN {m=int(t*0.01); if (m > 262144) m=262144; print m}'); sudo sed -i "s/^vm.min_free_kbytes=.*/vm.min_free_kbytes=$calc/" /etc/sysctl.d/99-memory-optimization.conf
sudo sysctl -p /etc/sysctl.d/99-memory-optimization.conf
```

Revert quickly (if anything feels off):

```bash
sudo sed -i 's/^vm.overcommit*memory=.*/vm.overcommit*memory=0/; s/^vm.overcommit_ratio=.*/vm.overcommit_ratio=50/; s/^vm.min_free_kbytes=.*/vm.min_free_kbytes=262144/' /etc/sysctl.d/99-memory-optimization.conf
sudo sysctl -p /etc/sysctl.d/99-memory-optimization.conf
```

### Blinking Cursor Instead of Hyprland

If your system boots to a **blinking cursor** instead of starting Hyprland:

1. **Get to a terminal:** Press `CTRL+ALT+F3`
2. **Login** with your username and password
3. **Re-run the installer** to fix GPU/graphics issues:
    ```bash
    curl -fsSL https://ArchRiot.org/setup.sh | bash
    ```
4. **Reboot** after the script completes

This issue is almost always GPU-related and the installer will detect and fix graphics driver problems automatically.

### ISO Wi‑Fi Troubleshooting (Mediatek MT7921K)

Live fix without rebuilding the ISO. Apply in this order:

- Enable iwd and avoid conflicts:
  sudo systemctl enable --now iwd
  sudo systemctl mask wpa_supplicant

- Unblock radio:
  rfkill list
  sudo rfkill unblock all

- Reload Mediatek driver with ASPM disabled (stability fix):
  sudo modprobe -r mt7921e
  sudo modprobe mt7921e disable_aspm=1
  If the module refuses to unload, disconnect the interface first (via iwctl), then retry.

- Check for missing firmware:
  dmesg | grep -i firmware
  If you see missing Mediatek firmware messages:
    - On an installed system with network:
      sudo pacman -Syu linux-firmware
      sudo modprobe -r mt7921e && sudo modprobe mt7921e disable_aspm=1
    - On live ISO without network:
      Copy required Mediatek firmware files to /lib/firmware/mediatek from another machine/USB, then reload the module.

- Connect using iwd (iwctl):
  iwctl
  device list
  station {device} scan
  station {device} get-networks
  station {device} connect {network}
  station {device} show
  exit

- Persist after install (recommended):
  Create /etc/modprobe.d/mt7921e.conf with:
  options mt7921e disable_aspm=1

- Verify:
    - PCIe device bound to mt7921e:
      lspci -k | grep -A3 -i 7921
    - Radio unblocked:
      rfkill list
    - Link visible/connected:
      iw dev

### Prevent Idle Sleep and Keep Downloads Active

If your system still sleeps after ~10 minutes and interrupts downloads, ensure idle suspend is disabled at logind and that Hypridle isn’t triggering suspend earlier than expected.

1. Verify logind idle policy (no idle suspend)

- Confirm the ArchRiot drop-in exists and disables idle action:

```bash
cat /etc/systemd/logind.conf.d/20-idle-ignore.conf
# Expect:
# [Login]
# IdleAction=ignore
```

- Check effective settings:

```bash
loginctl show-logind -p IdleAction -p IdleActionUSec
# Expect: IdleAction=ignore
```

- If missing, create the drop-in (requires reboot to take effect):

```bash
sudo mkdir -p /etc/systemd/logind.conf.d
printf "[Login]\nIdleAction=ignore\n" | sudo tee /etc/systemd/logind.conf.d/20-idle-ignore.conf
# Reboot to apply, or be aware that restarting logind can disrupt your session:
# sudo systemctl restart systemd-logind
```

2. Hypridle behavior (lock vs suspend)

- ArchRiot’s intended defaults:
    - Lock at 10 minutes
    - Suspend at 30 minutes, only when undocked (native: archriot --suspend-if-undocked; fallback: suspend-if-undocked.sh)
- Check if Hypridle is running:

```bash
pgrep -a hypridle || echo "hypridle not running"
```

- Temporarily stop Hypridle (for testing):

```bash
pkill hypridle
```

- If you want to keep downloads going (no auto-suspend):
    - Edit your Hypridle config (commonly at ~/.config/hypr/hypridle.conf) to remove or comment the suspend on-timeout action that calls suspend-if-undocked.sh.
    - Keep the 10-minute lock action if desired; remove only the suspend action.

3. Make suspend opt-in during long downloads (inhibitor)

- Use a one-off sleep inhibitor while running a download:

```bash
systemd-inhibit --what=sleep --why="Active download" your-download-command-here
```

- Or keep the session awake until you cancel:

```bash
systemd-inhibit --what=sleep --why="Keep awake for downloads" bash -c 'while :; do sleep 300; done'
```

4. Optional hard-disable of system sleep targets

- If you want to completely prevent sleep system-wide:

```bash
sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target
```

- Re-enable later:

```bash
sudo systemctl unmask sleep.target suspend.target hibernate.target hybrid-sleep.target
```

5. Verification checklist

- Logind shows idle ignore:

```bash
loginctl show-logind -p IdleAction -p IdleActionUSec
```

- No active Hypridle suspend action is present; only lock is configured (or Hypridle not running during long downloads).
- Inhibitors show expected blockers during downloads:

```bash
systemd-inhibit --list
```

- The system remains awake past your previous timeout (e.g., run a 20-minute no-op test):

```bash
timeout 1200 bash -c 'date; sleep 1200; date'
```

Note:

- ArchRiot does not restart systemd-logind during install/upgrade. After creating or modifying drop-ins under /etc/systemd/logind.conf.d, reboot to apply cleanly.
- Keeping the lock at 10 minutes while disabling suspend preserves security without interrupting transfers.

## 🛠️ Development Tools

_For contributors and power users who want to build from source_

```bash
# Build ArchRiot from source (in repo directory)
make                                 # Build the Go binary
make test                            # Run test suite
make dev                             # Development build

# Version checking
cat ~/.local/share/archriot/VERSION  # Show installed ArchRiot version
```

## 📂 Repository Information

- **Main Repository:** [https://github.com/CyphrRiot/ArchRiot](https://github.com/CyphrRiot/ArchRiot)
- **Maintenance:** Active, with regular updates and improvements
- **Community:** Open to issues, suggestions, and contributions

## 🌐 Connect with CyphrRiot

<div align="center">

**Stay connected for updates, tips, and Linux excellence**

[![Follow CyphrRiot on X](https://img.shields.io/badge/Follow-@CyphrRiot-1DA1F2?style=for-the-badge&logo=x&logoColor=white)](https://x.com/CyphrRiot)
[![GitHub Profile](https://img.shields.io/badge/GitHub-CyphrRiot-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/CyphrRiot)

_Building the Linux system you actually want to use_

</div>

## 📋 System Requirements

- **Fresh Arch Linux installation** (recommended)
- **Internet connection** for package downloads
- **4GB+ RAM** (8GB+ recommended for development)
- **8GB+ storage** (10GB+ for full development setup)
- **CPU:** Any modern processor (optimized for both Intel/AMD)
- **GPU:** Any modern graphics card (see GPU Support section for full compatibility details)

_Note: ArchRiot began as a unique rice[^1] and evolved from [DHH's Omarchy](https://omarchy.org/) installer, but has since become a completely distinct Linux distribution based on Arch. It features a custom installer, advanced Go-based package management system, and many custom applications and tools focused on privacy, development productivity, and clean aesthetics._

[^1]: In the context of Linux, "rice" is slang for customizing or tweaking a desktop environment or user interface to make it look aesthetically pleasing or highly personalized, often with a focus on minimalism, unique themes, or lightweight setups. It comes from the term "ricer," originally used in car culture to describe heavily modified cars (inspired by "rice burner" for Japanese cars).

## 🖥️ VM & Hardware Notes

**ArchRiot is designed for bare metal hardware.** While it works perfectly in VMs, you're missing the full experience. This system is built to replace whatever disappointing Linux distribution you're currently stuck with. Install it on real hardware where it belongs.

**For VM testing:**

- **VirtualBox/VMware:** Works out of the box with the ArchRiot ISO
- **QEMU/KVM:** Full acceleration support with virtio drivers
- **Recommended VM specs:** 4GB+ RAM, 20GB+ storage, EFI boot enabled

**Hardware compatibility:**

- **Multi-GPU systems:** Automatic detection and driver installation
- **High-DPI displays:** Proper scaling across all applications
- **Gaming hardware:** Full support for NVIDIA/AMD/Intel graphics
- **Modern laptops:** Battery management, backlight control, and power profiles

### Fractional Scaling (Wayland) — Behavior and Tips

- Per‑monitor scaling is supported and configured via Hyprland monitors (editable in `~/.config/hypr/monitors.conf`) or the Control Panel’s Display Settings.
- Wayland‑native apps render crisply at fractional scales (e.g., 0.9, 1.25). XWayland apps may look slightly soft at non‑integer scales.
- Brave/PWAs: ArchRiot’s Brave wrapper sets `--force-device-scale-factor=1` and `--high-dpi-support=1` to prevent content shrink at scales below 1.0. Use per‑site/page zoom in Brave for fine‑tuning at >1.0 scales.
- Mixed‑DPI (multi‑monitor): Prefer per‑monitor Wayland scaling; avoid global Xft.dpi hacks. Use kanshi profiles if you hot‑plug monitors frequently.
- Verification:
    - `hyprctl monitors` shows current scale per display
    - Waybar and GTK apps should scale cleanly as you adjust monitor scale
    - Screenshots reflect scaled geometry as expected
- Known quirks:
    - Some Electron apps under XWayland can appear soft at fractional scales — install Wayland builds when available
    - Legacy toolkits may ignore fractional scaling; consider alternatives or run with specific toolkit flags if needed

## QA Matrix — 3.5 Release Validation

Hardware

- Intel iGPU (modern/Arc), AMD Radeon, NVIDIA proprietary

Displays

- Single monitor
- Multi‑monitor (dock/undock)
- Fractional scaling (0.9, 1.25)

Locale

- en‑US
- One non‑Latin (e.g., ja‑JP or ar)

Validation checklist

- [ ] Waybar single‑instance works; reload is SIGUSR2‑first; ~/.cache/archriot/runtime.log present
- [ ] Portals OK: Kooha screencast prompt works; OBS shows “Screen Capture (PipeWire)”
- [ ] Help binds: SUPER+H opens website; SUPER+SHIFT+H shows local help with current binds
- [ ] Thunar “Open Terminal Here” launches Ghostty
- [ ] Zed focus‑or‑launch (native > Flatpak) via SUPER+Z
- [ ] Updater dialog: non‑blocking; copy includes logs path
- [ ] Brave wrapper: brave://gpu shows acceleration; safe‑mode fixes any crash on workspace/monitor switch
- [ ] Fractional scaling behaves as documented; Brave content not shrunken at <1.0
- [ ] Control Panel window scales (60% x 85%) across monitors/scales
- [ ] Screenshots (optional): capture confirmations for README refresh

**You're done!** If the installer finished, your ArchRiot system is ready to rock. Reboot and enjoy your perfectly configured desktop.

**Something not working?** Re-run the installer - it's designed to fix problems and maintain your system.

🎉 Thank you [Vaxryy](https://x.com/vaxryy) for creating Hyprland—the compositor that doesn't suck.

And, thank you to JaKoolIt for [your amazing scripts](https://github.com/JaKooLit/Arch-Hyprland)!

## 👥 Contributors

- Tarso Galvão (surtarso)
    - GitHub: https://github.com/surtarso
    - Notable contributions: Expanded Waybar workspace styles (ModulesWorkspaces), per-window icon mapping, related Waybar config polish.

## 🧱 Project Architecture

ArchRiot is structured so that the entrypoint (`source/main.go`) only parses CLI flags and delegates to modular packages. This keeps the core small, testable, and safe to evolve.

- Entrypoint (delegation-only)
    - `source/main.go`
        - Parses flags
        - Delegates to package functions
        - Contains no feature logic; only minimal glue

- Core session helpers
    - `source/session/`
        - Waybar lifecycle: `LaunchWaybar`, `ReloadWaybar`, `SweepWaybar`
        - Reload coalescer for Hyprland: `ReloadHyprland`, `ReloadHyprlandImmediate`
        - Power and stability: `PowerMenu`, `StabilizeSession`, `RebootNow`
        - Environment guards and flows: `SuspendGuard`, `MullvadStartup`
        - UX helpers: `PomodoroClick`, `PomodoroDelayToggle`, `WorkspaceClick`, `WelcomeLaunch`, `StartupBackground`

- Theming
    - `source/theming/` + `source/theming/applications/`
        - Theme registry and application (e.g., Waybar, Ghostty, Hyprland, Text Editor)
        - Integrated with the reload coalescer to avoid double reloads

- Display and workspace tools
    - `source/displays/` (e.g., `Autogen` for kanshi)
    - `source/waybartools/` (e.g., `SetupTemperature`)
    - `source/windows/` (e.g., `Switcher`, `FixOffscreen`)

- Upgrade, guardrails, and secure boot
    - `source/upgradeguard/` (e.g., `PreInstall` ABI guard for compositor/Wayland)
    - `source/upgrunner/` (e.g., `Bootstrap` upgrade flow orchestration)
    - `source/secboot/` (e.g., `RunContinuation`, `RestoreHyprlandConfig`)

- CLI utilities and installer helpers
    - `source/cli/` (e.g., `ShowHelp`, `ValidateConfig`)
    - `source/installer/` (e.g., `EnsureYay` bootstrap)
    - `source/executor/`, `source/orchestrator/` (installation orchestration)

Key rules (enforced in code and process)

- `main.go` is an entrypoint only: parse flags → delegate to packages.
- Every CLI flag maps to a dedicated package function (no feature logic in `main.go`).
- One change at a time; always run `make` after every change and stage exact files.
- Remove unreachable or duplicate code immediately to prevent drift.

## ✨ What’s New in v3.9.2

- New first‑class memory clean:
    - `~/.local/share/archriot/install/archriot --memory-clean`
    - Options:
        - `--if-low` (only clean when MemAvailable < 1024 MB by default)
        - `--threshold-mb N` (adjust threshold)
        - `--notify`
    - Elevation is automatic: root → pkexec (GUI) → sudo fallback

- Optional low‑memory stabilize (no compositor reload):
    - `~/.local/share/archriot/install/archriot --stabilize-session --memory-clean --notify`

- Optional keybinding (disabled by default unless you enable it):
    - `bind = $mod CTRL SHIFT, K, exec, $HOME/.local/share/archriot/install/archriot --stabilize-session --memory-clean --notify  # Stabilize (low-memory)`

## ✨ What’s New in v3.6

- Installer hang mitigation: per-command timeouts with output capture to prevent stalls; non‑critical commands continue, critical commands fail fast with clear logs. Extended timeouts for pacman/yay operations.
- Hypridle reliability: autostart now uses an absolute path (/usr/bin/hypridle); ensureIdleLockSetup appends the same to avoid PATH/env timing issues after upgrades.
- GTK Control Panel safety: settings reapply runs only when a GUI session is detected; forces GSK_RENDERER=gl; automatic cairo software-renderer fallback if GL fails to avoid GTK OpenGL context errors.
- Brave multi‑monitor flags: documented optional flags for Wayland multi‑monitor stability/performance under the Brave Wrapper section.
- Documentation polish: fixed TOC/deep links; removed duplication (Memory Tuning); updated version badge and Current Release to v3.6.1.

## ✨ What’s New in v3.5 (Docs & UX)

- 🔑 Keybindings Help (SUPER+SHIFT+H): Local, dynamic help window generated from your Hyprland config; opens as a compact Brave “app” window. GTK fallback available.
- ✉️ Telegram (SUPER+G): Focus-or-launch behavior now mirrors Fuzzel; includes a resilient delayed-retry on cold starts.
- 🛡️ Mullvad during upgrade: Auto-reconnect guard if you start the upgrade connected (no service restarts or NM changes).
- 🧩 Installer guard (ABI mismatch): Warns if Hyprland/wlroots/portal/Wayland upgrades are pending; use `--strict-abi` to block install until the system is upgraded.
- 🔋 Battery “Full”: Shows 100% when the controller reports “fully-charged” for clearer UX.
- 📚 README polish: Emoji “Navigate This Guide,” Quick Install panel, and deep technical content consolidated into Advanced Usage and Troubleshooting.

<a id="license"></a>

## 📄 License

ArchRiot is released under the [MIT License](https://opensource.org/licenses/MIT), enabling community contributions and modifications.

# 🛡️⚔️🪐 Hack the Planet 🪐⚔️🛡️

# ArchRiot 3.4 Release Notes

Focused on a unified default font, UI polish, and Wayland-friendly scaling.

## Highlights

- Paper Mono as the default system monospace font
    - Shipped in config/fonts and installed to ~/.local/share/fonts
    - Applied across GTK3/GTK4, Waybar, Fuzzel, Ghostty, Control Panel, Welcome, and Zed
    - Symbols coverage via ttf-nerd-fonts-symbols retained for icons/glyphs
- Brave webapps and windows scale correctly at fractional monitor scales
    - Added --force-device-scale-factor=1 --high-dpi-support=1 to Brave (Wayland)
    - Prevents content shrinking inside PWAs when monitor scale is < 1.0 (e.g., 0.9)
- Thunar “Open Terminal Here” launches Ghostty directly
    - Reliable even without exo preferred apps configured
- Control Panel wallpaper section no longer grows offscreen
    - Background slider wrapped in a horizontal scroller
    - Current wallpaper label wraps/ellipsizes to preserve layout

### Font credit and usage

Paper Mono by Paper Design is used as the default system monospace font in ArchRiot 3.4.

- Project: https://github.com/paper-design/paper-mono/tree/main
- The font files are included under `config/fonts` and are installed to `~/.local/share/fonts` during setup.
- Please refer to the upstream repository for license and attribution details.

## Upgrade

```
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

# ArchRiot 3.2 Release Notes

Focused on persistence and UX polish for volume controls, keybindings, and SSH keys.

## Highlights

- No duplicate volume notifications
    - Hardware XF86 audio keybindings now invoke Volume.sh with --no-notify to prevent double toasts
    - Waybar visuals remain unchanged (scroll/click still works as before)
- Persistent user keybindings
    - Hyprland uses config/hypr/hyprland.conf directly by default (no separate keybindings.conf).
    - Installer preserves this file (alongside monitors.conf) so your binds survive repairs
- Shell persistence and SSH
    - Advises putting personal aliases in ~/.config/fish/conf.d/local.fish (preserved by installer)
    - Provides ~/.config/fish/conf.d/ssh-agent.fish (not installed by default). To enable: ln -s ~/.local/share/archriot/config/fish/conf.d/ssh-agent.fish ~/.config/fish/conf.d/ssh-agent.fish
- Unified fullscreen examples
    - Documented, optional binds you can enable in keybindings.conf:
        - bind = $mod, F11, fullscreen
        - bind = ALT, RETURN, fullscreen

## Upgrade

```
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

# ArchRiot 3.1 Release Notes

Focused on stability and UX polish related to Waybar modules, audio, and app bindings.

## Highlights

- AMD CPU temperature detection for Waybar custom temp module (k10temp/zenpower), plus Intel and generic fallbacks
- Waybar audio resiliency: volume module degrades gracefully when pamixer/PipeWire aren't ready, preventing crashes
- Hardware volume keys in Hyprland: XF86AudioRaiseVolume, XF86AudioLowerVolume, XF86AudioMute, XF86AudioMicMute wired to Volume.sh
- Keyboard settings launcher now respects $VISUAL/$EDITOR and falls back across common editors
- Bluesky launcher (.desktop) and PWA window rules for a consistent floating panel experience
- Waybar Mullvad module definition present; shows disconnected when Mullvad isn't installed

## Upgrade

```
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

# ArchRiot 3.0 Release Notes

Focused on reliability, portability, and predictable behavior across diverse hardware — while staying true to ArchRiot’s privacy-first, dev-centric philosophy.

## Highlights

- Waybar portability
    - Network module no longer hardcodes interface names; auto-detects reliably.
    - Temperature module no longer hardcodes hwmon paths; auto-detects sensors across hardware.
    - Flicker-free reloads via SIGUSR2 (no kill/restart).

- Screen recording indicator
    - Kooha recording dot shown as the left-most item on the right; click-to-stop recording.
    - Clean, lightweight PipeWire-based detection; zero idle footprint when idle.

- Portals stack aligned with Hyprland
    - Added xdg-desktop-portal, xdg-desktop-portal-hyprland, xdg-desktop-portal-gtk for consistent screencast/screenshot/file chooser behavior.

- Audio reliability
    - rtkit added for PipeWire/WirePlumber realtime scheduling under load.

- Memory tuning (opt-in)
    - When enabled, vm.min_free_kbytes ≈ 1% of RAM (capped at 256MB) and overcommit heuristic (vm.overcommit_memory=0) to avoid exec starvation while maintaining responsiveness.

- GPU auto-detect
    - Installs the correct NVIDIA/AMD/Intel driver and VA stack based on lspci detection.

- Hyprlock polish
    - CPU and memory label refresh reduced to 5s to lower background CPU use while locked.

- Fuzzel stays
    - We do not use Walker; Fuzzel remains the launcher.

## Safer defaults

- TRIM disabled by default
    - No automatic fstrim.timer enablement (especially important for LUKS+btrfs). If you want weekly TRIM, enable `fstrim.timer` yourself and consider adding `discard` to your LUKS mapping (with the usual privacy trade-offs).

- Brave VA-API flags not forced
    - We didn’t add aggressive browser flags by default; if you want hardware decode, verify with `vainfo` and adjust your browser flags locally as needed.

## Deferred for a later update

- Secure Boot guided sbctl flow (end-to-end).
- Optional systemd-oomd integration (opt-in with protected slices).

## Upgrade

Run the one-liner:

```
curl -fsSL https://ArchRiot.org/setup.sh | bash
```

Thanks for the feedback and testing that shaped 3.0. This release focuses on rock-solid defaults and cross-hardware portability without compromising ArchRiot’s principles.

## ✨ What’s New in v3.9.3

- Installer safeguard:
    - monitors.conf backup/edits are now skipped when Hyprland is active (hyprctl -j monitors succeeds).
    - Prevents display resets and Brave crashes during upgrades.
- Process guardrails (docs):
    - plan.md updated with non‑negotiable rules to avoid compositor reloads in installer flows, never add new .sh for features (use Go flags), only edit README at the bottom, and gate systemd unit enable/start on unit existence.
- No user‑visible feature changes; this is a safety and reliability release.

## ✨ What’s New in v3.9.4

- Portal stall recovery:
    - New flag: `~/.local/share/archriot/install/archriot --portals-restart`
    - Safely restarts `xdg-desktop-portal-hyprland` and `xdg-desktop-portal` via systemd --user when available; falls back to best‑effort kill+start.
- Optional keybinding (commented suggestion):
    - `bind = $mod SHIFT, P, exec, $HOME/.local/share/archriot/install/archriot --portals-restart  # Restart portals`
- No compositor reloads; this is a targeted reliability improvement.

## ✨ What’s New in v3.9.5

- Hyprland config safety during upgrades:
    - Diff-aware config copy: skip writes when unchanged (prevents auto-reload).
    - No monitors.conf backup/sed during install (zero display churn).
    - Always enforce displays post-config updates (no Brave gating).
- Edit-time stability:
    - kanshi autostart disabled inside Hyprland to avoid reload-induced output flips while editing.
- No compositor reloads; this is a reliability-focused release.
