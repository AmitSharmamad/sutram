# Single stage - just copy the pre-built binary
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh sutram

# Set working directory
WORKDIR /app

# Copy the pre-built binary from the host
COPY bin/web ./web

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
