#!/bin/bash
set -e

cd "$(dirname "$0")/.."

GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
BINARY="vigilante"

# Add .exe for Windows
if [ "$GOOS" = "windows" ]; then
    BINARY="vigilante.exe"
fi

echo "Building $BINARY for $GOOS/$GOARCH..."
GOOS=$GOOS GOARCH=$GOARCH go build -o "$BINARY" .

echo "Build complete: $(pwd)/$BINARY"
