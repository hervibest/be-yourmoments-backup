#!/bin/sh
set -e

echo "🚀 Running migration..."
/goose -dir ./db/migrations postgres "$NOTIFICATION_DB_URL" up

echo "🎯 Starting service..."
exec ./notification-svc
