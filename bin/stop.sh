#!/bin/bash
set -e

cd "$(dirname "$0")/.."

PID_FILE="vigilante.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "No PID file found. Service may not be running."
    exit 1
fi

PID=$(cat "$PID_FILE")

if ! kill -0 $PID 2>/dev/null; then
    echo "Process $PID is not running. Removing stale PID file."
    rm -f "$PID_FILE"
    exit 1
fi

echo "Stopping vigilante (PID: $PID)..."
kill -TERM $PID

# Wait for process to stop
for i in {1..10}; do
    if ! kill -0 $PID 2>/dev/null; then
        echo "vigilante stopped"
        rm -f "$PID_FILE"
        exit 0
    fi
    sleep 1
done

echo "vigilante did not stop gracefully, forcing..."
kill -KILL $PID 2>/dev/null || true
rm -f "$PID_FILE"
echo "vigilante forcefully stopped"
