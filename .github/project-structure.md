# Definitive Algorithmic Trading Platform Structure in Golang

This document outlines a scalable and deterministic architecture for an algorithmic trading platform built with Go, capable of handling dozens of modules.

## Application Types

### Web Application (REST APIs)
- Entry point in `cmd/web/main.go`
- Handles HTTP requests and serves REST APIs
- Responsible for request validation, routing, and response formatting
- Should not contain business logic; delegates to application services

### Worker Application (Async Task Execution)
- Entry point in `cmd/worker/main.go`
- Processes background jobs and asynchronous tasks
- Should use a reliable message queue system (e.g., Kafka, RabbitMQ)
- Each worker should have a single responsibility

### Scheduler Application (Cron Jobs)
- Entry point in `cmd/scheduler/main.go`
- Responsible for scheduling and triggering time-based jobs
- Should use a robust scheduler library or framework
- Jobs should be defined declaratively with clear execution patterns

## Project Root Structure

```
/
├── cmd/                    # Entry points
├── internal/               # Private application code
│   ├── domain/             # Core domain models and business rules
│   ├── modules/            # Business modules
│   ├── infra/              # Infrastructure components
│   ├── app/                # Application services
│   └── common/             # Shared utilities
├── pkg/                    # Public libraries that can be used by external applications
├── api/                    # API definitions and documentation
├── configs/                # Configuration files
├── scripts/                # Automation scripts
├── deployments/            # Deployment configurations
├── docs/                   # Documentation files
└── test/                   # Test files and fixtures
```

## Core Domain Layer

### `internal/domain/`
- `entity/`          # Pure domain entities with business rules
  - `order.go`       # Order entity with validation logic
  - `trade.go`       # Trade entity
  - `position.go`    # Position entity
  - `instrument.go`  # Financial instrument entity
  - `user.go`        # User entity
  - `strategy.go`    # Strategy entity
- `valueobject/`     # Immutable value objects
  - `money.go`       # Money value object
  - `timerange.go`   # Time range value object
  - `tradingstatus.go` # Trading status enumeration
- `event/`           # Domain events
  - `ordercreated.go`
  - `tradeexecuted.go`
- `repository/`      # Repository interfaces
  - `order_repo.go`   # Order table repository
  - `trade_repo.go`   # Trade table repository
  - `user_repo.go`    # User table repository
  - `role_repo.go`    # Role table repository
  - `user_role_repo.go` # User-Role many-to-many table repository
  - `strategy_repo.go` # Strategy table repository
- `service/`         # Domain service interfaces
  - `portfolio_service.go`   # Portfolio-related business logic
  - `risk_service.go`        # Risk management business logic
  - `execution_service.go`   # Order execution business logic
  - `user_service.go`        # User management business logic
  - `role_service.go`        # Role management business logic
  - `user_role_service.go`   # User-Role relationship business logic
- `factory/`         # Factory interfaces for entity creation
  - `orderfactory.go`
  - `strategyfactory.go`
- `errors/`          # Domain-specific error types
  - `tradeerrors.go`
  - `ordererrors.go`

## Business Modules (Self-contained modules with their own layers)

### `internal/modules/`

```
modules/
├── user/                      # User module
│   ├── application/           # Application layer
│   │   ├── dto/               # Data Transfer Objects
│   │   │   ├── userdto.go
│   │   │   └── authdto.go
│   │   ├── service/           # Application services
│   │   │   ├── user_service.go     # User-related business logic
│   │   │   ├── role_service.go     # Role-related business logic
│   │   │   └── user_role_service.go # User-role relationship business logic
│   │   ├── query/             # Query handlers (CQRS pattern)
│   │   │   └── userqueries.go
│   │   └── command/           # Command handlers (CQRS pattern)
│   │       └── usercommands.go
│   ├── domain/                # Module-specific domain logic
│   │   ├── model/             # Domain models
│   │   │   └── user.go
│   │   └── repository/        # Repository interfaces
│   │       ├── user_repository.go # User table repository
│   │       ├── role_repository.go # Role table repository
│   │       └── user_role_repository.go # User-Role join table repository
│   ├── infrastructure/        # Infrastructure implementations
│   │   ├── persistence/       # Repository implementations
│   │   │   ├── user_repo.go    # User table repository implementation
│   │   │   ├── role_repo.go    # Role table repository implementation
│   │   │   └── user_role_repo.go # User-Role join table repository implementation
│   │   └── auth/              # Authentication providers
│   │       ├── oauth.go
│   │       └── jwt.go
│   ├── interfaces/            # Interface adapters
│   │   ├── http/              # HTTP handlers
│   │   │   ├── controller.go  # Route handlers
│   │   │   ├── dto.go         # Request/Response objects
│   │   │   └── router.go      # Route definitions
│   │   └── grpc/              # gRPC handlers
│   │       └── server.go
│   ├── di.go                  # Dependency injection for module
│   └── module.go              # Module registration
```

