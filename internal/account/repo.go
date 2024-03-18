package account

import (
	"context"
	"transaction-server/internal/common/db"
)

type IRepo interface {
	FindByID(ctx context.Context, receiver db.IModel, id string) error
	Create(ctx context.Context, receiver db.IModel) error
}
