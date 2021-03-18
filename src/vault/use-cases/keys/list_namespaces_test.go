package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListNamespaces_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewListNamespacesUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		expectedNamespaces := []string{"ns1/ns2", "_", "tenant0", ""}

		gomock.InOrder(
			mockStorage.EXPECT().List(ctx, "").Return([]string{"ns1/", "_/", "tenant0/", "keys/"}, nil),

			// ns1/ns2 with 1 account
			mockStorage.EXPECT().List(ctx, "ns1/").Return([]string{"ns2/"}, nil),
			mockStorage.EXPECT().List(ctx, "ns1/ns2/").Return([]string{"keys/"}, nil),

			// _ with 1 account
			mockStorage.EXPECT().List(ctx, "_/").Return([]string{"keys/"}, nil),

			// tenant0 with 2 accounts
			mockStorage.EXPECT().List(ctx, "tenant0/").Return([]string{"keys/", "keys/"}, nil),
		)

		namespaces, err := usecase.Execute(ctx)

		assert.NoError(t, err)
		assert.Contains(t, namespaces, expectedNamespaces[0])
		assert.Contains(t, namespaces, expectedNamespaces[1])
		assert.Contains(t, namespaces, expectedNamespaces[2])
		assert.Contains(t, namespaces, expectedNamespaces[3])
	})

	t.Run("should fail with same error if List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr)

		keys, err := usecase.Execute(ctx)

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail with same error if recursive List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		gomock.InOrder(
			mockStorage.EXPECT().List(ctx, "").Return([]string{"ns1/", "_/", "tenant0/", "keys/"}, nil),
			mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr),
		)
		keys, err := usecase.Execute(ctx)

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})
}
