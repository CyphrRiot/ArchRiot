# ArchRiot Development Roadmap

## Current Status: v2.0.17

- Installation system reliability: ✅ Complete
- Theme config nightmare: ✅ ELIMINATED
- Critical installation failures: ✅ FIXED
- Plymouth upgrade protection: ✅ FIXED
- Ivy Bridge Vulkan compatibility: ✅ FIXED

## COMPLETED ARCHITECTURAL OVERHAULS

### ✅ COMPLETED: Control Panel Architecture Overhaul

**Status**: ✅ COMPLETED IN v2.0.18 - All fundamental issues resolved
**Problem**: Control Panel used redundant config system that conflicted with system configs
**Impact**: Fixed settings persistence, eliminated config conflicts, improved user experience

**COMPLETED IMPLEMENTATION**:

**✅ New Architecture Implemented**:

- **Eliminated** redundant `archriot-config` dependencies for system components
- **Direct system config management**: Control Panel reads/writes actual system files
- **Single source of truth**: Each setting has ONE authoritative location
- **Preserved config system**: For ArchRiot-specific functionality (Pomodoro timer)

**COMPLETED CONVERSIONS**:

1. **BlueLightWidget**: Direct hyprland.conf management - reads/writes actual temperature settings
2. **DisplayWidget**: Direct hyprctl integration with smart aspect ratio resolution detection
3. **MullvadWidget**: Direct mullvad CLI integration - account numbers from actual VPN service
4. **PowerWidget**: Direct powerprofilesctl management (tested: balanced ↔ power-saver)
5. **AudioWidget**: Already using direct pactl system management (verified)
6. **CameraWidget**: Already using direct v4l2/device system management (verified)
7. **PomodoroWidget**: Preserved config system for ArchRiot-specific functionality

**SMART RESOLUTION DETECTION**:

- Intelligent aspect ratio detection prevents display distortion
- Supports 16:9, 16:10, 4:3, 21:9 ultrawide, 32:9 super ultrawide
- Eliminates nonsensical resolutions like 1920x1440

**SUCCESS CRITERIA ACHIEVED**:

- ✅ Control Panel changes persist across sessions/reboots
- ✅ No config conflicts between Control Panel and system
- ✅ User settings preserved during ArchRiot updates
- ✅ Simplified architecture with single source of truth per setting
- ✅ Smart resolution detection prevents display distortion
- ✅ Real-time system integration working perfectly

## IMMEDIATE PRIORITIES

### 1. ✅ COMPLETED: Theme System Optimization & Redundancy Elimination

**Status**: ✅ COMPLETED IN v2.0.16-2.0.17 - All critical issues resolved
**Problem**: Theme system had multiple redundant operations causing confusion and performance issues
**Impact**: Fixed slower installations, eliminated duplicate theme setups, resolved confusing error messages

**COMPLETED FIXES**:

- ✅ **Eliminated "Cannot find theme directory" false positives**
    - Updated validate.sh to check consolidated backgrounds directory instead of old theme structure
    - Theme verification now properly validates consolidated system (23 backgrounds found)
    - Zero validation errors during normal installation

- ✅ **Consolidated Hyprland restart operations**
    - Audited entire codebase - only single hyprctl reload exists in install.sh (already optimal)
    - No mid-installation Hyprland restarts disrupting user experience
    - Single restart maximum during entire installation

- ✅ **Optimized background service management**
    - Removed redundant swaybg startup calls in theming.sh setup_backgrounds()
    - Single background service initialization after all theme components ready
    - Background service starts exactly once and persists correctly

- ✅ **Streamlined theme setup workflow**
    - All theme operations consolidated in single theming.sh file
    - Updated swaybg-next for consolidated flat directory structure
    - Updated hyprlock-wrapper.sh to use consolidated system without symlinks
    - Eliminated all old theme structure dependencies

- ✅ **Implemented installation state management**
    - Created ~/.config/archriot/.install-state tracking system
    - Tracks cursor_theme, icon_theme, gtk_themes, backgrounds completion
    - Skips redundant operations when components already configured
    - Faster re-runs (10 seconds vs 90+ seconds for repeated installs)

- ✅ **Enhanced Fuzzel application launcher**
    - Maintained Kora icon theme for comprehensive application icon coverage
    - Added automatic cleanup of unwanted desktop entries (lstopo removed)
    - Clean application launcher with proper icon display

**TECHNICAL IMPLEMENTATION**:

1. ✅ **Theme State Tracking**: ~/.config/archriot/.install-state tracks component completion
2. ✅ **Background Service Management**: Single swaybg initialization, state file tracking
3. ✅ **Validation System**: Consolidated system validation with accurate reporting
4. ✅ **Directory Structure**: Flat ~/.config/archriot/backgrounds/ with 23 working backgrounds
5. ✅ **Icon Management**: Strategic Kora + Tela-purple-dark combination for optimal coverage

**SUCCESS CRITERIA ACHIEVED**:

- ✅ Zero "Cannot find theme directory" errors during normal installation
- ✅ Hyprland restarts maximum once during entire installation
- ✅ Background service starts exactly once and persists correctly
- ✅ No duplicate theme setup operations
- ✅ Faster installation with cleaner progress indication
- ✅ Consistent theme state across fresh installs and upgrades
- ✅ Clean application launcher without confusing entries

### 2. 🚨 CRITICAL: Process Detachment Audit - IN PROGRESS

