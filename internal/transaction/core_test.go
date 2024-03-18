package transaction_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"transaction-server/internal/dto"
	"transaction-server/internal/transaction/mock"

	"github.com/stretchr/testify/assert"
	"transaction-server/internal/transaction"
)

type testDependencies struct {
	mockRepoCtrl *gomock.Controller
	mockRepo     *mock.MockIRepo
	core         transaction.ICore
}

func setupTest(t *testing.T) *testDependencies {
	mockRepoCtrl := gomock.NewController(t)
	mockRepo := mock.NewMockIRepo(mockRepoCtrl)
	core := transaction.NewCore(mockRepo)
	return &testDependencies{
		mockRepoCtrl: mockRepoCtrl,
		mockRepo:     mockRepo,
		core:         core,
	}
}

func teardownTest(td *testDependencies) {
	td.mockRepoCtrl.Finish()
}

func TestCore_Create_Success(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{}

	td.mockRepo.EXPECT().Create(ctx, model).Return(nil)

	err := td.core.Create(ctx, model)
	assert.NoError(t, err)
}

func TestCore_Get_Success(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{}
	id := "some_id"

	td.mockRepo.EXPECT().FindByID(ctx, model, id).Return(nil)

	err := td.core.Get(ctx, model, id)
	assert.NoError(t, err)
}

func TestCore_List_Success(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	request := &dto.ListTransactionRequest{}

	td.mockRepo.EXPECT().FindManyWithFilters(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	transactions, err := td.core.List(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
}

func TestCore_List_RepoError(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	request := &dto.ListTransactionRequest{}

	td.mockRepo.EXPECT().FindManyWithFilters(ctx, gomock.Any(), gomock.Any()).Return(errors.New("repo error"))

	transactions, err := td.core.List(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, transactions)
}
