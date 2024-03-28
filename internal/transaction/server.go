package transaction

import (
	"github.com/gin-gonic/gin"
	"transaction-server/internal/common"
	"transaction-server/internal/dto"
	"transaction-server/internal/validator"
)

type IServer interface {
	Create(ctx *gin.Context, req *dto.CreateTransactionRequest) *dto.CreateTransactionResponse
	Get(ctx *gin.Context, id string) *dto.GetTransactionResponse
	List(ctx *gin.Context, req *dto.ListTransactionRequest) *dto.ListTransactionResponse
}

type Server struct {
	core ICore
}

func NewServer(core ICore) IServer {
	return &Server{core: core}
}

func (s *Server) Create(ctx *gin.Context, req *dto.CreateTransactionRequest) *dto.CreateTransactionResponse {
	if err := validator.NewValidTransaction(req, validator.CreateTransactionValidator); err != nil {
		return &dto.CreateTransactionResponse{Base: dto.GetErrorResponse(common.ErrValidationFailed, err.Error())}
	}
	transaction := new(Transaction)
	transaction.ApplyDto(req.Transaction)
	if err := s.core.Create(ctx, transaction); err != nil {
		return &dto.CreateTransactionResponse{Base: dto.GetErrorResponse(common.ErrDBPersistError, err.Error())}
	}
	return &dto.CreateTransactionResponse{Transaction: transaction.ToDto(), Base: &dto.Base{Success: true}}
}

func (s *Server) Get(ctx *gin.Context, id string) *dto.GetTransactionResponse {
	if err := validator.NewValidTransaction(id, validator.GetTransactionValidator); err != nil {
		return &dto.GetTransactionResponse{Base: dto.GetErrorResponse(common.ErrValidationFailed, err.Error())}
	}
	transaction := new(Transaction)
	if err := s.core.Get(ctx, transaction, id); err != nil {
		return &dto.GetTransactionResponse{Base: dto.GetErrorResponse(common.ErrDBQueryError, err.Error())}
	}
	return &dto.GetTransactionResponse{Transaction: transaction.ToDto(), Base: &dto.Base{Success: true}}
}

func (s *Server) List(ctx *gin.Context, req *dto.ListTransactionRequest) *dto.ListTransactionResponse {
	if err := validator.NewValidTransaction(req, validator.ListTransactionValidator); err != nil {
		return &dto.ListTransactionResponse{Base: dto.GetErrorResponse(common.ErrValidationFailed, err.Error())}
	}
	transactions, err := s.core.List(ctx, req)
	if err != nil {
		return &dto.ListTransactionResponse{Base: dto.GetErrorResponse(common.ErrDBQueryError, err.Error())}
	}
	transactionsDto := make([]*dto.Transaction, 0)
	for _, transaction := range *transactions {
		transactionsDto = append(transactionsDto, transaction.ToDto())
	}
	return &dto.ListTransactionResponse{Transactions: transactionsDto, Base: &dto.Base{Success: true}}
}
