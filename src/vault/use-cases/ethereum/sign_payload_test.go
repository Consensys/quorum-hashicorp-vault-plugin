package ethereum

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	mocks2 "github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignPayload_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	mockGetAccountUC := mocks.NewMockGetAccountUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC).AnyTimes()

	usecase := NewSignUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		account.PrivateKey = "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xda")

		assert.NoError(t, err)
		assert.Equal(t, "0x618b7f28507e1a1fe180005393ad2e61f6dca806c5bbd7426e3377fc23476a775644aea49fcc0e4a1c61f490979fc7e91043b22c58b441a075255eb88c034bc400", signature)
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, "0xda")

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of ECDSA private key fails", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		account.PrivateKey = "invalidPrivKey"

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xda")

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail with InvalidParameterError if data is not a hex string", func(t *testing.T) {
		signature, err := usecase.Execute(ctx, address, namespace, "invalid data")

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
