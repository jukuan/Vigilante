#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Restarting vigilante..."
./bin/stop.sh
sleep 2
./bin/start.sh "$@"
