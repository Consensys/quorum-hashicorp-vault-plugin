package zksnarks

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAccounts_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewListAccountsUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := apputils.FakeETHAccount()
		expectedKeys := []string{fakeAccount.Address}
		mockStorage.EXPECT().List(ctx, storage.ComputeZksStorageKey("", fakeAccount.Namespace)).Return(expectedKeys, nil)

		keys, err := usecase.Execute(ctx, fakeAccount.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, expectedKeys, keys)
	})

	t.Run("should fail with StorageError if List fails", func(t *testing.T) {
		mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))

		keys, err := usecase.Execute(ctx, "namespace")

		assert.Nil(t, keys)
		assert.True(t, errors.IsStorageError(err))
	})
}
