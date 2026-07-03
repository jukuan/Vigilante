#!/bin/bash
# Example email alert script
# Usage: ./email.sh "ALERT: 12 lines in logs for last 5 minutes with like FATAL"

MESSAGE="$1"
SUBJECT="Watchlog Alert"

# Replace with actual email sending logic
echo "$MESSAGE" | mail -s "$SUBJECT" admin@example.com

echo "Email sent: $SUBJECT - $MESSAGE"
