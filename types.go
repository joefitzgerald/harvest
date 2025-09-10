package harvest

import (
	"time"

	"github.com/shopspring/decimal"
)

// Company represents a company in Harvest.
type Company struct {
	BaseURI                          string    `json:"base_uri"`
	FullDomain                       string    `json:"full_domain"`
	Name                             string    `json:"name"`
	IsActive                         bool      `json:"is_active"`
	WeekStartDay                     string    `json:"week_start_day"`
	WantsTimestampTimers             bool      `json:"wants_timestamp_timers"`
	TimeFormat                       string    `json:"time_format"`
	DateFormat                       string    `json:"date_format"`
	PlanType                         string    `json:"plan_type"`
	Clock                            string    `json:"clock"`
	DecimalSymbol                    string    `json:"decimal_symbol"`
	ThousandsSeparator               string    `json:"thousands_separator"`
	ColorScheme                      string    `json:"color_scheme"`
	WeeklyCapacity                   int       `json:"weekly_capacity"`
	ExpenseFeature                   bool      `json:"expense_feature"`
	InvoiceFeature                   bool      `json:"invoice_feature"`
	EstimateFeature                  bool      `json:"estimate_feature"`
	ApprovalFeature                  bool      `json:"approval_feature"`
}

// Client represents a client in Harvest.
type Client struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	Address     string     `json:"address,omitempty"`
	Currency    string     `json:"currency,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Contact represents a client contact in Harvest.
type Contact struct {
	ID          int64      `json:"id"`
	ClientID    int64      `json:"client_id"`
	Client      *Client    `json:"client,omitempty"`
	Title       string     `json:"title,omitempty"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name,omitempty"`
	Email       string     `json:"email,omitempty"`
	PhoneOffice string     `json:"phone_office,omitempty"`
	PhoneMobile string     `json:"phone_mobile,omitempty"`
	Fax         string     `json:"fax,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Project represents a project in Harvest.
type Project struct {
	ID                          int64               `json:"id"`
	Client                      *Client      `json:"client"`
	Name                        string              `json:"name"`
	Code                        string              `json:"code,omitempty"`
	IsActive                    bool                `json:"is_active"`
	IsBillable                  bool                `json:"is_billable"`
	IsFixedFee                  bool                `json:"is_fixed_fee"`
	BillBy                      string              `json:"bill_by"`
	Budget                      *decimal.Decimal    `json:"budget,omitempty"`
	BudgetBy                    string              `json:"budget_by,omitempty"`
	BudgetIsMonthly             bool                `json:"budget_is_monthly"`
	NotifyWhenOverBudget        bool                `json:"notify_when_over_budget"`
	OverBudgetNotificationPercentage decimal.Decimal `json:"over_budget_notification_percentage,omitempty"`
	ShowBudgetToAll             bool                `json:"show_budget_to_all"`
	CostBudget                  *decimal.Decimal    `json:"cost_budget,omitempty"`
	CostBudgetIncludeExpenses   bool                `json:"cost_budget_include_expenses"`
	HourlyRate                  *decimal.Decimal    `json:"hourly_rate,omitempty"`
	Fee                         *decimal.Decimal    `json:"fee,omitempty"`
	Notes                       string              `json:"notes,omitempty"`
	StartsOn                    *Date               `json:"starts_on,omitempty"`
	EndsOn                      *Date               `json:"ends_on,omitempty"`
	CreatedAt                   time.Time           `json:"created_at"`
	UpdatedAt                   time.Time           `json:"updated_at"`
}

// ProjectUserAssignment represents a user assignment to a project.
type ProjectUserAssignment struct {
	ID              int64            `json:"id"`
	Project         *Project         `json:"project"`
	User            *User            `json:"user"`
	IsActive        bool             `json:"is_active"`
	IsProjectManager bool            `json:"is_project_manager"`
	UseDefaultRates bool             `json:"use_default_rates"`
	HourlyRate      *decimal.Decimal `json:"hourly_rate,omitempty"`
	Budget          *decimal.Decimal `json:"budget,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// ProjectTaskAssignment represents a task assignment to a project.
