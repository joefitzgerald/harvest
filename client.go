package harvest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.harvestapp.com/v2/"
	defaultTimeout = 30 * time.Second
	// DefaultPerPage is the default number of items to request per page for list operations.
	DefaultPerPage = 2000
)

type API struct {
	httpClient  *http.Client
	baseURL     *url.URL
	accessToken string
	accountID   string
	userAgent   string

	// Service endpoints
	Company     *CompanyService
	Clients     *ClientsService
	Contacts    *ContactsService
	Projects    *ProjectsService
	TimeEntries *TimeEntriesService
	Users       *UsersService
	Tasks       *TasksService
	Invoices    *InvoicesService
	Estimates   *EstimatesService
	Expenses    *ExpensesService
	Reports     *ReportsService
	Roles       *RolesService
}

// New creates a new Harvest API client with the given User-Agent.
// It reads HARVEST_ACCESS_TOKEN and HARVEST_ACCOUNT_ID from environment variables.
func New(userAgent string) (*API, error) {
	accessToken := os.Getenv("HARVEST_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("HARVEST_ACCESS_TOKEN environment variable is required")
	}

	accountID := os.Getenv("HARVEST_ACCOUNT_ID")
	if accountID == "" {
		return nil, fmt.Errorf("HARVEST_ACCOUNT_ID environment variable is required")
	}

	if userAgent == "" {
		return nil, fmt.Errorf("User-Agent is required (format: 'AppName (contact@example.com)')")
	}

	return NewWithConfig(accessToken, accountID, userAgent, nil)
}

// NewWithConfig creates a new Harvest API client with custom configuration.
func NewWithConfig(accessToken, accountID, userAgent string, httpClient *http.Client) (*API, error) {
	if accessToken == "" || accountID == "" || userAgent == "" {
		return nil, fmt.Errorf("accessToken, accountID, and userAgent are required")
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	c := &API{
		httpClient:  httpClient,
		baseURL:     baseURL,
		accessToken: accessToken,
		accountID:   accountID,
		userAgent:   userAgent,
	}

	// Initialize services
	c.Company = &CompanyService{client: c}
	c.Clients = &ClientsService{client: c}
	c.Contacts = &ContactsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.TimeEntries = &TimeEntriesService{client: c}
	c.Users = &UsersService{client: c}
	c.Tasks = &TasksService{client: c}
	c.Invoices = &InvoicesService{client: c}
	c.Estimates = &EstimatesService{client: c}
	c.Expenses = &ExpensesService{client: c}
	c.Reports = &ReportsService{client: c}
	c.Roles = &RolesService{client: c}

	return c, nil
}

// NewRequest creates an API request.
func (c *API) NewRequest(ctx context.Context, method, urlStr string, body any) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Harvest-Account-Id", c.accountID)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *API) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	// Check for API errors
	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// Generic CRUD methods using Go 1.25 generics

// ListPage performs a GET request to list resources with pagination, returning a single page.
func ListPage[T any](ctx context.Context, c *API, path string, opts *ListOptions) (*Paginated[T], error) {
	u, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var result Paginated[T]
	_, err = c.Do(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListPageFromURL performs a GET request using a full pagination URL.
// This is used for cursor-based pagination where the API provides full URLs in the links section.
func ListPageFromURL[T any](ctx context.Context, c *API, fullURL string) (*Paginated[T], error) {
	// Parse the full URL to extract just the path and query
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

	// Use the path and raw query from the parsed URL
	pathAndQuery := u.Path
	if u.RawQuery != "" {
		pathAndQuery += "?" + u.RawQuery
	}

	req, err := c.NewRequest(ctx, "GET", pathAndQuery, nil)
	if err != nil {
		return nil, err
	}

	var result Paginated[T]
	_, err = c.Do(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// List performs a GET request to list all resources across all pages.
// Supports both page-based and cursor-based pagination.
func List[T any](ctx context.Context, c *API, path string, opts *ListOptions) ([]T, error) {
	if opts == nil {
		opts = &ListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allItems []T

	// Fetch first page
	result, err := ListPage[T](ctx, c, path, opts)
	if err != nil {
		return nil, err
	}
	allItems = append(allItems, result.Items...)

	// Continue fetching remaining pages
	for result.HasNextPage() {
		// Check if using cursor-based pagination
		if nextURL := result.GetNextPageURL(); nextURL != "" {
			// Use cursor-based pagination (follow the Links.Next URL)
			result, err = ListPageFromURL[T](ctx, c, nextURL)
			if err != nil {
				return nil, err
			}
		} else if result.NextPage != nil {
			// Use page-based pagination
			opts.Page = *result.NextPage
			result, err = ListPage[T](ctx, c, path, opts)
			if err != nil {
				return nil, err
			}
		} else {
			// Should not reach here if HasNextPage is working correctly
			break
		}
		allItems = append(allItems, result.Items...)
	}

	return allItems, nil
}

// Get performs a GET request to retrieve a single resource.
func Get[T any](ctx context.Context, c *API, path string) (*T, error) {
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result T
	_, err = c.Do(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Create performs a POST request to create a new resource.
func Create[T any](ctx context.Context, c *API, path string, body any) (*T, error) {
	req, err := c.NewRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}

	var result T
	_, err = c.Do(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Update performs a PATCH request to update an existing resource.
func Update[T any](ctx context.Context, c *API, path string, body any) (*T, error) {
	req, err := c.NewRequest(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var result T
	_, err = c.Do(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete performs a DELETE request to remove a resource.
func Delete(ctx context.Context, c *API, path string) error {
	req, err := c.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(ctx, req, nil)
	return err
}

// addOptions adds the parameters in opts as URL query parameters to s.
func addOptions(s string, opts any) (string, error) {
	v, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}