Apply the same structure for each business module:

- `strategy/`     # Trading strategy module
- `backtesting/`  # Backtesting module
- `stock/`        # Stock data module
- `parser/`       # Data parsing module
- `portfolio/`    # Portfolio management
- `execution/`    # Order execution
- `risk/`         # Risk management
- `analytics/`    # Analytics and reporting
- `notification/` # Alerts and notifications
- `marketdata/`   # Market data processing

## Infrastructure Layer

### `internal/infra/`

```
infra/
├── database/                 # Database implementations
│   ├── postgres/             # PostgreSQL adapter
│   │   ├── connection.go     # DB connection handling
│   │   ├── migration.go      # DB migrations
│   │   └── transaction.go    # Transaction management
│   ├── clickhouse/           # ClickHouse adapter
│   │   └── connection.go
│   └── cache/                # Cache implementations
│       ├── redis.go          # Redis adapter
│       └── memory.go         # In-memory cache
├── broker/                   # Broker connections
│   ├── zerodha/              # Zerodha broker integration
│   │   ├── client.go         # API client
│   │   ├── mapper.go         # Data mappers
│   │   └── config.go         # Configuration
│   └── dhan/                 # Dhan broker integration
│       ├── client.go
│       ├── mapper.go
│       └── config.go
├── messagebus/               # Messaging infrastructure
│   ├── kafka/                # Kafka implementation
│   │   └── producer.go
│   └── inmemory/             # In-memory implementation
│       └── bus.go
├── storage/                  # Storage services
│   ├── s3.go                 # S3 compatible storage
│   └── local.go              # Local file storage
├── monitoring/               # Monitoring infrastructure
│   ├── prometheus.go         # Prometheus metrics
│   └── logging/              # Logging infrastructure
│       ├── zap.go            # Zap logger implementation
│       └── logrus.go         # Logrus logger implementation
└── http/                     # HTTP client infrastructure
    ├── client.go             # HTTP client with retry, circuit breaking
    └── middleware/           # HTTP middlewares
        ├── auth.go           # Authentication middleware
        ├── logging.go        # Logging middleware
        └── metrics.go        # Metrics collection middleware
```

## Application Layer

### `internal/app/`

```
app/
├── web/                     # Web application specific components
│   ├── middleware/          # HTTP middleware components
│   │   ├── auth.go          # Authentication middleware
│   │   ├── logging.go       # Request logging middleware
│   │   └── ratelimit.go     # Rate limiting middleware
│   ├── handler/             # Common request handlers
│   │   ├── error.go         # Error handling
│   │   └── response.go      # Response formatting
│   ├── server.go            # HTTP server setup
│   └── router.go            # Main router setup
├── worker/                  # Worker application components
│   ├── job/                 # Job definitions
│   │   ├── base_job.go      # Base job interface and common functionality
│   │   ├── order_job.go     # Order processing jobs
│   │   └── data_job.go      # Data processing jobs
│   ├── queue/               # Queue implementations
│   │   ├── kafka.go         # Kafka queue adapter
│   │   └── rabbitmq.go      # RabbitMQ queue adapter
│   ├── processor.go         # Job processor
│   └── manager.go           # Worker management
├── scheduler/               # Scheduler application components
│   ├── job/                 # Scheduled job definitions
│   │   ├── base_job.go      # Base scheduled job interface
│   │   ├── market_job.go    # Market data jobs
│   │   └── report_job.go    # Reporting jobs
│   ├── cron.go              # Cron job scheduler
│   └── registry.go          # Job registry
├── config/                  # Application configuration
│   ├── config.go            # Configuration loader
│   └── validation.go        # Configuration validation
├── bootstrap/               # Application bootstrapping
│   ├── web.go               # Web application setup
│   ├── worker.go            # Worker application setup
│   ├── scheduler.go         # Scheduler application setup
│   └── module.go            # Module registration
├── api/                     # API specifications
│   ├── rest/                # REST API
│   │   ├── server.go        # Server setup
│   │   └── router.go        # Main router setup
│   ├── grpc/                # gRPC API
│   │   ├── server.go        # Server setup
│   │   └── interceptor.go   # gRPC interceptors
│   └── websocket/           # WebSocket API
│       ├── server.go        # Server setup
│       └── handler.go       # WebSocket handlers
└── messaging/               # Inter-application messaging
    ├── event/               # Event definitions
    │   ├── base_event.go    # Base event structure
    │   └── user_event.go    # User-related events
    └── publisher.go         # Event publisher
```

