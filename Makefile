GOPATH := $(shell go env GOPATH)
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint

# Application commands with Docker Compose
up:
	@docker compose up

down:
	@docker compose down

down-data:
	@docker compose down -v --remove-orphans

# Setup Environments
setup-prod-env:
	@./scripts/setup-prod-env.sh

setup-dev-env:
	@./scripts/setup-dev-env.sh

# Development Tools
test:
	@go test -v -cover -short ./...

lint:
	@$(GOLANGCI_LINT) run ./...

.PHONY: up down down-data setup-prod-env setup-dev-env test lint
