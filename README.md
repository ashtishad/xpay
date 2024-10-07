# xPay: Digital Wallet

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
- [golang-jwt](https://github.com/golang-jwt/jwt): JSON Web Token handling.
- [golang-migrate](https://github.com/golang-migrate/migrate): Database migrations.
- [viper](https://github.com/spf13/viper): For configuration management. (config: config.yaml)

#### Development Tools
- [swaggo/swag](https://github.com/swaggo/swag): Swagger API documentation.
- [Air](https://github.com/cosmtrek/air): Live reloading. (config: .air.toml)
- [golangci-lint](https://golangci-lint.run/): Linting (config: .golangci.yaml)
