#!/bin/bash
# Generate test log lines into a directory

DIR="${1:-/tmp/vigilante-demo}"
mkdir -p "$DIR"

while true; do
  case $((RANDOM % 5)) in
    0) echo "[$(date)] INFO: all systems nominal" >> "$DIR/app.log" ;;
    1) echo "[$(date)] WARN: high memory usage" >> "$DIR/app.log" ;;
    2) echo "[$(date)] ERROR: connection timeout" >> "$DIR/app.log" ;;
    3) echo "[$(date)] FATAL: disk failure imminent" >> "$DIR/app.log" ;;
    4) echo "[$(date)] DEBUG: processing request" >> "$DIR/app.log" ;;
  esac
  sleep $((1 + RANDOM % 5))
done
