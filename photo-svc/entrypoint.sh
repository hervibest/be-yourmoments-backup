#!/bin/sh
set -e

echo "🚀 Running migration..."
/goose -dir ./db/migrations postgres "$PHOTO_DB_URL" up

echo "🎯 Starting service..."
exec ./photo-svc
