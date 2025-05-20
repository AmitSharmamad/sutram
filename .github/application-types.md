# Application Types Guidelines

This document provides detailed guidelines for the three main application types in the platform: Web (REST APIs), Worker (Async Task Processing), and Scheduler (Cron Jobs).

## Web Application (REST APIs)

### Responsibilities
- Handle HTTP requests and responses
- Route requests to appropriate controllers
- Validate incoming data
- Handle authentication and authorization
- Format and return responses
- Manage API versioning

### Architecture Guidelines
1. **Controller-Service Pattern**:
   - Controllers handle request/response and delegate to services
   - Services contain business logic and orchestrate domain operations
   - Use DTOs (Data Transfer Objects) for input/output

2. **Middleware Stack**:
   - Authentication middleware
   - Request logging middleware
   - Request ID middleware for tracing
   - CORS middleware
   - Rate limiting middleware
   - Error handling middleware

3. **API Versioning**:
   - Use URL path versioning (e.g., `/api/stock`)
   - Each version should have its own router group
   - Maintain backward compatibility

4. **Request Validation**:
   - Validate all incoming requests
   - Return descriptive validation errors
   - Use struct tags for validation rules

5. **Response Formatting**:
   - Consistent response structure
   - Proper HTTP status codes
   - Pagination for list endpoints
   - Include metadata in responses when appropriate

## Worker Application (Async Task Processing)

### Responsibilities
- Process background jobs
- Handle long-running tasks
- Provide reliable task execution
- Manage task retries and error handling
- Scale horizontally for high-throughput processing

### Architecture Guidelines
1. **Job Definition**:
   - Each job type should be defined as a struct
   - Jobs should be serializable
   - Jobs should contain all necessary execution data

2. **Worker Pattern**:
   - Workers should pull jobs from queues
   - Each worker should have a single responsibility
   - Workers should be stateless for horizontal scaling

3. **Queue Management**:
   - Use reliable messaging system (Kafka, RabbitMQ, etc.)
   - Implement dead-letter queues for failed jobs
   - Set appropriate timeouts and TTLs

4. **Error Handling**:
   - Implement exponential backoff for retries
   - Log detailed error information
   - Provide alerting for persistent failures

5. **Observability**:
   - Track job execution metrics (success, failure, duration)
   - Log job processing events
   - Monitor queue depth and processing latency

## Scheduler Application (Cron Jobs)

### Responsibilities
- Schedule time-based jobs
- Manage job execution intervals
- Handle job dependencies
- Ensure jobs run at correct times
- Manage execution history

### Architecture Guidelines
1. **Job Definition**:
   - Define jobs declaratively
   - Specify schedule using cron syntax
   - Include timeout and retry policies
   - Define dependencies between jobs

2. **Scheduling Patterns**:
   - Interval-based jobs (every N minutes)
   - Cron-based jobs (at specific times)
   - Event-triggered jobs (based on system events)
   - Dependent jobs (execute after another job finishes)

3. **Execution Management**:
   - Use distributed locking for cluster environments
   - Implement leader election for single-execution jobs
   - Track execution history and status

4. **Error Handling**:
   - Define retry policies per job
   - Log detailed error information
   - Alert on critical job failures

5. **Observability**:
   - Track execution times and durations
   - Monitor job success and failure rates
   - Alert on missed or consistently failing jobs

## Cross-Cutting Concerns

### Configuration
- Use environment variables for runtime configuration
- Load job-specific configurations from a centralized source
- Support different configurations per environment

### Logging
- Use structured logging
- Include context information in log entries
- Log at appropriate levels (debug, info, warning, error)

### Metrics
- Track application and job-specific metrics
- Monitor resource usage (CPU, memory, disk, network)
- Set up alerts for abnormal patterns

### Health Checks
- Implement liveness and readiness probes
- Include dependency health checks
- Expose health check endpoints

### Deployment
- Use containerization (Docker)
- Support orchestration platforms (Kubernetes)
- Implement proper resource requests and limits