# ArchRiot Development Roadmap

## Current Status: v1.9.2

- Installation system reliability: ✅ Complete
- Theme config fixes: ✅ Complete
- Waybar update notifications: ✅ Restored

## IMMEDIATE PRIORITIES

### 1. 🚨 CRITICAL: Eliminate Theme Config Nightmare

**Status**: Urgent - Causing maintenance hell and silent feature losses
**Problem**: Theme configs completely override main config instead of extending it
**Impact**: Every new feature must be manually added to every theme config or it disappears

**Actions**:

#### Phase 1: Analysis Complete ✅

- [x] Map all files being overridden by theme system
- [x] Identify scope of theme config replacements

**Files Currently Being Completely Replaced by Themes:**

- `waybar/config` - Main waybar configuration (CRITICAL: update notifications)
- `fuzzel.ini` - Fuzzel launcher configuration
- `neovim.lua` - Neovim theme plugin
- `btop.theme` - btop system monitor theme
- `mako.ini` - Notification daemon configuration
- `hyprlock.conf` - Lock screen configuration (sourced via main config)

**Theme-Only Files (not in main config):**

- `ghostty.conf` - Terminal configuration
- `hyprland.conf` - Hyprland window manager config
- `backgrounds/` - Theme-specific background images
- `text-editor/*.xml` - Text editor syntax themes

#### Phase 2: File Consolidation

- [ ] Copy CypherRiot configs to replace main configs (waybar, fuzzel, mako, btop)
- [ ] Move CypherRiot backgrounds to main backgrounds directory
- [ ] Move CypherRiot ghostty.conf to main config
- [ ] Move CypherRiot hyprland.conf to main config (if different)
- [ ] Move CypherRiot text-editor themes to main config

#### Phase 3: Theme System Removal

- [ ] Remove tokyo-night theme directory completely
- [ ] Remove cypherriot theme directory
- [ ] Remove theme linking logic from `install/desktop/theming.sh`
- [ ] Remove theme symlink creation (`~/.config/archriot/current/theme`)
- [ ] Update hyprlock.conf to use direct paths instead of theme sourcing

#### Phase 4: Installer Updates

- [ ] Remove `setup_archriot_theme_system()` function
- [ ] Remove `set_default_theme()` function
- [ ] Remove `link_theme_configs()` function
- [ ] Remove `link_waybar_config()` function
- [ ] Simplify theming.sh to only handle cursors, icons, GTK themes
- [ ] Test installation with simplified theme system

#### Phase 5: Documentation & Testing

- [ ] Update README.md to remove theme selection instructions
- [ ] Update INSTALLER.md to remove theme system warnings
- [ ] Test all applications work with consolidated configs
- [ ] Verify waybar update notifications work
- [ ] Test lock screen functionality
- [ ] Test fuzzel launcher functionality

### 2. 🔧 Fix Fuzzel Sudo Integration

**Status**: High Priority - Affects core user experience
**Problem**: Migrate command works in terminal but fails in Fuzzel launcher
**Impact**: Users cannot access migration tool from application launcher

**Actions**:

- [ ] Implement pkexec wrapper for migrate command
- [ ] Create GUI sudo helper for privileged operations
- [ ] Update desktop integration files
- [ ] Test fuzzel integration with sudo commands
- [ ] Verify migrate tool launches from application menu

### 3. 🚨 CRITICAL: Fix Installation Failures ✅ RESOLVED

**Status**: COMPLETED - Critical failures have been fixed
**Problem**: Multiple critical installation failures discovered during testing
**Impact**: Users cannot complete installation due to broken components

**Fixed Issues**:

- ✅ **Migrate tool download**: Fixed "text file busy" error with temp file approach
- ✅ **PipeWire audio system**: Fixed dependency conflicts with proper conflict detection
- ✅ **Local variable errors**: Fixed all `local` declarations outside functions
- ✅ **Audio system**: Complete audio stack now installs properly

