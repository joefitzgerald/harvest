package harvest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// ListOptions specifies optional parameters to List methods.
type ListOptions struct {
	Page        int       `url:"page,omitempty"`
	PerPage     int       `url:"per_page,omitempty"`
	UpdatedSince *time.Time `url:"updated_since,omitempty"`
}

// Paginated represents a paginated response from the Harvest API.
type Paginated[T any] struct {
	Items      []T              `json:"-"`
	Links      *PaginationLinks `json:"links"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
	TotalEntries int            `json:"total_entries"`
	NextPage   *int             `json:"next_page"`
	PreviousPage *int           `json:"previous_page"`
	Page       int              `json:"page"`
	
	// The actual items will be in a field named after the resource type
	// We'll handle this with custom unmarshaling or in resource-specific methods
}

// PaginationLinks represents pagination links in API responses.
type PaginationLinks struct {
	First    string `json:"first"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	Last     string `json:"last"`
}

// HasNextPage returns true if there is a next page of results.
func (p *Paginated[T]) HasNextPage() bool {
	return p.NextPage != nil
}

// HasPreviousPage returns true if there is a previous page of results.
func (p *Paginated[T]) HasPreviousPage() bool {
	return p.PreviousPage != nil
}

// Iterator provides iteration over paginated results.
type Iterator[T any] struct {
	client   *API
	ctx      context.Context
	path     string
	opts     *ListOptions
	current  *Paginated[T]
	index    int
	fetcher  func(context.Context, *API, string, *ListOptions) (*Paginated[T], error)
}

// NewIterator creates a new iterator for paginated results.
func NewIterator[T any](ctx context.Context, client *API, path string, opts *ListOptions, 
	fetcher func(context.Context, *API, string, *ListOptions) (*Paginated[T], error)) *Iterator[T] {
	if opts == nil {
		opts = &ListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = 100 // Default page size
	}
	
	return &Iterator[T]{
		client:  client,
		ctx:     ctx,
		path:    path,
		opts:    opts,
		fetcher: fetcher,
	}
}

// Next returns the next item in the iteration.
func (it *Iterator[T]) Next() (*T, error) {
	// Fetch first page if not loaded
	if it.current == nil {
		page, err := it.fetcher(it.ctx, it.client, it.path, it.opts)
		if err != nil {
			return nil, err
		}
		it.current = page
		it.index = 0
	}

	// Check if we need to fetch the next page
	if it.index >= len(it.current.Items) {
		if !it.current.HasNextPage() {
			return nil, nil // End of iteration
		}
		
		it.opts.Page = *it.current.NextPage
		page, err := it.fetcher(it.ctx, it.client, it.path, it.opts)
		if err != nil {
			return nil, err
		}
		it.current = page
		it.index = 0
	}

	// Return current item and advance
	if it.index < len(it.current.Items) {
		item := &it.current.Items[it.index]
		it.index++
		return item, nil
	}

	return nil, nil
}

// All fetches all pages and returns all items.
func (it *Iterator[T]) All() ([]T, error) {
	var allItems []T
	
	for {
		item, err := it.Next()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		allItems = append(allItems, *item)
	}
	
	return allItems, nil
}

// Rate represents the rate limit for the Harvest API.
type Rate struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     Timestamp `json:"reset"`
}

// Timestamp represents a time that can be unmarshalled from a JSON number.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Harvest sends timestamps as Unix timestamps
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

// ParseRate parses the rate limit headers from an HTTP response.
func ParseRate(r *http.Response) Rate {
	rate := Rate{}
	
	if limit := r.Header.Get("X-RateLimit-Limit"); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get("X-RateLimit-Reset"); reset != "" {
		if timestamp, err := strconv.ParseInt(reset, 10, 64); err == nil {
			rate.Reset = Timestamp{time.Unix(timestamp, 0)}
		}
	}
	
	return rate
}