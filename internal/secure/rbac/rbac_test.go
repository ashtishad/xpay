package rbac

import (
	"testing"

	"github.com/ashtishad/xpay/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	policy, err := LoadPolicy()
	assert.NoError(t, err)
	assert.NotNil(t, policy)

	rbac := New(policy)
	assert.NotNil(t, rbac)
	assert.NotEmpty(t, rbac.routes)
	assert.NotEmpty(t, rbac.permissions)
}

func TestRBAC_HasPermission(t *testing.T) {
	policy, _ := LoadPolicy()
	rbac := New(policy)

	tests := []struct {
		name     string
		role     string
		path     string
		method   string
		expected bool
	}{
		// Admin permissions (all allowed)
		{"Admin Create User", "admin", "/api/v1/users", "POST", true},
		{"Admin Create Wallet", "admin", "/api/v1/users/:user_uuid/wallets", "POST", true},
		{"Admin Get Wallet Balance", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "GET", true},
		{"Admin Update Wallet Status", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/status", "PATCH", true},
		{"Admin Add Card", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "POST", true},
		{"Admin Get Card", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "GET", true},
		{"Admin Update Card", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PATCH", true},
		{"Admin Delete Card", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "DELETE", true},
		{"Admin List Cards", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "GET", true},

		// User permissions
		{"User Create Wallet", "user", "/api/v1/users/:user_uuid/wallets", "POST", true},
		{"User Get Wallet Balance", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "GET", true},
		{"User Update Wallet Status", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/status", "PATCH", true},
		{"User Add Card", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "POST", true},
		{"User Get Card", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "GET", true},
		{"User Update Card", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PATCH", true},
		{"User Delete Card", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "DELETE", true},
		{"User List Cards", "user", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "GET", true},
		{"User Create User (Denied)", "user", "/api/v1/users", "POST", false},

		// Agent permissions
		{"Agent Create User", "agent", "/api/v1/users", "POST", true},
		{"Agent Get Wallet Balance", "agent", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "GET", true},
		{"Agent Update Wallet Status", "agent", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/status", "PATCH", true},
		{"Agent Get Card", "agent", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "GET", true},
		{"Agent List Cards", "agent", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "GET", true},
		{"Agent Create Wallet (Denied)", "agent", "/api/v1/users/:user_uuid/wallets", "POST", false},
		{"Agent Add Card (Denied)", "agent", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "POST", false},

		// Merchant permissions
		{"Merchant Create Wallet", "merchant", "/api/v1/users/:user_uuid/wallets", "POST", true},
		{"Merchant Get Wallet Balance", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "GET", true},
		{"Merchant Update Wallet Status", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/status", "PATCH", true},
		{"Merchant Add Card", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "POST", true},
		{"Merchant Get Card", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "GET", true},
		{"Merchant Update Card", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PATCH", true},
		{"Merchant Delete Card", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "DELETE", true},
		{"Merchant List Cards", "merchant", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "GET", true},
		{"Merchant Create User (Denied)", "merchant", "/api/v1/users", "POST", false},

		// Invalid routes (all denied)
		{"Invalid User Route", "admin", "/api/v1/users/:user_uuid", "GET", false},
		{"Invalid Wallet Route", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid", "GET", false},
		{"Invalid Card Route", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid/invalid", "GET", false},

		// Non-existent Methods (all denied)
		{"Non-existent Method for User Creation", "admin", "/api/v1/users", "PUT", false},
		{"Non-existent Method for Wallet Balance", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "POST", false},
		{"Non-existent Method for Card Update", "admin", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PUT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rbac.HasPermission(tt.role, tt.path, tt.method)
			assert.Equal(t, tt.expected, result, "Role: %s, Path: %s, Method: %s", tt.role, tt.path, tt.method)
		})
	}
}

func TestGetRouteName(t *testing.T) {
	policy, _ := LoadPolicy()
	rbac := New(policy)

	tests := []struct {
		name     string
		path     string
		method   string
		expected string
	}{
		// User Management
		{"Create User with Role", "/api/v1/users", "POST", "CreateUserWithRole"},

		// Wallet Management
		{"Create Wallet", "/api/v1/users/:user_uuid/wallets", "POST", "CreateWallet"},
		{"Get Wallet Balance", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "GET", "GetWalletBalance"},
		{"Update Wallet Status", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/status", "PATCH", "UpdateWalletStatus"},

		// Card Management
		{"Add Card to Wallet", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "POST", "AddCardToWallet"},
		{"Get Card", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "GET", "GetCard"},
		{"Update Card", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PATCH", "UpdateCard"},
		{"Delete Card", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "DELETE", "DeleteCard"},
		{"List Cards", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards", "GET", "ListCards"},

		// Invalid Routes
		{"Invalid User Route", "/api/v1/users/:user_uuid", "GET", ""},
		{"Invalid Wallet Route", "/api/v1/users/:user_uuid/wallets/:wallet_uuid", "GET", ""},
		{"Invalid Card Route", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid/invalid", "GET", ""},

		// Non-existent Methods
		{"Non-existent Method for User Creation", "/api/v1/users", "PUT", ""},
		{"Non-existent Method for Wallet Balance", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/balance", "POST", ""},
		{"Non-existent Method for Card Update", "/api/v1/users/:user_uuid/wallets/:wallet_uuid/cards/:card_uuid", "PUT", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRouteName(rbac, tt.path, tt.method)
			assert.Equal(t, tt.expected, result, "Path: %s, Method: %s", tt.path, tt.method)
			if result != tt.expected {
				t.Logf("Routes: %+v", rbac.routes)
			}
		})
	}
}

func TestCanCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		creatorRole string
		newUserRole string
		canCreate   bool
	}{
		{"Admin Create Admin", domain.UserRoleAdmin, domain.UserRoleAdmin, true},
		{"Admin Create User", domain.UserRoleAdmin, domain.UserRoleUser, true},
		{"Admin Create Agent", domain.UserRoleAdmin, domain.UserRoleAgent, true},
		{"Admin Create Merchant", domain.UserRoleAdmin, domain.UserRoleMerchant, true},
		{"Agent Create User", domain.UserRoleAgent, domain.UserRoleUser, true},
		{"Agent Create Merchant", domain.UserRoleAgent, domain.UserRoleMerchant, true},
		{"Agent Create Admin", domain.UserRoleAgent, domain.UserRoleAdmin, false},
		{"Agent Create Agent", domain.UserRoleAgent, domain.UserRoleAgent, false},
		{"User Create Any", domain.UserRoleUser, domain.UserRoleUser, false},
		{"Merchant Create Any", domain.UserRoleMerchant, domain.UserRoleUser, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanCreateUser(tt.creatorRole, tt.newUserRole)
			assert.Equal(t, tt.canCreate, result)
		})
	}
}

func TestLoadPolicy(t *testing.T) {
	policy, err := LoadPolicy()
	assert.NoError(t, err)
	assert.NotNil(t, policy)

	// Check if all expected roles are present
	expectedRoles := []string{"admin", "user", "agent", "merchant"}
	for _, role := range expectedRoles {
		assert.Contains(t, policy.Roles, role)
	}

	// Check if all expected route categories are present
	expectedCategories := []string{"users", "wallets", "cards"}
	for _, category := range expectedCategories {
		assert.Contains(t, policy.Routes, category)
	}

	// Check a specific permission
	assert.Contains(t, policy.Roles["admin"], "CreateUserWithRole")
	assert.Contains(t, policy.Roles["admin"]["CreateUserWithRole"], "POST")
}
