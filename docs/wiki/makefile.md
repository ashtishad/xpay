# Makefile Guide

This guide explains the Makefile commands used in the XPay project, their purposes, and how to use them.

## Table of Contents
1. [Environment Variables](#environment-variables)
2. [Application Commands](#application-commands)
3. [Setup Git Hooks](#setup-git-hooks)
4. [Docker Compose Commands](#docker-compose-commands)
5. [Database Commands](#database-commands)
6. [Docker Application Commands](#docker-application-commands)

## Environment Variables

```makefile
DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable
```

This variable sets the database connection string for local development. It's used in database migration commands.

## Application Commands

### Run the Application
```makefile
run:
	@go run main.go
```
**Purpose**: Starts the XPay application.
**Usage**: `make run`

### Watch for Changes (Hot Reload)
```makefile
watch:
	@if ! command -v air > /dev/null; then \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air; \
	rm -rf ./tmp
```
**Purpose**: Runs the application with hot reloading using `air`.
**Usage**: `make watch`
**Note**: Automatically installs `air` if not present and cleans up temporary files.

### Run Tests
```makefile
test:
	@go test -v -cover -short ./...
```
**Purpose**: Runs all tests with verbose output, coverage, and in short mode.
**Usage**: `make test`

### Lint Code
```makefile
lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...
```
**Purpose**: Runs the golangci-lint tool for code quality checks.
**Usage**: `make lint`
**Note**: Automatically installs `golangci-lint` if not present.

## Setup Git Hooks

```makefile
setup-hooks:
	@cp scripts/pre-push .git/hooks/
	@chmod +x .git/hooks/pre-push
	@echo "Git hooks set up successfully"
```
**Purpose**: Sets up Git hooks, specifically the pre-push hook.
**Usage**: `make setup-hooks`

## Docker Compose Commands

### Start Docker Compose Services
```makefile
up:
	@docker compose up -d
```
**Purpose**: Starts all services defined in docker-compose.yml in detached mode.
**Usage**: `make up`

### Stop Docker Compose Services
```makefile
down:
	@docker compose down
```
**Purpose**: Stops all running docker-compose services.
**Usage**: `make down`

### Stop and Remove Docker Compose Services and Volumes
```makefile
down-data:
	@docker compose down -v --remove-orphans
```
**Purpose**: Stops services, removes containers, networks, volumes, and orphan containers.
**Usage**: `make down-data`

## Database Commands

### Check and Install Migrate Tool
```makefile
check_and_install_migrate:
	@which migrate > /dev/null 2>&1 || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
**Purpose**: Ensures the `migrate` tool is installed.
**Usage**: This is a helper command used by other migration commands.

### Run Migrations Up
```makefile
migrate-up: check_and_install_migrate
	@migrate -path migrations -database "$(DB_URL)" -verbose up
```
**Purpose**: Applies all pending database migrations.
**Usage**: `make migrate-up`

### Revert Migrations
```makefile
migrate-down: check_and_install_migrate
	@migrate -path migrations -database "$(DB_URL)" -verbose down
```
**Purpose**: Reverts the last applied database migration.
**Usage**: `make migrate-down`

### Create New Migration
```makefile
migrate-create: check_and_install_migrate
	@migrate create -ext sql -dir migrations -seq $(name)
```
**Purpose**: Creates a new migration file.
**Usage**: `make migrate-create name=your_migration_name`

## Docker Application Commands

### Build Docker Image
```makefile
docker-build:
	@docker build -t xpay:latest .
```
**Purpose**: Builds a Docker image for the XPay application.
**Usage**: `make docker-build`

### Run Docker Container
```makefile
docker-run:
	@docker run --name xpay_app --network xpay_network \
		-e APP_ENV="dev" \
		-e DB_URL="postgres://ash:samplepass@postgres:5432/xpay?sslmode=disable&timezone=UTC" \
		-e SERVER_ADDRESS="0.0.0.0:8080" \
		-e GIN_MODE="release" \
		-p 8080:8080 xpay:latest
```
**Purpose**: Runs the XPay application in a Docker container.
**Usage**: `make docker-run`
**Note**: Configures environment variables and network settings for the container.

### Stop and Remove Docker Container
```makefile
docker-stop:
	@docker stop xpay_app || true
	@docker rm xpay_app || true
```
**Purpose**: Stops and removes the XPay Docker container.
**Usage**: `make docker-stop`

### Rebuild and Rerun Docker Container
```makefile
docker-rerun: docker-stop docker-build docker-run
```
**Purpose**: Stops the existing container, rebuilds the image, and runs a new container.
**Usage**: `make docker-rerun`

## Best Practices

1. Always use `make test` and `make lint` before committing changes.
2. Use `make watch` during development for faster iteration.
3. Run `make migrate-up` after pulling new changes to keep your database schema up-to-date.
4. Use `make docker-rerun` when you need to rebuild and restart the application in Docker.
5. Run `make setup-hooks` after cloning the repository to set up Git hooks.

## Troubleshooting

- If you encounter permission issues with Docker commands, ensure your user is part of the Docker group.
- For database connection issues, verify the `DB_URL` in the Makefile and ensure your PostgreSQL server is running.
- If `air` or `golangci-lint` fail to install, check your Go installation and `GOPATH` settings.
- If Git hooks are not working, make sure you've run `make setup-hooks` and that the scripts have execute permissions.
