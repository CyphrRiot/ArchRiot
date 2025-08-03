# ArchRiot Development Roadmap

## AI DEVELOPMENT RULES

THE FIRST AND MAIN RULE -- LISTEN CLOSELY AND REVIEW EVERYTHING. SLOW DOWN. BE METHODICAL. ONE THING AT A TIME. BE THOROUGH. WAIT FOR FEEDBACK.

# AI Development Guidelines

You are an expert developer specializing in Go, Python, and Bash on Arch Linux. Create reliable, efficient code with proper documentation and error handling.

## CRITICAL WORKFLOW RULES

### Required Confirmations

- **ALWAYS ask "Continue?" before ANY file modification**
- **ALWAYS ask "Commit?" before ANY git commit**
- **NEVER use `git add -A` or `git add .`** - Add specific files only

## **IMPORTANT** -- RULES

READ THE RULES. READ THE ROADMAP.

1. Propose ONE SMALL DIRECT CHANGE AT A TIME.

2. IF I AGREE to the change, then make it.

3. After you've made it, wait for me to confirm and review the code and ask "Continue?" and wait for "yes" or further instructions.

4. Once it compiles properly, **WAIT FOR ME TO TEST AND REPORT BACK**

5. If there is ANY new info, update the ROADMAP so we stay on task and keep up-to-date

6. After each working, tested fix, ask "Commit?"

DO NOT SKIP AHEAD. DO NOT DEVIATE FROM THIS BEHAVIOR.

### Version Management

- **ALWAYS update CHANGELOG.md when bumping VERSION**
- **NEVER commit version changes without CHANGELOG entry**
- **Use semantic versioning in CHANGELOG.md**
- **CHANGELOG.md is the ONLY place for version history**

### Final Actions (Only if 100% confident)

**ASK FIRST, THEN:**

1. Update VERSION and README.md version
2. Update CHANGELOG.md
3. Commit with descriptive comment
4. Push changes

**NEVER COMMIT WITHOUT ASKING FIRST**

### SHELL SCRIPT MIGRATION

- **COMPLETE MIGRATION ONLY**: Migrate ALL functionality from shell scripts to YAML
- **NO PARTIAL IMPLEMENTATIONS**: Every package, command, config goes into YAML
- **READ ENTIRE SCRIPT FIRST**: Then implement everything in one edit
- **CHECK FOR DUPLICATES**: Always verify packages/commands don't already exist in YAML
- **CHECK EXISTING SYSTEM**: Use system commands to understand how packages should be configured and what dependencies exist - NOT to skip installation
- **100% CERTAINTY REQUIRED**: DO NOT COMPLETE ANY TASK UNTIL YOU ARE 100% CERTAIN IT WORKS CORRECTLY
- **NO BABYSITTING**: Don't require step-by-step guidance when source script has all info

### IMPLEMENTATION APPROACH

- **ONE COMPLETE EDIT**: Make full implementations, not incremental changes
- **FOLLOW EXISTING PATTERNS**: Use established YAML structure and Go patterns
- **NO SELECTIVE EDITING**: If it's in the script, it goes in YAML

### FILE COPYING MIGRATION PATTERN

**SYSTEMATIC 3-STEP APPROACH FOR EACH FILE TYPE:**

1. **Fix the copy** - Remove/modify copying pattern in YAML
2. **Fix the primary reference** - Update main config file that uses it
3. **Find all references** - Grep entire codebase and update ALL references to point to `~/.local/share/archriot/config/{location}`

**EXAMPLE PATTERN:**

- Step 1: Remove `applications/icons/*` copying from YAML
- Step 2: Update desktop files to use `%h/.local/share/archriot/config/applications/icons/`
- Step 3: Grep for all icon references and update every single one

**CRITICAL:** Never break functionality - all three steps must be completed atomically for each file type.

### COMMIT RULES

- **ASK BEFORE COMMITTING**: Never commit without explicit permission
- **UPDATE VERSION**: Update VERSION and README.md version when releasing
- **UPDATE CHANGELOG**: Document all changes in CHANGELOG.md

## CURRENT OBJECTIVE: YAML ARCHITECTURE (v2.5.0)

**STATUS**: YAML architecture implemented and working
**MAJOR MIGRATIONS COMPLETED**: 01-config.sh migration complete with file structure consolidation
**CURRENT TASK**: Script reference architecture decision and remaining pending script migrations

### COMPLETED MIGRATIONS

**✅ apps.sh → desktop.apps (COMPLETE)**

- All desktop applications, thumbnails, system tools migrated to YAML
- Thumbnail fix system implemented with proper PDF disabling
- Script installation, desktop applications, icons all handled via config patterns

**✅ 01-config.sh → Multiple YAML modules (COMPLETE)**

