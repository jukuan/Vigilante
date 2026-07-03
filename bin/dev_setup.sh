#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Setting up vigilate development environment..."

# Download dependencies
echo "Running go mod tidy..."
go mod tidy

# Ensure all scripts are executable
echo "Making scripts executable..."
chmod +x bin/*.sh scripts/*.sh 2>/dev/null || true

echo "Development setup complete."
