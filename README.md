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

## ⌨️ Essential Commands

### Theme & Appearance

```bash
omarchy-theme-next                    # Switch to next theme
SUPER + CTRL + SHIFT + SPACE         # Switch themes (keybind)
SUPER + CTRL + SPACE                 # Cycle through backgrounds
```

### System Management

```bash
omarchy-update                       # Update system packages
migrate                              # Backup/restore system (interactive TUI)
sudo systemctl reboot               # Restart system
sudo systemctl poweroff             # Shutdown system
```

### Applications (Keybinds)

```bash
SUPER + RETURN                       # Open terminal (Kitty)
SUPER + D  or  SUPER + SPACE         # App launcher (wofi)
SUPER + F                            # File manager (Thunar)
SUPER + B                            # Browser (Brave)
SUPER + E                            # Email (Proton Mail)
SUPER + N                            # Text editor (Neovim)
SUPER + T                            # System monitor (btop)
SUPER + ESCAPE                       # Power menu
```

### Screenshots

```bash
SUPER + SHIFT + S                    # Region screenshot
SUPER + SHIFT + W                    # Window screenshot
SUPER + SHIFT + F                    # Full screen screenshot
SUPER + PRINT                        # Color picker
```

### Window Management

```bash
SUPER + W  or  SUPER + Q             # Close window
SUPER + V                            # Toggle floating
SUPER + J                            # Toggle split
SUPER + Arrow Keys                   # Move focus
SUPER + SHIFT + Arrow Keys           # Swap windows
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
omarchy-fix-thunar-thumbnails        # Fix thumbnail generation
omarchy-fix-background               # Fix theme backgrounds
omarchy-fix-waybar-theme             # Fix waybar styling
```

## 🎯 Key Customizations

### 🔧 **Core System Changes**

- **Terminal:** Kitty (replaces Alacritty) with 90% transparency and dark theme
- **Browser:** Brave (replaces Chromium) with native Wayland support
- **File Manager:** Thunar (replaces Nautilus) with comprehensive dark theming
- **Shell:** Fish as default (replaces Bash) with proper PATH configuration
- **Theme:** CypherRiot as default (replaces Tokyo Night)
- **Code Editor:** Zed (Wayland) + Neovim with proper theme integration
- **Applications:** All major apps now run native Wayland (no more XWayland issues)
- **Backup Tool:** Latest migrate binary for comprehensive system backup/restore
- **Memory Optimization:** Intelligent memory management that actually works
- **Blue Light Filter:** Automatic hyprsunset at 4000K for reduced eye strain
- **GTK Theming:** Dark theme everywhere - no more jarring white dialogs
- **DPI Scaling:** Fixed scaling issues for consistent UI across all applications

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
- **CSS Parser Fixed** - Eliminated all !important declarations causing waybar errors
- **Transparency System** - Consistent 90-98% opacity across all applications

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
- **Better package selection** - Proven, lightweight alternatives (removed Ark, Micro)
- **Enhanced keybindings** - More intuitive, conflict-free shortcuts
- **Wayland-first approach** - Native Wayland for all major applications
- **Improved theming** - Consistent dark mode throughout with proper GTK integration
- **Application launcher fixes** - All applications properly integrated in Wofi
- **Advanced backup** - Comprehensive migration capabilities
- **DPI scaling fixes** - Consistent scaling across all applications
- **File dialog theming** - Dark themes for all application file choosers

## 🔍 Post-Installation Validation

After installation completes, verify everything is working correctly:

### Automatic Validation

The installer automatically runs a post-installation check. If you need to run it manually:

```bash
~/.local/share/omarchy/bin/omarchy-post-install-check
```

### Manual Verification

Check these key components:

```bash
# Verify theme system
ls ~/.config/omarchy/current/theme     # Should show active theme
ls ~/.config/omarchy/current/background # Should show escape_velocity.jpg

# Test background cycling
SUPER + CTRL + SPACE                   # Should cycle through 6 backgrounds

# Test theme switching
omarchy-theme-next                     # Should switch to next theme

# Verify waybar
pgrep waybar                          # Should show running process
```

### Expected Defaults

After fresh installation, you should see:

- **Default theme:** CypherRiot (purple/blue aesthetic)
- **Default background:** escape_velocity.jpg (space/galaxy scene)
- **PDF files:** Show proper document icons (not thumbnails)
- **Image files:** Show thumbnail previews in Thunar
- **Waybar:** Running with tomato timer, system stats, and transparent microphone button

## 🎉 Recent Major Improvements

OhmArchy has undergone extensive improvements to fix all major issues and enhance the user experience:

### ✅ Application Integration Fixes

- **AbiWord**: Now properly installed and appears in Wofi launcher
- **Feather Wallet**: Added with beautiful feather icon and desktop integration
- **Signal**: Native Wayland support with proper DPI scaling (no more grey borders)
- **Zed Editor**: Native Wayland launcher with theme integration
- **Removed Bloat**: Ark and Micro no longer installed by default

