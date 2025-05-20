# Sutram Web Application

This is the web application component of the Sutram algorithmic trading platform.

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL
- ClickHouse (for analytics)

## Quick Start

### Running with Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
# Start all services (web, PostgreSQL, ClickHouse)
docker-compose up -d
```

### Running WebSocket Server for Live Data

The platform includes a WebSocket server for streaming live market data:

1. Configure your WebSocket environment variables:
```bash
cp .env.websocket.example .env
# Edit .env with your actual credentials
```

2. Run the WebSocket server using make:
```bash
# Using make directly
make run-live

# Or using the convenience script
./run-live.sh
```

The WebSocket server will run on port 9091 by default and will connect to the upstream
WebSocket provider (wss://api-feed.dhan.co) with the credentials provided in the .env file.

### Running Locally

1. Start the databases:

```bash
cd db
docker-compose up -d
```

2. Run the web application:

```bash
./run-web.sh
```

The application will be available at http://localhost:8080.

## API Endpoints

### User Module

- **POST /api/users/register**: Register a new user
  - Request: `{ "name": "...", "username": "...", "email": "...", "password": "..." }`
  - Response: `{ "success": true, "message": "User registered successfully", "data": { "token": "...", "user": {...}, "expires_at": 123456789 } }`

- **POST /api/users/login**: Login a user
  - Request: `{ "email": "...", "password": "..." }`
  - Response: `{ "success": true, "message": "Login successful", "data": { "token": "...", "user": {...}, "expires_at": 123456789 } }`

- **GET /api/users/profile**: Get user profile (requires authentication)
  - Response: `{ "success": true, "data": {...} }`

## OAuth Authentication

The application supports OAuth authentication with Google. See [OAUTH.md](OAUTH.md) for setup instructions.

## Development

### Project Structure

- `cmd/web`: Web application entry point
- `internal/modules`: Business modules (user, strategy, etc.)
- `internal/domain`: Core domain models and business rules
- `db`: Database configuration and SQL scripts

### Building from Source

```bash
# Generate wire_gen.go for dependency injection
wire ./cmd/web

# Build the binary
go build -o s-web cmd/web/main.go
```
