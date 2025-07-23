# Build stage
FROM golang:1.24.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the deployment application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deployment cmd/deployment/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh sutram

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/deployment .

# Change ownership
RUN chown sutram:sutram /app/deployment

# Switch to non-root user
USER sutram

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8082/health || exit 1

# Run the application
CMD ["./deployment"]
