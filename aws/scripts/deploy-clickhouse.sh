#!/bin/bash

# Deploy ClickHouse on ECS with EFS

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
POSTGRES_STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-postgres"
STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-clickhouse"

# Check if Infrastructure stack exists (for ECS cluster and service discovery)
if ! stack_exists "$INFRA_STACK_NAME"; then
    print_error "Infrastructure stack $INFRA_STACK_NAME does not exist. Deploy infrastructure first."
    exit 1
fi

# Check if PostgreSQL stack exists (for ECS security group)
if ! stack_exists "$POSTGRES_STACK_NAME"; then
    print_error "PostgreSQL stack $POSTGRES_STACK_NAME does not exist. Deploy PostgreSQL first."
    exit 1
fi

print_info "Deploying ClickHouse for environment: $ENVIRONMENT"

# Deploy ClickHouse stack
deploy_stack "$STACK_NAME" \
    "$SCRIPT_DIR/../templates/clickhouse.yaml" \
    "Environment=$ENVIRONMENT" \
    "ProjectName=$PROJECT_NAME" \
    "CPU=$CPU" \
    "Memory=$MEMORY"

print_success "ClickHouse deployment completed for $ENVIRONMENT environment"
