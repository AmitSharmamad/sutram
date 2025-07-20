.PHONY: help deploy-dev deploy-prod destroy-dev destroy-prod deploy-infra-dev deploy-infra-prod deploy-postgres-dev deploy-postgres-prod deploy-clickhouse-dev deploy-clickhouse-prod deploy-web-dev deploy-web-prod deploy-live-dev deploy-live-prod deploy-deployment-dev deploy-deployment-prod validate-setup setup-ecr

# Default environment and image tag
ENV ?= dev
TAG ?= latest

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Sutram AWS Deployment Commands"
	@echo "==============================="
	@echo ""
	@echo "Usage: make <target> [ENV=dev|prod] [TAG=image_tag]"
	@echo ""
	@echo "Main Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make deploy-dev                    # Deploy all services to dev"
	@echo "  make deploy-prod TAG=v1.2.3       # Deploy all services to prod with specific tag"
	@echo "  make deploy-web ENV=dev TAG=latest # Deploy only web service to dev"
	@echo "  make destroy-dev                   # Destroy all dev resources"

# Full deployment commands
deploy-dev: ## Deploy all services to dev environment
	@./scripts/deploy-all.sh dev $(TAG)

deploy-prod: ## Deploy all services to prod environment
	@./scripts/deploy-all.sh prod $(TAG)

# Individual service deployment commands
deploy-infra: ## Deploy infrastructure (VPC) to specified environment
	@./scripts/deploy-infra.sh $(ENV)

deploy-postgres: ## Deploy PostgreSQL to specified environment
	@./scripts/deploy-postgres.sh $(ENV)

deploy-clickhouse: ## Deploy ClickHouse to specified environment
	@./scripts/deploy-clickhouse.sh $(ENV)

deploy-web: ## Deploy web service to specified environment
	@./scripts/deploy-web.sh $(ENV) $(TAG)

deploy-live: ## Deploy live service to specified environment
	@./scripts/deploy-live.sh $(ENV) $(TAG)

deploy-deployment: ## Deploy deployment service to specified environment
	@./scripts/deploy-deployment.sh $(ENV) $(TAG)

# Quick access commands for dev environment
deploy-infra-dev: ## Deploy infrastructure to dev
	@./scripts/deploy-infra.sh dev

deploy-postgres-dev: ## Deploy PostgreSQL to dev
	@./scripts/deploy-postgres.sh dev

deploy-clickhouse-dev: ## Deploy ClickHouse to dev
	@./scripts/deploy-clickhouse.sh dev

deploy-web-dev: ## Deploy web service to dev
	@./scripts/deploy-web.sh dev $(TAG)

deploy-live-dev: ## Deploy live service to dev
	@./scripts/deploy-live.sh dev $(TAG)

deploy-deployment-dev: ## Deploy deployment service to dev
	@./scripts/deploy-deployment.sh dev $(TAG)

# Quick access commands for prod environment
deploy-infra-prod: ## Deploy infrastructure to prod
	@./scripts/deploy-infra.sh prod

deploy-postgres-prod: ## Deploy PostgreSQL to prod
	@./scripts/deploy-postgres.sh prod

deploy-clickhouse-prod: ## Deploy ClickHouse to prod
	@./scripts/deploy-clickhouse.sh prod

deploy-web-prod: ## Deploy web service to prod
	@./scripts/deploy-web.sh prod $(TAG)

deploy-live-prod: ## Deploy live service to prod
	@./scripts/deploy-live.sh prod $(TAG)

deploy-deployment-prod: ## Deploy deployment service to prod
	@./scripts/deploy-deployment.sh prod $(TAG)

# Destruction commands
destroy-dev: ## Destroy all dev environment resources
	@./scripts/destroy-all.sh dev

destroy-prod: ## Destroy all prod environment resources
	@./scripts/destroy-all.sh prod

# Utility commands
validate-setup: ## Validate the entire AWS deployment setup
	@./scripts/validate-setup.sh

setup-ecr: ## Setup ECR repositories for specified environment
	@./scripts/setup-ecr.sh $(ENV)

setup-ecr-dev: ## Setup ECR repositories for dev
	@./scripts/setup-ecr.sh dev

setup-ecr-prod: ## Setup ECR repositories for prod
	@./scripts/setup-ecr.sh prod

check-aws: ## Check if AWS CLI is configured
	@./scripts/common.sh && check_aws

validate-templates: ## Validate all CloudFormation templates
	@echo "Validating CloudFormation templates..."
	@for template in templates/*.yaml; do \
		echo "Validating $$template..."; \
		aws cloudformation validate-template --template-body file://$$template > /dev/null && echo "✓ $$template is valid" || echo "✗ $$template is invalid"; \
	done
