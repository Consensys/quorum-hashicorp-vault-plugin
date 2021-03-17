package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func TestCreateKey_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewCreateKeyUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully by generating a private key: BN256, EDDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.EDDSA
		fakeKey.Curve = entities.BN256

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, "", fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.NotEmpty(t, key.PublicKey)
	})

	t.Run("should execute use case successfully by importing a private key: BN256, EDDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.EDDSA
		fakeKey.Curve = entities.BN256
		privKey := "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191866e2a1e939a98bb32734cd6694c7ad58e3164ee215edc56307e9c59c8d3f1b4868507981bf553fd21c1d97b0c0d665cbcdb5adeed192607ca46763cb0ca03c7"

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, privKey, fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.Equal(t, "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191", key.PublicKey)
	})

	t.Run("should execute use case successfully by generating a private key: Secp256k1, ECDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.ECDSA
		fakeKey.Curve = entities.Secp256k1

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, "", fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.NotEmpty(t, key.PublicKey)
	})

	t.Run("should execute use case successfully by importing a private key: Secp256k1, ECDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.ECDSA
		fakeKey.Curve = entities.Secp256k1
		privKey := "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, privKey, fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", key.PublicKey)
	})

	t.Run("should fail with InvalidParameter if curve and algo pair are not supported", func(t *testing.T) {
		key, err := usecase.Execute(ctx, "namespace", "id", "invalidAlgo", entities.ECDSA, "", map[string]string{})

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if Store fails", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.ECDSA
		fakeKey.Curve = entities.Secp256k1
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(expectedErr)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, "", map[string]string{})

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
	})
}
