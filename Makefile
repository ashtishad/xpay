DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable

# Application commands
run:
	@go run main.go

watch:
	@if ! command -v air > /dev/null; then \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air; \
	rm -rf ./tmp

test:
	@go test -race ./...

lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

swagger:
	@which swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init


# Setup Environments
setup-prod-env:
	chmod +x scripts/setup-prod-env.sh
	@./scripts/setup-prod-env.sh

setup-dev-env:
	chmod +x scripts/setup-dev-env.sh
	@./scripts/setup-dev-env.sh

# Docker compose commands
up:
	@docker compose up -d

down:
	@docker compose down

down-data:
	@docker compose down -v --remove-orphans

# Database commands
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

.PHONY: run watch test lint up down down-data docker-build docker-run docker-stop docker-rerun \
        migrate-up migrate-down migrate-create
