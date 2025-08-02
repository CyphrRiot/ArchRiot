# ArchRiot Development Roadmap

## 🎯 NEXT: YAML ARCHITECTURE TRANSITION (v2.5.0)

## **IMPORTANT** -- RULES

READ THE RULES. READ THE ROADMAP.

1. Propose ONE SMALL DIRECT CHANGE AT A TIME.

2. IF I AGREE to the change, then make it.

3. After you've made it, wait for me to confirm and review the code and ask "Continue?" and wait for "yes" or further instructions.

4. Once it compiles properly, **WAIT FOR ME TO TEST AND REPORT BACK**

5. If there is ANY new info, update the ROADMAP so we stay on task and keep up-to-date

6. After each working, tested fix, ask "Commit?"

DO NOT SKIP AHEAD. DO NOT DEVIATE FROM THIS BEHAVIOR.

### PLAN

**TARGET**: Replace 30+ scattered .sh files with unified YAML-driven installation system
**STATUS**: MAJOR PIVOT - Moving from Bash to Go for performance and reliability
**SAFETY**: All development in separate branch until fully tested and approved

### 🔄 CURRENT PROGRESS (v2.5 BRANCH)

**✅ COMPLETED TUI ARCHITECTURE FIXES:**

- ✅ **v2.5 branch created**: Proper tracking for TUI fixes and architecture changes
- ✅ **TUI threading model fixed**: TUI now owns main thread, installation runs in goroutine (no more race conditions)
- ✅ **Responsive terminal sizing implemented**: Width/height calculations based on actual terminal size
- ✅ **Package installation spam removed**: Clean TUI showing only module-level operations, not individual package details
- ✅ **Command output truncation**: Limited to 200 chars to prevent massive dumps
- ✅ **Mirror fixing bullshit eliminated**: Removed hanging reflector commands, graceful pacman sync failure
- ✅ **Screen clearing issues resolved**: Removed WithAltScreen() to preserve output, added proper spacing
- ✅ **Input field system implemented**: Full
  user input handling for git credentials and reboot prompt
- ✅ **Beautiful Tokyo Night TUI**: ASCII art, proper styling, bordered scroll window working

**✅ COMPLETION INTERACTION FIXES:**

- ✅ **Beautiful reboot buttons**: Single-line YES/NO buttons below scroll window with arrow key navigation
- ✅ **Fixed button selection**: Clear visual indication with [►YES◄] and [ NO ] formatting
- ✅ **Proper quit handling**: tea.Quit instead of os.Exit(0) for clean TUI shutdown
- ✅ **No screen replacement**: Content stays visible, buttons appear below scroll window
- ✅ **Responsive layout**: All content fits within terminal bounds without scrolling off
- ✅ **Default to NO**: Safe default selection for reboot prompt

**🎉 TUI COMPLETION SUCCESS:**

- ✅ **Beautiful installation progress**: Tokyo Night themed interface with scroll window
- ✅ **Clean completion interaction**: Reboot prompt with arrow key selection
- ✅ **All content preserved**: Installation log stays visible throughout
- ✅ **Proper exit handling**: Clean shutdown without hanging
- ✅ **Perfect user experience**: Professional TUI matching Migrate quality

**CORE FUNCTIONALITY STATUS:**

- ✅ **Package installation working**: YAML-driven batch installation functional
- ✅ **Config copying working**: Preservation directives and file copying implemented
- ✅ **Module execution working**: Priority-based execution (core→development→desktop→media)
- ✅ **Logging system working**: File-based logging with proper error handling
- ✅ **No hanging commands**: Eliminated mirror fixing, graceful error handling

**🎯 IMMEDIATE PRIORITY: MODULAR REFACTORING**

**✅ MODULAR REFACTORING COMPLETE!**

**COMPLETED PACKAGES:**

