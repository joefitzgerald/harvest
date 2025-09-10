package harvest

import (
	"context"
	"fmt"
)

// TasksService handles communication with the task related
// methods of the Harvest API.
type TasksService struct {
	client *API
}

// TaskListOptions specifies optional parameters to the List method.
type TaskListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// TaskList represents a list of tasks.
type TaskList struct {
	Tasks []Task `json:"tasks"`
	Paginated[Task]
}

// List returns a list of tasks.
func (s *TasksService) List(ctx context.Context, opts *TaskListOptions) (*TaskList, error) {
	u, err := addOptions("tasks", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var tasks TaskList
	_, err = s.client.Do(ctx, req, &tasks)
	if err != nil {
		return nil, err
	}

	// Copy tasks to Items for pagination
	tasks.Items = tasks.Tasks
	
	return &tasks, nil
}

// Get retrieves a specific task.
func (s *TasksService) Get(ctx context.Context, taskID int64) (*Task, error) {
	return Get[Task](ctx, s.client, fmt.Sprintf("tasks/%d", taskID))
}

// TaskCreateRequest represents a request to create a task.
type TaskCreateRequest struct {
	Name              string  `json:"name"`
	BillableByDefault *bool   `json:"billable_by_default,omitempty"`
	DefaultHourlyRate float64 `json:"default_hourly_rate,omitempty"`
	IsDefault         *bool   `json:"is_default,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// Create creates a new task.
func (s *TasksService) Create(ctx context.Context, task *TaskCreateRequest) (*Task, error) {
	return Create[Task](ctx, s.client, "tasks", task)
}

// TaskUpdateRequest represents a request to update a task.
type TaskUpdateRequest struct {
	Name              string  `json:"name,omitempty"`
	BillableByDefault *bool   `json:"billable_by_default,omitempty"`
	DefaultHourlyRate float64 `json:"default_hourly_rate,omitempty"`
	IsDefault         *bool   `json:"is_default,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// Update updates a task.
func (s *TasksService) Update(ctx context.Context, taskID int64, task *TaskUpdateRequest) (*Task, error) {
	return Update[Task](ctx, s.client, fmt.Sprintf("tasks/%d", taskID), task)
}

// Delete deletes a task.
func (s *TasksService) Delete(ctx context.Context, taskID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("tasks/%d", taskID))
}