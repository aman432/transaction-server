package validator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"transaction-server/internal/common/db/datatype"
	"transaction-server/internal/dto"
)

type (
	AccountValidator string
)

const (
	CreateAccountValidator = "Create"
	GetAccountValidator    = "Get"
)

// NewValidAccount validates Account APIs and return error or nil
func NewValidAccount(ev interface{}, validator AccountValidator) error {
	var ve validation.Validatable
	switch validator {
	case CreateAccountValidator:
		ve = &ValidCreateAccount{ev.(*dto.CreateAccountRequest)}
	case GetAccountValidator:
		ve = &ValidGetAccount{ev.(string)}
	}
	err := ve.Validate()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// ValidCreateAccount wraps Create Plan struct
type ValidCreateAccount struct {
	*dto.CreateAccountRequest
}

func (v *ValidCreateAccount) Validate() error {
	return validation.ValidateStruct(
		v.Account,
		validation.Field(
			&v.Account.Name,
			validation.Required,
		),
		validation.Field(
			&v.Account.DocumentNumber,
			validation.Required,
		),
	)
}

// ValidGetAccount wraps Get Plan struct
type ValidGetAccount struct {
	id string
}

func (v *ValidGetAccount) Validate() error {
	return validation.ValidateStruct(
		v,
		validation.Field(
			&v.id,
			validation.Required,
			validation.By(datatype.IsUUID),
		),
	)
}
