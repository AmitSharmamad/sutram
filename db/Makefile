ch:
	docker-compose up -d sutram-clickhouse
pg:
	docker-compose up -d sutram-postgres

all:
	make down
	make clean
	make up

up:
	docker-compose up -d --build

down:
	docker-compose down -v

clean:
	rm -rf db/ch
	rm -rf db/pg

start:
	docker-compose down -v
	docker-compose up --build -d

# Migration targets
migrate-build:
	@echo "Building migration binary..."
	@cd .. && go build -o bin/db-migrate cmd/db-migrate/main.go

migrate-up: migrate-build
	@echo "Running all pending migrations..."
	@cd .. && ./bin/db-migrate -action=up

migrate-down: migrate-build
	@echo "Rolling back all migrations..."
	@cd .. && ./bin/db-migrate -action=down

migrate-steps: migrate-build
	@echo "Running migration steps (use STEPS=n to specify number)..."
	@cd .. && ./bin/db-migrate -action=steps -steps=$(or $(STEPS),1)

migrate-version: migrate-build
	@echo "Checking current migration version..."
	@cd .. && ./bin/db-migrate -action=version

migrate-validate: migrate-build
	@echo "Validating migrations..."
	@cd .. && ./bin/db-migrate -action=validate

migrate-force: migrate-build
	@echo "Forcing migration version (use VERSION=n to specify)..."
	@cd .. && ./bin/db-migrate -action=force -version=$(VERSION)

# Migration setup - runs database and applies migrations
setup: up migrate-up
	@echo "Database setup complete!"

# Fresh setup - clean slate with migrations
fresh: down clean up migrate-up
	@echo "Fresh database setup complete!"

.PHONY: migrate-build migrate-up migrate-down migrate-steps migrate-version migrate-validate migrate-force setup fresh