### 4. 📊 Implement Modern Progress Bar System

**Status**: FUTURE ENHANCEMENT - Clean installer experience
**Problem**: Current installer output is verbose and overwhelming
**Impact**: Users can't see actual progress through installation noise

**Desired Experience**:

```
┌─────────────────────────────────────────────────────────┐
│ ArchRiot Installation Progress                [73%]     │
├─────────────────────────────────────────────────────────┤
│ Installing: Desktop Environment                         │
│ ████████████████████████████████░░░░░░░░░░░░░░░░░░░░░░  │
└─────────────────────────────────────────────────────────┘
```

**Technical Requirements**:

- **Static progress bar**: Stays in place, updates percentage in real-time
- **Module names**: Show friendly names ("Desktop Environment" not "install/desktop/apps.sh")
- **Background logging**: All verbose output goes to log files
- **Error surfacing**: Only show actual errors on console, not warnings
- **Clean completion**: Beautiful summary at end

**Implementation Plan**:

- [ ] Create progress bar UI system (similar to migrate tool)
- [ ] Map technical modules to user-friendly names
- [ ] Redirect all package manager output to log files
- [ ] Implement error vs warning detection
- [ ] Add real-time percentage calculation based on module completion
- [ ] Test that errors still surface properly while maintaining clean UI
- [ ] Preserve all reliability improvements from v1.9.0+

**Success Criteria**:

- Beautiful, clean console output during installation
- All verbose logs written to files in background
- Real errors immediately visible to user
- Progress bar shows actual completion percentage
- User can see what's happening without noise

### 5. 🚀 System-wide v1.9.2 Deployment

**Status**: Medium Priority - Production rollout
**Problem**: Not all systems have latest theme fixes
**Impact**: Update notifications may be missing on some systems

**Actions**:

- [ ] Deploy v1.9.2 to all production systems
- [ ] Verify update notification functionality
- [ ] Monitor for edge cases or regressions
- [ ] Collect user feedback on theme fixes
- [ ] Document any deployment issues

### 6. ✅ Theme System Verification

**Status**: Medium Priority - Quality assurance
**Problem**: Need to verify theme consolidation doesn't break functionality
**Impact**: Ensure simplified theme system maintains all features

**Actions**:

- [ ] Test waybar theming with unified config
- [ ] Verify lock screen themes work correctly
- [ ] Check desktop theming elements
- [ ] Validate all color schemes and styling
- [ ] Update theme documentation

## NEXT DEVELOPMENT PHASE

### Short-term Goals (Next 2 weeks)

- Complete theme system consolidation
- Fix fuzzel sudo integration
- Restore progress bars
- Deploy v1.9.2 everywhere

### Medium-term Goals (Next month)

- Enhanced verification system
- Performance optimization
- Better error recovery tools
- Improved user onboarding

### Long-term Vision (Next quarter)

- Automated updates system
- Plugin architecture for extensibility
- Cross-distribution support
- Advanced customization options

## SUCCESS METRICS

### Theme System Success

- ✅ Single theme config file
- ✅ No feature loss when adding new capabilities
- ✅ All theming functionality preserved
- ✅ Reduced maintenance overhead

### Installation Experience Success

- ✅ Progress bars visible during installation
- ✅ All errors immediately visible
- ✅ No silent failures
- ✅ User feedback is positive

### Integration Success

- ✅ Migrate tool launches from Fuzzel
- ✅ All sudo operations work in GUI
- ✅ Desktop integration is seamless
- ✅ No privilege escalation failures

## NOTES

**Critical Learning**: Theme overrides are a maintenance nightmare. The discovery that theme configs completely replace main configs rather than extending them explains why features like update notifications disappeared. This must be the #1 priority to prevent future silent failures.

**Deployment Strategy**: v1.9.2 fixes are critical for production systems. Theme notification functionality is essential for user experience.

**Development Approach**: One focused change per iteration, always ask "Continue?" before modifications, never commit without explicit permission.
