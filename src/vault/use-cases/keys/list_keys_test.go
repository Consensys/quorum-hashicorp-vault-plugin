package keys

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListKeys_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewListKeysUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		key := apputils.FakeKey()
		expectedKeys := []string{key.ID}
		mockStorage.EXPECT().List(ctx, storage.ComputeKeysStorageKey("", key.Namespace)).Return(expectedKeys, nil)

		keys, err := usecase.Execute(ctx, key.Namespace)

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
