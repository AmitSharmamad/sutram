# Gateway Service

This service acts as a reverse proxy for the Sutram application, routing traffic to the appropriate backend services.

## Configuration

The gateway routes traffic as follows:
- All requests to `/api/*` are routed to the API server running on port `9090`
- All other requests are routed to the frontend server running on port `3000`

The gateway server runs on port `8010`.

## Usage

### Build the Gateway

```bash
cd gateway
make build
```

### Start the Gateway

```bash
cd gateway
make start
```

### Stop the Gateway

```bash
cd gateway
make stop
```

### View Logs

```bash
cd gateway
make logs
```

### Quick Run (Development)

For quick local development, you can use the run-gateway.sh script which uses Docker without requiring docker-compose:

```bash
cd gateway
make quick-run
```

or directly:

```bash
cd gateway
./run-gateway.sh
```

## Testing

You can test the gateway by accessing:

- Frontend: http://localhost:8010/
- API: http://localhost:8010/api/

## Configuration Files

- `Dockerfile` - Docker image configuration for Nginx
- `nginx.conf` - Main Nginx configuration
- `conf.d/default.conf` - Server and upstream configuration
