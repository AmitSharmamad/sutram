#!/bin/bash

# Setup ECR repositories for Docker images

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

print_info "Setting up ECR repositories for $ENVIRONMENT environment"

# List of services that need ECR repositories
SERVICES=("web" "live" "deployment")

for service in "${SERVICES[@]}"; do
    REPO_NAME="${PROJECT_NAME}-${ENVIRONMENT}-${service}"
    
    # Check if repository already exists
    if aws ecr describe-repositories --repository-names "$REPO_NAME" &> /dev/null; then
        print_info "ECR repository $REPO_NAME already exists"
    else
        print_info "Creating ECR repository: $REPO_NAME"
        aws ecr create-repository --repository-name "$REPO_NAME" --region "$AWS_REGION"
        
        if [[ $? -eq 0 ]]; then
            print_success "Created ECR repository: $REPO_NAME"
        else
            print_error "Failed to create ECR repository: $REPO_NAME"
            exit 1
        fi
    fi
done

# Get ECR login command
print_info "ECR repositories setup complete!"
print_info "To push images, first login to ECR:"
echo ""
echo "aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com"
echo ""
print_info "Then build and push your images:"
for service in "${SERVICES[@]}"; do
    REPO_NAME="${PROJECT_NAME}-${ENVIRONMENT}-${service}"
    echo "docker build -t $REPO_NAME:latest -f docker/$service/Dockerfile ."
    echo "docker tag $REPO_NAME:latest \$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com/$REPO_NAME:latest"
    echo "docker push \$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com/$REPO_NAME:latest"
    echo ""
done
