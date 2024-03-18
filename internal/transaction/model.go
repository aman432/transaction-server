package transaction

import (
	"time"
	"transaction-server/internal/common/db"
	"transaction-server/internal/dto"
)

type Transaction struct {
	db.Model
	AccountId     string            `json:"account_id"`
	OperationType dto.OperationType `json:"operation_type"`
	Amount        float64           `json:"amount"`
	EventDate     int64             `json:"event_date"`
}

func (e *Transaction) TableName() string {
	return "transactions"
}

func (e *Transaction) EntityName() string {
	return "Transaction"
}

func (e *Transaction) SetDefaults() error {
	return nil
}

func (e *Transaction) ToDto() *dto.Transaction {
	return &dto.Transaction{
		ID:            e.ID,
		AccountID:     e.AccountId,
		OperationType: e.OperationType.String(),
		Amount:        e.Amount,
		EventDate:     time.Unix(e.EventDate, 0).String(),
	}
}

func (e *Transaction) ApplyDto(val *dto.Transaction) {
	e.AccountId = val.AccountID
	e.OperationType = OperationFromString(val.OperationType)
	e.Amount = setAmountSign(e.OperationType, val.Amount)
	e.EventDate = time.Now().Unix()
}

func OperationFromString(o string) dto.OperationType {
	for i, operationType := range dto.OperationTypes {
		if operationType == o {
			return dto.OperationType(i)
		}
	}
	return 0
}

func setAmountSign(opType dto.OperationType, amount float64) float64 {
	if opType.IsPresent(dto.NegativeOperationTypes[:]) {
		return -amount
	}
	return amount
}

type IListRequest interface {
	GetLimit() uint32
	GetOffset() uint32
	GetAccountId() string
	GetOperationType() string
}
