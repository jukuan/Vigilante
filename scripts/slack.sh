#!/bin/bash
# Example Slack alert script
# Usage: ./slack.sh "ALERT: 12 lines in logs for last 5 minutes with like FATAL"

MESSAGE="$1"
WEBHOOK_URL="${SLACK_WEBHOOK_URL:-https://hooks.slack.com/services/your-webhook-url}"

# Replace with actual Slack webhook call
curl -X POST -H 'Content-type: application/json' \
  --data "{\"text\":\"$MESSAGE\"}" \
  "$WEBHOOK_URL" 2>/dev/null

echo "Slack notification sent: $MESSAGE"
