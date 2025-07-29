# Single stage - just copy the pre-built binary
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh sutram

# Set working directory
WORKDIR /app

# Copy the pre-built binaries from the host
COPY bin/web ./web
COPY bin/db-migrate ./db-migrate
COPY bin/ch-migrate ./ch-migrate

# Change ownership
RUN chown sutram:sutram /app/web /app/db-migrate /app/ch-migrate

# Switch to non-root user
USER sutram

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./web"]
