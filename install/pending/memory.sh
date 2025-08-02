#!/bin/bash

# ArchRiot Memory Optimization Configuration
# Intelligent memory management that prevents aggressive caching and preserves 1GB free RAM

set -e

echo "🧠 Configuring Memory Optimization..."

# Create sysctl configuration for memory management
echo "📝 Creating memory optimization configuration..."

# Create the sysctl configuration file
sudo tee /etc/sysctl.d/99-memory-optimization.conf >/dev/null <<EOF
# ArchRiot Memory Management Optimization
# Fixes Linux's aggressive caching behavior that consumes all available RAM

# Reserve 1GB of free memory (1048576 KB)
# Prevents system from caching everything and struggling to free memory when needed
vm.min_free_kbytes=1048576

# Reduce VFS cache pressure - be less aggressive about caching file system metadata
# Default is 100, setting to 50 makes kernel less likely to reclaim directory/inode caches
vm.vfs_cache_pressure=50

# Reduce swappiness - prefer to use RAM before swapping to disk
# Default is 60, setting to 10 makes system much less likely to swap
vm.swappiness=10

# Limit dirty pages to 5% of memory before forcing writes
# Prevents large write bursts that can cause system lag
vm.dirty_ratio=5

# Start background writeback at 2% dirty pages
# Earlier writeback prevents dirty page buildup
vm.dirty_background_ratio=2

# Disable zone reclaim mode for better NUMA performance
# Prevents unnecessary memory reclaim on NUMA systems
vm.zone_reclaim_mode=0
EOF

# Apply the configuration immediately
echo "⚡ Applying memory optimization settings..."
sudo sysctl -p /etc/sysctl.d/99-memory-optimization.conf

# Verify the settings were applied
echo "🔍 Verifying memory optimization settings..."

check_setting() {
    local setting="$1"
    local expected="$2"
    local current=$(sysctl -n "$setting" 2>/dev/null || echo "unknown")

    if [[ "$current" == "$expected" ]]; then
        echo "  ✓ $setting = $current"
    else
        echo "  ⚠ $setting = $current (expected: $expected)"
    fi
}

check_setting "vm.min_free_kbytes" "1048576"
check_setting "vm.vfs_cache_pressure" "50"
check_setting "vm.swappiness" "10"
check_setting "vm.dirty_ratio" "5"
check_setting "vm.dirty_background_ratio" "2"
check_setting "vm.zone_reclaim_mode" "0"

# Show current memory status
echo ""
echo "💾 Current memory status:"
free -h

# Calculate and show reserved memory
total_kb=$(awk '/MemTotal/ {print $2}' /proc/meminfo)
free_kb=$(awk '/MemFree/ {print $2}' /proc/meminfo)
reserved_mb=$((1048576 / 1024))

echo ""
echo "🛡️  Memory protection active:"
echo "   • ${reserved_mb}MB (1GB) reserved for system responsiveness"
echo "   • Reduced aggressive file caching"
echo "   • Optimized for interactive desktop performance"

echo ""
echo "✅ Memory optimization configured successfully!"
echo "   These settings will persist across reboots and prevent"
echo "   the system from consuming all RAM for file caches."
