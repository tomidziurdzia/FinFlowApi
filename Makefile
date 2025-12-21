.PHONY: help migrate migrate-up migrate-down migrate-status db-verify build run test test-integration clean

help:
	@echo "Available commands:"
	@echo "  make migrate        - Run all pending migrations"
	@echo "  make migrate-up     - Run all pending migrations (alias for migrate)"
	@echo "  make migrate-status  - Show migration status"
	@echo "  make db-verify      - Verify database connection and show users table structure"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests (skips integration tests)"
	@echo "  make test-integration - Run all tests including integration tests"
	@echo "  make clean          - Clean build artifacts"

migrate migrate-up:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@export $$(cat .env | grep -v '^#' | xargs); \
	if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL is not configured in .env"; \
		echo "Make sure you have DATABASE_URL in your .env file"; \
		exit 1; \
	fi; \
	echo "Running migrations..."; \
	for migration in internal/infrastructure/db/migrations/*.sql; do \
		if [ -f "$$migration" ]; then \
			echo "Executing $$migration..."; \
			psql "$$DATABASE_URL" -f "$$migration" || exit 1; \
		fi; \
	done; \
	echo "✅ All migrations executed successfully"

migrate-status:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@export $$(cat .env | grep -v '^#' | xargs); \
	if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL is not configured in .env"; \
		exit 1; \
	fi; \
	echo "Database connection: OK"; \
	echo ""; \
	echo "Tables in database:"; \
	psql "$$DATABASE_URL" -c "\dt" || exit 1; \
	echo ""; \
	echo "Users table structure:"; \
	psql "$$DATABASE_URL" -c "\d users" || echo "Users table does not exist"

db-verify:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@export $$(cat .env | grep -v '^#' | xargs); \
	if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL is not configured in .env"; \
		exit 1; \
	fi; \
	echo "Verifying database connection..."; \
	psql "$$DATABASE_URL" -c "SELECT version();" > /dev/null 2>&1 && echo "✅ Database connection successful" || (echo "❌ Database connection failed" && exit 1); \
	echo ""; \
	echo "Users table structure:"; \
	psql "$$DATABASE_URL" -c "\d users" 2>/dev/null || echo "Users table does not exist"

build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "✅ Build complete: bin/api"

run:
	@echo "Running application..."
	@go run cmd/api/main.go

test:
	@echo "Running tests (skipping integration tests)..."
	@SKIP_DB_TESTS=true go test -v ./...

test-integration:
	@echo "Running all tests including integration tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "✅ Clean complete"

