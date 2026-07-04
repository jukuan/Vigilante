#!/bin/bash
set -e

cd "$(dirname "$0")/.."

IMAGE="vigilante"
CONTAINER="vigilante"
CONFIG="${1:-$(pwd)/config.yaml}"
LOG_DIR="${2:-/var/log}"

# Check config exists
if [ ! -f "$CONFIG" ]; then
    echo "Error: Config file not found: $CONFIG"
    echo "Usage: $0 [config.yaml] [log-directory]"
    exit 1
fi

# Stop and remove existing container
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
    echo "Removing existing container..."
    docker stop "$CONTAINER" 2>/dev/null || true
    docker rm "$CONTAINER" 2>/dev/null || true
fi

# Build image
echo "Building Docker image..."
docker build -t "$IMAGE" .

# Run container
echo "Starting container..."
docker run -d \
    --name "$CONTAINER" \
    --restart unless-stopped \
    -v "$CONFIG:/app/config.yaml:ro" \
    -v "$LOG_DIR:$LOG_DIR:ro" \
    "$IMAGE"

echo ""
echo "Container started. Check logs with:"
echo "  docker logs -f $CONTAINER"
echo ""
echo "Stop with:"
echo "  docker stop $CONTAINER"
