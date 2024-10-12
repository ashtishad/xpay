# xPay: Digital Wallet
<a name="top"></a>

<p align="left">
  <a href="https://www.linkedin.com/in/ashef/">
    <img src="https://skillicons.dev/icons?i=go,docker,kubernetes,postgresql,aws,kafka&theme=light" alt="Skills" />
  </a>
</p>

## Table of Contents
- [Quick Start](#quick-start)
- [Tools/Libraries Used](#toolslibraries-used)
- [Progress](#progress)
- [Architecture and Request Flow](#architecture-and-request-flow)
- [Directory Structure](#directory-structure)
- [Wiki](#wiki)
- [API Documentation](#api-documentation)

## Quick Start

1. Clone and set up configuration:
   ```
   git clone git@github.com:ashtishad/xpay.git && cd xpay
   cp local-dev/config.yaml.example config.yaml
   ```

2. Run and manage the application:
   ```
   make up            # Start the application
   make down          # Stop the application
   make down-data     # Stop and remove data
   ```

3. Development commands (run locally):
   ```
   make test          # Run tests
   make lint          # Run linter
   make swagger       # Generate Swagger docs
   make migrate-create name=your_migration_name  # Create a migration
   ```

**For more details on available commands, refer to the `Makefile`**


## Tools/Libraries Used

#### Used in the Core API
- [Gin](https://github.com/gin-gonic/gin): HTTP routing, middleware and cors.
- [pgx](https://github.com/jackc/pgx): Database driver and connection pooling, using standard *sql.DB handle.
- [golang-migrate](https://github.com/golang-migrate/migrate): Database migrations.
- [golang-jwt](https://github.com/golang-jwt/jwt/): JSON Web Token handling.
- [viper](https://github.com/spf13/viper): For configuration management. (config: config.yaml)
- [swaggo/swag](https://github.com/swaggo/swag): Swagger API documentation (usage: make swagger).
- [golangci-lint](https://golangci-lint.run/): Linting (config: .golangci.yaml)

</details>

## Progress

✅ Implemented | 🔄 In Progress/Planned

| Area | Features and Best Practices | Status |
|------|------------------------------|--------|
| API Design & Architecture | • Domain Driven Design, Clean Architecure <br>• RESTful API<br>• Event streaming with Apache Kafka<br>• OpenAPI 2.0 specifications | ✅<br>✅<br>🔄<br>✅ |
| Security | • JWT-ES256 with ECDSA asymmetric key pairs<br>• AES-256-GCM for card data encryption<br>• SQL injection prevention with parameterized sql queries<br>• Role based access control (RBAC) <br>• DTO for controlled data to the client<br>• User input and query param validation<br>• IP-Based Rate limiting with Token Bucket algorithm | ✅<br>✅<br>✅<br>✅<br>✅<br>✅<br>✅ |
| Database | • ACID transactions with appropriate isolation levels<br>• Raw SQL for performance<br>• Connection pooling with pgx, exposing standard *sql.DB<br>• Optimized indexing and unique constraints<br>• Version-controlled schema changes with migrations | ✅<br>✅<br>✅<br>✅<br>✅ |
| Core Operations & Observability | • Custom AppError interface for error handling<br>• Centralized configuration management with Viper<br>• Structured logging with slog<br>• Context with timeout for each request <br>• Comprehensive test coverage<br>• Code quality with golangci-lint | ✅<br>✅<br>✅<br>✅<br>🔄<br>✅ |
| Payment Gateways | • Idempotent payment processing<br>• Stripe integration<br>• PayPal integration<br>• Webhook handling for asynchronous events | 🔄<br>🔄<br>🔄<br>🔄 |
| Deployment & Monitoring | • Multi-stage Docker builds for minimal image size <br>• GitHub Actions CI pipeline<br>• AWS RDS with PostgreSQL<br>• ECS Fargate for serverless container deployment<br>• Prometheus metrics and Grafana dashboards | ✅<br>✅<br>🔄<br>🔄<br>🔄 |

<a href="#top">Back to Top</a>

## Architecture and Request Flow:

The project follows domain-driven design, loosely coupled clean architecture, suited for large codebase.

```

┌─────────┐    JSON    ┌───────────────┐   Domain   ┌─────────────┐
│ Client  ├───────────►│   Handlers    ├───────────►│ Repositories│
└─────────┘    (DTO)   └───────────────┘   Models   └─────────────┘
                              │                           │
                              │                           │
                              │          Domain           │
                              │          Models           │
                              │                           │
┌─────────┐    JSON    ┌───────────────┐   Domain   ┌─────────────┐
│ Client   ◄───────────┤   Handlers     ◄───────────┤Repositories │
└─────────┘    (DTO)   └───────────────┘   Models   └─────────────┘

```
Example: Create a wallet

wallet.go (model) -> wallet_repository.go -> wallet_handlers.go (using DTOs)

1. Domain Models: `internal/domain/*.go`
2. Repositories: `internal/domain/*_repository.go`
3. HTTP Handlers: `internal/server/handlers/*.go`
4. DTOs: `internal/dto/*.go`
5. Routes: `internal/server/routes/*.go`

<a href="#top">Back to Top</a>

## Directory Structure

<details>
<summary>Click to expand Directory Structure</summary>

command: `tree -a -I '.git|.DS_Store|.gitignore|.idea|.vscode|docs'`

```bash
├── .github
│   └── workflows
│       └── test.yaml                 # CI/CD pipeline for running tests
├── internal
│   ├── domain
│   │   ├── card.go                   # Card domain model
│   │   ├── card_repository.go        # Card repository interface, database interactions
│   │   ├── helpers.go                # Domain-specific helper functions
│   │   ├── user.go                   # User domain model
│   │   ├── user_repository.go        # User repository interface, database interactions
│   │   ├── wallet.go                 # Wallet domain model
│   │   └── wallet_repository.go      # Wallet repository interface, database interactions
│   ├── dto
│   │   ├── auth.go                   # Authentication-related DTOs/REST API Request Response Structurers
│   │   ├── card.go                   # Card-related DTOs
│   │   ├── shared.go                 # Shared DTO structures
│   │   └── wallet.go                 # User-related DTOs
│   │   └── wallet.go                 # Wallet-related DTOs
│   ├── secure
│   │   ├── card_aes.go               # Card AES-256 with GCM mode, Validate, Encrypt and Decrypt
│   │   ├── jwt.go                    # JWT token handling, generate and validate tokens
│   │   ├── rbac.go                   # Role based access control (RBAC) policies
│   │   ├── password.go               # Password hashing and verification with bcrypt
│   │   └── password_test.go          # Password utility tests
│   ├── server
│   │   ├── handlers
│   │   │   ├── auth.go               # Login, Register handlers
│   │   │   ├── card.go               # Card http handlers
│   │   │   ├── helpers.go            # Handlers helper functions
│   │   │   └── user.go               # User HTTP handlers
│   │   │   └── wallet.go             # Wallet HTTP handlers
│   │   ├── middlewares
│   │   │   ├── auth.go               # Auth middleware (Validate token, Set Authorized user in req context)
│   │   │   ├── cors.go               # CORS middleware
│   │   │   ├── gin_logger.go         # Custom Logging middleware for gin
│   │   │   ├── middlewares.go        # Core Middleware setup
│   │   │   └── rate_limiter.go       # IP-Based rate limiter with token bucket algorithm
│   │   │   └── request_id.go         # Request ID middleware, sets X-Request-ID header
│   │   ├── routes
│   │   │   ├── auth.go               # Authentication routes
│   │   │   ├── card.go               # Card routes
│   │   │   ├── routes.go             # Core routes setup
│   │   │   └── user.go               # User  routes
│   │   │   └── wallet.go             # Wallet routes
│   │   └── server.go                  # HTTP server setup with gin
│   ├── infra
│   │   ├── postgres
│   │   │   ├── postgres_connection.go    # Postgres connection setup with pgx, returns *sql.DB
│   │   │   └── postgres_migrations.go    # Database migration handling with golang-migrate/v4
│   │   ├── kafka
│   │   │   └── sample.md                 # Placeholder for Kafka integration
│   ├── common
│   │   ├── app_errs.go               # Custom error types
│   │   ├── config.go                 # Configuration management
│   │   ├── constants.go              # Global constants
│   │   ├── context_keys.go           # Context key definitions
│   │   ├── custom_err_messages.go    # Error message definitions
│   │   ├── slog_config.go            # Structured logging configuration
│   │   ├── timeouts.go               # Context timeout constants
├── migrations
│   ├── 000001_create_users_table.down.sql   # User table rollback
│   ├── 000001_create_users_table.up.sql     # User table creation
│   ├── 000002_create_wallets_table.down.sql # Wallet table rollback
│   └── 000002_create_wallets_table.up.sql   # Wallet table creation
│   ├── 000003_create_cards_table.down.sql   # Cards table rollback
│   └── 000003_create_cards_table.up.sql     # Cards table creation
├── scripts
│   └── pre-push                      # Git pre-push hook (ensures run tests and lint before every push)
├── local-dev
│   └── config.yaml.example           # Example configuration file (place it to project root as `config.yaml`)
├── config.yaml                       # Application configuration
├── main.go                           # Application entry point
├── Makefile                          # Development commands and shortcuts
├── Dockerfile                        # Docker file with multi stage builds
├── .dockerignore                     # Directories to ignore in the Docker builds
├── README.md                         # Project documentation
├── compose.yaml                      # Docker Compose configuration
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksums
└── .air.toml                         # Live reload configuration with air
```

</details>

<a href="#top">Back to Top</a>

## Wiki

For detailed information on various aspects of the project, refer to the following guides:

<details>
<summary>Click to expand Wiki</summary>

- [Makefile Commands](https://github.com/ashtishad/xpay/blob/main/docs/wiki/makefile.md): Comprehensive guide to all Make commands used in development and deployment.
- [Configuration Management](https://github.com/ashtishad/xpay/blob/main/docs/wiki/config.md): Learn how to manage application configuration using Viper.
- [Dockerfile Guide](https://github.com/ashtishad/xpay/blob/main/docs/wiki/dockerfile.md): Instructions for building and running the XPay application in Docker.
- [Generating Secrets](https://github.com/ashtishad/xpay/blob/main/docs/wiki/generating_secrets.md): Procedures for generating and managing cryptographic keys and secrets.
- [GitHub Actions Test Workflow](https://github.com/ashtishad/xpay/blob/main/docs/wiki/github_actions_test_workflow.md): Understanding the CI/CD pipeline setup using GitHub Actions.
- [Linter Configuration](https://github.com/ashtishad/xpay/blob/main/docs/wiki/linter_config.md): Explanation of golangci-lint setup and usage in the project.
- [Configuration and Key Management in Production](https://github.com/ashtishad/xpay/blob/main/docs/wiki/configuration_key_management_in_production.md): Best practices for managing configs and secrets in production environments.
- [Zed/VSCode Shortcuts](https://github.com/ashtishad/xpay/blob/main/docs/wiki/zed_vscode_shortcuts.md): Helpful keyboard shortcuts for efficient coding in Zed or VSCode editors.

</details>

## API Documentation

<details>
<summary>Click to expand API Documentation</summary>

### Authentication Endpoints

#### Register User
- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Description**: Registers a new user with hashed password, generates JWT tokens, sets an HTTP-only cookie and X-Request-Id header.
- **Access**: Public
- **Request Body**:
  ```json
  {
    "fullName": "John Doe",
    "email": "someone@example.com",
    "password": "samplepass"
  }
  ```
- **Success Response**: `201 Created`
- **Error Responses**: `400 Bad Request`, `409 Conflict`, `500 Internal Server Error`

#### Login
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Description**: Authenticate a user, verifies password, generates JWT token, sets an HTTP-only cookie and X-Request-Id header.
- **Access**: Public
- **Request Body**:
  ```json
  {
    "email": "someone@example.com",
    "password": "samplepass"
  }
  ```
- **Success Response**: `200 OK`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `404 Not Found`, `500 Internal Server Error`

### User Management Endpoints

#### Create User with Specific Role
- **URL**: `/api/v1/users`
- **Method**: `POST`
- **Description**: Creates a new user with a specific role.
- **Access**: Admin (can create any role), Agent (can create user or merchant roles)
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "fullName": "Keanu Reeves",
    "email": "keanu@example.com",
    "password": "keanupass",
    "role": "admin"
  }
  ```
- **Success Response**: `201 Created`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `409 Conflict`, `500 Internal Server Error`

### Wallet Endpoints

#### Create a New Wallet
- **URL**: `/api/v1/users/{user_uuid}/wallets`
- **Method**: `POST`
- **Access**: Admin, Merchant, User
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "currency": "USD"
  }
  ```
- **Success Response**: `201 Created`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `409 Conflict`, `500 Internal Server Error`

#### Get Wallet Balance
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/balance`
- **Method**: `GET`
- **Access**: Admin, Agent, Merchant, User (own wallet only)
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
- **Error Responses**: `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `500 Internal Server Error`

#### Update Wallet Status
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/status`
- **Method**: `PATCH`
- **Access**: Admin, Agent, Merchant, User (own wallet only)
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "status": "inactive"
  }
  ```
- **Success Response**: `200 OK`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `500 Internal Server Error`

### Card Endpoints

#### Add a New Card to Wallet
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards`
- **Method**: `POST`
- **Access**: Admin, Merchant, User (own wallet only)
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "cardNumber": "4111111111111111",
    "provider": "visa",
    "type": "credit",
    "expiryDate": "12/25",
    "cvv": "123"
  }
  ```
- **Success Response**: `201 Created`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `409 Conflict`, `500 Internal Server Error`

#### Get Card Details
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `GET`
- **Access**: Admin, Agent (read-only), Merchant, User (own cards only)
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
- **Error Responses**: `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `500 Internal Server Error`

#### Update Card Details
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `PATCH`
- **Access**: Admin, Merchant, User (own cards only)
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "expiryDate": "12/26",
    "status": "inactive"
  }
  ```
- **Success Response**: `200 OK`
- **Error Responses**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `500 Internal Server Error`

#### Delete Card
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `DELETE`
- **Access**: Admin, Merchant, User (own cards only)
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
- **Error Responses**: `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `500 Internal Server Error`

#### List Cards
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards`
- **Method**: `GET`
- **Access**: Admin, Agent (read-only), Merchant, User (own wallet only)
- **Authentication**: Required (Bearer Token)
- **Query Parameters**:
  - `provider` (optional): Filter by card provider
  - `status` (optional): Filter by card status
- **Success Response**: `200 OK`
- **Error Responses**: `401 Unauthorized`, `403 Forbidden`, `500 Internal Server Error`

</details>

[Back to Top](#top)
