package rbac

import (
	"strings"

	"github.com/ashtishad/xpay/internal/domain"
)

// RBAC handles role-based access control permissions
type RBAC struct {
	// routes stores the policy routes for quick lookup
	// e.g., "users" -> "/api/v1/users" -> "POST" -> "CreateUserWithRole"
	routes map[string]map[string]map[string]string

	// permissions stores preprocessed role permissions for efficient checking
	// e.g., "admin" -> "CreateUserWithRole" -> "POST" -> true
	permissions map[string]map[string]map[string]bool
}

// New creates a new RBAC instance from a Policy
func New(policy *Policy) *RBAC {
	rbac := &RBAC{
		routes:      policy.Routes,
		permissions: make(map[string]map[string]map[string]bool),
	}

	for role, actions := range policy.Roles {
		rbac.permissions[role] = make(map[string]map[string]bool)
		for actionName, methods := range actions {
			rbac.permissions[role][actionName] = make(map[string]bool)
			for _, method := range methods {
				rbac.permissions[role][actionName][method] = true
			}
		}
	}

	return rbac
}

// HasPermission checks if a role has permission for a given path and method
// Used in the Auth Middleware
func (r *RBAC) HasPermission(role, path, method string) bool {
	if role == "admin" {
		return true
	}

	routeName := getRouteName(r, path, method)
	if routeName == "" {
		return false
	}

	rolePerm, ok := r.permissions[role]
	if !ok {
		return false
	}

	routePerm, ok := rolePerm[routeName]
	if !ok {
		return false
	}

	return routePerm[method]
}

// getRouteName resolves the route name from a path and method
func getRouteName(rbac *RBAC, path, method string) string {
	parts := strings.Split(path, "/")
	for i := len(parts); i > 0; i-- {
		template := strings.Join(parts[:i], "/")
		for _, routes := range rbac.routes {
			if routeActions, ok := routes[template]; ok {
				if routeName, ok := routeActions[method]; ok {
					return routeName
				}
			}
		}

		if i > 0 {
			parts[i-1] = ":param"
		}
	}

	return ""
}

// CanCreateUser checks if a given role can create a user with a specific role
// Helpful in Create User With Role Endpoint
func CanCreateUser(creatorRole, newUserRole string) bool {
	switch creatorRole {
	case domain.UserRoleAdmin:
		return true
	case domain.UserRoleAgent:
		return newUserRole == domain.UserRoleUser || newUserRole == domain.UserRoleMerchant
	default:
		return false
	}
}
