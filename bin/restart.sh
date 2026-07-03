#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Restarting vigilate..."
./bin/stop.sh
sleep 2
./bin/start.sh "$@"
