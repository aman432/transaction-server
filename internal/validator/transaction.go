package validator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"transaction-server/internal/common/db/datatype"
	"transaction-server/internal/dto"
)

type (
	TransactionValidator string
)

const (
	CreateTransactionValidator = "Create"
	GetTransactionValidator    = "Get"
	ListTransactionValidator   = "List"
)

// NewValidTransaction validates Transaction APIs and return error or nil
func NewValidTransaction(ev interface{}, validator TransactionValidator) error {
	var ve validation.Validatable
	switch validator {
	case CreateTransactionValidator:
		ve = &ValidCreateTransaction{ev.(*dto.CreateTransactionRequest)}
	case GetTransactionValidator:
		ve = &ValidGetTransaction{ev.(string)}
	case ListTransactionValidator:
		ve = &ValidListTransaction{ev.(*dto.ListTransactionRequest)}
	}
	err := ve.Validate()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// ValidCreateTransaction wraps Create Plan struct
type ValidCreateTransaction struct {
	*dto.CreateTransactionRequest
}

func (v *ValidCreateTransaction) Validate() error {
	return validation.ValidateStruct(
		v.Transaction,
		validation.Field(
			&v.Transaction.OperationType,
			validation.In(
				"Normal_Purchase",
				"Purchase_With_Installment",
				"Withdraw",
				"Credit_Voucher",
			),
			validation.Required,
		),
		validation.Field(
			&v.Transaction.AccountID,
			validation.Required,
			validation.By(datatype.IsUUID),
		),
		validation.Field(
			&v.Transaction.Amount,
			validation.Required,
		),
	)
}

// ValidGetTransaction wraps Get Plan struct
type ValidGetTransaction struct {
	id string
}

func (v *ValidGetTransaction) Validate() error {
	return validation.ValidateStruct(
		v,
		validation.Field(
			&v.id,
			validation.Required,
			validation.By(datatype.IsUUID),
		),
	)
}

// ValidListTransaction wraps List Transaction struct
type ValidListTransaction struct {
	*dto.ListTransactionRequest
}

func (v *ValidListTransaction) Validate() error {
	return validation.ValidateStruct(
		v,
		validation.Field(
			&v.Limit,
			validation.Max(uint32(20)),
		),
	)
}
