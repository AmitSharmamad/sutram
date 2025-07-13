# Sutram AWS Infrastructure Makefile

.PHONY: help deploy destroy status clean validate

# Default environment
ENV ?= dev

# Help target
help:
	@echo "Sutram AWS Infrastructure Management"
	@echo ""
	@echo "Usage:"
	@echo "  make deploy ENV=<env> SERVICE=<service>  - Deploy infrastructure"
	@echo "  make destroy ENV=<env> SERVICE=<service> - Destroy infrastructure"
	@echo "  make status ENV=<env>                    - Show deployment status"
	@echo "  make validate                            - Validate CloudFormation templates"
	@echo "  make clean                               - Clean temporary files"
	@echo ""
	@echo "Environments:"
	@echo "  dev   - Development environment (default)"
	@echo "  prod  - Production environment"
	@echo ""
	@echo "Services:"
	@echo "  all        - All services (default)"
	@echo "  infra      - Infrastructure (VPC, networking)"
	@echo "  pg         - PostgreSQL database"
	@echo "  ch         - ClickHouse database"
	@echo "  web        - Web service"
	@echo "  live       - Live service"
	@echo "  deployment - Deployment service"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy                              # Deploy all services to dev"
	@echo "  make deploy ENV=prod                     # Deploy all services to prod"
	@echo "  make deploy ENV=dev SERVICE=web          # Deploy only web service to dev"
	@echo "  make destroy ENV=dev SERVICE=all         # Destroy all dev services"
	@echo "  make status ENV=prod                     # Show prod deployment status"

# Default service
SERVICE ?= all

# Deploy infrastructure
deploy:
	@echo "Deploying $(SERVICE) to $(ENV) environment..."
	@./scripts/deploy.sh $(ENV) $(SERVICE)

# Destroy infrastructure
destroy:
	@echo "Destroying $(SERVICE) in $(ENV) environment..."
	@./scripts/destroy.sh $(ENV) $(SERVICE)

# Show deployment status
status:
	@echo "Checking deployment status for $(ENV) environment..."
	@./scripts/deploy.sh $(ENV) status

# Validate CloudFormation templates
validate:
	@echo "Validating CloudFormation templates..."
	@for template in templates/*.yaml; do \
		echo "Validating $$template..."; \
		aws cloudformation validate-template --template-body file://$$template >/dev/null && \
		echo "✓ $$template is valid" || \
		echo "✗ $$template is invalid"; \
	done

# Clean temporary files
clean:
	@echo "Cleaning temporary files..."
	@find . -name "*.tmp" -delete
	@find . -name ".DS_Store" -delete
	@echo "✓ Cleanup complete"

# Quick shortcuts
deploy-dev:
	@make deploy ENV=dev SERVICE=all

deploy-prod:
	@make deploy ENV=prod SERVICE=all

destroy-dev:
	@make destroy ENV=dev SERVICE=all

destroy-prod:
	@make destroy ENV=prod SERVICE=all

# Infrastructure only
infra-dev:
	@make deploy ENV=dev SERVICE=infra

infra-prod:
	@make deploy ENV=prod SERVICE=infra

# Databases only
db-dev:
	@make deploy ENV=dev SERVICE=pg
	@make deploy ENV=dev SERVICE=ch

db-prod:
	@make deploy ENV=prod SERVICE=pg
	@make deploy ENV=prod SERVICE=ch

# Services only (requires databases to be deployed first)
services-dev:
	@make deploy ENV=dev SERVICE=web
	@make deploy ENV=dev SERVICE=live
	@make deploy ENV=dev SERVICE=deployment

services-prod:
	@make deploy ENV=prod SERVICE=web
	@make deploy ENV=prod SERVICE=live
	@make deploy ENV=prod SERVICE=deployment
