package harvest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPTOAndRatesEndpoints(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/pto/types":
			fmt.Fprint(w, `{"pto_types":[{"id":101,"name":"Vacation","accrual_type":"fixed","default_days_per_year":25.0,"max_carryover_days":5.0,"requires_approval":true,"visible_to_all":true,"is_active":true,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"per_page":2000,"total_pages":1,"total_entries":1,"next_page":null,"previous_page":null,"page":1}`)
		case "/v2/pto/requests":
			fmt.Fprint(w, `{"pto_requests":[{"id":501,"status":"pending","start_date":"2026-10-01","end_date":"2026-10-02","days":2.0,"start_time_minutes":null,"reviewer":null,"user":{"id":1782959,"name":"Kim Allen"},"pto_type":{"id":101,"name":"Vacation"},"request_days":[{"id":1001,"date":"2026-10-01","hours":8.0,"day_fraction":1.0}],"created_at":"2026-09-01T12:00:00Z","updated_at":"2026-09-01T12:00:00Z"}],"next_page":null,"page":1}`)
		case "/v2/pto/balances":
			if r.URL.Query().Get("user_id") != "1782959" || r.URL.Query().Get("year") != "2026" {
				t.Errorf("balances query = %s", r.URL.RawQuery)
			}
			fmt.Fprint(w, `{"user":{"id":1782959,"name":"Kim Allen"},"year":2026,"as_of":"2026-09-03","balances":[{"pto_type":{"id":101,"name":"Vacation"},"is_requestable":true,"year":2026,"as_of":"2026-09-03","is_unlimited":false,"requires_approval":true,"effective_accrual_type":"fixed","accrued_days":"25.0","carryover_days":"2.0","available_days":"27.0","used_days":"10.0","pending_days":"3.0","remaining_days":"14.0","used_hours":"80.0","pending_hours":"24.0","used_fraction":"0.3704","pending_fraction":"0.1111"}]}`)
		case "/v2/pto/holiday_calendars":
			fmt.Fprint(w, `{"pto_holiday_calendars":[{"id":12,"name":"US Holidays","entries_count":1,"entries":[{"id":44,"name":"New Year's Day","date":"2026-01-01"}],"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"next_page":null}`)
		case "/v2/pto/work_schedules":
			fmt.Fprint(w, `{"pto_work_schedules":[{"id":8,"name":"Standard Full Time","is_default":true,"is_active":true,"weekly_hours":40.0,"monday_hours":8.0,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"next_page":null}`)
		case "/v2/pto/allocations":
			fmt.Fprint(w, `{"pto_allocations":[{"id":780,"year":2026,"accrual_type":null,"days_per_year":28.0,"carryover_days":2.0,"requires_approval":false,"effective_days_per_year":28.0,"effective_accrual_type":"fixed","effective_requires_approval":false,"user":{"id":1782959,"name":"Kim Allen"},"pto_type":{"id":101,"name":"Vacation"},"created_at":"2026-09-03T12:00:00Z","updated_at":"2026-09-03T12:00:00Z"}],"next_page":null}`)
		case "/v2/users/3226125/billable_rates":
			fmt.Fprint(w, `{"billable_rates":[{"id":1836493,"amount":8.25,"start_date":"2019-01-01","end_date":"2019-05-31","created_at":"2020-05-01T13:17:42Z","updated_at":"2020-05-01T13:17:50Z"},{"id":1836496,"amount":15.0,"start_date":"2020-05-01","end_date":null,"created_at":"2020-05-01T13:18:10Z","updated_at":"2020-05-01T13:18:10Z"}],"next_page":null}`)
		case "/v2/users/3226125/cost_rates":
			fmt.Fprint(w, `{"cost_rates":[{"id":825304,"amount":15.25,"start_date":"2020-05-01","end_date":null,"created_at":"2020-05-01T13:19:31Z","updated_at":"2020-05-01T13:19:31Z"}],"next_page":null}`)
		case "/v2/users/1782959/teammates":
			fmt.Fprint(w, `{"teammates":[{"id":3230547,"first_name":"Jim","last_name":"Allen","email":"jim@example.com"}],"next_page":null}`)
		case "/v2/invoices/13150378/payments":
			fmt.Fprint(w, `{"invoice_payments":[{"id":10112854,"amount":10700,"paid_at":"2017-02-21T00:00:00Z","paid_date":"2017-02-21","recorded_by":"Alice Doe","transaction_id":null,"created_at":"2017-06-27T16:24:57Z","updated_at":"2017-06-27T16:24:57Z","payment_gateway":{"id":1234,"name":"Linkpoint International"}}],"next_page":null}`)
		case "/v2/user_assignments":
			fmt.Fprint(w, `{"user_assignments":[{"id":1,"project":{"id":10,"name":"P","code":"C"},"user":{"id":20,"name":"U"},"is_active":true,"is_project_manager":false,"use_default_rates":true,"hourly_rate":null,"budget":null,"created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z"}],"next_page":null}`)
		case "/v2/task_assignments":
			fmt.Fprint(w, `{"task_assignments":[{"id":2,"project":{"id":10,"name":"P","code":"C"},"task":{"id":30,"name":"T"},"is_active":true,"billable":true,"created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z"}],"next_page":null}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	api := newTestAPI(t, server)
	ctx := context.Background()

	types, err := api.PTO.ListTypes(ctx, nil)
	if err != nil || len(types) != 1 || types[0].Name != "Vacation" || types[0].DefaultDaysPerYear.String() != "25" {
		t.Fatalf("types: %v %+v", err, types)
	}
	reqs, err := api.PTO.ListRequests(ctx, nil)
	if err != nil || len(reqs) != 1 || reqs[0].User.ID != 1782959 || len(reqs[0].RequestDays) != 1 || reqs[0].Reviewer != nil {
		t.Fatalf("requests: %v %+v", err, reqs)
	}
	bal, err := api.PTO.GetBalances(ctx, &PTOBalanceOptions{UserID: 1782959, Year: 2026})
	if err != nil || bal.User.ID != 1782959 || len(bal.Balances) != 1 || bal.Balances[0].RemainingDays.String() != "14" {
		t.Fatalf("balances: %v %+v", err, bal)
	}
	cals, err := api.PTO.ListHolidayCalendars(ctx, nil)
	if err != nil || len(cals) != 1 || len(cals[0].Entries) != 1 || cals[0].Entries[0].Date.Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("calendars: %v %+v", err, cals)
	}
	ws, err := api.PTO.ListWorkSchedules(ctx, nil)
	if err != nil || len(ws) != 1 || !ws[0].IsDefault {
		t.Fatalf("work schedules: %v %+v", err, ws)
	}
	allocs, err := api.PTO.ListAllocations(ctx, &PTOAllocationListOptions{Year: 2026})
	if err != nil || len(allocs) != 1 || allocs[0].AccrualType != "" || allocs[0].DaysPerYear.String() != "28" {
		t.Fatalf("allocations: %v %+v", err, allocs)
	}
	br, err := api.Users.ListBillableRates(ctx, 3226125, nil)
	if err != nil || len(br) != 2 || br[1].EndDate != nil || br[0].EndDate.Format("2006-01-02") != "2019-05-31" {
		t.Fatalf("billable rates: %v %+v", err, br)
	}
	cr, err := api.Users.ListCostRates(ctx, 3226125, nil)
	if err != nil || len(cr) != 1 || cr[0].Amount.String() != "15.25" {
		t.Fatalf("cost rates: %v %+v", err, cr)
	}
	tm, err := api.Users.ListTeammates(ctx, 1782959, nil)
	if err != nil || len(tm) != 1 || tm[0].Email != "jim@example.com" {
		t.Fatalf("teammates: %v %+v", err, tm)
	}
	pay, err := api.Invoices.ListPayments(ctx, 13150378, nil)
	if err != nil || len(pay) != 1 || pay[0].PaymentGateway.Name != "Linkpoint International" || pay[0].PaidDate.Format("2006-01-02") != "2017-02-21" {
		t.Fatalf("payments: %v %+v", err, pay)
	}
	ua, err := api.Projects.ListAllUserAssignments(ctx, nil)
	if err != nil || len(ua) != 1 || ua[0].Project.ID != 10 || ua[0].User.ID != 20 {
		t.Fatalf("user assignments: %v %+v", err, ua)
	}
	ta, err := api.Projects.ListAllTaskAssignments(ctx, nil)
	if err != nil || len(ta) != 1 || ta[0].Task.ID != 30 {
		t.Fatalf("task assignments: %v %+v", err, ta)
	}
}
