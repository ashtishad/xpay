# Makefile Guide

This guide explains how to use the Makefile commands in the XPay project for both production and development environments.

## Prerequisites

- Docker Desktop installed and running
- Git installed

## Production Environment

Use this setup if you want to run the project and interact with the API without making code changes.

### Setup and Run

1. Clone the repository:
   ```
   git clone git@github.com:ashtishad/xpay.git && cd xpay
   ```

2. Set up the production environment:
   ```
   make setup-prod-env
   ```

3. Start the application:
   ```
   make up
   ```

### Management Commands

- Stop the application:
  ```
  make down
  ```

- Stop and remove all data:
  ```
  make down-data
  ```

- Run tests:
  ```
  make test
  ```

## Development Environment

Use this setup if you intend to modify the code and contribute to the project.

### Setup and Run

1. Clone the repository (if not done already):
   ```
   git clone git@github.com:ashtishad/xpay.git && cd xpay
   ```

2. Set up the development environment:
   ```
   make setup-dev-env
   ```

3. Start the database:
   ```
   make up
   ```

4. Run the application (choose one):
   - With live reload:
     ```
     make watch
     ```
   - Without live reload:
     ```
     make run
     ```

### Development Commands

- Run tests:
  ```
  make test
  ```

- Run linter:
  ```
  make lint
  ```

- Generate Swagger documentation:
  ```
  make swagger
  ```

- Create a new database migration:
  ```
  make migrate-create name=your_migration_name
  ```

- Apply database migrations:
  ```
  make migrate-up
  ```

- Revert last database migration:
  ```
  make migrate-down
  ```

## Environment Variables

The `DB_URL` environment variable is used for database connections. It's defined in your `config.yaml` file or set in the Docker Compose configuration.

- Local development:
  ```
  DB_URL=postgres://ash:samplepass@localhost:5432/xpay?sslmode=disable&timezone=UTC
  ```

- Docker environment:
  ```
  DB_URL=postgres://ash:samplepass@postgres:5432/xpay?sslmode=disable&timezone=UTC
  ```

## Troubleshooting

- For database issues, check the `DB_URL` in your `config.yaml` or environment variables.
- Ensure Docker is running for Docker-related commands.
- For development tool issues, verify your Go installation and `GOPATH` settings.

For more detailed information on each command, refer to the comments in the Makefile.
