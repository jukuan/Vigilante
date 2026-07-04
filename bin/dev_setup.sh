#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Setting up vigilante development environment..."

# Ensure Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Download dependencies
echo "Running go mod tidy..."
go mod tidy

# Install linting tools
echo "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Ensure all scripts are executable
echo "Making scripts executable..."
chmod +x bin/*.sh scripts/*.sh 2>/dev/null || true

# Build to verify everything compiles
echo "Building to verify..."
go build -o vigilante .

echo "Development setup complete."
echo ""
echo "Quick start:"
echo "  ./bin/test.sh   - run tests"
echo "  ./bin/lint.sh   - run linter"
echo "  ./bin/build.sh  - build binary"
