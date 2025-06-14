# Sutram AWS Infrastructure

This directory contains the AWS infrastructure configuration for the Sutram Trading Platform using CloudFormation templates and deployment automation.

## 💰 Cost Optimization

**🎯 Development Environment: ~$78/month (71% savings from production)**

This infrastructure includes aggressive cost optimizations for rapid development:
- **FARGATE_SPOT**: 50-70% compute cost reduction
- **Single NAT Gateway**: $32/month savings
- **Minimal Resources**: t3.micro instances, minimal storage
- **Disabled Services**: CDN, WAF, monitoring for development

See [COST_OPTIMIZATION.md](./COST_OPTIMIZATION.md) for complete cost breakdown and optimization strategies.

## 🏗️ Architecture Overview

The infrastructure is organized into 5 separate CloudFormation stacks following best practices for maintainability, security, and scalability:

### Stack 1: Networking (`01-networking.yaml`)
- **VPC** with public/private/database subnets across 2 AZs
- **Internet Gateway** and **NAT Gateways** for internet access
- **Security Groups** for ALB, ECS, and databases
- **KMS Key** for encryption
- **Route Tables** and subnet associations

### Stack 2: Data Layer (`02-data.yaml`)
- **Amazon RDS PostgreSQL** (Multi-AZ in production)
- **ClickHouse on EC2** with EBS storage
- **ElastiCache Redis** for caching
- **S3 Buckets** for instrument data, app assets, and backups
- **Systems Manager Parameters** for configuration

### Stack 3: Application Infrastructure (`03-app-infrastructure.yaml`)
- **ECS Fargate Cluster** with auto-scaling
- **Application Load Balancer** with SSL termination
- **ECR Repositories** for Docker images
- **CloudWatch Log Groups** for application logs
- **IAM Roles** for ECS tasks

### Stack 4: Application Services (`04-app-services.yaml`)
- **ECS Task Definitions** for web app and live trading service
- **ECS Services** with auto-scaling policies
- **AWS Secrets Manager** for sensitive configuration
- **CloudWatch Alarms** for monitoring

### Stack 5: CDN & Domain (`05-cdn.yaml`)
- **CloudFront Distribution** for global content delivery
- **AWS WAF** for security protection
- **Route 53** for DNS management (optional)
- **ACM Certificates** for SSL/TLS

## 📁 Directory Structure

```
aws/
├── Makefile                    # Deployment automation
├── README.md                   # This file
├── config/                     # Environment configurations
│   ├── dev.env                # Development environment
│   ├── staging.env            # Staging environment
│   └── prod.env               # Production environment
├── scripts/                    # Deployment scripts
│   └── deploy.sh              # Main deployment functions
└── templates/                  # CloudFormation templates
    ├── 01-networking.yaml     # VPC and networking
    ├── 02-data.yaml          # Databases and storage
    ├── 03-app-infrastructure.yaml  # ECS and ALB
    ├── 04-app-services.yaml  # ECS services
    └── 05-cdn.yaml           # CloudFront and DNS
```

## 🚀 Quick Start

### Prerequisites

1. **AWS CLI** installed and configured
2. **Docker** installed (for building images)
3. **Make** utility available
4. **Bash** shell (macOS/Linux)

### Initial Setup

1. **Configure AWS Account ID** in environment files:
```bash
# Edit config files and add your AWS Account ID
vim config/dev.env
vim config/staging.env
vim config/prod.env
```

2. **Validate templates**:
```bash
make validate-all ENV=dev
```

3. **Deploy to development**:
```bash
make deploy-all ENV=dev
```

## 🎯 Common Commands

### Deployment

| Command | Description |
|---------|-------------|
| `make deploy-all ENV=dev` | Deploy all stacks to development |
| `make deploy-all ENV=staging` | Deploy all stacks to staging |
| `make deploy-all ENV=prod` | Deploy all stacks to production |
| `make quick-deploy ENV=dev` | Deploy everything + upload data |

### Individual Stack Management

| Command | Description |
|---------|-------------|
| `make deploy-networking ENV=dev` | Deploy only networking stack |
| `make deploy-data ENV=dev` | Deploy only data stack |
| `make deploy-app-infra ENV=dev` | Deploy only app infrastructure |
| `make deploy-app-services ENV=dev` | Deploy only app services |
| `make deploy-cdn ENV=dev` | Deploy only CDN stack |

### Monitoring & Information

| Command | Description |
|---------|-------------|
| `make status ENV=dev` | Show status of all stacks |
| `make outputs ENV=dev` | Show outputs from all stacks |
| `make setup-secrets ENV=dev` | Guide for configuring secrets |

### Application Deployment

| Command | Description |
|---------|-------------|
| `make build-images ENV=dev` | Build Docker images locally |
| `make push-images ENV=dev` | Push images to ECR |
| `make update-services ENV=dev` | Update ECS services |
| `make migrate-data ENV=dev` | Upload CSV data to S3 |

### Cleanup

| Command | Description |
|---------|-------------|
| `make delete-all ENV=dev` | Delete all stacks (DESTRUCTIVE) |
| `make delete-app-services ENV=dev` | Delete only app services |
| `make rollback-services ENV=dev` | Rollback to previous version |

## 🔧 Configuration

### Environment Files

Each environment has its own configuration file in the `config/` directory:

- **`dev.env`**: Development environment (minimal resources)
- **`staging.env`**: Staging environment (moderate resources)  
- **`prod.env`**: Production environment (high availability, larger instances)

### Key Configuration Options

```bash
# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCOUNT_ID=123456789012

# Domain Configuration (optional)
DOMAIN_NAME=sutram.com
CREATE_ROUTE53_HOSTED_ZONE=true

# Database Configuration
DATABASE_INSTANCE_CLASS=db.t3.micro  # dev
DATABASE_INSTANCE_CLASS=db.r5.large  # prod

# Application Configuration
WEB_APP_DESIRED_COUNT=2
LIVE_TRADING_DESIRED_COUNT=1
```

