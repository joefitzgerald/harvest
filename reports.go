package harvest

import (
	"context"

	"github.com/shopspring/decimal"
)

// ReportsService handles communication with the report related
// methods of the Harvest API.
type ReportsService struct {
	client *API
}

// TimeReportsOptions specifies optional parameters for time reports.
type TimeReportsOptions struct {
	From           string `url:"from"`
	To             string `url:"to"`
	ClientID       int64  `url:"client_id,omitempty"`
	ProjectID      int64  `url:"project_id,omitempty"`
	TaskID         int64  `url:"task_id,omitempty"`
	UserID         int64  `url:"user_id,omitempty"`
	IsBilled       *bool  `url:"is_billed,omitempty"`
	IsRunning      *bool  `url:"is_running,omitempty"`
	OnlyBillable   *bool  `url:"only_billable,omitempty"`
	OnlyUnbillable *bool  `url:"only_unbillable,omitempty"`
	Page           int    `url:"page,omitempty"`
	PerPage        int    `url:"per_page,omitempty"`
}

// TimeReport represents a time report entry.
type TimeReport struct {
	ClientID       int64           `json:"client_id"`
	ClientName     string          `json:"client_name"`
	ProjectID      int64           `json:"project_id"`
	ProjectName    string          `json:"project_name"`
	TaskID         int64           `json:"task_id"`
	TaskName       string          `json:"task_name"`
	UserID         int64           `json:"user_id"`
	UserName       string          `json:"user_name"`
	WeeklyCapacity int             `json:"weekly_capacity"`
	AvatarURL      string          `json:"avatar_url"`
	IsContractor   bool            `json:"is_contractor"`
	TotalHours     decimal.Decimal `json:"total_hours"`
	BillableHours  decimal.Decimal `json:"billable_hours"`
	Currency       string          `json:"currency"`
	BillableAmount decimal.Decimal `json:"billable_amount"`
}

// TimeReportResults represents time report results.
type TimeReportResults struct {
	Results      []TimeReport     `json:"results"`
	PerPage      int              `json:"per_page"`
	TotalPages   int              `json:"total_pages"`
	TotalEntries int              `json:"total_entries"`
	NextPage     *int             `json:"next_page"`
	PreviousPage *int             `json:"previous_page"`
	Page         int              `json:"page"`
	Links        *PaginationLinks `json:"links"`
}

