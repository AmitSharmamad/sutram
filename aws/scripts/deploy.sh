#!/bin/bash

# Deploy Sutram Infrastructure
# Usage: ./deploy.sh [environment] [service]
# Example: ./deploy.sh dev all
# Example: ./deploy.sh prod web

set -e

# Default values
ENVIRONMENT=${1:-dev}
SERVICE=${2:-all}

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_DIR="$PROJECT_ROOT/config"
TEMPLATES_DIR="$PROJECT_ROOT/templates"

# Source environment configuration
if [ -f "$CONFIG_DIR/${ENVIRONMENT}.env" ]; then
    source "$CONFIG_DIR/${ENVIRONMENT}.env"
else
    echo "Error: Configuration file for environment '$ENVIRONMENT' not found!"
    echo "Expected: $CONFIG_DIR/${ENVIRONMENT}.env"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

# Function to deploy a CloudFormation stack
deploy_stack() {
    local stack_name=$1
    local template_file=$2
    local parameters=$3
    
    log_info "Deploying stack: $stack_name"
    log_info "Template: $template_file"
    
    # Check if stack exists
    if aws cloudformation describe-stacks --stack-name "$stack_name" >/dev/null 2>&1; then
        log_info "Stack exists, updating..."
        aws cloudformation update-stack \
            --stack-name "$stack_name" \
            --template-body "file://$template_file" \
            --parameters "$parameters" \
            --capabilities CAPABILITY_NAMED_IAM \
            --no-cli-pager
        
        log_info "Waiting for stack update to complete..."
        aws cloudformation wait stack-update-complete --stack-name "$stack_name"
    else
        log_info "Stack does not exist, creating..."
        aws cloudformation create-stack \
            --stack-name "$stack_name" \
            --template-body "file://$template_file" \
            --parameters "$parameters" \
            --capabilities CAPABILITY_NAMED_IAM \
            --no-cli-pager
        
        log_info "Waiting for stack creation to complete..."
        aws cloudformation wait stack-create-complete --stack-name "$stack_name"
    fi
    
    log_success "Stack $stack_name deployed successfully!"
}

# Function to deploy infrastructure
deploy_infra() {
    log_info "Deploying infrastructure (VPC, networking)..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=UseSingleNATGateway,ParameterValue=$USE_SINGLE_NAT_GATEWAY"
    
    deploy_stack "$INFRA_STACK_NAME" "$TEMPLATES_DIR/infra.yaml" "$parameters"
}

# Function to deploy PostgreSQL
deploy_pg() {
    log_info "Deploying PostgreSQL database..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=InfraStackName,ParameterValue=$INFRA_STACK_NAME"
    parameters="$parameters ParameterKey=DatabaseInstanceClass,ParameterValue=$DATABASE_INSTANCE_CLASS"
    parameters="$parameters ParameterKey=DisableMultiAZ,ParameterValue=$DISABLE_MULTI_AZ"
    parameters="$parameters ParameterKey=MinimalStorage,ParameterValue=$MINIMAL_STORAGE"
    
    deploy_stack "$PG_STACK_NAME" "$TEMPLATES_DIR/pg.yaml" "$parameters"
}

# Function to deploy ClickHouse
deploy_ch() {
    log_info "Deploying ClickHouse database..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=InfraStackName,ParameterValue=$INFRA_STACK_NAME"
    parameters="$parameters ParameterKey=ClickHouseInstanceType,ParameterValue=$CLICKHOUSE_INSTANCE_TYPE"
    
    if [ -n "$KEY_NAME" ]; then
        parameters="$parameters ParameterKey=KeyName,ParameterValue=$KEY_NAME"
    fi
    
    deploy_stack "$CH_STACK_NAME" "$TEMPLATES_DIR/ch.yaml" "$parameters"
}

# Function to deploy web service
deploy_web() {
    log_info "Deploying web service..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=InfraStackName,ParameterValue=$INFRA_STACK_NAME"
    parameters="$parameters ParameterKey=PostgreSQLStackName,ParameterValue=$PG_STACK_NAME"
    parameters="$parameters ParameterKey=ClickHouseStackName,ParameterValue=$CH_STACK_NAME"
    parameters="$parameters ParameterKey=WebInstanceType,ParameterValue=$WEB_INSTANCE_TYPE"
    parameters="$parameters ParameterKey=DesiredCapacity,ParameterValue=$WEB_DESIRED_CAPACITY"
    parameters="$parameters ParameterKey=MinSize,ParameterValue=$WEB_MIN_SIZE"
    parameters="$parameters ParameterKey=MaxSize,ParameterValue=$WEB_MAX_SIZE"
    
    if [ -n "$KEY_NAME" ]; then
        parameters="$parameters ParameterKey=KeyName,ParameterValue=$KEY_NAME"
    fi
    
    deploy_stack "$WEB_STACK_NAME" "$TEMPLATES_DIR/web.yaml" "$parameters"
}

