#!/bin/bash

# Destroy all AWS resources for an environment

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

print_warning "This will destroy ALL AWS resources for the $ENVIRONMENT environment!"
print_warning "This action cannot be undone."
echo ""
read -p "Are you sure you want to continue? Type 'yes' to confirm: " confirm

if [[ "$confirm" != "yes" ]]; then
    print_info "Destruction cancelled."
    exit 0
fi

print_info "Starting destruction of $ENVIRONMENT environment resources..."

# Delete stacks in reverse order (application services first, infrastructure last)
STACKS=(
    "${PROJECT_NAME}-${ENVIRONMENT}-deployment"
    "${PROJECT_NAME}-${ENVIRONMENT}-live"
    "${PROJECT_NAME}-${ENVIRONMENT}-web"
    "${PROJECT_NAME}-${ENVIRONMENT}-clickhouse"
    "${PROJECT_NAME}-${ENVIRONMENT}-postgres"
    "${PROJECT_NAME}-${ENVIRONMENT}-infra"
)

for stack in "${STACKS[@]}"; do
    echo ""
    print_info "Deleting stack: $stack"
    delete_stack "$stack"
done

echo ""
print_success "All resources destroyed for $ENVIRONMENT environment"
print_warning "Please verify in the AWS Console that all resources have been removed"
