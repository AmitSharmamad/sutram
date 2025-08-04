#!/bin/bash

# Deploy Deployment Service on ECS

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/common.sh"

# Check arguments
if [[ $# -lt 1 || $# -gt 2 ]]; then
    print_error "Usage: $0 <environment> [image_tag]"
    print_error "Environment: dev or prod"
    print_error "Image tag: optional, defaults to 'latest'"
    exit 1
fi

ENVIRONMENT="$1"
IMAGE_TAG="${2:-latest}"

# Validate environment
if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
    print_error "Environment must be 'dev' or 'prod'"
    exit 1
fi

# Load environment variables
ENV_FILE="$SCRIPT_DIR/../config/${ENVIRONMENT}.env"
load_env "$ENV_FILE"

# Check AWS CLI
check_aws

# Stack names
CLICKHOUSE_STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-clickhouse"
STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-deployment"

# Check if ClickHouse stack exists (for ECS cluster)
if ! stack_exists "$CLICKHOUSE_STACK_NAME"; then
    print_error "ClickHouse stack $CLICKHOUSE_STACK_NAME does not exist. Deploy ClickHouse first."
    exit 1
fi

print_info "Deploying Deployment service for environment: $ENVIRONMENT with image tag: $IMAGE_TAG"

# Deploy Deployment service stack
deploy_stack "$STACK_NAME" \
    "$SCRIPT_DIR/../templates/deployment.yaml" \
    "Environment=$ENVIRONMENT" \
    "ProjectName=$PROJECT_NAME" \
    "CPU=$CPU" \
    "Memory=$MEMORY" \
    "DeploymentPort=$DEPLOYMENT_PORT" \
    "ImageTag=$IMAGE_TAG" \
    "BaseDomainName=${BASE_DOMAIN_NAME}" \
    "PublicHostedZoneId=$PUBLIC_HOSTED_ZONE_ID" \
    "ExistingCertificateArn=${EXISTING_CERTIFICATE_ARN:-}"

print_success "Deployment service deployment completed for $ENVIRONMENT environment"
