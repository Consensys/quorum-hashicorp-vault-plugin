package keys

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	mocks2 "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateKey_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockGetKey := mocks2.NewMockGetKeyUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	mockGetKey.EXPECT().WithStorage(mockStorage).Return(mockGetKey).AnyTimes()

	usecase := NewUpdateKeyUseCase(mockGetKey).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		newTags := map[string]string{
			"newTag1": "tagValue1",
			"newTag2": "tagValue2",
		}

		mockGetKey.EXPECT().Execute(ctx, fakeKey.ID, fakeKey.Namespace).Return(fakeKey, nil)
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, newTags)

		assert.NoError(t, err)
		assert.Equal(t, newTags, key.Tags)
	})

	t.Run("should fail with same error if GetKey fails", func(t *testing.T) {
		fakeKey := utils.FakeKey()
		expectedErr := fmt.Errorf("error")

		mockGetKey.EXPECT().Execute(ctx, fakeKey.ID, fakeKey.Namespace).Return(nil, expectedErr)

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, map[string]string{})
		assert.Nil(t, key)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail with StorageError if Put fails", func(t *testing.T) {
		fakeKey := utils.FakeKey()

		mockGetKey.EXPECT().Execute(ctx, fakeKey.ID, fakeKey.Namespace).Return(fakeKey, nil)
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(fmt.Errorf("error"))

		key, err := usecase.Execute(ctx, fakeKey.Namespace, fakeKey.ID, map[string]string{})
		assert.Nil(t, key)
		assert.True(t, errors.IsStorageError(err))
	})
}
