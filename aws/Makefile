# Sutram AWS Infrastructure Makefile
# Manages CloudFormation stacks for the Sutram trading platform

.PHONY: help init validate deploy-all deploy-networking deploy-data deploy-app-infra deploy-app-services deploy-cdn
.PHONY: delete-all delete-networking delete-data delete-app-infra delete-app-services delete-cdn
.PHONY: status outputs clean validate-all update-services build-images push-images
.PHONY: deploy-dev deploy-staging deploy-prod delete-dev delete-staging delete-prod
.PHONY: show-costs scale-up scale-down stop-services start-services cost-optimization-report

# Default environment
ENV ?= dev

# Default region
AWS_REGION ?= us-east-1

# Source the deployment functions
SHELL := /bin/bash
export PATH := $(shell pwd)/scripts:$(PATH)

# Colors
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

# Help target
help: ## Show this help message
	@echo "$(BLUE)Sutram AWS Infrastructure Management$(NC)"
	@echo "$(BLUE)=====================================$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Environment Options:$(NC)"
	@echo "  ENV=dev      Deploy to development environment"
	@echo "  ENV=staging  Deploy to staging environment" 
	@echo "  ENV=prod     Deploy to production environment"
	@echo ""
	@echo "$(GREEN)Examples:$(NC)"
	@echo "  make deploy-all ENV=dev"
	@echo "  make delete-all ENV=staging"
	@echo "  make status ENV=prod"

init: ## Initialize AWS CLI and validate configuration
	@echo "$(BLUE)Initializing AWS configuration...$(NC)"
	@aws sts get-caller-identity
	@echo "$(GREEN)AWS CLI configured successfully$(NC)"

validate-all: ## Validate all CloudFormation templates
	@echo "$(BLUE)Validating all templates...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@for template in templates/*.yaml; do \
		echo "Validating $$template..."; \
		aws cloudformation validate-template --template-body file://$$template --region $(AWS_REGION) > /dev/null; \
		echo "$(GREEN)✓ $$template is valid$(NC)"; \
	done

# Individual stack deployment targets
deploy-networking: ## Deploy networking stack (VPC, subnets, security groups)
	@echo "$(BLUE)Deploying networking stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "01-networking"

deploy-data: ## Deploy data stack (RDS, ClickHouse, Redis, S3)
	@echo "$(BLUE)Deploying data stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "02-data"

deploy-app-infra: ## Deploy application infrastructure (ECS, ALB, ECR)
	@echo "$(BLUE)Deploying application infrastructure stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "03-app-infrastructure"

deploy-app-services: ## Deploy application services (ECS services, task definitions)
	@echo "$(BLUE)Deploying application services stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "04-app-services"

deploy-cdn: ## Deploy CDN stack (CloudFront, Route 53, WAF)
	@echo "$(BLUE)Deploying CDN stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "05-cdn"

# Stack deletion targets
delete-networking: ## Delete networking stack
	@echo "$(RED)Deleting networking stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && delete_stack "01-networking"

delete-data: ## Delete data stack
	@echo "$(RED)Deleting data stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && delete_stack "02-data"

delete-app-infra: ## Delete application infrastructure stack
	@echo "$(RED)Deleting application infrastructure stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && delete_stack "03-app-infrastructure"

delete-app-services: ## Delete application services stack
	@echo "$(RED)Deleting application services stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && delete_stack "04-app-services"

delete-cdn: ## Delete CDN stack
	@echo "$(RED)Deleting CDN stack for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && delete_stack "05-cdn"

# Full deployment targets
deploy-all: validate-all deploy-networking deploy-data deploy-app-infra deploy-app-services deploy-cdn ## Deploy all stacks in correct order
	@echo "$(GREEN)All stacks deployed successfully for $(ENV) environment!$(NC)"
	@make outputs ENV=$(ENV)

# Full deletion target (in reverse order)
delete-all: ## Delete all stacks in reverse order
	@echo "$(RED)Deleting all stacks for $(ENV) environment in reverse order...$(NC)"
	@echo "$(YELLOW)WARNING: This will delete ALL infrastructure for $(ENV) environment!$(NC)"
	@read -p "Are you sure? Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ]
	@make delete-cdn ENV=$(ENV) || true
	@make delete-app-services ENV=$(ENV) || true
	@make delete-app-infra ENV=$(ENV) || true
	@make delete-data ENV=$(ENV) || true
	@make delete-networking ENV=$(ENV) || true
	@echo "$(GREEN)All stacks deleted for $(ENV) environment$(NC)"