# Function to deploy live service
deploy_live() {
    log_info "Deploying live service..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=InfraStackName,ParameterValue=$INFRA_STACK_NAME"
    parameters="$parameters ParameterKey=PostgreSQLStackName,ParameterValue=$PG_STACK_NAME"
    parameters="$parameters ParameterKey=ClickHouseStackName,ParameterValue=$CH_STACK_NAME"
    parameters="$parameters ParameterKey=LiveInstanceType,ParameterValue=$LIVE_INSTANCE_TYPE"
    parameters="$parameters ParameterKey=DesiredCapacity,ParameterValue=$LIVE_DESIRED_CAPACITY"
    parameters="$parameters ParameterKey=MinSize,ParameterValue=$LIVE_MIN_SIZE"
    parameters="$parameters ParameterKey=MaxSize,ParameterValue=$LIVE_MAX_SIZE"
    
    if [ -n "$KEY_NAME" ]; then
        parameters="$parameters ParameterKey=KeyName,ParameterValue=$KEY_NAME"
    fi
    
    deploy_stack "$LIVE_STACK_NAME" "$TEMPLATES_DIR/live.yaml" "$parameters"
}

# Function to deploy deployment service
deploy_deployment() {
    log_info "Deploying deployment service..."
    
    local parameters="ParameterKey=Environment,ParameterValue=$ENVIRONMENT"
    parameters="$parameters ParameterKey=ProjectName,ParameterValue=$PROJECT_NAME"
    parameters="$parameters ParameterKey=InfraStackName,ParameterValue=$INFRA_STACK_NAME"
    parameters="$parameters ParameterKey=PostgreSQLStackName,ParameterValue=$PG_STACK_NAME"
    parameters="$parameters ParameterKey=ClickHouseStackName,ParameterValue=$CH_STACK_NAME"
    parameters="$parameters ParameterKey=DeploymentInstanceType,ParameterValue=$DEPLOYMENT_INSTANCE_TYPE"
    parameters="$parameters ParameterKey=DesiredCapacity,ParameterValue=$DEPLOYMENT_DESIRED_CAPACITY"
    parameters="$parameters ParameterKey=MinSize,ParameterValue=$DEPLOYMENT_MIN_SIZE"
    parameters="$parameters ParameterKey=MaxSize,ParameterValue=$DEPLOYMENT_MAX_SIZE"
    
    if [ -n "$KEY_NAME" ]; then
        parameters="$parameters ParameterKey=KeyName,ParameterValue=$KEY_NAME"
    fi
    
    deploy_stack "$DEPLOYMENT_STACK_NAME" "$TEMPLATES_DIR/deployment.yaml" "$parameters"
}

# Function to show deployment status
show_status() {
    log_info "Deployment Status for environment: $ENVIRONMENT"
    echo
    
    stacks=("$INFRA_STACK_NAME" "$PG_STACK_NAME" "$CH_STACK_NAME" "$WEB_STACK_NAME" "$LIVE_STACK_NAME" "$DEPLOYMENT_STACK_NAME")
    
    for stack in "${stacks[@]}"; do
        if aws cloudformation describe-stacks --stack-name "$stack" >/dev/null 2>&1; then
            status=$(aws cloudformation describe-stacks --stack-name "$stack" --query 'Stacks[0].StackStatus' --output text)
            echo -e "${GREEN}✓${NC} $stack: $status"
        else
            echo -e "${RED}✗${NC} $stack: NOT_DEPLOYED"
        fi
    done
}

# Main deployment logic
case $SERVICE in
    "all")
        log_info "Deploying all services for environment: $ENVIRONMENT"
        deploy_infra
        deploy_pg
        deploy_ch
        deploy_web
        deploy_live
        deploy_deployment
        log_success "All services deployed successfully!"
        show_status
        ;;
    "infra")
        deploy_infra
        ;;
    "pg")
        deploy_pg
        ;;
    "ch")
        deploy_ch
        ;;
    "web")
        deploy_web
        ;;
    "live")
        deploy_live
        ;;
    "deployment")
        deploy_deployment
        ;;
    "status")
        show_status
        ;;
    *)
        log_error "Unknown service: $SERVICE"
        echo "Usage: $0 [environment] [service]"
        echo "Environment: dev, prod"
        echo "Service: all, infra, pg, ch, web, live, deployment, status"
        exit 1
        ;;
esac
