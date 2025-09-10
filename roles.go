package harvest

import (
	"context"
	"fmt"
)

// RolesService handles communication with the role related
// methods of the Harvest API.
type RolesService struct {
	client *API
}

// RoleListOptions specifies optional parameters to the List method.
type RoleListOptions struct {
	ListOptions
}

// RoleList represents a list of roles.
type RoleList struct {
	Roles []Role `json:"roles"`
	Paginated[Role]
}

// List returns a list of roles.
func (s *RolesService) List(ctx context.Context, opts *RoleListOptions) (*RoleList, error) {
	u, err := addOptions("roles", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var roles RoleList
	_, err = s.client.Do(ctx, req, &roles)
	if err != nil {
		return nil, err
	}

	// Copy roles to Items for pagination
	roles.Items = roles.Roles
	
	return &roles, nil
}

// Get retrieves a specific role.
func (s *RolesService) Get(ctx context.Context, roleID int64) (*Role, error) {
	return Get[Role](ctx, s.client, fmt.Sprintf("roles/%d", roleID))
}

// RoleCreateRequest represents a request to create a role.
type RoleCreateRequest struct {
	Name    string  `json:"name"`
	UserIDs []int64 `json:"user_ids,omitempty"`
}

// Create creates a new role.
func (s *RolesService) Create(ctx context.Context, role *RoleCreateRequest) (*Role, error) {
	return Create[Role](ctx, s.client, "roles", role)
}

// RoleUpdateRequest represents a request to update a role.
type RoleUpdateRequest struct {
	Name    string  `json:"name,omitempty"`
	UserIDs []int64 `json:"user_ids,omitempty"`
}

// Update updates a role.
func (s *RolesService) Update(ctx context.Context, roleID int64, role *RoleUpdateRequest) (*Role, error) {
	return Update[Role](ctx, s.client, fmt.Sprintf("roles/%d", roleID), role)
}

// Delete deletes a role.
func (s *RolesService) Delete(ctx context.Context, roleID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("roles/%d", roleID))
}