## Common Utilities

### `internal/common/`

```
common/
├── constants/               # Constants used across the app
│   └── appconstants.go
├── utils/                   # Utility functions
│   ├── timeutil/            # Time-related utilities
│   │   └── timeutil.go
│   ├── mathutil/            # Mathematical utilities
│   │   └── mathutil.go
│   ├── idgen/               # ID generation utilities
│   │   └── idgen.go
│   └── validation/          # Validation utilities
│       └── validation.go
├── types/                   # Common types
│   ├── decimal.go           # Decimal type for financial calculations
│   └── nullable.go          # Nullable types
├── errors/                  # Common error handling
│   ├── apperrors.go         # Application error types
│   └── errorhandling.go     # Error handling utilities
└── patterns/                # Common design patterns
    ├── observer.go          # Observer pattern implementation
    └── retry.go             # Retry pattern implementation
```

## Public Libraries

### `pkg/`

```
pkg/
├── indicator/               # Technical indicators library
│   ├── moving_average.go
│   ├── oscillators.go
│   └── volume.go
├── trading/                 # Trading utilities library
│   ├── ordertype.go
│   └── positionsizing.go
└── financial/               # Financial calculations library
    ├── returns.go
    └── risk.go
```

## Entry Points

### `cmd/`

```
cmd/
├── web/                     # REST API server
│   ├── main.go              # Entry point for the web application
│   ├── wire.go              # Dependency injection definitions
│   └── wire_gen.go          # Generated dependency injection code
├── worker/                  # Background worker
│   ├── main.go              # Entry point for the worker application
│   ├── wire.go              # Dependency injection definitions
│   └── wire_gen.go          # Generated dependency injection code
├── scheduler/               # Cron job scheduler
│   ├── main.go              # Entry point for the scheduler application
│   ├── wire.go              # Dependency injection definitions
│   └── wire_gen.go          # Generated dependency injection code
├── trader/                  # Trading engine
│   └── main.go
├── backtest/                # Backtesting application
│   └── main.go
├── cli/                     # Command-line tools
│   └── main.go
└── migrate/                 # Database migration tool
    ├── main.go
    ├── wire.go              # Dependency injection definitions
    └── wire_gen.go          # Generated dependency injection code
```

## Configuration

### `configs/`

```
configs/
├── local.yaml               # Local development config
├── dev.yaml                 # Development environment config
├── prod.yaml                # Production environment config
└── test.yaml                # Testing environment config
```

## Deployment

### `deployments/`

```
deployments/
├── docker/                  # Docker configurations
│   ├── Dockerfile           # Main application Dockerfile
│   ├── Dockerfile.dev       # Development Dockerfile
│   └── docker-compose.yml   # Docker Compose configuration
├── kubernetes/              # Kubernetes manifests
│   ├── base/                # Base configurations
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── overlays/            # Environment-specific overlays
│       ├── dev/
│       │   └── kustomization.yaml
│       └── prod/
│           └── kustomization.yaml
└── terraform/               # Infrastructure as Code
    ├── modules/             # Reusable Terraform modules
    │   ├── database/
    │   └── kubernetes/
    └── environments/        # Environment-specific configurations
        ├── dev/
        └── prod/
```

## Testing Structure

### `test/`

```
test/
├── unit/                    # Unit tests
│   ├── domain/              # Domain model tests
│   └── modules/             # Module-specific tests
├── integration/             # Integration tests
│   ├── database/            # Database integration tests
│   └── api/                 # API integration tests
├── e2e/                     # End-to-end tests
│   ├── scenarios/           # Test scenarios
│   └── utils/               # Test utilities
├── performance/             # Performance tests
│   ├── load/                # Load tests
│   └── benchmark/           # Benchmark tests
└── fixtures/                # Test fixtures
    ├── users.json           # User test data
    └── stocks.json          # Stock test data
```

## Documentation

### `docs/`

```
docs/
├── architecture/            # Architecture documentation
│   ├── overview.md          # System overview
│   └── diagrams/            # Architecture diagrams
│       ├── system.png
│       └── data_flow.png
├── api/                     # API documentation
│   ├── rest.md              # REST API documentation
│   └── grpc.md              # gRPC API documentation
├── development/             # Development guides
│   ├── setup.md             # Development setup guide
│   └── guidelines.md        # Coding guidelines
└── operations/              # Operations guides
    ├── deployment.md        # Deployment guide
    └── monitoring.md        # Monitoring guide
```

