#!/bin/bash
# Dummy alert script: writes the alert message to a local file

MESSAGE="$1"
LOGFILE="alerts.log"

if [ -z "$MESSAGE" ]; then
    echo "Usage: $0 <message>" >&2
    exit 1
fi

echo "[$(date '+%Y-%m-%d %H:%M:%S')] $MESSAGE" >> "$LOGFILE"
echo "Logged alert to $LOGFILE"
