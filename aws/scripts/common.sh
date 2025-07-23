#!/bin/bash

# Common functions and variables for deployment scripts

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to load environment variables
load_env() {
    local env_file="$1"
    if [[ ! -f "$env_file" ]]; then
        print_error "Environment file $env_file not found!"
        exit 1
    fi
    
    print_info "Loading environment from $env_file"
    set -a
    source "$env_file"
    set +a
}

# Function to check if AWS CLI is installed and configured
check_aws() {
    if ! command -v aws &> /dev/null; then
        print_error "AWS CLI not found. Please install AWS CLI."
        exit 1
    fi
    
    if ! aws sts get-caller-identity &> /dev/null; then
        print_error "AWS CLI not configured. Please run 'aws configure'."
        exit 1
    fi
    
    print_success "AWS CLI is configured"
}

# Function to check if CloudFormation stack exists
stack_exists() {
    local stack_name="$1"
    aws cloudformation describe-stacks --stack-name "$stack_name" &> /dev/null
}

# Function to wait for stack operation to complete
wait_for_stack() {
    local stack_name="$1"
    local operation="$2"
    
    print_info "Waiting for stack $stack_name to complete $operation..."
    
    aws cloudformation wait "stack-${operation}-complete" --stack-name "$stack_name"
    
    if [[ $? -eq 0 ]]; then
        print_success "Stack $stack_name $operation completed successfully"
    else
        print_error "Stack $stack_name $operation failed"
        exit 1
    fi
}

# Function to deploy CloudFormation stack
deploy_stack() {
    local stack_name="$1"
    local template_file="$2"
    shift 2
    local parameters=("$@")
    
    print_info "Deploying stack: $stack_name"
    
    local param_args=""
    for param in "${parameters[@]}"; do
        param_args+=" --parameter-overrides $param"
    done
    
    aws cloudformation deploy \
        --template-file "$template_file" \
        --stack-name "$stack_name" \
        --capabilities CAPABILITY_IAM \
        $param_args
    
    if [[ $? -eq 0 ]]; then
        print_success "Successfully deployed stack: $stack_name"
    else
        print_error "Failed to deploy stack: $stack_name"
        exit 1
    fi
}

# Function to delete CloudFormation stack
delete_stack() {
    local stack_name="$1"
    
    if stack_exists "$stack_name"; then
        print_warning "Deleting stack: $stack_name"
        aws cloudformation delete-stack --stack-name "$stack_name"
        wait_for_stack "$stack_name" "delete"
    else
        print_info "Stack $stack_name does not exist"
    fi
}

# Function to get stack output
get_stack_output() {
    local stack_name="$1"
    local output_key="$2"
    
    aws cloudformation describe-stacks \
        --stack-name "$stack_name" \
        --query "Stacks[0].Outputs[?OutputKey=='$output_key'].OutputValue" \
        --output text
}
