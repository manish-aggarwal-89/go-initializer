.PHONY: all build run test swagger migrate init tidy clean docker-build mock help mockV3

BINARY_NAME = {{SERVICE_NAME}}

# Default target
all: build
MAIN_CMD = ./cmd/{{PROVIDER}}
MIGRATION_CMD = ./cmd/migration

# Mockery
MOCKERY := $(shell go env GOPATH)/bin/mockery

# Build the application binary
build:
	@echo ">> Building $(BINARY_NAME)"
	@go build -ldflags "-s -w" -o $(BINARY_NAME) $(MAIN_CMD)

# Run the application
run:
	@go run $(MAIN_CMD)

# Run unit tests
test:
	@echo ">> Running unit tests"
	@go test ./... -cover -covermode=count

# Run unit tests with coverage report
test-coverage:
	@echo ">> Running tests with coverage"
	@if [ -f coverage.out ]; then rm coverage.out; fi
	@go test ./... -cover -coverprofile=coverage.out -covermode=count
	@go tool cover -func=coverage.out

# Generate Swagger docs via scripts (uses swagger_exclude.sh + swagger_gen.sh)
# Requires: go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	@chmod +x scripts/swagger_gen.sh scripts/swagger_exclude.sh
	@./scripts/swagger_gen.sh {{PROVIDER}}

# this will use .mockery.yml template, make sure your mockery version is 3.
mockV3:
	@mockery

migrate-up:
	@echo ">> Running pending migrations"
	@go run cmd/migration/main.go up

migrate-down:
	@echo ">> Rolling back last migration"
	@go run cmd/migration/main.go down

migrate-new:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-new NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo ">> Creating new migration: $(NAME)"
	@go run cmd/migration/main.go new $(NAME)

migrate-down-to:
	@if [ -z "$(word 2,$(MAKECMDGOALS))" ]; then \
		echo "Error: Migration ID required"; \
		echo "Usage: make migrate-down-to 20240101120000"; \
		exit 1; \
	fi
	@echo "Rolling back to migration ID: $(word 2,$(MAKECMDGOALS))"
	@go run cmd/migration/main.go down-to $(word 2,$(MAKECMDGOALS))

migrate-new-seed:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-new-seed NAME=your_seed_name"; \
		exit 1; \
	fi
	@echo ">> Creating new seed migration: $(NAME)"
	@go run cmd/migration/main.go new-seed $(NAME)

migrate-help:
	@echo "Migration Commands:"
	@echo "  migrate-up       - Run all pending migrations"
	@echo "  migrate-down     - Rollback last migration"
	@echo "  migrate-down-to  - Rollback to a specific migration ID (usage: make migrate-down-to 20240101120000)"
	@echo "  migrate-seed     - Run data seeding"
	@echo "  migrate-new      - Create new migration (usage: make migrate-new NAME=your_migration_name)"
	@echo "  migrate-new-seed - Create new seed migration (usage: make migrate-new-seed NAME=your_seed_name)"
	@echo "  migrate-help     - Show this help message"