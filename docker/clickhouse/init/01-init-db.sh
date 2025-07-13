#!/bin/bash
set -e

# Wait for ClickHouse to be ready
until clickhouse-client --query "SELECT 1" > /dev/null 2>&1; do
    echo "Waiting for ClickHouse to be ready..."
    sleep 2
done

echo "ClickHouse is ready, creating database and user..."

# Create database
clickhouse-client --query "CREATE DATABASE IF NOT EXISTS sutram"

# Create admin user if it doesn't exist
clickhouse-client --query "CREATE USER IF NOT EXISTS admin IDENTIFIED BY 'REDACTED_PASSWORD'"

# Grant permissions to admin user
clickhouse-client --query "GRANT ALL ON sutram.* TO admin"
clickhouse-client --query "GRANT ALL ON system.* TO admin"

echo "Database and user setup completed!"
