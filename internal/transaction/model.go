package transaction

import (
	"time"
	"transaction-server/internal/common/db"
	"transaction-server/internal/dto"
)

// Transaction represents a financial transaction entity.
type Transaction struct {
	db.Model                        // Embedding the common database model
	AccountId     string            `json:"account_id"`     // ID of the account associated with the transaction
	OperationType dto.OperationType `json:"operation_type"` // Type of operation (e.g., purchase, withdrawal)
	Amount        float64           `json:"amount"`         // Amount of the transaction
	EventDate     int64             `json:"event_date"`     // Date and time when the transaction occurred (Unix timestamp)
}

// TableName returns the name of the database table for the Transaction entity.
func (e *Transaction) TableName() string {
	return "transactions"
}

// EntityName returns the name of the entity.
func (e *Transaction) EntityName() string {
	return "Transaction"
}

// SetDefaults sets default values for the Transaction entity.
func (e *Transaction) SetDefaults() error {
	return nil
}

// ToDto converts the Transaction entity to its DTO (data transfer object) representation.
func (e *Transaction) ToDto() *dto.Transaction {
	return &dto.Transaction{
		ID:            e.ID,
		AccountID:     e.AccountId,
		OperationType: e.OperationType.String(),
		Amount:        e.Amount,
		EventDate:     time.Unix(e.EventDate, 0).String(),
	}
}

// ApplyDto updates the Transaction entity fields based on the values provided in the DTO.
func (e *Transaction) ApplyDto(val *dto.Transaction) {
	e.AccountId = val.AccountID
	e.OperationType = OperationFromString(val.OperationType)
	e.Amount = setAmountSign(e.OperationType, val.Amount)
	e.EventDate = time.Now().Unix()
}

// OperationFromString converts a string representation of an operation type to its corresponding OperationType.
func OperationFromString(o string) dto.OperationType {
	for i, operationType := range dto.OperationTypes {
		if operationType == o {
			return dto.OperationType(i)
		}
	}
	return 0
}

// setAmountSign sets the sign of the amount based on the operation type.
func setAmountSign(opType dto.OperationType, amount float64) float64 {
	if opType.IsPresent(dto.NegativeOperationTypes[:]) {
		return -amount
	}
	return amount
}

// IListRequest is the interface that wraps basic request attribute getters for list requests.
type IListRequest interface {
	GetLimit() uint32
	GetOffset() uint32
	GetAccountId() string
	GetOperationType() string
}