- Python dependencies → core.base (python, python-psutil + validation)
- All scripts → config/bin/ structure with automatic copying
- Desktop applications → config/applications/ with proper patterns
- Systemd services → core.base (version-check timer setup)
- Display manager removal + autologin → desktop.hyprland
- Git configuration → core.identity (aliases, pull.rebase, default branch)
- Environment setup → core.shell (tmux config)
- Dynamic progress tracking system implemented (per-module updates)
- GTK privacy settings moved from sed commands to proper config files
- Config preservation system via preserve_if_exists patterns

**✅ theming.sh → system.themes + system.backgrounds (COMPLETE)**

- 598 lines of complex bash reduced to 34 lines of clean YAML
- All 5 packages migrated: bibata-cursor-theme, kora-icon-theme, tela-icon-theme-purple-git, gnome-themes-extra, kvantum-qt5
- All gsettings commands migrated for proper GTK theming (dark theme, cursor, icons, window buttons)
- Background/wallpaper system simplified - eliminated redundant copying, direct path references
- Applications directory consolidated - hidden desktop files moved to main directory with hidden\_ prefix
- Fixed Go installer config pattern bug for /\* patterns with custom targets
- Desktop DPMS bug fixed - conditional DPMS only on laptops to prevent floating window positioning issues

**✅ communication.sh → Already covered by existing YAML (DELETED)**

- All packages (brave-bin, signal-desktop) already in desktop.apps
- All desktop files already handled by applications/\* pattern
- Script was completely redundant and removed

**✅ productivity.sh → desktop.editors + existing modules (COMPLETE)**

- Most functionality already existed: zed (development.tools), thunar (desktop.apps), papers (desktop.apps + mimetypes)
- Added missing packages: gnome-text-editor, abiword, papers, unzip, p7zip to desktop.apps
- Created new desktop.editors module for text editor configuration (gsettings, themes)
- Vulkan driver support ensured via system.hardware dependency on development.tools
- Complex Zed desktop integration simplified - existing patterns handle it
- Fixed major floating window positioning bug with cross-workspace window recovery

**✅ File Structure Consolidation (COMPLETE)**

- Moved applications/ → config/applications/ (desktop files, icons, hidden menu cleanup)
- Moved images/ → config/images/ (all project images, updated all references)
- Moved default/icons + default/plymouth → config/default/ (cursor themes, boot splash)
- Moved all bin/ → config/bin/ (all scripts and tools)
- Removed duplicate directories, updated all path references
- All files now in config/ structure for unified YAML copying

**CRITICAL ARCHITECTURE DECISION PENDING**

**⚠️ FILE COPYING VS DIRECT REFERENCE ARCHITECTURE**

- Current: Copy scripts, icons, backgrounds, images, desktop files (creates duplicates)
- Issue: Many files could be referenced directly from ~/.local/share/archriot/config/
- Decision: Audit ALL copied files to determine if direct reference is better
- Action needed: Systematic review of every copied asset
- Rationale: Eliminate unnecessary duplication, single source of truth

**✅ MIGRATION PROGRESS USING SYSTEMATIC PATTERN**

**KEEP COPYING (CORRECT DECISIONS):**

1. **Desktop Applications (config/applications/\*)**
    - Current: Copy .desktop files to ~/.local/share/applications/
    - Decision: KEEP COPYING (desktop environments require standard location)

2. **Config Files (hypr/_, waybar/_, etc.)**
    - Current: Copy to ~/.config/{app}/
    - Decision: KEEP COPYING (applications require standard XDG locations)

**AUDIT METHODOLOGY:**

For each file type:

1. **Identify current copying behavior** in YAML configs
2. **Check if applications support alternate paths** (CLI args, env vars, configs)
3. **Find all references** that would need updating (grep entire codebase)
4. **Evaluate pros/cons** of copying vs direct reference
5. **Test compatibility** with target applications
6. **Update references systematically** if direct reference is better

**COPYING VS REFERENCE DECISION MATRIX:**

**Copy when:**

- Application hardcodes config location (no alternatives)
- User customization expected (preserve_if_exists needed)
- Standard system locations required for compatibility

**Reference directly when:**

- Application supports alternate paths
- ArchRiot system files (not user-customizable)
- Eliminates duplication without breaking functionality
- Updates should be immediate (no copy lag)

**EXPECTED OUTCOMES:**

- Eliminate unnecessary file duplication
- Faster updates (no copying step)
- Cleaner filesystem (less scattered files)
- Single source of truth for ArchRiot assets
- Consistent architecture across all file types

**CURRENT STATUS:** Systematic 3-step pattern complete for script architecture. File copying architecture now optimized with minimal duplication.

### CRITICAL LESSONS LEARNED

**CRITICAL LESSONS LEARNED**

**❌ REPEATED MISTAKES BY AI:**

