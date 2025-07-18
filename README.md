# :: 𝔸𝕣𝕔𝕙ℝ𝕚𝕠𝕥 ::

![Version](https://img.shields.io/badge/version-1.1.37-4c1d95)
![License](https://img.shields.io/github/license/CyphrRiot/ArchRiot?color=1e293b)
![Arch Linux](https://img.shields.io/badge/Arch_Linux-0f172a?logo=arch-linux&logoColor=4c1d95)
![Hyprland](https://img.shields.io/badge/Hyprland-1e1e2e?logoColor=3730a3)
![Wayland](https://img.shields.io/badge/Wayland-313244?logo=wayland&logoColor=1e40af)

![Language](https://img.shields.io/badge/language-Shell-1e1e2e)
![Language](https://img.shields.io/badge/language-Python-313244)
![Language](https://img.shields.io/badge/language-CSS-45475a)
![Language](https://img.shields.io/badge/language-Lua-585b70)
![Language](https://img.shields.io/badge/language-HTML-6c7086)

![Maintained](https://img.shields.io/maintenance/yes/2025?color=4c1d95)
![Last Commit](https://img.shields.io/github/last-commit/CyphrRiot/ArchRiot?color=3730a3)
![Code Size](https://img.shields.io/github/languages/code-size/CyphrRiot/ArchRiot?color=1e40af)
![Issues](https://img.shields.io/github/issues/CyphrRiot/ArchRiot?color=64748b)

[![CyphrRiot on X](https://img.shields.io/badge/Follow-@CyphrRiot-1DA1F2?style=for-the-badge&logo=x&logoColor=white)](https://x.com/CyphrRiot)
[![GitHub Profile](https://img.shields.io/badge/GitHub-CyphrRiot-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/CyphrRiot)
[![Migrate Tool](https://img.shields.io/badge/Backup_Tool-Migrate-6B46C1?style=for-the-badge&logo=github&logoColor=white)](https://github.com/CyphrRiot/Migrate)

## **A customized Arch Linux Cypherpunk & Hacker setup**

Turn a fresh Arch installation into a fully-configured, beautiful, and modern development system based on Hyprland by running a single command.

ArchRiot is an even more opiniated setup and was originally a unique rice¹ and then a fork of [DHH's Omarchy](https://github.com/basecamp/archriot) installer with extensive customizations focused on privacy, development productivity, and clean aesthetics.

¹ _In the context of Linux, "rice" is slang for customizing or tweaking a desktop environment or user interface to make it look aesthetically pleasing or highly personalized, often with a focus on minimalism, unique themes, or lightweight setups. It comes from the term "ricer," originally used in car culture to describe heavily modified cars (inspired by "rice burner" for Japanese cars)._

_Created by a hacker, cypherpunk, and blockchain developer with decades of experience running Linux and Unix - this is the system I've always wanted._

![ArchRiot Screenshot](images/screenshot.png)

## 🆕 **What's New in v1.1.36**

**Comprehensive Branding & Optimization Update**

- ✅ **Complete Branding Consistency** - Eliminated all remaining OhmArchy references throughout the codebase
- ✅ **Fixed Critical Waybar Issues** - Corrected power-menu paths and module configurations
- ✅ **Hyprlock Background Fix** - Resolved missing background images on lock screen
- ✅ **Code Quality Improvements** - Enhanced documentation, error handling, and system validation
- ✅ **Performance Monitoring** - Added installation timing and optional debug logging
- ✅ **Clean Codebase** - Removed unused files and fixed filename typos
- ✅ **Enhanced User Experience** - Added customization hints and better completion feedback

This release ensures ArchRiot is fully consistent, optimized, and production-ready with complete ArchRiot branding throughout all 287+ project files.

## 🚀 Installation

### Prerequisites: Fresh Arch Linux Setup

Download the Arch Linux ISO, put it on a USB stick (use balenaEtcher on Mac/Windows), and boot from the stick.

**WiFi Setup** (skip if using ethernet):

1. Run `iwctl`
2. Type `station wlan0 scan`
3. Type `station wlan0 connect <tab>`
4. Pick your network from the list
5. Enter your WiFi password

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

**⚠️ Important:** You must setup disk encryption to use ArchRiot as designed! The setup relies exclusively on disk encryption to secure your device, as it'll auto-login the user after the disk has been decrypted at boot.

Once Arch has been installed, pick reboot, login with the user you just setup, and now you're ready to install ArchRiot.

### Method 1: One-Line Install or Upgrade (Recommended)

```bash
curl -fsSL https://archriot.org/setup.sh | bash
```

**Note: Upgrading is exactly the same command! Simple!**

### Method 2: Manual Clone (For Customization)

```bash
git clone https://github.com/CyphrRiot/ArchRiot.git ~/.local/share/archriot
~/.local/share/archriot/install.sh
```

### Optional: Pre-Installation Validation

Want confidence before installing? Run the validation script to test compatibility:

```bash
curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/validate-installation.sh | bash
```

This comprehensive test validates:

- System requirements (Arch Linux, internet, disk space)
- Package availability (Hyprland, Waybar, etc.)
- GPU compatibility and Wayland support
- Repository accessibility and theme integrity
- User environment and permissions

**Results**: Pass/fail report with specific guidance before you commit to installation.

### Installation Features

- **Automatic backup** - Creates timestamped backup of existing configs
- **Dependency verification** - Ensures Python3 and required packages
- **Script validation** - Verifies all waybar modules are functional
- **Error handling** - Clear feedback and rollback capability
- **100% confidence** - Comprehensive testing and validation

## ⌨️ Essential Commands

### Theme & Appearance

```bash
theme-next                            # Switch to next theme
SUPER + CTRL + SHIFT + SPACE         # Switch themes (keybind)
SUPER + CTRL + SPACE                 # Cycle through backgrounds
```

### System Management

```bash
update                               # Update system packages
migrate                              # Backup/restore system (interactive TUI)
sudo systemctl reboot               # Restart system
sudo systemctl poweroff             # Shutdown system
```

### Applications (Keybinds)

```bash
SUPER + RETURN                       # Open terminal (Ghostty)
SUPER + SHIFT + RETURN               # Centered floating terminal
SUPER + D  or  SUPER + SPACE         # App launcher (fuzzel)
SUPER + F                            # File manager (Thunar)
SUPER + B                            # Browser (Brave)
SUPER + E                            # Email (Proton Mail)
SUPER + G                            # Signal messenger
SUPER + M                            # Google Messages
SUPER + O                            # Text editor (Gnome Text Editor)
SUPER + N                            # Text editor (Neovim)
SUPER + T                            # System monitor (btop)
SUPER + Z                            # Code editor (Zed)
SUPER + X                            # X/Twitter
SUPER + L                            # Lock screen (CypherRiot theme)
SUPER + ESCAPE                       # Power menu
```

### Screenshots

```bash
SUPER + SHIFT + S                    # Region screenshot
SUPER + SHIFT + W                    # Window screenshot
SUPER + SHIFT + F                    # Full screen screenshot
SUPER + PRINT                        # Color picker
```

### Screen Recording

```bash
Kooha                                # GUI screen recorder (launch from SUPER+D)
```

### Window Management

```bash
SUPER + W  or  SUPER + Q             # Close window
SUPER + V                            # Toggle floating
SUPER + J                            # Toggle split
SUPER + Arrow Keys                   # Move focus
SUPER + SHIFT + Arrow Keys           # Swap windows
SUPER + CTRL + Arrow Keys            # Smart window movement
SUPER + 1-9                          # Switch workspace
SUPER + SHIFT + 1-9                  # Move window to workspace
```

### Audio & Media

```bash
XF86AudioRaiseVolume                 # Volume up
XF86AudioLowerVolume                 # Volume down
XF86AudioMute                        # Toggle mute
XF86AudioMicMute                     # Toggle microphone
XF86AudioPlay/Pause                  # Media play/pause

# Beautiful volume overlay appears for 1 second with progress bar
```

### Waybar Controls

```bash
Click tomato timer                   # Start/pause Pomodoro timer
Double-click tomato timer            # Reset timer to 25:00
Click microphone icon               # Toggle microphone mute
Click network icon                  # Open network manager
```

### Fix Scripts (If Needed)

```bash
fix-thunar-thumbnails                # Fix thumbnail generation
fix-background                       # Fix theme backgrounds
```

### 🔧 Troubleshooting

Quick fixes for common issues:

- **Initramfs errors**: Fixed in v1.0.1 - update if experiencing build issues
- **Reliable installer**: Rock-solid installation process with comprehensive error handling
- **Background issues**: Run `fix-background`
- **Thumbnail problems**: Run `fix-thunar-thumbnails`

## 🎯 Key Customizations

### 🔧 **Core System Changes**

- **Terminal:** Ghostty (replaces Kitty) with 90% transparency and dark theme
- **Browser:** Brave (replaces Chromium) with native Wayland support
- **File Manager:** Thunar (replaces Nautilus) with comprehensive dark theming
- **Shell:** Fish as default (replaces Bash) with proper PATH configuration
- **Theme:** CypherRiot as default (replaces Tokyo Night)
- **Code Editor:** Zed (Wayland) + Neovim with proper theme integration
- **Applications:** All major apps now run native Wayland (no more XWayland issues)
- **Backup Tool:** Latest migrate binary for comprehensive system backup/restore
- **Memory Optimization:** Intelligent memory management that actually works
- **Blue Light Filter:** Optional hyprsunset at 3500K for reduced eye strain (configurable during install)
- **GTK Theming:** Dark theme everywhere - no more jarring white dialogs
- **DPI Scaling:** Fixed scaling issues for consistent UI across all applications

#### 🧠 **Memory Management Fix**

Linux's default memory management is **aggressively stupid** about caching. The kernel will happily consume 90%+ of your RAM for file caches, then struggle to free it when applications actually need memory.

ArchRiot fixes this with intelligent memory management tuning:

- **No more lag spikes** when opening applications
- **Better responsiveness** under memory pressure
- **Reduced swap usage** with proper RAM utilization
- **Optimized caching** that doesn't hog system resources

**Result:** Your system stays fast and responsive even when running multiple applications.

### 📱 **Advanced Waybar Integration**

- **Tomato Timer** - Built-in Pomodoro timer with visual states (idle/running/break/finished)
- **Mullvad VPN Status** - Real-time VPN connection status with location display
- **System Monitoring** - CPU aggregate usage, accurate memory monitoring
- **Microphone Control** - Visual mic status with one-click toggle
- **Custom Separators** - Clean, organized module layout
- **CSS Parser Fixed** - Eliminated all !important declarations causing waybar errors
- **Transparency System** - Consistent 90-98% opacity across all applications

### 📱 **Clean Web Applications**

- **Proton Mail** (SUPER+E / XF86Mail) - Privacy-focused email in floating window
- **Signal** (SUPER+G) - Private messaging app
- **Google Messages** (SUPER+M) - Web-based messaging in floating window
- **X/Twitter** (SUPER+X) - Social platform in floating window
- **GitHub** - Development platform with proper icons from homarr-labs

### 🪟 **Responsive Window Management**

- **Percentage-based sizing** - Windows scale properly across different screen resolutions (1080p, 1440p, 4K, ultrawide)
- **Smart centering** - All floating windows automatically center regardless of monitor size
- **Cross-resolution compatibility** - No hardcoded pixel positions, works on any display setup
- **Optimized app windows**:
    - **X/Twitter**: `850x90%` (mobile-style layout with responsive height)
    - **Proton Mail**: `1000x75%` (perfect email reading dimensions)
    - **Google Messages**: `1000x75%` (comfortable messaging interface)
    - **Signal**: `1000x1080` (maintains native desktop experience)
- **Future-proof design** - Window rules adapt automatically to new monitor configurations

### 🎮 **GPU Support**

ArchRiot automatically detects and installs optimal drivers for all major GPUs:

- **NVIDIA**: Proprietary drivers with Wayland and hardware acceleration
- **AMD/Radeon**: Open-source Mesa drivers with Vulkan support
- **Intel**: Mesa drivers including Intel Arc support

All GPUs get proper Wayland integration and hardware video acceleration for optimal performance.

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

### 🎨 **Document & Media Handling**

- **Gnome Text Editor** - Default for text/markdown files (clean, modern text editing with Tokyo Night theme)
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

## 🔄 System Management

### Updates

```bash
update
```

Updates ArchRiot by pulling latest changes and re-running the installer. Simple, safe, and reliable - no dangerous migrations.

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

ArchRiot includes two carefully curated themes with CypherRiot as the default:

- **cypherriot** (default) - Custom purple/blue aesthetic with full waybar integration
- **tokyo-night** - Vibrant city lights with modern cyberpunk vibes

### Theme Management

- **Switch themes:** `theme-next` or manually symlink
- **Theme location:** `~/.config/archriot/current/theme`
- **Backgrounds:** Automatically matched to theme with time-based variants

## ⚡ Key Features & Performance

### Window Management

- **Hyprland compositor** - Smooth animations, efficient memory usage
- **GPU acceleration** - Comprehensive automatic GPU driver setup (see GPU Support section)
- **Tiling & floating** - Flexible window arrangements
- **Multi-workspace** - Organized workflow separation
- **Auto-login** - Direct to tty1 with Hyprland autostart
- **Blue light filter** - Optional hyprsunset reduces eye strain during evening use (installer asks user preference)

### Development Ready

- **Fish shell** - Modern, user-friendly command line with autocompletion
- **Modern CLI tools** - eza, bat, ripgrep, fzf, zoxide for enhanced productivity
- **Git integration** - GitHub CLI, lazygit, proper aliases
- **Code editors** - Zed (Wayland native), Neovim (power user), AbiWord (document editing)
- **Container support** - Docker, development environments
- **Wayland Native** - All development tools run with native Wayland for better performance
- **Theme Integration** - Consistent dark theme across all editors and development tools

### Privacy & Security Focus

- **Brave browser** - Ad blocking, privacy protection by default with native Wayland
- **Proton Mail** - End-to-end encrypted email with XF86Mail key support and proper icon
- **Mullvad VPN** - Anonymous browsing with live waybar status indicator
- **Feather Wallet** - Privacy-focused Monero wallet with beautiful feather icon
- **Signal** - Secure messaging with native Wayland support (no more scaling issues)
- **Local tools** - Reduced dependency on cloud services
- **Clean telemetry** - Minimal data collection
- **Media Downloads** - yt-dlp and spotdl for offline media privacy

### Health & Comfort Features

- **Optional blue light filtering** - installer asks if you want hyprsunset configured with `exec-once = hyprsunset -t 3500` for immediate warm temperature on startup
- **3500K color temperature** - Scientifically optimal warm setting reduces blue light exposure without color distortion
- **Simple management** - Enable/disable by editing `~/.config/hypr/hyprland.conf` (add/remove `exec-once = hyprsunset -t 3500`)
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
- **Better package selection** - Proven, lightweight alternatives (removed Ark, Micro)
- **Enhanced keybindings** - More intuitive, conflict-free shortcuts
- **Responsive window management** - Percentage-based sizing that works across all screen resolutions
- **Wayland-first approach** - Native Wayland for all major applications
- **Improved theming** - Consistent dark mode throughout with proper GTK integration
- **Fixed GTK selection highlighting** - No more jarring bright white highlights in file managers
- **Waybar improvements** - Better font sizing (100%) and improved date format (Sunday • July 13 • 01:49 PM)
- **Zed editor integration** - Native Wayland support with SUPER+Z keybinding
- **Application launcher fixes** - All applications properly integrated in Fuzzel
- **Advanced backup** - Comprehensive migration capabilities
- **DPI scaling fixes** - Consistent scaling across all applications
- **File dialog theming** - Dark themes for all application file choosers

## 🔍 Post-Installation Validation

After installation completes, verify everything is working correctly:

### Automatic Validation

The installer automatically runs a post-installation check. If you need to run it manually:

```bash
~/.local/share/archriot/bin/post-install-check
```

### Manual Verification

Check these key components:

```bash
# Verify theme system
ls ~/.config/archriot/current/theme     # Should show active theme
ls ~/.config/archriot/current/background # Should show escape_velocity.jpg

# Test background cycling
SUPER + CTRL + SPACE                   # Should cycle through 6 backgrounds

# Test theme switching
theme-next                             # Should switch to next theme

# Verify waybar
pgrep waybar                          # Should show running process
```

### Expected Defaults

After fresh installation, you should see:

- **Default theme:** CypherRiot (purple/blue aesthetic)
- **Default background:** riot_zero.png (riot-themed wallpaper)
- **PDF files:** Show proper document icons (not thumbnails)
- **Image files:** Show thumbnail previews in Thunar
- **Waybar:** Running with tomato timer, system stats, and transparent microphone button

## 🧪 Validation & Testing

### Pre-Installation Validation

Before installing, you can run a comprehensive validation script:

```bash
curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/validate-installation.sh | bash
```

**What it tests:**

- System compatibility (Arch Linux, hardware, drivers)
- Internet connectivity and package availability
- Repository accessibility and file integrity
- Wayland/Hyprland compatibility
- Theme files and expected installation outcome

**Results**: Detailed pass/fail report with 30+ validation checks

## 🛠️ Management Tools

```bash
update                               # Update system packages
theme-next                           # Switch to next theme
fix-background                       # Reset background system
validate-system                      # Check system health
```

## 📂 Repository Information

- **Main Repository:** https://github.com/CyphrRiot/ArchRiot
- **Original Upstream:** https://github.com/basecamp/archriot
- **Maintenance:** Active, with regular updates and improvements
- **Community:** Open to issues, suggestions, and contributions

## 📋 System Requirements

- **Fresh Arch Linux installation** (recommended)
- **Internet connection** for package downloads
- **4GB+ RAM** (8GB+ recommended for development)
- **10GB+ storage** (20GB+ for full development setup)
- **CPU:** Any modern processor (optimized for both Intel/AMD)
- **GPU:** Automatic detection and setup for NVIDIA, AMD/Radeon, and Intel graphics cards

## 📄 License

ArchRiot is released under the [MIT License](https://opensource.org/licenses/MIT), maintaining compatibility with the original Omarchy project while enabling community contributions and modifications.
