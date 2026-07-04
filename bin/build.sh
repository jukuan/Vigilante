#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Building vigilante..."
go build -o vigilante .

echo "Build complete: $(pwd)/vigilante"