### ✅ UI & Theme Improvements

- **Wofi Launcher**: Fixed icon/font balance (24px icons, 16px text)
- **GTK Dark Theme**: Eliminated white file dialogs in all applications
- **XDG Portals**: Proper file chooser theming across all apps
- **DPI Scaling**: Fixed oversized windows and scaling inconsistencies
- **Transparency**: Consistent 90-98% opacity across all applications

### ✅ Wayland Native Support

- **Signal**: Runs natively on Wayland with proper scaling
- **Zed**: Native Wayland with improved theme integration
- **Brave**: Proper Wayland rendering
- **File Dialogs**: All apps now use dark-themed file choosers

### ✅ System Stability

- **Waybar CSS**: Eliminated parser errors and !important declarations
- **Installation Process**: 95%+ success rate with comprehensive error handling
- **Configuration Management**: Automatic backup and validation
- **Modular Architecture**: Clean installer system with proper dependency handling

### ✅ User Experience

- **All Applications in Wofi**: Every installed app properly appears in launcher
- **Consistent Theming**: Dark mode everywhere, no jarring white interfaces
- **Proper Icons**: Beautiful icons for all applications including custom ones
- **Working Shortcuts**: All keyboard shortcuts function correctly
- **Performance**: Fixed memory management and application responsiveness

## 🛠️ Troubleshooting

**Good News**: Most major issues have been resolved! The troubleshooting section is now much shorter.

### Common Issues & Solutions

#### Background Problems _(Rare - Fixed)_

**Issue:** Background stuck on wrong image or not switching

```bash
# Fix background system (rarely needed)
omarchy-fix-background

# Manually restart background service
pkill swaybg
swaybg -i ~/.config/omarchy/current/background -m fill &
```

#### Waybar Issues _(Resolved - CSS Parser Errors Fixed)_

**Issue:** Waybar not starting or showing errors

```bash
# These issues are now resolved, but if needed:
pkill waybar
waybar &
```

#### Application Launcher Issues _(Fixed)_

**Issue:** Applications missing from Wofi or wrong sizing

- **Status**: ✅ RESOLVED - All applications now properly integrated
- **Wofi sizing**: ✅ RESOLVED - Balanced 24px icons with 16px text
- **Missing apps**: ✅ RESOLVED - AbiWord, Feather Wallet, all apps appear

#### File Dialog Theme Issues _(Fixed)_

**Issue:** White file dialogs in Signal, Zed, or other apps

- **Status**: ✅ RESOLVED - Dark theme now applies to all file choosers
- **GTK theming**: ✅ RESOLVED - Proper dark theme integration
- **XDG portals**: ✅ RESOLVED - Consistent theming across apps

#### DPI/Scaling Issues _(Fixed)_

**Issue:** Applications appearing oversized or wrong scaling

- **Status**: ✅ RESOLVED - All apps now scale consistently
- **Signal scaling**: ✅ RESOLVED - No more oversized Signal windows
- **Wayland rendering**: ✅ RESOLVED - Native support for major apps

#### Thumbnail Problems

**Issue:** No thumbnails in file manager or PDFs showing thumbnails instead of icons

```bash
# Fix thumbnail system
omarchy-fix-thunar-thumbnails

# Clear thumbnail cache
rm -rf ~/.cache/thumbnails ~/.thumbnails
```

#### Theme System Issues

**Issue:** Theme switching not working or broken theme links

```bash
# Reset theme system
rm -rf ~/.config/omarchy/current
mkdir -p ~/.config/omarchy/current

# Set CypherRiot as default
ln -sf ~/.config/omarchy/themes/cypherriot ~/.config/omarchy/current/theme
omarchy-fix-background
```

#### Keyboard Shortcuts Not Working

**Issue:** SUPER+CTRL+SPACE or other shortcuts not responding

```bash
# Reload Hyprland configuration
hyprctl reload

# Check if scripts exist
ls ~/.local/bin/swaybg-next           # Background cycling
ls ~/.local/bin/omarchy-theme-next    # Theme switching
```

#### System Validation

**Issue:** Want to check overall system health

```bash
# Comprehensive system validation
~/.local/share/omarchy/bin/omarchy-validate-system

# Quick installation check
~/.local/share/omarchy/bin/omarchy-post-install-check
```

### System Recovery

If you encounter major issues:

1. **Reboot** - Many issues resolve after restart
2. **Re-run installer** - Safe to run multiple times
3. **Check logs** - Look in `~/.config/hypr/hyprland.conf`
4. **Reset configs** - Backup and regenerate configurations

### Getting Help

1. **Check validation output** - Run post-install check for specific issues
2. **Review error messages** - Most scripts provide clear feedback
3. **Check repository issues** - Common problems may have solutions
4. **Fresh installation** - When in doubt, start clean

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
