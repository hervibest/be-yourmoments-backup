#!/bin/sh
set -e

echo "🚀 Running migration..."
/goose -dir ./db/migrations postgres "$TRANSACTION_DB_URL" up

echo "🎯 Starting service..."
exec ./transaction-svc
