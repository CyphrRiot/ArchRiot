# ArchRiot Development Roadmap

## 🎯 NEXT: YAML ARCHITECTURE TRANSITION (v2.5.0)

## **IMPORTANT** -- RULES

READ THE RULES. READ THE ROADMAP.

1. Propose ONE SMALL DIRECT CHANGE AT A TIME.

2. IF I AGREE to the change, then make it.

3. After you've made it, wait for me to confirm and review the code and ask "Continue?" and wait for "yes" or further instructions.

4. Test each change!

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

1. **🔄 IN PROGRESS: Shell Script Analysis**
    - 33 .sh files in install/pending/ awaiting conversion
    - Need to categorize by complexity (simple package lists vs complex logic)
    - Identify patterns for YAML structure optimization

2. **📋 TODO: YAML Module Expansion**
    - Add missing categories: system, applications, optional
    - Implement conditional installation based on hardware/user choice
    - Add pre/post installation hooks where needed

3. **📋 TODO: Handler Functions**
    - GPU detection and driver selection
    - yay AUR helper installation
    - Service configuration and enablement
    - Hardware-specific configurations

4. **📋 TODO: Testing & Release**
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
- **✅ **MAJOR SIZE REDUCTION\*\*: Reduced main.go from 810+ lines to 98 lines (88% reduction!)
- ✅ **MODULAR ARCHITECTURE WORKING**: Fixed channel communication by proper initialization order
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
