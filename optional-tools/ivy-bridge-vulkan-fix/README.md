# Ivy Bridge Vulkan Fix - SAFE VERSION

## 🛡️ Safe, Non-Destructive Solution for ThinkPad X230 and Ivy Bridge Systems

This is a **SAFE, TESTABLE** fix for Vulkan compatibility issues on Intel HD Graphics 4000 (Ivy Bridge) systems. It does **NOT** modify any core ArchRiot files or system configurations.

### What This Fix Does (SAFELY)

✅ **Creates user-specific configurations only**
✅ **No system file modifications**
✅ **Easy to test and remove**
✅ **No impact on other installations**

### Affected Hardware

- **ThinkPad X230** (Intel HD Graphics 4000)
- **ThinkPad T430** (Intel HD Graphics 4000)
- **Other Ivy Bridge systems** (2012-2013)
- **Sandy Bridge systems** (2011) for reference

### Symptoms This Fixes

- ❌ Zed editor fails to start or crashes
- ❌ "Ivy Bridge Vulkan support is incomplete" warnings
- ❌ Poor performance in graphics applications
- ❌ Applications falling back to software rendering

## Quick Installation

### One-Command Fix

```bash
# Download and run the safe fix
curl -fsSL https://raw.githubusercontent.com/CyphrRiot/ArchRiot/master/optional-tools/ivy-bridge-vulkan-fix/fix-ivy-bridge-vulkan.sh | bash
```

### Manual Installation

```bash
# Clone or download ArchRiot repository
git clone https://github.com/CyphrRiot/ArchRiot.git
cd ArchRiot/optional-tools/ivy-bridge-vulkan-fix

# Run the safe fix
./fix-ivy-bridge-vulkan.sh
```

## What Gets Created (All Safe & Removable)

### 1. 🚀 Safe Launcher

- **File**: `~/.local/bin/zed-ivy-bridge`
- **Purpose**: Launches Zed with Ivy Bridge compatibility
- **Usage**: `zed-ivy-bridge` instead of `zed`

### 2. 📝 Configuration Overlay

- **File**: `~/.config/zed/ivy-bridge-overlay.json`
- **Purpose**: Contains optimal settings for older Intel graphics
- **Usage**: Manually copy settings to your main config

### 3. 📊 Detailed Logs

- **File**: `~/.cache/archriot/ivy-bridge-fix.log`
- **Purpose**: Complete log of all operations and tests
- **Usage**: Debug any issues that arise

## Testing the Fix

### Step 1: Run the Fix

```bash
./fix-ivy-bridge-vulkan.sh
```

### Step 2: Test Zed

```bash
# Try the new launcher
zed-ivy-bridge

# Or test with specific settings
LIBGL_ALWAYS_SOFTWARE=0 zed-ivy-bridge
```

### Step 3: Verify Results

```bash
# Check if Vulkan is working
vulkaninfo --summary

# Check created files
ls -la ~/.local/bin/zed-ivy-bridge
ls -la ~/.config/zed/ivy-bridge-overlay.json
```

## Manual Configuration Steps

The fix creates configuration overlays but **requires manual activation**:

### For Ivy Bridge/Sandy Bridge Systems:

1. **Update Zed Configuration**:

    ```bash
    # View the overlay settings
    cat ~/.config/zed/ivy-bridge-overlay.json

    # Manually copy these settings to ~/.config/zed/settings.json:
    {
      "gpu": {
        "use_vulkan": false,
        "use_opengl": true
      },
      "experimental": {
        "renderer": "opengl"
      }
    }
    ```

2. **Use the Safe Launcher**:
    ```bash
    # Create an alias for convenience
    echo 'alias zed="zed-ivy-bridge"' >> ~/.bashrc
    source ~/.bashrc
    ```

## Testing Results by Hardware

| GPU Generation             | Expected Result           | Test Command              |
| -------------------------- | ------------------------- | ------------------------- |
| **Ivy Bridge** (HD 4000)   | ✅ Works with OpenGL      | `zed-ivy-bridge`          |
| **Sandy Bridge** (HD 3000) | ✅ Works with OpenGL      | `zed-ivy-bridge`          |
| **Haswell** (HD 4600)      | ✅ Works with minor fixes | `zed-ivy-bridge`          |
| **Broadwell+** (HD 5500+)  | ✅ Works normally         | `zed` or `zed-ivy-bridge` |

## Troubleshooting

### If Zed Still Doesn't Work

1. **Check the logs**:

    ```bash
    tail -f ~/.cache/archriot/ivy-bridge-fix.log
    ```

2. **Try software rendering**:

    ```bash
    LIBGL_ALWAYS_SOFTWARE=1 zed-ivy-bridge
    ```

3. **Verify your GPU**:
    ```bash
    lspci | grep -i "vga\|3d\|display"
    glxinfo | grep "OpenGL renderer"
    ```

### If Performance Is Poor

1. **Check hardware acceleration**:

    ```bash
    glxinfo | grep "direct rendering"
    ```

2. **Try Mesa optimizations**:
    ```bash
    export MESA_GL_VERSION_OVERRIDE=3.3
    zed-ivy-bridge
    ```

## Complete Removal

If you want to remove this fix entirely:

```bash
# Remove all created files
rm -f ~/.local/bin/zed-ivy-bridge
rm -f ~/.config/zed/ivy-bridge-overlay.json
rm -f ~/.cache/archriot/ivy-bridge-fix.log

# Remove any aliases you added
sed -i '/zed-ivy-bridge/d' ~/.bashrc
```

## Safety Guarantees

### ✅ What This Fix Does

- Creates user-specific configuration files
- Adds optional compatibility launcher
- Generates detailed logs for debugging
- Provides manual merge instructions

### ❌ What This Fix Does NOT Do

- Modify core ArchRiot installation files
- Change system-wide driver configurations
- Break existing installations
- Require root permissions for operation
- Make irreversible changes

## Advanced Testing

### Vulkan Capability Test

```bash
# Test basic Vulkan functionality
vulkaninfo 2>&1 | head -20

# Test cube rendering (safe 5-second test)
timeout 5 vkcube --c 30
```

### OpenGL Fallback Test

```bash
# Test OpenGL capabilities
glxinfo | grep -E "OpenGL version|OpenGL renderer"

# Test with forced OpenGL
export LIBGL_ALWAYS_SOFTWARE=0
export MESA_GL_VERSION_OVERRIDE=3.3
zed-ivy-bridge
```

## User Reports & Success Stories

### Expected Success Pattern:

1. **Before**: Zed crashes on startup, Vulkan warnings
2. **After**: Zed launches normally with `zed-ivy-bridge`
3. **Performance**: Adequate for code editing, some features may be limited

### Please Report:

- Hardware model and generation
- Before/after behavior
- Performance observations
- Any remaining issues

## Technical Details

### Hardware Detection Logic

```bash
# Safe detection method used
lspci | grep -i "vga\|3d\|display" | grep -i intel

# Conservative classification:
# - "ivybridge": HD 4000, HD 2500
# - "sandybridge": HD 3000, HD 2000
# - "haswell": HD 4600, HD 4400
# - "modern": Broadwell and newer
```

### Environment Variables Set

```bash
# For Ivy Bridge compatibility
export LIBGL_ALWAYS_SOFTWARE=0
export MESA_GL_VERSION_OVERRIDE=3.3
export MESA_GLSL_VERSION_OVERRIDE=330
export DISABLE_VULKAN=1  # Only when needed
```

---

**Note**: This fix is designed for ArchRiot v2.0.1+ and focuses on providing a safe, testable solution that doesn't break existing installations. All changes are user-specific and easily reversible.
