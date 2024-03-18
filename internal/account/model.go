package account

import (
	"transaction-server/internal/common/db"
	"transaction-server/internal/dto"
)

// Account represents the account entity.
type Account struct {
	db.Model              // Embedding the common database model
	Name           string `json:"name,omitempty" audit:"name"`              // Name of the account
	DocumentNumber string `json:"document_number,omitempty" audit:"doc_no"` // Document number associated with the account
}

// TableName returns the name of the database table for the Account entity.
func (e *Account) TableName() string {
	return "accounts"
}

// EntityName returns the name of the entity.
func (e *Account) EntityName() string {
	return "account"
}

// SetDefaults sets default values for the Account entity.
func (e *Account) SetDefaults() error {
	return nil
}

// ToDto converts the Account entity to its DTO (data transfer object) representation.
func (e *Account) ToDto() *dto.Account {
	return &dto.Account{
		ID:             e.ID,
		Name:           e.Name,
		DocumentNumber: e.DocumentNumber,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

// ApplyDto updates the Account entity fields based on the values provided in the DTO.
func (e *Account) ApplyDto(val *dto.Account) {
	e.Name = val.Name
	e.DocumentNumber = val.DocumentNumber
}
