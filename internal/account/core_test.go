package account_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"transaction-server/internal/account"
	"transaction-server/internal/account/mock"
)

type testDependencies struct {
	mockRepoCtrl *gomock.Controller
	mockRepo     *mock.MockIRepo
	core         account.ICore
}

func setupTest(t *testing.T) *testDependencies {
	mockRepoCtrl := gomock.NewController(t)
	mockRepo := mock.NewMockIRepo(mockRepoCtrl)
	core := account.NewCore(mockRepo)
	return &testDependencies{
		mockRepoCtrl: mockRepoCtrl,
		mockRepo:     mockRepo,
		core:         core,
	}
}

func teardownTest(td *testDependencies) {
	td.mockRepoCtrl.Finish()
}

func TestCore_Create(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	td.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	err := td.core.Create(context.Background(), &account.Account{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCore_Get(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	td.mockRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), "id").Return(nil)

	err := td.core.Get(context.Background(), &account.Account{}, "id")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCore_Create_Error(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	td.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("create error"))

	err := td.core.Create(context.Background(), &account.Account{})
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestCore_Get_Error(t *testing.T) {
	td := setupTest(t)
	defer teardownTest(td)

	td.mockRepo.EXPECT().FindByID(gomock.Any(), gomock.Any(), "id").Return(errors.New("get error"))

	err := td.core.Get(context.Background(), &account.Account{}, "id")
	if err == nil {
		t.Error("expected error but got nil")
	}
}
