package transaction

import (
	"context"
	db2 "transaction-server/internal/common/db"
)

type IRepo interface {
	FindByID(ctx context.Context, receiver db2.IModel, id string) error
	Create(ctx context.Context, receiver db2.IModel) error
	FindManyWithFilters(ctx context.Context, models interface{}, req db2.FindManyWithFiltersRequester) error
}
