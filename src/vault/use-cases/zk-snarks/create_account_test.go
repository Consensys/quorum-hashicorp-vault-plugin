package zksnarks

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
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
		fakeAccount := utils.FakeZksAccount()
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.NotEmpty(t, account.PublicKey)
	})

	t.Run("should fail with StorageError if Put fails", func(t *testing.T) {
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(fmt.Errorf("error"))

		account, err := usecase.Execute(ctx, "namespace")
		assert.Nil(t, account)
		assert.True(t, errors.IsStorageError(err))
	})
}
