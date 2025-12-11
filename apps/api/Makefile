# Define variables
GO := go
PROJECT_NAME := collab-board
BINARY_NAME := $(PROJECT_NAME)
DB_NAME := collab_board
# Local
# DATABASE_URL ?= "postgres://admin:admin@localhost:5432/collab_board?sslmode=disable"
# Docker
DATABASE_URL ?= "postgresql://root:password@localhost:5433/$(DB_NAME)?sslmode=disable"

# Define targets
.PHONY: help build postgres createdb dropdb psql run test clean migrateup migratedown setup start env lint fmt check

# Define commands
CREATE_CONTAINER := docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=root \
  -e POSTGRES_DB="$(DB_NAME)" \
  -e PGDATA=/var/lib/postgresql/data/pgdata \
  -v postgres-data:/var/lib/postgresql/data/pgdata \
  -p 5433:5432 \
  -d postgres
PERSIST_DATA := docker volume create pgdata
MIGRATE_UP := migrate -path ./db/migrations -database $(DATABASE_URL) -verbose up
MIGRATE_DOWN := migrate -path ./db/migrations -database $(DATABASE_URL) -verbose down
CREATE_DB := docker exec -it postgres createdb --username=root --owner=root $(DB_NAME)
DROP_DB := docker exec -it posgres dropdb $(DB_NAME)
DOCKER_PSQL := docker exec -it postgres psql -U root -d $(DB_NAME)

# Duuh
help:  ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/) ? substr(lastLine, RSTART + 2, RLENGTH) : ""; \
		printf "  %-20s %s\n", $$1, helpMessage; \
	} { lastLine = $$0 }' $(MAKEFILE_LIST)

# Create Docker container and PostgreSQL Database
postgres:
	@echo "Creating a Docker container and PostgreSQL Database..."
	$(CREATE_CONTAINER)

# Create PostgreSQL Database (In case the Docker container is already created)
createdb:
	@echo "Creating the PostgreSQL database..."
	$(CREATE_DB)

# Drop PostgreSQL Database
dropdb:
	@echo "Dropping the PostgreSQL database..."
	$(DROP_DB)

# PSQL into Docker Container Database
psql:
	@echo "PSQL-ing..."
	$(DOCKER_PSQL)

# Build and run targets
build:  ## Build the Go application
	@echo "Building the application..."
	$(GO) build -o $(BINARY_NAME) .

# Run application
run: build  ## Build and run the Go application
	@echo "Running the application..."
	./$(BINARY_NAME)

# Run migrate up
migrateup: 
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL is not set. Please set it before running this target."; \
		exit 1; \
	fi
	@echo "Applying migrations..."
	$(MIGRATE_UP)

# Run migrate down
migratedown: 
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL is not set. Please set it before running this target."; \
		exit 1; \
	fi
	@echo "Removing migrations..."
	$(MIGRATE_DOWN)

# Run tests
test:  ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

# Clean up generated files
clean:  ## Clean up generated files
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Setup and maintenance targets
setup:  ## Set up the project (download dependencies)
	@echo "Setting up the project..."
	$(GO) mod download
	$(MAKE) generate

# Display current environment variables
env: 
	@echo "Current environment settings:"
	@echo "DATABASE_URL=$(DATABASE_URL)"

# Run linting on the Go code
lint: 
	@echo "Running golint..."
	@command -v golint >/dev/null 2>&1 || { echo >&2 "golint is not installed. Install it by running: go install golang.org/x/lint/golint@latest"; exit 1; }
	golint ./...

# Format the Go code
fmt: 
	@echo "Formatting Go code..."
	$(GO) fmt ./...

# Run formatting, linting, and tests
check: fmt lint test
	@echo "Running code quality checks..."
