package account

import (
	"github.com/gin-gonic/gin"
	"transaction-server/internal/common"
	"transaction-server/internal/dto"
	"transaction-server/internal/validator"
)

type IServer interface {
	Create(ctx *gin.Context, req *dto.CreateAccountRequest) *dto.CreateAccountResponse
	Get(ctx *gin.Context, id string) *dto.GetAccountResponse
}

type Server struct {
	core ICore
}

func NewServer(core ICore) IServer {
	return &Server{core: core}
}

func (s *Server) Create(ctx *gin.Context, req *dto.CreateAccountRequest) *dto.CreateAccountResponse {
	if err := validator.NewValidAccount(req, validator.CreateAccountValidator); err != nil {
		return &dto.CreateAccountResponse{Base: dto.GetErrorResponse(common.ErrValidationFailed, err.Error())}
	}
	account := new(Account)
	account.ApplyDto(req.Account)
	if err := s.core.Create(ctx, account); err != nil {
		return &dto.CreateAccountResponse{Base: dto.GetErrorResponse(common.ErrDBPersistError, err.Error())}
	}
	return &dto.CreateAccountResponse{Account: account.ToDto(), Base: &dto.Base{Success: true}}
}

func (s *Server) Get(ctx *gin.Context, id string) *dto.GetAccountResponse {
	if err := validator.NewValidAccount(id, validator.GetAccountValidator); err != nil {
		return &dto.GetAccountResponse{Base: dto.GetErrorResponse(common.ErrValidationFailed, err.Error())}
	}
	account := new(Account)
	if err := s.core.Get(ctx, account, id); err != nil {
		return &dto.GetAccountResponse{Base: dto.GetErrorResponse(common.ErrDBQueryError, err.Error())}
	}
	return &dto.GetAccountResponse{Account: account.ToDto(), Base: &dto.Base{Success: true}}
}
