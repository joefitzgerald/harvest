package harvest

import (
	"context"
	"fmt"
)

// PTOService handles communication with the Paid Time Off (PTO) related methods of the
// Harvest API: https://help.getharvest.com/api-v2/pto-api/pto/
type PTOService struct {
	client *API
}

// PTOTypeListOptions specifies optional parameters for listing time off policies.
type PTOTypeListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// PTORequestListOptions specifies optional parameters for listing time off requests.
type PTORequestListOptions struct {
	ListOptions
	UserID       int64  `url:"user_id,omitempty"`
	Status       string `url:"status,omitempty"`
	PTOTypeID    int64  `url:"pto_type_id,omitempty"`
	From         string `url:"from,omitempty"`
	To           string `url:"to,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// WorkScheduleListOptions specifies optional parameters for listing work schedules.
type WorkScheduleListOptions struct {
	ListOptions
	IsActive *bool `url:"is_active,omitempty"`
}

// PTOAssignmentListOptions specifies optional parameters for listing PTO user assignments.
type PTOAssignmentListOptions struct {
	UserID  int64   `url:"user_id,omitempty"`
	UserIDs []int64 `url:"user_ids,omitempty,comma"`
}

// PTOAllocationListOptions specifies optional parameters for listing PTO allocations.
type PTOAllocationListOptions struct {
	ListOptions
	Year       int     `url:"year,omitempty"`
	UserIDs    []int64 `url:"user_ids,omitempty,comma"`
	PTOTypeIDs []int64 `url:"pto_type_ids,omitempty,comma"`
}

// PTOBalanceOptions specifies parameters for retrieving a user's time off balances.
type PTOBalanceOptions struct {
	UserID int64  `url:"user_id,omitempty"`
	Year   int    `url:"year,omitempty"`
	AsOf   string `url:"as_of,omitempty"`
}

// ListHolidayCalendars returns all holiday calendars (GET /v2/pto/holiday_calendars).
func (s *PTOService) ListHolidayCalendars(ctx context.Context, opts *ListOptions) ([]HolidayCalendar, error) {
	return listValues[HolidayCalendar](ctx, s.client, "pto/holiday_calendars", "pto_holiday_calendars", defaultListOptions(opts))
}

// ListHolidayCalendarEntries returns the entries of one holiday calendar.
func (s *PTOService) ListHolidayCalendarEntries(ctx context.Context, calendarID int64, opts *ListOptions) ([]HolidayCalendarEntry, error) {
	return listValues[HolidayCalendarEntry](ctx, s.client, fmt.Sprintf("pto/holiday_calendars/%d/entries", calendarID), "", defaultListOptions(opts))
}

// ListTypes returns all time off policies (GET /v2/pto/types).
func (s *PTOService) ListTypes(ctx context.Context, opts *PTOTypeListOptions) ([]PTOType, error) {
	if opts == nil {
		opts = &PTOTypeListOptions{}
	}
	opts.ListOptions = *defaultListOptions(&opts.ListOptions)
	return listValues[PTOType](ctx, s.client, "pto/types", "pto_types", opts)
}

// ListRequests returns time off requests (GET /v2/pto/requests).
func (s *PTOService) ListRequests(ctx context.Context, opts *PTORequestListOptions) ([]PTORequest, error) {
	if opts == nil {
		opts = &PTORequestListOptions{}
	}
	opts.ListOptions = *defaultListOptions(&opts.ListOptions)
	return listValues[PTORequest](ctx, s.client, "pto/requests", "pto_requests", opts)
}

// ListWorkSchedules returns work schedules (GET /v2/pto/work_schedules).
func (s *PTOService) ListWorkSchedules(ctx context.Context, opts *WorkScheduleListOptions) ([]WorkSchedule, error) {
	if opts == nil {
		opts = &WorkScheduleListOptions{}
	}
	opts.ListOptions = *defaultListOptions(&opts.ListOptions)
	return listValues[WorkSchedule](ctx, s.client, "pto/work_schedules", "pto_work_schedules", opts)
}

// ListAssignments returns PTO user assignments (GET /v2/pto/assignments).
func (s *PTOService) ListAssignments(ctx context.Context, opts *PTOAssignmentListOptions) ([]PTOAssignment, error) {
	return listValues[PTOAssignment](ctx, s.client, "pto/assignments", "", opts)
}

// ListAllocations returns PTO allocations for a year (GET /v2/pto/allocations).
func (s *PTOService) ListAllocations(ctx context.Context, opts *PTOAllocationListOptions) ([]PTOAllocation, error) {
	if opts == nil {
		opts = &PTOAllocationListOptions{}
	}
	opts.ListOptions = *defaultListOptions(&opts.ListOptions)
	return listValues[PTOAllocation](ctx, s.client, "pto/allocations", "pto_allocations", opts)
}

// GetBalances returns a user's time off balances (GET /v2/pto/balances).
func (s *PTOService) GetBalances(ctx context.Context, opts *PTOBalanceOptions) (*PTOBalances, error) {
	item, err := GetRaw[PTOBalances](ctx, s.client, "pto/balances", opts)
	if err != nil {
		return nil, err
	}
	return &item.Value, nil
}

// defaultListOptions returns opts with Page and PerPage defaulted, never nil.
func defaultListOptions(opts *ListOptions) *ListOptions {
	if opts == nil {
		opts = &ListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}
	return opts
}
