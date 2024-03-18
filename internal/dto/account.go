// Package dto provides request and response objects.
package dto

// swagger:model
type CreateAccountRequest struct {
	// The account information.
	Account *Account `json:"account"`
}

// swagger:model
type CreateAccountResponse struct {
	// Base response object.
	*Base
	// The created account information.
	Account *Account `json:"account,omitempty"`
}

// swagger:model
type GetAccountResponse struct {
	// Base response object.
	*Base
	// The retrieved account information.
	Account *Account `json:"account,omitempty"`
}

// Account represents account details.
type Account struct {
	// The ID of the account.
	ID string `json:"id"`
	// The name of the account.
	Name string `json:"name"`
	// The document number associated with the account.
	DocumentNumber string `json:"document_number"`
	// The timestamp when the account was created.
	CreatedAt int64 `json:"created_at"`
	// The timestamp when the account was last updated.
	UpdatedAt int64 `json:"updated_at"`
}
