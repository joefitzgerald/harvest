package harvest

import (
	"context"
	"fmt"
)

// ExpensesService handles communication with the expense related
// methods of the Harvest API.
type ExpensesService struct {
	client *API
}

// ExpenseListOptions specifies optional parameters to the List method.
type ExpenseListOptions struct {
	ListOptions
	UserID         int64  `url:"user_id,omitempty"`
	ClientID       int64  `url:"client_id,omitempty"`
	ProjectID      int64  `url:"project_id,omitempty"`
	IsBilled       *bool  `url:"is_billed,omitempty"`
	ApprovalStatus string `url:"approval_status,omitempty"`
	UpdatedSince   string `url:"updated_since,omitempty"`
	From           string `url:"from,omitempty"`
	To             string `url:"to,omitempty"`
}

// ExpenseList represents a list of expenses.
type ExpenseList struct {
	Expenses []Expense `json:"expenses"`
	Paginated[Expense]
}

// ListPage returns a single page of expenses.
func (s *ExpensesService) ListPage(ctx context.Context, opts *ExpenseListOptions) (*ExpenseList, error) {
	u, err := addOptions("expenses", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var expenses ExpenseList
	_, err = s.client.Do(ctx, req, &expenses)
	if err != nil {
		return nil, err
	}

	// Copy expenses to Items for pagination
	expenses.Items = expenses.Expenses

	return &expenses, nil
}

// List returns all expenses across all pages.
func (s *ExpensesService) List(ctx context.Context, opts *ExpenseListOptions) ([]Expense, error) {
	if opts == nil {
		opts = &ExpenseListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allExpenses []Expense

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allExpenses = append(allExpenses, result.Expenses...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allExpenses, nil
}

// Get retrieves a specific expense.
func (s *ExpensesService) Get(ctx context.Context, expenseID int64) (*Expense, error) {
	return Get[Expense](ctx, s.client, fmt.Sprintf("expenses/%d", expenseID))
}

// ExpenseCreateRequest represents a request to create an expense.
type ExpenseCreateRequest struct {
	ProjectID         int64   `json:"project_id"`
	ExpenseCategoryID int64   `json:"expense_category_id"`
	SpentDate         string  `json:"spent_date"`
	UserID            int64   `json:"user_id,omitempty"`
	Notes             string  `json:"notes,omitempty"`
	Units             float64 `json:"units,omitempty"`
	TotalCost         float64 `json:"total_cost,omitempty"`
	Billable          *bool   `json:"billable,omitempty"`
}

// Create creates a new expense.
func (s *ExpensesService) Create(ctx context.Context, expense *ExpenseCreateRequest) (*Expense, error) {
	return Create[Expense](ctx, s.client, "expenses", expense)
}

// ExpenseUpdateRequest represents a request to update an expense.
type ExpenseUpdateRequest struct {
	ProjectID         int64   `json:"project_id,omitempty"`
	ExpenseCategoryID int64   `json:"expense_category_id,omitempty"`
	SpentDate         string  `json:"spent_date,omitempty"`
	Notes             string  `json:"notes,omitempty"`
	Units             float64 `json:"units,omitempty"`
	TotalCost         float64 `json:"total_cost,omitempty"`
	Billable          *bool   `json:"billable,omitempty"`
}

// Update updates an expense.
func (s *ExpensesService) Update(ctx context.Context, expenseID int64, expense *ExpenseUpdateRequest) (*Expense, error) {
	return Update[Expense](ctx, s.client, fmt.Sprintf("expenses/%d", expenseID), expense)
}

// Delete deletes an expense.
func (s *ExpensesService) Delete(ctx context.Context, expenseID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("expenses/%d", expenseID))
}

// ExpenseCategoryListOptions specifies optional parameters for listing expense categories.
type ExpenseCategoryListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// ExpenseCategoryList represents a list of expense categories.
type ExpenseCategoryList struct {
	ExpenseCategories []ExpenseCategory `json:"expense_categories"`
	Paginated[ExpenseCategory]
}

// ListCategoriesPage returns a single page of expense categories.
func (s *ExpensesService) ListCategoriesPage(ctx context.Context, opts *ExpenseCategoryListOptions) (*ExpenseCategoryList, error) {
	u, err := addOptions("expense_categories", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var categories ExpenseCategoryList
	_, err = s.client.Do(ctx, req, &categories)
	if err != nil {
		return nil, err
	}

	// Copy categories to Items for pagination
	categories.Items = categories.ExpenseCategories

	return &categories, nil
}

// ListCategories returns all expense categories across all pages.
func (s *ExpensesService) ListCategories(ctx context.Context, opts *ExpenseCategoryListOptions) ([]ExpenseCategory, error) {
	if opts == nil {
		opts = &ExpenseCategoryListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allCategories []ExpenseCategory

	for {
		result, err := s.ListCategoriesPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allCategories = append(allCategories, result.ExpenseCategories...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allCategories, nil
}

// GetCategory retrieves a specific expense category.
func (s *ExpensesService) GetCategory(ctx context.Context, categoryID int64) (*ExpenseCategory, error) {
	return Get[ExpenseCategory](ctx, s.client, fmt.Sprintf("expense_categories/%d", categoryID))
}

// ExpenseCategoryCreateRequest represents a request to create an expense category.
type ExpenseCategoryCreateRequest struct {
	Name      string  `json:"name"`
	UnitName  string  `json:"unit_name,omitempty"`
	UnitPrice float64 `json:"unit_price,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// CreateCategory creates a new expense category.
func (s *ExpensesService) CreateCategory(ctx context.Context, category *ExpenseCategoryCreateRequest) (*ExpenseCategory, error) {
	return Create[ExpenseCategory](ctx, s.client, "expense_categories", category)
}

// ExpenseCategoryUpdateRequest represents a request to update an expense category.
type ExpenseCategoryUpdateRequest struct {
	Name      string  `json:"name,omitempty"`
	UnitName  string  `json:"unit_name,omitempty"`
	UnitPrice float64 `json:"unit_price,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UpdateCategory updates an expense category.
func (s *ExpensesService) UpdateCategory(ctx context.Context, categoryID int64, category *ExpenseCategoryUpdateRequest) (*ExpenseCategory, error) {
	return Update[ExpenseCategory](ctx, s.client, fmt.Sprintf("expense_categories/%d", categoryID), category)
}

// DeleteCategory deletes an expense category.
func (s *ExpensesService) DeleteCategory(ctx context.Context, categoryID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("expense_categories/%d", categoryID))
}
