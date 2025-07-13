# AWS Deployment Scripts

Simple, focused deployment scripts for the Sutram infrastructure.

## Quick Start

Deploy everything to dev:
```bash
./scripts/deploy-all.sh dev
```

Deploy to production:
```bash
./scripts/deploy-all.sh prod
```

Destroy everything:
```bash
./scripts/destroy.sh dev
```

## Individual Services

Deploy individual services:
```bash
./scripts/infra.sh dev        # VPC, networking
./scripts/pg.sh dev           # PostgreSQL database
./scripts/ch.sh dev           # ClickHouse database
./scripts/web.sh dev          # Web service
./scripts/live.sh dev         # Live service
./scripts/deployment.sh dev   # Deployment service
```

## Using Makefile

```bash
make deploy ENV=dev           # Deploy all to dev
make deploy-web ENV=prod      # Deploy web service to prod
make destroy ENV=dev          # Destroy all dev resources
make status ENV=prod          # Check deployment status
make validate                 # Validate CloudFormation templates
```

## Scripts Overview

- `infra.sh` - Deploys VPC, subnets, security groups
- `pg.sh` - Deploys PostgreSQL RDS instance  
- `ch.sh` - Deploys ClickHouse on ECS
- `web.sh` - Deploys web service on ECS
- `live.sh` - Deploys live service on ECS
- `deployment.sh` - Deploys deployment service on ECS
- `deploy-all.sh` - Deploys everything in correct order
- `destroy.sh` - Destroys all resources

Each script is simple, focused, and uses `aws cloudformation deploy` for deployment.
