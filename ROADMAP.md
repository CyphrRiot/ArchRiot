# ArchRiot Development Roadmap

## AI DEVELOPMENT RULES

### SHELL SCRIPT MIGRATION

- **COMPLETE MIGRATION ONLY**: Migrate ALL functionality from shell scripts to YAML
- **NO PARTIAL IMPLEMENTATIONS**: Every package, command, config goes into YAML
- **READ ENTIRE SCRIPT FIRST**: Then implement everything in one edit
- **CHECK FOR DUPLICATES**: Always verify packages/commands don't already exist in YAML
- **CHECK EXISTING SYSTEM**: Use system commands to understand how packages should be configured and what dependencies exist - NOT to skip installation
- **NO BABYSITTING**: Don't require step-by-step guidance when source script has all info

### IMPLEMENTATION APPROACH

- **ONE COMPLETE EDIT**: Make full implementations, not incremental changes
- **FOLLOW EXISTING PATTERNS**: Use established YAML structure and Go patterns
- **NO SELECTIVE EDITING**: If it's in the script, it goes in YAML

### COMMIT RULES

- **ASK BEFORE COMMITTING**: Never commit without explicit permission
- **UPDATE VERSION**: Update VERSION and README.md version when releasing
- **UPDATE CHANGELOG**: Document all changes in CHANGELOG.md

## CURRENT OBJECTIVE: YAML ARCHITECTURE (v2.5.0)

**STATUS**: YAML architecture implemented and working
**TASK**: Convert remaining shell scripts from `install/pending/` to YAML

### MIGRATION TARGET

Replace 30+ shell scripts with single `packages.yaml` configuration:

- `install/pending/*.sh` → YAML modules
- Maintain all functionality
- Improve reliability and maintainability

## YAML ARCHITECTURE

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

- `communication.sh` - Web applications
- `theming.sh` - Desktop themes
- `productivity.sh` - Office tools
- `utilities.sh` - System utilities
- `specialty.sh` - Specialized tools

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

1. Convert remaining pending scripts
2. Remove obsolete shell scripts
3. Update documentation
4. Final testing and validation
