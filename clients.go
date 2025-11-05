package harvest

import (
	"context"
	"fmt"
)

// ClientsService handles communication with the client related
// methods of the Harvest API.
type ClientsService struct {
	client *API
}

// ClientListOptions specifies optional parameters to the List method.
type ClientListOptions struct {
	ListOptions
	IsActive     *bool  `url:"is_active,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// ClientList represents a list of clients.
type ClientList struct {
	Clients []Client `json:"clients"`
	Paginated[Client]
}

// ListPage returns a single page of clients.
func (s *ClientsService) ListPage(ctx context.Context, opts *ClientListOptions) (*ClientList, error) {
	u, err := addOptions("clients", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var clients ClientList
	_, err = s.client.Do(ctx, req, &clients)
	if err != nil {
		return nil, err
	}

	// Copy clients to Items for pagination
	clients.Items = clients.Clients

	return &clients, nil
}

// List returns all clients across all pages.
func (s *ClientsService) List(ctx context.Context, opts *ClientListOptions) ([]Client, error) {
	if opts == nil {
		opts = &ClientListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allClients []Client

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allClients = append(allClients, result.Clients...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allClients, nil
}

// Get retrieves a specific client.
func (s *ClientsService) Get(ctx context.Context, clientID int64) (*Client, error) {
	return Get[Client](ctx, s.client, fmt.Sprintf("clients/%d", clientID))
}

// ClientCreateRequest represents a request to create a client.
type ClientCreateRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active,omitempty"`
	Address  string `json:"address,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// Create creates a new client.
func (s *ClientsService) Create(ctx context.Context, client *ClientCreateRequest) (*Client, error) {
	return Create[Client](ctx, s.client, "clients", client)
}

// ClientUpdateRequest represents a request to update a client.
type ClientUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Address  string `json:"address,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// Update updates a client.
func (s *ClientsService) Update(ctx context.Context, clientID int64, client *ClientUpdateRequest) (*Client, error) {
	return Update[Client](ctx, s.client, fmt.Sprintf("clients/%d", clientID), client)
}

// Delete deletes a client.
func (s *ClientsService) Delete(ctx context.Context, clientID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("clients/%d", clientID))
}

// ContactsService handles communication with the contact related
// methods of the Harvest API.
type ContactsService struct {
	client *API
}

// ContactListOptions specifies optional parameters to the List method.
type ContactListOptions struct {
	ListOptions
	ClientID     int64  `url:"client_id,omitempty"`
	UpdatedSince string `url:"updated_since,omitempty"`
}

// ContactList represents a list of contacts.
type ContactList struct {
	Contacts []Contact `json:"contacts"`
	Paginated[Contact]
}

// ListPage returns a single page of contacts.
func (s *ContactsService) ListPage(ctx context.Context, opts *ContactListOptions) (*ContactList, error) {
	u, err := addOptions("contacts", opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var contacts ContactList
	_, err = s.client.Do(ctx, req, &contacts)
	if err != nil {
		return nil, err
	}

	// Copy contacts to Items for pagination
	contacts.Items = contacts.Contacts

	return &contacts, nil
}

// List returns all contacts across all pages.
func (s *ContactsService) List(ctx context.Context, opts *ContactListOptions) ([]Contact, error) {
	if opts == nil {
		opts = &ContactListOptions{}
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultPerPage
	}

	var allContacts []Contact

	for {
		result, err := s.ListPage(ctx, opts)
		if err != nil {
			return nil, err
		}

		allContacts = append(allContacts, result.Contacts...)

		if !result.HasNextPage() {
			break
		}

		opts.Page = *result.NextPage
	}

	return allContacts, nil
}

// GetContact retrieves a specific contact.
func (s *ContactsService) Get(ctx context.Context, contactID int64) (*Contact, error) {
	return Get[Contact](ctx, s.client, fmt.Sprintf("contacts/%d", contactID))
}

// ContactCreateRequest represents a request to create a contact.
type ContactCreateRequest struct {
	ClientID    int64  `json:"client_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Title       string `json:"title,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneOffice string `json:"phone_office,omitempty"`
	PhoneMobile string `json:"phone_mobile,omitempty"`
	Fax         string `json:"fax,omitempty"`
}

// CreateContact creates a new contact.
func (s *ContactsService) Create(ctx context.Context, contact *ContactCreateRequest) (*Contact, error) {
	return Create[Contact](ctx, s.client, "contacts", contact)
}

// ContactUpdateRequest represents a request to update a contact.
type ContactUpdateRequest struct {
	ClientID    int64  `json:"client_id,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Title       string `json:"title,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneOffice string `json:"phone_office,omitempty"`
	PhoneMobile string `json:"phone_mobile,omitempty"`
	Fax         string `json:"fax,omitempty"`
}

// UpdateContact updates a contact.
func (s *ContactsService) Update(ctx context.Context, contactID int64, contact *ContactUpdateRequest) (*Contact, error) {
	return Update[Contact](ctx, s.client, fmt.Sprintf("contacts/%d", contactID), contact)
}

// DeleteContact deletes a contact.
func (s *ContactsService) Delete(ctx context.Context, contactID int64) error {
	return Delete(ctx, s.client, fmt.Sprintf("contacts/%d", contactID))
}