1. **✅ tui/**: Terminal UI components (model, view, update, buttons)
2. **✅ config/**: YAML parsing and configuration loading (types.go, loader.go)
3. **✅ logger/**: Logging system with file output
4. **✅ git/**: Git credential handling with TUI integration
5. **✅ installer/**: Package installation and config copying logic
6. **✅ version/**: Version file reading and management
7. **✅ executor/**: Module execution logic with dependency handling
8. **✅ orchestrator/**: Main installation orchestration and workflow

**REFACTORING ACHIEVEMENTS:**

- **✅ MAIN.GO REDUCED**: From 810+ lines to 98 lines (88% reduction!)
- **✅ CLEAN ARCHITECTURE**: Each package has single responsibility
- **✅ PROPER INITIALIZATION**: Fixed all initialization order issues
- **✅ NO UNUSED CODE**: Removed PackageResult struct and unused constants
- **✅ MODULAR DESIGN**: Easy to maintain and extend each component

**🎯 NEXT STEPS: YAML MIGRATION**

**CURRENT FOCUS: SHELL SCRIPT TO YAML MIGRATION**

**✅ COMPLETED SHELL SCRIPTS:**

1. **✅ 01-base.sh**: Converted to `core.base` module in packages.yaml
    - Base development tools (base-devel, git, rsync, bc)
    - Eliminated problematic mirror fixing and manual yay installation
    - Uses native package manager installation

2. **✅ 02-identity.sh**: Converted to `core.identity` module and Go `git/` package
    - Git credential handling with TUI integration
    - Interactive user input replaced with native Go TUI
    - Git configuration and aliases applied automatically

3. **✅ fonts.sh**: Converted to `system.fonts` module in packages.yaml
    - Essential fonts only: font-awesome, noto-fonts, hack-nerd, cascadia-mono-nerd, ia-writer
    - Removed unused font packages (emoji, cjk, liberation, jetbrains)
    - Font cache refresh handled by handler function

4. **✅ bluetooth.sh**: Converted to `system.bluetooth` module in packages.yaml
    - Simple package installation: blueberry
    - Service enablement handled by handler function
    - Bluetooth management with GUI interface

5. **✅ printer.sh**: Converted to `system.printer` module in packages.yaml
    - CUPS printing system packages: cups, cups-pdf, cups-filters, system-config-printer
    - Service enablement handled by handler function
    - Complete printing solution with GUI configuration

6. **✅ power.sh**: Converted to `system.power` module in packages.yaml
    - Power profiles daemon (provides powerprofilesctl) and powertop monitoring
    - Handler manages service enablement and conditional battery detection
    - Modern power management with profile switching support

7. **✅ filesystems.sh**: Converted to `system.filesystems` module in packages.yaml
    - Complete filesystem support: gvfs, udisks2, SMB, MTP, NTFS, FAT32, exFAT
    - GUI mounting and network shares with udisks service enablement
    - Universal compatibility for USB drives, SD cards, and mobile devices

8. **✅ media.sh**: Converted to expanded `media.players` module in packages.yaml
    - Media tools: ffmpeg, yt-dlp, python-opencv (webcam for Control Panel), spotdl
    - Handler manages spotdl special installation requirements
    - Comprehensive media consumption and creation support

9. **✅ tools.sh**: Converted to `development.helpers` module in packages.yaml
    - Development compilers: clang, cmake, ninja, rust, python-pip, go
    - CLI utilities: github-cli, jq, curl, wget, zip tools, tmux
    - Handler manages migrate tool installation

10. **✅ networking.sh**: Converted to `system.networking` module in packages.yaml
    - Wireless networking with iwd package
    - Service enablement handled by handler function
    - Works for both fresh installs and system updates

11. **✅ 04-shell.sh**: Converted to expanded `core.shell` module in packages.yaml
    - Complete terminal toolchain: fish, kitty, modern CLI tools, file search
    - Handler manages fish shell configuration and default shell setup
    - Established design pattern for complex script conversion

12. **✅ audio.sh**: Converted to `system.audio` module in packages.yaml
    - Complete PipeWire audio stack: pipewire, pipewire-alsa, pipewire-pulse, wireplumber
    - Audio controls: pavucontrol, pamixer, playerctl
    - Handler manages conflict resolution (removes pulseaudio/jack2) and service enablement
    - Tested: Successfully detects missing packages and executes handler functions

13. **✅ memory.sh**: Converted to `system.memory` module with clean config pattern
    - Created `config/system/99-memory-optimization.conf` with memory optimization settings
    - Uses YAML config copying with `target: /etc/sysctl.d/...` specification
    - Simple handler applies copied config with `sysctl -p`
    - **LESSON**: Configuration content belongs in config files, NOT hardcoded in Go handlers

14. **✅ mimetypes.sh**: Converted to `system.mimetypes` module in packages.yaml
    - Default application associations for all major file types
    - Handler manages xdg-mime commands for images, videos, PDFs, text, and browser
    - Comprehensive file type coverage: PNG→imv, MP4→mpv, PDF→Papers, HTTP→brave
    - Tested: Handler executes successfully and MIME associations are applied correctly

15. **✅ editors.sh**: Converted to `development.editors` module with YAML command system
    - **MAJOR ARCHITECTURAL BREAKTHROUGH**: Implemented YAML command system
    - Neovim + LSP servers: lua, pyright, typescript, bash, gopls, tree-sitter-cli
    - Commands field replaces hardcoded Go handlers: `sudo ln -sf /usr/bin/nvim /usr/bin/vi`
    - Comprehensive error handling: command failure detection, detailed logging, error propagation
    - **ACHIEVEMENT**: System configuration now in YAML, not hardcoded in Go code

**🎯 YAML COMMAND SYSTEM ARCHITECTURE (MAJOR BREAKTHROUGH):**

**✅ COMMANDS FIELD IMPLEMENTED:**

- **YAML-defined commands**: No more hardcoded system commands in Go handlers
- **Error handling**: Command failure detection with detailed logging
- **Command tracking**: Individual command logging with success/failure status
- **Module integration**: Commands execute after packages/configs, before handlers
- **Pattern established**: `commands: ["sudo ln -sf /usr/bin/nvim /usr/bin/vi"]`

**COMMAND EXTRACTION RULES:**

1. **NEVER hardcode system commands in Go handlers**
2. **ALWAYS use YAML commands field for shell operations**
3. **EXTRACT existing hardcoded commands** to YAML systematically
4. **TEST each extraction** to ensure error handling works
5. **ELIMINATE simple handlers** that only run commands

**NEXT PHASE: SYSTEMATIC COMMAND EXTRACTION**

Extract all hardcoded commands from existing handlers:

- Service enablement: `systemctl enable --now service`
- Cache updates: `fc-cache -f`, `update-desktop-database`
- System commands: `sysctl -p`, `xdg-mime`, `xdg-settings`
- Package operations: yay installs, curl downloads
- File operations: chmod, chown, mkdir

**🎯 HANDLER SYSTEM ARCHITECTURE (CRITICAL - DO NOT REPEAT MISTAKES):**

**✅ ABSTRACTED HANDLER SYSTEM IMPLEMENTED:**

- **NEVER write 300+ line handler files again!**
- **Use abstractions**: `enableService()`, `runCommand()`, `sendLog()`
- **Registry pattern**: Simple 1-3 line handlers in map

5. **Only complex logic gets full functions** (like setupFishShell)
6. **Reduced 350+ lines to 140 lines** (60% reduction)

**CRITICAL CONFIG PATTERN RULE:**

- **NEVER hardcode configuration content in Go handlers**
- **ALWAYS use config files + YAML copying pattern**
- **Keep system configurations separate from code**

**SYSTEM ARCHITECTURE PATTERNS (STOP BEING STUPID):**

**CONFIG COPYING FLOW:**

1. **Repository**: `ArchRiot/config/system/file.conf` (development)
2. **User Install**: `~/.local/share/archriot/config/system/file.conf` (archriot-installer copies here)
3. **Config Copy**: `~/.config/system/file.conf` (YAML `configs:` copies here)
4. **System Copy**: `/etc/system/file.conf` (commands copy here if needed)

**YAML PATTERNS:**

- `packages: []` - Install packages
- `configs: [pattern: system/file.conf]` - Copy from ~/.local/share/archriot/config/ to ~/.config/
- `commands: []` - Run shell commands AFTER packages and configs
- `depends: []` - Module dependencies

**NEVER FORGET:**

- When developing: we're in `~/Code/ArchRiot/` (repository)
- When users run installer: configs are in `~/.local/share/archriot/config/` (installed)
- Config copying goes: `~/.local/share/archriot/config/` → `~/.config/`
- Commands reference copied files in `~/.config/` or final system locations

**HANDLER DESIGN RULES:**

1. **Simple handlers**: Use abstractions in registry map
2. **Complex handlers**: Separate function only when needed
3. **Abstract common patterns**: service enablement, commands, logging
4. **Keep registry clean**: Max 3 lines per handler
5. **Document in code**: Clear examples for future developers

**PACKAGE DUPLICATION PREVENTION RULES:**

- **NEVER install the same package in multiple modules**
- **ALWAYS check existing modules before adding packages**: Use `grep -r "package_name" packages.yaml`
- **USE logical organization**: core→development→desktop→system→media
- **VERIFY dependencies**: modules inherit packages from dependencies
- **CLEAN UP immediately**: remove duplicates as soon as discovered

**PACKAGE ORGANIZATION GUIDELINES:**

1. **core.base**: Essential system packages (base-devel, git)
2. **core.shell**: Shell environment only (fish, terminal tools, basic utilities)
3. **development.tools**: Development applications (zed, btop, fastfetch, tree)
4. **development.helpers**: Development utilities (wget, curl, compilers, build tools)
5. **development.editors**: Editor-specific packages (LSP servers, tree-sitter)
6. **desktop.\***: Desktop environment packages (hyprland, apps)
7. **system.\***: System configuration packages (audio, fonts, networking)
8. **media.\***: Media applications (players, tools)

**DUPLICATION DETECTION PROCESS:**

Before adding any package:

1. Run: `grep -r "package_name" install/packages.yaml`
2. Check if package already exists in any module
3. If exists, verify it belongs in the correct logical category
4. If adding to wrong category, move it instead of duplicating
5. Update dependencies if module relationships change

**EXAMPLE - GOOD HANDLER:**

```
"enable_bluetooth_service": func() error {
    enableService("bluetooth")
    sendLog("📶", "Bluetooth", "Service enabled")
    return nil
},
```

**EXAMPLE - BAD HANDLER:**

- 25+ lines of repeated systemctl commands
- Manual TUI message formatting
- Copy-paste error handling code

**📋 TODO: GENERIC FILE COPY HANDLER SYSTEM**

**NEEDED**: Create abstracted file copy handler pattern like config copying but for system files:

- `system_files` section in YAML for `/etc/`, `/usr/`, etc.
- Generic `copySystemFile()` helper function
- Proper permissions and ownership handling
- Validation and backup capabilities

This would eliminate the need for custom handlers that just copy files to system locations.

**🔄 IN PROGRESS: Command Extraction & Script Analysis**

1. **🔄 IMMEDIATE PRIORITY: Extract hardcoded commands from existing handlers**
    - Move service enablement commands to YAML
    - Extract cache update commands (fc-cache, update-desktop-database)
    - Move system configuration commands (sysctl, xdg-mime, xdg-settings)
    - Replace custom package installations with YAML commands

2. **📋 REMAINING: 18 .sh files** in install/pending/ awaiting conversion
    - Need to categorize by complexity (simple package lists vs complex logic)
    - Identify patterns for YAML structure optimization

3. **📋 TODO: YAML Module Expansion**
    - Add missing categories: system, applications, optional
    - Implement conditional installation based on hardware/user choice
    - Add pre/post installation hooks where needed

4. **📋 TODO: Handler Functions**
    - GPU detection and driver selection
    - yay AUR helper installation
    - Service configuration and enablement
    - Hardware-specific configurations

5. **📋 TODO: Testing & Release**
    - Comprehensive testing on clean Arch systems
      12.5.0 release preparation

**CRITICAL LESSONS LEARNED:**

- ✅ **ALWAYS TEST FIRST**: External testing prevented major integration failures
- ✅ **COPY WORKING PATTERNS**: Migrate's simple cursor approach > complex custom systems
- ❌ **DON'T OVERCOMPLICATE**: Channels and async messaging caused hanging issues
- ❌ **WATCH SYSTEM CHANGES**: Accidentally overwrote user's git config during testing
- ✅ **FOLLOW THE RULES**: One change at a time prevents debugging nightmares
- ✅ **CONSISTENT FORMATTING**: Emoji alignment and spacing matter for professional appearance
- ✅ **MODULAR LOGGING**: Created sendFormattedLog() function for consistent message formatting
- ✅ **GIT CREDENTIAL FIXED**: Git YES/NO confirmation system working! Fixed initialization order (program was nil)
- ✅ **LOG ALIGNMENT**: Fixed misaligned completion messages by using sendFormattedLog consistently
- ✅ **CLEAN ORGANIZATION**: Moved all shell scripts to pending/ directory for systematic conversion tracking
- ✅ **MAJOR SIZE REDUCTION**: Reduced main.go from 810+ lines to 98 lines (88% reduction!)
- ✅ **MODULAR ARCHITECTURE WORKING**: Fixed channel communication by proper initialization order
- ✅ **CONFIG PATTERN ESTABLISHED**: Never hardcode system configuration in handlers - use config files
- ✅ **PACKAGE DUPLICATION ELIMINATED**: Systematic cleanup preventing installation redundancy
- ✅ **PACKAGE ORGANIZATION RULES**: Established logical categories and duplication detection process
- ✅ **YAML COMMAND SYSTEM**: Major architectural breakthrough - system commands in YAML, not Go code
- ✅ **COMMAND ERROR HANDLING**: Comprehensive failure detection, logging, and error propagation
- ✅ **HANDLERS ELIMINATED**: Entire handlers.go package deleted - 100% YAML-driven commands
- ✅ **CONFIG PATTERNS DOCUMENTED**: Explicit flow to prevent future confusion about paths
- ✅ **GO LESSONS LEARNED**: Always initialize variables before passing to other packages
- ✅ **CLEAN CODE ACHIEVED**: Removed all unused code, proper package separation
- ✅ **SYSTEMATIC APPROACH**: One change at a time prevented breaking working systems

**ARCHITECTURE INSIGHTS:**

- **Go Benefits Confirmed**: Native YAML, static typing, compiled performance, beautiful TUI
- **Modular Design Challenges**: Package communication and channel sharing requires careful design
- **Basic Go Knowledge Required**: Exports, channels, and package boundaries must be properly understood
- **Test Everything**: Each modular extraction must be tested immediately to prevent breaking working systems
- **Current State**: MODULAR REFACTORING COMPLETE - Ready for YAML migration phase

### 📋 CORRECT YAML PLAN: Actually Eliminate Shell Scripts

#### OBJECTIVE: Replace (not supplement) 30+ .sh files with YAML processing

**Current Reality Check:**

- We have 30+ shell scripts doing repetitive package installation + config copying
- Most scripts follow identical patterns with different package lists
- YAML can handle 90% of this with proper structure and handlers
- Go installer can process YAML natively without external dependencies
- Result: Single compiled binary replaces entire shell script ecosystem

#### YAML STRUCTURE (packages.yaml):

```yaml
core:
    base:
        packages: [list]
        configs: [pattern-based copying]
        handler: special_logic_function
    identity:
        packages: [list]
        configs: [pattern-based copying]
        depends: [core.base]

desktop:
    hyprland:
        packages: [list]
        configs: [pattern-based copying]
        depends: [core.base]
    apps:
        packages: [list]
        configs: [pattern-based copying]
        depends: [desktop.hyprland]

development:
    tools:
        packages: [list]
        configs: [pattern-based copying]

media:
    players:
        packages: [list]
        configs: [pattern-based copying]
```

#### IMPLEMENTATION PHASES:

**Phase 1: Core Infrastructure** ✅ COMPLETE

- YAML parsing system
- Package installation engine
- Config copying with patterns
- Dependency resolution
- Error handling and logging
- Beautiful TUI progress display

**Phase 2: Handler Functions** 🔄 IN PROGRESS

- Git credentials input (NEXT)
- GPU driver detection and installation
- Service configuration and startup
- yay AUR helper installation
- Hardware-specific configurations

**Phase 3: Complete YAML Migration** 📋 PLANNED

- Convert all remaining .sh files to YAML modules
- Add missing categories (system, applications, optional)
- Implement all discovered handler functions
- Remove shell script dependencies entirely

**Phase 4: Testing and Release** 📋 PLANNED

- Full end-to-end testing on clean systems
- Version 2.5.0 release preparation
- Documentation updates
- Merge to main branch

#### BENEFITS:

- **90% code reduction**: Replace 30+ shell scripts with single YAML file + handlers
- **Parallel execution**: Concurrent package installation with goroutines
- **Better error handling**: Proper error types and recovery
- **No external dependencies**: Single compiled binary, no bash/gum/yq needed
- **Consistent patterns**: All modules follow same YAML structure
- **Easier maintenance**: Add new software by editing YAML, not writing shell scripts
- **Beautiful UI**: Professional TUI instead of basic terminal output

**Current Status: v2.5.0 (IN DEVELOPMENT)**

- TUI architecture: ✅ COMPLETE
- Core YAML processing: ✅ COMPLETE
- Git credentials input: 🔄 NEXT PRIORITY
- Shell script analysis:
  📋 ONGOING
- Handler functions: 📋 IN PROGRESS
- Full YAML migration: 📋 PLANNED

### 🚀 PREPARATION: Understanding Current Architecture for YAML Transition

**Status**: Architecture analysis complete, proceeding with systematic YAML migration

**Goal**: Complete replacement of shell script architecture

#### Current Architecture Analysis:

**File Categories for YAML Migration:**

1. **Simple Package Lists** (Easy YAML conversion):
    - `install/applications/` - Direct package → YAML mapping
    - `install/development/` - Tools and packages → YAML mapping
    - `install/optional/` - Optional packages → YAML mapping
    - Most desktop components → YAML mapping

2. **Complex Logic** (Requires handlers):
    - GPU detection and driver installation
    - Service configuration and startup
    - yay AUR helper installation
    - Hardware-specific configurations

3. **Current Duplicated Patterns** (YAML will eliminate):
    - `install_packages` function calls in every script
    - Config file copying logic repeated everywhere
    - Inconsistent error handling patterns
    - Manual dependency management

#### YAML Processing Requirements:

- Package installation with pacman/yay
- Config file copying with pattern matching
- Dependency resolution and ordering
- Handler function execution for complex operations
- Error handling and rollback capability
- Progress reporting integration
