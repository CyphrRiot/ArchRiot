# ArchRiot Development Roadmap

## Current Status: v2.0.1

- Installation system reliability: ✅ Complete
- Theme config nightmare: ✅ ELIMINATED
- Critical installation failures: ✅ FIXED
- Waybar update notifications: ✅ Simplified and working
- Upgrade path compatibility: ✅ Fixed for v2.0.0+ transition

## IMMEDIATE PRIORITIES

### 1. ✅ COMPLETED: Theme Config Nightmare ELIMINATED

**Status**: COMPLETED in v2.0.0-2.0.1 - Theme system completely redesigned
**Problem**: Theme configs completely override main config instead of extending it
**Impact**: Every new feature must be manually added to every theme config or it disappears

**COMPLETED WORK**:

#### Phase 1: Analysis ✅ COMPLETE

- ✅ Mapped all files being overridden by theme system
- ✅ Identified scope of theme config replacements

#### Phase 2: File Consolidation ✅ COMPLETE

- ✅ Copied CypherRiot configs to replace main configs (waybar, fuzzel, mako, btop)
- ✅ Moved CypherRiot backgrounds to main backgrounds directory
- ✅ Moved CypherRiot ghostty.conf to main config
- ✅ Merged CypherRiot hyprland.conf styling into main config
- ✅ Moved CypherRiot text-editor themes to main config

#### Phase 3: Theme System Removal ✅ COMPLETE

- ✅ Removed tokyo-night theme directory completely
- ✅ Removed cypherriot theme directory
- ✅ Removed theme linking logic from `install/desktop/theming.sh`
- ✅ Removed theme symlink creation (`~/.config/archriot/current/theme`)
- ✅ Updated hyprlock.conf to use direct CypherRiot config

#### Phase 4: Installer Updates ✅ COMPLETE

- ✅ Removed `setup_archriot_theme_system()` function
- ✅ Removed `set_default_theme()` function
- ✅ Removed `link_theme_configs()` function
- ✅ Removed `link_waybar_config()` function
- ✅ Simplified theming.sh to only handle cursors, icons, GTK themes
- ✅ Tested installation with simplified theme system

#### Phase 5: Documentation & Testing ✅ COMPLETE

- ✅ Updated README.md to remove theme selection instructions
- ✅ Updated INSTALLER.md to remove theme system warnings
- ✅ Tested all applications work with consolidated configs
- ✅ Verified waybar update notifications work
- ✅ Tested lock screen functionality
- ✅ Tested fuzzel launcher functionality

#### v2.0.1 Upgrade Path Fixes ✅ COMPLETE

- ✅ Added missing backgrounds directory to repository
- ✅ Fixed hyprlock.conf upgrade detection and automatic replacement
- ✅ Fixed background service startup after upgrade
- ✅ Verified upgrade path works for v1.9.x → v2.0.1

### 2. ✅ COMPLETED: Critical Installation Failures FIXED

**Status**: COMPLETED in v2.0.0-2.0.1 - All critical failures resolved
**Problem**: Multiple critical installation failures discovered during testing
**Impact**: Users cannot complete installation due to broken components

**FIXED ISSUES**:

- ✅ **Migrate tool download**: Fixed "text file busy" error with temp file approach and process killing
- ✅ **PipeWire audio system**: Fixed dependency conflicts with proper conflict detection and removal
- ✅ **Local variable errors**: Fixed all `local` declarations outside functions
- ✅ **Audio system**: Complete audio stack now installs properly (pavucontrol, pamixer, playerctl)
- ✅ **Waybar update notifications**: Eliminated complex caching system, direct version comparison works
- ✅ **Background system**: Fixed upgrade path, backgrounds install and display correctly
- ✅ **Lock screen**: Fixed hyprlock.conf transition from theme sourcing to consolidated config

### 3. 🔧 Fix Fuzzel Sudo Integration

**Status**: High Priority - Affects core user experience
**Problem**: Migrate command works in terminal but fails in Fuzzel launcher
**Impact**: Users cannot access migration tool from application launcher

**Actions**:

- [ ] Implement pkexec wrapper for migrate command
- [ ] Create GUI sudo helper for privileged operations
- [ ] Update desktop integration files
- [ ] Test fuzzel integration with sudo commands
- [ ] Verify migrate tool launches from application menu

### 4. 📊 Implement Modern Progress Bar System

**Status**: FUTURE ENHANCEMENT - Clean installer experience
**Problem**: Current installer output is verbose and overwhelming
**Impact**: Users can't see actual progress through installation noise

