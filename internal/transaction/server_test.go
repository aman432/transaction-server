package transaction_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"transaction-server/internal/transaction/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"transaction-server/internal/common"
	"transaction-server/internal/dto"
	"transaction-server/internal/transaction"
)

type ServerTest struct {
	mockCoreCtrl *gomock.Controller
	core         *mock.MockICore
	server       transaction.IServer
}

func setupServerTest(t *testing.T) *ServerTest {
	mockCoreCtrl := gomock.NewController(t)
	mockCore := mock.NewMockICore(mockCoreCtrl)
	server := transaction.NewServer(mockCore)
	return &ServerTest{
		mockCoreCtrl: mockCoreCtrl,
		core:         mockCore,
		server:       server,
	}
}

func teardownServerTest(td *ServerTest) {
	td.mockCoreCtrl.Finish()
}

func TestServer_Create_Success(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.CreateTransactionRequest{
		Transaction: &dto.Transaction{
			OperationType: "Normal_Purchase",
			AccountID:     "0b0e0000000000",
			Amount:        100,
		},
	}

	td.core.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	resp := td.server.Create(ctx, req)
	assert.True(t, resp.Success)
}

func TestServer_Create_ValidationFailed(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.CreateTransactionRequest{
		Transaction: &dto.Transaction{},
	}

	resp := td.server.Create(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrValidationFailed, resp.Error.Code)
}

func TestServer_Create_DBPersistError(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.CreateTransactionRequest{
		Transaction: &dto.Transaction{
			OperationType: "Normal_Purchase",
			AccountID:     "0b0e0000000000",
			Amount:        100,
		},
	}

	td.core.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("DB error"))

	resp := td.server.Create(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrDBPersistError, resp.Error.Code)
}

func TestServer_Get_Success(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "0b0e0000000000"

	td.core.EXPECT().Get(ctx, gomock.Any(), id).Return(nil)

	resp := td.server.Get(ctx, id)
	assert.True(t, resp.Success)
}

func TestServer_Get_ValidationFailed(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "some_invalid_id"

	resp := td.server.Get(ctx, id)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrValidationFailed, resp.Error.Code)
}

func TestServer_Get_DBQueryError(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "0b0e0000000000"

	td.core.EXPECT().Get(ctx, gomock.Any(), id).Return(errors.New("DB error"))

	resp := td.server.Get(ctx, id)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrDBQueryError, resp.Error.Code)
}

func TestServer_List_Success(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.ListTransactionRequest{}

	transactions := make([]transaction.Transaction, 0)
	td.core.EXPECT().List(ctx, req).Return(&transactions, nil)

	resp := td.server.List(ctx, req)
	assert.True(t, resp.Success)
}

func TestServer_List_ValidationFailed(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.ListTransactionRequest{
		Limit: 40,
	}

	resp := td.server.List(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrValidationFailed, resp.Error.Code)
}

func TestServer_List_DBQueryError(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.ListTransactionRequest{}

	td.core.EXPECT().List(ctx, req).Return(nil, errors.New("DB error"))

	resp := td.server.List(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, common.ErrDBQueryError, resp.Error.Code)
}
