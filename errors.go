package harvest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse represents an error response from the Harvest API.
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"error"`
	Errors   []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"error_description,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if len(e.Errors) > 0 {
		msg := fmt.Sprintf("%v %v: %d %s", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message)
		for _, err := range e.Errors {
			msg += fmt.Sprintf("\n  %s: %s", err.Field, err.Message)
		}
		return msg
	}
	return fmt.Sprintf("%v %v: %d %s", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message)
}

// RateLimitError occurs when the API rate limit is exceeded.
type RateLimitError struct {
	Rate     Rate
	Response *http.Response
	Message  string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %s (rate limit: %d/%d, resets at %s)",
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.StatusCode,
		e.Message,
		e.Rate.Remaining,
		e.Rate.Limit,
		e.Rate.Reset.Time.Format("15:04:05"))
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Check for rate limit
	if r.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{
			Rate:     ParseRate(r),
			Response: r,
			Message:  "API rate limit exceeded",
		}
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	switch r.StatusCode {
	case http.StatusUnauthorized:
		errorResponse.Message = "Authentication failed. Check your access token and account ID."
	case http.StatusForbidden:
		errorResponse.Message = "Access forbidden. You don't have permission to access this resource."
	case http.StatusNotFound:
		errorResponse.Message = "Resource not found."
	case http.StatusUnprocessableEntity:
		if errorResponse.Message == "" {
			errorResponse.Message = "Invalid request. Check your input parameters."
		}
	default:
		if errorResponse.Message == "" {
			errorResponse.Message = fmt.Sprintf("Unexpected status code: %d", r.StatusCode)
		}
	}

	return errorResponse
}
