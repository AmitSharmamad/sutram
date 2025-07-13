#!/bin/bash

# Deploy PostgreSQL Database
# Usage: ./pg.sh [environment]

set -e

ENVIRONMENT=${1:-dev}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_DIR="$PROJECT_ROOT/config"
TEMPLATES_DIR="$PROJECT_ROOT/templates"

# Source environment configuration
source "$CONFIG_DIR/${ENVIRONMENT}.env"

# Deploy CloudFormation stack
aws cloudformation deploy \
    --template-file "$TEMPLATES_DIR/pg.yaml" \
    --stack-name "$PG_STACK_NAME" \
    --parameter-overrides \
        Environment="$ENVIRONMENT" \
        ProjectName="$PROJECT_NAME" \
        InfraStackName="$INFRA_STACK_NAME" \
        DatabaseInstanceClass="$DATABASE_INSTANCE_CLASS" \
        DisableMultiAZ="$DISABLE_MULTI_AZ" \
        MinimalStorage="$MINIMAL_STORAGE" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-cli-pager

echo "✅ PostgreSQL deployed successfully!"