type ProjectTaskAssignment struct {
	ID              int64            `json:"id"`
	Project         *Project         `json:"project"`
	Task            *Task            `json:"task"`
	IsActive        bool             `json:"is_active"`
	Billable        bool             `json:"billable"`
	HourlyRate      *decimal.Decimal `json:"hourly_rate,omitempty"`
	Budget          *decimal.Decimal `json:"budget,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// User represents a user in Harvest.
type User struct {
	ID                      int64            `json:"id"`
	FirstName               string           `json:"first_name"`
	LastName                string           `json:"last_name"`
	Email                   string           `json:"email"`
	Telephone               string           `json:"telephone,omitempty"`
	Timezone                string           `json:"timezone"`
	HasAccessToAllFutureProjects bool       `json:"has_access_to_all_future_projects"`
	IsContractor            bool             `json:"is_contractor"`
	IsActive                bool             `json:"is_active"`
	WeeklyCapacity          int              `json:"weekly_capacity"`
	DefaultHourlyRate       *decimal.Decimal `json:"default_hourly_rate,omitempty"`
	CostRate                *decimal.Decimal `json:"cost_rate,omitempty"`
	Roles                   []string         `json:"roles"`
	AvatarURL               string           `json:"avatar_url"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}

// Task represents a task in Harvest.
type Task struct {
	ID                int64            `json:"id"`
	Name              string           `json:"name"`
	BillableByDefault bool             `json:"billable_by_default"`
	DefaultHourlyRate *decimal.Decimal `json:"default_hourly_rate,omitempty"`
	IsDefault         bool             `json:"is_default"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// TimeEntry represents a time entry in Harvest.
type TimeEntry struct {
	ID               int64                  `json:"id"`
	SpentDate        Date                   `json:"spent_date"`
	User             *User                  `json:"user"`
	Client           *Client         `json:"client"`
	Project          *Project               `json:"project"`
	Task             *Task                  `json:"task"`
	UserAssignment   *ProjectUserAssignment `json:"user_assignment"`
	TaskAssignment   *ProjectTaskAssignment `json:"task_assignment"`
	Invoice          *Invoice               `json:"invoice,omitempty"`
	Hours            decimal.Decimal        `json:"hours"`
	RoundedHours     decimal.Decimal        `json:"rounded_hours"`
	Notes            string                 `json:"notes,omitempty"`
	IsLocked         bool                   `json:"is_locked"`
	LockedReason     string                 `json:"locked_reason,omitempty"`
	IsClosed         bool                   `json:"is_closed"`
	IsBilled         bool                   `json:"is_billed"`
	TimerStartedAt   *time.Time             `json:"timer_started_at,omitempty"`
	StartedTime      string                 `json:"started_time,omitempty"`
	EndedTime        string                 `json:"ended_time,omitempty"`
	IsRunning        bool                   `json:"is_running"`
	Billable         bool                   `json:"billable"`
	Budgeted         bool                   `json:"budgeted"`
	BillableRate     *decimal.Decimal       `json:"billable_rate,omitempty"`
	CostRate         *decimal.Decimal       `json:"cost_rate,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	ExternalReference *ExternalReference    `json:"external_reference,omitempty"`
}

// ExternalReference represents an external reference for a time entry.
type ExternalReference struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	AccountID string `json:"account_id"`
	Permalink string `json:"permalink"`
	Service   string `json:"service"`
	ServiceIconURL string `json:"service_icon_url"`
}

// Invoice represents an invoice in Harvest.
type Invoice struct {
	ID                 int64           `json:"id"`
	Client             *Client  `json:"client"`
	LineItems          []InvoiceItem   `json:"line_items"`
	Estimate           *Estimate       `json:"estimate,omitempty"`
	Number             string          `json:"number"`
	PurchaseOrder      string          `json:"purchase_order,omitempty"`
	Amount             decimal.Decimal `json:"amount"`
	DueAmount          decimal.Decimal `json:"due_amount"`
	Tax                *decimal.Decimal `json:"tax,omitempty"`
	TaxAmount          *decimal.Decimal `json:"tax_amount,omitempty"`
	Tax2               *decimal.Decimal `json:"tax2,omitempty"`
	Tax2Amount         *decimal.Decimal `json:"tax2_amount,omitempty"`
	Discount           *decimal.Decimal `json:"discount,omitempty"`
	DiscountAmount     *decimal.Decimal `json:"discount_amount,omitempty"`
	Subject            string          `json:"subject,omitempty"`
	Notes              string          `json:"notes,omitempty"`
	Currency           string          `json:"currency"`
	State              string          `json:"state"`
	PeriodStart        *Date           `json:"period_start,omitempty"`
	PeriodEnd          *Date           `json:"period_end,omitempty"`
	IssueDate          Date            `json:"issue_date"`
	DueDate            *Date           `json:"due_date,omitempty"`
	PaymentTerm        string          `json:"payment_term,omitempty"`
	SentAt             *time.Time      `json:"sent_at,omitempty"`
	PaidAt             *time.Time      `json:"paid_at,omitempty"`
	ClosedAt           *time.Time      `json:"closed_at,omitempty"`
	RecurringInvoiceID *int64          `json:"recurring_invoice_id,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// InvoiceItem represents a line item on an invoice.
type InvoiceItem struct {
	ID          int64           `json:"id"`
	Project     *Project        `json:"project,omitempty"`
	Kind        string          `json:"kind"`
	Description string          `json:"description"`
	Quantity    decimal.Decimal `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Amount      decimal.Decimal `json:"amount"`
	Taxed       bool            `json:"taxed"`
	Taxed2      bool            `json:"taxed2"`
}

