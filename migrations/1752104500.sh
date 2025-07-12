#!/bin/bash

# OhmArchy Migration - Update repository references
# This migration ensures users get updates from the correct OhmArchy repository

echo "Updating OhmArchy repository references..."

# Update git remote to point to OhmArchy fork
cd ~/.local/share/omarchy
git remote set-url origin https://github.com/CyphrRiot/OhmArchy.git

# Ensure we're on the main branch
git checkout main 2>/dev/null || git checkout -b main

echo "Repository updated to track OhmArchy fork"