**Desired Experience**:

**Static progress bar**: Stays in place, updates percentage in real-time
**Module names**: Show friendly names ("Desktop Environment" not "install/desktop/apps.sh")
**Background logging**: All verbose output goes to log files
**Error surfacing**: Only show actual errors on console, not warnings
**Clean completion**: Beautiful summary at end

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

### 5. 🚀 System-wide v2.0.1 Deployment

**Status**: Medium Priority - Production rollout
**Problem**: Ensure all systems get latest fixes and improvements
**Impact**: Users need v2.0.1 fixes for proper upgrade path

**Actions**:

- [ ] Deploy v2.0.1 to all production systems
- [ ] Verify upgrade path works properly (v1.9.x → v2.0.1)
- [ ] Monitor for edge cases or regressions
- [ ] Collect user feedback on theme consolidation
- [ ] Verify background and lock screen functionality after upgrades
- [ ] Document any deployment issues

### 6. ✅ COMPLETED: Theme System Verification

**Status**: COMPLETED in v2.0.1 - Quality assurance passed
**Problem**: Need to verify theme consolidation doesn't break functionality
**Impact**: Ensure simplified theme system maintains all features

**COMPLETED VERIFICATION**:

- ✅ Tested waybar theming with unified config
- ✅ Verified lock screen themes work correctly
- ✅ Checked desktop theming elements
- ✅ Validated all color schemes and styling
- ✅ Updated theme documentation
- ✅ Tested upgrade path from v1.9.x systems
- ✅ Verified background system works after upgrade
- ✅ Confirmed update notifications appear properly

## CRITICAL BUG FIXES

### ✅ COMPLETED: Ivy Bridge Vulkan Support Issue

**Status**: FIXED with Safe Optional Tool - Available in optional-tools/ivy-bridge-vulkan-fix/
**Problem**: Intel HD Graphics 4000 (Ivy Bridge) has incomplete Vulkan support causing Zed editor failures
**Impact**: ThinkPad X230, T430 and similar systems now have working Zed editor
**User Report**: "I have to get vulkan working for zed but that's like one of the few things that didn't work" - RESOLVED

**SAFE SOLUTION IMPLEMENTED**:

- ✅ **Safe Fix Created**: `optional-tools/ivy-bridge-vulkan-fix/fix-ivy-bridge-vulkan.sh`
- ✅ **No System Changes**: Does NOT modify core ArchRiot files or hardware.sh
- ✅ **User-Specific Configs**: Creates only user-specific configuration overlays
- ✅ **Zed Compatibility**: Provides OpenGL fallback launcher and configuration
- ✅ **Easy Testing**: Completely reversible, testable solution
- ✅ **Available in Launcher**: Integrated into optional-tools launcher menu

**What the Fix Provides**:

```bash
# Safe launcher wrapper
~/.local/bin/zed-ivy-bridge

# Configuration overlay (manual merge required)
~/.config/zed/ivy-bridge-overlay.json

# Detailed logs
~/.cache/archriot/ivy-bridge-fix.log
```

**Installation**:

- Run: `optional-tools/ivy-bridge-vulkan-fix/fix-ivy-bridge-vulkan.sh`
- Or: Use optional-tools launcher menu option 2
- Result: Working Zed editor on Ivy Bridge systems without breaking other installations

## SYSTEMATIC CODE OPTIMIZATION & BUG FIXES

### 7. 🔧 File-by-File Optimization Initiative

**Status**: NEW PRIORITY - Systematic codebase improvement
**Problem**: Need comprehensive analysis of every file for optimization opportunities and bug fixes
**Impact**: Improve installer performance, reliability, and maintainability

**Approach**: Go through every single file, one at a time, proposing focused optimizations and bug fixes

#### Current Focus: install.sh Optimizations

**File**: `install.sh` (615 lines)
**Status**: ANALYSIS COMPLETE - Multiple optimization opportunities identified

**Identified Issues**:

1. **🔧 READY: Eliminate Redundant Temp File Operations**
    - **Problem**: `execute_module()` function creates temp files for every module execution
    - **Impact**: Unnecessary I/O overhead, temp file cleanup, multiple read/write operations
    - **Location**: Lines 213-230 in `execute_module()` function
    - **Fix**: Replace `mktemp` + `tee` + `cat` + `rm` with direct process substitution
    - **Benefit**: Faster module execution, cleaner I/O, reduced disk operations