// Estimate represents an estimate in Harvest.
type Estimate struct {
	ID            int64           `json:"id"`
	Client        *Client  `json:"client"`
	LineItems     []EstimateItem  `json:"line_items"`
	Number        string          `json:"number"`
	PurchaseOrder string          `json:"purchase_order,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
	Tax           *decimal.Decimal `json:"tax,omitempty"`
	TaxAmount     *decimal.Decimal `json:"tax_amount,omitempty"`
	Tax2          *decimal.Decimal `json:"tax2,omitempty"`
	Tax2Amount    *decimal.Decimal `json:"tax2_amount,omitempty"`
	Discount      *decimal.Decimal `json:"discount,omitempty"`
	DiscountAmount *decimal.Decimal `json:"discount_amount,omitempty"`
	Subject       string          `json:"subject,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	Currency      string          `json:"currency"`
	State         string          `json:"state"`
	IssueDate     Date            `json:"issue_date"`
	SentAt        *time.Time      `json:"sent_at,omitempty"`
	AcceptedAt    *time.Time      `json:"accepted_at,omitempty"`
	DeclinedAt    *time.Time      `json:"declined_at,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// EstimateItem represents a line item on an estimate.
type EstimateItem struct {
	ID          int64           `json:"id"`
	Project     *Project        `json:"project,omitempty"`
	Kind        string          `json:"kind"`
	Description string          `json:"description"`
	Quantity    decimal.Decimal `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Amount      decimal.Decimal `json:"amount"`
	Taxed       bool            `json:"taxed"`
	Taxed2      bool            `json:"taxed2"`
}

// Expense represents an expense in Harvest.
type Expense struct {
	ID               int64           `json:"id"`
	Client           *Client  `json:"client"`
	Project          *Project        `json:"project"`
	ExpenseCategory  *ExpenseCategory `json:"expense_category"`
	User             *User           `json:"user"`
	UserAssignment   *ProjectUserAssignment `json:"user_assignment"`
	Invoice          *Invoice        `json:"invoice,omitempty"`
	Receipt          *Receipt        `json:"receipt,omitempty"`
	Notes            string          `json:"notes,omitempty"`
	IsLocked         bool            `json:"is_locked"`
	LockedReason     string          `json:"locked_reason,omitempty"`
	IsClosed         bool            `json:"is_closed"`
	IsBilled         bool            `json:"is_billed"`
	Billable         bool            `json:"billable"`
	SpentDate        Date            `json:"spent_date"`
	TotalCost        decimal.Decimal `json:"total_cost"`
	Units            *decimal.Decimal `json:"units,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// ExpenseCategory represents an expense category.
type ExpenseCategory struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	UnitName  string          `json:"unit_name,omitempty"`
	UnitPrice *decimal.Decimal `json:"unit_price,omitempty"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Receipt represents a receipt attachment.
type Receipt struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	FileSize int    `json:"file_size"`
	ContentType string `json:"content_type"`
}

// Role represents a role in Harvest.
type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UserIDs   []int64   `json:"user_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Date represents a date in YYYY-MM-DD format.
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Date.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // Remove quotes
	
	if s == "" || s == "null" {
		return nil
	}
	
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	
	d.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler for Date.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

// String returns the date as a string in YYYY-MM-DD format.
func (d Date) String() string {
	return d.Format("2006-01-02")
}