#!/bin/bash

# Destroy Sutram Infrastructure
# Usage: ./destroy.sh [environment]

set -e

ENVIRONMENT=${1:-dev}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_DIR="$PROJECT_ROOT/config"

# Source environment configuration
source "$CONFIG_DIR/${ENVIRONMENT}.env"

# Function to delete stack if it exists
delete_stack() {
    local stack_name=$1
    if aws cloudformation describe-stacks --stack-name "$stack_name" >/dev/null 2>&1; then
        echo "🗑️  Deleting $stack_name..."
        aws cloudformation delete-stack --stack-name "$stack_name"
        aws cloudformation wait stack-delete-complete --stack-name "$stack_name"
        echo "✅ $stack_name deleted"
    else
        echo "ℹ️  Stack $stack_name does not exist"
    fi
}

# Function to delete a CloudFormation stack
delete_stack() {
    local stack_name=$1
    
    # Check if stack exists
    if aws cloudformation describe-stacks --stack-name "$stack_name" >/dev/null 2>&1; then
        log_info "Deleting stack: $stack_name"
        aws cloudformation delete-stack --stack-name "$stack_name"
        
        log_info "Waiting for stack deletion to complete..."
        aws cloudformation wait stack-delete-complete --stack-name "$stack_name"
        log_success "Stack $stack_name deleted successfully!"
    else
        log_warning "Stack $stack_name does not exist, skipping..."
    fi
}

# Function to delete deployment service
delete_deployment() {
    log_info "Deleting deployment service..."
    delete_stack "$DEPLOYMENT_STACK_NAME"
}

echo "🚨 WARNING: This will permanently delete all AWS resources for environment: $ENVIRONMENT"
read -p "Type 'yes' to confirm deletion: " confirmation

if [ "$confirmation" != "yes" ]; then
    echo "❌ Deletion cancelled"
    exit 0
fi

echo "🗑️  Deleting all stacks in reverse order..."

# Delete services first (reverse order of deployment)
delete_stack "$DEPLOYMENT_STACK_NAME"
delete_stack "$LIVE_STACK_NAME"
delete_stack "$WEB_STACK_NAME"

# Delete databases
delete_stack "$CH_STACK_NAME"
delete_stack "$PG_STACK_NAME"

# Delete infrastructure last
delete_stack "$INFRA_STACK_NAME"

echo "✅ All resources deleted successfully!"
