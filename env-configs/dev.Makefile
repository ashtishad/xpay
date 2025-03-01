GOPATH := $(shell go env GOPATH)
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint
DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable

# Application commands
run:
	@go run main.go

watch:
	@air

test:
	@go test -v -cover -short ./...

lint:
	@$(GOLANGCI_LINT) run ./...

swagger:
	@swag init

# Setup Environments
setup-prod-env:
	@./scripts/setup-prod-env.sh

setup-dev-env:
	@./scripts/setup-dev-env.sh

# Docker compose commands
up:
	@docker compose up -d

down:
	@docker compose down

down-data:
	@docker compose down -v --remove-orphans

# Database Migration Commands
migrate-up:
	@migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	@migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-create:
	@migrate create -ext sql -dir migrations -seq $(name)

# Docker application commands
docker-build:
	@docker build -t xpay:latest .

docker-run:
	@docker run --name xpay_app --network xpay_network \
		-e APP_ENV="dev" \
		-e DB_URL="postgres://ash:samplepass@postgres:5432/xpay?sslmode=disable&timezone=UTC" \
		-e SERVER_ADDRESS="0.0.0.0:8080" \
		-e GIN_MODE="release" \
		-p 8080:8080 xpay:latest

docker-stop:
	@docker stop xpay_app || true
	@docker rm xpay_app || true

docker-rerun: docker-stop docker-build docker-run

.PHONY: run watch test lint swagger up down down-data \
        migrate-up migrate-down migrate-create \
        docker-build docker-run docker-stop docker-rerun setup-prod-env setup-dev-env
