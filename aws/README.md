# Sutram AWS Deployment

This directory contains all the necessary files and scripts to deploy the Sutram application to AWS using CloudFormation and ECS.

## Architecture Overview

The deployment creates the following AWS resources:

### 1. Infrastructure (VPC)
- VPC with public and private subnets across 2 AZs
- Internet Gateway and NAT Gateway
- Route tables and security groups

### 2. PostgreSQL (RDS)
- PostgreSQL RDS instance in private subnets
- Security groups for database access
- Automated backups and encryption

### 3. ClickHouse (ECS + EFS)
- ClickHouse running on ECS Fargate
- EFS for persistent data storage
- CloudWatch logging

### 4. Application Services (ECS)
- **Web Service**: Main application with load balancer
- **Live Service**: Real-time data processing
- **Deployment Service**: Deployment management

## Prerequisites

1. **AWS CLI**: Install and configure AWS CLI with appropriate credentials
   ```bash
   aws configure
   ```

2. **Permissions**: Ensure your AWS user/role has permissions for:
   - CloudFormation
   - ECS, ECR
   - RDS, EFS
   - VPC, EC2
   - IAM (for service roles)
   - CloudWatch Logs

3. **Docker Images**: Build and push your application images to ECR before deploying services

## Quick Start

### Deploy Everything

For development environment:
```bash
make deploy-dev
```

For production environment:
```bash
make deploy-prod TAG=v1.0.0
```

### Deploy Individual Services

```bash
# Deploy infrastructure first
make deploy-infra ENV=dev

# Deploy databases
make deploy-postgres ENV=dev
make deploy-clickhouse ENV=dev

# Deploy application services
make deploy-web ENV=dev TAG=latest
make deploy-live ENV=dev TAG=latest
make deploy-deployment ENV=dev TAG=latest
```

## Configuration

Environment-specific configurations are stored in:
- `config/dev.env` - Development environment
- `config/prod.env` - Production environment

### Key Configuration Parameters

- **AWS_REGION**: AWS region for deployment
- **VPC_CIDR**: VPC CIDR block
- **DB_INSTANCE_CLASS**: RDS instance type
- **CPU/MEMORY**: ECS task resources
- **DB_PASSWORD**: Database password (change these!)

## Scripts

All deployment scripts are in the `scripts/` directory:

- `common.sh` - Common functions and utilities
- `deploy-infra.sh` - Deploy VPC and networking
- `deploy-postgres.sh` - Deploy PostgreSQL RDS
- `deploy-clickhouse.sh` - Deploy ClickHouse on ECS
- `deploy-web.sh` - Deploy web service
- `deploy-live.sh` - Deploy live service
- `deploy-deployment.sh` - Deploy deployment service
- `deploy-all.sh` - Deploy everything in order
- `destroy-all.sh` - Remove all resources

## Deployment Order

Services must be deployed in this order due to dependencies:

1. **Infrastructure** (VPC, subnets)
2. **PostgreSQL** (creates ECS security group)
3. **ClickHouse** (creates ECS cluster and IAM roles)
4. **Application Services** (web, live, deployment)

## Docker Images

Before deploying application services, ensure Docker images are available in ECR:

```bash
# Create ECR repositories (one-time setup)
aws ecr create-repository --repository-name sutram-dev-web
aws ecr create-repository --repository-name sutram-dev-live
aws ecr create-repository --repository-name sutram-dev-deployment

# Build and push images
docker build -t sutram-dev-web:latest -f docker/web/Dockerfile .
docker tag sutram-dev-web:latest 123456789012.dkr.ecr.us-east-1.amazonaws.com/sutram-dev-web:latest
docker push 123456789012.dkr.ecr.us-east-1.amazonaws.com/sutram-dev-web:latest
```

## Environment Differences

### Development (dev)
- Smaller instance sizes
- Single AZ for some services
- Shorter log retention
- No deletion protection

### Production (prod)
- Larger instance sizes
- Multi-AZ deployment
- Longer log retention
- Deletion protection enabled

## Monitoring and Logs

- **CloudWatch Logs**: All services log to CloudWatch
- **ECS Console**: Monitor service health and tasks
- **RDS Console**: Database metrics and performance

## Security

- All databases are in private subnets
- Security groups follow least privilege principle
- Encryption enabled for RDS and EFS
- No public access to databases

## Cleanup

To remove all resources:

```bash
make destroy-dev    # For dev environment
make destroy-prod   # For prod environment
```

**Warning**: This will permanently delete all data and resources!

## Troubleshooting

### Common Issues

1. **Stack deployment fails**: Check CloudFormation events in AWS Console
2. **Permission denied**: Ensure AWS credentials have required permissions
3. **Resource conflicts**: Make sure resource names are unique
4. **Image not found**: Ensure Docker images are pushed to ECR

### Useful Commands

```bash
# Check AWS configuration
aws sts get-caller-identity

# Validate CloudFormation templates
make validate-templates

# View stack events
aws cloudformation describe-stack-events --stack-name sutram-dev-web

# View ECS service logs
aws logs tail /ecs/sutram-dev-web --follow
```

## Cost Optimization

For development:
- Use t3.micro instances
- Scale down when not in use
- Use GP3 storage with lower provisioned IOPS

For production:
- Monitor costs with AWS Cost Explorer
- Use Reserved Instances for predictable workloads
- Set up billing alerts

## Support

If you encounter issues:
1. Check CloudFormation stack events
2. Review CloudWatch logs
3. Verify security group rules
4. Ensure all dependencies are deployed
