#!/bin/sh
set -e

echo "🚀 Running migration..."
/goose -dir ./db/migrations postgres "$USER_DB_URL" up

echo "🎯 Starting service..."
exec ./user-svc
