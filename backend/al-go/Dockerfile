FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN go build -o main src/main.go

# Use a smaller base image for the final image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/main /main

# Expose the port your application listens on
EXPOSE 8080

# Command to run your application
CMD ["/main"]
