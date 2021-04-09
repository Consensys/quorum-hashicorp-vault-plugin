package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetKey_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewGetKeyUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeKey := apputils.FakeKey()
		expectedEntry, _ := logical.StorageEntryJSON(storage.ComputeKeysStorageKey(fakeKey.ID, fakeKey.Namespace), fakeKey)
		mockStorage.EXPECT().
			Get(ctx, storage.ComputeKeysStorageKey(fakeKey.ID, fakeKey.Namespace)).
			Return(expectedEntry, nil)

		key, err := usecase.Execute(ctx, fakeKey.ID, fakeKey.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, fakeKey.Namespace, key.Namespace)
		assert.NotEmpty(t, key.PublicKey)
		assert.Equal(t, fakeKey.ID, key.ID)
	})

	t.Run("should fail with same error if Get fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, expectedErr)

		key, err := usecase.Execute(ctx, "my-key", "namespace")

		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return CodedError with status 404 if nothing is found", func(t *testing.T) {
		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)

		key, err := usecase.Execute(ctx, "my-key", "namespace")

		assert.Nil(t, key)
		assert.Error(t, err)

		codedError, ok := err.(logical.HTTPCodedError)
		assert.True(t, ok)
		assert.Equal(t, 404, codedError.Code())
	})
}
