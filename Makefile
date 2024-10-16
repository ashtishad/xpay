# Application commands with Docker Compose
up:
	@docker compose up

down:
	@docker compose down

down-data:
	@docker compose down -v --remove-orphans

# Setup Environments
setup-prod-env:
	chmod +x scripts/setup-prod-env.sh
	@./scripts/setup-prod-env.sh

setup-dev-env:
	chmod +x scripts/setup-dev-env.sh
	@./scripts/setup-dev-env.sh

# Development Tools
test:
	@go test -v -cover -short ./...

lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# Declare PHONY targets
.PHONY: up down down-data setup-prod-env setup-dev-env test lint