# Environment-specific shortcuts
deploy-dev: ## Deploy all stacks to development environment
	@make deploy-all ENV=dev

deploy-staging: ## Deploy all stacks to staging environment
	@make deploy-all ENV=staging

deploy-prod: ## Deploy all stacks to production environment
	@make deploy-all ENV=prod

delete-dev: ## Delete all stacks from development environment
	@make delete-all ENV=dev

delete-staging: ## Delete all stacks from staging environment
	@make delete-all ENV=staging

delete-prod: ## Delete all stacks from production environment
	@make delete-all ENV=prod

# Status and information targets
status: ## Show status of all stacks
	@echo "$(BLUE)Stack Status for $(ENV) environment:$(NC)"
	@echo "$(BLUE)================================$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && \
	for stack in "01-networking" "02-data" "03-app-infrastructure" "04-app-services" "05-cdn"; do \
		status=$$(get_stack_status "$$stack"); \
		if [ "$$status" = "NOT_EXISTS" ]; then \
			echo "$(RED)❌ $$stack: NOT_EXISTS$(NC)"; \
		elif [[ "$$status" == *"COMPLETE"* ]]; then \
			echo "$(GREEN)✅ $$stack: $$status$(NC)"; \
		elif [[ "$$status" == *"PROGRESS"* ]]; then \
			echo "$(YELLOW)🔄 $$stack: $$status$(NC)"; \
		else \
			echo "$(RED)❌ $$stack: $$status$(NC)"; \
		fi; \
	done

outputs: ## Show outputs from all deployed stacks
	@echo "$(BLUE)Stack Outputs for $(ENV) environment:$(NC)"
	@echo "$(BLUE)===================================$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && \
	for stack in "01-networking" "02-data" "03-app-infrastructure" "04-app-services" "05-cdn"; do \
		echo ""; \
		echo "$(GREEN)Stack: $$stack$(NC)"; \
		echo "$(GREEN)==================$(NC)"; \
		show_outputs "$$stack" || echo "$(YELLOW)No outputs or stack not found$(NC)"; \
	done

# Application deployment targets
build-images: ## Build Docker images for web app and live trading
	@echo "$(BLUE)Building Docker images...$(NC)"
	@cd .. && docker build -f Dockerfile.web -t sutram-web:$(ENV) .
	@cd .. && docker build -f cmd/live/Dockerfile -t sutram-live:$(ENV) . || echo "$(YELLOW)Live service Dockerfile not found, skipping...$(NC)"
	@echo "$(GREEN)Docker images built successfully$(NC)"

push-images: ## Push Docker images to ECR
	@echo "$(BLUE)Pushing Docker images to ECR...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$(AWS_REGION).amazonaws.com
	@docker tag sutram-web:$(ENV) $$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$(AWS_REGION).amazonaws.com/sutram/$(ENV)/web-app:latest
	@docker push $$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$(AWS_REGION).amazonaws.com/sutram/$(ENV)/web-app:latest
	@echo "$(GREEN)Docker images pushed successfully$(NC)"

update-services: ## Update ECS services with new image versions
	@echo "$(BLUE)Updating ECS services for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV) && deploy_stack "04-app-services"

# Utility targets
clean: ## Clean up temporary files
	@echo "$(BLUE)Cleaning up temporary files...$(NC)"
	@rm -f *.tmp
	@echo "$(GREEN)Cleanup completed$(NC)"

migrate-data: ## Upload CSV data files to S3
	@echo "$(BLUE)Uploading instrument data to S3...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws s3 sync ../data/ s3://sutram-$(ENV)-instrument-data-$$(aws sts get-caller-identity --query Account --output text)/data/ --exclude "*" --include "*.csv"
	@echo "$(GREEN)Data upload completed$(NC)"

