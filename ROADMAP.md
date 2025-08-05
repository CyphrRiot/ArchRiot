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

# ArchRiot Development Roadmap

## Current Status: January 2025

ArchRiot has reached stable v2.5.6 with a robust Go-based installer, comprehensive Waybar integration, and solid system foundations. This roadmap outlines the next major development phases.

---

## 🎯 **Phase 1: Tools Integration (Q1 2025)**

### **Objective: Unified Command Interface**

Replace the current shell-based `archriot-tools` script with an integrated Go module that provides a beautiful, consistent TUI for optional tools.

### **Current Implementation Status**

- ✅ **Basic tools module created** (`source/tools/tools.go`)
- ✅ **Command line argument parsing** added to main.go
- ✅ **TUI framework integration** using Bubble Tea + Lipgloss
- ✅ **Tool registry system** (Go-implemented tools, no shell scripts)
- ✅ **Professional TUI interface** (matches Migrate's clean design pattern)
- ✅ **Four tools implemented** (Secure Boot, Memory Optimizer, Performance Tuner, Dev Environment)
- 🔄 **Testing and refinement** (basic functionality working)
- ❌ **Documentation updates** (not started)

### **Target Command Structure**

```bash
archriot              # Standard installation (unchanged)
archriot --tools      # Launch tools interface
archriot --version    # Show version info
archriot --help       # Show help message
```

### **Proposed Tools Interface Design**

#### **Visual Design Requirements**

- **Consistent styling**: Match installer's purple/cyan theme
- **Professional layout**: Same TUI framework as main installer
- **Intuitive navigation**: Arrow keys, Enter, 'q' to quit
- **Status indicators**: Show tool availability and readiness
- **Category organization**: Group tools by function

#### **Interface Layout**

```
🔧 ArchRiot Optional Tools 🔧

Select a tool to run (Enter) • Toggle details (i/d) • Refresh (r) • Quit (q)

> ✓ 🛡️  Secure Boot Setup          [Security]
    Clean UEFI Secure Boot implementation using sbctl + shim-signed
    ⚠ Advanced Tool - Use with caution
    Category: Security

  ✓ ⚡ Performance Tuner           [Performance]
  ✗ 🧠 Memory Optimizer            [System]
  ✓ 💻 Development Environment     [Development]

Navigate: ↑↓ • Select: Enter/Space • Details: i/d • Refresh: r • Quit: q
```

### **Implementation Tasks**

#### **Phase 1.1: Core Infrastructure** (COMPLETED)

- [x] Complete tool registry system with Go implementations
- [x] Implement proper error handling and validation
- [x] Add tool execution with native Go functions
- [x] Create comprehensive tool metadata system
- [x] Professional TUI matching Migrate's design pattern

#### **Phase 1.2: Tool Categories** (PARTIALLY COMPLETED)

- [x] **Security Tools**
    - [x] Secure Boot Setup (native Go implementation)
    - [ ] Firewall Configuration
    - [ ] System Hardening
- [x] **Performance Tools**
    - [x] Memory Optimizer (native Go implementation)
    - [x] Performance Tuner (native Go implementation)
    - [ ] I/O Scheduler Optimization
- [x] **Development Tools**
    - [x] Development Environment Setup (native Go implementation)
    - [ ] Docker/Container Configuration
    - [ ] IDE Integration Tools
- [ ] **System Tools**
    - [ ] Backup Configuration
    - [ ] Network Optimization
    - [ ] Power Management

#### **Phase 1.3: User Experience** (IN PROGRESS)

- [x] Clean single-column menu design (matches Migrate pattern)
- [x] Tool availability checking and validation
- [x] Professional selection indicators and styling
- [x] Detailed tool information views
- [ ] Progress indication for long-running tools
- [ ] Enhanced tool execution logs and output handling
- [ ] Confirmation dialogs for destructive operations

#### **Phase 1.4: Integration & Testing** (NEXT PRIORITY)

- [x] Build system supports tools module (working)
- [x] Basic functionality tested and working
- [ ] Comprehensive testing on fresh systems
- [ ] Documentation updates (README.md, tool help text)
- [ ] Migration path for existing `archriot-tools` users
- [ ] Tool implementation completion and refinement

### **Technical Requirements**

#### **Tool Discovery System**

- Scan `optional-tools/` directory structure
- Auto-detect executable permissions
- Validate tool prerequisites
- Support for nested tool categories

#### **Tool Metadata Format**

```go
type Tool struct {
    ID          string    // Unique identifier
    Name        string    // Display name with emoji
    Description string    // Brief description
    Category    string    // Security, Performance, Development, System
    Path        string    // Full path to executable
    Executable  bool      // Permission check result
    Advanced    bool      // Requires warning dialog
    Prerequisites []string // Required packages/conditions
}
```

#### **Directory Structure**

```
optional-tools/
├── security/
│   ├── setup-secure-boot.sh
│   ├── configure-firewall.sh
│   └── system-hardening.sh
├── performance/
│   ├── memory-optimizer.sh
│   ├── cpu-tuning.sh
│   └── io-optimization.sh
├── development/
│   ├── dev-environment.sh
│   ├── docker-setup.sh
│   └── ide-integration.sh
└── system/
    ├── backup-config.sh
    ├── network-optimization.sh
    └── power-management.sh
```

### **Success Criteria**

- [ ] All existing tools accessible through new interface
- [ ] Same visual quality and user experience as main installer
- [ ] Zero breaking changes to existing installation process
- [ ] Tool execution success rate >95%
- [ ] Documentation complete and up-to-date

---

## 🚀 **Phase 2: Advanced Tool Development (Q2 2025)**

### **New Tool Categories**

#### **System Administration**

- Advanced backup and restore system
- System monitoring and alerting
- Log management and analysis
- Service management interface

#### **Gaming Optimization**

- GPU driver optimization
- Game-specific performance tuning
- Gaming peripheral configuration
- Steam/Lutris integration tools

#### **Privacy & Security**

- VPN configuration wizard
- Privacy-focused browser setup
- Encrypted storage management
- Network security auditing

### **Tool Framework Enhancements**

- Plugin system for third-party tools
- Tool dependency management
- Automated tool updates
- Configuration backup/restore

---

## 🔧 **Phase 3: System Enhancement (Q3 2025)**

### **Installation Improvements**

- Hardware detection and optimization
- Automated driver installation
- Post-install optimization wizard
- System validation and health checks

### **Configuration Management**

- Centralized settings management
- Profile-based configurations
- Cloud sync for settings
- Migration tools between systems

---

## 📱 **Phase 4: Ecosystem Expansion (Q4 2025)**

### **Platform Support**

- ArchRiot Server Edition
- Container/Docker variants
- ARM64 support expansion
- Virtual machine optimizations

### **Community Features**

- Tool marketplace/repository
- Community configuration sharing
- Documentation wiki
- User contribution system

---

## 🎯 **Immediate Next Steps (Current Week)**

### **Critical Path Items**

1. **Complete tools module implementation**
    - Finish tool scanning system
    - Add proper error handling
    - Test tool execution workflow

2. **Build system integration**
    - Update Makefile for new module
    - Verify compilation on clean system
    - Test packaging process

3. **User testing**
    - Test on development system
    - Verify existing workflows unchanged
    - Validate tool discovery and execution

4. **Documentation**
    - Update README.md with new commands
    - Create tools documentation
    - Update help text and examples

### **Risk Mitigation**

- **Breaking changes**: Maintain backward compatibility
- **Tool failures**: Implement proper error handling and recovery
- **User confusion**: Clear documentation and intuitive interface
- **Performance impact**: Ensure tools interface is responsive

---

## 📊 **Success Metrics**

### **Phase 1 Targets**

- **User Adoption**: >80% of users try tools interface within 30 days
- **Tool Success Rate**: >95% of tool executions complete successfully
- **User Satisfaction**: >90% positive feedback on tools interface
- **Performance**: Tools interface loads in <2 seconds
- **Reliability**: Zero crashes in tools interface during testing

### **Long-term Targets**

- **Tool Library**: 20+ high-quality tools across all categories
- **Community Engagement**: 50+ community-contributed tools
- **Documentation Quality**: 100% of tools have comprehensive docs
- **Platform Coverage**: Support for 95% of common hardware configurations

---

## 🛠️ **Development Guidelines**

### **Code Quality Standards**

- All new code must include comprehensive error handling
- TUI components must be responsive and accessible
- Tool execution must be sandboxed and safe
- All tools must include help text and prerequisites

### **Testing Requirements**

- Unit tests for all tool discovery and execution logic
- Integration tests on fresh Arch Linux installations
- User acceptance testing with real-world scenarios
- Performance testing under various system loads

### **Documentation Standards**

- All tools must include README.md with usage examples
- Code must include inline documentation
- User-facing help text must be clear and actionable
- Installation and setup procedures must be step-by-step

---

## 📝 **Notes**

### **Current Codebase Status**

As of January 2025, the tools integration work has been **substantially completed** for Phase 1. The following components are implemented and working:

- `source/tools/tools.go` - Complete native Go tool registry system
- Main.go argument parsing - Fully implemented and tested (`archriot --tools`)
- Tool execution system - Native Go functions, no shell script dependencies
- TUI components - Professional interface matching Migrate's design pattern
- Four working tools - Secure Boot, Memory Optimizer, Performance Tuner, Dev Environment

### **Lessons Learned**

- Always complete planning phase before implementation
- Maintain strict separation between proposal and implementation
- Follow established development workflows (successfully applied)
- Prioritize testing and validation
- **UI/UX Design**: Study existing successful patterns (Migrate) before implementing
- **Code Quality**: Native Go implementations are superior to shell script wrappers

### **Dependencies**

- Bubble Tea TUI framework (already integrated)
- Lipgloss styling library (already integrated)
- Go 1.19+ (current requirement)
- Existing optional-tools directory structure

## 🚨 **CRITICAL DEVELOPMENT RULES**

### **Build System Requirements**

**NEVER use manual `go build` commands. ALWAYS use the Makefile.**

#### ❌ FORBIDDEN:

```bash
cd source && go build -o archriot-test .
cd source && go build .
go build -o anything
```

#### ✅ REQUIRED:

```bash
make build          # Build the installer
make test          # Test the installer
make clean         # Clean build artifacts
```

#### **Why This Matters:**

- Makefile handles proper versioning and build flags
- Prevents creation of unwanted `source/archriot-test` files
- Ensures consistent build process across development
- Maintains clean repository state
- Follows established project conventions

#### **Enforcement:**

- Any build commands MUST use `make`
- No exceptions for testing or development
- Repository should never contain `source/archriot-test` or similar files
- All documentation and scripts should reference `make build`

This roadmap will be updated as development progresses and requirements evolve.

---

## 📋 **Go Tool Implementation Proposal**

### **Purpose: Replace Shell Scripts with Native Go Tools**

Based on analysis of the existing `optional-tools/` directory structure and shell scripts, this proposal outlines the complete reimplementation of all tools as native Go functions within the tools registry system.

### **Current Shell Script Analysis**

#### **Existing Tool Structure**

```
optional-tools/
├── README.md                 # Comprehensive tool documentation
├── IMPLEMENTATION.md         # Implementation status and details
├── launcher.sh              # Shell-based tool launcher
├── secure-boot/
│   ├── README.md            # Secure Boot documentation
│   └── setup-secure-boot.sh # Secure Boot implementation
├── dell-sleep-fix/
│   ├── README.md            # Dell XPS sleep fix documentation
│   └── setup-dell-sleep-fix.sh # Dell sleep crash fix
└── ivy-bridge-vulkan-fix/
    ├── README.md            # Ivy Bridge Vulkan compatibility
    └── fix-ivy-bridge-vulkan.sh # Vulkan fix for old Intel graphics
```

#### **Shell Script Entry Point**

```bash
# archriot-tools (4.4KB shell script)
# - Interactive menu system for optional tools
# - Safety warnings and user confirmation
# - System requirements checking
# - Documentation viewer
# - Tool execution wrapper
```

### **Tool Implementation Specifications**

#### **1. Secure Boot Setup Tool**

**Current Implementation:** `secure-boot/setup-secure-boot.sh` (13.3KB)

**Go Implementation Requirements:**

```go
type SecureBootTool struct {
    Name        string
    Description string
    Category    string
    Advanced    bool
    Available   bool
}

func executeSecureBoot() error {
    // Implementation phases:
    // 1. System compatibility checks (UEFI mode, internet, Arch Linux)
    // 2. Package installation (sbctl, shim-signed from AUR)
    // 3. UEFI setup guidance (interactive prompts)
    // 4. Key management (creation and enrollment)
    // 5. File signing (bootloader, kernels, initramfs)
    // 6. Verification and setup validation
    // 7. Pacman hook configuration for automatic signing
}
```

**Features to Implement:**

- ✅ UEFI mode detection and validation
- ✅ Internet connectivity checks
- ✅ Package installation (sbctl, shim-signed)
- ✅ Interactive UEFI setup guidance
- ✅ Custom Secure Boot key creation
- ✅ Bootloader and kernel signing
- ✅ Comprehensive verification system
- ✅ Automatic kernel signing on updates
- ✅ Multi-mode setup (quick, manual, verification)
- ✅ Recovery guidance and troubleshooting

**Risk Level:** ⚠️ MODERATE - Modifies boot components
**Compatibility:** Universal UEFI systems (Intel, AMD)

#### **2. Dell XPS Sleep Fix Tool**

**Current Implementation:** `dell-sleep-fix/setup-dell-sleep-fix.sh`

**Go Implementation Requirements:**

```go
type DellSleepFixTool struct {
    Name        string
    Description string
    Category    string
    Advanced    bool
    Available   bool
}

func executeDellSleepFix() error {
    // Implementation phases:
    // 1. Hardware compatibility detection (Dell XPS + Intel Arc Graphics)
    // 2. XE graphics driver parameter configuration
    // 3. systemd sleep configuration setup
    // 4. Pre/post sleep hook installation
    // 5. Kernel parameter addition for stability
    // 6. Diagnostic tool creation
    // 7. Verification and testing guidance
}
```

**Features to Implement:**

- ✅ Dell XPS hardware detection
- ✅ Intel Lunar Lake Arc Graphics validation (130V/140V)
- ✅ XE driver stability configuration
- ✅ S2idle sleep mode enforcement
- ✅ systemd sleep hook management
- ✅ Kernel parameter optimization
- ✅ Sleep diagnostic utilities
- ✅ Automatic log rotation setup
- ✅ Recovery and uninstall procedures

**Files to Create/Modify:**

- `/etc/modprobe.d/xe-graphics-fix.conf`
- `/etc/systemd/sleep.conf.d/dell-sleep-fix.conf`
- `/lib/systemd/system-sleep/dell-sleep-hook.sh`
- `/boot/loader/entries/*.conf` (kernel parameters)
- `/etc/logrotate.d/dell-sleep`

**Risk Level:** ⚠️ MODERATE - Modifies kernel parameters and power management
**Compatibility:** Dell XPS laptops with Intel Lunar Lake Arc Graphics

#### **3. Ivy Bridge Vulkan Fix Tool**

**Current Implementation:** `ivy-bridge-vulkan-fix/fix-ivy-bridge-vulkan.sh`

**Go Implementation Requirements:**

```go
type IvyBridgeVulkanTool struct {
    Name        string
    Description string
    Category    string
    Advanced    bool
    Available   bool
}

func executeIvyBridgeVulkanFix() error {
    // Implementation phases:
    // 1. Intel graphics generation detection
    // 2. User-specific configuration creation
    // 3. Application launcher setup
    // 4. Vulkan compatibility testing
    // 5. OpenGL fallback configuration
    // 6. Environment variable optimization
}
```

**Features to Implement:**

- ✅ Intel graphics hardware detection (HD 3000/4000)
- ✅ Safe user-specific configuration (no system modification)
- ✅ Application compatibility launchers
- ✅ Vulkan capability testing
- ✅ OpenGL fallback configuration
- ✅ Performance optimization settings
- ✅ Easy removal and testing procedures

**Files to Create:**

- `~/.local/bin/zed-ivy-bridge` (compatibility launcher)
- `~/.config/zed/ivy-bridge-overlay.json` (configuration overlay)
- `~/.cache/archriot/ivy-bridge-fix.log` (detailed logs)

**Risk Level:** ✅ LOW - User-specific configurations only
**Compatibility:** Intel Ivy Bridge/Sandy Bridge systems (2011-2013)

#### **4. Performance Optimization Tools**

**New Go Implementation Requirements:**

```go
type PerformanceTools struct {
    MemoryOptimizer    Tool
    CPUGovernorTuner   Tool
    IOSchedulerTuner   Tool
    SystemCacheTuner   Tool
}

func executeMemoryOptimizer() error {
    // - zswap configuration
    // - swappiness tuning
    // - memory pressure handling
    // - OOM killer optimization
}

func executeCPUGovernorTuner() error {
    // - CPU frequency governor selection
    // - Per-core configuration
    // - Power profile optimization
    // - Thermal management
}

func executeIOSchedulerTuner() error {
    // - Storage device optimization
    // - SSD/HDD specific tuning
    // - Scheduler selection (mq-deadline, kyber, bfq)
    // - Queue depth optimization
}
```

#### **5. Development Environment Tools**

**New Go Implementation Requirements:**

```go
type DevelopmentTools struct {
    DevEnvironmentSetup Tool
    DockerConfiguration Tool
    IDEIntegration     Tool
    GitOptimization    Tool
}

func executeDevEnvironment() error {
    // - Essential development packages
    // - Compiler toolchains
    // - Version managers (nvm, rustup, etc.)
    // - Development utilities
}

func executeDockerConfiguration() error {
    // - Docker installation and setup
    // - Rootless Docker configuration
    // - Container optimization
    // - Development workflow setup
}
```

### **Implementation Strategy**

#### **Phase 1: Core Infrastructure (COMPLETED)**

- [x] Go tool registry system
- [x] Professional TUI interface
- [x] Tool execution framework
- [x] Error handling and logging

#### **Phase 2: Critical Tool Implementation (NEXT)**

- [ ] Secure Boot Setup (highest priority)
- [ ] Dell XPS Sleep Fix (hardware-specific)
- [ ] Ivy Bridge Vulkan Fix (compatibility)

#### **Phase 3: Performance and Development Tools**

- [ ] Memory Optimizer implementation
- [ ] CPU Governor Tuner implementation
- [ ] Development Environment Setup
- [ ] Docker Configuration Tools

#### **Phase 4: Advanced Features**

- [ ] Tool dependency management
- [ ] Configuration backup/restore
- [ ] Tool update system
- [ ] User preference persistence

### **Technical Implementation Details**

#### **Tool Registration System**

```go
func GetAvailableTools() ([]Tool, error) {
    tools := []Tool{
        {
            ID:          "secure-boot",
            Name:        "🛡️  Secure Boot Setup",
            Description: "Clean UEFI Secure Boot implementation using sbctl + shim-signed",
            Category:    "Security",
            ExecuteFunc: executeSecureBoot,
            Advanced:    true,
            Available:   checkSecureBootAvailable(),
        },
        {
            ID:          "dell-sleep-fix",
            Name:        "🛠️  Dell XPS Sleep Fix",
            Description: "Fix sleep/suspend crashes on Dell XPS laptops with Intel Arc Graphics",
            Category:    "Hardware",
            ExecuteFunc: executeDellSleepFix,
            Advanced:    true,
            Available:   checkDellHardware(),
        },
        {
            ID:          "ivy-bridge-vulkan",
            Name:        "📺 Ivy Bridge Vulkan Fix",
            Description: "Vulkan compatibility for Intel HD Graphics 3000/4000",
            Category:    "Compatibility",
            ExecuteFunc: executeIvyBridgeVulkanFix,
            Advanced:    false,
            Available:   checkIvyBridgeHardware(),
        },
        // Additional tools...
    }
    return tools, nil
}
```

#### **Hardware Detection System**

```go
func checkSecureBootAvailable() bool {
    // Check UEFI mode and Secure Boot capability
    if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
        return false
    }
    return true
}

func checkDellHardware() bool {
    // Check for Dell system with Intel Arc Graphics
    cmd := exec.Command("lspci")
    output, err := cmd.Output()
    if err != nil {
        return false
    }

    isDell := strings.Contains(string(output), "Dell")
    hasIntelArc := strings.Contains(string(output), "Intel Arc Graphics")
    return isDell && hasIntelArc
}

func checkIvyBridgeHardware() bool {
    // Check for Intel HD Graphics 3000/4000
    cmd := exec.Command("lspci")
    output, err := cmd.Output()
    if err != nil {
        return false
    }

    return strings.Contains(string(output), "Intel.*HD Graphics") &&
           (strings.Contains(string(output), "3000") || strings.Contains(string(output), "4000"))
}
```

### **Migration Benefits**

#### **Advantages of Go Implementation**

1. **Unified Codebase**: Single language for entire project
2. **Better Error Handling**: Comprehensive error management
3. **Type Safety**: Compile-time error detection
4. **Cross-Platform**: Easier portability if needed
5. **Performance**: Faster execution than shell scripts
6. **Testing**: Unit testable components
7. **Maintenance**: Easier to maintain and extend
8. **Integration**: Better integration with main installer

#### **User Experience Improvements**

1. **Consistent Interface**: Same TUI framework as main installer
2. **Real-time Feedback**: Better progress indication
3. **Error Recovery**: More sophisticated error handling
4. **Validation**: Better pre-flight checks
5. **Documentation**: Integrated help system
6. **Safety**: More comprehensive safety checks

### **File Removal Strategy**

After Go implementation is complete and tested:

#### **Files to Remove**

```bash
# Shell script infrastructure
rm -rf optional-tools/
rm -f archriot-tools

# Legacy documentation (will be integrated into main docs)
# Note: Keep key documentation content for reference
```

#### **Documentation Migration**

- Move essential documentation to main README.md
- Integrate tool help into Go help system
- Preserve safety warnings and hardware compatibility info
- Maintain troubleshooting guides in new format

### **Testing and Validation**

#### **Testing Requirements**

- [ ] Unit tests for each tool implementation
- [ ] Hardware compatibility testing
- [ ] Error condition testing
- [ ] Recovery procedure validation
- [ ] User acceptance testing

#### **Safety Validation**

- [ ] Backup/restore testing
- [ ] Rollback procedure verification
- [ ] Error recovery testing
- [ ] Hardware damage prevention checks

### **Success Criteria**

1. **Functional Parity**: All shell script functionality preserved
2. **Improved Safety**: Better error handling and recovery
3. **Better UX**: Consistent interface and feedback
4. **Maintainability**: Easier to extend and modify
5. **Documentation**: Comprehensive integrated help
6. **Testing**: Full test coverage for critical operations

### **Timeline Estimate**

- **Week 1-2**: Secure Boot implementation and testing
- **Week 3**: Dell Sleep Fix implementation
- **Week 4**: Ivy Bridge Vulkan Fix implementation
- **Week 5**: Performance tools implementation
- **Week 6**: Development tools implementation
- **Week 7**: Integration testing and documentation
- **Week 8**: User acceptance testing and refinement

This proposal provides a complete roadmap for eliminating the shell script dependencies while preserving and improving all existing functionality through native Go implementations.