2. **🔍 IDENTIFIED: Inefficient Multiple Tee Operations**
    - **Problem**: `log_message()` function uses multiple `tee` calls for each log level
    - **Impact**: Multiple process spawns for simple logging operations
    - **Location**: Lines 100-118 in `log_message()` function
    - **Fix**: Consolidate logging with single output redirection approach

3. **🔍 IDENTIFIED: Complex Module Discovery Logic**
    - **Problem**: `get_installation_modules()` has nested loops and complex path handling
    - **Impact**: Slower startup, harder to maintain module ordering
    - **Location**: Lines 40-55 in `get_installation_modules()` function
    - **Fix**: Simplify with direct array population and cleaner sorting

4. **🔍 IDENTIFIED: Hardcoded Path Redundancy**
    - **Problem**: Multiple hardcoded paths repeated throughout script
    - **Impact**: Maintenance overhead, potential for inconsistencies
    - **Location**: Throughout script (LOG_FILE, ERROR_LOG, INSTALL_DIR paths)
    - **Fix**: Centralize path management in configuration section

5. **🔍 IDENTIFIED: Missing Error Handling in Critical Paths**
    - **Problem**: Some operations lack proper error handling
    - **Impact**: Silent failures, unclear error states
    - **Location**: Font cache updates, icon cache updates (lines 410-415)
    - **Fix**: Add explicit error checking and logging for all critical operations

**Next Files for Analysis**:

- [ ] **PRIORITY 1**: `install/applications/productivity.sh` - Zed installation and config (77 lines)
- [ ] **PRIORITY 2**: `install/system/hardware.sh` - GPU driver installation (29 lines)
- [ ] **PRIORITY 3**: `setup.sh` - Secondary installer entry point
- [ ] `validate.sh` - System validation script
- [ ] `install/core/*.sh` - Core installation modules
- [ ] `install/system/*.sh` - System configuration modules
- [ ] `install/desktop/*.sh` - Desktop environment modules
- [ ] `install/applications/*.sh` - Application installation modules

**Optimization Principles**:

- One focused fix per interaction
- Measure performance impact before/after
- Maintain all existing functionality
- Improve error handling and logging
- Reduce I/O operations where possible
- Simplify complex logic without losing capability

#### Next Proposed Fix: productivity.sh Zed Configuration

**File**: `install/applications/productivity.sh`
**Issue**: Lines 55-85 contain redundant Zed setup operations
**Problem**: Multiple file operations and path checks that could be consolidated
**Fix**: Combine file operations, eliminate redundant path resolution, simplify config installation

## NEXT DEVELOPMENT PHASE

### Short-term Goals (Next 2 weeks)

- ✅ Complete theme system consolidation (DONE in v2.0.1)
- ✅ Fix Ivy Bridge Vulkan compatibility (DONE - Safe optional tool created)
- Fix fuzzel sudo integration
- Restore clean installer with progress bars
- Deploy v2.0.1 everywhere

### Medium-term Goals (Next month)

- Enhanced verification system
- Performance optimization
- Better error recovery tools
- Improved user onboarding
- Cross-platform installer testing

### Long-term Vision (Next quarter)

- Automated updates system
- Plugin architecture for extensibility
- Cross-distribution support
- Advanced customization options

## SUCCESS METRICS

### Theme System Success ✅ ACHIEVED

- ✅ Single unified theme config (no more overrides)
- ✅ No feature loss when adding new capabilities
- ✅ All theming functionality preserved (CypherRiot integrated)
- ✅ Massive maintenance overhead reduction (eliminated 1000+ lines of theme code)
- ✅ Upgrade path works correctly
- ✅ Background and lock screen function properly

### Installation Experience Success ✅ PARTIALLY ACHIEVED

- ⚠️ Progress bars need restoration (currently verbose output)
- ✅ All errors immediately visible (no more silent failures)
- ✅ Critical installation failures resolved
- ✅ User feedback positive for reliability improvements

### Integration Success ⚠️ PARTIALLY ACHIEVED

- ⚠️ Migrate tool launches from Fuzzel (needs pkexec wrapper)
- ✅ Desktop integration is seamless
- ✅ Update notifications work reliably
- ✅ Background cycling and management works

## NOTES

**Critical Learning**: Theme overrides are a maintenance nightmare. The discovery that theme configs completely replace main configs rather than extending them explains why features like update notifications disappeared. This must be the #1 priority to prevent future silent failures.

**Deployment Strategy**: v1.9.2 fixes are critical for production systems. Theme notification functionality is essential for user experience.

**Development Approach**: One focused change per iteration, always ask "Continue?" before modifications, never commit without explicit permission.
