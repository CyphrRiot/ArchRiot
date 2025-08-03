# ArchRiot Development Roadmap

## AI DEVELOPMENT RULES

THE FIRST AND MAIN RULE -- LISTEN CLOSELY AND REVIEW EVERYTHING. SLOW DOWN. BE METHODICAL. ONE THING AT A TIME. BE THOROUGH. WAIT FOR FEEDBACK.

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

**COMPLETED MIGRATIONS:**

1. **Scripts (config/bin/\*)** - ✅ COMPLETE
    - Step 1: ✅ Modified YAML to copy only `upgrade-system` (user CLI command)
    - Step 2: ✅ All internal scripts stay in `~/.local/share/archriot/config/bin/`
    - Step 3: ⚠️ PENDING - Must update all references in hyprland.conf, waybar, etc.

2. **tmux Config (config/default/tmux.conf)** - ✅ COMPLETE
    - Step 1: ✅ Changed YAML target to `~/.config/tmux/tmux.conf` (XDG location)
    - Step 2: ✅ tmux automatically finds config in XDG location
    - Step 3: ✅ No additional references to update

3. **Icons (config/applications/icons/\*)** - 🔄 IN PROGRESS
    - Step 1: ✅ Removed copying pattern from YAML
    - Step 2: 🔄 Updating X.desktop to use `%h/.local/share/archriot/config/applications/icons/`
    - Step 3: ⚠️ PENDING - Must update ALL desktop files and any other icon references

**PENDING MIGRATIONS:**

4. **Desktop Applications (config/applications/\*)**
    - Current: Copy .desktop files to ~/.local/share/applications/
    - Decision: KEEP COPYING (desktop environments require standard location)

5. **Images (config/images/\*)**
    - Current: Copy to ~/.local/share/archriot/images/
    - Decision: Already correct (staying in ArchRiot directory)

6. **Config Files (hypr/_, waybar/_, etc.)**
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

**CURRENT STATUS:** Following systematic 3-step pattern for each file type to ensure no functionality breaks.

### NEXT IMMEDIATE TASKS

**PHASE 1: ARCHITECTURE AUDIT (CRITICAL)**

1. **Complete file copying audit** using methodology above
2. **Identify all unnecessary copying** in current YAML configs
3. **Test application compatibility** with direct references
4. **Update all references systematically** (scripts, icons, images, etc.)

**PHASE 2: PENDING MIGRATIONS** 5. Convert remaining pending scripts to YAML modules: - communication.sh → communication module - theming.sh → theming module - productivity.sh → productivity module - utilities.sh → utilities module - specialty.sh → specialty module

**PHASE 3: VALIDATION** 6. **Audit previous migration work** for sloppy implementations 7. **Verify all 01-config.sh functionality** was properly migrated 8. **Check for missing or broken functionality** from hasty changes

### MIGRATION TARGET

Replace 30+ shell scripts with single `packages.yaml` configuration:

- `install/pending/*.sh` → YAML modules
- Maintain all functionality
- Improve reliability and maintainability

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

- `communication.sh` - Web applications (brave, signal, web app desktop files)
- `theming.sh` - Desktop themes (cursor, icon, GTK themes, backgrounds)
- `productivity.sh` - Office tools (text editor, zed with vulkan detection)
- `utilities.sh` - System utilities (btop, fastfetch, system tools)
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
