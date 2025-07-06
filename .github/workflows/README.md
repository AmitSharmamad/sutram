# GitHub Workflows for Sutram Trading Platform

This directory contains GitHub Actions workflows for automated deployment and management of the Sutram Trading Platform. The workflows are designed to handle different aspects of the application lifecycle including testing, building, deployment, and infrastructure management.

## 🚀 Workflows Overview

### 1. **Main Deployment Workflow** (`deploy.yml`)
- **Triggers**: Push to `main` or `develop` branches
- **Purpose**: Build and deploy the main web application
- **Environments**: 
  - Development (from `develop` branch)
  - Staging (from `main` branch)
  - Production (from `main` branch, after staging)

### 2. **Live Trading Service Deployment** (`deploy-live.yml`)
- **Triggers**: Push to `main`/`develop` branches with changes to live service, or manual dispatch
- **Purpose**: Deploy the live trading WebSocket service
- **Features**: 
  - Separate Docker image for live service
  - Manual deployment option with environment selection
  - Health checks after deployment

### 3. **Infrastructure Deployment** (`deploy-infrastructure.yml`)
- **Triggers**: Push to `main`/`develop` branches with AWS changes, or manual dispatch
- **Purpose**: Deploy and manage AWS infrastructure using CloudFormation
- **Features**:
  - Template validation
  - Cost monitoring
  - Infrastructure destruction (manual)
  - Environment-specific deployments

### 4. **Database and Deployment Tasks** (`deploy-db-tasks.yml`)
- **Triggers**: Push with database changes, or manual dispatch
- **Purpose**: Handle database migrations and deployment tasks
- **Features**:
  - Database migrations
  - Backup and restore operations
  - Deployment tool execution

### 5. **PR Testing** (`pr-test.yml`)
- **Triggers**: Pull requests to `main` branch
- **Purpose**: Run tests and ensure code quality
- **Features**: 
  - Test coverage analysis
  - 100% coverage requirement for changed files

## 🏗️ Architecture

### Branching Strategy
- **`develop`**: Development environment deployments
- **`main`**: Staging and production deployments
- **Feature branches**: PR testing only

### Environment Flow
```
develop → Development Environment
main → Staging Environment → Production Environment
```

### Container Images
All services are built as Docker images and pushed to GitHub Container Registry:
- `ghcr.io/your-org/sutram:latest` - Web application
- `ghcr.io/your-org/sutram-live:latest` - Live trading service
- `ghcr.io/your-org/sutram-migrate:latest` - Migration tools
- `ghcr.io/your-org/sutram-deployment:latest` - Deployment tools

## 🔧 Setup Instructions

### 1. Repository Secrets

Add the following secrets to your GitHub repository:

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID       # AWS Access Key
AWS_SECRET_ACCESS_KEY   # AWS Secret Key

# Optional: Additional secrets for specific services
DHAN_ACCESS_TOKEN       # Trading API tokens
ZERODHA_API_KEY
ZERODHA_ACCESS_TOKEN
ANGELONE_API_KEY
```

### 2. Repository Variables

Add the following variables to your repository:

```bash
AWS_REGION              # Default: us-east-1
```

### 3. Environment Protection Rules

Set up environment protection rules in your repository settings:

#### Development Environment
- No protection rules needed
- Auto-deploy from `develop` branch

#### Staging Environment
- Required reviewers: 1+
- Auto-deploy from `main` branch after tests pass

#### Production Environment
- Required reviewers: 2+
- Wait timer: 5 minutes
- Deploy only after staging succeeds

### 4. AWS Account Setup

Ensure your AWS account has the following:

1. **IAM Role/User** with appropriate permissions:
   - ECS Full Access
   - CloudFormation Full Access
   - RDS Full Access
   - EC2 Full Access
   - S3 Full Access
   - Systems Manager Parameter Store Access

2. **AWS Account ID** configured in environment files:
   ```bash
   # Update aws/config/dev.env
   AWS_ACCOUNT_ID=123456789012
   ```

## 🚦 Usage Guide

### Automatic Deployments

1. **Development**: Push to `develop` branch
2. **Staging**: Push to `main` branch
3. **Production**: Automatically follows staging (with approval)

### Manual Deployments

#### Deploy Specific Service
```bash
# Go to Actions tab → Deploy Live Trading Service → Run workflow
# Select environment and trigger
```

#### Deploy Infrastructure
```bash
# Go to Actions tab → Deploy Infrastructure → Run workflow
# Select environment and action (deploy/destroy/validate)
```

#### Database Operations
```bash
# Go to Actions tab → Database and Deployment Tasks → Run workflow
# Select task (migrate/deployment/backup/restore) and environment
```

### Local Development

For local development and testing:

```bash
# Test the build locally
docker build -f Dockerfile.web -t sutram-web .
docker build -f Dockerfile.live -t sutram-live .

