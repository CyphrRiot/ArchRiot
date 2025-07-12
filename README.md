# OhmArchy

**A customized Arch Linux setup based on Omarchy, optimized for the CypherRiot workflow.**

Turn a fresh Arch installation into a fully-configured, beautiful, and modern development system based on Hyprland by running a single command.

OhmArchy is an even more opiniated personal fork of [Basecamp's Omarchy](https://github.com/basecamp/omarchy) with extensive customizations focused on privacy, development productivity, and clean aesthetics.

![OhmArchy Screenshot](images/screenshot.png)

## 🚀 Installation

### Prerequisites: Fresh Arch Linux Setup

Download the Arch Linux ISO, put it on a USB stick (use balenaEtcher on Mac/Windows), and boot from the stick.

If you're on wifi, start by running `iwctl`, then type `station wlan0 scan`, then `station wlan0 connect <tab>`, pick your network, and enter the password. If you're on ethernet, you don't need this.

Run `archinstall` and pick these options (and leave anything not mentioned as-is):

| Section                  | Option                                                                         |
| ------------------------ | ------------------------------------------------------------------------------ |
| Mirrors and repositories | Select regions > Your country                                                  |
| Disk configuration       | Partitioning > Default partitioning layout > Select disk (with space + return) |
| Disk > File system       | btrfs (default structure: yes + use compression)                               |
| Disk > Disk encryption   | Encryption type: LUKS + Encryption password + Partitions (select the one)      |
| Hostname                 | Give your computer a name                                                      |
| Root password            | Set yours                                                                      |
| User account             | Add a user > Superuser: Yes > Confirm and exit                                 |
| Audio                    | pipewire                                                                       |
| Network configuration    | Copy ISO network config                                                        |
| Additional packages      | Add wget (type "/wget" to filter list)                                         |
| Timezone                 | Set yours                                                                      |

**⚠️ Important:** You must setup disk encryption to use OhmArchy as designed! The setup relies exclusively on disk encryption to secure your device, as it'll auto-login the user after the disk has been decrypted at boot.

Once Arch has been installed, pick reboot, login with the user you just setup, and now you're ready to install OhmArchy.

### Method 1: One-Line Install (Recommended)

```bash
curl -fsSL https://cyphrriot.github.io/OhmArchy/setup.sh | bash
```

### Method 2: Manual Clone (For Customization)

```bash
git clone https://github.com/CyphrRiot/OhmArchy.git ~/.local/share/omarchy
~/.local/share/omarchy/install.sh
```

### Installation Features

- **Automatic backup** - Creates timestamped backup of existing configs
- **Dependency verification** - Ensures Python3 and required packages
- **Script validation** - Verifies all waybar modules are functional
- **Error handling** - Clear feedback and rollback capability
- **100% confidence** - Comprehensive testing and validation

## 🎯 Key Customizations

### 🔧 **Core System Changes**

- **Terminal:** Kitty (replaces Alacritty)
- **Browser:** Brave (replaces Chromium)
- **File Manager:** Thunar (replaces Nautilus)
- **Shell:** Fish as default (replaces Bash) with proper PATH configuration
- **Theme:** CypherRiot as default (replaces Tokyo Night)
- **Code Editor:** Added Zed alongside Neovim
- **Backup Tool:** Latest migrate binary for comprehensive system backup/restore
- **Memory Optimization:** Intelligent memory management that actually works
- **Blue Light Filter:** Automatic hyprsunset at 4000K for reduced eye strain

#### 🧠 **Memory Management Fix**

Linux's default memory management is **aggressively stupid** about caching. The kernel will happily consume 90%+ of your RAM for file caches, then struggle to free it when applications actually need memory. This causes:

- **Lag spikes** when opening new applications
- **Swap thrashing** even with plenty of "free" RAM
- **Poor responsiveness** during memory pressure
- **Misleading memory reports** (cached ≠ available)

**OhmArchy fixes this with intelligent sysctl tuning:**

```bash
vm.min_free_kbytes=1048576    # Always keep 1GB truly free
vm.vfs_cache_pressure=50      # Be less aggressive about caching
vm.swappiness=10              # Prefer RAM over swap usage
vm.dirty_ratio=5              # Limit dirty page cache buildup
```

**Result:** Your system maintains responsive performance with proper memory pressure handling, ensuring applications get the RAM they need without the kernel being stubborn about giving up its precious caches.

### 📱 **Advanced Waybar Integration**

- **Tomato Timer** - Built-in Pomodoro timer with visual states (idle/running/break/finished)
- **Mullvad VPN Status** - Real-time VPN connection status with location display
- **System Monitoring** - CPU aggregate usage, accurate memory monitoring
- **Microphone Control** - Visual mic status with one-click toggle
- **Custom Separators** - Clean, organized module layout

### 📱 **Clean Web Applications**

- **Proton Mail** (SUPER+E / XF86Mail) - Privacy-focused email in floating window
- **Google Messages** (SUPER+ALT+G) - Communication
- **X/Twitter** (SUPER+X) - Social platform
- **GitHub** - Development platform with proper icons from homarr-labs

### ⌨️ **Enhanced Keybindings & Productivity**

- **SUPER+D** = **SUPER+SPACE** (Unified app launcher)
- **Left-click Arch icon** - nwg-drawer app grid
- **Right-click Arch icon** - wofi app launcher
- **XF86Mail** - Floating Proton Mail window
- **SUPER+SHIFT+S** - Region screenshot (primary)
- **SUPER+SHIFT+W** - Window screenshot
- **SUPER+SHIFT+F** - Full screen screenshot
- **Key repeat enabled** (40 rate, 600 delay for responsive typing)
- **All media keys** - Volume, brightness, playback controls

### 🎨 **Document & Media Handling**

- **Apostrophe** - Default for text/markdown files (clean, distraction-free writing)
- **Papers** - Default PDF viewer (GNOME's modern document viewer)
- **MPV** - Video playback with optimal performance
- **Better waybar network** - nmtui instead of impala for reliable WiFi management
- **Screenshot tools** - grim/slurp/hyprshot integration for all capture needs

### 🚫 **Removed Bloat & Corporate Apps**

- **Removed 37signals/Basecamp tools** - Hey, Basecamp web apps
- **Removed corporate social** - Discord, proprietary messaging
- **Removed heavy productivity** - Obsidian, LibreOffice, OBS Studio, KDEnlive, Pinta
- **Removed proprietary services** - 1Password, Typora, Dropbox, Spotify, Zoom
- **Removed entertainment** - YouTube webapp, WhatsApp webapp

## 🔄 System Management

### Updates

```bash
omarchy-update
```

Updates system packages and applies any available OhmArchy migrations.

### Backup & Restore

```bash
migrate
```

**Note:** `migrate` is a TUI (Text User Interface) with **no command-line options**. Simply run the command and use the interactive menu to:

- Create comprehensive system backups
- Restore from previous backups
- Migrate configurations between installations
- Preserve all your customizations

The migrate tool automatically downloads the latest version during installation from [CypherRiot/Migrate](https://github.com/CyphrRiot/Migrate).

## 🎨 Themes & Customization

### Available Themes

OhmArchy includes multiple themes with CypherRiot as the default:

- **cypherriot** (default) - Custom purple/blue aesthetic with full waybar integration
- **catppuccin** - Pastel perfection
- **everforest** - Green nature vibes
- **gruvbox** - Retro warm colors
- **kanagawa** - Japanese ink painting
- **nord** - Arctic cool tones
- **tokyo-night** - Vibrant city lights

### Theme Management

- **Switch themes:** `omarchy-theme-next` or manually symlink
- **Theme location:** `~/.config/omarchy/current/theme`
- **Backgrounds:** Automatically matched to theme with time-based variants

## ⚡ Key Features & Performance

### Window Management

- **Hyprland compositor** - Smooth animations, efficient memory usage
- **GPU acceleration** - Automatic NVIDIA, AMD/Radeon, and Intel driver setup with Vulkan support
- **Tiling & floating** - Flexible window arrangements
- **Multi-workspace** - Organized workflow separation
- **Auto-login** - Direct to tty1 with Hyprland autostart
- **Blue light filter** - Automatic hyprsunset reduces eye strain during evening use

### Development Ready

- **Fish shell** - Modern, user-friendly command line with autocompletion
- **Modern CLI tools** - eza, bat, ripgrep, fzf, zoxide for enhanced productivity
- **Git integration** - GitHub CLI, lazygit, proper aliases
- **Code editors** - Zed (modern), Neovim (power user)
- **Container support** - Docker, development environments

### Privacy & Security Focus

- **Brave browser** - Ad blocking, privacy protection by default
- **Proton Mail** - End-to-end encrypted email with XF86Mail key support
- **Mullvad VPN** - Anonymous browsing with live waybar status indicator
- **Feather Wallet** - Privacy-focused Monero wallet for secure cryptocurrency transactions
- **Local tools** - Reduced dependency on cloud services
- **Clean telemetry** - Minimal data collection

### Health & Comfort Features

- **Automatic blue light filtering** - hyprsunset configured with `exec-once = hyprsunset -t 4000` for immediate warm temperature on startup
- **4000K color temperature** - Scientifically optimal warm setting reduces blue light exposure without color distortion
- **No manual switching needed** - Runs continuously from login, unlike redshift/f.lux time-based switching
- **GPU accelerated filtering** - Native Wayland compositor integration for smooth, lag-free color adjustment
- **Memory pressure relief** - Intelligent VM tuning prevents system lag and swap thrashing
- **Responsive performance** - Conservative memory management keeps applications snappy
- **Clean, minimal UI** - Reduced visual clutter and distractions for focused work

### Audio & Media

- **PipeWire/WirePlumber** - Modern audio stack with low latency
- **MPV** - Lightweight, powerful video player with hardware acceleration
- **Screenshot integration** - Multiple capture methods with clipboard support
- **Media key support** - Volume, brightness, and playback controls work out of the box

## 🔀 Differences from Original Omarchy

This is a **heavily customized fork** optimized for:

### Philosophy Changes

- **Privacy over convenience** - Proton Mail vs. corporate email
- **Performance over features** - Lightweight apps vs. feature-heavy alternatives
- **Development focus** - Tools for coding vs. general productivity
- **Clean aesthetics** - Minimal, distraction-free environment

### Technical Changes

- **Modern shell** - Fish with intelligent defaults
- **Better package selection** - Proven, lightweight alternatives
- **Enhanced keybindings** - More intuitive, conflict-free shortcuts
- **Improved theming** - Consistent dark mode throughout
- **Advanced backup** - Comprehensive migration capabilities

## 📂 Repository Information

- **Main Repository:** https://github.com/CyphrRiot/OhmArchy
- **Original Upstream:** https://github.com/basecamp/omarchy
- **Maintenance:** Active, with regular updates and improvements
- **Community:** Open to issues, suggestions, and contributions

## 📋 System Requirements

- **Fresh Arch Linux installation** (recommended)
- **Internet connection** for package downloads
- **4GB+ RAM** (8GB+ recommended for development)
- **20GB+ storage** (50GB+ for full development setup)
- **CPU:** Any modern processor (optimized for both Intel/AMD)
- **GPU:** Automatic detection and setup for NVIDIA, AMD/Radeon, and Intel graphics cards

## 📄 License

OhmArchy is released under the [MIT License](https://opensource.org/licenses/MIT), maintaining compatibility with the original Omarchy project while enabling community contributions and modifications.
