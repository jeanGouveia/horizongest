.PHONY: help build run test clean db-reset db-setup

# Default target
help:
	@echo "HorizonGest Backend - Available Commands:"
	@echo ""
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make db-setup       - Setup PostgreSQL database and user"
	@echo "  make db-reset       - Reset database (truncate all tables, preserve schema)"
	@echo ""

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Setup PostgreSQL database and user
db-setup:
	@echo "Setting up PostgreSQL database..."
	@./setup-postgres.sh

# Reset database (truncate all tables, preserve schema)
db-reset:
	@echo "Resetting database..."
	@PGPASSWORD=horizongest_secure_password psql -U horizongest_user -h localhost -d horizongest -f scripts/reset_database.sql
	@echo "Database reset completed successfully"