1. **False Claims**: Multiple times claimed "100% certainty" about completing work that wasn't done
2. **Sloppy Verification**: Said "fixed all references" without actually checking every single one
3. **Incomplete Audits**: Missed obvious issues like applications path patterns in YAML
4. **Relative Path Issues**: Claimed to remove all `../` paths but left several behind
5. **Invalid YAML Types**: Created non-existent type "Config" instead of checking valid types
6. **Doubled Path Patterns**: Added extra config/ prefixes causing source not found errors
7. **Assumption-Based Fixes**: Making changes without verifying the actual root cause
8. **Rush to Fix Symptoms**: Focusing on workarounds instead of finding root causes
9. **Cross-Workspace Ignorance**: Not considering that window fixes need to work across all workspaces
10. **Hidden File Consolidation Oversight**: Having settings but forgetting to install the actual packages

**✅ SUCCESSFUL PATTERNS:**

1. **Systematic 3-Step Approach**: Fix copy → Fix primary reference → Find ALL references
2. **Atomic Changes**: One change at a time with review before proceeding
3. **Actual Testing**: Running installer revealed real issues vs theoretical fixes
4. **Architecture Simplification**: Removing unnecessary copying and config noise
5. **Root Cause Analysis**: Finding DPMS as the actual cause of window positioning bugs
6. **File Consolidation**: Moving hidden apps to main directory with clear naming
7. **Direct Path References**: Eliminating unnecessary file duplication
8. **Multiple Verification Passes**: Going through scripts multiple times to catch missed items
9. **Cross-Workspace Bug Fixes**: Proper floating window recovery across all workspaces
10. **Modular Architecture**: Creating focused modules like desktop.editors for clean separation

### NEXT IMMEDIATE TASKS

**PHASE 1: PENDING SCRIPT MIGRATIONS**

1. Convert remaining pending scripts to YAML modules:
    - ✅ communication.sh → DELETED (redundant, already covered)
    - ✅ theming.sh → system.themes + system.backgrounds (COMPLETE)
    - ✅ productivity.sh → desktop.editors + existing modules (COMPLETE)
    - ✅ utilities.sh → desktop.utilities + existing modules (COMPLETE)
    - specialty.sh → specialty module

**PHASE 2: CRITICAL ISSUES REQUIRING IMMEDIATE ATTENTION**

2. **FRESH INSTALL VALIDATION (CRITICAL)**
    - **Missing core dependencies**: setup.sh installs git but NOT yay or go
    - **Dependency chain gaps**: Fresh Arch systems need base-devel, sudo setup
    - **Module execution order**: Some modules may execute before dependencies available
    - **Network requirements**: No validation that internet/AUR access works
    - **User permission issues**: sudo/wheel group setup happens mid-install

3. **FILE CLEANUP AND INTEGRITY (HIGH PRIORITY)**
    - **Stray backup files**: install/pending/install.sh.backup should be removed
    - **Broken script references**: utilities.sh still references deleted scripts
    - **Validation script outdated**: validate.sh references old file locations
    - **Documentation drift**: README/INSTALLER.md reference old script paths
    - **Remaining relative paths**: 1 remaining ../path reference (btop config comment)

4. **REMAINING SCRIPT MIGRATIONS (PENDING)**
    - **utilities.sh**: Contains essential system utilities and tools
    - **specialty.sh**: Specialized applications and configurations
    - **14 other pending scripts**: Various functionality not yet migrated
    - **Complex interdependencies**: Some scripts reference each other

**PHASE 3: COMPREHENSIVE TESTING AND VALIDATION**

5. **Fresh Install Testing**: Test complete installation on clean Arch Linux system
6. **Dependency validation**: Verify all required packages available in repos
7. **Module execution order**: Test dependency chain works correctly
8. **Error handling**: Ensure graceful failures when packages/repos unavailable
9. **Performance validation**: Verify reference-based approach performs well
10. **Documentation update**: Update all guides for new architecture

### MIGRATION TARGET

Replace 30+ shell scripts with single `packages.yaml` configuration:

- **✅ COMPLETED**: 4 major scripts migrated (communication.sh, theming.sh, productivity.sh, utilities.sh)
- **⏳ IN PROGRESS**: specialty.sh (1 major script remaining)
- **📋 PENDING**: 13 additional scripts in install/pending/ directory
- **🎯 GOAL**: Maintain all functionality while improving reliability and maintainability

### ESTIMATED SCOPE AND TIMELINE

**🚀 IMMEDIATE PRIORITY (1-2 sessions):**

- Fix fresh install dependency gaps (yay, go installation in setup.sh)
- Complete utilities.sh migration → utilities module
- Complete specialty.sh migration → specialty module
- Remove obsolete files and fix broken references

**⚙️ HIGH PRIORITY (2-3 sessions):**

- End-to-end fresh install testing on clean Arch system
- Fix remaining documentation drift and validation scripts
- Remove 14 remaining pending scripts (mostly duplicates/obsolete)
- Comprehensive module dependency validation

