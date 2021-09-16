package keys

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

func TestDestroyKey_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewDestroyKeyUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeKey := utils.FakeKey()

		mockStorage.EXPECT().Delete(ctx, gomock.Any()).Return(nil)

		err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID)
		assert.NoError(t, err)
	})

	t.Run("should fail with StorageError if Delete fails", func(t *testing.T) {
		fakeKey := utils.FakeKey()

		mockStorage.EXPECT().Delete(ctx, gomock.Any()).Return(fmt.Errorf("error"))

		err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID)
		assert.True(t, errors.IsStorageError(err))
	})
}
