# OhmArchy GPU Support

OhmArchy provides comprehensive GPU support for all major graphics vendors with automatic detection and optimized driver installation.

## Supported Graphics Cards

### NVIDIA GPUs

- **Script**: `nvidia.sh`
- **Detection**: Automatic via `lspci`
- **Drivers**: NVIDIA proprietary drivers (open/legacy as appropriate)
- **Features**:
    - Wayland support with proper environment variables
    - Hardware video acceleration (NVENC/NVDEC)
    - Early KMS (Kernel Mode Setting)
    - Automatic driver selection based on GPU generation

### AMD/Radeon GPUs

- **Script**: `amd-vulkan.sh`
- **Detection**: Automatic via `lspci`
- **Drivers**: Open-source Mesa/RADV drivers
- **Features**:
    - Vulkan API support with `vulkan-radeon`
    - Hardware video acceleration (VA-API/VDPAU)
    - 32-bit gaming support
    - Optimized for Wayland/Hyprland

### Intel GPUs

- **Script**: `intel-vulkan.sh`
- **Detection**: Automatic via `lspci`
- **Drivers**: Open-source Mesa/ANV drivers
- **Features**:
    - Vulkan API support with `vulkan-intel`
    - Hardware video acceleration (VA-API)
    - Intel Arc discrete GPU support
    - GuC firmware optimization
    - GPU monitoring tools (`intel_gpu_top`)

## Installation Process

During OhmArchy installation, each GPU script:

1. **Detects** the appropriate GPU hardware
2. **Enables** multilib repository if needed (for 32-bit support)
3. **Installs** vendor-specific drivers and tools
4. **Configures** Hyprland environment variables
5. **Verifies** installation success
6. **Optimizes** system settings for performance

## Manual Installation

Individual GPU support can be installed manually:

```bash
# For AMD/Radeon
~/.local/share/omarchy/install/amd-vulkan.sh

# For Intel
~/.local/share/omarchy/install/intel-vulkan.sh

# For NVIDIA
~/.local/share/omarchy/install/nvidia.sh
```

## Verification

After installation, verify GPU support:

```bash
# Check Vulkan support
vulkaninfo --summary

# Check hardware acceleration (VA-API)
vainfo

# Monitor GPU usage
# Intel: intel_gpu_top or gpu-monitor
# AMD: radeontop (if installed)
# NVIDIA: nvidia-smi
```

## Troubleshooting

- **Reboot** after GPU driver installation for full functionality
- **Check logs** with journalctl if issues persist
- **Verify kernel modules** are loaded correctly
- **Check environment variables** in Hyprland config

The GPU detection is safe - scripts only install packages when the appropriate hardware is detected.
