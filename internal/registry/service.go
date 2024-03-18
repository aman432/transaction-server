package registry

import (
	"context"
	"transaction-server/app"
	"transaction-server/internal/account"
	"transaction-server/internal/common/db"
	"transaction-server/internal/transaction"
)

type IRegistry interface {
	GetAccountsServer() account.IServer
	GetTransactionsServer() transaction.IServer
}

type Registry struct {
	accountServer     account.IServer
	transactionServer transaction.IServer
}

func (r Registry) GetTransactionsServer() transaction.IServer {
	return r.transactionServer
}

func (r Registry) GetAccountsServer() account.IServer {
	return r.accountServer
}

func NewRegistry(ctx context.Context) IRegistry {
	commonRepo := db.NewRepo(app.Context().DB())
	accountCore := account.NewCore(commonRepo)
	accountServer := account.NewServer(accountCore)

	transactionCore := transaction.NewCore(commonRepo)
	transactionServer := transaction.NewServer(transactionCore)
	return &Registry{
		accountServer:     accountServer,
		transactionServer: transactionServer,
	}
}
