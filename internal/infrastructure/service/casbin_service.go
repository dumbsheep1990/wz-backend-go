package service

import (
	"log"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"

	"wz-backend-go/internal/domain/admin/valueobject"
)

// CasbinService provides integration with Casbin for authorization
type CasbinService struct {
	enforcer *casbin.Enforcer
}

// NewCasbinService creates a new CasbinService with the given model and policy paths
func NewCasbinService(modelPath, policyPath string) (*CasbinService, error) {
	// Load model from file or string
	var m model.Model
	var err error
	
	if modelPath != "" {
		m, err = model.NewModelFromFile(modelPath)
	} else {
		// Default RBAC model with resource domains
		m, err = model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")
`)
	}
	if err != nil {
		return nil, err
	}

	// Create adapter
	var adapter fileadapter.Adapter
	if policyPath != "" {
		adapter = fileadapter.NewAdapter(policyPath)
	} else {
		// In-memory adapter if no policy path provided
		adapter = fileadapter.NewAdapter("") 
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return &CasbinService{
		enforcer: enforcer,
	}, nil
}

// GetEnforcer returns the Casbin enforcer instance
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}

// LoadPoliciesFromPermissions loads policies from a list of permissions
func (s *CasbinService) LoadPoliciesFromPermissions(roleID string, permissions []valueobject.Permission) error {
	// Clear existing policies for this role
	err := s.enforcer.RemoveFilteredPolicy(0, roleID)
	if err != nil {
		return err
	}

	// Add new policies
	for _, permission := range permissions {
		parts := strings.Split(permission.Value(), ":")
		if len(parts) != 2 {
			// Skip invalid permissions
			log.Printf("跳过无效权限格式: %s", permission.Value())
			continue
		}

		resource := parts[0]
		action := parts[1]

		// For wildcard permissions like "admin:*"
		if action == "*" {
			// Add a policy with the wildcard
			_, err = s.enforcer.AddPolicy(roleID, "万知", resource, "*")
		} else {
			// Add specific permission
			_, err = s.enforcer.AddPolicy(roleID, "万知", resource, action)
		}

		if err != nil {
			return err
		}
	}

	// Save policies to storage
	return s.enforcer.SavePolicy()
}

// Enforce checks if the role has permission to perform the action on the resource
func (s *CasbinService) Enforce(roleID, resource, action string) (bool, error) {
	// Use the domain "万知" for all permissions
	return s.enforcer.Enforce(roleID, "万知", resource, action)
}

// AddRoleForUser adds a role inheritance relationship
func (s *CasbinService) AddRoleForUser(user, role, domain string) (bool, error) {
	return s.enforcer.AddRoleForUserInDomain(user, role, domain)
}

// DeleteRoleForUser removes a role inheritance relationship
func (s *CasbinService) DeleteRoleForUser(user, role, domain string) (bool, error) {
	return s.enforcer.DeleteRoleForUserInDomain(user, role, domain)
}

// SavePolicy saves the current policy to storage
func (s *CasbinService) SavePolicy() error {
	return s.enforcer.SavePolicy()
}
