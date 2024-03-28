package transaction

import (
	"context"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
	"math"
	"transaction-server/internal/common/db"
	"transaction-server/internal/dto"
)

type ICore interface {
	Create(ctx context.Context, model *Transaction) error
	Get(ctx context.Context, model *Transaction, id string) error
	List(ctx context.Context, request IListRequest) (*[]Transaction, error)
}

type Core struct {
	repo IRepo
}

func (c Core) Create(ctx context.Context, model *Transaction) error {
	return c.repo.Transaction(ctx, func(ctx context.Context) error {
		if !utils.Contains(dto.NegativeOperationTypesString, model.OperationType.String()) {
			if err2 := c.checkAndUpdateExistingBalances(ctx, model); err2 != nil {
				return err2
			}
		}
		return c.repo.Create(ctx, model)
	})
}

func (c Core) checkAndUpdateExistingBalances(ctx context.Context, model *Transaction) error {
	listRequest := &dto.ListTransactionRequest{
		AccountId:      model.AccountId,
		OperationTypes: dto.NegativeOperationTypesString,
		Limit:          100,
	}
	transactionResponseList, err := c.List(ctx, listRequest)
	if err != nil {
		return err
	}
	if len(*transactionResponseList) == 0 {
		model.Balance = model.Amount
	} else {
		remainingBal, err := c.updateExistingBalances(ctx, model, transactionResponseList, err)
		if err != nil {
			return err
		}
		model.Balance = remainingBal
	}
	return nil
}

func (c Core) updateExistingBalances(ctx context.Context, model *Transaction, transactionResponseList *[]Transaction, err error) (float64, error) {
	remainingBal := model.Balance
	for _, transaction := range *transactionResponseList {
		if remainingBal == 0 {
			break
		}
		if math.Abs(transaction.Balance) <= remainingBal {
			remainingBal = remainingBal - math.Abs(transaction.Balance)
			transaction.Balance = 0
		} else {
			transaction.Balance = remainingBal - math.Abs(transaction.Balance)
			remainingBal = 0
		}
		if err = c.repo.Update(ctx, &transaction, "balance"); err != nil {
			return 0, err
		}
	}
	return remainingBal, nil
}

func (c Core) Get(ctx context.Context, model *Transaction, id string) error {
	return c.repo.FindByID(ctx, model, id)
}

func (c Core) List(ctx context.Context, request IListRequest) (*[]Transaction, error) {
	conditions := make([]clause.Expression, 0)
	if request.GetAccountId() != "" {
		conditions = append(conditions, clause.Eq{Column: "account_id", Value: request.GetAccountId()})
	}
	if request.GetOperationType() != "" {
		conditions = append(conditions, clause.Eq{Column: "operation_type", Value: OperationFromString(request.GetOperationType())})
	}
	if len(request.GetOperationTypes()) > 0 {
		conditions = append(conditions, clause.IN{Column: "operation_type", Values: OperationFromStrings(request.GetOperationTypes())})
	}
	repoRequest := &db.FindManyWithConditionsRequest{
		FindManyRequest: db.FindManyRequest{
			Limit:  request.GetLimit(),
			Offset: request.GetOffset(),
		},
		Conditions: conditions,
	}
	listResponse := make([]Transaction, 0)
	if err := c.repo.FindManyWithFilters(ctx, &listResponse, repoRequest); err != nil {
		return nil, err
	}
	return &listResponse, nil
}

func NewCore(repo IRepo) ICore {
	return &Core{repo: repo}
}
