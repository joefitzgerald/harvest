package harvest

import (
	"context"
	"fmt"
	"net/url"
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

// ListPage returns a single page of estimates.
func (s *EstimatesService) ListPage(ctx context.Context, opts *EstimateListOptions) (*EstimateList, error) {
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

// List returns all estimates across all pages.
func (s *EstimatesService) List(ctx context.Context, opts *EstimateListOptions) ([]Estimate, error) {
	if opts == nil {
		opts = &EstimateListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allEstimates []Estimate

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allEstimates = append(allEstimates, result.Estimates...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allEstimates, nil
}

// Get retrieves a specific estimate.
func (s *EstimatesService) Get(ctx context.Context, estimateID int64) (*Estimate, error) {
	return Get[Estimate](ctx, s.client, fmt.Sprintf("estimates/%d", estimateID))
}

// EstimateCreateRequest represents a request to create an estimate.
type EstimateCreateRequest struct {
	ClientID      int64                     `json:"client_id"`
	Number        string                    `json:"number,omitempty"`
	PurchaseOrder string                    `json:"purchase_order,omitempty"`
	Tax           float64                   `json:"tax,omitempty"`
	Tax2          float64                   `json:"tax2,omitempty"`
	Discount      float64                   `json:"discount,omitempty"`
	Subject       string                    `json:"subject,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Currency      string                    `json:"currency,omitempty"`
	IssueDate     string                    `json:"issue_date,omitempty"`
	LineItems     []EstimateLineItemRequest `json:"line_items,omitempty"`
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
	ClientID      int64                     `json:"client_id,omitempty"`
	Number        string                    `json:"number,omitempty"`
	PurchaseOrder string                    `json:"purchase_order,omitempty"`
	Tax           float64                   `json:"tax,omitempty"`
	Tax2          float64                   `json:"tax2,omitempty"`
	Discount      float64                   `json:"discount,omitempty"`
	Subject       string                    `json:"subject,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Currency      string                    `json:"currency,omitempty"`
	IssueDate     string                    `json:"issue_date,omitempty"`
	LineItems     []EstimateLineItemRequest `json:"line_items,omitempty"`
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

// EstimateItemCategoryListOptions specifies optional parameters for listing estimate item categories.
type EstimateItemCategoryListOptions struct {
	ListOptions
	UpdatedSince string `url:"updated_since,omitempty"`
}

// EstimateItemCategoryList represents a list of estimate item categories.
type EstimateItemCategoryList struct {
	EstimateItemCategories []EstimateItemCategory `json:"estimate_item_categories"`
	Paginated[EstimateItemCategory]
}

// ListItemCategoriesPage returns a single page of estimate item categories.
func (s *EstimatesService) ListItemCategoriesPage(ctx context.Context, opts *EstimateItemCategoryListOptions) (*EstimateItemCategoryList, error) {
	u, err := addOptions("estimate_item_categories", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var categories EstimateItemCategoryList
	_, err = s.client.Do(ctx, req, &categories)
	if err != nil {
		return nil, err
	}

	// Copy categories to Items for pagination
	categories.Items = categories.EstimateItemCategories

	return &categories, nil
}

// ListItemCategories returns all estimate item categories across all pages.
// This endpoint uses cursor-based pagination.
func (s *EstimatesService) ListItemCategories(ctx context.Context, opts *EstimateItemCategoryListOptions) ([]EstimateItemCategory, error) {
	if opts == nil {
		opts = &EstimateItemCategoryListOptions{}
	}
	// Don't set Page - it's deprecated for cursor-based pagination
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allCategories []EstimateItemCategory

	// Fetch first page
	result, err := s.ListItemCategoriesPage(ctx, opts)
	if err != nil {
		return nil, err
	}
	allCategories = append(allCategories, result.EstimateItemCategories...)

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

			var categories EstimateItemCategoryList
			_, err = s.client.Do(ctx, req, &categories)
			if err != nil {
				return nil, err
			}
			categories.Items = categories.EstimateItemCategories
			result = &categories
			allCategories = append(allCategories, categories.EstimateItemCategories...)
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = s.ListItemCategoriesPage(ctx, opts)
			if err != nil {
				return nil, err
			}
			allCategories = append(allCategories, result.EstimateItemCategories...)
		} else {
			break
		}
	}

	return allCategories, nil
}

// GetItemCategory retrieves a specific estimate item category.
func (s *EstimatesService) GetItemCategory(ctx context.Context, categoryID int64) (*EstimateItemCategory, error) {
	return Get[EstimateItemCategory](ctx, s.client, fmt.Sprintf("estimate_item_categories/%d", categoryID))
}

// EstimateItemCategoryCreateRequest represents a request to create an estimate item category.
type EstimateItemCategoryCreateRequest struct {
	Name string `json:"name"`
}

// CreateItemCategory creates a new estimate item category.
func (s *EstimatesService) CreateItemCategory(ctx context.Context, category *EstimateItemCategoryCreateRequest) (*EstimateItemCategory, error) {
	return Create[EstimateItemCategory](ctx, s.client, "estimate_item_categories", category)
}

// EstimateItemCategoryUpdateRequest represents a request to update an estimate item category.
type EstimateItemCategoryUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

// UpdateItemCategory updates an estimate item category.
func (s *EstimatesService) UpdateItemCategory(ctx context.Context, categoryID int64, category *EstimateItemCategoryUpdateRequest) (*EstimateItemCategory, error) {
	return Update[EstimateItemCategory](ctx, s.client, fmt.Sprintf("estimate_item_categories/%d", categoryID), category)
}

// DeleteItemCategory deletes an estimate item category.
func (s *EstimatesService) DeleteItemCategory(ctx context.Context, categoryID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("estimate_item_categories/%d", categoryID))
}
