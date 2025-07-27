#!/bin/bash

# ==============================================================================
# ArchRiot Filesystem Support Setup
# ==============================================================================
# Simple filesystem support - GUI mounting and network shares
# ==============================================================================

# Install GUI mounting support and network filesystems
yay -S --noconfirm --needed \
    gvfs \
    udisks2 \
    gvfs-smb \
    gvfs-mtp \
    ntfs-3g

# Enable automatic mounting service
sudo systemctl enable udisks2.service

echo "✅ Filesystem support complete!"
