#!/bin/bash
set -e

cd "$(dirname "$0")/.."

PID_FILE="vigilante.pid"
CONFIG_FILE="${1:-config.yaml}"

if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
    echo "vigilante is already running (PID: $(cat $PID_FILE))"
    exit 1
fi

echo "Starting vigilante with config: $CONFIG_FILE"
nohup ./vigilante "$CONFIG_FILE" > vigilante.log 2>&1 &
echo $! > "$PID_FILE"
echo "Started with PID: $!"
