#!/bin/bash

# Deploy PostgreSQL RDS

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/common.sh"

# Check arguments
if [[ $# -ne 1 ]]; then
    print_error "Usage: $0 <environment>"
    print_error "Environment: dev or prod"
    exit 1
fi

ENVIRONMENT="$1"

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
INFRA_STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-infra"
STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-postgres"

# Check if infrastructure stack exists
if ! stack_exists "$INFRA_STACK_NAME"; then
    print_error "Infrastructure stack $INFRA_STACK_NAME does not exist. Deploy infrastructure first."
    exit 1
fi

print_info "Deploying PostgreSQL for environment: $ENVIRONMENT"

# Deploy PostgreSQL stack
deploy_stack "$STACK_NAME" \
    "$SCRIPT_DIR/../templates/postgres.yaml" \
    "Environment=$ENVIRONMENT" \
    "ProjectName=$PROJECT_NAME" \
    "DBInstanceClass=$DB_INSTANCE_CLASS" \
    "DBAllocatedStorage=$DB_ALLOCATED_STORAGE" \
    "DBName=$DB_NAME" \
    "DBUsername=$DB_USERNAME" \
    "DBPassword=$DB_PASSWORD" \
    "BaseDomainName=${BASE_DOMAIN_NAME}" \
    "PublicHostedZoneId=$PUBLIC_HOSTED_ZONE_ID"

# Get database endpoint
DB_ENDPOINT=$(get_stack_output "$STACK_NAME" "DatabaseEndpoint")
print_success "PostgreSQL deployment completed for $ENVIRONMENT environment"
print_info "Database endpoint: $DB_ENDPOINT"