## Module Interactions and Dependency Rule

1. **Dependency Flow**:
   - `cmd` → `internal/app` → `internal/modules` → `internal/domain`
   - External dependencies point inward, never outward

2. **Module Independence**:
   - Each module is self-contained with its own layers
   - Modules communicate through well-defined interfaces
   - No circular dependencies between modules

3. **Domain Integrity**:
   - Domain models contain business rules and validation
   - Domain services implement business processes
   - Domain layer has no dependencies on external frameworks

## Module Definition Template

For each new module, follow this structure:

```
modules/{module_name}/
├── application/              # Application services and DTOs
│   ├── service/              # Application services
│   │   └── {entity}_service.go # Business logic for each entity
│   ├── dto/                  # Data Transfer Objects
│   │   └── {entity}_dto.go   # DTOs for each entity
│   ├── command/              # Command handlers (CQRS pattern)
│   │   └── {entity}_commands.go # Command handlers for each entity
│   └── query/                # Query handlers (CQRS pattern)
│       └── {entity}_queries.go  # Query handlers for each entity
├── domain/                   # Domain models and interfaces
│   ├── model/                # Domain models
│   │   └── {entity}.go       # Entity models
│   ├── repository/           # Repository interfaces
│   │   └── {entity}_repository.go # Repository interface for each entity
│   └── service/              # Domain service interfaces
│       └── {entity}_service.go # Domain service for each entity
├── infrastructure/           # Repository and external service implementations
│   ├── persistence/          # Repository implementations
│   │   └── {entity}_repo.go  # Repository implementation for each table
│   └── external/             # External service integrations
│       └── {service}_adapter.go # Adapter for external service
├── interfaces/               # API controllers and handlers
│   ├── http/                 # HTTP handlers
│   │   ├── controller.go     # REST API controllers
│   │   ├── dto.go            # Request/Response objects
│   │   └── router.go         # Route definitions
│   ├── worker/               # Worker job handlers
│   │   └── {job}_handler.go  # Handler for worker jobs
│   └── scheduler/            # Scheduler job handlers
│       └── {job}_handler.go  # Handler for scheduled jobs
├── di.go                     # Dependency injection setup
└── module.go                 # Module registration
```

## File Naming Conventions

1. **Singular Nouns**: Use singular form for entities (`user.go` not `users.go`)
2. **Lowercase with Underscores**: For multi-word file names (`order_service.go`)
3. **Descriptive Suffixes**: Indicate file purpose with suffixes:
   - `*_repository.go`: Repository implementations
   - `*_service.go`: Service implementations
   - `*_controller.go`: API controllers
   - `*_test.go`: Test files
4. **One File Per Table/Entity**: Create one repository and service file for each database table or domain entity:
   - Repositories:
     - `user_repo.go`: For the user table
     - `role_repo.go`: For the role table
     - `user_role_repo.go`: For the user_roles join table
   - Services:
     - `user_service.go`: For user-related business logic
     - `role_service.go`: For role-related business logic
     - `user_role_service.go`: For user-role relationship business logic

## Implementation Considerations

1. **Dependency Injection**: Use constructor injection for dependencies
2. **Error Handling**: Return domain-specific errors that can be translated to appropriate responses
3. **Configuration**: Use environment-based configuration with validation
4. **Observability**: Implement logging, metrics, and tracing across all layers
5. **State Management**: Prefer immutable state and explicit state transitions
6. **Concurrency**: Use goroutines and channels with explicit error handling
7. **Testing**: Write unit tests for business logic and integration tests for infrastructure

This structure promotes:
- Clear separation of concerns
- Deterministic file locations
- Scalability for dozens of modules
- Testability at all levels
- Maintainable and extendable architecture


### `deployments/`
- `docker/`: Docker configurations
- `kubernetes/`: Kubernetes manifests

## Implementation Considerations

1. Use dependency injection for modular component testing
2. Consider concurrency patterns for high-throughput operations
3. Implement circuit breakers for external service calls
4. Use proper error handling with context
5. Develop observability patterns for production monitoring
6. Follow specific guidelines for web, worker, and scheduler applications (see `application-types.md`)
7. Use consistent messaging patterns between applications
8. Implement proper retry mechanisms for all async operations
9. Ensure proper configuration management across all application types

This structure promotes:
- Clear separation of concerns
- Testable components
- Maintainable codebase
- Scalable architecture
- Easy extension with new strategies
- Reliable application interactions