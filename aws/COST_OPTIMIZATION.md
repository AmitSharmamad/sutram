# AWS Cost Optimization Guide for Sutram Trading Platform

## Overview

This document outlines the aggressive cost optimization strategies implemented for the Sutram Trading Platform's AWS infrastructure during the rapid development phase. All optimizations maintain full functionality while minimizing costs.

## Cost Optimization Summary

### Monthly Cost Breakdown (US East 1)

#### 🔴 Original Architecture (Production Settings)
| Service | Configuration | Monthly Cost |
|---------|---------------|--------------|
| **Networking** |
| NAT Gateway (2x) | 2 NAT Gateways + Data Transfer | ~$64.00 |
| **Compute** |
| ECS Fargate (Web App) | 2 tasks × 1024 CPU × 2048 MB | ~$45.00 |
| ECS Fargate (Live Trading) | 1 task × 2048 CPU × 4096 MB | ~$45.00 |
| Container Insights | Enabled | ~$5.00 |
| **Data Layer** |
| RDS PostgreSQL | db.t3.small, Multi-AZ, 20GB, Backups | ~$35.00 |
| ElastiCache Redis | cache.t3.small | ~$15.00 |
| ClickHouse EC2 | t3.small, 50GB root + 100GB data, encrypted | ~$25.00 |
| S3 Storage | Multiple buckets with KMS encryption | ~$10.00 |
| **CDN & Security** |
| CloudFront | Standard distribution | ~$5.00 |
| Route 53 | Hosted zone + health checks | ~$5.00 |
| WAF | Basic rules | ~$10.00 |
| KMS | Multiple keys | ~$5.00 |
| **TOTAL ORIGINAL** | | **~$269.00/month** |

#### 🟢 Optimized Architecture (Development Settings)
| Service | Configuration | Monthly Cost |
|---------|---------------|--------------|
| **Networking** |
| NAT Gateway (1x) | Single NAT Gateway shared | ~$32.00 |
| **Compute** |
| ECS Fargate SPOT (Web App) | 1 task × 256 CPU × 512 MB (SPOT) | ~$7.00 |
| ECS Fargate SPOT (Live Trading) | 1 task × 256 CPU × 512 MB (SPOT) | ~$7.00 |
| Container Insights | Disabled | $0.00 |
| **Data Layer** |
| RDS PostgreSQL | db.t3.micro, Single-AZ, 10GB, No backups | ~$12.00 |
| ElastiCache Redis | cache.t3.micro, No snapshots | ~$8.00 |
| ClickHouse EC2 | t3.micro, 20GB root + 20GB data, unencrypted | ~$8.00 |
| S3 Storage | Single bucket with AES256 encryption | ~$3.00 |
| **CDN & Security** |
| CloudFront | Disabled for development | $0.00 |
| Route 53 | Disabled for development | $0.00 |
| WAF | Disabled for development | $0.00 |
| KMS | Minimal usage | ~$1.00 |
| **TOTAL OPTIMIZED** | | **~$78.00/month** |

### 💰 **Total Monthly Savings: ~$191.00 (71% reduction)**

## Optimization Strategies Applied

### 1. Networking Layer (`01-networking.yaml`)
```yaml
# Cost Optimizations Applied:
- Single NAT Gateway shared across AZs (saves ~$32/month)
- Minimal public IP allocations
- Optimized security group rules
```

**Savings**: ~$32/month

### 2. Data Layer (`02-data.yaml`)
```yaml
# PostgreSQL Optimizations:
- Instance: db.t3.micro (instead of db.t3.small)
- Single-AZ deployment (no Multi-AZ)
- Storage: 10GB minimum (instead of 20GB)
- No automated backups in development
- No encryption for development
- No enhanced monitoring

# Redis Optimizations:
- Instance: cache.t3.micro (smallest available)
- No snapshot backups for development
- Simplified maintenance windows

# ClickHouse Optimizations:
- Instance: t3.micro (instead of t3.small)
- Storage: 20GB root + 20GB data (instead of 50GB + 100GB)
- No encryption for development
- Minimal CloudWatch monitoring

# S3 Optimizations:
- Single bucket for development
- AES256 encryption (instead of KMS)
- No backup bucket for development
```

**Savings**: ~$55/month

### 3. Application Infrastructure (`03-app-infrastructure.yaml`)
```yaml
# ECS Optimizations:
- FARGATE_SPOT as primary capacity provider (50-70% savings)
- Container Insights disabled for development
- Minimal CloudWatch logging
- Optimized task roles and policies
```

**Savings**: ~$88/month

### 4. Application Services (`04-app-services.yaml`)
```yaml
# Task Definition Optimizations:
- Web App: 256 CPU + 512 MB memory (instead of 1024 CPU + 2048 MB)
- Live Trading: 256 CPU + 512 MB memory (instead of 2048 CPU + 4096 MB)
- Single instance for each service (instead of 2+ replicas)
- Minimal health check configurations
```

**Savings**: ~$76/month

### 5. CDN Layer (`05-cdn.yaml`)
```yaml
# Development Simplifications:
- CloudFront disabled for development
- Route 53 disabled for development
- WAF disabled for development
- Direct ALB access for development testing
```

**Savings**: ~$20/month

## Environment-Specific Configuration

