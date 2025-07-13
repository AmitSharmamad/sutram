# Sutram AWS Infrastructure Makefile

.PHONY: help deploy destroy status clean validate

# Default environment
ENV ?= dev

# Help target
help:
	@echo "Sutram AWS Infrastructure Management"
	@echo ""
	@echo "Usage:"
	@echo "  make deploy ENV=<env>       - Deploy all infrastructure"
	@echo "  make <service> ENV=<env> - Deploy specific service"
	@echo "  make destroy ENV=<env>      - Destroy all infrastructure"
	@echo "  make status ENV=<env>       - Show deployment status"
	@echo "  make validate               - Validate CloudFormation templates"
	@echo ""
	@echo "Environments:"
	@echo "  dev   - Development environment (default)"
	@echo "  prod  - Production environment"
	@echo ""
	@echo "Services:"
	@echo "  infra      - Infrastructure (VPC, networking)"
	@echo "  pg         - PostgreSQL database"
	@echo "  ch         - ClickHouse database"
	@echo "  web        - Web service"
	@echo "  live       - Live service"
	@echo "  deployment - Deployment service"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy                     # Deploy all services to dev"
	@echo "  make deploy ENV=prod            # Deploy all services to prod"
	@echo "  make web ENV=dev         # Deploy only web service to dev"
	@echo "  make destroy ENV=dev            # Destroy all dev services"
	@echo "  make status ENV=prod            # Show prod deployment status"

# Deploy all infrastructure
deploy:
	@echo "Deploying all services to $(ENV) environment..."
	@./scripts/deploy-all.sh $(ENV)

# Deploy individual services
infra:
	@./scripts/infra.sh $(ENV)

pg:
	@./scripts/pg.sh $(ENV)

ch:
	@./scripts/ch.sh $(ENV)

web:
	@./scripts/web.sh $(ENV)

live:
	@./scripts/live.sh $(ENV)

deployment:
	@./scripts/deployment.sh $(ENV)

# Destroy infrastructure
destroy:
	@echo "Destroying all services in $(ENV) environment..."
	@./scripts/destroy.sh $(ENV)

# Show deployment status
status:
	@echo "Checking deployment status for $(ENV) environment..."
	@aws cloudformation list-stacks --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE --query "StackSummaries[?contains(StackName, '$(ENV)')].{Name:StackName,Status:StackStatus}" --output table

# Validate CloudFormation templates
validate:
	@echo "Validating CloudFormation templates..."
	@for template in templates/*.yaml; do \
		echo "Validating $$template..."; \
		aws cloudformation validate-template --template-body file://$$template >/dev/null && \
		echo "✓ $$template is valid" || \
		echo "✗ $$template is invalid"; \
	done
# Quick shortcuts
dev:
	@make deploy ENV=dev

prod:
	@make deploy ENV=prod

destroy-dev:
	@make destroy ENV=dev

destroy-prod:
	@make destroy ENV=prod
