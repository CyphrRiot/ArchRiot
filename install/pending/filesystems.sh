#!/bin/bash

# ==============================================================================
# ArchRiot Filesystem Support Setup
# ==============================================================================
# Simple filesystem support - GUI mounting and network shares
# ==============================================================================

# Install GUI mounting support and network filesystems
install_packages "gvfs udisks2" "essential"
install_packages "gvfs-smb gvfs-mtp" "essential"
install_packages "ntfs-3g" "essential"

# Enable automatic mounting service
sudo systemctl enable udisks2.service

echo "✅ Filesystem support complete!"
