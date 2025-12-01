# wa-service Makefile

.PHONY: help build run test test-unit test-integration test-all test-coverage test-integration-coverage clean mocks regen-mocks migrate-create migrate-up migrate-down migrate-reset migrate-version migrate-force migrate-goto migrate-help

# Default target
help:
	@echo "Available targets:"
	@echo "  help           Show this help message"
	@echo "  build          Build the application"
	@echo "  run            Run the application"
	@echo "  test           Run unit tests (fast feedback loop)"
	@echo "  test-unit      Run unit tests only"
	@echo "  test-integration Run integration tests only"
	@echo "  test-all       Run both unit and integration tests"
	@echo "  test-coverage  Run unit tests with coverage"
	@echo "  test-integration-coverage Run integration tests with coverage"
	@echo "  mocks          Generate mocks"
	@echo "  regen-mocks    Clean and regenerate mocks"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Database Migration Commands:"
	@echo "  migrate-create  Create new migration template"
	@echo "  migrate-up      Run all pending migrations"
	@echo "  migrate-down    Rollback last migration"
	@echo "  migrate-reset   Reset all migrations"
	@echo "  migrate-version Show current migration version"
	@echo "  migrate-force   Force migration version"
	@echo "  migrate-goto    Migrate to specific version"
	@echo "  migrate-help    Show migration commands help"

# Application targets
build:
	go build -o bin/github.com/fajarAnd/workshop-brin/wa-service cmd/bot/main.go

run:
	go run cmd/bot/main.go

# Testing targets
test:
	go test ./internal/app/...

test-unit:
	go test -tags=unit ./internal/app/...

test-integration:
	go test -tags=integration ./internal/app/...

test-all:
	go test ./internal/app/...

test-coverage:
	go test -cover ./internal/app/...

test-integration-coverage:
	go test -tags=integration -cover ./internal/app/...

# Mock generation
mocks:
	mockery

regen-mocks:
	rm -rf mocks/
	mockery

clean:
	rm -rf bin/ mocks/

# Migration Management Commands
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Run all pending migrations
migrate-up:
	go run cmd/migrate/main.go up

# Rollback last migration
migrate-down:
	go run cmd/migrate/main.go down

# Reset all migrations (down to zero)
migrate-reset:
	go run cmd/migrate/main.go reset

# Show current migration version
migrate-version:
	go run cmd/migrate/main.go version

# Force migration to specific version (for fixing broken state)
migrate-force:
	@read -p "Enter version to force: " version; \
	go run cmd/migrate/main.go force $$version

# Migrate to specific version
migrate-goto:
	@read -p "Enter target version: " version; \
	go run cmd/migrate/main.go goto $$version

# Show migration commands help
migrate-help:
	@echo "Database Migration Commands:"
	@echo "  make migrate-create  - Create new migration template"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-reset   - Reset all migrations"
	@echo "  make migrate-version - Show current migration version"
	@echo "  make migrate-force   - Force migration version"
	@echo "  make migrate-goto    - Migrate to specific version"
	@echo "  make migrate-help    - Show this help"