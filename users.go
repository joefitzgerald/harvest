package harvest

import (
	"context"
	"fmt"
	"net/url"
)

// UsersService handles communication with the user related
// methods of the Harvest API.
type UsersService struct {
	client *API
}

// UserListOptions specifies optional parameters to the List method.
type UserListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// UserList represents a list of users.
type UserList struct {
	Users []User `json:"users"`
	Paginated[User]
}

// ListPage returns a single page of users.
func (s *UsersService) ListPage(ctx context.Context, opts *UserListOptions) (*UserList, error) {
	u, err := addOptions("users", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var users UserList
	_, err = s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, err
	}

	// Copy users to Items for pagination
	users.Items = users.Users

	return &users, nil
}

// List returns all users across all pages.
// This endpoint uses cursor-based pagination.
func (s *UsersService) List(ctx context.Context, opts *UserListOptions) ([]User, error) {
	if opts == nil {
		opts = &UserListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allUsers []User

	// Fetch first page
	result, err := s.ListPage(ctx, opts)
	if err != nil {
		return nil, err
	}
	allUsers = append(allUsers, result.Users...)

	// Continue fetching remaining pages
	for result.HasNextPage() {
		// Check if using cursor-based pagination
		if nextURL := result.GetNextPageURL(); nextURL != "" {
			// Parse the URL to get path and query
			u, err := url.Parse(nextURL)
			if err != nil {
				return nil, err
			}
			pathAndQuery := u.Path
			if u.RawQuery != "" {
				pathAndQuery += "?" + u.RawQuery
			}

			req, err := s.client.NewRequest(ctx, "GET", pathAndQuery, nil)
			if err != nil {
				return nil, err
			}

			var users UserList
			_, err = s.client.Do(ctx, req, &users)
			if err != nil {
				return nil, err
			}
			users.Items = users.Users
			result = &users
			allUsers = append(allUsers, users.Users...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListPage(ctx, opts)
			if err != nil {
				return nil, err
			}
			allUsers = append(allUsers, result.Users...)
		} else {
			break
		}
	}

	return allUsers, nil
}

// Get retrieves a specific user.
func (s *UsersService) Get(ctx context.Context, userID int64) (*User, error) {
	return Get[User](ctx, s.client, fmt.Sprintf("users/%d", userID))
}

// Me retrieves the currently authenticated user.
func (s *UsersService) Me(ctx context.Context) (*User, error) {
	return Get[User](ctx, s.client, "users/me")
}

// UserCreateRequest represents a request to create a user.
type UserCreateRequest struct {
	FirstName                    string   `json:"first_name"`
	LastName                     string   `json:"last_name"`
	Email                        string   `json:"email"`
	Telephone                    string   `json:"telephone,omitempty"`
	Timezone                     string   `json:"timezone,omitempty"`
	HasAccessToAllFutureProjects *bool    `json:"has_access_to_all_future_projects,omitempty"`
	IsContractor                 *bool    `json:"is_contractor,omitempty"`
	IsActive                     *bool    `json:"is_active,omitempty"`
	WeeklyCapacity               int      `json:"weekly_capacity,omitempty"`
	DefaultHourlyRate            float64  `json:"default_hourly_rate,omitempty"`
	CostRate                     float64  `json:"cost_rate,omitempty"`
	Roles                        []string `json:"roles,omitempty"`
}

// Create creates a new user.
func (s *UsersService) Create(ctx context.Context, user *UserCreateRequest) (*User, error) {
	return Create[User](ctx, s.client, "users", user)
}

// UserUpdateRequest represents a request to update a user.
type UserUpdateRequest struct {
	FirstName                    string   `json:"first_name,omitempty"`
	LastName                     string   `json:"last_name,omitempty"`
	Email                        string   `json:"email,omitempty"`
	Telephone                    string   `json:"telephone,omitempty"`
	Timezone                     string   `json:"timezone,omitempty"`
	HasAccessToAllFutureProjects *bool    `json:"has_access_to_all_future_projects,omitempty"`
	IsContractor                 *bool    `json:"is_contractor,omitempty"`
	IsActive                     *bool    `json:"is_active,omitempty"`
	WeeklyCapacity               int      `json:"weekly_capacity,omitempty"`
	DefaultHourlyRate            float64  `json:"default_hourly_rate,omitempty"`
	CostRate                     float64  `json:"cost_rate,omitempty"`
	Roles                        []string `json:"roles,omitempty"`
}

// Update updates a user.
func (s *UsersService) Update(ctx context.Context, userID int64, user *UserUpdateRequest) (*User, error) {
	return Update[User](ctx, s.client, fmt.Sprintf("users/%d", userID), user)
}

// Delete archives a user.
func (s *UsersService) Delete(ctx context.Context, userID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("users/%d", userID))
}

// UserProjectAssignmentListOptions specifies optional parameters for listing user project assignments.
type UserProjectAssignmentListOptions struct {
	ListOptions
	UpdatedSince string `url:"updated_since,omitempty"`
}

// UserProjectAssignmentList represents a list of user project assignments.
type UserProjectAssignmentList struct {
	ProjectAssignments []ProjectUserAssignment `json:"project_assignments"`
	Paginated[ProjectUserAssignment]
}

// ListProjectAssignmentsPage returns a single page of project assignments for a user.
func (s *UsersService) ListProjectAssignmentsPage(ctx context.Context, userID int64, opts *UserProjectAssignmentListOptions) (*UserProjectAssignmentList, error) {
	u, err := addOptions(fmt.Sprintf("users/%d/project_assignments", userID), opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var assignments UserProjectAssignmentList
	_, err = s.client.Do(ctx, req, &assignments)
	if err != nil {
		return nil, err
	}

	// Copy assignments to Items for pagination
	assignments.Items = assignments.ProjectAssignments

	return &assignments, nil
}

// ListProjectAssignments returns all project assignments for a user across all pages.
// This endpoint uses cursor-based pagination.
func (s *UsersService) ListProjectAssignments(ctx context.Context, userID int64, opts *UserProjectAssignmentListOptions) ([]ProjectUserAssignment, error) {
	if opts == nil {
		opts = &UserProjectAssignmentListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allAssignments []ProjectUserAssignment

	// Fetch first page
	result, err := s.ListProjectAssignmentsPage(ctx, userID, opts)
	if err != nil {
		return nil, err
	}
	allAssignments = append(allAssignments, result.ProjectAssignments...)

	// Continue fetching remaining pages
	for result.HasNextPage() {
		// Check if using cursor-based pagination
		if nextURL := result.GetNextPageURL(); nextURL != "" {
			// Parse the URL to get path and query
			u, err := url.Parse(nextURL)
			if err != nil {
				return nil, err
			}
			pathAndQuery := u.Path
			if u.RawQuery != "" {
				pathAndQuery += "?" + u.RawQuery
			}

			req, err := s.client.NewRequest(ctx, "GET", pathAndQuery, nil)
			if err != nil {
				return nil, err
			}

			var assignments UserProjectAssignmentList
			_, err = s.client.Do(ctx, req, &assignments)
			if err != nil {
				return nil, err
			}
			assignments.Items = assignments.ProjectAssignments
			result = &assignments
			allAssignments = append(allAssignments, assignments.ProjectAssignments...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListProjectAssignmentsPage(ctx, userID, opts)
			if err != nil {
				return nil, err
			}
			allAssignments = append(allAssignments, result.ProjectAssignments...)
		} else {
			break
		}
	}

	return allAssignments, nil
}

// ListMyProjectAssignmentsPage returns a single page of project assignments for the currently authenticated user.
func (s *UsersService) ListMyProjectAssignmentsPage(ctx context.Context, opts *UserProjectAssignmentListOptions) (*UserProjectAssignmentList, error) {
	u, err := addOptions("users/me/project_assignments", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var assignments UserProjectAssignmentList
	_, err = s.client.Do(ctx, req, &assignments)
	if err != nil {
		return nil, err
	}

	// Copy assignments to Items for pagination
	assignments.Items = assignments.ProjectAssignments

	return &assignments, nil
}

// ListMyProjectAssignments returns all project assignments for the currently authenticated user across all pages.
func (s *UsersService) ListMyProjectAssignments(ctx context.Context, opts *UserProjectAssignmentListOptions) ([]ProjectUserAssignment, error) {
	if opts == nil {
		opts = &UserProjectAssignmentListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allAssignments []ProjectUserAssignment

	for {
		result, err := s.ListMyProjectAssignmentsPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allAssignments = append(allAssignments, result.ProjectAssignments...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allAssignments, nil
}
