#!/bin/bash
set -e

cd "$(dirname "$0")/.."

PID_FILE="vigilate.pid"
CONFIG_FILE="${1:-config.yaml}"

if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
    echo "vigilate is already running (PID: $(cat $PID_FILE))"
    exit 1
fi

echo "Starting vigilate with config: $CONFIG_FILE"
nohup ./vigilate "$CONFIG_FILE" > vigilate.log 2>&1 &
echo $! > "$PID_FILE"
echo "Started with PID: $!"
