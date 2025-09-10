package harvest

import (
	"context"
	"fmt"
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

// List returns a list of invoices.
func (s *InvoicesService) List(ctx context.Context, opts *InvoiceListOptions) (*InvoiceList, error) {
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

// Get retrieves a specific invoice.
func (s *InvoicesService) Get(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Get[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID))
}

// InvoiceCreateRequest represents a request to create an invoice.
type InvoiceCreateRequest struct {
	ClientID      int64                     `json:"client_id"`
	EstimateID    int64                     `json:"estimate_id,omitempty"`
	Number        string                    `json:"number,omitempty"`
	PurchaseOrder string                    `json:"purchase_order,omitempty"`
	Tax           float64                   `json:"tax,omitempty"`
	Tax2          float64                   `json:"tax2,omitempty"`
	Discount      float64                   `json:"discount,omitempty"`
	Subject       string                    `json:"subject,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Currency      string                    `json:"currency,omitempty"`
	IssueDate     string                    `json:"issue_date,omitempty"`
	DueDate       string                    `json:"due_date,omitempty"`
	PaymentTerm   string                    `json:"payment_term,omitempty"`
	LineItems     []InvoiceLineItemRequest  `json:"line_items,omitempty"`
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
	ClientID      int64                     `json:"client_id,omitempty"`
	EstimateID    int64                     `json:"estimate_id,omitempty"`
	Number        string                    `json:"number,omitempty"`
	PurchaseOrder string                    `json:"purchase_order,omitempty"`
	Tax           float64                   `json:"tax,omitempty"`
	Tax2          float64                   `json:"tax2,omitempty"`
	Discount      float64                   `json:"discount,omitempty"`
	Subject       string                    `json:"subject,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Currency      string                    `json:"currency,omitempty"`
	IssueDate     string                    `json:"issue_date,omitempty"`
	DueDate       string                    `json:"due_date,omitempty"`
	PaymentTerm   string                    `json:"payment_term,omitempty"`
	LineItems     []InvoiceLineItemRequest  `json:"line_items,omitempty"`
}

// Update updates an invoice.
func (s *InvoicesService) Update(ctx context.Context, invoiceID int64, invoice *InvoiceUpdateRequest) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID), invoice)
}

// Delete deletes an invoice.
func (s *InvoicesService) Delete(ctx context.Context, invoiceID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("invoices/%d", invoiceID))
}

// MarkAsSent marks an invoice as sent.
func (s *InvoicesService) MarkAsSent(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/messages", invoiceID), nil)
}

// MarkAsClosed marks an invoice as closed.
func (s *InvoicesService) MarkAsClosed(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/close", invoiceID), nil)
}

// Reopen reopens a closed invoice.
func (s *InvoicesService) Reopen(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/reopen", invoiceID), nil)
}

// MarkAsDraft marks an invoice as draft.
func (s *InvoicesService) MarkAsDraft(ctx context.Context, invoiceID int64) (*Invoice, error) {
	return Update[Invoice](ctx, s.client, fmt.Sprintf("invoices/%d/draft", invoiceID), nil)
}