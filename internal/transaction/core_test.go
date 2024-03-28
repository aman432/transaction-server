package transaction_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	db2 "transaction-server/internal/common/db"
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

func TestCore_Create_Success_With_Negative_Balance(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{
		AccountId:     "some_id",
		OperationType: dto.OperationTypeWithdraw,
	}

	td.mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fc func(ctx context.Context) error) error {
			return fc(ctx)
		},
	)
	td.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	err := td.core.Create(ctx, model)
	assert.NoError(t, err)
}

func TestCore_Create_Success_With_Positive_Balance_And_NO_Existing_Balance(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{
		AccountId:     "some_id",
		OperationType: dto.OperationTypePurchaseWithInstallment,
		Amount:        100,
	}

	td.mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fc func(ctx context.Context) error) error {
			return fc(ctx)
		},
	)
	td.mockRepo.EXPECT().FindManyWithFilters(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	td.mockRepo.EXPECT().Create(ctx, model).Return(nil)

	err := td.core.Create(ctx, model)
	assert.NoError(t, err)
}

func TestCore_Create_Success_With_Positive_Balance_And_Existing_Negative_Balance_Example2(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{
		AccountId:     "some_id",
		OperationType: dto.OperationTypePurchaseWithInstallment,
		Amount:        60,
		Balance:       60,
	}

	td.mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fc func(ctx context.Context) error) error {
			return fc(ctx)
		},
	)
	td.mockRepo.EXPECT().FindManyWithFilters(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, models interface{}, req db2.FindManyWithFiltersRequester) error {
			receiver := models.(*[]transaction.Transaction)
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -50,
				Balance:       -50,
			})
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -23.5,
				Balance:       -23.5,
			})
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -18.7,
				Balance:       -18.7,
			})
			return nil
		})
	td.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
	td.mockRepo.EXPECT().Create(ctx, model).Return(nil)

	err := td.core.Create(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), model.Balance)
}

func TestCore_Create_Success_With_Positive_Balance_And_Existing_Negative_Balance_Example3(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	ctx := context.Background()
	model := &transaction.Transaction{
		AccountId:     "some_id",
		OperationType: dto.OperationTypePurchaseWithInstallment,
		Amount:        60,
		Balance:       60,
	}

	td.mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fc func(ctx context.Context) error) error {
			return fc(ctx)
		},
	)
	td.mockRepo.EXPECT().FindManyWithFilters(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, models interface{}, req db2.FindManyWithFiltersRequester) error {
			receiver := models.(*[]transaction.Transaction)
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -50,
				Balance:       -50,
			})
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -23.5,
				Balance:       -23.5,
			})
			*receiver = append(*receiver, transaction.Transaction{
				AccountId:     "some_id",
				OperationType: dto.OperationTypeWithdraw,
				Amount:        -18.7,
				Balance:       -18.7,
			})
			return nil
		})
	td.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
	td.mockRepo.EXPECT().Create(ctx, model).Return(nil)

	err := td.core.Create(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), model.Balance)
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
