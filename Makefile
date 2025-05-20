.PHONY: build run clean test-coverage test coverage docker-build docker-run db-start db-stop build-live run-live

# Build the web application
build:
	@echo "Building the web application..."
	@mkdir -p bin
	@wire ./cmd/web
	@go build -o bin/web-server cmd/web/main.go cmd/web/wire_gen.go

# Run the web application
run: build
	@echo "Running the web server..."
	@./bin/web-server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

# Run tests with coverage analysis
coverage:
	@echo "Running tests with coverage analysis..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated in coverage.html"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t sutram-web -f Dockerfile.web .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env sutram-web

# Start the database
db-start:
	@echo "Starting the database..."
	@cd db && docker-compose up -d

# Stop the database
db-stop:
	@echo "Stopping the database..."
	@cd db && docker-compose down

# Start everything with docker-compose
all-up:
	@echo "Starting all services..."
	@docker-compose up -d

# Stop everything
all-down:
	@echo "Stopping all services..."
	@docker-compose down

# Build the live WebSocket server
build-live:
	@echo "Building the live WebSocket server..."
	@go build -gcflags="all=-N -l" -o s-live cmd/live/main.go

# Run the live WebSocket server
live: build-live
	@echo "Running the live WebSocket server..."
	@./s-live
	@echo "Stopping all services..."
	@docker-compose down