setup-secrets: ## Guide for setting up secrets in AWS Secrets Manager
	@echo "$(BLUE)Secrets Setup Guide for $(ENV) environment:$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@echo ""
	@echo "$(YELLOW)After deployment, you need to update the following secrets in AWS Secrets Manager:$(NC)"
	@echo ""
	@echo "$(GREEN)1. Database Secret:$(NC)"
	@echo "   aws secretsmanager update-secret --secret-id sutram/$(ENV)/database/credentials --secret-string '{\"username\":\"admin\",\"password\":\"YOUR_DB_PASSWORD\"}'"
	@echo ""
	@echo "$(GREEN)2. Broker API Secret:$(NC)"
	@echo "   aws secretsmanager update-secret --secret-id sutram/$(ENV)/broker/api-keys --secret-string '{\"dhan_access_token\":\"YOUR_TOKEN\",\"zerodha_api_key\":\"YOUR_KEY\",...}'"
	@echo ""
	@echo "$(GREEN)3. Application Secret:$(NC)"
	@echo "   aws secretsmanager update-secret --secret-id sutram/$(ENV)/application/secrets --secret-string '{\"jwt_secret\":\"YOUR_JWT_SECRET\",\"encryption_key\":\"YOUR_KEY\",...}'"
	@echo ""
	@echo "$(YELLOW)Use the AWS Console or CLI to update these secrets with your actual values.$(NC)"

# Quick deployment target that does everything
quick-deploy: validate-all deploy-all migrate-data ## Quick deployment: validate, deploy all stacks, and upload data
	@echo "$(GREEN)Quick deployment completed for $(ENV) environment!$(NC)"
	@make setup-secrets ENV=$(ENV)

# Rollback target
rollback-services: ## Rollback application services to previous version
	@echo "$(YELLOW)Rolling back application services for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecs update-service --cluster sutram-$(ENV)-cluster --service sutram-$(ENV)-web-app --force-new-deployment --region $(AWS_REGION)
	@aws ecs update-service --cluster sutram-$(ENV)-cluster --service sutram-$(ENV)-live-trading --force-new-deployment --region $(AWS_REGION)
	@echo "$(GREEN)Rollback initiated$(NC)"

# Cost monitoring and optimization targets
show-costs: ## Show current AWS costs for the environment
	@echo "$(BLUE)Fetching current costs for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ce get-cost-and-usage \
		--time-period Start=`date -d '30 days ago' +%Y-%m-%d`,End=`date +%Y-%m-%d` \
		--granularity MONTHLY \
		--metrics BlendedCost \
		--group-by Type=DIMENSION,Key=SERVICE \
		--filter file://<(echo '{"Tags":{"Key":"Environment","Values":["$(ENV)"]}}') \
		--region $(AWS_REGION) \
		--query 'ResultsByTime[0].Groups[*].[Keys[0],Metrics.BlendedCost.Amount]' \
		--output table

scale-up: ## Scale up resources for testing (increases costs)
	@echo "$(YELLOW)Scaling up $(ENV) environment for testing...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-web-app \
		--desired-count 2 --region $(AWS_REGION)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-live-trading \
		--desired-count 1 --region $(AWS_REGION)

scale-down: ## Scale down resources to minimum (reduces costs)
	@echo "$(GREEN)Scaling down $(ENV) environment to minimum...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-web-app \
		--desired-count 1 --region $(AWS_REGION)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-live-trading \
		--desired-count 1 --region $(AWS_REGION)

stop-services: ## Stop ECS services to minimize costs during off-hours
	@echo "$(RED)Stopping ECS services for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-web-app \
		--desired-count 0 --region $(AWS_REGION)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-live-trading \
		--desired-count 0 --region $(AWS_REGION)

start-services: ## Start ECS services after off-hours
	@echo "$(GREEN)Starting ECS services for $(ENV) environment...$(NC)"
	@source scripts/deploy.sh && load_config $(ENV)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-web-app \
		--desired-count 1 --region $(AWS_REGION)
	@aws ecs update-service --cluster $(PROJECT_NAME)-$(ENV)-cluster \
		--service $(PROJECT_NAME)-$(ENV)-live-trading \
		--desired-count 1 --region $(AWS_REGION)

cost-optimization-report: ## Generate cost optimization report
	@echo "$(BLUE)Generating cost optimization report...$(NC)"
	@echo "See COST_OPTIMIZATION.md for detailed cost breakdown and savings strategies"
	@echo "Current environment: $(ENV)"
	@echo "Estimated monthly cost: $(if $(filter dev,$(ENV)),~$$78 (71% savings),$(if $(filter staging,$(ENV)),~$$150,~$$300))"

# Default target
.DEFAULT_GOAL := help
