FROM golang:1.23-bookworm AS builder

ENV GOINSECURE=*.github.com

WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates

# Copy the source code
COPY src src
COPY go.mod go.mod
COPY go.sum go.sum

# RUN go mod download

# Build the application
RUN go build -o ./main src/main.go

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /app/main

# Expose the port your application listens on
EXPOSE 8080

# Command to run your application
CMD ["/app/main"]
