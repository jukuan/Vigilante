#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Running CI checks locally..."

if ! command -v act &> /dev/null; then
    echo "Error: 'act' is not installed."
    echo "Install it from: https://github.com/nektos/act"
    exit 1
fi

act push -P ubuntu-24.04=catthehacker/ubuntu:act-22.04 --container-architecture linux/amd64
