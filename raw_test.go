package harvest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// newTestAPI returns an API pointed at server with the /v2/ base path Harvest uses.
func newTestAPI(t *testing.T, server *httptest.Server) *API {
	t.Helper()
	api, err := NewWithConfig("token", "12345", "harvest-test (test@example.com)", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	base, err := url.Parse(server.URL + "/v2/")
	if err != nil {
		t.Fatal(err)
	}
	api.baseURL = base
	return api
}

func TestListRawFollowsPagesAndKeepsRaw(t *testing.T) {
	var gotPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.RequestURI())
		if r.Header.Get("Harvest-Account-Id") != "12345" {
			t.Errorf("missing account header")
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("page") {
		case "", "1":
			fmt.Fprintf(w, `{"time_entries":[{"id":1,"hours":1.5,"brand_new_attr":"kept"}],
				"per_page":1,"total_pages":2,"total_entries":2,"next_page":2,"previous_page":null,"page":1,
				"links":{"first":"%[1]s/v2/time_entries?page=1&per_page=1","next":"%[1]s/v2/time_entries?page=2&per_page=1","previous":null,"last":"%[1]s/v2/time_entries?page=2&per_page=1"}}`, "http://api.example")
		case "2":
			fmt.Fprint(w, `{"time_entries":[{"id":2,"hours":2}],
				"per_page":1,"total_pages":2,"total_entries":2,"next_page":null,"previous_page":1,"page":2,
				"links":{"first":"x","next":null,"previous":"y","last":"z"}}`)
		default:
			t.Errorf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	defer server.Close()
	api := newTestAPI(t, server)

	items, err := ListRaw[TimeEntry](context.Background(), api, "time_entries", "time_entries", &ListOptions{PerPage: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 items, got %d", len(items))
	}
	if items[0].Value.ID != 1 || items[1].Value.ID != 2 {
		t.Errorf("ids = %d,%d", items[0].Value.ID, items[1].Value.ID)
	}
	var extra map[string]any
	if err := json.Unmarshal(items[0].Raw, &extra); err != nil {
		t.Fatal(err)
	}
	if extra["brand_new_attr"] != "kept" {
		t.Errorf("raw lost unknown attribute: %s", items[0].Raw)
	}
	if len(gotPaths) != 2 {
		t.Errorf("want 2 requests, got %v", gotPaths)
	}
}

func TestListRawPageNumberFallback(t *testing.T) {
	pages := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages++
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "2" {
			fmt.Fprint(w, `{"roles":[{"id":2,"name":"b"}],"next_page":null,"page":2}`)
			return
		}
		fmt.Fprint(w, `{"roles":[{"id":1,"name":"a"}],"next_page":2,"page":1}`)
	}))
	defer server.Close()
	api := newTestAPI(t, server)

	items, err := ListRaw[Role](context.Background(), api, "roles", "roles", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || pages != 2 {
		t.Fatalf("items=%d pages=%d", len(items), pages)
	}
}

func TestListRawAutoDetectsKeyAndBareLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"pto_assignments":[{"user":{"id":7,"name":"Kim"},"holiday_calendar":null,"work_schedule":{"id":8,"name":"Std"},"effective_work_schedule":{"name":"Std","weekly_hours":40.0}}]}`)
	}))
	defer server.Close()
	api := newTestAPI(t, server)

	items, err := api.PTO.ListAssignments(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].User.ID != 7 || items[0].WorkSchedule.ID != 8 || items[0].HolidayCalendar != nil {
		t.Fatalf("unexpected %+v", items)
	}
}

func TestListRawMissingKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"foo":[],"bar":[]}`)
	}))
	defer server.Close()
	api := newTestAPI(t, server)

	if _, err := ListRaw[Role](context.Background(), api, "roles", "roles", nil); err == nil {
		t.Fatal("expected error for missing key")
	}
	if _, err := ListRaw[Role](context.Background(), api, "roles", "", nil); err == nil {
		t.Fatal("expected error for ambiguous auto-detect")
	}
}

func TestListRawStopsOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "2" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"error":"forbidden"}`)
			return
		}
		fmt.Fprint(w, `{"roles":[{"id":1}],"next_page":2}`)
	}))
	defer server.Close()
	api := newTestAPI(t, server)

	items, err := ListRaw[Role](context.Background(), api, "roles", "roles", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if items != nil {
		t.Errorf("partial results must not be returned on error, got %d", len(items))
	}
}

func TestDateUnmarshalNull(t *testing.T) {
	var v struct {
		A Date  `json:"a"`
		B *Date `json:"b"`
		C Date  `json:"c"`
	}
	if err := json.Unmarshal([]byte(`{"a":null,"b":null,"c":"2026-09-26"}`), &v); err != nil {
		t.Fatal(err)
	}
	if !v.A.IsZero() || v.B != nil || v.C.Format("2006-01-02") != "2026-09-26" {
		t.Errorf("unexpected %+v", v)
	}
}

func TestInvoiceMessageRecipientsAreObjects(t *testing.T) {
	var m InvoiceMessage
	if err := json.Unmarshal([]byte(`{"id":1,"recipients":[{"name":"A","email":"a@example.com"}],"subject":"s","body":"b"}`), &m); err != nil {
		t.Fatal(err)
	}
	if len(m.Recipients) != 1 || m.Recipients[0].Email != "a@example.com" {
		t.Fatalf("recipients = %+v", m.Recipients)
	}
}
