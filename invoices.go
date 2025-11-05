package harvest

import (
	"context"
	"fmt"
	"net/url"
)

// InvoicesService handles communication with the invoice related
// methods of the Harvest API.
type InvoicesService struct {
	client *API
}

// InvoiceListOptions specifies optional parameters to the List method.
type InvoiceListOptions struct {
	ListOptions
	ClientID     int64  `url:"client_id,omitempty"`
	ProjectID    int64  `url:"project_id,omitempty"`
	State        string `url:"state,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
	From         string `url:"from,omitempty"`
	To           string `url:"to,omitempty"`
}

// InvoiceList represents a list of invoices.
type InvoiceList struct {
	Invoices []Invoice `json:"invoices"`
	Paginated[Invoice]
}

// ListPage returns a single page of invoices.
func (s *InvoicesService) ListPage(ctx context.Context, opts *InvoiceListOptions) (*InvoiceList, error) {
	u, err := addOptions("invoices", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var invoices InvoiceList
	_, err = s.client.Do(ctx, req, &invoices)
	if err != nil {
		return nil, err
	}

	// Copy invoices to Items for pagination
	invoices.Items = invoices.Invoices

	return &invoices, nil
}

// List returns all invoices across all pages.
func (s *InvoicesService) List(ctx context.Context, opts *InvoiceListOptions) ([]Invoice, error) {
	if opts == nil {
		opts = &InvoiceListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allInvoices []Invoice

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allInvoices = append(allInvoices, result.Invoices...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allInvoices, nil
}

// Get retrieves a specific invoice.
func (s *InvoicesService) Get(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Get[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID))
}

// InvoiceCreateRequest represents a request to create an invoice.
type InvoiceCreateRequest struct {
	ClientID      int64                    `json:"client_id"`
	EstimateID    int64                    `json:"estimate_id,omitempty"`
	Number        string                   `json:"number,omitempty"`
	PurchaseOrder string                   `json:"purchase_order,omitempty"`
	Tax           float64                  `json:"tax,omitempty"`
	Tax2          float64                  `json:"tax2,omitempty"`
	Discount      float64                  `json:"discount,omitempty"`
	Subject       string                   `json:"subject,omitempty"`
	Notes         string                   `json:"notes,omitempty"`
	Currency      string                   `json:"currency,omitempty"`
	IssueDate     string                   `json:"issue_date,omitempty"`
	DueDate       string                   `json:"due_date,omitempty"`
	PaymentTerm   string                   `json:"payment_term,omitempty"`
	LineItems     []InvoiceLineItemRequest `json:"line_items,omitempty"`
}

// InvoiceLineItemRequest represents a line item in an invoice request.
type InvoiceLineItemRequest struct {
	ProjectID   int64   `json:"project_id,omitempty"`
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Taxed       *bool   `json:"taxed,omitempty"`
	Taxed2      *bool   `json:"taxed2,omitempty"`
}

// Create creates a new invoice.
func (s *InvoicesService) Create(ctx context.Context, invoice *InvoiceCreateRequest) (*Invoice, error) {
	return Create[Invoice](ctx, s.client, "invoices", invoice)
}

// InvoiceUpdateRequest represents a request to update an invoice.
type InvoiceUpdateRequest struct {
	ClientID      int64                    `json:"client_id,omitempty"`
	EstimateID    int64                    `json:"estimate_id,omitempty"`
	Number        string                   `json:"number,omitempty"`
	PurchaseOrder string                   `json:"purchase_order,omitempty"`
	Tax           float64                  `json:"tax,omitempty"`
	Tax2          float64                  `json:"tax2,omitempty"`
	Discount      float64                  `json:"discount,omitempty"`
	Subject       string                   `json:"subject,omitempty"`
	Notes         string                   `json:"notes,omitempty"`
	Currency      string                   `json:"currency,omitempty"`
	IssueDate     string                   `json:"issue_date,omitempty"`
	DueDate       string                   `json:"due_date,omitempty"`
	PaymentTerm   string                   `json:"payment_term,omitempty"`
	LineItems     []InvoiceLineItemRequest `json:"line_items,omitempty"`
}

// Update updates an invoice.
func (s *InvoicesService) Update(ctx context.Context, invoiceID int64, invoice *InvoiceUpdateRequest) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID), invoice)
}

// Delete deletes an invoice.
func (s *InvoicesService) Delete(ctx context.Context, invoiceID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID))
}

// InvoiceMessageRequest represents a request to create an invoice message.
type InvoiceMessageRequest struct {
	EventType string `json:"event_type"`
}

// InvoiceMessageListOptions specifies optional parameters for listing invoice messages.
type InvoiceMessageListOptions struct {
	ListOptions
	UpdatedSince string `url:"updated_since,omitempty"`
}

// InvoiceMessageList represents a list of invoice messages.
type InvoiceMessageList struct {
	InvoiceMessages []InvoiceMessage `json:"invoice_messages"`
	Paginated[InvoiceMessage]
}

// ListMessagesPage returns a single page of messages for an invoice.
func (s *InvoicesService) ListMessagesPage(ctx context.Context, invoiceID int64, opts *InvoiceMessageListOptions) (*InvoiceMessageList, error) {
	u, err := addOptions(fmt.Sprintf("invoices/%d/messages", invoiceID), opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var messages InvoiceMessageList
	_, err = s.client.Do(ctx, req, &messages)
	if err != nil {
		return nil, err
	}

	// Copy messages to Items for pagination
	messages.Items = messages.InvoiceMessages

	return &messages, nil
}

// ListMessages returns all messages for an invoice across all pages.
func (s *InvoicesService) ListMessages(ctx context.Context, invoiceID int64, opts *InvoiceMessageListOptions) ([]InvoiceMessage, error) {
	if opts == nil {
		opts = &InvoiceMessageListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allMessages []InvoiceMessage

	for {
		result, err := s.ListMessagesPage(ctx, invoiceID, opts)
		if err != nil {
			return nil, err
		}

		allMessages = append(allMessages, result.InvoiceMessages...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allMessages, nil
}

// MarkAsSent marks a draft invoice as sent.
func (s *InvoicesService) MarkAsSent(ctx context.Context, invoiceID int64) (*InvoiceMessage, error) {
	req := &InvoiceMessageRequest{EventType: "send"}
	return Create[InvoiceMessage](ctx, s.client, fmt.Sprintf("invoices/%d/messages", invoiceID), req)
}

// MarkAsClosed marks an invoice as closed.
func (s *InvoicesService) MarkAsClosed(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/close", invoiceID), nil)
}

// Reopen reopens a closed invoice.
func (s *InvoicesService) Reopen(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/reopen", invoiceID), nil)
}

// MarkAsDraft marks an open invoice as a draft.
func (s *InvoicesService) MarkAsDraft(ctx context.Context, invoiceID int64) (*InvoiceMessage, error) {
	req := &InvoiceMessageRequest{EventType: "draft"}
	return Create[InvoiceMessage](ctx, s.client, fmt.Sprintf("invoices/%d/messages", invoiceID), req)
}

// InvoiceItemCategoryListOptions specifies optional parameters for listing invoice item categories.
type InvoiceItemCategoryListOptions struct {
	ListOptions
	UpdatedSince string `url:"updated_since,omitempty"`
}

// InvoiceItemCategoryList represents a list of invoice item categories.
type InvoiceItemCategoryList struct {
	InvoiceItemCategories []InvoiceItemCategory `json:"invoice_item_categories"`
	Paginated[InvoiceItemCategory]
}

// ListItemCategoriesPage returns a single page of invoice item categories.
func (s *InvoicesService) ListItemCategoriesPage(ctx context.Context, opts *InvoiceItemCategoryListOptions) (*InvoiceItemCategoryList, error) {
	u, err := addOptions("invoice_item_categories", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var categories InvoiceItemCategoryList
	_, err = s.client.Do(ctx, req, &categories)
	if err != nil {
		return nil, err
	}

	// Copy categories to Items for pagination
	categories.Items = categories.InvoiceItemCategories

	return &categories, nil
}

// ListItemCategories returns all invoice item categories across all pages.
// This endpoint uses cursor-based pagination.
func (s *InvoicesService) ListItemCategories(ctx context.Context, opts *InvoiceItemCategoryListOptions) ([]InvoiceItemCategory, error) {
	if opts == nil {
		opts = &InvoiceItemCategoryListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allCategories []InvoiceItemCategory

	// Fetch first page
	result, err := s.ListItemCategoriesPage(ctx, opts)
	if err != nil {
		return nil, err
	}
	allCategories = append(allCategories, result.InvoiceItemCategories...)

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

			var categories InvoiceItemCategoryList
			_, err = s.client.Do(ctx, req, &categories)
			if err != nil {
				return nil, err
			}
			categories.Items = categories.InvoiceItemCategories
			result = &categories
			allCategories = append(allCategories, categories.InvoiceItemCategories...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListItemCategoriesPage(ctx, opts)
			if err != nil {
				return nil, err
			}
			allCategories = append(allCategories, result.InvoiceItemCategories...)
		} else {
			break
		}
	}

	return allCategories, nil
}

// GetItemCategory retrieves a specific invoice item category.
func (s *InvoicesService) GetItemCategory(ctx context.Context, categoryID int64) (*InvoiceItemCategory, error) {
	return Get[InvoiceItemCategory](ctx, s.client, fmt.Sprintf("invoice_item_categories/%d", categoryID))
}

// InvoiceItemCategoryCreateRequest represents a request to create an invoice item category.
type InvoiceItemCategoryCreateRequest struct {
	Name string `json:"name"`
}

// CreateItemCategory creates a new invoice item category.
func (s *InvoicesService) CreateItemCategory(ctx context.Context, category *InvoiceItemCategoryCreateRequest) (*InvoiceItemCategory, error) {
	return Create[InvoiceItemCategory](ctx, s.client, "invoice_item_categories", category)
}

// InvoiceItemCategoryUpdateRequest represents a request to update an invoice item category.
type InvoiceItemCategoryUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

// UpdateItemCategory updates an invoice item category.
func (s *InvoicesService) UpdateItemCategory(ctx context.Context, categoryID int64, category *InvoiceItemCategoryUpdateRequest) (*InvoiceItemCategory, error) {
	return Update[InvoiceItemCategory](ctx, s.client, fmt.Sprintf("invoice_item_categories/%d", categoryID), category)
}

// DeleteItemCategory deletes an invoice item category.
func (s *InvoicesService) DeleteItemCategory(ctx context.Context, categoryID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("invoice_item_categories/%d", categoryID))
}
