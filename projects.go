package harvest

import (
	"context"
	"fmt"
	"net/url"
)

// ProjectsService handles communication with the project related
// methods of the Harvest API.
type ProjectsService struct {
	client *API
}

// ProjectListOptions specifies optional parameters to the List method.
type ProjectListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	ClientID     int64  `url:"client_id,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// ProjectList represents a list of projects.
type ProjectList struct {
	Projects []Project `json:"projects"`
	Paginated[Project]
}

// ListPage returns a single page of projects.
func (s *ProjectsService) ListPage(ctx context.Context, opts *ProjectListOptions) (*ProjectList, error) {
	u, err := addOptions("projects", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var projects ProjectList
	_, err = s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, err
	}

	// Copy projects to Items for pagination
	projects.Items = projects.Projects

	return &projects, nil
}

// List returns all projects across all pages.
func (s *ProjectsService) List(ctx context.Context, opts *ProjectListOptions) ([]Project, error) {
	if opts == nil {
		opts = &ProjectListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allProjects []Project

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, result.Projects...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allProjects, nil
}

// Get retrieves a specific project.
func (s *ProjectsService) Get(ctx context.Context, projectID int64) (*Project, error) {
	return Get[Project](ctx, s.client, fmt.Sprintf("projects/%d", projectID))
}

// ProjectCreateRequest represents a request to create a project.
type ProjectCreateRequest struct {
	ClientID                         int64   `json:"client_id"`
	Name                             string  `json:"name"`
	Code                             string  `json:"code,omitempty"`
	IsActive                         *bool   `json:"is_active,omitempty"`
	IsBillable                       *bool   `json:"is_billable,omitempty"`
	IsFixedFee                       *bool   `json:"is_fixed_fee,omitempty"`
	BillBy                           string  `json:"bill_by,omitempty"`
	Budget                           float64 `json:"budget,omitempty"`
	BudgetBy                         string  `json:"budget_by,omitempty"`
	BudgetIsMonthly                  *bool   `json:"budget_is_monthly,omitempty"`
	NotifyWhenOverBudget             *bool   `json:"notify_when_over_budget,omitempty"`
	OverBudgetNotificationPercentage float64 `json:"over_budget_notification_percentage,omitempty"`
	ShowBudgetToAll                  *bool   `json:"show_budget_to_all,omitempty"`
	CostBudget                       float64 `json:"cost_budget,omitempty"`
	CostBudgetIncludeExpenses        *bool   `json:"cost_budget_include_expenses,omitempty"`
	HourlyRate                       float64 `json:"hourly_rate,omitempty"`
	Fee                              float64 `json:"fee,omitempty"`
	Notes                            string  `json:"notes,omitempty"`
	StartsOn                         string  `json:"starts_on,omitempty"`
	EndsOn                           string  `json:"ends_on,omitempty"`
}

// Create creates a new project.
func (s *ProjectsService) Create(ctx context.Context, project *ProjectCreateRequest) (*Project, error) {
	return Create[Project](ctx, s.client, "projects", project)
}

// ProjectUpdateRequest represents a request to update a project.
type ProjectUpdateRequest struct {
	ClientID                         int64   `json:"client_id,omitempty"`
	Name                             string  `json:"name,omitempty"`
	Code                             string  `json:"code,omitempty"`
	IsActive                         *bool   `json:"is_active,omitempty"`
	IsBillable                       *bool   `json:"is_billable,omitempty"`
	IsFixedFee                       *bool   `json:"is_fixed_fee,omitempty"`
	BillBy                           string  `json:"bill_by,omitempty"`
	Budget                           float64 `json:"budget,omitempty"`
	BudgetBy                         string  `json:"budget_by,omitempty"`
	BudgetIsMonthly                  *bool   `json:"budget_is_monthly,omitempty"`
	NotifyWhenOverBudget             *bool   `json:"notify_when_over_budget,omitempty"`
	OverBudgetNotificationPercentage float64 `json:"over_budget_notification_percentage,omitempty"`
	ShowBudgetToAll                  *bool   `json:"show_budget_to_all,omitempty"`
	CostBudget                       float64 `json:"cost_budget,omitempty"`
	CostBudgetIncludeExpenses        *bool   `json:"cost_budget_include_expenses,omitempty"`
	HourlyRate                       float64 `json:"hourly_rate,omitempty"`
	Fee                              float64 `json:"fee,omitempty"`
	Notes                            string  `json:"notes,omitempty"`
	StartsOn                         string  `json:"starts_on,omitempty"`
	EndsOn                           string  `json:"ends_on,omitempty"`
}

// Update updates a project.
func (s *ProjectsService) Update(ctx context.Context, projectID int64, project *ProjectUpdateRequest) (*Project, error) {
	return Update[Project](ctx, s.client, fmt.Sprintf("projects/%d", projectID), project)
}

// Delete deletes a project.
func (s *ProjectsService) Delete(ctx context.Context, projectID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("projects/%d", projectID))
}

// UserAssignmentListOptions specifies optional parameters for listing user assignments.
type UserAssignmentListOptions struct {
	ListOptions
	UserID       int64  `url:"user_id,omitempty"`
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// UserAssignmentList represents a list of user assignments.
type UserAssignmentList struct {
	UserAssignments []ProjectUserAssignment `json:"user_assignments"`
	Paginated[ProjectUserAssignment]
}

// ListUserAssignmentsPage returns a single page of user assignments for a project.
func (s *ProjectsService) ListUserAssignmentsPage(ctx context.Context, projectID int64, opts *UserAssignmentListOptions) (*UserAssignmentList, error) {
	u, err := addOptions(fmt.Sprintf("projects/%d/user_assignments", projectID), opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var assignments UserAssignmentList
	_, err = s.client.Do(ctx, req, &assignments)
	if err != nil {
		return nil, err
	}

	// Copy assignments to Items for pagination
	assignments.Items = assignments.UserAssignments

	return &assignments, nil
}

// ListUserAssignments returns all user assignments for a project across all pages.
// This endpoint uses cursor-based pagination.
func (s *ProjectsService) ListUserAssignments(ctx context.Context, projectID int64, opts *UserAssignmentListOptions) ([]ProjectUserAssignment, error) {
	if opts == nil {
		opts = &UserAssignmentListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allAssignments []ProjectUserAssignment

	// Fetch first page
	result, err := s.ListUserAssignmentsPage(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}
	allAssignments = append(allAssignments, result.UserAssignments...)

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

			var assignments UserAssignmentList
			_, err = s.client.Do(ctx, req, &assignments)
			if err != nil {
				return nil, err
			}
			assignments.Items = assignments.UserAssignments
			result = &assignments
			allAssignments = append(allAssignments, assignments.UserAssignments...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListUserAssignmentsPage(ctx, projectID, opts)
			if err != nil {
				return nil, err
			}
			allAssignments = append(allAssignments, result.UserAssignments...)
		} else {
			break
		}
	}

	return allAssignments, nil
}

// GetUserAssignment retrieves a specific user assignment.
func (s *ProjectsService) GetUserAssignment(ctx context.Context, projectID, userAssignmentID int64) (*ProjectUserAssignment, error) {
	return Get[ProjectUserAssignment](ctx, s.client, fmt.Sprintf("projects/%d/user_assignments/%d", projectID, userAssignmentID))
}

// UserAssignmentCreateRequest represents a request to create a user assignment.
type UserAssignmentCreateRequest struct {
	UserID           int64   `json:"user_id"`
	IsActive         *bool   `json:"is_active,omitempty"`
	IsProjectManager *bool   `json:"is_project_manager,omitempty"`
	UseDefaultRates  *bool   `json:"use_default_rates,omitempty"`
	HourlyRate       float64 `json:"hourly_rate,omitempty"`
	Budget           float64 `json:"budget,omitempty"`
}

// CreateUserAssignment creates a new user assignment for a project.
func (s *ProjectsService) CreateUserAssignment(ctx context.Context, projectID int64, assignment *UserAssignmentCreateRequest) (*ProjectUserAssignment, error) {
	return Create[ProjectUserAssignment](ctx, s.client, fmt.Sprintf("projects/%d/user_assignments", projectID), assignment)
}

// UserAssignmentUpdateRequest represents a request to update a user assignment.
type UserAssignmentUpdateRequest struct {
	IsActive         *bool   `json:"is_active,omitempty"`
	IsProjectManager *bool   `json:"is_project_manager,omitempty"`
	UseDefaultRates  *bool   `json:"use_default_rates,omitempty"`
	HourlyRate       float64 `json:"hourly_rate,omitempty"`
	Budget           float64 `json:"budget,omitempty"`
}

// UpdateUserAssignment updates a user assignment.
func (s *ProjectsService) UpdateUserAssignment(ctx context.Context, projectID, userAssignmentID int64, assignment *UserAssignmentUpdateRequest) (*ProjectUserAssignment, error) {
	return Update[ProjectUserAssignment](ctx, s.client, fmt.Sprintf("projects/%d/user_assignments/%d", projectID, userAssignmentID), assignment)
}

// DeleteUserAssignment deletes a user assignment.
func (s *ProjectsService) DeleteUserAssignment(ctx context.Context, projectID, userAssignmentID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("projects/%d/user_assignments/%d", projectID, userAssignmentID))
}

// TaskAssignmentListOptions specifies optional parameters for listing task assignments.
type TaskAssignmentListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// TaskAssignmentList represents a list of task assignments.
type TaskAssignmentList struct {
	TaskAssignments []ProjectTaskAssignment `json:"task_assignments"`
	Paginated[ProjectTaskAssignment]
}

// ListTaskAssignmentsPage returns a single page of task assignments for a project.
func (s *ProjectsService) ListTaskAssignmentsPage(ctx context.Context, projectID int64, opts *TaskAssignmentListOptions) (*TaskAssignmentList, error) {
	u, err := addOptions(fmt.Sprintf("projects/%d/task_assignments", projectID), opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var assignments TaskAssignmentList
	_, err = s.client.Do(ctx, req, &assignments)
	if err != nil {
		return nil, err
	}

	// Copy assignments to Items for pagination
	assignments.Items = assignments.TaskAssignments

	return &assignments, nil
}

// ListTaskAssignments returns all task assignments for a project across all pages.
// This endpoint uses cursor-based pagination.
func (s *ProjectsService) ListTaskAssignments(ctx context.Context, projectID int64, opts *TaskAssignmentListOptions) ([]ProjectTaskAssignment, error) {
	if opts == nil {
		opts = &TaskAssignmentListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allAssignments []ProjectTaskAssignment

	// Fetch first page
	result, err := s.ListTaskAssignmentsPage(ctx, projectID, opts)
	if err != nil {
		return nil, err
	}
	allAssignments = append(allAssignments, result.TaskAssignments...)

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

			var assignments TaskAssignmentList
			_, err = s.client.Do(ctx, req, &assignments)
			if err != nil {
				return nil, err
			}
			assignments.Items = assignments.TaskAssignments
			result = &assignments
			allAssignments = append(allAssignments, assignments.TaskAssignments...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListTaskAssignmentsPage(ctx, projectID, opts)
			if err != nil {
				return nil, err
			}
			allAssignments = append(allAssignments, result.TaskAssignments...)
		} else {
			break
		}
	}

	return allAssignments, nil
}

// GetTaskAssignment retrieves a specific task assignment.
func (s *ProjectsService) GetTaskAssignment(ctx context.Context, projectID, taskAssignmentID int64) (*ProjectTaskAssignment, error) {
	return Get[ProjectTaskAssignment](ctx, s.client, fmt.Sprintf("projects/%d/task_assignments/%d", projectID, taskAssignmentID))
}

// TaskAssignmentCreateRequest represents a request to create a task assignment.
type TaskAssignmentCreateRequest struct {
	TaskID     int64   `json:"task_id"`
	IsActive   *bool   `json:"is_active,omitempty"`
	Billable   *bool   `json:"billable,omitempty"`
	HourlyRate float64 `json:"hourly_rate,omitempty"`
	Budget     float64 `json:"budget,omitempty"`
}

// CreateTaskAssignment creates a new task assignment for a project.
func (s *ProjectsService) CreateTaskAssignment(ctx context.Context, projectID int64, assignment *TaskAssignmentCreateRequest) (*ProjectTaskAssignment, error) {
	return Create[ProjectTaskAssignment](ctx, s.client, fmt.Sprintf("projects/%d/task_assignments", projectID), assignment)
}

// TaskAssignmentUpdateRequest represents a request to update a task assignment.
type TaskAssignmentUpdateRequest struct {
	IsActive   *bool   `json:"is_active,omitempty"`
	Billable   *bool   `json:"billable,omitempty"`
	HourlyRate float64 `json:"hourly_rate,omitempty"`
	Budget     float64 `json:"budget,omitempty"`
}

// UpdateTaskAssignment updates a task assignment.
func (s *ProjectsService) UpdateTaskAssignment(ctx context.Context, projectID, taskAssignmentID int64, assignment *TaskAssignmentUpdateRequest) (*ProjectTaskAssignment, error) {
	return Update[ProjectTaskAssignment](ctx, s.client, fmt.Sprintf("projects/%d/task_assignments/%d", projectID, taskAssignmentID), assignment)
}

// DeleteTaskAssignment deletes a task assignment.
func (s *ProjectsService) DeleteTaskAssignment(ctx context.Context, projectID, taskAssignmentID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("projects/%d/task_assignments/%d", projectID, taskAssignmentID))
}
