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

# Declare PHONY targets
.PHONY: up down down-data setup-prod-env setup-dev-env test
