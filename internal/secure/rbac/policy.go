package rbac

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed policy.json
var policyFS embed.FS

// Policy represents the RBAC policy structure
type Policy struct {
	// Routes maps resource types to their URL patterns and allowed methods
	// e.g., "users" -> "/api/v1/users" -> "POST" -> "CreateUserWithRole"
	Routes map[string]map[string]map[string]string `json:"routes"`

	// Roles maps role names to their allowed actions and HTTP methods
	// e.g., "admin" -> "CreateUserWithRole" -> ["POST"]
	Roles map[string]map[string][]string `json:"roles"`
}

// LoadPolicy reads and parses the RBAC policy from the embedded policy.json file
func LoadPolicy() (*Policy, error) {
	data, err := policyFS.ReadFile("policy.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read policy.json: %w", err)
	}

	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy data: %w", err)
	}

	return &policy, nil
}