# Run infrastructure validation
cd aws
make validate-all ENV=dev

# Run local tests
go test ./...
```

## 📊 Monitoring and Troubleshooting

### Cost Monitoring
- Infrastructure workflow includes cost monitoring
- Check AWS Cost Explorer for detailed analysis
- Development environment optimized for ~$78/month

### Deployment Status
- Check Actions tab for workflow status
- AWS ECS console for service health
- CloudWatch logs for application logs

### Common Issues

1. **ECS Task Definition Not Found**
   - Run infrastructure deployment first
   - Check if all CloudFormation stacks are deployed

2. **Image Not Found**
   - Ensure Docker images are built and pushed
   - Check GitHub Container Registry

3. **Database Connection Issues**
   - Verify RDS instance is running
   - Check security group rules
   - Verify environment variables/secrets

4. **High AWS Costs**
   - Run cost optimization report
   - Consider stopping dev environment during off-hours
   - Use spot instances for non-critical workloads

## 🔍 Workflow Details

### Build Process
1. **Code Checkout**: Latest code from branch
2. **Go Setup**: Install Go 1.23 and dependencies
3. **Wire Generation**: Generate dependency injection code
4. **Testing**: Run unit tests with coverage
5. **Docker Build**: Multi-platform image build
6. **Registry Push**: Push to GitHub Container Registry

### Deployment Process
1. **Image Pull**: Get latest image from registry
2. **Task Definition Update**: Update ECS task definition
3. **Service Update**: Rolling update of ECS service
4. **Health Check**: Verify service is running
5. **Notification**: Send deployment status

### Infrastructure Process
1. **Template Validation**: Validate CloudFormation templates
2. **Parameter Setup**: Configure environment-specific parameters
3. **Stack Deployment**: Deploy/update CloudFormation stacks
4. **Status Check**: Verify stack deployment
5. **Cost Monitoring**: Generate cost reports

## 📝 Best Practices

### Code Changes
- Always create feature branches
- Ensure tests pass before merging
- Update documentation for significant changes
- Use semantic commit messages

### Database Changes
- Always run migrations in order: dev → staging → prod
- Test migrations on development first
- Create backups before major changes
- Document migration steps

### Infrastructure Changes
- Validate templates before deployment
- Test in development environment first
- Monitor costs after changes
- Document infrastructure changes

### Security
- Never commit secrets to repository
- Use environment-specific configurations
- Rotate AWS credentials regularly
- Monitor access logs

## 🆘 Emergency Procedures

### Rollback Application
```bash
# Find previous working image
docker image ls ghcr.io/your-org/sutram

# Update ECS service manually
aws ecs update-service --cluster sutram-prod-cluster --service sutram-prod-web-app --task-definition sutram-prod-web-app:PREVIOUS_REVISION
```

### Stop All Services (Cost Saving)
```bash
# Run Infrastructure workflow with "destroy" action
# Or use AWS CLI
aws ecs update-service --cluster sutram-dev-cluster --service sutram-dev-web-app --desired-count 0
```

### Database Recovery
```bash
# List available snapshots
aws rds describe-db-snapshots --db-instance-identifier sutram-prod-postgres

# Restore from snapshot (requires manual intervention)
aws rds restore-db-instance-from-db-snapshot --db-instance-identifier sutram-prod-postgres-restored --db-snapshot-identifier sutram-prod-backup-TIMESTAMP
```

## 📞 Support

For issues with workflows:
1. Check GitHub Actions logs
2. Review AWS CloudFormation events
3. Check ECS service events
4. Review application logs in CloudWatch

For urgent production issues:
1. Use manual workflow dispatch for quick fixes
2. Monitor AWS console for service health
3. Check cost alerts for unexpected charges

---

*This documentation is maintained alongside the workflow files. Please update when making changes to the deployment process.*