**Status**: CRITICAL PRIORITY - System stability issue
**Problem**: Background processes not properly detached from installer terminal sessions
**Impact**: Services disappear when installer GUI terminal closes, breaking user experience

**Root Cause Discovered**: Today we found swaybg (background service) was never properly detached with `nohup` and `disown`. Only waybar was fixed in v1.6.2. This suggests there may be OTHER services with the same issue.

**URGENT ACTIONS**:

- [ ] **AUDIT EVERY SINGLE FILE** in the installation process for background process starts
- [ ] **Search for ALL instances** of `&` without proper `nohup` and `disown`
- [ ] **Identify ALL services** that should persist after installer terminal closes
- [ ] **Fix ALL process detachment issues** found during audit
- [ ] **Test ALL background services** survive terminal closure
- [ ] **Document proper patterns** for future background service starts

**Files to Audit** (MUST check every single one):

- [ ] `install.sh` - Main installer
- [ ] `setup.sh` - Secondary installer
- [ ] `install/core/*.sh` - Core system setup
- [ ] `install/system/*.sh` - System configuration
- [ ] `install/desktop/*.sh` - Desktop environment (theming.sh FIXED)
- [ ] `install/applications/*.sh` - Application installation
- [ ] `bin/*` - All utility scripts (swaybg-next FIXED)
- [ ] Any other scripts that start background processes

**Search Patterns to Find**:

- `.*&$` - Lines ending with & (background processes)
- `systemctl.*start` - Service starts that may need detachment
- `pkill.*sleep.*command.*&` - Kill/restart patterns
- Any daemon or service startup commands

**Success Criteria**:

- ALL background services properly detached with `nohup` and `disown`
- NO services disappear when installer terminal closes
- Comprehensive documentation of proper background process patterns
- Zero process orphaning or service loss during installation

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

### 5. 🚀 System-wide v2.0.5 Deployment

**Status**: Medium Priority - Production rollout
**Problem**: Ensure all systems get latest fixes and improvements
**Impact**: Users need v2.0.5 fixes for Plymouth protection and Ivy Bridge support

**Actions**:

- [ ] Deploy v2.0.5 to all production systems
- [ ] Verify upgrade path works properly (v1.9.x → v2.0.5)
- [ ] Monitor for edge cases or regressions
- [ ] Collect user feedback on theme consolidation
- [ ] Verify background and lock screen functionality after upgrades
- [ ] Document any deployment issues

## COMPLETED MAJOR FIXES

### v2.0.16-2.0.17: Theme System & Fuzzel Optimization ✅

- ✅ **Complete Theme System Optimization** - Eliminated all redundant operations and false errors
- ✅ **Installation State Management** - Smart component tracking prevents duplicate operations
- ✅ **Background System Consolidation** - Single service initialization, 23 backgrounds working
- ✅ **Fuzzel Application Launcher** - Clean interface with proper Kora icon coverage
- ✅ **Validation System Update** - Accurate consolidated system validation
- ✅ **Desktop Entry Cleanup** - Automatic removal of unwanted system utilities

### v2.0.15: Critical Installation Fixes ✅

- ✅ **Media Applications Fix** - Resolved "print_status: command not found" errors
- ✅ **Spotdl Installation** - Proper --nocheck functionality for problematic packages

### v2.0.5: Core System Stability ✅

- ✅ **Plymouth Theme Protection** - upgrade-system now preserves custom themes during package updates
- ✅ **Ivy Bridge Vulkan Support** - Safe optional tool for ThinkPad X230 and similar systems
- ✅ **Theme System Consolidation** - Eliminated config override nightmare (v2.0.0-2.0.1)
- ✅ **Critical Installation Failures** - All major installer bugs resolved (v2.0.0-2.0.1)
- ✅ **Missing Plymouth Files** - Fixed logo.png and installer detection issues

## SYSTEMATIC CODE OPTIMIZATION & BUG FIXES

### 6. 🔧 File-by-File Optimization Initiative

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

- Fix fuzzel sudo integration
- Restore clean installer with progress bars
- Deploy v2.0.5 everywhere
- Complete systematic file optimization initiative

### Medium-term Goals (Next month)

- Enhanced verification system
- Better error recovery tools
- Improved user onboarding
- Cross-platform installer testing

### Long-term Vision (Next quarter)

- Automated updates system
- Plugin architecture for extensibility
- Cross-distribution support
- Advanced customization options

## SUCCESS METRICS

### Current Success Metrics

### Installation Experience Success ⚠️ PARTIALLY ACHIEVED

- ⚠️ Progress bars need restoration (currently verbose output)
- ✅ All errors immediately visible (no more silent failures)
- ✅ Critical installation failures resolved
- ✅ Plymouth themes protected during upgrades
- ✅ User feedback positive for reliability improvements

### Integration Success ⚠️ PARTIALLY ACHIEVED

- ⚠️ Migrate tool launches from Fuzzel (needs pkexec wrapper)
- ✅ Desktop integration is seamless
- ✅ Update notifications work reliably
- ✅ Background cycling and management works
- ✅ Ivy Bridge systems have working Zed editor

## NOTES

**Critical Learning**: Theme overrides are a maintenance nightmare. The discovery that theme configs completely replace main configs rather than extending them explains why features like update notifications disappeared. This must be the #1 priority to prevent future silent failures.

**Deployment Strategy**: v1.9.2 fixes are critical for production systems. Theme notification functionality is essential for user experience.

**Development Approach**: One focused change per iteration, always ask "Continue?" before modifications, never commit without explicit permission.
