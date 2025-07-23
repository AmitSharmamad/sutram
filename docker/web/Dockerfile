# Build stage
FROM golang:1.24.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install git, ca-certificates, and wire
RUN apk add --no-cache git ca-certificates && \
    go install github.com/google/wire/cmd/wire@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate wire code and build the web application
RUN wire cmd/web/wire.go && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o web cmd/web/wire_gen.go cmd/web/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh sutram

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/web .

# Change ownership
RUN chown sutram:sutram /app/web

# Switch to non-root user
USER sutram

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./web"]
