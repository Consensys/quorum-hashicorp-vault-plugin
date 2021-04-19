package keys

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignPayload_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	mockGetKeyUC := mocks.NewMockGetKeyUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"

	mockGetKeyUC.EXPECT().WithStorage(mockStorage).Return(mockGetKeyUC).AnyTimes()

	usecase := NewSignUseCase(mockGetKeyUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully: ECDSA", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA
		key.PrivateKey = "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("should execute use case successfully: EDDSA", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.BN254
		key.Algorithm = entities.EDDSA

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetKeyUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of EDDSA private key fails", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.BN254
		key.Algorithm = entities.EDDSA
		key.PrivateKey = "account.PrivateKey"

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "0xdaaa")

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail if creation of ECDSA private key fails", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA
		key.PrivateKey = "account.PrivateKey"

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

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
