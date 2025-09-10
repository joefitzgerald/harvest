# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Testing

```bash
# Download dependencies
go mod download

# Build the library
go build ./...

# Run static analysis
go vet ./...

# Format code (check which files need formatting)
gofmt -l .

# Format code (apply formatting)
gofmt -w .

# Run tests (when tests are added)
go test -v ./...
go test -v -run TestName ./...  # Run specific test
```

### Linting and Quality Checks

```bash
# Install golangci-lint if needed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run comprehensive linting
golangci-lint run ./...
```

## Code Architecture

This is a Go client library for the Harvest API v2 that uses modern Go patterns including generics (Go 1.25+) to reduce code duplication. The architecture follows a service-oriented pattern where each Harvest resource has its own service struct.

### Core Components

1. **Client (`client.go`)**: Main API client that manages authentication, HTTP requests, and service initialization. Uses generic CRUD methods (`List`, `Get`, `Create`, `Update`, `Delete`) to reduce boilerplate.

2. **Services**: Each Harvest resource (Projects, TimeEntries, Users, etc.) has a dedicated service file that implements resource-specific methods while leveraging the generic CRUD operations.

3. **Types (`types.go`)**: Contains all struct definitions for API requests and responses. Uses `shopspring/decimal` for precise financial calculations.

4. **Pagination (`pagination.go`)**: Generic pagination system with iterator pattern for efficient data retrieval.

5. **Error Handling (`errors.go`)**: Custom error types including rate limit errors with automatic parsing of rate limit headers.

### Key Design Patterns

- **Generic CRUD Operations**: The library uses Go 1.25 generics to implement reusable CRUD methods that work with any resource type
- **Service Pattern**: Each API resource has its own service struct attached to the main client
- **Builder Pattern for Options**: List methods use option structs for filtering and pagination
- **Context Support**: All API methods accept `context.Context` for cancellation and timeout control

### Authentication

The library requires two environment variables:

- `HARVEST_ACCESS_TOKEN`: Personal access token from Harvest
- `HARVEST_ACCOUNT_ID`: Harvest account ID

User-Agent header is mandatory per Harvest API requirements (format: "AppName (contact@example.com)").

### API Coverage

The library provides complete coverage of Harvest API v2 endpoints:

- Company, Clients, Contacts
- Projects (with User/Task Assignments)
- Users (with Project/Expense Category Assignments, Billable Rates, Cost Rates)
- Tasks, Time Entries
- Invoices (with Line Items, Messages, Payments)
- Estimates (with Line Items, Messages)
- Expenses (with Categories)
- Reports (Time, Expense, Uninvoiced, Project Budget)
- Roles

### Postman Collection

The `docs/harvest-api-v2.postman_collection.json` file contains the complete API specification and can be used for:

- API endpoint reference
- Request/response examples
- Testing API calls directly

## Important Conventions

- Use `shopspring/decimal.Decimal` for all monetary values and rates
- Date-only fields use custom `Date` type (format: "YYYY-MM-DD")
- All list endpoints return paginated results with metadata
- Rate limiting is handled automatically with proper error types
- JSON tags use snake_case matching Harvest API conventions
- Omit empty fields in requests using `omitempty` tags
- Boolean pointers are used for optional boolean fields (using pattern: `&[]bool{true}[0]`)

## Adding New Features

When implementing new Harvest API endpoints:

1. Add request/response types to `types.go`
2. Create or update the appropriate service file
3. Use generic CRUD methods where applicable
4. Ensure proper error handling including rate limit errors
5. Add appropriate JSON tags with `omitempty` for optional fields
6. Use `decimal.Decimal` for financial values
7. Support context cancellation
