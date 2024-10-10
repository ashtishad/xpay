package secure

import (
	"net/http"
	"strings"

	"github.com/ashtishad/xpay/internal/domain"
)

// Action represents the type of operation being performed
type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// Resource represents the type of resource being accessed
type Resource string

const (
	ResourceUser    Resource = "user"
	ResourceWallet  Resource = "wallet"
	ResourceCard    Resource = "card"
	ResourcePayment Resource = "payment"
)

// Policy defines the permissions for each role
var Policy = map[string]map[Resource][]Action{
	domain.UserRoleAdmin: {
		ResourceUser:    {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourceWallet:  {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourceCard:    {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourcePayment: {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
	},
	domain.UserRoleAgent: {
		ResourceUser:    {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourceWallet:  {ActionRead, ActionUpdate},
		ResourceCard:    {ActionRead},
		ResourcePayment: {ActionCreate, ActionRead},
	},
	domain.UserRoleMerchant: {
		ResourceWallet:  {ActionCreate, ActionRead, ActionUpdate},
		ResourceCard:    {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourcePayment: {ActionCreate, ActionRead},
	},
	domain.UserRoleUser: {
		ResourceWallet:  {ActionCreate, ActionRead, ActionUpdate},
		ResourceCard:    {ActionCreate, ActionRead, ActionUpdate, ActionDelete},
		ResourcePayment: {ActionCreate, ActionRead},
	},
}

// CanPerform checks if a given role can perform an action on a resource
func CanPerform(role string, resource Resource, action Action) bool {
	if actions, ok := Policy[role][resource]; ok {
		for _, a := range actions {
			if a == action {
				return true
			}
		}
	}
	return false
}

// CanCreateUser checks if a given role can create a user with a specific role
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

// CheckRBAC determines if the given role can perform the action on the resource for the given path and method
func CheckRBAC(path, method, userRole string) bool {
	resource, action := determineResourceAndAction(path, method)
	if resource == "" || action == "" {
		return false
	}

	return CanPerform(userRole, resource, action)
}

// determineResourceAndAction identifies the resource and action based on the request path and method
func determineResourceAndAction(path, method string) (Resource, Action) {
	switch {
	case path == "/api/v1/users":
		return ResourceUser, ActionCreate
	case strings.HasPrefix(path, "/api/v1/users/:user_uuid/wallets"):
		return handleWalletPath(path, method)
	case strings.HasPrefix(path, "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards"):
		return handleCardPath(path, method)
	default:
		return "", ""
	}
}

// handleWalletPath determines the resource and action for wallet-related paths
func handleWalletPath(path, method string) (Resource, Action) {
	switch {
	case strings.HasSuffix(path, "/balance"):
		if method == http.MethodGet {
			return ResourceWallet, ActionRead
		}

	case strings.HasSuffix(path, "/status"):
		if method == http.MethodPatch {
			return ResourceWallet, ActionUpdate
		}

	default:
		switch method {
		case http.MethodPost:
			return ResourceWallet, ActionCreate
		case http.MethodGet:
			return ResourceWallet, ActionRead
		case http.MethodPatch, http.MethodPut:
			return ResourceWallet, ActionUpdate
		case http.MethodDelete:
			return ResourceWallet, ActionDelete
		}
	}

	return "", ""
}

// handleCardPath determines the resource and action for card-related paths
func handleCardPath(path, method string) (Resource, Action) {
	if strings.HasSuffix(path, "/cards") {
		if method == http.MethodGet {
			return ResourceCard, ActionRead
		}

		return ResourceCard, ActionCreate
	}

	switch method {
	case http.MethodGet:
		return ResourceCard, ActionRead
	case http.MethodPatch:
		return ResourceCard, ActionUpdate
	case http.MethodDelete:
		return ResourceCard, ActionDelete
	default:
		return "", ""
	}
}
