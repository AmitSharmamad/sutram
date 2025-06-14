#!/bin/bash

# Sutram AWS CloudFormation Deployment Scripts
# This script provides functions to deploy and manage CloudFormation stacks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT="dev"
CONFIG_FILE=""
DRY_RUN=false

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Load environment configuration
load_config() {
    local env=$1
    CONFIG_FILE="config/${env}.env"
    
    if [[ ! -f "$CONFIG_FILE" ]]; then
        log_error "Configuration file $CONFIG_FILE not found!"
        exit 1
    fi
    
    log_info "Loading configuration from $CONFIG_FILE"
    set -a
    source "$CONFIG_FILE"
    set +a
    
    # Validate required variables
    if [[ -z "$PROJECT_NAME" || -z "$ENVIRONMENT" || -z "$AWS_REGION" ]]; then
        log_error "Missing required configuration variables!"
        exit 1
    fi
}

# Get stack name with environment prefix
get_stack_name() {
    local stack_type=$1
    echo "${PROJECT_NAME}-${ENVIRONMENT}-${stack_type}"
}

# Check if stack exists
stack_exists() {
    local stack_name=$1
    aws cloudformation describe-stacks --stack-name "$stack_name" --region "$AWS_REGION" > /dev/null 2>&1
}

# Wait for stack operation to complete
wait_for_stack() {
    local stack_name=$1
    local operation=$2
    
    log_info "Waiting for stack $operation to complete: $stack_name"
    
    if [[ "$operation" == "DELETE" ]]; then
        aws cloudformation wait stack-delete-complete --stack-name "$stack_name" --region "$AWS_REGION"
    else
        aws cloudformation wait stack-create-complete --stack-name "$stack_name" --region "$AWS_REGION" 2>/dev/null || \
        aws cloudformation wait stack-update-complete --stack-name "$stack_name" --region "$AWS_REGION"
    fi
}

# Deploy a single stack
deploy_stack() {
    local stack_type=$1
    local template_file="templates/${stack_type}.yaml"
    local stack_name=$(get_stack_name "$stack_type")
    
    if [[ ! -f "$template_file" ]]; then
        log_error "Template file $template_file not found!"
        return 1
    fi
    
    log_info "Deploying stack: $stack_name"
    
    # Prepare parameters
    local params=""
    
    case "$stack_type" in
        "01-networking")
            params="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
            params="$params ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
            # Cost optimization parameters
            if [[ -n "$USE_SINGLE_NAT_GATEWAY" ]]; then
                params="$params ParameterKey=UseSingleNATGateway,ParameterValue=$USE_SINGLE_NAT_GATEWAY"
            fi
            ;;
        "02-data")
            params="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
            params="$params ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
            params="$params ParameterKey=NetworkingStackName,ParameterValue=$(get_stack_name '01-networking')"
            params="$params ParameterKey=DatabaseInstanceClass,ParameterValue=$DATABASE_INSTANCE_CLASS"
            params="$params ParameterKey=ClickHouseInstanceType,ParameterValue=$CLICKHOUSE_INSTANCE_TYPE"
            params="$params ParameterKey=RedisNodeType,ParameterValue=$REDIS_NODE_TYPE"
            # Cost optimization parameters
            if [[ -n "$DISABLE_MULTI_AZ" ]]; then
                params="$params ParameterKey=DisableMultiAZ,ParameterValue=$DISABLE_MULTI_AZ"
            fi
            if [[ -n "$MINIMAL_STORAGE" ]]; then
                params="$params ParameterKey=MinimalStorage,ParameterValue=$MINIMAL_STORAGE"
            fi
            ;;
        "03-app-infrastructure")
            params="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
            params="$params ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
            params="$params ParameterKey=NetworkingStackName,ParameterValue=$(get_stack_name '01-networking')"
            if [[ -n "$CERTIFICATE_ARN" ]]; then
                params="$params ParameterKey=CertificateArn,ParameterValue=$CERTIFICATE_ARN"
            fi
            if [[ -n "$DOMAIN_NAME" ]]; then
                params="$params ParameterKey=DomainName,ParameterValue=$DOMAIN_NAME"
            fi
            # Cost optimization parameters
            if [[ -n "$USE_FARGATE_SPOT" ]]; then
                params="$params ParameterKey=UseFargateSpot,ParameterValue=$USE_FARGATE_SPOT"
            fi
            if [[ -n "$DISABLE_CONTAINER_INSIGHTS" ]]; then
                params="$params ParameterKey=DisableContainerInsights,ParameterValue=$DISABLE_CONTAINER_INSIGHTS"
            fi
            ;;
        "04-app-services")
            params="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
            params="$params ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
            params="$params ParameterKey=NetworkingStackName,ParameterValue=$(get_stack_name '01-networking')"
            params="$params ParameterKey=DataStackName,ParameterValue=$(get_stack_name '02-data')"
            params="$params ParameterKey=AppInfraStackName,ParameterValue=$(get_stack_name '03-app-infrastructure')"
            params="$params ParameterKey=WebAppImageTag,ParameterValue=$WEB_APP_IMAGE_TAG"
            params="$params ParameterKey=LiveTradingImageTag,ParameterValue=$LIVE_TRADING_IMAGE_TAG"
            params="$params ParameterKey=WebAppDesiredCount,ParameterValue=$WEB_APP_DESIRED_COUNT"
            params="$params ParameterKey=LiveTradingDesiredCount,ParameterValue=$LIVE_TRADING_DESIRED_COUNT"
            # Cost optimization parameters
            if [[ -n "$MINIMAL_RESOURCES" ]]; then
                params="$params ParameterKey=MinimalResources,ParameterValue=$MINIMAL_RESOURCES"
            fi
            ;;
        "05-cdn")
            params="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
            params="$params ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
            params="$params ParameterKey=AppInfraStackName,ParameterValue=$(get_stack_name '03-app-infrastructure')"
            params="$params ParameterKey=DataStackName,ParameterValue=$(get_stack_name '02-data')"
            if [[ -n "$DOMAIN_NAME" ]]; then
                params="$params ParameterKey=DomainName,ParameterValue=$DOMAIN_NAME"
            fi
            if [[ -n "$CREATE_ROUTE53_HOSTED_ZONE" ]]; then
                params="$params ParameterKey=CreateRoute53HostedZone,ParameterValue=$CREATE_ROUTE53_HOSTED_ZONE"
            fi
            if [[ -n "$ROUTE53_HOSTED_ZONE_ID" ]]; then
                params="$params ParameterKey=Route53HostedZoneId,ParameterValue=$ROUTE53_HOSTED_ZONE_ID"
            fi
            # Cost optimization parameters
            if [[ -n "$DISABLE_CDN" ]]; then
                params="$params ParameterKey=DisableCDN,ParameterValue=$DISABLE_CDN"
            fi
            if [[ -n "$DISABLE_WAF" ]]; then
                params="$params ParameterKey=DisableWAF,ParameterValue=$DISABLE_WAF"
            fi
            ;;
    esac
    
    # Check if this is a dry run
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would deploy stack $stack_name with parameters:"
        echo "$params" | tr ' ' '\n'
        return 0
    fi
    
    # Check if stack exists to determine operation
    local operation="create-stack"
    if stack_exists "$stack_name"; then
        operation="update-stack"
        log_info "Stack exists, updating..."
    else
        log_info "Stack does not exist, creating..."
    fi
    
    # Deploy the stack
    local param_args=""
    if [[ -n "$params" ]]; then
        param_args="--parameters $params"
    fi
    
    aws cloudformation "$operation" \
        --stack-name "$stack_name" \
        --template-body "file://$template_file" \
        --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
        --region "$AWS_REGION" \
        $param_args
    
    # Wait for completion
    wait_for_stack "$stack_name" "DEPLOY"
    log_success "Stack deployed successfully: $stack_name"
}

