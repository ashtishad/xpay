DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable

# Application commands with Docker Compose
up:
	@docker compose up

down:
	@docker compose down

down-data:
	@docker compose down -v --remove-orphans

# Development Tools
test:
	@go test -v -cover -short ./...

lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

swagger:
	@which swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init

# Git Hooks
setup-hooks:
	@cp scripts/pre-push .git/hooks/
	@chmod +x .git/hooks/pre-push
	@echo "Git hooks set up successfully"

# Database Migration Commands
check_and_install_migrate:
	@which migrate > /dev/null 2>&1 || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: check_and_install_migrate
	@migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down: check_and_install_migrate
	@migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-create: check_and_install_migrate
	@migrate create -ext sql -dir migrations -seq $(name)

# Declare PHONY targets
.PHONY: up down down-data test lint swagger setup-hooks \
        migrate-up migrate-down migrate-create check_and_install_migrate
