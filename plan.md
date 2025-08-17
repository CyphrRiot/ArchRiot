# ArchRiot Development Plan

## OVERVIEW

This document tracks outstanding development tasks for the ArchRiot installer and system components.

## 🚧 OUTSTANDING TASKS

### TASK 1: Terminal Emoji Fallback System (TOP PRIORITY)

### Problem

When installing from a basic terminal (before Wayland/Hyprland), emojis don't display properly, making the installer format look broken and unprofessional.

### Requirements

- Automatic detection of terminal emoji support
- ASCII fallbacks for all emojis in installer
- Maintain exact visual width/spacing
- No user configuration required
- Test mode for verification

### Root Cause Analysis

The installer uses emojis in two places:

1. **Logger system** - Already has centralized emoji handling via `getStatusEmoji()` and `typeEmojis` map
2. **TUI hardcoded strings** - Direct emoji strings like "🎯 Current Step:" that bypass logger

### Implementation Strategy (SIMPLE)

1. **Enhance existing logger emoji system**
    - Add terminal detection to `InitLogging()`
    - Add ASCII fallback maps (single characters only)
    - Modify existing `getStatusEmoji()` to check terminal support
    - Add `getTypeEmoji()` function for type emojis

2. **Fix TUI hardcoded emojis**
    - Replace hardcoded "🎯 Current Step:" with logger-based approach
    - Replace hardcoded "📁 Log File:" with logger-based approach
    - All other emojis already use logger system

3. **Add testing support**
    - Add `--ascii-only` flag to force ASCII mode
    - Add terminal detection logic (check TERM, LANG variables)

### Technical Details

```go
// Terminal detection
func detectTerminalEmojiSupport() bool {
    term := os.Getenv("TERM")
    lang := os.Getenv("LANG")

    // Basic terminals don't support emojis
    if term == "linux" || term == "console" || term == "dumb" {
        return false
    }

    // No UTF-8 = no emojis
    if !strings.Contains(lang, "UTF-8") {
        return false
    }

    return true
}

// Single-character fallbacks (maintain spacing)
var statusFallbacks = map[string]string{
    "Success":  "+",
    "Error":    "X",
    "Warning":  "?",
    "Progress": ".",
    "Complete": "!",
    "Info":     "i",
}

var typeFallbacks = map[string]string{
    "Package":  "*",
    "Git":      "+",
    "File":     "-",
    "System":   "~",
    "Module":   "=",
    "Database": "#",
}
```

### Files to Modify

1. `source/logger/logger.go` - Add detection and fallbacks
2. `source/tui/model.go` - Replace hardcoded emojis (2 lines only)
3. `source/main.go` - Add `--ascii-only` flag

### Testing Plan

1. Test with `--ascii-only` flag in normal terminal
2. Test in actual TTY (Ctrl+Alt+F2)
3. Test with `TERM=linux` environment variable
4. Verify spacing/alignment maintained

### TASK 2: Secure Boot Enablement

### Problem

Users want Secure Boot configuration for enhanced security with LUKS encryption.

### Requirements

- Detect current Secure Boot status
- Educational prompts about benefits
- Automated key generation and installation

### Implementation Challenges

- Requires UEFI Setup Mode for custom keys
- User must manually enter BIOS settings
- Coordination between installer phases

### Technical Approach

- Detection functions already implemented
- TUI prompts for user education
- Multi-stage setup with reboot coordination

## NEXT STEPS

**Remaining Tasks:**

1. Complete emoji fallback system
2. Finish Secure Boot implementation
3. Additional features as requested

**Development Priorities:**

1. User experience improvements (emoji fallbacks)
2. Security enhancements (Secure Boot)
3. Feature requests from community

**Code Architecture:**

- Modular design with clear separation of concerns
- TUI-based user interaction with semantic logging
- File system checks: `os.Stat()` for directories/files (e.g., `/sys/firmware/efi`)
- State management: Global variables with TUI callbacks for user decisions

```go
// Example detection patterns
func DetectSecureBootStatus() (bool, bool, error) {
    // Implementation for Secure Boot detection
}

func DetectLuksEncryption() (bool, []string, error) {
    // Implementation for LUKS detection
}

type SecureBootStatusMsg struct {
    Enabled     bool
    Supported   bool
    LuksUsed    bool
    LuksDevices []string
}

type SecureBootPromptMsg struct{}

var secureBootCallback func(bool)

func SetSecureBootCallback(callback func(bool)) {
    secureBootCallback = callback
}

type InstallModel struct {
    // TUI model fields
    secureBootEnabled   bool
    secureBootSupported bool
    luksDetected        bool
    luksDevices         []string
}

func checkSecureBootRecommendation() (bool, string) {
    // Logic for determining if Secure Boot should be recommended
}

func getSecureBootEducation(luksDevices []string) string {
    // Educational content about Secure Boot benefits
}
```
