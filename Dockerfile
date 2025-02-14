FROM golang:1.23-bookworm AS builder

WORKDIR /app

RUN go install github.com/google/wire/cmd/wire@latest

COPY src src
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

RUN make build

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/s-web s-web
COPY .env .env

# Expose the port your application listens on
EXPOSE 9090

# Command to run your application
CMD ["/app/s-web"]
