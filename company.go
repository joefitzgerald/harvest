package harvest

import (
	"context"
)

// CompanyService handles communication with the company related
// methods of the Harvest API.
type CompanyService struct {
	client *API
}

// Get retrieves the company for the currently authenticated user.
func (s *CompanyService) Get(ctx context.Context) (*Company, error) {
	return Get[Company](ctx, s.client, "company")
}