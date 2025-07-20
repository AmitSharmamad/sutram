#!/bin/bash

# Validate AWS deployment setup

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/common.sh"

print_info "Validating AWS deployment setup..."

# Check if AWS CLI is installed and configured
print_info "Checking AWS CLI configuration..."
check_aws

# Validate CloudFormation templates
print_info "Validating CloudFormation templates..."
TEMPLATES_DIR="$SCRIPT_DIR/../templates"
TEMPLATE_FILES=("infra.yaml" "postgres.yaml" "clickhouse.yaml" "web.yaml" "live.yaml" "deployment.yaml")

for template in "${TEMPLATE_FILES[@]}"; do
    template_path="$TEMPLATES_DIR/$template"
    if [[ -f "$template_path" ]]; then
        aws cloudformation validate-template --template-body "file://$template_path" > /dev/null
        if [[ $? -eq 0 ]]; then
            print_success "✓ $template is valid"
        else
            print_error "✗ $template is invalid"
            exit 1
        fi
    else
        print_error "Template file $template not found"
        exit 1
    fi
done

# Check environment files
print_info "Checking environment configuration files..."
ENV_FILES=("dev.env" "prod.env")
CONFIG_DIR="$SCRIPT_DIR/../config"

for env_file in "${ENV_FILES[@]}"; do
    env_path="$CONFIG_DIR/$env_file"
    if [[ -f "$env_path" ]]; then
        print_success "✓ $env_file exists"
        
        # Check for critical variables
        if grep -q "DB_PASSWORD=change_me" "$env_path"; then
            print_warning "⚠ Default password detected in $env_file - please change it!"
        fi
    else
        print_error "Environment file $env_file not found"
        exit 1
    fi
done

# Check script permissions
print_info "Checking script permissions..."
SCRIPT_FILES=(
    "common.sh"
    "deploy-infra.sh"
    "deploy-postgres.sh"
    "deploy-clickhouse.sh"
    "deploy-web.sh"
    "deploy-live.sh"
    "deploy-deployment.sh"
    "deploy-all.sh"
    "destroy-all.sh"
    "setup-ecr.sh"
)

for script in "${SCRIPT_FILES[@]}"; do
    script_path="$SCRIPT_DIR/$script"
    if [[ -f "$script_path" && -x "$script_path" ]]; then
        print_success "✓ $script is executable"
    elif [[ -f "$script_path" ]]; then
        print_error "✗ $script is not executable"
        exit 1
    else
        print_error "Script file $script not found"
        exit 1
    fi
done

echo ""
print_success "All validations passed!"
echo ""
print_info "You can now deploy using:"
print_info "  make setup-ecr-dev     # Setup ECR repositories for dev"
print_info "  make deploy-dev        # Deploy all services to dev"
print_info ""
print_info "Don't forget to:"
print_info "  1. Change default passwords in config files"
print_info "  2. Build and push Docker images to ECR"
print_info "  3. Review the README.md for detailed instructions"
