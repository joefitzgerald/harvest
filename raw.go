package harvest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// RawItem pairs a decoded API object with the exact JSON the API returned for it.
// Callers that cache Harvest data can persist Raw so no attribute is lost even when
// Value's struct has no field for it.
type RawItem[T any] struct {
	Value T
	Raw   json.RawMessage
}

// rawPage decodes a list response twice: once for the pagination metadata and once as a
// map so the resource array can be located by key without knowing the resource type.
type rawPage struct {
	meta   Paginated[json.RawMessage]
	fields map[string]json.RawMessage
}

func (p *rawPage) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &p.meta); err != nil {
		return err
	}
	return json.Unmarshal(b, &p.fields)
}

// items returns the raw elements of the resource array. When key is empty the sole
// top-level array (other than pagination fields) is used.
func (p *rawPage) items(key string) ([]json.RawMessage, error) {
	if key == "" {
		var found string
		for k, v := range p.fields {
			if len(v) > 0 && v[0] == '[' && k != "links" {
				if found != "" {
					return nil, fmt.Errorf("harvest: response has multiple arrays (%q, %q); specify a key", found, k)
				}
				found = k
			}
		}
		if found == "" {
			return nil, fmt.Errorf("harvest: response has no array to list")
		}
		key = found
	}
	arr, ok := p.fields[key]
	if !ok {
		return nil, fmt.Errorf("harvest: response has no %q array", key)
	}
	var items []json.RawMessage
	if err := json.Unmarshal(arr, &items); err != nil {
		return nil, fmt.Errorf("harvest: decoding %q array: %w", key, err)
	}
	return items, nil
}

// ListRaw fetches every page of a list endpoint and returns each element decoded into T
// together with its raw JSON. path is relative to the API base (for example
// "time_entries" or "users/42/billable_rates"); key names the resource array in the
// response ("time_entries") and may be empty to auto-detect it; opts is any
// go-querystring-tagged options struct or nil for no query parameters.
//
// Both page-number and cursor (links.next) pagination are followed. Endpoints that
// return a bare object array without pagination metadata yield a single page.
func ListRaw[T any](ctx context.Context, c *API, path, key string, opts any) ([]RawItem[T], error) {
	u, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}

	var all []RawItem[T]
	for {
		req, err := c.NewRequest(ctx, "GET", u, nil)
		if err != nil {
			return nil, err
		}

		var page rawPage
		if _, err := c.Do(ctx, req, &page); err != nil {
			return nil, err
		}

		raws, err := page.items(key)
		if err != nil {
			return nil, err
		}
		for _, r := range raws {
			var v T
			if err := json.Unmarshal(r, &v); err != nil {
				return nil, fmt.Errorf("harvest: decoding %s item: %w", path, err)
			}
			all = append(all, RawItem[T]{Value: v, Raw: r})
		}

		next, err := nextPageURL(u, &page.meta)
		if err != nil {
			return nil, err
		}
		if next == "" {
			return all, nil
		}
		u = next
	}
}

// nextPageURL derives the request URL (path + query, relative to the API base) for the
// page after the one described by meta, or "" when there is none.
func nextPageURL(current string, meta *Paginated[json.RawMessage]) (string, error) {
	if meta.Links != nil && meta.Links.Next != "" {
		nu, err := url.Parse(meta.Links.Next)
		if err != nil {
			return "", err
		}
		if nu.RawQuery != "" {
			return nu.Path + "?" + nu.RawQuery, nil
		}
		return nu.Path, nil
	}
	if meta.NextPage != nil {
		cu, err := url.Parse(current)
		if err != nil {
			return "", err
		}
		q := cu.Query()
		q.Set("page", strconv.Itoa(*meta.NextPage))
		cu.RawQuery = q.Encode()
		return cu.String(), nil
	}
	return "", nil
}

// GetRaw fetches a single object and returns it decoded into T with its raw JSON.
func GetRaw[T any](ctx context.Context, c *API, path string, opts any) (RawItem[T], error) {
	var item RawItem[T]
	u, err := addOptions(path, opts)
	if err != nil {
		return item, err
	}
	req, err := c.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return item, err
	}
	var raw json.RawMessage
	if _, err := c.Do(ctx, req, &raw); err != nil {
		return item, err
	}
	if err := json.Unmarshal(raw, &item.Value); err != nil {
		return item, fmt.Errorf("harvest: decoding %s: %w", path, err)
	}
	item.Raw = raw
	return item, nil
}

// listValues is ListRaw without the raw payloads, for typed convenience methods.
func listValues[T any](ctx context.Context, c *API, path, key string, opts any) ([]T, error) {
	items, err := ListRaw[T](ctx, c, path, key, opts)
	if err != nil {
		return nil, err
	}
	out := make([]T, len(items))
	for i := range items {
		out[i] = items[i].Value
	}
	return out, nil
}
