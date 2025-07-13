#!/bin/bash

# Deploy Infrastructure (VPC, networking)
# Usage: ./infra.sh [environment]

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
    --template-file "$TEMPLATES_DIR/infra.yaml" \
    --stack-name "$INFRA_STACK_NAME" \
    --parameter-overrides \
        Environment="$ENVIRONMENT" \
        ProjectName="$PROJECT_NAME" \
        UseSingleNATGateway="$USE_SINGLE_NAT_GATEWAY" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-cli-pager

echo "✅ Infrastructure deployed successfully!"