### Development Environment (`config/dev.env`)
```bash
# Ultra-aggressive cost optimization
USE_SINGLE_NAT_GATEWAY=true
USE_FARGATE_SPOT=true
DISABLE_MULTI_AZ=true
MINIMAL_STORAGE=true
DISABLE_ENHANCED_MONITORING=true
DISABLE_CONTAINER_INSIGHTS=true
MINIMAL_RESOURCES=true
DISABLE_CDN=true
DISABLE_WAF=true

# Minimal instance sizes
DATABASE_INSTANCE_CLASS=db.t3.micro
REDIS_NODE_TYPE=cache.t3.micro
CLICKHOUSE_INSTANCE_TYPE=t3.micro

# Single instance deployment
WEB_APP_DESIRED_COUNT=1
LIVE_TRADING_DESIRED_COUNT=1
```

### Staging Environment (`config/staging.env`)
```bash
# Moderate optimization for testing
USE_SINGLE_NAT_GATEWAY=true
USE_FARGATE_SPOT=true
DISABLE_MULTI_AZ=false
MINIMAL_STORAGE=false
DISABLE_ENHANCED_MONITORING=false
DISABLE_CONTAINER_INSIGHTS=false
MINIMAL_RESOURCES=false
DISABLE_CDN=false
DISABLE_WAF=true

# Small instance sizes
DATABASE_INSTANCE_CLASS=db.t3.small
REDIS_NODE_TYPE=cache.t3.small
CLICKHOUSE_INSTANCE_TYPE=t3.small

# Limited scaling
WEB_APP_DESIRED_COUNT=2
LIVE_TRADING_DESIRED_COUNT=1
```

### Production Environment (`config/prod.env`)
```bash
# Production settings (full features)
USE_SINGLE_NAT_GATEWAY=false
USE_FARGATE_SPOT=false
DISABLE_MULTI_AZ=false
MINIMAL_STORAGE=false
DISABLE_ENHANCED_MONITORING=false
DISABLE_CONTAINER_INSIGHTS=false
MINIMAL_RESOURCES=false
DISABLE_CDN=false
DISABLE_WAF=false

# Production instance sizes
DATABASE_INSTANCE_CLASS=db.r5.large
REDIS_NODE_TYPE=cache.r5.large
CLICKHOUSE_INSTANCE_TYPE=c5n.large

# Full scaling
WEB_APP_DESIRED_COUNT=3
LIVE_TRADING_DESIRED_COUNT=2
```

## Deployment Commands

### Quick Development Deployment
```bash
# Deploy optimized development environment
make deploy-dev

# Or individual stacks
make deploy-networking-dev
make deploy-data-dev
make deploy-app-infra-dev
make deploy-app-services-dev
```

### Cost Monitoring
```bash
# Check current costs
make show-costs-dev

# Validate optimizations
make validate-all-dev

# Show stack outputs
make outputs-dev
```

## Performance Impact Assessment

### ⚠️ Development Limitations
- **CPU/Memory**: Minimal allocation may cause slower response times
- **Storage**: Limited storage may require cleanup of logs/data
- **Networking**: Single NAT Gateway may be a bottleneck for high traffic
- **Monitoring**: Reduced observability for debugging

### ✅ Mitigation Strategies
- **Auto-scaling**: Configured to scale up when needed
- **Spot Instance Fallback**: Falls back to regular Fargate if spot unavailable
- **Quick Scaling**: Easy commands to scale up for testing
- **Environment Promotion**: Simple commands to deploy to staging/prod

## Additional Cost Optimization Tips

### 1. Resource Tagging
All resources are tagged for cost allocation:
```yaml
Tags:
  - Key: Project
    Value: sutram
  - Key: Environment
    Value: dev
  - Key: CostCenter
    Value: development
```

### 2. Automated Cleanup
```bash
# Schedule cleanup jobs
make cleanup-dev          # Clean unused resources
make stop-dev            # Stop services during off-hours
make start-dev           # Restart services
```

### 3. Monitoring and Alerts
```bash
# Set up cost alerts (one-time setup)
make setup-billing-alerts

# Check daily costs
make daily-cost-report
```

## Scaling Back Up

When ready to scale for production testing:

```bash
# Quick scale-up for load testing
make scale-up-dev

# Deploy staging environment
make deploy-staging

# Deploy production environment
make deploy-prod
```

## Cost Optimization Checklist

- [x] ✅ Single NAT Gateway for development
- [x] ✅ FARGATE_SPOT for 50-70% compute savings
- [x] ✅ Minimal CPU/Memory allocations
- [x] ✅ Single-AZ database deployment
- [x] ✅ Smallest instance types
- [x] ✅ Disabled monitoring and insights
- [x] ✅ Simplified storage configurations
- [x] ✅ Disabled CDN for development
- [x] ✅ Disabled WAF for development
- [x] ✅ AES256 instead of KMS encryption
- [x] ✅ Single service instances
- [x] ✅ Conditional resource creation
- [x] ✅ Environment-specific parameters

## Total Estimated Monthly Costs

| Environment | Monthly Cost | Use Case |
|-------------|--------------|----------|
| Development | **~$78** | Day-to-day development, feature testing |
| Staging | **~$150** | Integration testing, client demos |
| Production | **~$300** | Live trading, full monitoring, HA |

**Development savings: 71% reduction from production costs**

---

*Note: Costs are estimates based on us-east-1 pricing as of June 2025. Actual costs may vary based on usage patterns, data transfer, and AWS pricing changes.*
