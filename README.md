# Harvest API v2 Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/joefitzgerald/harvest.svg)](https://pkg.go.dev/github.com/joefitzgerald/harvest)
[![Go Report Card](https://goreportcard.com/badge/github.com/joefitzgerald/harvest)](https://goreportcard.com/report/github.com/joefitzgerald/harvest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern, comprehensive Go client library for the [Harvest API v2](https://help.getharvest.com/api-v2/). This library provides complete coverage of all Harvest API v2 endpoints with a clean, type-safe interface.

## Why This Library?

- **Full API v2 Support**: Unlike older libraries that only support Harvest API v1 (deprecated), this library is built exclusively for the current Harvest API v2
- **Modern Go**: Uses Go 1.25+ generics to reduce code duplication and provide a cleaner API
- **Complete Coverage**: Implements all 69+ endpoints from the Harvest API v2
- **Type Safety**: Strongly typed request and response models with proper JSON marshaling
- **Production Ready**: Automatic rate limiting, context support, and comprehensive error handling

> **Note**: If you're looking for a Harvest API v1 client, see [adlio/harvest](https://github.com/adlio/harvest). However, be aware that Harvest API v1 is deprecated and that library hasn't been updated since 2018.

## Features

- ✅ **Complete API Coverage**: All Harvest API v2 endpoints implemented
- ✅ **Type Safety**: Strongly typed request and response models
- ✅ **Generic CRUD Operations**: Leverages Go generics for reduced code duplication
- ✅ **Automatic Pagination**: Built-in pagination support with iterator pattern
- ✅ **Rate Limiting**: Automatic rate limit detection and error handling
- ✅ **Decimal Precision**: Uses `shopspring/decimal` for accurate financial calculations
- ✅ **Context Support**: All methods accept `context.Context` for cancellation
- ✅ **Configurable User-Agent**: Required by Harvest API for identification

## Installation

```bash
go get github.com/joefitzgerald/harvest
```

Requires Go 1.25 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/joefitzgerald/harvest"
)

func main() {
    // Set environment variables:
    // export HARVEST_ACCESS_TOKEN="your-token"
    // export HARVEST_ACCOUNT_ID="your-account-id"
    
    // Create client with required User-Agent
    client, err := harvest.New("MyApp (contact@example.com)")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Get company info
    company, err := client.Company.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Company: %s\n", company.Name)
    
    // List active projects
    projects, err := client.Projects.List(ctx, &harvest.ProjectListOptions{
        ListOptions: harvest.ListOptions{
            Page:    1,
            PerPage: 100,
        },
        IsActive: &[]bool{true}[0],
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, project := range projects.Projects {
        fmt.Printf("Project: %s (Client: %s)\n", project.Name, project.Client.Name)
    }
}
```

## Authentication

The client requires two environment variables:
- `HARVEST_ACCESS_TOKEN`: Your personal access token from Harvest
- `HARVEST_ACCOUNT_ID`: Your Harvest account ID

To get these:
1. Log in to Harvest
2. Go to [Developers > Personal Access Tokens](https://id.getharvest.com/developers)
3. Create a new token
4. Note your Account ID from the same page

You can also create a client with explicit credentials:

```go
client, err := harvest.NewWithConfig(
    "your-access-token",
    "your-account-id", 
    "MyApp (contact@example.com)",
    nil, // optional custom HTTP client
)
```

## User-Agent Requirement

Harvest requires a User-Agent header that includes:
1. The name of your application
2. A link to your application or email address

Examples:
- `"MyTimeTracker (https://example.com)"`
- `"John's Integration (john@example.com)"`

## Available Services

### Core Resources
- **Company**: Get company information
- **Clients**: Manage clients and contacts
- **Projects**: Manage projects and assignments
- **Users**: Manage users and their assignments
- **Tasks**: Manage tasks
- **TimeEntries**: Track time entries

### Financial Resources
- **Invoices**: Create and manage invoices
- **Estimates**: Create and manage estimates
- **Expenses**: Track expenses and expense categories

### Reporting
- **Reports**: Access various reports
  - Time reports
  - Expense reports
  - Uninvoiced reports
  - Project budget reports

### Administration
- **Roles**: Manage user roles

## Examples

### Time Entry Management

```go
// Create time entry via duration
entry, err := client.TimeEntries.CreateViaDuration(ctx, &harvest.TimeEntryCreateViaDurationRequest{
    ProjectID:  12345,
    TaskID:     67890,
    SpentDate:  "2025-01-10",
    Hours:      2.5,
    Notes:      "Working on feature X",
})

// Start a timer
timer, err := client.TimeEntries.CreateViaStartEndTime(ctx, &harvest.TimeEntryCreateViaStartEndTimeRequest{
    ProjectID: 12345,
    TaskID:    67890,
    SpentDate: "2025-01-10",
})

// Stop the timer
stopped, err := client.TimeEntries.Stop(ctx, timer.ID)
```

### Project Management

```go
// Create a project
project, err := client.Projects.Create(ctx, &harvest.ProjectCreateRequest{
    ClientID:   123,
    Name:       "New Website",
    IsBillable: &[]bool{true}[0],
    BillBy:     "project",
    Budget:     50000,
})

// Assign user to project
assignment, err := client.Projects.CreateUserAssignment(ctx, project.ID, 
    &harvest.UserAssignmentCreateRequest{
        UserID:           456,
        IsProjectManager: &[]bool{true}[0],
        HourlyRate:       150.00,
    })

// Update project budget
updated, err := client.Projects.Update(ctx, project.ID, &harvest.ProjectUpdateRequest{
    Budget: 75000,
})
```

### Pagination

```go
// Manual pagination
opts := &harvest.ProjectListOptions{
    ListOptions: harvest.ListOptions{
        Page:    1,
        PerPage: 50,
    },
}

for {
    projects, err := client.Projects.List(ctx, opts)
    if err != nil {
        return err
    }
    
    for _, project := range projects.Projects {
        // Process project
    }
    
    if !projects.HasNextPage() {
        break
    }
    opts.Page = *projects.NextPage
}
```

### Error Handling

```go
projects, err := client.Projects.List(ctx, nil)
if err != nil {
    switch e := err.(type) {
    case *harvest.RateLimitError:
        fmt.Printf("Rate limit exceeded. Resets at: %s\n", e.Rate.Reset)
        // Wait and retry
    case *harvest.ErrorResponse:
        fmt.Printf("API error: %s\n", e.Message)
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## API Coverage

This library provides complete coverage of the Harvest API v2. All endpoints documented in the [official API documentation](https://help.getharvest.com/api-v2/) are implemented.

### Implemented Endpoints

| Resource | Operations |
|----------|------------|
| **Company** | Get |
| **Clients** | List, Get, Create, Update, Delete |
| **Contacts** | List, Get, Create, Update, Delete |
| **Projects** | List, Get, Create, Update, Delete |
| **Project User Assignments** | List, Get, Create, Update, Delete |
| **Project Task Assignments** | List, Get, Create, Update, Delete |
| **Tasks** | List, Get, Create, Update, Delete |
| **Users** | List, Get, Create, Update, Delete, Get Current |
| **User Project Assignments** | List, Get Current |
| **User Billable Rates** | List, Get, Create, Update, Delete |
| **User Cost Rates** | List, Get, Create, Update, Delete |
| **Time Entries** | List, Get, Create, Update, Delete, Restart, Stop |
| **Invoices** | List, Get, Create, Update, Delete, Send |
| **Invoice Line Items** | Create, Update, Delete |
| **Invoice Messages** | List, Get, Create, Delete |
| **Invoice Payments** | List, Get, Create, Delete |
| **Estimates** | List, Get, Create, Update, Delete, Send |
| **Estimate Line Items** | Create, Update, Delete |
| **Estimate Messages** | List, Get, Create, Delete |
| **Expenses** | List, Get, Create, Update, Delete |
| **Expense Categories** | List, Get, Create, Update, Delete |
| **Roles** | List, Get, Create, Update, Delete |
| **Reports** | Time, Expenses, Uninvoiced, Project Budget |

## Development

### Building

```bash
# Download dependencies
go mod download

# Build
go build ./...

# Run tests
go test -v ./...

# Format code
gofmt -w .

# Lint
golangci-lint run ./...
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

When adding new endpoints:
1. Add types to `types.go`
2. Add service methods following the existing pattern
3. Use generics for standard CRUD operations
4. Ensure proper error handling
5. Add comprehensive tests

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Harvest](https://www.getharvest.com/) for providing a comprehensive API
- [shopspring/decimal](https://github.com/shopspring/decimal) for precise decimal handling
- [google/go-querystring](https://github.com/google/go-querystring) for query string encoding

## Related Projects

- [adlio/harvest](https://github.com/adlio/harvest) - Harvest API v1 client (deprecated API, unmaintained)
- [Harvest API Documentation](https://help.getharvest.com/api-v2/) - Official API documentation
- [Harvest Postman Collection](https://www.postman.com/harvest-api) - Official Postman collection

## Support

For issues, questions, or contributions, please use the [GitHub issue tracker](https://github.com/joefitzgerald/harvest/issues).