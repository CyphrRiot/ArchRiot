#!/bin/bash

# ==============================================================================
# OhmArchy Development Tools Setup
# ==============================================================================
# Simple development tools - essential compilers and language tools
# ==============================================================================

# Install core build tools
yay -S --noconfirm --needed \
    base-devel \
    git \
    clang \
    cmake \
    ninja \
    rust \
    python-pip \
    go

# Install common development utilities
yay -S --noconfirm --needed \
    github-cli \
    jq \
    curl \
    wget \
    unzip \
    zip

# Install migrate backup tool (simple binary download)
echo "📦 Installing migrate backup tool..."
migrate_url="https://raw.githubusercontent.com/CyphrRiot/Migrate/main/bin/migrate"
mkdir -p ~/.local/bin
curl -L -o ~/.local/bin/migrate "$migrate_url"
chmod +x ~/.local/bin/migrate

echo "✅ Development tools setup complete!"
