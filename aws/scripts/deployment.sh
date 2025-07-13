#!/bin/bash

# Deploy Deployment Service (ECS)
# Usage: ./deployment.sh [environment]

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
    --template-file "$TEMPLATES_DIR/deployment-ecs.yaml" \
    --stack-name "$DEPLOYMENT_STACK_NAME" \
    --parameter-overrides \
        Environment="$ENVIRONMENT" \
        ProjectName="$PROJECT_NAME" \
        InfraStackName="$INFRA_STACK_NAME" \
        ImageTag="$IMAGE_TAG" \
        TaskCpu="$TASK_CPU" \
        TaskMemory="$TASK_MEMORY" \
        DesiredCapacity="$DEPLOYMENT_DESIRED_CAPACITY" \
        MinSize="$DEPLOYMENT_MIN_SIZE" \
        MaxSize="$DEPLOYMENT_MAX_SIZE" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-cli-pager

echo "✅ Deployment service deployed successfully!"
