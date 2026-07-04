#!/bin/bash
set -e
cd "$(dirname "$0")/.."

DEMO_DIR="/tmp/vigilante-demo"
CONFIG="config.test.yaml"

echo "=== Integration Test ==="

# Setup
mkdir -p "$DEMO_DIR"
rm -f "$DEMO_DIR/app.log" alerts.log state.json

# Create test config with short cooldown
cat > "$CONFIG" <<EOF
inactivity_seconds: 60
state_file: state.json
rules:
  - name: test-errors
    log_dir: $DEMO_DIR
    file_glob: "*.log"
    patterns:
      - "ERROR"
      - "FATAL"
    actions:
      - scripts/dummy-alert.sh
    cooldown_seconds: 5
EOF

# Build
echo "Building..."
go build -o vigilante .

# Start vigilante
echo "Starting vigilante..."
./vigilante "$CONFIG" > test-vigilante.log 2>&1 &
PID=$!
echo "vigilante PID: $PID"

# Give it a moment to start
sleep 2

# Write some log lines
echo "Writing test log lines..."
echo "INFO: everything ok" >> "$DEMO_DIR/app.log"
echo "ERROR: something went wrong" >> "$DEMO_DIR/app.log"
echo "DEBUG: processing" >> "$DEMO_DIR/app.log"
echo "FATAL: system crash" >> "$DEMO_DIR/app.log"
echo "ERROR: another error" >> "$DEMO_DIR/app.log"

# Wait for cooldown to fire
echo "Waiting for alert cooldown (7 seconds)..."
sleep 7

# Stop vigilante gracefully
echo "Stopping vigilante..."
kill -TERM $PID 2>/dev/null || true
sleep 1
kill -KILL $PID 2>/dev/null || true
rm -f vigilante.pid

# Check results
echo ""
if [ -f alerts.log ]; then
    echo "=== Alerts generated ==="
    cat alerts.log
    echo ""
    echo "✓ Integration test passed"
else
    echo "✗ No alerts generated — test failed"
    echo "vigilante output:"
    cat test-vigilante.log
    exit 1
fi

# Cleanup
rm -f "$CONFIG" test-vigilante.log state.json alerts.log
rm -rf "$DEMO_DIR"
