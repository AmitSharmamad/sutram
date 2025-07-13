#!/bin/bash

# Destroy Sutram Infrastructure
# Usage: ./destroy.sh [environment] [service]
# Example: ./destroy.sh dev all
# Example: ./destroy.sh prod web

set -e

# Default values
ENVIRONMENT=${1:-dev}
SERVICE=${2:-all}

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_DIR="$PROJECT_ROOT/config"

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

# Function to delete live service
delete_live() {
    log_info "Deleting live service..."
    delete_stack "$LIVE_STACK_NAME"
}

# Function to delete web service
delete_web() {
    log_info "Deleting web service..."
    delete_stack "$WEB_STACK_NAME"
}

# Function to delete ClickHouse
delete_ch() {
    log_info "Deleting ClickHouse database..."
    delete_stack "$CH_STACK_NAME"
}

# Function to delete PostgreSQL
delete_pg() {
    log_info "Deleting PostgreSQL database..."
    delete_stack "$PG_STACK_NAME"
}

# Function to delete infrastructure
delete_infra() {
    log_info "Deleting infrastructure (VPC, networking)..."
    delete_stack "$INFRA_STACK_NAME"
}

# Confirmation function
confirm_deletion() {
    echo -e "${RED}WARNING: This will permanently delete AWS resources!${NC}"
    echo "Environment: $ENVIRONMENT"
    echo "Service(s): $SERVICE"
    echo
    read -p "Are you sure you want to continue? (type 'yes' to confirm): " confirmation
    
    if [ "$confirmation" != "yes" ]; then
        log_info "Deletion cancelled."
        exit 0
    fi
}

# Main deletion logic
case $SERVICE in
    "all")
        confirm_deletion
        log_info "Deleting all services for environment: $ENVIRONMENT"
        log_warning "Deleting in reverse order to respect dependencies..."
        delete_deployment
        delete_live
        delete_web
        delete_ch
        delete_pg
        delete_infra
        log_success "All services deleted successfully!"
        ;;
    "infra")
        confirm_deletion
        delete_infra
        ;;
    "pg")
        confirm_deletion
        delete_pg
        ;;
    "ch")
        confirm_deletion
        delete_ch
        ;;
    "web")
        confirm_deletion
        delete_web
        ;;
    "live")
        confirm_deletion
        delete_live
        ;;
    "deployment")
        confirm_deletion
        delete_deployment
        ;;
    *)
        log_error "Unknown service: $SERVICE"
        echo "Usage: $0 [environment] [service]"
        echo "Environment: dev, prod"
        echo "Service: all, infra, pg, ch, web, live, deployment"
        exit 1
        ;;
esac
