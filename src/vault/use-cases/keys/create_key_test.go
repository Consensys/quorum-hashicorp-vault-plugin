package keys

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateKey_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewCreateKeyUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully by generating a private key: BN254, EDDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.EDDSA
		fakeKey.Curve = entities.BN254

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, "", fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.NotEmpty(t, key.PublicKey)
	})

	t.Run("should execute use case successfully by importing a private key: BN254, EDDSA", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.EDDSA
		fakeKey.Curve = entities.BN254
		privKey := "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZGGbioek5qYuzJzTNZpTHrVjjFk7iFe3FYwfpxZyNPxtIaFB5gb9VP9IcHZewwNZly821re7RkmB8pGdjywygPH"

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, privKey, fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.Equal(t, "X9Yz_5-O42-eOodHCUBhA4VMD2ZQy5CMAQ6lXqvDUZE=", key.PublicKey)
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
		privKey := "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, privKey, fakeKey.Tags)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.Equal(t, "BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=", key.PublicKey)
	})

	t.Run("should fail with InvalidParameter if curve and algo pair are not supported", func(t *testing.T) {
		key, err := usecase.Execute(ctx, "namespace", "id", "invalidAlgo", entities.ECDSA, "", map[string]string{})

		assert.Nil(t, key)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with StorageError if Put fails", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		fakeKey.Algorithm = entities.ECDSA
		fakeKey.Curve = entities.Secp256k1

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(fmt.Errorf("error"))

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, fakeKey.Algorithm, fakeKey.Curve, "", map[string]string{})
		assert.Nil(t, key)
		assert.True(t, errors.IsStorageError(err))
	})
}
