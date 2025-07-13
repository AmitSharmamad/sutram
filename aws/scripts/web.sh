#!/bin/bash

# Deploy Web Service (ECS)
# Usage: ./web.sh [environment]

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
    --template-file "$TEMPLATES_DIR/web-ecs.yaml" \
    --stack-name "$WEB_STACK_NAME" \
    --parameter-overrides \
        Environment="$ENVIRONMENT" \
        ProjectName="$PROJECT_NAME" \
        InfraStackName="$INFRA_STACK_NAME" \
        ImageTag="$IMAGE_TAG" \
        TaskCpu="$TASK_CPU" \
        TaskMemory="$TASK_MEMORY" \
        DesiredCapacity="$WEB_DESIRED_CAPACITY" \
        MinSize="$WEB_MIN_SIZE" \
        MaxSize="$WEB_MAX_SIZE" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-cli-pager

echo "✅ Web service deployed successfully!"