## 🔐 Security Configuration

### Secrets Management

After deployment, you need to configure secrets in AWS Secrets Manager:

1. **Database Credentials**:
```bash
aws secretsmanager update-secret \
  --secret-id sutram/dev/database/credentials \
  --secret-string '{"username":"admin","password":"your-secure-password"}'
```

2. **Broker API Keys**:
```bash
aws secretsmanager update-secret \
  --secret-id sutram/dev/broker/api-keys \
  --secret-string '{
    "dhan_access_token":"your-token",
    "zerodha_api_key":"your-key",
    "zerodha_access_token":"your-token"
  }'
```

3. **Application Secrets**:
```bash
aws secretsmanager update-secret \
  --secret-id sutram/dev/application/secrets \
  --secret-string '{
    "jwt_secret":"your-jwt-secret",
    "encryption_key":"your-encryption-key",
    "google_oauth_client_id":"your-client-id",
    "google_oauth_client_secret":"your-client-secret"
  }'
```

### Access Control

- **ECS Tasks** run with minimal IAM permissions
- **Databases** are in private subnets with security group restrictions
- **S3 Buckets** have encryption at rest and restricted access
- **KMS Keys** encrypt all sensitive data
- **WAF** protects against common web attacks

## 📊 Monitoring & Observability

### CloudWatch Integration

- **Application Logs**: Automatically sent to CloudWatch Logs
- **Metrics**: CPU, memory, and custom trading metrics
- **Alarms**: High CPU utilization alerts
- **Dashboards**: Available through CloudWatch console

### Key Metrics to Monitor

- **ECS Service Health**: Task count, CPU, memory usage
- **Database Performance**: Connection count, query performance
- **Trading Metrics**: Order execution time, data processing latency
- **Cost Optimization**: AWS Cost Explorer integration

## 🚨 Troubleshooting

### Common Issues

1. **Stack Deployment Fails**:
   - Check CloudFormation events in AWS Console
   - Verify IAM permissions
   - Check resource limits (EIPs, VPCs, etc.)

2. **Application Not Accessible**:
   - Verify security group rules
   - Check ECS service health
   - Confirm ALB target health

3. **Database Connection Issues**:
   - Verify security group allows ECS access
   - Check database credentials in Secrets Manager
   - Confirm database is in AVAILABLE state

4. **High Costs**:
   - Consider using FARGATE_SPOT for development
   - Right-size database instances
   - Implement S3 lifecycle policies

### Useful AWS CLI Commands

```bash
# Check stack status
aws cloudformation describe-stacks --stack-name sutram-dev-networking

# View ECS service status
aws ecs describe-services --cluster sutram-dev-cluster --services sutram-dev-web-app

# Check application logs
aws logs tail /ecs/sutram/dev/web-app --follow

# Monitor real-time metrics
aws cloudwatch get-metric-statistics --namespace AWS/ECS --metric-name CPUUtilization
```

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy to AWS
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      
      - name: Deploy infrastructure
        run: |
          cd aws
          make deploy-all ENV=staging
          
      - name: Build and push images
        run: |
          cd aws
          make build-images ENV=staging
          make push-images ENV=staging
          make update-services ENV=staging
```

## 💰 Cost Estimation

### Development Environment (~$150-300/month)
- **ECS Fargate**: ~$50-100
- **RDS t3.micro**: ~$15
- **ClickHouse t3.small**: ~$15
- **ElastiCache t3.micro**: ~$10
- **ALB**: ~$20
- **NAT Gateway**: ~$45
- **Data Transfer**: ~$10-50

### Production Environment (~$800-2000/month)
- **ECS Fargate**: ~$200-500
- **RDS r5.large Multi-AZ**: ~$200
- **ClickHouse c5n.xlarge**: ~$150
- **ElastiCache r5.large**: ~$100
- **ALB**: ~$20
- **NAT Gateway**: ~$90
- **CloudFront**: ~$50-100
- **Data Transfer**: ~$100-300

## 📈 Scaling Considerations

### Horizontal Scaling
- **ECS Auto Scaling**: Based on CPU/memory metrics
- **Database Read Replicas**: For read-heavy workloads
- **Multi-Region**: For global latency optimization

### Vertical Scaling
- **Instance Sizing**: Start small, scale based on metrics
- **Storage**: EBS GP3 with auto-scaling for ClickHouse
- **Memory**: Increase ECS task memory for complex strategies

## 🔒 Compliance & Security

### Data Protection
- **Encryption at Rest**: All databases and S3 buckets
- **Encryption in Transit**: TLS 1.2+ everywhere
- **Key Management**: AWS KMS with rotation
- **Backup Strategy**: Automated RDS backups, S3 versioning

### Access Control
- **IAM Principle of Least Privilege**: Minimal permissions
- **VPC Network Isolation**: Private subnets for databases
- **WAF Protection**: Against OWASP Top 10
- **Security Groups**: Whitelist-based access

## 📞 Support

For issues with this infrastructure:

1. Check the troubleshooting section above
2. Review CloudFormation events in AWS Console
3. Check the application logs in CloudWatch
4. Verify configuration in the environment files

## 🚀 Next Steps

After successful deployment:

1. **Configure Secrets**: Update all secrets in AWS Secrets Manager
2. **Upload Data**: Run `make migrate-data ENV=dev` to upload CSV files
3. **Set up Monitoring**: Create CloudWatch dashboards
4. **Configure CI/CD**: Set up automated deployments
5. **Test Trading**: Verify broker integrations work correctly
6. **Set up Alerts**: Configure SNS notifications for critical events
