.PHONY: build run clean test-coverage test coverage docker-build docker-run db-start db-stop build-live run-live build-migrate migrate build-client client build-test-client test-client test-simulator test-live deployment web live

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

# ---------------------------------------------------------- #

web:
	@echo "Generating wire dependencies..."
	@wire cmd/web/wire.go
	@echo "Building the web application..."
	@go build -gcflags="all=-N -l" -o bin/web cmd/web/wire_gen.go cmd/web/main.go
	@echo "Running the web application..."
	@./bin/web

run: web

# ---------------------------------------------------------- #

live:
	@echo "Building the live WebSocket server..."
	@go build -gcflags="all=-N -l" -o bin/live cmd/live/main.go
	@echo "Running the live WebSocket server..."
	@./bin/live

# ---------------------------------------------------------- #

client:
	@echo "Building the market data simulator client..."
	@go build -o bin/client cmd/client/main.go
	@echo "Running the market data simulator server..."
	@./bin/client

# ---------------------------------------------------------- #

migrate:
	@echo "Building the data migration tool..."
	@go build -o bin/migrate cmd/migrate/main.go
	@echo "Running the data migration tool..."
	@./bin/migrate

# ---------------------------------------------------------- #

deployment:
	@echo "Building the deployment tool..."
	@go build -gcflags="all=-N -l" -o bin/deployment cmd/deployment/main.go
	@echo "Running the deployment tool..."
	@./bin/deployment

# ---------------------------------------------------------- #

sync:
	git checkout main
	git pull --rebase
	git checkout -
	git merge origin main
	git push

# ---------------------------------------------------------- #

git:
	cd ~/workspace/sutram
	sudo lsof -i :9090 -t | xargs sudo kill -9
	git pull --rebase

	cd ~/workspace/mukham
	sudo lsof -i :3000 -t | xargs sudo kill -9
	git pull --rebase


# ---------------------------------------------------------- #

db:
	@make -C db all

# Database Migration Commands
db-migrate-up:
	@echo "Running all pending database migrations..."
	@make -C db migrate-up

db-migrate-down:
	@echo "Rolling back database migrations..."
	@make -C db migrate-down

db-migrate-version:
	@echo "Checking current migration version..."
	@make -C db migrate-version

db-migrate-validate:
	@echo "Validating database migrations..."
	@make -C db migrate-validate

db-setup:
	@echo "Setting up database with migrations..."
	@make -C db setup

db-fresh:
	@echo "Fresh database setup..."
	@make -C db fresh

.PHONY: db-migrate-up db-migrate-down db-migrate-version db-migrate-validate db-setup db-fresh
