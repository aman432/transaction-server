package transaction

import (
	"context"
	db2 "transaction-server/internal/common/db"
)

type IRepo interface {
	FindByID(ctx context.Context, receiver db2.IModel, id string) error
	Create(ctx context.Context, receiver db2.IModel) error
	Update(ctx context.Context, receiver db2.IModel, selectiveList ...string) error
	FindManyWithFilters(ctx context.Context, models interface{}, req db2.FindManyWithFiltersRequester) error
	Transaction(ctx context.Context, fc func(ctx context.Context) error) error
}
