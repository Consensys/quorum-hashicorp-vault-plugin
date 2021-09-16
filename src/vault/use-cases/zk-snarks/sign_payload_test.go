package zksnarks

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
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
	mockGetAccountUC := mocks.NewMockGetZksAccountUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC).AnyTimes()

	usecase := NewSignUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		account := apputils.FakeZksAccount()

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of EDDSA private key fails", func(t *testing.T) {
		account := apputils.FakeZksAccount()
		account.PrivateKey = "account.PrivateKey"

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail with InvalidParameterError if data is not a hex string", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA
		key.PrivateKey = "account.PrivateKey"

		signature, err := usecase.Execute(ctx, address, namespace, "invalid data")

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
