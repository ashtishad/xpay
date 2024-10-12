# Makefile Guide

This guide explains the Makefile commands used in the XPay project, their purposes, and how to use them.

## Table of Contents
1. [Environment Variables](#environment-variables)
2. [Standard Mode Commands](#standard-mode-commands)
3. [Live Reload Mode Commands](#live-reload-mode-commands)
4. [Development Tools](#development-tools)
5. [Git Hooks](#git-hooks)
6. [Database Migration Commands](#database-migration-commands)

## Environment Variables

```makefile
DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable
```

This variable sets the database connection string for local development. It's used in database migration commands.

## Standard Mode Commands

### Start Application
```makefile
up:
	@docker compose up --build
```
**Purpose**: Starts all services defined in compose.yaml, building images if necessary.
**Usage**: `make up`

### Stop Application
```makefile
down:
	@docker compose down
```
**Purpose**: Stops all running docker compose services.
**Usage**: `make down`

### Stop and Remove Data
```makefile
down-data:
	@docker compose down -v --remove-orphans
```
**Purpose**: Stops services, removes containers, networks, volumes, and orphan containers.
**Usage**: `make down-data`

## Live Reload Mode Commands

### Start Application with Live Reload
```makefile
watch:
	@docker compose -f compose.yaml -f compose.dev.yaml up --build
```
**Purpose**: Starts the application with live reloading for development.
**Usage**: `make watch`

### Stop Live Reload Application
```makefile
watch-down:
	@docker compose -f compose.yaml -f compose.dev.yaml down
```
**Purpose**: Stops all running docker compose services in live reload mode.
**Usage**: `make watch-down`

### Stop and Remove Data (Live Reload Mode)
```makefile
watch-down-data:
	@docker compose -f compose.yaml -f compose.dev.yaml down -v --remove-orphans
```
**Purpose**: Stops services, removes containers, networks, volumes, and orphan containers in live reload mode.
**Usage**: `make watch-down-data`

## Development Tools

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

### Generate Swagger Documentation
```makefile
swagger:
	@which swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init
```
**Purpose**: Generates Swagger documentation for the API.
**Usage**: `make swagger`
**Note**: Automatically installs the `swag` tool if not present.

## Git Hooks

```makefile
setup-hooks:
	@cp scripts/pre-push .git/hooks/
	@chmod +x .git/hooks/pre-push
	@echo "Git hooks set up successfully"
```
**Purpose**: Sets up Git hooks, specifically the pre-push hook.
**Usage**: `make setup-hooks`

## Database Migration Commands

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

## Best Practices

1. Use `make watch` during development for live reloading.
2. Use `make up` for running the application in standard mode.
3. Always run `make test` and `make lint` before committing changes.
4. Run `make migrate-up` after pulling new changes to keep your database schema up-to-date.
5. Use `make swagger` to update API documentation when endpoints change.
6. Run `make setup-hooks` after cloning the repository to set up Git hooks.

## Troubleshooting

- If you encounter permission issues with Docker commands, ensure your user is part of the Docker group.
- For database connection issues, verify the `DB_URL` in the Makefile and ensure your PostgreSQL server is running.
- If `golangci-lint` or `swag` fail to install, check your Go installation and `GOPATH` settings.
- If Git hooks are not working, make sure you've run `make setup-hooks` and that the scripts have execute permissions.
