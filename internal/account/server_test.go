package account_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"transaction-server/internal/account/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"transaction-server/internal/account"
	"transaction-server/internal/common"
	"transaction-server/internal/dto"
)

type ServerTest struct {
	mockCoreCtrl *gomock.Controller
	mockCore     *mock.MockICore
	server       account.IServer
}

func setupServerTest(t *testing.T) *ServerTest {
	mockCoreCtrl := gomock.NewController(t)
	mockCore := mock.NewMockICore(mockCoreCtrl)
	server := account.NewServer(mockCore)
	return &ServerTest{
		mockCoreCtrl: mockCoreCtrl,
		mockCore:     mockCore,
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
	req := &dto.CreateAccountRequest{
		Account: &dto.Account{
			Name:           "John Doe",
			DocumentNumber: "123456789",
		},
	}

	td.mockCore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	resp := td.server.Create(ctx, req)
	assert.True(t, resp.Success)
}

func TestServer_Create_ValidationFailed(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.CreateAccountRequest{
		Account: &dto.Account{},
	}

	resp := td.server.Create(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, resp.Error.Code, common.ErrValidationFailed)
}

func TestServer_Create_DBPersistError(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	req := &dto.CreateAccountRequest{
		Account: &dto.Account{
			Name:           "John Doe",
			DocumentNumber: "123456789",
		},
	}

	td.mockCore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("DB error"))

	resp := td.server.Create(ctx, req)
	assert.False(t, resp.Success)
	assert.Equal(t, resp.Error.Code, common.ErrDBPersistError)
}

func TestServer_Get_Success(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "0b0e0000000000"

	td.mockCore.EXPECT().Get(gomock.Any(), gomock.Any(), id).Return(nil)

	resp := td.server.Get(ctx, id)
	assert.True(t, resp.Success)
}

func TestServer_Get_ValidationFailed(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "invalid-id"

	resp := td.server.Get(ctx, id)
	assert.False(t, resp.Success)
	assert.Equal(t, resp.Error.Code, common.ErrValidationFailed)
}

func TestServer_Get_DBQueryError(t *testing.T) {
	td := setupServerTest(t)
	defer teardownServerTest(td)

	ctx := &gin.Context{}
	id := "0b0e0000000000"

	td.mockCore.EXPECT().Get(gomock.Any(), gomock.Any(), id).Return(errors.New("DB error"))

	resp := td.server.Get(ctx, id)
	assert.False(t, resp.Success)
	assert.Equal(t, resp.Error.Code, common.ErrDBQueryError)
}
