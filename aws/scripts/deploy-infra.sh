#!/bin/bash

# Deploy infrastructure (VPC, subnets, etc.)

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

# Stack name
STACK_NAME="${PROJECT_NAME}-${ENVIRONMENT}-infra"

print_info "Deploying infrastructure for environment: $ENVIRONMENT"

# Deploy infrastructure stack
deploy_stack "$STACK_NAME" \
    "$SCRIPT_DIR/../templates/infra.yaml" \
    "Environment=$ENVIRONMENT" \
    "ProjectName=$PROJECT_NAME" \
    "VpcCidr=$VPC_CIDR" \
    "PublicSubnet1Cidr=$PUBLIC_SUBNET_1_CIDR" \
    "PublicSubnet2Cidr=$PUBLIC_SUBNET_2_CIDR" \
    "PrivateSubnet1Cidr=$PRIVATE_SUBNET_1_CIDR" \
    "PrivateSubnet2Cidr=$PRIVATE_SUBNET_2_CIDR"

print_success "Infrastructure deployment completed for $ENVIRONMENT environment"
