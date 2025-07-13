#!/bin/bash

# Deploy all services in order
# Usage: ./deploy-all.sh [environment]

set -e

ENVIRONMENT=${1:-dev}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Deploying all services for environment: $ENVIRONMENT"

# Deploy in order (infrastructure first, then databases, then services)
echo "1. Deploying infrastructure..."
"$SCRIPT_DIR/infra.sh" "$ENVIRONMENT"

echo "2. Deploying PostgreSQL..."
"$SCRIPT_DIR/pg.sh" "$ENVIRONMENT"

echo "3. Deploying ClickHouse..."
"$SCRIPT_DIR/ch.sh" "$ENVIRONMENT"

echo "4. Deploying Web service..."
"$SCRIPT_DIR/web.sh" "$ENVIRONMENT"

echo "5. Deploying Live service..."
"$SCRIPT_DIR/live.sh" "$ENVIRONMENT"

echo "6. Deploying Deployment service..."
"$SCRIPT_DIR/deployment.sh" "$ENVIRONMENT"

echo "✅ All services deployed successfully!"
