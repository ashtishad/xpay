# xPay: Digital Wallet
<a name="top"></a>

## Table of Contents
- [Quick Start](#quick-start)
- [Tools/Libraries Used](#toolslibraries-used)
- [Progress](#progress)
- [Architecture and Request Flow](#architecture-and-request-flow)
- [Directory Structure](#directory-structure)
- [Wiki](#wiki)
- [API Documentation](#api-documentation)

## Quick Start

1. Clone your new repository locally with ssh:
   ```
   git clone git@github.com:ashtishad/xpay.git
   ```

2. Copy `config.yaml.example` from the `/local-dev` directory to the project root as `config.yaml`:
   ```
   cp local-dev/config.yaml.example config.yaml
   ```

3. Run `make up` to start the Docker services in the background.

4. Run `make run` to start the application.

5. (Or) Run `make watch` to live reload the application with air.

Refer to **Makefile** for more details on development commands. Example: `make migrate-create name=create-cards-table`

## Tools/Libraries Used

#### Used in the Core API
- [Gin](https://github.com/gin-gonic/gin): HTTP routing and middleware.
- [pgx](https://github.com/jackc/pgx): Database driver and connection pooling, using standard *sql.DB handle.
- [golang-migrate](https://github.com/golang-migrate/migrate): Database migrations.
- [golang-jwt](https://github.com/golang-jwt/jwt/): JSON Web Token handling.
- [viper](https://github.com/spf13/viper): For configuration management. (config: config.yaml)

<details>
<summary>Click to see development tools</summary>

- [swaggo/swag](https://github.com/swaggo/swag): Swagger API documentation.
- [Air](https://github.com/cosmtrek/air): Live reloading. (config: .air.toml)
- [golangci-lint](https://golangci-lint.run/): Linting (config: .golangci.yaml)

</details>

## Progress

✅ Implemented | 🔄 In Progress/Planned

| Area | Features and Best Practices | Status |
|------|------------------------------|--------|
| API Design & Architecture | • Domain Driven Design with clear bounded contexts<br>• RESTful API design<br>• Event streaming with Apache Kafka<br>• OpenAPI 3.0 specifications | ✅<br>✅<br>🔄<br>✅ |
| Security | • ES256 JWT with asymmetric key pairs<br>• AES-256-GCM for card data encryption<br>• SQL injection prevention with parameterized queries<br>• DTO for controlled data to the client<br>• Route protection with Auth middleware<br>• Input and query param validation<br>• Rate limiting with Leaky Bucket algorithm | ✅<br>✅<br>✅<br>✅<br>✅<br>✅<br>🔄 |
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
│   │   ├── common.go                 # Shared DTO structures
│   │   └── wallet.go                 # Wallet-related DTOs
│   ├── secure
│   │   ├── card_aes.go               # Card AES-256 with GCM mode, Validate, Encrypt and Decrypt
│   │   ├── jwt.go                    # JWT token handling, generate and validate tokens
│   │   ├── password.go               # Password hashing and verification with bcrypt
│   │   └── password_test.go          # Password utility tests
│   ├── server
│   │   ├── handlers
│   │   │   ├── auth.go               # Login, Register handlers
│   │   │   ├── auth.go               # Card http handlers
│   │   │   ├── helpers.go            # Handlers helper functions
│   │   │   └── wallet.go             # Wallet HTTP handlers
│   │   ├── middlewares
│   │   │   ├── auth.go               # Auth middleware (Validate token, Set Authorized user in req context)
│   │   │   ├── cors.go               # CORS middleware
│   │   │   ├── gin_logger.go         # Custom Logging middleware for gin
│   │   │   ├── middlewares.go        # Core Middleware setup
│   │   │   └── request_id.go         # Request ID middleware, sets X-Request-ID header
│   │   ├── routes
│   │   │   ├── auth.go               # Authentication routes
│   │   │   ├── auth.go               # Card routes
│   │   │   ├── routes.go             # Core routes setup
│   │   │   └── wallet.go             # Wallet routes
│   │   └── server.go                  # HTTP server setup with gin
│   ├── infra
│   │   ├── docker
│   │   │   └── init-db.sql               # Initial database setup script for docker compose
│   │   ├── postgres
│   │   │   ├── postgres_connection.go    # Postgres connection setup with pgx, returns *sql.DB
│   │   │   └── postgres_migrations.go    # Database migration handling with golang-migrate/v4
│   │   ├── kafka
│   │   │   └── sample.md                 # Placeholder for Kafka integration
│   └── common
│       ├── app_errs.go               # Custom error types
│       ├── config.go                 # Configuration management
│       ├── constants.go              # Global constants
│       ├── context_keys.go           # Context key definitions
│       ├── custom_err_messages.go    # Error message definitions
│       ├── slog_config.go            # Structured logging configuration
│       └── timeouts.go               # Timeout constants
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
- **Request Body**:
  ```json
  {
    "fullName": "John Doe",
    "email": "someone@example.com",
    "password": "samplepass"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "user": {
      "uuid": "92e275af-4803-4929-968c-3feb25e038d3",
      "fullName": "John Doe",
      "email": "someone@example.com",
      "status": "active",
      "role": "user",
      "createdAt": "2024-10-07T06:18:54.980941Z",
      "updatedAt": "2024-10-07T06:18:54.980941Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "FullName must be at least 3 characters. Email must be a valid email. Password must be at least 8 characters"}`
  - `409 Conflict`: `{"error": "user with this email already exists"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Login
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Description**: Authenticate a user, verifies password, generates JWT token, sets an HTTP-only cookie and X-Request-Id header.
- **Request Body**:
  ```json
  {
    "email": "someone@example.com",
    "password": "samplepass"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "user": {
      "uuid": "92e275af-4803-4929-968c-3feb25e038d3",
      "fullName": "John Doe",
      "email": "someone@example.com",
      "status": "active",
      "role": "user",
      "createdAt": "2024-10-07T06:18:54.980941Z",
      "updatedAt": "2024-10-07T06:18:54.980941Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "Email must be a valid email. Password must be at least 8 characters"}`
  - `401 Unauthorized`: `{"error": "Invalid credentials"}`
  - `404 Not Found`: `{"error": "user not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

### Wallet Endpoints

#### Create a New Wallet
- **URL**: `/api/v1/users/{user_uuid}/wallets`
- **Method**: `POST`
- **Description**: Creates a new wallet for the specified user.
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "currency": "USD"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "wallet": {
      "uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "balance": 0,
      "currency": "USD",
      "status": "active",
      "createdAt": "2024-10-07T06:18:54.980941Z",
      "updatedAt": "2024-10-07T06:18:54.980941Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "Invalid currency"}`
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only create a wallet for yourself"}`
  - `409 Conflict`: `{"error": "User already has a wallet for this currency"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Get Wallet Balance
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/balance`
- **Method**: `GET`
- **Description**: Retrieves the balance of a specific wallet for a user.
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
  ```json
  {
    "balance": 1000,
    "currency": "USD"
  }
  ```
- **Error Responses**:
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only access your own wallet"}`
  - `404 Not Found`: `{"error": "Wallet not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Update Wallet Status
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/status`
- **Method**: `PUT`
- **Description**: Updates the status of a specific wallet for a user.
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "status": "inactive"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Wallet status updated successfully"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "Invalid status"}`
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only update your own wallet"}`
  - `404 Not Found`: `{"error": "Wallet not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`


### Card Endpoints

#### Add a New Card to Wallet
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards`
- **Method**: `POST`
- **Description**: Adds a new card to the specified wallet, encrypting sensitive data.
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
  ```json
  {
    "card": {
      "uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "provider": "visa",
      "type": "credit",
      "lastFour": "1111",
      "expiryDate": "12/25",
      "status": "active",
      "createdAt": "2024-10-07T06:18:54.980941Z",
      "updatedAt": "2024-10-07T06:18:54.980941Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "Invalid card details"}`
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only add cards to your own wallet"}`
  - `404 Not Found`: `{"error": "Wallet not found"}`
  - `409 Conflict`: `{"error": "A card of this type and provider already exists"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Get Card Details
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `GET`
- **Description**: Retrieves details of a specific card.
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
  ```json
  {
    "uuid": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "provider": "visa",
    "type": "credit",
    "lastFour": "1111",
    "expiryDate": "12/25",
    "status": "active",
    "createdAt": "2024-10-07T06:18:54.980941Z",
    "updatedAt": "2024-10-07T06:18:54.980941Z"
  }
  ```
- **Error Responses**:
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only access your own cards"}`
  - `404 Not Found`: `{"error": "Card not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Update Card Details
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `PATCH`
- **Description**: Updates the details of a specific card.
- **Authentication**: Required (Bearer Token)
- **Request Body**:
  ```json
  {
    "expiryDate": "12/26",
    "status": "inactive"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Card updated successfully"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: `{"error": "Invalid update details"}`
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only update your own cards"}`
  - `404 Not Found`: `{"error": "Card not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### Delete Card
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards/{card_uuid}`
- **Method**: `DELETE`
- **Description**: Soft deletes a specific card.
- **Authentication**: Required (Bearer Token)
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Card deleted successfully"
  }
  ```
- **Error Responses**:
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only delete your own cards"}`
  - `404 Not Found`: `{"error": "Card not found"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

#### List Cards
- **URL**: `/api/v1/users/{user_uuid}/wallets/{wallet_uuid}/cards`
- **Method**: `GET`
- **Description**: Retrieves a list of cards for a specific wallet.
- **Authentication**: Required (Bearer Token)
- **Query Parameters**:
  - `provider` (optional): Filter by card provider
  - `status` (optional): Filter by card status
- **Success Response**: `200 OK`
  ```json
  {
      "cards": [
          {
              "uuid": "93289d24-1c46-4a05-b92c-2ce2284e6462",
              "provider": "mastercard",
              "type": "credit",
              "lastFour": "4444",
              "expiryDate": "11/26",
              "status": "active",
              "createdAt": "2024-10-10T06:04:31.807741+06:00",
              "updatedAt": "2024-10-10T06:05:59.505754+06:00"
          },
          {
              "uuid": "790dc5de-6d9f-44be-9edb-2579ab8bfb5a",
              "provider": "amex",
              "type": "credit",
              "lastFour": "8431",
              "expiryDate": "12/25",
              "status": "active",
              "createdAt": "2024-10-10T06:04:27.99682+06:00",
              "updatedAt": "2024-10-10T06:04:27.99682+06:00"
          },
          {
              "uuid": "77c9c9b9-4fe6-4d51-af7f-72ae5a181fc6",
              "provider": "visa",
              "type": "credit",
              "lastFour": "1111",
              "expiryDate": "12/25",
              "status": "active",
              "createdAt": "2024-10-10T06:04:20.158023+06:00",
              "updatedAt": "2024-10-10T06:04:20.158023+06:00"
          }
      ]
  }
  ```
- **Error Responses**:
  - `401 Unauthorized`: `{"error": "Authentication required"}`
  - `403 Forbidden`: `{"error": "You can only list cards from your own wallet"}`
  - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`

  ### User Management Endpoints

  #### Create User with Specific Role
  - **URL**: `/api/v1/users`
  - **Method**: `POST`
  - **Description**: Creates a new user with a specific role (admin, agent, or merchant). Only admins can perform this action and an admin can create another admin.
  - **Authentication**: Required (Bearer Token)
  - **Authorization**: Admin only
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
    ```json
    {
      "user": {
          "uuid": "3414e162-62f1-40fa-b66f-e077587b74a4",
          "fullName": "Keanu Reeves",
          "email": "keanu@example.com",
          "status": "active",
          "role": "admin",
          "createdAt": "2024-10-10T17:52:41.777578Z",
          "updatedAt": "2024-10-10T17:52:41.777578Z"
        }
    }
    ```
  - **Error Responses**:
    - `400 Bad Request`: `{"error": "Invalid input. FullName must be at least 3 characters. Email must be valid. Password must be at least 8 characters. Role must be admin, agent, or merchant."}`
    - `401 Unauthorized`: `{"error": "Authentication required"}`
    - `403 Forbidden`: `{"error": "Admin access required"}`
    - `409 Conflict`: `{"error": "User with this email already exists"}`
    - `500 Internal Server Error`: `{"error": "An unexpected error occurred"}`


</details>

[Back to Top](#top)
