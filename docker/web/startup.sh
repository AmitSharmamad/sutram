#!/bin/sh

set -e

echo "🚀 Starting web service with migrations..."

# Run PostgreSQL migrations
echo "📊 Running PostgreSQL migrations..."
./db-migrate -action=up
if [ $? -eq 0 ]; then
    echo "✅ PostgreSQL migrations completed successfully"
else
    echo "❌ PostgreSQL migrations failed"
    exit 1
fi

# Run ClickHouse migrations
echo "📈 Running ClickHouse migrations..."
./ch-migrate -action=up
if [ $? -eq 0 ]; then
    echo "✅ ClickHouse migrations completed successfully"
else
    echo "❌ ClickHouse migrations failed"
    exit 1
fi

# Start the web server
echo "🌐 Starting web server..."
exec ./web