// TimeReports retrieves time reports.
func (s *ReportsService) TimeReports(ctx context.Context, opts *TimeReportsOptions) (*TimeReportResults, error) {
	u, err := addOptions("reports/time/team", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var report TimeReportResults
	_, err = s.client.Do(ctx, req, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// ExpenseReportsOptions specifies optional parameters for expense reports.
type ExpenseReportsOptions struct {
	From      string `url:"from"`
	To        string `url:"to"`
	ClientID  int64  `url:"client_id,omitempty"`
	ProjectID int64  `url:"project_id,omitempty"`
	UserID    int64  `url:"user_id,omitempty"`
	IsBilled  *bool  `url:"is_billed,omitempty"`
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
}

// ExpenseReport represents an expense report entry.
type ExpenseReport struct {
	ClientID            int64           `json:"client_id"`
	ClientName          string          `json:"client_name"`
	ProjectID           int64           `json:"project_id"`
	ProjectName         string          `json:"project_name"`
	ExpenseCategoryID   int64           `json:"expense_category_id"`
	ExpenseCategoryName string          `json:"expense_category_name"`
	UserID              int64           `json:"user_id"`
	UserName            string          `json:"user_name"`
	IsContractor        bool            `json:"is_contractor"`
	TotalAmount         decimal.Decimal `json:"total_amount"`
	BillableAmount      decimal.Decimal `json:"billable_amount"`
	Currency            string          `json:"currency"`
}

// ExpenseReportResults represents expense report results.
type ExpenseReportResults struct {
	Results      []ExpenseReport  `json:"results"`
	PerPage      int              `json:"per_page"`
	TotalPages   int              `json:"total_pages"`
	TotalEntries int              `json:"total_entries"`
	NextPage     *int             `json:"next_page"`
	PreviousPage *int             `json:"previous_page"`
	Page         int              `json:"page"`
	Links        *PaginationLinks `json:"links"`
}

// ExpenseReports retrieves expense reports.
func (s *ReportsService) ExpenseReports(ctx context.Context, opts *ExpenseReportsOptions) (*ExpenseReportResults, error) {
	u, err := addOptions("reports/expenses/team", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var report ExpenseReportResults
	_, err = s.client.Do(ctx, req, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// UninvoicedReportOptions specifies optional parameters for uninvoiced reports.
type UninvoicedReportOptions struct {
	From      string `url:"from"`
	To        string `url:"to"`
	ClientID  int64  `url:"client_id,omitempty"`
	ProjectID int64  `url:"project_id,omitempty"`
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
}

// UninvoicedReport represents an uninvoiced report entry.
type UninvoicedReport struct {
	ClientID           int64           `json:"client_id"`
	ClientName         string          `json:"client_name"`
	ProjectID          int64           `json:"project_id"`
	ProjectName        string          `json:"project_name"`
	Currency           string          `json:"currency"`
	TotalHours         decimal.Decimal `json:"total_hours"`
	UninvoicedHours    decimal.Decimal `json:"uninvoiced_hours"`
	UninvoicedExpenses decimal.Decimal `json:"uninvoiced_expenses"`
	UninvoicedAmount   decimal.Decimal `json:"uninvoiced_amount"`
}

// UninvoicedReportResults represents uninvoiced report results.
type UninvoicedReportResults struct {
	Results      []UninvoicedReport `json:"results"`
	PerPage      int                `json:"per_page"`
	TotalPages   int                `json:"total_pages"`
	TotalEntries int                `json:"total_entries"`
	NextPage     *int               `json:"next_page"`
	PreviousPage *int               `json:"previous_page"`
	Page         int                `json:"page"`
	Links        *PaginationLinks   `json:"links"`
}

// UninvoicedReports retrieves uninvoiced reports.
func (s *ReportsService) UninvoicedReports(ctx context.Context, opts *UninvoicedReportOptions) (*UninvoicedReportResults, error) {
	u, err := addOptions("reports/uninvoiced", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var report UninvoicedReportResults
	_, err = s.client.Do(ctx, req, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// ProjectBudgetReportOptions specifies optional parameters for project budget reports.
type ProjectBudgetReportOptions struct {
	Page         int    `url:"page,omitempty"`
	PerPage      int    `url:"per_page,omitempty"`
	IsActive     *bool  `url:"is_active,omitempty"`
	ClientID     int64  `url:"client_id,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// ProjectBudgetReport represents a project budget report entry.
type ProjectBudgetReport struct {
	ClientID         int64            `json:"client_id"`
	ClientName       string           `json:"client_name"`
	ProjectID        int64            `json:"project_id"`
	ProjectName      string           `json:"project_name"`
	ProjectCode      string           `json:"project_code"`
	ProjectStartDate *Date            `json:"project_start_date"`
	ProjectEndDate   *Date            `json:"project_end_date"`
	IsBillable       bool             `json:"is_billable"`
	IsActive         bool             `json:"is_active"`
	BudgetBy         string           `json:"budget_by"`
	Budget           *decimal.Decimal `json:"budget"`
	BudgetSpent      decimal.Decimal  `json:"budget_spent"`
	BudgetRemaining  *decimal.Decimal `json:"budget_remaining"`
}

// ProjectBudgetReportResults represents project budget report results.
type ProjectBudgetReportResults struct {
	Results      []ProjectBudgetReport `json:"results"`
	PerPage      int                   `json:"per_page"`
	TotalPages   int                   `json:"total_pages"`
	TotalEntries int                   `json:"total_entries"`
	NextPage     *int                  `json:"next_page"`
	PreviousPage *int                  `json:"previous_page"`
	Page         int                   `json:"page"`
	Links        *PaginationLinks      `json:"links"`
}

// ProjectBudgetReports retrieves project budget reports.
func (s *ReportsService) ProjectBudgetReports(ctx context.Context, opts *ProjectBudgetReportOptions) (*ProjectBudgetReportResults, error) {
	u, err := addOptions("reports/project_budget", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var report ProjectBudgetReportResults
	_, err = s.client.Do(ctx, req, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
