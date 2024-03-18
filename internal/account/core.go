package account

import (
	"context"
)

type ICore interface {
	Create(ctx context.Context, account *Account) error
	Get(ctx context.Context, account *Account, id string) error
}

type Core struct {
	repo IRepo
}

func NewCore(repo IRepo) *Core {
	return &Core{repo: repo}
}

func (c *Core) Create(ctx context.Context, account *Account) error {
	return c.repo.Create(ctx, account)
}

func (c *Core) Get(ctx context.Context, account *Account, id string) error {
	return c.repo.FindByID(ctx, account, id)
}
