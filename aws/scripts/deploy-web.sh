#!/bin/bash

# Deploy Web Service on ECS

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
FORCE_UPDATE="$(date +%Y%m%d-%H%M%S)"

# Validate environment
if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
    print_error "Environment must be 'dev' or 'prod'"
    exit 1
fi

# Load environment variables
ENV_FILE="$SCRIPT_DIR/../config/${ENVIRONMENT}.env"
load_env "$ENV_FILE"

# Set default values for optional OAuth parameters
GOOGLE_OAUTH_CLIENT_ID="${GOOGLE_OAUTH_CLIENT_ID:-527787289499-ae0fqthcusd5jpti5d2liq5mgr667llv.apps.googleusercontent.com}"
GOOGLE_OAUTH_CLIENT_SECRET="${GOOGLE_OAUTH_CLIENT_SECRET:-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}"
GOOGLE_OAUTH_REDIRECT_URL="${GOOGLE_OAUTH_REDIRECT_URL:-}"

# Check AWS CLI
check_aws

# Stack names
CLICKHOUSE_STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-clickhouse"
STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-web"

# Check if ClickHouse stack exists (for ECS cluster)
if ! stack_exists "$CLICKHOUSE_STACK_NAME"; then
    print_error "ClickHouse stack $CLICKHOUSE_STACK_NAME does not exist. Deploy ClickHouse first."
    exit 1
fi

print_info "Deploying Web service for environment: $ENVIRONMENT with image tag: $IMAGE_TAG"

# Validate required parameters
if [[ -z "$PUBLIC_HOSTED_ZONE_ID" ]]; then
    print_error "PUBLIC_HOSTED_ZONE_ID is required but not set"
    exit 1
fi

if [[ -z "$BASE_DOMAIN_NAME" ]]; then
    print_error "BASE_DOMAIN_NAME is required but not set"
    exit 1
fi

# Deploy Web service stack
deploy_stack "$STACK_NAME" \
    "$SCRIPT_DIR/../templates/web.yaml" \
    "Environment=$ENVIRONMENT" \
    "ProjectName=$PROJECT_NAME" \
    "CPU=$CPU" \
    "Memory=$MEMORY" \
    "WebPort=$WEB_PORT" \
    "ImageTag=$IMAGE_TAG" \
    "ForceUpdate=$FORCE_UPDATE" \
    "GoogleOAuthClientId=$GOOGLE_OAUTH_CLIENT_ID" \
    "GoogleOAuthClientSecret=$GOOGLE_OAUTH_CLIENT_SECRET" \
    "GoogleOAuthRedirectURL=$GOOGLE_OAUTH_REDIRECT_URL" \
    "PublicHostedZoneId=$PUBLIC_HOSTED_ZONE_ID" \
    "BaseDomainName=${BASE_DOMAIN_NAME}" \
    "ExistingCertificateArn=${EXISTING_CERTIFICATE_ARN:-}"

# Force ECS service update to ensure new deployment
SERVICE_NAME="${PROJECT_NAME}-${ENVIRONMENT}-web"
CLUSTER_NAME="${PROJECT_NAME}-${ENVIRONMENT}-cluster"

print_info "Forcing ECS service update for $SERVICE_NAME on cluster $CLUSTER_NAME..."

# Verify cluster exists first
if ! aws ecs describe-clusters --clusters "$CLUSTER_NAME" --query 'clusters[0].status' --output text --no-cli-pager | grep -q ACTIVE; then
    print_error "ECS cluster $CLUSTER_NAME not found or not active. Please check cluster name."
    exit 1
fi

# Verify service exists
if ! aws ecs describe-services --cluster "$CLUSTER_NAME" --services "$SERVICE_NAME" --query 'services[0].status' --output text --no-cli-pager | grep -q ACTIVE; then
    print_error "ECS service $SERVICE_NAME not found or not active in cluster $CLUSTER_NAME."
    exit 1
fi

aws ecs update-service \
    --cluster "$CLUSTER_NAME" \
    --service "$SERVICE_NAME" \
    --force-new-deployment \
    --no-cli-pager > /dev/null

print_info "Waiting for service to stabilize..."
aws ecs wait services-stable \
    --cluster "$CLUSTER_NAME" \
    --services "$SERVICE_NAME"

# Get service URL
SERVICE_URL=$(get_stack_output "$STACK_NAME" "ServiceURL")
print_success "Web service deployment completed for $ENVIRONMENT environment"
print_info "Service URL: $SERVICE_URL"
