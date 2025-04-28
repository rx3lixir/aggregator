# Binary output name
BINARY_NAME=api

# Migration settings
DB_URL=postgres://rx3lixir:password@localhost:5432/aggregator?sslmode=disable
MIGRATIONS_PATH=./migrations/

# Default make command
all: build

# Build the binary
build:
	@echo "Building..."
	go build -o ./bin/$(BINARY_NAME) ./cmd/api/main.go

# Run the server
run: build
	@echo "Running..."
	./bin/$(BINARY_NAME)

# Test your application
test:
	@echo "Testing..."
	go test ./...

# Clean up binaries
clean:
	@echo "Cleaning..."
	go clean
	rm -f ./bin/$(BINARY_NAME)

# Create a new migration
migrate-new:
ifndef name
	$(error name is not set. Usage: make migrate-new name=your_migration_name)
endif
	@echo "Creating new migration '$(name)'..."
	docker run --rm -v $(shell pwd)/$(MIGRATIONS_PATH):/migrations migrate/migrate create -ext sql -dir /migrations $(name)

# Apply all up migrations
migrate-up:
	@echo "Applying migrations..."
	docker run --rm -v $(shell pwd)/$(MIGRATIONS_PATH):/migrations --network host migrate/migrate -path=/migrations -database "$(DB_URL)" up

# Help command to display available commands
help:
	@echo "Makefile commands:"
	@echo "all                  - Build the application"
	@echo "build                - Build the binary"
	@echo "run                  - Build and run the application"
	@echo "test                 - Run tests"
	@echo "clean                - Remove binaries"
	@echo "migrate-up           - Apply all up migrations"
	@echo "migrate-down         - Rollback one migration"
	@echo "migrate-force-down-all - Rollback all migrations"
	@echo "migrate-new name=xxx - Create a new migration file"
	@echo "help                 - Display this help"

# Mark commands that don't correspond to files as .PHONY
.PHONY: all build run test clean help migrate-up migrate-down migrate-force-down-all migrate-new
