package harvest

import (
	"context"
	"fmt"
)

// TimeEntriesService handles communication with the time entry related
// methods of the Harvest API.
type TimeEntriesService struct {
	client *API
}

// TimeEntryListOptions specifies optional parameters to the List method.
type TimeEntryListOptions struct {
	ListOptions
	UserID              int64  `url:"user_id,omitempty"`
	ClientID            int64  `url:"client_id,omitempty"`
	ProjectID           int64  `url:"project_id,omitempty"`
	TaskID              int64  `url:"task_id,omitempty"`
	ExternalReferenceID string `url:"external_reference_id,omitempty"`
	IsBilled            *bool  `url:"is_billed,omitempty"`
	IsRunning           *bool  `url:"is_running,omitempty"`
	ApprovalStatus      string `url:"approval_status,omitempty"`
	UpdatedSince        string `url:"updated_since,omitempty"`
	From                string `url:"from,omitempty"`
	To                  string `url:"to,omitempty"`
}

// TimeEntryList represents a list of time entries.
type TimeEntryList struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	Paginated[TimeEntry]
}

// ListPage returns a single page of time entries.
func (s *TimeEntriesService) ListPage(ctx context.Context, opts *TimeEntryListOptions) (*TimeEntryList, error) {
	u, err := addOptions("time_entries", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var entries TimeEntryList
	_, err = s.client.Do(ctx, req, &entries)
	if err != nil {
		return nil, err
	}

	// Copy entries to Items for pagination
	entries.Items = entries.TimeEntries

	return &entries, nil
}

// List returns all time entries across all pages.
func (s *TimeEntriesService) List(ctx context.Context, opts *TimeEntryListOptions) ([]TimeEntry, error) {
	if opts == nil {
		opts = &TimeEntryListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allEntries []TimeEntry

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allEntries = append(allEntries, result.TimeEntries...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allEntries, nil
}

// Get retrieves a specific time entry.
func (s *TimeEntriesService) Get(ctx context.Context, timeEntryID int64) (*TimeEntry, error) {
	return Get[TimeEntry](ctx, s.client, fmt.Sprintf("time_entries/%d", timeEntryID))
}

// TimeEntryCreateViaDurationRequest represents a request to create a time entry via duration.
type TimeEntryCreateViaDurationRequest struct {
	ProjectID         int64                     `json:"project_id"`
	TaskID            int64                     `json:"task_id"`
	SpentDate         string                    `json:"spent_date"`
	Hours             float64                   `json:"hours"`
	UserID            int64                     `json:"user_id,omitempty"`
	Notes             string                    `json:"notes,omitempty"`
	ExternalReference *ExternalReferenceRequest `json:"external_reference,omitempty"`
}

// ExternalReferenceRequest represents an external reference in a request.
type ExternalReferenceRequest struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	AccountID string `json:"account_id"`
	Permalink string `json:"permalink"`
}

// CreateViaDuration creates a new time entry via duration.
func (s *TimeEntriesService) CreateViaDuration(ctx context.Context, entry *TimeEntryCreateViaDurationRequest) (*TimeEntry, error) {
	return Create[TimeEntry](ctx, s.client, "time_entries", entry)
}

// TimeEntryCreateViaStartEndRequest represents a request to create a time entry via start and end time.
type TimeEntryCreateViaStartEndRequest struct {
	ProjectID         int64                     `json:"project_id"`
	TaskID            int64                     `json:"task_id"`
	SpentDate         string                    `json:"spent_date"`
	StartedTime       string                    `json:"started_time"`
	EndedTime         string                    `json:"ended_time"`
	UserID            int64                     `json:"user_id,omitempty"`
	Notes             string                    `json:"notes,omitempty"`
	ExternalReference *ExternalReferenceRequest `json:"external_reference,omitempty"`
}

// CreateViaStartEnd creates a new time entry via start and end time.
func (s *TimeEntriesService) CreateViaStartEnd(ctx context.Context, entry *TimeEntryCreateViaStartEndRequest) (*TimeEntry, error) {
	return Create[TimeEntry](ctx, s.client, "time_entries", entry)
}

// TimeEntryUpdateRequest represents a request to update a time entry.
type TimeEntryUpdateRequest struct {
	ProjectID         int64                     `json:"project_id,omitempty"`
	TaskID            int64                     `json:"task_id,omitempty"`
	SpentDate         string                    `json:"spent_date,omitempty"`
	StartedTime       string                    `json:"started_time,omitempty"`
	EndedTime         string                    `json:"ended_time,omitempty"`
	Hours             float64                   `json:"hours,omitempty"`
	Notes             string                    `json:"notes,omitempty"`
	ExternalReference *ExternalReferenceRequest `json:"external_reference,omitempty"`
}

// Update updates a time entry.
func (s *TimeEntriesService) Update(ctx context.Context, timeEntryID int64, entry *TimeEntryUpdateRequest) (*TimeEntry, error) {
	return Update[TimeEntry](ctx, s.client, fmt.Sprintf("time_entries/%d", timeEntryID), entry)
}

// Delete deletes a time entry.
func (s *TimeEntriesService) Delete(ctx context.Context, timeEntryID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("time_entries/%d", timeEntryID))
}

// RestartRequest represents a request to restart a time entry.
type RestartRequest struct {
	ID int64 `json:"id"`
}

// Restart restarts a stopped time entry.
func (s *TimeEntriesService) Restart(ctx context.Context, timeEntryID int64) (*TimeEntry, error) {
	req := RestartRequest{ID: timeEntryID}
	return Update[TimeEntry](ctx, s.client, fmt.Sprintf("time_entries/%d/restart", timeEntryID), req)
}

// Stop stops a running time entry.
func (s *TimeEntriesService) Stop(ctx context.Context, timeEntryID int64) (*TimeEntry, error) {
	return Update[TimeEntry](ctx, s.client, fmt.Sprintf("time_entries/%d/stop", timeEntryID), nil)
}

// DeleteExternalReference deletes an external reference from a time entry.
func (s *TimeEntriesService) DeleteExternalReference(ctx context.Context, timeEntryID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("time_entries/%d/external_reference", timeEntryID))
}
