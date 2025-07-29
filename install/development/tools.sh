#!/bin/bash

# ==============================================================================
# ArchRiot Development Tools Setup
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
    zip \
    tmux

# Install migrate backup tool (simple binary download)
echo "📦 Installing migrate backup tool..."
migrate_url="https://raw.githubusercontent.com/CyphrRiot/Migrate/main/bin/migrate"
mkdir -p ~/.local/bin

# Stop any running migrate processes before updating
pkill -f migrate 2>/dev/null || true
sleep 1

# Download to temporary location first to avoid "text file busy" error
temp_migrate="/tmp/migrate-$$"
if curl -L -o "$temp_migrate" "$migrate_url"; then
    chmod +x "$temp_migrate"
    mv "$temp_migrate" ~/.local/bin/migrate
    echo "✓ Migrate backup tool installed successfully"
else
    echo "❌ Failed to download migrate tool"
    rm -f "$temp_migrate"
fi

echo "✅ Development tools setup complete!"
