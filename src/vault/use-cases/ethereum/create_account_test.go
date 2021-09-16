package ethereum

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities/testutils"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/consensys/quorum/common"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewCreateAccountUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully by generating a private key", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount.Namespace, "")

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.True(t, common.IsHexAddress(account.Address))
	})

	t.Run("should execute use case successfully by importing a private key", func(t *testing.T) {
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		fakeAccount := testutils.FakeETHAccount()
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount.Namespace, privKey)

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})

	t.Run("should fail with StorageError if Put fails", func(t *testing.T) {
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(fmt.Errorf("error"))

		account, err := usecase.Execute(ctx, "namespace", "")
		assert.Nil(t, account)
		assert.True(t, errors.IsStorageError(err))
	})
}
