FROM golang:1.23

WORKDIR /app

# Copy the source code
COPY src src
COPY data data
COPY go.mod go.mod
COPY go.sum go.sum

# Build the application
RUN go build -o main src/main.go

# Expose the port your application listens on
EXPOSE 8080

# Command to run your application
CMD ["/app/main"]
