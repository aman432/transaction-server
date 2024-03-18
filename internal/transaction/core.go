package transaction

import (
	"context"
	"gorm.io/gorm/clause"
	"transaction-server/internal/common/db"
)

type ICore interface {
	Create(ctx context.Context, model *Transaction) error
	Get(ctx context.Context, model *Transaction, id string) error
	List(ctx context.Context, request IListRequest) ([]*Transaction, error)
}

type Core struct {
	repo IRepo
}

func (c Core) Create(ctx context.Context, model *Transaction) error {
	return c.repo.Create(ctx, model)
}

func (c Core) Get(ctx context.Context, model *Transaction, id string) error {
	return c.repo.FindByID(ctx, model, id)
}

func (c Core) List(ctx context.Context, request IListRequest) ([]*Transaction, error) {
	conditions := make([]clause.Expression, 0)
	if request.GetAccountId() != "" {
		conditions = append(conditions, clause.Eq{Column: "account_id", Value: request.GetAccountId()})
	}
	if request.GetOperationType() != "" {
		conditions = append(conditions, clause.Eq{Column: "operation_type", Value: OperationFromString(request.GetOperationType())})
	}
	repoRequest := &db.FindManyWithConditionsRequest{
		FindManyRequest: db.FindManyRequest{
			Limit:  request.GetLimit(),
			Offset: request.GetOffset(),
		},
		Conditions: conditions,
	}
	listResponse := make([]*Transaction, 0)
	if err := c.repo.FindManyWithFilters(ctx, &listResponse, repoRequest); err != nil {
		return nil, err
	}
	return listResponse, nil
}

func NewCore(repo IRepo) ICore {
	return &Core{repo: repo}
}
