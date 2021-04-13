package zksnarks

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewGetAccountUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := apputils.FakeZksAccount()
		expectedEntry, _ := logical.StorageEntryJSON(
			storage.ComputeZksStorageKey(fakeAccount.PublicKey, fakeAccount.Namespace), fakeAccount)
		mockStorage.EXPECT().
			Get(ctx, storage.ComputeZksStorageKey(fakeAccount.PublicKey, fakeAccount.Namespace)).
			Return(expectedEntry, nil)

		account, err := usecase.Execute(ctx, fakeAccount.PublicKey, fakeAccount.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.NotEmpty(t, account.PublicKey)
	})

	t.Run("should fail with StorageError if Get fails", func(t *testing.T) {
		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))

		account, err := usecase.Execute(ctx, "0xaddress", "namespace")

		assert.Nil(t, account)
		assert.True(t, errors.IsStorageError(err))
	})

	t.Run("should fail with NotFoundError if nothing is found", func(t *testing.T) {
		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)

		account, err := usecase.Execute(ctx, "0xaddress", "namespace")
		assert.Nil(t, account)
		assert.True(t, errors.IsNotFoundError(err))
	})
}