**📊 ESTIMATED COMPLETION:**

- **Core migrations**: 90% complete (3/5 major scripts done)
- **Architecture implementation**: 85% complete (YAML system working)
- **Bug fixes and cleanup**: 70% complete (major bugs found and fixed)
- **Documentation and testing**: 60% complete (needs fresh install validation)
- **OVERALL PROJECT**: ~80% complete, estimated 4-6 more sessions to completion

**🎯 SUCCESS CRITERIA:**

- Single-command fresh Arch Linux install works flawlessly
- All 30+ shell scripts replaced with YAML configuration
- Zero broken references or obsolete files
- Comprehensive documentation updated
- Performance validated on clean system

### CRITICAL FRESH INSTALL CONCERNS

**⚠️ DEPENDENCY GAPS IDENTIFIED:**

- **yay AUR helper**: Not installed by setup.sh, required for many packages
- **go compiler**: Not guaranteed present, needed to build installer
- **sudo configuration**: Happens mid-install, may cause permission issues
- **internet connectivity**: No validation before attempting package downloads

**🔧 EXECUTION ORDER ISSUES:**

- development.tools depends on system.hardware (Vulkan) - ✅ FIXED
- Some modules execute commands before packages fully installed
- User environment may not be properly configured early in process

**📁 FILE CLEANUP REQUIRED:**

- Remove install/pending/install.sh.backup
- Update validation scripts for new file locations
- Fix documentation references to old script paths
- Remove any remaining obsolete files or references

## YAML ARCHITECTURE

First principles:

- The installer is in ~/.local/share/archriot
- All files are installed into ~/.local/bin and ~/.config/{location} from the installer
- The installer should be in bin and MUST BE RUN from setup.sh
- Users will ALWAYS run `curl -fsSL https://ArchRiot.org/setup.sh | bash` for an install or upgrade
- This SHOULD clone the repo into ~/.local/share/archriot
- The relative path for the `configs` copy pattern is ~/.local/share/archriot/configs/{name} for the YAML
- The only files the user should be running for installs/upgrades are setup.sh (duh), the archriot-installer, and the packages.yaml file

### Structure

```yaml
core:
    base: { packages, configs, commands, depends }
    identity: { git configuration }
    shell: { terminal tools }

desktop:
    hyprland: { compositor and tools }
    apps: { desktop applications }

development:
    tools: { editors and utilities }
    helpers: { compilers and languages }

system:
    fonts: { system fonts }
    audio: { pipewire setup }
    power: { power management }

media:
    players: { media applications }
```

### Module Properties

- `packages`: Array of packages to install
- `configs`: Configuration file patterns to copy
- `commands`: Shell commands to execute
- `depends`: Module dependencies
- `start/end`: Progress messages
- `type`: Module category

## REMAINING WORK

### Pending Scripts to Convert

**REMAINING TO MIGRATE:**

- ✅ `communication.sh` - DELETED (redundant, already covered)
- ✅ `theming.sh` - Desktop themes (cursor, icon, GTK themes, backgrounds)
- ✅ `productivity.sh` - Office tools (text editor, zed with vulkan detection)
- ✅ `utilities.sh` - System utilities (btop, fastfetch, system tools)
- `specialty.sh` - Specialized tools (papers, specialty applications)

**MIGRATION APPROACH:**

- Read entire script first, identify all packages/commands/configs
- Check for duplicates in existing YAML modules
- Use config patterns instead of shell commands where possible
- Preserve all functionality, no partial implementations

### Implementation Files

- `packages.yaml` - Main configuration
- `main.go` - Entry point
- `orchestrator/` - Execution logic
- `config/` - YAML loading
- `executor/` - Module execution
- `installer/` - Package/config handling

## BENEFITS ACHIEVED

- Single configuration source
- Proper dependency management
- Professional TUI with progress tracking
- Comprehensive logging
- Batch package installation
- Intelligent config preservation
- Module-based architecture
- Validation and error handling

## NEXT STEPS

**IMMEDIATE (CRITICAL):**

1. **Fix script architecture** - Remove copying commands, update references to config/bin/
2. **Continue pending migrations** - Convert remaining 5 scripts to YAML modules
3. **Remove obsolete shell scripts** - Delete completed migration scripts
4. **Update documentation** - Reflect new architecture and completed migrations
5. **Final testing and validation** - Ensure all functionality works with new structure

**ARCHITECTURE BENEFITS ACHIEVED:**

- Single YAML configuration source for all installations
- Proper dependency management between modules
- Dynamic progress tracking with per-module updates
- Config preservation system (preserve_if_exists patterns)
- Unified file structure in config/ directory
- Elimination of duplicate sed/echo script generation
- Professional TUI with comprehensive logging
- Module-based architecture with clear separation
