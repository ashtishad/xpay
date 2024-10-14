# Implementing Role-Based Access Control (RBAC) in XPay

This guide outlines the implementation of RBAC across the XPay project, focusing on the core components and their integration.

## 1. RBAC Package

Location: `internal/secure/rbac/`

Purpose: Define and enforce access control policies.

Key Components:
- `policy.json`: Defines routes and role permissions.
- `policy.go`: Loads and parses the RBAC policy.
- `rbac.go`: Implements RBAC logic and permission checks.

Implementation:
```go
// rbac.go
type RBAC struct {
    routes      map[string]map[string]map[string]string
    permissions map[string]map[string]map[string]bool
}

func (r *RBAC) HasPermission(role, path, method string) bool {
    // Implementation details
}

func getRouteName(rbac *RBAC, path, method string) string {
    // Implementation details
}
```

## 2. Auth Middleware

Location: `internal/server/middlewares/auth.go`

Purpose: Authenticate requests and enforce RBAC policies.

Implementation:
```go
func AuthMiddleware(userRepo domain.UserRepository, jwtPublicKey *ecdsa.PublicKey, rbac *rbac.RBAC) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract and validate JWT token
        // 2. Fetch user details
        // 3. Check RBAC permissions
        if !rbac.HasPermission(user.Role, c.FullPath(), c.Request.Method) {
            c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Access denied"})
            c.Abort()
            return
        }
        // 4. Set user in context and continue
    }
}
```

## 3. Routes Package and Server Initialization

Location:
- `internal/server/routes/`
- `internal/server/server.go`

Purpose: Define API routes and integrate RBAC with the server setup.

Implementation:
```go
// routes/routes.go
func InitRoutes(rg *gin.RouterGroup, db *sql.DB, config *common.AppConfig, jm *secure.JWTManager, cardEncryptor *secure.CardEncryptor, rbac *rbac.RBAC) {
    // Initialize repositories
    // Set up public routes
    // Set up authenticated routes with RBAC
    authGroup := rg.Group("/users")
    authGroup.Use(middlewares.AuthMiddleware(userRepo, jm.GetPublicKey(), rbac))
    // Define protected routes
}

// server.go
func NewServer(ctx context.Context) (*Server, error) {
    // Load configuration
    // Initialize database connection
    // Set up JWT manager and card encryptor
    policy, err := rbac.LoadPolicy()
    if err != nil {
        return nil, fmt.Errorf("failed to load rbac policy: %w", err)
    }
    rbacInstance := rbac.New(policy)
    // Initialize server and routes
    s.setupRoutes(jwtManager, cardEncryptor, rbacInstance)
    return s, nil
}
```

## Integration Flow

1. Server Initialization:
   - Load RBAC policy from `policy.json`
   - Create RBAC instance

2. Route Setup:
   - Define routes with appropriate handlers
   - Apply AuthMiddleware to protected routes

3. Request Processing:
   - AuthMiddleware authenticates the user
   - RBAC checks permissions for the requested route and method
   - Request is allowed or denied based on RBAC decision

## Testing

Comprehensive unit tests should be written for:
- RBAC package (policy loading, permission checks)
- Auth middleware (token validation, RBAC integration)
- Integration tests for protected routes

## Maintenance and Updates

- Regularly review and update `policy.json` as new routes or roles are added
- Ensure RBAC tests are updated alongside policy changes
- Monitor logs for unauthorized access attempts to identify potential security issues or missing permissions

## Further Considerations

- Implement fine-grained permissions beyond role-based access if needed
- Consider caching RBAC decisions for performance in high-traffic scenarios
- Regularly audit RBAC logs and permissions to ensure principle of least privilege
