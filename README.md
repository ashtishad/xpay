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


## API Documentation

### Authentication Endpoints

### Register User
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
  - `400 Bad Request`:
    ```json
    {
      "error": "FullName must be at least 3 characters. Email must be a valid email. Password must be at least 8 characters"
    }
    ```
  - `409 Conflict`:
    ```json
    {
      "error": "user with this email already exists"
    }
    ```
  - `500 Internal Server Error`:
    ```json
    {
      "error": "An unexpected error occurred"
    }
    ```

### Login
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Description**: Authenticate a user, verfies password, generates JWT token, sets an HTTP-only cookie and X-Request-Id header.
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
  - `400 Bad Request`:
    ```json
    {
      "error": "Email must be a valid email. Password must be at least 8 characters"
    }
    ```
  - `401 Unauthorized`:
    ```json
    {
      "error": "Invalid credentials"
    }
    ```
  - `404 Not Found`:
    ```json
    {
      "error": "user not found"
    }
    ```
  - `500 Internal Server Error`:
    ```json
    {
      "error": "An unexpected error occurred"
    }
    ```



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
      - `400 Bad Request`:
        ```json
        {
          "error": "Invalid currency"
        }
        ```
      - `401 Unauthorized`:
        ```json
        {
          "error": "Authentication required"
        }
        ```
      - `403 Forbidden`:
        ```json
        {
          "error": "You can only create a wallet for yourself"
        }
        ```
      - `409 Conflict`:
        ```json
        {
          "error": "User already has a wallet for this currency"
        }
        ```
      - `500 Internal Server Error`:
        ```json
        {
          "error": "An unexpected error occurred"
        }
        ```

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
      - `401 Unauthorized`:
        ```json
        {
          "error": "Authentication required"
        }
        ```
      - `403 Forbidden`:
        ```json
        {
          "error": "You can only access your own wallet"
        }
        ```
      - `404 Not Found`:
        ```json
        {
          "error": "Wallet not found"
        }
        ```
      - `500 Internal Server Error`:
        ```json
        {
          "error": "An unexpected error occurred"
        }
        ```

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
      - `400 Bad Request`:
        ```json
        {
          "error": "Invalid status"
        }
        ```
      - `401 Unauthorized`:
        ```json
        {
          "error": "Authentication required"
        }
        ```
      - `403 Forbidden`:
        ```json
        {
          "error": "You can only update your own wallet"
        }
        ```
      - `404 Not Found`:
        ```json
        {
          "error": "Wallet not found"
        }
        ```
      - `500 Internal Server Error`:
        ```json
        {
          "error": "An unexpected error occurred"
        }
        ```
