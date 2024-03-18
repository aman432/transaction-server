package account

import (
	"transaction-server/internal/common/db"
	"transaction-server/internal/dto"
)

type Account struct {
	db.Model
	Name           string `json:"name,omitempty" audit:"name"`
	DocumentNumber string `json:"document_number,omitempty" audit:"doc_no"`
}

func (e *Account) TableName() string {
	return "accounts"
}

func (e *Account) EntityName() string {
	return "account"
}

func (e *Account) SetDefaults() error {
	return nil
}

func (e *Account) ToDto() *dto.Account {
	return &dto.Account{
		ID:             e.ID,
		Name:           e.Name,
		DocumentNumber: e.DocumentNumber,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func (e *Account) ApplyDto(val *dto.Account) {
	e.Name = val.Name
	e.DocumentNumber = val.DocumentNumber
}
