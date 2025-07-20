# ArchRiot Critical Issues & Implementation Plan

## 🚨 CRITICAL INSTALLER BUGS (URGENT)

### 1. Hidden Desktop Files Not Installing ❌ BROKEN

**Problem:** Fresh installs still show unwanted applications in Fuzzel menu

- "About Xfce"
- "btop++"
- "GNOME System Monitor" (should be hidden, only custom "System Monitor" should show)
- "mpv Media Player" (should be hidden, only custom "Media Player" should show)

**Root Cause:** `utilities.sh` script runs but hidden files aren't being copied

- ✅ Hidden files exist: `~/.local/share/archriot/applications/hidden/*.desktop`
- ❌ Files not copied to: `~/.local/share/applications/`
- ✅ Script executes: Shows "Running: utilities" and "Successfully completed"
- ❌ No debug output visible (installer swallows stderr/stdout)

**Failed Attempts:**

1. Fixed `script_dir` variable definition - didn't work
2. Added debugging output - not visible in installer
3. Multiple path fixes - still broken
4. Fixed utilities.sh with absolute paths and simplified logic - script runs but files don't appear
5. Added comprehensive debugging output - shows copy commands succeed (result: 0) but files vanish

**Current Status:**

- ✅ Script finds hidden directory correctly
- ✅ Copy commands execute successfully (exit code 0)
- ❌ Files disappear immediately after copy or are copied to wrong location
- ❌ ~/.local/share/applications/ shows no hidden files despite successful copy

**Next Steps:**

- Check if files are being removed immediately after copy
- Verify target directory exists and is writable
- Test manual copy to isolate the issue
- Check for competing processes overwriting files

### 2. Missing Icons on Fresh Installs ❌ BROKEN

**Problem:** Web applications show without icons after fresh install

- Twitter (X)
- Volume Control
- Zed
- Wifi Manager
- Proton Mail
- Google Messages

**Root Cause:** Icon files not being copied during installation

- Icons exist in repo: `applications/icons/*.png`
- Copy logic in utilities.sh may be failing
- Desktop files reference icons that don't exist in system

### 3. Thunar Bookmarks Using Literal $HOME ❌ BROKEN

**Problem:** Thunar bookmarks show "Failed to open file://$HOME/Videos"

- Bookmarks file contains literal `$HOME` instead of expanded paths
- Should be: `file:///home/username/Videos`
- Actually is: `file://$HOME/Videos`

**Failed Fix:** Changed heredoc from `'EOF'` to `EOF` and explicit echo commands
**Status:** Still broken after fresh install

### 4. Backup System Chaos ✅ FIXED

**Problem:** Installer created 100+ scattered backup directories
**Solution:** Consolidated to single `~/.archriot/backups/` that overwrites previous

## 🔧 RECENT FIXES COMPLETED

### Intel Graphics Boot Hang ✅ FIXED

- **Problem:** Dell XPS 13 9350 hanging on "ARCHRIOT" screen
- **Cause:** `LIBGL_ALWAYS_SOFTWARE=1` in both graphics.conf and fish config
- **Solution:** Removed software rendering, uses hardware acceleration
- **Status:** ✅ System boots properly now

### Window Rules for System Monitor ✅ IMPLEMENTED

- Added centering rules for gnome-system-monitor
- Matches Feather wallet behavior when clicked from Waybar

### Application Name Cleanup ✅ IMPLEMENTED

- "GNOME System Monitor" → "System Monitor"
- "mpv Media Player" → "Media Player"
- "Thunar File Manager" → "File Manager"
- "Removable Drives and Media" → "Removable Drives"

### GNOME Secrets Integration ✅ IMPLEMENTED

- Added password manager to utilities installer
- Perfect GTK integration for ArchRiot aesthetic

## 🎯 IMMEDIATE PRIORITIES

1. **Fix hidden desktop files installer** (blocks clean Fuzzel menu)
2. **Fix missing icons** (professional appearance)
3. **Fix Thunar bookmarks** (basic functionality)
4. **Remove debugging output** (clean up after fixes)

## 📋 TESTING CHECKLIST

For next installer fix test on fresh system:

### Hidden Applications Check

```bash
# Should NOT appear in Fuzzel:
- About Xfce
- btop++
- GNOME System Monitor (original)
- mpv Media Player (original)
- Qt5 Settings
- Qt6 Settings
- Thunar Preferences

# Should appear ONCE each:
- System Monitor (custom)
- Media Player (custom)
- File Manager (custom)
```

### Icons Check

```bash
# Should have proper icons:
- Twitter (X)
- Volume Control
- Zed
- Wifi Manager
- Proton Mail
- Google Messages
```

### Thunar Bookmarks Check

```bash
# Should work without errors:
cat ~/.config/gtk-3.0/bookmarks
# Should contain: file:///home/username/Downloads (not file://$HOME/Downloads)
```

## 🔄 VERSION TRACKING

- **Current:** 1.1.55 (with installer fixes)
- **Next:** 1.1.56 (after hidden files mystery is solved)

## 📝 LESSONS LEARNED

1. **Test on fresh systems** - fixes that work locally may not work on clean installs
2. **Installer output is hidden** - debug messages don't show during install
3. **Relative paths are fragile** - use absolute paths in installers
4. **One change at a time** - multiple simultaneous changes make debugging impossible
5. **Ask before implementing** - prevents breaking working systems
6. **Copy success ≠ file persistence** - files can copy successfully but still disappear
7. **Debug everything** - even successful operations may have hidden issues
8. **Track progress properly** - plan.md should be in git for collaboration
