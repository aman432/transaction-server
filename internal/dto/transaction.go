package dto

// OperationType represents transaction operation types.
// swagger:model
type OperationType int

const (
	OperationTypeNormalPurchase          OperationType = 1
	OperationTypePurchaseWithInstallment OperationType = 2
	OperationTypeWithdraw                OperationType = 3
	OperationTypeCreditVoucher           OperationType = 4
)

var OperationTypes = [...]string{
	"Unknown",
	"Normal_Purchase",
	"Purchase_With_Installment",
	"Withdraw",
	"Credit_Voucher",
}

var NegativeOperationTypes = [...]OperationType{
	OperationTypeNormalPurchase,
	OperationTypeWithdraw,
}

func (o OperationType) IsPresent(arr []OperationType) bool {
	for _, a := range arr {
		if a == o {
			return true
		}
	}
	return false
}

func (o OperationType) String() string {
	return OperationTypes[o]
}

// Transaction represents a transaction object.
// swagger:model
type Transaction struct {
	// The ID of the transaction.
	ID string `json:"id"`
	// The ID of the account associated with the transaction.
	AccountID string `json:"account_id"`
	// The type of operation.
	OperationType string `json:"operation_type"`
	// The amount of the transaction.
	Amount float64 `json:"amount"`
	// The event date of the transaction.
	EventDate string `json:"event_date"`
}

// CreateTransactionRequest represents the request object for creating a transaction.
// swagger:model
type CreateTransactionRequest struct {
	// The transaction to be created.
	Transaction *Transaction `json:"transaction"`
}

// CreateTransactionResponse represents the response object for creating a transaction.
// swagger:model
type CreateTransactionResponse struct {
	// The base response object.
	*Base
	// The created transaction.
	Transaction *Transaction `json:"transaction,omitempty"`
}

// GetTransactionResponse represents the response object for retrieving a transaction.
// swagger:model
type GetTransactionResponse struct {
	// The base response object.
	*Base
	// The retrieved transaction.
	Transaction *Transaction `json:"transaction,omitempty"`
}

// ListTransactionRequest represents the request object for listing transactions.
// swagger:model
type ListTransactionRequest struct {
	// The limit for the number of transactions.
	Limit uint32 `json:"limit"`
	// The offset for pagination.
	Offset uint32 `json:"offset"`
	// The account ID to filter transactions.
	AccountId string `json:"account_id"`
	// The operation type to filter transactions.
	OperationType string `json:"operation_type"`
}

// GetLimit returns the limit value for pagination.
func (l *ListTransactionRequest) GetLimit() uint32 {
	return l.Limit
}

// GetOffset returns the offset value for pagination.
func (l *ListTransactionRequest) GetOffset() uint32 {
	return l.Offset
}

// GetAccountId returns the account ID.
func (l *ListTransactionRequest) GetAccountId() string {
	return l.AccountId
}

// GetOperationType returns the operation type.
func (l *ListTransactionRequest) GetOperationType() string {
	return l.OperationType
}

// ListTransactionResponse represents the response object for listing transactions.
// swagger:model
type ListTransactionResponse struct {
	// The base response object.
	*Base
	// The list of transactions.
	Transactions []*Transaction `json:"transactions,omitempty"`
}