# Delete a single stack
delete_stack() {
    local stack_type=$1
    local stack_name=$(get_stack_name "$stack_type")
    
    if ! stack_exists "$stack_name"; then
        log_warning "Stack does not exist: $stack_name"
        return 0
    fi
    
    log_info "Deleting stack: $stack_name"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would delete stack $stack_name"
        return 0
    fi
    
    aws cloudformation delete-stack --stack-name "$stack_name" --region "$AWS_REGION"
    wait_for_stack "$stack_name" "DELETE"
    log_success "Stack deleted successfully: $stack_name"
}

# Validate template
validate_template() {
    local template_file=$1
    
    log_info "Validating template: $template_file"
    aws cloudformation validate-template --template-body "file://$template_file" --region "$AWS_REGION" > /dev/null
    log_success "Template is valid: $template_file"
}

# Show stack outputs
show_outputs() {
    local stack_type=$1
    local stack_name=$(get_stack_name "$stack_type")
    
    if ! stack_exists "$stack_name"; then
        log_error "Stack does not exist: $stack_name"
        return 1
    fi
    
    log_info "Stack outputs for: $stack_name"
    aws cloudformation describe-stacks \
        --stack-name "$stack_name" \
        --region "$AWS_REGION" \
        --query 'Stacks[0].Outputs[*].[OutputKey,OutputValue,Description]' \
        --output table
}

# Get stack status
get_stack_status() {
    local stack_type=$1
    local stack_name=$(get_stack_name "$stack_type")
    
    if ! stack_exists "$stack_name"; then
        echo "NOT_EXISTS"
        return 0
    fi
    
    aws cloudformation describe-stacks \
        --stack-name "$stack_name" \
        --region "$AWS_REGION" \
        --query 'Stacks[0].StackStatus' \
        --output text
}

# Export functions for use by Makefile
export -f log_info log_success log_warning log_error
export -f load_config get_stack_name stack_exists wait_for_stack
export -f deploy_stack delete_stack validate_template show_outputs get_stack_status
