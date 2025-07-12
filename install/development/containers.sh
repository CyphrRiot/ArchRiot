#!/bin/bash

# ==============================================================================
# OhmArchy Container Tools Setup
# ==============================================================================
# Simple container tools - Podman and Distrobox
# ==============================================================================

# Install Podman (Docker alternative)
yay -S --noconfirm --needed \
    podman \
    buildah \
    distrobox

# Enable Podman socket for compatibility
systemctl --user enable --now podman.socket

echo "✅ Container tools setup complete!"
