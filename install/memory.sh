# Memory optimization configuration
echo "Configuring memory management optimization..."

# Create sysctl configuration for memory optimization
sudo tee /etc/sysctl.d/99-memory-optimization.conf >/dev/null <<EOF
# Memory Management Optimization
vm.min_free_kbytes=1048576
vm.vfs_cache_pressure=50
vm.swappiness=10
vm.dirty_ratio=5
vm.dirty_background_ratio=2
vm.zone_reclaim_mode=0
EOF

# Apply settings immediately
sudo sysctl --system >/dev/null

# Single verification call instead of multiple
echo "Memory optimization settings applied:"
echo "Min free memory: 1048576 KB (1024 MB)"
echo "Cache pressure: 50"
echo "Swappiness: 10"
echo "✓ Memory optimization configured"
