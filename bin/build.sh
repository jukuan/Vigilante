#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Building vigilate..."
go build -o vigilate .

echo "Build complete: $(pwd)/vigilate"
