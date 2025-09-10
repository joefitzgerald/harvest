package harvest

import (
	"context"
	"fmt"
)

// EstimatesService handles communication with the estimate related
// methods of the Harvest API.
type EstimatesService struct {
	client *API
}

// EstimateListOptions specifies optional parameters to the List method.
type EstimateListOptions struct {
	ListOptions
	ClientID     int64  `url:"client_id,omitempty"`
	State        string `url:"state,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
	From         string `url:"from,omitempty"`
	To           string `url:"to,omitempty"`
}

// EstimateList represents a list of estimates.
type EstimateList struct {
	Estimates []Estimate `json:"estimates"`
	Paginated[Estimate]
}

// List returns a list of estimates.
func (s *EstimatesService) List(ctx context.Context, opts *EstimateListOptions) (*EstimateList, error) {
	u, err := addOptions("estimates", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var estimates EstimateList
	_, err = s.client.Do(ctx, req, &estimates)
	if err != nil {
		return nil, err
	}

	// Copy estimates to Items for pagination
	estimates.Items = estimates.Estimates
	
	return &estimates, nil
}

// Get retrieves a specific estimate.
func (s *EstimatesService) Get(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Get[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d", estimateID))
}

// EstimateCreateRequest represents a request to create an estimate.
type EstimateCreateRequest struct {
	ClientID      int64                      `json:"client_id"`
	Number        string                     `json:"number,omitempty"`
	PurchaseOrder string                     `json:"purchase_order,omitempty"`
	Tax           float64                    `json:"tax,omitempty"`
	Tax2          float64                    `json:"tax2,omitempty"`
	Discount      float64                    `json:"discount,omitempty"`
	Subject       string                     `json:"subject,omitempty"`
	Notes         string                     `json:"notes,omitempty"`
	Currency      string                     `json:"currency,omitempty"`
	IssueDate     string                     `json:"issue_date,omitempty"`
	LineItems     []EstimateLineItemRequest  `json:"line_items,omitempty"`
}

// EstimateLineItemRequest represents a line item in an estimate request.
type EstimateLineItemRequest struct {
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Taxed       *bool   `json:"taxed,omitempty"`
	Taxed2      *bool   `json:"taxed2,omitempty"`
}

// Create creates a new estimate.
func (s *EstimatesService) Create(ctx context.Context, estimate *EstimateCreateRequest) (*Estimate, error) {
	return Create[Estimate](ctx, s.client, "estimates", estimate)
}

// EstimateUpdateRequest represents a request to update an estimate.
type EstimateUpdateRequest struct {
	ClientID      int64                      `json:"client_id,omitempty"`
	Number        string                     `json:"number,omitempty"`
	PurchaseOrder string                     `json:"purchase_order,omitempty"`
	Tax           float64                    `json:"tax,omitempty"`
	Tax2          float64                    `json:"tax2,omitempty"`
	Discount      float64                    `json:"discount,omitempty"`
	Subject       string                     `json:"subject,omitempty"`
	Notes         string                     `json:"notes,omitempty"`
	Currency      string                     `json:"currency,omitempty"`
	IssueDate     string                     `json:"issue_date,omitempty"`
	LineItems     []EstimateLineItemRequest  `json:"line_items,omitempty"`
}

// Update updates an estimate.
func (s *EstimatesService) Update(ctx context.Context, estimateID int64, estimate *EstimateUpdateRequest) (*Estimate, error) {
	return Update[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d", estimateID), estimate)
}

// Delete deletes an estimate.
func (s *EstimatesService) Delete(ctx context.Context, estimateID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("estimates/%d", estimateID))
}

// MarkAsSent marks an estimate as sent.
func (s *EstimatesService) MarkAsSent(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Update[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d/messages", estimateID), nil)
}

// MarkAsAccepted marks an estimate as accepted.
func (s *EstimatesService) MarkAsAccepted(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Update[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d/accept", estimateID), nil)
}

// MarkAsDeclined marks an estimate as declined.
func (s *EstimatesService) MarkAsDeclined(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Update[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d/decline", estimateID), nil)
}

// Reopen reopens a closed estimate.
func (s *EstimatesService) Reopen(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Update[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d/reopen", estimateID), nil)
}